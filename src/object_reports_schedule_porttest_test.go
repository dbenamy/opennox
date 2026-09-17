//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"testing"
)

func TestObjectReportsMinimapSchedule(t *testing.T) {
	s, units, _ := objectReportsPlayers(t)
	nodes, free := alloc.Make([]server.MinimapItem{}, 65)
	t.Cleanup(free)
	var rows []struct {
		Name   string
		Return int
	}
	defer func() {
		spellbookCapture(t, "object-reports-schedule", rows, "63550214ba3acf3a830ccb7687722e8c7f877336861cc49d80d3f626c016f90e")
	}()
	for _, count := range []int{0, 1, 2, 3, 7, 30, 59, 60, 61, 65} {
		for _, frame := range []uint32{0, 1, 60, 1000, 0xffffffff} {
			for _, elapsed := range []uint32{0, 1, 2, 3, 8, 20, 59, 60, 61, 0x80000000, 0xffffffff} {
				name := fmt.Sprintf("n%d/frame%d/elapsed%d", count, frame, elapsed)
				t.Run(name, func(t *testing.T) {
					clear(nodes)
					s.SetFrame(frame)
					for i := 0; i < count; i++ {
						nodes[i].Field8 = &nodes[(i+1)%count]
						nodes[i].Field12 = &nodes[(i+count-1)%count]
					}
					for i := range units {
						u := &units[i]
						ud := u.UpdateDataPlayer()
						ud.Player.Field4580 = nil
						if count != 0 {
							ud.Player.Field4580 = &nodes[0]
						}
						ud.Field67 = frame - elapsed
						got := legacy.PortTestObjectReports(7, u, nil, 0, 0, 0, nil)
						want := 0
						if count > 60 || count > 0 && elapsed > uint32(60/count) {
							want = 1
						}
						if got != want || ud.Field67 != frame-elapsed {
							t.Fatalf("slot%d got%d want%d or timestamp mutated", i, got, want)
						}
						rows = append(rows, struct {
							Name   string
							Return int
						}{fmt.Sprintf("%s/slot%d", name, ud.Player.PlayerIndex()), got})
					}
				})
			}
		}
	}
}
