//go:build porttest

package opennox

import (
	"bytes"
	flags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/internal/netlist"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestPlayerFilesSaveRequest(t *testing.T) {
	o := newReliableReportsOwner(t)
	defer flags.PortTestGameFlags(0)()
	env := legacy.PortTestNewMeterEnvironment()
	t.Cleanup(env.Restore)
	code := env.NamedWord("nox_player_netCode_85319C")
	var rows []map[string]any
	for _, gf := range []flags.GameFlag{0, 1} {
		for _, value := range []uint32{0, 1, 65535, 65536, 0xffffffff} {
			o.reset()
			flags.ResetGame()
			flags.SetGame(gf)
			*code = value
			ret := legacy.PortTestPlayerFileCall("nox_xxx_netSavePlayer_41CE00")
			want := []byte{0xc1, byte(value), byte(value >> 8)}
			state := o.state()
			local := o.s.NetList.CopyPacketsA(31, netlist.Kind0)
			if ret != 1 || *code != value {
				t.Fatal("save request return/code", ret, *code)
			}
			if gf == 1 {
				if len(state.Nodes) != 0 || !bytes.Equal(local, want) {
					t.Fatal("host save queue", local, state)
				}
			} else {
				if len(local) != 0 || len(state.Nodes) != 1 || !bytes.Equal(state.Nodes[0].Data, want) || state.Nodes[0].To != 31 || state.Nodes[0].Priority != 1 || state.Nodes[0].Ordered != 0 {
					t.Fatal("client save queue", local, state)
				}
			}
			rows = append(rows, map[string]any{"flags": uint32(gf), "code": value, "return": ret, "queue": state, "local": local})
		}
	}
	spellbookCapture(t, "player-files-save-request", rows, "")
}
