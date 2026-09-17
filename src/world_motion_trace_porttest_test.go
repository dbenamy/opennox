//go:build porttest

package opennox

import (
	"github.com/opennox/libs/object"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/legacy"
	"math"
	"testing"
)

func TestWorldMotionProjectileTrace(t *testing.T) {
	o := newCollisionCoreOwner(t)
	o.s.PortTestAIEmptyMap()
	t.Cleanup(o.s.Map.Free)
	configure, unchanged, free := o.s.PortTestPathWalls()
	t.Cleanup(free)
	globals, restore := legacy.PortTestWorldMotionGlobals()
	t.Cleanup(restore)
	a, b := newObjectXferSimple(t, o.s), newObjectXferSimple(t, o.s)
	sa, sb := *a, *b
	t.Cleanup(func() { *a = sa; *b = sb })
	ids := collisionCoreIDs(a, b)
	type row struct {
		Wall, Target        int
		Step                uint32
		Return              int8
		Hit                 uint32
		Normal, Proposed    [2]uint32
		Trace, GridX, GridY uint32
	}
	var rows []row
	for wall := 0; wall < 3; wall++ {
		for target := 0; target < 4; target++ {
			for _, step := range []float32{0, 5.999, 6, 6.001, 12, 100} {
				o.s.PortTestAIEmptyMap()
				configure(wall)
				*a = sa
				*b = sb
				worldGeometryResetObject(a, 1001, 100, 100, false)
				worldGeometryResetObject(b, 1002, 105, 100, false)
				a.ObjFlags = 4
				a.Collide = o.callback
				a.NewPos = types.Pointf{100 + step, 100}
				b.ObjClass = 0
				b.ObjFlags = 4
				b.Collide = o.callback
				b.Shape.Circle.R = 2
				b.Shape.Circle.R2 = 4
				if target == 2 {
					b.PosVec = types.Pointf{150, 100}
					b.NewPos = b.PosVec
				}
				if target == 3 {
					b.PosVec = types.Pointf{180, 100}
					b.NewPos = b.PosVec
				}
				if target != 0 {
					o.s.Map.AddObjectToIndex(b)
				}
				*globals["trace"] = 0
				*globals["grid-x"] = 0x1234
				*globals["grid-y"] = 0x5678
				hit := uint32(0xabcdef01)
				normal := types.Pointf{17, -19}
				rv := legacy.PortTestWorldMotionTrace(a, &hit, &normal)
				if wall == 0 && target == 0 && (rv != 0 || hit != 0xabcdef01 || normal != (types.Pointf{17, -19}) || a.NewPos != (types.Pointf{100 + step, 100})) {
					t.Fatal("unobstructed projectile trace modified outputs")
				}
				if wall == 1 && target == 0 && step == 100 && (rv != 1 || hit != 0 || a.NewPos != a.PosVec || *globals["trace"] != 1) {
					t.Fatal("wall trace must report normal and reset proposed position")
				}
				if wall == 0 && target == 1 && step >= 5.999 && (rv != 1 || hit != uint32(uintptr(b.CObj()))) {
					t.Fatal("projectile sampling missed indexed object", step, rv, hit)
				}
				if hit != 0xabcdef01 {
					hit = collisionCoreID(t, ids, hit)
				}
				if !unchanged() {
					t.Fatal("projectile trace changed wall")
				}
				rows = append(rows, row{wall, target, math.Float32bits(step), rv, hit, motionBits(normal), motionBits(a.NewPos), *globals["trace"], *globals["grid-x"], *globals["grid-y"]})
			}
		}
	}
	spellbookCapture(t, "world-motion-projectile-trace", rows, "1049b231aa90410dfac004bc87a860a2dd75f47d94ac84ddb80ca80061ceec3c")
}

func TestWorldMotionProjectileDispatch(t *testing.T) {
	o := newCollisionCoreOwner(t)
	o.s.PortTestAIEmptyMap()
	t.Cleanup(o.s.Map.Free)
	configure, unchanged, free := o.s.PortTestPathWalls()
	t.Cleanup(free)
	globals, restore := legacy.PortTestWorldMotionGlobals()
	t.Cleanup(restore)
	a, b := newObjectXferSimple(t, o.s), newObjectXferSimple(t, o.s)
	sa, sb := *a, *b
	t.Cleanup(func() { *a = sa; *b = sb })
	ids := collisionCoreIDs(a, b)
	type row struct {
		Wall     int
		Type     string
		Flags    uint32
		Calls    [][4]uint32
		Proposed [2]uint32
		Trace    uint32
		Cache    [3]uint32
	}
	var rows []row
	for wall := 0; wall < 2; wall++ {
		for _, typ := range []string{"none", "normal", "SmallFist", "MediumFist", "LargeFist"} {
			for _, flags := range []object.Flags{4, 4 | 0x20, 4 | 0x40, 4 | 0x60} {
				o.s.PortTestAIEmptyMap()
				configure(wall)
				o.resetCalls()
				*a = sa
				*b = sb
				for _, name := range []string{"fist-small", "fist-medium", "fist-large"} {
					*globals[name] = 0
				}
				*globals["trace"] = 0x1234
				worldGeometryResetObject(a, 1001, 100, 100, false)
				worldGeometryResetObject(b, 1002, 105, 100, false)
				a.ObjFlags = flags
				a.Collide = o.callback
				a.NewPos = types.Pointf{200, 100}
				b.ObjClass = 0
				b.ObjFlags = 4
				b.Collide = o.callback
				b.Shape.Circle.R = 3
				b.Shape.Circle.R2 = 9
				if typ != "normal" && typ != "none" {
					b.TypeInd = uint16(o.s.Types.IndByID(typ))
				}
				if typ != "none" {
					o.s.Map.AddObjectToIndex(b)
				}
				legacy.PortTestWorldMotionDispatch(a)
				calls := o.calls(ids)
				if flags&0x60 != 0 {
					if len(calls) != 0 || *globals["trace"] != 0x1234 {
						t.Fatal("disabled projectile dispatch")
					}
				} else if typ == "normal" {
					if len(calls) != 2 || calls[0][0] != 1001 || calls[0][1] != 1002 || calls[1][0] != 1002 || calls[1][1] != 1001 || calls[1][2] != calls[0][2]^0x80000000 || calls[1][3] != calls[0][3]^0x80000000 {
						t.Fatal("projectile bilateral callback order and reversed normal", calls)
					}
					if *globals["trace"] != 0 {
						t.Fatal("dispatch must clear wall trace marker after first callback")
					}
				} else if typ != "none" && len(calls) != 0 {
					t.Fatal("fist target must suppress callbacks", typ, calls)
				}
				if !unchanged() {
					t.Fatal("projectile dispatch changed wall")
				}
				cache := [3]uint32{*globals["fist-small"], *globals["fist-medium"], *globals["fist-large"]}
				for i, name := range []string{"SmallFist", "MediumFist", "LargeFist"} {
					if cache[i] == 0 || cache[i] != uint32(o.s.Types.IndByID(name)) {
						t.Fatal("projectile cache initialization precedes disabled guard")
					}
				}
				rows = append(rows, row{wall, typ, uint32(flags), calls, motionBits(a.NewPos), *globals["trace"], cache})
			}
		}
	}
	spellbookCapture(t, "world-motion-projectile-dispatch", rows, "478ccf83f4cf473d4d318cba128c7d4b619d471ac915ba6248a14c6e0b441645")
}
