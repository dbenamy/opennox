//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestGameStatisticsEncoding(t *testing.T) {
	words, state := statisticsRandomOwner(t)
	statisticsOwnReportTags(t)
	setClock, restore := legacy.PortTestStatisticsClock()
	defer restore()
	type row struct {
		Outer  bool
		Size   int
		Clock  uint32
		Bytes  []byte
		Result int32
		State  []byte
		A, B   uint32
	}
	var rows []row
	for _, outer := range []bool{false, true} {
		for _, size := range []int{1, 2, 3, 4, 13, 14, 15, 16, 17, 239, 240, 241, 242, 255, 256, 257, 512, 1024} {
			for _, clock := range []uint32{0, 1, 1700000000, 0x7fffffff, 0x80000000, 0xffffffff} {
				clear(state)
				*words["random-a"], *words["random-b"] = 0, 0
				setClock(clock)
				data := make([]byte, size)
				for i := range data {
					data[i] = byte(i*71 + 13)
				}
				source := data
				if outer {
					source, _ = legacy.PortTestStatisticsRuns(data)
				}
				seedCheck := int32(clock)
				if seedCheck > 0 {
					seedCheck = -seedCheck
				}
				first, _ := legacy.PortTestStatisticsRandom(seedCheck, 1)
				if first[0] < 0 || first[0] >= 1 {
					t.Fatal("selected clock seed has invalid insertion fraction", clock, first)
				}
				clear(state)
				*words["random-a"], *words["random-b"] = 0, 0
				got, result := legacy.PortTestStatisticsEncoding(data, outer)
				finalState := bytes.Clone(state)
				a, b := *words["random-a"], *words["random-b"]
				if len(source) < 15 {
					if got != nil || result != -2 {
						t.Fatal("short report encoding", outer, size, result)
					}
				} else {
					if int(result) != len(got) {
						t.Fatal("encoding length/clock calls", outer, size, result)
					}
					frame := got
					if outer {
						fields := statisticsDecode(t, got)
						if len(fields) != 1 || fields[0].Tag != "CNTL" || fields[0].Kind != 20 {
							t.Fatal("outer report record", fields)
						}
						frame = fields[0].Data
					}
					if len(frame) != len(source)+5 {
						t.Fatal("encoded frame size")
					}
					insert := int(frame[5])
					if insert < 10 || insert+4 > len(frame) {
						t.Fatal("embedded seed position")
					}
					seed := int32(clock)
					if seed > 0 {
						seed = -seed
					}
					if int32(binary.BigEndian.Uint32(frame[insert:])) != seed {
						t.Fatal("embedded clock seed")
					}
					payload := append([]byte{}, frame[:5]...)
					payload = append(payload, frame[6:insert]...)
					payload = append(payload, frame[insert+4:]...)
					clear(state)
					*words["random-a"], *words["random-b"] = 0, 0
					mask := legacy.PortTestStatisticsRandomBytes(seed, len(source))
					for i := range payload {
						payload[i] ^= mask[i]
					}
					if !bytes.Equal(payload, source) {
						t.Fatal("decoded local report differs")
					}
				}
				rows = append(rows, row{outer, size, clock, got, result, finalState, a, b})
			}
		}
	}
	spellbookCapture(t, "game-statistics-encoding", rows, "fe7fcd3e74257057ea9c51c1092bd53a1f6ac19d13b73d9b38c26b365aba40bc")
}
