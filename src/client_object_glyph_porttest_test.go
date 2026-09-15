//go:build porttest

package opennox

import (
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"testing"
	"unsafe"
)

func TestClientObjectDrawingGlyph(t *testing.T) {
	o := newObjectDrawingOwner(t)
	c := o.c
	var out []objectDrawingResult
	id := 0
	for mode := 0; mode < 5; mode++ {
		for _, delta := range []image.Point{{0, 0}, {1, 0}, {100, 100}, {149, 0}, {150, 0}, {151, 0}, {-149, 1}} {
			for _, kind := range []uint32{2, 4, 5} {
				id++
				o.reset(uint32(id), 120)
				o.data[3] = kind
				dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(48, 48))
				dr.DrawData = unsafe.Pointer(&o.data[0])
				dr.DrawFuncPtr = legacy.PortTestObjectDrawCallback(5)
				if mode > 0 {
					noxflags.SetGame(noxflags.GameFlag(2))
				}
				if mode > 1 {
					player := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(48, 48).Add(delta))
					*memmap.PtrUint32(0x852978, 8) = uint32(uintptr(player.C()))
					if mode == 3 {
						player.Buffs |= 1 << 21
					}
					if mode == 4 {
						w := unsafe.Slice((*uint32)(dr.C()), 128)
						w[30] |= 0x40000000
					}
				}
				for step := 0; step < 3; step++ {
					dr.AnimFrameSlave = uint32(step * 15)
					c.srv.SetFrame(uint32(120 + step))
					clear(o.pix.Pix)
					o.drawTrace = nil
					ret := dr.CallDraw(c.Viewport())
					out = append(out, o.result(t, id, step, ret))
					far := delta.X*delta.X+delta.Y*delta.Y >= 22500
					if mode == 2 && far && len(o.drawTrace) != 0 {
						t.Fatal("out-of-range glyph drew")
					}
					if !(mode == 2 && far) && len(o.drawTrace) != 1 {
						t.Fatal("visible glyph failed to draw")
					}
				}
			}
		}
	}
	effectsCapture(t, "object-drawing-glyph", out, len(out), "ea5fad0b65bbb9ab218b2173abc4ae7e368f910ece2adf78dafef4bb3fd0fb62")
}
