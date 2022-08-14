package handlers

import (
	"net"
	"pangya/src/domain/account"
	"pangya/src/internal/logger"
	"pangya/src/internal/packet"
	"pangya/src/internal/server/pangya"

	"github.com/pangbox/pangcrypt"
	"go.uber.org/zap"
)

type P0x0001_ClientLogin struct {
}

func NewP0x0001_ClientLogin() pangya.PacketHandler {
	return &P0x0001_ClientLogin{}
}

func (h *P0x0001_ClientLogin) Action(conn net.Conn, pak packet.Packet, key uint16) error {
	r := packet.NewReader(&pak)
	username, err := r.ReadLString()
	if err != nil {
		return err
	}
	password, err := r.ReadLString()
	if err != nil {
		return err
	}

	logger.Log.Debug(
		"try login",
		zap.String("username", username),
		zap.String("password", password),
	)

	w := packet.NewPacket(0x0001)

	if acc, found := account.Svc().FindAccountByUsernameAndPassword(username, password); found {
		w.PutUint8(0x00) // status ok
		w.PutLString(acc.Username)
		w.PutUint32(uint32(acc.ID))
		w.PutBytes([]byte{0x00, 0x00, 0x00, 0x00, 0x25, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}) // unknown
		w.PutUint16(0x00)                                                                                      // nickname?
	} else {
		w.PutUint8(227)      // status
		w.PutUint32(5100143) // invalid credentials
	}

	res, err := pangcrypt.ServerEncrypt(w.ToBytes(), byte(key), 0x00)
	if err != nil {
		return err
	}

	conn.Write(res)
	return nil
}
