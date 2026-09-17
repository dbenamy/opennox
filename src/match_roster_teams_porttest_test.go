//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unicode/utf16"
)

func TestMatchRosterTeamRoster(t *testing.T) {
	o := newMatchRosterOwner(t)
	var rows []matchRosterMessageRow
	defer func() {
		spellbookCapture(t, "match-roster-team-roster", rows, "867df32de94332f6cfda59532899f5d0b702b9ed20d54a76ca572eb654b48777")
	}()
	for _, to := range []uint32{1, 7, 31, 159, 255} {
		for count := 0; count <= 3; count++ {
			for _, mode := range []uint32{1, 513} {
				label := fmt.Sprintf("to%d/count%d/mode%x", to, count, mode)
				t.Run(label, func(t *testing.T) {
					o.reset()
					restore := noxflags.PortTestGameFlags(noxflags.GameFlag(mode))
					defer restore()
					o.s.Teams.Reset()
					o.s.Teams.ActiveCnt = 0
					for i := range o.units {
						u := &o.units[i]
						u.TeamVal = server.ObjectTeam{}
						u.UpdateDataPlayer().Player.NetCodeVal = u.NetCode
						if count > 0 {
							u.TeamVal.ID = server.TeamID(1 + i%count)
						}
					}
					var want [][]byte
					for i := 0; i < count; i++ {
						tm := o.s.Teams.Create(server.TeamID(i + 1))
						name := []string{"", "A", "Team Ω"}[i]
						tm.SetNameAnd68(name, i+3)
						tm.Lessons = i - 2
						objectXferSetWord(tm.C(), 60, uint32(0x80000000)+uint32(i))
						for j := range o.units {
							u := &o.units[j]
							if u.TeamVal.ID == tm.IDVal {
								u.TeamVal.ID = 0
								legacy.Nox_xxx_createAtImpl_4191D0(tm.IDVal, &u.TeamVal, 1, int(u.NetCode), 0)
							}
						}
						data := make([]byte, 18)
						data[0] = 196
						binary.LittleEndian.PutUint32(data[2:], uint32(tm.IDVal))
						binary.LittleEndian.PutUint32(data[6:], uint32(0x80000000)+uint32(i))
						binary.LittleEndian.PutUint32(data[10:], uint32(tm.Lessons))
						if mode&512 != 0 {
							data[14] = 1
						}
						text := utf16.Encode([]rune(name))
						data[15] = byte(len(text))
						data[16] = byte(tm.ColorInd)
						data[17] = byte(i + 3)
						for _, c := range text {
							data = append(data, byte(c), byte(c>>8))
						}
						want = append([][]byte{data}, want...)
						for j := range o.units {
							u := &o.units[j]
							if u.TeamVal.ID != tm.IDVal {
								continue
							}
							data := []byte{196, 1, byte(tm.IDVal), 0, 0, 0, 0, 0, 0, 0}
							binary.LittleEndian.PutUint16(data[6:], uint16(u.NetCode))
							binary.LittleEndian.PutUint16(data[8:], u.TypeInd)
							want = append([][]byte{data}, want...)
						}
					}
					if to == 31 {
						want = nil
					}
					o.reset()
					rv := matchRosterCall("team-roster", nil, nil, nil, to)
					state := o.state()
					if rv != 0 || len(state.Nodes) != len(want) {
						t.Fatal("team roster count", rv, len(state.Nodes), len(want))
					}
					for i, data := range want {
						if !bytes.Equal(state.Nodes[i].Data, data) {
							t.Fatalf("team roster payload %d: %x want %x", i, state.Nodes[i].Data, data)
						}
					}
					rows = append(rows, matchRosterMessageRow{label, rv, state})
				})
			}
		}
	}
}

func TestMatchRosterTeamAssignment(t *testing.T) {
	o := newMatchRosterOwner(t)
	ball := legacy.PortTestMapDrawableTeamWord()
	oldBall := *ball
	*ball = 0
	t.Cleanup(func() { *ball = oldBall })
	type teamRow struct {
		ID      byte
		Group   uint32
		Name    string
		Members uint32
	}
	type row struct {
		Name     string
		Assigned byte
		Teams    []teamRow
		Queue    legacy.PortTestReliableReportState
	}
	var rows []row
	defer func() {
		spellbookCapture(t, "match-roster-team-assignment", rows, "68d1ddb58d0504dd589b59dd6f70a869f15dfa696bfd12a1d1b61c49cc919a46")
	}()
	for _, mode := range []uint32{1, 0x81, 0x8001, 0x8021, 0x8041, 0x8011} {
		for count := 0; count <= 3; count++ {
			for _, group := range []uint32{0, 7, 0x80000001} {
				for _, cap := range []uint32{0, 1, 2, 4} {
					for _, assigned := range []bool{false, true} {
						label := fmt.Sprintf("mode%x/count%d/group%x/cap%d/assigned%t", mode, count, group, cap, assigned)
						t.Run(label, func(t *testing.T) {
							o.reset()
							restore := noxflags.PortTestGameFlags(noxflags.GameFlag(mode))
							defer restore()
							o.s.Teams.Reset()
							o.s.Teams.ActiveCnt = 0
							*memmap.PtrUint8(0x5D4594, 371568) = byte(cap)
							*o.roster["team-cap"] = cap / 2
							*memmap.PtrUint8(0x5D4594, 371433) = 0
							if mode == 0x81 {
								*memmap.PtrUint8(0x5D4594, 371433) = 128
							}
							noxflags.UnsetGamePlay(^noxflags.GameplayFlag(0))
							noxflags.SetGamePlay(4)
							for i := range o.units {
								u := &o.units[i]
								u.TeamVal = server.ObjectTeam{}
								u.UpdateDataPlayer().Player.NetCodeVal = u.NetCode
							}
							for i := 0; i < count; i++ {
								tm := o.s.Teams.Create(server.TeamID(i + 1))
								tm.SetNameAnd68(fmt.Sprintf("Old%d", i), 0)
								if i > 0 {
									objectXferSetWord(tm.C(), 60, uint32(6+i))
								}
							}
							if count > 0 {
								legacy.Nox_xxx_createAtImpl_4191D0(1, &o.units[1].TeamVal, 1, int(o.units[1].NetCode), 0)
								if assigned {
									legacy.Nox_xxx_createAtImpl_4191D0(1, &o.units[0].TeamVal, 1, int(o.units[0].NetCode), 0)
								}
							}
							o.reset()
							pl := o.units[0].UpdateDataPlayer().Player
							pl.Field2068 = group
							clear(pl.Field2072[:])
							copy(pl.Field2072[:], utf16.Encode([]rune("Group")))
							expected := byte(o.units[0].TeamVal.ID)
							if mode&0x8000 == 0 && mode != 0x81 {
								if expected == 0 && count > 0 {
									expected = 1
									if count > 1 {
										expected = 2
									}
								}
							} else if group != 0 {
								found := byte(0)
								for i := 1; i < count; i++ {
									if uint32(6+i) == group {
										found = byte(i + 1)
									}
								}
								if found != 0 {
									expected = found
								} else {
									max := cap
									if (mode&96 != 0 || mode&16 != 0) && max > 2 {
										max = 2
									}
									used := count - 1
									if used < 0 {
										used = 0
									}
									if uint32(used) < max && (mode&96 == 0 || uint32(used) < cap/2) {
										expected = 1
									}
								}
							}
							matchRosterCall("assign-team", nil, pl, nil)
							if byte(o.units[0].TeamVal.ID) != expected {
								t.Fatal("assigned team", o.units[0].TeamVal.ID, expected)
							}
							teams := []teamRow{}
							for tm := o.s.Teams.First(); tm != nil; tm = o.s.Teams.Next(tm) {
								teams = append(teams, teamRow{byte(tm.IDVal), uint32(tm.Ind60()), tm.Name(), objectXferGetWord(tm.C(), 48)})
							}
							rows = append(rows, row{label, byte(o.units[0].TeamVal.ID), teams, o.state()})
						})
					}
				}
			}
		}
	}
	// Team links are embedded in the fixture units; clear them before teardown.
	for i := range o.units {
		o.units[i].TeamVal = server.ObjectTeam{}
	}
}

func TestMatchRosterUnitlessAssignment(t *testing.T) {
	o := newMatchRosterOwner(t)
	pl := o.s.Players.ByInd(3)
	if pl == nil || pl.PlayerUnit != nil {
		t.Fatal("unitless roster slot")
	}
	for _, mode := range []uint32{1, 0x81, 0x8001, 0x8021} {
		restore := noxflags.PortTestGameFlags(noxflags.GameFlag(mode))
		o.reset()
		pl.Field2068 = 7
		matchRosterCall("assign-team", nil, pl, nil)
		if o.s.Teams.Count() != 0 || len(o.state().Nodes) != 0 {
			t.Fatal("unitless assignment mutated state")
		}
		restore()
	}
}
