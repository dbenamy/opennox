//go:build porttest

package opennox

import (
	"fmt"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/ntype"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

func TestSessionEntryPlayerReset(t *testing.T) {
	type row struct {
		Count       int
		Mask, Flags uint32
		Stages      [3]byte
		Scores      [3][2]uint32
		HP          [3]uint16
		Queue       legacy.PortTestReliableReportState
	}
	var rows []row
	for count := 0; count <= 3; count++ {
		t.Run(fmt.Sprint(count), func(t *testing.T) {
			o := newMatchRosterOwner(t)
			oldAbilities := noxServer.abilities
			noxServer.abilities.Init(noxServer)
			t.Cleanup(func() { noxServer.abilities = oldAbilities })
			defer noxflags.PortTestGameFlags(0x20000)()
			for i, slot := range []ntype.PlayerInd{1, 7, 31} {
				pl := o.s.Players.ByInd(slot)
				u := pl.PlayerUnit
				if i >= count {
					pl.Active = 0
				}
				health := (*server.HealthData)(o.record(t, int(unsafe.Sizeof(server.HealthData{}))))
				health.Cur = 11
				health.Max = 99
				u.HealthData = health
				objectXferSetWord(pl.C(), 4700, 1) // Established player state retains its loadout; downstream item generation has separate contracts.
				objectXferSetWord(pl.C(), 2140, 17)
				objectXferSetWord(pl.C(), 2136, 19)
				*(*byte)(unsafe.Add(pl.C(), 3676)) = 8
			}
			o.reset()
			legacy.Nox_xxx_servResetPlayers_4D23C0()
			r := row{Count: count, Mask: *o.words["mask"], Flags: uint32(noxflags.GetGame()), Queue: o.state()}
			wantMask := uint32(0x80000082)
			for i, slot := range []ntype.PlayerInd{1, 7, 31} {
				pl := o.s.Players.ByIndRaw(slot)
				r.Stages[i] = *(*byte)(unsafe.Add(pl.C(), 3676))
				r.Scores[i] = [2]uint32{objectXferGetWord(pl.C(), 2140), objectXferGetWord(pl.C(), 2136)}
				r.HP[i] = o.units[i].HealthData.Cur
				if i < count {
					wantMask &^= 1 << uint(slot)
					if r.Stages[i] != 2 || r.Scores[i] != [2]uint32{} || r.HP[i] != 99 {
						t.Fatalf("active player%d %+v", slot, r)
					}
				} else if r.Stages[i] != 8 || r.Scores[i] != [2]uint32{17, 19} || r.HP[i] != 11 {
					t.Fatal("inactive player changed")
				}
			}
			if r.Mask != wantMask || r.Flags&0x20000 != 0 {
				t.Fatalf("reset mask/flags %+v", r)
			}
			rows = append(rows, r)
		})
	}
	spellbookCapture(t, "session-entry-player-reset", rows, "46cf498a147ca4c7360a79c29c2de4b61c5620e4ad58d996fa7efbad8821ee59")
}
