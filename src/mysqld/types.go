package mysqld

import "net"

type Ok struct {
	NumRows      uint64
	LastInsertId uint64
	StatusFlags  uint16
	Warnings     uint16
	Info         string
}

type Conn struct {
	Conn            net.Conn
	Sequence        uint8
	PluginAuth0     string
	PluginAuth1     string
	CapabilityFlags CapabilityFlag
	StatusFlags     StatusFlag
	Server          *Server
	Username        string
	MaxPacketSize   uint32
	Database        string
	CharacterSet    byte
}

type Error struct {
	Code    uint16
	Message string
	State   string
}

type Server struct {
	CapabilityFlags CapabilityFlag
	Queries         chan *Query
}
