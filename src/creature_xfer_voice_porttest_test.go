//go:build porttest

package opennox

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func TestCreatureXferVoiceRecords(t *testing.T) {
	s := newCreatureXferOwner(t)
	path := filepath.Join(t.TempDir(), "voices.bin")
	type row struct {
		Case              string
		Wire              []byte
		Present           bool
		ReadCRC, WriteCRC uint32
	}
	var rows []row
	defer func() { spellbookCapture(t, "creature-xfer-voices", rows, "") }()
	for _, name := range []string{"", "PortVoice", strings.Repeat("v", 255)} {
		for _, known := range []bool{false, true} {
			t.Run(fmt.Sprintf("len%d-known%v", len(name), known), func(t *testing.T) {
				u := newCreatureXferObject(t, s, "NPC")
				set, free := alloc.New([20]uint32{})
				defer free()
				cstr, freeName := alloc.CString(name)
				defer freeName()
				set[0] = uint32(uintptr(unsafe.Pointer(cstr)))
				head, freeHead := alloc.New([20]uint32{})
				defer freeHead()
				other, freeOther := alloc.CString("OtherVoice")
				defer freeOther()
				head[0] = uint32(uintptr(unsafe.Pointer(other)))
				if known {
					head[19] = uint32(uintptr(unsafe.Pointer(set)))
				}
				restore := legacy.PortTestCreatureXferVoiceSets(unsafe.Pointer(head))
				defer restore()
				objectXferSetWord(u.UpdateData, 488, uint32(uintptr(unsafe.Pointer(set))))
				defer objectXferSetWord(u.UpdateData, 488, 0)
				p := new(mapDrawableStream)
				itemXferRewardName(p, name)
				if err := cryptfile.OpenGlobal(path, cryptfile.WriteOnly, -1); err != nil {
					t.Fatal(err)
				}
				if ret := legacy.PortTestCreatureXferHelper(3, u, nil, 0); ret != 1 {
					t.Fatalf("writer=%d", ret)
				}
				writeCRC := cryptfile.Global().PortTestChecksum()
				cryptfile.Close()
				got, err := os.ReadFile(path)
				if err != nil || !bytes.Equal(got, p.Bytes()) {
					t.Fatal("voice wire")
				}
				v := newCreatureXferObject(t, s, "NPC")
				// An unknown name clears a previously installed voice set.
				objectXferSetWord(v.UpdateData, 488, uint32(uintptr(unsafe.Pointer(head))))
				defer objectXferSetWord(v.UpdateData, 488, 0)
				if err := cryptfile.OpenGlobal(path, cryptfile.ReadOnly, -1); err != nil {
					t.Fatal(err)
				}
				defer cryptfile.Close()
				ret := legacy.PortTestCreatureXferHelper(3, v, nil, 0)
				want := uint32(0)
				if known {
					want = uint32(uintptr(unsafe.Pointer(set)))
				}
				pos, err := cryptfile.Global().File.Seek(0, 1)
				if ret != 1 || objectXferGetWord(v.UpdateData, 488) != want || err != nil || pos != int64(p.Len()) {
					t.Fatal("voice identity/status/position")
				}
				rows = append(rows, row{t.Name(), got, known, cryptfile.Global().PortTestChecksum(), writeCRC})
			})
		}
	}
	t.Run("nil-set", func(t *testing.T) {
		u := newCreatureXferObject(t, s, "NPC")
		if err := cryptfile.OpenGlobal(path, cryptfile.WriteOnly, -1); err != nil {
			t.Fatal(err)
		}
		defer cryptfile.Close()
		if ret := legacy.PortTestCreatureXferHelper(3, u, nil, 0); ret != 1 {
			t.Fatalf("writer=%d", ret)
		}
		cryptfile.Close()
		got, err := os.ReadFile(path)
		if err != nil || !bytes.Equal(got, []byte{0}) {
			t.Fatal("nil voice record")
		}
	})
}
