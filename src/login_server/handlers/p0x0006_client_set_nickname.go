package handlers

import (
	"pangya/src/domain/account"
	"pangya/src/internal/pangya"
	"pangya/src/internal/utils"
	loginserver "pangya/src/login_server"
)

type P0x0006_ClientSetNickname struct {
	ls *loginserver.LoginServer
}

func NewP0x0006_ClientSetNickname(ls *loginserver.LoginServer) pangya.PacketHandler {
	return &P0x0006_ClientSetNickname{ls: ls}
}

func (h *P0x0006_ClientSetNickname) Action(c *pangya.ConnAccount, r *pangya.PacketReader, key uint16) error {
	nickname, err := r.ReadLString()
	if err != nil {
		return err
	}

	if err = account.Svc().SetNicknamebyId(nickname, c.Account.ID); err != nil {
		return err
	}

	p := pangya.NewPacket(0x0001)
	p.PutUint16(217)
	return utils.SendEncryptedPacketToClient(p, c.Conn, key)
}
