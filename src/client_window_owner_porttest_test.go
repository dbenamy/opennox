//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"golang.org/x/image/font/basicfont"
	"image"
	"testing"
	"unsafe"
)

type windowHelperNotice struct {
	Window   uint32
	Code     int
	A, B     uint32
	Geometry [6]int32
}
type windowHelperResult struct {
	Case, Step, Op int
	Return         uint32
	XY             [2]uint32
	Windows        [][]uint32
	Notices        []windowHelperNotice
}
type windowHelperOwner struct {
	font unsafe.Pointer
	*entryOwner
	windows []*gui.Window
	calls   []windowHelperNotice
	text    []uint16
}

func newWindowHelperOwner(t *testing.T) *windowHelperOwner {
	o := &windowHelperOwner{entryOwner: newEntryOwner(t)}
	fontPtr, freeFont := o.c.Render().GetFonts().PortTestWindowFont(basicfont.Face7x13)
	o.font = fontPtr
	t.Cleanup(freeFont)
	var free func()
	o.text, free = alloc.Make([]uint16{}, 16)
	t.Cleanup(free)
	copy(o.text, []uint16{'h', 'e', 'l', 'l', 'o'})
	t.Cleanup(func() {
		for _, w := range o.windows {
			w.Flags &^= gui.StatusDestroyed
		}
	})
	return o
}
func (o *windowHelperOwner) normalize(v uint32) uint32 {
	if v == uint32(uintptr(o.font)) {
		return 0xe3000001
	}
	for i, w := range o.windows {
		if v == uint32(uintptr(w.C())) {
			return 0xe1000001 + uint32(i)
		}
	}
	if v == uint32(uintptr(unsafe.Pointer(&o.text[0]))) {
		return 0xe2000001
	}
	if n, ok := o.c.imageRefs[v]; ok {
		return n
	}
	return v
}
func (o *windowHelperOwner) create() {
	// Some raw-flag cases set the destruction bit without queueing destruction.
	// Restore it before invoking the real owner cleanup.
	for _, w := range o.windows {
		w.Flags &^= gui.StatusDestroyed
	}
	o.entryOwner.destroy()
	o.parent.Flags = 8
	o.parent.SetPos(image.Pt(3, 4))
	o.windows = []*gui.Window{o.parent}
	o.calls = nil
	callback := func(w *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
		a, b := ev.EventArgsC()
		o.calls = append(o.calls, windowHelperNotice{Window: o.normalize(uint32(uintptr(w.C()))), Code: ev.EventCode(), A: o.normalize(uint32(a)), B: o.normalize(uint32(b)), Geometry: [6]int32{int32(w.Off.X), int32(w.Off.Y), int32(w.SizeVal.X), int32(w.SizeVal.Y), int32(w.EndPos.X), int32(w.EndPos.Y)}})
		if ev.EventCode() == 16386 || ev.EventCode() == 16413 {
			return gui.RawEventResp(uintptr(unsafe.Pointer(&o.text[0])))
		}
		return gui.RawEventResp(0x12345678)
	}
	o.win = o.c.GUI.NewWindowRaw(o.parent, 8, 10, 12, 60, 60, callback)
	o.win.SetID(100)
	o.win.DrawData().FontPtr = o.font
	o.windows = append(o.windows, o.win)
	for _, cfg := range [][6]int{{1, 101, 2, 3, 30, 30}, {1, 102, 10, 10, 40, 40}, {2, 103, 1, 1, 10, 10}, {3, 104, 2, 2, 12, 12}} {
		w := o.c.GUI.NewWindowRaw(o.windows[cfg[0]], 8, cfg[2], cfg[3], cfg[4], cfg[5], callback)
		w.SetID(uint(cfg[1]))
		o.windows = append(o.windows, w)
	}
	o.calls = nil
}
func (o *windowHelperOwner) invoke(id, step, op int, w *gui.Window, a, b int, p unsafe.Pointer, flag int) windowHelperResult {
	ret, xy := legacy.PortTestWindowHelper(op, w, a, b, p, flag)
	result := windowHelperResult{Case: id, Step: step, Op: op, Return: o.normalize(uint32(ret)), XY: xy}
	for _, w := range o.windows {
		words := append([]uint32(nil), unsafe.Slice((*uint32)(w.C()), 101)...)
		for _, i := range []int{8, 13, 15, 17, 19, 21, 23, 59, 93, 94, 95, 96, 97, 98, 99, 100} {
			words[i] = o.normalize(words[i])
		}
		result.Windows = append(result.Windows, words)
	}
	result.Notices = append([]windowHelperNotice(nil), o.calls...)
	o.calls = nil
	return result
}
