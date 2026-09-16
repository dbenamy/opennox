//go:build porttest

package opennox

import (
	"github.com/opennox/libs/client/keybind"
	"github.com/opennox/libs/strman"
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"strings"
	"testing"
	"unsafe"
)

type bindingOwner struct {
	*listboxOwner
	words   map[int]*uint32
	columns [2]*gui.Window
	modal   *gui.Window
	menu    bool
}

func newBindingOwner(t *testing.T, menu bool) *bindingOwner {
	o := &bindingOwner{listboxOwner: newListboxOwner(t), menu: menu}
	o.c.ctrl = new(CtrlEventHandler)
	animWord := legacy.PortTestBindingAnimationWord()
	oldAnimWord := *animWord
	*animWord = 0
	t.Cleanup(func() { *animWord = oldAnimWord })
	var restore func()
	o.words, restore = legacy.PortTestBindingWords()
	t.Cleanup(restore)
	// C uses these literal UTF-16 spaces when removing duplicate assignments.
	for _, off := range []int{185444, 185448, 185452, 185456, 187824, 187828, 187832, 187836} {
		b := memmap.BlobByAddr(0x587000).Data[off : off+4]
		old := append([]byte(nil), b...)
		t.Cleanup(func() { copy(b, old) })
		copy(b, []byte{32, 0, 0, 0})
	}
	for i := range o.columns {
		o.create(t, 8, 180, 80, func(_ *gui.WindowData, d *gui.ScrollListBoxData) { d.Count = 16 })
		o.columns[i] = o.win
		o.win = nil
	}
	o.modal = o.c.GUI.NewWindowRaw(nil, 8, 10, 10, 100, 40, nil)
	o.modal.SetID(980)
	o.modal.SetFunc94(func(_ *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
		if ev.EventCode() == 23 {
			return gui.RawEventResp(1)
		}
		return nil
	})
	base := 1321244
	modal := 1321232
	if menu {
		base = 1522624
		modal = 1522612
	}
	for i, w := range o.columns {
		*o.words[base+4*i] = uint32(uintptr(w.C()))
	}
	*o.words[modal] = uint32(uintptr(o.modal.C()))
	var entries []strman.Entry
	for _, k := range keybind.ListKeys() {
		entries = append(entries, strman.Entry{ID: k.TitleID(), Vals: []strman.Variant{{Str: "Key " + k.String()}}})
	}
	configure, restore := o.c.srv.Server.PortTestMeterStrings(entries...)
	t.Cleanup(restore)
	configure(0)
	return o
}
func (o *bindingOwner) selectedWord() *uint32 {
	if o.menu {
		return o.words[1522632]
	}
	return o.words[1321252]
}
func bindingEvent(w *gui.Window, event int, a, b uintptr) int {
	return gui.EventRespInt(w.Func94(&gui.RawEvent{Event: event, Arg1: a, Arg2: b}))
}
func bindingText(w *gui.Window, event int, s string, row int) {
	bindingEvent(w, event, uintptr(unsafe.Pointer(alloc.InternCString16(s))), uintptr(row))
}
func (o *bindingOwner) resetRows(rows [2][]string, column, row int) {
	*o.selectedWord() = 0
	for i, w := range o.columns {
		bindingEvent(w, 0x400f, 0, 0)
		for _, s := range rows[i] {
			bindingText(w, 0x400d, s, -1)
		}
		bindingEvent(w, 16403, ^uintptr(0), 0)
	}
	if column >= 0 {
		*o.selectedWord() = uint32(uintptr(o.columns[column].C()))
		bindingEvent(o.columns[column], 16403, uintptr(row), 0)
	}
}

type bindingResult struct {
	Menu, Mouse         bool
	Key                 uint32
	Column, Row, Result int
	Text                [2][]string
	Selection           [2]int32
	Selected            int
}

func (o *bindingOwner) snapshot(mouse bool, key uint32, column, row, result int) bindingResult {
	r := bindingResult{Menu: o.menu, Mouse: mouse, Key: key, Column: column, Row: row, Result: result, Selected: -1}
	for i, w := range o.columns {
		if *o.selectedWord() == uint32(uintptr(w.C())) {
			r.Selected = i
		}
		d := (*gui.ScrollListBoxData)(w.WidgetData)
		r.Selection[i] = int32(d.Field_12)
		for _, v := range unsafe.Slice(d.Items, int(d.Field_11_0)) {
			r.Text[i] = append(r.Text[i], strings.TrimRight(legacy.GoWStringSlice(v.Text[:]), "\x00"))
		}
	}
	return r
}
