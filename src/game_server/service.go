package gameserver

import (
	"encoding/binary"
	"net"
	"pangya/src/internal/pangya"
)

type GameServer struct {
	srv pangya.Server
}

type gameServerConfig struct {
}

func (gsc *gameServerConfig) OnClientConnect(conn net.Conn) uint16 {
	key := uint16(1)
	keyBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(keyBytes, key)

	p := pangya.NewPacket(0x0000)
	p.PutBytes([]byte{0x3f, 0x00, 0x01, 0x01})
	p.PutBytes(keyBytes)
	p.PutLString(conn.RemoteAddr().String())

	conn.Write(p.ToBytes())
	return key
}

func New() pangya.Server {
	return &GameServer{srv: pangya.NewServer(&gameServerConfig{})}
}

func (ls *GameServer) Listen(port int) error {
	return ls.srv.Listen(port)
}

func (ls *GameServer) AddHandler(id uint16, ph pangya.PacketHandler) {
	ls.srv.AddHandler(id, ph)
}

func (ls *GameServer) ServerName() string {
	return "GameServer"
}
