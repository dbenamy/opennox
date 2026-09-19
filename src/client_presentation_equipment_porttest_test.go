//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"testing"
)

func TestClientPresentationEquipment(t *testing.T) {
	o := newObjectRenderOwner(t)
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
	for _, op := range []byte{80, 81, 82, 83, 255} {
		weapon := op == 80 || op == 81
		capacity := len(npc.Armor)
		if weapon {
			capacity = len(npc.Weapon)
		}
		for _, used := range []int{0, 1, capacity - 1, capacity, capacity + 1} {
			for _, mask := range []uint32{0, 1, 0x80000000, 0xffffffff} {
				for _, mods := range [][4]byte{{1, 2, 3, 4}, {255, 0, 99, 1}, {4, 4, 4, 4}} {
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
					legacy.PortTestPresentationEquip(op, 7, mask, &mods)
					name := fmt.Sprintf("op=%d/used=%d/mask=%08x/mods=%v", op, used, mask, mods)
					if *npc != want {
						t.Fatalf("%s: insertion, capacity, fields or modifier identity", name)
					}
					rows = append(rows, record{name, npc.WeaponEquip, npc.ArmorEquip, normalize(npc.Weapon[:]), normalize(npc.Armor[:])})
				}
			}
		}
		before := *npc
		legacy.PortTestPresentationEquip(op, 999, 1, &[4]byte{1, 2, 3, 4})
		if *npc != before {
			t.Fatal("missing NPC mutated an existing record")
		}
	}
	spellbookCapture(t, "client-presentation-equipment", rows, "754bb1b502147ebfd3e05e75fd959b94c367ee9cc9b8c25012f11d9edf843489")
}
