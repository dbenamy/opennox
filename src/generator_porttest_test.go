//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"math"
	"testing"
)

func TestGeneratorTileFixture(t *testing.T) {
	n, intact := legacy.PortTestGeneratorTileContracts()
	if n != 128 || !intact {
		t.Fatalf("tile fixture cases=%d intact=%v", n, intact)
	}
}

func generatorBase(op int) legacy.PortTestRoamSpec {
	s := callbackBase(59 + op)
	s.Callbacks.Generator = &legacy.PortTestGeneratorSpec{Op: op, Flags: 2, Frame: 128, Sources: 1, SourceSpellWord: 0x12345678, DefHealth: 100, PlayerFlags: 4, PlayerPos: [2]uint32{math.Float32bits(150), math.Float32bits(100)}, View: [2]uint16{200, 200}, HealthScale: math.Float32bits(1.5), Balance: map[string]float64{"MaxOnscreenMonsterCount": 100}, Radius: math.Float32bits(45), Point: [2]uint32{math.Float32bits(140), math.Float32bits(100)}}
	return s
}
func TestGeneratorPlacementSmoke(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for op := 0; op < 3; op++ {
		specs = append(specs, generatorBase(op))
	}
	callbackHash(t, "generator-placement-smoke", legacy.PortTestRoam(specs), generatorHashes["generator-placement-smoke"])
}

func TestGeneratorPlacementContracts(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, op := range []int{1, 2} {
		for _, seed := range []int{1, 2, 3, 7, 15, 31, 63, 127} {
			for _, tile := range []int{0, 1, 5, 175} {
				for wall := 0; wall < 3; wall++ {
					for mode := 0; mode < 4; mode++ {
						s := generatorBase(op)
						s.Seed = seed
						s.Combat.Wall = wall
						g := s.Callbacks.Generator
						g.Tile = tile
						g.Cold = mode&1 != 0
						g.TowardPlayer = mode&2 != 0
						if mode == 3 {
							g.BlockerRadius = math.Float32bits(50)
							g.BlockerAtSource = true
						}
						if seed&1 == 0 {
							s.Spells.TargetFlags |= 0x4000
						}
						specs = append(specs, s)
					}
				}
			}
		}
	}
	r := legacy.PortTestRoam(specs)
	for i, v := range r {
		if specs[i].Callbacks.Generator.BlockerAtSource && v.Callbacks.Return != 0 {
			t.Fatalf("case %d escaped full blocker", i)
		}
		if v.Other != specs[i].Seed+1 {
			t.Fatalf("case %d consumed Other RNG", i)
		}
	}
	callbackHash(t, "generator-placement", r, generatorHashes["generator-placement"])
}
func TestGeneratorVacancyContracts(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, radius := range []float32{0, 1, 5, 15, 50} {
		for _, dx := range []float32{-16, -15, -5, -1, 0, 1, 5, 15, 16} {
			s := generatorBase(0)
			g := s.Callbacks.Generator
			g.BlockerRadius = math.Float32bits(radius)
			g.Point = [2]uint32{math.Float32bits(140 + dx), math.Float32bits(100)}
			specs = append(specs, s)
		}
	}
	callbackHash(t, "generator-vacancy", legacy.PortTestRoam(specs), generatorHashes["generator-vacancy"])
}

func TestGeneratorUpdateGates(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, flags := range []uint32{0, 4, 0x1000004, 0x1000024, 0x1008004} {
		for _, stage := range []uint32{0, 9, 10, 11, 100, 0xffffffff} {
			for rate := byte(0); rate < 5; rate++ {
				for mode := 0; mode < 4; mode++ {
					s := generatorBase(3)
					g := s.Callbacks.Generator
					g.GenObjFlags = flags
					g.Stage = stage
					g.RateClass = rate
					g.Level = uint32(rate % 3)
					g.Cache = [8]uint32{10, 2, 30, 60, 90, 120, 150, math.Float32bits(.5)}
					g.Balance = map[string]float64{"QuestHardcoreStage": 10, "QuestHardcoreSpawnRateIncrease": 2, "QuestHardcoreSpawnCap": .5, "SpawnRateHighValue": 30, "SpawnRateNormalValue": 60, "SpawnRateLowValue": 90, "SpawnRateVeryLowValue": 120, "SpawnRateVeryVeryLowValue": 150}
					if mode&1 != 0 {
						g.Cache = [8]uint32{}
					}
					if mode&2 != 0 {
						g.GenXStatus = 0x800
					}
					s.Owner.Frame = []uint32{0, 7, 8, 128, 0xffffffff}[rate]
					g.Frame = s.Owner.Frame
					g.LastSpawn = []uint32{0, 1, 7, 128, 0xffffffff}[rate]
					// No source creature is selected while the live count is at its cap.
					g.Current = 0
					g.Limit = 0
					specs = append(specs, s)
				}
			}
		}
	}
	callbackHash(t, "generator-update-gates", legacy.PortTestRoam(specs), generatorHashes["generator-update-gates"])
}

func TestGeneratorObjectsSmoke(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for op := 4; op < 8; op++ {
		s := generatorBase(op)
		specs = append(specs, s)
	}
	callbackHash(t, "generator-objects-smoke", legacy.PortTestRoam(specs), generatorHashes["generator-objects-smoke"])
}

func TestGeneratorDeathContracts(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, quest := range []uint32{0, 4096} {
		for _, killer := range []bool{false, true} {
			for _, state := range []uint32{0, 1, 2, 0xffffffff} {
				for _, flags := range []uint32{0, 4, 0x24} {
					s := generatorBase(4)
					g := s.Callbacks.Generator
					g.Killer = killer
					g.QuestState = state
					g.PlayerFlags = flags
					s.Lifecycle.GameFlags = quest
					specs = append(specs, s)
				}
			}
		}
	}
	callbackHash(t, "generator-death", legacy.PortTestRoam(specs), generatorHashes["generator-death"])
}
func TestGeneratorPickContracts(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, flags := range []uint32{0, 1, 2, 3, 4} {
		for _, distance := range []float32{0, 44, 45, 200, 300, 301, 1000} {
			for mode := 0; mode < 8; mode++ {
				s := generatorBase(5)
				g := s.Callbacks.Generator
				g.Flags = flags
				g.PlayerPos = [2]uint32{math.Float32bits(100 + distance), math.Float32bits(100)}
				if mode&1 != 0 {
					g.PlayerStatus = 1
				}
				if mode&2 != 0 {
					g.PlayerFlags |= 0x20
				}
				if mode&4 != 0 {
					s.Combat.Wall = 1
				}
				specs = append(specs, s)
			}
		}
	}
	callbackHash(t, "generator-pick", legacy.PortTestRoam(specs), generatorHashes["generator-pick"])
}
func TestGeneratorSpawnContracts(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, scale := range []float32{0, 0.000001, .5, 1, 1.5, -1, 655.36, float32(math.Inf(1)), float32(math.NaN())} {
		for _, beholder := range []bool{false, true} {
			for mode := 0; mode < 4; mode++ {
				s := generatorBase(6)
				g := s.Callbacks.Generator
				g.HealthScale = math.Float32bits(scale)
				g.Beholder = beholder
				g.Joined = mode&1 != 0
				if mode&2 != 0 {
					g.Balance["MaxOnscreenMonsterCount"] = 0
				}
				specs = append(specs, s)
			}
		}
	}
	r := legacy.PortTestRoam(specs)
	for i, v := range r {
		want := 1
		if i%4 == 3 {
			want = 0
		}
		if got := len(v.Lifecycle.Created); got != want {
			t.Fatalf("case %d created=%d want%d", i, got, want)
		}
	}
	callbackHash(t, "generator-spawn", r, generatorHashes["generator-spawn"])
}
func TestGeneratorCopyContracts(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for i := uint32(0); i < 64; i++ {
		s := generatorBase(7)
		s.Callbacks.Generator.SourceFill = i
		specs = append(specs, s)
	}
	r := legacy.PortTestRoam(specs)
	for i, v := range r {
		o := v.Callbacks.Generator.Objects
		for j, x := range o.SourceData {
			if o.DestinationData[j] != x {
				t.Fatalf("copy case%d word%d differs", i, j)
			}
		}
	}
	callbackHash(t, "generator-copy", r, generatorHashes["generator-copy"])
}

func TestGeneratorUpdateSpawn(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for sources := byte(1); sources <= 4; sources++ {
		for level := uint32(0); level < 3; level++ {
			for _, frame := range []uint32{7, 8, 135, 136} {
				for mode := 0; mode < 4; mode++ {
					s := generatorBase(8)
					g := s.Callbacks.Generator
					g.Sources = sources
					g.Level = level
					g.Flags = 1
					g.GenObjFlags = 0x1000004
					g.Limit = 1
					s.Owner.Frame = frame
					g.Frame = frame
					g.Cache = [8]uint32{10, 2, 30, 60, 90, 120, 150, math.Float32bits(.5)}
					if mode&1 != 0 {
						g.Current = 1
					}
					if mode&2 != 0 {
						g.Joined = true
						g.Balance["MaxOnscreenMonsterCount"] = 0
					}
					specs = append(specs, s)
				}
			}
		}
	}
	r := legacy.PortTestRoam(specs)
	spawned := 0
	for i, v := range r {
		if len(v.Lifecycle.Created) > 0 {
			spawned++
			if specs[i].Owner.Frame != 136 || specs[i].Callbacks.Generator.Current != 0 || specs[i].Callbacks.Generator.Joined {
				t.Fatalf("case%d unexpected spawn", i)
			}
		}
	}
	if spawned != 12 {
		t.Fatalf("integrated spawns=%d want12", spawned)
	}
	callbackHash(t, "generator-update-spawn", r, generatorHashes["generator-update-spawn"])
}
func TestGeneratorSpawnDefinitionHealth(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, health := range []uint32{0, 1, 100, 65535, 0x80000000, 0xffffffff} {
		for _, scale := range []float32{0, .5, 1, 1.5, -1} {
			for _, beholder := range []bool{false, true} {
				s := generatorBase(6)
				g := s.Callbacks.Generator
				g.UseDef = true
				g.DefHealth = health
				g.HealthScale = math.Float32bits(scale)
				g.Beholder = beholder
				specs = append(specs, s)
			}
		}
	}
	callbackHash(t, "generator-definition-health", legacy.PortTestRoam(specs), generatorHashes["generator-definition-health"])
}

func TestGeneratorCopyInventory(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for count := 0; count <= 3; count++ {
		for _, subclass := range []uint32{0, 1, 0x10, 0x11} {
			for mask := 0; mask < 8; mask++ {
				s := generatorBase(7)
				inv := &legacy.PortTestGeneratorInventorySpec{Count: count, Subclass: subclass}
				for i := 0; i < 3; i++ {
					inv.Equipped[i] = mask&(1<<i) != 0
				}
				s.Callbacks.Generator.Inventory = inv
				specs = append(specs, s)
			}
		}
	}
	r := legacy.PortTestRoam(specs)
	for i, v := range r {
		inv := v.Callbacks.Generator.Objects.Inventory
		sp := specs[i].Callbacks.Generator.Inventory
		if sp.Subclass&0x10 == 0 {
			continue
		}
		for j, init := range inv.Init {
			sourceIndex := sp.Count - 1 - j
			if init[0] != 970 || init[4] != 0x12340000+uint32(sourceIndex) {
				t.Fatalf("case %d clone %d modifier copy %v", i, j, init)
			}
			item := inv.Items[sp.Count+j]
			wantEquipped := sp.Equipped[sourceIndex]
			if sourceIndex == 0 && sp.Count == 3 && sp.Equipped[2] {
				wantEquipped = false
			}
			if (item[4]&0x100 != 0) != wantEquipped {
				t.Fatalf("case %d clone %d equip flag %#x", i, j, item[4])
			}
		}
	}
	callbackHash(t, "generator-inventory", r, generatorHashes["generator-inventory"])
}

var generatorHashes = map[string]string{
	"generator-copy":              "01a36456409f301f758261c405ef130659c1ddef2a545ecb168844e510463482",
	"generator-death":             "675ffa6eb8f0c25655b2c1dcde61b4579c7f6e09648d9816f9e484e2868e8b68",
	"generator-definition-health": "fd1d25c857c0861f94d75af6f8fa6f1f3df65e7510efcfaa59a2b45698ac4beb",
	"generator-inventory":         "61413002da7df898fc791e1c5a5dc3e6a65040cdaf8c85b8747df566a984c74d",
	"generator-objects-smoke":     "965d7ee18be6e3938cba44c841343c6af2ca2c8ba36a3fd76b59faf56831a3ad",
	"generator-pick":              "118dc7f1d15e4cf5a545cd398e1735b43d28490816973eb24356f4d66b38597e",
	"generator-placement-smoke":   "2f29897fc7e259740439fe83b74769fe77d86ed21a3cdcd10966c8cbbdc00790",
	"generator-placement":         "75b3a073171b0d732bb6552627766398b744fddf7b527b30a45f0a149d49381f",
	"generator-spawn":             "8226c9ed25d47e605ec38052e5bb11fd885139036e516af79d0dd2d1ad8c4eb6",
	"generator-update-gates":      "a6b5834d316d9fc03f87306494c2f9614b50786e23eb48437fb1eaa8b88b40f5",
	"generator-update-spawn":      "c98e6cde0a9821bf35660bc2543b00e0044dfe7eb614b44bc76afdaec90e69af",
	"generator-vacancy":           "425f42651bf578fedb9ca3f18f73827b5e516b32153328ae0e7e8a596c35e943",
}
