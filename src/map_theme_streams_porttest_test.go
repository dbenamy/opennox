//go:build porttest

package opennox

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc/handles"
)

type themeInputCase struct {
	Text string
	Spec legacy.PortTestPaintSpec
}

func themeCapture(t *testing.T, label string, inputs []themeInputCase) []legacy.PortTestPaintResult {
	t.Helper()
	handles.Init()
	defer handles.Release()
	t.Chdir(t.TempDir())
	cases := make([]legacy.PortTestPaintSpec, len(inputs))
	for i, input := range inputs {
		name := fmt.Sprintf("theme-input-%04d.dat", i)
		if err := os.WriteFile(name, []byte(input.Text), 0600); err != nil {
			t.Fatal(err)
		}
		s := input.Spec
		if len(s.Records) == 0 {
			s = themeBase()
		}
		paintString(&s.Records[7], 0, name)
		cases[i] = s
	}
	return themeCheckCapture(t, label, themeRun(cases))
}

func themeCheckCapture(t *testing.T, label string, out []legacy.PortTestPaintResult) []legacy.PortTestPaintResult {
	t.Helper()
	for i, r := range out {
		if !r.Intact || !r.ControlOK {
			t.Fatalf("%s case %d guard/control state", label, i)
		}
	}
	data, err := json.Marshal(out)
	if err != nil {
		t.Fatal(err)
	}
	if prefix := os.Getenv("OPENNOX_MAP_THEME_CAPTURE"); prefix != "" {
		if err := os.WriteFile(prefix+"-"+label+".json", data, 0600); err != nil {
			t.Fatal(err)
		}
	}
	if want, ok := themeCExpected[label]; ok {
		got := fmt.Sprintf("%x", sha256.Sum256(data))
		if got != want {
			t.Fatalf("%s complete capture differs from C: got %s want %s", label, got, want)
		}
	} else {
		t.Fatalf("missing locked C capture for %s", label)
	}
	t.Logf("%s: %d cases %x", label, len(out), sha256.Sum256(data))
	return out
}
func themeStream(text string, op, steps int) themeInputCase {
	s := themeBase()
	for i := 0; i < steps; i++ {
		s.Actions = append(s.Actions, paintAction(op, roomValue(legacy.PortTestThemeFile), roomValue(legacy.PortTestThemeToken)))
	}
	return themeInputCase{text, s}
}
func themeWithPlayers(s *legacy.PortTestPaintSpec, classes []byte) {
	r := roomRecord(132)
	r.Words[0] = uint32(len(classes))
	for i, v := range classes {
		r.Words[4+4*i] = uint32(v)
	}
	s.Records = append(s.Records, r)
	s.Globals["themePlayers"] = roomArg(len(s.Records))
}
func TestMapThemeRawTokens(t *testing.T) {
	texts := []string{"", " ", "\n\t\r\v\f", "one", "one ", "one two ", "abcdefgh second", "// comment\nlast ", "// missing newline", "a// comment\nlast ", "/ /tail ", "///\nnext ", "a\r\nb\n", "a\x00b ", "IF 100% yes ELSE no ENDIF tail "}
	for c := 0; c < 256; c++ {
		texts = append(texts, "a"+string([]byte{byte(c)})+"b tail ")
	}
	for _, n := range []int{1, 2, 59, 60, 254, 255} {
		texts = append(texts, strings.Repeat("x", n)+" next ")
	}
	var cases []themeInputCase
	for _, text := range texts {
		cases = append(cases, themeStream(text, 2, 5))
	}
	themeCapture(t, "raw-tokens", cases)
}
func TestMapThemeFilteredTokens(t *testing.T) {
	var cases []themeInputCase
	for _, text := range []string{"", "word ", "IF 100% yes ELSE no ENDIF tail ", "IF 0% yes ELSE no ENDIF tail ", "IF 100% yes ENDIF tail ", "IF 0% no ENDIF tail ", "IF WIZARD wizard ELSE other ENDIF tail ", "IF CONJURER conjurer ELSE other ENDIF tail ", "IF WARRIOR warrior ELSE other ENDIF tail ", "IF invalid body ELSE alternate ENDIF tail ", "ELSE skipped ENDIF tail ", "ENDIF tail ", "IF 100%", "IF 0% body", "IF 100% IF 0% hidden ELSE nested ENDIF ELSE outer ENDIF tail "} {
		for _, classes := range [][]byte{nil, {0}, {1}, {2}, {0, 1, 2}} {
			c := themeStream(text, 1, 6)
			themeWithPlayers(&c.Spec, classes)
			cases = append(cases, c)
		}
	}
	for depth := 1; depth <= 6; depth++ {
		for _, truth := range []string{"0%", "100%"} {
			text := strings.Repeat("IF "+truth+" ", depth) + "value " + strings.Repeat("ENDIF ", depth) + "tail "
			cases = append(cases, themeStream(text, 1, 5))
		}
	}
	themeCapture(t, "filtered-tokens", cases)
}
func TestMapThemeSkips(t *testing.T) {
	var cases []themeInputCase
	for _, op := range []int{3, 5, 6} {
		for _, text := range []string{"", "plain", "plain\nnext ", "\n\nlast ", "ENDIF tail ", "ELSE alternate ENDIF tail ", "IF nested ENDIF ENDIF tail ", "IF nested ELSE skipped ENDIF ELSE tail ", "IF unterminated ", "//comment\nENDIF tail "} {
			for _, line := range []uint32{0, 100, 0xfffffffe} {
				c := themeStream(text, op, 1)
				c.Spec.Globals["themeLine"] = roomValue(int32(line))
				c.Spec.Actions = append(c.Spec.Actions, paintAction(2, roomValue(legacy.PortTestThemeFile), roomValue(legacy.PortTestThemeToken)))
				cases = append(cases, c)
			}
		}
	}
	themeCapture(t, "skips", cases)
}
func TestMapThemeOperators(t *testing.T) {
	var cases []themeInputCase
	for _, text := range []string{"< ", "> ", "== ", "!= ", "<= ", ">= ", "= ", "? ", "", "<", "// comment\n== ", "IF 100% != ELSE < ENDIF "} {
		s := themeBase()
		s.Records[2].Words[0] = 0x12345678
		s.Actions = []legacy.PortTestPaintAction{paintAction(8, roomValue(legacy.PortTestThemeFile), roomArg(3))}
		cases = append(cases, themeInputCase{text, s})
	}
	out := themeCapture(t, "operators", cases)
	for i := 0; i < 4; i++ {
		step := out[i].Steps[0]
		if step.Return != 1 {
			t.Fatal("valid relation rejected")
		}
		for _, r := range step.Records {
			if r.ID == step.Slots[3] && r.Words[0] != uint32(i) {
				t.Fatal("relation code")
			}
		}
	}
}
func TestMapThemeConditions(t *testing.T) {
	var cases []themeInputCase
	classSets := [][]byte{nil, {0}, {1}, {2}, {0, 1, 2}, {2, 2, 1}, make([]byte, 32)}
	for _, classes := range classSets {
		var texts []string
		texts = append(texts, "WIZARD ", "wizard ", "CONJURER ", "WARRIOR ", "unknown ", "", "NUMPLAYERS <= 1 ", "NUMPLAYERſ == 0 ")
		for _, key := range []string{"NUMPLAYERS", "EXPERIENCE_LEVEL"} {
			for _, op := range []string{"<", ">", "==", "!="} {
				for _, n := range []int{-1, 0, 1, 2, 3, 4, 31, 32, 33} {
					texts = append(texts, fmt.Sprintf("%s %s %d ", key, op, n))
				}
			}
		}
		for _, text := range texts {
			s := themeBase()
			themeWithPlayers(&s, classes)
			s.Records[2].Words[0] = 0x12345678
			s.Actions = []legacy.PortTestPaintAction{paintAction(7, roomValue(legacy.PortTestThemeFile), roomArg(3))}
			cases = append(cases, themeInputCase{text, s})
		}
	}
	for seed := 0; seed < 16; seed++ {
		for _, pct := range []int{-1, 0, 1, 49, 50, 99, 100, 101} {
			s := themeBase()
			s.Seed = seed
			s.Actions = []legacy.PortTestPaintAction{paintAction(7, roomValue(legacy.PortTestThemeFile), roomArg(3))}
			cases = append(cases, themeInputCase{fmt.Sprintf("%d%% ", pct), s})
		}
	}
	themeCapture(t, "conditions", cases)
}
