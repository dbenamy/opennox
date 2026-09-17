//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func TestObjectXferNamesAndPickupHandlers(t *testing.T) {
	core := newObjectXferOwner(t)
	path := filepath.Join(t.TempDir(), "names.bin")
	var captures []objectXferCaptureRow
	defer func() { spellbookCapture(t, "object-xfer-names", captures, "128626fc068730414326b95ea2478602099ed42e92690ad3fafb0f97961809ce") }()
	for _, nameLen := range []int{0, 3, 255} {
		for _, scriptLen := range []int{0, 1, 31, 127} {
			t.Run(fmt.Sprintf("name%d-handler%d", nameLen, scriptLen), func(t *testing.T) {
				u := newObjectXferSimple(t, core)
				objectXferSetCommon(u)
				u.TeamVal.ID = 7
				u.Field34 = 222
				u.ScriptPickup.Flags = 0x87654321
				name, handler := strings.Repeat("n", nameLen), strings.Repeat("s", scriptLen)
				ptr, free := alloc.CString(name)
				u.IDPtr = unsafe.Pointer(ptr)
				t.Cleanup(func() { u.IDPtr = nil; free() })
				copy(unsafe.Slice((*byte)(u.Field189), 2572), handler)
				var want mapDrawableStream
				want.u16(64)
				want.u32(u.Extent)
				want.u32(88)
				want.f32(64.25)
				want.f32(128.75)
				want.u8(255)
				want.u32(0)
				want.u8(byte(nameLen))
				want.WriteString(name)
				want.u8(7)
				want.u8(0)
				want.u16(0)
				want.u32(0)
				want.u16(1)
				want.u32(uint32(scriptLen))
				want.WriteString(handler)
				want.u32(u.ScriptPickup.Flags)
				want.u32(99)
				if err := cryptfile.OpenGlobal(path, cryptfile.WriteOnly, -1); err != nil {
					t.Fatal(err)
				}
				if legacy.Nox_xxx_mapReadWriteObjData_4F4530(u, 60) != 1 {
					t.Fatal("named writer failed")
				}
				if err := cryptfile.Close(); err != nil {
					t.Fatal(err)
				}
				got, err := os.ReadFile(path)
				if err != nil {
					t.Fatal(err)
				}
				if !bytes.Equal(got, want.Bytes()) {
					t.Fatalf("named bytes=%x want=%x", got, want.Bytes())
				}
				v := newObjectXferSimple(t, core)
				defer objectXferCaptureCase(t, &captures, v)
				if err := cryptfile.OpenGlobal(path, cryptfile.ReadOnly, -1); err != nil {
					t.Fatal(err)
				}
				defer cryptfile.Close()
				if legacy.Nox_xxx_mapReadWriteObjData_4F4530(v, 60) != 1 {
					t.Fatal("named reader failed")
				}
				if v.ID() != name || alloc.GoString((*byte)(v.Field189)) != handler || v.ScriptPickup.Flags != 0x87654321 || v.ScriptPickup.Func != -1 {
					t.Fatal("name/pickup-handler round trip")
				}
			})
		}
	}
}

func TestObjectXferFuturePickupHandlerVersion(t *testing.T) {
	core := newObjectXferOwner(t)
	u := newObjectXferSimple(t, core)
	stream := objectXferCommonStream(60, 64, true, 0)
	binary.LittleEndian.PutUint16(stream.Bytes()[stream.Len()-14:], 2)
	path := filepath.Join(t.TempDir(), "handler-version.bin")
	if err := os.WriteFile(path, stream.Bytes(), 0600); err != nil {
		t.Fatal(err)
	}
	if err := cryptfile.OpenGlobal(path, cryptfile.ReadOnly, -1); err != nil {
		t.Fatal(err)
	}
	defer cryptfile.Close()
	if legacy.Nox_xxx_mapReadWriteObjData_4F4530(u, 60) != 0 {
		t.Fatal("future handler version accepted")
	}
	pos, err := cryptfile.Global().File.Seek(0, io.SeekCurrent)
	if err != nil || pos != int64(stream.Len()-12) {
		t.Fatalf("rejection stream position=%d want=%d err=%v", pos, stream.Len()-12, err)
	}
}
