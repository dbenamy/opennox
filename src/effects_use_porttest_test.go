//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"math"
	"testing"
)

var effectsUseHashes = map[string]string{
	"effects-flags":              "49b44e5e18a52a0390546deb1f8fc7c55a702c4a7fbe9acab3d643e013c29155",
	"effects-arithmetic":         "e51485e0fc9fc7be0d4bfdffdee01ba94330d0bd727732cfeaff3b3a0f908d66",
	"effects-grip-inversion":     "7230841eacf523078380cb1504c8190dedb13cdfa5b628400efcdcf7ecf70951",
	"effects-readiness-recharge": "3adb36ce83c80bd762f8106d095206e980a9271a69586e1f8af4c6f92ccb7add",
	"effects-recharge-percent":   "2fcabc1f127b311d34f806f475c66686510b1fa1380cc38f561d8d4e6c9c5426",
	"effects-protection":         "bf26996b3a31ad6a9b2ee1a4f8a15d82278c74464bbd6ded950ab20a0db05148",
	"effects-regeneration":       "a6ad82700eb43b762cc2bdc59c457308833a4f3f068439b0fdc585e86c09637d",
	"effects-replenishment":      "c6cf5d5dab9abf7399fd0128e9d308560fda5bbd9d9ac9d11dd1886092d9773b",
	"effects-resource-transfers": "ceefac4eeccc214a6f4615f387ef039024a57939232333490ce2b87e42eaf961",
	"effects-wand-acceptance":    "36a082bb7ade6805f904c168b66e3e80fc936a6c3ce3da375c8055204006b22f",
	"effects-use-dispatch":       "9599becd1091ce32f63b6a18ed3747e90e1ce389b29abfc1c4b54e1003e82a22",
	"effects-status-force":       "48960c82a2fd5a4a5e2e7bf43a4035b9d55a89baff01f8b9ca493cc6c2c19876",
	"effects-projectiles":        "e4bc352bfdc93311b885d054341e9d973dbd74d14b3bb4130170f07f7133c543",
	"effects-fire-wand":          "89b37203dd0cf0ca981b7ade2772ab19b574083d7f0c6fc269234602bed21553",
	"effects-inventory-lookup":   "9e2eb648c66bdc1a68109420e59d3e9bb3ffec3a8f55938ecaabc562d95aae9a",
	"effects-speed-rounding":     "aa442a64e52b00de68e40988a3a9417dda6ae86ccf46690bddfefe587e90f511",
	"effects-wand-classes":       "923eacb65d34f28b2d78e02dc61c8b5e66a198f20677a2b826df1d54aa1eb14c",
	"effects-null-guards":        "1520fec6df341e5d086cc049c99dc8aa94ca3ba1efcc49e865a3f9342ce5da42",
	"effects-fireball-wall":      "5effe12089787b193edbbb8ad0904000a97f4c95364633acdd29217973ad73a0",
}

func effectsUseBase() legacy.PortTestRoamSpec {
	s := equipmentBase()
	p := s.Callbacks.Shop
	p.EffectsUse = &legacy.PortTestEffectsUseSpec{}
	for i := range p.Items {
		p.Items[i].Flags = 0x100
	}
	return s
}
func effectsIdentity(op int) int { return op - 499 }
func TestEffectsFlags(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	kinds := [][3]int{{legacy.PortTestEffects4DFB50, legacy.PortTestEffects4DFB80, 8}, {legacy.PortTestEffects4DFC30, legacy.PortTestEffects4DFCA0, 16}, {legacy.PortTestEffects4DFD10, legacy.PortTestEffects4DFD40, 1}, {legacy.PortTestEffects4DFD80, legacy.PortTestEffects4DFDB0, 4}, {legacy.PortTestEffects4DFDE0, legacy.PortTestEffects4DFE10, 2}, {legacy.PortTestEffects4E0140, legacy.PortTestEffects4E0170, 32}}
	for _, kind := range kinds {
		for subject := 1; subject <= 3; subject++ {
			for mask := 0; mask < 16; mask++ {
				for _, equipped := range []uint32{0, 0x100} {
					s := effectsUseBase()
					p := s.Callbacks.Shop
					p.Resources.Subject = subject
					for i := range p.Items {
						p.Items[i].Flags = equipped
					}
					p.EffectsUse.UnitWords = map[int]uint32{440: 0xaabbccff, 552: math.Float32bits(1)}
					for i := 0; i < 4; i++ {
						p.Items[0].Mods[i] = mask&(1<<i) != 0
						p.EffectsUse.Modifiers[i] = legacy.PortTestEffectsModifier{Engage: effectsIdentity(kind[0]), Words: map[int]uint32{120: math.Float32bits(.25)}}
					}
					p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestEffects4DFBB0, Value: uint32(kind[2])}, {Op: kind[0]}, {Op: kind[1]}, {Op: legacy.PortTestEffects4DFBB0, Value: uint32(kind[2])}}
					specs = append(specs, s)
				}
			}
		}
	}
	callbackHash(t, "effects-flags", legacy.PortTestRoam(specs), effectsUseHashes["effects-flags"])
}
func TestEffectsArithmetic(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	bits := []uint32{0, 0x80000000, 1, 0x3e800000, 0x3f800000, 0x40000000, 0xbf800000, 0x7f800000, 0xff800000, 0x7fc12345}
	for _, op := range []int{legacy.PortTestEffects4E0370, legacy.PortTestEffects4E0380, legacy.PortTestEffects4E04C0, legacy.PortTestEffects4E09B0} {
		for _, a := range bits {
			for _, b := range bits {
				s := effectsUseBase()
				p := s.Callbacks.Shop
				p.EffectsUse.Scalar = b
				p.EffectsUse.TargetWords = map[int]uint32{544: b}
				p.EffectsUse.Modifiers[0].Words = map[int]uint32{44: a, 80: a}
				p.Sequence = []legacy.PortTestShopAction{{Op: op}}
				specs = append(specs, s)
			}
		}
	}
	callbackHash(t, "effects-arithmetic", legacy.PortTestRoam(specs), effectsUseHashes["effects-arithmetic"])
}
func TestEffectsGripAndInversion(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, v := range []uint32{0, 1, 2, 0x7fffffff, 0x80000000, 0xffffffff} {
		s := effectsUseBase()
		p := s.Callbacks.Shop
		p.EffectsUse.Scalar = 0x12345678
		p.EffectsUse.Modifiers[0].Words = map[int]uint32{96: v}
		p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestEffects4E03D0}, {Op: legacy.PortTestEffects4E0480}}
		specs = append(specs, s)
	}
	for _, class := range []uint32{0, 0x1000, 0x1000000, 0x2000000, 0x10000000} {
		for mask := 0; mask < 16; mask++ {
			for _, v := range []uint32{0, 1, 0xffffffff} {
				s := effectsUseBase()
				p := s.Callbacks.Shop
				p.Items[0].Class = class
				for i := 0; i < 4; i++ {
					p.Items[0].Mods[i] = mask&(1<<i) != 0
					p.EffectsUse.Modifiers[i] = legacy.PortTestEffectsModifier{Collide: effectsIdentity(legacy.PortTestEffects4E0480), Words: map[int]uint32{96: v}}
				}
				p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestEffects4E03F0}}
				specs = append(specs, s)
			}
		}
	}
	callbackHash(t, "effects-grip-inversion", legacy.PortTestRoam(specs), effectsUseHashes["effects-grip-inversion"])
}
func TestEffectsReadinessRecharge(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, class := range []uint32{0, 0x1000, 0x1000000, 0x2000000, 0x10000000} {
		for mask := 0; mask < 16; mask++ {
			for _, fn := range []int{0, 42, 43} {
				s := effectsUseBase()
				p := s.Callbacks.Shop
				p.Items[0].Class = class
				for i := 0; i < 4; i++ {
					p.Items[0].Mods[i] = mask&(1<<i) != 0
					p.EffectsUse.Modifiers[i] = legacy.PortTestEffectsModifier{Attack: fn, Words: map[int]uint32{48: uint32(3 + i)}}
				}
				p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestEffects4E0960}, {Op: legacy.PortTestEffects53C940}}
				specs = append(specs, s)
			}
		}
	}
	callbackHash(t, "effects-readiness-recharge", legacy.PortTestRoam(specs), effectsUseHashes["effects-readiness-recharge"])
}
func TestEffectsRechargePercent(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, percent := range []uint32{0, 1, 49, 99, 100, 101, 0xffffffff, 0x7fffffff} {
		for _, delta := range []uint32{0, 1, 50, 100, 0xffffffff} {
			for _, maximum := range []byte{0, 1, 3, 100, 255} {
				s := effectsUseBase()
				p := s.Callbacks.Shop
				p.Items[0].Class = 0x1000
				p.Items[0].Use[108] = 2
				p.Items[0].Use[109] = maximum
				for i := 0; i < 4; i++ {
					p.Items[0].Use[112+i] = byte(percent >> (8 * i))
				}
				p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestEffects53C520, Value: delta}, {Op: legacy.PortTestEffects53C520, Value: delta}}
				specs = append(specs, s)
			}
		}
	}
	callbackHash(t, "effects-recharge-percent", legacy.PortTestRoam(specs), effectsUseHashes["effects-recharge-percent"])
}

func TestEffectsProtection(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, kind := range [][3]int{{legacy.PortTestEffects4DFE40, legacy.PortTestEffects4DFD10, 17}, {legacy.PortTestEffects4DFF40, legacy.PortTestEffects4DFD80, 20}, {legacy.PortTestEffects4E0040, legacy.PortTestEffects4DFDE0, 18}} {
		for _, value := range []float32{-.25, 0, .01, .015625, .1, .125, .25, .5, .6, .7, 1, float32(math.Inf(1)), math.Float32frombits(0x7fc12345)} {
			for _, power := range []byte{0, 1, 3} {
				for _, item := range []bool{false, true} {
					for mask := 0; mask < 4; mask++ {
						s := effectsUseBase()
						p := s.Callbacks.Shop
						p.EffectsUse.UnitItem = item
						p.EffectsUse.BuffPower = power
						if power != 0 {
							p.Resources.Buffs = 1 << kind[2]
							p.EffectsUse.ItemWords = []map[int]uint32{{340: 1 << kind[2], 408 + 4*(kind[2]/4): uint32(power) << (8 * (kind[2] % 4))}}
						}
						p.EffectsUse.Balance = map[string][]float64{"FireSpellProtection": {.1, .2, .3}, "ElectricitySpellProtection": {.15, .25, .35}, "PoisonSpellProtection": {.2, .3, .4}}
						for i := 0; i < 4; i++ {
							p.EffectsUse.Modifiers[i] = legacy.PortTestEffectsModifier{Engage: effectsIdentity(kind[1]), Words: map[int]uint32{120: math.Float32bits(value)}}
							for j := range p.Items {
								p.Items[j].Mods[i] = mask&(1<<(i%2)) != 0
							}
						}
						p.Sequence = []legacy.PortTestShopAction{{Op: kind[0]}}
						specs = append(specs, s)
					}
				}
			}
		}
	}
	callbackHash(t, "effects-protection", legacy.PortTestRoam(specs), effectsUseHashes["effects-protection"])
}
func TestEffectsRegeneration(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, frame := range []uint32{0, 1, 29, 30, 60, 90, 180, 0xffffffff} {
		for _, flags := range []uint32{0, 0x20, 0x8000} {
			for _, hp := range []uint16{0, 50, 100} {
				for _, armor := range []bool{false, true} {
					for _, linked := range []bool{false, true} {
						s := effectsUseBase()
						p := s.Callbacks.Shop
						s.Owner.Frame = frame
						s.Owner.FPS = 30
						p.Resources.HP = hp
						p.Resources.Flags = flags
						p.EffectsUse.UnitWords = map[int]uint32{536: 0}
						p.EffectsUse.Modifiers[0].Words = map[int]uint32{108: 300}
						if !linked {
							p.Inventory.Linked = nil
						}
						if armor {
							p.Items[0].Class = 0x2000000
							p.Inventory.ArmorBits[23] = 0x4000
						}
						p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestEffects4E01D0}, {Op: legacy.PortTestEffects4E01D0}}
						specs = append(specs, s)
					}
				}
			}
		}
	}
	callbackHash(t, "effects-regeneration", legacy.PortTestRoam(specs), effectsUseHashes["effects-regeneration"])
}
func TestEffectsReplenishment(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, frame := range []uint32{0, 1, 2, 3, 0xffffffff} {
		for _, current := range []byte{0, 1, 254, 255} {
			for _, maximum := range []byte{0, 1, 100, 255} {
				for _, linked := range []bool{false, true} {
					s := effectsUseBase()
					p := s.Callbacks.Shop
					s.Owner.Frame = frame
					p.Items[0].Class = 0x1000
					p.Items[0].Use[108] = current
					p.Items[0].Use[109] = maximum
					p.EffectsUse.Modifiers[0].Words = map[int]uint32{108: 3}
					if !linked {
						p.Inventory.Linked = nil
					}
					p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestEffects4E02C0}, {Op: legacy.PortTestEffects4E02C0}}
					specs = append(specs, s)
				}
			}
		}
	}
	callbackHash(t, "effects-replenishment", legacy.PortTestRoam(specs), effectsUseHashes["effects-replenishment"])
}
func TestEffectsResourceTransfers(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, op := range []int{legacy.PortTestEffects4E0740, legacy.PortTestEffects4E07C0, legacy.PortTestEffects4E08E0} {
		for _, damage := range []uint32{0, 1, 9, 50, 200, 0xffffffff} {
			for _, factor := range []float32{0, .1, .5, 1, 2, -.5} {
				for _, target := range []int{0, 1, 4} {
					s := effectsUseBase()
					p := s.Callbacks.Shop
					p.EffectsUse.Target = target
					p.EffectsUse.Scalar = damage
					p.EffectsUse.Modifiers[0].Words = map[int]uint32{68: math.Float32bits(factor)}
					p.Sequence = []legacy.PortTestShopAction{{Op: op}}
					specs = append(specs, s)
				}
			}
		}
	}
	callbackHash(t, "effects-resource-transfers", legacy.PortTestRoam(specs), effectsUseHashes["effects-resource-transfers"])
}
func TestEffectsWandAcceptance(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, subject := range []int{1, 3} {
		for _, frame := range []uint32{0, 1, 9, 10, 100, 0xffffffff} {
			for _, charge := range [][2]byte{{0, 0}, {0, 10}, {1, 10}, {255, 255}} {
				for _, accept := range []bool{false, true} {
					for _, target := range []bool{false, true} {
						s := effectsUseBase()
						p := s.Callbacks.Shop
						p.Resources.Subject = subject
						s.Owner.Frame = frame
						p.Items[0].Class = 0x1000
						p.Items[0].Use[100] = 10
						p.Items[0].Use[92] = 5
						p.Items[0].Use[108], p.Items[0].Use[109] = charge[0], charge[1]
						p.EffectsUse.AcceptSpell = accept
						p.EffectsUse.CastTarget = target
						p.EffectsUse.PlayerWords = map[int]uint32{2284: 0xffffffff, 2288: 16777217}
						p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestEffects53F4F0}, {Op: legacy.PortTestEffects53F4F0}}
						specs = append(specs, s)
					}
				}
			}
		}
	}
	callbackHash(t, "effects-wand-acceptance", legacy.PortTestRoam(specs), effectsUseHashes["effects-wand-acceptance"])
}
func TestEffectsUseDispatch(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, blocked := range []bool{false, true} {
		for _, nilItem := range []bool{false, true} {
			for _, nilUse := range []bool{false, true} {
				for _, deletes := range []bool{false, true} {
					s := effectsUseBase()
					p := s.Callbacks.Shop
					p.Inventory.Blocked = blocked
					p.Inventory.NilItem = nilItem
					p.EffectsUse.NilUse = nilUse
					p.Inventory.UseDelete = deletes
					p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestEffects53F8E0}}
					specs = append(specs, s)
				}
			}
		}
	}
	callbackHash(t, "effects-use-dispatch", legacy.PortTestRoam(specs), effectsUseHashes["effects-use-dispatch"])
}

func TestEffectsStatusAndForce(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, op := range []int{legacy.PortTestEffects4E04D0, legacy.PortTestEffects4E0640, legacy.PortTestEffects4E0670, legacy.PortTestEffects4E06F0, legacy.PortTestEffects4E0850} {
		for _, target := range []int{0, 1} {
			for _, value := range []uint32{0, 1, 3, 255} {
				for _, flags := range []uint32{0, 0x20, 0x8000} {
					s := effectsUseBase()
					p := s.Callbacks.Shop
					p.EffectsUse.Target = target
					p.EffectsUse.Modifiers[0].Words = map[int]uint32{56: math.Float32bits(float32(value)), 60: value, 72: value}
					p.EffectsUse.Balance = map[string][]float64{"StunEnchantDuration": {90}, "ConfuseEnchantDuration": {75}}
					if target == 1 {
						p.Resources.Flags = flags
					} else {
						s.Spells.TargetFlags = flags
					}
					p.Sequence = []legacy.PortTestShopAction{{Op: op}}
					specs = append(specs, s)
				}
			}
		}
	}
	callbackHash(t, "effects-status-force", legacy.PortTestRoam(specs), effectsUseHashes["effects-status-force"])
}
func TestEffectsProjectiles(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, dir := range []uint32{0, 1, 7, 8, 63, 64, 127, 128, 247, 248, 255} {
		for _, triple := range []byte{0, 1} {
			for _, charge := range [][2]byte{{0, 0}, {0, 10}, {1, 10}, {3, 10}} {
				for _, blocked := range []bool{false, true} {
					s := effectsUseBase()
					p := s.Callbacks.Shop
					s.Owner.Frame = 100
					p.Inventory.Position.X = 512
					p.Inventory.Position.Y = 512
					p.EffectsUse.Projectiles = true
					p.EffectsUse.ProjectileSpeed = math.Float32bits(2.5)
					p.EffectsUse.UnitWords = map[int]uint32{124: dir, 176: math.Float32bits(5), 80: math.Float32bits(.25), 84: math.Float32bits(-.5)}
					p.Items[0].Use[100] = 10
					p.Items[0].Use[88] = 200
					p.Items[0].Use[96] = triple
					p.Items[0].Use[108], p.Items[0].Use[109] = charge[0], charge[1]
					if blocked {
						p.Inventory.WallMode = 1
					}
					p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestEffects53F290}, {Op: legacy.PortTestEffects53F290}}
					specs = append(specs, s)
				}
			}
		}
	}
	for _, dir := range []int{0, 1, 63, 64, 127, 128, 255} {
		for _, speed := range []float32{0, .125, 2.5, 100} {
			s := effectsUseBase()
			p := s.Callbacks.Shop
			p.EffectsUse.Projectiles = true
			p.EffectsUse.ProjectileSpeed = math.Float32bits(speed)
			p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestEffects53F480, Value: 0xfffffffe, Side: dir}}
			specs = append(specs, s)
		}
	}
	callbackHash(t, "effects-projectiles", legacy.PortTestRoam(specs), effectsUseHashes["effects-projectiles"])
}
func TestEffectsFireWand(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, dir := range []uint32{0, 1, 63, 64, 127, 128, 255} {
		for _, frame := range []uint32{0, 1, 30, 31, 100, 0xffffffff} {
			for _, disabled := range []bool{false, true} {
				s := effectsUseBase()
				p := s.Callbacks.Shop
				s.Owner.Frame = frame
				s.Owner.FPS = 30
				p.EffectsUse.Projectiles = true
				p.EffectsUse.DisableProjectiles = disabled
				p.EffectsUse.UnitWords = map[int]uint32{124: dir, 176: math.Float32bits(5)}
				p.EffectsUse.ItemWords = []map[int]uint32{{136: 0}}
				p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestEffects53F670}, {Op: legacy.PortTestEffects53F670}}
				specs = append(specs, s)
			}
		}
	}
	callbackHash(t, "effects-fire-wand", legacy.PortTestRoam(specs), effectsUseHashes["effects-fire-wand"])
}

func TestEffectsInventoryLookup(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for item := 0; item < 3; item++ {
		for slot := 0; slot < 4; slot++ {
			for _, class := range []uint32{0, 0x1000, 0x1000000, 0x2000000, 0x10000000} {
				for _, mask := range []uint32{0, 1, 2, 4, 8, 16, 32, 3, 255, 256} {
					s := effectsUseBase()
					p := s.Callbacks.Shop
					p.Items[item].Class = class
					p.Items[item].Mods[slot] = true
					p.EffectsUse.Modifiers[slot].Engage = effectsIdentity(legacy.PortTestEffects4DFD10)
					p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestEffects4DFBB0, Value: mask}}
					specs = append(specs, s)
				}
			}
		}
	}
	r := legacy.PortTestRoam(specs)
	for i, v := range r {
		sp := specs[i].Callbacks.Shop
		allowed := sp.Sequence[0].Value == 1
		for _, it := range sp.Items {
			for _, yes := range it.Mods {
				if yes {
					allowed = allowed && it.Class&0x13001000 != 0
				}
			}
		}
		if (v.Callbacks.Shop.Sequence[0].Return != 0) != allowed {
			t.Fatalf("inventory lookup %d", i)
		}
	}
	callbackHash(t, "effects-inventory-lookup", r, effectsUseHashes["effects-inventory-lookup"])
}
func TestEffectsSpeedRounding(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, speed := range []uint32{0, 0x80000000, 1, 0x3dcccccd, 0x3f800000, 0x4b800000, 0xcb800000, 0x7f800000, 0x7fc12345} {
		for _, delta := range []uint32{0, 0x80000000, 1, 0x3dcccccd, 0x3f000000, 0x3f800000, 0xbf800000} {
			s := effectsUseBase()
			p := s.Callbacks.Shop
			p.EffectsUse.UnitWords = map[int]uint32{552: speed}
			p.EffectsUse.Modifiers[0].Words = map[int]uint32{120: delta}
			p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestEffects4DFC30}, {Op: legacy.PortTestEffects4DFCA0}}
			specs = append(specs, s)
		}
	}
	callbackHash(t, "effects-speed-rounding", legacy.PortTestRoam(specs), effectsUseHashes["effects-speed-rounding"])
}
func TestEffectsWandClasses(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, class := range []uint32{0x1000, 0x1000000} {
		for _, sub := range []uint32{0, 0x40000, 0x4000000, 0x4040000} {
			for _, fn := range []int{0, 42} {
				for _, rate := range []float64{0, 2.9, -1, 2147483648} {
					s := effectsUseBase()
					p := s.Callbacks.Shop
					s.Owner.Frame = 100
					p.Items[0].Class = class
					p.Items[0].Subclass = sub
					p.Items[0].Use[108] = 2
					p.Items[0].Use[109] = 3
					p.Items[0].Use[100] = 10
					p.Items[0].Use[92] = 5
					p.EffectsUse.AcceptSpell = true
					p.EffectsUse.Balance = map[string][]float64{"OblivionStaffRechargeRate": {rate}}
					p.Items[0].Mods[2] = true
					p.EffectsUse.Modifiers[2] = legacy.PortTestEffectsModifier{Attack: fn, Words: map[int]uint32{48: 6}}
					p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestEffects53C940}, {Op: legacy.PortTestEffects53F4F0}, {Op: legacy.PortTestEffects53F4F0}}
					specs = append(specs, s)
				}
			}
		}
	}
	callbackHash(t, "effects-wand-classes", legacy.PortTestRoam(specs), effectsUseHashes["effects-wand-classes"])
}
func TestEffectsNullGuards(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, op := range []int{legacy.PortTestEffects4DFBB0, legacy.PortTestEffects4DFC30, legacy.PortTestEffects4DFCA0, legacy.PortTestEffects4DFD40, legacy.PortTestEffects4DFE40, legacy.PortTestEffects4DFF40, legacy.PortTestEffects4E0040, legacy.PortTestEffects4E0170, legacy.PortTestEffects4E0740, legacy.PortTestEffects4E07C0, legacy.PortTestEffects4E08E0} {
		s := effectsUseBase()
		p := s.Callbacks.Shop
		p.EffectsUse.NilUnit = true
		p.Sequence = []legacy.PortTestShopAction{{Op: op}}
		specs = append(specs, s)
	}
	for _, op := range []int{legacy.PortTestEffects4E01D0, legacy.PortTestEffects4E02C0, legacy.PortTestEffects4E03F0, legacy.PortTestEffects4E0640, legacy.PortTestEffects53C940, legacy.PortTestEffects53F8E0} {
		s := effectsUseBase()
		p := s.Callbacks.Shop
		p.Inventory.NilItem = true
		p.Sequence = []legacy.PortTestShopAction{{Op: op}}
		specs = append(specs, s)
	}
	callbackHash(t, "effects-null-guards", legacy.PortTestRoam(specs), effectsUseHashes["effects-null-guards"])
}

func TestEffectsFireballWall(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, wall := range []int{0, 1, 2} {
		s := effectsUseBase()
		p := s.Callbacks.Shop
		s.Owner.Frame = 100
		p.EffectsUse.Projectiles = true
		p.EffectsUse.ProjectileSpeed = math.Float32bits(2)
		p.EffectsUse.UnitWords = map[int]uint32{124: 0, 176: math.Float32bits(96), 80: 0, 84: 0}
		p.Inventory.WallMode = wall
		p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestEffects53F290}}
		specs = append(specs, s)
	}
	r := legacy.PortTestRoam(specs)
	for i, v := range r {
		if len(v.Lifecycle.Created) != 1 {
			t.Fatalf("wall %d projectile count", i)
		}
		x, y := v.Lifecycle.Created[0][14], v.Lifecycle.Created[0][15]
		want := float32(200)
		if i == 1 {
			want = 100
		}
		if x != math.Float32bits(want) || y != math.Float32bits(100) {
			t.Fatalf("wall %d origin: %08x,%08x", i, x, y)
		}
	}
	callbackHash(t, "effects-fireball-wall", r, effectsUseHashes["effects-fireball-wall"])
}
