//go:build porttest

package opennox

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"unsafe"

	"github.com/opennox/libs/spell"
	"github.com/opennox/opennox/v1/internal/cryptfile"
)

func TestCreatureXferMerchantRecords(t *testing.T) {
	s := newCreatureXferOwner(t)
	path := filepath.Join(t.TempDir(), "merchant.bin")
	var rows []itemXferCaptureRow
	defer func() {
		spellbookCapture(t, "creature-xfer-merchant", rows, "d1820f99d8a6a53c84f9313b792a1dcf73eb4cc65c83ceecf2cc3a6bb743fb7e")
	}()
	for _, version := range []int{43, 46, 47, 48, 49, 50, 60, 61, 64} {
		for _, count := range []int{0, 1, 2, 60} {
			t.Run(fmt.Sprintf("v%d-count%d", version, count), func(t *testing.T) {
				u := newCreatureXferObject(t, s, "Monster")
				u.ObjSubClass |= 8
				base := new(mapDrawableStream)
				base.u16(uint16(version))
				objectXferEmptyBase(base, int16(version))
				off := creatureXferMonsterHistorical(base, version)
				shop := new(mapDrawableStream)
				if version >= 50 {
					shop.u32(0x12345678)
				}
				if version >= 61 {
					shop.u32(0x87654321)
				}
				if version >= 48 {
					itemXferRewardName(shop, "PortVendor")
				}
				shop.u8(byte(count))
				for i := 0; i < count; i++ {
					if version < 50 {
						shop.u32(0)
					}
					shop.u8(byte(i + 1))
					itemXferRewardName(shop, "portgold")
					if version >= 47 {
						itemXferRewardName(shop, spell.ID(1).String())
					}
					shop.Write([]byte{0, 0, 0, 0})
				}
				p := new(mapDrawableStream)
				p.Write(base.Bytes()[:off])
				p.Write(shop.Bytes())
				p.Write(base.Bytes()[off:])
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
				pos, err := cryptfile.Global().File.Seek(0, 1)
				if err != nil || pos != int64(p.Len()) || *(*byte)(u.InitData) != byte(count) {
					t.Fatal("merchant count/position")
				}
				for i := 0; i < count; i++ {
					entry := unsafe.Add(u.InitData, 4+28*i)
					if objectXferGetWord(entry, 0) != uint32(s.Types.IndByID("portgold")) || *(*byte)(unsafe.Add(entry, 4)) != byte(i+1) {
						t.Fatalf("stock%d type/quantity", i)
					}
					want := uint32(0)
					if version >= 47 {
						want = 1
					}
					if objectXferGetWord(entry, 8) != want {
						t.Fatalf("stock%d parameter", i)
					}
				}
				if version >= 50 && objectXferGetWord(u.InitData, 1716) != 0x12345678 {
					t.Fatal("vendor first tail")
				}
				if version >= 61 && objectXferGetWord(u.InitData, 1720) != 0x87654321 {
					t.Fatal("vendor second tail")
				}
				itemXferCaptureCase(t, &rows, u, cryptfile.Global().PortTestChecksum(), 0)
			})
		}
	}
}
