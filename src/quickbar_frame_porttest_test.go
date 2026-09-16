//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/things"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

func TestQuickbarGUIFrames(t *testing.T) {
	q := newQuickbarOwner(t)
	q.ownFullQuickbar(t)
	var rows []quickbarResult
	for class := uint32(0); class < 3; class++ {
		for _, expanded := range []bool{false, true} {
			q.reset(t)
			*(*byte)(unsafe.Add(unsafe.Pointer(&q.players[0]), 2251)) = byte(class)
			q.configure([]server.PortTestSpellClassDef{{Index: 1, Flags: uint32(things.SpellClassAny), Valid: true}})
			q.c.srv.Spells.EnableAll()
			q.call("nox_xxx_quickBarCreate_45E190")
			if expanded && class != 0 {
				q.call("sub_460920")
			}
			for frame := 0; frame < 3; frame++ {
				clear(q.pix.Pix)
				blank := effectsPixelHash(q.pix)
				q.c.GUI.Draw()
				q.check(t, effectsPixelHash(q.pix) != blank, "complete GUI frame draws actual quickbar windows")
				rows = append(rows, q.snapshot(fmt.Sprintf("class%d-expanded%v-frame%d", class, expanded, frame), 0))
			}
		}
	}
	spellbookCapture(t, "quickbar-gui-frames", rows, "9ad38bc8097c13398dec76466c080c6d2a73266c06b44657aa240eae71ecf60c")
}
