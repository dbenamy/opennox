//go:build porttest

package opennox

import (
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

// Locked from repeated original-C captures before conversion.
var inventoryHashes = map[string]string{
	"inventory-team-members":         "041bc138a731c8780f34a350212d874c25abadad419dd08c5f20697e0138117a",
	"inventory-links":                "1edf30144b9c04c839998e201dead78cd0ce74154be62b36e031e5156d9256e1",
	"inventory-drop-eligibility":     "307f6239ef734594d372da5bc4af1914310ba4b97f441dfb4d804feafeaaf66a",
	"inventory-pickup-boundaries":    "d53014d23cc948e63362559a6523c3cf7b24a6e107c43c53e82c3d1ae4a13f4e",
	"inventory-default-drops":        "bfd558e106f59d2838bf6175cf588be92310f933deebd876223c5fcbad9f369c",
	"inventory-dispatch-placement":   "b56c96c4eab82752739d70260f098f3a361776336a7919229b592ec578a7fa7c",
	"inventory-chest":                "e5a0889898bd20c65c736a1a76af1a5bd372d90995415bf526b417f74f384e9d",
	"inventory-equipment-pickup":     "7dee69cabd7d66eeeab4bd430f494bb0145fa6d65a205136ac2a4e50c6c84e95",
	"inventory-crown-treasure":       "b01343352cd4b3d867c2abe0b2d9adf2c97f32b92346b0ae1d6e88ed6ff2b29b",
	"inventory-drop-rejections":      "b2f84da3080b2167bdfcc697fed07a25ca7b47419d44a3e28a54f572ff4c6fbd",
	"inventory-ammo-merge":           "cf1e3547630470ee73cae35a26d5e4cd2f21f009f5a1564b44867ba009ef8c83",
	"inventory-food-audio":           "9ce469d1508ed2f8190c09f2bcda341545d8be14d75fc924552c170ca16ecccc",
	"inventory-drop-special-items":   "13555ce7040eb218d5d7410a4c8afe656726f4fec1e7e0875f0289c628b01ca5",
	"inventory-geometry-boundaries":  "85bc742d746565f946d1dbf88a0159c13450cac1650cce8363c1995957ecabff",
	"inventory-drop-all-integration": "a5f70649fe39dcf02e65becbeb2cfd6525827e77d2021ac42e7797067b6491b3",
	"inventory-armor-replacement":    "4a31bdede69f15c2e47b4274c7640334a6be4570b7e71bfee1de34721b8b2e8f",
	"inventory-team-objectives":      "06cdb38d017296960d2ca764e4d05930934680afbc1820b993e5baa5d2d2f372",
}

func inventoryBase() legacy.PortTestRoamSpec {
	s := resourceBase()
	p := s.Callbacks.Shop
	p.CaptureData = true
	p.Resources.PickupResult = true
	p.Inventory = &legacy.PortTestInventorySpec{Carry: 64, Position: types.Pointf{X: 100, Y: 100}, Target: types.Pointf{X: 110, Y: 120}, Radius: 50, DropResult: true}
	p.Items = []legacy.PortTestShopItem{{Type: 23, Class: 8}, {Type: 23, Class: 8}, {Type: 23, Class: 8}}
	return s
}
func TestInventoryLinks(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for subject := 1; subject < 4; subject++ {
		for _, uf := range []uint32{0, 0x20, 0x8000} {
			for _, flags := range []uint32{0, 0x20, 0x100, 0x10000000} {
				for _, weight := range []byte{0, 1, 255} {
					for _, protected := range []bool{false, true} {
						s := inventoryBase()
						p := s.Callbacks.Shop
						p.Resources.Subject = subject
						p.Resources.Flags = uf
						p.Resources.Protected = protected
						p.Inventory.Linked = []int{0, 1}
						p.Inventory.Weights = []byte{weight, weight, weight}
						p.Items[2].Flags = flags
						p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestInventory4F3070, Item: 2, Value: 1}, {Op: legacy.PortTestInventory4ED0C0, Item: 1}, {Op: legacy.PortTestInventory4ED0C0, Item: 0}}
						if uf&0x20 == 0 && flags&0x20 == 0 {
							p.Sequence = append(p.Sequence, legacy.PortTestShopAction{Op: legacy.PortTestInventory4ED0C0, Item: 2})
						}
						specs = append(specs, s)
					}
				}
			}
		}
	}
	callbackHash(t, "inventory-links", legacy.PortTestRoam(specs), inventoryHashes["inventory-links"])
}
func TestInventoryDropEligibility(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for subject := 1; subject < 4; subject++ {
		for _, flags := range []uint32{0, 0x20, 0x10000000, 0x10000020, 0xffffffff} {
			s := inventoryBase()
			p := s.Callbacks.Shop
			p.Resources.Subject = subject
			p.Items[0].Flags = flags
			p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestInventory4EDCD0}, {Op: legacy.PortTestInventory53EBF0}}
			specs = append(specs, s)
		}
	}
	r := legacy.PortTestRoam(specs)
	for i, v := range r {
		s := specs[i].Callbacks.Shop
		want := uint32(0)
		if s.Items[0].Flags&0x20 != 0 || s.Resources.Subject == 2 || s.Items[0].Flags&0x10000000 == 0 {
			want = 1
		}
		if v.Callbacks.Shop.Sequence[0].Return != want {
			t.Fatalf("drop eligibility case %d", i)
		}
	}
	callbackHash(t, "inventory-drop-eligibility", r, inventoryHashes["inventory-drop-eligibility"])
}
func TestInventoryPickupBoundaries(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	ops := []int{legacy.PortTestInventory4F3350, legacy.PortTestInventory4F34D0, legacy.PortTestInventory4F3510, legacy.PortTestInventory4F3C60, legacy.PortTestInventory4F3CE0, legacy.PortTestInventory4F3DD0}
	for _, op := range ops {
		for subject := 1; subject < 4; subject++ {
			for _, success := range []bool{false, true} {
				for _, used := range []bool{false, true} {
					for _, gf := range []uint32{0, 2048, 4096} {
						s := inventoryBase()
						p := s.Callbacks.Shop
						p.Resources.Subject = subject
						p.Resources.PickupResult = success
						p.Inventory.UseDelete = used
						s.Lifecycle.GameFlags = gf
						p.Sequence = []legacy.PortTestShopAction{{Op: op, Value: 0xffffffff}}
						specs = append(specs, s)
					}
				}
			}
		}
	}
	callbackHash(t, "inventory-pickup-boundaries", legacy.PortTestRoam(specs), inventoryHashes["inventory-pickup-boundaries"])
}

func TestInventoryDefaultDrops(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	ops := []int{legacy.PortTestInventory4ED290, legacy.PortTestInventory4ED500, legacy.PortTestInventory4ED580, legacy.PortTestInventory4ED710, legacy.PortTestInventory4EDDE0, legacy.PortTestInventory4EDE50, legacy.PortTestInventory4EE370, legacy.PortTestInventory53AB10, legacy.PortTestInventory53EB70}
	for _, op := range ops {
		for subject := 1; subject < 4; subject++ {
			for _, linked := range []bool{false, true} {
				for _, gf := range []uint32{0, 2048, 4096} {
					s := inventoryBase()
					p := s.Callbacks.Shop
					p.Resources.Subject = subject
					s.Lifecycle.GameFlags = gf
					if linked {
						p.Inventory.Linked = []int{0, 1, 2}
						p.Inventory.Owned = []int{0, 1, 2}
					}
					p.Sequence = []legacy.PortTestShopAction{{Op: op, Item: 1}}
					specs = append(specs, s)
				}
			}
		}
	}
	callbackHash(t, "inventory-default-drops", legacy.PortTestRoam(specs), inventoryHashes["inventory-default-drops"])
}
func TestInventoryDispatchPlacement(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, op := range []int{legacy.PortTestInventory4ED790, legacy.PortTestInventory4ED810, legacy.PortTestInventory4ED930, legacy.PortTestInventory4ED970, legacy.PortTestInventory4EDA40} {
		for _, wall := range []int{0, 1, 2} {
			for _, seed := range []int{1, 3, 17, 91} {
				for _, success := range []bool{false, true} {
					s := inventoryBase()
					p := s.Callbacks.Shop
					p.Inventory.Linked = []int{0, 1, 2}
					p.Inventory.DropResult = success
					p.Inventory.Target = types.Pointf{X: 200, Y: 100}
					p.Inventory.WallMode = wall
					s.Seed = seed
					p.Sequence = []legacy.PortTestShopAction{{Op: op, Item: 1}}
					specs = append(specs, s)
				}
			}
		}
	}
	callbackHash(t, "inventory-dispatch-placement", legacy.PortTestRoam(specs), inventoryHashes["inventory-dispatch-placement"])
}
func TestInventoryChest(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, sub := range []uint32{0, 0x100, 0x200, 0x400, 0x800} {
		for n := 0; n <= 3; n++ {
			for wall := 0; wall < 3; wall++ {
				s := inventoryBase()
				p := s.Callbacks.Shop
				p.Resources.Subject = 2
				p.Resources.Subclass = sub
				p.Inventory.Shape = [5]uint32{2, 0x41200000}
				p.Inventory.WallMode = wall
				for i := 0; i < n; i++ {
					p.Inventory.Linked = append(p.Inventory.Linked, i)
				}
				p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestInventory4EE2A0}, {Op: legacy.PortTestInventory4EDF00, Item: 2}}
				specs = append(specs, s)
			}
		}
	}
	callbackHash(t, "inventory-chest", legacy.PortTestRoam(specs), inventoryHashes["inventory-chest"])
}

func TestInventoryEquipmentPickup(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, op := range []int{legacy.PortTestInventory53A720, legacy.PortTestInventory4F3B00, legacy.PortTestInventory53E7F0, legacy.PortTestInventory53A9C0} {
		for subject := 1; subject < 4; subject++ {
			for _, eligible := range []bool{false, true} {
				for _, gf := range []uint32{0, 2048, 4096} {
					for _, success := range []bool{false, true} {
						for _, blocked := range []bool{false, true} {
							s := inventoryBase()
							p := s.Callbacks.Shop
							p.Resources.Subject = subject
							p.Resources.PickupResult = success
							s.Lifecycle.Eligible = eligible
							s.Lifecycle.GameFlags = gf
							p.Inventory.Blocked = blocked
							p.Inventory.PickupInsert = true
							p.Inventory.WeaponBits = map[uint16]uint32{23: 4}
							p.Inventory.ArmorBits = map[uint16]uint32{23: 2}
							p.Items[0].Class = 0x1000000
							p.Items[0].Subclass = 4
							if op == legacy.PortTestInventory53E7F0 {
								p.Items[0].Class = 0x2000000
								p.Items[0].Subclass = 2
							}
							if op == legacy.PortTestInventory53A9C0 {
								p.Items[0].Subclass = 0x800000
							}
							p.Sequence = []legacy.PortTestShopAction{{Op: op, Value: 1, Side: 1}}
							specs = append(specs, s)
						}
					}
				}
			}
		}
	}
	callbackHash(t, "inventory-equipment-pickup", legacy.PortTestRoam(specs), inventoryHashes["inventory-equipment-pickup"])
}
func TestInventoryCrownTreasure(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, op := range []int{legacy.PortTestInventory4F3400, legacy.PortTestInventory4ED5E0, legacy.PortTestInventory4F3580, legacy.PortTestInventory4ED710} {
		for subject := 1; subject < 4; subject++ {
			for _, gf := range []uint32{0, 16, 64} {
				for _, success := range []bool{false, true} {
					s := inventoryBase()
					p := s.Callbacks.Shop
					p.Resources.Subject = subject
					p.Resources.PickupResult = success
					p.Inventory.PickupInsert = true
					s.Lifecycle.GameFlags = gf
					p.Inventory.Treasure = 1
					p.Inventory.TreasureMax = 3
					if op == legacy.PortTestInventory4ED5E0 || op == legacy.PortTestInventory4ED710 {
						p.Inventory.Linked = []int{0, 1, 2}
						p.Inventory.Owned = []int{0, 1, 2}
					}
					p.Sequence = []legacy.PortTestShopAction{{Op: op, Value: 1}}
					specs = append(specs, s)
				}
			}
		}
	}
	callbackHash(t, "inventory-crown-treasure", legacy.PortTestRoam(specs), inventoryHashes["inventory-crown-treasure"])
}
func TestInventoryDropRejections(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, flags := range []uint32{0, 0x20, 0x8000, 0x8020} {
		for _, mask := range []uint32{0, 1, 2, 3} {
			for _, typ := range []uint32{23, 24} {
				s := inventoryBase()
				p := s.Callbacks.Shop
				p.Resources.Flags = flags
				p.Inventory.Linked = []int{0}
				p.Inventory.DropTable = [][3]uint32{{1, typ, mask}}
				p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestInventory53EBF0}, {Op: legacy.PortTestInventory4ED290}}
				specs = append(specs, s)
			}
		}
	}
	callbackHash(t, "inventory-drop-rejections", legacy.PortTestRoam(specs), inventoryHashes["inventory-drop-rejections"])
}

func TestInventoryAmmoMerge(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, stored := range []byte{0, 1, 100, 249, 250, 255} {
		for _, incoming := range []byte{0, 1, 150, 151, 250, 255} {
			for variant := 0; variant < 4; variant++ {
				s := inventoryBase()
				p := s.Callbacks.Shop
				p.Inventory.Linked = []int{0, 1}
				p.Inventory.Owned = []int{0, 1}
				p.Inventory.WeaponBits = map[uint16]uint32{23: 2}
				p.Inventory.PickupInsert = true
				p.Inventory.Blocked = true
				for i := range p.Items {
					p.Items[i].Class = 0x1000000
					p.Items[i].Subclass = 2
					p.Items[i].Use[1] = 200
				}
				p.Items[0].Use[0] = stored
				p.Items[1].Use[0] = 250
				p.Items[2].Use[0] = incoming
				if variant&1 != 0 {
					p.Items[0].Use[2] = 1
				}
				if variant&2 != 0 {
					p.Items[0].Mods[0] = true
				}
				p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestInventory4F3B00, Item: 2, Value: 0xffffffff, Side: 0x12345678}}
				specs = append(specs, s)
			}
		}
	}
	callbackHash(t, "inventory-ammo-merge", legacy.PortTestRoam(specs), inventoryHashes["inventory-ammo-merge"])
}
func TestInventoryFoodAudio(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, op := range []int{legacy.PortTestInventory4F3350, legacy.PortTestInventory4EDE50} {
		for _, mat := range []uint16{0, 1, 2, 3, 0xffff} {
			for _, sub := range []uint32{0, 1, 2, 4, 128, 130, 132, 255} {
				for variant := 0; variant < 4; variant++ {
					s := inventoryBase()
					p := s.Callbacks.Shop
					p.Inventory.FoodDrop = [][3]uint32{{0, 1, 835}, {2, 0, 837}, {4, 0, 833}, {128, 0, 839}}
					p.Inventory.FoodPickup = [][3]uint32{{0, 1, 834}, {2, 0, 836}, {4, 0, 832}, {128, 0, 838}}
					p.Inventory.Materials = []uint16{mat}
					p.Items[0].Subclass = sub
					p.Inventory.UseDelete = variant&1 != 0
					p.Resources.PickupResult = variant&2 != 0
					if op == legacy.PortTestInventory4EDE50 {
						p.Inventory.Linked = []int{0}
					}
					p.Sequence = []legacy.PortTestShopAction{{Op: op, Value: 1}}
					specs = append(specs, s)
				}
			}
		}
	}
	callbackHash(t, "inventory-food-audio", legacy.PortTestRoam(specs), inventoryHashes["inventory-food-audio"])
}
func TestInventoryDropSpecialItems(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, class := range []uint32{8, 64, 2, 0x10000000, 0x1000000} {
		for variant := 0; variant < 4; variant++ {
			for _, gf := range []uint32{0, 32, 2048, 4096} {
				s := inventoryBase()
				p := s.Callbacks.Shop
				p.Inventory.Linked = []int{0, 1}
				p.Inventory.Owned = []int{0, 1}
				p.Items[0].Class = class
				p.Inventory.WeaponBits = map[uint16]uint32{23: 4}
				p.Inventory.ServerFlags = 2
				s.Lifecycle.GameFlags = gf
				if variant&1 != 0 {
					p.Items[0].Flags = 0x10000
				}
				if variant&2 != 0 {
					p.Items[0].Flags |= 0x80000
				}
				p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestInventory4ED290}}
				specs = append(specs, s)
			}
		}
	}
	for _, name := range []string{"Glyph", "Torch", "Lantern"} {
		s := inventoryBase()
		p := s.Callbacks.Shop
		p.Inventory.ItemTypes = []string{name}
		p.Inventory.Linked = []int{0}
		p.Resources.Buffs = 1 << 15
		p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestInventory4ED290}}
		specs = append(specs, s)
	}
	callbackHash(t, "inventory-drop-special-items", legacy.PortTestRoam(specs), inventoryHashes["inventory-drop-special-items"])
}
func TestInventoryGeometryBoundaries(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	bits := []uint32{0, 0x80000000, 1, 0xbf800000, 0x42c80000, 0x7f7fffff, 0x7f800000, 0xff800000, 0x7fc12345, 0x7f812345}
	for _, kind := range []uint32{0, 1, 2, 3, 0xffffffff} {
		for i, b := range bits {
			s := inventoryBase()
			p := s.Callbacks.Shop
			p.Inventory.Shape = [5]uint32{kind, b, 0, b, bits[(i+3)%len(bits)]}
			p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestInventory4EE2A0}}
			specs = append(specs, s)
		}
	}
	for _, radius := range []float32{0, 1, -1, 50, 75.25, 256} {
		for _, pos := range []types.Pointf{{X: 512, Y: 512}, {X: 2048, Y: 2048}, {X: 4000, Y: 4000}} {
			for _, seed := range []int{1, 7, 256} {
				s := inventoryBase()
				p := s.Callbacks.Shop
				p.Inventory.Position = pos
				p.Inventory.Radius = radius
				s.Seed = seed
				p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestInventory4ED970}}
				specs = append(specs, s)
			}
		}
	}
	callbackHash(t, "inventory-geometry-boundaries", legacy.PortTestRoam(specs), inventoryHashes["inventory-geometry-boundaries"])
}
func TestInventoryDropAllIntegration(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for n := 0; n <= 3; n++ {
		for variant := 0; variant < 8; variant++ {
			for _, gf := range []uint32{0, 4096, 8192} {
				s := inventoryBase()
				p := s.Callbacks.Shop
				s.Lifecycle.GameFlags = gf
				p.Inventory.DefaultDrop = true
				for i := 0; i < n; i++ {
					p.Inventory.Linked = append(p.Inventory.Linked, i)
					p.Inventory.Owned = append(p.Inventory.Owned, i)
				}
				if variant&1 != 0 {
					p.Items[0].Flags = 0x10000000
				}
				if variant&2 != 0 {
					p.Items[0].Flags |= 0x20
				}
				if variant&4 != 0 {
					p.Inventory.Position = types.Pointf{X: -1000, Y: -1000}
				}
				p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestInventory4EDA40}}
				specs = append(specs, s)
			}
		}
	}
	callbackHash(t, "inventory-drop-all-integration", legacy.PortTestRoam(specs), inventoryHashes["inventory-drop-all-integration"])
}

func TestInventoryArmorReplacement(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, old := range []string{"StreetSneakers", "WizardRobe", "WoodenShield", "SteelShield", "Torch"} {
		for _, incoming := range []string{"WoodenShield", "SteelShield", "WizardRobe"} {
			for _, sub := range []uint32{2, 4, 0x20} {
				for mods := 0; mods < 4; mods++ {
					s := inventoryBase()
					p := s.Callbacks.Shop
					s.Lifecycle.Eligible = true
					p.Inventory.PickupInsert = true
					p.Inventory.Linked = []int{0}
					p.Inventory.Owned = []int{0}
					p.Inventory.ItemTypes = []string{old, incoming}
					p.Inventory.ArmorBits = map[uint16]uint32{28: 2, 31: 2, 32: 2, 33: 2, 34: 2}
					p.Inventory.Materials = []uint16{0, []uint16{2, 4, 8, 16}[mods]}
					for i := 0; i < 2; i++ {
						p.Items[i].Class = 0x2000000
						p.Items[i].Subclass = sub
					}
					p.Items[0].Flags = 0x100
					p.Items[0].Mods[0] = mods&1 != 0
					p.Items[1].Mods[3] = mods&2 != 0
					p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestInventory53E7F0, Item: 1, Value: 1, Side: 1}}
					specs = append(specs, s)
				}
			}
		}
	}
	callbackHash(t, "inventory-armor-replacement", legacy.PortTestRoam(specs), inventoryHashes["inventory-armor-replacement"])
}

func TestInventoryTeamObjectives(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, teams := range [][4]byte{{}, {1, 1, 2, 1}, {1, 2, 1, 1}, {2, 1, 1, 1}, {1, 1, 1, 1}, {1, 1, 2, 3}} {
		for _, op := range []int{legacy.PortTestInventory4F3400, legacy.PortTestInventory4ED5E0, legacy.PortTestInventory4F3580} {
			for _, treasure := range []uint32{0, 1, 2, 3, 0xffffffff} {
				s := inventoryBase()
				p := s.Callbacks.Shop
				p.Inventory.Teams = teams
				p.Inventory.CrownTimes = [3]uint32{30, 20, 10}
				p.Inventory.Gameplay = 4
				p.Inventory.Treasure = treasure
				p.Inventory.TreasureMax = 3
				s.Lifecycle.GameFlags = 16 | 64
				p.Inventory.PickupInsert = true
				if op == legacy.PortTestInventory4ED5E0 {
					p.Inventory.Linked = []int{0}
					p.Inventory.Owned = []int{0}
				}
				p.Sequence = []legacy.PortTestShopAction{{Op: op, Value: 1}}
				specs = append(specs, s)
			}
		}
	}
	callbackHash(t, "inventory-team-objectives", legacy.PortTestRoam(specs), inventoryHashes["inventory-team-objectives"])
}

func TestInventoryTeamMembers(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, teams := range [][4]byte{{}, {1, 1, 2, 1}, {1, 2, 1, 1}, {2, 1, 1, 1}, {1, 1, 1, 1}, {1, 1, 2, 3}} {
		for _, op := range []int{legacy.PortTestInventory4F3400, legacy.PortTestInventory4ED5E0, legacy.PortTestInventory4F3580} {
			for _, treasure := range []uint32{0, 1, 2, 3, 0xffffffff} {
				s := inventoryBase()
				p := s.Callbacks.Shop
				p.Inventory.Teams = teams
				p.Inventory.TeamMembers = true
				p.Inventory.CrownTimes = [3]uint32{30, 20, 10}
				p.Inventory.Gameplay = 4
				p.Inventory.Treasure = treasure
				p.Inventory.TreasureMax = 3
				s.Lifecycle.GameFlags = 16 | 64
				p.Inventory.PickupInsert = true
				if op == legacy.PortTestInventory4ED5E0 {
					p.Inventory.Linked = []int{0}
					p.Inventory.Owned = []int{0}
				}
				p.Sequence = []legacy.PortTestShopAction{{Op: op, Value: 1}}
				specs = append(specs, s)
			}
		}
	}
	callbackHash(t, "inventory-team-members", legacy.PortTestRoam(specs), inventoryHashes["inventory-team-members"])
}
