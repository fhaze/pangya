package pangya

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPacketID(t *testing.T) {
	p := NewPacket(0x07)
	assert.Equal(t, uint16(0x07), p.ID)
}

func TestPacketNoID(t *testing.T) {
	p := NewPacket()
	p.PutUint16(0x25)
	assert.Equal(t, uint16(0x00), p.ID)
	assert.Equal(t, []byte{0x25, 0x00}, p.ToBytes())
}

func TestPutUint8(t *testing.T) {
	p := NewPacket(0x04)
	p.PutUint8(0x08)
	assert.Equal(t, []byte{0x08}, p.Payload)
}

func TestPutUint16(t *testing.T) {
	p := NewPacket(0x88)
	p.PutUint16(0x0f)
	assert.Equal(t, []byte{0x0f, 0x00}, p.Payload)
}

func TestPutUint32(t *testing.T) {
	p := NewPacket(0x55)
	p.PutUint32(0x0f)
	assert.Equal(t, []byte{0x0f, 0x00, 0x00, 0x00}, p.Payload)
}

func TestPutString(t *testing.T) {
	p := NewPacket(0x10)
	p.PutString("hoge", 8)
	b := make([]byte, 8)
	copy(b, "hoge")
	assert.Equal(t, b, p.Payload)
}

func TestPutLString(t *testing.T) {
	p := NewPacket(0x13)
	p.PutLString("hoge")
	var exp []byte
	exp = append(exp, 0x04, 0x00)
	exp = append(exp, []byte("hoge")...)
	assert.Equal(t, exp, p.Payload)
}

func TestPacketFromBytes(t *testing.T) {
	p, err := PacketFromBytes([]byte{0x01, 0x00, 0x07, 0x00})
	assert.NoError(t, err)
	assert.Equal(t, uint16(0x01), p.ID)
	assert.Equal(t, []byte{0x07, 0x00}, p.Payload)
}

func TestPacketToBytes(t *testing.T) {
	p := NewPacket(0x01)
	p.PutUint16(0xf0)
	assert.Equal(t, []byte{0x01, 0x00, 0xf0, 0x00}, p.ToBytes())
}
