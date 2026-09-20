//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
	"unsafe"
)

func TestGameMessageClientSessionVote(t *testing.T) {
	o := newEntryOwner(t)
	words, restore := legacy.PortTestClientInteractionWords()
	t.Cleanup(restore)
	voteWords, restoreVote := legacy.PortTestVoteGUIOwner()
	t.Cleanup(restoreVote)
	icon := o.c.GUI.NewWindowRaw(nil, 8, 10, 20, 30, 40, nil)
	t.Cleanup(icon.Destroy)
	*words["dword_5d4594_1321216"] = uint32(uintptr(unsafe.Pointer(icon)))
	mode := serverConfigOwnBytes(t, 0x5D4594, 1197304, 4)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	type row struct {
		On, Kind, Value, Return int
		Mode, Choice, Previous  uint32
		Hidden                  bool
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for kind := 0; kind < 256; kind++ {
			values := []int{0}
			if kind == 6 {
				values = make([]int, 256)
				for i := range values {
					values[i] = i
				}
			}
			for _, value := range values {
				binary.LittleEndian.PutUint32(connected, uint32(on))
				binary.LittleEndian.PutUint32(mode, 0xcafe1234)
				*voteWords["choice"], *voteWords["previousChoice"] = 17, 29
				icon.Show()
				data := []byte{238, byte(kind), byte(value)}
				input := bytes.Clone(data)
				n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(238), data)
				wantMode, wantChoice, wantPrevious := uint32(0xcafe1234), uint32(17), uint32(29)
				wantN := -1
				hidden := false
				if kind == 6 {
					wantN = 3
					wantMode = uint32(value)
					hidden = value != 1
				}
				if kind == 7 {
					wantN = 2
					wantChoice = 0
					wantPrevious = 0
				}
				if n != wantN || !bytes.Equal(input, data) || binary.LittleEndian.Uint32(mode) != wantMode || *voteWords["choice"] != wantChoice || *voteWords["previousChoice"] != wantPrevious || icon.GetFlags().IsHidden() != hidden {
					t.Fatal("vote submessage/gate", on, kind, value)
				}
				rows = append(rows, row{on, kind, value, n, wantMode, wantChoice, wantPrevious, hidden})
			}
		}
	}
	interactionCapture(t, "game-client-session-vote", rows)
}
