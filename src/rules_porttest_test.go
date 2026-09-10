//go:build porttest

package opennox

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"testing"
	"unicode/utf16"

	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

var ruleTestModes = []struct {
	Name  string
	Flags uint32
}{
	{"[ELIMINATION]", 1024}, {"[DEATHMATCH]", 256}, {"[CAPTURE_THE_FLAG]", 32},
	{"[KING_OF_THE_REALM]", 16}, {"[FLAGBALL]", 64}, {"[COMMON]", 6128}, {"[QUEST]", 4096},
}

func ruleTestMode(v uint16) uint32 {
	var out uint32
	for _, bit := range []uint{4, 5, 6, 7, 8, 9, 10, 12} {
		if (uint32(v)/(1<<bit))%2 == 1 {
			out += 1 << bit
		}
	}
	return out
}

func ruleTestInitial(seed int64) (out [60]byte) {
	gen := rand.New(rand.NewSource(seed))
	gen.Read(out[:])
	copy(out[:], "porttst1\x00")
	return out
}

func ruleTestNul(s string) string {
	if i := strings.IndexByte(s, 0); i >= 0 {
		return s[:i]
	}
	return s
}

func ruleTestNarrow(s string) string {
	var out []byte
	for _, v := range utf16.Encode([]rune(ruleTestNul(s))) {
		if byte(v) == 0 {
			break
		}
		out = append(out, byte(v))
	}
	return string(out)
}

func ruleTestKeyword(s, want string) bool {
	s = ruleTestNul(s)
	b := []byte(s)
	for i, c := range b {
		if c >= 'A' && c <= 'Z' {
			b[i] = c + 32
		}
	}
	return string(b) == want
}

// The model consumes already-known tokens. It independently computes settings
// bit positions; parser tests specify token expectations rather than calling
// the production tokenizer to determine their own expected result.
func ruleTestCommand(settings *[60]byte, context *uint32, mode uint32, host bool, tokens []string, catalog server.PortTestRuleServer) uint8 {
	if len(tokens) == 0 {
		return 0
	}
	first := ruleTestNarrow(tokens[0])
	for _, m := range ruleTestModes {
		if first == m.Name {
			*context = m.Flags
			if mode == m.Flags {
				return 1
			}
			return 0
		}
	}
	if *context&mode == 0 || !ruleTestKeyword(tokens[0], "set") || len(tokens) != 4 {
		return 0
	}
	off := ruleTestKeyword(tokens[3], "off")
	switch {
	case ruleTestKeyword(tokens[1], "spell"):
		var found *server.PortTestRuleSpell
		id := ruleTestNarrow(tokens[2])
		for i := range catalog.Spells {
			if catalog.Spells[i].ID == id {
				found = &catalog.Spells[i]
				break
			}
		}
		if found == nil {
			for i := range catalog.Spells {
				if catalog.Spells[i].Title == ruleTestNul(tokens[2]) {
					found = &catalog.Spells[i]
					break
				}
			}
		}
		if found == nil || !found.Valid {
			return 0
		}
		if found.RuleFlags&0x7000000 != 0 && off {
			settings[24+found.Index/8] &^= byte(1 << uint(found.Index%8))
		}
	case ruleTestKeyword(tokens[1], "weapon"), ruleTestKeyword(tokens[1], "armor"):
		if !host {
			return 0
		}
		items, offset := catalog.Weapons, 44
		if ruleTestKeyword(tokens[1], "armor") {
			items, offset = catalog.Armors, 48
		}
		var bit uint32
		for _, it := range items {
			if strings.ToLower(it.Name) == strings.ToLower(ruleTestNarrow(tokens[2])) {
				bit = it.Bit
				break
			}
		}
		if bit == 0 {
			return 0
		}
		old := binary.LittleEndian.Uint32(settings[offset : offset+4])
		if off {
			old &^= bit
		} else {
			old |= bit
		}
		binary.LittleEndian.PutUint32(settings[offset:offset+4], old)
	default:
		return 0
	}
	if mode == *context {
		return 1
	}
	return 0
}

func ruleTestAssert(t *testing.T, label string, got legacy.PortTestRulesState, want [60]byte, context uint32, result uint8, rejected []string) {
	t.Helper()
	if got.Settings != want || got.Context != context || got.Result != result || !reflect.DeepEqual(got.Rejected, rejected) || !got.LinksValid || !got.GuardsValid || !got.TableUnchanged || !got.HandlesUnchanged {
		t.Fatalf("%s got=%+v want settings=%x context=%x result=%x rejected=%q", label, got, want, context, result, rejected)
	}
}

func TestRulesModeHeadersABI(t *testing.T) {
	values := make([]uint16, 65536)
	for i := range values {
		values[i] = uint16(i)
	}
	got := legacy.PortTestRuleHeaders(values)
	if len(got) != len(values) {
		t.Fatal("missing header results")
	}
	for i, v := range values {
		want := -1
		for j, m := range ruleTestModes {
			if m.Flags == ruleTestMode(v) {
				want = j
				break
			}
		}
		if got[i] != want {
			t.Fatalf("flags=%04x got index=%d want=%d", v, got[i], want)
		}
	}
}

func TestRulesDirectivesABI(t *testing.T) {
	_, catalog := server.PortTestRuleServerSetup()
	commands := [][]string{{"unknown"}, {"set"}, {"set", "spell"}, {"set", "spell", "SPELL_FIREBALL"}, {"set", "spell", "SPELL_FIREBALL", "off", "extra"}, {"set", "other", "thing", "off"}, {"[deathmatch]"}, {"set", "spell", "missing", "off"}, {"set", "weapon", "missing", "off"}, {"set", "armor", "missing", "off"}}
	commands = append(commands,
		[]string{"ſet", "spell", "SPELL_FIREBALL", "off"},
		[]string{"\u0173et", "spell", "SPELL_FIREBALL", "off"},
		[]string{"set", "ſpell", "SPELL_FIREBALL", "off"},
		[]string{"set", "spell", "SPELL_FIREBALL", "\u016fff"},
		[]string{"[COMMON]\u0100suffix"},
		[]string{"set", "spell", "SPELL_FIREBALL\u0100suffix", "off"},
	)
	for _, m := range ruleTestModes {
		commands = append(commands, []string{m.Name}, []string{m.Name, "ignored"})
	}
	for _, s := range catalog.Spells {
		for _, name := range []string{s.ID, s.Title, strings.ToLower(s.ID), strings.ToLower(s.Title)} {
			for _, op := range []string{"off", "OFF", "on", "anything"} {
				commands = append(commands, []string{"SeT", "SpElL", name, op})
			}
		}
	}
	for _, category := range []string{"weapon", "armor"} {
		items := catalog.Weapons
		if category == "armor" {
			items = catalog.Armors
		}
		for _, it := range items {
			if it.Bit == 0 || it.Bit&(it.Bit-1) != 0 {
				t.Fatal("fixture expected single-bit item", it)
			}
			for _, op := range []string{"off", "on", "OFF", "anything"} {
				commands = append(commands, []string{"SET", category, strings.ToUpper(it.Name), op})
			}
		}
	}
	// Each directive is preceded by the same section to keep its eligibility
	// independent of previous header cases, while settings changes accumulate.
	for _, mode := range []uint16{0, 16, 32, 64, 256, 1024, 4096, 6128, 0x1400, 0xffff} {
		for _, header := range []string{"[COMMON]", "[DEATHMATCH]", "[QUEST]"} {
			for _, host := range []bool{false, true} {
				var tokens [][]string
				for _, cmd := range commands {
					tokens = append(tokens, []string{header}, cmd)
				}
				spec := legacy.PortTestRulesSpec{Kind: "tokens", Initial: ruleTestInitial(int64(mode)), Flags: mode, Context: 0xdeadbeef, Host: host, WithRejected: true, SeedRejected: []string{"preserved"}, Tokens: tokens}
				got, err := legacy.PortTestRules(spec)
				if err != nil {
					t.Fatal(err)
				}
				if !got.FilesUnchanged || len(got.Steps) != len(tokens) {
					t.Fatal("fixture shape or files changed")
				}
				want, context := spec.Initial, spec.Context
				for i, cmd := range tokens {
					result := ruleTestCommand(&want, &context, uint32(mode), host, cmd, catalog)
					ruleTestAssert(t, strings.Join(cmd, " "), got.Steps[i], want, context, result, spec.SeedRejected)
				}
			}
		}
	}
}

func TestRulesLinesABI(t *testing.T) {
	_, catalog := server.PortTestRuleServerSetup()
	cases := []struct {
		line   string
		tokens []string
	}{
		{"", nil}, {" \t\r\n", nil}, {"# comment", []string{"#", "comment"}},
		{"[COMMON]", []string{"[COMMON]"}}, {"[deathmatch]", []string{"[deathmatch]"}},
		{"\"[DEATHMATCH]\"", []string{"[DEATHMATCH]"}},
		{"set spell SPELL_FIREBALL off", []string{"set", "spell", "SPELL_FIREBALL", "off"}},
		{"SET\tSpElL\tSPELL_FIREBALL\tOFF", []string{"SET", "SpElL", "SPELL_FIREBALL", "OFF"}},
		{"set spell \"Rule Fireball\" off", []string{"set", "spell", "Rule Fireball", "off"}},
		{"set spell \"Fíre Burst\" off", []string{"set", "spell", "Fíre Burst", "off"}},
		{"set spell \"rule fireball\" off", []string{"set", "spell", "rule fireball", "off"}},
		{"\"set\" spell SPELL_FIREBALL off", []string{"set", "spell", "SPELL_FIREBALL", "off"}},
		{"set \"spell\" \"Rule Fireball\" \"off\"", []string{"set", "spell", "Rule Fireball", "off"}},
		{"set  \"spell\" \"Rule Fireball\" off", []string{"set", "\"spell\"", "Rule Fireball", "off"}},
		{" \"set\" spell SPELL_FIREBALL off", []string{"\"set\"", "spell", "SPELL_FIREBALL", "off"}},
		{" set \"spell\" \"Rule Fireball\" off", []string{"set", "spell", "Rule Fireball", "off"}},
		{"set spell \"Rule Fireball\"off", []string{"set", "spell", "Rule Fireball", "off"}},
		{"set spell \"Rule Fireball", []string{"set", "spell", "Rule Fireball"}},
		{"set spell \"\" off", []string{"set", "spell", " off"}},
		{"set spell SPELL_FIREBALL off\x00ignored", []string{"set", "spell", "SPELL_FIREBALL", "off"}},
		{"set weapon Bow off", []string{"set", "weapon", "Bow", "off"}},
		{"set armor LeatherArmor off", []string{"set", "armor", "LeatherArmor", "off"}},
		{"set spell SPELL_ANCHOR off", []string{"set", "spell", "SPELL_ANCHOR", "off"}},
		{"set spell SPELL_BURN off", []string{"set", "spell", "SPELL_BURN", "off"}},
		{"unknown " + strings.Repeat("x", 240), []string{"unknown", strings.Repeat("x", 240)}},
	}
	many := append([]string{"set"}, strings.Fields(strings.Repeat("x ", 31))...)
	cases = append(cases, struct {
		line   string
		tokens []string
	}{strings.Join(many, " "), many})
	for _, mode := range []uint16{0, 256, 4096, 6128} {
		header := "[COMMON]"
		if mode == 256 {
			header = "[DEATHMATCH]"
		} else if mode == 4096 {
			header = "[QUEST]"
		}
		for _, host := range []bool{false, true} {
			for index, c := range cases {
				for _, trailing := range []string{"", " "} {
					if strings.ContainsRune(c.line, 0) && trailing != "" {
						continue
					}
					line := c.line + trailing
					initial := ruleTestInitial(int64(index))
					for i := 24; i < 52; i++ {
						initial[i] = 255
					}
					spec := legacy.PortTestRulesSpec{Kind: "lines", Initial: initial, Flags: mode, Context: 0xdeadbeef, Host: host, WithRejected: true, SeedRejected: []string{"before"}, Lines: []string{header, line}}
					got, err := legacy.PortTestRules(spec)
					if err != nil {
						t.Fatal(err)
					}
					if len(got.Steps) != 2 {
						t.Fatal("missing line states")
					}
					want, context := initial, spec.Context
					rejected := append([]string(nil), spec.SeedRejected...)
					for i, tokens := range [][]string{{header}, c.tokens} {
						result := ruleTestCommand(&want, &context, uint32(mode), host, tokens, catalog)
						if len(tokens) > 0 && result == 0 {
							rejected = append(rejected, ruleTestNul(spec.Lines[i]))
						}
						ruleTestAssert(t, spec.Lines[i], got.Steps[i], want, context, 0, rejected)
					}
				}
			}
		}
	}
}

// Documents pair physical file bytes with explicit token expectations. This
// oracle deliberately has no tokenizer in common with the loader under test.
type ruleTestFileLine struct {
	raw, ending string
	tokens      []string
}
type ruleTestDocument []ruleTestFileLine

func (d ruleTestDocument) data() string {
	var b strings.Builder
	for _, l := range d {
		b.WriteString(l.raw)
		b.WriteString(l.ending)
	}
	return b.String()
}

func ruleTestRead(d ruleTestDocument, found bool, settings *[60]byte, context *uint32, mode uint32, host, retain bool, rejected *[]string, catalog server.PortTestRuleServer) bool {
	*context = 6128
	if !found {
		return false
	}
	for _, line := range d {
		raw := line.raw
		if line.ending != "" {
			raw += "\n"
		}
		if len(raw) > 255 {
			raw = raw[:255]
		}
		raw = ruleTestNul(raw)
		if i := strings.IndexByte(raw, '\n'); i >= 0 {
			raw = raw[:i]
		}
		if raw == "" {
			continue
		}
		var wide []rune
		for _, b := range []byte(raw) {
			wide = append(wide, rune(b))
		}
		if len(line.tokens) != 0 && ruleTestCommand(settings, context, mode, host, line.tokens, catalog) == 0 && retain {
			*rejected = append(*rejected, string(wide))
		}
	}
	return true
}

func ruleTestDocs() (ruleTestDocument, ruleTestDocument, ruleTestDocument) {
	makeDoc := func(lines ...string) ruleTestDocument {
		var d ruleTestDocument
		for _, s := range lines {
			d = append(d, ruleTestFileLine{s, "\n", strings.Fields(s)})
		}
		return d
	}
	return makeDoc("[COMMON]", "# map", "set spell SPELL_FIREBALL off", "[DEATHMATCH]", "set weapon Bow off", "set armor LeatherArmor off"),
		makeDoc("[COMMON]", "# user", "set spell SPELL_BURN off", "[DEATHMATCH]", "set spell SPELL_FIREBALL off", "set armor LeatherArmor off", "set weapon Bow on"),
		makeDoc("[COMMON]", "# internet", "set weapon Bow off", "set armor LeatherArmor on", "set spell SPELL_BLINK off", "[QUEST]", "set spell SPELL_FIREBALL off")
}

func ruleTestLoad(t *testing.T, spec legacy.PortTestRulesSpec, docs map[string]ruleTestDocument) {
	t.Helper()
	spec.Dir = t.TempDir()
	spec.Files = make(map[string]string)
	for path, doc := range docs {
		spec.Files[path] = doc.data()
	}
	got, err := legacy.PortTestRules(spec)
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Steps) != 1 || !got.FilesUnchanged {
		t.Fatal("file mutation or missing snapshot")
	}
	want, context := spec.Initial, spec.Context
	for i := 24; i < 52; i++ {
		want[i] = 255
	}
	var rejected []string // Top-level load always clears the supplied list.
	read := func(path string) bool {
		var doc ruleTestDocument
		found := false
		for p, d := range docs {
			if strings.EqualFold(p, path) {
				doc, found = d, true
				break
			}
		}
		return ruleTestRead(doc, found, &want, &context, ruleTestMode(spec.Flags), spec.Host, spec.WithRejected, &rejected, got.Catalog)
	}
	mapName := ruleTestNul(string(spec.Initial[:8]))
	base := "maps/" + mapName + "/"
	loaded := false
	if spec.Selection&2 != 0 {
		user := spec.User
		if spec.UserNil {
			user = "user.rul"
		}
		loaded = read(base + user)
	}
	if spec.Selection&1 != 0 && !loaded {
		read(base + mapName + ".rul")
	}
	if spec.Selection&4 != 0 && spec.Online != 0 {
		read("internet.rul")
	}
	result := uint8(spec.Flags)
	if spec.Flags&64 != 0 {
		want[40] &^= 16
		result = 239
	}
	if spec.Wrapper {
		result = 0
	}
	ruleTestAssert(t, fmt.Sprintf("load selection=%x flags=%x online=%x files=%v", spec.Selection, spec.Flags, spec.Online, reflect.ValueOf(docs).MapKeys()), got.Steps[0], want, context, result, rejected)
}

func TestRulesLoadPrecedenceABI(t *testing.T) {
	mapDoc, userDoc, internetDoc := ruleTestDocs()
	for _, flags := range []uint16{0, 256, 4096, 6128, 64} {
		for selection := uint8(0); selection < 8; selection++ {
			for presence := 0; presence < 8; presence++ {
				for _, online := range []uint32{0, 0x80000000} {
					docs := make(map[string]ruleTestDocument)
					if presence&1 != 0 {
						docs["maps/porttst1/porttst1.rul"] = mapDoc
					}
					if presence&2 != 0 {
						docs["maps/porttst1/user.rul"] = userDoc
					}
					if presence&4 != 0 {
						docs["internet.rul"] = internetDoc
					}
					spec := legacy.PortTestRulesSpec{Kind: "load", Initial: ruleTestInitial(int64(presence)), Selection: selection, Flags: flags, Context: 0xdeadbeef, Online: online, Host: (int(selection)+presence)%2 == 0, WithRejected: true, UserNil: true, SeedRejected: []string{"must be cleared", "also removed"}}
					ruleTestLoad(t, spec, docs)
				}
			}
		}
	}
	for _, variant := range []string{"custom", "wrapper", "empty-user", "uppercase", "eight-bytes", "short-name", "no-list", "high-selection"} {
		t.Run(variant, func(t *testing.T) {
			spec := legacy.PortTestRulesSpec{Kind: "load", Initial: ruleTestInitial(17), Selection: 7, Flags: 0xffdf, Context: 0xdeadbeef, Online: 1, Host: true, WithRejected: true, UserNil: true, SeedRejected: []string{"old"}}
			mapName, userName := "porttst1", "user.rul"
			u := userDoc
			switch variant {
			case "custom", "wrapper":
				spec.UserNil = false
				spec.User = "custom.rul"
				userName = spec.User
				spec.Wrapper = variant == "wrapper"
			case "empty-user":
				u = ruleTestDocument{}
			case "eight-bytes":
				copy(spec.Initial[:], "porttst1EXTRA\x00")
			case "short-name":
				copy(spec.Initial[:], "abc\x00rest")
				mapName = "abc"
			case "no-list":
				spec.WithRejected = false
			case "high-selection":
				spec.Selection = 255
			}
			docs := map[string]ruleTestDocument{"maps/" + mapName + "/" + mapName + ".rul": mapDoc, "maps/" + mapName + "/" + userName: u, "internet.rul": internetDoc}
			if variant == "uppercase" {
				upper := make(map[string]ruleTestDocument)
				for p, d := range docs {
					upper[strings.ToUpper(p)] = d
				}
				docs = upper
			}
			ruleTestLoad(t, spec, docs)
		})
	}
}

func TestRulesFileBytesABI(t *testing.T) {
	doc := ruleTestDocument{
		{"", "\n", nil}, {" \t", "\r\n", nil},
		{"[COMMON]", "\r\n", []string{"[COMMON]"}},
		{"set spell \"F\xedre Burst\" off", "\r\n", []string{"set", "spell", "Fíre Burst", "off"}},
		{"set spell \"Fíre Burst\" off", "\n", []string{"set", "spell", "FÃ\u00adre Burst", "off"}},
		{"\x00set weapon Bow off", "\n", nil},
		{"set armor LeatherArmor off\x00ignored", "\n", []string{"set", "armor", "LeatherArmor", "off"}},
		{"unknown " + strings.Repeat("x", 400), "\n", []string{"unknown", strings.Repeat("x", 247)}},
		{"[DEATHMATCH]", "\n", []string{"[DEATHMATCH]"}},
		{"set spell SPELL_FIREBALL off", "", []string{"set", "spell", "SPELL_FIREBALL", "off"}},
	}
	crOnly := ruleTestDocument{{"[COMMON]\rset weapon Bow off\rset armor LeatherArmor off", "", []string{"[COMMON]", "set", "weapon", "Bow", "off", "set", "armor", "LeatherArmor", "off"}}}
	for _, file := range []struct {
		name  string
		doc   ruleTestDocument
		found bool
	}{{"bytes", doc, true}, {"cr-only", crOnly, true}, {"empty", nil, true}, {"missing", nil, false}} {
		for _, mode := range []uint16{0, 256, 4096, 6128, 64} {
			for _, host := range []bool{false, true} {
				spec := legacy.PortTestRulesSpec{Kind: "file", Initial: ruleTestInitial(29), Flags: mode, Context: 0xdeadbeef, Host: host, WithRejected: true, SeedRejected: []string{"retained"}, Dir: t.TempDir(), Path: "rules.rul", Files: make(map[string]string)}
				if file.found {
					spec.Files[spec.Path] = file.doc.data()
				}
				got, err := legacy.PortTestRules(spec)
				if err != nil {
					t.Fatal(err)
				}
				want, context := spec.Initial, spec.Context
				rejected := append([]string(nil), spec.SeedRejected...)
				result := uint8(0)
				if ruleTestRead(file.doc, file.found, &want, &context, uint32(mode), host, true, &rejected, got.Catalog) {
					result = 1
				}
				if len(got.Steps) != 1 || !got.FilesUnchanged {
					t.Fatal("file mutation or missing snapshot")
				}
				ruleTestAssert(t, file.name, got.Steps[0], want, context, result, rejected)
			}
		}
	}
}
