//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/object"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/ntype"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

func TestSessionEntryPlayerDeparture(t *testing.T) {
	type row struct {
		Flags                                 uint32
		Index                                 int
		Roster                                []int
		UnitCleared, Deleted, LifetimeMatches bool
		Events                                int
		Queue                                 legacy.PortTestReliableReportState
	}
	var rows []row
	for _, flags := range []uint32{0, 2, 1024, 4096} {
		for _, index := range []int{1, 7, 31} {
			t.Run(fmt.Sprintf("flags%x/index%d", flags, index), func(t *testing.T) {
				o := newMatchRosterOwner(t)
				defer noxflags.PortTestGameFlags(noxflags.GameFlag(flags))()
				t.Cleanup(legacy.PortTestSessionEntryRosterOwner())
				pl := o.s.Players.ByInd(ntype.PlayerInd(index))
				u := pl.PlayerUnit
				legacy.PortTestSessionEntryRoster("add", index)
				events := 0
				o.s.OnPlayerLeave(func(p *server.Player) {
					events++
					if p != pl || p.PlayerUnit != u || p.Active == 0 {
						t.Error("departure callback must precede player cleanup")
					}
				})
				oldDeleted := o.s.Objs.DeletedList
				t.Cleanup(func() { o.s.Objs.DeletedList = oldDeleted; u.DeletedNext = nil })
				for i := range o.units {
					d := o.units[i].UpdateData
					objectXferSetWord(d, 324+4*index, 0x12345678)
					for _, off := range []int{452, 484, 516} {
						*(*byte)(unsafe.Add(d, off+index)) = 0xa5
					}
				}
				o.reset()
				legacy.Nox_xxx_playerForceDisconnect_4DE7C0(ntype.PlayerInd(index))
				wantActive := uint32(0)
				if flags&2 != 0 && !isDedicatedServer {
					wantActive = 1
				}
				r := row{Flags: flags, Index: index, Roster: legacy.PortTestSessionEntryRosterMembers(), UnitCleared: pl.PlayerUnit == nil, Deleted: u.Flags().Has(object.FlagDestroyed), LifetimeMatches: uint32(pl.Active) == wantActive, Events: events, Queue: o.state()}
				if !r.UnitCleared || !r.Deleted || !r.LifetimeMatches || r.Events != 1 || len(r.Roster) != 0 {
					t.Fatalf("departure state %+v", r)
				}
				for i := range o.units {
					if &o.units[i] == u {
						continue
					}
					d := o.units[i].UpdateData
					if objectXferGetWord(d, 324+4*index) != 0 {
						t.Fatal("peer reference word")
					}
					for _, off := range []int{452, 484, 516} {
						if *(*byte)(unsafe.Add(d, off+index)) != 0 {
							t.Fatal("peer reference byte")
						}
					}
				}
				rows = append(rows, r)
			})
		}
	}
	spellbookCapture(t, "session-entry-player-departure", rows, "6c03b5abf6b84b11620fb8e901fbb22dd6076489c60ae46e3e24c0e8c642e7d9")
}
