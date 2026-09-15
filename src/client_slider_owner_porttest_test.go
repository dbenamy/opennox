//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"image"
	"testing"
	"unsafe"
)

type sliderNotice struct {
	Code int
	A, B uint32
}
type sliderOwner struct {
	*objectRenderOwner
	input              *gui.SliderData
	draw               *gui.WindowData
	parent, win, thumb *gui.Window
	notices            []sliderNotice
}

func newSliderOwner(t *testing.T) *sliderOwner {
	o := &sliderOwner{objectRenderOwner: newObjectRenderOwner(t)}
	o.c.GUI = gui.New(o.c.Render())
	var free func()
	o.input, free = alloc.New(gui.SliderData{})
	t.Cleanup(free)
	o.draw, free = gui.NewWindowData()
	t.Cleanup(free)
	o.parent = o.c.GUI.NewWindowRaw(nil, 8, 3, 4, 90, 90, func(_ *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
		a, b := ev.EventArgsC()
		av := uint32(a)
		// Normalize only notification fields whose actual payload is a window.
		switch ev.EventCode() {
		case 16384, 16389, 16390, 16393, 16396:
			av = o.norm(av)
		}
		o.notices = append(o.notices, sliderNotice{ev.EventCode(), av, uint32(b)})
		return nil
	})
	o.parent.SetID(55)
	t.Cleanup(o.c.GUI.DestroyAll)
	return o
}
func (o *sliderOwner) norm(v uint32) uint32 {
	if v == 0 {
		return 0
	}
	for i, w := range []*gui.Window{o.parent, o.win, o.thumb} {
		if w != nil && v == uint32(uintptr(w.C())) {
			return 0xe1000001 + uint32(i)
		}
	}
	if v == uint32(uintptr(o.input.CWidgetData())) {
		return 0xe2000001
	}
	if o.win != nil && v == uint32(uintptr(o.win.WidgetData)) {
		return 0xe2000002
	}
	if n, ok := o.c.imageRefs[v]; ok {
		return n
	}
	return v
}
func (o *sliderOwner) rawWindow(w *gui.Window) []uint32 {
	if w == nil {
		return nil
	}
	out := append([]uint32(nil), unsafe.Slice((*uint32)(w.C()), 101)...)
	for _, i := range []int{8, 13, 15, 17, 19, 21, 23, 59, 93, 94, 95, 96, 97, 98, 99, 100} {
		out[i] = o.norm(out[i])
	}
	return out
}
func (o *sliderOwner) create(t *testing.T, id, horizontal, mode, parent, owner int, flags gui.StatusFlags, min, max uint32, size image.Point, configure ...func(*gui.WindowData)) {
	t.Helper()
	if o.win != nil {
		o.win.Destroy()
		o.c.GUI.FreeDestroyed()
		o.win = nil
		o.thumb = nil
	}
	o.resetRender(uint32(id), 120)
	o.notices = nil
	*o.input = gui.SliderData{Min: min, Max: max, Field3: min}
	style := gui.StyleVertSlider
	if horizontal != 0 {
		style = gui.StyleHorizSlider
	}
	*o.draw = gui.WindowData{Style: style, Field0: 0, BgColorVal: 0x21082108, EnColorVal: 0x03e003e0, HlColorVal: 0x7fff7fff, DisColorVal: 0x42104210, SelColorVal: 0x7c007c00, TextColorVal: 0x7fff7fff}
	images := []*noxrender.Image{o.images[1], o.images[3], o.images[5], o.images[7], o.images[9]}
	o.draw.BgImageHnd = noxrender.ImageHandle(images[0].C())
	o.draw.EnImageHnd = noxrender.ImageHandle(images[1].C())
	o.draw.DisImageHnd = noxrender.ImageHandle(images[2].C())
	o.draw.SelImageHnd = noxrender.ImageHandle(images[3].C())
	o.draw.HlImageHnd = noxrender.ImageHandle(images[4].C())
	for _, f := range configure {
		f(o.draw)
	}
	if owner != 0 {
		o.draw.Window = o.parent
	}
	var p *gui.Window
	if parent != 0 {
		p = o.parent
	}
	o.win = legacy.Nox_gui_newSlider_4B4EE0(p, flags|gui.StatusFlags(mode*128), 10, 12, size.X, size.Y, o.draw, o.input)
	if o.win == nil {
		t.Fatal("valid slider was not created")
	}
	o.win.SetID(77)
	o.thumb = o.win.Field100Ptr
	if o.thumb == nil || o.thumb.Parent() != o.win {
		t.Fatal("slider has no actual thumb child")
	}
	o.thumb.SetID(88)
	if o.win.WidgetData == o.input.CWidgetData() {
		t.Fatal("slider did not own its copied data")
	}
}

type sliderResult struct {
	Case, Step, Return int
	Input, Data        gui.SliderData
	Window, Thumb      []uint32
	Notices            []sliderNotice
	Focus              uint32
	Render             objectRenderResult
}

func (o *sliderOwner) snapshot(t *testing.T, id, step, ret int) sliderResult {
	o.win.Draw()
	o.thumb.Draw()
	return o.state(t, id, step, ret)
}
func (o *sliderOwner) state(t *testing.T, id, step, ret int) sliderResult {
	return sliderResult{id, step, ret, *o.input, *(*gui.SliderData)(o.win.WidgetData), o.rawWindow(o.win), o.rawWindow(o.thumb), append([]sliderNotice(nil), o.notices...), o.norm(uint32(uintptr(o.c.GUI.Focused().C()))), o.renderResult(t, id, step, ret)}
}
