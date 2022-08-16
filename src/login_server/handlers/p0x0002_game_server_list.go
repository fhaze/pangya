package handlers

import (
	"net"
	"pangya/src/internal/pangya"
	"pangya/src/internal/sync"
	syncclient "pangya/src/sync_client"
)

type P0x0002_GameServerList struct {
	syc syncclient.Client
}

func NewP0x0002_GameServerList(syc syncclient.Client) pangya.PacketHandler {
	return &P0x0002_GameServerList{syc: syc}
}

func (h *P0x0002_GameServerList) Action(conn net.Conn, req pangya.Packet, key uint16) error {
	return h.syc.Request(&sync.PacketBase{ID: sync.PacketGameServerList})
}
