//go:build porttest

package opennox

import (
	"encoding/binary"
	"fmt"
	"testing"

	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

type matchRosterWinnerRow struct {
	Name   string
	Return uint32
	Flags  uint32
	Queue  legacy.PortTestReliableReportState
}

func TestMatchRosterWinnerSelection(t *testing.T) {
	o := newMatchRosterOwner(t)
	var rows []matchRosterWinnerRow
	defer func() {
		spellbookCapture(t, "match-roster-winners", rows, "89e18421d5042d28bfc721fc3d8770d2098f81b65dcfc4efaf39c23c598521ab")
	}()
	scoreSets := [][3]int{{0, 1, 1}, {5, 2, 3}, {-2, -1, -1}, {-2147483648, 0, 2147483647}, {3, 3, 3}}
	for _, op := range []string{"winner-min", "winner-max"} {
		for count := 0; count <= 3; count++ {
			for _, scores := range scoreSets {
				for _, players := range []uint32{0, 1, 3, 7} {
					for _, affiliated := range []bool{false, true} {
						name := fmt.Sprintf("%s/teams%d/scores%v/players%x/affiliated%t", op, count, scores, players, affiliated)
						t.Run(name, func(t *testing.T) {
							resetFlags := noxflags.PortTestGameFlags(0)
							defer resetFlags()
							o.reset()
							o.s.Teams.Reset()
							o.s.Teams.ActiveCnt = 0
							eligible := map[uint32]int{}
							for i := 0; i < count; i++ {
								tm := o.s.Teams.Create(server.TeamID(i + 1))
								tm.Lessons = scores[i]
								eligible[0x10000+uint32(tm.IDVal)] = scores[i]
							}
							for i := range o.units {
								u := &o.units[i]
								u.TeamVal.ID = 0
								if affiliated && count > 0 {
									u.TeamVal.ID = 1
								}
								pl := u.UpdateDataPlayer().Player
								status := uint32(1)
								if players&(1<<i) != 0 {
									status = 0
								}
								objectXferSetWord(pl.C(), 3680, status)
								score := scores[(i+2)%3]
								objectXferSetWord(pl.C(), 2136, uint32(score))
								objectXferSetWord(pl.C(), 2140, uint32(score))
								if status == 0 && u.TeamVal.ID == 0 {
									eligible[uint32(uint16(u.NetCode))] = score
								}
							}
							rv := matchRosterCall(op, nil, nil, nil)
							state := o.state()
							if !noxflags.HasGame(8) {
								t.Fatal("match end flag")
							}
							if len(eligible) == 0 {
								if rv != 0 || len(state.Nodes) != 0 {
									t.Fatal("empty selection reported a winner")
								}
							} else {
								if rv != 1 || len(state.Nodes) != 1 {
									t.Fatal("winner report count", rv, len(state.Nodes))
								}
								data := state.Nodes[0].Data
								if len(data) != 8 || data[3] != 1 || binary.LittleEndian.Uint32(data[4:]) != o.s.Frame() {
									t.Fatal("winner format", data)
								}
								code := uint32(binary.LittleEndian.Uint16(data[1:]))
								best := int(2147483647)
								if op == "winner-max" {
									best = -2147483648
								}
								for _, score := range eligible {
									if op == "winner-min" && score < best || op == "winner-max" && score > best {
										best = score
									}
								}
								if data[0] == 89 && code == 0 {
									ties := 0
									for _, score := range eligible {
										if score == best {
											ties++
										}
									}
									if ties < 2 {
										t.Fatal("draw without tied best scores")
									}
								} else {
									if data[0] == 89 {
										code += 0x10000
									} else if data[0] != 88 {
										t.Fatal("winner opcode", data[0])
									}
									score, ok := eligible[code]
									if !ok || score != best {
										t.Fatal("selected result is not an eligible best score", code, score, best)
									}
								}
							}
							rows = append(rows, matchRosterWinnerRow{name, rv, uint32(noxflags.GetGame()), state})
						})
					}
				}
			}
		}
	}
}
func TestMatchRosterFlagWinnerSelection(t *testing.T) {
	o := newMatchRosterOwner(t)
	var rows []matchRosterWinnerRow
	defer func() {
		spellbookCapture(t, "match-roster-flag-winners", rows, "8df3f438d14fa6a9e9141a6546509bd97e3aba42c0cc6dbe4656ac34f19a26e8")
	}()
	for _, mode := range []uint32{32, 64, 96} {
		for count := 0; count <= 3; count++ {
			for _, scores := range [][3]int{{0, 1, 1}, {5, 2, 3}, {-2, -1, -1}, {-2147483648, 0, 2147483647}, {3, 3, 3}} {
				name := fmt.Sprintf("mode%x/teams%d/scores%v", mode, count, scores)
				t.Run(name, func(t *testing.T) {
					restore := noxflags.PortTestGameFlags(noxflags.GameFlag(mode))
					defer restore()
					o.reset()
					o.s.Teams.Reset()
					o.s.Teams.ActiveCnt = 0
					best := -1
					winner := uint16(0)
					ties := 0
					for i := 0; i < count; i++ {
						tm := o.s.Teams.Create(server.TeamID(i + 1))
						tm.Lessons = scores[i]
						if scores[i] > best {
							best = scores[i]
							winner = uint16(i + 1)
							ties = 1
						} else if scores[i] == best {
							winner = uint16(i + 1)
							ties++
						}
					}
					if ties > 1 {
						winner = 0
					}
					rv := matchRosterCall("winner-flag", nil, nil, nil)
					state := o.state()
					if rv != 1 || len(state.Nodes) != 1 {
						t.Fatal("flag result count")
					}
					want := []byte{87, 255, 255, 1, 0, 0, 0, 0}
					if winner != 0 {
						binary.LittleEndian.PutUint16(want[1:], winner)
						if mode&64 != 0 {
							want[0] = 86
							want[3] = 0
						}
					}
					binary.LittleEndian.PutUint32(want[4:], o.s.Frame())
					if string(state.Nodes[0].Data) != string(want) {
						t.Fatal("flag result", state.Nodes[0].Data, want)
					}
					rows = append(rows, matchRosterWinnerRow{name, rv, uint32(noxflags.GetGame()), state})
				})
			}
		}
	}
}
