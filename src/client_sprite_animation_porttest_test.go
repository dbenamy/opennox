//go:build porttest

package opennox

import (
	"encoding/binary"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/legacy/common/alloc/handles"
	"image"
	"testing"
	"unsafe"
)

func spriteAnimationTestImage(frame int) []byte {
	const width, height = 7, 5
	b := make([]byte, 17)
	binary.LittleEndian.PutUint32(b, width)
	binary.LittleEndian.PutUint32(b[4:], height)
	for y := 0; y < height; y++ {
		b = append(b, 2, width)
		for x := 0; x < width; x++ {
			color := uint16(0x8001 | ((frame+1)*47+x*131+y*997)&0x7ffe)
			b = binary.LittleEndian.AppendUint16(b, color)
		}
	}
	return b
}

func TestClientSpriteAnimationImageOwner(t *testing.T) {
	t.Cleanup(handles.PortTestInit())
	if legacy.PortTestSpriteAnimationCallback(3) == legacy.PortTestSpriteAnimationCallback(5) {
		t.Fatal("static-random and slave callbacks must retain distinct identities")
	}
	c, pix, _ := newEffectsFullOwner(t)
	t.Cleanup(legacy.PortTestSpriteAnimationEnvironment())
	images, free := c.r.GetBag().PortTestSpriteImages([][]byte{spriteAnimationTestImage(0), spriteAnimationTestImage(1)})
	t.Cleanup(free)
	data, freeData := alloc.Make([]uint32{}, 2)
	defer freeData()
	data[0] = 8
	dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(48, 48))
	words := unsafe.Slice((*uint32)(dr.C()), 128)
	words[30] = 0x40000000 // fixed lighting, no world-light owner needed
	words[76] = uint32(uintptr(unsafe.Pointer(&data[0])))
	dr.DrawFuncPtr = legacy.PortTestSpriteAnimationCallback(2)
	blank := effectsPixelHash(pix)
	var prior string
	for i, img := range images {
		clear(pix.Pix)
		data[1] = uint32(uintptr(img.C()))
		if got := dr.CallDraw(c.Viewport()); got != 1 {
			t.Fatalf("draw return %d", got)
		}
		if words[2] != data[1] {
			t.Fatal("draw metadata did not retain owned image handle")
		}
		hash := effectsPixelHash(pix)
		if hash == blank || hash == prior {
			t.Fatalf("frame %d did not render distinct pixels", i)
		}
		prior = hash
	}
}
