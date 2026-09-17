//go:build porttest

package opennox

import (
	"github.com/opennox/libs/object"
	"github.com/opennox/libs/types"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
	"unsafe"
)

func TestCollisionCoreCircleWallScan(t *testing.T) {
	w := newWorldCollisionOwner(t)
	worldGeometryTables(t)
	t.Cleanup(noxflags.PortTestGameFlags(0))
	w.s.PortTestAIEmptyMap()
	configure, guard, free := w.s.PortTestGeometryWalls()
	t.Cleanup(free)
	reset, hits, restore := legacy.PortTestWorldGeometryHitOwner()
	t.Cleanup(restore)
	a := newObjectXferSimple(t, w.s)
	ids := map[unsafe.Pointer]uint32{a.CObj(): 1001}
	type row struct {
		Op                 string
		Direction          int
		Neighbours, Window bool
		Offset             [2]float32
		Variant            int
		Return             int
		Object             worldGeometryObjectState
		Hits               [][5]uint32
	}
	var rows []row
	nonzero := map[string]int{}
	for _, op := range []string{"circle-walls"} {
		for dir := -1; dir < 11; dir++ {
			for _, neighbours := range []bool{false, true} {
				for _, window := range []bool{false, true} {
					for _, offset := range [][2]float32{{0, 0}, {11.5, 11.5}, {23, 23}, {0, 23}, {23, 0}, {-5, 11.5}, {28, 11.5}, {11.5, -5}, {11.5, 28}, {100, 100}} {
						for variant := 0; variant < 4; variant++ {
							configure(dir, neighbours, window)
							reset()
							worldGeometryResetObject(a, 1001, 138+offset[0], 92+offset[1], false)
							a.VelVec = types.Pointf{3, -4}
							switch variant {
							case 1:
								a.ObjFlags = 0x4000
								a.Shape.Circle.R = 8
							case 2:
								a.ObjClass = object.Class(0x400000)
							case 3:
								a.PrevPos = types.Pointf{138 - offset[0], 92 - offset[1]}
								a.Mass = 7
							}
							a.UpdateCollider(a.NewPos)
							rv := int(legacy.PortTestCollisionCore(op, a, nil, nil, 0))
							h := hits(ids)
							if len(h) > 0 {
								nonzero[op]++
							}
							if !guard() {
								t.Fatal("wall record guards")
							}
							if dir == -1 && (len(h) != 0 || a.Pos24 != (types.Pointf{})) {
								t.Fatal("empty wall grid changed physics", op)
							}
							rows = append(rows, row{op, dir, neighbours, window, offset, variant, rv, worldGeometryState(a), h})
						}
					}
				}
			}
		}
	}
	for _, op := range []string{"circle-walls"} {
		if nonzero[op] == 0 {
			t.Fatal("no actual wall contact exercised", op)
		}
	}
	spellbookCapture(t, "collision-core-circle-wall-scan", rows, "c1316eaa461d5273388d9989192314c6fae14657ab90bef60f4d949dc55bc8ef")
}
