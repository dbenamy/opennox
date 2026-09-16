//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
	"unsafe"
)

func TestBindingEditorRouting(t *testing.T) {
	type record struct {
		Menu                       bool
		Event, Column, Row, Result int
		Selected                   int
		Hidden, Focused, Stacked   bool
		Prompt                     string
		Scroll                     [4]uint16
	}
	var records []record
	for _, menu := range []bool{false, true} {
		t.Run(fmt.Sprintf("menu=%v", menu), func(t *testing.T) {
			o := newBindingOwner(t, menu)
			boxes := [4]*gui.Window{nil, nil, o.columns[0], o.columns[1]}
			for i := 0; i < 2; i++ {
				o.create(t, 8, 180, 35, func(_ *gui.WindowData, d *gui.ScrollListBoxData) { d.Count = 16; d.Field_3 = 1 })
				boxes[i] = o.win
				o.win = nil
			}
			rootOff, base, buffer := 1321228, 1321236, 1321256
			op := "sub_4C3CD0"
			if menu {
				rootOff, base, buffer = 1522604, 1522616, 1522636
				op = "sub_4CBF60"
			}
			*o.words[rootOff] = uint32(uintptr(o.parent.C()))
			for i, w := range boxes {
				w.SetID(uint(910 + i))
				*o.words[base+4*i] = uint32(uintptr(w.C()))
			}
			sd := (*gui.ScrollListBoxData)(boxes[0].WidgetData)
			(*gui.Window)(sd.Field_7).SetID(921)
			(*gui.Window)(sd.Field_8).SetID(922)
			for _, w := range boxes[1:] {
				bindingEvent(w, 16408, uintptr(sd.Field_7), 0)
				bindingEvent(w, 16409, uintptr(sd.Field_8), 0)
			}
			buf := memmap.BlobByAddr(0x5D4594).Data[buffer : buffer+512]
			old := append([]byte(nil), buf...)
			t.Cleanup(func() { copy(buf, old) })
			for _, event := range []int{23, 16384, 16391, 16393, 16400} {
				for col := 0; col < 2; col++ {
					for _, row := range []int{-1, 0, 3} {
						rows := [2][]string{{"one", "two", "three", "four"}, {"five", "six", "seven", "eight"}}
						o.resetRows(rows, col, row)
						*o.selectedWord() = 0
						for _, w := range boxes[:2] {
							bindingEvent(w, 0x400f, 0, 0)
							for i := 0; i < 4; i++ {
								bindingText(w, 0x400d, fmt.Sprintf("Action %d", i), -1)
							}
						}
						o.modal.StackPop()
						o.c.GUI.Focus(nil)
						o.modal.SetHidden(true)
						clear(buf)
						a := uint32(uintptr(o.columns[col].C()))
						b := uint32(0)
						switch event {
						case 16384, 16391:
							if col == 0 {
								a = uint32(uintptr(sd.Field_7))
							} else {
								a = uint32(uintptr(sd.Field_8))
							}
						case 16393:
							a = uint32(uintptr(sd.Field_9))
							if row > 0 {
								b = uint32(row)
							}
						}
						result := legacy.PortTestBindingInvoke(op, [4]uint32{uint32(uintptr(boxes[0].C())), uint32(event), a, b})
						r := record{Menu: menu, Event: event, Column: col, Row: row, Result: int(result), Selected: o.snapshot(false, 0, col, row, 0).Selected, Hidden: o.modal.Flags.IsHidden(), Focused: o.c.GUI.Focused() == o.modal, Stacked: o.c.GUI.StackHead() == o.modal, Prompt: legacy.GoWStringP(unsafe.Pointer(&buf[0]))}
						for i, w := range boxes {
							r.Scroll[i] = (*gui.ScrollListBoxData)(w.WidgetData).Field_13_1
						}
						if event == 16400 {
							opened := row >= 0
							if r.Hidden == opened || r.Focused != opened || r.Stacked != opened {
								t.Fatalf("capture routing %+v", r)
							}
							if opened && r.Selected != col {
								t.Fatalf("selected column %+v", r)
							}
						}
						if event == 23 && r.Result != 1 {
							t.Fatalf("focus response %+v", r)
						}
						records = append(records, r)
					}
				}
			}
		})
	}
	spellbookCapture(t, "binding-routing", records, "22a667044908725ad803326ae28269fda5f4797c3f1fa8d06b6ea174d937285b")
}
