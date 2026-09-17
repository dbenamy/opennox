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
	"github.com/opennox/opennox/v1/common/unit/ai"
	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func TestCreatureXferActionArguments(t *testing.T) {
	s := newCreatureXferOwner(t)
	path := filepath.Join(t.TempDir(), "arguments.bin")
	var rows []struct {
		Case              string
		Writer, Reader    [6]uint32
		Wire              []byte
		ReadCRC, WriteCRC uint32
	}
	defer func() {
		spellbookCapture(t, "creature-xfer-action-arguments", rows, "f42e6c0ca02e3a05d6290f05d502f23c850e7d3208df22d1def7b8f70d1ae331")
	}()
	for id := 0; id < 72; id++ {
		for _, present := range []bool{false, true} {
			t.Run(fmt.Sprintf("id%d-present%v", id, present), func(t *testing.T) {
				u := newCreatureXferObject(t, s, "Monster")
				u.ScriptIDVal = 0x12345678
				wp, free := alloc.New(server.Waypoint{})
				t.Cleanup(free)
				wp.Index = 0x87654321
				entry := [6]uint32{uint32(id), 0xfffffffd, 0x12345678, 0x87654321, 0xfedcba98, 0xaabbccdd}
				n := int(*memmap.PtrUint32(0x587000, uintptr(255604+16*id)))
				if n > 2 {
					t.Fatal("shipped action exceeds two argument slots")
				}
				want := new(mapDrawableStream)
				itemXferRewardName(want, ai.ActionType(id).String())
				want.u8(byte(n))
				normalized := entry
				expected := [6]uint32{uint32(id), 0x51515151, 0x51515151, 0x51515151, 0x51515151, entry[5]}
				for j := 0; j < n; j++ {
					kind := *memmap.PtrUint32(0x587000, uintptr(255608+16*id+4*j))
					word := 1 + 2*j
					switch kind {
					case 0:
						want.u32(entry[word])
						want.u32(entry[word+1])
						expected[word], expected[word+1] = entry[word], entry[word+1]
					case 1, 2:
						entry[word] = 0
						value := uint32(0)
						if present {
							if kind == 1 {
								entry[word] = uint32(uintptr(u.CObj()))
								value = uint32(u.ScriptIDVal)
							} else {
								entry[word] = uint32(uintptr(unsafe.Pointer(wp)))
								value = wp.Index
							}
						}
						normalized[word] = 0
						if present {
							normalized[word] = 1
						}
						want.u32(value)
						expected[word] = value
					case 3, 4, 6:
						want.u32(entry[word])
						expected[word] = entry[word]
					case 5:
						want.u32(entry[word])
						expected[word] = creatureXferClampTimer(entry[word], 5)
					case 7:
						want.u8(byte(entry[word]))
						expected[word] = expected[word]&0xffffff00 | uint32(byte(entry[word]))
					default:
						t.Fatalf("unexpected shipped kind %d", kind)
					}
				}
				want.u32(entry[5])
				before := entry
				if err := cryptfile.OpenGlobal(path, cryptfile.WriteOnly, -1); err != nil {
					t.Fatal(err)
				}
				if ret := legacy.PortTestCreatureXferHelper(1, u, unsafe.Pointer(&entry[0]), 5); ret != 1 {
					t.Fatalf("writer result=%d", ret)
				}
				writeCRC := cryptfile.Global().PortTestChecksum()
				cryptfile.Close()
				got, err := os.ReadFile(path)
				if err != nil || !bytes.Equal(got, want.Bytes()) || entry != before {
					t.Fatal("action argument writer bytes/state")
				}
				loaded := [6]uint32{0, 0x51515151, 0x51515151, 0x51515151, 0x51515151, 0}
				if err := cryptfile.OpenGlobal(path, cryptfile.ReadOnly, -1); err != nil {
					t.Fatal(err)
				}
				defer cryptfile.Close()
				if ret := legacy.PortTestCreatureXferHelper(1, u, unsafe.Pointer(&loaded[0]), 5); ret != 1 {
					t.Fatalf("reader result=%d", ret)
				}
				pos, err := cryptfile.Global().File.Seek(0, 1)
				if err != nil || pos != int64(want.Len()) || loaded != expected {
					t.Fatalf("argument state=%08x want=%08x", loaded, expected)
				}
				rows = append(rows, struct {
					Case              string
					Writer, Reader    [6]uint32
					Wire              []byte
					ReadCRC, WriteCRC uint32
				}{t.Name(), normalized, loaded, got, cryptfile.Global().PortTestChecksum(), writeCRC})
			})
		}
	}
}
