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

	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/legacy/common/alloc/handles"
	"github.com/opennox/opennox/v1/server"
)

func TestPrefabRuntimeLoader(t *testing.T) {
	handles.Init()
	defer handles.Release()
	s := newObjectXferOwner(t)
	s.MapGroups.Init()
	defer s.MapGroups.Free()
	defer noxflags.PortTestGameFlags(0x600000)()
	words, restore := legacy.PortTestPrefabRuntimeGlobals()
	defer restore()
	prefabRuntimeTables(t)
	prefabRuntimePaths(t, words)
	for _, span := range [][3]uintptr{{0x973F18, 35912, 72}, {0x5D4594, 739980, 8}} {
		p := unsafe.Slice((*byte)(memmap.PtrOff(span[0], span[1])), int(span[2]))
		old := bytes.Clone(p)
		defer copy(p, old)
	}
	metadata, free := alloc.Make([]byte{}, 76)
	defer free()
	*words["metadata"] = uint32(uintptr(unsafe.Pointer(&metadata[0])))
	*words["count"] = 1
	path := filepath.Join(t.TempDir(), "prefab.lib")
	copy(unsafe.Slice((*byte)(unsafe.Pointer(uintptr(*words["path"]))), 2048), path)
	type row struct {
		Name                     string
		Return                   uint64
		Loaded, Selected, Placed uint32
		Closed                   bool
		Waypoint                 []uint32
		Groups                   []prefabGroupRecord
	}
	var rows []row
	for _, version := range []byte{0, 1, 2, 255} {
		for _, kind := range []string{"empty", "debug", "waypoint", "group", "combined", "bad-magic", "unknown-object", "rejected-section"} {
			var b bytes.Buffer
			put := func(v any) { binary.Write(&b, binary.LittleEndian, v) }
			put(uint32(0))
			b.WriteByte(4)
			b.WriteString("room")
			b.WriteByte(1)
			b.WriteByte(version)
			put(math.Float32bits(100))
			put(math.Float32bits(200))
			if version > 1 {
				put(uint32(7))
				b.WriteString("skip-me")
			}
			magic := uint32(0xcafedead)
			if kind == "bad-magic" {
				magic = 0
			}
			put(magic)
			put(uint32(32))
			put(uint32(32))
			for _, v := range []uint32{46, 0, 46, 92, 0, 46, 92, 46} {
				put(v)
			}
			var sections bytes.Buffer
			section := func(name string, data []byte) {
				sections.WriteByte(byte(len(name) + 1))
				sections.WriteString(name)
				sections.WriteByte(0)
				binary.Write(&sections, binary.LittleEndian, uint32(len(data)))
				sections.Write(data)
			}
			if kind == "debug" || kind == "combined" {
				section("DebugData", prefabDebugWire(1, [][2]string{{"room", "north"}}))
			}
			if kind == "rejected-section" {
				section("DebugData", prefabDebugWire(0, nil))
			}
			if kind == "unknown-object" {
				section("PrefabUnknownType", nil)
			}
			if kind == "waypoint" || kind == "combined" {
				var data bytes.Buffer
				for _, v := range []any{uint16(4), uint32(1), uint32(101), float32(10), float32(15), byte(1), byte('w'), uint32(1), byte(0)} {
					binary.Write(&data, binary.LittleEndian, v)
				}
				section("WayPoints", data.Bytes())
			}
			if kind == "group" || kind == "combined" {
				var data bytes.Buffer
				for _, v := range []any{uint16(3), uint32(1), byte(2), byte('g'), byte(0), byte(1), uint32(201), uint32(1), uint32(101)} {
					binary.Write(&data, binary.LittleEndian, v)
				}
				section("GroupData", data.Bytes())
			}
			sections.WriteByte(0)
			for _, v := range sections.Bytes() {
				b.WriteByte(v ^ 126)
			}
			if err := os.WriteFile(path, b.Bytes(), 0600); err != nil {
				t.Fatal(err)
			}
			*words["loaded"] = math.MaxUint32
			*words["placed"] = 0
			ret := legacy.PortTestPrefabCall(14, [6]uint32{})
			closed := *words["file"] == 0 && cryptfile.Global() == nil
			r := row{Name: fmt.Sprintf("v%d/%s", version, kind), Return: ret, Loaded: *words["loaded"], Selected: *words["selected"], Placed: *words["placed"], Closed: closed}
			success := kind != "bad-magic" && kind != "unknown-object" && kind != "rejected-section"
			if ret != uint64(bool2int(success)) {
				t.Fatalf("%s return %d", r.Name, ret)
			}
			if !closed {
				t.Errorf("%s left prefab file open", r.Name)
				legacy.PortTestPrefabCall(10, [6]uint32{})
			}
			if success && (*words["loaded"] != 0 || *words["selected"] != 0 || *words["placed"] != 0) {
				t.Fatal("loaded cache state")
			}
			if kind == "waypoint" || kind == "combined" {
				if *words["waypoints"] == 0 {
					t.Fatalf("%s missing waypoint cache", r.Name)
				}
				wp := *(**server.Waypoint)(unsafe.Pointer(uintptr(*words["waypoints"])))
				r.Waypoint = []uint32{wp.Index, math.Float32bits(wp.PosVec.X), math.Float32bits(wp.PosVec.Y), wp.Flags}
				if wp.Index != 101 || wp.PosVec.X != -737 || wp.PosVec.Y != -732 || wp.ID() != "w" {
					t.Fatalf("%s waypoint bounds %+v", r.Name, r.Waypoint)
				}
			}
			for p := s.MapGroups.Refs; p != nil; p = p.Next4 {
				r.Groups = append(r.Groups, prefabGroupSnapshot(p.Field0))
			}
			if (kind == "group" || kind == "combined") && len(r.Groups) != 1 {
				t.Fatalf("%s missing group cache", r.Name)
			}
			rows = append(rows, r)
			noxServer.Nox_xxx_free503F40()
			inverted := legacy.PortTestPrefabCall(8, [6]uint32{})
			if inverted != uint64(bool2int(!success)) || *words["file"] != 0 || cryptfile.Global() != nil {
				t.Fatalf("%s selection wrapper result/close differs", r.Name)
			}
			noxServer.Nox_xxx_free503F40()
			for _, invalid := range []uint32{1, math.MaxUint32, 0x80000000} {
				*words["selected"] = 0x12345678
				if legacy.PortTestPrefabCall(8, [6]uint32{invalid}) != 0 || *words["selected"] != 0x12345678 {
					t.Fatal("invalid selection changed cache index")
				}
			}
		}
	}
	spellbookCapture(t, "prefab-runtime-loader", rows, "cef02ad3be4b23437be6a875e7673108a0709803c50c1c669565aae9855ce0b6")
}
