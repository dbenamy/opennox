//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/spell"
	"github.com/opennox/libs/things"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

func TestQuickbarSlotRendering(t *testing.T) {
	q := newQuickbarOwner(t)
	q.ownFullQuickbar(t)
	meters := legacy.PortTestNewMeterEnvironment()
	t.Cleanup(meters.Restore)
	abilities := q.c.srv.abilities.defs
	t.Cleanup(func() { q.c.srv.abilities.defs = abilities })
	frame := q.c.srv.Frame()
	t.Cleanup(func() { q.c.srv.SetFrame(frame) })
	var rows []quickbarResult
	for _, class := range []uint32{0, 1} {
		for _, icon := range []bool{false, true} {
			for _, enabled := range []bool{false, true} {
				for _, since := range []uint32{0, 1, 9, 10, 40, 1360} {
					for _, timer := range []uint32{0, 1, 127, 128, 255} {
						q.reset(t)
						q.configure([]server.PortTestSpellClassDef{{Index: 1, Flags: uint32(things.SpellClassAny), Valid: true}})
						sd := q.c.srv.Spells.DefByInd(spell.ID(1))
						sd.Title = "Quickbar spell"
						sd.Def.ManaCost = 20
						sd.Enabled = enabled
						q.c.srv.abilities.defs[1] = AbilityDef{name: "Quickbar ability", field24: 1, delay: 1200}
						if icon {
							sd.Icon = unsafe.Pointer(q.images[1])
							sd.IconEnabled = unsafe.Pointer(q.images[2])
							q.c.srv.abilities.defs[1].icon8 = q.images[1]
							q.c.srv.abilities.defs[1].icon12 = q.images[2]
						}
						*(*byte)(unsafe.Add(unsafe.Pointer(&q.players[0]), 2251)) = byte(class)
						q.bar[0] = 1
						q.c.srv.SetFrame(2000)
						*memmap.PtrUint32(0x587000, 133484) = 0
						*memmap.PtrUint32(0x5D4594, 1049540) = 2000 - since
						*memmap.PtrUint32(0x5D4594, 1047764+24+8) = uint32Bool(enabled)
						*memmap.PtrUint32(0x5D4594, 1047764+24+20) = 2000 - since
						*memmap.PtrUint8(0x5D4594, 1049545) = byte(timer)
						meters.Records[1].Current = 10
						blank := effectsPixelHash(q.pix)
						op := "nox_xxx_quickBarDrawFn_45FBD0"
						if class != 0 {
							op = "nox_xxx_quickBarWarriorDraw_45FDE0"
						}
						ret := q.call(op, q.bar[53])
						q.check(t, ret == 1 && effectsPixelHash(q.pix) != blank, "populated slot renders pixels or fallback text")
						wantTimer := byte(timer)
						if class != 0 && timer > 0 && timer < 128 {
							wantTimer--
						}
						q.check(t, *memmap.PtrUint8(0x5D4594, 1049545) == wantTimer, "spell flash timer uses signed byte decrement")
						rows = append(rows, q.snapshot(fmt.Sprintf("class%d-icon%v-enabled%v-since%d-timer%d", class, icon, enabled, since, timer), ret))
					}
				}
			}
		}
	}
	spellbookCapture(t, "quickbar-slot-rendering", rows, "c60518093674943097f94b5fb7a3f42e8e9a64e33ec6587ae8e1dbc3ab2ddab1")
}
