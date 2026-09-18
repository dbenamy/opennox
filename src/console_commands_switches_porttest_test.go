//go:build porttest

package opennox

import (
	"context"
	"github.com/opennox/libs/console"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestConsoleCommandsSwitches(t *testing.T) {
	o := newConsoleCommandOwner(t)
	oldCycle := legacy.Sub_4D0D70()
	t.Cleanup(func() { legacy.Sub_4D0D90(oldCycle) })
	type row struct {
		Input  string
		Before int
		Result bool
		After  int
		Output []consoleCommandLine
	}
	var rows []row
	for _, before := range []int{0, 1} {
		for _, input := range []string{"on", "ON", "off", "OFF", "invalid"} {
			legacy.Sub_4D0D90(before)
			got := o.call(t, "set cycle", false, input)
			want := before
			switch input {
			case "on", "ON":
				want = 1
			case "off", "OFF":
				want = 0
			}
			if got != (input != "invalid") || legacy.Sub_4D0D70() != want {
				t.Fatal(input, before, got, legacy.Sub_4D0D70())
			}
			rows = append(rows, row{input, before, got, legacy.Sub_4D0D70(), o.printer.lines})
		}
	}
	oldEngine := noxflags.GetEngine()
	t.Cleanup(func() { noxflags.ResetEngine(); noxflags.SetEngine(oldEngine) })
	noxflags.ResetEngine()
	ctx := console.WithCheats(console.AsServer(context.Background()))
	// Exercise actual registry parsing, permissions and dispatch for both forms.
	for _, command := range []string{"set netdebug", "unset netdebug", "set netdebug"} {
		if !o.console.Exec(ctx, command) {
			t.Fatal("registered execution", command)
		}
		if noxflags.HasEngine(noxflags.EngineNetDebug) != (command == "set netdebug") {
			t.Fatal("debug flag", command)
		}
	}
	for _, path := range []string{"set spellpoints", "allow user", "allow IP"} {
		if !o.call(t, path, false, "ignored") || len(o.printer.lines) != 1 || o.printer.lines[0].Text != "not implemented" {
			t.Fatal(path, o.printer.lines)
		}
	}
	spellbookCapture(t, "console-commands-switches", rows, "98b2c315db97febcfb2c081aa9b0336c768e88397d4940f8b9ba07262f9306dc")
}
