//go:build porttest

package opennox

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc/handles"
)

func TestObjectXferPendingReferences(t *testing.T) {
	handles.Init()
	t.Cleanup(handles.Release)
	core := newObjectXferOwner(t)
	path := filepath.Join(t.TempDir(), "references.bin")
	for _, flags := range []uint32{0, 0x200000, 0x400000} {
		for _, inner := range []int16{1, 2, 3, 4, 5, 60, 61, 62, 63, 64} {
			t.Run(fmt.Sprintf("flags%x-inner%d", flags, inner), func(t *testing.T) {
				t.Cleanup(noxflags.PortTestGameFlags(noxflags.GameFlag(flags)))
				pending, closePending := legacy.PortTestObjectXferPendingOwners()
				defer closePending()
				u := newObjectXferSimple(t, core)
				stream := objectXferCommonStream(60, inner, true, 0)
				if err := os.WriteFile(path, stream.Bytes(), 0600); err != nil {
					t.Fatal(err)
				}
				if err := cryptfile.OpenGlobal(path, cryptfile.ReadOnly, -1); err != nil {
					t.Fatal(err)
				}
				defer cryptfile.Close()
				if legacy.Nox_xxx_mapReadWriteObjData_4F4530(u, 60) != 1 {
					t.Fatal("valid reference record rejected")
				}
				rows := pending()
				want := 0
				if flags == 0 && inner >= 2 {
					want = 2
				}
				if len(rows) != want {
					t.Fatalf("pending refs=%v want count=%d", rows, want)
				}
				if want == 2 && (rows[0] != [2]uint32{0x2468ace0, 456} || rows[1] != [2]uint32{0x2468ace0, 123}) {
					t.Fatalf("pending order/IDs=%v", rows)
				}
			})
		}
	}
}

func TestObjectXferGeneratedScriptIDs(t *testing.T) {
	handles.Init()
	t.Cleanup(handles.Release)
	core := newObjectXferOwner(t)
	path := filepath.Join(t.TempDir(), "script-id.bin")
	for _, flags := range []uint32{0, 0x200000, 0x400000} {
		for _, present := range []bool{false, true} {
			t.Run(fmt.Sprintf("flags%x-present%v", flags, present), func(t *testing.T) {
				t.Cleanup(noxflags.PortTestGameFlags(noxflags.GameFlag(flags)))
				u := newObjectXferSimple(t, core)
				before := u.ScriptIDVal
				stream := objectXferCommonStream(60, 64, present, 0)
				binary.LittleEndian.PutUint32(stream.Bytes()[6:], 0)
				// Use the real pending-reference owner for records with optional fields;
				// the companion test checks the list contents independently.
				_, closePending := legacy.PortTestObjectXferPendingOwners()
				defer closePending()
				if err := os.WriteFile(path, stream.Bytes(), 0600); err != nil {
					t.Fatal(err)
				}
				if err := cryptfile.OpenGlobal(path, cryptfile.ReadOnly, -1); err != nil {
					t.Fatal(err)
				}
				defer cryptfile.Close()
				if legacy.Nox_xxx_mapReadWriteObjData_4F4530(u, 60) != 1 {
					t.Fatal("zero-ID record rejected")
				}
				want := 0
				if flags == 0 {
					want = before + 1
				}
				if u.ScriptIDVal != want {
					t.Fatalf("script ID=%d want=%d", u.ScriptIDVal, want)
				}
			})
		}
	}
}
