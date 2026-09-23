//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"math"
	"testing"
)

func updateRegistryGeneratorBase(op int) legacy.PortTestRoamSpec {
	s := generatorBase(op)
	s.Callbacks.Generator.RegisteredUpdate = true
	return s
}
func TestUpdateRegistryGeneratorUpdateGates(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, flags := range []uint32{0, 4, 0x1000004, 0x1000024, 0x1008004} {
		for _, stage := range []uint32{0, 9, 10, 11, 100, 0xffffffff} {
			for rate := byte(0); rate < 5; rate++ {
				for mode := 0; mode < 4; mode++ {
					s := updateRegistryGeneratorBase(3)
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
	updateRegistryHash(t, "update-registry-generator-update-gates", legacy.PortTestRoam(specs), updateRegistryHashes["update-registry-generator-update-gates"])
}

func TestUpdateRegistryGeneratorUpdateSpawn(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for sources := byte(1); sources <= 4; sources++ {
		for level := uint32(0); level < 3; level++ {
			for _, frame := range []uint32{7, 8, 135, 136} {
				for mode := 0; mode < 4; mode++ {
					s := updateRegistryGeneratorBase(8)
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
	updateRegistryHash(t, "update-registry-generator-update-spawn", r, updateRegistryHashes["update-registry-generator-update-spawn"])
}
