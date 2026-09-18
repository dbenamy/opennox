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

	"github.com/opennox/libs/datapath"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc/handles"
)

func TestPrefabRuntimeIntro(t *testing.T) {
	prefabRuntimeTables(t)
	words, restore := legacy.PortTestPrefabRuntimeGlobals()
	defer restore()
	defer handles.PortTestInit()()
	defer noxflags.PortTestGameFlags(0)()
	dir := t.TempDir()
	oldData := datapath.Data()
	datapath.SetData(dir)
	defer datapath.SetData(oldData)
	name := unsafe.Slice((*byte)(memmap.PtrOff(0x85B3FC, 36)), 80)
	oldName := append([]byte(nil), name...)
	clear(name)
	copy(name, "fixture")
	defer copy(name, oldName)
	oldFile := cryptfile.Global()
	cryptfile.SetGlobal(nil)
	defer func() { cryptfile.Close(); cryptfile.SetGlobal(oldFile) }()
	path := filepath.Join(dir, "intro.bin")
	introPath := filepath.Join(dir, "maps", "fixture", "fixture.txt")
	if err := os.MkdirAll(filepath.Dir(introPath), 0700); err != nil {
		t.Fatal(err)
	}
	type row struct {
		Name               string
		Return             uint64
		Position           int64
		Memory, File, Wire []byte
		FileExists         bool
		Freed              bool
	}
	var rows []row
	for _, version := range []uint16{0, 1, 2, 0x7fff, 0x8000} {
		for _, flags := range []uint32{0, 0x200000, 0x400000, 0x600000} {
			for _, length := range []int{0, 1, 16, 1024} {
				if err := os.Remove(introPath); err != nil && !os.IsNotExist(err) {
					t.Fatal(err)
				}
				noxflags.ResetGame()
				noxflags.SetGame(noxflags.GameFlag(flags))
				data := make([]byte, length)
				for i := range data {
					data[i] = byte(17 + 31*i)
				}
				var wire bytes.Buffer
				binary.Write(&wire, binary.LittleEndian, version)
				binary.Write(&wire, binary.LittleEndian, uint32(length))
				wire.Write(data)
				if err := os.WriteFile(path, wire.Bytes(), 0600); err != nil {
					t.Fatal(err)
				}
				if err := cryptfile.OpenGlobal(path, cryptfile.ReadOnly, -1); err != nil {
					t.Fatal(err)
				}
				ret := legacy.PortTestPrefabCall(32, [6]uint32{})
				pos, err := cryptfile.Global().File.File.Seek(0, io.SeekCurrent)
				if err != nil {
					t.Fatal(err)
				}
				r := row{Name: fmt.Sprintf("read/v%d/f%x/n%d", version, flags, length), Return: ret, Position: pos}
				if p := *words["intro"]; p != 0 {
					r.Memory = append([]byte(nil), unsafe.Slice((*byte)(unsafe.Pointer(uintptr(p))), length)...)
				}
				r.File, err = os.ReadFile(introPath)
				r.FileExists = err == nil
				if err != nil && !os.IsNotExist(err) {
					t.Fatal(err)
				}
				valid := int16(version) >= 1
				wantRet := uint64(0)
				wantPos := int64(2)
				if valid {
					wantRet = 1
					wantPos = int64(6 + length)
				}
				if ret != wantRet || pos != wantPos {
					t.Fatalf("%s return/cursor", r.Name)
				}
				if valid && length > 0 && flags&0x400000 == 0 {
					if flags&0x200000 != 0 {
						if !r.FileExists || !bytes.Equal(r.File, data) || r.Memory != nil {
							t.Fatal("intro editor file")
						}
					} else if !bytes.Equal(r.Memory, data) || r.FileExists {
						t.Fatal("intro memory")
					}
				} else if r.Memory != nil || r.FileExists {
					t.Fatal("intro empty/rejected/skip side effect")
				}
				p := *words["intro"]
				freed := legacy.PortTestPrefabCall(31, [6]uint32{})
				r.Freed = freed == uint64(p) && *words["intro"] == 0 && legacy.PortTestPrefabCall(31, [6]uint32{}) == 0
				if !r.Freed {
					t.Fatal("intro release state")
				}
				rows = append(rows, r)
				cryptfile.Close()
			}
		}
	}
	for _, flags := range []uint32{0, 0x200000, 0x400000, 0x600000} {
		for _, length := range []int{-1, 0, 1, 257} {
			noxflags.ResetGame()
			noxflags.SetGame(noxflags.GameFlag(flags))
			data := bytes.Repeat([]byte{'x'}, max(0, length))
			if length < 0 {
				if err := os.Remove(introPath); err != nil && !os.IsNotExist(err) {
					t.Fatal(err)
				}
			} else if err := os.WriteFile(introPath, data, 0600); err != nil {
				t.Fatal(err)
			}
			if err := cryptfile.OpenGlobal(path, cryptfile.WriteOnly, -1); err != nil {
				t.Fatal(err)
			}
			ret := legacy.PortTestPrefabCall(32, [6]uint32{})
			cryptfile.Close()
			got, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			var want bytes.Buffer
			binary.Write(&want, binary.LittleEndian, uint16(1))
			n := 0
			if flags&0x200000 != 0 && length >= 0 {
				n = length
			}
			binary.Write(&want, binary.LittleEndian, uint32(n))
			want.Write(data[:n])
			if ret != 1 || !bytes.Equal(got, want.Bytes()) {
				t.Fatal("intro write schema")
			}
			rows = append(rows, row{Name: fmt.Sprintf("write/f%x/n%d", flags, length), Return: ret, Wire: got, Freed: *words["intro"] == 0})
		}
	}
	spellbookCapture(t, "prefab-runtime-intro", rows, "2c445d2eb2b9856f886c209f3174c26e3df5825104572c62c118da7166fd6135")
}
