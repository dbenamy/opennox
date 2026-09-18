//go:build porttest

package opennox

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/legacy/common/alloc/handles"
)

func TestPrefabRuntimeFileControls(t *testing.T) {
	handles.Init()
	defer handles.Release()
	old := cryptfile.Global()
	cryptfile.SetGlobal(nil)
	defer func() { cryptfile.Close(); cryptfile.SetGlobal(old) }()
	words, restore := legacy.PortTestPrefabRuntimeGlobals()
	defer restore()
	metadata, free := alloc.Make([]byte{}, 3*76)
	defer free()
	*words["metadata"] = uint32(uintptr(unsafe.Pointer(&metadata[0])))
	*words["count"] = 3
	offsets := []uint32{0, 7, 4096}
	for i, off := range offsets {
		binary.LittleEndian.PutUint32(metadata[i*76+72:], off)
	}
	path := filepath.Join(t.TempDir(), "library.bin")
	if err := os.WriteFile(path, []byte("first-file"), 0600); err != nil {
		t.Fatal(err)
	}
	pathArg, freePath := alloc.CString(path)
	defer freePath()
	missingArg, freeMissing := alloc.CString(path + ".missing")
	defer freeMissing()
	arg := func(p *byte) [6]uint32 { return [6]uint32{uint32(uintptr(unsafe.Pointer(p)))} }
	type row struct {
		Name     string
		Return   uint32
		Position int64
		Closed   bool
	}
	var rows []row
	if legacy.PortTestPrefabCall(9, arg(missingArg)) != 0 || *words["file"] != 0 || cryptfile.Global() != nil {
		t.Fatal("missing file opened")
	}
	rows = append(rows, row{Name: "missing", Closed: true})
	if legacy.PortTestPrefabCall(11, [6]uint32{}) != 0 || legacy.PortTestPrefabCall(10, [6]uint32{}) != 0 {
		t.Fatal("closed file controls")
	}
	for round := 0; round < 2; round++ {
		if legacy.PortTestPrefabCall(9, arg(pathArg)) != 1 {
			t.Fatal("open failed")
		}
		handle := *words["file"]
		cf := cryptfile.Global()
		if handle == 0 || cf == nil {
			t.Fatal("file owner missing")
		}
		pos := func() int64 {
			n, e := cf.File.File.Seek(0, io.SeekCurrent)
			if e != nil {
				t.Fatal(e)
			}
			return n
		}
		for _, index := range []int32{-1, 0, 1, 2, 3, math.MinInt32, math.MaxInt32} {
			before := pos()
			got := uint32(legacy.PortTestPrefabCall(11, [6]uint32{uint32(index)}))
			want, position := uint32(0), before
			if index >= 0 && index < 3 {
				want = handle
				position = int64(offsets[index])
			}
			if got != want || pos() != position {
				t.Fatalf("seek %d: ret=%x pos=%d want=%x/%d", index, got, pos(), want, position)
			}
			normalized := uint32(0)
			if got != 0 {
				normalized = 1
			}
			rows = append(rows, row{fmt.Sprintf("round%d/index%d", round, index), normalized, pos(), false})
		}
		// An already open file wins over the new path and is rewound.
		if legacy.PortTestPrefabCall(9, arg(missingArg)) != 1 || cryptfile.Global() != cf || *words["file"] != handle || pos() != 0 {
			t.Fatal("reopen did not preserve and reset current file")
		}
		b, err := cf.ReadU8()
		if err != nil || b != 'f' {
			t.Fatal("reopened wrong data")
		}
		rows = append(rows, row{fmt.Sprintf("round%d/reopen", round), 1, pos(), false})
		if uint32(legacy.PortTestPrefabCall(10, [6]uint32{})) != handle || *words["file"] != 0 || cryptfile.Global() != nil {
			t.Fatal("close return or ownership")
		}
		if legacy.PortTestPrefabCall(10, [6]uint32{}) != 0 {
			t.Fatal("second close")
		}
		rows = append(rows, row{Name: fmt.Sprintf("round%d/close", round), Return: 1, Closed: true})
	}
	spellbookCapture(t, "prefab-runtime-file-controls", rows, "344ab19b7f1347660d72d0a2df78b69e61b6b0d5254ebec6b086a3a54cd4f104")
}
