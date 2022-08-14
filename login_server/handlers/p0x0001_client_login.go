package handlers

import (
	"net"
	"pangya/internal/logger"
	"pangya/internal/packet"
	"pangya/internal/server"

	"github.com/pangbox/pangcrypt"
	"go.uber.org/zap"
)

type P0x0001_ClientLogin struct{}

func NewP0x0001_ClientLogin() server.PacketHandler {
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

	w := packet.NewPacket(0x01)
	w.PutUint8(227)      // status
	w.PutUint32(5100143) // invalid credentials

	logger.Log.Debug(
		"try login",
		zap.String("username", username),
		zap.String("password", password),
	)

	res, err := pangcrypt.ServerEncrypt(w.ToBytes(), byte(key), 0x00)
	if err != nil {
		return err
	}

	conn.Write(res)
	return nil
}
