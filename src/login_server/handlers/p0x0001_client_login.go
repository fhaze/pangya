package handlers

import (
	"pangya/src/domain/account"
	"pangya/src/internal/logger"
	"pangya/src/internal/pangya"
	"pangya/src/internal/utils"
	loginserver "pangya/src/login_server"

	"go.uber.org/zap"
)

type P0x0001_ClientLogin struct {
	ls *loginserver.LoginServer
}

func NewP0x0001_ClientLogin(ls *loginserver.LoginServer) pangya.PacketHandler {
	return &P0x0001_ClientLogin{ls: ls}
}

func (h *P0x0001_ClientLogin) Action(c *pangya.ConnAccount, r *pangya.PacketReader, key uint16) error {
	username, err := r.ReadLString()
	if err != nil {
		return err
	}
	password, err := r.ReadLString()
	if err != nil {
		return err
	}

	logger.Log.Debug(
		"trying login",
		zap.String("username", username),
		zap.String("password", password),
	)

	w := pangya.NewPacket(0x0001)

	if acc, found := account.Svc().FindAccountByUsernameAndPassword(username, password); found {
		c.Account = &acc

		if acc.Nickname == nil {
			appendNeedSetNickname(&w)
			return utils.SendEncryptedPacketToClient(w, c.Conn, key)
		}

		appendLoginSuccess(&w, acc)
		if err := utils.SendEncryptedPacketToClient(w, c.Conn, key); err != nil {
			return err
		}

		w = packetGameServerList(h.ls)
		return utils.SendEncryptedPacketToClient(w, c.Conn, key)
	}

	appendInvalidCredentials(&w)
	return utils.SendEncryptedPacketToClient(w, c.Conn, key)
}
