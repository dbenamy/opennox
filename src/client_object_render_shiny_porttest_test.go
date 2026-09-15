//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"testing"
)

func TestClientObjectRenderShinyPeriods(t *testing.T) {
	o := newObjectRenderOwner(t)
	c := o.c
	type result struct {
		Count, Position, Origin, Z, Step int
		Start, Code, Age                 uint32
		Return                           uint16
		Render                           objectRenderResult
	}
	var out []result
	id := 0
	for _, count := range []int{1, 4, 16, 32} {
		for _, start := range []uint32{0, 120, 0xfffffff0} {
			for _, code := range []uint32{0, 7, 65535, 0xffffffff} {
				for pi, pos := range []image.Point{image.Pt(96, 96), image.Pt(0, 0)} {
					for origin := 0; origin < 2; origin++ {
						for _, z := range []int{0, 2} {
							id++
							o.resetRender(uint32(id), start)
							o.animation.ImagesSz = uint8(count)
							if origin != 0 {
								c.Viewport().World.Min = image.Pt(7, 9)
								c.Viewport().Screen.Min = image.Pt(3, 5)
							}
							dr := o.drawable(int(code), pos)
							dr.ZVal = uint16(z)
							dr.ZVal2 = uint16(z + 1)
							for step, age := range []uint32{0, 1, 2, uint32(count*2 - 1), uint32(count * 2), uint32(count*8 - 1), uint32(count * 8), uint32(count*8 + 1)} {
								c.srv.SetFrame(start + age)
								before := len(o.drawTrace)
								got := legacy.PortTestObjectRenderShiny(c.Viewport(), dr)
								frame := start + age + code
								if got != uint16(frame/uint32(8*count)) {
									t.Fatal("shiny period return does not preserve16bit quotient")
								}
								want := 0
								if (frame%uint32(8*count))>>1 < uint32(count) {
									want = 1
								}
								if len(o.drawTrace)-before != want {
									t.Fatal("shiny active/off-period draw count")
								}
								out = append(out, result{count, pi, origin, z, step, start, code, age, got, o.renderResult(t, id, step, int(got))})
							}
						}
					}
				}
			}
		}
	}
	effectsCapture(t, "object-render-shiny-periods", out, len(out), "941ab7b2b50e71a88fea225a990320f2d464b8be5044e58203f1adb36a4f3f64")
}
