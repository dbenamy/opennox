//go:build porttest

package opennox

import (
	"encoding/binary"
	"fmt"
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"testing"
	"unsafe"
)

func TestColorLightTargetDirection(t *testing.T) {
	c, pix, effects := newEffectsFullOwner(t)
	env := legacy.PortTestNewClientParticleEnvironment()
	t.Cleanup(env.Restore)
	type record struct {
		XY               image.Point
		Enabled, Present bool
		Light            []uint32
	}
	var rows []record
	for _, xy := range []image.Point{{100, 0}, {0, 100}, {-100, 0}, {0, -100}, {100, 100}, {-100, 100}, {-100, -100}, {100, -100}, {5000, 1}, {1, 5000}} {
		for _, enabled := range []bool{false, true} {
			for _, present := range []bool{false, true} {
				t.Run(fmt.Sprintf("%d-%d-enabled%v-present%v", xy.X, xy.Y, enabled, present), func(t *testing.T) {
					c.resetCase(effects, pix, 17, 123)
					env.Reset()
					d := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(200, 200))
					target := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, d.PosVec.Add(xy))
					target.ObjClass = 0x400000
					target.NetCode32 = 77
					raw := unsafe.Slice((*byte)(d.C()), int(unsafe.Sizeof(*d)))
					clear(raw[136:276])
					binary.LittleEndian.PutUint32(raw[164:], 0x12345678)
					binary.LittleEndian.PutUint32(raw[168:], 9)
					if enabled {
						raw[176] = 0x40
					}
					code := uint32(88)
					if present {
						code = 77
					}
					binary.LittleEndian.PutUint32(raw[264:], code)
					legacy.PortTestColorLight(3, c.Viewport(), d)
					words := make([]uint32, 10)
					for i := range words {
						words[i] = binary.LittleEndian.Uint32(raw[136+4*i:])
					}
					if !enabled || !present {
						if words[7] != 0x12345678 || words[8] != 9 {
							t.Fatal("missing/disabled target changed direction")
						}
					} else {
						if words[8] != 0 || words[7]>>16 != 0x1234 {
							t.Fatal("direction setter ownership")
						}
						if xy == image.Pt(100, 0) && uint16(words[7]) != 0 {
							t.Fatal("east direction")
						}
					}
					rows = append(rows, record{xy, enabled, present, words})
				})
			}
		}
	}
	spellbookCapture(t, "color-light-direction", rows, "9c4da0c65f2933d2ed77f8bdd73104a06110434c2d19526d11cbca198afbc33f")
}
