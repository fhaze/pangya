package pangya

type PacketHandler interface {
	Action(conn *ConnAccount, req *PacketReader, key uint16) error
}
