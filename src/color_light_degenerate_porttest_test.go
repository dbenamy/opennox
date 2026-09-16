//go:build porttest

package opennox

import (
	"encoding/binary"
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"testing"
	"unsafe"
)

// A coincident target or zero rotation arc has no direction change to apply.
// Preserve the existing light instead of converting NaN to an integer angle.
func TestColorLightDegenerateDirections(t *testing.T) {
	c, pix, effects := newEffectsFullOwner(t)
	env := legacy.PortTestNewClientParticleEnvironment()
	t.Cleanup(env.Restore)
	oldFPS := c.srv.TickRate()
	t.Cleanup(func() { c.srv.SetTickRate(oldFPS) })
	for _, op := range []int{3, 4} {
		for _, frame := range []uint32{0, 1, 30, 0xffffffff} {
			c.resetCase(effects, pix, 17, frame)
			env.Reset()
			c.srv.SetTickRate(30)
			d := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(200, 200))
			target := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, d.PosVec)
			target.ObjClass = 0x400000
			target.NetCode32 = 77
			raw := unsafe.Slice((*byte)(d.C()), int(unsafe.Sizeof(*d)))
			clear(raw[136:276])
			binary.LittleEndian.PutUint32(raw[164:], 0x12345678)
			binary.LittleEndian.PutUint32(raw[168:], 1)
			binary.LittleEndian.PutUint16(raw[176:], 0x1c0)
			binary.LittleEndian.PutUint32(raw[264:], 77)
			binary.LittleEndian.PutUint16(raw[268:], 90)
			binary.LittleEndian.PutUint16(raw[270:], 10)
			binary.LittleEndian.PutUint16(raw[272:], 90)
			legacy.PortTestColorLight(op, c.Viewport(), d)
			if binary.LittleEndian.Uint32(raw[164:]) != 0x12345678 || binary.LittleEndian.Uint32(raw[168:]) != 1 {
				t.Errorf("op%d frame%d: undefined direction changed light: angle=%x mode=%x", op, frame, binary.LittleEndian.Uint32(raw[164:]), binary.LittleEndian.Uint32(raw[168:]))
			}
		}
	}
}
