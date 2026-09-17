//go:build porttest

package opennox

import (
	"github.com/opennox/libs/object"
	"github.com/opennox/libs/types"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"math"
	"testing"
)

func TestSpatialTargetingCursorCandidates(t *testing.T) {
	o := newCollisionCoreOwner(t)
	o.s.PortTestAIEmptyMap()
	t.Cleanup(o.s.Map.Free)
	t.Cleanup(o.s.PortTestRewardTypes([]string{"Polyp", "SpatialTarget"}, nil, true, 0, 0))
	globals, restore := legacy.PortTestSpatialGlobals()
	t.Cleanup(restore)
	engine := noxflags.GetEngine()
	noxflags.ResetEngine()
	t.Cleanup(func() { noxflags.ResetEngine(); noxflags.SetEngine(engine) })
	a, b := &o.units[0], newObjectXferSimple(t, o.s)
	sa, sb := *a, *b
	t.Cleanup(func() { *a = sa; *b = sb })
	ids := collisionCoreIDs(a, b)
	pl := o.units[1].UpdateDataPlayer().Player
	savedPlayer := *pl
	t.Cleanup(func() { *pl = savedPlayer })
	configure, guard, free := o.s.PortTestPathWalls()
	t.Cleanup(free)
	type row struct {
		Variant              string
		Box                  bool
		Z                    uint32
		Offset               [2]float32
		Chosen, Score, Polyp uint32
	}
	var rows []row
	variants := []string{"ordinary", "self", "disabled", "hidden", "invisible", "detect-invisible", "class-none", "polyp", "generator-off", "generator-on", "high-class", "player", "player-missing", "player-spectator", "player-local-headless", "player-local-rendered", "blocked"}
	selected := map[string]int{}
	for _, variant := range variants {
		for _, box := range []bool{false, true} {
			for _, z := range []float32{-101, 0, .25, 10, 100.00001} {
				for _, offset := range [][2]float32{{0, 0}, {9.999, 0}, {10, 0}, {10.001, 0}, {0, 9.999}, {0, 10}, {0, 10.001}, {0, -9.999}, {0, -10}, {7, 7}, {-7, -7}} {
					*a = sa
					*b = sb
					*pl = savedPlayer
					noxflags.ResetEngine()
					configure(0)
					worldGeometryResetObject(a, 1001, 100, 100, false)
					worldGeometryResetObject(b, 1002, 200, 100, box)
					a.ObjClass = 4
					b.ObjClass = 2
					a.Buffs = 0
					b.Buffs = 0
					b.TypeInd = uint16(o.s.Types.ByID("SpatialTarget").Ind())
					b.ZVal = z
					b.Shape.Circle.R2 = 100
					if box {
						b.Shape.Box.W = 20
						b.Shape.Box.H = 20
						b.Shape.Box.Calc()
					}
					target := b
					switch variant {
					case "self":
						target = a
					case "disabled":
						b.ObjFlags = 0x20
					case "hidden":
						b.ObjFlags = 0x8000
					case "invisible":
						b.Buffs = 1
					case "detect-invisible":
						b.Buffs = 1
						a.Buffs = 1 << 21
					case "class-none":
						b.ObjClass = 8
					case "polyp":
						b.ObjClass = 8
						b.TypeInd = uint16(o.s.Types.ByID("Polyp").Ind())
					case "generator-off":
						b.ObjClass = 0x200
					case "generator-on":
						b.ObjClass = 0x200
						b.ObjFlags = 0x4000
					case "high-class":
						b.ObjClass = 0x80000000
					case "player", "player-missing", "player-spectator", "player-local-headless", "player-local-rendered":
						b.ObjClass = 4
						pl.NetCodeVal = 1002
						pl.Field3680 = 0
						if variant == "player-missing" {
							b.NetCode = 987654
						}
						if variant == "player-spectator" {
							pl.Field3680 = 1
						}
						if variant == "player-local-headless" {
							noxflags.SetEngine(noxflags.EngineNoRendering)
						}
					case "blocked":
						configure(1)
					}
					*globals["owner"] = uint32(uintptr(a.CObj()))
					*globals["local"] = b.NetCode
					*globals["chosen"] = 0
					*globals["score"] = 0
					*globals["polyp"] = 0
					mouse := types.Pointf{200 + offset[0], 100 - z + offset[1]}
					before := mouse
					legacy.PortTestSpatialCursorCandidate(target, &mouse)
					chosen := collisionCoreID(t, ids, *globals["chosen"])
					if chosen != 0 {
						selected[variant]++
					}
					if mouse != before || !guard() {
						t.Fatal("cursor changed inputs")
					}
					switch variant {
					case "self", "disabled", "hidden", "invisible", "class-none", "generator-off", "player-missing", "player-spectator", "player-local-headless", "blocked":
						if chosen != 0 {
							t.Fatal("cursor filter admitted target", variant)
						}
					}
					if !box && variant == "ordinary" && z == 0 && offset == ([2]float32{}) && chosen != 1002 {
						t.Fatal("ordinary center not selected")
					}
					if !box && variant == "ordinary" && z == 0 && (offset == ([2]float32{10, 0}) || offset == ([2]float32{0, 10})) && chosen != 0 {
						t.Fatal("circle boundary must be strict")
					}
					if *globals["polyp"] != uint32(o.s.Types.ByID("Polyp").Ind()) {
						t.Fatal("Polyp lookup not initialized")
					}
					rows = append(rows, row{variant, box, math.Float32bits(z), offset, chosen, *globals["score"], *globals["polyp"]})
				}
			}
		}
	}
	for _, name := range []string{"ordinary", "detect-invisible", "polyp", "generator-on", "high-class", "player", "player-local-rendered"} {
		if selected[name] == 0 {
			t.Fatal("admission branch unexercised", name)
		}
	}
	spellbookCapture(t, "spatial-targeting-cursor-candidates", rows, "a7c8cc3cef150a574af4b18d04dea976fc0ee87388486fdd6408af8ecafbc426")
}

func TestSpatialTargetingCursorOrder(t *testing.T) {
	o := newCollisionCoreOwner(t)
	o.s.PortTestAIEmptyMap()
	t.Cleanup(o.s.Map.Free)
	t.Cleanup(o.s.PortTestRewardTypes([]string{"Polyp"}, nil, true, 0, 0))
	globals, restore := legacy.PortTestSpatialGlobals()
	t.Cleanup(restore)
	a, b, c := &o.units[0], newObjectXferSimple(t, o.s), newObjectXferSimple(t, o.s)
	sa, sb, sc := *a, *b, *c
	t.Cleanup(func() { *a = sa; *b = sb; *c = sc })
	ids := collisionCoreIDs(a, b, c)
	player := a.UpdateDataPlayer().Player
	sp := *player
	t.Cleanup(func() { *player = sp })
	type row struct {
		Indexed, Reverse bool
		Height           uint32
		Chosen, Score    uint32
	}
	var rows []row
	for _, indexed := range []bool{false, true} {
		for _, reverse := range []bool{false, true} {
			for _, z := range []float32{-1, 0, 1, .000001} {
				o.s.PortTestAIEmptyMap()
				*a = sa
				*b = sb
				*c = sc
				worldGeometryResetObject(a, 1001, 100, 100, false)
				worldGeometryResetObject(b, 1002, 100, 100, false)
				worldGeometryResetObject(c, 1003, 100, 100, false)
				a.ObjClass = object.ClassPlayer
				b.ObjClass = object.ClassMonster
				c.ObjClass = object.ClassMonster
				b.Shape.Circle.R2 = 100
				c.Shape.Circle.R2 = 100
				c.ZVal = z
				a.ObjFlags = 4
				b.ObjFlags = 4
				c.ObjFlags = 4
				a.Buffs = 0
				b.Buffs = 0
				c.Buffs = 0
				player.CursorVec = image.Pt(100, 100)
				*globals["owner"] = uint32(uintptr(a.CObj()))
				*globals["chosen"] = 0
				*globals["score"] = 0
				first, second := b, c
				if reverse {
					first, second = second, first
				}
				var chosen uint32
				if indexed {
					o.s.Map.AddObjectToIndex(first)
					o.s.Map.AddObjectToIndex(second)
					chosen = legacy.PortTestSpatialCursor(a)
				} else {
					mouse := types.Pointf{100, 100}
					legacy.PortTestSpatialCursorCandidate(first, &mouse)
					legacy.PortTestSpatialCursorCandidate(second, &mouse)
					chosen = *globals["chosen"]
				}
				chosen = collisionCoreID(t, ids, chosen)
				if !indexed && z == 0 {
					want := uint32(1002)
					if reverse {
						want = 1003
					}
					if chosen != want {
						t.Fatal("equal cursor scores must preserve first candidate", chosen, want)
					}
				}
				if z == 1 && chosen != 1003 {
					t.Fatal("higher visible candidate should win", chosen)
				}
				if z == -1 && chosen != 1002 {
					t.Fatal("lower candidate replaced higher", chosen)
				}
				rows = append(rows, row{indexed, reverse, math.Float32bits(z), chosen, *globals["score"]})
			}
		}
	}
	spellbookCapture(t, "spatial-targeting-cursor-order", rows, "0c8f27e61e16bb654395e030667d7adc8a7a002d0b92729228694a8cadac6a25")
}

func TestSpatialTargetingCursorRadius(t *testing.T) {
	o := newCollisionCoreOwner(t)
	o.s.PortTestAIEmptyMap()
	t.Cleanup(o.s.Map.Free)
	t.Cleanup(o.s.PortTestRewardTypes([]string{"Polyp"}, nil, true, 0, 0))
	globals, restore := legacy.PortTestSpatialGlobals()
	t.Cleanup(restore)
	a, b := &o.units[0], newObjectXferSimple(t, o.s)
	sa, sb := *a, *b
	t.Cleanup(func() { *a = sa; *b = sb })
	type row struct {
		Radius, Y     uint32
		Mouse         [2]uint32
		Chosen, Score uint32
	}
	var rows []row
	for _, radius := range []float32{0, .1, .999, 1, 1.001, 1.4142135, 9.999, 10, 10.001, 46340, 65536} {
		for _, y := range []float32{100, 16777216, -16777216} {
			for _, off := range []types.Pointf{{0, 0}, {0, .001}, {0, .5}, {0, 1}, {1, 0}, {0, 1.414}, {0, 9.999}, {0, 10}, {10, 0}, {0, 10.001}} {
				*a = sa
				*b = sb
				worldGeometryResetObject(a, 1001, 100, 100, false)
				worldGeometryResetObject(b, 1002, 100, 100, false)
				// Change height, not map position, to exercise large projected coordinates without invalid map indexing.
				b.ObjClass = 2
				b.Buffs = 0
				b.ZVal = 100 - y
				b.Shape.Circle.R = radius
				*globals["owner"] = uint32(uintptr(a.CObj()))
				*globals["chosen"] = 0
				*globals["score"] = 0
				mouse := types.Pointf{100 + off.X, y + off.Y}
				legacy.PortTestSpatialCursorCandidate(b, &mouse)
				chosen := uint32(0)
				if *globals["chosen"] != 0 {
					if *globals["chosen"] != uint32(uintptr(b.CObj())) {
						t.Fatal("unexpected cursor target")
					}
					chosen = 1002
				}
				if y == 100 && radius == .1 && off == (types.Pointf{0, .001}) && chosen != 0 {
					t.Fatal("fractional squared radius must truncate")
				}
				if y == 100 && radius == .1 && off == (types.Pointf{}) && chosen != 1002 {
					t.Fatal("center-line radius should use strict horizontal interval")
				}
				rows = append(rows, row{math.Float32bits(radius), math.Float32bits(y), motionBits(mouse), chosen, *globals["score"]})
			}
		}
	}
	spellbookCapture(t, "spatial-targeting-cursor-radius", rows, "fa2a8886d6ef60592db60f41126ff31fc5bfd3e1bdba3d430844cacbc08d1a47")
}
