//go:build porttest

package opennox

import (
	"image"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
)

func (o *inventoryDisplayOwner) displayImages() {
	for i := 0; i < 23; i++ {
		*memmap.PtrUint32(0x5D4594, 1049908+uintptr(4*i)) = uint32(uintptr(o.images[i].C()))
	}
	o.parent.DrawData().BgImageHnd = noxrender.ImageHandle(o.images[3].C())
	o.parent.DrawData().HlImageHnd = noxrender.ImageHandle(o.images[7].C())
}

func TestClientInventoryDisplayButtons(t *testing.T) {
	o := newInventoryDisplayOwner(t)
	var rows []inventoryDisplayResult
	for _, op := range []int{0, 9, 10, 13} {
		for _, mode := range []uint32{0, 5, 6} {
			for _, present := range []bool{false, true} {
				o.reset(t)
				o.displayImages()
				*o.displayWords["dword_5d4594_1049864"] = mode
				child := o.c.GUI.NewWindowRaw(o.parent, gui.StatusFlags(8), 3, 5, 60, 70, nil)
				if child == nil {
					t.Fatal("button child")
				}
				child.DrawData().BgImageHnd = noxrender.ImageHandle(o.images[8].C())
				child.DrawData().HlImageHnd = noxrender.ImageHandle(o.images[11].C())
				child.DrawData().ImgPtVal = image.Pt(2, 3)
				if present {
					dr := o.item(t, "Bow", 123)
					cells := legacy.PortTestInventoryCells()
					cells[21].Drawable, cells[21].Count, cells[21].Codes[0] = dr, 1, 123
					*o.displayWords["dword_5d4594_1062480"] = o.cell(21)
					*o.displayWords["dword_5d4594_1063116"] = uint32(uintptr(dr.C()))
					*memmap.PtrUint32(0x5D4594, 1049848) = uint32(uintptr(dr.C()))
					o.equipment[8] = uint32(uintptr(dr.C()))
				}
				for _, flag := range []byte{0, 1, 2} {
					*memmap.PtrUint8(0x5D4594, 1049868) = flag
					r := o.call(t, len(rows), op, txptr(child.C()), txptr(child.DrawData().C()), 0)
					if r.Return != 1 {
						t.Fatalf("draw op%d return%d", op, r.Return)
					}
					rows = append(rows, r)
				}
			}
		}
	}
	inventoryDisplayCapture(t, "buttons", rows, "5a14b6a0384a02db75ae1232334eeae53c204cb99d5447f7cd1080698f7d518b")
}

func TestClientInventoryDisplayTray(t *testing.T) {
	o := newInventoryDisplayOwner(t)
	var rows []inventoryDisplayResult
	for _, scroll := range []uint32{0, 1, 49, 50, 51, 799, 800, 849, 850} {
		for _, mode := range []uint32{0, 5, 6} {
			for _, health := range [][2]uint16{{0, 0}, {0, 100}, {24, 100}, {25, 100}, {49, 100}, {50, 100}, {99, 100}, {100, 100}} {
				o.reset(t)
				o.displayImages()
				*o.displayWords["dword_5d4594_1062512"] = scroll
				*o.displayWords["dword_5d4594_1049864"] = mode
				*o.displayWords["dword_5d4594_1062552"] = 12345
				for col := 0; col < 4; col++ {
					dr := o.item(t, "Bow", uint32(100+col))
					*(*uint16)(unsafe.Add(dr.C(), 292)), *(*uint16)(unsafe.Add(dr.C(), 294)) = health[0], health[1]
					*(*uint16)(unsafe.Add(dr.C(), 448)) = []uint16{0, 1, 32767, 65535}[col]
					*(*uint16)(unsafe.Add(dr.C(), 450)) = 2
					// Include first/last visible rows and the hidden row in every case.
					for _, row := range []int{0, 1, 16, 17, 18, 19, 20} {
						cell := &legacy.PortTestInventoryCells()[col*21+row]
						cell.Drawable, cell.Count, cell.Codes[0] = dr, byte(1+col*10), uint32(100+col)
						if col == 0 {
							cell.Equipped = 1
						}
						if col == 1 {
							cell.Alternate = 1
						}
					}
				}
				*o.displayWords["dword_5d4594_1062480"] = o.cell(21)
				*o.displayWords["dword_5d4594_1062488"] = 102
				r := o.call(t, len(rows), 8, 0, 0, 0)
				if r.Return != 150 {
					t.Fatalf("tray bottom return%d", r.Return)
				}
				if len(r.Text) == 0 || r.Text[0].Text != "12345" {
					t.Fatal("tray currency label")
				}
				rows = append(rows, r)
			}
		}
	}
	inventoryDisplayCapture(t, "tray", rows, "6a654df45a9433eb1423308dbdefd1ebd739b3c950409874b0ef84ecb4056770")
}
