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
