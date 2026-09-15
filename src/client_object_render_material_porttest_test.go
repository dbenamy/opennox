//go:build porttest

package opennox

import (
	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/client/noxrender"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"testing"
)

func TestClientObjectRenderMaterialsAndCrop(t *testing.T) {
	o := newObjectRenderOwner(t)
	c := o.c
	type result struct {
		Class, Flags, Buffs uint32
		Z, Mode, Step       int
		Render              objectRenderResult
	}
	var out []result
	id := 0
	for _, class := range []uint32{0, 2, 0x80000, 0x80002} {
		for _, flags := range []uint32{0, 0x40000000, 0x40008020, 0x1000000, 0x4000000, 0x5000000} {
			for _, buffs := range []uint32{0, 1 << 11, 1 << 23, 1 << 25, 1<<11 | 1<<23, 1<<11 | 1<<25, 1<<23 | 1<<25} {
				for _, z := range []int{-2, -1, 0, 1, 10} {
					for mode := 0; mode < 2; mode++ {
						id++
						o.resetRender(uint32(id), 120)
						if mode == 1 {
							noxflags.SetGame(noxflags.GameFlag(2048))
							for x := range c.tiles.nox_arr2_853BC0 {
								for y := range c.tiles.nox_arr2_853BC0[x] {
									c.tiles.nox_arr2_853BC0[x][y] = noxrender.RGB{R: 127 << 16, G: 90 << 16, B: 210 << 16}
								}
							}
						}
						dr := o.drawable(8, image.Pt(48, 48))
						dr.ObjClass = object.Class(class)
						dr.ObjFlags = object.Flags(flags)
						dr.Buffs = buffs
						dr.ZVal = uint16(z)
						dr.ZVal2 = uint16(mode * 3)
						dr.Field_0 = uint32(id%8) | uint32(id%5)<<8
						for step := 0; step < 3; step++ {
							c.srv.SetFrame(uint32(120 + step))
							legacy.Nox_xxx_drawObject_4C4770_draw(c.Viewport(), dr, noxrender.ImageHandle(o.images[(id+step)%32].C()))
							out = append(out, result{class, flags, buffs, z, mode, step, o.renderResult(t, id, step, 0)})
						}
					}
				}
			}
		}
	}
	effectsCapture(t, "object-render-material-crop", out, len(out), "764d3fae9e5c55e8deb071e8af27ab98ea8295bf294a03dfb638c5bcc8e50737")
}
