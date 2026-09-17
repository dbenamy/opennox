//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/prand"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"testing"
)

func TestTeamRuntimeGroupBalance(t *testing.T) {
	o := newMatchRosterOwner(t)
	defer noxflags.PortTestGameFlags(0x81)()
	*memmap.PtrUint8(0x5D4594, 371433) = 128
	groups := []uint32{0, 7, 0x80000001, 99}
	oldRNG := o.s.Rand.Logic
	t.Cleanup(func() { o.s.Rand.Logic = oldRNG })
	type row struct {
		Choice int
		Reset  bool
		IDs    []byte
		Counts [2]uint32
		RNG    int
		Queue  legacy.PortTestReliableReportState
	}
	var rows []row
	for choice := 0; choice < 64; choice++ {
		for _, reset := range []bool{false, true} {
			t.Run(fmt.Sprintf("choice%d/reset%t", choice, reset), func(t *testing.T) {
				o.s.Teams.Reset()
				o.s.Teams.ActiveCnt = 0
				a := o.s.Teams.Create(1)
				b := o.s.Teams.Create(2)
				objectXferSetWord(a.C(), 60, 7)
				objectXferSetWord(b.C(), 60, 0x80000001)
				v := choice
				var want [3]byte
				var counts [2]uint32
				for i := range o.units {
					u := &o.units[i]
					u.TeamVal = server.ObjectTeam{}
					pl := u.UpdateDataPlayer().Player
					pl.NetCodeVal = u.NetCode
					pl.Field3680 = 0
					pl.Field2068 = groups[v%4]
					v /= 4
					if pl.Field2068 == 7 {
						want[i] = 1
						counts[0]++
					}
					if pl.Field2068 == 0x80000001 {
						want[i] = 2
						counts[1]++
					}
					if reset {
						legacy.Nox_xxx_createAtImpl_4191D0(1, u.TeamPtr(), 0, int(u.NetCode), 0)
					}
				}
				o.reset()
				o.s.Rand.Logic = prand.New(7)
				arg := 0
				if reset {
					arg = 1
				}
				legacy.PortTestTeamRuntimeMap("balance", arg)
				r := row{Choice: choice, Reset: reset, Counts: [2]uint32{objectXferGetWord(a.C(), 48), objectXferGetWord(b.C(), 48)}, RNG: o.s.Rand.Logic.Index(), Queue: o.state()}
				for i := range o.units {
					got := byte(o.units[i].TeamVal.ID)
					if got != want[i] {
						t.Fatal("group mapping", i, got, want[i])
					}
					r.IDs = append(r.IDs, got)
				}
				if r.Counts != counts {
					t.Fatalf("membership counters %v do not match assigned members %v", r.Counts, counts)
				}
				rng := prand.New(7)
				for i := 0; i < 100; i++ {
					rng.Int(0, 1)
				}
				if r.RNG != rng.Index() {
					t.Fatal("group shuffle random consumption")
				}
				rows = append(rows, r)
			})
		}
	}
	spellbookCapture(t, "team-runtime-group-balance", rows, "4180ac8575ac5469aafe61690045e978cbc81087ca76637627e2644a4b9b1e2c")
}
