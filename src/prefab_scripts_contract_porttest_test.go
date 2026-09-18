//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"unsafe"

	"github.com/opennox/noxscript/ns/asm"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/internal/binfile"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc/handles"
)

func prefabScriptsWords(words ...uint32) []byte {
	b := make([]byte, 4*len(words))
	for i, v := range words {
		binary.LittleEndian.PutUint32(b[4*i:], v)
	}
	return b
}
func prefabScriptsFiles(t *testing.T, data ...[]byte) ([]unsafe.Pointer, []*binfile.File) {
	t.Helper()
	var raw []unsafe.Pointer
	var files []*binfile.File
	dir := t.TempDir()
	for i, b := range data {
		p := filepath.Join(dir, fmt.Sprintf("script-%d.bin", i))
		if err := os.WriteFile(p, b, 0600); err != nil {
			t.Fatal(err)
		}
		f, err := os.OpenFile(p, os.O_RDWR, 0600)
		if err != nil {
			t.Fatal(err)
		}
		bf := binfile.NewFile(f)
		h := legacy.NewFileHandle(bf)
		t.Cleanup(func() { legacy.Nox_fs_close(h) })
		files = append(files, bf)
		raw = append(raw, unsafe.Pointer(h))
	}
	return raw, files
}
func prefabScriptsOutput(t *testing.T, f *binfile.File) []byte {
	t.Helper()
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		t.Fatal(err)
	}
	b, err := io.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	return b
}
func TestPrefabScriptsNames(t *testing.T) {
	p := unsafe.Slice((*byte)(memmap.PtrOff(0x5D4594, 2489164)), 256)
	old := append([]byte(nil), p...)
	defer copy(p, old)
	for _, name := range []string{"", "Door", "OnEnter", "Waypoint-2", strings.Repeat("n", 200), strings.Repeat("n", 250), strings.Repeat("n", 255), strings.Repeat("n", 256)} {
		for _, id := range []int32{0, 1, -1, 2147483647, -2147483648} {
			for _, objectOnly := range []bool{false, true} {
				got := legacy.PortTestPrefabScriptsName(name, id, -46, 92, objectOnly)
				want := fmt.Sprintf("%s%%%d%%-46%%92", name, id)
				if objectOnly {
					want = fmt.Sprintf("%s%%%d", name, id)
				}
				if len(want) >= 256 {
					want = "ERROR_NAME_TOO_LONG!"
				}
				if got != want {
					t.Errorf("name %q id %d object=%t: got %q want %q", name, id, objectOnly, got, want)
				}
			}
		}
	}
}
func TestPrefabScriptsOperandWidths(t *testing.T) {
	handles.Init()
	t.Cleanup(handles.Release)
	gl, restore := legacy.PortTestPrefabScriptsGlobals()
	defer restore()
	*gl[0], *gl[1], *gl[2] = 1000, 1000, 1000
	// These operands are 32-bit words in the real reader and VM instruction format.
	for _, op := range []uint32{3, 4, 6, 19, 20, 21} {
		for _, value := range []uint32{0, 127, 128, 255, 256, 65536, 0x12345678, 0xffffffff} {
			t.Run(fmt.Sprintf("op%d/%08x", op, value), func(t *testing.T) {
				words := []uint32{op, value, 72}
				if op == 3 {
					words = []uint32{op, 1, value, 72}
				}
				input := prefabScriptsWords(words...)
				if _, err := asm.Decode(words); err != nil {
					t.Fatalf("invalid independent instruction fixture: %v", err)
				}
				raw, f := prefabScriptsFiles(t, input, nil)
				if ret := legacy.PortTestPrefabScriptsCall(4, raw[0], raw[1], nil, 0); ret != 1 {
					t.Fatalf("copy return %d", ret)
				}
				if got := prefabScriptsOutput(t, f[1]); !bytes.Equal(got, input) {
					t.Errorf("instruction changed without remapping: got %x want %x", got, input)
				}
			})
		}
	}
}
