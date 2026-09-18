//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/types"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/ntype"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"math"
	"testing"
)

func TestSessionEntryPlayerReady(t *testing.T) {
	o := newMatchRosterOwner(t)
	type row struct {
		Index    int
		Flags    uint32
		Observer bool
		Status   uint32
		Position [2]uint32
		Buffs    uint32
		Duration uint16
		Power    byte
		Queue    legacy.PortTestReliableReportState
	}
	var rows []row
	for _, flags := range []uint32{0, 512, 8192, 8192 | 128, 4096} {
		for _, index := range []int{-1, 0, 1, 3, 7, 31, 32} {
			if index == 31 && flags&128 != 0 {
				continue
			} // Local cooperative startup opens its own GUI owner.
			for _, observer := range []bool{false, true} {
				if observer && (index == 3 || flags&512 != 0) {
					continue
				}
				t.Run(fmt.Sprintf("flags%x/index%d/observer%v", flags, index, observer), func(t *testing.T) {
					defer noxflags.PortTestGameFlags(noxflags.GameFlag(flags))()
					o.reset()
					r := row{Index: index, Flags: flags, Observer: observer}
					var pl *server.Player
					if index >= 0 && index < 32 {
						pl = o.s.Players.ByInd(ntype.PlayerInd(index))
					}
					if pl != nil {
						pl.Field3680 = 0
						if observer {
							pl.Field3680 = 1
						}
						objectXferSetWord(pl.C(), 3632, 0x11223344)
						objectXferSetWord(pl.C(), 3636, 0x55667788)
						if u := pl.PlayerUnit; u != nil {
							u.PosVec = types.Pointf{X: 12.5, Y: -3.25}
							u.ObjFlags = 0
							u.Buffs = 0
							clear(u.BuffsDur[:])
							clear(u.BuffsPower[:])
						}
					}
					legacy.Nox_xxx_gameServerReadyMB_4DD180(index)
					if pl != nil {
						r.Status = pl.Field3680
						r.Position = [2]uint32{objectXferGetWord(pl.C(), 3632), objectXferGetWord(pl.C(), 3636)}
						wantStatus := uint32(16)
						if observer {
							wantStatus |= 1
						}
						if r.Status != wantStatus {
							t.Fatalf("status%x want%x", r.Status, wantStatus)
						}
						wantPos := [2]uint32{0x11223344, 0x55667788}
						if observer {
							wantPos = [2]uint32{math.Float32bits(12.5), math.Float32bits(-3.25)}
						}
						if r.Position != wantPos {
							t.Fatal("ready position")
						}
						if u := pl.PlayerUnit; u != nil {
							r.Buffs = u.Buffs
							r.Duration = u.BuffsDur[23]
							r.Power = u.BuffsPower[23]
							wantBuff := flags&8192 != 0 && flags&128 == 0
							if wantBuff && (r.Buffs != 1<<23 || r.Duration != 150 || r.Power != 5) || !wantBuff && (r.Buffs != 0 || r.Duration != 0 || r.Power != 0) {
								t.Fatalf("ready buff %+v", r)
							}
						}
					}
					r.Queue = o.state()
					rows = append(rows, r)
				})
			}
		}
	}
	spellbookCapture(t, "session-entry-player-ready", rows, "e4fd775d3d66977001d334542d21dd64ad7d2ae7558b60ffc940198836acac56")
}
