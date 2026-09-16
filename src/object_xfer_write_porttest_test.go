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

	"github.com/opennox/libs/object"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func TestObjectXferCommonWriteContracts(t *testing.T) {
	core := newObjectXferOwner(t)
	path := filepath.Join(t.TempDir(), "written.bin")
	for _, outer := range []int{0, 9, 10, 19, 20, 29, 30, 39, 40, 60, 64} {
		for _, reason := range []string{"default", "name", "flags", "extra", "team"} {
			for _, nameLen := range []int{0, 3, 255, 256, 257} {
				if reason != "name" && nameLen != 0 {
					continue
				}
				t.Run(fmt.Sprintf("outer%d-%s-len%d", outer, reason, nameLen), func(t *testing.T) {
					u := newObjectXferSimple(t, core)
					u.ObjFlags = 0
					u.Field5 = 0
					u.Extent = 0x12345678
					u.ScriptIDVal = 0x2468ace0
					u.PosVec = types.Pointf{X: 64.25, Y: 128.75}
					u.Field34 = 222
					name := strings.Repeat("x", nameLen)
					switch reason {
					case "name":
						p, free := alloc.CString(name)
						u.IDPtr = unsafe.Pointer(p)
						t.Cleanup(func() { u.IDPtr = nil; free() })
					case "flags":
						u.ObjFlags = object.Flags(0x01408162)
					case "extra":
						u.Field5 = 0xde
					case "team":
						u.TeamVal.ID = 7
					}
					modern := outer >= 40
					present := name != "" || reason == "flags" || reason == "extra" || reason == "team"
					var want mapDrawableStream
					want.u16(64)
					want.u32(u.Extent)
					if modern {
						want.u32(uint32(u.ScriptIDVal))
					} else {
						want.u32(uint32(u.ObjFlags) & 0x11408162)
					}
					want.f32(u.PosVec.X)
					want.f32(u.PosVec.Y)
					if modern {
						if present {
							want.u8(255)
						} else {
							want.u8(0)
						}
					}
					if !modern || present {
						if modern {
							want.u32(uint32(u.ObjFlags) & 0x11408162)
						}
						if outer >= 10 {
							want.u8(byte(nameLen))
							want.WriteString(name[:int(byte(nameLen))])
						}
						if outer >= 20 {
							want.u8(byte(u.TeamVal.ID))
						}
						if outer >= 30 {
							want.u8(0)
						}
						if modern {
							want.u16(0) // No owned-object references.
							want.u32(u.Field5 & 0x5e)
							want.u16(1)
							want.u32(0)
							want.u32(0)  // Empty pickup script handler.
							want.u32(99) // Lifetime 222 - frame 123.
						}
					}
					if err := cryptfile.OpenGlobal(path, cryptfile.WriteOnly, -1); err != nil {
						t.Fatal(err)
					}
					result := legacy.Nox_xxx_mapReadWriteObjData_4F4530(u, outer)
					if err := cryptfile.Close(); err != nil {
						t.Fatal(err)
					}
					got, err := os.ReadFile(path)
					if err != nil {
						t.Fatal(err)
					}
					if result != 1 || !bytes.Equal(got, want.Bytes()) {
						t.Fatalf("result=%d bytes=%x want=%x", result, got, want.Bytes())
					}
					wantName := name
					if outer >= 10 {
						wantName = name[:int(byte(nameLen))]
					}
					if u.ID() != wantName || u.Field34 != 222 {
						t.Fatalf("writer mutation name=%q lifetime=%d", u.ID(), u.Field34)
					}
				})
			}
		}
	}
}
