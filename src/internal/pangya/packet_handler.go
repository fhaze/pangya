package pangya

import "net"

type PacketHandler interface {
	Action(conn net.Conn, req Packet, key uint16) error
}
