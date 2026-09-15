//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"testing"
	"unicode/utf16"
	"unsafe"

	"github.com/opennox/libs/client/seat"
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/client/input"
	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/legacy"
)

type entrySeat struct {
	callbacks seat.InputConfig
	enabled   bool
	changes   []bool
}

func (s *entrySeat) InputTick() {}
func (s *entrySeat) ReplaceInputs(c seat.InputConfig) seat.InputConfig {
	old := s.callbacks
	s.callbacks = c
	return old
}
func (s *entrySeat) OnInput(f func(seat.InputEvent)) { s.callbacks = append(s.callbacks, f) }
func (s *entrySeat) SetTextInput(v bool)             { s.enabled = v; s.changes = append(s.changes, v) }
func (s *entrySeat) send(ev seat.InputEvent) {
	for _, f := range s.callbacks {
		f(ev)
	}
}

type entryOwner struct {
	*objectRenderOwner
	parent, win *gui.Window
	seat        *entrySeat
	env         *legacy.PortTestEntryEnvironment
	language    func(int)
	notices     [][3]uint32
	named       []string
}

func newEntryOwner(t *testing.T, extraNames ...string) *entryOwner {
	o := &entryOwner{objectRenderOwner: newObjectRenderOwner(t, extraNames...)}
	o.c.GUI = gui.New(o.c.Render())
	o.env = legacy.PortTestNewEntryEnvironment()
	t.Cleanup(o.env.Restore)
	var restore func()
	o.language, restore = o.c.srv.Server.PortTestEntryStrings()
	t.Cleanup(restore)
	oldLoad := legacy.Nox_xxx_gLoadImg
	t.Cleanup(func() { legacy.Nox_xxx_gLoadImg = oldLoad })
	names := map[string]int{"DefaultLBUpButton": 0, "DefaultLBUpButtonLit": 1, "DefaultLBUpButtonDis": 2, "DefaultLBDownButton": 3, "DefaultLBDownButtonLit": 4, "DefaultLBDownButtonDis": 5, "DefaultSliderThumb": 6, "DefaultSliderThumbLit": 7, "DefaultSliderThumbDis": 8}
	legacy.Nox_xxx_gLoadImg = func(name string) *noxrender.Image {
		i, ok := names[name]
		if !ok {
			panic("unowned entry listbox image: " + name)
		}
		o.named = append(o.named, name)
		return o.images[i]
	}
	o.parent = o.c.GUI.NewWindowRaw(nil, 8, 3, 4, 90, 90, func(_ *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
		a, b := ev.EventArgsC()
		av := uint32(a)
		switch ev.EventCode() {
		case 16384, 16389, 16390, 16415:
			if o.win != nil && av == uint32(uintptr(o.win.C())) {
				av = 0xe1000002
			}
		}
		o.notices = append(o.notices, [3]uint32{uint32(ev.EventCode()), av, uint32(b)})
		return nil
	})
	o.parent.SetID(55)
	t.Cleanup(func() { o.destroy(); o.c.GUI.DestroyAll(); o.c.GUI.FreeDestroyed() })
	return o
}
func (o *entryOwner) destroy() {
	if o.win != nil {
		o.win.Destroy()
		o.c.GUI.FreeDestroyed()
		o.c.GUI.FreeDestroyed()
		o.win = nil
	}
	// Restore fixture state even when the original destructor leaves an active entry.
	o.env.ClearActive()
}
func (o *entryOwner) create(t *testing.T, lang int, flags gui.StatusFlags, configure func(*gui.WindowData, *gui.EntryFieldData)) *gui.EntryFieldData {
	t.Helper()
	o.destroy()
	o.language(lang)
	o.env.Context(true)
	o.env.Blink(0)
	o.seat = &entrySeat{}
	o.c.Inp = input.New(o.c.Log, o.seat, false, lang)
	o.c.Inp.OnInputString(func(s string) {
		for _, v := range utf16.Encode([]rune(s)) {
			legacy.NoxInputOnChar(v)
		}
	})
	o.named = nil
	o.notices = nil
	draw := gui.WindowData{Style: gui.StyleEntryField | 0x100, Window: o.parent, BgColorVal: 0x21082108, EnColorVal: 0x03e003e0, HlColorVal: 0x7fff7fff, DisColorVal: 0x42104210, TextColorVal: 0x7fff7fff}
	draw.BgImageHnd = noxrender.ImageHandle(o.images[1].C())
	draw.DisImageHnd = noxrender.ImageHandle(o.images[5].C())
	data := gui.EntryFieldData{Field_1040: 256}
	if configure != nil {
		configure(&draw, &data)
	}
	o.win = legacy.Nox_gui_newEntryField_488500(o.parent, flags, 10, 12, 70, 19, &draw, &data)
	if o.win == nil {
		t.Fatal("entry constructor")
	}
	o.win.SetID(77)
	return (*gui.EntryFieldData)(o.win.WidgetData)
}
func (o *entryOwner) key(key, state uint32) int {
	return gui.EventRespInt(o.win.Func93(&gui.RawEvent{Event: 21, Arg1: uintptr(key), Arg2: uintptr(state)}))
}
func (o *entryOwner) text(v []uint16) {
	// A separately allocated terminated slice preserves UTF-16 code units exactly.
	buf, free := alloc.Make([]uint16{}, len(v)+1)
	defer free()
	copy(buf, v)
	o.win.Func94(&gui.RawEvent{Event: 16414, Arg1: uintptr(unsafe.Pointer(&buf[0]))})
}
