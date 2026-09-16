//go:build porttest

package cnxz

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type mapDecodeCapture struct {
	Name                 string
	Size                 int
	Compressed, Expanded [32]byte
}

func captureMapDecode(t *testing.T, label string, rows []mapDecodeCapture, want string) {
	t.Helper()
	raw, err := json.Marshal(rows)
	if err != nil {
		t.Fatal(err)
	}
	if base := os.Getenv("OPENNOX_MAP_DECODE_CAPTURE"); base != "" {
		if err := os.WriteFile(base+"-"+label+".json", raw, 0600); err != nil {
			t.Fatal(err)
		}
	}
	got := fmt.Sprintf("%x", sha256.Sum256(raw))
	t.Logf("%s %s", label, got)
	if want != "" && got != want {
		t.Fatalf("capture mismatch %s got %s want %s", label, got, want)
	}
}
func decodeMapContract(t *testing.T, dir, name string, compressed, expected []byte) mapDecodeCapture {
	t.Helper()
	src, dst := filepath.Join(dir, "input.nxz"), filepath.Join(dir, "output.map")
	if err := os.WriteFile(src, compressed, 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(dst, []byte("old destination must be truncated"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := DecompressFile(src, dst); err != nil {
		t.Fatalf("%s: %v", name, err)
	}
	got, err := os.ReadFile(dst)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, expected) {
		t.Fatalf("%s: expanded bytes differ got %d want %d", name, len(got), len(expected))
	}
	source, _ := os.ReadFile(src)
	if !bytes.Equal(source, compressed) {
		t.Fatal("compressed input changed")
	}
	return mapDecodeCapture{name, len(got), sha256.Sum256(compressed), sha256.Sum256(got)}
}
func TestMapDecodeSynthetic(t *testing.T) {
	dir := t.TempDir()
	var rows []mapDecodeCapture
	for _, size := range []int{1, 2, 3, 4, 7, 11, 12, 31, 255, 511, 512, 521, 522, 8191, 8192, 8193, 65535, 65536, 65537, 499999, 500000, 500001, 1000001} {
		for pattern := 0; pattern < 5; pattern++ {
			data := make([]byte, size)
			rng := uint32(0x7ac01823)
			for i := range data {
				rng = rng*1664525 + 1013904223
				switch pattern {
				case 0:
					data[i] = 0
				case 1:
					data[i] = "aabacaba"[i%8]
				case 2:
					data[i] = byte(i)
				case 3:
					data[i] = byte(rng >> 24)
				case 4:
					if i%17 == 0 {
						data[i] = byte(rng >> 24)
					} else {
						data[i] = 32
					}
				}
			}
			src, comp := filepath.Join(dir, "original.map"), filepath.Join(dir, "compressed.nxz")
			if err := os.WriteFile(src, data, 0600); err != nil {
				t.Fatal(err)
			}
			if err := CompressFile(src, comp); err != nil {
				t.Fatal(err)
			}
			encoded, err := os.ReadFile(comp)
			if err != nil {
				t.Fatal(err)
			}
			if len(encoded) < 4 || int(binary.LittleEndian.Uint32(encoded)) != len(data) {
				t.Fatal("wrong compressed header")
			}
			rows = append(rows, decodeMapContract(t, dir, fmt.Sprintf("%d/%d", size, pattern), encoded, data))
		}
	}
	captureMapDecode(t, "synthetic", rows, "56d40cb714d8176aede00dc93fff70021ca6b80c28b02687c6b114041a742244")
}
func TestMapDecodeRealMaps(t *testing.T) {
	assets := os.Getenv("OPENNOX_MAP_DECODE_ASSETS")
	if assets == "" {
		t.Skip("asset path not configured")
	}
	files, err := filepath.Glob(filepath.Join(assets, "maps", "*", "*.nxz"))
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	var rows []mapDecodeCapture
	for _, path := range files {
		encoded, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		expected, err := os.ReadFile(strings.TrimSuffix(path, ".nxz") + ".map")
		if err != nil {
			t.Fatal(err)
		}
		name, err := filepath.Rel(filepath.Join(assets, "maps"), path)
		if err != nil {
			t.Fatal(err)
		}
		rows = append(rows, decodeMapContract(t, dir, name, encoded, expected))
	}
	if len(rows) == 0 {
		t.Fatal("no real compressed map fixtures")
	}
	captureMapDecode(t, "real-maps", rows, "20b230bc3111761f92159aba009dadf144259e3b97aa60e1b7ba01db7edb590e")
}
func TestMapDecodeFiles(t *testing.T) {
	dir := t.TempDir()
	src, dst := filepath.Join(dir, "source"), filepath.Join(dir, "destination")
	sentinel := []byte("unchanged")
	for _, paths := range [][2]string{{"", dst}, {src, ""}, {src, dst}} {
		if err := os.WriteFile(dst, sentinel, 0600); err != nil {
			t.Fatal(err)
		}
		if DecompressFile(paths[0], paths[1]) == nil {
			t.Fatal("invalid file paths accepted")
		}
		got, _ := os.ReadFile(dst)
		if !bytes.Equal(got, sentinel) {
			t.Fatal("failed input changed destination")
		}
	}
	for size := 0; size < 4; size++ {
		if err := os.WriteFile(src, make([]byte, size), 0600); err != nil {
			t.Fatal(err)
		}
		if DecompressFile(src, dst) == nil {
			t.Fatal("short header accepted")
		}
		got, _ := os.ReadFile(dst)
		if !bytes.Equal(got, sentinel) {
			t.Fatal("short header changed destination")
		}
	}
	var w mapDecodeBits
	w.symbol(t, 65)
	w.symbol(t, 273)
	encoded := append([]byte{1, 0, 0, 0}, w.data...)
	if err := os.WriteFile(src, encoded, 0600); err != nil {
		t.Fatal(err)
	}
	if DecompressFile(src, dir) == nil {
		t.Fatal("directory output accepted")
	}
	// All input is read before destination creation, including the same path.
	if err := DecompressFile(src, src); err != nil {
		t.Fatal(err)
	}
	got, _ := os.ReadFile(src)
	if !bytes.Equal(got, []byte{65}) {
		t.Fatal("same-path expansion failed")
	}
}

// Independent bit-coded files exercise end-block padding and cross-block history
// without asking the production compressor to choose those encodings.
type mapDecodeBits struct {
	data []byte
	used uint8
}

func (w *mapDecodeBits) put(v uint32, n int) {
	for bit := n - 1; bit >= 0; bit-- {
		if w.used == 0 {
			w.data = append(w.data, 0)
		}
		w.data[len(w.data)-1] |= byte((v>>bit)&1) << (7 - w.used)
		w.used = (w.used + 1) % 8
	}
}
func (w *mapDecodeBits) symbol(t *testing.T, sym int) {
	t.Helper()
	order := []int{}
	for i := 256; i < 272; i++ {
		order = append(order, i)
	}
	order = append(order, 0, 32, 48, 255)
	for i := 1; i < 274; i++ {
		if i == 32 || i == 48 || i >= 255 && i < 272 {
			continue
		}
		order = append(order, i)
	}
	index := -1
	for i, v := range order {
		if v == sym {
			index = i
			break
		}
	}
	if index < 0 {
		t.Fatal("fixture symbol")
	}
	offset := 0
	for prefix, n := range []int{2, 3, 3, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 5, 5, 5} {
		if index < offset+(1<<n) {
			w.put(uint32(prefix), 4)
			w.put(uint32(index-offset), n)
			return
		}
		offset += 1 << n
	}
	t.Fatal("fixture index")
}
func TestMapDecodeBlockHistory(t *testing.T) {
	dir := t.TempDir()
	var rows []mapDecodeCapture
	for _, prefix := range []string{"A", "abcd", strings.Repeat("abcdefg", 9363)} {
		for _, blocks := range []int{1, 2, 3} {
			var w mapDecodeBits
			expected := []byte(prefix)
			for _, b := range []byte(prefix) {
				w.symbol(t, int(b))
			}
			for i := 0; i < blocks; i++ {
				w.symbol(t, 273)
				w.used = 0       // C compression block padding is discarded.
				w.symbol(t, 263) // Eleven-byte back-reference, overlapping distance1.
				w.put(0, 3)
				w.put(1, 9)
				expected = append(expected, bytes.Repeat(expected[len(expected)-1:], 11)...)
			}
			w.symbol(t, 273)
			encoded := binary.LittleEndian.AppendUint32(nil, uint32(len(expected)))
			encoded = append(encoded, w.data...)
			rows = append(rows, decodeMapContract(t, dir, fmt.Sprintf("%d/%d", len(prefix), blocks), encoded, expected))
		}
	}
	captureMapDecode(t, "block-history", rows, "13fd63dd1cacaf77797269309a068b2c938cbe6eef97baf7038b56e00c494e8a")
}
