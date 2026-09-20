//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/libs/strman"
	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestGameMessageClientSessionObjectives(t *testing.T) {
	o := newTeamUIOwner(t)
	entries := []strman.Entry{}
	for _, key := range []string{"BallHomeTT", "BallAwayTT", "BallRedTT", "BallBlueTT"} {
		entries = append(entries, strman.Entry{ID: strman.ID("guifb.c:" + key), Vals: []strman.Variant{{Str: key}}})
	}
	// Supply both CTF and ball texts because the temporary string owner replaces
	// its table; tooltip contracts must observe actual localized values.
	for _, pair := range [][2]string{{"FlagHomeTT", "home"}, {"FlagAwayTT", "away"}, {"TheirFlagCarriedTT", "their flag carried"}, {"YourFlagCarriedTT", "your flag carried"}} {
		entries = append(entries, strman.Entry{ID: strman.ID("GUI_CTF.c:" + pair[0]), Vals: []strman.Variant{{Str: pair[1]}}})
	}
	configure, restore := o.c.srv.PortTestMeterStrings(entries...)
	t.Cleanup(restore)
	configure(0)
	oldLoad := legacy.Nox_xxx_gLoadImg
	t.Cleanup(func() { legacy.Nox_xxx_gLoadImg = oldLoad })
	names := []string{"BallAtHome", "BallAway", "BallRed", "BallBlue"}
	legacy.Nox_xxx_gLoadImg = func(name string) *noxrender.Image {
		for i, n := range names {
			if name == n {
				return o.images[i+1]
			}
		}
		return oldLoad(name)
	}
	if teamUICall("ctf-construct", 0, 0) != 1 || teamUICall("ball-construct", 0, 0) != 1 {
		t.Fatal("HUD owners")
	}
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	carriers := serverConfigOwnBytes(t, 0x5D4594, 1090128, 6)
	statuses := serverConfigOwnBytes(t, 0x5D4594, 1045612, 16)
	ballState := serverConfigOwnBytes(t, 0x5D4594, 1045644, 1)
	type row struct {
		On, Kind, ID, Action, Team, Code, Selected int
		Carriers                                   []byte
		Statuses                                   []byte
		Ball                                       byte
		Text                                       string
		Image                                      int
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for _, kind := range []int{216, 217} {
			ids := []int{1, 16}
			teams := []int{0, 1, 2, 255}
			selecteds := []int{0, 1}
			actions := []int{0, 1, 2, 3, 4, 127, 128, 255}
			if kind == 217 {
				ids = []int{0}
				teams = []int{0}
				selecteds = []int{0}
				actions = make([]int, 256)
				for i := range actions {
					actions[i] = i
				}
			}
			for _, id := range ids {
				for _, team := range teams {
					for _, selected := range selecteds {
						for _, action := range actions {
							for _, code := range []int{0, 0x8001, 0xffff} {
								binary.LittleEndian.PutUint32(connected, uint32(on))
								for i := range carriers {
									carriers[i] = 0xa5
								}
								for i := range statuses {
									statuses[i] = 77
								}
								ballState[0] = 88
								w := o.window("ball")
								if kind == 216 {
									w = o.window("ctf").ChildByID(uint(8810 + id))
								}
								w.Flags &^= 32
								if selected != 0 {
									w.Flags |= 32
								}
								w.DrawData().SetTooltip(nil, "unchanged")
								w.DrawData().SetBackgroundImage(o.images[0])
								data := []byte{byte(kind), byte(action)}
								if kind == 216 {
									data = append(data, byte(id), byte(team))
								}
								data = binary.LittleEndian.AppendUint16(data, uint16(code))
								input := bytes.Clone(data)
								wantCarriers := bytes.Clone(carriers)
								wantStatuses := bytes.Clone(statuses)
								wantBall := byte(88)
								wantText := "unchanged"
								wantImage := 0
								if on != 0 {
									if kind == 216 {
										wantStatuses[id-1] = byte(action)
										switch action {
										case 0:
											wantText = "home"
										case 1:
											wantText = "their flag carried"
											if selected != 0 {
												wantText = "your flag carried"
											}
										case 2:
											wantText = "away"
										}
										if team == 1 || team == 2 {
											if action == 0 || action == 2 {
												binary.LittleEndian.PutUint16(wantCarriers[(team-1)*2:], 0)
											} else if action == 1 {
												binary.LittleEndian.PutUint16(wantCarriers[(team-1)*2:], uint16(code))
											}
										}
									} else {
										wantBall = byte(action)
										switch action {
										case 0:
											wantImage = 1
											wantText = "BallHomeTT"
										case 1:
											wantImage = 2
											wantText = "BallAwayTT"
										case 2:
											wantImage = 3
											wantText = "BallRedTT"
										case 4:
											wantImage = 4
											wantText = "BallBlueTT"
										}
										if action == 0 || action == 1 {
											binary.LittleEndian.PutUint16(wantCarriers[4:], 0)
										} else if action == 2 || action == 4 {
											binary.LittleEndian.PutUint16(wantCarriers[4:], uint16(code))
										}
									}
								}
								n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(kind), data)
								if n != len(data) || !bytes.Equal(data, input) || !bytes.Equal(carriers, wantCarriers) || !bytes.Equal(statuses, wantStatuses) || ballState[0] != wantBall || w.DrawData().Tooltip() != wantText {
									t.Fatal("objective status", on, kind, id, action, team, code, selected, w.DrawData().Tooltip(), wantText)
								}
								if w.DrawData().BgImageHnd != o.images[wantImage].C() {
									t.Fatal("objective image", kind, action, wantImage)
								}
								rows = append(rows, row{on, kind, id, action, team, code, selected, bytes.Clone(carriers), bytes.Clone(statuses), ballState[0], w.DrawData().Tooltip(), wantImage})
							}
						}
					}
				}
			}
		}
	}
	interactionCapture(t, "game-client-session-objectives", rows)
}
