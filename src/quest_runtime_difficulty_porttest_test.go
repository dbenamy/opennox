//go:build porttest

package opennox

import (
	"fmt"
	"math"
	"testing"

	noxflags "github.com/opennox/opennox/v1/common/flags"
)

func TestQuestRuntimeDifficultyCache(t *testing.T) {
	o := newQuestRuntimeOwner(t)
	defer noxflags.PortTestGameFlags(0)()
	type row struct {
		Name               string
		Return             [2]uint64
		Difficulty         [2]uint32
		Cache, Initialized uint32
	}
	var rows []row
	defer func() {
		spellbookCapture(t, "quest-runtime-difficulty-cache", rows, "581304947a6bf8200963c184b8f8c80ec8f2374093ed87425830b76457f35423")
	}()
	for count := 0; count <= 3; count++ {
		for _, stage := range []uint32{0, 1, 100, 0x7fffffff, 0x80000000, 0xffffffff} {
			for _, delta := range []float64{-1, 0, 0.1, 0.25, 2} {
				for _, initialized := range []uint32{0, 1, 0xffffffff} {
					name := fmt.Sprintf("count%d/stage%x/delta%g/initialized%x", count, stage, delta, initialized)
					t.Run(name, func(t *testing.T) {
						for i := range o.units {
							p := uint32(0)
							if i < count {
								p = 1
							}
							objectXferSetWord(o.units[i].UpdateDataPlayer().Player.C(), 4792, p)
						}
						*o.quest["202028"] = stage
						*o.quest["1563928"] = initialized
						*o.quest["1563912"] = math.Float32bits(0.75)
						used := float32(0.75)
						if initialized == 0 {
							used = float32(delta)
						}
						want := math.Float32bits(float32(float64(stage) * ((float64(count)-1)*float64(used) + 1)))
						var got row
						got.Name = name
						for pass := 0; pass < 2; pass++ {
							value := delta
							if pass == 1 {
								value = 99
							}
							o.balance(map[string]float64{"PlayerDifficultyDelta": value})
							got.Return[pass] = questRuntimeCall("sub_4E3D50", nil)
							got.Difficulty[pass] = *o.quest["202024"]
							if got.Return[pass] != uint64(want) || got.Difficulty[pass] != want {
								t.Fatalf("pass%d difficulty%x return%x want%x", pass, got.Difficulty[pass], got.Return[pass], want)
							}
						}
						got.Cache = *o.quest["1563912"]
						got.Initialized = *o.quest["1563928"]
						wantInit := initialized
						if initialized == 0 {
							wantInit = 1
						}
						if got.Cache != math.Float32bits(used) || got.Initialized != wantInit {
							t.Fatal("cache initialization")
						}
						rows = append(rows, got)
					})
				}
			}
		}
	}
}
