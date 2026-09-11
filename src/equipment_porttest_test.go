//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"math"
	"testing"
)

var equipmentHashes = map[string]string{
	"equipment-switch-sequences":   "e947b54116e5d4ce9fe7473fa900227784c84615cb368622e3ec65c52e30ec64",
	"equipment-modifier-callbacks": "902a43fa179091d807aaeb0576b7130b942a232fa6958db18faecb2a1edfd067",
	"equipment-strength":           "4d06204d75e0a3de6d47c9f4cc7f751e94ce025b32dd834220c21edff6e92827",
	"equipment-armor-values":       "3af26646a01219e74336411f50622800ec999b800bf3ae8e1935cd66e075630d",
	"equipment-admission":          "cbbeeefec69b27582242e426070a2899c538d5d586a0c7141c43209f08784952",
	"equipment-npc-owners":         "a9e9d482a4880f3ea47d36ea69427ad50522a6fdfe030c5f76f01974de2a2fc0",
	"equipment-bow-ammo":           "6d79f8a7b48e8665320b26b58ef6fcdda4e8e828336cd5aefe969e830baf179d",
	"equipment-shield-selection":   "202430f11e7cf430319b27a39f80da4ef9e8d25224fb59a2685d6a12f3d7a613",
	"equipment-inventory-policies": "4700973a030f6826566898f2bdf029e7e33ef838245ffdae82d577d699064e42",
	"equipment-sound-secondary":    "a4a2bc9de8c0b34382e9789a88be4fc281ae601c4da979410951685f8d354677",
	"equipment-drop-table":         "d04a7e73f2a0b105f8b9e1033b3ed83b434af1839b720fc4c7ba8b9eeb25eeb9",
	"equipment-try-owners":         "76f065eeed6e46d8753badbf5034043d137f57a3773a5b0756d4201ffcd06ef6",
	"equipment-null-guards":        "e3daa48abc12be6c6eccfd01962f754f50092737a2f5566775f7fca7b7bb8394",
	"equipment-armor-masks":        "381d410a3e25096f5b907ed4072a9aeee63cab4cac278f7a39efecf5689ff50e",
	"equipment-cold-table":         "35061941d78ed068a0a3649c70dca8cb6d07bae1adce497ba5460f99192678ff",
	"equipment-armor-search":       "4c4dcbf2674fb5d8576f8953ef8adaf1cf98ee027e65ee8ae21626623c11cc8a",
	"equipment-npc-sync":           "c0c08d1f8138e1e4d3be1565a8dfa2c348ebae900e2cf987c889fb6fd8142569",
}

func equipmentBase() legacy.PortTestRoamSpec {
	s := inventoryBase()
	p := s.Callbacks.Shop
	s.Lifecycle.Eligible = true
	p.Equipment = &legacy.PortTestEquipmentSpec{GameEx: 0x1e, Strength: 30, NPCStrength: 30, Threshold: math.Float64bits(1), Definitions: []legacy.PortTestEquipmentDef{{Type: 23, Strength: 20, Coeff: math.Float32bits(.25)}, {Type: 23, Armor: true, Strength: 20, Coeff: math.Float32bits(.25)}}}
	p.Inventory.Linked = []int{0, 1, 2}
	p.Inventory.Owned = []int{0, 1, 2}
	p.Inventory.WeaponBits = map[uint16]uint32{23: 4}
	p.Inventory.ArmorBits = map[uint16]uint32{23: 2}
	for i := range p.Items {
		p.Items[i].Class = 0x1000000
		p.Items[i].Subclass = 4
	}
	return s
}
func TestEquipmentSwitchSequences(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, armor := range []bool{false, true} {
		for subject := 1; subject <= 3; subject++ {
			for flags := 0; flags < 4; flags++ {
				for _, ex := range []uint32{0, 0x1e} {
					s := equipmentBase()
					p := s.Callbacks.Shop
					p.Resources.Subject = subject
					p.Equipment.GameEx = ex
					equip, dequip := legacy.PortTestEquipment53A420, legacy.PortTestEquipment53A140
					if armor {
						equip, dequip = legacy.PortTestEquipment53E650, legacy.PortTestEquipment53E430
						for i := range p.Items {
							p.Items[i].Class = 0x2000000
							p.Items[i].Subclass = 2
						}
					}
					if flags&1 != 0 {
						p.Resources.Subclass = 0x10
					}
					if flags&2 != 0 {
						p.Items[0].Mods[2] = true
						p.Equipment.Effects[2] = legacy.PortTestEquipmentEffect{Engage: true, Disengage: true, Return: 0x12345678}
					}
					p.Sequence = []legacy.PortTestShopAction{{Op: equip, Item: 0, Value: 1, Side: 1}, {Op: equip, Item: 1, Value: 1, Side: 1}, {Op: dequip, Item: 1, Value: 1, Side: 1}, {Op: dequip, Item: 0, Value: 1, Side: 1}}
					specs = append(specs, s)
				}
			}
		}
	}
	callbackHash(t, "equipment-switch-sequences", legacy.PortTestRoam(specs), equipmentHashes["equipment-switch-sequences"])
}
func TestEquipmentModifierCallbacks(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for mask := 0; mask < 16; mask++ {
		for variant := 0; variant < 4; variant++ {
			s := equipmentBase()
			p := s.Callbacks.Shop
			for i := 0; i < 4; i++ {
				p.Items[0].Mods[i] = mask&(1<<i) != 0
				p.Equipment.Effects[i] = legacy.PortTestEquipmentEffect{Engage: variant&1 != 0, Disengage: variant&2 != 0, Return: uint32(0x80000000) + uint32(i)}
			}
			p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestEquipment4F2FF0}, {Op: legacy.PortTestEquipment4F3030}}
			specs = append(specs, s)
		}
	}
	r := legacy.PortTestRoam(specs)
	for i, result := range r {
		sp := specs[i].Callbacks.Shop
		for j, enabled := range []bool{sp.Equipment.Effects[3].Engage, sp.Equipment.Effects[3].Disengage} {
			var want uint32
			if sp.Items[0].Mods[3] {
				want = 50103
				if enabled {
					want = 0x80000003
				}
			}
			if result.Callbacks.Shop.Sequence[j].Return != want {
				t.Fatalf("modifier return case %d step %d", i, j)
			}
		}
	}
	callbackHash(t, "equipment-modifier-callbacks", r, equipmentHashes["equipment-modifier-callbacks"])
}
func TestEquipmentStrength(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for subject := 1; subject <= 3; subject++ {
		for _, strength := range []uint32{0, 19, 20, 30, 0xffffffff} {
			for _, armor := range []bool{false, true} {
				for variant := 0; variant < 4; variant++ {
					s := equipmentBase()
					p := s.Callbacks.Shop
					p.Resources.Subject = subject
					p.Resources.Subclass = 0x10
					p.Equipment.Strength = strength
					p.Equipment.NPCStrength = byte(strength)
					p.Equipment.Cheat = variant&1 != 0
					if variant&2 != 0 {
						p.Equipment.Definitions = nil
					}
					if armor {
						p.Items[0].Class = 0x2000000
					}
					p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestEquipment4F9FD0}, {Op: legacy.PortTestEquipment4F3180}}
					specs = append(specs, s)
				}
			}
		}
	}
	r := legacy.PortTestRoam(specs)
	for i, result := range r {
		sp := specs[i].Callbacks.Shop
		var want uint32
		switch sp.Resources.Subject {
		case 1:
			want = sp.Equipment.Strength
		case 3:
			want = uint32(sp.Equipment.NPCStrength)
		}
		if result.Callbacks.Shop.Sequence[0].Return != want {
			t.Fatalf("strength getter case %d", i)
		}
		allowed := sp.Equipment.Cheat || sp.Resources.Subject == 1 && sp.Equipment.Definitions != nil && int32(sp.Equipment.Strength) >= 20
		if (result.Callbacks.Shop.Sequence[1].Return != 0) != allowed {
			t.Fatalf("strength admission case %d", i)
		}
	}
	callbackHash(t, "equipment-strength", r, equipmentHashes["equipment-strength"])
}
func TestEquipmentArmorValues(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for subject := 1; subject <= 3; subject++ {
		for _, bits := range []uint32{0, 0x80000000, 1, 0x3e800000, 0x3f800000, 0x40000000, 0xbf800000, 0x7f800000, 0x7fc12345} {
			for effect := 0; effect < 3; effect++ {
				s := equipmentBase()
				p := s.Callbacks.Shop
				p.Resources.Subject = subject
				p.Equipment.Definitions[1].Coeff = bits
				for i := range p.Items {
					p.Items[i].Class = 0x2000000
					p.Items[i].Flags = 0x100
				}
				if effect > 0 {
					p.Items[0].Mods[0] = true
				}
				if effect == 2 {
					p.Equipment.Effects[0] = legacy.PortTestEquipmentEffect{Defend: true, Output: math.Float32bits(.75)}
				}
				p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestEquipment415C00}, {Op: legacy.PortTestEquipment53E300}}
				specs = append(specs, s)
			}
		}
	}
	callbackHash(t, "equipment-armor-values", legacy.PortTestRoam(specs), equipmentHashes["equipment-armor-values"])
}

func TestEquipmentAdmission(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, armor := range []bool{false, true} {
		for _, linked := range []bool{false, true} {
			for _, eligible := range []bool{false, true} {
				for _, strength := range []uint32{19, 20} {
					for _, active := range []uint32{0, 2, 4, 6} {
						for _, flags := range []uint32{0, 0x100} {
							s := equipmentBase()
							p := s.Callbacks.Shop
							s.Lifecycle.Eligible = eligible
							p.Equipment.Strength = strength
							p.Equipment.ActiveAbilities = active
							p.Items[0].Flags = flags
							if !linked {
								p.Inventory.Linked = nil
								p.Inventory.Owned = nil
								p.Equipment.HolderOnly = true
							}
							op := legacy.PortTestEquipment53A420
							if armor {
								p.Items[0].Class = 0x2000000
								p.Items[0].Subclass = 2
								op = legacy.PortTestEquipment53E650
							}
							p.Sequence = []legacy.PortTestShopAction{{Op: op, Value: 1, Side: 1}}
							specs = append(specs, s)
						}
					}
				}
			}
		}
	}
	callbackHash(t, "equipment-admission", legacy.PortTestRoam(specs), equipmentHashes["equipment-admission"])
}
func TestEquipmentNPCOwners(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, armor := range []bool{false, true} {
		for _, linked := range []bool{false, true} {
			for _, sub := range []uint32{0, 0x10} {
				for variant := 0; variant < 4; variant++ {
					s := equipmentBase()
					p := s.Callbacks.Shop
					p.Resources.Subject = 3
					p.Resources.Subclass = sub
					if !linked {
						p.Inventory.Linked = nil
						p.Inventory.Owned = nil
					}
					eq, deq := legacy.PortTestEquipment53A2C0, legacy.PortTestEquipment53A030
					if armor {
						eq, deq = legacy.PortTestEquipment53E520, legacy.PortTestEquipment53E3A0
						for i := range p.Items {
							p.Items[i].Class = 0x2000000
							p.Items[i].Subclass = 2
						}
					}
					if variant&1 != 0 {
						p.Items[0].Flags = 0x100
					}
					if variant&2 != 0 {
						p.Items[0].Mods[3] = true
						p.Equipment.Effects[3] = legacy.PortTestEquipmentEffect{Engage: true, Disengage: true, Return: 43}
					}
					p.Sequence = []legacy.PortTestShopAction{{Op: eq, Item: 0}, {Op: eq, Item: 1}, {Op: deq, Item: 1}, {Op: legacy.PortTestEquipment4E4B20, Item: 0, Value: 1}, {Op: legacy.PortTestEquipment4E4B20, Item: 0, Value: 2}}
					specs = append(specs, s)
				}
			}
		}
	}
	callbackHash(t, "equipment-npc-owners", legacy.PortTestRoam(specs), equipmentHashes["equipment-npc-owners"])
}
func TestEquipmentBowAmmo(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for subject := 1; subject <= 3; subject += 2 {
		for variant := 0; variant < 8; variant++ {
			for _, bits := range []uint32{2, 4, 8, 0x80} {
				s := equipmentBase()
				p := s.Callbacks.Shop
				p.Resources.Subject = subject
				p.Resources.Subclass = 0x10
				p.Items[0].Type = 23
				p.Items[0].Subclass = bits
				p.Items[1].Type = 24
				p.Items[1].Subclass = 4
				p.Items[2].Type = 25
				p.Items[2].Subclass = 2
				p.Inventory.WeaponBits = map[uint16]uint32{23: bits, 24: 4, 25: 2}
				if variant&1 != 0 {
					p.Inventory.Linked = []int{0, 2}
					p.Inventory.Owned = []int{0, 2}
					p.Equipment.HolderOnly = true
				}
				if variant&2 != 0 {
					p.Items[1].Flags = 0x100
					p.Equipment.ActiveWeapon = 2
					p.Equipment.WeaponFlags = 4
				}
				if variant&4 != 0 {
					p.Items[2].Flags = 0x100
					p.Equipment.WeaponFlags |= 2
				}
				p.Equipment.Definitions = append(p.Equipment.Definitions, legacy.PortTestEquipmentDef{Type: 24, Strength: 20}, legacy.PortTestEquipmentDef{Type: 25, Strength: 20})
				p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestEquipment53A420}, {Op: legacy.PortTestEquipment53A0F0, Value: 1, Side: 1}, {Op: legacy.PortTestEquipment53A680}, {Op: legacy.PortTestEquipment53A140, Item: 1, Value: 1, Side: 1}}
				specs = append(specs, s)
			}
		}
	}
	callbackHash(t, "equipment-bow-ammo", legacy.PortTestRoam(specs), equipmentHashes["equipment-bow-ammo"])
}
func TestEquipmentShieldSelection(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, ex := range []uint32{0, 2, 0x1e} {
		for _, flags := range []uint32{0, 16, 0x100, 0x110, 0x10000010} {
			for _, secondary := range []int{0, 1} {
				for saved := 0; saved < 3; saved++ {
					s := equipmentBase()
					p := s.Callbacks.Shop
					p.Equipment.GameEx = ex
					p.Items[0].Flags = 0x100
					p.Items[0].Subclass = 0x10
					p.Equipment.ActiveWeapon = 1
					p.Equipment.WeaponFlags = 0x10
					p.Inventory.WeaponBits = map[uint16]uint32{23: 0x10}
					for i := 1; i < 3; i++ {
						p.Items[i].Class = 0x2000000
						p.Items[i].Subclass = 2
						p.Items[i].Type = uint16(23 + i)
						p.Items[i].Flags = flags
					}
					p.Inventory.ArmorBits = map[uint16]uint32{24: 0x1000000, 25: 0x2000000}
					p.Equipment.Definitions = append(p.Equipment.Definitions, legacy.PortTestEquipmentDef{Type: 24, Armor: true, Strength: 20}, legacy.PortTestEquipmentDef{Type: 25, Armor: true, Strength: 20})
					p.Equipment.SecondaryWeapon = secondary
					if saved != 0 {
						p.Equipment.SavedShield = saved + 1
					}
					p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestEquipment980523}, {Op: legacy.PortTestEquipment9805EB}, {Op: legacy.PortTestEquipment53A140, Value: 1, Side: 1}, {Op: legacy.PortTestEquipment53A3D0}, {Op: legacy.PortTestEquipment53E600}}
					specs = append(specs, s)
				}
			}
		}
	}
	callbackHash(t, "equipment-shield-selection", legacy.PortTestRoam(specs), equipmentHashes["equipment-shield-selection"])
}
func TestEquipmentInventoryPolicies(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for n := 0; n <= 3; n++ {
		for _, typ := range []int{0, 23, 24, -1} {
			for variant := 0; variant < 8; variant++ {
				s := equipmentBase()
				p := s.Callbacks.Shop
				p.Inventory.Linked = nil
				p.Inventory.Owned = nil
				for i := 0; i < n; i++ {
					p.Inventory.Linked = append(p.Inventory.Linked, i)
					p.Inventory.Owned = append(p.Inventory.Owned, i)
				}
				p.Items[1].Type = 24
				if variant&1 != 0 {
					p.Items[0].Flags = 0x20
				}
				if variant&2 != 0 {
					p.Items[2].Mods[0] = true
				}
				if variant&4 != 0 {
					p.Inventory.NilItem = true
				}
				p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestEquipment4E7D30, Value: uint32(typ)}, {Op: legacy.PortTestEquipment4E7EC0, Item: 2}}
				specs = append(specs, s)
			}
		}
	}
	callbackHash(t, "equipment-inventory-policies", legacy.PortTestRoam(specs), equipmentHashes["equipment-inventory-policies"])
}
func TestEquipmentSoundAndSecondary(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, class := range []uint32{8, 0x1000, 0x1000000, 0x2000000} {
		for _, material := range []uint16{0, 1, 2, 4, 8, 16, 31} {
			for _, sub := range []uint32{0, 2, 0x20, 0x40, 0x4000000} {
				s := equipmentBase()
				p := s.Callbacks.Shop
				p.Items[0].Class = class
				p.Items[0].Subclass = sub
				p.Inventory.Materials = []uint16{material}
				p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestEquipment53A6C0}, {Op: legacy.PortTestEquipment53AAB0}, {Op: legacy.PortTestEquipment53EAE0}, {Op: legacy.PortTestEquipment53E2D0}, {Op: legacy.PortTestEquipment53AB90}}
				specs = append(specs, s)
			}
		}
	}
	callbackHash(t, "equipment-sound-secondary", legacy.PortTestRoam(specs), equipmentHashes["equipment-sound-secondary"])
}
func TestEquipmentDropTable(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, names := range [][]string{{}, {"Sword"}, {"WoodenShield", "Sword", "SteelShield"}, {"missing-equipment-type"}} {
		for _, flags := range []uint32{0, 1, 2, 3, 0xffffffff} {
			s := equipmentBase()
			p := s.Callbacks.Shop
			p.Equipment.TableNames = names
			for range names {
				p.Equipment.TableFlags = append(p.Equipment.TableFlags, flags)
			}
			p.Inventory.ItemTypes = []string{"Sword"}
			p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestEquipment53EC40}, {Op: legacy.PortTestEquipment53EC80, Value: 1}, {Op: legacy.PortTestEquipment53EC80, Value: 2}, {Op: legacy.PortTestEquipment53EC40}}
			specs = append(specs, s)
		}
	}
	callbackHash(t, "equipment-drop-table", legacy.PortTestRoam(specs), equipmentHashes["equipment-drop-table"])
}
func TestEquipmentTryOwners(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for subject := 1; subject <= 3; subject++ {
		for _, class := range []uint32{8, 0x1000000, 0x2000000} {
			for _, state := range []byte{0, 1, 15, 16, 17, 18} {
				s := equipmentBase()
				p := s.Callbacks.Shop
				p.Resources.Subject = subject
				p.Items[0].Class = class
				p.Equipment.PlayerState = state
				p.Equipment.ItemTeamWord = []uint32{2}
				p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestEquipment4F2F70, Value: 1, Side: 1}, {Op: legacy.PortTestEquipment53E7B0}, {Op: legacy.PortTestEquipment4F2FB0, Value: 1, Side: 1}}
				specs = append(specs, s)
			}
		}
	}
	callbackHash(t, "equipment-try-owners", legacy.PortTestRoam(specs), equipmentHashes["equipment-try-owners"])
}

func TestEquipmentNullGuards(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, op := range []int{legacy.PortTestEquipment53A3D0, legacy.PortTestEquipment53AB90, legacy.PortTestEquipment53E600, legacy.PortTestEquipment4E7D30, legacy.PortTestEquipment4E7EC0, legacy.PortTestEquipment980523, legacy.PortTestEquipment9805EB, legacy.PortTestEquipment4F9FD0, legacy.PortTestEquipment53A6C0} {
		s := equipmentBase()
		p := s.Callbacks.Shop
		p.Equipment.NilUnit = true
		p.Sequence = []legacy.PortTestShopAction{{Op: op}}
		specs = append(specs, s)
	}
	for _, op := range []int{legacy.PortTestEquipment53A6C0, legacy.PortTestEquipment53AAB0, legacy.PortTestEquipment53EAE0, legacy.PortTestEquipment4E7EC0, legacy.PortTestEquipment53AB90, legacy.PortTestEquipment53EC80} {
		s := equipmentBase()
		p := s.Callbacks.Shop
		p.Inventory.NilItem = true
		p.Equipment.SecondaryWeapon = 2
		p.Sequence = []legacy.PortTestShopAction{{Op: op}}
		specs = append(specs, s)
	}
	s := equipmentBase()
	p := s.Callbacks.Shop
	p.Equipment.NilUnit = true
	p.Inventory.NilItem = true
	p.Equipment.Cheat = true
	p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestEquipment4F3180}}
	specs = append(specs, s)
	callbackHash(t, "equipment-null-guards", legacy.PortTestRoam(specs), equipmentHashes["equipment-null-guards"])
}
func TestEquipmentArmorMasks(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, class := range []uint32{8, 0x1000000, 0x2000000} {
		for _, mask := range []uint32{0, 1, 4, 8, 0xc00, 0xc0d, 0xffffffff} {
			s := equipmentBase()
			p := s.Callbacks.Shop
			p.Items[0].Class = class
			p.Inventory.ArmorBits = map[uint16]uint32{23: mask}
			p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestEquipment53E650}, {Op: legacy.PortTestEquipment53E2D0}, {Op: legacy.PortTestEquipment53E430}}
			specs = append(specs, s)
		}
	}
	r := legacy.PortTestRoam(specs)
	for i, result := range r {
		sp := specs[i].Callbacks.Shop
		want := sp.Items[0].Class&0x2000000 == 0 || sp.Inventory.ArmorBits[23]&0xc0d == 0
		if (result.Callbacks.Shop.Sequence[1].Return != 0) != want {
			t.Fatalf("armor mask case %d", i)
		}
	}
	callbackHash(t, "equipment-armor-masks", r, equipmentHashes["equipment-armor-masks"])
}
func TestEquipmentColdTable(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, names := range [][]string{{}, {"Sword"}, {"WoodenShield", "Sword", "SteelShield"}, {"missing-equipment-type"}} {
		for _, flags := range []uint32{0, 1, 2, 3, 0xffffffff} {
			s := equipmentBase()
			p := s.Callbacks.Shop
			p.Equipment.TableNames = names
			p.Equipment.ColdTable = true
			for range names {
				p.Equipment.TableFlags = append(p.Equipment.TableFlags, flags)
			}
			p.Inventory.ItemTypes = []string{"Sword"}
			p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestEquipment53EC80, Value: 1}, {Op: legacy.PortTestEquipment53EC80, Value: 2}, {Op: legacy.PortTestEquipment53EC40}}
			specs = append(specs, s)
		}
	}
	callbackHash(t, "equipment-cold-table", legacy.PortTestRoam(specs), equipmentHashes["equipment-cold-table"])
}
func TestEquipmentArmorSearch(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, linked := range [][]int{{}, {0}, {0, 1}, {0, 1, 2}, {2, 1, 0}} {
		for _, flags := range []uint32{0, 0x100} {
			s := equipmentBase()
			p := s.Callbacks.Shop
			p.Inventory.Linked = linked
			p.Inventory.Owned = linked
			p.Items[0].Flags = 0x100
			for i := 1; i < 3; i++ {
				p.Items[i].Class = 0x2000000
				p.Items[i].Flags = flags
				p.Items[i].Subclass = uint32(1 << i)
			}
			p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestEquipment53E7B0, Item: 2}}
			specs = append(specs, s)
		}
	}
	callbackHash(t, "equipment-armor-search", legacy.PortTestRoam(specs), equipmentHashes["equipment-armor-search"])
}

func TestEquipmentNPCSync(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, extra := range []uint32{0, 0x400000, 0x20000000} {
		for _, value := range []uint32{0, 1, 2, 0xffffffff} {
			for _, class := range []uint32{8, 0x1000000, 0x2000000} {
				for _, mask := range []uint32{0, 0x80000000, 0xffffffff} {
					s := equipmentBase()
					p := s.Callbacks.Shop
					p.Resources.Subject = 3
					p.Resources.Subclass = 0x10
					p.Resources.ExtraClass = extra
					p.Resources.SyncSeed = 0x12345678
					p.Items[0].Class = class
					p.Inventory.WeaponBits = map[uint16]uint32{23: mask}
					p.Inventory.ArmorBits = map[uint16]uint32{23: mask}
					p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestEquipment4E4B20, Value: value}}
					specs = append(specs, s)
				}
			}
		}
	}
	r := legacy.PortTestRoam(specs)
	for i, result := range r {
		if specs[i].Callbacks.Shop.Resources.ExtraClass != 0 && result.Callbacks.Shop.Sequence[0].Return != 65003 {
			t.Fatalf("NPC sync return identity case %d", i)
		}
	}
	callbackHash(t, "equipment-npc-sync", r, equipmentHashes["equipment-npc-sync"])
}
