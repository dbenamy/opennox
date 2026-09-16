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

func TestColorLightRotation(t *testing.T) {
	c, pix, effects := newEffectsFullOwner(t)
	env := legacy.PortTestNewClientParticleEnvironment()
	t.Cleanup(env.Restore)
	oldFPS := c.srv.TickRate()
	t.Cleanup(func() { c.srv.SetTickRate(oldFPS) })
	type record struct {
		FPS               uint32
		Speed             int16
		Flags, Start, End uint16
		Frame             uint32
		Light             []uint32
	}
	var rows []record
	for _, fps := range []uint32{1, 30, 60} {
		for _, speed := range []int16{-32768, -360, -1, 0, 1, 90, 360, 32767} {
			for _, flags := range []uint16{0, 0x80, 0x180} {
				for _, ends := range [][2]uint16{{0, 90}, {90, 270}, {359, 0}, {45, 360}} {
					for _, frame := range []uint32{0, 1, 29, 30, 31, 59, 60, 61, 1000, 0x7fffffff, 0xffffffff} {
						// The legacy cycle narrows through int32. Exclude extreme
						// combinations whose subsequent floating-to-int cast is
						// outside C's defined range; they are not an oracle.
						if flags&0x80 != 0 && speed != 0 {
							span := 360.0
							if flags&0x100 != 0 {
								span = float64(int(ends[1]) - int(ends[0]))
							}
							cycle := float64(frame) / (span / float64(speed) * float64(fps))
							whole := int64(cycle)
							if whole < -2147483648 || whole > 2147483647 {
								continue
							}
						}
						t.Run(fmt.Sprintf("fps%d-speed%d-flags%x-start%d-end%d-frame%d", fps, speed, flags, ends[0], ends[1], frame), func(t *testing.T) {
							c.resetCase(effects, pix, 17, frame)
							env.Reset()
							c.srv.SetTickRate(fps)
							d := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(48, 48))
							raw := unsafe.Slice((*byte)(d.C()), int(unsafe.Sizeof(*d)))
							clear(raw[136:276])
							binary.LittleEndian.PutUint32(raw[164:], 0x12345678)
							binary.LittleEndian.PutUint32(raw[168:], 1)
							binary.LittleEndian.PutUint16(raw[176:], flags)
							binary.LittleEndian.PutUint16(raw[268:], ends[0])
							binary.LittleEndian.PutUint16(raw[270:], uint16(speed))
							binary.LittleEndian.PutUint16(raw[272:], ends[1])
							legacy.PortTestColorLight(4, c.Viewport(), d)
							words := make([]uint32, 10)
							for i := range words {
								words[i] = binary.LittleEndian.Uint32(raw[136+4*i:])
							}
							if flags&0x80 == 0 || speed == 0 {
								if words[7] != 0x12345678 || words[8] != 1 {
									t.Fatal("disabled rotation changed direction")
								}
							} else {
								if words[7]>>16 != 0x1234 || words[8] != 0 {
									t.Fatal("rotation setter ownership")
								}
								if frame == 0 {
									want := uint16((uint32(ends[0])*65536 + 180) / 360)
									if uint16(words[7]) != want {
										t.Fatalf("initial angle=%d want=%d", uint16(words[7]), want)
									}
								}
							}
							rows = append(rows, record{fps, speed, flags, ends[0], ends[1], frame, words})
						})
					}
				}
			}
		}
	}
	spellbookCapture(t, "color-light-rotation", rows, "1d3d2d5924bd55310a9786918decfa6ba2b53cc080a782f699ac22777b0e193e")
}
