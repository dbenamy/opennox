//go:build porttest

package opennox

import (
	"math"
	"os"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
)

// Run each prerequisite probe in its own process before repairs: legacy stack
// corruption must not invalidate captures from unrelated cases.
func TestMapPopulationPrerequisiteProbe(t *testing.T) {
	probe := os.Getenv("OPENNOX_POPULATION_PROBE")
	if probe == "" {
		t.Skip("isolated prerequisite diagnostic")
	}
	s := populationBase()
	switch probe {
	case "spell-name":
		paintString(&s.Records[4], 0, "fireball")
		s.Actions = []legacy.PortTestPaintAction{paintAction(1, roomArg(5))}
	case "spellbook-invalid":
		s.Globals["gameFlags"] = roomValue(0)
		s.Records[0].Words[1096] = 1
		paintString(&s.Records[4], 0, "not_a_spell")
		s.Actions = []legacy.PortTestPaintAction{paintAction(5, roomArg(1), roomArg(5))}
	case "point-output":
		s.Records[0].Words[1096] = 1
		s.Records[0].Words[548] = 10
		s.Records[5].Words[0] = 5
		paintString(&s.Records[5], 4, "*")
		s.Actions = []legacy.PortTestPaintAction{paintAction(10, roomArg(1), roomArg(2), roomArg(6))}
	case "finale-point":
		s.Globals["dword_5d4594_1550916"] = roomArg(3)
		s.Actions = []legacy.PortTestPaintAction{paintAction(12, roomArg(1))}
	default:
		t.Fatal("unknown probe")
	}
	out := populationCapture(t, "prerequisite-"+probe, []legacy.PortTestPaintSpec{s})
	if probe == "spellbook-invalid" {
		if out[0].Steps[0].Return != 0 || out[0].Steps[0].Objects[0] != 0 {
			t.Fatal("invalid spellbook leaked pooled object")
		}
	}
	if probe == "finale-point" {
		room := s.Records[2].Words
		x := float32((float64(math.Float32frombits(room[44])) + float64(math.Float32frombits(room[36]))) * 0.5)
		y := float32((float64(math.Float32frombits(room[48])) + float64(math.Float32frombits(room[40]))) * 0.5)
		want := [2]uint32{math.Float32bits(float32((float64(y)+float64(x))*0.70710677 + 2957)), math.Float32bits(float32((float64(y)-float64(x))*0.70710677 + 2956))}
		found := false
		for _, r := range out[0].Steps[0].Records {
			if r.Kind == "object" {
				found = true
				if r.Words[14] != want[0] || r.Words[15] != want[1] {
					t.Fatalf("PlayerStart position %x want %x", r.Words[14:16], want)
				}
			}
		}
		if !found {
			t.Fatal("missing PlayerStart")
		}
	}
	if probe == "spell-name" && out[0].Steps[0].Return != 27 {
		t.Fatal("known spell not found")
	}
}

func TestMapPopulationPrerequisiteRegressions(t *testing.T) {
	for _, probe := range []string{"spell-name", "spellbook-invalid", "point-output", "finale-point"} {
		t.Run(probe, func(t *testing.T) { t.Setenv("OPENNOX_POPULATION_PROBE", probe); TestMapPopulationPrerequisiteProbe(t) })
	}
}
