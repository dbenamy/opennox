//go:build porttest

package opennox

import (
	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"math"
	"testing"
)

var attackHashes = map[string]string{
	"attack-abilities":         "fc930c5ba5e25a15fc41cc14a3f3ea8fa20227e2f86a50e5d3193412f5b2a2e9",
	"attack-chakram-depletion": "d7e423bc8d121223789c797a5505aa744541a3238560e8c4373b4f443a74042f",
	"attack-effects":           "0249d643248744f0cf8ab48aa3aa448ef3453a91ead0398d599df25c7506b4bc",
	"attack-hit-filters":       "ef79a01d0bd1f6464892da9ea1a7a671498ca282690ec8b1a0e28f9b08aaffa4",
	"attack-nearest":           "a1c35ff6760f0ad7a72077d8f2425ea0d2a7de46b8e3ecdab722f01a7a123e11",
	"attack-npc-shots":         "863dabcf5849eb026636ed4f9462b533eeb07ef14e49551d7d7f1cc688825f05",
	"attack-npc":               "98069714d8c428f91ce7b67200b972a0fe827f92ba0cef7d918de1bbc25f9db2",
	"attack-positive":          "39a8c627130c6908ded87f24df85b16c51b99390fa1145ca22961a68010e942a",
	"attack-reloads":           "757a7f7fc7bd3a46f9d19c39ebf26695b7264ac4a283688f9ca93862a8bfe108",
	"attack-shot-effects":      "a47c9df15621506ac8b5a86d4100c0d382c8e7b88b2a700743aab899e0eae94a",
	"attack-shot-failures":     "b58f7b9a443283d3f3da39e1bbb5c328ec57fc184c868faadff7f691e6b35b01",
	"attack-shots":             "5d367157a718bfc5e066a361469e1bdc2a956bfa7e19aea3c8f175670ee31d88",
	"attack-trace":             "ae515767329735c450f37a38108585b8507abcc2501b8da692bb05d6d57320bc",
	"attack-unarmed-timing":    "1cc8ebd7b2ba10cfe20ea140f2ba5f1364d55a7aba3f009ebe7959f235b526f4",
	"attack-warcry":            "5cb29d4f6f47731d8586d51bbde39cde50152ee9d1b0fa566f8333719517ced9",
	"attack-weapon-timing":     "97ba40e00bf7631eff2057818eac915ef8b08e637f19901db17e0f1044682a36",
}

func attackBase() legacy.PortTestRoamSpec {
	s := objectiveBase()
	p := s.Callbacks.Shop
	w := p.TemporaryUpdates
	o := w.World.Objectives
	o.Attack = &legacy.PortTestAttackSpec{Actor: 100, RecordWords: map[int]uint32{0: math.Float32bits(12.5), 8: math.Float32bits(50), 16: math.Float32bits(512), 20: math.Float32bits(512), 32: 1}, RecordRefs: map[int]int{12: 100, 28: 3}, ProjectileSpeed: 10, NearestRange: math.Float32bits(50)}
	w.Target = 4
	p.Inventory.Linked = []int{0}
	p.Inventory.Owned = []int{0}
	p.Inventory.WeaponBits = map[uint16]uint32{23: 0x100, 24: 2, 25: 0x100}
	p.Equipment.WeaponFlags = 0x100
	p.Equipment.ActiveWeapon = 1
	p.Equipment.Definitions[0].Words = map[int]uint32{68: math.Float32bits(30), 72: 5, 76: 0}
	p.Items[0].Class = uint32(object.ClassWeapon)
	p.Items[0].Flags = 0x100
	p.Items[0].Subclass = 0x100
	p.Items[1].Class = uint32(object.ClassSimple)
	p.Items[1].Type = 24
	p.Items[2].Class = uint32(object.ClassWeapon)
	p.Items[2].Type = 25
	w.ItemWords[1][56] = math.Float32bits(540)
	w.ItemWords[1][60] = math.Float32bits(512)
	o.PlayerWords[0][136] = 100
	p.EffectsUse.Balance["ItemDamagePercentage"] = []float64{0}
	return s
}
func attackHash(t *testing.T, name string, s []legacy.PortTestRoamSpec) {
	t.Helper()
	callbackHash(t, name, effectsTimedRun(t, s), attackHashes[name])
}
func TestAttackUnarmedTiming(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, anim := range []uint32{0, 1, 23, 24, 25, 26} {
		for _, frame := range []uint32{0, 99, 100, 101, 102, 103, 104, 105, 106, 110, 115, 0xffffffff} {
			for _, prior := range []uint32{0, 1, 3, 255} {
				s := attackBase()
				p := s.Callbacks.Shop
				s.Owner.Frame = frame
				p.Equipment.ActiveWeapon = 0
				p.Equipment.WeaponFlags = 0
				p.Inventory.Linked = nil
				p.Inventory.Owned = nil
				p.TemporaryUpdates.World.Objectives.PlayerDataWords[0] = map[int]uint32{8: anim}
				p.TemporaryUpdates.World.Objectives.PlayerUpdateWords[0] = map[int]uint32{236: prior}
				p.Sequence = []legacy.PortTestShopAction{{Op: 905}, {Op: 905}}
				cases = append(cases, s)
			}
		}
	}
	attackHash(t, "attack-unarmed-timing", cases)
}
func TestAttackHitFilters(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, flags := range []uint32{0, 1, 4, 8, 0x10, 0x20, 0x40, 0x80, 0x8000, 0x8049} {
		for _, mask := range []uint32{0, 1, 255} {
			for _, back := range []bool{false, true} {
				for _, hitStill := range []uint32{0, 1} {
					s := attackBase()
					p := s.Callbacks.Shop
					w := p.TemporaryUpdates
					p.Items[1].Flags = flags
					w.World.Objectives.Attack.RecordWords[32] = mask
					w.World.Objectives.Attack.RecordWords[24] = hitStill
					if back {
						w.ItemWords[1][56] = math.Float32bits(480)
					}
					p.Sequence = []legacy.PortTestShopAction{{Op: 902}}
					cases = append(cases, s)
				}
			}
		}
	}
	attackHash(t, "attack-hit-filters", cases)
}
func TestAttackNearest(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, x := range []float32{480, 512, 513, 540, 561, 562, 563} {
		for _, y := range []float32{512, 520, 560} {
			for _, radius := range []float32{0, 5, 32} {
				s := attackBase()
				p := s.Callbacks.Shop
				w := p.TemporaryUpdates
				w.ItemWords[1][56] = math.Float32bits(x)
				w.ItemWords[1][60] = math.Float32bits(y)
				w.ItemWords[1][176] = math.Float32bits(radius)
				p.Sequence = []legacy.PortTestShopAction{{Op: 903}}
				cases = append(cases, s)
			}
		}
	}
	attackHash(t, "attack-nearest", cases)
}
func TestAttackEffectDispatch(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, op := range []int{900, 904} {
		for mask := uint32(0); mask < 16; mask++ {
			for _, buffs := range []uint32{0, 1 << 23, 1 << 27} {
				s := attackBase()
				p := s.Callbacks.Shop
				sp := p.TemporaryUpdates.World.Objectives.Attack
				sp.AttackMask = mask
				sp.PreMask = mask
				sp.SetOutput = true
				sp.Output = math.Float32bits(25)
				p.Items[1].Class = uint32(object.ClassSimple)
				p.TemporaryUpdates.ItemWords[1][340] = buffs
				for i := 0; i < 4; i++ {
					p.Items[0].Mods[i] = true
				}
				p.Sequence = []legacy.PortTestShopAction{{Op: op}}
				cases = append(cases, s)
			}
		}
	}
	attackHash(t, "attack-effects", cases)
}
func TestAttackShots(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, kind := range []uint32{4, 8} {
		for _, ammo := range []int{0, 4} {
			for _, count := range []uint32{0, 1, 2, 255} {
				for _, infinite := range []uint32{0, 1} {
					for _, dir := range []uint32{0, 64, 128, 255} {
						s := attackBase()
						p := s.Callbacks.Shop
						o := p.TemporaryUpdates.World.Objectives
						o.Attack.Ammo = ammo
						o.PlayerWords[0][124] = dir | dir<<16
						o.UseWords = []map[int]uint32{nil, {0: 20 | count<<8 | infinite<<16}}
						p.Sequence = []legacy.PortTestShopAction{{Op: 908, Value: kind}}
						cases = append(cases, s)
					}
				}
			}
		}
	}
	attackHash(t, "attack-shots", cases)
}
func TestAttackWeaponTiming(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, bits := range []uint32{0, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024, 2048, 4096, 8192, 16384, 32768, 0x10000, 0x20000, 0x40000, 0x80000, 0x100000, 0x200000, 0x400000, 0x800000, 0x1000000, 0x2000000, 0x4000000, 0x7800000} {
		for _, delta := range []uint32{0, 1, 2, 3, 4, 5, 6, 8, 10, 15, 30, 0xffffffff} {
			for _, prior := range []uint32{0, 3} {
				s := attackBase()
				p := s.Callbacks.Shop
				w := p.TemporaryUpdates
				s.Owner.Frame = 100 + delta
				p.Equipment.WeaponFlags = bits
				p.Inventory.WeaponBits[23] = bits
				p.Items[0].Subclass = bits
				w.World.Objectives.PlayerUpdateWords[0] = map[int]uint32{236: prior}
				w.World.Objectives.UseWords = []map[int]uint32{{0: 20 | 2<<8, 108: 10 | 20<<8, 112: 50}}
				w.Indexed = []int{4}
				p.Sequence = []legacy.PortTestShopAction{{Op: 905}}
				cases = append(cases, s)
			}
		}
	}
	attackHash(t, "attack-weapon-timing", cases)
}
func TestAttackReloads(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, op := range []int{907, 910, 911} {
		for _, bits := range []uint32{4, 8, 128} {
			for _, linked := range [][]int{nil, {0}, {0, 1}, {1, 0}, {2, 1, 0}} {
				for _, state := range []uint32{0, 1, 2} {
					s := attackBase()
					p := s.Callbacks.Shop
					w := p.TemporaryUpdates
					p.Inventory.Linked = linked
					p.Inventory.Owned = linked
					p.Inventory.WeaponBits[23] = bits
					p.Inventory.WeaponBits[25] = 128
					p.Items[1].Class = uint32(object.ClassWeapon)
					p.Items[1].Flags = 0x100
					p.Items[1].Subclass = 2
					w.World.Objectives.UseWords = []map[int]uint32{{0: state}, {0: 20 | 2<<8}}
					p.Sequence = []legacy.PortTestShopAction{{Op: op}}
					cases = append(cases, s)
				}
			}
		}
	}
	attackHash(t, "attack-reloads", cases)
}
func TestAttackTrace(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, indexed := range [][]int{nil, {4}, {4, 5}} {
		for _, area := range []bool{false, true} {
			for _, material := range []uint32{0, 0x4000} {
				for _, radius := range []float32{0, 20, 50} {
					s := attackBase()
					p := s.Callbacks.Shop
					w := p.TemporaryUpdates
					w.Indexed = indexed
					w.ItemWords[1][24] = material
					w.ItemWords[2][56] = math.Float32bits(530)
					w.ItemWords[2][60] = math.Float32bits(515)
					p.Items[2].Class = uint32(object.ClassSimple)
					w.World.Objectives.Attack.RecordWords[8] = math.Float32bits(radius)
					if area {
						p.Items[0].Subclass |= 0x4000
					}
					p.Sequence = []legacy.PortTestShopAction{{Op: 901}}
					cases = append(cases, s)
				}
			}
		}
	}
	attackHash(t, "attack-trace", cases)
}
func TestAttackWarcry(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, class := range []uint32{uint32(object.ClassSimple), uint32(object.ClassMonster)} {
		for _, flags := range []uint32{0, 0x20, 0x8000, 0x8020} {
			for _, power := range []uint32{0, 1, 2} {
				s := attackBase()
				p := s.Callbacks.Shop
				w := p.TemporaryUpdates
				p.Items[1].Class = class
				p.Items[1].Flags = flags
				w.UpdateWords[1][0] = power
				p.Sequence = []legacy.PortTestShopAction{{Op: 906}}
				cases = append(cases, s)
			}
		}
	}
	attackHash(t, "attack-warcry", cases)
}
func attackItemUse(t *testing.T, r legacy.PortTestRoamResult, id uint32) []uint32 {
	t.Helper()
	for _, d := range r.Callbacks.Shop.Sequence[0].ObjectData {
		if d[0] == id {
			off := 3 + int(d[1])/4
			return d[off : off+int(d[2])/4]
		}
	}
	t.Fatalf("missing use data %d", id)
	return nil
}
func TestAttackPositive(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	s := attackBase()
	p := s.Callbacks.Shop
	p.TemporaryUpdates.Indexed = []int{4}
	p.Sequence = []legacy.PortTestShopAction{{Op: 901}}
	cases = append(cases, s)
	s = attackBase()
	p = s.Callbacks.Shop
	o := p.TemporaryUpdates.World.Objectives
	o.Attack.Ammo = 4
	o.UseWords = []map[int]uint32{nil, {0: 20 | 2<<8}}
	p.Sequence = []legacy.PortTestShopAction{{Op: 908, Value: 4}}
	cases = append(cases, s)
	s = attackBase()
	p = s.Callbacks.Shop
	p.Items[1].Class = uint32(object.ClassMonster)
	p.Items[1].Subclass = 0x20000
	p.Items[1].Flags = 4
	p.Sequence = []legacy.PortTestShopAction{{Op: 906}}
	cases = append(cases, s)
	s = attackBase()
	p = s.Callbacks.Shop
	s.Owner.Frame = 106
	p.Equipment.WeaponFlags = 64
	p.Inventory.WeaponBits[23] = 64
	p.Items[0].Subclass = 64
	p.Sequence = []legacy.PortTestShopAction{{Op: 905}}
	cases = append(cases, s)
	s = attackBase()
	p = s.Callbacks.Shop
	p.TemporaryUpdates.Indexed = []int{4}
	p.Items[0].Mods[0] = true
	p.TemporaryUpdates.World.Objectives.Attack.AttackMask = 1
	p.TemporaryUpdates.World.Objectives.Attack.SetOutput = true
	p.TemporaryUpdates.World.Objectives.Attack.Output = math.Float32bits(25)
	p.Sequence = []legacy.PortTestShopAction{{Op: 904}, {Op: 901}}
	cases = append(cases, s)
	s = attackBase()
	p = s.Callbacks.Shop
	p.Inventory.Gameplay = 1
	p.Items[0].Mods[2] = true
	p.TemporaryUpdates.World.Objectives.Attack.PreMask = 4
	p.TemporaryUpdates.World.Objectives.Attack.SetOutput = true
	p.TemporaryUpdates.World.Objectives.Attack.Output = math.Float32bits(25)
	p.Sequence = []legacy.PortTestShopAction{{Op: 902}}
	cases = append(cases, s)
	r := effectsTimedRun(t, cases)
	if d := r[0].Callbacks.Damage; len(d) != 5 || d[0] != 70001 || d[1] != 54000 || d[2] != 70000 || d[3] != 13 {
		t.Fatalf("melee damage %v", d)
	}
	if r[0].Callbacks.Shop.Sequence[0].Return != 1 {
		t.Fatal("melee hit not reported")
	}
	if len(r[1].Lifecycle.Created) != 1 {
		t.Fatal("arrow not created")
	}
	shot := r[1].Lifecycle.Created[0]
	if shot[56/4] != math.Float32bits(521) || shot[60/4] != math.Float32bits(512) || shot[80/4] != math.Float32bits(10) || shot[84/4] != 0 {
		t.Fatal("arrow trajectory")
	}
	if got := (attackItemUse(t, r[1], 70001)[0] >> 8) & 255; got != 1 {
		t.Fatalf("ammo remaining %d want 1", got)
	}
	if u := worldObject(t, r[2], 0, 70001); u[340/4]&(1<<5) == 0 || u[(344+2*5)/4]>>16 != 90 {
		t.Fatal("warcry did not apply stun duration")
	}
	if len(r[3].Lifecycle.Created) != 1 {
		t.Fatal("round chakram not created")
	}
	chakram := r[3].Lifecycle.Created[0]
	if chakram[504/4] != 70000 {
		t.Fatalf("chakram inventory=%d want weapon", chakram[504/4])
	}
	for i := 4; i < 6; i++ {
		d := r[i].Callbacks.Damage
		if len(d) != 5 || d[3] != 25 {
			t.Fatalf("modifier case %d damage=%v", i, d)
		}
	}
	callbackHash(t, "attack-positive", r, attackHashes["attack-positive"])
}
func TestAttackNPC(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, sub := range []uint32{0, 0x10} {
		for _, weapon := range []int{0, 3} {
			for _, bits := range []uint32{0, 0x100, 0x400, 4, 8, 0x10000} {
				for _, delta := range []uint32{0, 1, 2, 3, 5, 10, 20, 0xffffffff} {
					s := attackBase()
					p := s.Callbacks.Shop
					sp := p.TemporaryUpdates.World.Objectives.Attack
					sp.Actor = 1
					p.Resources.Subject = 3
					p.Resources.Subclass = sub
					s.Owner.Frame = 100 + delta
					sp.ActorWords = map[int]uint32{136: 100}
					sp.UpdateWords = map[int]uint32{2056: bits, 2068: 23}
					sp.UpdateRefs = map[int]int{2064: weapon}
					sp.RecordRefs[12] = 1
					p.Sequence = []legacy.PortTestShopAction{{Op: 905}}
					cases = append(cases, s)
				}
			}
		}
	}
	attackHash(t, "attack-npc", cases)
}
func TestAttackAbilities(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, abilities := range []uint32{2, 4, 6} {
		for _, delta := range []uint32{0, 1, 2, 3, 4, 5, 6, 7, 8, 12, 14, 15, 16, 22} {
			for _, prior := range []uint32{0, 2} {
				for _, buffs := range []uint32{0, 1 << 5, 1 << 25} {
					s := attackBase()
					p := s.Callbacks.Shop
					s.Owner.Frame = 100 + delta
					p.Resources.Buffs = buffs
					p.Equipment.ActiveAbilities = abilities
					p.TemporaryUpdates.World.Objectives.PlayerUpdateWords[0] = map[int]uint32{236: prior}
					s.Spells.DurationEmpty = true
					p.TemporaryUpdates.World.Objectives.SpellDefinitions = []server.PortTestSpellClassDef{{Index: 13, Valid: true}}
					p.EffectsUse.Balance["CounterspellRange"] = []float64{100}
					p.Sequence = []legacy.PortTestShopAction{{Op: 905}, {Op: 905}}
					cases = append(cases, s)
				}
			}
		}
	}
	attackHash(t, "attack-abilities", cases)
}
func TestAttackShotEffects(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, op := range []int{908, 909} {
		for _, mask := range []uint32{0, 4, 8, 12} {
			for _, recoil := range []uint32{0, 4, 8, 12} {
				s := attackBase()
				p := s.Callbacks.Shop
				o := p.TemporaryUpdates.World.Objectives
				o.Attack.RecoilMask = recoil
				o.Attack.Ammo = 4
				o.UseWords = []map[int]uint32{nil, {0: 20 | 2<<8}}
				for i := 2; i < 4; i++ {
					p.Items[0].Mods[i] = mask&(1<<i) != 0 || recoil&(1<<i) != 0
					p.EffectsUse.Modifiers[i].Attack = 34
					p.EffectsUse.Modifiers[i].Words = map[int]uint32{44: math.Float32bits(1.5)}
				}
				p.Sequence = []legacy.PortTestShopAction{{Op: op, Value: 4}}
				cases = append(cases, s)
			}
		}
	}
	attackHash(t, "attack-shot-effects", cases)
}
func TestAttackShotFailures(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, wall := range []int{0, 1, 2} {
		for _, missing := range []bool{false, true} {
			for _, ammo := range []int{0, 4} {
				for _, kind := range []uint32{4, 8} {
					s := attackBase()
					p := s.Callbacks.Shop
					o := p.TemporaryUpdates.World.Objectives
					p.Inventory.WallMode = wall
					o.PlayerWords[0][56] = math.Float32bits(100)
					o.PlayerWords[0][60] = math.Float32bits(100)
					o.PlayerWords[0][176] = math.Float32bits(50)
					o.Attack.Ammo = ammo
					o.UseWords = []map[int]uint32{nil, {0: 20 | 2<<8}}
					if missing {
						o.Attack.MissingTypes = []string{"ArcherArrow", "ArcherBolt", "WeakArcherArrow"}
					}
					p.Sequence = []legacy.PortTestShopAction{{Op: 908, Value: kind}}
					cases = append(cases, s)
				}
			}
		}
	}
	r := effectsTimedRun(t, cases)
	for i, x := range r {
		wall := i / 8
		missing := (i/4)%2 != 0
		ammo := (i/2)%2 != 0
		created := 0
		if wall != 1 && !missing {
			created = 1
		}
		if len(x.Lifecycle.Created) != created {
			t.Fatalf("shot failure %d creation count", i)
		}
		want := uint32(2)
		if wall != 1 && ammo {
			want = 1
		}
		if got := (attackItemUse(t, x, 70001)[0] >> 8) & 255; got != want {
			t.Fatalf("shot failure %d ammo=%d want %d", i, got, want)
		}
	}
	callbackHash(t, "attack-shot-failures", r, attackHashes["attack-shot-failures"])
}
func TestAttackChakramDepletion(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, bits := range []uint32{64, 128} {
		for _, count := range []uint32{0, 1, 2} {
			for _, reserve := range []bool{false, true} {
				for _, missing := range []bool{false, true} {
					s := attackBase()
					p := s.Callbacks.Shop
					o := p.TemporaryUpdates.World.Objectives
					s.Owner.Frame = 106
					p.Equipment.WeaponFlags = bits
					p.Inventory.WeaponBits[23] = bits
					p.Inventory.WeaponBits[25] = 128
					p.Items[0].Subclass = bits
					p.Items[2].Subclass = 128
					p.Equipment.Definitions = append(p.Equipment.Definitions, legacy.PortTestEquipmentDef{Type: 25, Strength: 20, Coeff: math.Float32bits(.25)})
					o.UseWords = []map[int]uint32{{0: 20 | count<<8}, nil, {0: 20 | 5<<8}}
					if reserve {
						p.Inventory.Linked = []int{0, 2}
						p.Inventory.Owned = []int{0, 2}
					}
					if missing {
						o.Attack.MissingTypes = []string{"RoundChakramInMotion", "FanChakramInMotion"}
					}
					p.Sequence = []legacy.PortTestShopAction{{Op: 905}}
					cases = append(cases, s)
				}
			}
		}
	}
	r := effectsTimedRun(t, cases)
	// Fan, final charge, reserve weapon, available projectile.
	u, ud, _ := objectivePlayerData(r[18].Callbacks.Shop.Sequence[0], 0)
	if u[504/4] != 70002 || ud[104/4] != 70002 {
		t.Fatalf("fan depletion inventory=%d equipped=%d", u[504/4], ud[104/4])
	}
	deleted := false
	for i := 0; i+1 < len(r[18].Trace); i++ {
		if r[18].Trace[i] == 32 && r[18].Trace[i+1] == 70000 {
			deleted = true
		}
	}
	if !deleted {
		t.Fatal("depleted fan deletion not requested")
	}
	callbackHash(t, "attack-chakram-depletion", r, attackHashes["attack-chakram-depletion"])
}
func TestAttackNPCShots(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, player := range []bool{true, false} {
		for _, ammo := range []int{0, 4} {
			for _, kind := range []uint32{4, 8} {
				s := attackBase()
				p := s.Callbacks.Shop
				o := p.TemporaryUpdates.World.Objectives
				o.Attack.Ammo = ammo
				o.UseWords = []map[int]uint32{nil, {0: 20 | 2<<8}}
				if !player {
					p.Resources.Subject = 3
					p.Resources.Subclass = 0x10
					o.Attack.Actor = 1
				}
				p.Sequence = []legacy.PortTestShopAction{{Op: 908, Value: kind}}
				cases = append(cases, s)
			}
		}
	}
	r := effectsTimedRun(t, cases)
	for i, x := range r {
		want := uint32(2)
		if i < 4 && i%4 >= 2 {
			want = 1
		}
		if got := (attackItemUse(t, x, 70001)[0] >> 8) & 255; got != want {
			t.Fatalf("actor shot %d ammo=%d want %d", i, got, want)
		}
	}
	callbackHash(t, "attack-npc-shots", r, attackHashes["attack-npc-shots"])
}
