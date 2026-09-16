//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/spell"
	"github.com/opennox/libs/things"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

func TestQuickbarManaShading(t *testing.T) {
	q := newQuickbarOwner(t)
	q.ownFullQuickbar(t)
	var rows []quickbarResult
	for _, cost := range []int{0, 20} {
		for _, mana := range []uint32{0, 1, 9, 19, 20, 21} {
			for _, blocked := range []int{0, 1, 2, 3} {
				q.reset(t)
				q.configure([]server.PortTestSpellClassDef{{Index: 1, Flags: uint32(things.SpellClassAny), Valid: true}})
				sd := q.c.srv.Spells.DefByInd(spell.ID(1))
				sd.Enabled = true
				sd.Title = "Mana shading"
				sd.Icon = unsafe.Pointer(q.images[1])
				q.bar[0] = 1
				*memmap.PtrPtr(0x852978, 8) = q.items[0].C()
				*(*uint32)(unsafe.Add(q.items[0].C(), 276)) = 0
				*(*uint32)(unsafe.Add(q.items[0].C(), 124)) = 0
				q.call("nox_xxx_quickBarWarriorDraw_45FDE0", q.bar[53])
				plain := effectsPixelHash(q.pix)
				clear(q.pix.Pix)
				sd.Def.ManaCost = cost
				q.meters.Records[1].Current = mana
				switch blocked {
				case 1:
					sd.Enabled = false
				case 2:
					*(*uint32)(unsafe.Add(q.items[0].C(), 276)) = 1
				case 3:
					*(*uint32)(unsafe.Add(q.items[0].C(), 124)) = 1 << 29
				}
				ret := q.call("nox_xxx_quickBarWarriorDraw_45FDE0", q.bar[53])
				shaded := blocked != 0 || (cost > 0 && mana < uint32(cost))
				q.check(t, ret == 1 && (effectsPixelHash(q.pix) != plain) == shaded, "mana deficit, animation, disabled spell and buff shading")
				rows = append(rows, q.snapshot(fmt.Sprintf("cost%d-mana%d-blocked%d", cost, mana, blocked), ret))
			}
		}
	}
	spellbookCapture(t, "quickbar-mana-shading", rows, "fc09adc9000339d24f0e25eebc9a68a8aa5e8d39c2a83dc08eb5d0abbf67e052")
}
