package sync

type PacketBase struct {
	ID string `json:"id"`
}

type ClientPacketHandshake struct {
	PacketBase
	Server string `json:"server"`
	IP     string `json:"ip"`
	Port   int    `json:"port"`
}

type ServerPacketHandshake struct {
	PacketBase
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}
