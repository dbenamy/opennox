//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/opennox/v1/common/memmap"
	"testing"
	"unsafe"
)

func TestSpellbookAbilityRemoval(t *testing.T) {
	o, bar, _ := newSpellbookAdditionOwner(t)
	oldAbilities := o.c.srv.abilities.defs
	t.Cleanup(func() { o.c.srv.abilities.defs = oldAbilities })
	for id := 1; id <= 5; id++ {
		o.c.srv.abilities.defs[id] = AbilityDef{name: fmt.Sprint("Ability", id), field24: 1}
	}
	records := unsafe.Slice(memmap.PtrUint32(0x5D4594, 1047788), 30)
	oldRecords := append([]uint32(nil), records...)
	t.Cleanup(func() { copy(records, oldRecords) })
	var rows []struct {
		Book             spellbookResult
		Slots, Abilities []uint32
	}
	for _, id := range []uint32{1, 3, 5} {
		for _, rank := range []uint32{0, 1, 7} {
			prepareSpellbookAddition(t, o, bar, 640)
			clear(records)
			for i := 1; i <= 5; i++ {
				records[(i-1)*6] = uint32(i)
				records[(i-1)*6+4] = rank
				*(*uint32)(unsafe.Add(unsafe.Pointer(&o.players[0]), 3696+4*i)) = rank
			}
			for i := 0; i < 25; i++ {
				bar[2*i] = uint32(i%5 + 1)
				bar[2*i+1] = 0xa1000000 + uint32(i)
			}
			ret := o.bookCall("nox_xxx_clientQuestDisableAbility_45D4A0", id)
			for i := 1; i <= 5; i++ {
				want := rank
				if uint32(i) == id {
					want = 0
				}
				if records[(i-1)*6+4] != want || *(*uint32)(unsafe.Add(unsafe.Pointer(&o.players[0]), 3696+4*i)) != want {
					t.Fatal("ability state and knowledge must be cleared together")
				}
			}
			for i := 0; i < 25; i++ {
				want := uint32(i%5 + 1)
				if want == id {
					want = 0
				}
				if bar[2*i] != want || bar[2*i+1] != 0xa1000000+uint32(i) {
					t.Fatal("ability removal must clear duplicate slots and retain flags")
				}
			}
			rows = append(rows, struct {
				Book             spellbookResult
				Slots, Abilities []uint32
			}{o.bookSnapshot(fmt.Sprintf("ability%d-rank%d", id, rank), ret), append([]uint32(nil), bar[:50]...), append([]uint32(nil), records...)})
		}
	}
	spellbookCapture(t, "ability-removal", rows, "faa69b3a9daf912bb23c7ab25b32172ccc23a7079ee1589bb7a3503dc8af46af")
}
