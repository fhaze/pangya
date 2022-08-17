package handlers

import (
	"pangya/src/domain/account"
	"pangya/src/internal/pangya"
	"pangya/src/internal/utils"
	loginserver "pangya/src/login_server"
)

type P0x0007_ClientCheckNickname struct {
	ls *loginserver.LoginServer
}

func NewP0x0007_ClientCheckNickname(ls *loginserver.LoginServer) pangya.PacketHandler {
	return &P0x0007_ClientCheckNickname{ls: ls}
}

func (h *P0x0007_ClientCheckNickname) Action(c *pangya.ConnAccount, r *pangya.PacketReader, key uint16) error {
	nickname, err := r.ReadLString()
	if err != nil {
		return nil
	}

	w := pangya.NewPacket(0x00E)

	if len(nickname) < 4 && len(nickname) > 16 {
		w.PutUint32(2)
		return utils.SendEncryptedPacketToClient(w, c.Conn, key)
	}

	exists, err := account.Svc().ExistsNickname(nickname)
	if err != nil {
		return err
	}

	if exists {
		w.PutUint32(2)
	} else {
		w.PutUint32(0)
		w.PutLString(nickname)
	}

	return utils.SendEncryptedPacketToClient(w, c.Conn, key)
}
