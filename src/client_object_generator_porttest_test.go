//go:build porttest

package opennox

import (
	"image"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/legacy"
)

type objectDrawingResult struct {
	Case, Step, Return  int
	Pixels              string
	Drawables           [][]uint32
	Render, Globals     []uint32
	Images              []objectImageDraw
	Names               []string
	Calls               []effectsSpawnCall
	Deleted, RawDeleted []uint32
	Logic, Other        int
}

func (o *objectDrawingOwner) result(t *testing.T, id, step, ret int) objectDrawingResult {
	return objectDrawingResult{id, step, ret, effectsPixelHash(o.pix), o.snapshot(t),
		append([]uint32(nil), unsafe.Slice((*uint32)(unsafe.Pointer(o.c.r.Data())), 264)...),
		o.globals(t), append([]objectImageDraw(nil), o.drawTrace...), append([]string(nil), o.namedCalls...),
		append([]effectsSpawnCall(nil), o.c.Calls...), append([]uint32(nil), o.c.Deleted...),
		append([]uint32(nil), o.rawDeleted...), o.c.srv.Rand.Logic.Index(), o.c.srv.Rand.Other.Index()}
}
func (o *objectDrawingOwner) generatorData(count, delay, kind int) {
	clear(o.data)
	o.data[0] = 56
	b := unsafe.Slice((*byte)(unsafe.Pointer(&o.data[0])), 56)
	for i := 0; i < 5; i++ {
		o.data[i+1] = uint32(uintptr(unsafe.Pointer(&o.frames[32*i])))
		b[24+i], b[29+i] = byte(count), byte(delay)
		o.data[9+i] = uint32(kind)
	}
}
func TestClientObjectDrawingGenerator(t *testing.T) {
	o := newObjectDrawingOwner(t)
	c := o.c
	var out []objectDrawingResult
	id := 0
	for _, start := range []uint32{0, 120, 0xfffffffc} {
		for _, count := range []int{1, 2, 7, 32} {
			for _, delay := range []int{0, 1, 7, 255} {
				for _, kind := range []int{0, 2, 4, 5, 1} {
					for _, state := range []uint32{0, 0x100, 0x200, 0x400, 0x800, 0x300} {
						if kind == 0 && state != 0x400 && state != 0x800 {
							continue
						}
						id++
						o.reset(uint32(id), start)
						o.generatorData(count, delay, kind)
						dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(48, 48))
						dr.DrawData = unsafe.Pointer(&o.data[0])
						dr.DrawFuncPtr = legacy.PortTestObjectDrawCallback(12)
						w := unsafe.Slice((*uint32)(dr.C()), 128)
						w[70] = state
						w[32] = 0 - start
						w[77] = uint32(count - 1)
						// The completed transition must clear both lighting flags.
						w[28] |= 0x80000
						w[30] |= 0x21000000
						if state == 0x400 {
							p := legacy.PortTestObjectDrawHelper(5, dr, nil, 0)
							if p != uint32(uintptr(dr.DrawData)) || w[108] != uint32(count*(delay+1)) {
								t.Fatal("countdown initialization contract")
							}
						}
						if legacy.PortTestObjectDrawHelper(6, nil, nil, 0) != 1 {
							t.Fatal("generator update contract")
						}
						duration := uint32(count * (delay + 1))
						for step, age := range []uint32{0, 1, 2, duration - 1, duration, duration + 1} {
							c.srv.SetFrame(start + age)
							// Jump through the actual remaining countdown, including its final tick.
							if state == 0x400 && age < duration {
								w[108] = duration - age
							}
							if state == 0x400 && age >= duration {
								w[108] = 0
								w[70] = 0x800
							}
							clear(o.pix.Pix)
							o.drawTrace = nil
							ret := dr.CallDraw(c.Viewport())
							if kind == 1 {
								if ret != 0 || len(o.drawTrace) != 0 {
									t.Fatal("unsupported generator mode drew")
								}
							} else if ret != 1 {
								t.Fatal("valid generator failed")
							}
							if state == 0x400 && age == duration-1 && kind != 1 && (w[108] != 0 || w[70]&0xc00 != 0x800) {
								t.Fatal("final countdown transition")
							}
							out = append(out, o.result(t, id, step, ret))
						}
					}
				}
			}
		}
	}
	effectsCapture(t, "object-drawing-generator", out, len(out), "22a5ccdafb27a8fe752425ed4b60276767b40a11aaa821febe2b4e0cd07d1eee")
}
