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

	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func TestItemXferModifierRecords(t *testing.T) {
	s := newItemXferOwner(t)
	path := filepath.Join(t.TempDir(), "modifiers.bin")
	cache := memmap.PtrUint32(0x5D4594, 1564960)
	oldCache := *cache
	t.Cleanup(func() { *cache = oldCache })
	names := []string{"PortPower", "PortMaterial", strings.Repeat("m", 255), ""}
	var mods []*server.ModifierEff
	var cname []*byte
	for _, name := range names {
		m, free := alloc.New(server.ModifierEff{})
		t.Cleanup(free)
		p, freeName := alloc.CString(name)
		t.Cleanup(freeName)
		mods = append(mods, m)
		cname = append(cname, p)
	}
	t.Cleanup(s.PortTestControlsModifiers(mods, cname))
	var rows []itemXferCaptureRow
	defer func() {
		spellbookCapture(t, "item-xfer-modifiers", rows, "5dec507c6294b4c666c4f14a4a43a21b5a75e97471e7d97b8b0d05d256f35b83")
	}()
	for _, name := range []string{"Weapon", "Armor", "Ammo", "Team"} {
		for mask := 0; mask < 16; mask++ {
			t.Run(fmt.Sprintf("%s-mask%x", name, mask), func(t *testing.T) {
				u := newItemXferObject(t, s, name)
				objectXferSetCommon(u)
				ptrs := (*[4]*server.ModifierEff)(u.InitData)
				for i, m := range mods {
					if mask&(1<<i) != 0 {
						ptrs[i] = m
					}
				}
				version := uint16(60)
				if name == "Weapon" {
					version = 64
				} else if name == "Armor" {
					version = 62
				}
				want := objectXferCurrentRecord(version)
				for i, n := range names {
					if mask&(1<<i) == 0 {
						n = ""
					}
					want.u8(byte(len(n)))
					want.WriteString(n)
				}
				if name == "Weapon" || name == "Armor" {
					want.u16(80)
					want.u32(0)
				} else if name == "Ammo" {
					want.u8(0)
					want.u8(0)
				}
				if err := cryptfile.OpenGlobal(path, cryptfile.WriteOnly, -1); err != nil {
					t.Fatal(err)
				}
				if err := u.CallXfer(nil); err != nil {
					t.Fatal(err)
				}
				written := cryptfile.Global().PortTestChecksum()
				if err := cryptfile.Close(); err != nil {
					t.Fatal(err)
				}
				data, err := os.ReadFile(path)
				if err != nil {
					t.Fatal(err)
				}
				if !bytes.Equal(data, want.Bytes()) {
					t.Fatalf("modifier bytes=%x want=%x", data, want.Bytes())
				}
				v := newItemXferObject(t, s, name)
				objectXferSetWord(v.InitData, 16, 0x12345678)
				checksum := itemXferReadRecord(t, v, path)
				got := (*[4]*server.ModifierEff)(v.InitData)
				for i, m := range mods {
					var expected *server.ModifierEff
					if mask&(1<<i) != 0 && names[i] != "" {
						expected = m
					}
					if got[i] != expected {
						t.Fatalf("modifier slot%d identity", i)
					}
				}
				tail := uint32(0x12345678)
				if mask&7 != 0 {
					tail = 0xffffffff
				}
				if objectXferGetWord(v.InitData, 16) != tail {
					t.Fatal("modifier attribute tail")
				}
				if name == "Team" && (objectXferGetWord(v.UpdateData, 0) != objectXferGetWord(unsafe.Pointer(v), 56) || objectXferGetWord(v.UpdateData, 4) != objectXferGetWord(unsafe.Pointer(v), 60)) {
					t.Fatal("team marker position")
				}
				itemXferCaptureCase(t, &rows, v, checksum, written)
			})
		}
	}
}
