//go:build porttest

package server

import (
	"bytes"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"image"
	"unsafe"
)

// PortTestPathWalls supplies a real diagonal wall across the horizontal ray
// (100,100) to (200,100), plus an explicitly transparent definition variant.
func (s *Server) PortTestPathWalls() (configure func(int), unchanged func() bool, free func()) {
	raw, release := alloc.Make([]byte{}, int(unsafe.Sizeof(Wall{}))+16)
	wl := (*Wall)(unsafe.Pointer(&raw[8]))
	oldDef, oldCount := s.Walls.defs[0], s.Walls.defsCnt
	var before []byte
	configure = func(mode int) {
		clear(s.Walls.byPos)
		clear(raw)
		for i := 0; i < 8; i++ {
			raw[i] = 0xa5
			raw[len(raw)-8+i] = 0x5a
		}
		*wl = Wall{X5: 6, Y6: 4, Dir0: 2}
		s.Walls.defsCnt = 1
		s.Walls.defs[0] = WallDef{}
		if mode == 2 {
			s.Walls.defs[0].Flags32 = 2
		}
		if mode != 0 {
			s.Walls.addByPos(wl, image.Pt(6, 4))
		}
		before = bytes.Clone(raw)
	}
	unchanged = func() bool { return bytes.Equal(raw, before) }
	free = func() { clear(s.Walls.byPos); s.Walls.defs[0], s.Walls.defsCnt = oldDef, oldCount; release() }
	return
}
