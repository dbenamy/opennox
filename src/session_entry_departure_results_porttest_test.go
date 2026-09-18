//go:build porttest

package opennox

import (
	"fmt"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/ntype"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestSessionEntryDepartureResults(t *testing.T) {
	type row struct {
		Remaining   int
		Elimination bool
		Flags       uint32
		Lessons     [2]int32
		Score       [2]uint32
		Status      [2]uint32
		ResetFrame  uint32
		Queue       legacy.PortTestReliableReportState
	}
	var rows []row
	for remaining := 0; remaining <= 2; remaining++ {
		for _, elimination := range []bool{false, true} {
			t.Run(fmt.Sprintf("remaining%d/elimination%v", remaining, elimination), func(t *testing.T) {
				o := newMatchRosterOwner(t)
				t.Cleanup(legacy.PortTestSessionEntryRosterOwner())
				flags := uint32(0x4000000)
				if elimination {
					flags |= 1024
				}
				defer noxflags.PortTestGameFlags(noxflags.GameFlag(flags))()
				o.balance(map[string]float64{"SuddenDeathPlayerThreshold": 3, "SuddenDeathLifeTime": 9})
				for _, off := range []uintptr{3520, 3536, 3476, 1392} {
					clear(serverConfigOwnBytes(t, 0x5D4594, off, 4))
				}
				o.s.Players.ByInd(3).Active = 0
				for i, slot := range []ntype.PlayerInd{7, 31} {
					p := o.s.Players.ByInd(slot)
					if i >= remaining {
						p.Active = 0
					}
					p.Lessons = int32(11 + i)
					p.Field2140 = uint32(21 + i)
					p.Field3680 = 256
					if remaining == 2 && i == 1 {
						p.Field3680 = 33
					} // Participates in the match; ineligible to win.
				}
				oldDeleted := o.s.Objs.DeletedList
				t.Cleanup(func() { o.s.Objs.DeletedList = oldDeleted; o.units[0].DeletedNext = nil })
				oldClosing := nox_xxx_serverIsClosing_825764
				nox_xxx_serverIsClosing_825764 = false
				t.Cleanup(func() { nox_xxx_serverIsClosing_825764 = oldClosing })
				o.reset()
				legacy.Nox_xxx_playerForceDisconnect_4DE7C0(1)
				r := row{Remaining: remaining, Elimination: elimination, Flags: uint32(noxflags.GetGame()), ResetFrame: memmap.Uint32(0x5D4594, 3520), Queue: o.state()}
				for i, slot := range []ntype.PlayerInd{7, 31} {
					p := o.s.Players.ByIndRaw(slot)
					r.Lessons[i] = p.Lessons
					r.Score[i] = p.Field2140
					r.Status[i] = p.Field3680
					if remaining < 2 && i < remaining {
						if p.Lessons != 0 || p.Field2140 != 0 || p.Field3680&256 != 0 {
							t.Fatal("remaining player reset")
						}
					}
				}
				if remaining < 2 {
					if r.Flags&0x4000000 != 0 || r.ResetFrame != 123 {
						t.Fatal("match reset on departure")
					}
				} else {
					if r.Flags&0x4000000 == 0 || r.ResetFrame != 0 {
						t.Fatal("running match unexpectedly reset")
					}
					if (r.Flags&8 != 0) != elimination {
						t.Fatal("elimination winner transition")
					}
				}
				rows = append(rows, r)
			})
		}
	}
	spellbookCapture(t, "session-entry-departure-results", rows, "fb0743a3accb9c608ca55cb1db4f78339b2bcb99a128e45f47f0e6f6a4ce89d6")
}
