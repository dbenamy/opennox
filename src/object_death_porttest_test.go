//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"strings"
	"testing"
)

func deathBase(op int) legacy.PortTestRoamSpec {
	s := callbackBase(51 + op)
	s.Callbacks.Death = &legacy.PortTestDeathSpec{Name: "PortTestDebrisA", Flags: 4, DropThreshold: 100, Description: "Boots"}
	return s
}
func TestObjectDeathSmoke(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for op := 0; op < 8; op++ {
		specs = append(specs, deathBase(op))
	}
	callbackHash(t, "death-smoke", legacy.PortTestRoam(specs), deathHashes["death-smoke"])
}

func TestObjectDeathCreationContracts(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, op := range []int{1, 2} {
		for _, name := range []string{"", "PortTestDebrisA", "PortTestDebrisB", "missing", strings.Repeat("x", 127)} {
			for _, sound := range []uint32{0, 286, 757, 0xffffffff} {
				for _, flags := range []uint32{0, 4, 0x8000, 0xffff7fff, 0xffffffff} {
					s := deathBase(op)
					s.Callbacks.Death.Name = name
					s.Callbacks.Death.Sound = sound
					s.Callbacks.Death.Flags = flags
					specs = append(specs, s)
				}
			}
		}
	}
	callbackHash(t, "death-creation", legacy.PortTestRoam(specs), deathHashes["death-creation"])
}
func TestObjectDeathMarkerContracts(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, owner := range []bool{false, true} {
		for markers := uint32(0); markers < 16; markers++ {
			s := deathBase(3)
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
	callbackHash(t, "death-marker", r, deathHashes["death-marker"])
}
func TestObjectDeathBarrelContracts(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, barrel := range []bool{false, true} {
		for _, cache := range []uint32{0, 13, 23} {
			for _, count := range []uint32{0, 1, 3} {
				for _, threshold := range []uint32{0, 1, 50, 99, 100} {
					for _, seed := range []int{1, 7, 31, 63} {
						s := deathBase(0)
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
	callbackHash(t, "death-barrel", legacy.PortTestRoam(specs), deathHashes["death-barrel"])
}
func TestObjectDeathBoulderContracts(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for seed := 1; seed <= 64; seed++ {
		for rotation := uint32(0); rotation < 2; rotation++ {
			for mode := 0; mode < 3; mode++ {
				s := deathBase(4)
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
	callbackHash(t, "death-boulder", r, deathHashes["death-boulder"])
}

func TestObjectDeathEquipmentContracts(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, op := range []int{6, 7} {
		for _, language := range []int{0, 1, 2, 7} {
			for _, desc := range []string{"Boots", "BOOTS", "Boot", "s", "Glove\u00e9"} {
				for _, material := range []uint32{0, 2, 4, 8, 16, 30, 12, 6, 0x8000} {
					for _, owner := range []bool{false, true} {
						s := deathBase(op)
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
	callbackHash(t, "death-equipment", legacy.PortTestRoam(specs), deathHashes["death-equipment"])
}

var deathHashes = map[string]string{
	"death-smoke":     "1785f0bf1f7cfe72b397ea3b062054e3ff9d0af9ab4f11c6b23085ce673aaa3a",
	"death-creation":  "f991f5cff0f14409395cf9f6703a52063f763ed4c08dde028ce54f1cb84ea16e",
	"death-marker":    "1cdbc2b4cd5824f0ae75034dab6083d13bd1d1d9370ad215c4273150f2e42fa0",
	"death-barrel":    "cbced0b08ac7093301c6910e500e03482f9b51e34b2a335f6765c2ecd942d2c0",
	"death-boulder":   "a508a003c1d8d26b3f2e88cf5c5fdc7fd7a535c10c8d5e7dde7c1cad7da094e0",
	"death-equipment": "87e5b0d48d3333a54606023a0ce1d66967ff63e080452f15329445aed8b80b38",
}
