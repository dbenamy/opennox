//go:build porttest

package opennox

import (
	"encoding/binary"
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"testing"
	"unsafe"
)

func TestClientScreenEffectsRopeStaticEndpoints(t *testing.T) {
	o := newObjectDrawingOwner(t)
	env := legacy.PortTestNewScreenEnvironment()
	t.Cleanup(env.Restore)
	c := o.c
	c.callbackRefs[legacy.PortTestScreenEffectCallback(1)] = 0xec000001
	for mask := 0; mask < 4; mask++ {
		o.reset(1, 120)
		from := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(20, 35))
		from.NetCode32 = 7
		to := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(70, 65))
		to.NetCode32 = 8
		codes := [2]uint32{7, 8}
		if mask&1 != 0 {
			w := unsafe.Slice((*uint32)(from.C()), 128)
			w[28] |= 0x20000000
			codes[0] |= 0x8000
		}
		if mask&2 != 0 {
			w := unsafe.Slice((*uint32)(to.C()), 128)
			w[28] |= 0x20000000
			codes[1] |= 0x8000
		}
		dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(45, 50))
		dr.DrawFuncPtr = legacy.PortTestScreenEffectCallback(1)
		b := unsafe.Slice((*byte)(dr.C()), 512)
		b[432] = 1
		// The real ray-effect constructor stores the encoded uint16 endpoint IDs
		// in these unaligned32-bit slots; drawable NetCode32 stores cleared IDs.
		binary.LittleEndian.PutUint32(b[437:], codes[0])
		binary.LittleEndian.PutUint32(b[441:], codes[1])
		blank := effectsPixelHash(o.pix)
		if dr.CallDraw(c.Viewport()) != 1 || effectsPixelHash(o.pix) == blank {
			t.Fatalf("rope omitted valid endpoints static mask%d", mask)
		}
		c.snapshotDrawables(t)
	}
}
