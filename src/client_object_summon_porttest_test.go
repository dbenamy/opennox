//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"image"
	"testing"
	"unsafe"
)

func TestClientObjectDrawingSummon(t *testing.T) {
	o := newObjectDrawingOwner(t)
	c := o.c
	childData, free := alloc.Make([]uint32{}, 2)
	defer free()
	childData[0], childData[1] = 8, o.frames[9]
	c.dataRefs[uint32(uintptr(unsafe.Pointer(&childData[0])))] = 0xe9000002
	var out []objectDrawingResult
	id := 0
	for _, start := range []uint32{0, 120, 0xfffffffc} {
		for _, duration := range []uint32{0, 1, 2, 7, 31, 65535} {
			for _, count := range []byte{1, 7, 32} {
				for _, pos := range []image.Point{{48, 48}, {0, 0}, {5887, 5887}} {
					id++
					o.reset(uint32(id), start)
					b := unsafe.Slice((*byte)(unsafe.Pointer(&o.data[0])), 56)
					b[8], b[9] = count, byte(id%3)
					dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, pos)
					dr.DrawData = unsafe.Pointer(&o.data[0])
					dr.DrawFuncPtr = legacy.PortTestObjectDrawCallback(6)
					child := o.newUnlinked(t)
					child.PosVec = pos
					child.DrawData = unsafe.Pointer(&childData[0])
					child.DrawFuncPtr = legacy.PortTestSpriteAnimationCallback(2)
					w := unsafe.Slice((*uint32)(dr.C()), 128)
					w[108] = uint32(uintptr(child.C()))
					w[109] = duration | 0xa1230000
					w[77] = uint32(count - 1)
					w[32] = 0 - start
					ages := []uint32{0}
					if duration > 1 {
						ages = append(ages, duration/2)
					}
					if duration > 0 {
						ages = append(ages, duration-1, duration)
					}
					priorAge := ^uint32(0)
					for step, age := range ages {
						if age == priorAge {
							continue
						}
						priorAge = age
						c.srv.SetFrame(start + age)
						clear(o.pix.Pix)
						o.drawTrace = nil
						ret := dr.CallDraw(c.Viewport())
						if age < duration {
							if ret != 1 || dr.PosVec != pos || w[77] != uint32(count-1) || o.data[3] != 2 {
								t.Fatal("summon did not restore parent animation/position")
							}
							if len(o.unlinked) != 1 || len(o.rawDeleted) != 0 {
								t.Fatal("live summon lost child ownership")
							}
							if len(o.drawTrace) < 2 || o.drawTrace[len(o.drawTrace)-1].Image != c.imageRefs[childData[1]] {
								t.Fatal("summon did not dispatch child draw")
							}
						} else {
							if ret != 0 || len(o.unlinked) != 0 || len(o.rawDeleted) != 1 || len(c.Deleted) != 1 {
								t.Fatal("summon expiration did not delete parent and child exactly once")
							}
						}
						out = append(out, o.result(t, id, step, ret))
					}
				}
			}
		}
	}
	effectsCapture(t, "object-drawing-summon", out, len(out), "28d66af6b2db0802b7e6810d7a5430332d4ac5b885a8d32595ecc7c41752d560")
}
