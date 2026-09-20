//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"reflect"
	"testing"
)

func TestGameMessageClientSessionSequence(t *testing.T) {
	reset, snapshot, restore := legacy.PortTestClientSequenceOwner()
	t.Cleanup(restore)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	oldTicks := legacy.PlatformTicks
	t.Cleanup(func() { legacy.PlatformTicks = oldTicks })
	type row struct {
		On, Current, Sequence, Length, Return, Calls int
		State                                        legacy.PortTestClientSequenceState
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for _, current := range []uint16{0, 1, 9999, 10000, 10001, 65535} {
			seqs := []uint16{0, 1, 9999, 10000, 55534, 55535, 55536, 65535, current, current - 1, current + 1, current - 10000, current - 10001}
			for _, seq := range seqs {
				for _, size := range []int{0, 1, 127, 255} {
					reset()
					binary.LittleEndian.PutUint32(connected, uint32(on))
					*memmap.PtrUint16(0x5D4594, 1197360) = current
					calls := 0
					legacy.PlatformTicks = func() uint64 { calls++; return 0x100000002 }
					data := binary.LittleEndian.AppendUint16([]byte{204}, seq)
					data = append(data, byte(size))
					for i := 0; i < size; i++ {
						data = append(data, byte(31*i+7))
					}
					input := bytes.Clone(data)
					n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(204), data)
					got := snapshot()
					// Original helper rejects the previous 10,000 sequence values, including
					// its 65535-based wrap boundary, and leaves future/current values queued.
					age := int(current) - int(seq)
					old := age > 0 && age <= 10000
					if current < 10000 && int(seq) >= 65535-(10000-int(current)) {
						old = true
					}
					accepted := on != 0 && !old
					wantCount := 0
					if accepted {
						wantCount = 1
					}
					if n != len(data) || !bytes.Equal(data, input) || len(got.Nodes) != wantCount || calls != wantCount || got.Current != current || got.Pending != 0 {
						t.Fatal("sequence insertion/gate", on, current, seq, size, n, calls, got)
					}
					wantReady := uint32(0)
					if accepted && seq == current {
						wantReady = uint32(seq) + 1
					}
					if got.Ready != wantReady {
						t.Fatal("sequence ready cursor")
					}
					if accepted {
						node := got.Nodes[0]
						if node.Sequence != uint32(seq) || node.Ticks != 2 || !bytes.Equal(node.Data, data[4:]) {
							t.Fatal("sequence copied payload/timestamp")
						}
					}
					rows = append(rows, row{on, int(current), int(seq), size, n, calls, got})
					for i := 4; i < len(data); i++ {
						data[i] ^= 255
					}
					input = bytes.Clone(data)
					n = legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(204), data)
					if n != len(data) || !bytes.Equal(data, input) || !reflect.DeepEqual(snapshot(), got) || calls != wantCount {
						t.Fatal("duplicate sequence preserves first payload and time", on, current, seq, size)
					}
				}
			}
		}
	}
	interactionCapture(t, "game-client-session-sequence", rows)
}
