//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"testing"
	"unsafe"
)

func TestClientObjectDrawingArrowRetry(t *testing.T) {
	o := newObjectDrawingOwner(t)
	c := o.c
	for _, op := range []int{1, 2} {
		o.reset(1, 120)
		o.data[0], o.data[2] = 12, 32
		dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(48, 48))
		dr.DrawData = unsafe.Pointer(&o.data[0])
		dr.DrawFuncPtr = legacy.PortTestObjectDrawCallback(op)
		dr.Field_81, dr.Field_82 = 28, 48
		c.FailEvery = 2
		blank := effectsPixelHash(o.pix)
		if got := dr.CallDraw(c.Viewport()); got != 1 {
			t.Fatalf("arrow%d failed main draw: %d", op, got)
		}
		if c.Objs.Count != 1 || c.Objs.DeadlineList != nil || dr.Field_81 != 28 || dr.Field_82 != 48 {
			t.Fatal("failed optional tail changed ownership or retry anchor")
		}
		if len(c.Calls) != 2 || c.Calls[1].Ref != 0 {
			t.Fatal("tail allocation failure was not exercised")
		}
		if effectsPixelHash(o.pix) == blank {
			t.Fatal("failed tail suppressed main sprite pixels")
		}
		c.FailEvery = 0
		if got := dr.CallDraw(c.Viewport()); got != 1 {
			t.Fatal("arrow retry did not draw")
		}
		tail := c.Objs.DeadlineList
		if tail == nil || c.Objs.Count != 2 || tail.Deadline != 140 || dr.Field_81 != 48 || dr.Field_82 != 48 {
			t.Fatal("successful retry did not establish tail lifetime/anchor")
		}
		if tail.PosVec != image.Pt(28, 48) {
			t.Fatal("retry lost original tail anchor")
		}
		words := unsafe.Slice((*uint32)(tail.C()), 128)
		if words[108] != 48 || words[109] != 48 {
			t.Fatal("tail endpoint does not match the sprite")
		}
		o.snapshot(t)
	}
}

func TestClientObjectDrawingGlyphOpaque(t *testing.T) {
	o := newObjectDrawingOwner(t)
	c := o.c
	var prior string
	for padding := 0; padding < 4; padding++ {
		o.reset(1, 120)
		o.data[3] = 5
		for i := 0; i < padding; i++ {
			c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(10, 10))
		}
		dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(48, 48))
		dr.DrawData = unsafe.Pointer(&o.data[0])
		dr.DrawFuncPtr = legacy.PortTestObjectDrawCallback(5)
		if got := dr.CallDraw(c.Viewport()); got != 1 {
			t.Fatal("glyph draw failed")
		}
		actual := effectsPixelHash(o.pix)
		clear(o.pix.Pix)
		c.r.Data().SetAlphaEnabled(true)
		c.r.Data().SetAlpha(255)
		dr.DrawFuncPtr = legacy.PortTestSpriteAnimationCallback(0)
		dr.CallDraw(c.Viewport())
		if actual != effectsPixelHash(o.pix) {
			t.Fatal("default glyph is not opaque")
		}
		if padding > 0 && actual != prior {
			t.Fatal("glyph opacity depends on allocation position")
		}
		prior = actual
	}
}

func TestClientObjectDrawingGeneratorRandomBounds(t *testing.T) {
	o := newObjectDrawingOwner(t)
	c := o.c
	for state := 0; state < 4; state++ {
		for seed := uint32(1); seed <= 128; seed++ {
			o.reset(seed, 120)
			bytes := unsafe.Slice((*byte)(unsafe.Pointer(&o.data[0])), 56)
			clear(o.data)
			o.data[0] = 16
			for i := 0; i < 5; i++ {
				o.data[i+1] = uint32(uintptr(unsafe.Pointer(&o.frames[32*i])))
				bytes[24+i] = 1
				bytes[29+i] = 0
				o.data[9+i] = 4
			}
			dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(48, 48))
			dr.DrawData = unsafe.Pointer(&o.data[0])
			dr.DrawFuncPtr = legacy.PortTestObjectDrawCallback(12)
			words := unsafe.Slice((*uint32)(dr.C()), 128)
			words[70] = []uint32{0, 0x100, 0x200, 0x400}[state]
			if state == 3 {
				words[108] = 1
			} // valid final countdown; completed state clamps frame
			before := c.srv.Rand.Other.Index()
			if got := dr.CallDraw(c.Viewport()); got != 1 {
				t.Fatal("generator draw failed")
			}
			if len(o.drawTrace) == 0 || o.drawTrace[0].Position != image.Pt(48+state, 48) {
				t.Fatalf("one-frame generator state=%d seed=%d trace=%v expected=%v", state, seed, o.drawTrace, image.Pt(48+state, 48))
			}
			if c.srv.Rand.Other.Index() != before+1 {
				t.Fatal("generator random must consume exactly one RNG value")
			}
			o.snapshot(t)
		}
	}
}
