//go:build porttest

package opennox

import (
	"encoding/binary"
	noxcolor "github.com/opennox/libs/color"
	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"testing"
	"unsafe"
)

func TestClientPresentationBubbles(t *testing.T) {
	names := []string{"RedBubbleParticle", "WhiteBubbleParticle", "LightBlueBubbleParticle", "OrangeBubbleParticle", "GreenBubbleParticle", "VioletBubbleParticle", "LightVioletBubbleParticle", "YellowBubbleParticle", "OtherBubbleParticle"}
	palette := [][3]byte{{255, 128, 128}, {255, 255, 255}, {200, 200, 255}, {255, 100, 50}, {64, 255, 64}, {255, 100, 255}, {255, 200, 255}, {255, 255, 200}, {200, 200, 255}}
	c, pix, env := newEffectsFullOwner(t, names...)
	cache := serverConfigOwnBytes(t, 0x5D4594, 1217512, 32)
	type record struct {
		Kind              int
		Frame             uint32
		Z                 int16
		Args              [5]byte
		Lifetime, Failure int
		Data              [15]byte
		Deadline          uint32
	}
	var rows []record
	for kind, name := range names {
		typ := c.Things.TypeByID(name)
		typ.LightColor = noxrender.RGB{R: 0x12345621, G: 0x12345683, B: 0x123456e7}
		for _, frame := range []uint32{120, 0xfffffff0} {
			for _, z := range []int16{-32768, 0, 32767} {
				for _, args := range [][5]byte{{0, 0, 0, 0, 0}, {1, 2, 3, 4, 5}, {255, 254, 253, 252, 251}} {
					for _, lifetime := range []int{-1, 0, 1, 255} {
						for _, fail := range []int{0, 1} {
							c.resetCase(env, pix, 31, frame)
							clear(cache)
							c.FailEvery = fail
							legacy.PortTestPresentationBubble(typ.Index(), 123, -456, z, args[0], args[1], args[2], args[3], args[4], lifetime)
							if len(c.Calls) != 1 || c.Calls[0].Type != typ.Index() || c.Calls[0].Position != image.Pt(123, -456) || c.srv.Rand.Logic.Index() != 31 || c.srv.Rand.Other.Index() != 32 {
								t.Fatal("bubble spawn/RNG")
							}
							for i, n := range names[:8] {
								if memmap.Uint32(0x5D4594, 1217512+4*uintptr(i)) != uint32(c.Things.IndByID(n)) {
									t.Fatal("bubble color cache")
								}
							}
							dr := c.Objs.List1
							if (dr == nil) != (fail == 1) {
								t.Fatal("bubble allocation failure")
							}
							var data [15]byte
							var deadline uint32
							if dr != nil {
								copy(data[:], unsafe.Slice((*byte)(unsafe.Pointer(&dr.Union)), 15))
								var want [15]byte
								binary.LittleEndian.PutUint32(want[:], uint32(noxcolor.RGB5551Color(0x21, 0x83, 0xe7).Color32()))
								rgb := palette[kind]
								binary.LittleEndian.PutUint32(want[4:], uint32(noxcolor.RGB5551Color(rgb[0], rgb[1], rgb[2]).Color32()))
								copy(want[8:], []byte{args[0], 1, args[1], args[1], args[3], args[4], args[2]})
								deadline = dr.Deadline
								if data != want || dr.ZVal != uint16(z) || deadline != frame+uint32(lifetime) {
									t.Fatal("bubble color/parameters/lifetime")
								}
								if c.Objs.List6 != dr || c.Objs.List4 != dr || c.Objs.DeadlineList != dr || dr.ObjFlags&0x600000 != 0x600000 {
									t.Fatal("bubble list ownership")
								}
							}
							rows = append(rows, record{kind, frame, z, args, lifetime, fail, data, deadline})
						}
					}
				}
			}
		}
	}
	spellbookCapture(t, "client-presentation-bubbles", rows, "4425d8e44275b3bcaf2929870ee96f2af55e83dd291943381a26641b974df65b")
}
