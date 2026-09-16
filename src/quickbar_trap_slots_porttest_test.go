//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/things"
	"github.com/opennox/opennox/v1/client/gui"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/server"
	"image"
	"testing"
	"unsafe"
)

func TestQuickbarTrapSlotEligibility(t *testing.T) {
	q := newQuickbarOwner(t)
	q.ownFullQuickbar(t)
	var rows []quickbarResult
	for _, id := range []uint32{1, 4} {
		for _, restricted := range []bool{false, true} {
			for _, quest := range []bool{false, true} {
				for _, duplicate := range []bool{false, true} {
					for _, inside := range []bool{false, true} {
						q.reset(t)
						*(*byte)(unsafe.Add(unsafe.Pointer(&q.players[0]), 2251)) = 1
						*(*uint32)(unsafe.Add(unsafe.Pointer(&q.players[0]), 3832)) = 1
						flags := uint32(things.SpellClassAny)
						if restricted {
							flags |= uint32(things.SpellNoTrap)
						}
						q.configure([]server.PortTestSpellClassDef{{Index: id, Flags: flags, Valid: true}})
						if quest {
							noxflags.SetGame(noxflags.GameModeQuest)
						}
						q.call("nox_xxx_quickBarCreate_45E190")
						q.call("sub_461060")
						trap := unsafe.Slice(memmap.PtrUint32(0x5D4594, 1047940), 64)
						if duplicate {
							trap[0] = id
						}
						trap[3] = 0x12345680
						w := (*gui.Window)(unsafe.Pointer(uintptr(trap[54])))
						q.check(t, w != nil, "actual second trap slot")
						pos := image.Pt(1, 1)
						if inside {
							pos = w.GlobalPos().Add(image.Pt(5, 5))
						}
						ret := q.call("nox_xxx_spellPutInBox_45DEB0", uint32(uintptr(unsafe.Pointer(&trap[0]))), id, uint32(pos.X), uint32(pos.Y))
						allowed := inside && !restricted && !duplicate && !(quest && id == 4)
						q.check(t, ret == uint32Bool(allowed), "trap slot rejects forbidden, duplicate and quest Blink spells")
						want := uint32(0)
						if allowed {
							want = id
						}
						q.check(t, trap[2] == want && trap[3] == 0x12345680, "trap insertion changes id without changing flag bytes")
						rows = append(rows, q.snapshot(fmt.Sprintf("id%d-restricted%v-quest%v-duplicate%v-inside%v", id, restricted, quest, duplicate, inside), ret))
					}
				}
			}
		}
	}
	spellbookCapture(t, "quickbar-trap-slots", rows, "e434b48ff462f2c2b5ced25e5feae1272e6ea68401cc1bcbb6d0396c712ef48b")
}
