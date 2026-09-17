//go:build porttest

package server

import (
	"github.com/opennox/libs/wall"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"image"
	"unsafe"
)

// PortTestGeometryWalls owns actual wall records and the lookup buckets used by
// collision. It supplies all eleven shipped direction shapes and optional neighbours.
func (s *Server) PortTestGeometryWalls() (configure func(int, bool, bool), check func() bool, restore func()) {
	oldPos, oldDef, oldCount := s.Walls.byPos, s.Walls.defs[0], s.Walls.defsCnt
	s.Walls.byPos = make([]*Wall, wallsPerBucket*WallGridSize)
	s.Walls.defsCnt = 1
	s.Walls.defs[0] = WallDef{}
	raw, free := alloc.Make([]byte{}, 16+9*int(unsafe.Sizeof(Wall{})))
	records := unsafe.Slice((*Wall)(unsafe.Pointer(&raw[8])), 9)
	configure = func(dir int, neighbours, window bool) {
		clear(s.Walls.byPos)
		clear(raw)
		for i := 0; i < 8; i++ {
			raw[i] = 0xa5
			raw[len(raw)-8+i] = 0x5a
		}
		k := 0
		for dy := -1; dy <= 1; dy++ {
			for dx := -1; dx <= 1; dx++ {
				wl := &records[k]
				k++
				*wl = Wall{X5: uint8(6 + dx), Y6: uint8(4 + dy), Dir0: uint8(max(dir, 0))}
				if window {
					wl.Flags4 = wall.FlagWindow
				}
				if dir >= 0 && (neighbours || dx == 0 && dy == 0) {
					s.Walls.addByPos(wl, image.Pt(6+dx, 4+dy))
				}
			}
		}
	}
	check = func() bool {
		for i := 0; i < 8; i++ {
			if raw[i] != 0xa5 || raw[len(raw)-8+i] != 0x5a {
				return false
			}
		}
		return true
	}
	restore = func() { s.Walls.byPos, s.Walls.defs[0], s.Walls.defsCnt = oldPos, oldDef, oldCount; free() }
	return
}
