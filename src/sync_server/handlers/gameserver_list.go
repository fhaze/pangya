package handlers

import (
	"encoding/json"
	"net"
	"pangya/src/internal/pangya"
	"pangya/src/internal/sync"
	syncserver "pangya/src/sync_server"
)

type GameServerList struct {
	svc syncserver.Server
}

func NewGameServerList(svc syncserver.Server) sync.PacketHandler {
	return &ClientHandshake{svc: svc}
}

func (ch *GameServerList) Action(conn net.Conn, _ []byte) error {
	var buf []byte
	var res []pangya.ServerInfo
	for _, s := range ch.svc.GameServerList() {
		res = append(res, s.Info)
	}

	buf, err := json.Marshal(res)
	if err != nil {
		return err
	}

	_, err = conn.Write(buf)
	return err
}
