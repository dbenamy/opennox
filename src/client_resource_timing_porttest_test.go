//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
)

func TestClientResourceFrameSamples(t *testing.T) {
	samples := serverConfigOwnBytes(t, 0x5D4594, 815220, 480)
	last := serverConfigOwnBytes(t, 0x5D4594, 815740, 8)
	cursor := serverConfigOwnBytes(t, 0x5D4594, 815756, 8)
	words, restore := legacy.PortTestClientResourceWords()
	t.Cleanup(restore)
	old := legacy.PlatformTicks
	t.Cleanup(func() { legacy.PlatformTicks = old })
	ticks := uint64(0)
	reads := 0
	legacy.PlatformTicks = func() uint64 { reads++; return ticks }
	type row struct {
		Tick, Last, Cursor, Next, Sample uint64
		Count, Stored                    uint32
	}
	var rows []row
	values := []uint64{0, 1, 33, 0xffffffff, 0x100000000, 0x100000001, 0x7fffffffffffffff, 0xffffffffffffffff}
	for _, now := range values {
		for _, before := range values {
			for _, index := range []uint64{0, 1, 59} {
				for _, high := range []uint64{0, 1, 0xffffffff} {
					for _, count := range []uint32{0, 60, 0xffffffff} {
						ticks = now
						// The existing platform C boundary returns a 32-bit unsigned tick.
						narrowed := uint64(uint32(now))
						reads = 0
						for i := range samples {
							samples[i] = 0xa5
						}
						binary.LittleEndian.PutUint64(last, before)
						initial := index | (high << 32)
						binary.LittleEndian.PutUint64(cursor, initial)
						*words["frameCount"] = count
						expected := bytes.Clone(samples)
						binary.LittleEndian.PutUint64(expected[index*8:], narrowed-before)
						legacy.Sub_43C650()
						if reads != 1 || !bytes.Equal(samples, expected) || binary.LittleEndian.Uint64(last) != narrowed || binary.LittleEndian.Uint64(cursor) != (initial+1)%60 || *words["frameCount"] != count+1 {
							t.Fatal("frame sample", now, before, initial, count)
						}
						rows = append(rows, row{now, before, initial, (initial + 1) % 60, narrowed - before, count, count + 1})
					}
				}
			}
		}
	}
	interactionCapture(t, "client-resource-frame-samples", rows)
}

func TestClientResourceFrameAverage(t *testing.T) {
	samples := serverConfigOwnBytes(t, 0x5D4594, 815220, 480)
	average := serverConfigOwnBytes(t, 0x587000, 91880, 8)
	words, restore := legacy.PortTestClientResourceWords()
	t.Cleanup(restore)
	type row struct {
		Count   int32
		Pattern int
		Average uint64
	}
	var rows []row
	for _, count := range []int32{-2147483648, -1, 0, 1, 10, 11, 59, 60, 61, 2147483647} {
		for pattern := 0; pattern < 6; pattern++ {
			var sum uint64
			n := int(count)
			if n > 60 {
				n = 60
			}
			for i := 0; i < 60; i++ {
				v := uint64(i + 1)
				switch pattern {
				case 0:
					v = 0
				case 1:
					v = 33
				case 2:
					v = uint64(i*17 + 3)
				case 3:
					v = 0xffffffffffffffff
				case 4:
					v = uint64(i)<<32 | uint64(uint32(0xffffffff)-uint32(i))
				case 5:
					if i%2 == 0 {
						v = 0x8000000000000000
					}
				}
				binary.LittleEndian.PutUint64(samples[i*8:], v)
				if i < n {
					sum += v
				}
			}
			input := bytes.Clone(samples)
			*words["frameCount"] = uint32(count)
			binary.LittleEndian.PutUint64(average, 0xabcdef9876543210)
			legacy.PortTestClientFrameAverage()
			want := uint64(33)
			if n > 10 {
				want = sum / uint64(n)
			}
			got := binary.LittleEndian.Uint64(average)
			if got != want || !bytes.Equal(samples, input) || *words["frameCount"] != uint32(count) {
				t.Fatal("frame average", count, pattern, got, want)
			}
			rows = append(rows, row{count, pattern, got})
		}
	}
	interactionCapture(t, "client-resource-frame-average", rows)
}

func TestClientResourceScalarState(t *testing.T) {
	shell := serverConfigOwnBytes(t, 0x5D4594, 815092, 4)
	drawing := serverConfigOwnBytes(t, 0x5D4594, 816408, 4)
	bindings := serverConfigOwnBytes(t, 0x5D4594, 1193128, 4)
	for _, v := range []uint32{0, 1, 2, 255, 256, 0x7fffffff, 0x80000000, 0xffffffff} {
		binary.LittleEndian.PutUint32(shell, v)
		if legacy.PortTestClientShellState() != int32(v) || binary.LittleEndian.Uint32(shell) != v {
			t.Fatal("shell state", v)
		}
		legacy.Sub_43E8C0(int(int32(v)))
		if binary.LittleEndian.Uint32(drawing) != v {
			t.Fatal("drawing state", v)
		}
	}
	for v := 0; v < 256; v++ {
		binary.LittleEndian.PutUint32(bindings, 0xabcdef00|uint32(v))
		if got := legacy.PortTestClientBindingCount(); got != byte(v) || binary.LittleEndian.Uint32(bindings) != 0xabcdef00|uint32(v) {
			t.Fatal("binding count", v, got)
		}
	}
}
