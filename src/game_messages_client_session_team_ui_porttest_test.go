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
	"unicode/utf16"
)

func TestGameMessageClientSessionTeamUI(t *testing.T) {
	o := newTeamUIOwner(t)
	t.Cleanup(o.c.srv.PortTestMapDrawableTeamMessages())
	t.Cleanup(o.c.srv.PortTestGameMessageTeamTitles())
	t.Cleanup(noxflags.PortTestGameFlags(0))
	o.openPlayers(t)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	type state struct {
		Rows  []legacy.PortTestTeamUIRow
		Names []string
		Count int
	}
	snapshot := func() state {
		return state{legacy.PortTestTeamUIRows(true), teamUIRowNames(o.window("players").ChildByID(10502)), o.c.srv.Teams.Count()}
	}
	type row struct {
		On, Kind, Title, Color int
		State                  state
		Return                 int
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for _, kind := range []int{0, 4, 6, 7} {
			for title := 0; title < 2; title++ {
				for _, color := range []int{1, 2} {
					setup := func() {
						o.c.srv.Teams.Reset()
						o.c.srv.Teams.ActiveCnt = 0
						for i, name := range []string{"Before", "Other"} {
							tm := o.c.srv.Teams.Create(server.TeamID(i + 1))
							legacy.PortTestTeamRuntimeName(tm, name, 0)
							tm.ColorInd = server.TeamColor(i + 1)
						}
						teamUICall("team-clear", 0, 0)
						binary.LittleEndian.PutUint32(connected, uint32(on))
					}
					name := "New Ω"
					text := utf16.Encode([]rune(name))
					data := make([]byte, 46)
					data[0], data[1] = 196, byte(kind)
					binary.LittleEndian.PutUint32(data[2:], 1)
					length := 46
					switch kind {
					case 0:
						length = 18 + 2*len(text)
						binary.LittleEndian.PutUint32(data[6:], 0x80000000)
						binary.LittleEndian.PutUint32(data[10:], 0xffffffff)
						data[15], data[16], data[17] = byte(len(text)), byte(color), byte(title)
						for i, c := range text {
							binary.LittleEndian.PutUint16(data[18+2*i:], c)
						}
						if title != 0 {
							name = o.c.srv.Teams.TeamTitle(server.TeamColor(color))
						}
					case 4:
						for i, c := range text {
							binary.LittleEndian.PutUint16(data[6+2*i:], c)
						}
					case 6:
						length = 6
					case 7:
						length = 2
					}
					setup()
					tm := o.c.srv.Teams.ByID(1)
					if on != 0 {
						switch kind {
						case 0:
							legacy.PortTestTeamRuntimeName(tm, name, 0)
							legacy.PortTestTeamRuntimeOther("group", tm, nil, -2147483648, 0)
							legacy.PortTestTeamRuntimeOther("lessons", tm, nil, -1, 0)
							tm.ColorInd = server.TeamColor(color)
							legacy.PortTestTeamUI("team-add", nil, tm, 0, 0, "")
						case 4:
							legacy.PortTestTeamRuntimeMessage("rename", tm, nil, 0, name)
						case 6:
							o.c.srv.TeamRemove(tm, false)
							legacy.Sub_456EA0("Before")
						case 7:
							o.c.srv.TeamsRemoveActive(false)
							legacy.Sub_456FA0()
						}
					}
					want := snapshot()
					setup()
					input := bytes.Clone(data)
					n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(196), data[:length])
					got := snapshot()
					if n != length || !bytes.Equal(input, data) || !reflect.DeepEqual(got, want) {
						t.Fatal("team UI message", on, kind, title, color, n, got, want)
					}
					rows = append(rows, row{on, kind, title, color, got, n})
				}
			}
		}
	}
	interactionCapture(t, "game-client-session-team-ui", rows)
}
