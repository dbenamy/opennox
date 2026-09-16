//go:build porttest

package opennox

import (
	"fmt"
	"image"
	"reflect"
	"testing"
	"unsafe"

	"github.com/opennox/libs/things"
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/server"
)

func TestQuickbarCrossRowDragging(t *testing.T) {
	q := newQuickbarOwner(t)
	var rows []quickbarResult
	for from := 0; from < 5; from++ {
		for to := 0; to < 5; to++ {
			if from == to {
				continue
			}
			for slot := 0; slot < 5; slot++ {
				q.reset(t)
				var defs []server.PortTestSpellClassDef
				for i := 0; i < 25; i++ {
					defs = append(defs, server.PortTestSpellClassDef{Index: uint32(i + 1), Flags: uint32(things.SpellClassAny), Valid: true})
					q.bar[2*i] = uint32(i + 1)
					q.bar[2*i+1] = 0xaabb0000 + uint32(i)
				}
				q.configure(defs)
				q.c.srv.Spells.EnableAll()
				*(*byte)(unsafe.Add(unsafe.Pointer(&q.players[0]), 2251)) = 1
				q.call("nox_xxx_clientUpdateButtonRow_45E110", uint32(from))
				w := (*gui.Window)(unsafe.Pointer(uintptr(q.bar[53+slot])))
				target := (slot + 2) % 5
				pos := (*gui.Window)(unsafe.Pointer(uintptr(q.bar[53+target]))).GlobalPos().Add(image.Pt(5, 5))
				packed := uint32(uint16(pos.X)) | uint32(uint16(pos.Y))<<16
				ret := q.call("nox_xxx_quickBarWnd_45EF50", uint32(uintptr(w.C())), 5, packed)
				q.check(t, ret == 1 && q.c.GUI.Captured() == w, "cross-row drag starts on original row")
				q.call("nox_xxx_clientUpdateButtonRow_45E110", uint32(to))
				want := append([]uint32(nil), q.bar[:50]...)
				a, b := 10*from+2*slot, 10*to+2*target
				want[a], want[b] = want[b], want[a]
				want[a+1], want[b+1] = want[b+1], want[a+1]
				ret = q.call("nox_xxx_quickBarWnd_45EF50", uint32(uintptr(w.C())), 7, packed)
				q.check(t, ret == 1 && reflect.DeepEqual(want, q.bar[:50]), "release swaps original-row source with newly selected destination")
				q.check(t, q.c.GUI.Captured() == nil && q.c.dragndropSpellType == 0, "cross-row release clears capture and drag")
				rows = append(rows, q.snapshot(fmt.Sprintf("from%d-to%d-source%d", from, to, slot), ret))
			}
		}
	}
	spellbookCapture(t, "quickbar-cross-row-dragging", rows, "ef9eca9f221ac464ab4fe9c7e487ff84b536e4678bb103e7dadea9fa4838993b")
}
