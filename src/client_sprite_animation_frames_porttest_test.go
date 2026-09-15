//go:build porttest

package opennox

import (
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/legacy/common/alloc/handles"
	"image"
	"testing"
	"unsafe"
)

func TestClientSpriteAnimationFrames(t *testing.T) {
	t.Cleanup(handles.PortTestInit())
	c, pix, env := newEffectsFullOwner(t)
	t.Cleanup(legacy.PortTestSpriteAnimationEnvironment())
	oldFlags := noxflags.GetGame()
	t.Cleanup(func() { noxflags.ResetGame(); noxflags.SetGame(oldFlags) })
	raw := make([][]byte, 255)
	for i := range raw {
		raw[i] = spriteAnimationTestImage(i)
	}
	images, free := c.r.GetBag().PortTestSpriteImages(raw)
	t.Cleanup(free)
	c.imageRefs = make(map[uint32]uint32)
	for i, img := range images {
		c.imageRefs[uint32(uintptr(img.C()))] = 0xe3000000 + uint32(i)
	}
	data, freeData := alloc.Make([]uint32{}, 14)
	defer freeData()
	frames, freeFrames := alloc.Make([]uint32{}, 510)
	defer freeFrames()
	for i := range frames {
		frames[i] = uint32(uintptr(images[(i+(i/255)*113)%255].C()))
	}
	c.dataRefs = map[uint32]uint32{uint32(uintptr(unsafe.Pointer(&data[0]))): 0xe4000001}
	c.callbackRefs = make(map[unsafe.Pointer]uint32)
	for _, op := range []int{0, 1, 2, 3, 5} {
		c.callbackRefs[legacy.PortTestSpriteAnimationCallback(op)] = 0xe5000000 + uint32(op)
	}
	type result struct {
		Op, Mode, Count, Delay, Variant, Step, Return int
		Frame                                         uint32
		Pixels                                        string
		Render, Globals                               []uint32
		Drawables                                     [][]uint32
		Deleted                                       []uint32
		Logic, Other                                  int
	}
	var out []result
	for _, op := range []int{0, 1, 2, 3, 5} {
		modes := []int{0, 1, 2, 3, 4, 5, 6}
		if op >= 2 {
			modes = []int{0}
		}
		for _, mode := range modes {
			for _, count := range []int{1, 2, 7, 31, 255} {
				for _, delay := range []int{0, 1, 7, 255} {
					for variant := 0; variant < 8; variant++ {
						c.resetCase(env, pix, uint32(1+31*variant+op*127+mode*13), 120)
						noxflags.ResetGame()
						if variant == 3 {
							noxflags.SetGame(noxflags.GameFlag(32))
						}
						clear(data)
						data[0] = 16
						data[1] = uint32(uintptr(unsafe.Pointer(&frames[0])))
						bytes := unsafe.Slice((*byte)(unsafe.Pointer(&data[0])), 56)
						bytes[8], bytes[9] = byte(count), byte(delay)
						data[3] = uint32(mode)
						if op == 1 {
							data[2] = uint32(uintptr(unsafe.Pointer(&frames[255])))
							bytes[24], bytes[25] = byte(count), byte(count)
							bytes[29], bytes[30] = byte(delay), byte(255-delay)
							data[9], data[10] = uint32(mode), uint32(mode)
						} else if op == 2 {
							data[0], data[1] = 8, frames[variant%count]
						}
						dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, []image.Point{{48, 48}, {94, 94}, {1, 1}, {50, 30}}[variant%4])
						words := unsafe.Slice((*uint32)(dr.C()), 128)
						words[30] = 0x40000000
						if variant&1 == 0 {
							words[30] |= 0x1000000
						}
						if variant >= 2 {
							words[28] |= 0x10000000
						}
						if op == 2 && variant >= 2 {
							words[28] |= 0x40000
						}
						words[76] = uint32(uintptr(unsafe.Pointer(&data[0])))
						words[77] = uint32(variant % count)
						start := uint32(120)
						if variant >= 4 {
							start = 0xffffff00
						}
						words[79] = start
						words[32] = 0
						if variant >= 4 {
							words[32] = 0 - start
						}
						dr.DrawFuncPtr = legacy.PortTestSpriteAnimationCallback(op)
						period := uint32(delay + 1)
						for step, elapsed := range []uint32{0, period - 1, period, uint32(count)*period - 1, uint32(count) * period, uint32(2*count) * period} {
							frame := start + elapsed
							c.srv.SetFrame(frame)
							clear(pix.Pix)
							*memmap.PtrUint32(0x5D4594, 1321512) = 0
							*memmap.PtrUint32(0x5D4594, 1321516) = 0
							got := dr.CallDraw(c.Viewport())
							last := *memmap.PtrUint32(0x5D4594, 1321516)
							if last != 0 {
								if last != uint32(uintptr(dr.C())) {
									t.Fatal("unowned last-drawn object")
								}
								last = c.refs[dr]
							}
							out = append(out, result{op, mode, count, delay, variant, step, got, frame, effectsPixelHash(pix), append([]uint32(nil), unsafe.Slice((*uint32)(unsafe.Pointer(c.r.Data())), 264)...), []uint32{*memmap.PtrUint32(0x5D4594, 1321512), last}, c.snapshotDrawables(t), append([]uint32(nil), c.Deleted...), c.srv.Rand.Logic.Index(), c.srv.Rand.Other.Index()})
							if c.Objs.Count == 0 {
								break
							}
						}
					}
				}
			}
		}
	}
	effectsCapture(t, "sprite-animation-frames", out, len(out), "648484bf7f0f989ee84d59ff9e38f35a8251a8fefc5b3d5cbf9ddb4f574e8d55")
}
