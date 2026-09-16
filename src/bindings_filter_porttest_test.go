//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestBindingEditorWheelFilter(t *testing.T) {
	type record struct {
		Menu               bool
		Event, Key, Result int
		Scroll             uint16
	}
	var records []record
	for _, menu := range []bool{false, true} {
		t.Run(fmt.Sprintf("menu=%v", menu), func(t *testing.T) {
			o := newBindingOwner(t, menu)
			op := "sub_4C3A60"
			if menu {
				op = "sub_4CC140"
			}
			for _, event := range []int{19, 20, 21, 23} {
				for _, key := range []int{0, 15, 30} {
					w := o.columns[0]
					d := (*gui.ScrollListBoxData)(w.WidgetData)
					d.Field_13_1 = 7
					result := legacy.PortTestBindingInvoke(op, [4]uint32{uint32(uintptr(w.C())), uint32(event), uint32(key), 1})
					if event == 19 || event == 20 {
						if result != 0 || d.Field_13_1 != 7 {
							t.Fatalf("wheel event %d was not filtered", event)
						}
					}
					records = append(records, record{menu, event, key, int(result), d.Field_13_1})
				}
			}
		})
	}
	spellbookCapture(t, "binding-wheel-filter", records, "8cd46a4e235f8a962d94fad92080f19d2ae76021dab3d52fcbfb56d40e333cf3")
}
