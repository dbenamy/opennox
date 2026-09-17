//go:build porttest

package opennox

import (
	"encoding/binary"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
	"unsafe"
)

func TestSpatialTargetingObjectProbe(t *testing.T) {
	o := newCollisionCoreOwner(t)
	o.s.PortTestAIEmptyMap()
	t.Cleanup(o.s.Map.Free)
	a, b := newObjectXferSimple(t, o.s), newObjectXferSimple(t, o.s)
	sa, sb := *a, *b
	t.Cleanup(func() { *a = sa; *b = sb })
	ids := collisionCoreIDs(a, b)
	type row struct {
		Box, Indexed, Disabled bool
		Point                  [2]float32
		Return                 uint32
		Next, Prev             [2]uint32
	}
	var rows []row
	hits := 0
	for _, box := range []bool{false, true} {
		for _, indexed := range []bool{false, true} {
			for _, disabled := range []bool{false, true} {
				for _, p := range []types.Pointf{{100, 100}, {109.999, 100}, {110, 100}, {110.001, 100}, {100, 109.999}, {100, 110}, {95, 95}, {105, 105}, {150, 150}} {
					o.s.PortTestAIEmptyMap()
					*a = sa
					*b = sb
					worldGeometryResetObject(a, 1001, 90, 90, false)
					worldGeometryResetObject(b, 1002, 100, 100, box)
					a.ObjClass = 1
					a.ObjFlags = 4
					b.ObjFlags = 4
					b.ObjClass = 0
					a.Collide = o.callback
					b.Collide = o.callback
					a.ObjOwner = nil
					b.ObjOwner = nil
					if disabled {
						b.ObjFlags = 4 | 0x20
					}
					if indexed {
						o.s.Map.AddObjectToIndex(b)
					}
					next, prev := p, types.Pointf{80, 80}
					var rv uint32
					if indexed {
						rv = legacy.PortTestSpatialProbe(a, &next, &prev)
					} else {
						rv = legacy.PortTestSpatialCandidate(a, b, &next, &prev)
					}
					if next != p || prev != (types.Pointf{80, 80}) {
						t.Fatal("ordinary contact changed path")
					}
					if disabled && rv != 0 {
						t.Fatal("disabled object admitted")
					}
					if !disabled && p == (types.Pointf{100, 100}) && rv != uint32(uintptr(b.CObj())) {
						t.Fatal("center contact missed", box, indexed)
					}
					if rv != 0 {
						hits++
					}
					rows = append(rows, row{box, indexed, disabled, [2]float32{p.X, p.Y}, collisionCoreID(t, ids, rv), motionBits(next), motionBits(prev)})
				}
			}
		}
	}
	if hits == 0 {
		t.Fatal("no indexed/ordinary contact")
	}
	spellbookCapture(t, "spatial-targeting-object-probe", rows, "b5cafc1e98549054a825e98ca98988f66cf5e147adf4be6e93293c9b1db52031")
}

func TestSpatialTargetingDoorProbe(t *testing.T) {
	o := newCollisionCoreOwner(t)
	o.s.PortTestAIEmptyMap()
	t.Cleanup(o.s.Map.Free)
	mapDrawableTables(t)
	a, b := newObjectXferSimple(t, o.s), newObjectXferSimple(t, o.s)
	sa, sb := *a, *b
	t.Cleanup(func() { *a = sa; *b = sb })
	ids := collisionCoreIDs(a, b)
	data := collisionCoreGuarded(t, o, 32)
	ud := unsafe.Slice((*byte)(data), 32)
	table := unsafe.Slice((*int32)(memmap.PtrOff(0x587000, 196184)), 64)
	type row struct {
		Direction, Path int
		Indexed, Open   bool
		Delta           [2]int32
		Return          uint32
		Next, Prev      [2]uint32
	}
	var rows []row
	contacts := 0
	for dir := 0; dir < 32; dir++ {
		for path := 0; path < 4; path++ {
			for _, indexed := range []bool{false, true} {
				for _, open := range []bool{false, true} {
					o.s.PortTestAIEmptyMap()
					*a = sa
					*b = sb
					worldGeometryResetObject(a, 1001, 100, 100, false)
					worldGeometryResetObject(b, 1002, 100, 100, false)
					a.ObjClass = 1
					a.ObjFlags = 4
					b.ObjFlags = 4
					b.ObjClass = 0x80
					b.ObjSubClass = 0
					if open {
						b.ObjSubClass = 4
					}
					b.UpdateData = data
					clear(ud)
					binary.LittleEndian.PutUint32(ud[12:], uint32(dir))
					dx, dy := float32(table[2*dir]), float32(table[2*dir+1])
					if dx == 0 && dy == 0 {
						t.Fatal("missing shipped door direction", dir)
					}
					center := types.Pointf{100 + dx/2, 100 + dy/2}
					prev := types.Pointf{center.X - 1, center.Y}
					next := types.Pointf{center.X + 1, center.Y}
					if dy == 0 {
						prev = types.Pointf{center.X, center.Y - 1}
						next = types.Pointf{center.X, center.Y + 1}
					}
					if path == 3 {
						prev = types.Pointf{center.X - dy, center.Y + dx}
						next = types.Pointf{center.X + dy, center.Y - dx}
					}
					if path == 1 {
						prev = types.Pointf{300, 300}
						next = types.Pointf{310, 310}
					}
					if path == 2 {
						next = types.Pointf{100, 100}
					}
					before, previous := next, prev
					if indexed {
						o.s.Map.AddObjectToIndex(b)
					}
					var rv uint32
					if indexed {
						rv = legacy.PortTestSpatialProbe(a, &next, &prev)
					} else {
						rv = legacy.PortTestSpatialCandidate(a, b, &next, &prev)
					}
					if prev != previous {
						t.Fatal("door probe changed previous point")
					}
					if rv != 0 {
						contacts++
						if next != prev {
							t.Fatal("door contact must restore previous point")
						}
					} else if next != before {
						t.Fatal("miss changed proposed point")
					}
					if open && rv != 0 {
						t.Fatal("open door contact")
					}
					if !open && !indexed && path == 0 && rv == 0 {
						t.Fatal("axis midpoint crossing missed", dir)
					}
					rows = append(rows, row{dir, path, indexed, open, [2]int32{table[2*dir], table[2*dir+1]}, collisionCoreID(t, ids, rv), motionBits(next), motionBits(prev)})
				}
			}
		}
	}
	if contacts == 0 {
		t.Fatal("no door intersections")
	}
	spellbookCapture(t, "spatial-targeting-door-probe", rows, "a5ee959f5dd55969a065cc80dca19b986493c33ad73bf2a03849e74312f1bce6")
}

func TestSpatialTargetingRay(t *testing.T) {
	o := newCollisionCoreOwner(t)
	o.s.PortTestAIEmptyMap()
	t.Cleanup(o.s.Map.Free)
	configure, unchanged, free := o.s.PortTestPathWalls()
	t.Cleanup(free)
	type row struct {
		Wall   int
		Ray    [4]float32
		Return int32
	}
	var rows []row
	for mode := 0; mode < 3; mode++ {
		configure(mode)
		for _, ray := range [][4]float32{{100, 100, 200, 100}, {200, 100, 100, 100}, {100, 100, 100, 100}, {100, 50, 200, 50}, {100, 80, 200, 120}} {
			before := ray
			rv := legacy.PortTestSpatialRay(&ray)
			if ray != before || !unchanged() {
				t.Fatal("ray changed inputs")
			}
			if mode == 0 && rv != 1 {
				t.Fatal("empty map rejected ray")
			}
			if mode == 1 && before == ([4]float32{100, 100, 200, 100}) && rv != 0 {
				t.Fatal("solid wall admitted ray")
			}
			rows = append(rows, row{mode, ray, rv})
		}
	}
	spellbookCapture(t, "spatial-targeting-ray", rows, "8a1ed724c234c39e043ce590d5d1fa1d37e3655d30966a40dcdc3438b9306336")
}
