package server

import (
	"net"
	"pangya/internal/packet"
)

type PacketHandler interface {
	Action(conn net.Conn, pak packet.Packet, key uint16) error
}
