//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/things"
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/server"
	"image"
	"testing"
	"unsafe"
)

func TestQuickbarSlotEvents(t *testing.T) {
	q := newQuickbarOwner(t)
	var rows []quickbarResult
	for source := 0; source < 5; source++ {
		for target := 0; target < 6; target++ {
			for _, event := range []uint32{6, 7, 8, 12, 16} {
				q.reset(t)
				defs := make([]server.PortTestSpellClassDef, 5)
				for i := 0; i < 5; i++ {
					defs[i] = server.PortTestSpellClassDef{Index: uint32(i + 1), Flags: uint32(things.SpellClassAny), Valid: true}
					q.bar[2*i] = uint32(i + 1)
					q.bar[2*i+1] = 0x76543200 | uint32(i&1)
				}
				q.configure(defs)
				q.c.srv.Server.Spells.EnableAll()
				*(*byte)(unsafe.Add(unsafe.Pointer(&q.players[0]), 2251)) = 1
				w := (*gui.Window)(unsafe.Pointer(uintptr(q.bar[53+source])))
				pos := image.Pt(600, 300)
				if target < 5 {
					pos = (*gui.Window)(unsafe.Pointer(uintptr(q.bar[53+target]))).GlobalPos().Add(image.Pt(5, 5))
				}
				packed := uint32(uint16(pos.X)) | uint32(uint16(pos.Y))<<16
				ret := q.call("nox_xxx_quickBarWnd_45EF50", uint32(uintptr(w.C())), 5, packed)
				q.check(t, ret == 1 && q.c.GUI.Captured() == w && q.c.dragndropSpellType == 2, "slot press captures actual window")
				before := append([]uint32(nil), q.bar[:50]...)
				ret = q.call("nox_xxx_quickBarWnd_45EF50", uint32(uintptr(w.C())), event, packed)
				if event == 6 || event == 7 {
					q.check(t, ret == 1 && q.c.dragndropSpellType == 0, "slot release ends drag")
					if target < 5 && target != source {
						q.check(t, q.bar[2*source] == before[2*target] && q.bar[2*target] == before[2*source], "same-row drag swaps ids")
						q.check(t, q.bar[2*source+1] == before[2*target+1] && q.bar[2*target+1] == before[2*source+1], "same-row drag swaps full flags")
					}
					if target == 5 {
						q.check(t, q.bar[2*source] == 0 && q.bar[2*source+1] == 0x76543200, "outside drop clears id and low flag byte")
					}
				} else {
					q.check(t, ret == 0 && q.c.dragndropSpellType == 2 && q.c.GUI.Captured() == w, "unhandled event preserves drag")
				}
				rows = append(rows, q.snapshot(fmt.Sprintf("source%d-target%d-event%d", source, target, event), ret))
			}
		}
	}
	spellbookCapture(t, "quickbar-slot-events", rows, "b825d46dbc297d029353da80d9adcd293298693edbd2b0ad1601e8d04d74d63b")
}
