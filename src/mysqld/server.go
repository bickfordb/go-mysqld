package mysqld

import "encoding/binary"
import "fmt"
import "net"
import "os"

func (server *Server) Listen(port int) (err error) {
  addr := fmt.Sprintf(":%d", port)
  ln, err := net.Listen("tcp", addr)
  if err != nil {
    return
  }
  server.CapabilityFlags.Set(LONG_PASSWORD,
    LONG_FLAG,
    FOUND_ROWS,
    CONNECT_WITH_DB,
    IGNORE_SPACE,
    PROTOCOL_41,
    INTERACTIVE,
    TRANSACTIONS,
    MULTI_STATEMENTS,
    CONNECT_ATTRS,
    PLUGIN_AUTH,
    SECURE_CONNECTION,
    PLUGIN_AUTH_LENENC_CLIENT_DATA)

  for {
    conn, err := ln.Accept()
    if err != nil {
      println("err:", err)
    } else {
      c := &Conn{}
      c.Conn = conn
      c.Server = server
      go HandleConnection(c)
    }
  }
  return
}


func (c *Conn) WriteHeader(length uint32) (err error) {
  header := make([]byte, 4)
  header[0] = byte(length & 0xff)
  header[1] = byte((length >> 8) & 0xff)
  header[2] = byte((length >> 16) & 0xff)
  header[3] = byte(c.Sequence)
  c.Sequence += 1
  err = WriteAll(header, c.Conn)
  return
}

func (c *Conn) WriteBuffer(buffer *MySQLBuffer) (err error) {
  bs := buffer.Bytes()
  err = c.WriteHeader(uint32(len(bs)))
  if err != nil { return }
  err = WriteAll(bs, c.Conn)
  return
}

func WriteInitialHandshake(conn *Conn) (err error) {
  buf := &MySQLBuffer{}
  // capability flags (lower 2 bytes)
  flags := uint32(conn.Server.CapabilityFlags)
  salt := []byte{
    0x77, 0x63, 0x6a, 0x6d, 0x61, 0x22, 0x23, 0x27, // first part
    0x38, 0x26, 0x55, 0x58, 0x3b, 0x5d, 0x44, 0x78, 0x53, 0x73, 0x6b, 0x41,
    0x00}
  // protocol version
  buf.WriteByte(0x0a)
  // server version
  buf.WriteNullString("5.5.29-debug")
  // connection id
  binary.Write(buf, binary.LittleEndian, uint32(4))
  buf.Write(salt[:8])
  buf.WriteByte(0x0)

  var flagsL = uint16(flags)
  var flagsH = uint16(flags >> 16)

  // lower 2 bytes of flags
  binary.Write(buf, binary.LittleEndian, flagsL)
  // character set
  buf.WriteByte(byte(CHARSET_UTF8))
  // status flags
  binary.Write(buf, binary.LittleEndian, uint16(conn.StatusFlags))
  // capability flags (upper 2 bytes)
  binary.Write(buf, binary.LittleEndian, flagsH)

  // auth plugin data length
  buf.WriteByte(byte(len(salt)))
  // 10 byte reserved
  for i := 0; i < 10; i += 1{
    buf.WriteByte(0x0)
  }
  for i := 8; i < len(salt); i += 1 {
    buf.WriteByte(salt[i])
  }
  // auth plugin name
  buf.WriteNullString("mysql_native_password")
  err = conn.WriteBuffer(buf)
  return
}

func (conn *Conn) ReadPacket() (packet []byte, err error) {
  header := make([]byte, 4)
  err = ReadAll(header, conn.Conn)
  if err != nil { return }
  var length uint32
  length |= uint32(header[0])
  length |= uint32(header[1] << 8)
  length |= uint32(header[2] << 16)
  //sequence := header[3]
  packet = make([]byte, int(length))
  err = ReadAll(packet, conn.Conn)
  conn.Sequence += 1
  return
}

func (conn *Conn) SendOK() (err error) {
  buf := &MySQLBuffer{}
  // OK
  buf.WriteByte(0x00)
  // affected rows
  buf.WriteLenc(0)
  // last insert id
  buf.WriteLenc(0)
  if (conn.CapabilityFlags.Test(PROTOCOL_41)) {
    // status
    binary.Write(buf, binary.LittleEndian, uint16(conn.StatusFlags))
    // warnings
    binary.Write(buf, binary.LittleEndian, uint16(0))
  } else if (conn.CapabilityFlags.Test(TRANSACTIONS)) {
    // status
    binary.Write(buf, binary.LittleEndian, uint16(0))
  }
  // extra info
  conn.WriteBuffer(buf)
  conn.Sequence = 0
  return
}

func ReadAuth(conn *Conn) (err error) {
  packet, err := conn.ReadPacket()
  buf := NewMySQLBuffer(packet)
  flagsOrd, err := buf.ReadU32()
  if err != nil { return }
  conn.CapabilityFlags = CapabilityFlag(flagsOrd)
  if err != nil { return }
  conn.MaxPacketSize, err = buf.ReadU32()
  if err != nil { return }
  conn.CharacterSet, err = buf.ReadByte()
  if err != nil { return }
  // reserved
  buf.Next(23)
  conn.Username, err = buf.ReadNullString()
  if err != nil { return }
  auth, err := buf.ReadNullString()
  conn.PluginAuth0 = string(auth)
  if err != nil { return }
  if conn.CapabilityFlags.Test(CONNECT_WITH_DB) && buf.Len() > 0 {
    conn.Database, err = buf.ReadNullString()
    if err != nil { return }
  }
  log("read plugin auth")
  if conn.CapabilityFlags.Test(PLUGIN_AUTH) {
    conn.PluginAuth1, err = buf.ReadNullStringOrEOF()
  }
  if err != nil { return }
  log("done read plugin auth")
  if conn.CapabilityFlags.Test(CONNECT_ATTRS) {
    err = fmt.Errorf("dont know how to read connect attrs")
  }
  if err != nil { return }
  if buf.Len() != 0 {
    err = fmt.Errorf("%d unexpected bytes", buf.Len())
  }
  return
}

func handleConnection(conn *Conn) (err error) {
  conn.StatusFlags.Set(SERVER_STATUS_AUTOCOMMIT)
  err = WriteInitialHandshake(conn)
  if err != nil { return }
  err = ReadAuth(conn)
  if err != nil { return }
  err = conn.SendOK()
  if err != nil { return }
  err = HandleCommands(conn)
  if err != nil { return }
  log("done")
  return
}

func HandleCommands(conn *Conn) (err error) {
  for {
    var pkt []byte
    conn.Sequence = 0
    pkt, err = conn.ReadPacket()
    if err != nil { return }
    if len(pkt) < 1 {
      err = fmt.Errorf("expecting a command")
      return
    }
    cmd := Command(pkt[0])
    switch cmd {
    case COM_QUIT:
      err = conn.SendOK()
      return
    case COM_QUERY:
      query := string(pkt[1:])
      results := make(chan map[string]interface{})
      errors := make(chan Error)
      handler := conn.Server.OnQuery
      if handler == nil {
        handler = defaultOnQuery
      }
      go handler(conn, query, results, errors)
      err = conn.SendResultSet(results, errors)
      if err != nil { return }
      continue
    default:
      conn.SendError(NotImplemented)
      continue
    }
  }
  return
}

func (conn *Conn) SendNumFields(numFields uint64) (err error) {
  b := &MySQLBuffer{}
  b.WriteLenc(numFields)
  err = conn.WriteBuffer(b)
  return
}

func (conn *Conn) SendColumnDef(name string, columnType byte) (err error) {
  flags := uint16(0)
  decimals := byte(0x51)
  buf := &MySQLBuffer{}
  // catalog
  buf.WriteLencString("def")
  // schema
  buf.WriteLencString("")
  // table
  buf.WriteLencString("")
  // org_table
  buf.WriteLencString("")
  // name
  buf.WriteLencString(name)
  // org_name
  buf.WriteLencString(name)
  // next length
  buf.WriteLenc(0x0c)
  // character_set
  binary.Write(buf, binary.LittleEndian, uint16(CHARSET_UTF8))
  // column length
  binary.Write(buf, binary.LittleEndian, uint32(1 << 16))
  // column type
  buf.WriteByte(columnType)
  // flags
  binary.Write(buf, binary.LittleEndian, flags)
  // decimals
  binary.Write(buf, binary.LittleEndian, decimals)
  // Reserved
  buf.WriteByte(0)
  buf.WriteByte(0)
  err = conn.WriteBuffer(buf)
  return
}

func (conn *Conn) SendRow(cols []string, row map[string]interface{}) (err error) {
  buf := &MySQLBuffer{}
  for _, col := range cols {
    err = buf.WriteLencString(fmt.Sprintf("%v", row[col]))
    if err != nil { return }
  }
  err = conn.WriteBuffer(buf)
  return
}

func (conn *Conn) SendError(err Error) (e error) {
  buf := &MySQLBuffer{}
  buf.WriteByte(ERR)
  binary.Write(buf, binary.LittleEndian, err.Code)
  // sql-state marker #
  buf.WriteString(`#`)
  // sql-state (?) 5 ascii bytes
  var state string
  if len(err.State) < 5 {
    state = `S1000`
  } else {
    state = err.State[:5]
  }
  buf.WriteString(state)
  buf.WriteString(err.Message)
  e = conn.WriteBuffer(buf)
  return
}

func (conn *Conn) SendResultSet(rows chan map[string]interface{}, errors chan Error) (err error) {
  if rows == nil {
    err = conn.SendOK()
    return
  }
  i := -1
  var cols []string
  for {
    i += 1
    var errOk, rowOk bool
    var row map[string]interface{}
    var anError Error
    select {
    case row, rowOk = <-rows:
    case anError, errOk = <-errors:
      if errOk {
        conn.SendError(anError)
      }
      continue
      //return
    }
    fmt.Println("Got row")
    if !rowOk {
      fmt.Println("done")
      break
    }
    if i == 0 {
      conn.SendNumFields(uint64(len(row)))
      // send the column definition
      for col, val := range row {
        cols = append(cols, col)
        columnType := byte(MYSQL_TYPE_STRING)
        switch val.(type) {
        case byte, uint, uint16, uint32, uint64, int, int8, int16, int32, int64:
          columnType = MYSQL_TYPE_LONGLONG
        }
        err = conn.SendColumnDef(col, columnType)
        if err != nil { return }
      }
      err = conn.SendEOF(0, 0)
      if err != nil { return }
    }
    err = conn.SendRow(cols, row)
    if err != nil { return }
  }
  err = conn.SendEOF(0, 0)
  return
}

func (conn *Conn) SendEOF(warnings uint16, flags StatusFlag) (err error) {
  b := &MySQLBuffer{}
  b.WriteByte(EOF)
  binary.Write(b, binary.LittleEndian, warnings)
  binary.Write(b, binary.LittleEndian, uint16(flags))
  err = conn.WriteBuffer(b)
  return
}

func HandleConnection(conn *Conn) {
  defer conn.Conn.Close()
  fmt.Println("connection from:", conn.Conn.RemoteAddr(), "on", conn.Conn.LocalAddr())
  err := handleConnection(conn)
  if err != nil {
    fmt.Fprintf(os.Stderr, "unexpected error: %s\n", err.Error())
  }
  return
}

func defaultOnQuery(conn *Conn, query string, results chan map[string]interface{}, errors chan Error) {
  defer close(errors)
  defer close(results)
  errors<-NotImplemented
}

var NotImplemented Error = Error{
  Code:1,
  Message:"Not Implemented"}


