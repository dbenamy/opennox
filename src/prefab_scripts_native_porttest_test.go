//go:build porttest

package opennox

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/legacy/common/alloc/handles"
)

func TestPrefabScriptsNativeFileFailures(t *testing.T) {
	handles.Init()
	t.Cleanup(handles.Release)
	newObjectXferOwner(t)
	if ret := legacy.PortTestPrefabScriptsCall(15, nil, nil, nil, 73); ret != 73 {
		t.Fatal("empty pending list changed index", ret)
	}
	if got := legacy.PortTestPrefabScriptsName("Door\x00ignored", 7, -46, 92, false); got != "Door%7%-46%92" {
		t.Fatal("name C-string boundary", got)
	}
	prefabScriptsTables(t)
	_, restore := legacy.PortTestPrefabRuntimeGlobals()
	defer restore()
	prefix := unsafe.Slice((*byte)(memmap.PtrOff(0x973F18, 42152)), 2048)
	old := append([]byte(nil), prefix...)
	defer copy(prefix, old)
	dir := t.TempDir()
	clear(prefix)
	copy(prefix, dir)
	backup := filepath.Join(dir, "backup.obj")
	stale := []byte("old backup must not become a new script")
	if err := os.WriteFile(backup, stale, 0600); err != nil {
		t.Fatal(err)
	}
	input := filepath.Join(dir, "missing.obj")
	p, free := alloc.CString(input)
	defer free()
	xy := [2]int32{}
	if ret := legacy.PortTestPrefabScriptsCall(10, unsafe.Pointer(p), unsafe.Pointer(&xy), nil, 0); ret != 0 {
		t.Fatal("missing source accepted")
	}
	if _, err := os.Stat(input); !os.IsNotExist(err) {
		t.Fatal("source created from old backup", err)
	}
	if got, err := os.ReadFile(backup); err != nil || !bytes.Equal(got, stale) {
		t.Fatal("old backup changed", err)
	}
	for _, dir := range []string{"", strings.Repeat("x", 2048)} {
		p, free := alloc.CString(dir)
		ret := legacy.PortTestPrefabScriptsCall(16, unsafe.Pointer(p), nil, nil, 0)
		free()
		if ret != 0 {
			t.Fatal("invalid generation directory accepted")
		}
	}
	// These streams had no terminating instruction and could loop forever in C.
	for _, data := range [][]byte{nil, {4}, {4, 0, 0, 0}, {4, 0, 0, 0, 1, 0, 0, 0}} {
		raw, _ := prefabScriptsFiles(t, data, nil)
		if ret := legacy.PortTestPrefabScriptsCall(4, raw[0], raw[1], nil, 0); ret != 0 {
			t.Fatal("incomplete instruction accepted", ret)
		}
	}
}
