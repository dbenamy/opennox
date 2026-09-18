//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/console"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"reflect"
	"testing"
)

func TestConsoleCommandsMute(t *testing.T) {
	type row struct {
		Command                                 string
		Client, ServerToken, Connected, Missing bool
		Slot                                    int
		Before, After                           uint32
		Result                                  bool
		Output                                  []consoleCommandLine
	}
	var rows []row
	for _, command := range []string{"mute", "unmute"} {
		for _, client := range []bool{false, true} {
			for _, serverToken := range []bool{false, true} {
				for _, connected := range []bool{false, true} {
					for _, slot := range []int{1, 31} {
						for _, missing := range []bool{false, true} {
							t.Run(fmt.Sprintf("%s/%t/%t/%t/%d/%t", command, client, serverToken, connected, slot, missing), func(t *testing.T) {
								o := newConsoleCommandOwner(t)
								for i := range o.players {
									o.players[i].Active = 0
								}
								p := &o.players[0]
								p.Active = 1
								p.PlayerInd = byte(slot)
								p.SetName("Player Ω")
								before := uint32(0x80000040)
								if command == "unmute" {
									before |= 12
								}
								p.Field3680 = before
								flags := noxflags.GameFlag(0)
								if connected {
									flags = 2
								}
								t.Cleanup(noxflags.PortTestGameFlags(flags))
								name := "player Ω"
								if missing {
									name = "missing"
								}
								args := []string{name}
								if serverToken {
									args = []string{"server", name}
								}
								got := o.call(t, command, client, args...)
								success := !missing
								bit := uint32(4)
								if client || !serverToken {
									bit = 8
									success = success && connected && (!serverToken) && (command == "unmute" || slot != 31)
								}
								want := before
								if success {
									if command == "mute" {
										want |= bit
									} else {
										want &^= bit
									}
								}
								label := "missing"
								if success {
									label = "muted"
									if command == "unmute" {
										label = "unmuted"
									}
								}
								printedName := name
								if client && serverToken {
									printedName = "server"
								}
								output := []consoleCommandLine{{console.ColorRed, label + " " + printedName}}
								if !got || p.Field3680 != want || !reflect.DeepEqual(o.printer.lines, output) {
									t.Fatalf("result=%t status=%x want=%x lines=%v want=%v", got, p.Field3680, want, o.printer.lines, output)
								}
								rows = append(rows, row{command, client, serverToken, connected, missing, slot, before, p.Field3680, got, o.printer.lines})
							})
						}
					}
				}
			}
		}
	}
	spellbookCapture(t, "console-commands-mute", rows, "c2635c4a49873cf125ff792127bf01fd2abea90fb94bd2ce389830c7e7e17ade")
}
func TestConsoleCommandsArity(t *testing.T) {
	o := newConsoleCommandOwner(t)
	for _, path := range []string{"set weapons", "set staffs", "set monsters", "set name", "set sysop", "set players", "set time", "set lessons", "set cycle", "set spell", "ban", "kick", "mute", "unmute", "show motd", "show seq", "set netdebug", "unset netdebug"} {
		// One excess argument count beyond every accepted arity, without dereferencing tokens.
		if o.call(t, path, false, "a", "b", "c", "d") {
			if path != "set name" {
				t.Errorf("accepted excess arguments: %s", path)
			}
		}
	}
	for _, path := range []string{"set weapons", "set staffs", "set monsters", "set name", "set sysop", "set players", "set time", "set lessons", "set cycle", "set spell", "ban", "kick", "mute", "unmute"} {
		if o.call(t, path, false) {
			t.Errorf("accepted missing argument: %s", path)
		}
	}
}
