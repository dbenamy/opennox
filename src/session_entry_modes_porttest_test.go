//go:build porttest

package opennox

import (
	"encoding/binary"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestSessionEntryModeIndex(t *testing.T) {
	raw := serverConfigOwnBytes(t, 0x587000, 4704, 24)
	type extra struct{ Input, Result int32 }
	type row struct {
		Table   [6]uint32
		Results []byte
		Extra   []extra
	}
	var rows []row
	for _, table := range [][6]uint32{{256, 1024, 32, 16, 64, 128}, {16, 32, 16, 64, 128, 256}, {17, 33, 65, 129, 257, 4096}} {
		first := map[uint32]int32{}
		for i, v := range table {
			binary.LittleEndian.PutUint32(raw[4*i:], v)
			if _, ok := first[v]; !ok {
				first[v] = int32(i)
			}
		}
		r := row{Table: table, Results: make([]byte, 65536)}
		for value := uint32(0); value < 65536; value++ {
			got := legacy.PortTestSessionEntryScalar("mode-index", int32(value))
			want := first[value&0x17f0]
			if got != want {
				t.Fatalf("mode%x table%v result%d want%d", value, table, got, want)
			}
			r.Results[value] = byte(got)
		}
		for _, value := range []int32{-2147483648, -1, 65536, 65552, 2147483647} {
			got := legacy.PortTestSessionEntryScalar("mode-index", value)
			if got != first[uint32(value)&0x17f0] {
				t.Fatal("short narrowing", value, got)
			}
			r.Extra = append(r.Extra, extra{value, got})
		}
		rows = append(rows, r)
	}
	spellbookCapture(t, "session-entry-mode-index", rows, "910112cb9c2efc60e78f2dace2836fef1e3f8dfd6323c331fe73fc53b894d7bc")
}
