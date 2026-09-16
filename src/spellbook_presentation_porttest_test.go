//go:build porttest

package opennox

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/opennox/libs/things"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

func TestSpellbookRewardPresentation(t *testing.T) {
	o, bar, configure := newSpellbookAdditionOwner(t)
	oldAbilities := o.c.srv.abilities.defs
	t.Cleanup(func() { o.c.srv.abilities.defs = oldAbilities })
	o.c.srv.abilities.defs[1] = AbilityDef{name: "Berserk", field24: 1}
	guides := unsafe.Slice(memmap.PtrUint32(0x5D4594, 740076), 41*7)
	oldGuides := append([]uint32(nil), guides...)
	t.Cleanup(func() { copy(guides, oldGuides) })
	clear(guides)
	guides[8] = 1
	timestamp := memmap.PtrUint32(0x5D4594, 1217504)
	oldTimestamp := *timestamp
	t.Cleanup(func() { *timestamp = oldTimestamp })
	var rows []struct {
		Book            spellbookResult
		Particles       [][]uint32
		RNG             [2]int
		Slot, Timestamp uint32
	}
	for _, kind := range []uint32{2, 3, 4} {
		for _, notify := range []uint32{0, 1} {
			for _, auto := range []uint32{0, 1, 2} {
				func() {
					prepareSpellbookAddition(t, o, bar, 640)
					*timestamp = 0
					configure([]server.PortTestSpellClassDef{{Index: 1, Flags: uint32(things.SpellClassAny), Valid: true}})
					class := byte(1)
					if kind == 3 {
						class = 0
					} else if kind == 4 {
						class = 2
					}
					*(*byte)(unsafe.Add(unsafe.Pointer(&o.players[0]), 2251)) = class
					*(*uint32)(unsafe.Add(unsafe.Pointer(&o.players[0]), 3700)) = 1
					particles, freeParticles := legacy.PortTestEffectsScreenParticles(768)
					defer freeParticles()
					var ret uint32
					switch kind {
					case 2:
						ret = o.bookCall("nox_xxx_netSpellRewardCli_45CFE0", 1, 1, notify, auto)
					case 3:
						ret = o.bookCall("nox_xxx_abilityReward_45D290", 1, notify, auto)
					case 4:
						ret = o.bookCall("nox_xxx_netGuideRewardCli_45D140", 1, notify)
					}
					wantParticles := 0
					if notify != 0 && (kind != 3 || auto != 0) {
						wantParticles = 600
					}
					add := notify != 0 && kind != 4 && auto == 1
					if add {
						wantParticles += 50
					}
					if len(particles()) != wantParticles || (bar[0] != 0) != add {
						t.Fatalf("kind%d notify%d auto%d: particles%d want%d slot%d", kind, notify, auto, len(particles()), wantParticles, bar[0])
					}
					if notify != 0 && o.bookWindow().GetFlags().IsHidden() {
						t.Fatal("reward must open book")
					}
					rows = append(rows, struct {
						Book            spellbookResult
						Particles       [][]uint32
						RNG             [2]int
						Slot, Timestamp uint32
					}{o.bookSnapshot(fmt.Sprintf("kind%d-notify%d-auto%d", kind, notify, auto), ret), particles(), [2]int{o.c.srv.Rand.Logic.Index(), o.c.srv.Rand.Other.Index()}, bar[0], *timestamp})
				}()
			}
		}
	}
	spellbookCapture(t, "reward-presentation", rows, "9a8999bbd4fafd2f203f7ba506ae1cf62e8a11fa3fb4198feb2c4b9384670134")
}
