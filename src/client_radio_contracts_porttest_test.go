//go:build porttest

package opennox

import (
	"image"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func TestClientRadioSelectionContract(t *testing.T) {
	o := newRadioOwner(t)
	o.create(t, 1, 0, 1, 1, 1, 0, 8)
	w := o.windows[0]
	ret := gui.EventRespInt(w.Func93(&gui.RawEvent{Event: 6, Arg1: 123}))
	if ret != 1 || w.DrawData().Field0&4 == 0 {
		t.Fatal("mouse did not select radio")
	}
	if o.windows[1].DrawData().Field0&4 != 0 || o.windows[2].DrawData().Field0&4 == 0 || o.windows[3].DrawData().Field0&4 != 0 {
		t.Fatal("selection did not preserve exact sibling group behavior")
	}
	if len(o.notices) != 1 || o.notices[0].Code != 16391 || o.notices[0].A != 0xe1000002 || o.notices[0].B != 123 {
		t.Fatal("selection notification")
	}
	before := o.notices[0].Flags
	if len(before) != 4 || before[0]&4 != 0 || before[1]&4 == 0 || before[3]&4 == 0 {
		t.Fatal("notification must precede group selection mutations")
	}
	o.notices = nil
	if gui.EventRespInt(w.Func93(&gui.RawEvent{Event: 6})) != 0 || len(o.notices) != 0 {
		t.Fatal("selected unhighlighted mouse result")
	}
	w.DrawData().Field0 |= 2
	if gui.EventRespInt(w.Func93(&gui.RawEvent{Event: 7})) != 1 || len(o.notices) != 0 {
		t.Fatal("selected highlighted mouse result")
	}
}

func TestClientRadioOwnedData(t *testing.T) {
	o := newRadioOwner(t)
	input := gui.RadioButtonData{Field0: 7}
	draw := gui.WindowData{Style: gui.StyleRadioButton, Window: o.parent}
	w := newRadioButton(o.c.GUI, o.parent, 8, 1, 1, 20, 20, &draw, &input)
	if w == nil {
		t.Fatal("radio constructor")
	}
	ptr := w.WidgetData
	if ptr == unsafe.Pointer(&input) || *(*gui.RadioButtonData)(ptr) != input || !alloc.PortTestAllocationLive(ptr) {
		t.Fatal("radio does not own its copied data")
	}
	defer func() {
		if alloc.PortTestAllocationLive(ptr) {
			alloc.FreePtr(ptr)
		}
	}()
	w.Destroy()
	o.c.GUI.FreeDestroyed()
	if alloc.PortTestAllocationLive(ptr) {
		t.Fatal("destroyed radio leaked its owned data")
	}
}

func TestClientRadioDeferredWindowCleanup(t *testing.T) {
	g := gui.New(nil)
	t.Cleanup(g.DestroyAll)
	var destroyed []int
	other := 0
	callback := func(id int) gui.WindowFunc {
		return func(w *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
			switch ev.EventCode() {
			case 2:
				if !w.Flags.Has(gui.StatusDestroyed) {
					t.Error("cleanup must observe queued destruction")
				}
				destroyed = append(destroyed, id)
			case 999:
				other++
			}
			return nil
		}
	}
	parent := g.NewWindowRaw(nil, 8, 0, 0, 20, 20, callback(1))
	child := g.NewWindowRaw(parent, 8, 0, 0, 10, 10, callback(2))
	parent.Destroy()
	parent.Destroy()
	child.Destroy()
	parent.Func94(&gui.RawEvent{Event: 999})
	if len(destroyed) != 0 || other != 0 {
		t.Fatal("queued windows must defer cleanup and reject ordinary events")
	}
	g.FreeDestroyed()
	if len(destroyed) != 2 || destroyed[0] != 1 || destroyed[1] != 2 {
		t.Fatalf("deferred cleanup order/count: %v", destroyed)
	}
	g.FreeDestroyed()
	if len(destroyed) != 2 {
		t.Fatal("cleanup delivered more than once")
	}
	var second *gui.Window
	first := g.NewWindowRaw(nil, 8, 0, 0, 10, 10, func(w *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
		callback(3)(w, ev)
		if ev.EventCode() == 2 {
			second.Destroy()
		}
		return nil
	})
	second = g.NewWindowRaw(nil, 8, 0, 0, 10, 10, callback(4))
	first.Destroy()
	g.FreeDestroyed()
	if len(destroyed) != 3 || destroyed[2] != 3 {
		t.Fatal("newly queued cleanup should wait for the next pass")
	}
	g.FreeDestroyed()
	if len(destroyed) != 4 || destroyed[3] != 4 {
		t.Fatal("cleanup queued by another callback was lost")
	}
}

func TestClientRadioSliderOwnedData(t *testing.T) {
	o := newSliderOwner(t)
	for horizontal := 0; horizontal < 2; horizontal++ {
		o.create(t, horizontal+1, horizontal, 0, 1, 1, 8, 0, 100, image.Pt(20, 70))
		ptr := o.win.WidgetData
		if !alloc.PortTestAllocationLive(ptr) {
			t.Fatal("slider allocation is not owned")
		}
		o.win.Destroy()
		if !alloc.PortTestAllocationLive(ptr) {
			t.Fatal("slider data released before deferred cleanup")
		}
		o.c.GUI.FreeDestroyed()
		if alloc.PortTestAllocationLive(ptr) {
			t.Fatal("slider data survived deferred cleanup")
		}
	}
}
