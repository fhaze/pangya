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
	logger.Log.Sugar().Infof("registered %v from %s", req.Server, conn.RemoteAddr())

	res := sync.ServerPacketGameServerList{
		PacketBase: sync.PacketBase{
			ID: sync.PacketGameServerList,
		},
	}
	for _, gameServer := range ch.svc.GameServerList() {
		res.Servers = append(res.Servers, gameServer.Info)
	}

	var l int
	for _, loginServer := range ch.svc.LoginServerList() {
		var buf []byte
		buf, err := json.Marshal(res)
		if err != nil {
			return err
		}
		if _, err = loginServer.Conn.Write(buf); err != nil {
			return err
		}
		l++
	}
	logger.Log.Sugar().Infof("broadcasted %v from %s to %d LoginServer(s)", req.Server, conn.RemoteAddr(), l)

	hsRes := sync.ServerPacketHandshake{
		PacketBase: sync.PacketBase{
			ID: sync.PacketHandshake,
		},
		Status: "OK",
	}
	var buf []byte
	buf, err := json.Marshal(hsRes)
	if err != nil {
		return err
	}

	_, err = conn.Write(buf)
	return err
}
