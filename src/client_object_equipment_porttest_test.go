//go:build porttest

package opennox

import (
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"testing"
	"unsafe"
)

func TestClientObjectDrawingEquipment(t *testing.T) {
	o := newObjectDrawingOwner(t)
	c := o.c
	var out []objectDrawingResult
	id := 0
	for _, op := range []int{7, 8, 9, 10, 16, 17} {
		for mask := 0; mask < 16; mask++ {
			for _, typ := range []int{4, 5} {
				for _, kind := range []uint32{2, 4, 5} {
					id++
					o.reset(uint32(id), 120)
					if op == 7 || op == 9 || op == 16 {
						o.data[0], o.data[1] = 8, o.frames[id%32]
					} else {
						o.data[3] = kind
					}
					dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(typ, image.Pt(48, 48))
					dr.DrawData = unsafe.Pointer(&o.data[0])
					dr.DrawFuncPtr = legacy.PortTestObjectDrawCallback(op)
					w := unsafe.Slice((*uint32)(dr.C()), 128)
					w[77] = uint32(id % 32)
					for i := 0; i < 4; i++ {
						if mask&(1<<i) != 0 {
							w[108+i] = uint32(uintptr(o.mods[i].C()))
						}
					}
					if op == 17 {
						noxflags.SetGame(noxflags.GameFlag(128))
						w[28] |= 0x10000000
						w[30] |= 0x1000000
						if mask%3 == 0 {
							w[109] = uint32(uintptr(o.mods[3].C()))
						}
					}
					for step := 0; step < 3; step++ {
						c.srv.SetFrame(uint32(120 + step*17))
						clear(o.pix.Pix)
						o.drawTrace = nil
						ret := dr.CallDraw(c.Viewport())
						out = append(out, o.result(t, id, step, ret))
					}
				}
			}
		}
	}
	effectsCapture(t, "object-drawing-equipment", out, len(out), "93efdb84ec27ae3124beb33f0006089b31a04e1e7cddc4ff01568c710a177879")
}
