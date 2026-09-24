//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"image"
	"testing"
	"unsafe"

	"github.com/opennox/libs/prand"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

// Exercise the creation bridge independently of the existing rain-orb draw tests.
// Arguments deliberately cross the C ushort/char boundaries; failed allocation
// must consume no random numbers and return a null pointer word.
func TestClientEffectsRainCreation(t *testing.T) {
	c, pix, env := newEffectsFullOwner(t)
	data, free := alloc.Make([]byte{}, 32)
	defer free()
	for _, fail := range []bool{false, true} {
		for _, z := range []int32{-1, 0, 32768, 65535, 65536, 65537} {
			for _, velocity := range []int32{-129, -128, -1, 0, 127, 128, 255, 256} {
				c.resetCase(env, pix, 31, 123)
				if fail {
					c.FailEvery = 1
				}
				for i := range data {
					data[i] = 0xa5
				}
				points := (*[4]int32)(unsafe.Pointer(&data[8]))
				*points = [4]int32{-17, 83, 257, -513}
				before := bytes.Clone(data)
				got := legacy.PortTestClientEffects(4, c.Viewport(), nil, [8]int32{4, z, velocity}, unsafe.Pointer(points))
				if !bytes.Equal(data, before) {
					t.Fatal("rain creation changed points or guards")
				}
				if len(c.Calls) != 1 || c.Calls[0].Type != 4 || c.Calls[0].Position != image.Pt(-17, 83) {
					t.Fatalf("spawn calls: %+v", c.Calls)
				}
				if c.srv.Rand.Logic.Index() != 31 {
					t.Fatal("logic RNG consumed")
				}
				if fail {
					if got != 0 || c.Objs.Count != 0 || c.srv.Rand.Other.Index() != 32 {
						t.Fatal("failed creation changed state")
					}
					continue
				}
				dr := c.Objs.List1
				if dr == nil || got != uint32(uintptr(unsafe.Pointer(dr))) || c.Objs.Count != 1 || c.Objs.List4 != dr || dr.Flags()&0x400000 == 0 {
					t.Fatal("created drawable pointer/list mismatch")
				}
				raw := unsafe.Slice((*byte)(unsafe.Pointer(dr)), 512)
				if dr.ZVal != uint16(z) || dr.ZVal2 != 0 || dr.VelZ != int8(velocity) || binary.LittleEndian.Uint16(raw[440:]) != uint16(z) {
					t.Fatalf("height/velocity narrowing: %d %d", z, velocity)
				}
				if int32(binary.LittleEndian.Uint32(raw[432:])) != 257 || int32(binary.LittleEndian.Uint32(raw[436:])) != -513 {
					t.Fatal("destination changed")
				}
				expected := prand.New(32)
				want := expected.Int(3, 10)
				if raw[442] != byte(want) || c.srv.Rand.Other.Index() != expected.Index() {
					t.Fatal("creation RNG value/consumption")
				}
			}
		}
	}
}
