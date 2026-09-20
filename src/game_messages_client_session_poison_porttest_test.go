//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"testing"
)

func TestGameMessageClientSessionPoison(t *testing.T) {
	o := newObjectRenderOwner(t)
	t.Cleanup(o.c.srv.PortTestScreenNPCs(2))
	a, b := o.c.srv.NPCs.New(7), o.c.srv.NPCs.New(999)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	allies := combatAllyStorage(t)
	type row struct {
		On, Mode, Ally, Marked, Value int
		NPC                           uint32
		Allies                        []byte
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for mode := 0; mode < 3; mode++ {
			for ally := 0; ally < 2; ally++ {
				for marked := 0; marked < 2; marked++ {
					for value := 0; value < 256; value++ {
						binary.LittleEndian.PutUint32(connected, uint32(on))
						*a = server.NPC{LiveVal: 1, IDVal: 7, Field1312: 0xcafe1234}
						if mode == 1 {
							a.LiveVal = 0
						}
						if mode == 2 {
							a.IDVal = 8
						}
						*b = server.NPC{LiveVal: 1, IDVal: 999, Field1312: 0xdeadbeef}
						want, other := *a, *b
						legacy.PortTestCombatAllyClear()
						legacy.PortTestCombatAllyAdd(999, 123, 456)
						if ally != 0 {
							legacy.PortTestCombatAllyAdd(7, 234, 567)
						}
						if on != 0 {
							legacy.PortTestCombatAllyFlag(7, byte(value))
							if mode == 0 {
								want.Field1312 = uint32(value)
							}
						}
						wantAllies := bytes.Clone(allies)
						legacy.PortTestCombatAllyClear()
						legacy.PortTestCombatAllyAdd(999, 123, 456)
						if ally != 0 {
							legacy.PortTestCombatAllyAdd(7, 234, 567)
						}
						code := uint16(7 | marked<<15)
						data := []byte{218, byte(code), byte(code >> 8), byte(value)}
						input := bytes.Clone(data)
						n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(218), data)
						if n != 4 || !bytes.Equal(input, data) || *a != want || *b != other || !bytes.Equal(allies, wantAllies) {
							t.Fatal("poison routing/isolation", on, mode, ally, marked, value)
						}
						rows = append(rows, row{on, mode, ally, marked, value, a.Field1312, bytes.Clone(allies)})
					}
				}
			}
		}
	}
	interactionCapture(t, "game-client-session-poison", rows)
}
