//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/prand"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"testing"
)

func TestTeamRuntimeBalance(t *testing.T) {
	o := newMatchRosterOwner(t)
	oldRNG := o.s.Rand.Logic
	t.Cleanup(func() { o.s.Rand.Logic = oldRNG })
	oldHost := legacy.ClientPlayerNetCode()
	t.Cleanup(func() { legacy.ClientSetPlayerNetCode(oldHost) })
	type row struct {
		Name   string
		IDs    []byte
		Counts [2]int
		RNG    int
		Queue  legacy.PortTestReliableReportState
	}
	var rows []row
	defer func() {
		spellbookCapture(t, "team-runtime-balance", rows, "0ca54e7884d56b30742ae0b586da579ab9725c5a4fda1162eedb016b1c60a106")
	}()
	for mask := 0; mask < 8; mask++ {
		for _, observerAllowed := range []bool{false, true} {
			for _, headless := range []bool{false, true} {
				for _, seed := range []int{0, 7, 99} {
					name := fmt.Sprintf("mask%d/observer%t/headless%t/seed%d", mask, observerAllowed, headless, seed)
					t.Run(name, func(t *testing.T) {
						defer noxflags.PortTestGameFlags(1)()
						noxflags.ResetEngine()
						if headless {
							noxflags.SetEngine(noxflags.EngineNoRendering)
						}
						o.s.Teams.Reset()
						o.s.Teams.ActiveCnt = 0
						o.s.Teams.Create(1)
						o.s.Teams.Create(2)
						legacy.ClientSetPlayerNetCode(int(o.units[0].NetCode))
						var eligible [3]bool
						count := 0
						for i := range o.units {
							u := &o.units[i]
							u.TeamVal = server.ObjectTeam{}
							pl := u.UpdateDataPlayer().Player
							pl.NetCodeVal = u.NetCode
							pl.Field3680 = 0
							if mask&(1<<i) != 0 {
								pl.Field3680 = 1
								if observerAllowed {
									pl.Field3680 |= 0x20
								}
							}
							eligible[i] = (mask&(1<<i) == 0 || observerAllowed) && !(headless && i == 0)
							if eligible[i] {
								count++
							}
						}
						o.reset()
						o.s.Rand.Logic = prand.New(seed)
						legacy.PortTestTeamRuntimeMap("balance", 0)
						r := row{Name: name, RNG: o.s.Rand.Logic.Index(), Queue: o.state()}
						wantRNG := prand.New(seed)
						if count > 1 {
							for i := 0; i < 100; i++ {
								wantRNG.Int(0, 1)
							}
						}
						if r.RNG != wantRNG.Index() {
							t.Fatalf("random consumption %d want %d", r.RNG, wantRNG.Index())
						}
						for i := range o.units {
							id := byte(o.units[i].TeamVal.ID)
							r.IDs = append(r.IDs, id)
							if !eligible[i] {
								if id != 0 {
									t.Fatal("ineligible player assigned")
								}
								continue
							}
							if id != 1 && id != 2 {
								t.Fatal("eligible player unassigned")
							}
							r.Counts[id-1]++
						}
						if r.Counts[0]+r.Counts[1] != count || r.Counts[0] < r.Counts[1] || r.Counts[0] > r.Counts[1]+1 {
							t.Fatal("unbalanced teams", r.Counts, count)
						}
						rows = append(rows, r)
					})
				}
			}
		}
	}
}
