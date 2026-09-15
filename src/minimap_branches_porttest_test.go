//go:build porttest

package opennox

import (
	"github.com/opennox/libs/types"
	"github.com/opennox/libs/wall"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"image"
	"testing"
	"unsafe"
)

func TestMinimapVisibilityWrappers(t *testing.T) {
	o := newMinimapOwner(t)
	var rows []minimapResult
	for _, guiState := range []uint32{0, 1, 2} {
		for _, visible := range []uint32{0, 1} {
			for _, poison := range []bool{false, true} {
				for _, message := range []bool{false, true} {
					o.resetMinimap(t)
					dr := o.player(t, 0, image.Pt(230, 230))
					if poison {
						dr.Buffs = 1 << 2
					}
					*o.words["gui"] = guiState
					*memmap.PtrUint32(0x5D4594, 1096424) = visible
					if message {
						alloc.StrCopy16(unsafe.Slice(memmap.PtrUint16(0x5D4594, 823804), 320), "Map message")
						*memmap.PtrUint32(0x5D4594, 824440) = 200
					}
					before := effectsPixelHash(o.pix)
					ret := legacy.PortTestMinimap(14, 0, 0, 0, 0, 0)
					want := uint32(0)
					if guiState != 1 && visible != 0 && !poison {
						want = 1
					}
					if got := memmap.Uint32(0x5D4594, 1096316); got != want {
						t.Fatalf("wrapper floor%d want%d", got, want)
					}
					if guiState == 1 && (ret != 1 || effectsPixelHash(o.pix) != before) {
						t.Fatal("GUI suppression changed drawing or return")
					}
					rows = append(rows, o.capture(t, 14, ret))
				}
			}
		}
	}
	for _, poison := range []bool{false, true} {
		o.resetMinimap(t)
		dr := o.player(t, 0, image.Pt(230, 230))
		if poison {
			dr.Buffs = 1 << 2
		}
		clip := image.Rect(10, 11, 85, 86)
		o.c.r.Data().SetClipRect(clip)
		o.c.r.Data().SetClipRect2(image.Rect(10, 11, 84, 85))
		before := effectsPixelHash(o.pix)
		ret := legacy.PortTestMinimap(12, int(uintptr(unsafe.Pointer(dr))), 0, 0, 0, 0)
		if poison {
			if effectsPixelHash(o.pix) != before || o.c.r.Data().ClipRect() != clip {
				t.Fatal("poisoned minimap wrapper changed pixels or clipping")
			}
		} else if o.c.r.Data().ClipRect() != o.pix.Rect {
			t.Fatal("wrapper did not restore full-screen clipping")
		}
		rows = append(rows, o.capture(t, 12, ret))
	}
	minimapCapture(t, "visibility-wrappers", rows, "7ddc6839a902f95722c7f11c77dccb8905664423d1bd891d18ebaa2b4595fc3f")
}
func TestMinimapDebugOverlay(t *testing.T) {
	o := newMinimapOwner(t)
	configure, restore := o.c.srv.PortTestMinimapDebug()
	t.Cleanup(restore)
	var rows []minimapResult
	paths := [][]types.Pointf{nil, {{X: 215.75, Y: 228.25}}, {{X: 215.75, Y: 228.25}, {X: 237.5, Y: 239.75}}, {{X: -4.5, Y: 210.25}, {X: 215.75, Y: 228.25}, {X: 237.5, Y: 239.75}, {X: 254.125, Y: 218.5}}}
	for _, enabled := range []bool{false, true} {
		for _, path := range paths {
			for _, count := range []int{-1, 0, 1, 3} {
				for _, zoom := range []uint32{500, 2300} {
					o.resetMinimap(t)
					dr := o.player(t, 0, image.Pt(230, 230))
					*o.words["zoom"] = zoom
					noxflags.UnsetEngine(noxflags.EngineShowAI)
					if enabled {
						noxflags.SetEngine(noxflags.EngineShowAI)
					}
					configure(path, count)
					ret := legacy.PortTestMinimap(13, int(uintptr(unsafe.Pointer(dr))), 1, 0, 0, 0)
					if *o.words["debugIterator"] != 0 {
						t.Fatal("debug monster traversal retained a live iterator")
					}
					if got := o.c.srv.AI.Paths.Points(); len(got) != len(path) {
						t.Fatal("drawing changed debug path length")
					} else {
						for i, p := range path {
							if got[i] != p {
								t.Fatal("drawing changed debug path point")
							}
						}
					}
					rows = append(rows, o.capture(t, 13, ret))
				}
			}
		}
	}
	noxflags.UnsetEngine(noxflags.EngineShowAI)
	minimapCapture(t, "debug-overlay", rows, "d19b1faffd6908755042403f75924a39f57871485d0e0d819f5eb2e72019f50b")
}
func TestMinimapDoorWalls(t *testing.T) {
	o := newMinimapOwner(t)
	var rows []minimapResult
	for _, dir := range []byte{0, 8, 16, 24} {
		for _, neighbor := range []int{-1, 0, 1, 2, 3, 4, 5, 6, 7, 8} {
			for _, revealed := range []bool{false, true} {
				for _, zoom := range []uint32{500, 1200, 2000, 2300} {
					o.resetMinimap(t)
					local := o.player(t, 0, image.Pt(230, 230))
					*o.words["zoom"] = zoom
					door := o.drawable(40, image.Pt(230, 230))
					door.Field_74_4 = dir
					wl := o.c.srv.Walls.CreateAtGrid(image.Pt(10, 10))
					wl.Flags4 = wall.Flags(0x10)
					wl.Field32 = uint32(uintptr(unsafe.Pointer(door)))
					if neighbor >= 0 {
						x, y := 2, 0 // Unrelated neighbor is a negative control.
						if neighbor >= 5 {
							p := []image.Point{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}[neighbor-5]
							x, y = p.X, p.Y
						}
						if neighbor < 4 {
							off := uintptr(149240 + 8*neighbor)
							x, y = int(memmap.Int32(0x587000, off)), int(memmap.Int32(0x587000, off+4))
						}
						adj := o.c.srv.Walls.CreateAtGrid(image.Pt(10+x, 10+y))
						if adj == nil {
							t.Fatal("door neighbor allocation")
						}
						if revealed {
							adj.Field12 = 1
						}
					}
					ret := legacy.PortTestMinimap(13, int(uintptr(unsafe.Pointer(local))), 1, 0, 0, 0)
					rows = append(rows, o.capture(t, 13, ret))
				}
			}
		}
	}
	minimapCapture(t, "door-walls", rows, "2cce62e05865184c7694799a5dea64349632846a464d4ba51cd77566c157157f")
}
