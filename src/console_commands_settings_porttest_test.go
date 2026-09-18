//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/console"
	"github.com/opennox/opennox/v1/legacy"
	"reflect"
	"testing"
)

func TestConsoleCommandsToggles(t *testing.T) {
	type row struct {
		Command                string
		Args                   []string
		Client, Result         bool
		Before, After, Updated uint32
		Output                 []consoleCommandLine
	}
	var rows []row
	for _, spec := range []struct {
		command, label string
		bit            uint32
		prefix         []string
	}{
		{"set weapons", "weapons", 1, nil}, {"set staffs", "staffs", 16, nil}, {"set monsters", "monsters", 4, nil},
		{"set monsters", "respawn", 8, []string{"respawn"}},
	} {
		for _, client := range []bool{false, true} {
			for _, before := range []uint32{0, 0xffffffff} {
				for _, arg := range []string{"on", "ON", "off", "oFf", "invalid"} {
					t.Run(fmt.Sprintf("%s/%s/%t/%x/%s", spec.command, spec.label, client, before, arg), func(t *testing.T) {
						o := newConsoleCommandOwner(t)
						legacy.Sub_409EC0(-1)
						legacy.Sub_409E70(int(before))
						*o.optionWords["settings-updated"] = 0
						args := append(append([]string(nil), spec.prefix...), arg)
						got := o.call(t, spec.command, client, args...)
						on := arg == "on" || arg == "ON"
						valid := arg != "invalid"
						want := before
						var output []consoleCommandLine
						if valid {
							if on {
								want |= spec.bit
							} else {
								want &^= spec.bit
							}
							value := "OFF"
							if on {
								value = "ON"
							}
							output = []consoleCommandLine{{console.ColorRed, spec.label + " " + value}}
						}
						after := legacy.Nox_xxx_getServerSubFlags_409E60()
						if got != valid || after != want || (*o.optionWords["settings-updated"] != 0) != valid || !reflect.DeepEqual(o.printer.lines, output) || legacy.PortTestConsoleServer() == client {
							t.Fatalf("result=%t flags=%x want=%x updated=%d output=%v want=%v", got, after, want, *o.optionWords["settings-updated"], o.printer.lines, output)
						}
						rows = append(rows, row{spec.command, args, client, got, before, after, *o.optionWords["settings-updated"], o.printer.lines})
					})
				}
			}
		}
	}
	spellbookCapture(t, "console-commands-toggles", rows, "83241a0df860517c0338e8b6ef6b1f919ce8b8ed00e9fcf65ff87813249e75e2")
}
