package utils

import (
	"encoding/hex"
	"net"
	"pangya/src/internal/logger"
	"pangya/src/internal/pangya"

	"github.com/pangbox/pangcrypt"
	"go.uber.org/zap"
)

func SendEncryptedPacketToClient(pak pangya.Packet, conn net.Conn, key uint16) error {
	data, err := pangcrypt.ServerEncrypt(pak.ToBytes(), byte(key), 0x0001)
	if err != nil {
		logger.Log.Error(
			"error encrypting packet",
			zap.String("client", conn.RemoteAddr().String()),
			zap.Uint16("packetID", pak.ID),
			zap.String("packetPayload", hex.EncodeToString(pak.Payload)),
			zap.String("packetPayloadStr", string(pak.Payload)),
			zap.Error(err),
		)
		return err
	}
	l, err := conn.Write(data)
	if err != nil {
		logger.Log.Error(
			"error sending packet",
			zap.String("client", conn.RemoteAddr().String()),
			zap.Uint16("packetID", pak.ID),
			zap.String("packetPayload", hex.EncodeToString(pak.Payload)),
			zap.String("packetPayloadStr", string(pak.Payload)),
			zap.Error(err),
		)
		return err
	}
	logger.Log.Debug(
		"sent packet",
		zap.String("client", conn.RemoteAddr().String()),
		zap.Int("length", l),
		zap.Uint16("packetID", pak.ID),
		zap.String("packetPayload", hex.EncodeToString(pak.Payload)),
		zap.String("packetPayloadStr", string(pak.Payload)),
		zap.Error(err),
	)
	return nil
}
