//go:build porttest

package opennox

import (
	"github.com/opennox/libs/spell"
	"strings"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
)

func TestMapPopulationInventories(t *testing.T) {
	var cases []legacy.PortTestPaintSpec
	for seed := 0; seed < 64; seed++ {
		for _, weight := range []uint32{0, 1, 50, 99, 100} {
			s := populationBase()
			s.Seed = seed
			s.Objects = []legacy.PortTestPaintObject{{Slot: 100, Type: 1}}
			s.Records = append(s.Records, roomRecord(2060), roomRecord(2060))
			s.Records[8].Words[0] = weight
			s.Records[8].Words[2052] = 2
			s.Records[8].Refs[2056] = roomArg(10)
			s.Records[9].Words[0] = 100
			s.Records[9].Words[2052] = 1
			paintString(&s.Records[8], 8, "PaintObject")
			paintString(&s.Records[8], 72, "PaintBook")
			paintString(&s.Records[9], 8, "PaintDoor")
			s.Actions = []legacy.PortTestPaintAction{paintAction(4, roomArg(1), roomArg(100), roomArg(9))}
			cases = append(cases, s)
		}
	}
	out := populationCapture(t, "inventories", cases)
	for i, r := range out {
		if r.Steps[0].Objects[0] == 0 {
			t.Fatalf("case %d missing live owner", i)
		}
		count := 0
		for _, rec := range r.Steps[0].Records {
			if rec.Kind == "object" {
				count++
			}
		}
		if count < 2 || count > 3 {
			t.Fatalf("case %d expected owner and 1–2 children, got %d", i, count)
		}
	}
}
func TestMapPopulationSpellbooks(t *testing.T) {
	var cases []legacy.PortTestPaintSpec
	for seed := 0; seed < 64; seed++ {
		for _, count := range []uint32{0, 1, 2, 5} {
			s := populationBase()
			s.Seed = seed
			s.Records[0].Words[1096] = count
			for j := 0; j < 5; j++ {
				s.Records[0].Words[548+4*j] = uint32(10 + j)
			}
			paintString(&s.Records[4], 0, "*")
			s.Actions = []legacy.PortTestPaintAction{paintAction(5, roomArg(1), roomArg(5))}
			cases = append(cases, s)
		}
	}
	out := populationCapture(t, "spellbooks", cases)
	for i, r := range out {
		positive := cases[i].Records[0].Words[1096] > 0
		if (r.Steps[0].Return != 0) != positive {
			t.Fatal("book admission")
		}
	}
}
func TestMapPopulationPlainItems(t *testing.T) {
	var cases []legacy.PortTestPaintSpec
	for _, name := range []string{"PaintObject", "paintdoor", "PaintBook", "missing"} {
		s := populationBase()
		paintString(&s.Records[4], 0, name)
		s.Actions = []legacy.PortTestPaintAction{paintAction(6, roomArg(5), roomValue(0), roomValue(0))}
		cases = append(cases, s)
	}
	out := populationCapture(t, "plain-items", cases)
	for i, r := range out {
		if (r.Steps[0].Return != 0) != (i < 3) {
			t.Fatal("item admission")
		}
	}
}
func TestMapPopulationExitNames(t *testing.T) {
	var cases []legacy.PortTestPaintSpec
	for _, typ := range []int{1, 5} {
		for _, n := range []int{0, 1, 79, 80, 81, 120} {
			s := populationBase()
			s.Objects = []legacy.PortTestPaintObject{{Slot: 100, Type: typ}}
			name := make([]byte, n)
			for i := range name {
				name[i] = byte('a' + i%26)
			}
			paintString(&s.Records[4], 0, string(name))
			s.Actions = []legacy.PortTestPaintAction{paintAction(37, roomArg(100), roomArg(5)), paintAction(37, roomValue(0), roomArg(5)), paintAction(37, roomArg(100), roomValue(0))}
			cases = append(cases, s)
		}
	}
	out := populationCapture(t, "exit-names", cases)
	for i, r := range out {
		want := uint32(0)
		if cases[i].Objects[0].Type == 5 {
			want = 1
		}
		if r.Steps[0].Return != want || r.Steps[1].Return != 0 || r.Steps[2].Return != 0 {
			t.Fatal("exit type/null gates")
		}
	}
}

func TestMapPopulationSpellNames(t *testing.T) {
	var cases []legacy.PortTestPaintSpec
	var wants []uint32
	for id := spell.ID(1); id.Valid(); id++ {
		name := strings.TrimPrefix(id.String(), "SPELL_")
		for _, v := range []string{name, strings.ToLower(name)} {
			s := populationBase()
			paintString(&s.Records[4], 0, v)
			s.Actions = []legacy.PortTestPaintAction{paintAction(1, roomArg(5))}
			cases = append(cases, s)
			wants = append(wants, uint32(id))
		}
	}
	for _, n := range []int{0, 1, 53, 54, 59, 60, 61, 120} {
		s := populationBase()
		paintString(&s.Records[4], 0, strings.Repeat("x", n))
		s.Actions = []legacy.PortTestPaintAction{paintAction(1, roomArg(5))}
		cases = append(cases, s)
		wants = append(wants, 0)
	}
	out := populationCapture(t, "spell-names", cases)
	for i, r := range out {
		if r.Steps[0].Return != wants[i] {
			t.Fatalf("case %d spell %d want %d", i, r.Steps[0].Return, wants[i])
		}
	}
}
