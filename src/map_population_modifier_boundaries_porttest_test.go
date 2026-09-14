//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestMapPopulationModifierBoundaries(t *testing.T) {
	var cases []legacy.PortTestPaintSpec
	for _, typ := range []string{"PopulationWeapon", "PopulationArmor"} {
		for _, selector := range []string{"*", "test-item", "TEST-ITEM"} {
			for _, modifier := range []string{"PopulationModifier", "none", "unknown"} {
				for _, chance := range []uint32{0, 100} {
					for mask := 0; mask < 16; mask++ {
						s := populationBase()
						s.Records = append(s.Records, roomRecord(156), roomRecord(144), roomRecord(64), roomRecord(64))
						paintString(&s.Records[4], 0, selector)
						paintString(&s.Records[8], 0, "test-item")
						paintString(&s.Records[8], 60, typ)
						paintString(&s.Records[10], 0, "PopulationModifier")
						paintString(&s.Records[11], 0, modifier)
						s.Globals["modifier"] = roomArg(10)
						s.Globals["modifierName"] = roomArg(11)
						for j := 0; j < 4; j++ {
							s.Globals[fmt.Sprintf("populationTable%d", 254688+4*j)] = roomValue(int32(chance))
							if mask&(1<<j) != 0 {
								s.Records[8].Words[136+4*j] = 1
								s.Records[8].Refs[120+4*j] = roomArg(12)
							}
						}
						s.Actions = []legacy.PortTestPaintAction{paintAction(6, roomArg(5), roomArg(9), roomValue(1))}
						cases = append(cases, s)
					}
				}
			}
		}
	}
	out := populationCapture(t, "modifier-boundaries", cases)
	for i, r := range out {
		step := r.Steps[0]
		if step.Return == 0 {
			t.Fatalf("case %d missing item", i)
		}
		// The first 32 cases in each 96-case group use the registered modifier.
		registered := (i/32)%3 == 0
		certain := (i/16)%2 == 1
		mask := i % 16
		found := false
		for _, rec := range step.Records {
			if rec.Kind != "objectInitData" {
				continue
			}
			found = true
			for j := 0; j < 4; j++ {
				want := registered && certain && mask&(1<<j) != 0
				if j == 3 && mask&4 == 0 {
					want = false
				}
				if (rec.Words[j] != 0) != want {
					t.Fatalf("case %d modifier slot %d", i, j)
				}
			}
			if rec.Words[4] != 0 {
				t.Fatalf("case %d trailing attributes", i)
			}
		}
		if !found {
			t.Fatalf("case %d missing attribute storage", i)
		}
	}
}
