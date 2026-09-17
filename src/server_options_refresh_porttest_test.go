//go:build porttest

package opennox

import (
	"encoding/binary"
	"fmt"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
	"unsafe"
)

func TestServerOptionsRefresh(t *testing.T) {
	type row struct {
		Mode               uint16
		Open, Host         bool
		Title, Score, Time string
		SpellMask          []byte
		Weapons, Armor     uint32
	}
	var rows []row
	for _, tc := range []struct {
		mode  uint16
		title string
	}{{0x100, "Arena battle"}, {0x20, "Capture flag"}, {0x400, "Last survivor"}, {0x1000, "Quest adventure"}} {
		for _, open := range []bool{false, true} {
			for _, host := range []bool{false, true} {
				t.Run(fmt.Sprintf("%x-%t-%t", tc.mode, open, host), func(t *testing.T) {
					catalog := legacy.PortTestMapCatalogOpen(9)
					t.Cleanup(catalog.Close)
					o := newServerOptionsOwner(t)
					flags := noxflags.GameFlag(0)
					if host {
						flags = 1
					}
					defer noxflags.PortTestGameFlags(flags)()
					legacy.PortTestServerOptionsModeName(0)
					st := unsafe.Slice(memmap.PtrUint8(0x5D4594, 371438), 58)
					copy(st, "Arena")
					copy(st[9:], "Refresh")
					for i := 24; i < 52; i++ {
						st[i] = byte(5*i + 3)
					}
					binary.LittleEndian.PutUint16(st[52:], tc.mode)
					binary.LittleEndian.PutUint16(st[54:], 91)
					st[56] = 17
					if !open {
						*o.optionWords["root"] = 0
					}
					o.call("refresh", 0, "")
					mask := append([]byte(nil), unsafe.Slice(memmap.PtrUint8(0x5D4594, 1045488), 20)...)
					if string(mask) != string(st[24:44]) || *memmap.PtrUint32(0x5D4594, 1045452) != binary.LittleEndian.Uint32(st[44:]) || *memmap.PtrUint32(0x5D4594, 1045456) != binary.LittleEndian.Uint32(st[48:]) {
						t.Fatal("refresh lost selection masks")
					}
					got := row{tc.mode, open, host, o.options.ChildByID(10119).DrawData().Text(), o.entry(10134), o.entry(10135), mask, *memmap.PtrUint32(0x5D4594, 1045452), *memmap.PtrUint32(0x5D4594, 1045456)}
					if open && (got.Title != tc.title || got.Score != "91" || got.Time != "17") {
						t.Fatalf("refresh labels %+v", got)
					}
					rows = append(rows, got)
				})
			}
		}
	}
	spellbookCapture(t, "server-options-refresh", rows, "e5890fb5155289c0f90cf18f967429e2b526096b206725df72e7ab4041a8415d")
}

func TestServerOptionsTryHide(t *testing.T) {
	o := newServerOptionsOwner(t)
	if o.call("try-close", 0, "") != 1 || *o.optionWords["root"] != 0 {
		t.Fatal("open try-close")
	}
	if o.call("try-close", 0, "") != 0 {
		t.Fatal("closed try-close")
	}
}
