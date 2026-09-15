//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"image"
	"runtime"
	"testing"
	"unsafe"
)

func TestClientWindowGeometryAndState(t *testing.T) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	o := newWindowHelperOwner(t)
	var out []windowHelperResult
	id := 0
	for _, flags := range []gui.StatusFlags{0, 8, 16, 0x408, 0x800, 0xffffffff} {
		for _, value := range []int{-17, -1, 0, 1, 8, 16, 0x7fffffff, -0x80000000} {
			for target := 0; target < 3; target++ {
				id++
				o.create()
				o.parent.SetPos(image.Pt(value, value^0x1234))
				o.win.Flags = flags
				w := o.win
				if target == 1 {
					w = o.windows[4]
				} else if target == 2 {
					w = nil
				}
				step := 0
				save := func(op, a, b int, p unsafe.Pointer, flag int) {
					out = append(out, o.invoke(id, step, op, w, a, b, p, flag))
					step++
				}
				for _, op := range []int{0, 1, 2, 3, 4, 0, 1, 2, 3, 5, 6, 8, 7, 8, 11, 12, 24, 25} {
					save(op, value, value^0x1234, nil, 0)
				}
				for _, p := range []unsafe.Pointer{nil, o.parent.C(), o.win.C(), o.windows[4].C()} {
					save(9, 0, 0, p, 0)
					save(10, 0, 0, p, 0)
				}
				for _, p := range []unsafe.Pointer{nil, unsafe.Pointer(o.images[1].C()), unsafe.Pointer(o.images[3].C())} {
					save(13, 0, 0, p, 0)
					save(14, 0, 0, p, 0)
					save(15, 0, 0, p, 0)
				}
			}
		}
	}
	effectsCapture(t, "window-geometry-state", out, 5040, "12647edec7ad804df12d0b174570dd7d80882c150631686fe1da44e9d6b275dc")
}
func TestClientWindowTreeAndLabels(t *testing.T) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	o := newWindowHelperOwner(t)
	data, free := alloc.New(gui.WindowData{})
	defer free()
	var out []windowHelperResult
	id := 0
	for hidden := 0; hidden < 16; hidden++ {
		for _, style := range []gui.StyleFlags{0, 0x80, 0x800, 0x880} {
			id++
			o.create()
			o.win.DrawData().Style = style
			if id%2 == 0 {
				o.windows[3].Flags &^= 8
			}
			for i, w := range o.windows[1:5] {
				if hidden&(1<<i) != 0 {
					w.Flags |= gui.StatusHidden
				}
			}
			step := 0
			save := func(op int, w *gui.Window, a, b int, p unsafe.Pointer, flag int) {
				out = append(out, o.invoke(id, step, op, w, a, b, p, flag))
				step++
			}
			for _, pt := range [][2]int{{0, 0}, {13, 16}, {15, 19}, {16, 20}, {23, 26}, {25, 28}, {37, 40}, {73, 76}, {74, 77}} {
				save(20, o.win, pt[0], pt[1], nil, 0)
			}
			for _, w := range append([]*gui.Window{nil}, o.windows...) {
				save(21, w, 0, 0, nil, 0)
			}
			save(16, o.win, 0, 0, unsafe.Pointer(&o.text[0]), 0)
			save(17, o.win, 0, 0, nil, 0)
			save(18, o.win, 0, 0, nil, 0)
			*data = gui.WindowData{Field0: 0xdeadbeef, Style: style, Window: o.parent, BgColorVal: 0x11223344, EnColorVal: 0x55667788, TextColorVal: 0x7fff7fff}
			if id%2 != 0 {
				data.FontPtr = o.font
			}
			data.SetText("raw label")
			raw := unsafe.Slice((*uint16)(unsafe.Add(data.C(), 72)), 64)
			raw[3] = 0xd800
			raw[63] = 0x6789
			save(19, o.win, 0, 0, data.C(), 0)
			save(18, o.win, 0, 0, nil, 0)
			save(19, o.win, 0, 0, nil, 0)
			save(19, nil, 0, 0, data.C(), 0)
			for _, interval := range [][2]int{{99, 105}, {102, 103}, {110, 100}} {
				for _, flag := range []int{0, 1, -1} {
					save(22, o.win, interval[0], interval[1], nil, flag)
					save(23, o.win, interval[0], interval[1], nil, flag)
				}
			}
		}
	}
	effectsCapture(t, "window-tree-labels", out, 2624, "bc2ddd44d57f341dd03d26294971eead58ce055a5a94f4086a8ce8cbcda3c047")
}
