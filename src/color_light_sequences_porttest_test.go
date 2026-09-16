//go:build porttest

package opennox

import (
	"encoding/binary"
	"fmt"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"testing"
	"unsafe"
)

func TestColorLightUpdateSequences(t *testing.T) {
	c, pix, effects := newEffectsFullOwner(t)
	env := legacy.PortTestNewClientParticleEnvironment()
	t.Cleanup(env.Restore)
	t.Cleanup(noxflags.PortTestGameFlags(0))
	oldFPS := c.srv.TickRate()
	t.Cleanup(func() { c.srv.SetTickRate(oldFPS) })
	c.srv.SetTickRate(30)
	vp := c.Viewport()
	oldVP := *vp
	t.Cleanup(func() { *vp = oldVP })
	vp.World = image.Rect(0, 0, 1000, 1000)
	vp.Size = image.Pt(1000, 1000)
	type record struct {
		Flags, Mode, Count int
		Present            bool
		Frame              uint32
		Light              []uint32
		Linked             bool
	}
	var rows []record
	for _, flags := range []int{0, 0x40, 0x80, 0xc0, 0x180, 0x1c0} {
		for _, mode := range []int{0, 1} {
			for _, count := range []int{0, 1, 2} {
				for _, present := range []bool{false, true} {
					t.Run(fmt.Sprintf("flags%x-mode%d-count%d-present%v", flags, mode, count, present), func(t *testing.T) {
						c.resetCase(effects, pix, 17, 0)
						env.Reset()
						d := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(200, 200))
						target := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(300, 200))
						target.ObjClass = 0x400000
						target.NetCode32 = 77
						raw := unsafe.Slice((*byte)(d.C()), int(unsafe.Sizeof(*d)))
						clear(raw[136:276])
						clear(raw[432:435])
						binary.LittleEndian.PutUint32(raw[164:], 0x12345678)
						binary.LittleEndian.PutUint32(raw[168:], uint32(mode))
						binary.LittleEndian.PutUint16(raw[176:], uint16(flags|0x15))
						for i := 0; i < 2; i++ {
							raw[178+3*i] = byte(20 + 100*i)
							raw[179+3*i] = byte(30 + 100*i)
							raw[180+3*i] = byte(40 + 100*i)
							raw[226+i] = byte(5 + 30*i)
							raw[242+i] = byte(10 + 60*i)
						}
						for i := 0; i < 3; i++ {
							raw[432+i] = byte(count)
							binary.LittleEndian.PutUint16(raw[258+2*i:], 3)
						}
						code := uint32(88)
						if present {
							code = 77
						}
						binary.LittleEndian.PutUint32(raw[264:], code)
						binary.LittleEndian.PutUint16(raw[268:], 45)
						binary.LittleEndian.PutUint16(raw[270:], 90)
						binary.LittleEndian.PutUint16(raw[272:], 180)
						c.Objs.List5Add(d)
						for _, frame := range []uint32{0, 1, 2, 3, 5, 6, 7, 30} {
							c.srv.SetFrame(frame)
							got := legacy.PortTestColorLight(5, vp, d)
							if got != 1 || (d.InClientUpdateList != 0) != (count != 0) {
								t.Fatal("sequence list/return")
							}
							words := make([]uint32, 10)
							for i := range words {
								words[i] = binary.LittleEndian.Uint32(raw[136+4*i:])
							}
							if count == 0 {
								if words[7] != 0x12345678 || words[8] != uint32(mode) {
									t.Fatal("empty arrays applied direction")
								}
							} else if flags&0x40 != 0 && present {
								if uint16(words[7]) != 0 || words[8] != 0 {
									t.Fatal("target direction must take precedence over rotation")
								}
							}
							if count < 2 && (words[0] != 0 || words[1] != 0 || words[2] != 0 || words[3] != 0) {
								t.Fatal("single/empty arrays interpolated")
							}
							if c.Objs.Count != 2 || c.srv.Rand.Logic.Index() != 17 || c.srv.Rand.Other.Index() != 18 {
								t.Fatal("update ownership/RNG")
							}
							rows = append(rows, record{flags, mode, count, present, frame, words, d.InClientUpdateList != 0})
						}
					})
				}
			}
		}
	}
	spellbookCapture(t, "color-light-sequences", rows, "7583afa0841860903f89a98cd2c6ca99e5b3d3a8992d593ac863df141bb7b300")
}
