//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"github.com/opennox/libs/noxnet/netmsg"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"strings"
	"testing"
	"unicode/utf16"
)

func TestGameMessageClientSessionTeamCreate(t *testing.T) {
	o := newMatchRosterOwner(t)
	t.Cleanup(noxflags.PortTestGameFlags(0))
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	type row struct {
		On, Existing     int
		ID, Group, Score uint32
		Color            byte
		Input, Name      string
		Count, Return    int
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for existing := 0; existing < 2; existing++ {
			for _, id := range []uint32{1, 2, 255, 0x10001} {
				for _, name := range []string{"", "Team Ω", strings.Repeat("x", 18), strings.Repeat("y", 19)} {
					for _, color := range []byte{1, 2, 4} {
						for _, value := range []uint32{0, 1, 0x80000000, 0xffffffff} {
							o.s.Teams.Reset()
							o.s.Teams.ActiveCnt = 0
							if existing != 0 {
								tm := o.s.Teams.Create(server.TeamID(id))
								legacy.PortTestTeamRuntimeName(tm, "Before", 99)
							}
							binary.LittleEndian.PutUint32(connected, uint32(on))
							text := utf16.Encode([]rune(name))
							data := make([]byte, 18+len(text)*2)
							data[0], data[1] = 196, 0
							binary.LittleEndian.PutUint32(data[2:], id)
							binary.LittleEndian.PutUint32(data[6:], value)
							binary.LittleEndian.PutUint32(data[10:], ^value)
							data[15], data[16] = byte(len(text)), color
							for i, c := range text {
								binary.LittleEndian.PutUint16(data[18+2*i:], c)
							}
							input := bytes.Clone(data)
							n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(196), data)
							tm := o.s.Teams.ByID(server.TeamID(id))
							wantName := ""
							group, score := uint32(0), uint32(0)
							if n != len(data) || !bytes.Equal(input, data) {
								t.Fatal("team create length/input")
							}
							count := existing
							if on != 0 {
								count = 1
								if len(text) > 20 {
									text = text[:20]
								}
								wantName = string(utf16.Decode(text))
								group = value
								score = ^value
								if tm == nil || tm.Name() != wantName || uint32(tm.Ind60()) != group || uint32(tm.Lessons) != score || byte(tm.ColorInd) != color {
									t.Fatal("team create fields", existing, id, name, color, value)
								}
							} else if existing != 0 {
								wantName = "Before"
								if tm == nil || tm.Name() != wantName {
									t.Fatal("disconnected existing team")
								}
							} else if tm != nil {
								t.Fatal("disconnected created team")
							}
							if o.s.Teams.Count() != count {
								t.Fatal("team count", o.s.Teams.Count(), count)
							}
							rows = append(rows, row{on, existing, id, group, score, color, name, wantName, count, n})
						}
					}
				}
			}
		}
	}
	interactionCapture(t, "game-client-session-team-create", rows)
}
