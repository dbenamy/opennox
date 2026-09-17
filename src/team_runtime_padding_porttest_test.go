//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"testing"
	"unicode/utf16"
	"unsafe"

	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/internal/netlist"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

// Fixed-width team messages must initialize every byte, including the unused
// tail after a short name and the reserved final word of a team-change request.
func TestTeamRuntimeMessagePadding(t *testing.T) {
	o := newMatchRosterOwner(t)
	for _, op := range []string{"rename", "change", "join-request", "change-request"} {
		for _, name := range []string{"abcdefghijklmnopqrs", "A", "", "Team Ω"} {
			t.Run(fmt.Sprintf("%s/%q", op, name), func(t *testing.T) {
				restore := noxflags.PortTestGameFlags(1)
				defer restore()
				o.s.Teams.Reset()
				o.s.Teams.ActiveCnt = 0
				a := o.s.Teams.Create(1)
				b := o.s.Teams.Create(2)
				// An ordinary non-player member exercises the real list without
				// requiring the unrelated player-join text/GUI owner.
				member := (*server.ObjectTeam)(o.record(t, int(unsafe.Sizeof(server.ObjectTeam{}))))
				code := 10001
				legacy.Nox_xxx_createAtImpl_4191D0(a.ID(), member, 0, code, 0)
				o.reset()
				if op == "change" {
					defer noxflags.PortTestGameFlags(0x2001)()
				}
				legacy.PortTestTeamRuntimeMessage(op, b, member, code, name)
				var got []byte
				if op == "join-request" || op == "change-request" {
					got = o.s.NetList.CopyPacketsA(31, netlist.Kind0)
				} else {
					state := o.state()
					if len(state.Nodes) != 1 {
						t.Fatalf("message count: %d", len(state.Nodes))
					}
					got = state.Nodes[0].Data
				}
				want := make([]byte, 10)
				if op == "rename" {
					want = make([]byte, 46)
					want[1] = 4
					for i, c := range utf16.Encode([]rune(name)) {
						binary.LittleEndian.PutUint16(want[6+2*i:], c)
					}
					if b.Name() != name {
						t.Fatal("renamed team differs")
					}
				} else {
					want[1] = map[string]byte{"change": 3, "join-request": 10, "change-request": 11}[op]
					binary.LittleEndian.PutUint16(want[6:], uint16(code))
				}
				want[0] = 196
				binary.LittleEndian.PutUint32(want[2:], uint32(b.ID()))
				if !bytes.Equal(got, want) {
					t.Fatal("fixed-width team message has unspecified bytes")
				}
			})
		}
	}
}
