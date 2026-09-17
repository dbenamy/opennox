//go:build porttest

package opennox

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"unsafe"

	"github.com/opennox/libs/object"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/internal/cryptfile"
)

func TestItemXferWeaponCharges(t *testing.T) {
	s := newItemXferOwner(t)
	path := filepath.Join(t.TempDir(), "charges.bin")
	cache := memmap.PtrUint32(0x5D4594, 1564960)
	old := *cache
	t.Cleanup(func() { *cache = old })
	var rows []itemXferCaptureRow
	defer func() {
		spellbookCapture(t, "item-xfer-weapon-charges", rows, "78ce522d53ec6949f2a27ff33e3a3fe7b1cad0a2041537b4db82dc63d39752d9")
	}()
	cases := []struct {
		count, max byte
		amount     uint32
	}{
		{0, 5, 0}, {3, 5, 50}, {5, 5, 100}, {6, 5, 100}, {255, 5, 100},
		{3, 0, 50}, {3, 7, 50}, {3, 255, 50}, {3, 5, 101}, {3, 5, 0x7fffffff}, {3, 5, 0x80000000}, {3, 5, 0xffffffff},
	}
	for _, v := range []int16{40, 41, 60, 61, 62, 63, 64} {
		for _, quest := range []bool{false, true} {
			for _, special := range []bool{false, true} {
				for i, tc := range cases {
					t.Run(fmt.Sprintf("v%d-quest%v-special%v-case%d", v, quest, special, i), func(t *testing.T) {
						flags := noxflags.GameFlag(0x200000)
						if quest {
							flags |= 0x1000
						}
						t.Cleanup(noxflags.PortTestGameFlags(flags))
						u := newItemXferObject(t, s, "Weapon")
						u.ObjClass = object.ClassWand
						u.ObjSubClass = object.SubClass(0x10000)
						if special {
							u.ObjSubClass = object.SubClass(0x4000000)
						}
						*(*byte)(unsafe.Add(u.UseData.Ptr, 108)) = 3
						*(*byte)(unsafe.Add(u.UseData.Ptr, 109)) = 5
						objectXferSetWord(u.UseData.Ptr, 112, 100)
						var p mapDrawableStream
						p.u16(uint16(v))
						objectXferEmptyBase(&p, v)
						p.Write(make([]byte, 4))
						hasCharges := v >= 41 && (!special || v >= 62)
						if hasCharges {
							p.u8(tc.count)
							p.u8(tc.max)
							if v >= 61 {
								p.u32(tc.amount)
							}
						}
						if v >= 42 {
							p.u16(37)
						}
						if v == 63 {
							p.u8(19)
						}
						if v >= 64 {
							p.u32(123)
						}
						if err := os.WriteFile(path, p.Bytes(), 0600); err != nil {
							t.Fatal(err)
						}
						if err := cryptfile.OpenGlobal(path, cryptfile.ReadOnly, -1); err != nil {
							t.Fatal(err)
						}
						defer cryptfile.Close()
						if err := u.CallXfer(nil); err != nil {
							t.Fatal(err)
						}
						count, max, amount := byte(3), byte(5), uint32(100)
						if hasCharges {
							count, max = tc.count, tc.max
							if v >= 61 {
								amount = tc.amount
							}
							if quest && (count > 5 || max != 5 || int32(amount) < 0 || amount > 100) {
								count, max, amount = 0, 5, 0
							}
						}
						if *(*byte)(unsafe.Add(u.UseData.Ptr, 108)) != count || *(*byte)(unsafe.Add(u.UseData.Ptr, 109)) != max || objectXferGetWord(u.UseData.Ptr, 112) != amount {
							t.Fatalf("charges=%d/%d/%08x want=%d/%d/%08x", *(*byte)(unsafe.Add(u.UseData.Ptr, 108)), *(*byte)(unsafe.Add(u.UseData.Ptr, 109)), objectXferGetWord(u.UseData.Ptr, 112), count, max, amount)
						}
						pos, err := cryptfile.Global().File.Seek(0, 1)
						if err != nil || pos != int64(p.Len()) {
							t.Fatalf("position=%d want=%d err=%v", pos, p.Len(), err)
						}
						itemXferCaptureCase(t, &rows, u)
					})
				}
			}
		}
	}
}
