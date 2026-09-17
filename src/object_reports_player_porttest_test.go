//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/prand"
	"github.com/opennox/opennox/v1/internal/netlist"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

func TestObjectReportsPlayerAnimation(t *testing.T) {
	s, units, _ := objectReportsPlayers(t)
	t.Cleanup(s.PortTestObjectReportAnimations(3, 1))
	a, u := &units[0], &units[1]
	ud := u.UpdateDataPlayer()
	pl := ud.Player
	var rows []struct {
		Name     string
		Return   int
		Packet   []byte
		Reaction uint32
		Random   int
		Reported uint32
	}
	defer func() {
		spellbookCapture(t, "object-reports-player-animation", rows, "62a6ef467578fd4f7d37eb8878c5b9e0ab0722fa55b508c2f7789612e52bffdd")
	}()
	states := map[int]byte{0: 4, 1: 0, 2: 21, 3: 1, 4: 2, 5: 6, 10: 21, 12: 3, 13: 0, 14: 0, 15: 40, 16: 40, 17: 40, 18: 48, 19: 49, 20: 47, 21: 30, 22: 0, 23: 50, 24: 19, 25: 20, 26: 15, 27: 16, 28: 16, 29: 16, 30: 52, 32: 54}
	animated := map[int]bool{1: true, 10: true, 2: true, 15: true, 16: true, 17: true, 14: true, 20: true, 18: true, 19: true, 21: true, 22: true, 24: true, 25: true, 27: true, 28: true, 29: true, 26: true, 30: true, 32: true}
	for state := 0; state < 34; state++ {
		for _, confused := range []bool{false, true} {
			for _, weapon := range []uint32{0, 16} {
				for _, elapsed := range []uint32{0, 1, 2, 3, 4, 5, 0xffffffff} {
					for _, seed := range []int{1, 7, 12345} {
						name := fmt.Sprintf("state%d/confused%t/weapon%d/elapsed%d/seed%d", state, confused, weapon, elapsed, seed)
						t.Run(name, func(t *testing.T) {
							*ud = server.PlayerUpdateData{Player: pl}
							s.SetFrame(100)
							s.NetList.ResetAll()
							s.Rand.Logic = prand.New(seed)
							expected := prand.New(seed)
							u.NetCode = 123
							u.Field131 = weapon
							u.Direction1 = 0
							ud.State = server.PlayerState(state)
							*(*byte)(unsafe.Add(u.UpdateData, 236)) = 91
							*(*uint32)(unsafe.Add(u.UpdateData, 164)) = 100 - elapsed
							if confused {
								*(*uint16)(unsafe.Add(u.UpdateData, 160)) = 1
							}
							rv := legacy.PortTestObjectReports(5, a, u, 1, 1, 0, nil)
							got := s.NetList.CopyPacketsA(1, netlist.Kind2)
							anim := states[state]
							frame := byte(255)
							if animated[state] {
								frame = 91
							}
							if confused && expected.IntClamp(0, 10) >= 8 {
								anim = 50
								frame = 255
							}
							if (state == 3 || state == 4) && weapon == 16 {
								anim = 51
								frame = 255
							}
							reaction := uint32(100 - elapsed)
							if state != 30 {
								if elapsed/2 >= 3 || elapsed >= 4 {
									reaction = 0
								} else {
									anim = 52
									frame = byte(elapsed / 2)
								}
							}
							if rv != 1 || len(got) != 12 {
								t.Fatalf("return%d packet%x", rv, got)
							}
							if got[0] != 195 || got[9] != 64 || got[10] != frame || got[11] != anim {
								t.Fatalf("packet%x anim%d frame%d", got, anim, frame)
							}
							stored := *(*uint32)(unsafe.Add(u.UpdateData, 164))
							reported := *(*uint32)(unsafe.Add(unsafe.Pointer(a.UpdateDataPlayer().Player), 4452+4*7))
							if stored != reaction || s.Rand.Logic.Index() != expected.Index() || reported != 100 {
								t.Fatal("reaction/random/recipient timestamp")
							}
							rows = append(rows, struct {
								Name     string
								Return   int
								Packet   []byte
								Reaction uint32
								Random   int
								Reported uint32
							}{name, rv, got, stored, s.Rand.Logic.Index(), reported})
						})
					}
				}
			}
		}
	}
}
