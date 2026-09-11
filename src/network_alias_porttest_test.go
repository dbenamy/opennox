//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
)

func TestNetworkAliasSelection(t *testing.T) {
	put := func(data []byte, slot int, a, b uint16, frame uint32) {
		binary.LittleEndian.PutUint16(data[slot*8:], a)
		binary.LittleEndian.PutUint16(data[slot*8+2:], b)
		binary.LittleEndian.PutUint32(data[slot*8+4:], frame)
	}
	for start := uint32(0); start < 256; start++ {
		for _, keys := range [][2]uint32{{start, 123}, {start + 65536, 123}, {start, 0xffffffff}, {0xffffff00 | start, 123}} {
			for _, frame := range []uint32{0, 1, 0x80000000, 0xffffffff} {
				for mode := 0; mode < 6; mode++ {
					data := make([]byte, 255*8)
					for slot := 0; slot < 255; slot++ {
						put(data, slot, 0xaaaa, 0xbbbb, frame)
					}
					first := int(byte(keys[0]))
					if first == 0 || first == 255 {
						first = 1
					}
					switch mode {
					case 1:
						put(data, first, uint16(keys[0]), uint16(keys[1]), frame)
					case 2:
						put(data, 254, uint16(keys[0]), uint16(keys[1]), frame)
						put(data, 1, uint16(keys[0]), uint16(keys[1]), frame)
					case 3:
						put(data, (first+126)%254+1, 0xaaaa, 0xbbbb, frame-1)
					case 4:
						put(data, first, 0xaaaa, 0xbbbb, frame+1)
					case 5:
						put(data, 0, uint16(keys[0]), uint16(keys[1]), 0)
					}
					// Rank every eligible slot by circular distance, independent of the C loop.
					chosen, distance := 255, 255
					for slot := 1; slot < 255; slot++ {
						rec := data[slot*8:]
						match := int32(binary.LittleEndian.Uint16(rec)) == int32(keys[0]) && int32(binary.LittleEndian.Uint16(rec[2:])) == int32(keys[1])
						expired := binary.LittleEndian.Uint32(rec[4:]) < frame
						rank := (slot - first + 254) % 254
						if (match || expired) && rank < distance {
							chosen, distance = slot, rank
						}
					}
					got := legacy.PortTestAliasTable(data, []legacy.PortTestAliasCall{{Kind: "select", Key1: keys[0], Key2: keys[1], Frame: frame}})[0]
					if got.Return != int(int8(byte(chosen))) || !got.GuardsOK || !bytes.Equal(data, got.Data) {
						t.Fatalf("keys=%x frame=%x mode=%d got=%d want=%d guards=%v", keys, frame, mode, got.Return, int(int8(byte(chosen))), got.GuardsOK)
					}
				}
			}
		}
	}
}

func TestNetworkAliasWriteReset(t *testing.T) {
	values := []uint32{0, 1, 0x7fff, 0x8000, 0xffff, 0x10000, 0x80000000, 0xffffffff}
	for slot := 0; slot < 255; slot++ {
		for i, a := range values {
			b, frame := values[(i+3)%len(values)], values[(i+5)%len(values)]
			data := bytes.Repeat([]byte{0x96}, 255*8)
			want := append([]byte(nil), data...)
			binary.LittleEndian.PutUint16(want[slot*8:], uint16(a))
			binary.LittleEndian.PutUint16(want[slot*8+2:], uint16(b))
			binary.LittleEndian.PutUint32(want[slot*8+4:], frame)
			got := legacy.PortTestAliasTable(data, []legacy.PortTestAliasCall{{Kind: "write", Slot: slot, Key1: a, Key2: b, Frame: frame}, {Kind: "reset", Wrapper: i%2 != 0}})
			if !got[0].PointerSame || !got[0].GuardsOK || !bytes.Equal(got[0].Data, want) {
				t.Fatalf("write slot=%d values=%x/%x/%x failed", slot, a, b, frame)
			}
			if got[1].Return != 0 || !got[1].GuardsOK || !bytes.Equal(got[1].Data, make([]byte, 255*8)) {
				t.Fatalf("reset slot=%d failed", slot)
			}
		}
	}
}
