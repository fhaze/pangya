package handlers

import (
	"pangya/src/internal/logger"
	"pangya/src/internal/pangya"
	"pangya/src/internal/utils"
	loginserver "pangya/src/login_server"

	"go.uber.org/zap"
)

type P0x0003_ClientSelectServer struct {
	ls *loginserver.LoginServer
}

func NewP0x0003_ClientSelectServer(ls *loginserver.LoginServer) pangya.PacketHandler {
	return &P0x0003_ClientSelectServer{ls: ls}
}

func (h *P0x0003_ClientSelectServer) Action(c *pangya.ConnAccount, r *pangya.PacketReader, key uint16) error {
	gameServerID, err := r.ReadUint16()
	if err != nil {
		return err
	}
	logger.Log.Info("player selected server", zap.Uint16("id", gameServerID))

	res := pangya.NewPacket(0x0003)
	res.PutUint16(0x0000) // unknown
	res.PutLString("key")
	return utils.SendEncryptedPacketToClient(res, c.Conn, key)
}
