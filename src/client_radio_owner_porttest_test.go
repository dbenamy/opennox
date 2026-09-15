//go:build porttest

package opennox

import (
	"image"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/client/noxrender"
)

type radioNotice struct {
	Code  int
	A, B  uint32
	Flags []uint32
}
type radioOwner struct {
	*objectRenderOwner
	parent  *gui.Window
	windows []*gui.Window
	notices []radioNotice
}

func newRadioOwner(t *testing.T) *radioOwner {
	o := &radioOwner{objectRenderOwner: newObjectRenderOwner(t)}
	o.c.GUI = gui.New(o.c.Render())
	o.parent = o.c.GUI.NewWindowRaw(nil, 8, 3, 4, 90, 90, func(_ *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
		a, b := ev.EventArgsC()
		av := uint32(a)
		switch ev.EventCode() {
		case 16384, 16389, 16390, 16391:
			av = o.norm(av)
		}
		var flags []uint32
		for _, w := range o.windows {
			flags = append(flags, w.DrawData().Field0)
		}
		o.notices = append(o.notices, radioNotice{ev.EventCode(), av, uint32(b), flags})
		return nil
	})
	o.parent.SetID(55)
	t.Cleanup(func() { o.destroy(); o.c.GUI.DestroyAll() })
	return o
}
func (o *radioOwner) norm(v uint32) uint32 {
	if v == 0 {
		return 0
	}
	if v == uint32(uintptr(o.parent.C())) {
		return 0xe1000001
	}
	for i, w := range o.windows {
		if v == uint32(uintptr(w.C())) {
			return 0xe1000002 + uint32(i)
		}
		if w.WidgetData != nil && v == uint32(uintptr(w.WidgetData)) {
			return 0xe2000001 + uint32(i)
		}
	}
	if n, ok := o.c.imageRefs[v]; ok {
		return n
	}
	return v
}
func (o *radioOwner) raw(w *gui.Window) []uint32 {
	out := append([]uint32(nil), unsafe.Slice((*uint32)(w.C()), 101)...)
	for _, i := range []int{8, 13, 15, 17, 19, 21, 23, 59, 93, 94, 95, 96, 97, 98, 99, 100} {
		out[i] = o.norm(out[i])
	}
	return out
}
func (o *radioOwner) destroy() {
	for _, w := range o.windows {
		w.Destroy()
	}
	o.c.GUI.FreeDestroyed()
	o.windows = nil
}
func (o *radioOwner) create(t *testing.T, id, mode, parent, owner, track, center int, flags gui.StatusFlags, configure ...func(*gui.WindowData)) {
	t.Helper()
	o.destroy()
	o.resetRender(uint32(id), 120)
	o.notices = nil
	for i := 0; i < 4; i++ {
		draw := gui.WindowData{Style: gui.StyleRadioButton, Window: o.parent, BgColorVal: 0x21082108, EnColorVal: 0x03e003e0, HlColorVal: 0x7fff7fff, DisColorVal: 0x42104210, SelColorVal: 0x7c007c00, TextColorVal: 0x7fff7fff}
		draw.SetText("Radio")
		draw.SetGroup(7)
		if i == 2 {
			draw.SetGroup(8)
		}
		if track != 0 {
			draw.Style |= 0x100
		}
		draw.BgImageHnd = noxrender.ImageHandle(o.images[1].C())
		draw.EnImageHnd = noxrender.ImageHandle(o.images[3].C())
		draw.DisImageHnd = noxrender.ImageHandle(o.images[5].C())
		draw.SelImageHnd = noxrender.ImageHandle(o.images[7].C())
		draw.HlImageHnd = noxrender.ImageHandle(o.images[9].C())
		if owner == 0 {
			draw.Window = nil
		}
		if i == 0 {
			for _, f := range configure {
				f(&draw)
			}
		} else {
			draw.Field0 = 4
		}
		par := o.parent
		if i == 0 && parent == 0 {
			par = nil
		}
		var w *gui.Window
		size := image.Pt(61, 17)
		if i == 0 {
			size = []image.Point{image.Pt(61, 17), image.Pt(11, 10), image.Pt(40, 24)}[id%3]
		}
		if i == 3 {
			draw.Style = 0
			w = o.c.GUI.NewWindowRaw(par, 8, 10, 12+18*i, 61, 17, nil)
			w.CopyDrawData(&draw)
		} else {
			data := gui.RadioButtonData{Field0: uint32(center)}
			w = newRadioButton(o.c.GUI, par, flags|gui.StatusFlags(mode*128), 10, 12+18*i, size.X, size.Y, &draw, &data)
		}
		if w == nil {
			t.Fatal("radio fixture window allocation")
		}
		w.SetID(uint(77 + i))
		o.windows = append(o.windows, w)
	}
}

type radioResult struct {
	Case, Step, Return int
	Windows            [][]uint32
	Data               []uint32
	Notices            []radioNotice
	Focus              uint32
	Render             objectRenderResult
}

func (o *radioOwner) snapshot(t *testing.T, id, step, ret int) radioResult {
	r := radioResult{Case: id, Step: step, Return: ret}
	for _, w := range o.windows {
		w.Draw()
		r.Windows = append(r.Windows, o.raw(w))
		if w.WidgetData != nil {
			r.Data = append(r.Data, *(*uint32)(w.WidgetData))
		}
	}
	r.Notices = append([]radioNotice(nil), o.notices...)
	r.Focus = o.norm(uint32(uintptr(o.c.GUI.Focused().C())))
	r.Render = o.renderResult(t, id, step, ret)
	return r
}
