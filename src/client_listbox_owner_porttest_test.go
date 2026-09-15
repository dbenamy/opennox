//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/client/input"
	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"testing"
	"unsafe"
)

type listboxNotice struct {
	Code          int
	A, B          uint32
	Count, Insert uint16
	Selection     []int32
}
type listboxOwner struct {
	*entryOwner
	notices []listboxNotice
}

func newListboxOwner(t *testing.T) *listboxOwner {
	o := &listboxOwner{entryOwner: newEntryOwner(t)}
	// Reuse the real language/input setup and then release its temporary entry.
	// All listbox cases below use the actual listbox constructor and owned rows.
	o.entryOwner.create(t, 0, 8, nil)
	o.entryOwner.destroy()
	installListboxPalette(t)
	o.parent.SetFunc94(func(_ *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
		a, b := ev.EventArgsC()
		n := listboxNotice{Code: ev.EventCode(), A: uint32(a), B: uint32(b)}
		if o.win != nil && o.win.WidgetData != nil {
			d := o.data()
			n.Count = d.Field_11_0
			n.Insert = d.Field_11_1
			n.Selection = o.selection()
			switch ev.EventCode() {
			case 16384, 16400, 16401:
				n.A = o.norm(n.A)
				n.B = o.norm(n.B)
			}
		}
		o.notices = append(o.notices, n)
		return nil
	})
	return o
}
func (o *listboxOwner) data() *gui.ScrollListBoxData {
	return (*gui.ScrollListBoxData)(o.win.WidgetData)
}
func (o *listboxOwner) selection() []int32 {
	d := o.data()
	word := *(*uint32)(unsafe.Add(o.win.WidgetData, 48))
	if d.Field_4 == 0 {
		return []int32{int32(word)}
	}
	return append([]int32(nil), unsafe.Slice((*int32)(unsafe.Pointer(uintptr(word))), int(d.Count)+1)...)
}
func (o *listboxOwner) norm(v uint32) uint32 {
	if v == 0 {
		return 0
	}
	if v == uint32(uintptr(o.parent.C())) {
		return 0xe1000001
	}
	if o.win == nil {
		return v
	}
	windows := []*gui.Window{o.win}
	var add func(*gui.Window)
	add = func(w *gui.Window) {
		for c := w.Field100Ptr; c != nil; c = c.Prev() {
			windows = append(windows, c)
			add(c)
		}
	}
	add(o.win)
	for i, w := range windows {
		if v == uint32(uintptr(w.C())) {
			return 0xe1000002 + uint32(i)
		}
		if w.WidgetData != nil && v == uint32(uintptr(w.WidgetData)) {
			return 0xe2000001 + uint32(i)
		}
	}
	if d := o.data(); d != nil {
		if v == uint32(uintptr(unsafe.Pointer(d.Items))) {
			return 0xe3000001
		}
		if d.Field_4 != 0 && v == *(*uint32)(unsafe.Add(o.win.WidgetData, 48)) {
			return 0xe4000001
		}
		for i := 0; i < int(d.Count); i++ {
			if v == uint32(uintptr(unsafe.Add(unsafe.Pointer(d.Items), i*524+4))) {
				return 0xe5000001 + uint32(i)
			}
		}
	}
	if n, ok := o.c.imageRefs[v]; ok {
		return n
	}
	return v
}
func (o *listboxOwner) create(t *testing.T, flags gui.StatusFlags, width, height int, configure func(*gui.WindowData, *gui.ScrollListBoxData)) {
	t.Helper()
	o.entryOwner.destroy()
	o.seat = &entrySeat{}
	o.c.Inp = input.New(o.c.Log, o.seat, false, 0)
	o.named = nil
	o.notices = nil
	draw := gui.WindowData{Style: gui.StyleScrollListBox | 0x100, Window: o.parent, BgColorVal: 0x21082108, EnColorVal: 0x03e003e0, HlColorVal: 0x7fff7fff, DisColorVal: 0x42104210, SelColorVal: 0x7c007c00, TextColorVal: 0x7fff7fff}
	draw.BgImageHnd = noxrender.ImageHandle(o.images[1].C())
	draw.EnImageHnd = noxrender.ImageHandle(o.images[3].C())
	draw.DisImageHnd = noxrender.ImageHandle(o.images[5].C())
	draw.HlImageHnd = noxrender.ImageHandle(o.images[7].C())
	draw.SelImageHnd = noxrender.ImageHandle(o.images[9].C())
	data := gui.ScrollListBoxData{Count: 8, Line_height: 10}
	if configure != nil {
		configure(&draw, &data)
	}
	o.win = legacy.Nox_gui_newScrollListBox_4A4310(o.parent, flags, 10, 12, width, height, &draw, &data)
	if o.win == nil {
		t.Fatal("listbox constructor")
	}
	o.win.SetID(77)
}
func (o *listboxOwner) event(code int, a, b uint32) int {
	return gui.EventRespInt(o.win.Func94(&gui.RawEvent{Event: code, Arg1: uintptr(a), Arg2: uintptr(b)}))
}
func (o *listboxOwner) textEvent(code int, text []uint16, arg int32) int {
	if text == nil {
		return o.event(code, 0, uint32(arg))
	}
	p, free := alloc.Make([]uint16{}, len(text)+1)
	defer free()
	copy(p, text)
	return o.event(code, uint32(uintptr(unsafe.Pointer(&p[0]))), uint32(arg))
}
