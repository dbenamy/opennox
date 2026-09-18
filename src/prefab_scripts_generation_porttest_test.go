//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/legacy/common/alloc/handles"
)

func TestPrefabScriptsGeneration(t *testing.T) {
	handles.Init()
	defer handles.Release()
	newObjectXferOwner(t)
	prefabScriptsTables(t)
	words, restore := legacy.PortTestPrefabRuntimeGlobals()
	defer restore()
	paint, restorePaint := legacy.PortTestPrefabScriptsPaintGlobals()
	defer restorePaint()
	prefabRuntimePaths(t, words)
	metadata, free := alloc.Make([]byte{}, 2048*76)
	defer free()
	*words["metadata"] = uint32(uintptr(unsafe.Pointer(&metadata[0])))
	oldFile := cryptfile.Global()
	cryptfile.SetGlobal(nil)
	defer func() { cryptfile.Close(); cryptfile.SetGlobal(oldFile) }()
	var restores []func()
	defer func() {
		for i := len(restores) - 1; i >= 0; i-- {
			restores[i]()
		}
	}()
	snapshot := func(base, off uintptr, n int) []byte {
		p := unsafe.Slice((*byte)(memmap.PtrOff(base, off)), n)
		old := append([]byte(nil), p...)
		restores = append(restores, func() { copy(p, old) })
		return p
	}
	snapshot(0x973F18, 35880, 44200-35880)
	snapshot(0x5D4594, 1599472, 1599648-1599472)
	snapshot(0x5D4594, 1549772, 8)
	oldTicks := legacy.PlatformTicks
	legacy.PlatformTicks = func() uint64 { return 0x10000002a }
	defer func() { legacy.PlatformTicks = oldTicks }()
	selection := unsafe.Slice((*byte)(memmap.PtrOff(0x973F18, 35912)), 72)
	prefix := unsafe.Slice((*byte)(memmap.PtrOff(0x973F18, 42152)), 2048)
	type row struct {
		Name                                                    string
		Return                                                  uint64
		Selection                                               []byte
		Paint                                                   map[string]uint32
		Count, Loaded, Placed, Selected, Instance, LastWaypoint uint32
		Metadata                                                []byte
		Ticks                                                   uint64
		Closed                                                  bool
	}
	var rows []row
	for _, kind := range []string{"missing", "bad-magic", "valid"} {
		dir := t.TempDir()
		for _, name := range []string{"oldObj.obj", "blend.obj"} {
			if err := os.WriteFile(filepath.Join(dir, name), []byte("stale"), 0600); err != nil {
				t.Fatal(err)
			}
		}
		var wantMeta []byte
		if kind != "missing" {
			data := []byte{0, 0, 0, 0}
			if kind == "valid" {
				var b bytes.Buffer
				binary.Write(&b, binary.LittleEndian, uint32(0xcafedead))
				binary.Write(&b, binary.LittleEndian, uint32(15))
				b.WriteByte(4)
				b.WriteString("Room")
				b.Write([]byte{7, 1})
				binary.Write(&b, binary.LittleEndian, float32(12.25))
				binary.Write(&b, binary.LittleEndian, float32(-7.5))
				binary.Write(&b, binary.LittleEndian, uint32(0))
				data = b.Bytes()
				wantMeta = make([]byte, 76)
				copy(wantMeta, "Room")
				binary.LittleEndian.PutUint32(wantMeta[64:], math.Float32bits(12.25))
				binary.LittleEndian.PutUint32(wantMeta[68:], math.Float32bits(-7.5))
				binary.LittleEndian.PutUint32(wantMeta[72:], 4)
			}
			if err := os.WriteFile(filepath.Join(dir, "AreaMap.lib"), data, 0600); err != nil {
				t.Fatal(err)
			}
		}
		clear(metadata)
		for i := range selection {
			selection[i] = 0xa5
		}
		for _, p := range paint {
			*p = 0xa5a5a5a5
		}
		for _, key := range []string{"selected", "loaded", "placed", "instance", "lastWaypoint"} {
			*words[key] = 17
		}
		*memmap.PtrUint32(0x973F18, 35880) = 19
		*memmap.PtrUint32(0x5D4594, 1599580) = 21
		p, fp := alloc.CString(dir)
		ret := legacy.PortTestPrefabScriptsCall(16, unsafe.Pointer(p), nil, nil, 0)
		fp()
		wantReturn := uint64(0)
		wantCount := uint32(0)
		if kind == "valid" {
			wantReturn = 1
			wantCount = 1
		}
		if ret != wantReturn || *words["count"] != wantCount {
			t.Fatalf("%s return/count %d/%d", kind, ret, *words["count"])
		}
		wantSelection := make([]byte, 72)
		wantSelection[60] = 2
		if !bytes.Equal(selection, wantSelection) {
			t.Fatalf("selection %x", selection)
		}
		gotPaint := map[string]uint32{}
		for key, p := range paint {
			want := uint32(0)
			switch key {
			case "dword_5d4594_3835356":
				want = 255
			case "dword_5d4594_3835364", "dword_5d4594_3835368", "dword_5d4594_3835372", "dword_5d4594_3835392":
				want = 1
			}
			if *p != want {
				t.Fatalf("%s %d want%d", key, *p, want)
			}
			gotPaint[key] = *p
		}
		if *words["selected"] != math.MaxUint32 || *words["loaded"] != math.MaxUint32 || *words["placed"] != 0 || *words["instance"] != 0 || *words["lastWaypoint"] != 0 {
			t.Fatal("generation counters")
		}
		if memmap.Uint32(0x973F18, 35880) != 0 || memmap.Uint32(0x5D4594, 1599580) != 0 {
			t.Fatal("script counters")
		}
		if !bytes.Equal(metadata[:len(wantMeta)], wantMeta) {
			t.Fatal("library metadata")
		}
		for _, name := range []string{"oldObj.obj", "blend.obj"} {
			if _, err := os.Stat(filepath.Join(dir, name)); !os.IsNotExist(err) {
				t.Fatal("stale script remains", err)
			}
		}
		if alloc.GoString(&prefix[0]) != dir {
			t.Fatal("generation prefix")
		}
		for _, key := range []string{"path", "alternate"} {
			if got := alloc.GoString((*byte)(unsafe.Pointer(uintptr(*words[key])))); got != dir+"\\AreaMap.lib" {
				t.Fatalf("%s path %q", key, got)
			}
		}
		ticks := *memmap.PtrUint64(0x5D4594, 1549772)
		// The live C clock bridge returns uint32 before widening into this field.
		if ticks != 42 {
			t.Fatal("timestamp", ticks)
		}
		closed := *words["file"] == 0 && cryptfile.Global() == nil
		if !closed {
			t.Fatal("library handle left open")
		}
		rows = append(rows, row{fmt.Sprint(kind), ret, append([]byte(nil), selection...), gotPaint, *words["count"], *words["loaded"], *words["placed"], *words["selected"], *words["instance"], *words["lastWaypoint"], append([]byte(nil), metadata[:len(wantMeta)]...), ticks, closed})
	}
	spellbookCapture(t, "prefab-scripts-generation", rows, "ce1620e8a963bba41f4e1913d30ae4123fd0214f2cb6d0126d9cc4d0d4e46ee6")
}
