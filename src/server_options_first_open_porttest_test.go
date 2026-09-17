//go:build porttest

package opennox

import (
	"context"
	"fmt"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestServerOptionsFirstOpen(t *testing.T) {
	type row struct {
		First, Headless bool
		Quest           int
		Commands        []string
	}
	var rows []row
	for _, first := range []bool{false, true} {
		for _, headless := range []bool{false, true} {
			for _, quest := range []int{0, 1, 2} {
				t.Run(fmt.Sprintf("%t-%t-%d", first, headless, quest), func(t *testing.T) {
					o := newServerOptionsOwner(t)
					defer noxflags.PortTestGameFlags(0)()
					oldEngine := noxflags.GetEngine()
					noxflags.ResetEngine()
					if headless {
						noxflags.SetEngine(noxflags.EngineNoRendering)
					}
					t.Cleanup(func() { noxflags.ResetEngine(); noxflags.SetEngine(oldEngine) })
					oldQuest := questFlag_1556156
					questFlag_1556156 = quest == 1
					t.Cleanup(func() { questFlag_1556156 = oldQuest })
					qp := memmap.PtrUint32(0x5D4594, 1556160)
					oldWord := *qp
					*qp = 0
					if quest == 2 {
						*qp = 1
					}
					t.Cleanup(func() { *qp = oldWord })
					var commands []string
					oldExec := legacy.ExecConsoleCmd
					legacy.ExecConsoleCmd = func(_ context.Context, cmd string) bool { commands = append(commands, cmd); return true }
					t.Cleanup(func() { legacy.ExecConsoleCmd = oldExec })
					o.installConstructor(t)
					if first {
						*o.optionWords["first-open"] = 1
					}
					for i := 0; i < 2; i++ {
						if o.call("construct", 0, "") != 1 || *o.optionWords["first-open"] != 0 {
							t.Fatal("first-open state")
						}
						o.call("close", 0, "")
						o.options = nil
						o.c.GUI.FreeDestroyed()
					}
					var want []string
					if first {
						if quest != 0 {
							want = []string{"execrul OTQuest.rul"}
						} else if headless {
							want = []string{"execrul server.rul"}
						}
					}
					if fmt.Sprint(commands) != fmt.Sprint(want) {
						t.Fatalf("startup commands %q want %q", commands, want)
					}
					rows = append(rows, row{first, headless, quest, commands})
				})
			}
		}
	}
	spellbookCapture(t, "server-options-first-open", rows, "5d92618b553926b6630771c5bb2e01dc4c90814d18dcbc496a6eedd60deb2ef0")
}
