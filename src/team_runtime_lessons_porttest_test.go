//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"reflect"
	"testing"
)

func TestTeamRuntimeLessons(t *testing.T) {
	o := newMatchRosterOwner(t)
	tm := o.s.Teams.Create(1)
	type row struct {
		Name    string
		Lessons uint32
		Queue   legacy.PortTestReliableReportState
	}
	var rows []row
	for _, host := range []bool{false, true} {
		for _, value := range []uint32{0, 1, 0x7fffffff, 0x80000000, 0xffffffff} {
			name := fmt.Sprintf("host%t/value%x", host, value)
			t.Run(name, func(t *testing.T) {
				mode := noxflags.GameFlag(0)
				if host {
					mode = 1
				}
				defer noxflags.PortTestGameFlags(mode)()
				o.reset()
				legacy.PortTestTeamRuntimeOther("lessons", tm, nil, int(value), 0)
				state := o.state()
				if uint32(tm.Lessons) != value {
					t.Fatal("team score width")
				}
				if host {
					want := []byte{196, 8, 1, 0, 0, 0, 0, 0, 0, 0}
					binary.LittleEndian.PutUint32(want[6:], value)
					if len(state.Nodes) != 1 || !bytes.Equal(state.Nodes[0].Data, want) {
						t.Fatal("team score message")
					}
				} else if len(state.Nodes) != 0 {
					t.Fatal("non-host score message")
				}
				o.reset()
				o.s.TeamChangeLessons(tm, int(value))
				native := o.state()
				if !reflect.DeepEqual(native, state) {
					t.Fatal("existing Go score setter differs")
				}
				rows = append(rows, row{name, value, state})
			})
		}
	}
	legacy.PortTestTeamRuntimeOther("lessons", nil, nil, 1, 0)
	spellbookCapture(t, "team-runtime-lessons", rows, "18b3efc86e89b18d120954905c2cdb22d751372158decc9ffdf89cb3298c8b61")
}
