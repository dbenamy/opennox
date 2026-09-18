//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/legacy/common/alloc/handles"
)

func TestPrefabScriptsScalars(t *testing.T) {
	handles.Init()
	t.Cleanup(handles.Release)
	raw, f := prefabScriptsFiles(t, nil)
	type row struct {
		Name     string
		Return   uint64
		Position int64
		Output   []byte
	}
	var rows []row
	for _, v := range []uint32{0, 1, 127, 128, 255, 256, 0x12345678, 0x80000000, 0xffffffff, 0x7fc00000, 0x7f800000, 0xff800000} {
		for n := 0; n <= 4; n++ {
			for _, op := range []int{0, 1} {
				data := prefabScriptsWords(v)[:n]
				prefabScriptsResetFile(t, f[0], data)
				ret := legacy.PortTestPrefabScriptsCall(op, raw[0], nil, nil, 0)
				var padded [4]byte
				copy(padded[:], data)
				word := binary.LittleEndian.Uint32(padded[:])
				want := uint64(word)
				if op == 1 {
					want = math.Float64bits(float64(math.Float32frombits(word)))
				}
				pos, err := f[0].Seek(0, io.SeekCurrent)
				if err != nil {
					t.Fatal(err)
				}
				if ret != want || pos != int64(n) {
					t.Fatalf("read op%d/%x/%d got %x position%d want%x", op, v, n, ret, pos, want)
				}
				rows = append(rows, row{fmt.Sprintf("read%d/%x/%d", op, v, n), ret, pos, nil})
			}
		}
		for _, op := range []int{2, 3} {
			prefabScriptsResetFile(t, f[0], nil)
			ret := legacy.PortTestPrefabScriptsCall(op, raw[0], nil, nil, v)
			pos, err := f[0].Seek(0, io.SeekCurrent)
			if err != nil {
				t.Fatal(err)
			}
			out := prefabScriptsOutput(t, f[0])
			if ret != 1 || pos != 4 || !bytes.Equal(out, prefabScriptsWords(v)) {
				t.Fatalf("write op%d/%x return%x output%x", op, v, ret, out)
			}
			rows = append(rows, row{fmt.Sprintf("write%d/%x", op, v), ret, pos, out})
		}
	}
	spellbookCapture(t, "prefab-scripts-scalars", rows, "dbd866060474f624e7fde5101562bb32ac302c7c73bbbfe3d2297d1d23256796")
}

func TestPrefabScriptsMergeLifecycle(t *testing.T) {
	_, restoreGlobals := legacy.PortTestPrefabScriptsGlobals()
	defer restoreGlobals()
	handles.Init()
	t.Cleanup(handles.Release)
	first := prefabScriptFunction{Name: "GLOBAL", Code: []uint32{72}}
	second := prefabScriptFunction{Name: "GLOBAL", Vars: []uint32{1, 1, 1, 1}, Code: []uint32{72}}
	funcs := []prefabScriptFunction{first, second}
	for _, empty := range []int{-1, 0, 1, 2} {
		for _, existing := range []bool{false, true} {
			t.Run(fmt.Sprintf("empty%d/existing%t", empty, existing), func(t *testing.T) {
				count, cleanup := legacy.PortTestPrefabScriptsHandleBalance()
				defer cleanup()
				dir := t.TempDir()
				paths := []string{filepath.Join(dir, "a.obj"), filepath.Join(dir, "b.obj"), filepath.Join(dir, "merged.obj")}
				a := prefabScriptsEncode([]string{"A"}, funcs)
				b := prefabScriptsEncode([]string{"B"}, funcs)
				if empty == 0 || empty == 2 {
					a = nil
				}
				if empty == 1 || empty == 2 {
					b = nil
				}
				for i, data := range [][]byte{a, b} {
					if err := os.WriteFile(paths[i], data, 0600); err != nil {
						t.Fatal(err)
					}
				}
				if existing {
					if err := os.WriteFile(paths[2], bytes.Repeat([]byte{0x5a}, 4096), 0600); err != nil {
						t.Fatal(err)
					}
				}
				x, fx := alloc.CString(paths[0])
				defer fx()
				y, fy := alloc.CString(paths[1])
				defer fy()
				z, fz := alloc.CString(paths[2])
				defer fz()
				ret := nox_script_readWriteZzz_541670(x, y, z)
				if ret != 1 {
					t.Errorf("valid merge failed: %d", ret)
				}
				want := prefabScriptsEncode([]string{"A", "B"}, funcs)
				if len(a) == 0 {
					want = b
				} else if len(b) == 0 {
					want = a
				}
				out, err := os.ReadFile(paths[2])
				if err != nil {
					t.Error(err)
				} else if !bytes.Equal(out, want) {
					t.Errorf("merge output length=%d want=%d; stale tail or changed bytes", len(out), len(want))
				}
				if n := count(); n != 0 {
					t.Errorf("merge left %d file handles registered", n)
				}
			})
		}
	}
}

func TestPrefabScriptsBorrowedHandles(t *testing.T) {
	handles.Init()
	t.Cleanup(handles.Release)
	gl, restore := legacy.PortTestPrefabScriptsGlobals()
	defer restore()
	_ = gl
	funcs := []prefabScriptFunction{{Name: "GLOBAL", Code: []uint32{72}}, {Name: "GLOBAL", Vars: []uint32{1, 1, 1, 1}, Code: []uint32{72}}}
	a := prefabScriptsEncode([]string{"A"}, funcs)
	b := prefabScriptsEncode([]string{"B"}, funcs)
	raw, f := prefabScriptsFiles(t, a, b, nil)
	count, cleanup := legacy.PortTestPrefabScriptsHandleBalance()
	defer cleanup()
	legacy.Nox_script_readWriteWww_5417C0(f[0], f[1], f[2])
	if count() != 0 {
		t.Fatal("borrowed adapter registered extra handles")
	}
	for i := 0; i < 2; i++ {
		if _, err := f[i].Seek(0, io.SeekStart); err != nil {
			t.Fatalf("caller-owned input closed: %v", err)
		}
		if ret := legacy.PortTestPrefabScriptsCall(0, raw[i], nil, nil, 0); ret != uint64(binary.LittleEndian.Uint32([]byte("SCRI"))) {
			t.Fatal("borrowed handle invalid", ret)
		}
	}
	out := prefabScriptsOutput(t, f[2])
	if !bytes.Equal(out, prefabScriptsEncode([]string{"A", "B"}, funcs)) {
		t.Fatal("borrowed output differs")
	}
}
