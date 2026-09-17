//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"unsafe"

	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/legacy"
)

func TestCreatureXferActionReferences(t *testing.T) {
	s := newCreatureXferOwner(t)
	t.Cleanup(noxflags.PortTestGameFlags(1))
	path := filepath.Join(t.TempDir(), "refs.bin")
	cache := memmap.PtrUint32(0x5d4594, 2487688)
	oldCache := *cache
	defer func() { *cache = oldCache }()
	type row struct {
		Case              string
		Wire              []byte
		Values            []uint32
		ReadCRC, WriteCRC uint32
	}
	var rows []row
	defer func() {
		spellbookCapture(t, "creature-xfer-action-references", rows, "bef27611f488889453daac9b85fd6f793cc968da9cd95121aa2f06b61a8c9d91")
	}()
	for _, count := range []int{0, 1, 16} {
		for _, netCount := range []int{0, 1, 8} {
			for _, mode := range []string{"absent", "live", "destroyed"} {
				t.Run(fmt.Sprintf("count%d-net%d-%s", count, netCount, mode), func(t *testing.T) {
					restore := legacy.PortTestCreatureXferLookupOwner()
					defer restore()
					u := newCreatureXferObject(t, s, "Monster")
					creatureXferSetTimers(u, 1)
					target := newCreatureXferObject(t, s, "Monster")
					target.NetCode = 222
					target.ScriptIDVal = 0x12345678
					old := s.Objs.List
					if mode != "absent" {
						s.Objs.List = target
					} else {
						s.Objs.List = nil
					}
					defer func() { s.Objs.List = old }()
					if mode == "destroyed" {
						target.ObjFlags |= 0x20
					}
					ptr := uint32(uintptr(target.CObj()))
					for _, off := range []int{1196, 1216} {
						objectXferSetWord(u.UpdateData, off, ptr)
					}
					for _, off := range []int{392, 1200} {
						objectXferSetWord(u.UpdateData, off, target.NetCode)
					}
					*(*byte)(unsafe.Add(u.UpdateData, 1129)) = byte(count)
					*(*byte)(unsafe.Add(u.UpdateData, 2172)) = byte(netCount)
					for i := 0; i < count; i++ {
						objectXferSetWord(u.UpdateData, 1132+4*i, ptr)
					}
					for i := 0; i < netCount; i++ {
						objectXferSetWord(u.UpdateData, 2140+4*i, target.NetCode)
					}
					base := creatureXferEmptyActionRecord(u, 4, s.Frame()).Bytes()
					if len(base) != 195 {
						t.Fatal("unexpected base layout")
					}
					sid := uint32(target.ScriptIDVal)
					nid := uint32(0)
					if mode == "live" {
						nid = sid
					}
					for _, off := range []int{106, 190} {
						binary.LittleEndian.PutUint32(base[off:], sid)
					}
					for _, off := range []int{137, 178} {
						binary.LittleEndian.PutUint32(base[off:], nid)
					}
					p := new(mapDrawableStream)
					p.Write(base[:106])
					for i := 0; i < count; i++ {
						p.u32(sid)
					}
					p.Write(base[106:])
					for i := 0; i < netCount; i++ {
						p.u32(nid)
					}
					if err := cryptfile.OpenGlobal(path, cryptfile.WriteOnly, -1); err != nil {
						t.Fatal(err)
					}
					if ret := legacy.PortTestCreatureXferHelper(0, u, nil, 0); ret != 1 {
						t.Fatalf("write=%d", ret)
					}
					writeCRC := cryptfile.Global().PortTestChecksum()
					cryptfile.Close()
					got, err := os.ReadFile(path)
					if err != nil || !bytes.Equal(got, p.Bytes()) {
						t.Fatal("reference writer bytes")
					}
					v := newCreatureXferObject(t, s, "Monster")
					if err := cryptfile.OpenGlobal(path, cryptfile.ReadOnly, -1); err != nil {
						t.Fatal(err)
					}
					defer cryptfile.Close()
					if ret := legacy.PortTestCreatureXferHelper(0, v, nil, 0); ret != 1 {
						t.Fatalf("read=%d", ret)
					}
					pos, err := cryptfile.Global().File.Seek(0, 1)
					if err != nil || pos != int64(p.Len()) {
						t.Fatal("reference position")
					}
					var values []uint32
					for _, group := range []struct {
						offset, count int
						want          uint32
					}{{1196, 1, sid}, {1216, 1, sid}, {392, 1, nid}, {1200, 1, nid}, {1132, count, sid}, {2140, netCount, nid}} {
						for i := 0; i < group.count; i++ {
							value := objectXferGetWord(v.UpdateData, group.offset+4*i)
							if value != group.want {
								t.Fatalf("reference offset%d=%x want=%x", group.offset+4*i, value, group.want)
							}
							values = append(values, value)
						}
					}
					rows = append(rows, row{t.Name(), got, values, cryptfile.Global().PortTestChecksum(), writeCRC})
				})
			}
		}
	}
}
