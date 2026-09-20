//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"github.com/opennox/libs/noxnet/netmsg"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"image"
	"testing"
)

func TestGameMessageClientSessionTeamLocalRequest(t *testing.T) {
	o := newTeamUIOwner(t)
	t.Cleanup(o.c.srv.PortTestMapDrawableTeamMessages())
	t.Cleanup(noxflags.PortTestGameFlags(0))
	reset, packets, free := legacy.PortTestVisibilityReliableOwner()
	t.Cleanup(free)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	oldCode := legacy.ClientPlayerNetCode()
	t.Cleanup(func() { legacy.ClientSetPlayerNetCode(oldCode) })
	dr := o.c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(10, 20))
	if dr == nil {
		t.Fatal("local drawable")
	}
	dr.ObjClass = 0
	type row struct {
		On, Present, Same, Flag, Local int
		Packets                        []legacy.PortTestShopPacketResult
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for present := 0; present < 2; present++ {
			for same := 0; same < 2; same++ {
				for _, flag := range []byte{0, 1, 2, 3, 255} {
					for _, local := range []int{17, 0x8011, 0x10011} {
						reset()
						o.c.srv.Teams.Reset()
						o.c.srv.Teams.ActiveCnt = 0
						o.c.srv.Teams.Create(1)
						binary.LittleEndian.PutUint32(connected, uint32(on))
						legacy.ClientSetPlayerNetCode(local)
						dr.NetCode32 = uint32(local)
						if present == 0 {
							dr.NetCode32 = 999
						}
						*dr.TeamPtr() = server.ObjectTeam{}
						if same != 0 {
							dr.TeamPtr().ID = 1
						}
						data := make([]byte, 18)
						data[0], data[1], data[2], data[14], data[16] = 196, 0, 1, flag, 1
						input := bytes.Clone(data)
						n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(196), data)
						got := packets()
						sent := on != 0 && present != 0 && same == 0 && flag&1 != 0
						if n != 18 || !bytes.Equal(input, data) || (len(got) != 0) != sent || len(got) > 1 || dr.TeamPtr().ID != server.TeamID(same) {
							t.Fatal("local team request", on, present, same, flag, local, n, got)
						}
						if sent {
							want := []byte{196, 10, 1, 0, 0, 0, byte(local), byte(local >> 8), 0, 0}
							if got[0].Recipient != 31 || !bytes.Equal(got[0].Data, want) {
								t.Fatal("local team request bytes", got[0])
							}
						}
						rows = append(rows, row{on, present, same, int(flag), local, got})
					}
				}
			}
		}
	}
	interactionCapture(t, "game-client-session-team-local-request", rows)
}
func TestGameMessageClientSessionTeamLocalHost(t *testing.T) {
	o := newMatchRosterOwner(t)
	t.Cleanup(noxflags.PortTestGameFlags(1))
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	oldCode := legacy.ClientPlayerNetCode()
	t.Cleanup(func() { legacy.ClientSetPlayerNetCode(oldCode) })
	type row struct {
		On, Present, Same, Flag int
		Member                  byte
		Count                   uint32
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for present := 0; present < 2; present++ {
			for same := 0; same < 2; same++ {
				for _, flag := range []byte{0, 1, 2, 3, 255} {
					o.s.Teams.Reset()
					o.s.Teams.ActiveCnt = 0
					tm := o.s.Teams.Create(1)
					for i := range o.units {
						o.units[i].TeamVal = server.ObjectTeam{}
					}
					u := &o.units[0]
					code := int(u.NetCode)
					legacy.ClientSetPlayerNetCode(9999)
					if same != 0 {
						legacy.Nox_xxx_createAtImpl_4191D0(1, u.TeamPtr(), 0, code, 0)
					}
					if present == 0 {
						code = 65535
					}
					legacy.ClientSetPlayerNetCode(code)
					binary.LittleEndian.PutUint32(connected, uint32(on))
					o.reset()
					data := make([]byte, 18)
					data[0], data[1], data[2], data[14], data[16] = 196, 0, 1, flag, 1
					input := bytes.Clone(data)
					n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(196), data)
					want := same
					if on != 0 && present != 0 && flag&1 != 0 {
						want = 1
					}
					if n != 18 || !bytes.Equal(input, data) || int(u.TeamVal.ID) != want || objectXferGetWord(tm.C(), 48) != uint32(want) {
						t.Fatal("host team assignment", on, present, same, flag, n, u.TeamVal, objectXferGetWord(tm.C(), 48))
					}
					rows = append(rows, row{on, present, same, int(flag), byte(u.TeamVal.ID), objectXferGetWord(tm.C(), 48)})
				}
			}
		}
	}
	interactionCapture(t, "game-client-session-team-local-host", rows)
}
