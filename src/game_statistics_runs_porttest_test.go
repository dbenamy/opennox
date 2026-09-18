//go:build porttest

package opennox

import (
	"bytes"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func statisticsExpandRuns(t *testing.T, b []byte) []byte {
	t.Helper()
	if len(b) == 0 {
		t.Fatal("missing escape marker")
	}
	marker := b[0]
	var out []byte
	for i := 1; i < len(b); i++ {
		if b[i] != marker {
			out = append(out, b[i])
			continue
		}
		i++
		if i == len(b) {
			t.Fatal("missing escape payload")
		}
		if b[i] == marker {
			out = append(out, marker)
			continue
		}
		n := int(b[i])
		i++
		if i == len(b) || n < 4 {
			t.Fatal("invalid run")
		}
		out = append(out, bytes.Repeat([]byte{b[i]}, n)...)
	}
	return out
}
func TestGameStatisticsRunEncoding(t *testing.T) {
	type row struct {
		Input, Output []byte
		Result        int32
	}
	var rows []row
	var inputs [][]byte
	for _, v := range []byte{0, 1, 127, 128, 255} {
		for _, n := range []int{0, 1, 2, 3, 4, 5, 127, 128, 254, 255, 256, 257, 511, 512} {
			inputs = append(inputs, bytes.Repeat([]byte{v}, n))
		}
	}
	for _, n := range []int{1, 2, 3, 4, 16, 255, 256, 257, 1024} {
		b := make([]byte, n)
		for i := range b {
			b[i] = byte(i)
		}
		inputs = append(inputs, b)
		b = make([]byte, n)
		for i := range b {
			b[i] = byte((i*71 + 13) % 7)
		}
		inputs = append(inputs, b)
	}
	for n := 0; n < 256; n++ {
		b := make([]byte, 256+n)
		for i := range b {
			b[i] = byte(i)
		}
		inputs = append(inputs, b)
	}
	for _, input := range inputs {
		got, ret := legacy.PortTestStatisticsRuns(input)
		if int(ret) != len(got) {
			t.Fatalf("return %d != size %d", ret, len(got))
		}
		maxRun, run := 0, 0
		var last byte
		for i, v := range input {
			if i > 0 && v == last {
				run++
			} else {
				run = 1
				last = v
			}
			if run > maxRun {
				maxRun = run
			}
		}
		if maxRun < 256 {
			if !bytes.Equal(statisticsExpandRuns(t, got), input) {
				t.Fatalf("run roundtrip input %x output %x", input, got)
			}
		} else {
			// The byte run counter wraps in C; preserve this observed legacy limitation.
			want := input[:len(input)%256]
			if input[0] >= 128 {
				want = input
			}
			if !bytes.Equal(statisticsExpandRuns(t, got), want) {
				t.Fatalf("long-run byte counter: input length %d byte %d, encoded %x, decoded length %d expected %d", len(input), input[0], got, len(statisticsExpandRuns(t, got)), len(want))
			}
		}
		freq := [256]int{}
		for _, v := range input {
			freq[v]++
		}
		least := 0
		for i := 1; i < 256; i++ {
			if freq[i] < freq[least] {
				least = i
			}
		}
		if int(got[0]) != least {
			t.Fatalf("escape selection %d != %d", got[0], least)
		}
		rows = append(rows, row{input, got, ret})
	}
	spellbookCapture(t, "game-statistics-run-encoding", rows, "24e68020e3d553eb34b28c8b62b0d7aca7f97aab1afa0445ed09f1d7a5801c6d")
}
