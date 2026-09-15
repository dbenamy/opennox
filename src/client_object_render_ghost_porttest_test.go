//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"testing"
)

func TestClientObjectRenderGhostAndFade(t *testing.T) {
	o := newObjectRenderOwner(t)
	c := o.c
	type result struct {
		See, Ghost, Fade, Semi, FPS int
		Age, Start                  uint32
		Origin                      int
		Opacity                     byte
		Render                      objectRenderResult
	}
	var out []result
	id := 0
	for see := 0; see < 2; see++ {
		for ghost := 0; ghost < 2; ghost++ {
			for fade := 0; fade < 2; fade++ {
				for semi := 0; semi < 2; semi++ {
					for _, fps := range []int{1, 30, 60, 120} {
						for _, age := range []uint32{0, 1, 30, 59, 60, 61, 120, 0xffffffff} {
							for _, origin := range []int{0, 1, 47, 100, -32768, -2147483648} {
								for _, start := range []uint32{120, 0xfffffff0} {
									id++
									o.resetRender(uint32(id), start)
									c.srv.SetTickRate(uint32(fps))
									local := o.drawable(7, image.Pt(48, 48))
									local.Buffs = uint32(see) << 21
									*memmap.PtrPtr(0x852978, 8) = local.C()
									dr := o.drawable(8, image.Pt(48, 48))
									dr.Field_8, dr.Field_9 = 48, 48
									// Fixed lighting isolates viewport/opacity arithmetic from world-grid bounds.
									dr.ObjFlags = 0x40000000
									if semi != 0 {
										dr.ObjFlags |= 0x4000000
									}
									dr.Field_120 = uint32(fade)
									dr.Field_85 = start - age
									if ghost != 0 {
										o.renderEnv.GhostType(dr.TypeIDVal)
									}
									c.Viewport().World.Min = image.Pt(origin, origin)
									opacity := legacy.PortTestObjectRenderGhost(c.Viewport(), dr)
									if see != 0 && opacity != 255 {
										t.Fatal("see-invisible must make ghost fully visible")
									}
									if see == 0 && origin == 0 && opacity != 128 {
										t.Fatal("centered ghost opacity must be128")
									}
									legacy.Nox_xxx_drawObject_4C4770_draw(c.Viewport(), dr, noxrender.ImageHandle(o.images[id%32].C()))
									out = append(out, result{see, ghost, fade, semi, fps, age, start, origin, opacity, o.renderResult(t, id, 0, 0)})
								}
							}
						}
					}
				}
			}
		}
	}
	effectsCapture(t, "object-render-ghost-fade", out, len(out), "b5c83febbb7e0aadb9221326e67892bb22f1e40794b8097cb5a3c26da1ea4049")
}
