//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

func TestGameMessageClientNPCAppearance(t *testing.T) {
	o := newObjectRenderOwner(t)
	t.Cleanup(o.c.srv.PortTestScreenNPCs(2))
	a, b := o.c.srv.NPCs.New(7), o.c.srv.NPCs.New(999)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	type row struct {
		On, Mode, High, ID, Value, Return, Lookup int
		Live                                      byte
		Colors                                    [6]uint32
		Kind, WeaponMask, ArmorMask               uint32
		Input                                     []byte
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for mode := 0; mode < 3; mode++ {
			for high := 0; high < 2; high++ {
				for _, id := range []uint16{0, 7, 32767} {
					for value := 0; value < 256; value++ {
						binary.LittleEndian.PutUint32(connected, uint32(on))
						*a = server.NPC{IDVal: int32(id), WeaponEquip: 0x12345678, ArmorEquip: 0x87654321, Field1312: 99}
						for i := range a.Color8 {
							a.Color8[i] = 0x13572468
						}
						a.Weapon[0].Field0 = 17
						a.Weapon[0].Field20 = 99
						a.Armor[0].Field0 = 23
						if mode != 0 {
							a.LiveVal = 1
						}
						if mode == 2 {
							a.IDVal = 2222
						}
						*b = server.NPC{LiveVal: 1, IDVal: 3333, WeaponEquip: 0xabcdef01}
						want, other := *a, *b
						code := id | uint16(high<<15)
						data := []byte{105, byte(code), byte(code >> 8)}
						for i := 0; i < 18; i++ {
							data = append(data, byte(value+i*37))
						}
						wantInput := bytes.Clone(data)
						if on != 0 {
							binary.LittleEndian.PutUint16(wantInput[1:], id)
							if mode != 2 {
								want = server.NPC{LiveVal: 1, IDVal: int32(id), Field1312: uint32(high)}
								for i := range want.Color8 {
									r, g, b := uint32(data[3+i*3]>>3), uint32(data[4+i*3]>>3), uint32(data[5+i*3]>>3)
									rgb := r<<10 | g<<5 | b
									want.Color8[i] = rgb | rgb<<16
								}
							}
						}
						n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(105), data)
						if n != 21 || !bytes.Equal(data, wantInput) || *a != want || *b != other {
							t.Fatalf("NPC on%d mode%d high%d id%d value%d return%d", on, mode, high, id, value, n)
						}
						found := o.c.srv.NPCs.ByID(int(id))
						lookup := -1
						if want.LiveVal != 0 && want.IDVal == int32(id) {
							if found != a {
								t.Fatal("NPC lookup identity")
							}
							lookup = 0
						} else if found != nil {
							t.Fatal("unexpected NPC allocation")
						}
						rows = append(rows, row{on, mode, high, int(id), value, n, lookup, a.LiveVal, a.Color8, a.Field1312, a.WeaponEquip, a.ArmorEquip, data})
					}
				}
			}
		}
	}
	gameMessageCapture(t, "game-client-npc-appearance", rows)
}
