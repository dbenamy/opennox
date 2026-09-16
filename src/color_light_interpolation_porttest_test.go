//go:build porttest

package opennox

import (
	"encoding/binary"
	"fmt"
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"math"
	"testing"
	"unsafe"
)

func TestColorLightInterpolation(t *testing.T) {
	c, pix, effects := newEffectsFullOwner(t)
	env := legacy.PortTestNewClientParticleEnvironment()
	t.Cleanup(env.Restore)
	type record struct {
		Op, Count, Period, Mode int
		Frame                   uint32
		Return                  int
		Light                   []uint32
	}
	var records []record
	for op := 0; op < 3; op++ {
		for _, count := range []int{0, 1, 2, 3, 16, 128, 255} {
			for _, period := range []int{0, 1, 2, 3, 255, 256, 257, 65535} {
				for _, mode := range []int{0, 1} {
					frames := []uint32{0, 1, 2, 3, 254, 255, 256, 257, uint32(period*count - 1), uint32(period * count), uint32(2*period*count - 1), uint32(2 * period * count), 0x7fffffff, 0x80000000, 0xffffffff}
					for _, frame := range frames {
						t.Run(fmt.Sprintf("op%d-count%d-period%d-mode%d-frame%d", op, count, period, mode, frame), func(t *testing.T) {
							c.resetCase(effects, pix, 17, frame)
							env.Reset()
							env.Constants(0.5, 65536)
							d := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(48, 48))
							raw := unsafe.Slice((*byte)(d.C()), int(unsafe.Sizeof(*d)))
							clear(raw[136:276])
							clear(raw[432:435])
							binary.LittleEndian.PutUint32(raw[140:], math.Float32bits(7))
							binary.LittleEndian.PutUint32(raw[144:], 123)
							for i := 0; i < 16; i++ {
								raw[178+3*i] = byte(7 + 13*i)
								raw[179+3*i] = byte(239 - 11*i)
								raw[180+3*i] = byte(3 + 7*i)
								raw[226+i] = byte(3 + 4*i)
								raw[242+i] = byte(9 + 13*i)
							}
							raw[432+op] = byte(count)
							binary.LittleEndian.PutUint16(raw[258+2*op:], uint16(period))
							raw[176] = byte(mode << (2 * op))
							before := append([]byte(nil), raw[136:176]...)
							wholeBefore := append([]byte(nil), raw...)
							got := legacy.PortTestColorLight(op, c.Viewport(), d)
							light := make([]uint32, 10)
							for i := range light {
								light[i] = binary.LittleEndian.Uint32(raw[136+4*i:])
							}
							if count <= 1 || count >= 128 || period == 0 {
								for i, v := range before {
									if raw[136+i] != v {
										t.Fatalf("disabled animation changed byte %d", i)
									}
								}
							} else if frame == 0 {
								switch op {
								case 0:
									if light[0] != 2 || light[4] != 7 || light[5] != 239 || light[6] != 3 {
										t.Fatal("initial RGB entry")
									}
								case 1:
									if light[1] != math.Float32bits(3) {
										t.Fatal("initial intensity entry")
									}
								case 2:
									if uint16(light[7]>>16) != uint16(9*65536/360+0) {
										t.Fatal("initial penumbra entry")
									}
								}
							}
							if c.Objs.Count != 1 || len(c.Calls) != 1 || c.srv.Rand.Logic.Index() != 17 || c.srv.Rand.Other.Index() != 18 {
								t.Fatal("animation changed ownership or RNG")
							}
							if op == 1 && count > 1 && count < 128 && period != 0 {
								wantFixed := uint32(int64(float64(math.Float32frombits(light[1]))*65536 + 0.5))
								if light[3] != wantFixed {
									t.Fatal("animated intensity did not update fixed-point intensity")
								}
							}
							for off, value := range wholeBefore {
								allowed := op == 0 && (off >= 136 && off < 140 || off >= 152 && off < 164) || op == 1 && off >= 140 && off < 152 || op == 2 && off >= 166 && off < 168
								if !allowed && raw[off] != value {
									t.Fatalf("changed unrelated drawable byte %d", off)
								}
							}
							records = append(records, record{op, count, period, mode, frame, got, light})
						})
					}
				}
			}
		}
	}
	spellbookCapture(t, "color-light-interpolation", records, "48941b54fc091b4f24354f8cb0c24cf52a63ddbc8055f8f285e5fcdcd44b1d4e")
}
