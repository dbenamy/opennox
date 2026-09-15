//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"runtime"
	"testing"
)

func TestClientObjectRenderFadePrecision(t *testing.T) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if legacy.PortTestProtectionFloatCW()&0x0f00 != 0x0200 {
		t.Fatal("fade fixture requires hosted x87 precision 53 and nearest rounding")
	}
	o := newObjectRenderOwner(t)
	c := o.c
	type result struct {
		FPS, Semi  int
		Age, Start uint32
		Render     objectRenderResult
	}
	var out []result
	id := 0
	for _, fps := range []uint32{1, 2, 3, 7, 10, 30, 59, 60, 61, 120, 127, 255} {
		var ages []uint32
		for age := uint32(0); age <= fps+1; age++ {
			ages = append(ages, age)
		}
		ages = append(ages, 2*fps, 0x80000000, 0xffffffff)
		for _, age := range ages {
			for semi := 0; semi < 2; semi++ {
				for _, start := range []uint32{120, 0xfffffff0} {
					id++
					o.resetRender(uint32(id), start)
					c.srv.SetTickRate(fps)
					dr := o.drawable(8, image.Pt(48, 48))
					dr.ObjFlags = 0x40000000
					if semi != 0 {
						dr.ObjFlags |= 0x4000000
					}
					dr.Field_120 = 1
					dr.Field_85 = start - age
					legacy.Nox_xxx_drawObject_4C4770_draw(c.Viewport(), dr, noxrender.ImageHandle(o.images[id%32].C()))
					out = append(out, result{int(fps), semi, age, start, o.renderResult(t, id, 0, 0)})
				}
			}
		}
	}
	effectsCapture(t, "object-render-fade-precision", out, len(out), "b7172234edc7ad017e466d0f6a79510d27456404d1533a9cd293ea278941b882")
}
