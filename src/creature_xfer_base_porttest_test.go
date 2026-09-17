//go:build porttest

package opennox

import (
	"encoding/binary"
	"fmt"
	"github.com/opennox/opennox/v1/common/unit/ai"
	"os"
	"path/filepath"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/memmap/nox/blobdata"
	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func newCreatureXferOwner(t *testing.T) *server.Server {
	s := newItemXferOwner(t)
	t.Cleanup(legacy.PortTestCreatureXferLookupOwner())
	t.Cleanup(legacy.PortTestCreatureXferDefinitions(nil))
	t.Cleanup(legacy.PortTestCreatureXferVoiceSets(nil))
	for off, data := range blobdata.PortTestCreatureXferTables() {
		dst := unsafe.Slice((*byte)(memmap.PtrOff(0x587000, off)), len(data))
		old := append([]byte(nil), dst...)
		copy(dst, data)
		t.Cleanup(func() { copy(dst, old) })
	}
	table := unsafe.Slice(memmap.PtrUint32(0x587000, 261768), 72)
	oldTable := append([]uint32(nil), table...)
	for i := range table {
		p, free := alloc.CString(ai.ActionType(i).String())
		t.Cleanup(free)
		table[i] = uint32(uintptr(unsafe.Pointer(p)))
	}
	t.Cleanup(func() { copy(table, oldTable) })
	t.Cleanup(s.PortTestCreatureXferTypes())
	cache := memmap.PtrUint32(0x5D4594, 2487692)
	old := *cache
	*cache = 0
	t.Cleanup(func() { *cache = old })
	for i := 1; i < s.Types.Count(); i++ {
		objectTypeCode16ByInd[i] = uint16(1000 + i)
	}
	return s
}
func newCreatureXferObject(t *testing.T, s *server.Server, name string) *server.Object {
	u := newObjectXferTyped(t, s, "Creature"+name)
	hp, free := alloc.New(server.HealthData{})
	*hp = server.HealthData{Cur: 80, Field2: 100, Max: 100}
	u.HealthData = hp
	t.Cleanup(func() { u.HealthData = nil; free() })
	return u
}
func TestCreatureXferCurrentRecords(t *testing.T) {
	s := newCreatureXferOwner(t)
	path := filepath.Join(t.TempDir(), "creature.bin")
	var rows []struct {
		State itemXferCaptureRow
		Wire  []byte
	}
	defer func() { spellbookCapture(t, "creature-xfer-current", rows, "") }()
	for _, name := range []string{"Monster", "NPC"} {
		t.Run(name, func(t *testing.T) {
			u := newCreatureXferObject(t, s, name)
			objectXferSetCommon(u)
			u.Field34 = 777
			if err := cryptfile.OpenGlobal(path, cryptfile.WriteOnly, -1); err != nil {
				t.Fatal(err)
			}
			if err := u.CallXfer(nil); err != nil {
				t.Fatal(err)
			}
			written := cryptfile.Global().PortTestChecksum()
			if err := cryptfile.Close(); err != nil {
				t.Fatal(err)
			}
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			version := uint16(64)
			if name == "NPC" {
				version = 62
			}
			if len(data) < 21 || binary.LittleEndian.Uint16(data) != version {
				t.Fatal("outer version/common framing")
			}
			v := newCreatureXferObject(t, s, name)
			v.Field34 = 999
			if err := cryptfile.OpenGlobal(path, cryptfile.ReadOnly, -1); err != nil {
				t.Fatal(err)
			}
			defer cryptfile.Close()
			if err := v.CallXfer(nil); err != nil {
				t.Fatal(err)
			}
			pos, err := cryptfile.Global().File.Seek(0, 1)
			if err != nil || pos != int64(len(data)) {
				t.Fatalf("position=%d length=%d err=%v", pos, len(data), err)
			}
			if v.Extent != u.Extent || v.ScriptIDVal != u.ScriptIDVal || v.PosVec != u.PosVec || u.Field34 != 777 || v.Field34 != 999 {
				t.Fatal("common fields/lifetime")
			}
			var state []itemXferCaptureRow
			itemXferCaptureCase(t, &state, v, cryptfile.Global().PortTestChecksum(), written)
			rows = append(rows, struct {
				State itemXferCaptureRow
				Wire  []byte
			}{state[0], data})
			t.Logf("current record bytes=%d", len(data))
		})
	}
}
func TestCreatureXferFutureVersion(t *testing.T) {
	s := newCreatureXferOwner(t)
	path := filepath.Join(t.TempDir(), "future.bin")
	for _, name := range []string{"Monster", "NPC"} {
		max := uint16(64)
		if name == "NPC" {
			max = 62
		}
		for _, version := range []uint16{max + 1, 32767} {
			t.Run(fmt.Sprintf("%s-v%d", name, version), func(t *testing.T) {
				u := newCreatureXferObject(t, s, name)
				u.Field34 = 999
				p := new(mapDrawableStream)
				p.u16(version)
				p.u32(0x12345678)
				if err := os.WriteFile(path, p.Bytes(), 0600); err != nil {
					t.Fatal(err)
				}
				if err := cryptfile.OpenGlobal(path, cryptfile.ReadOnly, -1); err != nil {
					t.Fatal(err)
				}
				defer cryptfile.Close()
				if err := u.CallXfer(nil); err == nil {
					t.Fatal("accepted future version")
				}
				pos, err := cryptfile.Global().File.Seek(0, 1)
				if err != nil || pos != 2 || u.Field34 != 999 {
					t.Fatal("future version consumed/mutated state")
				}
			})
		}
	}
}
