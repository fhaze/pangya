package handlers

import (
	"net"
	"pangya/src/domain/account"
	"pangya/src/internal/logger"
	"pangya/src/internal/pangya"
	"pangya/src/internal/utils"

	"go.uber.org/zap"
)

type P0x0001_ClientLogin struct {
}

func NewP0x0001_ClientLogin() pangya.PacketHandler {
	return &P0x0001_ClientLogin{}
}

func (h *P0x0001_ClientLogin) Action(conn net.Conn, req pangya.Packet, key uint16) error {
	r := pangya.NewPacketReader(&req)
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
		w.PutUint8(0x00)                                                                                       // status ok
		w.PutLString(acc.Username)                                                                             // username
		w.PutUint32(uint32(acc.ID))                                                                            // id
		w.PutBytes([]byte{0x00, 0x00, 0x00, 0x00, 0x25, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}) // unknown
		w.PutUint16(0x00)                                                                                      // nickname?
	} else {
		w.PutUint8(227)      // status
		w.PutUint32(5100143) // invalid credentials
	}

	return utils.SendEncryptedPacketToClient(w, conn, key)
}
