//go:build porttest

package opennox

import (
	"encoding/binary"
	"fmt"
	"github.com/opennox/libs/strman"
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"testing"
	"unsafe"
)

func TestServerPanelsAdvancedEdits(t *testing.T) {
	type row struct {
		Event, Focus int
		Input, Shown string
		Value        uint32
		Result       int
	}
	var rows []row
	o := newServerOptionsOwner(t)
	o.installSubpanels(t, true)
	legacy.PortTestServerPanelsConstruct("advserv", o.options, unsafe.Pointer(&o.settings[0]))
	w := (*gui.Window)(unsafe.Pointer(uintptr(*o.optionWords["panel-1316972"])))
	if w == nil {
		t.Fatal("advanced edit root")
	}
	field := unsafe.Slice(memmap.PtrUint8(0x5D4594, 371586), 4)
	child := w.ChildByID(2110)
	for _, event := range []int{16415, 16387} {
		for _, focus := range []int{0, 1} {
			if event == 16415 && focus == 1 {
				continue
			}
			for _, tc := range []struct {
				text  string
				value uint32
			}{{"", 123}, {"-1", 0}, {"0", 0}, {"1", 1}, {"100", 100}, {"2147483647", 2147483647}, {" 27", 27}, {"42tail", 42}, {"bad", 0}} {
				binary.LittleEndian.PutUint32(field, 123)
				serverPanelsSetText(child, tc.text)
				a, b := uintptr(child.C()), uintptr(0)
				wantResult := 1
				if event == 16387 {
					a, b = uintptr(focus), 2110
					if focus == 1 {
						wantResult = 0
					}
				}
				result := gui.EventRespInt(w.Func94(gui.AsWindowEvent(event, a, b)))
				want := tc.value
				if focus == 1 {
					want = 123
				}
				shown := serverPanelsGetText(child)
				value := binary.LittleEndian.Uint32(field)
				if value != want || shown != tc.text || result != wantResult {
					t.Fatalf("advanced edit %d/%d/%q got %d/%q/%d", event, focus, tc.text, value, shown, result)
				}
				rows = append(rows, row{event, focus, tc.text, shown, value, result})
			}
		}
	}
	w.Func94(gui.AsWindowEvent(16391, uintptr(w.ChildByID(2130).C()), 0))
	spellbookCapture(t, "server-panels-advanced-edits", rows, "e2e8d49a34c1567893e1f3b166fe3e14d0841e81c6413193cb9d84589c7bfca6")
}

func TestServerPanelsAudioSetting(t *testing.T) {
	type row struct {
		Value, Threshold, Result int
		Stored                   uint32
		Text                     string
	}
	var rows []row
	o := newServerOptionsOwner(t)
	o.installSubpanels(t, true)
	configure, restore := o.c.srv.PortTestMeterStrings(strman.Entry{ID: "advserv.c:AudCullDesc", Vals: []strman.Variant{{Str: "Sound range %d"}}}, strman.Entry{ID: "WindowDir:Blank", Vals: []strman.Variant{{Str: ""}}})
	t.Cleanup(restore)
	configure(0)
	oldStrings := strMan
	strMan = o.c.srv.Strings()
	t.Cleanup(func() { strMan = oldStrings })
	oldThreshold := o.c.srv.ai.soundMuteThreshold
	t.Cleanup(func() { o.c.srv.ai.soundMuteThreshold = oldThreshold })
	buf := unsafe.Slice(memmap.PtrUint8(0x5D4594, 1316716), 256)
	old := append([]byte(nil), buf...)
	t.Cleanup(func() { copy(buf, old) })
	legacy.PortTestServerPanelsConstruct("advserv", o.options, unsafe.Pointer(&o.settings[0]))
	w := (*gui.Window)(unsafe.Pointer(uintptr(*o.optionWords["panel-1316972"])))
	if w == nil {
		t.Fatal("audio panel root")
	}
	for _, value := range []int{-2147483648, -1, 0, 1, 50, 100, 101, 2147483647} {
		result := gui.EventRespInt(w.Func94(gui.AsWindowEvent(16393, uintptr(w.ChildByID(2119).C()), uintptr(uint32(value)))))
		want := value
		if want < 0 {
			want = 0
		}
		if want > 100 {
			want = 100
		}
		stored := *memmap.PtrUint32(0x5D4594, 371590)
		p := gui.EventRespInt(w.ChildByID(2120).Func94(gui.AsWindowEvent(16386, 0, 0)))
		text := alloc.GoString16((*uint16)(unsafe.Pointer(uintptr(uint32(p)))))
		if result != 1 || stored != uint32(value) || o.c.srv.ai.soundMuteThreshold != want || text != fmt.Sprintf("Sound range %d", value) {
			t.Fatalf("audio %d result %d stored %x threshold %d text %q", value, result, stored, o.c.srv.ai.soundMuteThreshold, text)
		}
		rows = append(rows, row{value, o.c.srv.ai.soundMuteThreshold, result, stored, text})
	}
	w.Func94(gui.AsWindowEvent(16391, uintptr(w.ChildByID(2130).C()), 0))
	spellbookCapture(t, "server-panels-audio-setting", rows, "f30e45e4d4c627faef79abe75c3b39fcb8960fde69a321ac67ac5ce2f8482bba")
}
