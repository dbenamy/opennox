//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestGameplayReportsNPC(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	var checks [][][]byte
	for _, id := range []uint32{0, 123, 32767, 32768, 65535} {
		for _, poison := range []uint32{0, 1, 255} {
			for _, seed := range []uint32{0, 0x01020304, 0xffffffff} {
				for equipped := uint32(0); equipped < 4; equipped++ {
					s := gameplayReportsBase(52, reportValue(7), reportObject(1))
					p := s.Callbacks.Shop
					a := p.TemporaryUpdates.World.Objectives.Attack
					p.Resources.Subject = 3
					p.Equipment.HolderOnly = true
					p.Inventory.Linked = []int{0, 1}
					p.Inventory.Owned = []int{0, 1}
					p.Items[0].Class = 0x1000000
					p.Items[0].Type = 23
					p.Items[0].Mods = [4]bool{}
					p.Items[0].Flags = 0
					p.Items[1].Class = 0x2000000
					p.Items[1].Type = 24
					p.Items[1].Mods = [4]bool{}
					p.Items[1].Flags = 0
					p.Inventory.WeaponBits[23] = 0x100
					p.Inventory.ArmorBits = map[uint16]uint32{24: 4}
					if equipped&1 != 0 {
						p.Items[0].Flags = 0x100
					}
					if equipped&2 != 0 {
						p.Items[1].Flags = 0x100
					}
					a.ActorWords = map[int]uint32{36: id, 540: poison}
					a.UpdateWords = map[int]uint32{}
					color := make([]byte, 20)
					for j := 0; j < 5; j++ {
						v := seed + uint32(j*13)
						a.UpdateWords[2076+j*4] = v
						binary.LittleEndian.PutUint32(color[j*4:], v)
					}
					code := uint16(id)
					if poison != 0 {
						code |= 0x8000
					}
					expected := [][]byte{append([]byte{105, byte(code), byte(code >> 8)}, color[:18]...)}
					if equipped&1 != 0 {
						expected = append([][]byte{{80, byte(id), byte(id >> 8), 0, 1, 0, 0}}, expected...)
					}
					if equipped&2 != 0 {
						expected = append([][]byte{{79, byte(id), byte(id >> 8), 4, 0, 0, 0}}, expected...)
					}
					cases = append(cases, s)
					checks = append(checks, expected)
				}
			}
		}
	}
	out := controlsRun(t, cases)
	for i, r := range out {
		ps := r.Callbacks.Shop.Sequence[0].Packets
		if len(ps) != len(checks[i]) {
			t.Fatalf("NPC case%d count%d want%d", i, len(ps), len(checks[i]))
		}
		for j, p := range ps {
			priority := uint32(0)
			if p.Data[0] == 105 {
				priority = 1
			}
			if !bytes.Equal(p.Data, checks[i][j]) || p.Recipient != 7 || p.Ordered != 1 || p.A4 != 0 || p.A5 != priority {
				t.Fatalf("NPC case%d message%d colors/equipment/routing", i, j)
			}
		}
	}
	gameplayReportsCapture(t, "npc", out)
}
