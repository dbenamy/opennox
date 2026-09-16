//go:build porttest

package opennox

import (
	"testing"
	"unsafe"

	"github.com/opennox/libs/strman"
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/client/audio/ail"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/legacy/dialog"
	"github.com/opennox/opennox/v1/legacy/timer"
)

type optionsAudioOwner struct {
	*optionsOwner
	menu    bool
	words   map[int]*uint32
	timers  [3]*timer.Timer
	buttons [3][2]*gui.Window
	sounds  [][2]int
}

func newOptionsAudioOwner(t *testing.T, menu bool) *optionsAudioOwner {
	o := &optionsAudioOwner{optionsOwner: newOptionsOwner(t), menu: menu}
	var restore func()
	o.words, restore = legacy.PortTestOptionsWords()
	t.Cleanup(restore)
	o.timers, restore = legacy.PortTestOptionsTimers()
	t.Cleanup(restore)
	*o.words[831092] = 1
	*o.words[816376] = 1
	rootOff := 1309820
	if menu {
		rootOff = 1309720
	}
	*o.words[rootOff] = uint32(uintptr(o.root.C()))
	o.root.SetFunc94(func(w *gui.Window, e gui.WindowEvent) gui.WindowEventResp {
		a, b := e.EventArgsC()
		return gui.RawEventResp(legacy.PortTestOptionsEvent(menu, w, e.EventCode(), (*gui.Window)(unsafe.Pointer(a)), int(b)))
	})
	for ch := 0; ch < 3; ch++ {
		draw := gui.WindowData{Style: gui.StyleCheckBox, Window: o.root}
		w := gui.NewCheckBoxRaw(o.c.GUI, o.root, 8, 0, 0, 10, 10, &draw)
		w.SetID(uint(361 + ch))
		o.buttons[ch][1] = w
		w = o.c.GUI.NewWindowRaw(o.root, 8, 0, 0, 10, 10, nil)
		w.SetID(uint(361 + ch))
		o.buttons[ch][0] = w
	}
	t.Cleanup(legacy.PortTestClientSoundObserver(func(id, volume int) { o.sounds = append(o.sounds, [2]int{id, volume}) }))
	oldDialog := legacy.Dialogs
	t.Cleanup(func() { legacy.Dialogs = oldDialog })
	var state [5]uint32
	var driver ail.Driver
	var timers [4]timer.TimerGroup
	legacy.Dialogs = dialog.NewDialog("dialog", o.words[122848], &state[0], &state[1], &state[2], &state[3], &driver, &state[4], o.c.srv.Strings, &timers[0], &timers[1], &timers[2], &timers[3], nil, nil, nil, nil, nil, nil, nil)
	var entries []strman.Entry
	for i, name := range []string{"OptionsPreviewA", "OptionsPreviewB", "OptionsPreviewC"} {
		off := 172892 + 4*i
		b := memmap.BlobByAddr(0x587000).Data[off : off+4]
		old := append([]byte(nil), b...)
		t.Cleanup(func() { copy(b, old) })
		*(*uint32)(unsafe.Pointer(&b[0])) = uint32(uintptr(unsafe.Pointer(alloc.InternCString(name))))
		entries = append(entries, strman.Entry{ID: strman.ID("Options.c:" + name), Vals: []strman.Variant{{Str: name, Str2: name + ".wav"}}})
	}
	configure, restore := o.c.srv.Server.PortTestMeterStrings(entries...)
	t.Cleanup(restore)
	configure(0)
	cursor := memmap.BlobByAddr(0x5D4594).Data[1309744:1309748]
	old := append([]byte(nil), cursor...)
	t.Cleanup(func() { copy(cursor, old) })
	return o
}
func (o *optionsAudioOwner) buttonOffset(ch int) int {
	if o.menu {
		return []int{1309728, 1309732, 1309736}[ch]
	}
	return []int{1309828, 1309836, 1309832}[ch]
}
