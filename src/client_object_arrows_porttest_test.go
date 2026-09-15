//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"testing"
	"unsafe"
)

func TestClientObjectDrawingArrows(t *testing.T) {
	o := newObjectDrawingOwner(t)
	c := o.c
	var out []objectDrawingResult
	id := 0
	for _, op := range []int{1, 2} {
		for _, fps := range []int{30, 60, 120} {
			for _, start := range []uint32{0, 120, 0xfffffffc} {
				for _, delta := range []image.Point{{0, 0}, {14, 0}, {10, 10}, {11, 9}, {14, 3}, {-20, 0}, {0, -20}} {
					for _, fail := range []int{0, 2, 3} {
						id++
						o.reset(uint32(id), start)
						c.srv.SetTickRate(uint32(fps))
						o.data[0], o.data[2] = 12, 32
						dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(48, 48))
						dr.DrawData = unsafe.Pointer(&o.data[0])
						dr.DrawFuncPtr = legacy.PortTestObjectDrawCallback(op)
						dr.Field_81, dr.Field_82 = uint32(48-delta.X), uint32(48-delta.Y)
						c.FailEvery = fail
						for step := 0; step < 3; step++ {
							c.srv.SetFrame(start + uint32(step))
							dr.AnimFrameSlave = uint32((id + step) % 32)
							clear(o.pix.Pix)
							o.drawTrace = nil
							ret := dr.CallDraw(c.Viewport())
							out = append(out, o.result(t, id, step, ret))
						}
					}
				}
			}
		}
	}
	for _, op := range []int{3, 4} {
		for _, fps := range []int{30, 60, 120} {
			for _, start := range []uint32{0, 120, 0xfffffffc} {
				for _, remaining := range []uint32{0, 1, 2, 9, 10, 19, 20, 39, 40} {
					id++
					o.reset(uint32(id), start)
					c.srv.SetTickRate(uint32(fps))
					dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(18, 40))
					dr.DrawFuncPtr = legacy.PortTestObjectDrawCallback(op)
					w := unsafe.Slice((*uint32)(dr.C()), 128)
					w[108], w[109] = 78, 60
					dr.Deadline = start + remaining
					for step := 0; step < 2; step++ {
						halves := unsafe.Slice((*int16)(dr.C()), 256)
						halves[52], halves[53] = int16(step*3), int16(step*2)
						for i := range o.pix.Pix {
							o.pix.Pix[i] = 0xaaaa
						}
						ret := dr.CallDraw(c.Viewport())
						out = append(out, o.result(t, id, step, ret))
					}
				}
			}
		}
	}
	effectsCapture(t, "object-drawing-arrows", out, len(out), "c7cedab0fa2d041920816e900230bb2e5df0fb0d04af3a18313fa6400b4576c1")
}
