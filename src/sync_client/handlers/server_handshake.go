package handlers

import (
	"encoding/json"
	"net"
	"pangya/src/internal/logger"
	"pangya/src/internal/sync"
	syncclient "pangya/src/sync_client"

	"go.uber.org/zap"
)

type ServerHandshake struct {
	svc syncclient.Client
}

func (*ServerHandshake) Action(conn net.Conn, pak []byte) error {
	var res sync.ServerPacketHandshake
	if err := json.Unmarshal(pak, &res); err != nil {
		return err
	}

	if res.Status == "OK" {
		logger.Log.Info("successfully registered on sync server")
	} else {
		logger.Log.Error(
			"error trying to register on sync server",
			zap.String("status", res.Status),
			zap.String("message", res.Message),
		)
	}
	return nil
}

func NewServerhandshake(svc syncclient.Client) sync.PacketHandler {
	return &ServerHandshake{svc: svc}
}
