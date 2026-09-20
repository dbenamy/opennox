//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"github.com/opennox/libs/noxnet/netmsg"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"reflect"
	"testing"
	"unsafe"
)

func TestGameMessageClientSessionTeamMembers(t *testing.T) {
	o := newMatchRosterOwner(t)
	printer, messageIndex, messageRegion := teamRuntimeJoinTextOwner(t, o.s)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	type state struct {
		Centered     []byte
		MessageIndex uint32
		Messages     []string
		Teams        [][]byte
		Members      [][7]uint32
		Queue        legacy.PortTestReliableReportState
		Count        int
		Play         uint32
	}
	identities := map[uint32]uint32{0: 0}
	for i := range o.units {
		identities[uint32(uintptr(o.units[i].TeamPtr().C()))] = uint32(i + 1)
	}
	norm := func(p uint32) uint32 {
		v, ok := identities[p]
		if !ok {
			t.Fatal("unknown team link", p)
		}
		return v
	}
	snapshot := func() state {
		r := state{Centered: bytes.Clone(messageRegion), MessageIndex: *messageIndex, Messages: append([]string(nil), printer.lines...), Queue: o.state(), Count: o.s.Teams.Count(), Play: uint32(noxflags.GetGamePlay())}
		for _, tm := range o.s.Teams.Teams() {
			b := bytes.Clone(unsafe.Slice((*byte)(tm.C()), 80))
			binary.LittleEndian.PutUint32(b[44:], norm(binary.LittleEndian.Uint32(b[44:])))
			r.Teams = append(r.Teams, b)
		}
		for i := range o.units {
			u := &o.units[i]
			r.Members = append(r.Members, [7]uint32{uint32(u.TeamVal.ID), norm(u.TeamVal.Field0), u.Field35, u.Field36, u.Field37, u.Field38, u.UpdateDataPlayer().Player.Field4792})
		}
		return r
	}
	type row struct {
		On, Kind, Missing int
		Mode, ID          uint32
		Return            int
		State             state
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for _, mode := range []uint32{1, 8193} {
			for _, kind := range []int{2, 3, 5, 6, 7, 9} {
				for missing := 0; missing < 2; missing++ {
					for _, id := range []uint32{1, 2, 0x10002, 255} {
						setup := func() {
							noxflags.ResetGame()
							noxflags.SetGame(noxflags.GameFlag(mode))
							noxflags.SetGamePlay(4)
							o.s.Teams.Reset()
							o.s.Teams.ActiveCnt = 0
							o.s.Teams.Create(1)
							o.s.Teams.Create(2)
							for i := range o.units {
								u := &o.units[i]
								u.TeamVal = server.ObjectTeam{}
								u.UpdateDataPlayer().Player.NetCodeVal = u.NetCode
								legacy.Nox_xxx_createAtImpl_4191D0(1, u.TeamPtr(), 0, int(u.NetCode), 0)
								u.Field35, u.Field36, u.Field37, u.Field38 = 11, 22, 0xffffffff, 0
							}
							binary.LittleEndian.PutUint32(connected, uint32(on))
							o.reset()
							printer.lines = nil
							clear(messageRegion)
							*messageIndex = 0
						}
						setup()
						code := uint32(o.units[0].NetCode)
						if missing != 0 {
							code = 65535
						}
						data := make([]byte, 10)
						data[0], data[1] = 196, byte(kind)
						binary.LittleEndian.PutUint32(data[2:], id)
						binary.LittleEndian.PutUint16(data[6:], uint16(code))
						length := 6
						if kind == 2 {
							binary.LittleEndian.PutUint32(data[2:], code)
						}
						if kind == 3 {
							length = 10
						}
						if kind == 7 || kind == 9 {
							length = 2
						}
						if on != 0 {
							tm := o.s.Teams.ByID(server.TeamID(id))
							switch kind {
							case 2:
								if missing == 0 {
									legacy.PortTestTeamRuntimeOther("leave", nil, o.units[0].TeamPtr(), int(code), 0)
								}
							case 3:
								if missing == 0 && tm != nil {
									legacy.PortTestTeamRuntimeMessage("change", tm, o.units[0].TeamPtr(), int(code), "")
								}
							case 5:
								legacy.PortTestTeamRuntimeOther("clear", tm, nil, 0, 0)
							case 6:
								if tm != nil {
									name := tm.Name()
									noxServer.TeamRemove(tm, false)
									legacy.Sub_456EA0(name)
								}
							case 7:
								noxServer.TeamsRemoveActive(false)
								legacy.Sub_456FA0()
							case 9:
								o.s.TeamsResetYyy()
							}
						}
						want := snapshot()
						setup()
						input := bytes.Clone(data)
						n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(196), data[:length])
						got := snapshot()
						if n != length || !bytes.Equal(input, data) || !reflect.DeepEqual(got, want) {
							t.Fatal("team membership dispatch", on, mode, kind, missing, id, n)
						}
						rows = append(rows, row{on, kind, missing, mode, id, n, got})
					}
				}
			}
		}
	}
	interactionCapture(t, "game-client-session-team-members", rows)
}
