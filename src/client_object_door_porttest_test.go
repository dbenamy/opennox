//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/client/noxrender"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"testing"
	"unsafe"
)

func TestClientObjectDrawingDoor(t *testing.T) {
	o := newObjectDrawingOwner(t)
	c := o.c
	var out []objectDrawingResult
	id := 0
	for _, dir := range []byte{0, 1, 8, 16, 24, 31} {
		for _, lock := range []byte{0, 1, 2, 255} {
			for mode := 0; mode < 3; mode++ {
				for _, pos := range []image.Point{{48, 48}, {80, 90}, {0, 0}} {
					id++
					o.reset(uint32(id), 120)
					o.data[0], o.data[2] = 12, 32
					if mode > 0 {
						noxflags.SetGame(noxflags.GameFlag(4096))
					}
					dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, pos)
					dr.DrawData = unsafe.Pointer(&o.data[0])
					dr.DrawFuncPtr = legacy.PortTestObjectDrawCallback(0)
					b := unsafe.Slice((*byte)(dr.C()), 512)
					b[299] = dir
					b[432] = byte(mode & 1)
					b[433] = lock
					for step := 0; step < 2; step++ {
						if step == 1 {
							for x := range c.tiles.nox_arr2_853BC0 {
								for y := range c.tiles.nox_arr2_853BC0[x] {
									c.tiles.nox_arr2_853BC0[x][y] = noxrender.RGB{R: (50 + x%100) << 16, G: (100 + y%100) << 16, B: 175 << 16}
								}
							}
						}
						clear(o.pix.Pix)
						o.drawTrace = nil
						ret := dr.CallDraw(c.Viewport())
						out = append(out, o.result(t, id, step, ret))
						if mode == 1 && len(o.namedCalls) != 4 {
							t.Fatal("door should resolve four lock images once")
						}
					}
				}
			}
		}
	}
	effectsCapture(t, "object-drawing-door", out, len(out), "cad98de7e6c05e4f3cdd9a42c5cb763bee1c30683a89dfb2dd60d48cb346c21f")
}
