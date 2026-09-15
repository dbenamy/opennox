//go:build porttest && !server

package opennox

import (
	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"testing"
	"unsafe"
)

// The server target intentionally has no client scanline occlusion implementation.
func TestClientObjectRenderOcclusion(t *testing.T) {
	o := newObjectRenderOwner(t)
	c := o.c
	type result struct {
		Mode, Direction, Position, Origin, Step int
		Flags                                   uint32
		Render                                  objectRenderResult
	}
	var out []result
	id := 0
	spans := [][]int{nil, {0, 95}, {0, 40}, {55, 95}, {0, 30, 60, 95}, {40, 50}}
	for mode, span := range spans {
		for dir := 0; dir < 32; dir++ {
			for pi, pos := range []image.Point{image.Pt(1, 1), image.Pt(48, 48), image.Pt(94, 94)} {
				for _, flags := range []uint32{0, 1, 0x40000000} {
					for origin := 0; origin < 2; origin++ {
						id++
						o.resetRender(uint32(id), 120)
						*memmap.PtrUint32(0x587000, 80808) = 1
						for y := 0; y < 96; y++ {
							c.tiles.nox_arr_956A00[y] = len(span)
							copy(c.tiles.nox_arr_957820[y].arr[:], span)
						}
						if origin != 0 {
							c.Viewport().World.Min = image.Pt(7, 9)
							c.Viewport().Screen.Min = image.Pt(3, 5)
						}
						dr := o.drawable(8, pos)
						dr.ObjClass = 0x80
						dr.ObjFlags = object.Flags(flags)
						*(*byte)(unsafe.Add(dr.C(), 299)) = byte(dir)
						for step := 0; step < 2; step++ {
							legacy.Nox_xxx_drawObject_4C4770_draw(c.Viewport(), dr, noxrender.ImageHandle(o.images[(id+step)%32].C()))
							// Empty horizontal spans are rejected. Short nonhorizontal paths can
							// be accepted by the existing helper; capture that behavior unchanged.
							if mode == 0 && flags == 0 && (dir == 12 || dir == 28) && len(o.drawTrace) != 0 {
								t.Fatalf("empty spans reached renderer: mode%d dir%d pos%d origin%d step%d", mode, dir, pi, origin, step)
							}
							out = append(out, result{mode, dir, pi, origin, step, flags, o.renderResult(t, id, step, 0)})
						}
					}
				}
			}
		}
	}
	effectsCapture(t, "object-render-occlusion", out, len(out), "b10776df02ca123de504afd8c6022350c482a3fb73a0181aec53b273f8b5816c")
}
