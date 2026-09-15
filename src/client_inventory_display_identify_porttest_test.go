//go:build porttest

package opennox

import (
	"image"
	"strings"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/client/gui"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func (o *inventoryDisplayOwner) identifyWindows(t *testing.T) (*gui.Window, *gui.Window, *gui.Window) {
	t.Helper()
	label := o.c.GUI.NewStaticText(o.parent, 9151, 0, 0, 200, 20, false, false, "")
	icon := o.c.GUI.NewWindowRaw(o.parent, 8, 5, 5, 30, 30, nil)
	icon.SetID(9155)
	draw := gui.WindowData{Style: gui.StyleScrollListBox | 0x100, Window: o.parent, TextColorVal: 0x7fff7fff}
	data := gui.ScrollListBoxData{Count: 64, Line_height: 13}
	list := legacy.Nox_gui_newScrollListBox_4A4310(o.parent, 8, 10, 30, 195, 180, &draw, &data)
	if list == nil {
		t.Fatal("identify list constructor")
	}
	list.SetID(9156)
	return label, icon, list
}

func TestClientInventoryDisplayIdentification(t *testing.T) {
	o := newInventoryDisplayOwner(t)
	pos, free := alloc.Make([]int32{}, 2)
	defer free()
	var rows []inventoryDisplayResult
	for _, name := range []string{"RedApple", "Bow", "Quiver"} {
		for _, detailed := range []bool{false, true} {
			for _, health := range [][2]uint16{{0, 0}, {0, 100}, {24, 100}, {25, 100}, {49, 100}, {50, 100}, {74, 100}, {75, 100}, {99, 100}, {100, 100}} {
				o.reset(t)
				label, _, list := o.identifyWindows(t)
				o.c.Mouse = image.Pt(11, 16)
				if detailed {
					noxflags.SetGame(2048)
				}
				dr := o.item(t, name, 123)
				*(*uint16)(unsafe.Add(dr.C(), 292)), *(*uint16)(unsafe.Add(dr.C(), 294)) = health[0], health[1]
				*(*byte)(unsafe.Add(dr.C(), 298)) = 7
				*o.displayWords["dword_5d4594_1063116"] = uint32(uintptr(dr.C()))
				r := o.call(t, len(rows), 3, txptr(unsafe.Pointer(&pos[0])), 0, 0)
				d := (*gui.ScrollListBoxData)(list.WidgetData)
				for _, item := range unsafe.Slice(d.Items, int(d.Count)) {
					if item.Text[0] != 0 {
						r.WidgetText = append(r.WidgetText, alloc.GoString16(&item.Text[0]))
					}
				}
				r.WidgetText = append(r.WidgetText, alloc.GoString16((*gui.StaticTextData)(label.WidgetData).Text))
				if !strings.HasPrefix(alloc.GoString16((*gui.StaticTextData)(label.WidgetData).Text), "Identify ") {
					t.Fatalf("identify label %q", alloc.GoString16((*gui.StaticTextData)(label.WidgetData).Text))
				}
				if len(r.WidgetText) < 3 {
					t.Fatalf("missing item details %v", r.WidgetText)
				}
				if *o.displayWords["dword_5d4594_1063120"] != uint32(uintptr(dr.C())) {
					t.Fatal("description cache did not select item")
				}
				rows = append(rows, r)
				// Cached redraw leaves the actual list and label unchanged.
				before := append([]gui.ScrollListBoxItem(nil), unsafe.Slice(d.Items, int(d.Count))...)
				r = o.call(t, len(rows), 3, txptr(unsafe.Pointer(&pos[0])), 0, 0)
				after := unsafe.Slice(d.Items, int(d.Count))
				for i := range before {
					if before[i] != after[i] {
						t.Fatal("cached identify rewrote rows")
					}
				}
				rows = append(rows, r)
				*o.displayWords["dword_5d4594_1063116"] = 0
				r = o.call(t, len(rows), 3, txptr(unsafe.Pointer(&pos[0])), 0, 0)
				if *o.displayWords["dword_5d4594_1063120"] != 0 {
					t.Fatal("empty selection retained cache")
				}
				rows = append(rows, r)
			}
		}
	}
	inventoryDisplayCapture(t, "identify", rows, "1e16c910126c894d88e21b5f0a39c65154f7dd142fb674bdcb48f93627c787b7")
}
