//go:build porttest

package opennox

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
)

func commandRulesAssert(t *testing.T, got legacy.PortTestCommandRulesResult, result int, commands []string) {
	t.Helper()
	if got.Result != result || !reflect.DeepEqual(got.Commands, commands) || !got.FilesUnchanged || !got.HeaderTableUnchanged || !got.FilenameBlobsSame || !got.HandlesUnchanged || !got.SettingsContextSame {
		t.Fatalf("got=%+v want result=%d commands=%q", got, result, commands)
	}
}

func TestRulesCommandHeadersABI(t *testing.T) {
	headers := []string{"", "unknown"}
	for _, mode := range ruleTestModes {
		headers = append(headers, mode.Name, strings.ToLower(mode.Name), " "+mode.Name, mode.Name+" ", mode.Name+"\r", mode.Name+"\x00ignored")
		for i := range []byte(mode.Name) {
			for _, b := range []byte{'a', '[', ']', ' ', 'X', 0xff} {
				s := []byte(mode.Name)
				s[i] = b
				headers = append(headers, string(s))
			}
		}
	}
	got, err := legacy.PortTestCommandRules(legacy.PortTestCommandRulesSpec{Mode: "header", Dir: t.TempDir(), Headers: headers})
	if err != nil {
		t.Fatal(err)
	}
	commandRulesAssert(t, got, 0, nil)
	if len(got.HeaderValues) != len(headers) {
		t.Fatal("missing header results")
	}
	for i, h := range headers {
		want := uint32(0)
		for _, m := range ruleTestModes {
			if ruleTestNul(h) == m.Name {
				want = m.Flags
				break
			}
		}
		if got.HeaderValues[i] != want {
			t.Fatalf("header=%q got=%x want=%x", h, got.HeaderValues[i], want)
		}
	}
}

type commandRuleLine struct{ raw, ending string }

func commandRuleDocument(lines []commandRuleLine) string {
	var out strings.Builder
	for _, l := range lines {
		out.WriteString(l.raw)
		out.WriteString(l.ending)
	}
	return out.String()
}

// Independent physical-line model; no production reader, header lookup or
// tokenizer participates in calculating the expected callback stream.
func commandRuleExpected(lines []commandRuleLine, flags uint32, after []uint32) []string {
	section := uint32(6128)
	var commands []string
	for _, l := range lines {
		raw := l.raw + l.ending
		if strings.HasSuffix(raw, "\r\n") {
			raw = raw[:len(raw)-2] + "\n"
		}
		if len(raw) > 254 {
			raw = raw[:254]
		}
		raw = ruleTestNul(raw)
		if i := strings.IndexByte(raw, '\n'); i >= 0 {
			raw = raw[:i]
		}
		if raw == "" {
			continue
		}
		header := false
		for _, m := range ruleTestModes {
			if raw == m.Name {
				section = m.Flags
				header = true
				break
			}
		}
		if header {
			continue
		}
		eligible := false
		for bit := uint(0); bit < 32; bit++ {
			if section/(1<<bit)%2 == 1 && flags/(1<<bit)%2 == 1 {
				eligible = true
				break
			}
		}
		if !eligible {
			continue
		}
		wide := make([]rune, len(raw))
		for i, b := range []byte(raw) {
			wide[i] = rune(b)
		}
		commands = append(commands, string(wide))
		if len(commands) <= len(after) {
			flags = after[len(commands)-1]
		}
	}
	return commands
}

func TestRulesCommandFilesABI(t *testing.T) {
	base := []commandRuleLine{
		{"before header", "\n"}, {"", "\n"}, {" \t ", "\n"}, {"# comment", "\n"},
		{"[DEATHMATCH]", "\n"}, {"first command", "\n"}, {"[deathmatch]", "\n"}, {"[COMMON] ", "\n"},
		{"[COMMON]", "\n"}, {"any mode", "\n"}, {"\x00ignored", "\n"}, {"prefix\x00ignored", "\n"},
		{"latin \xed", "\n"}, {"UTF8 í", "\n"}, {strings.Repeat("x", 300), "\n"}, {"after long line", "\n"},
		{"[QUEST]", "\n"}, {"quest command", "\n"}, {"[ELIMINATION]", "\n"}, {"elimination command", "\n"},
		{"[FLAGBALL]", "\n"}, {"flagball command", "\n"}, {"[CAPTURE_THE_FLAG]", "\n"}, {"ctf command", "\n"},
		{"[KING_OF_THE_REALM]", "\n"}, {"king command", "\n"}, {"[DEATHMATCH]\rbare CR text", "\n"},
	}
	for _, flags := range []uint32{0, 1, 16, 32, 64, 256, 1024, 4096, 6128, 0xffff, 0x80000000} {
		for _, format := range []string{"lf", "crlf", "final-no-lf"} {
			for _, dynamic := range []bool{false, true} {
				lines := append([]commandRuleLine(nil), base...)
				if format == "crlf" {
					for i := range lines {
						lines[i].ending = "\r\n"
					}
				}
				if format == "final-no-lf" {
					lines[len(lines)-1].ending = ""
				}
				var after []uint32
				if dynamic {
					after = []uint32{256, 4096, 6128, 0, 0xffff}
				}
				returns := make([]bool, 64)
				for i := range returns {
					returns[i] = i%2 != 0
				}
				spec := legacy.PortTestCommandRulesSpec{Mode: "file", Dir: t.TempDir(), Path: "rules.rul", Flags: flags, AfterCommandFlags: after, CallbackReturns: returns, Files: map[string]string{"rules.rul": commandRuleDocument(lines)}}
				got, err := legacy.PortTestCommandRules(spec)
				if err != nil {
					t.Fatal(err)
				}
				want := commandRuleExpected(lines, flags, after)
				commandRulesAssert(t, got, 1, want)
				finalFlags := flags
				n := len(want)
				if n > len(after) {
					n = len(after)
				}
				if n > 0 {
					finalFlags = after[n-1]
				}
				if got.FinalFlags != finalFlags {
					t.Fatalf("final flags=%x want=%x", got.FinalFlags, finalFlags)
				}
				if len(got.CallbackResults) != len(want) {
					t.Fatal("callback results missing")
				}
				for i, v := range got.CallbackResults {
					if v != returns[i] {
						t.Fatal("callback result mismatch")
					}
				}
			}
		}
	}
	for _, exists := range []bool{false, true} {
		files := map[string]string{}
		if exists {
			files["empty.rul"] = ""
		}
		got, err := legacy.PortTestCommandRules(legacy.PortTestCommandRulesSpec{Mode: "file", Dir: t.TempDir(), Path: "empty.rul", Files: files})
		if err != nil {
			t.Fatal(err)
		}
		result := 0
		if exists {
			result = 1
		}
		commandRulesAssert(t, got, result, nil)
	}
}

func TestRulesCommandSelectionABI(t *testing.T) {
	cases := []struct{ mode, path, mapName, user, fallback string }{
		{"path", `maps\Arena\Arena.map`, "", "maps/Arena/user.rul", "maps/Arena/Arena.rul"},
		{"wrapper", `maps\Arena\Arena.map`, "", "maps/Arena/user.rul", "maps/Arena/Arena.rul"},
		{"map", "", "Arena.map", "maps/Arena/user.rul", "maps/Arena/Arena.rul"},
		{"map", "", "Arena.xyz", "maps/Arena/user.rul", "maps/Arena/Arena.rul"},
		{"map", "", "a", "user.rul", ".rul"},
		{"map", "", "ab", "mauser.rul", "ma.rul"},
		{"map", "", "abc", "mapsuser.rul", "maps.rul"},
		{"path", "four", "", "user.rul", ".rul"},
		{"path", `maps\Arena\Arena.map`, "", "MAPS/ARENA/USER.RUL", "MAPS/ARENA/ARENA.RUL"},
		{"path", "Arena.map", "", "Arenauser.rul", "Arena.rul"},
		{"path", "maps/Arena/Arena.map", "", "maps/Arena/Arenauser.rul", "maps/Arena/Arena.rul"},
		{"path", `maps\Arena\Arena.map` + "\x00ignored", "", "maps/Arena/user.rul", "maps/Arena/Arena.rul"},
	}
	for _, c := range cases {
		for user := 0; user < 3; user++ {
			for _, fallback := range []bool{false, true} {
				files := map[string]string{"internet.rul": "internet must never run\n"}
				var want []string
				if user > 0 {
					files[c.user] = ""
					if user == 2 {
						files[c.user] = "user command\n"
						want = []string{"user command"}
					}
				}
				if fallback {
					files[c.fallback] = "map command\n"
					if user == 0 {
						want = []string{"map command"}
					}
				}
				got, err := legacy.PortTestCommandRules(legacy.PortTestCommandRulesSpec{Mode: c.mode, Dir: t.TempDir(), Path: c.path, Map: c.mapName, Flags: 256, Files: files})
				if err != nil {
					t.Fatal(err)
				}
				result := 1
				if c.mode == "wrapper" {
					result = 0
				}
				t.Run(fmt.Sprintf("%s/%s/user%d/map%v", c.mode, c.path, user, fallback), func(t *testing.T) { commandRulesAssert(t, got, result, want) })
			}
		}
	}
	got, err := legacy.PortTestCommandRules(legacy.PortTestCommandRulesSpec{Mode: "path", Dir: t.TempDir(), NilPath: true, Flags: 6128, Files: map[string]string{"user.rul": "must not run\n"}})
	if err != nil {
		t.Fatal(err)
	}
	commandRulesAssert(t, got, 0, nil)
}

// These inputs formerly underflowed a C buffer or caused a non-EOF read loop;
// test their bounded Go behavior without executing the undefined C cases.
func TestRulesCommandInputGuards(t *testing.T) {
	for _, mode := range []string{"path", "wrapper"} {
		for _, path := range []string{"", "a", "ab", "abc"} {
			got, err := legacy.PortTestCommandRules(legacy.PortTestCommandRulesSpec{Mode: mode, Dir: t.TempDir(), Path: path, Flags: 256, Files: map[string]string{"user.rul": "must not run\n"}})
			if err != nil {
				t.Fatal(err)
			}
			commandRulesAssert(t, got, 0, nil)
		}
	}
	got, err := legacy.PortTestCommandRules(legacy.PortTestCommandRulesSpec{Mode: "map", Dir: t.TempDir(), Map: "", Flags: 256})
	if err != nil {
		t.Fatal(err)
	}
	commandRulesAssert(t, got, 0, nil)
	for _, mode := range []string{"file", "path"} {
		spec := legacy.PortTestCommandRulesSpec{Mode: mode, Dir: t.TempDir(), Flags: 256}
		if mode == "file" {
			spec.Path = "rulesdir"
			spec.Files = map[string]string{"rulesdir/keep": "unchanged"}
		} else {
			spec.Path = `maps\Arena\Arena.map`
			spec.Files = map[string]string{"maps/Arena/user.rul/keep": "unchanged", "maps/Arena/Arena.rul": "must not run\n"}
		}
		got, err := legacy.PortTestCommandRules(spec)
		if err != nil {
			t.Fatal(err)
		}
		commandRulesAssert(t, got, 1, nil)
	}
	// Input truncation occurs before removing four bytes, so the fallback name
	// has 251 x's followed by .rul, within the filesystem's filename limit.
	got, err = legacy.PortTestCommandRules(legacy.PortTestCommandRulesSpec{Mode: "path", Dir: t.TempDir(), Path: strings.Repeat("x", 300), Flags: 256, Files: map[string]string{strings.Repeat("x", 251) + ".rul": "bounded fallback\n"}})
	if err != nil {
		t.Fatal(err)
	}
	commandRulesAssert(t, got, 1, []string{"bounded fallback"})
}
