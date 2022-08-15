package sync

import "net"

type PacketHandler interface {
	Action(conn net.Conn, pak []byte) error
}
