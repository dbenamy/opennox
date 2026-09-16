//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/things"
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

func TestQuickbarDirectionControls(t *testing.T) {
	q := newQuickbarOwner(t)
	q.ownFullQuickbar(t)
	var rows []quickbarResult
	for slot := 0; slot < 5; slot++ {
		for _, flags := range []uint32{0, 0x400, 0x200000, 0x200400} {
			for _, event := range []uint32{0, 5, 6, 7, 8, 12} {
				q.reset(t)
				*(*byte)(unsafe.Add(unsafe.Pointer(&q.players[0]), 2251)) = 1
				q.configure([]server.PortTestSpellClassDef{{Index: 1, Flags: uint32(things.SpellClassAny) | flags, Valid: true}})
				q.call("nox_xxx_quickBarCreate_45E190")
				q.bar[2*slot] = 1
				q.bar[2*slot+1] = 0x12345680
				arrow := (*gui.Window)(unsafe.Pointer(uintptr(q.bar[58+slot])))
				var control *gui.Window
				for child := arrow.Field100Ptr; child != nil; child = child.Prev() {
					if (*[101]uint32)(child.C())[92]>>16 == 4 {
						control = child
						break
					}
				}
				q.check(t, control != nil, "direction control is an actual constructor child")
				ret := q.call("sub_45F520", uint32(uintptr(control.C())), event)
				wantFlag := uint32(0x12345680)
				if event == 5 && flags == 0 {
					wantFlag ^= 1
				}
				q.check(t, q.bar[2*slot+1] == wantFlag, "direction toggle changes only permitted low flag bit")
				wantRet := uint32Bool(event >= 5 && event <= 7)
				q.check(t, ret == wantRet, "direction event handling")
				if event == 5 {
					wantSound := 798
					if flags == 0 {
						wantSound = 921
					} else {
						wantSound = 925
					}
					q.check(t, len(q.sounds) > 0 && q.sounds[len(q.sounds)-1][0] == wantSound, "direction feedback sound")
				}
				rows = append(rows, q.snapshot(fmt.Sprintf("slot%d-flags%d-event%d", slot, flags, event), ret))
			}
		}
	}
	spellbookCapture(t, "quickbar-direction-controls", rows, "6a12e44daa4dd0a3a66ea393025b998a645688d49f676b31c8b07f208933fa96")
}

func TestQuickbarRowAndTrapControls(t *testing.T) {
	q := newQuickbarOwner(t)
	q.ownFullQuickbar(t)
	var rows []quickbarResult
	for action := uint32(0); action < 5; action++ {
		for _, event := range []uint32{0, 5, 6, 7, 8, 12} {
			q.reset(t)
			*(*byte)(unsafe.Add(unsafe.Pointer(&q.players[0]), 2251)) = 1
			*(*uint32)(unsafe.Add(unsafe.Pointer(&q.players[0]), 3832)) = 1
			q.call("nox_xxx_quickBarCreate_45E190")
			parent := (*gui.Window)(unsafe.Pointer(uintptr(*q.quickWords["dword_5d4594_1049508"])))
			control := q.c.GUI.NewWindowRaw(parent, gui.StatusFlags(1032), 0, 0, 10, 10, nil)
			(*[101]uint32)(control.C())[92] = action
			ret := q.call("nox_xxx_quickbarTrapUpDownProc_45F630", uint32(uintptr(control.C())), event)
			q.check(t, ret == uint32Bool(event >= 5 && event <= 7), "row/trap control event handling")
			if event == 5 {
				q.check(t, memmap.Uint32(0x5D4594, 1049700) == 2 && memmap.Uint32(0x5D4594, 1049704) == action, "control starts two-frame pressed indicator")
				wantRow, wantTrap := uint32(0), uint32(0)
				switch action {
				case 0:
					wantRow = 4
				case 1:
					wantRow = 1
				case 2:
					wantRow = 4
				case 3:
					wantTrap = 2
				case 4:
					wantTrap = 1
				}
				q.check(t, byte(q.bar[50]) == byte(wantRow) && q.call("sub_4604E0") == wantTrap, "control changes correct row owner")
			}
			rows = append(rows, q.snapshot(fmt.Sprintf("action%d-event%d", action, event), ret))
		}
	}
	for _, initial := range []uint32{0, 1} {
		q.reset(t)
		*(*byte)(unsafe.Add(unsafe.Pointer(&q.players[0]), 2251)) = 1
		*(*uint32)(unsafe.Add(unsafe.Pointer(&q.players[0]), 3832)) = 1
		q.call("nox_xxx_quickBarCreate_45E190")
		if initial != 0 {
			q.call("sub_460920")
		}
		q.call("sub_461060")
		q.check(t, *q.quickWords["dword_5d4594_1049484"] == 1 && memmap.Uint32(0x5D4594, 1049476) == 0, "opening trap closes expanded rows")
		rows = append(rows, q.snapshot(fmt.Sprintf("expanded%d-trap-open", initial), 0))
		q.call("sub_461060")
		q.check(t, *q.quickWords["dword_5d4594_1049484"] == 0, "second trap toggle closes panel")
		rows = append(rows, q.snapshot(fmt.Sprintf("expanded%d-trap-closed", initial), 0))
	}
	spellbookCapture(t, "quickbar-row-trap-controls", rows, "599846a679536cd8b2395b67e3a7fa8758fa6923c10f1ddb79bb5f92ff666739")
}
