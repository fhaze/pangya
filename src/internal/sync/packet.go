package sync

import "pangya/src/internal/pangya"

type PacketBase struct {
	ID string `json:"id"`
}

type ClientPacketHandshake struct {
	PacketBase
	Server pangya.ServerInfo `json:"server"`
	IP     string            `json:"ip"`
	Port   int               `json:"port"`
}

type ServerPacketHandshake struct {
	PacketBase
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}
