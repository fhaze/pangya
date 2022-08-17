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

	i := pangya.NewPacket()
	i.PutBytes([]byte{0x00, 0x3f, 0x00, 0x01, 0x01}) // unknown
	i.PutBytes(keyBytes)
	i.PutLString("127.0.0.1")

	p := pangya.NewPacket()
	p.PutUint8(0x00)
	p.PutLString(string(i.ToBytes()))

	conn.Write(p.ToBytes())
	return key
}

func New() pangya.Server {
	return &GameServer{
		srv: pangya.NewServer(&gameServerConfig{}),
		info: pangya.ServerInfo{
			Type:     "GameServer",
			ID:       utils.GetUint16Env("GAME_ID"),
			Name:     utils.GetStringEnv("GAME_NAME"),
			IP:       utils.GetStringEnv("GAME_HOST"),
			Port:     utils.GetUint16Env("GAME_PORT"),
			MaxUsers: utils.GetUint32Env("GAME_MAX_USERS"),
			Flags:    utils.GetUint16Env("GAME_FLAGS"),
			Boosts:   utils.GetUint16Env("GAME_BOOSTS"),
			Icon:     utils.GetUint16Env("GAME_ICON"),
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
		Type:     ls.info.Type,
		ID:       ls.info.ID,
		Name:     ls.info.Name,
		IP:       ls.info.IP,
		Port:     ls.info.Port,
		MaxUsers: ls.info.MaxUsers,
		Flags:    ls.info.Flags,
		Boosts:   ls.info.Boosts,
		Icon:     ls.info.Icon,
	}
}
