//go:build porttest

package opennox

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"unsafe"

	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/unit/ai"
	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/legacy"
)

func TestCreatureXferActionStack(t *testing.T) {
	s := newCreatureXferOwner(t)
	t.Cleanup(noxflags.PortTestGameFlags(1))
	path := filepath.Join(t.TempDir(), "stack.bin")
	type row struct {
		Case              string
		Wire              []byte
		Entries           [][6]uint32
		ReadCRC, WriteCRC uint32
	}
	var rows []row
	defer func() {
		spellbookCapture(t, "creature-xfer-action-stack", rows, "4649d58d02b79be244555bef4a4e51bf1b45add93f1ed7472eff6c27b71ade2c")
	}()
	for length := 0; length <= 24; length++ {
		t.Run(fmt.Sprintf("length%d", length), func(t *testing.T) {
			u := newCreatureXferObject(t, s, "Monster")
			creatureXferSetTimers(u, 1)
			*(*byte)(unsafe.Add(u.UpdateData, 544)) = byte(length - 1)
			entries := unsafe.Slice((*[6]uint32)(unsafe.Add(u.UpdateData, 552)), length)
			stack := new(mapDrawableStream)
			for i := range entries {
				// Alternate idle and wait actions using shipped timestamp metadata.
				id := i % 2
				count := *memmap.PtrUint32(0x587000, uintptr(255604+16*id))
				if count != 1 || *memmap.PtrUint32(0x587000, uintptr(255608+16*id)) != 5 {
					t.Fatal("unexpected idle/wait metadata")
				}
				entries[i] = [6]uint32{uint32(id), uint32(i + 7), 0, 0, 0, uint32(100 + i)}
				itemXferRewardName(stack, ai.ActionType(id).String())
				stack.u8(byte(count))
				stack.u32(entries[i][1])
				stack.u32(entries[i][5])
			}
			base := creatureXferEmptyActionRecord(u, 4, s.Frame()).Bytes()
			p := new(mapDrawableStream)
			p.Write(base[:105])
			p.Write(stack.Bytes())
			p.Write(base[105:])
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
				t.Fatalf("stack writer wire lengths=%d/%d", len(got), p.Len())
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
			if err != nil || pos != int64(p.Len()) || *(*byte)(unsafe.Add(v.UpdateData, 544)) != byte(length-1) {
				t.Fatal("stack position/length")
			}
			loaded := unsafe.Slice((*[6]uint32)(unsafe.Add(v.UpdateData, 552)), length)
			for i, e := range loaded {
				want := entries[i]

				if e != want {
					t.Fatalf("entry%d=%x want=%x", i, e, want)
				}
			}
			rows = append(rows, row{t.Name(), got, append([][6]uint32(nil), loaded...), cryptfile.Global().PortTestChecksum(), writeCRC})
		})
	}
}
