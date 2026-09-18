//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/memmap/nox/blobdata"
	"github.com/opennox/opennox/v1/internal/binfile"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/legacy/common/alloc/handles"
)

func prefabRuntimeTables(t *testing.T) {
	t.Helper()
	for off, data := range blobdata.PortTestPrefabRuntimeTables() {
		p := unsafe.Slice((*byte)(memmap.PtrOff(0x587000, off)), len(data))
		old := append([]byte(nil), p...)
		copy(p, data)
		t.Cleanup(func() { copy(p, old) })
	}
}

func TestPrefabRuntimeScripts(t *testing.T) {
	prefabRuntimeTables(t)
	handles.Init()
	defer handles.Release()
	path := filepath.Join(t.TempDir(), "instructions.bin")
	type row struct {
		Name     string
		Return   uint64
		Word     uint32
		Position int64
	}
	var rows []row
	for opcode, kinds := range prefabScriptKinds {
		count := *memmap.PtrUint8(0x587000, 218640+uintptr(opcode)*268)
		if int(count) != len(kinds) {
			t.Fatalf("shipped opcode %d argument count", opcode)
		}
		for j, kind := range kinds {
			if *memmap.PtrUint32(0x587000, 218648+uintptr(opcode)*268+uintptr(j)*8) != kind {
				t.Fatalf("shipped opcode %d kind %d", opcode, j)
			}
		}
		for _, nameLen := range []int{0, 1, 63, 1024} {
			for _, strLen := range []int{0, 1, 255} {
				for _, count := range []int32{-1, 0, 1, 3} {
					name := fmt.Sprintf("op%d/name%d/string%d/count%d", opcode, nameLen, strLen, count)
					var b bytes.Buffer
					binary.Write(&b, binary.LittleEndian, uint32(nameLen))
					b.Write(bytes.Repeat([]byte{'n'}, nameLen))
					const output uint32 = 0x87654321
					binary.Write(&b, binary.LittleEndian, output)
					binary.Write(&b, binary.LittleEndian, count)
					for i := int32(0); i < count; i++ {
						b.WriteByte(byte(opcode))
						b.WriteByte(0xa5)
						for _, kind := range kinds {
							switch kind {
							case 0, 3, 4, 5, 6:
								b.Write(bytes.Repeat([]byte{0x5a}, 4))
							case 1:
								b.Write(bytes.Repeat([]byte{0x6b}, 8))
							case 2, 7:
								b.WriteByte(byte(strLen))
								b.Write(bytes.Repeat([]byte{'s'}, strLen))
							default:
								t.Fatal("unexpected shipped argument kind")
							}
						}
					}
					expectedPos := int64(b.Len())
					b.Write([]byte{0xde, 0xad, 0xbe, 0xef})
					if err := os.WriteFile(path, b.Bytes(), 0600); err != nil {
						t.Fatal(err)
					}
					f, err := binfile.BinfileOpen(path, binfile.ReadOnly)
					if err != nil {
						t.Fatal(err)
					}
					if err = f.SetKey(-1); err != nil {
						t.Fatal(err)
					}
					handle := legacy.NewFileHandle(f.File)
					out, free := alloc.Make([]uint32{}, 3)
					out[0] = 0xa5a5a5a5
					out[2] = 0x5a5a5a5a
					result := legacy.PortTestPrefabCall(1, [6]uint32{uint32(uintptr(unsafe.Pointer(handle))), uint32(uintptr(unsafe.Pointer(&out[1])))})
					pos, err := f.Seek(0, io.SeekCurrent)
					if err != nil {
						t.Fatal(err)
					}
					if result != uint64(uint32(count)) || out[1] != output || pos != expectedPos || out[0] != 0xa5a5a5a5 || out[2] != 0x5a5a5a5a {
						t.Fatalf("%s result=%x word=%x position=%d want=%d", name, result, out[1], pos, expectedPos)
					}
					rows = append(rows, row{name, result, out[1], pos})
					free()
					legacy.Nox_fs_close(handle)
				}
			}
		}
	}
	spellbookCapture(t, "prefab-runtime-scripts", rows, "0f6a38490f447ac5b9dc12fe662c6155de5b61117b24e0c446e1006bfd888db4")
}

// Audited independently from the shipped descriptor table.
var prefabScriptKinds = [][]uint32{
	{1},    // OpenSecretWall
	{1},    // CloseSecretWall
	{0},    // Wait
	{},     // End
	{1},    // ToggleSecretWall
	{5, 6}, // MoveTo
	{5, 3}, // Orient
	{5},    // EnableObject
	{5},    // DisableObject
	{6},    // EnableWaypoint
	{6},    // DisableWaypoint
	{4, 6}, // AudioEvent
	{2},    // Print
	{},     // Repeat
	{5, 2}, // Say
	{5},    // ToggleObject
	{6},    // ToggleWaypoint
	{7},    // KillScript
	{2},    // PrintToAll
	{5},    // FollowNearestWaypointPath
	{5},    // WaitForEnable
	{1},    // BreakDestructableWall
	{5},    // UnlockDoor
	{5},    // PushLocation
	{5},    // PopLocation
	{5},    // GotoHome
	{5},    // DeleteObject
	{5},    // LockDoor
	{5},    // EnableGroup
	{5},    // DisableGroup
	{5},    // ToggleGroup
	{5},    // OpenSecretWallGroup
	{5},    // CloseSecretWallGroup
	{5},    // ToggleSecretWallGroup
	{5},    // BreakWallGroup
	{0},    // WaitFrames
}
