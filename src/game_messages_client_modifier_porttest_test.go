//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"image"
	"testing"
	"unsafe"

	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/legacy"
)

func TestGameMessageClientModifierReport(t *testing.T) {
	o := newObjectDrawingOwner(t)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	names := make([]*byte, len(o.mods))
	for i, m := range o.mods {
		names[i] = *(**byte)(m.C())
	}
	t.Cleanup(o.c.srv.Server.PortTestControlsModifiers(o.mods, names))
	type row struct {
		On, Present, Static, WireStatic, Value, Return int
		Words                                          [6]uint32
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for present := 0; present < 2; present++ {
			for static := 0; static < 2; static++ {
				for wireStatic := 0; wireStatic < 2; wireStatic++ {
					for value := 0; value < 256; value++ {
						o.reset(1, 100)
						binary.LittleEndian.PutUint32(connected, uint32(on))
						dr := o.c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(300, 400))
						if dr == nil {
							t.Fatal("drawable allocation")
						}
						dr.NetCode32 = 17
						dr.ObjClass = 0
						if static != 0 {
							dr.ObjClass = object.Class(0x20400000)
						}
						for i := 0; i < 4; i++ {
							*txword(dr, 432+uintptr(4*i)) = uint32(uintptr(o.mods[i].C()))
						}
						*txword(dr, 448) = 0x12345678
						*txword(dr, 452) = 0x98765432
						raw := unsafe.Slice((*byte)(dr.C()), int(unsafe.Sizeof(*dr)))
						want := bytes.Clone(raw)
						code := uint16(17)
						if present == 0 {
							code = 18
						}
						if wireStatic != 0 {
							code |= 0x8000
						}
						data := []byte{103, byte(code), byte(code >> 8), byte(value), byte(value + 1), byte(value + 2), byte(value + 3)}
						input := bytes.Clone(data)
						if on != 0 && present != 0 && static == wireStatic {
							for i, id := range data[3:] {
								var p uint32
								if id >= 1 && id <= 4 {
									p = uint32(uintptr(o.mods[id-1].C()))
								}
								binary.LittleEndian.PutUint32(want[432+4*i:], p)
							}
							binary.LittleEndian.PutUint32(want[448:], 0xffffffff)
						}
						n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(103), data)
						if n != 7 || !bytes.Equal(data, input) || !bytes.Equal(raw, want) {
							t.Fatalf("modifier on%d present%d static%d wire%d value%d", on, present, static, wireStatic, value)
						}
						r := row{On: on, Present: present, Static: static, WireStatic: wireStatic, Value: value, Return: n}
						for i := range r.Words {
							v := *txword(dr, 432+uintptr(i*4))
							if i < 4 {
								if ref, ok := o.modRefs[v]; ok {
									v = ref
								}
							}
							r.Words[i] = v
						}
						rows = append(rows, r)
					}
				}
			}
		}
	}
	interactionCapture(t, "game-progress-modifier-report", rows)
}
