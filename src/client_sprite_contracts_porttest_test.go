//go:build porttest

package opennox

import (
	"github.com/opennox/libs/prand"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/legacy/common/alloc/handles"
	"image"
	"testing"
	"unsafe"
)

func TestClientSpriteAnimationRandomBounds(t *testing.T) {
	t.Cleanup(handles.PortTestInit())
	c, pix, env := newEffectsFullOwner(t)
	t.Cleanup(legacy.PortTestSpriteAnimationEnvironment())
	images, free := c.r.GetBag().PortTestSpriteImages([][]byte{spriteAnimationTestImage(0), spriteAnimationTestImage(1)})
	t.Cleanup(free)
	data, freeData := alloc.Make([]uint32{}, 14)
	defer freeData()
	// An owned sentinel makes a one-past selection observable without reading
	// outside fixture memory. Both active states declare exactly one frame.
	frames, freeFrames := alloc.Make([]uint32{}, 2)
	defer freeFrames()
	for i, img := range images {
		frames[i] = uint32(uintptr(img.C()))
	}
	data[0] = 16
	data[1] = uint32(uintptr(unsafe.Pointer(&frames[0])))
	data[2] = data[1]
	bytes := unsafe.Slice((*byte)(unsafe.Pointer(&data[0])), 56)
	bytes[24], bytes[25] = 1, 1
	data[9], data[10] = 4, 4
	oldBoundHits := 0
	for seed := uint32(1); seed <= 128; seed++ {
		if prand.New(int(seed+1)).Int(0, 1) == 1 {
			oldBoundHits++
		}
		c.resetCase(env, pix, seed, 120)
		dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(48, 48))
		words := unsafe.Slice((*uint32)(dr.C()), 128)
		words[30] = 0x40000000
		if seed&1 != 0 {
			words[30] |= 0x1000000
		}
		words[76] = uint32(uintptr(unsafe.Pointer(&data[0])))
		dr.DrawFuncPtr = legacy.PortTestSpriteAnimationCallback(1)
		before := c.srv.Rand.Other.Index()
		if got := dr.CallDraw(c.Viewport()); got != 1 {
			t.Fatalf("return %d", got)
		}
		if words[2] != frames[0] {
			t.Fatal("one-frame conditional animation selected its sentinel")
		}
		if c.srv.Rand.Other.Index() != before+1 {
			t.Fatal("conditional random must consume exactly one RNG value")
		}
	}
	if oldBoundHits == 0 {
		t.Fatal("boundary seeds do not distinguish the old inclusive bound")
	}
}
