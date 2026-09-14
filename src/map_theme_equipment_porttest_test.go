//go:build porttest

package opennox

import (
	"os"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc/handles"
)

func themeStoredString(words []uint32, off int) string {
	var out []byte
	for i := off; i < len(words)*4; i++ {
		v := byte(words[i/4] >> uint(8*(i%4)))
		if v == 0 {
			break
		}
		out = append(out, v)
	}
	return string(out)
}

func TestMapThemeTemplateRemovalProbe(t *testing.T) {
	handles.Init()
	defer handles.Release()
	t.Chdir(t.TempDir())
	input := "WEAPON TEMPLATE QUALITY alpha beta gamma END END WEAPON tier Sword QUALITY -alpha END END END "
	if err := os.WriteFile("theme-fixture.dat", []byte(input), 0600); err != nil {
		t.Fatal(err)
	}
	s := themeBase()
	s.Actions = []legacy.PortTestPaintAction{paintAction(11, roomArg(1), roomValue(legacy.PortTestThemeFile))}
	out := themeRun([]legacy.PortTestPaintSpec{s})
	step := out[0].Steps[0]
	if !out[0].Intact || !out[0].ControlOK || step.Return != 1 {
		t.Fatal("equipment template parsing")
	}
	rows := map[uint32][]uint32{}
	for _, r := range step.Records {
		if r.Alive {
			rows[r.ID] = r.Words
		}
	}
	item := rows[rows[step.Slots[1]][275]]
	if len(item) != 39 || item[34] != 2 {
		t.Fatal("template removal count")
	}
	values := rows[item[30]]
	if themeStoredString(values, 0) != "beta" || themeStoredString(values, 60) != "gamma" {
		t.Fatal("template removal did not preserve the two remaining values")
	}
}

func TestMapThemeTemplateRemovalBoundaries(t *testing.T) {
	changes := []struct {
		input string
		want  []string
	}{
		{"-alpha", []string{"beta", "gamma"}}, {"-beta", []string{"alpha", "gamma"}}, {"-gamma", []string{"alpha", "beta"}},
		{"-ALPHA", []string{"beta", "gamma"}}, {"-missing", []string{"alpha", "beta", "gamma"}}, {"-alpha delta", []string{"beta", "gamma", "delta"}},
		{"-alpha beta BETA", []string{"beta", "gamma"}}, {"-alpha -alpha", []string{"beta", "gamma"}}, {"-alpha -beta -gamma", nil},
	}
	names := []string{"QUALITY", "MATERIAL", "PRIMARY_ENCHANTMENT", "SECONDARY_ENCHANTMENT"}
	var cases []themeInputCase
	for _, kind := range []string{"WEAPON", "ARMOR"} {
		for _, slot := range names {
			for _, change := range changes {
				text := kind + " TEMPLATE " + slot + " alpha beta gamma END END " + kind + " tier item " + slot + " " + change.input + " END END END "
				s := themeBase()
				op := 11
				if kind == "ARMOR" {
					op = 14
				}
				s.Actions = []legacy.PortTestPaintAction{paintAction(op, roomArg(1), roomValue(legacy.PortTestThemeFile)), paintAction(32, roomArg(1))}
				cases = append(cases, themeInputCase{text, s})
			}
		}
	}
	out := themeCapture(t, "template-removals", cases)
	for i, r := range out {
		step := r.Steps[0]
		if step.Return != 1 || step.Globals["themeTemplate"] != 0 {
			t.Fatalf("case %d template completion", i)
		}
		rows := map[uint32][]uint32{}
		for _, rec := range step.Records {
			if rec.Alive {
				rows[rec.ID] = rec.Words
			}
		}
		head := 275
		if i >= 4*len(changes) {
			head = 277
		}
		item := rows[rows[step.Slots[1]][head]]
		slot := i / len(changes) % 4
		want := changes[i%len(changes)].want
		if len(item) != 39 || item[34+slot] != uint32(len(want)) {
			t.Fatalf("case %d template result count", i)
		}
		for j, name := range want {
			if themeStoredString(rows[item[30+slot]], 60*j) != name {
				t.Fatalf("case %d remaining item %d", i, j)
			}
		}
		for _, rec := range r.Steps[1].Records {
			if rec.Kind == "themeAllocation" && rec.Alive {
				t.Fatalf("case %d equipment cleanup", i)
			}
		}
	}
}
