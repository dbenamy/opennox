//go:build porttest

package opennox

import (
	"github.com/opennox/libs/object"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"testing"
	"unsafe"
)

func TestWorldGeometryPlayerWallContact(t *testing.T) {
	w := newWorldCollisionOwner(t)
	worldGeometryTables(t)
	t.Cleanup(noxflags.PortTestGameFlags(0))
	w.s.PortTestAIEmptyMap()
	configure, guards, free := w.s.PortTestGeometryWalls()
	t.Cleanup(free)
	reset, hits, restore := legacy.PortTestWorldGeometryHitOwner()
	t.Cleanup(restore)
	a := newObjectXferSimple(t, w.s)
	ids := map[unsafe.Pointer]uint32{a.CObj(): 1001}
	raw := unsafe.Slice((*byte)(w.record(t, 336)), 336)
	ud := unsafe.Pointer(&raw[8])
	oldUD := a.UpdateData
	t.Cleanup(func() { a.UpdateData = oldUD; a.ObjClass = object.ClassSimple })
	type row struct {
		Direction int
		Offset    [2]float32
		HasData   bool
		Return    int
		Wall      bool
		Object    worldGeometryObjectState
		Hits      [][5]uint32
	}
	var rows []row
	stored, withoutHit := 0, 0
	for dir := -1; dir < 11; dir++ {
		for _, offset := range [][2]float32{{0, 0}, {11.5, 11.5}, {23, 0}, {0, 23}, {100, 100}} {
			for _, hasData := range []bool{false, true} {
				configure(dir, false, false)
				reset()
				clear(raw)
				for i := 0; i < 8; i++ {
					raw[i] = 0xa5
					raw[328+i] = 0x5a
				}
				worldGeometryResetObject(a, 1001, 138+offset[0], 92+offset[1], false)
				a.ObjClass = object.ClassPlayer
				a.UpdateData = nil
				if hasData {
					a.UpdateData = ud
				}
				var scratch [16]uint32
				scratch[0] = 6
				scratch[1] = 4
				rv := legacy.PortTestWorldGeometryPhysics("circle-wall", a, nil, &scratch, 0, 0)
				wallPtr := *(*unsafe.Pointer)(unsafe.Add(ud, 296))
				if wallPtr != nil {
					expected := w.s.Walls.GetWallAtGrid(image.Pt(6, 4))
					if expected == nil || wallPtr != expected.C() {
						t.Fatal("player wall pointer not the real indexed wall")
					}
					stored++
					if rv == 0 {
						withoutHit++
					}
				}
				for i, v := range raw[8:328] {
					if i >= 296 && i < 300 {
						continue
					}
					if v != 0 {
						t.Fatal("unexpected player update mutation", i, v)
					}
				}
				for i := 0; i < 8; i++ {
					if raw[i] != 0xa5 || raw[328+i] != 0x5a {
						t.Fatal("player data guard")
					}
				}
				if !guards() {
					t.Fatal("wall guards")
				}
				rows = append(rows, row{dir, offset, hasData, rv, wallPtr != nil, worldGeometryState(a), hits(ids)})
			}
		}
	}
	if stored == 0 || withoutHit == 0 {
		t.Fatal("wall association must cover projection before contact response", stored, withoutHit)
	}
	spellbookCapture(t, "world-geometry-player-wall", rows, "6603a57f77b51e154df0f0f23ec6ea284fd9a71f55944b2dbd67f09633c592ee")
}
