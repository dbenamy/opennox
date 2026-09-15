//go:build porttest

package opennox

import (
	"encoding/binary"
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"testing"
	"unsafe"
)

func TestClientScreenEffectsDrawables(t *testing.T) {
	o := newObjectDrawingOwner(t)
	c := o.c
	env := legacy.PortTestNewScreenEnvironment()
	defer env.Restore()
	for op := 0; op < 5; op++ {
		c.callbackRefs[legacy.PortTestScreenEffectCallback(op)] = 0xec000000 + uint32(op)
	}
	type result struct {
		Op, Variant, Step int
		Start, Seed       uint32
		Draw              objectDrawingResult
		Globals           []uint32
	}
	var out []result
	id := 0
	for _, op := range []int{0, 2, 3} {
		for variant := 0; variant < 8; variant++ {
			for _, start := range []uint32{0, 120, 0xffffffe0} {
				for _, seed := range []uint32{1, 2, 17, 31} {
					id++
					o.reset(seed, start)
					env.Reset()
					o.data[0], o.data[2] = 12, 32
					pos := []image.Point{{48, 48}, {1, 1}, {95, 95}, {48, 1}}[variant%4]
					if variant >= 4 {
						c.Viewport().Screen = image.Rect(6, 9, 90, 93)
						c.Viewport().World = image.Rect(10, 20, 94, 104)
						c.Viewport().Size = c.Viewport().Screen.Size()
					}
					dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, pos)
					dr.DrawFuncPtr = legacy.PortTestScreenEffectCallback(op)
					if op == 0 {
						dr.DrawData = unsafe.Pointer(&o.data[0])
					}
					dr.Field_81, dr.Field_82 = uint32(pos.X+[]int{0, 39, 40, 79, 319, 320, 359, 360}[variant]), uint32(pos.Y)
					ages := []uint32{0, 1, 63, 127, 128, 255}
					if op == 2 {
						ages = []uint32{0, 1, 69, 70, 71}
					}
					for step, age := range ages {
						c.srv.SetFrame(start + age)
						dr.AnimFrameSlave = uint32((variant + step*7) % 32)
						clear(o.pix.Pix)
						o.drawTrace = nil
						ret := dr.CallDraw(c.Viewport())
						out = append(out, result{op, variant, step, start, seed, o.result(t, id, step, ret), env.State()})
						if ret == 0 {
							if op != 2 || age != 71 {
								t.Fatal("unexpected effect lifetime")
							}
							break
						}
					}
				}
			}
		}
	}
	effectsCapture(t, "screen-effects-drawables", out, len(out), "b3c8897faa5ff57c99b2d71ea2f446e4df2ac4830d1eafb501b4471bbad6c628")
}
func TestClientScreenEffectsRope(t *testing.T) {
	o := newObjectDrawingOwner(t)
	c := o.c
	env := legacy.PortTestNewScreenEnvironment()
	defer env.Restore()
	callback := legacy.PortTestScreenEffectCallback(1)
	c.callbackRefs[callback] = 0xec000001
	type result struct {
		Mode, Missing, Direction, Step int
		Delta                          image.Point
		Draw                           objectDrawingResult
		Globals                        []uint32
	}
	var out []result
	id := 0
	for mode := 0; mode < 5; mode++ {
		for missing := 0; missing < 3; missing++ {
			for _, dir := range []int{0, 1, 3, 5, 7, 8} {
				for _, delta := range []image.Point{{0, 0}, {30, 0}, {0, 30}, {30, 30}, {30, -30}, {-30, 15}, {15, -30}} {
					id++
					o.reset(uint32(id), 120)
					env.Reset()
					from := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(48, 48))
					from.NetCode32 = 7
					from.AnimDir = byte(dir)
					to := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(48, 48).Add(delta))
					to.NetCode32 = 8
					codes := [2]uint32{7, 8}
					mask := mode - 1
					if mode > 0 && mask&1 != 0 {
						w := unsafe.Slice((*uint32)(from.C()), 128)
						w[28] |= 0x20000000
						codes[0] |= 0x8000
					}
					if mode > 0 && mask&2 != 0 {
						w := unsafe.Slice((*uint32)(to.C()), 128)
						w[28] |= 0x20000000
						codes[1] |= 0x8000
					}
					if missing == 1 {
						codes[0] = 99
					}
					if missing == 2 {
						codes[1] = 99
					}
					dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(48, 48))
					dr.DrawFuncPtr = callback
					b := unsafe.Slice((*byte)(dr.C()), 512)
					if mode == 0 {
						for i, v := range []int{48, 48, 48 + delta.X, 48 + delta.Y} {
							binary.LittleEndian.PutUint16(b[437+i*2:], uint16(v))
						}
					} else {
						b[432] = 1
						binary.LittleEndian.PutUint32(b[437:], codes[0])
						binary.LittleEndian.PutUint32(b[441:], codes[1])
					}
					for step := 0; step < 2; step++ {
						to.ZVal, to.ZVal2 = uint16(step*3), uint16(step*2)
						if step > 0 {
							c.Viewport().World = image.Rect(10, 20, 106, 116)
						}
						clear(o.pix.Pix)
						ret := dr.CallDraw(c.Viewport())
						if ret != 1 {
							t.Fatal("rope return")
						}
						out = append(out, result{mode, missing, dir, step, delta, o.result(t, id, step, ret), env.State()})
					}
				}
			}
		}
	}
	effectsCapture(t, "screen-effects-rope", out, len(out), "237972bd0664fcf9a99e92a7943e20416b8404a32676ca1a206f9bfc49595305")
}
