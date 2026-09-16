//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"math"
)

type mapDrawableStream struct{ bytes.Buffer }

func (s *mapDrawableStream) u8(v byte) { s.WriteByte(v) }
func (s *mapDrawableStream) u16(v uint16) {
	var p [2]byte
	binary.LittleEndian.PutUint16(p[:], v)
	s.Write(p[:])
}
func (s *mapDrawableStream) u32(v uint32) {
	var p [4]byte
	binary.LittleEndian.PutUint32(p[:], v)
	s.Write(p[:])
}
func (s *mapDrawableStream) f32(v float32) { s.u32(math.Float32bits(v)) }
func (s *mapDrawableStream) fill(n int)    { s.Write(bytes.Repeat([]byte{0x6d}, n)) }

// A complete legacy base body, without the inner-version prefix.
func (s *mapDrawableStream) oldBase(outer, inner int16, team byte, extra uint32) {
	s.oldBaseCount(outer, inner, team, extra, 2)
}
func (s *mapDrawableStream) oldBaseCount(outer, inner int16, team byte, extra uint32, count uint16) {
	s.u32(0x12345678)
	s.u32(0x541)
	if outer < 40 || inner < 4 {
		s.u32(64)
		s.u32(128)
	} else {
		s.f32(64.25)
		s.f32(128.75)
	}
	if outer >= 10 {
		s.u8(3)
		s.WriteString("old")
	}
	if outer >= 20 {
		s.u8(team)
	}
	if outer >= 30 {
		s.u8(0x7a)
	}
	if outer >= 40 {
		s.fill(4)
		if inner >= 2 {
			if inner < 5 {
				s.u32(0x3f800000 | uint32(count))
			} else {
				s.u16(count)
			}
			// Legacy count uses the low 16 bits of the raw float field for inner 2–4.
			s.fill(int(uint16(4 * uint32(count))))
		}
		if inner >= 3 {
			s.u32(extra)
		}
	}
}
func (s *mapDrawableStream) modernBase(inner int16, present bool, team byte, extra uint32) {
	s.modernBaseCount(inner, present, team, extra, 2)
}
func (s *mapDrawableStream) modernBaseCount(inner int16, present bool, team byte, extra uint32, count uint16) {
	s.u16(uint16(inner))
	s.u32(0x12345678)
	s.u32(0x2468ace0)
	s.f32(64.25)
	s.f32(128.75)
	if !present {
		s.u8(0)
		return
	}
	s.u8(1)
	s.u32(0x541)
	s.u8(3)
	s.WriteString("new")
	s.u8(team)
	s.u8(0x7a)
	s.u16(count)
	s.fill(4 * int(count))
	s.u32(extra)
	if inner >= 63 {
		s.fill(2)
		s.u32(3)
		s.fill(3)
		s.fill(4)
	}
	if inner >= 64 {
		s.fill(4)
	}
}
