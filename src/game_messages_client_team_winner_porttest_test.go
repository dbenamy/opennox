//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/libs/object"
	"github.com/opennox/libs/strman"
	"github.com/opennox/opennox/v1/client"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"image"
	"slices"
	"testing"
	"unsafe"
)

func TestGameMessageClientTeamWinner(t *testing.T) {
	o := newInventoryTransactionOwner(t)
	o.reset(t)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	config := serverConfigOwnBytes(t, 0x5D4594, 371380, 58)
	hostConfig := serverConfigOwnBytes(t, 0x5D4594, 3488, 12)
	textbuf := serverConfigOwnBytes(t, 0x5D4594, 811376, 512)
	mode := serverConfigOwnBytes(t, 0x5D4594, 811060, 4)
	timer := serverConfigOwnBytes(t, 0x5D4594, 811908, 8)
	oldGate := legacy.Get_dword_5d4594_1200804()
	t.Cleanup(func() { legacy.Set_dword_5d4594_1200804(uint32(oldGate)) })
	oldCode := legacy.ClientPlayerNetCode()
	t.Cleanup(func() { legacy.ClientSetPlayerNetCode(oldCode) })
	legacy.ClientSetPlayerNetCode(7)
	oldTick := legacy.PlatformTicks
	t.Cleanup(func() { legacy.PlatformTicks = oldTick })
	legacy.PlatformTicks = func() uint64 { return 0x123456789 }
	oldUI := legacy.Sub_470510
	t.Cleanup(func() { legacy.Sub_470510 = oldUI })
	var ui []uint64
	legacy.Sub_470510 = func() { ui = append(ui, binary.LittleEndian.Uint64(timer)) }
	oldEngine := noxflags.GetEngine()
	t.Cleanup(func() { noxflags.ResetEngine(); noxflags.SetEngine(oldEngine) })
	t.Cleanup(noxflags.PortTestGameFlags(0))
	options, restore := legacy.PortTestServerOptionsWords()
	t.Cleanup(restore)
	for _, p := range options {
		*p = 0
	}
	formats := map[string]string{"TimeLimitReached": "Time;", "TeamWon": "Your team won", "teamformat": "Team %s", "FB_Victory": "%s won", "CTF_Victory": "Flag %s won", "CTF_Tie": "Flag tie", "DM_Loss": "Lost to %s", "DM_TeamVictory": "Your team won", "DM_Tie": "Tie", "HL_Tie": "Elimination tie", "HL_Header": "Elimination;", "HL_YourTeam": "Your team", "HL_Victory": "%s won", "Team": "Team %s"}
	var entries []strman.Entry
	for k, v := range formats {
		entries = append(entries, strman.Entry{ID: strman.ID("cdecode.c:" + k), Vals: []strman.Variant{{Str: v}}})
	}
	set, restoreStrings := o.c.srv.Server.PortTestMeterStrings(entries...)
	t.Cleanup(restoreStrings)
	set(0)
	units, freeUnits := alloc.Make([]server.Object{}, 4)
	t.Cleanup(freeUnits)
	update, freeUpdate := alloc.Make([]server.PlayerUpdateData{}, 3)
	t.Cleanup(freeUpdate)
	t.Cleanup(legacy.PortTestCreatureXferLookupOwner())
	oldList := o.c.srv.Objs.List
	o.c.srv.Objs.List = &units[0]
	t.Cleanup(func() { o.c.srv.Objs.List = oldList })
	for i := range o.players {
		o.players[i].Active = 0
	}
	for i := 0; i < 3; i++ {
		p := &o.players[i]
		p.Active = 1
		p.NetCodeVal = uint32(7 + i)
		p.PlayerInd = byte(i)
		p.PlayerUnit = &units[i]
		p.SetName([]string{"Ada", "Bob", "Observer"}[i])
		units[i].UpdateData = unsafe.Pointer(&update[i])
		update[i].Player = p
	}
	original := o.c.Objs.List1
	t.Cleanup(func() {
		for o.c.Objs.List1 != original {
			o.c.Nox_xxx_spriteDeleteStatic_45A4E0_drawable(o.c.Objs.List1)
		}
	})
	var sprites [4]*client.Drawable
	typ := o.c.Things.TypeByInd(4)
	oldClass := typ.ObjClass
	typ.ObjClass = object.ClassPlayer
	t.Cleanup(func() { typ.ObjClass = oldClass })
	for i := 0; i < 4; i++ {
		units[i].NetCode = uint32(7 + i)
		units[i].ObjClass = object.ClassPlayer
		if i < 3 {
			units[i].ObjNext = &units[i+1]
		}
		sprites[i] = o.c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(100+i, 200))
		sprites[i].NetCode32 = uint32(7 + i)
	}
	teams := []*server.Team{o.c.srv.Teams.ByID(1), o.c.srv.Teams.ByID(2)}
	for i, tm := range teams {
		if tm == nil {
			t.Fatal("winner team owner")
		}
		tm.SetNameAnd68([]string{"Ruby", "Azure"}[i], 0)
	}
	type row struct {
		Kind, On, Host, Elimination, Graphics, Reason, Code, Local int
		Frame, Gate, Limit, Score                                  uint32
		Text                                                       string
		Mode, Flags                                                uint32
		Time                                                       uint64
		UI                                                         []uint64
		Scores                                                     [3][2]uint32
		Teams                                                      [2]uint32
		Sounds                                                     [][2]int
	}
	var rows []row
	for _, kind := range []byte{86, 87, 89} {
		for on := 0; on < 2; on++ {
			for host := 0; host < 2; host++ {
				for elimination := 0; elimination < 2; elimination++ {
					for graphics := 0; graphics < 2; graphics++ {
						for _, reason := range []byte{0, 1, 2} {
							for _, code := range []uint16{0, 1, 2, 257, 65535} {
								for local := 0; local < 3; local++ {
									for _, time := range [][2]uint32{{100, 100}, {101, 100}, {0xffffffff, 0x80000000}} {
										for _, pair := range [][2]uint32{{0, 0}, {1, 1}, {65535, 0xffffffff}} {
											flags := noxflags.GameFlag(host | elimination*1024)
											reset := noxflags.PortTestGameFlags(flags)
											noxflags.ResetEngine()
											if graphics == 0 {
												noxflags.SetEngine(noxflags.EngineNoRendering)
											}
											binary.LittleEndian.PutUint32(connected, uint32(on))
											binary.LittleEndian.PutUint16(config[54:], uint16(pair[0]))
											for i := 0; i < 6; i++ {
												binary.LittleEndian.PutUint16(hostConfig[2*i:], uint16(pair[0]))
											}
											legacy.Set_dword_5d4594_1200804(time[1])
											clear(textbuf)
											copy(textbuf, []byte{'O', 0, 'l', 0, 'd', 0})
											binary.LittleEndian.PutUint32(mode, 9)
											binary.LittleEndian.PutUint64(timer, 17)
											ui = nil
											o.sounds = nil
											for _, tm := range teams {
												*(*uint32)(unsafe.Add(tm.C(), 44)) = 0
												tm.Lessons = int(pair[1])
											}
											membership := [4]byte{byte(local), 1, 2, 0}
											for i, id := range membership {
												u := units[i].TeamPtr()
												d := sprites[i].TeamPtr()
												*u = server.ObjectTeam{ID: server.TeamID(id)}
												*d = server.ObjectTeam{ID: server.TeamID(id)}
												if id != 0 {
													m := u
													if host == 0 {
														m = d
													}
													head := (*uint32)(unsafe.Add(teams[id-1].C(), 44))
													m.Field0 = *head
													*head = uint32(uintptr(m.C()))
												}
											}
											var scores [3][2]uint32
											for i := 0; i < 3; i++ {
												p := &o.players[i]
												p.Lessons = int32(pair[1])
												p.Field2140 = pair[1]
												p.Field3680 = 0
												if i == 2 {
													p.Field3680 = 1
												}
												scores[i] = [2]uint32{pair[1], pair[1]}
											}
											data := []byte{kind, byte(code), byte(code >> 8), reason, 0, 0, 0, 0}
											binary.LittleEndian.PutUint32(data[4:], time[0])
											input := bytes.Clone(data)
											n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(kind), data)
											accepted := on != 0 && time[0] > time[1]
											wantText, wantMode, wantTime, wantFlags := "Old", uint32(9), uint64(17), uint32(flags)
											wantTeams := [2]uint32{pair[1], pair[1]}
											var sounds [][2]int
											var wantUI []uint64
											if accepted {
												sounds = append(sounds, [2]int{309, 100})
												wantTime = 0x23456789
												if host == 0 {
													wantFlags |= 8
												}
												if graphics != 0 {
													wantUI = append(wantUI, wantTime)
												}
												id := byte(code)
												valid := id == 1 || id == 2
												own := valid && int(id) == local
												name := "(null)"
												if valid {
													name = []string{"Ruby", "Azure"}[id-1]
												}
												wantText = ""
												if reason == 1 && kind != 86 {
													wantText = "Time;"
												}
												wantMode = 0
												switch kind {
												case 86:
													if own {
														wantText += "Your team won"
													} else {
														wantText += "Team " + name + " won"
														wantMode = 1
													}
												case 87:
													if valid {
														wantText += "Flag " + name + " won"
														if !own {
															wantMode = 1
														}
													} else {
														wantText += "Flag tie"
													}
												case 89:
													if elimination != 0 {
														if !valid {
															wantText += "Elimination tie"
														} else {
															wantText += "Elimination;"
															if own {
																wantText += "Your team won"
															} else {
																wantText += "Team " + name + " won"
																wantMode = 1
															}
														}
													} else {
														if !valid {
															wantText += "Tie"
														} else if own {
															wantText += "Your team won"
														} else {
															wantText += "Lost to Team " + name
															wantMode = 1
														}
													}
												}
												if kind == 89 && valid && reason == 0 {
													if elimination == 0 {
														for i := 0; i < 2; i++ {
															if byte(i+1) == id {
																wantTeams[i] = pair[0]
															} else if wantTeams[i] >= pair[0] {
																wantTeams[i] = pair[0] - 1
															}
														}
													}
													for i := 0; i < 2; i++ {
														if membership[i] == id {
															continue
														}
														if elimination != 0 {
															if scores[i][1] < pair[0] {
																scores[i][1] = pair[0]
															}
														} else if scores[i][0] >= pair[0] {
															scores[i][0] = pair[0] - 1
														}
													}
												}
											}
											gotText := alloc.GoString16((*uint16)(unsafe.Pointer(&textbuf[0])))
											var gotScores [3][2]uint32
											for i := 0; i < 3; i++ {
												gotScores[i] = [2]uint32{uint32(o.players[i].Lessons), o.players[i].Field2140}
											}
											gotTeams := [2]uint32{uint32(teams[0].Lessons), uint32(teams[1].Lessons)}
											gotFlags := uint32(noxflags.GetGame())
											reset()
											if n != 8 || !bytes.Equal(data, input) || gotText != wantText || binary.LittleEndian.Uint32(mode) != wantMode || binary.LittleEndian.Uint64(timer) != wantTime || gotFlags != wantFlags || !slices.Equal(ui, wantUI) || !slices.Equal(o.sounds, sounds) || gotScores != scores || gotTeams != wantTeams {
												t.Fatalf("team winner kind%d on%d host%d elimination%d graphics%d reason%d code%x local%d time%v score%v text%q want%q scores%v want%v teams%v want%v", kind, on, host, elimination, graphics, reason, code, local, time, pair, gotText, wantText, gotScores, scores, gotTeams, wantTeams)
											}
											if memmap.Uint32(0x5D4594, 811908) != uint32(wantTime) {
												t.Fatal("winner clock width")
											}
											rows = append(rows, row{int(kind), on, host, elimination, graphics, int(reason), int(code), local, time[0], time[1], pair[0], pair[1], gotText, wantMode, gotFlags, wantTime, slices.Clone(ui), gotScores, gotTeams, slices.Clone(o.sounds)})
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	interactionCapture(t, "game-progress-team-winner", rows)
}
