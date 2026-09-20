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

func TestGameMessageClientObjectFields(t *testing.T) {
	c, pix, env := newEffectsFullOwner(t)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	type row struct {
		On, Present, Static, Special, Kind int
		Value                              uint32
		Return                             int
		Z                                  uint16
		Slave, Previous                    uint32
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for present := 0; present < 2; present++ {
			for static := 0; static < 2; static++ {
				for special := 0; special < 2; special++ {
					for _, kind := range []int{94, 95, 107} {
						values := []uint32{0, 1, 255, 65535, 0x80000000, 0xffffffff}
						if kind != 107 {
							values = nil
							for n := 0; n < 256; n++ {
								values = append(values, uint32(n))
							}
						}
						for _, value := range values {
							c.resetCase(env, pix, 1, 100)
							binary.LittleEndian.PutUint32(connected, uint32(on))
							dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(300, 400))
							if dr == nil {
								t.Fatal("allocation")
							}
							dr.ObjClass = 0
							if static != 0 {
								dr.ObjClass = object.Class(0x20400000)
							}
							if special != 0 {
								dr.ObjClass |= object.Class(2)
							}
							dr.ObjSubClass = 0x40000
							dr.NetCode32 = 17
							dr.AnimInd = 8
							dr.AnimFrameSlave = 91
							dr.Field_78 = 92
							raw := unsafe.Slice((*byte)(unsafe.Pointer(dr)), int(unsafe.Sizeof(*dr)))
							binary.LittleEndian.PutUint16(raw[104:], 0x1234)
							want := bytes.Clone(raw)
							code := uint16(17)
							if present == 0 {
								code = 18
							}
							if static != 0 {
								code |= 0x8000
							}
							data := []byte{byte(kind), byte(code), byte(code >> 8), byte(value)}
							if kind == 107 {
								data = data[:3]
								data = binary.LittleEndian.AppendUint32(data, value)
							}
							before := bytes.Clone(data)
							if on != 0 && present != 0 {
								switch kind {
								case 94:
									binary.LittleEndian.PutUint16(want[104:], uint16(value))
								case 95:
									binary.LittleEndian.PutUint16(want[104:], uint16(-int32(value)))
								case 107:
									if special == 0 {
										binary.LittleEndian.PutUint32(want[308:], value)
										binary.LittleEndian.PutUint32(want[312:], 91)
									}
								}
							}
							n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(kind), data)
							if n != len(data) || !bytes.Equal(data, before) || !bytes.Equal(raw, want) || c.Objs.Count != 1 {
								t.Fatalf("field on%d present%d static%d special%d kind%d value%x ret%d", on, present, static, special, kind, value, n)
							}
							rows = append(rows, row{on, present, static, special, kind, value, n, binary.LittleEndian.Uint16(raw[104:]), dr.AnimFrameSlave, dr.Field_78})
						}
					}
				}
			}
		}
	}
	interactionCapture(t, "game-client-object-fields", rows)
}
