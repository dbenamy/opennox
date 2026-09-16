//go:build porttest

package cnxz

import (
	"bytes"
	"encoding/binary"
	"os"
	"path/filepath"
	"testing"
)

func TestMapDecodeMalformed(t *testing.T) {
	var cases [][]byte
	cases = append(cases, nil, []byte{0})
	var w mapDecodeBits
	w.put(15, 4)
	w.put(31, 5)
	cases = append(cases, w.data) // Initial index275.
	w = mapDecodeBits{}
	w.symbol(t, 256)
	w.put(0, 3)
	w.put(1, 9)
	cases = append(cases, w.data) // Four bytes exceed declared3.
	w = mapDecodeBits{}
	w.symbol(t, 256)
	w.put(0, 3)
	w.put(0, 9)
	cases = append(cases, w.data) // Zero-distance reference.
	w = mapDecodeBits{}
	w.symbol(t, 272)
	w.put(0, 24)
	w.put(0, 8)
	cases = append(cases, w.data) // Unbounded table code.
	dir := t.TempDir()
	src, dst := filepath.Join(dir, "bad.nxz"), filepath.Join(dir, "keep.map")
	for i, body := range cases {
		encoded := binary.LittleEndian.AppendUint32(nil, 3)
		encoded = append(encoded, body...)
		if err := os.WriteFile(src, encoded, 0600); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(dst, []byte("keep"), 0600); err != nil {
			t.Fatal(err)
		}
		if err := DecompressFile(src, dst); err == nil {
			t.Fatalf("malformed case%d accepted", i)
		}
		got, err := os.ReadFile(dst)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(got, []byte("keep")) {
			t.Fatal("malformed decode changed destination")
		}
	}
}
func TestMapDecodeFrequencyState(t *testing.T) {
	var w mapDecodeBits
	previous := 0
	for _, n := range []int{2, 3, 3, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 5, 5, 5} {
		w.put(0, n-previous)
		w.put(1, 1)
		previous = n
	}
	d := newMapDecoder(w.data)
	d.counts[0] = -32768
	d.counts[1] = -1
	d.counts[2] = 32767
	d.counts[3] = 1
	d.counts[4] = 1
	if err := d.rebuild(); err != nil {
		t.Fatal(err)
	}
	if d.symbols[0] != 2 || d.symbols[1] != 4 || d.symbols[2] != 3 || d.symbols[272] != 1 || d.symbols[273] != 0 {
		t.Fatal("signed frequency/tie order")
	}
	if d.counts[0] != -16384 || d.counts[1] != -1 || d.counts[2] != 16383 || d.counts[3] != 0 {
		t.Fatal("arithmetic frequency halving")
	}
	w = mapDecodeBits{}
	w.symbol(t, 65)
	d = newMapDecoder(w.data)
	d.counts[65] = 32767
	var out [1]byte
	if err := d.decode(out[:]); err != nil {
		t.Fatal(err)
	}
	if out[0] != 65 || d.counts[65] != -32768 {
		t.Fatal("frequency increment must wrap at16 bits")
	}
}

func TestMapDecodeEmptyOutput(t *testing.T) {
	dir := t.TempDir()
	src, dst := filepath.Join(dir, "empty.nxz"), filepath.Join(dir, "empty.map")
	if err := os.WriteFile(src, []byte{0, 0, 0, 0}, 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(dst, []byte("old"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := DecompressFile(src, dst); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(dst)
	if err != nil || len(got) != 0 {
		t.Fatal("zero-size map must expand to an empty file", err)
	}
}
