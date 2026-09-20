//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/libs/strman"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"slices"
	"testing"
	"unsafe"
)

func TestGameMessageClientPlayerWinner(t *testing.T) {
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
	formats := map[string]string{"TimeLimitReached": "Time;", "DM_Loss": "Lost to %s", "DM_MaleVictory": "Male won", "DM_FemaleVictory": "Female won", "DM_Tie": "Tie", "HL_Tie": "Elimination tie", "HL_Header": "Elimination;", "HL_You": "You", "HL_Victory": "%s won"}
	var entries []strman.Entry
	for k, v := range formats {
		entries = append(entries, strman.Entry{ID: strman.ID("cdecode.c:" + k), Vals: []strman.Variant{{Str: v}}})
	}
	set, restoreStrings := o.c.srv.Server.PortTestMeterStrings(entries...)
	t.Cleanup(restoreStrings)
	set(0)
	for i := range o.players {
		o.players[i].Active = 0
	}
	for i := 0; i < 3; i++ {
		p := &o.players[i]
		p.Active = 1
		p.PlayerInd = byte(i)
		p.NetCodeVal = uint32(7 + i)
		p.SetName([]string{"Ada", "Bob", "Observer"}[i])
	}
	type row struct {
		On, Host, Elimination, Graphics, Reason, Code, Gender int
		Frame, Gate, Limit, Score                             uint32
		Text                                                  string
		Mode, Flags                                           uint32
		Time                                                  uint64
		UI                                                    []uint64
		Scores                                                [3][2]uint32
		Sounds                                                [][2]int
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for host := 0; host < 2; host++ {
			for elimination := 0; elimination < 2; elimination++ {
				for graphics := 0; graphics < 2; graphics++ {
					for _, reason := range []byte{0, 1, 2} {
						for _, code := range []uint16{0, 7, 8, 0x8007, 0xffff} {
							for gender := 0; gender < 2; gender++ {
								for _, time := range [][2]uint32{{0, 0}, {100, 100}, {101, 100}, {0xffffffff, 0x80000000}} {
									for _, pair := range [][2]uint32{{0, 0}, {0, 0xffffffff}, {1, 0}, {1, 1}, {65535, 65534}, {65535, 0xffffffff}} {
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
										var scores [3][2]uint32
										for i := 0; i < 3; i++ {
											p := &o.players[i]
											p.Lessons = int32(pair[1])
											p.Field2140 = pair[1]
											p.Field3680 = 0
											if i == 2 {
												p.Field3680 = 1
											}
											*(*byte)(unsafe.Add(unsafe.Pointer(p), 2252)) = byte(gender)
											scores[i] = [2]uint32{pair[1], pair[1]}
										}
										data := []byte{88, byte(code), byte(code >> 8), reason, 0, 0, 0, 0}
										binary.LittleEndian.PutUint32(data[4:], time[0])
										input := bytes.Clone(data)
										n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(88), data)
										accepted := on != 0 && time[0] > time[1]
										wantText, wantMode, wantTime, wantFlags := "Old", uint32(9), uint64(17), uint32(flags)
										var sounds [][2]int
										var wantUI []uint64
										if accepted {
											sounds = append(sounds, [2]int{309, 100})
											wantTime = 0x23456789
											if host == 0 {
												wantFlags |= 8
											}
											if graphics != 0 {
												wantUI = append(wantUI, 17)
											}
											wantText = ""
											if reason == 1 {
												wantText = "Time;"
											}
											id := code & 0x7fff
											valid := id == 7 || id == 8
											wantMode = 1
											if elimination == 0 {
												if code == 0 {
													wantText += "Tie"
													wantMode = 0
												} else if valid {
													if id == 7 {
														wantMode = 0
														if gender == 0 {
															wantText += "Male won"
														} else {
															wantText += "Female won"
														}
													} else {
														wantText += "Lost to Bob"
													}
												}
											} else {
												if code == 0 {
													wantText += "Elimination tie"
													wantMode = 0
												} else {
													wantText += "Elimination;"
													if id == 7 {
														wantText += "You won"
														wantMode = 0
													} else if valid {
														wantText += "Bob won"
													}
												}
											}
											adjust := valid && reason == 0 && (elimination == 0 || id != 7)
											if adjust {
												for i := 0; i < 2; i++ {
													win := uint16(7+i) == id
													if elimination == 0 {
														if win {
															scores[i][0] = pair[0]
														} else if scores[i][0] >= pair[0] {
															scores[i][0] = pair[0] - 1
														}
													} else {
														if win {
															if scores[i][1] >= pair[0] {
																scores[i][1] = pair[0] - 1
															}
														} else if scores[i][1] < pair[0] {
															scores[i][1] = pair[0]
														}
													}
												}
											}
										}
										gotText := alloc.GoString16((*uint16)(unsafe.Pointer(&textbuf[0])))
										var gotScores [3][2]uint32
										for i := 0; i < 3; i++ {
											gotScores[i] = [2]uint32{uint32(o.players[i].Lessons), o.players[i].Field2140}
										}
										gotFlags := uint32(noxflags.GetGame())
										reset()
										if n != 8 || !bytes.Equal(data, input) || gotText != wantText || binary.LittleEndian.Uint32(mode) != wantMode || binary.LittleEndian.Uint64(timer) != wantTime || gotFlags != wantFlags || !slices.Equal(ui, wantUI) || !slices.Equal(o.sounds, sounds) || gotScores != scores {
											t.Fatalf("player winner on%d host%d elimination%d graphics%d reason%d code%x gender%d time%v score%v text%q want%q scores%v want%v", on, host, elimination, graphics, reason, code, gender, time, pair, gotText, wantText, gotScores, scores)
										}
										if memmap.Uint32(0x5D4594, 811908) != uint32(wantTime) {
											t.Fatal("winner clock width")
										}
										rows = append(rows, row{on, host, elimination, graphics, int(reason), int(code), gender, time[0], time[1], pair[0], pair[1], gotText, wantMode, gotFlags, wantTime, slices.Clone(ui), gotScores, slices.Clone(o.sounds)})
									}
								}
							}
						}
					}
				}
			}
		}
	}
	interactionCapture(t, "game-progress-player-winner", rows)
}
