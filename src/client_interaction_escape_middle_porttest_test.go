//go:build porttest

package opennox

import (
	"encoding/binary"
	"slices"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func TestClientInteractionEscapeMiddleControllers(t *testing.T) {
	o := newInventoryWindowOwner(t)
	o.reset(t)
	o.construct(t)
	o.initAmount(t)
	words, restore := legacy.PortTestClientInteractionWords()
	defer restore()
	quick, restoreQuick := legacy.PortTestQuickbarWords()
	defer restoreQuick()
	for _, p := range quick {
		*p = 0
	}
	*words["dword_5d4594_1064856"] = 0
	*o.windowWords["dword_5d4594_1047520"] = 0
	capturedWord := serverConfigOwnBytes(t, 0x5D4594, 1047928, 4)
	clear(capturedWord)
	expanded := serverConfigOwnBytes(t, 0x5D4594, 1049476, 4)
	savedRow := serverConfigOwnBytes(t, 0x5D4594, 1047912, 1)
	savedRow[0] = 2
	panels := serverConfigOwnBytes(t, 0x5D4594, 1048196, 1024)
	clear(panels)
	cursor := serverConfigOwnBytes(t, 0x5D4594, 1096672, 4)
	clear(cursor)
	cancel := serverConfigOwnBytes(t, 0x5D4594, 1319100, 4)
	clear(cancel)
	record, free := alloc.New([64]uint32{})
	defer free()
	*record = [64]uint32{}
	restoreMain := legacy.PortTestBookQuickbar(unsafe.Pointer(record))
	defer restoreMain()
	oldFade := legacy.Get_nox_gameDisableMapDraw_5d4594_2650672()
	legacy.Set_nox_gameDisableMapDraw_5d4594_2650672(0)
	defer legacy.Set_nox_gameDisableMapDraw_5d4594_2650672(oldFade)
	amount := legacy.Get_nox_gui_itemAmount_dialog_1319228()
	identify := o.mainWindow().ChildByID(9150)
	count, freeCount := alloc.CString16("3")
	defer freeCount()
	amount.ChildByID(3601).Func94(gui.AsWindowEvent(16385, uintptr(unsafe.Pointer(count)), 0))
	oldDialog := nox_gui_curDialog_830224
	defer func() { nox_gui_curDialog_830224 = oldDialog }()
	var calls int
	var currentMask int
	nox_gui_curDialog_830224 = o.c.GUI.NewWindowRaw(nil, 8, 0, 0, 30, 30, func(_ *gui.Window, e gui.WindowEvent) gui.WindowEventResp {
		if e.EventCode() == 16391 {
			calls++
			if *o.windowWords["dword_5d4594_1319268"] != 0 || binary.LittleEndian.Uint32(expanded) != 0 || !identify.GetFlags().IsHidden() {
				t.Fatal("dialog callback before amount/expanded bar/identify close", currentMask)
			}
		}
		return nil
	})
	defer nox_gui_curDialog_830224.Destroy()
	button := o.c.GUI.NewWindowRaw(nox_gui_curDialog_830224, 8, 0, 0, 10, 10, nil)
	button.SetID(4002)
	type row struct {
		Mask         int
		DialogCalls  int
		Sounds       [][2]int
		Selected     byte
		IdentifyMode uint32
	}
	var captured []row
	for mask := 0; mask < 8; mask++ {
		currentMask = mask
		calls = 0
		o.sounds = nil
		amount.SetHidden(mask&1 == 0)
		*o.windowWords["dword_5d4594_1319268"] = uint32(mask & 1)
		binary.LittleEndian.PutUint32(expanded, uint32((mask>>1)&1))
		record[50] = 0 // selected row before closing.
		identify.SetHidden(mask&4 == 0)
		*o.windowWords["dword_5d4594_1049864"] = 0
		if mask&4 != 0 {
			*o.windowWords["dword_5d4594_1049864"] = 5
		}
		*o.windowWords["dword_5d4594_1063116"], *o.windowWords["dword_5d4594_1063120"] = 123, 456
		interactionCall("nox_xxx_consoleEsc_49B7A0")
		var wantSounds [][2]int
		selected := byte(0)
		if mask&2 != 0 {
			wantSounds = append(wantSounds, [2]int{800, 100})
			selected = 2
		}
		if calls != 1 || !slices.Equal(o.sounds, wantSounds) || byte(record[50]) != selected || *o.windowWords["dword_5d4594_1049864"] != 0 {
			t.Fatal("Escape middle controller batch", mask, calls, o.sounds, record[50])
		}
		if mask&4 != 0 && (*o.windowWords["dword_5d4594_1063116"] != 0 || *o.windowWords["dword_5d4594_1063120"] != 0) {
			t.Fatal("identify control words")
		}
		captured = append(captured, row{mask, calls, append([][2]int(nil), o.sounds...), byte(record[50]), *o.windowWords["dword_5d4594_1049864"]})
	}
	interactionCapture(t, "escape-middle", captured)
}
