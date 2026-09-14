//go:build porttest

package opennox

import (
	"math"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
)

func TestMapPopulationExitPlacement(t *testing.T) {
	var cases []legacy.PortTestPaintSpec
	for dir := 0; dir < 4; dir++ {
		for _, blocked := range []bool{false, true} {
			for seed := 0; seed < 8; seed++ {
				s := populationBase()
				s.Seed = seed
				s.Globals["dword_5d4594_2487576"] = roomArg(2)
				s.Records[0].Words[472] = 1
				s.Records[0].Words[216+60] = uint32(dir)
				paintString(&s.Records[0], 216, "PopulationExit")
				paintString(&s.Records[0], 476, "next.map")
				if blocked {
					s.Records[1].Words[216] = 1 << uint(8*dir)
				}
				s.Actions = []legacy.PortTestPaintAction{paintAction(13, roomArg(1))}
				cases = append(cases, s)
			}
		}
	}
	out := populationCapture(t, "exit-placement", cases)
	for i, r := range out {
		want := uint32(1)
		if cases[i].Records[1].Words[216] != 0 {
			want = 0
		}
		if r.Steps[0].Return != want {
			t.Fatalf("case %d exit admission", i)
		}
		found := false
		for _, rec := range r.Steps[0].Records {
			if rec.Kind == "objectCollideData" {
				found = true
				if rec.Words[0] != 0x7478656e {
					t.Fatal("exit destination")
				}
			}
		}
		if found != (want == 1) {
			t.Fatal("exit object count")
		}
	}
}
func TestMapPopulationDecodedPlacement(t *testing.T) {
	var cases []legacy.PortTestPaintSpec
	for _, objects := range []bool{false, true} {
		for seed := 0; seed < 8; seed++ {
			s := populationCacheBase()
			s.Seed = seed
			s.Globals["occupancyEnabled"] = roomValue(1)
			s.Records[0].Words[68] = 32
			s.Records[0].Words[84] = 1
			s.Records[0].Refs[80] = roomArg(6)
			if objects {
				s.Objects = []legacy.PortTestPaintObject{{Slot: 100, Type: 1, Words: map[int]uint32{56: math.Float32bits(23), 60: math.Float32bits(46), 40: 123}}}
				s.Records[3].Refs[0] = roomArg(100)
				s.Globals["dword_5d4594_1599540"] = roomArg(4)
			}
			s.Actions = []legacy.PortTestPaintAction{paintAction(25, roomArg(1)), paintAction(32, roomArg(1))}
			cases = append(cases, s)
		}
	}
	out := populationCapture(t, "decoded-placement", cases)
	for i, r := range out {
		step := r.Steps[1]
		if step.Return != 1 || step.Globals["dword_5d4594_1599476"] != 1 {
			t.Fatalf("case %d prefab not applied", i)
		}
		if len(cases[i].Objects) > 0 && step.Objects[5] == 0 {
			t.Fatalf("case %d object not activated", i)
		}
	}
}
func TestMapPopulationEnchantments(t *testing.T) {
	var cases []legacy.PortTestPaintSpec
	for seed := 0; seed < 32; seed++ {
		for mask := 0; mask < 16; mask++ {
			s := populationBase()
			s.Seed = seed
			s.Records = append(s.Records, roomRecord(156), roomRecord(144), roomRecord(64), roomRecord(64))
			paintString(&s.Records[4], 0, "*")
			paintString(&s.Records[8], 0, "test-item")
			paintString(&s.Records[8], 60, "PopulationWeapon")
			paintString(&s.Records[10], 0, "PopulationModifier")
			paintString(&s.Records[11], 0, "PopulationModifier")
			s.Globals["modifier"] = roomArg(10)
			s.Globals["modifierName"] = roomArg(11)
			for j := 0; j < 4; j++ {
				if mask&(1<<j) != 0 {
					s.Records[8].Words[136+4*j] = 1
					s.Records[8].Refs[120+4*j] = roomArg(12)
				}
			}
			s.Actions = []legacy.PortTestPaintAction{paintAction(6, roomArg(5), roomArg(9), roomValue(1))}
			cases = append(cases, s)
		}
	}
	out := populationCapture(t, "enchantments", cases)
	positive := 0
	for i, r := range out {
		if r.Steps[0].Return == 0 {
			t.Fatalf("case %d no enchanted item", i)
		}
		for _, rec := range r.Steps[0].Records {
			if rec.Kind == "objectInitData" {
				if rec.Words[4] != 0 {
					t.Fatalf("case %d trailing item-attribute word is not initialized", i)
				}
				for j := 0; j < 4; j++ {
					if rec.Words[j] != 0 {
						positive++
					}
				}
				if rec.Words[2] == 0 && rec.Words[3] != 0 {
					t.Fatal("dependent modifier without predecessor")
				}
			}
		}
	}
	if positive == 0 {
		t.Fatal("no modifier applied")
	}
}
