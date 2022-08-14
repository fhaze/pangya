package pangya

import (
	"net"
	"pangya/src/internal/packet"
)

type PacketHandler interface {
	Action(conn net.Conn, pak packet.Packet, key uint16) error
}
