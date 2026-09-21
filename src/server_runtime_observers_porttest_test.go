//go:build porttest

package opennox

import (
	"testing"

	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
)

func TestServerRuntimeObserverCount(t *testing.T) {
	o := newObjectDrawingOwner(t)
	players, free := o.c.srv.PortTestObjectRenderPlayers()
	t.Cleanup(free)
	type row struct {
		Enabled bool
		Pattern int
		Status  uint32
		Count   int
	}
	var rows []row
	for _, enabled := range []bool{false, true} {
		flags := noxflags.GameFlag(0)
		if enabled {
			flags = 0x8000
		}
		restore := noxflags.PortTestGameFlags(flags)
		t.Cleanup(restore)
		for pattern := 0; pattern < 6; pattern++ {
			for high := uint32(0); high < 2; high++ {
				for status := uint32(0); status < 256; status++ {
					bits := status | high<<31
					want := 0
					for i := range players {
						p := &players[i]
						p.PlayerInd = byte(i)
						p.Active = 0
						p.Field3680 = bits
						active := pattern == 1 || pattern == 2 && i%2 == 0 || pattern == 3 && i%2 == 1 || pattern == 4 && i == 31 || pattern == 5 && i%7 == 0
						if active {
							p.Active = []byte{1, 128, 255}[i%3]
						}
						if enabled && active && i != 31 && bits&0x21 == 1 {
							want++
						}
					}
					got := legacy.Nox_xxx_countObserverPlayers_425BF0()
					if got != want {
						t.Fatal("observer count", enabled, pattern, bits, got, want)
					}
					for i := range players {
						if players[i].Field3680 != bits || players[i].PlayerInd != byte(i) {
							t.Fatal("observer traversal mutated player")
						}
					}
					rows = append(rows, row{enabled, pattern, bits, got})
				}
			}
		}
		restore()
	}
	interactionCapture(t, "server-runtime-observer-count", rows)
}
