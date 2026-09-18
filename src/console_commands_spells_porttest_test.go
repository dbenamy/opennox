//go:build porttest

package opennox

import (
	"encoding/binary"
	"fmt"
	"github.com/opennox/libs/console"
	"github.com/opennox/libs/spell"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/server"
	"reflect"
	"testing"
)

func TestConsoleCommandsSpells(t *testing.T) {
	type row struct {
		Input, Action         string
		Mode                  uint32
		Selected              uint16
		Before, Result, After bool
		Updated               uint32
		Output                []consoleCommandLine
	}
	var rows []row
	for _, input := range []struct {
		name string
		id   spell.ID
	}{{"Fire title", spell.SPELL_FIREBALL}, {spell.SPELL_FIREBALL.String(), spell.SPELL_FIREBALL}, {"Special title", 132}, {"missing", 0}} {
		for _, mode := range []uint32{0, 64, 128} {
			for _, selected := range []uint16{0, 64} {
				for _, before := range []bool{false, true} {
					for _, action := range []string{"on", "ON", "off", "invalid"} {
						t.Run(fmt.Sprintf("%s/%x/%x/%t/%s", input.name, mode, selected, before, action), func(t *testing.T) {
							o := newConsoleCommandOwner(t)
							t.Cleanup(noxflags.PortTestGameFlags(noxflags.GameFlag(mode)))
							s := o.c.srv.S()
							t.Cleanup(s.PortTestSpellLifecycle([]server.PortTestSpellLifecycleDef{{Index: int(spell.SPELL_FIREBALL), Valid: true, Enabled: before}, {Index: 132, Valid: true, Enabled: before}}, nil))
							s.Spells.DefByInd(spell.SPELL_FIREBALL).Title = "Fire title"
							s.Spells.DefByInd(132).Title = "Special title"
							binary.LittleEndian.PutUint16(o.settings[52:], selected)
							*o.optionWords["settings-updated"] = 0
							got := o.call(t, "set spell", false, input.name, action)
							want, after, updated := false, before, false
							var output []consoleCommandLine
							say := func(s string) { output = []consoleCommandLine{{console.ColorRed, s + " " + input.name}} }
							switch {
							case mode == 128:
								want = true
								say("not in chat")
							case input.id == 0:
								say("invalid")
							case action == "on" || action == "ON":
								want = true
								if !(input.id == 132 && (mode&64 != 0 || selected&64 != 0)) && !before {
									after = true
									updated = true
									say("enabled")
								}
							case action == "off":
								want = true
								if before {
									after = false
									updated = true
									say("disabled")
								}
							}
							observed := before
							if input.id != 0 {
								observed = s.Spells.DefByInd(input.id).Enabled
							}
							if got != want || observed != after || (*o.optionWords["settings-updated"] != 0) != updated || !reflect.DeepEqual(o.printer.lines, output) {
								t.Fatalf("result=%t want=%t enabled=%t want=%t updated=%d lines=%v want=%v", got, want, observed, after, *o.optionWords["settings-updated"], o.printer.lines, output)
							}
							rows = append(rows, row{input.name, action, mode, selected, before, got, observed, *o.optionWords["settings-updated"], o.printer.lines})
						})
					}
				}
			}
		}
	}
	spellbookCapture(t, "console-commands-spells", rows, "2c55920c2b0850b4bc47366b8fa8d620abd3ef1272105e3da1c5be22fd948c4c")
}
