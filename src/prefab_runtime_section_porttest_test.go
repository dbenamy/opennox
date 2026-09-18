//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"unsafe"

	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func TestPrefabRuntimeSectionUnknown(t *testing.T) {
	for _, name := range []string{"", "PrefabFixtureUnknown", "debugdata"} {
		for _, initial := range []uint32{0, 1, 0xa5a5a5a5} {
			ok, failed := legacy.PortTestPrefabReadSection(nil, name, initial)
			if ok != 0 || failed != 0 {
				t.Fatalf("unknown %q: recognized=%d failed=%d", name, ok, failed)
			}
		}
	}
}

func prefabDebugWire(version uint16, pairs [][2]string) []byte {
	var b bytes.Buffer
	binary.Write(&b, binary.LittleEndian, version)
	binary.Write(&b, binary.LittleEndian, uint32(len(pairs)))
	for _, pair := range pairs {
		for _, s := range pair {
			binary.Write(&b, binary.LittleEndian, uint32(len(s)))
			b.WriteString(s)
		}
	}
	return b.Bytes()
}

func TestPrefabRuntimeSectionDispatch(t *testing.T) {
	oldServer, oldFile := noxServer, cryptfile.Global()
	noxServer = &Server{Server: &server.Server{}}
	cryptfile.SetGlobal(nil)
	t.Cleanup(func() { cryptfile.Close(); cryptfile.SetGlobal(oldFile); noxServer = oldServer })
	t.Cleanup(noxflags.PortTestGameFlags(noxflags.GameHost))
	path := filepath.Join(t.TempDir(), "debug.bin")
	bounds, free := alloc.Make([]uint32{}, 8)
	defer free()
	for i := range bounds {
		bounds[i] = uint32(0xa5000000) + uint32(i)
	}
	wantBounds := append([]uint32(nil), bounds...)
	valid := prefabDebugWire(1, [][2]string{{"room", "north"}, {"room", "south"}, {"exit", "west"}})
	for _, tc := range []struct {
		name   string
		data   []byte
		ok     int
		failed uint32
	}{
		{"valid", valid, 1, 0},
		{"empty", prefabDebugWire(1, nil), 1, 0},
		{"version-zero", prefabDebugWire(0, nil), 0, 1},
		{"missing-version", nil, 0, 1},
		{"short-version", []byte{1}, 1, 0},
		{"short-value", valid[:len(valid)-1], 1, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			noxServer.Map.Debug.Reset()
			if err := os.WriteFile(path, tc.data, 0600); err != nil {
				t.Fatal(err)
			}
			if err := cryptfile.OpenGlobal(path, cryptfile.ReadOnly, -1); err != nil {
				t.Fatal(err)
			}
			ok, failed := legacy.PortTestPrefabReadSection(unsafe.Pointer(&bounds[0]), "DebugData", 0xa5a5a5a5)
			if ok != tc.ok || failed != tc.failed {
				t.Fatalf("recognized=%d failed=%d", ok, failed)
			}
			if !reflect.DeepEqual(bounds, wantBounds) {
				t.Fatal("section changed prefab bounds")
			}
			// The existing Go DebugData reader accepts a partial read with no I/O
			// error: a one-byte version supplies 1, and a short value is padded
			// then NUL-trimmed. This bridge preserves that reader's behavior.
			if tc.name == "short-value" && !reflect.DeepEqual(noxServer.Map.Debug.Get("exit"), []string{"wes"}) {
				t.Fatal("existing short-read behavior changed")
			}
			if tc.name == "valid" {
				if !reflect.DeepEqual(noxServer.Map.Debug.Get("room"), []string{"north", "south"}) || !reflect.DeepEqual(noxServer.Map.Debug.Get("exit"), []string{"west"}) {
					t.Fatal("real debug owner did not receive decoded entries")
				}
			}
			cryptfile.Close()
		})
	}
}

func TestPrefabRuntimeObjectCleanup(t *testing.T) {
	loaded := legacy.PortTestPrefabLoadedWord()
	oldLoaded := *loaded
	t.Cleanup(func() { *loaded = oldLoaded })
	// Use the real C constructor and Go cache destructor. Mark instantiated so the
	// destructor releases cache wrappers while objects retain their external owner.
	oldObject, oldWall, oldTile, oldWP := legacy.Get_dword_5d4594_1599540(), legacy.Get_dword_5d4594_1599532(), legacy.Get_dword_5d4594_1599556(), legacy.Get_dword_5d4594_1599548()
	oldState := legacy.Get_dword_5d4594_1599476()
	oldPath, oldOther := legacy.Get_dword_5d4594_1599588(), legacy.Get_dword_5d4594_1599592()
	path, freePath := alloc.Make([]byte{}, 2048)
	defer freePath()
	other, freeOther := alloc.Make([]byte{}, 2048)
	defer freeOther()
	legacy.Set_dword_5d4594_1599588(unsafe.Pointer(&path[0]))
	legacy.Set_dword_5d4594_1599592(unsafe.Pointer(&other[0]))
	legacy.Set_dword_5d4594_1599540(nil)
	legacy.Set_dword_5d4594_1599532(nil)
	legacy.Set_dword_5d4594_1599556(nil)
	legacy.Set_dword_5d4594_1599548(nil)
	legacy.Set_dword_5d4594_1599476(1)
	t.Cleanup(func() {
		legacy.Set_dword_5d4594_1599540(oldObject)
		legacy.Set_dword_5d4594_1599532(oldWall)
		legacy.Set_dword_5d4594_1599556(oldTile)
		legacy.Set_dword_5d4594_1599548(oldWP)
		legacy.Set_dword_5d4594_1599476(oldState)
		legacy.Set_dword_5d4594_1599588(oldPath)
		legacy.Set_dword_5d4594_1599592(oldOther)
	})
	// Own the reset routine's mapped words separately from extracted C globals.
	for _, off := range []uintptr{1599484, 1599488, 1599492, 1599496, 1599500, 1599504, 1599508, 1599512, 1599516, 1599520, 1599524, 1599528, 1599536, 1599544, 1599552, 1599560, 1599568, 1599572} {
		p := memmap.PtrUint32(0x5D4594, off)
		old := *p
		t.Cleanup(func() { *p = old })
	}
	obj, free := alloc.Make([]byte{}, 512)
	defer free()
	if legacy.PortTestPrefabObjectNode(unsafe.Pointer(&obj[0])) == nil {
		t.Fatal("C cache allocation failed")
	}
	s := &Server{Server: &server.Server{}}
	// Cover the other three C node constructors with their respective payload owners.
	for kind := 0; kind < 3; kind++ {
		node := legacy.PortTestPrefabCacheNode(kind)
		if node == nil {
			t.Fatal("cache node allocation failed")
		}
		payload := *(*unsafe.Pointer)(node)
		if kind == 2 {
			defer alloc.FreePtr(payload)
		}
	}
	s.Nox_xxx_free503F40()
	if legacy.Get_dword_5d4594_1599540() != nil {
		t.Fatal("cache wrapper survived cleanup")
	}
	s.Nox_xxx_free503F40()
}
