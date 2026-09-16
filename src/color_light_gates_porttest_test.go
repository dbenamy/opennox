//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"testing"
	"unsafe"
)

func TestColorLightPropertyGates(t *testing.T) {
	c, pix, effects := newEffectsFullOwner(t)
	env := legacy.PortTestNewClientParticleEnvironment()
	t.Cleanup(env.Restore)
	for _, op := range []int{2, 4} {
		for _, frame := range []uint32{0, 1, 30, 0xffffffff} {
			c.resetCase(effects, pix, 17, frame)
			env.Reset()
			d := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(48, 48))
			raw := unsafe.Slice((*byte)(d.C()), int(unsafe.Sizeof(*d)))
			clear(raw[136:276])
			clear(raw[432:435])
			binary.LittleEndian.PutUint32(raw[164:], 0x12345678)
			if op == 2 {
				binary.LittleEndian.PutUint32(raw[168:], 1)
			}
			raw[176] = 0x90
			raw[434] = 2
			raw[242] = 20
			raw[243] = 40
			binary.LittleEndian.PutUint16(raw[262:], 3)
			binary.LittleEndian.PutUint16(raw[268:], 90)
			binary.LittleEndian.PutUint16(raw[270:], 10)
			binary.LittleEndian.PutUint16(raw[272:], 180)
			before := append([]byte(nil), raw...)
			legacy.PortTestColorLight(op, c.Viewport(), d)
			if !bytes.Equal(before, raw) {
				t.Fatalf("op%d frame%d changed a disabled light property", op, frame)
			}
		}
	}
}
