package mysqld

import "bytes"
import "io"
import "os"
import "fmt"
import "encoding/binary"

func WriteAll(b []byte, writer io.Writer) (err error) {
  n := len(b)
  offset := 0
  for offset < n {
    var amt int
    amt, err = writer.Write(b[offset:])
    if err != nil { return }
    offset += amt
  }
  return
}

func ReadAll(b []byte, reader io.Reader) (err error) {
  var read int
  for read < len(b) {
    var n int
    n, err = reader.Read(b[read:])
    read += n
    if err != nil { return }
  }
  return
}

func log(msg string, args ...interface{}) {
  msg = fmt.Sprintf(msg, args...)
  fmt.Fprintf(os.Stderr, "%s\n", msg)
}
func (b *MySQLBuffer) WriteLencString(s string) (err error) {
  err = b.WriteLenc(uint64(len(s)))
  if err != nil { return }
  _, err = b.WriteString(s)
  return
}

func (b *MySQLBuffer) WriteLenc(n uint64) (err error) {
  if n < 251 {
    err = b.WriteByte(byte(n))
  } else {
    err = fmt.Errorf("dont know how to write: %d", n)
  }
  return
}

func (b *MySQLBuffer) WriteOK(ok *Ok, conn *Conn) {
  b.WriteByte(0x00)
  b.WriteLenc(ok.NumRows)
  b.WriteLenc(ok.LastInsertId)
  if (conn.CapabilityFlags.Test(PROTOCOL_41)) {
    binary.Write(b, binary.LittleEndian, ok.StatusFlags)
    binary.Write(b, binary.LittleEndian, ok.Warnings)
  } else if (conn.CapabilityFlags.Test(TRANSACTIONS)) {
    binary.Write(b, binary.LittleEndian, ok.StatusFlags)
  }
  b.WriteString(ok.Info)
}

type MySQLBuffer struct {
  bytes.Buffer
}

func (b *MySQLBuffer) ReadNullString() (s string, err error) {
  s, err = b.ReadString(0x0)
  return
}

func (buf *MySQLBuffer) ReadNullStringOrEOF() (s string, err error) {
  var ret []byte = make([]byte, 0)
  for {
    var aByte byte
    aByte, err = buf.ReadByte()
    if err == io.EOF {
      s = string(ret)
      err = nil
      return
    } else if err != nil {
      return
    } else if aByte == 0x0 {
      s = string(ret)
      return
    }
    ret = append(ret, aByte)
  }
  return
}

func (b *MySQLBuffer) ReadU32() (result uint32, err error) {
  var o byte
  o, err = b.ReadByte()
  if err != nil { return }
  result |= uint32(o)
  o, err = b.ReadByte()
  if err != nil { return }
  result |= uint32(o) << 8
  o, err = b.ReadByte()
  if err != nil { return }
  result |= uint32(o) << 16
  o, err = b.ReadByte()
  if err != nil { return }
  result |= uint32(o) << 24
  return
}

func (b *MySQLBuffer) WriteNullString(s string) (err error) {
  _, err = b.WriteString(s)
  if err != nil { return }
  err = b.WriteByte(0x0)
  if err != nil { return }
  return
}

func NewMySQLBuffer(buf []byte) (result *MySQLBuffer) {
  b := bytes.NewBuffer(buf)
  result = &MySQLBuffer{*b}
  return
}

