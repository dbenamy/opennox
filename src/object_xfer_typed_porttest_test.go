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

	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/internal/cryptfile"
)

var objectXferKinds = []struct {
	name    string
	version uint16
	payload int
}{
	{"SpellPagePedestal", 60, 4}, {"Readable", 60, 5}, {"Exit", 60, 13},
	{"Door", 60, 12}, {"Trigger", 61, 64}, {"Hole", 60, 28},
	{"Transporter", 60, 4}, {"Elevator", 61, 9}, {"ElevatorShaft", 60, 4},
	{"Mover", 60, 29}, {"Glyph", 60, 10}, {"Sentry", 61, 12},
}

func TestObjectXferTypedVersionRejection(t *testing.T) {
	core := newObjectXferOwner(t)
	path := filepath.Join(t.TempDir(), "version.bin")
	for _, sp := range objectXferKinds {
		for _, v := range []uint16{sp.version + 1, 32767} {
			t.Run(fmt.Sprintf("%s-%d", sp.name, v), func(t *testing.T) {
				u := newObjectXferTyped(t, core, sp.name)
				var stream mapDrawableStream
				stream.u16(v)
				stream.u32(0x12345678)
				if err := os.WriteFile(path, stream.Bytes(), 0600); err != nil {
					t.Fatal(err)
				}
				if err := cryptfile.OpenGlobal(path, cryptfile.ReadOnly, -1); err != nil {
					t.Fatal(err)
				}
				defer cryptfile.Close()
				raw := unsafe.Slice((*byte)(unsafe.Pointer(u)), int(unsafe.Sizeof(*u)))
				before := append([]byte(nil), raw...)
				if err := u.CallXfer(nil); err == nil {
					t.Fatal("future version accepted")
				}
				pos, err := cryptfile.Global().File.Seek(0, io.SeekCurrent)
				if err != nil || pos != 2 || !bytes.Equal(before, raw) {
					t.Fatalf("rejection position=%d err=%v object changed=%v", pos, err, !bytes.Equal(before, raw))
				}
			})
		}
	}
}

func TestObjectXferTypedDefaultRecords(t *testing.T) {
	var captures []objectXferCaptureRow
	defer func() { spellbookCapture(t, "object-xfer-typed-default", captures, "d7afffaf6c7195147aa27a1f4ea8ed0c792cb5f27de697480693611d4e0082b3") }()
	core := newObjectXferOwner(t)
	mapDrawableTables(t)
	path := filepath.Join(t.TempDir(), "typed.bin")
	for _, sp := range objectXferKinds {
		t.Run(sp.name, func(t *testing.T) {
			u := newObjectXferTyped(t, core, sp.name)
			u.ObjFlags = 0
			u.Extent = 0x12345678
			u.ScriptIDVal = 88
			u.PosVec = types.Pointf{X: 64.25, Y: 128.75}
			u.Field34 = 777
			var want mapDrawableStream
			want.u16(sp.version)
			want.u16(64)
			want.u32(u.Extent)
			want.u32(88)
			want.f32(64.25)
			want.f32(128.75)
			want.u8(0)
			payload := make([]byte, sp.payload)
			switch sp.name {
			case "Readable", "Exit":
				binary.LittleEndian.PutUint32(payload, 1)
			case "Trigger":
				binary.LittleEndian.PutUint32(payload, 20)
				binary.LittleEndian.PutUint32(payload[4:], 30)
				for _, off := range []int{18, 28, 38} {
					binary.LittleEndian.PutUint16(payload[off:], 1)
				}
			case "Hole":
				binary.LittleEndian.PutUint16(payload[4:], 1)
			}
			want.Write(payload)
			if err := cryptfile.OpenGlobal(path, cryptfile.WriteOnly, -1); err != nil {
				t.Fatal(err)
			}
			if err := u.CallXfer(nil); err != nil {
				t.Fatal(err)
			}
			if err := cryptfile.Close(); err != nil {
				t.Fatal(err)
			}
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(data, want.Bytes()) {
				t.Fatalf("written=%x want=%x", data, want.Bytes())
			}
			if u.Field34 != 777 {
				t.Fatal("writer changed saved lifetime")
			}
			v := newObjectXferTyped(t, core, sp.name)
			defer objectXferCaptureCase(t, &captures, v)
			v.Field34 = 999
			if err := cryptfile.OpenGlobal(path, cryptfile.ReadOnly, -1); err != nil {
				t.Fatal(err)
			}
			defer cryptfile.Close()
			if err := v.CallXfer(nil); err != nil {
				t.Fatal(err)
			}
			pos, err := cryptfile.Global().File.Seek(0, io.SeekCurrent)
			if err != nil || pos != int64(len(data)) {
				t.Fatalf("read position=%d len=%d err=%v", pos, len(data), err)
			}
			if v.Extent != u.Extent || v.ScriptIDVal != 88 || v.PosVec != u.PosVec || v.NewPos != u.PosVec || v.Field34 != 999 {
				t.Fatal("typed common field round trip")
			}
			if sp.name == "Trigger" && (v.Shape.Box.W != 20 || v.Shape.Box.H != 30) {
				t.Fatal("trigger callback pointer/shape round trip")
			}
		})
	}
}
