//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestGameplayReportsInventoryFields(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	type msg struct {
		data               []byte
		recipient, ordered byte
		priority           uint32
	}
	var checks [][]msg
	for _, op := range []int{11, 12, 13, 28, 29, 50} {
		for _, class := range []uint32{0x1000, 0x2000000, 0x100000} {
			for mask := uint32(0); mask < 16; mask++ {
				for _, recipient := range []uint32{1, 31} {
					for _, player := range []bool{false, true} {
						for _, health := range []bool{false, true} {
							s := gameplayReportsBase(op, reportValue(recipient), reportObject(3), reportValue(1))
							p := s.Callbacks.Shop
							a := p.TemporaryUpdates.World.Objectives.Attack
							a.ActorWords = map[int]uint32{36: 123}
							if !player {
								p.TemporaryUpdates.ItemRefs = []map[int]int{{492: 4}}
								p.TemporaryUpdates.ItemWords[1][36] = 123
							}
							p.Equipment.HolderOnly = true
							p.Inventory.WeaponBits[23] = 0x01020304
							p.Inventory.ArmorBits = map[uint16]uint32{23: 0xa0b0c0d0}
							p.Items[0].Class = class
							p.Items[0].Type = 23
							p.Items[0].Health = health
							p.Items[0].HP = 37
							p.Items[0].MaxHP = 100
							p.TemporaryUpdates.ItemWords[0][36] = 321
							a.Controls.Reports.Caches[0] = 23
							mods := []byte{255, 255, 255, 255}
							for j := range mods {
								p.Items[0].Mods[j] = mask&(1<<j) != 0
								if p.Items[0].Mods[j] {
									mods[j] = byte(40 + j)
								}
							}
							var expected []msg
							add := func(data []byte, to, ordered byte, priority uint32) {
								expected = append([]msg{{data, to, ordered, priority}}, expected...)
							}
							equip := func(to byte) {
								if class == 0x100000 {
									return
								}
								opcode := byte(80)
								flags := uint32(0x01020304)
								if class == 0x2000000 {
									opcode = 79
									flags = 0xa0b0c0d0
								}
								if mask != 0 {
									opcode = 81
									if class == 0x2000000 {
										opcode = 82
									}
								}
								id := uint16(123)
								if player {
									id |= 0x8000
								}
								data := []byte{opcode, byte(id), byte(id >> 8), 0, 0, 0, 0}
								binary.LittleEndian.PutUint32(data[3:], flags)
								if mask != 0 {
									data = append(data, mods...)
								}
								add(data, to, 1, 0)
							}
							switch op {
							case 11:
								equip(byte(recipient))
							case 12:
								if class != 0x100000 {
									opcode := byte(84)
									flags := uint32(0x01020304)
									if class == 0x2000000 {
										opcode = 83
										flags = 0xa0b0c0d0
									}
									data := []byte{opcode, 123, 0, 0, 0, 0, 0}
									binary.LittleEndian.PutUint32(data[3:], flags)
									add(data, byte(recipient), 1, 0)
								}
							case 13:
								add([]byte{96, 65, 1}, byte(recipient), 1, 0)
								equip(255)
							case 28, 29:
								data := []byte{75, 65, 1, 23, 0}
								if op == 29 || class != 0x100000 {
									data[0] = 76
									data = append(data, mods...)
								}
								add(data, byte(recipient), 1, 0)
								if health {
									add([]byte{68, 65, 1, 37, 0, 100, 0}, byte(recipient), 1, 0)
								}
							case 50:
								data := append([]byte{103, 65, 1}, mods...)
								add(data, byte(recipient), 0, 1)
							}
							cases = append(cases, s)
							checks = append(checks, expected)
						}
					}
				}
			}
		}
	}
	out := controlsRun(t, cases)
	for i, r := range out {
		ps := r.Callbacks.Shop.Sequence[0].Packets
		want := checks[i]
		if len(ps) != len(want) {
			t.Fatalf("inventory case%d message count%d want%d", i, len(ps), len(want))
		}
		for j, p := range ps {
			c := want[j]
			if !bytes.Equal(p.Data, c.data) || p.Recipient != c.recipient || p.Ordered != c.ordered || p.A4 != 0 || p.A5 != c.priority {
				t.Fatalf("inventory case%d message%d fields or order", i, j)
			}
		}
	}
	gameplayReportsCapture(t, "inventory-fields", out)
}
