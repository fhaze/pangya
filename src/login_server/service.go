package loginserver

import (
	"encoding/binary"
	"net"
	"pangya/src/internal/pangya"
)

type LoginServer struct {
	srv pangya.Server
}

func New() pangya.Server {
	return &LoginServer{
		srv: pangya.NewServer(func(conn net.Conn) uint16 {
			key := uint16(1)
			keyBytes := make([]byte, 2)
			binary.LittleEndian.PutUint16(keyBytes, key)

			var res []byte
			res = append(res, []byte{0x00, 0x0b, 0x00, 0x00, 0x00, 0x00}...)
			res = append(res, keyBytes...)
			res = append(res, []byte{0x00, 0x00, 0x75, 0x27, 0x00, 0x00}...)

			conn.Write(res)
			return key
		}),
	}
}

func (ls *LoginServer) Listen(port int) error {
	return ls.srv.Listen(port)
}

func (ls *LoginServer) AddHandler(id uint16, ph pangya.PacketHandler) {
	ls.srv.AddHandler(id, ph)
}

func (ls *LoginServer) ServerName() string {
	return "LoginServer"
}
