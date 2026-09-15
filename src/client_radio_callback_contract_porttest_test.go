//go:build porttest

package opennox

import (
	"testing"

	"github.com/opennox/opennox/v1/client/gui"
)

func TestClientRadioSelectionCallbackMutation(t *testing.T) {
	o := newRadioOwner(t)
	for mode := 0; mode < 3; mode++ {
		o.create(t, mode+1, 0, 1, 1, 1, 0, 8)
		w := o.windows[0]
		calls := 0
		o.parent.SetFunc94(func(_ *gui.Window, ev gui.WindowEvent) gui.WindowEventResp {
			if ev.EventCode() == 16391 {
				calls++
				w.DrawData().SetGroup(8)
				w.DrawData().Field0 |= 4
			}
			return nil
		})
		switch mode {
		case 0:
			w.Func93(&gui.RawEvent{Event: 6})
		case 1:
			w.Func93(&gui.RawEvent{Event: 21, Arg1: 28, Arg2: 2})
		case 2:
			w.Func94(&gui.RawEvent{Event: 16392, Arg1: 1})
		}
		if calls != 1 {
			t.Fatal("selection must notify once")
		}
		selected := w.DrawData().Field0&4 != 0
		if selected != (mode != 1) {
			t.Fatal("keyboard toggles after notification; mouse/programmatic selection set the bit")
		}
		if o.windows[1].DrawData().Field0&4 == 0 || o.windows[2].DrawData().Field0&4 != 0 || o.windows[3].DrawData().Field0&4 == 0 {
			t.Fatal("group changes in the callback must govern sibling selection")
		}
	}
}
