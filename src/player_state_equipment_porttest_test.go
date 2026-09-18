//go:build porttest

package opennox

import (
	"fmt"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"reflect"
	"testing"
	"unsafe"
)

type playerStateEquipmentOwner struct {
	*matchRosterOwner
	modifiers map[unsafe.Pointer]uint32
}

func newPlayerStateEquipmentOwner(t *testing.T, missingColors ...bool) *playerStateEquipmentOwner {
	o := &playerStateEquipmentOwner{matchRosterOwner: newMatchRosterOwner(t), modifiers: map[unsafe.Pointer]uint32{nil: 0}}
	var mods []*server.ModifierEff
	var names []*byte
	for i := 0; i < 16; i++ {
		m, free := alloc.New(server.ModifierEff{})
		t.Cleanup(free)
		name := fmt.Sprintf("Color%d", i)
		switch i {
		case 0:
			name = "UserColor1"
			if len(missingColors) > 0 && missingColors[0] {
				name = "MissingColor"
			}
		case 9:
			name = "ArmorQuality1"
		case 10:
			name = "Material1"
		case 11:
			name = "Replenishment1"
		}
		n, free := alloc.CString(name)
		t.Cleanup(free)
		mods = append(mods, m)
		names = append(names, n)
		o.modifiers[unsafe.Pointer(m)] = uint32(i + 1)
	}
	t.Cleanup(o.s.PortTestControlsModifiers(mods, names))
	return o
}

type playerStateEquipmentSlot struct {
	Mask uint32
	Mods [4]uint32
	Tail uint32
}
type playerStateEquipmentSnapshot struct {
	WeaponMask, ArmorMask uint32
	Weapons, Armor        []playerStateEquipmentSlot
}

func (o *playerStateEquipmentOwner) snapshot(t *testing.T, pl *server.Player) playerStateEquipmentSnapshot {
	t.Helper()
	out := playerStateEquipmentSnapshot{WeaponMask: pl.WeaponEquip, ArmorMask: pl.ArmorEquip}
	slots := func(in []server.EquipmentData) []playerStateEquipmentSlot {
		var out []playerStateEquipmentSlot
		for _, s := range in {
			r := playerStateEquipmentSlot{Mask: s.Field0, Tail: s.Field20}
			for i, p := range s.Field4 {
				id, ok := o.modifiers[p]
				if !ok {
					t.Fatal("unknown modifier pointer")
				}
				r.Mods[i] = id
			}
			out = append(out, r)
		}
		return out
	}
	out.Weapons = slots(pl.Weapon[:])
	out.Armor = slots(pl.Armor[:])
	return out
}
func TestPlayerStateEquipment(t *testing.T) {
	o := newPlayerStateEquipmentOwner(t)
	pl := o.units[0].UpdateDataPlayer().Player
	pl.NetCodeVal = 0x12345678
	type row struct {
		Name  string
		State playerStateEquipmentSnapshot
	}
	var rows []row
	for _, cmd := range []byte{0, 79, 80, 81, 82, 255} {
		for _, fill := range []int{0, 1, 25, 26, 27} {
			for _, mask := range []uint32{0, 1, 0x80000000, 0xffffffff} {
				for _, known := range []bool{false, true} {
					name := fmt.Sprintf("cmd%d/fill%d/mask%x/known%t", cmd, fill, mask, known)
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
						code := uint32(0x87654321)
						if known {
							code = pl.NetCodeVal
						}
						mods := [4]byte{1, 16, 17, 255}
						legacy.PortTestPlayerStateEquipment("equip", nil, cmd, code, mask, &mods)
						if known {
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
	spellbookCapture(t, "player-state-equipment", rows, "f6ddedbaaff8eb114dcdb94cea890753c8d9004d40e065a04733e05083e40655")
}
func TestPlayerStateUnequip(t *testing.T) {
	o := newPlayerStateEquipmentOwner(t)
	pl := o.units[0].UpdateDataPlayer().Player
	pl.NetCodeVal = 123
	type row struct {
		Name  string
		State playerStateEquipmentSnapshot
	}
	var rows []row
	for _, cmd := range []byte{0, 83, 84, 255} {
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
					legacy.PortTestPlayerStateEquipment("unequip", nil, cmd, code, mask, nil)
					if known {
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
					rows = append(rows, row{fmt.Sprintf("cmd%d/index%d/mask%x/known%t", cmd, index, mask, known), got})
				}
			}
		}
	}
	spellbookCapture(t, "player-state-unequip", rows, "8936d5370316146f30b40954bafae5aa0a659baeafdca3fc537dfd33878ffdcf")
}

func TestPlayerStateRespawn(t *testing.T) {
	o := newPlayerStateEquipmentOwner(t)
	pl := o.units[0].UpdateDataPlayer().Player
	pl.NetCodeVal = 123
	type row struct {
		Name  string
		State playerStateEquipmentSnapshot
	}
	var rows []row
	legacy.PortTestPlayerStateEquipment("respawn", nil, 255, 0, 0, nil)
	for _, host := range []bool{false, true} {
		for _, mode := range []uint32{0, 2048, 4096, 6144} {
			for _, class := range []byte{0, 1, 2, 3} {
				for _, mask := range []byte{0, 1, 2, 4, 8, 16, 32, 64, 128, 255} {
					name := fmt.Sprintf("host%t/mode%x/class%d/mask%x", host, mode, class, mask)
					t.Run(name, func(t *testing.T) {
						flags := mode
						if host {
							flags |= 1
						}
						defer noxflags.PortTestGameFlags(noxflags.GameFlag(flags))()
						pl.WeaponEquip = 0x40000000
						pl.ArmorEquip = 0x20000000
						for i := range pl.Weapon {
							pl.Weapon[i] = server.EquipmentData{Field0: 0x123, Field4: [4]unsafe.Pointer{unsafe.Pointer(o.s.Modif.Nox_xxx_modifGetDescById413330(1)), nil, nil, nil}, Field20: uint32(100 + i)}
						}
						for i := range pl.Armor {
							pl.Armor[i] = server.EquipmentData{Field0: 0x456, Field4: [4]unsafe.Pointer{nil, unsafe.Pointer(o.s.Modif.Nox_xxx_modifGetDescById413330(16)), nil, nil}, Field20: uint32(200 + i)}
						}
						*(*byte)(unsafe.Add(pl.C(), 2251)) = class
						for i := 0; i < 5; i++ {
							*(*byte)(unsafe.Add(pl.C(), 2268+i)) = byte(i + 1)
						}
						want := o.snapshot(t, pl)
						if !host {
							want.WeaponMask = 0
							want.ArmorMask = 0
						}
						for i := range want.Weapons {
							want.Weapons[i].Mask = 0
							want.Weapons[i].Mods = [4]uint32{}
						}
						for i := range want.Armor {
							want.Armor[i].Mask = 0
							want.Armor[i].Mods = [4]uint32{}
						}
						wi, ai := 0, 0
						armor := func(bit uint32, mods [4]uint32) {
							want.Armor[ai].Mask = bit
							want.Armor[ai].Mods = mods
							ai++
							want.ArmorMask |= bit
						}
						weapon := func(bit uint32, mods [4]uint32) {
							want.Weapons[wi].Mask = bit
							want.Weapons[wi].Mods = mods
							wi++
							want.WeaponMask |= bit
						}
						if (class != 0 || mode&2048 != 0) && mask&1 != 0 {
							armor(1024, [4]uint32{0, 3, 4, 0})
						}
						if mask&2 != 0 {
							armor(4, [4]uint32{0, 2, 0, 0})
						}
						if mask&4 != 0 {
							armor(1, [4]uint32{6, 5, 0, 0})
						}
						switch class {
						case 0:
							if mode&2048 != 0 {
								if mask&32 != 0 {
									weapon(256, [4]uint32{11, 0, 0, 0})
								}
							} else if mode&4096 != 0 {
								weapon(256, [4]uint32{})
							} else {
								if mask&64 != 0 {
									weapon(512, [4]uint32{})
								}
								if mask&128 != 0 {
									armor(0x1000000, [4]uint32{})
								}
							}
						case 1:
							if mode&2048 != 0 {
								if mask&8 != 0 {
									weapon(0x8000, [4]uint32{10, 0, 0, 0})
								}
							} else if mode&4096 != 0 {
								weapon(0x10000, [4]uint32{0, 0, 12, 0})
							} else if mask&16 != 0 {
								armor(0x4000, [4]uint32{})
							}
						case 2:
							if mode&2048 != 0 {
								if mask&8 != 0 {
									weapon(0x8000, [4]uint32{10, 0, 0, 0})
								}
							} else if mode&4096 != 0 {
								weapon(4, [4]uint32{})
							}
						}
						legacy.PortTestPlayerStateEquipment("respawn", pl, mask, 0, 0, nil)
						got := o.snapshot(t, pl)
						if !reflect.DeepEqual(got, want) {
							t.Fatal("respawn slots", got, want)
						}
						rows = append(rows, row{name, got})
					})
				}
			}
		}
	}
	spellbookCapture(t, "player-state-respawn", rows, "400939b1312e2cfe1badda6a3d6cffb96772c3b4f72161ae28bc76702dc0063b")
}

func TestPlayerStateRespawnMissingColors(t *testing.T) {
	o := newPlayerStateEquipmentOwner(t, true)
	pl := o.units[0].UpdateDataPlayer().Player
	type row struct {
		Host  bool
		State playerStateEquipmentSnapshot
	}
	var rows []row
	for _, host := range []bool{false, true} {
		restore := noxflags.PortTestGameFlags(noxflags.GameFlag(bool2int(host)))
		pl.WeaponEquip = 0x1234
		pl.ArmorEquip = 0x5678
		for i := range pl.Weapon {
			pl.Weapon[i] = server.EquipmentData{Field0: 1, Field20: 99}
		}
		for i := range pl.Armor {
			pl.Armor[i] = server.EquipmentData{Field0: 2, Field20: 88}
		}
		legacy.PortTestPlayerStateEquipment("respawn", pl, 255, 0, 0, nil)
		got := o.snapshot(t, pl)
		if got.WeaponMask != uint32(bool2int(host))*0x1234 || got.ArmorMask != uint32(bool2int(host))*0x5678 {
			t.Fatal("mask reset before missing color lookup")
		}
		for _, slot := range got.Weapons {
			if slot.Mask != 0 || slot.Mods != [4]uint32{} || slot.Tail != 99 {
				t.Fatal("weapon reset", slot)
			}
		}
		for _, slot := range got.Armor {
			if slot.Mask != 0 || slot.Mods != [4]uint32{} || slot.Tail != 88 {
				t.Fatal("armor reset", slot)
			}
		}
		rows = append(rows, row{host, got})
		restore()
	}
	spellbookCapture(t, "player-state-respawn-missing-colors", rows, "08c8df8ca6b3d107768611375df4807478f3b115d85ce852641b3bc02cf11992")
}
