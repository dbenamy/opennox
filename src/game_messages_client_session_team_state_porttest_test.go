//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"github.com/opennox/libs/noxnet/netmsg"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"reflect"
	"testing"
)

func TestGameMessageClientSessionTeamState(t *testing.T) {
	o := newMatchRosterOwner(t)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	tm := o.s.Teams.Create(1)
	type row struct {
		On, Host, Kind int
		ID, Value      uint32
		Name           string
		Score          uint32
		Return         int
		Queue          legacy.PortTestReliableReportState
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for host := 0; host < 2; host++ {
			for _, kind := range []int{4, 8} {
				for _, id := range []uint32{0, 1, 2, 0x10001, 0xffffffff} {
					for _, value := range []uint32{0, 1, 0x7fffffff, 0x80000000, 0xffffffff} {
						noxflags.ResetGame()
						noxflags.SetGame(noxflags.GameFlag(host))
						binary.LittleEndian.PutUint32(connected, uint32(on))
						data := make([]byte, 46)
						data[0], data[1] = 196, byte(kind)
						binary.LittleEndian.PutUint32(data[2:], id)
						name := "Team Ω"
						length := 46
						if kind == 4 {
							for i, v := range []uint16{'T', 'e', 'a', 'm', ' ', 0x3a9, 0} {
								binary.LittleEndian.PutUint16(data[6+i*2:], v)
							}
						} else {
							length = 10
							binary.LittleEndian.PutUint32(data[6:], value)
						}
						setup := func() { legacy.PortTestTeamRuntimeName(tm, "Before", 0x123); tm.Lessons = 37; o.reset() }
						setup()
						if on != 0 && byte(id) == 1 {
							if kind == 4 {
								legacy.PortTestTeamRuntimeMessage("rename", tm, nil, 0, name)
							} else {
								legacy.PortTestTeamRuntimeOther("lessons", tm, nil, int(value), 0)
							}
						}
						want, queue := *tm, o.state()
						setup()
						input := bytes.Clone(data)
						n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(196), data[:length])
						got := o.state()
						if n != length || !bytes.Equal(input, data) || *tm != want || !reflect.DeepEqual(got, queue) {
							t.Fatal("team rename/score", on, host, kind, id, value, n)
						}
						rows = append(rows, row{on, host, kind, id, value, tm.Name(), uint32(tm.Lessons), n, got})
					}
				}
			}
		}
	}
	interactionCapture(t, "game-client-session-team-state", rows)
}
