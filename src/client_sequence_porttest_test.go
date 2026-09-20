//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"testing"
	"unsafe"

	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/ntype"
	"github.com/opennox/opennox/v1/legacy"
)

func TestClientSequenceDelivery(t *testing.T) {
	reset, snapshot, restore := legacy.PortTestClientSequenceOwner()
	t.Cleanup(restore)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	binary.LittleEndian.PutUint32(connected, 1)
	oldTicks, oldRecv := legacy.PlatformTicks, legacy.Nox_xxx_netOnPacketRecvCli_48EA70
	t.Cleanup(func() { legacy.PlatformTicks = oldTicks; legacy.Nox_xxx_netOnPacketRecvCli_48EA70 = oldRecv })
	type delivery struct {
		Index int
		Data  []byte
	}
	type row struct {
		Name      string
		Step      int
		Tick      uint64
		Calls     int
		Delivered []delivery
		State     legacy.PortTestClientSequenceState
	}
	var rows []row
	var delivered []delivery
	var tick uint64
	calls := 0
	legacy.PlatformTicks = func() uint64 { calls++; return tick }
	legacy.Nox_xxx_netOnPacketRecvCli_48EA70 = func(ind ntype.PlayerInd, p *byte, size int) int {
		delivered = append(delivered, delivery{int(ind), bytes.Clone(unsafe.Slice(p, size))})
		return -1 // Queue consumes each record independently of the callback result.
	}
	enqueue := func(seq uint16) {
		b := binary.LittleEndian.AppendUint16([]byte{204}, seq)
		b = append(b, 2, byte(seq), byte(seq>>8))
		if n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(204), b); n != len(b) {
			t.Fatalf("enqueue %d: %d", seq, n)
		}
	}
	cases := []struct {
		Name    string
		Current uint16
		Insert  []uint16
		Ticks   []uint64
		Want    [][]uint16
	}{
		{"empty", 0, nil, []uint64{0, 30001}, [][]uint16{nil, nil}},
		{"sorted", 0, []uint16{2, 0, 1}, []uint64{1, 2}, [][]uint16{{0, 1, 2}, nil}},
		{"wrap", 65534, []uint16{0, 65535, 65534, 1}, []uint64{1}, [][]uint16{{65534, 65535, 0, 1}}},
		{"gap_expired", 0, []uint16{3, 0, 4}, []uint64{0, 30001}, [][]uint16{{0}, {3, 4}}},
		{"gap_exact", 0, []uint16{3, 0, 4}, []uint64{0, 30000, 30001}, [][]uint16{{0}, nil, nil}},
		{"gap_early", 0, []uint16{3, 0, 4}, []uint64{0, 29999, 30001}, [][]uint16{{0}, nil, nil}},
		{"future_only", 0, []uint16{3, 4}, []uint64{30001, 60002}, [][]uint16{nil, nil}},
		{"multiple_gaps", 0, []uint16{0, 2, 4}, []uint64{0, 30001, 30002}, [][]uint16{{0}, {2}, {4}}},
		{"tick_narrowing", 0, []uint16{0, 3}, []uint64{0, 0x100000001, 0x100007531}, [][]uint16{{0}, nil, nil}},
	}
	for _, tc := range cases {
		reset()
		*memmap.PtrUint16(0x5D4594, 1197360) = tc.Current
		tick = 0
		for _, seq := range tc.Insert {
			enqueue(seq)
		}
		for step, now := range tc.Ticks {
			tick = now
			calls = 0
			delivered = nil
			legacy.Sub_48D660()
			var got []uint16
			for _, d := range delivered {
				if d.Index != 31 || len(d.Data) != 2 {
					t.Fatal("callback arguments", tc.Name, d)
				}
				got = append(got, binary.LittleEndian.Uint16(d.Data))
			}
			if !reflect.DeepEqual(got, tc.Want[step]) {
				t.Fatalf("%s step%d: got%v want%v", tc.Name, step, got, tc.Want[step])
			}
			rows = append(rows, row{tc.Name, step, tick, calls, delivered, snapshot()})
		}
	}
	interactionCapture(t, "client-sequence-delivery", rows)
}

func TestClientSequenceDeliveryLengths(t *testing.T) {
	reset, snapshot, restore := legacy.PortTestClientSequenceOwner()
	t.Cleanup(restore)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	binary.LittleEndian.PutUint32(connected, 1)
	oldTicks, oldRecv := legacy.PlatformTicks, legacy.Nox_xxx_netOnPacketRecvCli_48EA70
	t.Cleanup(func() { legacy.PlatformTicks = oldTicks; legacy.Nox_xxx_netOnPacketRecvCli_48EA70 = oldRecv })
	tick := uint64(0xffffffff)
	legacy.PlatformTicks = func() uint64 { return tick }
	type row struct {
		Length, Return int
		Delivered      []byte
		State          legacy.PortTestClientSequenceState
	}
	var rows []row
	for size := 0; size <= 255; size++ {
		for _, ret := range []int{-1, 0, 1, 0x7fffffff} {
			reset()
			data := []byte{204, 0, 0, byte(size)}
			for i := 0; i < size; i++ {
				data = append(data, byte(17*i+9))
			}
			n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(204), data)
			if n != len(data) {
				t.Fatal("insertion length", n, size)
			}
			var got []byte
			calls := 0
			legacy.Nox_xxx_netOnPacketRecvCli_48EA70 = func(ind ntype.PlayerInd, p *byte, n int) int {
				calls++
				if int(ind) != 31 || n != size {
					t.Fatal("delivery arguments", ind, n, size)
				}
				got = bytes.Clone(unsafe.Slice(p, n))
				return ret
			}
			legacy.Sub_48D660()
			state := snapshot()
			if calls != 1 || !bytes.Equal(got, data[4:]) || state.Current != 1 || len(state.Nodes) != 0 || state.Ready != 0 || state.Pending != 0 {
				t.Fatal("delivery length/state", size, ret, calls, state)
			}
			rows = append(rows, row{size, ret, got, state})
		}
	}
	interactionCapture(t, "client-sequence-delivery-lengths", rows)
}

func TestClientSequenceTickWrap(t *testing.T) {
	reset, snapshot, restore := legacy.PortTestClientSequenceOwner()
	t.Cleanup(restore)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	binary.LittleEndian.PutUint32(connected, 1)
	oldTicks, oldRecv := legacy.PlatformTicks, legacy.Nox_xxx_netOnPacketRecvCli_48EA70
	t.Cleanup(func() { legacy.PlatformTicks = oldTicks; legacy.Nox_xxx_netOnPacketRecvCli_48EA70 = oldRecv })
	var tick uint64
	legacy.PlatformTicks = func() uint64 { return tick }
	type row struct {
		Start, Now uint64
		State      legacy.PortTestClientSequenceState
		Delivered  []byte
	}
	var rows []row
	for _, start := range []uint64{1, 0x80000000, 0xfffffffe, 0xffffffff, 0x100000001} {
		for _, now := range []uint64{0, 1, 30000, 30001, 0xffffffff, 0x100000000} {
			reset()
			tick = start
			for _, seq := range []byte{0, 3} {
				legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(204), []byte{204, seq, 0, 1, seq})
			}
			var got []byte
			legacy.Nox_xxx_netOnPacketRecvCli_48EA70 = func(_ ntype.PlayerInd, p *byte, n int) int {
				if n != 1 {
					t.Fatal(n)
				}
				got = append(got, *p)
				return 0
			}
			legacy.Sub_48D660()
			if !bytes.Equal(got, []byte{0}) {
				t.Fatal("initial delivery", got)
			}
			tick = now
			got = nil
			legacy.Sub_48D660()
			// C promotes the truncated clock to uint64 before subtracting; clock
			// wrap/backward movement therefore becomes a large unsigned elapsed time.
			elapsed := uint64(uint32(now)) - uint64(uint32(start))
			expired := elapsed > 30000
			if expired != bytes.Equal(got, []byte{3}) || !expired && len(got) != 0 {
				t.Fatal("clock arithmetic", start, now, elapsed, got)
			}
			rows = append(rows, row{start, now, snapshot(), got})
		}
	}
	interactionCapture(t, "client-sequence-tick-wrap", rows)
}
