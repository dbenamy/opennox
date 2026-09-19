//go:build porttest

package opennox

import (
	"bytes"
	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

func TestPlayerFilesInventoryGridCount(t *testing.T) {
	raw, _, restore := legacy.PortTestMeterInventory()
	t.Cleanup(restore)
	cells := legacy.PortTestInventoryCells()
	var rows []map[string]any
	for _, pattern := range []string{"empty", "one", "edges", "all", "byte-max"} {
		clear(raw)
		var want uint32
		for col := 0; col < 4; col++ {
			for row := 0; row < 21; row++ {
				count := byte(0)
				if pattern == "one" && row == 7 && col == 2 {
					count = 5
				}
				if pattern == "edges" && (row == 0 || row == 19 || row == 20) {
					count = byte(col + 1)
				}
				if pattern == "all" {
					count = byte(row + col + 1)
				}
				if pattern == "byte-max" {
					count = 255
				}
				cells[row+21*col].Count = count
				if row < 20 {
					want += uint32(count)
				}
			}
		}
		before := bytes.Clone(raw)
		got := legacy.PortTestPlayerFileCall("sub_41B3B0")
		if got != want || !bytes.Equal(raw, before) {
			t.Fatal("inventory count", pattern, got, want)
		}
		rows = append(rows, map[string]any{"pattern": pattern, "count": got})
	}
	spellbookCapture(t, "player-files-inventory-grid-count", rows, "720f7af495eb06338fd3d3ef49031797e42cb98cde34f5ab76223987f88c4adb")
}

func TestPlayerFilesInventoryFilter(t *testing.T) {
	s := newItemXferOwner(t)
	t.Cleanup(s.PortTestInventoryEnvironment(false, nil, nil))
	glyph := s.Types.IndByID("Glyph")
	if glyph == 0 {
		t.Fatal("missing real glyph type")
	}
	cached := memmap.PtrUint32(0x5D4594, 527724)
	old := *cached
	t.Cleanup(func() { *cached = old })
	u, free := alloc.New(server.Object{})
	t.Cleanup(free)
	var rows []map[string]any
	for _, cold := range []bool{false, true} {
		for _, class := range []object.Class{0, object.ClassPlayer, object.ClassImmobile, object.ClassImmobile | object.ClassPlayer, 0xffffffff} {
			for _, typ := range []uint16{0, uint16(glyph), uint16(glyph) + 1, 65535} {
				*u = server.Object{ObjClass: class, TypeInd: typ}
				*cached = uint32(glyph)
				if cold {
					*cached = 0
				}
				raw := unsafe.Slice((*byte)(unsafe.Pointer(u)), int(unsafe.Sizeof(*u)))
				before := bytes.Clone(raw)
				want := uint32(1)
				if uint32(class)&0x40 != 0 || typ == uint16(glyph) {
					want = 0
				}
				got := legacy.PortTestPlayerFileCall("sub_41B3E0", uint32(uintptr(unsafe.Pointer(u))))
				if got != want || *cached != uint32(glyph) || !bytes.Equal(raw, before) {
					t.Fatal("inventory filter", cold, class, typ, got, want, *cached, glyph)
				}
				rows = append(rows, map[string]any{"cold": cold, "class": uint32(class), "type": typ, "glyph": glyph, "return": got})
			}
		}
	}
	spellbookCapture(t, "player-files-inventory-filter", rows, "0b9a189afb810d9cc23469d5bb4c00b5a2f2191107225f198b10c21d1f48b09e")
}
