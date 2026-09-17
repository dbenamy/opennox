//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"testing"
	"unsafe"

	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

func TestTeamRuntimeLeave(t *testing.T) {
	o := newMatchRosterOwner(t)
	empty := unsafe.Slice((*byte)(memmap.PtrOff(0x5D4594, 527664)), 4)
	oldEmpty := bytes.Clone(empty)
	clear(empty)
	t.Cleanup(func() { copy(empty, oldEmpty) })
	type row struct {
		Name         string
		Active       bool
		Group, Count uint32
		TeamName     string
		ID           byte
		Queue        legacy.PortTestReliableReportState
	}
	var rows []row
	for _, mode := range []int{1, 0x2001, 0x8001, 0xa001, 0x8021, 0x8061} {
		for _, linked := range []bool{false, true} {
			name := fmt.Sprintf("mode%x/linked%t", mode, linked)
			t.Run(name, func(t *testing.T) {
				defer noxflags.PortTestGameFlags(noxflags.GameFlag(mode))()
				o.s.Teams.Reset()
				o.s.Teams.ActiveCnt = 0
				tm := o.s.Teams.Create(1)
				tm.SetNameAnd68("Group", 0)
				objectXferSetWord(tm.C(), 60, 77)
				member := (*server.ObjectTeam)(o.record(t, 8))
				if linked {
					legacy.Nox_xxx_createAtImpl_4191D0(1, member, 0, 10001, 0)
				} else {
					member.ID = 1
				}
				o.reset()
				legacy.PortTestTeamRuntimeOther("leave", nil, member, 10001, 0)
				state := o.state()
				wantID := byte(1)
				if linked {
					wantID = 0
				}
				if byte(member.ID) != wantID || member.Field0 != 0 {
					t.Fatal("leave membership")
				}
				wantActive := !(linked && mode&0x8000 != 0 && mode&96 == 0)
				if tm.Active() != wantActive {
					t.Fatal("empty group team lifetime")
				}
				wantGroup := uint32(77)
				wantName := "Group"
				if linked && mode&0x8000 != 0 {
					wantGroup = 0
					wantName = ""
				}
				if uint32(tm.Ind60()) != wantGroup || tm.Name() != wantName || objectXferGetWord(tm.C(), 48) != 0 {
					t.Fatal("empty group state")
				}
				leaveMessages := 0
				for _, node := range state.Nodes {
					if len(node.Data) > 1 && node.Data[0] == 196 && node.Data[1] == 2 {
						leaveMessages++
						want := []byte{196, 2, 0, 0, 0, 0}
						binary.LittleEndian.PutUint32(want[2:], 10001)
						if !bytes.Equal(node.Data, want) {
							t.Fatal("member leave message")
						}
					}
				}
				wantMessages := 0
				if linked && mode&0x2000 != 0 {
					wantMessages = 1
				}
				if leaveMessages != wantMessages {
					t.Fatal("leave message filter")
				}
				rows = append(rows, row{name, tm.Active(), uint32(tm.Ind60()), objectXferGetWord(tm.C(), 48), tm.Name(), byte(member.ID), state})
			})
		}
	}
	legacy.PortTestTeamRuntimeOther("leave", nil, nil, 0, 0)
	spellbookCapture(t, "team-runtime-leave", rows, "1971dff97076665a8580799994cf788f619e1659e25d8b6c680bf441aa87b304")
}
