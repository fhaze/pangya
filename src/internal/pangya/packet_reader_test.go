package pangya

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadUint16(t *testing.T) {
	p := NewPacket(0x01)
	p.PutUint16(0xfe)

	r := NewPacketReader(&p)
	val, err := r.ReadUint16()
	assert.NoError(t, err)
	assert.Equal(t, uint16(0xfe), val)
}

func TestReadUint32(t *testing.T) {
	p := NewPacket(0x01)
	p.PutUint32(0xf8)

	r := NewPacketReader(&p)
	val, err := r.ReadUint32()
	assert.NoError(t, err)
	assert.Equal(t, uint32(0xf8), val)
}

func TestReadLString(t *testing.T) {
	p := NewPacket(0x02)
	p.PutLString("hoge")

	r := NewPacketReader(&p)
	val, err := r.ReadLString()
	assert.NoError(t, err)
	assert.Equal(t, "hoge", val)
}

func TestReadString(t *testing.T) {
	p := NewPacket(0x12)
	p.PutString("hoge", 8)

	r := NewPacketReader(&p)
	val, err := r.ReadString(8)
	assert.NoError(t, err)
	assert.Equal(t, "hoge", val)
}

func TestComplex(t *testing.T) {
	fugai := "fuga"

	p := NewPacket(0x07)
	p.PutUint16(0xf0)
	p.PutLString(fugai)
	p.PutUint32(0x09)

	r := NewPacketReader(&p)
	u2, err := r.ReadUint16()
	assert.NoError(t, err)
	fuga, err := r.ReadLString()
	assert.NoError(t, err)
	u4, err := r.ReadUint32()
	assert.NoError(t, err)

	assert.Equal(t, uint16(0xf0), u2)
	assert.Equal(t, fuga, "fuga")
	assert.Equal(t, uint32(0x09), u4)
}

func TestOverflow(t *testing.T) {
	p := NewPacket(0x00)
	p.PutUint16(0x88)

	r := NewPacketReader(&p)
	_, err := r.ReadUint32()
	assert.Error(t, err)
}
