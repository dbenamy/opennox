//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"strings"
	"testing"
)

// Original registered C callbacks through the existing projectile death owner.
var deathRegistryHashes = map[string]string{
	"death-registry-barrel":             "f653f0db4697ac6ab05ca5dd8faf122136cd78f4cb00628f7ddf25fdd595e4f6",
	"death-registry-boulder":            "0bb0f6f1c8dac6a7f5a2eed700a7068bac8ed35696fd08bced2a7afd22de593c",
	"death-registry-creation-callbacks": "0b77c6bc0194e0333533bc8c21f90687f9aac9a9dc735798488e038cbba1a3ec",
	"death-registry-creation":           "693075fcb3d25e890cea776bf42f59dc241ba78a295e914c0e28f8bf81cd6014",
	"death-registry-equipment":          "4e556710ad551b9f3b44e9047c128df0682493751438c802f2c29938291ac06c",
	"death-registry-generator":          "6a6e9e50caf845500e46fc31e8ecf1a7affe3e6e80ba3ae54b8685d800838eb1",
	"death-registry-marker":             "ab259edb80430608b6c6548825c501b3b88a9b881edeb49f59f3ba55abc4c113",
	"death-registry-smoke":              "4548a2a4b345000ba344c24e7ad412df8fed192b8977589781acc6cb995ad385",
}

func deathRegistryBase(op int) legacy.PortTestRoamSpec {
	s := deathBase(op)
	s.Callbacks.DeathRegistry = []string{"BarrelDie", "CreateObjectDie", "SpawnObjectDie", "MarkerDie", "BoulderDie", "GameBallDie", "ArmorDie", "WeaponDie"}[op]
	return s
}

func TestDeathRegistryIdentities(t *testing.T) {
	names := legacy.PortTestDeathRegistryNames()
	if len(names) != 14 {
		t.Fatalf("registry names: %v", names)
	}
	seen := make(map[string]bool)
	for _, name := range names {
		if seen[name] {
			t.Fatalf("duplicate name %s", name)
		}
		seen[name] = true
	}
}

func TestDeathRegistrySmoke(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for op := 0; op < 8; op++ {
		specs = append(specs, deathRegistryBase(op))
	}
	callbackHash(t, "death-registry-smoke", legacy.PortTestRoam(specs), deathRegistryHashes["death-registry-smoke"])
}

func TestDeathRegistryCreationContracts(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, op := range []int{1, 2} {
		for _, name := range []string{"", "PortTestDebrisA", "PortTestDebrisB", "missing", strings.Repeat("x", 127)} {
			for _, sound := range []uint32{0, 286, 757, 0xffffffff} {
				for _, flags := range []uint32{0, 4, 0x8000, 0xffff7fff, 0xffffffff} {
					s := deathRegistryBase(op)
					s.Callbacks.Death.Name = name
					s.Callbacks.Death.Sound = sound
					s.Callbacks.Death.Flags = flags
					specs = append(specs, s)
				}
			}
		}
	}
	callbackHash(t, "death-registry-creation", legacy.PortTestRoam(specs), deathRegistryHashes["death-registry-creation"])
}
func TestDeathRegistryMarkerContracts(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, owner := range []bool{false, true} {
		for markers := uint32(0); markers < 16; markers++ {
			s := deathRegistryBase(3)
			s.Callbacks.Death.Owner = owner
			s.Callbacks.Death.Markers = markers
			specs = append(specs, s)
		}
	}
	r := legacy.PortTestRoam(specs)
	for i, v := range r {
		mask := specs[i].Callbacks.Death.Markers
		if specs[i].Callbacks.Death.Owner && mask != 0 {
			mask &= mask - 1
		}
		for j := 0; j < 4; j++ {
			want := uint32(0)
			if mask&(1<<j) != 0 {
				want = 101
			}
			if got := v.Callbacks.Death.OwnerData[29+j]; got != want {
				t.Fatalf("case %d marker %d=%d want %d", i, j, got, want)
			}
		}
	}
	callbackHash(t, "death-registry-marker", r, deathRegistryHashes["death-registry-marker"])
}
func TestDeathRegistryBarrelContracts(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, barrel := range []bool{false, true} {
		for _, cache := range []uint32{0, 13, 23} {
			for _, count := range []uint32{0, 1, 3} {
				for _, threshold := range []uint32{0, 1, 50, 99, 100} {
					for _, seed := range []int{1, 7, 31, 63} {
						s := deathRegistryBase(0)
						s.Seed = seed
						d := s.Callbacks.Death
						d.Barrel = barrel
						d.Cache = cache
						d.DropCount = count
						d.DropThreshold = threshold
						specs = append(specs, s)
					}
				}
			}
		}
	}
	callbackHash(t, "death-registry-barrel", legacy.PortTestRoam(specs), deathRegistryHashes["death-registry-barrel"])
}
func TestDeathRegistryBoulderContracts(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for seed := 1; seed <= 64; seed++ {
		for rotation := uint32(0); rotation < 2; rotation++ {
			for mode := 0; mode < 3; mode++ {
				s := deathRegistryBase(4)
				s.Seed = seed
				s.Callbacks.Death.Rotation = rotation
				if mode == 1 {
					s.Callbacks.Enabled = map[string]bool{"porttestdebrisa": false}
				}
				if mode == 2 {
					s.Callbacks.Enabled = map[string]bool{"porttestdebrisb": false}
				}
				specs = append(specs, s)
			}
		}
	}
	r := legacy.PortTestRoam(specs)
	for i, v := range r {
		if v.Other != specs[i].Seed+1 {
			t.Fatalf("case %d consumed Other RNG", i)
		}
	}
	callbackHash(t, "death-registry-boulder", r, deathRegistryHashes["death-registry-boulder"])
}

func TestDeathRegistryEquipmentContracts(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, op := range []int{6, 7} {
		for _, language := range []int{0, 1, 2, 7} {
			for _, desc := range []string{"Boots", "BOOTS", "Boot", "s", "Glove\u00e9"} {
				for _, material := range []uint32{0, 2, 4, 8, 16, 30, 12, 6, 0x8000} {
					for _, owner := range []bool{false, true} {
						s := deathRegistryBase(op)
						d := s.Callbacks.Death
						d.Language = language
						d.Description = desc
						d.Material = material
						d.Owner = owner
						specs = append(specs, s)
					}
				}
			}
		}
	}
	callbackHash(t, "death-registry-equipment", legacy.PortTestRoam(specs), deathRegistryHashes["death-registry-equipment"])
}

func TestDeathRegistryCreationCallbacks(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, op := range []int{8, 9, 10} {
		for _, flags := range []uint32{0, 4, 0x20, 0x40, 0x8000, 0x8044} {
			for _, disabled := range []bool{false, true} {
				for _, fps := range []uint32{1, 30, 60} {
					s := creationBase(op)
					s.Callbacks.DeathRegistry = []string{"ImpEggDie", "PolypDie", "PotionDie"}[op-8]
					s.Callbacks.Creation.Flags = flags
					s.Callbacks.Creation.Disabled = map[string]bool{"ToxicCloud": disabled}
					s.Owner.FPS = fps
					cases = append(cases, s)
				}
			}
		}
	}
	callbackHash(t, "death-registry-creation-callbacks", legacy.PortTestRoam(cases), deathRegistryHashes["death-registry-creation-callbacks"])
}

func TestDeathRegistryGenerator(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, quest := range []uint32{0, 4096} {
		for _, killer := range []bool{false, true} {
			for _, state := range []uint32{0, 1, 2, 0xffffffff} {
				for _, flags := range []uint32{0, 4, 0x24} {
					s := generatorBase(4)
					s.Callbacks.DeathRegistry = "MonsterGeneratorDie"
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
	callbackHash(t, "death-registry-generator", legacy.PortTestRoam(specs), deathRegistryHashes["death-registry-generator"])
}
