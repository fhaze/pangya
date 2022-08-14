package packet

import "encoding/binary"

type Packet struct {
	ID      uint16
	Payload []byte
}

func (p *Packet) PutUint8(v uint8) *Packet {
	p.Payload = append(p.Payload, v)
	return p
}

func (p *Packet) PutUint16(v uint16) *Packet {
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, v)
	p.Payload = append(p.Payload, b...)
	return p
}

func (p *Packet) PutUint32(v uint32) *Packet {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, v)
	p.Payload = append(p.Payload, b...)
	return p
}

func (p *Packet) PutLString(v string) *Packet {
	p.PutUint16(uint16(len(v)))
	b := make([]byte, len(v))
	copy(b, []byte(v))
	p.Payload = append(p.Payload, b...)
	return p
}

func (p *Packet) PutString(v string, l uint) *Packet {
	b := make([]byte, l)
	copy(b, []byte(v))
	p.Payload = append(p.Payload, b...)
	return p
}

func (p *Packet) PutBytes(b []byte) *Packet {
	p.Payload = append(p.Payload, b...)
	return p
}

func (p *Packet) ToBytes() []byte {
	id := make([]byte, 2)
	binary.LittleEndian.PutUint16(id, p.ID)
	var b []byte
	b = append(b, id...)
	b = append(b, p.Payload...)
	return b
}

func NewPacket(id uint16) Packet {
	return Packet{ID: id}
}

func FromBytes(b []byte) (Packet, error) {
	if len(b) < 2 {
		return Packet{}, nil
	}
	return Packet{
		ID:      binary.LittleEndian.Uint16(b[:2]),
		Payload: b[2:],
	}, nil
}
