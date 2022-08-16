package gameserver

import (
	"encoding/binary"
	"net"
	"pangya/src/internal/pangya"
	"pangya/src/internal/utils"
)

type GameServer struct {
	srv  pangya.Server
	info pangya.ServerInfo
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
	return &GameServer{
		srv: pangya.NewServer(&gameServerConfig{}),
		info: pangya.ServerInfo{
			Name:     utils.GetStringEnv("GAME_NAME"),
			IP:       utils.GetStringEnv("GAME_HOST"),
			Port:     utils.GetUint16Env("GAME_PORT"),
			MaxUsers: utils.GetUint32Env("GAME_MAX_USERS"),
			Flags:    utils.GetUint32Env("GAME_FLAGS"),
			Boosts:   utils.GetUint16Env("GAME_BOOSTS"),
		},
	}
}

func (ls *GameServer) Listen(port int) error {
	return ls.srv.Listen(port)
}

func (ls *GameServer) AddHandler(id uint16, ph pangya.PacketHandler) {
	ls.srv.AddHandler(id, ph)
}

func (ls *GameServer) ServerInfo() pangya.ServerInfo {
	return pangya.ServerInfo{
		Type:     "GameServer",
		Name:     ls.info.Name,
		IP:       ls.info.IP,
		Port:     ls.info.Port,
		MaxUsers: ls.info.MaxUsers,
		Flags:    ls.info.Flags,
		Boosts:   ls.info.Boosts,
	}
}
