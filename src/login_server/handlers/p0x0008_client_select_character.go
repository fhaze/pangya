package handlers

import (
	"pangya/src/internal/logger"
	"pangya/src/internal/pangya"
	"pangya/src/internal/utils"
	loginserver "pangya/src/login_server"

	"go.uber.org/zap"
)

type P0x0008_ClientSelectCharacter struct {
	ls *loginserver.LoginServer
}

func NewP0x0008_ClientSelectCharacter(ls *loginserver.LoginServer) pangya.PacketHandler {
	return &P0x0008_ClientSelectCharacter{ls: ls}
}

func (h *P0x0008_ClientSelectCharacter) Action(c *pangya.ConnAccount, r *pangya.PacketReader, key uint16) error {
	characterId, err := r.ReadUint32()
	if err != nil {
		return err
	}
	hairColor, err := r.ReadUint16()
	if err != nil {
		return err
	}

	logger.Log.Info(
		"selected character",
		zap.Uint32("characterId", characterId),
		zap.Uint16("hairColor", hairColor),
	)

	w := pangya.NewPacket(0x0011)
	w.PutUint16(0)
	err = utils.SendEncryptedPacketToClient(w, c.Conn, key)
	if err != nil {
		return err
	}

	w = pangya.NewPacket(0x0001)
	appendLoginSuccess(&w, *c.Account)
	err = utils.SendEncryptedPacketToClient(w, c.Conn, key)
	if err != nil {
		return err
	}

	w = packetGameServerList(h.ls)
	return utils.SendEncryptedPacketToClient(w, c.Conn, key)
}
