package pangya

import "net"

type PacketHandler interface {
	Action(conn net.Conn, pak Packet, key uint16) error
}
