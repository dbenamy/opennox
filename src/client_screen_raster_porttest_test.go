//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
	"unsafe"
)

func TestClientScreenEffectsRasterPrimitives(t *testing.T) {
	o := newObjectDrawingOwner(t)
	c := o.c
	env := legacy.PortTestNewScreenEnvironment()
	defer env.Restore()
	type result struct {
		Op      int
		Input   [4]int32
		Alpha   byte
		Enabled bool
		Return  uint32
		Pixels  string
		Render  []uint32
	}
	var out []result
	capture := func(op int, a [4]int32, alpha byte, enabled bool) {
		o.reset(1, 120)
		env.Reset()
		*memmap.PtrUint32(0x5D4594, 1312492) = 0xabcd
		*memmap.PtrUint32(0x5D4594, 1312496) = 0x8421
		for i := range o.pix.Pix {
			o.pix.Pix[i] = 0xaaaa
		}
		c.r.Data().SetAlpha(alpha)
		c.r.Data().SetAlphaEnabled(enabled)
		ret := legacy.PortTestScreenPrimitive(op, a)
		out = append(out, result{op, a, alpha, enabled, ret, effectsPixelHash(o.pix), append([]uint32(nil), unsafe.Slice((*uint32)(unsafe.Pointer(c.r.Data())), 264)...)})
	}
	for _, alpha := range []byte{0, 128, 255} {
		for _, enabled := range []bool{false, true} {
			for _, a := range [][4]int32{{48, 48, 48, 48}, {0, 0, 95, 95}, {95, 0, 0, 95}, {48, -5, 48, 100}, {-5, 48, 100, 48}, {20, 20, 21, 21}, {20, 20, 21, 22}, {20, 20, 22, 21}, {-5, -5, -1, -1}} {
				capture(3, a, alpha, enabled)
			}
			for _, xy := range []int32{-5, 0, 48, 95, 100} {
				for _, radius := range []int32{0, 1, 2, 10, 31, 128} {
					capture(4, [4]int32{xy, xy, radius, 0xffff}, alpha, enabled)
				}
			}
		}
	}
	effectsCapture(t, "screen-effects-raster", out, len(out), "bb8f569d31fb9889f81a593e53441795df1747cc292e13c889c24cb53d83faf2")
}
