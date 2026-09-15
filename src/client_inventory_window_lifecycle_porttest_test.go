//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
	"unsafe"
)

func TestClientInventoryWindowReset(t *testing.T) {
	o := newInventoryWindowOwner(t)
	var rows []inventoryWindowResult
	for _, full := range []bool{false, true} {
		for _, identify := range []bool{false, true} {
			o.reset(t)
			o.construct(t)
			o.openInput()
			cells := legacy.PortTestInventoryCells()
			for col := 0; col < 4; col++ {
				for row := 0; row < 21; row++ {
					cell := &cells[col*21+row]
					if cell.Drawable == nil && (full || row == col) {
						o.stack(t, col, row, 1, "RedApple", uint32(100+col*21+row))
					}
					cell.Equipped, cell.Alternate = 1, 1
				}
			}
			if identify {
				o.invoke(12, 0, 0, 0, 0)
			}
			for _, name := range []string{"dword_5d4594_1062512", "dword_5d4594_1062516", "dword_5d4594_1062520"} {
				*o.windowWords[name] = 125
			}
			var want []uint32
			for row := 0; row < 21; row++ {
				for col := 0; col < 4; col++ {
					if dr := cells[col*21+row].Drawable; dr != nil {
						want = append(want, o.norm(uint32(txptr(dr.C()))))
					}
				}
			}
			before := len(o.events)
			rows = append(rows, o.capture(t, len(rows), 0, 27, 0, 0, 0, 0))
			var got []uint32
			for _, event := range o.events[before:] {
				if event[0] == 2 {
					got = append(got, event[2])
				}
			}
			if len(got) != len(want) {
				t.Fatalf("reset deleted%d want%d", len(got), len(want))
			}
			for i := range got {
				if got[i] != want[i] {
					t.Fatalf("reset deletion%d got%#x want%#x", i, got[i], want[i])
				}
			}
			for i, cell := range cells {
				if cell.Drawable != nil || cell.Count != 0 || cell.Equipped != 0 || cell.Alternate != 0 {
					t.Fatalf("reset retained cell%d", i)
				}
			}
			for _, v := range o.equipment {
				if v != 0 {
					t.Fatal("reset retained equipment")
				}
			}
			if *memmap.PtrUint8(0x5D4594, 1049868) != 0 || int32(*o.windowWords["dword_587000_136184"]) != -225 {
				t.Fatal("reset did not close panel")
			}
			before = len(o.events)
			rows = append(rows, o.capture(t, len(rows), 1, 27, 0, 0, 0, 0))
			for _, event := range o.events[before:] {
				if event[0] == 2 {
					t.Fatal("second reset deleted again")
				}
			}
		}
	}
	inventoryWindowCapture(t, "reset", rows, "5cd609c0411c6bf819b974681e0568f4cdd57613b782bfd61b896d52cd03eca0")
}

func TestClientInventoryWindowAlternateEvents(t *testing.T) {
	o := newInventoryWindowOwner(t)
	var rows []inventoryWindowResult
	for _, mode := range []uint32{0, 5, 6} {
		for _, event := range []uintptr{0, 5, 6, 7, 8, 9, 19} {
			for _, hasAlt := range []bool{false, true} {
				o.reset(t)
				o.construct(t)
				o.openInput()
				*o.windowWords["dword_5d4594_1049864"] = mode
				if hasAlt {
					o.stack(t, 0, 0, 1, "Bow", 100)
					cell := &legacy.PortTestInventoryCells()[0]
					cell.Alternate = 1
					*o.windowWords["dword_5d4594_1062480"] = uint32(txptr(unsafe.Pointer(cell)))
				}
				rows = append(rows, o.capture(t, len(rows), 0, 3, 0, 0, 0, 0))
				r := o.capture(t, len(rows), 1, 6, txptr(o.mainWindow().C()), event, inventoryWindowPoint(100, 100), 0)
				rows = append(rows, r)
				want := uint32(0)
				if mode == 6 || (event >= 5 && event <= 8) {
					want = 1
				}
				if r.Return != want {
					t.Fatalf("alternate mode%d event%d return%d want%d", mode, event, r.Return, want)
				}
			}
		}
	}
	inventoryWindowCapture(t, "alternate-events", rows, "7fda5462b286a5241afb3fd23c34344645e7cd750e439a5a1a64a11bf1d516f0")
}

func TestClientInventoryWindowAlternateDrag(t *testing.T) {
	o := newInventoryWindowOwner(t)
	var rows []inventoryWindowResult
	for _, name := range []string{"RedApple", "Bow", "Quiver"} {
		for _, outside := range []bool{false, true} {
			o.reset(t)
			o.construct(t)
			o.openInput()
			o.stack(t, 0, 0, 2, name, 100)
			rows = append(rows, o.capture(t, len(rows), 0, 1, txptr(o.mainWindow().C()), 5, inventoryWindowPoint(339, 38), 0))
			pos := inventoryWindowPoint(100, 100)
			if outside {
				pos = inventoryWindowPoint(600, 400)
			}
			rows = append(rows, o.capture(t, len(rows), 1, 6, txptr(o.mainWindow().C()), 6, pos, 0))
			if *memmap.PtrUint32(0x5D4594, 1049848) != 0 || o.c.dragndropItem != nil {
				t.Fatal("alternate release retained drag")
			}
		}
	}
	inventoryWindowCapture(t, "alternate-drag", rows, "d53380bcb3ad2a267a2d24571874e37476df01bd46d2750b65fa8723b465f307")
}
