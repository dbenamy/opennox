//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"math"
	"testing"
)

var effectsTimingHashes = map[string]string{
	"effects-timing-fire-wand":             "3127e652a8c819e68f63ae1f225b25f5caee4b0f5034cb68ba406d642bada177",
	"effects-timing-fireball-wall":         "990131db56f4194c120d84b3a889d15705eee09b2c2ca2956ca8a75d79d15e0b",
	"effects-timing-positive-regeneration": "964c39d309ff992f6923b8be373ab7acadaf3d640681b8a2f4160582d4349089",
	"effects-timing-projectiles":           "a054fcaef1c3e2c1931711f69fcaa9c21363793eed0b79aafe591bd834c6c86b",
	"effects-timing-regeneration":          "dabf6aa5da32d052253b97f42d04e516f72ad5be957e50a9c111fe2502ffeb02",
	"effects-timing-replenishment":         "ba4d30ed8d9440095b7d8ccb4cea2144d83515253c856d9aee0fda992eb30d7e",
	"effects-timing-wand-acceptance":       "2888bef90bbb9fd09ebc661bc0a6d3df064ea68f51a6a9657bac814864e6aa20",
	"effects-timing-wand-classes":          "9d4896afda911209988db7c08da1535ebc28065d412d30ad861a18ba05e0c357",
}

func TestEffectsTimingRegeneration(t *testing.T) {
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
	callbackHash(t, "effects-timing-regeneration", effectsTimedRun(t, specs), effectsTimingHashes["effects-timing-regeneration"])
}

func TestEffectsTimingReplenishment(t *testing.T) {
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
	callbackHash(t, "effects-timing-replenishment", effectsTimedRun(t, specs), effectsTimingHashes["effects-timing-replenishment"])
}

func TestEffectsTimingWandAcceptance(t *testing.T) {
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
	callbackHash(t, "effects-timing-wand-acceptance", effectsTimedRun(t, specs), effectsTimingHashes["effects-timing-wand-acceptance"])
}

func TestEffectsTimingProjectiles(t *testing.T) {
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
	callbackHash(t, "effects-timing-projectiles", effectsTimedRun(t, specs), effectsTimingHashes["effects-timing-projectiles"])
}

func TestEffectsTimingFireWand(t *testing.T) {
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
	callbackHash(t, "effects-timing-fire-wand", effectsTimedRun(t, specs), effectsTimingHashes["effects-timing-fire-wand"])
}

func TestEffectsTimingWandClasses(t *testing.T) {
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
	callbackHash(t, "effects-timing-wand-classes", effectsTimedRun(t, specs), effectsTimingHashes["effects-timing-wand-classes"])
}

func TestEffectsTimingFireballWall(t *testing.T) {
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
	r := effectsTimedRun(t, specs)
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
	callbackHash(t, "effects-timing-fireball-wall", r, effectsTimingHashes["effects-timing-fireball-wall"])
}

func effectsTimedRun(t *testing.T, specs []legacy.PortTestRoamSpec) []legacy.PortTestRoamResult {
	t.Helper()
	for i := range specs {
		s := &specs[i]
		s.MonsterState.Frame = s.Owner.Frame
		s.MonsterState.FPS = s.Owner.FPS
		s.Callbacks.Shop.EffectsUse.ExpectedClock = &[2]uint32{s.Owner.Frame, uint32(s.Owner.FPS)}
	}
	results := legacy.PortTestRoam(specs)
	for i, r := range results {
		for j, step := range r.Callbacks.Shop.Sequence {
			d := step.EffectsUseData
			want := specs[i].Callbacks.Shop.EffectsUse.ExpectedClock
			if len(d) < 2 || d[len(d)-2] != want[0] || d[len(d)-1] != want[1] {
				t.Fatalf("case %d step %d clock", i, j)
			}
		}
	}
	return results
}
func TestEffectsTimingPositiveRegeneration(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, frame := range []uint32{89, 90, 91, 179, 180, 181} {
		s := effectsUseBase()
		s.Owner.Frame = frame
		s.Owner.FPS = 30
		p := s.Callbacks.Shop
		p.Resources.HP = 50
		p.Resources.MaxHP = 100
		p.Resources.Flags = 0
		p.EffectsUse.UnitWords = map[int]uint32{536: 0}
		p.EffectsUse.Modifiers[0].Words = map[int]uint32{108: 300}
		p.Sequence = []legacy.PortTestShopAction{{Op: legacy.PortTestEffects4E01D0}, {Op: legacy.PortTestEffects4E01D0}}
		specs = append(specs, s)
	}
	results := effectsTimedRun(t, specs)
	for i, r := range results {
		for j, step := range r.Callbacks.Shop.Sequence {
			want := uint32(50)
			if i == 1 || i == 4 {
				want += uint32(j + 1)
			}
			if got := step.ResourceData[2][0] & 0xffff; got != want {
				t.Fatalf("case %d step %d HP %d want %d", i, j, got, want)
			}
		}
	}
	callbackHash(t, "effects-timing-positive-regeneration", results, effectsTimingHashes["effects-timing-positive-regeneration"])
}
