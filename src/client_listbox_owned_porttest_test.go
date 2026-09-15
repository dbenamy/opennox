//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"testing"
	"unsafe"
)

func TestClientListboxOwnedDataContract(t *testing.T) {
	o := newListboxOwner(t)
	for multi := 0; multi < 2; multi++ {
		for scroll := 0; scroll < 2; scroll++ {
			var input *gui.ScrollListBoxData
			o.create(t, 8, 70, 65, func(draw *gui.WindowData, d *gui.ScrollListBoxData) {
				input = d
				d.Field_4 = uint32(multi)
				d.Field_3 = uint32(scroll)
				draw.Window = nil // The constructor must establish its default owner itself.
			})
			d := o.data()
			if unsafe.Pointer(d) == unsafe.Pointer(input) || *d != *input {
				t.Fatal("constructor must own a copy of the updated caller options")
			}
			if o.win.DrawData().Window != o.win {
				t.Fatal("missing default notification owner")
			}
			input.Count = 99
			if d.Count != 8 {
				t.Fatal("widget data aliases caller options")
			}
			pointers := []unsafe.Pointer{o.win.WidgetData, unsafe.Pointer(d.Items)}
			if multi != 0 {
				pointers = append(pointers, unsafe.Pointer(uintptr(d.Field_12)))
			}
			if scroll != 0 {
				pointers = append(pointers, (*gui.Window)(d.Field_9).WidgetData)
			}
			for _, p := range pointers {
				if !alloc.PortTestAllocationLive(p) {
					t.Fatal("missing owned allocation")
				}
			}
			o.win.Destroy()
			for _, p := range pointers {
				if !alloc.PortTestAllocationLive(p) {
					t.Fatal("allocation freed before deferred cleanup")
				}
			}
			o.c.GUI.FreeDestroyed()
			o.c.GUI.FreeDestroyed()
			o.win = nil
			for _, p := range pointers {
				if alloc.PortTestAllocationLive(p) {
					t.Fatal("allocation survived cleanup")
				}
			}
		}
	}
}
