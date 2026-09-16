//go:build porttest

package cnxz

import (
	"bytes"
	"encoding/binary"
	"os"
	"path/filepath"
	"testing"
)

func TestMapEncodeEmptyAndSize(t *testing.T) {
	dir := t.TempDir()
	src, dst := filepath.Join(dir, "input"), filepath.Join(dir, "output")
	if err := os.WriteFile(src, nil, 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(dst, bytes.Repeat([]byte{123}, 100), 0600); err != nil {
		t.Fatal(err)
	}
	if err := CompressFile(src, dst); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(dst)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, []byte{0, 0, 0, 0}) {
		t.Fatalf("empty file = %x", got)
	}
	if err := DecompressFile(dst, src); err != nil {
		t.Fatal(err)
	}
	if got, err := os.ReadFile(src); err != nil || len(got) != 0 {
		t.Fatalf("empty roundtrip: %d, %v", len(got), err)
	}
	f, err := os.OpenFile(src, os.O_WRONLY, 0600)
	if err != nil {
		t.Fatal(err)
	}
	// Sparse metadata only: reject before allocating or reading the body.
	err = f.Truncate(int64(^uint(0) >> 1))
	closeErr := f.Close()
	if err != nil {
		t.Fatal(err)
	}
	if closeErr != nil {
		t.Fatal(closeErr)
	}
	if err := CompressFile(src, dst); err == nil {
		t.Fatal("oversized file accepted")
	}
	got, err = os.ReadFile(dst)
	if err != nil || !bytes.Equal(got, []byte{0, 0, 0, 0}) {
		t.Fatal("oversized input changed destination")
	}
}

func TestMapEncodeAllMatchCodes(t *testing.T) {
	e := newMapEncoder()
	want := make([]byte, 65535)
	for i := range want {
		want[i] = byte(i*37 + i/251)
	}
	e.literals(want)
	// Exercise every length and every distance prefix boundary, including overlapping
	// copies, across many adaptive table rebuilds. The decoder independently expands.
	for length := 4; length <= 521; length++ {
		for _, distance := range []int{1, 2, 511, 512, 513, 1023, 1024, 2047, 2048, 4095, 4096, 8191, 8192, 16383, 16384, 32767, 32768, 65534, 65535} {
			e.match(nil, length, distance)
			for i := 0; i < length; i++ {
				want = append(want, want[len(want)-distance])
			}
		}
	}
	// The decoder stops at the declared output size. Include pending bits in a
	// full word; it must neither need nor consume the unused suffix.
	encoded := append([]byte(nil), e.output...)
	var tail [4]byte
	binary.BigEndian.PutUint32(tail[:], e.word)
	encoded = append(encoded, tail[:]...)
	got := make([]byte, len(want))
	if err := newMapDecoder(encoded).decode(got); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Fatal("length/distance code roundtrip differs")
	}
}

func TestMapEncodeSmallBlocks(t *testing.T) {
	for _, size := range []int{1, 2, 3, 4, 5, 6, 7, 64, 65, 521, 4096, 65535} {
		input := make([]byte, 150001+5)
		for i := range input[:len(input)-5] {
			input[i] = byte((i % 997) * 13)
		}
		e := newMapEncoder()
		var encoded []byte
		for i := 0; i < len(input)-5; i += size {
			encoded = append(encoded, e.block(input[i:], min(size, len(input)-5-i))...)
		}
		got := make([]byte, len(input)-5)
		if err := newMapDecoder(encoded).decode(got); err != nil {
			t.Fatalf("block size %d: %v", size, err)
		}
		if !bytes.Equal(got, input[:len(got)]) {
			t.Fatalf("block size %d differs", size)
		}
	}
}
