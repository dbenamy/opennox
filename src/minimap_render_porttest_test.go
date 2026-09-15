//go:build porttest

package opennox

import (
	"bytes"
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"testing"
	"unsafe"
)

type minimapResult struct {
	Op     int
	Return uint32
	Words  map[string]uint32
	Cache  []uint32
	Render objectRenderResult
}

func (o *minimapOwner) capture(t *testing.T, op int, ret uint32) minimapResult {
	named := map[string]uint32{}
	for n, p := range o.words {
		named[n] = *p
	}
	cache := append([]uint32(nil), unsafe.Slice(memmap.PtrUint32(0x5D4594, 1096300), 5)...)
	cache = append(cache, memmap.Uint32(0x5D4594, 1096424), memmap.Uint32(0x5D4594, 588080))
	return minimapResult{op, ret, named, cache, o.renderResult(t, op, 0, int(ret))}
}
func (o *minimapOwner) player(t *testing.T, index int, pos image.Point) *client.Drawable {
	t.Helper()
	dr := o.drawable(7+index, pos)
	if dr == nil {
		t.Fatal("actual player drawable allocation")
	}
	dr.ObjClass = 4
	o.c.Objs.PlayerListAdd(dr)
	if index == 0 {
		*memmap.PtrPtr(0x852978, 8) = unsafe.Pointer(dr)
	}
	return dr
}
func TestMinimapFloorCache(t *testing.T) {
	o := newMinimapOwner(t)
	var rows []minimapResult
	for _, local := range []bool{false, true} {
		o.resetMinimap(t)
		o.polygon(t, 0, image.Rect(191, 190, 271, 270), 2)
		o.polygon(t, 1, image.Rect(401, 400, 501, 500), 3)
		dr := o.player(t, 0, image.Pt(230, 230))
		if !local {
			*memmap.PtrPtr(0x852978, 8) = nil
		}
		for _, tc := range []struct {
			pos  image.Point
			want uint32
		}{{image.Pt(230, 230), 2}, {image.Pt(450, 450), 3}, {image.Pt(800, 800), 3}, {image.Pt(230, 230), 2}} {
			if local {
				o.c.Viewport().World.Max = tc.pos
			} else { // Move through the real spatial owner rather than changing indexed coordinates.
				o.c.Nox_xxx_updateSpritePosition_49AA90(dr, tc.pos.X, tc.pos.Y)
			}
			ret := legacy.PortTestMinimap(11, int(uintptr(unsafe.Pointer(dr))), 0, 0, 0, 0)
			if ret != tc.want {
				t.Fatalf("local%v pos%v floor%d want%d", local, tc.pos, ret, tc.want)
			}
			rows = append(rows, o.capture(t, 11, ret))
		}
	}
	minimapCapture(t, "floor-cache", rows, "c725cbf3b5c1e911e9b7bdb1f5358d020ff92bd58012a9fba5b6e784afe1a719")
}
func TestMinimapFullRendering(t *testing.T) {
	o := newMinimapOwner(t)
	var rows []minimapResult
	for _, zoom := range []uint32{500, 1200, 1750, 2000, 2300, 4000} {
		for _, viewportLeft := range []int{0, 8, 16} {
			for scene := 0; scene < 8; scene++ {
				o.resetMinimap(t)
				*o.words["zoom"] = zoom
				o.c.Viewport().Screen.Min.X = viewportLeft
				local := o.player(t, 0, image.Pt(230, 230))
				remote := o.player(t, 1, image.Pt(245, 235))
				switch scene {
				case 1:
					local.TeamVal.ID = 1
					remote.TeamVal.ID = 1
				case 2:
					local.TeamVal.ID = 1
					remote.TeamVal.ID = 2
				case 3:
					o.players[0].Field3680 = 1
				case 4:
					remote.Buffs = 1 << 30
				case 5:
					o.players[1].Field3680 = 1
				case 6:
					o.polygon(t, 0, image.Rect(221, 220, 241, 240), 2)
				case 7:
					local.Buffs = 1 << 2
				}
				for i := 0; i < 11; i++ {
					wl := o.c.srv.Walls.CreateAtGrid(image.Pt(7+i%5, 8+i/5))
					if wl == nil {
						t.Fatal("actual wall allocation")
					}
					wl.Dir0 = byte(i)
					wl.Tile1 = byte(i % 3)
					wl.Field12 = uint32(i % 2)
					wl.Field8 = uint16(i % 3)
				}
				for i, name := range []string{"Crown", "GameBall"} {
					dr := o.c.Nox_xxx_spriteLoadAdd_45A360_drawable(o.c.Things.IndByID(name), image.Pt(220+i*20, 225))
					if dr == nil {
						t.Fatal("actual objective drawable allocation")
					}
					dr.NetCode32 = uint32(20 + i)
					o.c.Objs.MinimapAdd(dr, 1)
				}
				walls := o.c.srv.Walls.All()
				wallBytes := make([][]byte, len(walls))
				for i, w := range walls {
					wallBytes[i] = append([]byte(nil), unsafe.Slice((*byte)(unsafe.Pointer(w)), int(unsafe.Sizeof(*w)))...)
				}
				before, clip2 := o.c.r.Data().ClipRect(), o.c.r.Data().ClipRect2()
				ret := legacy.PortTestMinimap(13, int(uintptr(unsafe.Pointer(local))), 1, 0, 0, 0)
				if o.c.r.Data().ClipRect() != before || o.c.r.Data().ClipRect2() != clip2 {
					t.Fatal("minimap did not restore original clip state")
				}
				for i, w := range walls {
					if !bytes.Equal(wallBytes[i], unsafe.Slice((*byte)(unsafe.Pointer(w)), int(unsafe.Sizeof(*w)))) {
						t.Fatal("minimap draw mutated wall ownership or state")
					}
				}
				rows = append(rows, o.capture(t, 13, ret))
			}
		}
	}
	minimapCapture(t, "full-rendering", rows, "ea3c79d1e250c8cc01d63544be4422f2aa1809c9723c0538c8b0e6dcc8efac8e")
}

// The existing polygon predicate counts both incident edges when a ray passes
// through a vertex. Preserve the minimap's cold miss and subsequent cached level.
func TestMinimapCornerRayCompatibility(t *testing.T) {
	o := newMinimapOwner(t)
	o.polygon(t, 0, image.Rect(190, 190, 270, 270), 2)
	dr := o.player(t, 0, image.Pt(230, 230))
	var rows []minimapResult
	for _, want := range []uint32{1, 2, 2, 2, 2} {
		ret := legacy.PortTestMinimap(11, int(uintptr(unsafe.Pointer(dr))), 0, 0, 0, 0)
		if ret != want {
			t.Fatalf("corner ray floor%d want%d", ret, want)
		}
		rows = append(rows, o.capture(t, 11, ret))
	}
	minimapCapture(t, "corner-ray", rows, "78f6e192d052e1837d2eb3da7d04104173e11279cee02c9aeabfb2d081c53604")
}
