//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"reflect"
	"testing"
)

func TestTileCompositionFull(t *testing.T) {
	o := newTileCompositionOwner(t)
	type row struct {
		Fill     bool
		Pattern  int
		Position image.Point
		State    tileCompositionState
	}
	var rows []row
	positions := []image.Point{{-80, -80}, {-1, -1}, {0, 0}, {10, 10}, {11, 11}, {56, 56}, {57, 57}, {230, 276}, {5500, 5500}, {5800, 5800}, {6000, 6000}}
	for _, fill := range []bool{false, true} {
		for _, pattern := range []int{0, 1, 2, 3} {
			for _, pos := range positions {
				o.populate(pattern)
				o.resetBuffer()
				*o.counter = 77
				for i := range o.usage {
					o.usage[i] = 0xa5
				}
				for i := range o.edgeUsage {
					o.edgeUsage[i] = 0x5a
				}
				nox_client_texturedFloors_154956 = !fill
				nox_xxx_tileSetDrawFn_481420()
				vp := tileCompositionViewport(pos.X, pos.Y)
				before := *vp
				gridBefore := o.gridHash()
				legacy.PortTestTileCompositionFull(vp)
				if o.gridHash() != gridBefore {
					t.Fatal("full redraw mutated grid")
				}
				if *vp != before {
					t.Fatal("full redraw mutated viewport")
				}
				if *o.counter != 0xffffffff || *o.words["scrollX"] != 0 || *o.words["scrollY"] != 0 {
					t.Fatal("full redraw did not reset counter/scroll")
				}
				if pattern == 0 {
					for _, n := range o.usage {
						if n != 0 {
							t.Fatal("empty grid has used tiles")
						}
					}
				}
				for _, n := range o.edgeUsage {
					if n != 0 {
						t.Fatal("edgeless grid has used edges")
					}
				}
				rows = append(rows, row{fill, pattern, pos, o.state()})
			}
		}
	}
	tileCompositionCapture(t, "full", rows, "504daeb5c996ed08ddb9f5df460b449c0589b3645baf67beef0d9b235fe62d82")
}
func TestTileCompositionScroll(t *testing.T) {
	o := newTileCompositionOwner(t)
	type row struct {
		Fill        bool
		Start       image.Point
		Axis, Delta int
		State       tileCompositionState
	}
	var rows []row
	starts := []image.Point{{0, 0}, {230, 276}, {5500, 5500}, {5800, 5800}}
	deltas := []int{-100, -47, -46, -24, -23, -22, -1, 0, 22, 23, 24, 45, 46, 47, 100, 180, 181, 226, 227, 280}
	for _, fill := range []bool{false, true} {
		for _, start := range starts {
			for axis := 0; axis < 2; axis++ {
				for _, delta := range deltas {
					o.populate(1)
					o.resetBuffer()
					nox_client_texturedFloors_154956 = !fill
					nox_xxx_tileSetDrawFn_481420()
					vp := tileCompositionViewport(start.X, start.Y)
					legacy.PortTestTileCompositionFull(vp)
					if axis == 0 {
						vp.World = vp.World.Add(image.Pt(delta, 0))
						legacy.PortTestTileCompositionHorizontal(vp, start.X+delta)
					} else {
						vp.World = vp.World.Add(image.Pt(0, delta))
						legacy.PortTestTileCompositionVertical(vp, start.Y+delta)
					}
					rows = append(rows, row{fill, start, axis, delta, o.state()})
				}
			}
		}
	}
	tileCompositionCapture(t, "scroll", rows, "f084ce81b75b785fa1cf2340951eadaf101c926d2b2f1f6ffeaad189f1f3e8aa")
}
func TestTileCompositionPredicate(t *testing.T) {
	o := newTileCompositionOwner(t)
	type row struct {
		Position, Cell              image.Point
		Half, Enabled, Flag, Return int
	}
	var rows []row
	positions := []image.Point{{-80, -80}, {0, 0}, {57, 103}, {230, 276}, {5800, 5800}, {6000, 6000}}
	for _, pos := range positions {
		loX := int((uint32(pos.X) - 11) / 46)
		hiX := int(*o.words["columns"]) + loX - 1
		if hiX >= 128 {
			hiX = 127
			loX = 127 - int(*o.words["columns"])
		}
		loY := max(0, (pos.Y-11)/46-1)
		hiY := int(*o.words["rows"]) + loY
		if hiY >= 128 {
			hiY = 127
			loY = 127 - int(*o.words["rows"])
		}
		cells := []image.Point{{loX, loY}, {hiX - 1, hiY - 1}, {hiX, loY}, {loX, hiY}, {max(0, loX-1), loY}, {loX, max(0, loY-1)}}
		for _, cell := range cells {
			for half := 0; half < 2; half++ {
				for enabled := 0; enabled < 4; enabled++ {
					for flag := 0; flag < 4; flag++ {
						o.populate(0)
						c := &o.grid[cell.X][cell.Y]
						c[0] = uint32(enabled)
						c[1+5*half] = 1
						o.defs[1].Field58 = byte(flag)
						before := o.state()
						ret := legacy.PortTestTileCompositionRedraw(tileCompositionViewport(pos.X, pos.Y))
						want := 0
						if cell.X >= loX && cell.X < hiX && cell.Y >= loY && cell.Y < hiY && enabled&(1<<half) != 0 && flag&1 != 0 {
							want = 1
						}
						if ret != want {
							t.Fatalf("position%v cell%v half%d enabled%d flag%d got%d want%d", pos, cell, half, enabled, flag, ret, want)
						}
						if !reflect.DeepEqual(before, o.state()) {
							t.Fatal("predicate mutated renderer state")
						}
						rows = append(rows, row{pos, cell, half, enabled, flag, ret})
					}
				}
			}
		}
	}
	tileCompositionCapture(t, "predicate", rows, "0d23c8d5d142fff19f2ab5df3844d0a30fd7663b18530c0464d1c27006836064")
}

func TestTileCompositionScrollSequence(t *testing.T) {
	o := newTileCompositionOwner(t)
	type row struct {
		Fill            bool
		Direction, Step int
		State           tileCompositionState
	}
	var rows []row
	for _, fill := range []bool{false, true} {
		o.populate(2)
		o.resetBuffer()
		nox_client_texturedFloors_154956 = !fill
		nox_xxx_tileSetDrawFn_481420()
		vp := tileCompositionViewport(460, 460)
		legacy.PortTestTileCompositionFull(vp)
		for direction := 0; direction < 4; direction++ {
			for step := 0; step < 8; step++ {
				ox, oy := int(int32(*o.words["originX"])), int(int32(*o.words["originY"]))
				before := o.state()
				switch direction {
				case 0:
					vp = tileCompositionViewport(ox+150, oy+40)
					legacy.PortTestTileCompositionHorizontal(vp, ox+150)
				case 1:
					vp = tileCompositionViewport(ox+40, oy+150)
					legacy.PortTestTileCompositionVertical(vp, oy+150)
				case 2:
					vp = tileCompositionViewport(ox, oy+40)
					legacy.PortTestTileCompositionHorizontal(vp, ox)
				case 3:
					vp = tileCompositionViewport(ox+40, oy)
					legacy.PortTestTileCompositionVertical(vp, oy)
				}
				dx, dy := 0, 0
				switch direction {
				case 0:
					dx = 46
				case 1:
					dy = 46
				case 2:
					dx = -46
				case 3:
					dy = -46
				}
				if int(int32(*o.words["originX"])) != ox+dx || int(int32(*o.words["originY"])) != oy+dy {
					t.Fatalf("scroll transition direction%d step%d", direction, step)
				}
				if before.Pixels == o.state().Pixels {
					t.Fatalf("populated scroll did not write pixels direction%d step%d", direction, step)
				}
				rows = append(rows, row{fill, direction, step, o.state()})
			}
		}
	}
	tileCompositionCapture(t, "scroll-sequence", rows, "b6014d8bee88c4473779b42fb6cb38d9de4d04c2ac1e462e2ccc6f1bb7824952")
}
