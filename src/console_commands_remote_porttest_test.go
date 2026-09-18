//go:build porttest

package opennox

import (
	"fmt"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"testing"
)

func TestConsoleCommandsRemoteContext(t *testing.T) {
	for _, flags := range []uint32{0x4000, 0x8000, 0xc000} {
		for _, action := range []int{1, 2, 3, 6, 255} {
			t.Run(fmt.Sprintf("%x/%d", flags, action), func(t *testing.T) {
				o := newConsoleCommandOwner(t)
				t.Cleanup(noxflags.PortTestGameFlags(noxflags.GameFlag(flags)))
				serverConfigOwnBytes(t, 0x5D4594, 818228, 1024)
				unit, free := alloc.New(server.Object{})
				t.Cleanup(free)
				pl := &o.players[0]
				pl.PlayerUnit = unit
				text, freeText := alloc.CString16("ignored command")
				defer freeText()
				if legacy.PortTestConsoleRemote(nil, action, text) != 0 {
					t.Fatal("nil player accepted")
				}
				pl.PlayerUnit = nil
				if legacy.PortTestConsoleRemote(pl.C(), action, text) != 0 {
					t.Fatal("missing unit accepted")
				}
				pl.PlayerUnit = unit
				if legacy.PortTestConsoleRemote(pl.C(), action, text) != 1 {
					t.Fatal("mode-gated command return")
				}
				if legacy.PortTestConsoleSender() != nil {
					t.Fatal("mode-gated command left current sender set")
				}
			})
		}
	}
}
