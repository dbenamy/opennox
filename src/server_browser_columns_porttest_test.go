//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/legacy"
	"reflect"
	"testing"
)

func TestServerBrowserColumnForwarding(t *testing.T) {
	o := newListboxOwner(t)
	o.create(t, 8, 180, 100, nil)
	words, restore := legacy.PortTestServerBrowserWords()
	defer restore()
	*words["nox_wol_wnd_gameList_815012"] = uint32(uintptr(o.win.C()))
	type call struct {
		Column, Code int
		A, B         uint32
	}
	var calls []call
	for i, name := range []string{"dword_5d4594_815016", "dword_5d4594_815020", "dword_5d4594_815024", "dword_5d4594_815028", "dword_5d4594_815032"} {
		w := o.c.GUI.NewWindowRaw(o.parent, 8, 0, 0, 10, 10, nil)
		*words[name] = uint32(uintptr(w.C()))
		w.SetFunc94(func(_ *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
			a, b := ev.EventArgsC()
			calls = append(calls, call{i, ev.EventCode(), uint32(a), uint32(b)})
			return nil
		})
	}
	for _, code := range []int{22, 23, 16398, 16399, 16403, 16412} {
		for _, a := range []uint32{0, 1, 0xffffffff} {
			calls = nil
			b := uint32(91)
			legacy.PortTestServerBrowserColumnEvent(o.win.C(), code, a, b, false)
			var want []call
			if code != 22 && code != 23 {
				outb := b
				if code == 16403 || code == 16412 {
					outb = 0
				}
				for i := 0; i < 5; i++ {
					want = append(want, call{i, code, a, outb})
				}
			}
			if !reflect.DeepEqual(calls, want) {
				t.Fatal("forwarding", code, a, calls, want)
			}
		}
	}
	sender := o.c.GUI.NewWindowRaw(o.parent, 8, 0, 0, 10, 10, nil)
	for _, id := range []uint{10037, 10038, 10042, 10043} {
		sender.SetID(id)
		calls = nil
		legacy.PortTestServerBrowserColumnEvent(o.win.C(), 16400, uint32(uintptr(sender.C())), 0xffffffff, false)
		var want []call
		if id >= 10038 && id <= 10042 {
			for i := 0; i < 5; i++ {
				want = append(want, call{i, 16403, 0xffffffff, 0})
			}
		}
		if !reflect.DeepEqual(calls, want) {
			t.Fatal("selection", id, calls, want)
		}
	}
	for _, code := range []int{19, 20} {
		calls = nil
		up := o.c.GUI.NewWindowRaw(o.parent, 8, 0, 0, 10, 10, nil)
		d := (*gui.ScrollListBoxData)(o.win.WidgetData)
		d.Field_7 = up.C()
		d.Field_8 = up.C()
		o.win.SetFunc94(func(w *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
			if ev.EventCode() == 2 {
				return gui.RawEventResp(legacy.PortTestServerBrowserColumnEvent(w.C(), 2, 0, 0, false))
			}
			a, b := ev.EventArgsC()
			calls = append(calls, call{-1, ev.EventCode(), uint32(a), uint32(b)})
			return nil
		})
		legacy.PortTestServerBrowserColumnEvent(o.win.C(), code, 0, 0, true)
		var want []call
		for i := -1; i < 5; i++ {
			want = append(want, call{i, 16391, uint32(uintptr(up.C())), 0})
		}
		if !reflect.DeepEqual(calls, want) {
			t.Fatal("scroll", code, calls, want)
		}
		d.Field_7 = nil
		d.Field_8 = nil
	}
}
func TestServerBrowserInfoPosition(t *testing.T) {
	words, restore2 := legacy.PortTestServerBrowserWords()
	defer restore2()
	for _, mode := range []uint32{0, 1, 2, 0xffffffff} {
		*words["dword_587000_87408"] = mode
		for _, x := range []int32{-2147483648, -66, -1, 0, 216, 281, 535, 600, 2147483647} {
			for _, y := range []int32{-2147483648, -21, -1, 0, 27, 55, 75, 351, 451, 2147483647} {
				out := [2]uint32{}
				legacy.PortTestServerBrowserInfoPosition(x, y, &out)
				wrap := func(v int64) int64 { return int64(int32(uint32(v))) }
				px, py := wrap(int64(x)-65), wrap(int64(y)-20)
				if wrap(px+130) > 600 {
					px = 470
				}
				if wrap(py+120) > 451 {
					py = 331
				}
				low := uint32(27)
				if mode == 1 {
					low = 55
				}
				if uint32(py) < low {
					py = int64(low)
				}
				if px < 216 {
					px = 216
				}
				if out != [2]uint32{uint32(px), uint32(py)} {
					t.Fatal("position", mode, x, y, out, px, py)
				}
			}
		}
	}
}
