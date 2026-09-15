//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/legacy"
	"runtime"
	"testing"
)

func TestClientEntryEventsAndContext(t *testing.T) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	o := newEntryOwner(t)
	var out []entryResult
	id := 0
	for _, lang := range []int{0, 6, 8} {
		for _, flags := range []gui.StatusFlags{0, 8, 0x408, 0x88} {
			for track := 0; track < 2; track++ {
				for owner := 0; owner < 2; owner++ {
					id++
					d := o.create(t, lang, flags, func(draw *gui.WindowData, d *gui.EntryFieldData) {
						if track == 0 {
							draw.Style &^= 0x100
						}
						if owner == 0 {
							draw.Window = nil
						}
						draw.Field0 = 0xabcdef00
					})
					for step, code := range []int{5, 7, 8, 17, 18, 7, 21, -1} {
						ret := gui.EventRespInt(o.win.Func93(&gui.RawEvent{Event: code, Arg1: 0x1000e, Arg2: 0xffffffff}))
						out = append(out, o.snapshot(t, id, step, ret, false))
					}
					for step, a := range []uintptr{0, 1, 2, 0xffffffff, 0} {
						ret := gui.EventRespInt(o.win.Func94(&gui.RawEvent{Event: 23, Arg1: a}))
						out = append(out, o.snapshot(t, id, 8+step, ret, false))
					}
					if gui.EventRespInt(o.win.Func94(&gui.RawEvent{Event: 16413})) != int(uintptr(o.win.WidgetData)) {
						t.Fatal("entry get-text did not return owned buffer")
					}
					o.c.GUI.Focus(nil)
					o.c.GUI.Focus(o.win)
					o.env.Context(false)
					legacy.NoxInputOnChar('A')
					out = append(out, o.snapshot(t, id, 13, 0, false))
					o.env.Context(true)
					legacy.NoxInputOnChar('B')
					out = append(out, o.snapshot(t, id, 14, 0, false))
					d.Field_1044 = 1
					ret := o.key(28, 2)
					out = append(out, o.snapshot(t, id, 15, ret, false))
				}
			}
		}
	}
	effectsCapture(t, "entry-events-context", out, len(out), "11ba21681640d6958bca1eab0c30a8afe559daeaf012dae2d4e8fe732f5e11c0")
}
