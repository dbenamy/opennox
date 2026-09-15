//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/legacy/common/alloc/handles"
	"image"
	"testing"
	"unsafe"
)

func TestClientSpriteAnimationBoulder(t *testing.T) {
	t.Cleanup(handles.PortTestInit())
	c, pix, env := newEffectsFullOwner(t)
	t.Cleanup(legacy.PortTestSpriteAnimationEnvironment())
	raw := make([][]byte, 32)
	for i := range raw {
		raw[i] = spriteAnimationTestImage(i)
	}
	imgs, free := c.r.GetBag().PortTestSpriteImages(raw)
	t.Cleanup(free)
	c.imageRefs = map[uint32]uint32{}
	for i, img := range imgs {
		c.imageRefs[uint32(uintptr(img.C()))] = 0xe3000000 + uint32(i)
	}
	frames, freeFrames := alloc.Make([]uint32{}, 32)
	defer freeFrames()
	for i, img := range imgs {
		frames[i] = uint32(uintptr(img.C()))
	}
	data, freeData := alloc.Make([]uint32{}, 3)
	defer freeData()
	data[0], data[1], data[2] = 12, uint32(uintptr(unsafe.Pointer(&frames[0]))), 32
	c.dataRefs = map[uint32]uint32{uint32(uintptr(unsafe.Pointer(&data[0]))): 0xe4000001}
	c.callbackRefs = map[unsafe.Pointer]uint32{legacy.PortTestSpriteAnimationCallback(4): 0xe5000004}
	type result struct {
		Delta                       image.Point
		Initial, Bank, Step, Return int
		Pixels                      string
		Drawables                   [][]uint32
		Render                      []uint32
	}
	var out []result
	for _, delta := range []image.Point{{0, 0}, {9, 0}, {10, 0}, {11, 0}, {0, 9}, {0, 10}, {0, -10}, {-10, 0}, {6, 8}, {6, 7}, {-6, 8}, {6, -8}, {-6, -8}} {
		for _, initial := range []int{0, 1, 14, 15} {
			for _, bank := range []int{0, 16} {
				c.resetCase(env, pix, 1, 120)
				dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(48, 48))
				words := unsafe.Slice((*uint32)(dr.C()), 128)
				words[30] = 0x40000000
				words[76] = uint32(uintptr(unsafe.Pointer(&data[0])))
				words[110], words[111] = uint32(initial), uint32(bank)
				dr.DrawFuncPtr = legacy.PortTestSpriteAnimationCallback(4)
				for step := 0; step < 5; step++ {
					if step > 0 {
						pos := dr.PosVec.Add(delta)
						c.Nox_xxx_updateSpritePosition_49AA90(dr, pos.X, pos.Y)
					}
					clear(pix.Pix)
					got := dr.CallDraw(c.Viewport())
					out = append(out, result{delta, initial, bank, step, got, effectsPixelHash(pix), c.snapshotDrawables(t), append([]uint32(nil), unsafe.Slice((*uint32)(unsafe.Pointer(c.r.Data())), 264)...)})
				}
			}
		}
	}
	effectsCapture(t, "sprite-animation-boulder", out, len(out), "a5aa6343cc7b85d21c13caf0d2bf49d6f838ab47e94480bb5fd8404bd1647a0f")
}
