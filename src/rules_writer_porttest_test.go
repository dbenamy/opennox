//go:build porttest

package opennox

import (
	"encoding/binary"
	"fmt"
	"sort"
	"testing"
	"unicode/utf16"

	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

// This baseline deliberately covers offline writing. The online C writer reads
// uninitialized/overlapping temporary buffers; its output is not a stable oracle.
func TestRulesWriterOfflineABI(t *testing.T) {
	_, catalog := server.PortTestRuleServerSetup()
	for _, flags := range []uint16{0, 16, 32, 64, 256, 1024, 4096, 6128, 128, 255} {
		for pattern := 0; pattern < 4; pattern++ {
			for _, retain := range []bool{false, true} {
				initial := ruleTestInitial(int64(pattern + 90))
				binary.LittleEndian.PutUint16(initial[52:54], flags)
				for i := 24; i < 52; i++ {
					switch pattern {
					case 0:
						initial[i] = 255
					case 1:
						initial[i] = 0
					case 2:
						initial[i] = 0xaa
					}
				}
				spec := legacy.PortTestRulesSpec{Kind: "write", Initial: initial, Context: 1234, Host: true, WithRejected: retain, Dir: t.TempDir(), User: "saved.rul", Path: "maps/porttst1/saved.rul", SeedRejected: []string{"# preserved", "set spell \"Fíre Burst\" off", "prefix\u0100invisible"}, Files: map[string]string{"maps/porttst1/keep.txt": "untouched"}}
				got, err := legacy.PortTestRules(spec)
				if err != nil {
					t.Fatal(err)
				}
				var rejected []string
				if retain {
					rejected = spec.SeedRejected
				}
				result := uint8(0)
				if flags&128 != 0 {
					result = uint8(flags)
				}
				if len(got.Steps) != 1 || !got.FilesUnchanged {
					t.Fatal("changed input files or missing snapshot")
				}
				ruleTestAssert(t, "offline writer", got.Steps[0], initial, spec.Context, result, rejected)
				want := ""
				if flags&128 == 0 {
					if retain {
						for _, line := range rejected {
							var b []byte
							for _, c := range utf16.Encode([]rune(line)) {
								b = append(b, byte(c))
							}
							b = append(b, '\n')
							want += ruleTestNul(string(b))
						}
					}
					for _, m := range ruleTestModes {
						if m.Flags == ruleTestMode(flags) {
							want += m.Name
							break
						}
					}
					want += "\n"
					spells := append([]server.PortTestRuleSpell(nil), catalog.Spells...)
					sort.Slice(spells, func(i, j int) bool { return spells[i].Index < spells[j].Index })
					for _, s := range spells {
						if s.Index > 0 && s.Index <= 136 && s.Valid && s.RuleFlags&0x7000000 != 0 && initial[24+s.Index/8]&(1<<uint(s.Index%8)) == 0 {
							want += fmt.Sprintf("set spell \"%s\" off\n", s.ID)
						}
					}
					for _, kind := range []string{"Armor", "weapon"} {
						items, offset, count := catalog.Armors, 48, 26
						if kind == "weapon" {
							items, offset, count = catalog.Weapons, 44, 27
						}
						for bit := 0; bit < count; bit++ {
							if initial[offset+bit/8]&(1<<uint(bit%8)) != 0 {
								continue
							}
							name := ""
							for _, it := range items {
								if it.Bit == 1<<uint(bit) {
									name = it.Name
									break
								}
							}
							want += fmt.Sprintf("set %s \"%s\" off\n", kind, name)
						}
					}
				}
				if got.WriteExists != (flags&128 == 0) || got.Written != want {
					t.Fatalf("flags=%x pattern=%d retain=%v got exists=%v output=%q want=%q", flags, pattern, retain, got.WriteExists, got.Written, want)
				}
			}
		}
	}
	// Opening a path whose parent does not exist fails without changing settings,
	// the input file, list or context. The writer's historical return is still 0.
	spec := legacy.PortTestRulesSpec{Kind: "write", Initial: ruleTestInitial(100), Context: 91, Dir: t.TempDir(), User: "missing/saved.rul", Path: "maps/porttst1/missing/saved.rul", Files: map[string]string{"maps/porttst1/keep.txt": "untouched"}}
	binary.LittleEndian.PutUint16(spec.Initial[52:54], 256)
	got, err := legacy.PortTestRules(spec)
	if err != nil {
		t.Fatal(err)
	}
	if got.WriteExists || got.Written != "" || !got.FilesUnchanged || len(got.Steps) != 1 {
		t.Fatal("failed create changed files")
	}
	ruleTestAssert(t, "failed create", got.Steps[0], spec.Initial, spec.Context, 0, nil)
}

// The online oracle models the intended filter, not the unstable C stack reads.
func TestRulesWriterOnlineFiltering(t *testing.T) {
	for trial := 0; trial < 3; trial++ {
		for internetMask := 0; internetMask < 4; internetMask++ {
			for localMask := 0; localMask < 4; localMask++ {
				for currentMask := 0; currentMask < 4; currentMask++ {
					for _, source := range []string{"user", "map", "empty-user"} {
						initial := ruleTestInitial(int64(trial))
						for i := 24; i < 52; i++ {
							initial[i] = 255
						}
						binary.LittleEndian.PutUint16(initial[52:54], 256)
						names := []string{"SPELL_BLINK", "SPELL_FIREBALL"}
						ids := []int{4, 27}
						document := func(mask int) string {
							out := "[DEATHMATCH]\n"
							for i, name := range names {
								if mask&(1<<uint(i)) != 0 {
									out += fmt.Sprintf("set spell %s off\n", name)
								}
							}
							return out
						}
						files := map[string]string{"internet.rul": document(internetMask), "maps/porttst1/porttst1.rul": document(localMask)}
						effective := localMask
						if source == "user" {
							files["maps/porttst1/user.rul"] = document(localMask)
							files["maps/porttst1/porttst1.rul"] = document(localMask ^ 3)
						}
						if source == "empty-user" {
							files["maps/porttst1/user.rul"] = ""
							effective = 0
						}
						want := "[DEATHMATCH]\n"
						for i, id := range ids {
							bit := 1 << uint(i)
							if currentMask&bit == 0 {
								continue
							}
							initial[24+id/8] &^= 1 << uint(id%8)
							// Emit a local restriction unless solely inherited from internet rules.
							if internetMask&bit == 0 || effective&bit != 0 {
								want += fmt.Sprintf("set spell \"%s\" off\n", names[i])
							}
						}
						spec := legacy.PortTestRulesSpec{Kind: "write", Wrapper: trial%2 == 1, Initial: initial, Context: 1234, Online: 0x80000000, Host: true, WithRejected: true, Dir: t.TempDir(), User: "saved.rul", Path: "maps/porttst1/saved.rul", Files: files}
						got, err := legacy.PortTestRules(spec)
						if err != nil {
							t.Fatal(err)
						}
						context := uint32(256)
						if source == "empty-user" {
							context = 6128
						}
						if len(got.Steps) != 1 || !got.FilesUnchanged || !got.WriteExists || got.Written != want {
							t.Fatalf("trial=%d internet=%d local=%d current=%d source=%s got=%q want=%q", trial, internetMask, localMask, currentMask, source, got.Written, want)
						}
						ruleTestAssert(t, "online writer", got.Steps[0], initial, context, 0, nil)
					}
				}
			}
		}
	}
}

func TestRulesWriterOnlineFileEdges(t *testing.T) {
	for _, variant := range []string{"missing-inputs", "overwrite-user", "blocked", "failed-create"} {
		for _, wrapper := range []bool{false, true} {
			initial := ruleTestInitial(101)
			for i := 24; i < 52; i++ {
				initial[i] = 255
			}
			initial[24] &^= 16
			binary.LittleEndian.PutUint16(initial[52:54], 256)
			spec := legacy.PortTestRulesSpec{Kind: "write", Wrapper: wrapper, Initial: initial, Context: 1234, Online: 1, Host: true, Dir: t.TempDir(), User: "saved.rul", Path: "maps/porttst1/saved.rul", Files: map[string]string{"maps/porttst1/keep": "unchanged"}}
			want := "[DEATHMATCH]\nset spell \"SPELL_BLINK\" off\n"
			context, result, exists, unchanged := uint32(6128), uint8(0), true, true
			switch variant {
			case "overwrite-user":
				spec.User = "user.rul"
				spec.Path = "maps/porttst1/user.rul"
				spec.Files[spec.Path] = "[DEATHMATCH]\nset spell SPELL_BLINK off\n"
				spec.Files["internet.rul"] = "[DEATHMATCH]\nset spell SPELL_BLINK off\n"
				// Creation truncates user.rul before the local load, so that input is empty.
				want = "[DEATHMATCH]\n"
				unchanged = false
			case "blocked":
				binary.LittleEndian.PutUint16(spec.Initial[52:54], 128)
				want = ""
				exists = false
				context = 1234
				if !wrapper {
					result = 128
				}
			case "failed-create":
				spec.User = "missing/saved.rul"
				spec.Path = "maps/porttst1/missing/saved.rul"
				want = ""
				exists = false
				context = 1234
			}
			got, err := legacy.PortTestRules(spec)
			if err != nil {
				t.Fatal(err)
			}
			if len(got.Steps) != 1 || got.FilesUnchanged != unchanged || got.WriteExists != exists || got.Written != want {
				t.Fatalf("%s wrapper=%v got=%+v want=%q", variant, wrapper, got, want)
			}
			ruleTestAssert(t, variant, got.Steps[0], spec.Initial, context, result, nil)
		}
	}
}
