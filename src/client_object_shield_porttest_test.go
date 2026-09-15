//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"testing"
	"unsafe"
)

func TestClientObjectDrawingShield(t *testing.T) {
	o := newObjectDrawingOwner(t)
	c := o.c
	var out []objectDrawingResult
	id := 0
	for _, static := range []bool{false, true} {
		for _, present := range []bool{false, true} {
			for _, start := range []uint32{0, 120, 0xfffffffc} {
				for _, kind := range []uint32{2, 4, 5} {
					id++
					o.reset(uint32(id), start)
					o.data[3] = kind
					owner := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(48, 48))
					ow := unsafe.Slice((*uint32)(owner.C()), 128)
					ow[32] = 7
					if static {
						ow[28] |= 0x20000000
					}
					dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(3, 3))
					dr.DrawData = unsafe.Pointer(&o.data[0])
					dr.DrawFuncPtr = legacy.PortTestObjectDrawCallback(11)
					w := unsafe.Slice((*uint32)(dr.C()), 128)
					w[32] = 0 - start
					w[77] = 31
					w[108] = 7
					if !present {
						w[108] = 99
					}
					if static {
						w[108] |= 0x8000
					}
					for step := 0; step < 4; step++ {
						c.srv.SetFrame(start + uint32(step))
						clear(o.pix.Pix)
						o.drawTrace = nil
						pos := image.Pt(35+step*7, 35+step*5)
						c.Nox_xxx_updateSpritePosition_49AA90(owner, pos.X, pos.Y)
						if step == 3 {
							c.Nox_xxx_spriteDeleteStatic_45A4E0_drawable(owner)
						}
						ret := dr.CallDraw(c.Viewport())
						if present && step < 3 {
							if ret != 1 || dr.PosVec != pos.Add(image.Pt(0, 3)) {
								t.Fatal("shield did not follow network owner")
							}
						} else if ret != 0 {
							t.Fatal("orphan shield not removed")
						}
						out = append(out, o.result(t, id, step, ret))
						if ret == 0 {
							break
						}
					}
				}
			}
		}
	}
	effectsCapture(t, "object-drawing-shield", out, len(out), "fbd4a997d5f8f1afd70ba2f4771d496016c0ef2ffe70825e08f91743b7d6c71e")
}
