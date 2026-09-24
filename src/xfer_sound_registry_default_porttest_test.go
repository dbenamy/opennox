//go:build porttest

package opennox

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"
	"unsafe"

	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

// These are the same real type constructors used by the existing serialized
// record goldens; prove that their slots select every nondefault registration.
func TestXferSoundRegistryFixtureIdentities(t *testing.T) {
	s := newCreatureXferOwner(t)
	count := 0
	check := func(u *server.Object, name string) {
		t.Helper()
		if u.Xfer != legacy.PortTestXferRegistryPointer(name+"Xfer") {
			t.Fatal("fixture bypasses registration", name)
		}
		count++
	}
	for _, sp := range objectXferKinds {
		check(newObjectXferTyped(t, s, sp.name), sp.name)
	}
	check(newObjectXferTyped(t, s, "InvisibleLight"), "InvisibleLight")
	for _, sp := range itemXferKinds {
		check(newItemXferObject(t, s, sp.name), sp.name)
	}
	for _, name := range []string{"Monster", "NPC"} {
		check(newCreatureXferObject(t, s, name), name)
	}
	if count != 27 {
		t.Fatal("nondefault fixture count", count)
	}
}

func TestXferSoundRegistryDefaultRecords(t *testing.T) {
	s := newObjectXferOwner(t)
	mapDrawableTables(t)
	path := filepath.Join(t.TempDir(), "default.bin")
	for _, outer := range []bool{false, true} {
		t.Run(fmt.Sprintf("outer-%t", outer), func(t *testing.T) {
			call := func(u *server.Object) error {
				if outer {
					return asObjectS(u).CallXfer(nil)
				}
				return u.CallXfer(nil)
			}
			u := newObjectXferSimple(t, s)
			u.Xfer = legacy.PortTestXferRegistryPointer("DefaultXfer")
			u.ObjFlags = 0
			u.Extent = 0x12345678
			u.ScriptIDVal = 88
			u.PosVec = types.Pointf{X: 64.25, Y: 128.75}
			u.Field34 = 777
			var want mapDrawableStream
			want.u16(60)
			want.u16(64)
			want.u32(u.Extent)
			want.u32(88)
			want.f32(64.25)
			want.f32(128.75)
			want.u8(0)
			if err := cryptfile.OpenGlobal(path, cryptfile.WriteOnly, -1); err != nil {
				t.Fatal(err)
			}
			if err := call(u); err != nil {
				t.Fatal(err)
			}
			if err := cryptfile.Close(); err != nil {
				t.Fatal(err)
			}
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(data, want.Bytes()) || u.Field34 != 777 {
				t.Fatalf("default wire/lifetime: %x lifetime=%d", data, u.Field34)
			}
			v := newObjectXferSimple(t, s)
			v.Xfer = u.Xfer
			v.Field34 = 999
			if err := cryptfile.OpenGlobal(path, cryptfile.ReadOnly, -1); err != nil {
				t.Fatal(err)
			}
			defer cryptfile.Close()
			if err := call(v); err != nil {
				t.Fatal(err)
			}
			pos, err := cryptfile.Global().File.Seek(0, io.SeekCurrent)
			if err != nil || pos != int64(len(data)) || v.Extent != u.Extent || v.ScriptIDVal != 88 || v.PosVec != u.PosVec || v.NewPos != u.PosVec || v.Field34 != 999 {
				t.Fatal("default read state/position", pos, err)
			}
		})
	}
}

func TestXferSoundRegistryDefaultOuterError(t *testing.T) {
	s := newObjectXferOwner(t)
	u := newObjectXferSimple(t, s)
	u.Xfer = legacy.PortTestXferRegistryPointer("DefaultXfer")
	path := filepath.Join(t.TempDir(), "future.bin")
	old := legacy.Nox_xxx_XFerDefault4F49A0
	defer func() { legacy.Nox_xxx_XFerDefault4F49A0 = old }()
	calls := 0
	legacy.Nox_xxx_XFerDefault4F49A0 = func(*cryptfile.CryptFile, *server.Object, unsafe.Pointer) error { calls++; return nil }
	for _, version := range []uint16{61, 32767} {
		var stream mapDrawableStream
		stream.u16(version)
		stream.u32(0x12345678)
		if err := os.WriteFile(path, stream.Bytes(), 0600); err != nil {
			t.Fatal(err)
		}
		if err := cryptfile.OpenGlobal(path, cryptfile.ReadOnly, -1); err != nil {
			t.Fatal(err)
		}
		raw := unsafe.Slice((*byte)(u.CObj()), int(unsafe.Sizeof(*u)))
		before := append([]byte(nil), raw...)
		err := asObjectS(u).CallXfer(nil)
		want := fmt.Sprintf("default xfer: unexpected value 1: %d", version)
		if err == nil || err.Error() != want || calls != 0 {
			t.Fatal("outer default must keep local implementation and detailed error", err, calls)
		}
		pos, seekErr := cryptfile.Global().File.Seek(0, io.SeekCurrent)
		if seekErr != nil || pos != 2 || !bytes.Equal(raw, before) {
			t.Fatal("future-version rejection state", pos, seekErr)
		}
		if err := cryptfile.Close(); err != nil {
			t.Fatal(err)
		}
	}
}
