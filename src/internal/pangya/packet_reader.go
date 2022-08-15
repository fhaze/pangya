package pangya

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type PacketReader struct {
	i uint
	p *Packet
}

func (pr *PacketReader) ReadUint16() (uint16, error) {
	payload := pr.p.Payload[pr.i : pr.i+2]
	if err := pr.move(2); err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint16(payload), nil
}

func (pr *PacketReader) ReadUint32() (uint32, error) {
	payload := pr.p.Payload[pr.i : pr.i+4]
	if err := pr.move(4); err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint32(payload), nil
}

func (pr *PacketReader) ReadLString() (string, error) {
	payloadLength := binary.LittleEndian.Uint16(pr.p.Payload[pr.i : pr.i+2])
	if err := pr.move(2); err != nil {
		return "", err
	}
	payloadData := pr.p.Payload[pr.i : pr.i+uint(payloadLength)]
	if err := pr.move(uint(payloadLength)); err != nil {
		return "", err
	}
	return string(payloadData), nil
}

func (pr *PacketReader) ReadString(l uint) (string, error) {
	payload := pr.p.Payload[pr.i : pr.i+l]
	if err := pr.move(l); err != nil {
		return "", err
	}
	return string(bytes.Trim(payload, "\x00")), nil
}

func (pr *PacketReader) move(i uint) error {
	pr.i += i
	if pr.i > uint(len(pr.p.Payload)) {
		return fmt.Errorf("read pointer overflow %d of %d", i, len(pr.p.Payload))
	}
	return nil
}

func NewPacketReader(p *Packet) *PacketReader {
	return &PacketReader{0, p}
}
