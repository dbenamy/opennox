//go:build porttest

package opennox

import (
	"encoding/binary"
	"fmt"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"image"
	"testing"
	"unsafe"
)

func TestCombatOverlayEffectHistory(t *testing.T) {
	type record struct {
		Name    string
		Step    int
		Count   uint32
		History [12]uint32
		Pixels  string
		Alpha   byte
	}
	var rows []record
	for _, kind := range []int{1, 2} {
		for _, direction := range []uint32{0, 7, 15, 23, 31} {
			name := fmt.Sprintf("kind=%d/direction=%d", kind, direction)
			t.Run(name, func(t *testing.T) {
				o := newCombatOverlayOwner(t)
				light := serverConfigOwnBytes(t, 0x587000, 185472, 12)
				for i := 0; i < 3; i++ {
					binary.LittleEndian.PutUint32(light[4*i:], 255)
				}
				dr := o.drawable(7, image.Pt(300, 300))
				dr.ObjFlags |= 0x40000000
				dr.AnimFrameSlave = direction
				raw := spriteAnimationTestImage(0)
				// Opaque pixels are required when the trail enables alpha blending.
				for y := 0; y < 5; y++ {
					for x := 0; x < 7; x++ {
						raw[17+y*16+2+x*2+1] &= 0x7f
					}
				}
				images, freeImages := o.c.r.GetBag().PortTestSpriteImages([][]byte{raw})
				t.Cleanup(freeImages)
				data, freeData := alloc.New([2]uint32{})
				*data = [2]uint32{8, uint32(uintptr(images[0].C()))}
				t.Cleanup(freeData)
				dr.DrawData = unsafe.Pointer(data)
				dr.DrawFuncPtr = legacy.PortTestSpriteAnimationCallback(2)
				clear(o.pix.Pix)
				blankBefore := effectsPixelHash(o.pix)
				dr.CallDraw(o.c.Viewport())
				if kind == 1 && effectsPixelHash(o.pix) == blankBefore {
					t.Fatalf("static sprite prerequisite does not draw: class=%x flags=%x pos=%v metadata=%v", uint32(dr.ObjClass), uint32(dr.ObjFlags), dr.PosVec, dr.Field_2)
				}
				p, free := alloc.New([20]uint32{})
				t.Cleanup(free)
				p[0] = uint32(kind)
				p[14] = 0xaabbcc00
				legacy.PortTestCombatEffectAttach(unsafe.Pointer(p), dr)
				t.Cleanup(func() { legacy.PortTestCombatEffectDetach(unsafe.Pointer(p)); dr.Field_114 = nil })
				previous := image.Pt(290, 290)
				for step := 0; step < 15; step++ {
					pos := image.Pt(300+6*min(step, 6), 300+3*min(step, 6))
					dr.Field_8, dr.Field_9 = uint32(previous.X), uint32(previous.Y)
					dr.PosVec = pos
					clear(o.pix.Pix)
					blank := effectsPixelHash(o.pix)
					legacy.PortTestCombatEffectDraw(o.c.Viewport(), dr)
					count := p[14] & 255
					if p[14]&0xffffff00 != 0xaabbcc00 || count > 5 {
						t.Fatal("effect count/reserved bytes")
					}
					if step < 7 {
						if count != uint32(min(step+1, 5)) || p[2] != uint32(pos.X) || p[3] != uint32(pos.Y) {
							t.Fatal("moving effect history")
						}
						if step > 0 && effectsPixelHash(o.pix) == blank {
							t.Fatalf("moving effect did not draw: step=%d pos=%v flags=%x class=%x metadata=%v", step, dr.PosVec, dr.ObjFlags, dr.ObjClass, dr.Field_2)
						}
					}
					if dr.PosVec != pos {
						t.Fatal("effect drawing did not restore drawable position")
					}
					if o.c.r.Data().IsAlphaEnabled() {
						t.Fatal("effect drawing left alpha enabled")
					}
					var history [12]uint32
					copy(history[:], p[2:14])
					rows = append(rows, record{name, step, count, history, effectsPixelHash(o.pix), o.c.r.Data().Alpha()})
					previous = pos
				}
				if p[14]&255 != 0 {
					t.Fatal("stationary effect history did not drain")
				}
			})
		}
	}
	spellbookCapture(t, "combat-overlay-effect-history", rows, "18d29386e73d4cb08d44a912b63b843801ac3afff7eb48c4f358a54a5b87849b")
}
