//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"math"
	"testing"
	"unsafe"
)

func TestClientObjectDrawingGeometry(t *testing.T) {
	o := newObjectDrawingOwner(t)
	c := o.c
	var out []objectDrawingResult
	id := 0
	for _, op := range []int{13, 14, 15} {
		for _, pos := range []image.Point{{48, 48}, {0, 0}, {95, 95}, {1, 94}, {96, 48}} {
			for _, frac := range []float32{-0.51, -0.5, -0.49, 0, 0.49, 0.5, 0.51} {
				for color := 0; color < 4; color++ {
					id++
					o.reset(uint32(id), 120)
					dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, pos)
					dr.DrawFuncPtr = legacy.PortTestObjectDrawCallback(op)
					w := unsafe.Slice((*uint32)(dr.C()), 128)
					for i, x := range []float32{-10, -5, 10, -5, -10, 5, 10, 5} {
						w[16+i] = math.Float32bits(x + frac)
					}
					b := unsafe.Slice((*byte)(dr.C()), 512)
					switch color {
					case 1:
						copy(b[432:438], []byte{255, 0, 0, 0, 255, 0})
					case 2:
						copy(b[432:438], []byte{0, 0, 255, 255, 255, 255})
					case 3:
						copy(b[432:438], []byte{1, 2, 3, 4, 5, 6})
					}
					for step, alpha := range []byte{0, 128, 255} {
						// Nonzero background makes opaque black and alpha edges observable.
						for i := range o.pix.Pix {
							o.pix.Pix[i] = 0xaaaa
						}
						c.r.Data().SetAlpha(alpha)
						ret := dr.CallDraw(c.Viewport())
						out = append(out, o.result(t, id, step, ret))
					}
				}
			}
		}
	}
	effectsCapture(t, "object-drawing-geometry", out, len(out), "9a2075112d571e8505b385373147974d0742aac82785085de10708a64e9e0370")
}
