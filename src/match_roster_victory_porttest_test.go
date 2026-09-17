//go:build porttest

package opennox

import (
	"encoding/binary"
	"fmt"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

func TestMatchRosterVictoryLimits(t *testing.T) {
	o := newMatchRosterOwner(t)
	var rows []matchRosterWinnerRow
	defer func() {
		spellbookCapture(t, "match-roster-victory", rows, "697fbdb316d558744ee3c57d06ef7f793c7cb67cd9e0b1a6d01d71cb43579fad")
	}()
	for _, mode := range []uint32{0, 16, 32, 64, 256, 512, 1024, 1536} {
		for _, limit := range []uint16{0, 1, 2, 32768, 65535} {
			for _, scores := range [][3]int32{{0, 0, 0}, {1, 2, 3}, {3, 2, 1}, {-1, 2147483647, -2147483648}} {
				for _, status := range []uint32{0, 1, 0x20, 0x21} {
					for _, teams := range []bool{false, true} {
						label := fmt.Sprintf("mode%x/limit%d/scores%v/status%x/teams%t", mode, limit, scores, status, teams)
						t.Run(label, func(t *testing.T) {
							o.reset()
							restore := noxflags.PortTestGameFlags(noxflags.GameFlag(mode))
							defer restore()
							o.s.Teams.Reset()
							o.s.Teams.ActiveCnt = 0
							for i := range unsafe.Slice(memmap.PtrUint16(0x5D4594, 3488), 6) {
								*memmap.PtrUint16(0x5D4594, 3488+uintptr(i*2)) = limit
							}
							for i := range o.units {
								u := &o.units[i]
								u.TeamVal = server.ObjectTeam{}
								pl := u.UpdateDataPlayer().Player
								pl.Lessons = scores[i]
								pl.Field2140 = uint32(scores[i])
								pl.Field3680 = status
								if teams {
									tm := o.s.Teams.Create(server.TeamID(i + 1))
									tm.Lessons = int(scores[i])
									u.TeamVal.ID = tm.IDVal
								}
							}
							rv := matchRosterCall("check-victory", nil, nil, nil)
							state := o.state()
							eligible := []int{}
							for i, score := range scores {
								if mode&1024 != 0 {
									if status&1 == 0 && uint32(score) < uint32(limit) {
										eligible = append(eligible, i)
									}
								} else if (teams || status&1 == 0) && score >= int32(limit) {
									eligible = append(eligible, i)
								}
							}
							ended := false
							if limit != 0 {
								if mode&1024 != 0 {
									ended = len(eligible) < 2 && (status&1 == 0 || status&0x20 != 0)
								} else {
									ended = mode&512 == 0 && len(eligible) > 0
								}
							}
							if noxflags.HasGame(8) != ended || len(state.Nodes) != bool2int(ended) {
								t.Fatal("victory gate", noxflags.GetGame(), len(state.Nodes), ended)
							}
							if ended {
								data := state.Nodes[0].Data
								wantCode := uint16(0)
								wantOpcode := byte(88)
								if len(eligible) > 0 {
									i := eligible[0]
									wantCode = uint16(o.units[i].NetCode)
									if teams {
										wantCode = uint16(i + 1)
										wantOpcode = 89
									}
								}
								if len(data) != 8 || data[0] != wantOpcode || binary.LittleEndian.Uint16(data[1:]) != wantCode || data[3] != 0 || binary.LittleEndian.Uint32(data[4:]) != o.s.Frame() {
									t.Fatal("victory result", data)
								}
							}
							rows = append(rows, matchRosterWinnerRow{label, rv, uint32(noxflags.GetGame()), state})
						})
					}
				}
			}
		}
	}
}
