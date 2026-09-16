//go:build porttest

package cnxz

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Shipped compressed bytes are an independent oracle for the writer, including
// the three mixed-case filenames missed by the original directory-name loop.
func TestMapEncodeRealMaps(t *testing.T) {
	assets := os.Getenv("OPENNOX_MAP_DECODE_ASSETS")
	if assets == "" {
		t.Skip("asset path not configured")
	}
	files, err := filepath.Glob(filepath.Join(assets, "maps", "*", "*.nxz"))
	if err != nil {
		t.Fatal(err)
	}
	if len(files) == 0 {
		t.Fatal("no compressed map fixtures")
	}
	dir := t.TempDir()
	for _, path := range files {
		t.Run(filepath.Base(filepath.Dir(path)), func(t *testing.T) {
			src := strings.TrimSuffix(path, ".nxz") + ".map"
			want, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			out := filepath.Join(dir, "encoded.nxz")
			if err := CompressFile(src, out); err != nil {
				t.Fatal(err)
			}
			got, err := os.ReadFile(out)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(got, want) {
				t.Fatalf("compressed bytes differ: got %d bytes, want %d", len(got), len(want))
			}
		})
	}
}

func TestMapEncodeFiles(t *testing.T) {
	dir := t.TempDir()
	src, dst := filepath.Join(dir, "source"), filepath.Join(dir, "destination")
	sentinel := []byte("preserve on input error")
	for _, paths := range [][2]string{{"", dst}, {src, ""}, {src, dst}} {
		if err := os.WriteFile(dst, sentinel, 0600); err != nil {
			t.Fatal(err)
		}
		if CompressFile(paths[0], paths[1]) == nil {
			t.Fatal("invalid input accepted")
		}
		got, _ := os.ReadFile(dst)
		if !bytes.Equal(got, sentinel) {
			t.Fatal("destination changed on input error")
		}
	}
	data := bytes.Repeat([]byte("map compression same-path and truncation contract\x00"), 15000)
	if err := os.WriteFile(src, data, 0600); err != nil {
		t.Fatal(err)
	}
	if err := CompressFile(src, dst); err != nil {
		t.Fatal(err)
	}
	want, err := os.ReadFile(dst)
	if err != nil {
		t.Fatal(err)
	}
	unchanged, _ := os.ReadFile(src)
	if !bytes.Equal(unchanged, data) {
		t.Fatal("source changed")
	}
	if err := CompressFile(src, src); err != nil {
		t.Fatal(err)
	}
	got, _ := os.ReadFile(src)
	if !bytes.Equal(got, want) {
		t.Fatal("same-path output differs")
	}
	if err := DecompressFile(src, src); err != nil {
		t.Fatal(err)
	}
	got, _ = os.ReadFile(src)
	if !bytes.Equal(got, data) {
		t.Fatal("same-path round trip differs")
	}
}

func TestMapEncodeMatchBoundaries(t *testing.T) {
	dir := t.TempDir()
	var rows []mapDecodeCapture
	for seed := 0; seed < 48; seed++ {
		size := []int{4095, 4096, 4097, 65534, 65535, 65536, 65537, 131073, 499998, 500002, 1000003, 1048576}[seed%12]
		data := make([]byte, size)
		rng := uint32(seed + 1)
		for i := range data {
			rng = rng*1664525 + 1013904223
			data[i] = byte(rng >> 24)
			distance := []int{1, 3, 4, 5, 511, 512, 513, 1023, 1024, 32767, 65534, 65535}[i/523%12]
			if i >= distance && i%1103 < []int{3, 4, 5, 11, 12, 41, 42, 520, 521, 522, 1024, 1100}[seed%12] {
				data[i] = data[i-distance]
			}
			if seed >= 24 && i%10007 < 9000 {
				data[i] &= 3
			}
		}
		src, dst := filepath.Join(dir, "source.map"), filepath.Join(dir, "encoded.nxz")
		if err := os.WriteFile(src, data, 0600); err != nil {
			t.Fatal(err)
		}
		if err := CompressFile(src, dst); err != nil {
			t.Fatal(err)
		}
		encoded, err := os.ReadFile(dst)
		if err != nil {
			t.Fatal(err)
		}
		rows = append(rows, decodeMapContract(t, dir, fmt.Sprintf("matches/%d", seed), encoded, data))
	}
	captureMapDecode(t, "encode-matches", rows, "c8c09c83ad7cd5b40da7986ececb7f59fef83ca902d91a1491fcf659bf19728f")
}
