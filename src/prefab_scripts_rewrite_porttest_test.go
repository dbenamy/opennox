//go:build porttest

package opennox

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/memmap/nox/blobdata"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/legacy/common/alloc/handles"
	"github.com/opennox/opennox/v1/server"
)

func prefabScriptsTables(t *testing.T) {
	t.Helper()
	for off, data := range blobdata.PortTestPrefabScriptsTables() {
		dst := unsafe.Slice((*byte)(memmap.PtrOff(0x587000, off)), len(data))
		old := append([]byte(nil), dst...)
		copy(dst, data)
		t.Cleanup(func() { copy(dst, old) })
	}
}

func TestPrefabScriptsRewrite(t *testing.T) {
	handles.Init()
	defer handles.Release()
	newObjectXferOwner(t)
	prefabScriptsTables(t)
	words, restore := legacy.PortTestPrefabRuntimeGlobals()
	defer restore()
	prefix := unsafe.Slice((*byte)(memmap.PtrOff(0x973F18, 42152)), 2048)
	old := append([]byte(nil), prefix...)
	defer copy(prefix, old)
	nameBuf := unsafe.Slice((*byte)(memmap.PtrOff(0x5D4594, 2489164)), 256)
	oldName := append([]byte(nil), nameBuf...)
	defer copy(nameBuf, oldName)
	balance, cleanup := legacy.PortTestPrefabScriptsHandleBalance()
	defer cleanup()
	type row struct {
		Instance int32
		XY       [2]int32
		Missing  bool
		Return   uint64
		Output   []byte
		Handles  int
	}
	var rows []row
	for _, instance := range []int32{0, 7, -1, 2147483647} {
		for _, xy := range [][2]int32{{0, 0}, {-46, 92}, {2147483647, -2147483648}} {
			for _, missing := range []bool{false, true} {
				dir := t.TempDir()
				clear(prefix)
				copy(prefix, dir)
				path := filepath.Join(dir, "input.obj")
				funcs := []prefabScriptFunction{
					{Name: "GLOBAL", Code: []uint32{72}},
					{Name: "GLOBAL", Vars: []uint32{1, 1, 1, 1}, Code: []uint32{72}},
					{Name: "OnEnter", Vars: []uint32{1, 2}, Code: []uint32{4, 0x12345678, 72}},
					{Name: "OnExit", Code: []uint32{4, 0xffffff80, 72}},
				}
				strs := []string{"door", "route"}
				if !missing {
					if err := os.WriteFile(path, prefabScriptsEncode(strs, funcs), 0600); err != nil {
						t.Fatal(err)
					}
				}
				*words["instance"] = uint32(instance)
				p, free := alloc.CString(path)
				ret := legacy.PortTestPrefabScriptsCall(10, unsafe.Pointer(p), unsafe.Pointer(&xy), nil, 0)
				free()
				var output []byte
				if missing {
					if ret != 0 {
						t.Fatal("missing input return", ret)
					}
				} else {
					if ret != 1 {
						t.Fatal("rewrite return", ret)
					}
					var err error
					output, err = os.ReadFile(path)
					if err != nil {
						t.Fatal(err)
					}
					for i := 2; i < len(funcs); i++ {
						funcs[i].Name += fmt.Sprintf("%%%d%%%d%%%d", instance, xy[0], xy[1])
					}
					if want := prefabScriptsEncode(strs, funcs); !bytes.Equal(output, want) {
						t.Fatalf("rewrite differs: instance %d xy %v", instance, xy)
					}
					var vm server.NoxScriptVM
					if err := vm.ReadScript(bytes.NewReader(output)); err != nil {
						t.Fatal(err)
					}
					fn := vm.Funcs()[2]
					if fn.PosOff.X != int(xy[0]) || fn.PosOff.Y != int(xy[1]) {
						t.Fatal("VM coordinates", fn.PosOff)
					}
				}
				if _, err := os.Stat(filepath.Join(dir, "backup.obj")); !os.IsNotExist(err) {
					t.Fatal("backup remains", err)
				}
				if balance() != 0 {
					t.Fatal("file registrations remain", balance())
				}
				rows = append(rows, row{instance, xy, missing, ret, output, balance()})
			}
		}
	}
	spellbookCapture(t, "prefab-scripts-rewrite", rows, "9001bcea8d67d8ebd5311cc78d92443f99dbb29349a15b8776be9cdd1083c68b")
}
