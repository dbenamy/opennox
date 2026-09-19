//go:build porttest

package opennox

import (
	"context"
	"fmt"
	"github.com/opennox/libs/log"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/unit/ai"
	"github.com/opennox/opennox/v1/legacy"
	"log/slog"
	"strings"
	"testing"
)

type unitDebugHandler struct{ messages []string }

func (h *unitDebugHandler) Enabled(context.Context, slog.Level) bool { return true }
func (h *unitDebugHandler) Handle(_ context.Context, r slog.Record) error {
	h.messages = append(h.messages, r.Message)
	return nil
}
func TestUnitGameplayDebug(t *testing.T) {
	oldLog := ai.Log
	enabled := noxflags.HasEngine(noxflags.EngineShowAI)
	defer func() {
		ai.Log = oldLog
		if enabled {
			noxflags.SetEngine(noxflags.EngineShowAI)
		} else {
			noxflags.UnsetEngine(noxflags.EngineShowAI)
		}
	}()
	h := new(unitDebugHandler)
	ai.Log = log.NewSlog(slog.New(log.NewHandler(h)))
	for _, on := range []bool{false, true} {
		if on {
			noxflags.SetEngine(noxflags.EngineShowAI)
		} else {
			noxflags.UnsetEngine(noxflags.EngineShowAI)
		}
		for kind := 0; kind < 2; kind++ {
			for _, frame := range []uint32{0, 1, 0x7fffffff, 0x80000000, 0xffffffff} {
				for _, code := range []uint32{0, 65535, 0x80000000, 0xffffffff} {
					for _, name := range []string{"", "Goblin", "A%sd", "name\x00suffix", "René"} {
						before := len(h.messages)
						legacy.PortTestUnitDebug(kind, frame, code, name)
						if !on {
							if len(h.messages) != before {
								t.Fatal("disabled AI debug emitted message")
							}
							continue
						}
						format := "%d: Lost sight of %s(#%d)\n"
						if kind == 1 {
							format = "%d: %s(#%d) FRUSTRATED\n"
						}
						want := fmt.Sprintf(format, int32(frame), strings.SplitN(name, "\x00", 2)[0], int32(code))
						if len(h.messages) != before+1 || h.messages[before] != want {
							t.Fatalf("debug message want %q got %v", want, h.messages[before:])
						}
					}
				}
			}
		}
	}
	spellbookCapture(t, "unit-gameplay-debug", h.messages, "672cbc0a6ad1a41af6fd0ce84350ada3528e897b176690bd8decc5151f9d3292")
}
