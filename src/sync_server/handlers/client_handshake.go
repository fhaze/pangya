package handlers

import (
	"encoding/json"
	"net"
	"pangya/src/internal/logger"
	"pangya/src/internal/sync"
	syncserver "pangya/src/sync_server"
)

type ClientHandshake struct {
	svc syncserver.Server
}

func NewClientHandshake(svc syncserver.Server) sync.PacketHandler {
	return &ClientHandshake{svc: svc}
}

func (ch *ClientHandshake) Action(conn net.Conn, pak []byte) error {
	var req sync.ClientPacketHandshake
	if err := json.Unmarshal(pak, &req); err != nil {
		return err
	}

	ch.svc.AddClient(req.Server, conn)
	logger.Log.Sugar().Infof("registered %s from %s", req.Server, conn.RemoteAddr())

	res := sync.ServerPacketHandshake{
		PacketBase: sync.PacketBase{
			ID: sync.PacketHandshake,
		},
		Status: "OK",
	}
	var buf []byte
	buf, err := json.Marshal(res)
	if err != nil {
		return err
	}

	conn.Write(buf)
	return nil
}
