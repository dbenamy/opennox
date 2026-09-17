//go:build porttest

package opennox

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"unsafe"

	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/common/memmap"
)

func TestItemXferOldWeaponAttributeDefaults(t *testing.T) {
	s := newItemXferOwner(t)
	cache := memmap.PtrUint32(0x5D4594, 1564960)
	oldCache := *cache
	t.Cleanup(func() { *cache = oldCache })
	path := filepath.Join(t.TempDir(), "old-attributes.bin")
	for _, version := range []int16{-1, 0, 1, 10, 11} {
		for _, charged := range []bool{false, true} {
			t.Run(fmt.Sprintf("v%d-charged%v", version, charged), func(t *testing.T) {
				u := newItemXferObject(t, s, "Weapon")
				if charged {
					u.ObjClass = object.ClassWand
					u.ObjSubClass = object.SubClass(0x10000)
				}
				*(*uint32)(unsafe.Add(u.InitData, 16)) = 0x12345678
				var stream mapDrawableStream
				stream.u16(uint16(version))
				objectXferEmptyBase(&stream, version)
				if version >= 11 {
					stream.Write(make([]byte, 4))
				}
				if err := os.WriteFile(path, stream.Bytes(), 0600); err != nil {
					t.Fatal(err)
				}
				objectXferReadRecord(t, u, path)
				want := uint32(0x12345678)
				if charged {
					// Modern records without modifiers supply two 0xffff words.
					// Old records must define that same copied attribute tail.
					want = 0xffffffff
				}
				if got := objectXferGetWord(u.InitData, 16); got != want {
					t.Fatalf("attribute tail=%08x want=%08x", got, want)
				}
			})
		}
	}
}
