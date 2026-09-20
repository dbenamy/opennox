//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"reflect"
	"testing"
	"unsafe"
)

// Reuse the qualified equipment owners and their independent capacity, mask,
// modifier and canary contracts at the original message-dispatch boundary.
func gameMessageClientEquipmentCall(t *testing.T, on int, op byte, code uint16, mask uint32, mods [4]byte) {
	t.Helper()
	data := []byte{op, byte(code), byte(code >> 8)}
	data = binary.LittleEndian.AppendUint32(data, mask)
	if op == 81 || op == 82 {
		data = append(data, mods[:]...)
	}
	want := bytes.Clone(data)
	if on != 0 && op >= 79 && op <= 82 {
		binary.LittleEndian.PutUint16(want[1:], code&0x7fff)
	}
	n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(op), data)
	if n != len(data) || !bytes.Equal(data, want) {
		t.Fatalf("equipment op%d on%d code%x return%d input%v want%v", op, on, code, n, data, want)
	}
}

func TestGameMessageClientNPCEquipment(t *testing.T) {
	o := newObjectRenderOwner(t)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	t.Cleanup(o.c.srv.PortTestScreenNPCs(2))
	npc := o.c.srv.NPCs.New(7)
	var names []*byte
	for i := range o.mods {
		p, free := alloc.CString(fmt.Sprintf("PresentationModifier%d", i))
		t.Cleanup(free)
		names = append(names, p)
	}
	t.Cleanup(o.c.srv.PortTestControlsModifiers(o.mods, names))
	type record struct {
		Name                  string
		WeaponMask, ArmorMask uint32
		Weapon, Armor         [][6]uint32
	}
	var rows []record
	normalize := func(items []server.EquipmentData) [][6]uint32 {
		out := make([][6]uint32, len(items))
		for i, v := range items {
			out[i][0], out[i][5] = v.Field0, v.Field20
			for j, p := range v.Field4 {
				if p == nil {
					continue
				}
				found := false
				for k, m := range o.mods {
					if p == m.C() {
						out[i][j+1] = uint32(k + 1)
						found = true
						break
					}
				}
				if !found {
					t.Fatal("unowned equipment modifier")
				}
			}
		}
		return out
	}
	for on := 0; on < 2; on++ {
		binary.LittleEndian.PutUint32(connected, uint32(on))
		for _, op := range []byte{79, 80, 81, 82} {
			weapon := op == 80 || op == 81
			capacity := len(npc.Armor)
			if weapon {
				capacity = len(npc.Weapon)
			}
			for _, used := range []int{0, 1, capacity - 1, capacity, capacity + 1} {
				for _, mask := range []uint32{0, 1, 0x80000000, 0xffffffff} {
					for _, mods := range [][4]byte{{1, 2, 3, 4}, {255, 0, 99, 1}, {4, 4, 4, 4}} {
						if op == 79 || op == 80 {
							mods = [4]byte{255, 255, 255, 255}
						}
						o.c.srv.NPCs.Set(npc, 7)
						npc.WeaponEquip, npc.ArmorEquip = 0x40000000, 0x20000000
						slots := npc.Armor[:]
						if weapon {
							slots = npc.Weapon[:]
						}
						for i := range slots {
							slots[i].Field20 = uint32(0xcafe0000) + uint32(i)
							if i < used {
								slots[i].Field0 = 1
							}
						}
						if used > capacity {
							slots[capacity/2].Field0 = 0
						}
						want := *npc
						target := want.Armor[:]
						aggregate := &want.ArmorEquip
						if weapon {
							target = want.Weapon[:]
							aggregate = &want.WeaponEquip
						}
						if on != 0 {
							for i := range target {
								if target[i].Field0 != 0 {
									continue
								}
								target[i].Field0 = mask
								*aggregate |= mask
								for j, id := range mods {
									if id >= 1 && id <= 4 {
										target[i].Field4[j] = o.mods[id-1].C()
									} else {
										target[i].Field4[j] = nil
									}
								}
								break
							}
						}
						gameMessageClientEquipmentCall(t, on, op, 7, mask, mods)
						name := fmt.Sprintf("on=%d/op=%d/used=%d/mask=%08x/mods=%v", on, op, used, mask, mods)
						if *npc != want {
							t.Fatalf("%s: insertion, capacity, fields or modifier identity", name)
						}
						rows = append(rows, record{name, npc.WeaponEquip, npc.ArmorEquip, normalize(npc.Weapon[:]), normalize(npc.Armor[:])})
					}
				}
			}
			before := *npc
			gameMessageClientEquipmentCall(t, on, op, 999, 1, [4]byte{1, 2, 3, 4})
			if *npc != before {
				t.Fatal("missing NPC mutated an existing record")
			}
		}
	}
	gameMessageCapture(t, "game-client-npc-equipment", rows)
}

func TestGameMessageClientPlayerEquipment(t *testing.T) {
	o := newPlayerStateEquipmentOwner(t)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	pl := o.units[0].UpdateDataPlayer().Player
	pl.NetCodeVal = 7
	type row struct {
		Name  string
		State playerStateEquipmentSnapshot
	}
	var rows []row
	for on := 0; on < 2; on++ {
		binary.LittleEndian.PutUint32(connected, uint32(on))
		for _, cmd := range []byte{79, 80, 81, 82} {
			for _, fill := range []int{0, 1, 25, 26, 27} {
				for _, mask := range []uint32{0, 1, 0x80000000, 0xffffffff} {
					for _, known := range []bool{false, true} {
						name := fmt.Sprintf("on%d/cmd%d/fill%d/mask%x/known%t", on, cmd, fill, mask, known)
						t.Run(name, func(t *testing.T) {
							pl.WeaponEquip = 0x10
							pl.ArmorEquip = 0x20
							for i := range pl.Weapon {
								pl.Weapon[i] = server.EquipmentData{Field4: [4]unsafe.Pointer{unsafe.Pointer(o.s.Modif.Nox_xxx_modifGetDescById413330(1)), nil, nil, nil}, Field20: uint32(100 + i)}
								if i < fill {
									pl.Weapon[i].Field0 = 0x100
								}
							}
							for i := range pl.Armor {
								pl.Armor[i] = server.EquipmentData{Field4: [4]unsafe.Pointer{nil, unsafe.Pointer(o.s.Modif.Nox_xxx_modifGetDescById413330(16)), nil, nil}, Field20: uint32(200 + i)}
								if i < fill {
									pl.Armor[i].Field0 = 0x200
								}
							}
							want := o.snapshot(t, pl)
							code := uint32(999)
							if known {
								code = pl.NetCodeVal
							}
							mods := [4]byte{1, 16, 17, 255}
							if cmd == 79 || cmd == 80 {
								mods = [4]byte{255, 255, 255, 255}
							}
							gameMessageClientEquipmentCall(t, on, cmd, uint16(code)|0x8000, mask, mods)
							if known && on != 0 {
								slots := want.Armor
								if cmd == 80 || cmd == 81 {
									slots = want.Weapons
									want.WeaponMask |= mask
								} else {
									want.ArmorMask |= mask
								}
								if fill < len(slots) {
									slots[fill].Mask = mask
									slots[fill].Mods = [4]uint32{1, 16, 0, 0}
									if cmd == 79 || cmd == 80 {
										slots[fill].Mods = [4]uint32{}
									}
								}
							}
							got := o.snapshot(t, pl)
							if !reflect.DeepEqual(got, want) {
								t.Fatal("equipment slot/mask mutation", got, want)
							}
							rows = append(rows, row{name, got})
						})
					}
				}
			}
		}
	}
	gameMessageCapture(t, "game-client-player-equipment", rows)
}

func TestGameMessageClientPlayerUnequip(t *testing.T) {
	o := newPlayerStateEquipmentOwner(t)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	pl := o.units[0].UpdateDataPlayer().Player
	pl.NetCodeVal = 123
	type row struct {
		Name  string
		State playerStateEquipmentSnapshot
	}
	var rows []row
	for on := 0; on < 2; on++ {
		binary.LittleEndian.PutUint32(connected, uint32(on))
		for high := 0; high < 2; high++ {
			for _, cmd := range []byte{83, 84} {
				for _, index := range []int{0, 13, 25, 26, 27} {
					for _, mask := range []uint32{0, 1, 0x80000000, 0xffffffff} {
						for _, known := range []bool{false, true} {
							pl.WeaponEquip = 0xffffffff
							pl.ArmorEquip = 0xffffffff
							for i := range pl.Weapon {
								pl.Weapon[i] = server.EquipmentData{Field0: 0x42, Field4: [4]unsafe.Pointer{unsafe.Pointer(o.s.Modif.Nox_xxx_modifGetDescById413330(1)), nil, nil, nil}, Field20: uint32(i)}
								if i == index {
									pl.Weapon[i].Field0 = mask
								}
							}
							for i := range pl.Armor {
								pl.Armor[i] = server.EquipmentData{Field0: 0x42, Field4: [4]unsafe.Pointer{unsafe.Pointer(o.s.Modif.Nox_xxx_modifGetDescById413330(1)), nil, nil, nil}, Field20: uint32(i)}
								if i == index {
									pl.Armor[i].Field0 = mask
								}
							}
							want := o.snapshot(t, pl)
							code := uint32(456)
							if known {
								code = 123
							}
							gameMessageClientEquipmentCall(t, on, cmd, uint16(code)|uint16(high<<15), mask, [4]byte{})
							if known && on != 0 && high == 0 {
								slots := want.Armor
								if cmd == 84 {
									slots = want.Weapons
									want.WeaponMask &^= mask
								} else {
									want.ArmorMask &^= mask
								}
								if index < len(slots) {
									slots[index].Mask = 0
								}
							}
							got := o.snapshot(t, pl)
							if !reflect.DeepEqual(got, want) {
								t.Fatal("unequip mutation", cmd, index, mask, known)
							}
							rows = append(rows, row{fmt.Sprintf("on%d/high%d/cmd%d/index%d/mask%x/known%t", on, high, cmd, index, mask, known), got})
						}
					}
				}
			}
		}
	}
	gameMessageCapture(t, "game-client-player-unequip", rows)
}
