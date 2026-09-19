//go:build porttest

package opennox

import (
	flags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

func TestBookAwardsQuestSingleLevelSpells(t *testing.T) {
	o := newReliableReportsOwner(t)
	t.Cleanup(flags.PortTestGameFlags(0))
	bookAwardStrings(t, o.s)
	configure, restore := o.s.PortTestAISpellDefs()
	t.Cleanup(restore)
	t.Cleanup(o.s.PortTestBookSpellOwner())
	var defs []server.PortTestSpellClassDef
	for id := uint32(1); id <= 136; id++ {
		defs = append(defs, server.PortTestSpellClassDef{Index: id, Flags: 0x100, Valid: true})
	}
	configure(defs)
	u := &o.units[0]
	p := u.UpdateDataPlayer().Player
	ids, reset, records, free := legacy.PortTestPlayerFileRecords()
	t.Cleanup(free)
	p.Prot4636 = ids[0]
	single := map[uint32]bool{19: true, 34: true, 45: true, 46: true, 47: true, 48: true, 49: true, 117: true, 118: true, 119: true, 120: true, 121: true, 122: true, 123: true, 124: true, 125: true, 134: true}
	var rows []map[string]any
	for _, gf := range []uint32{0, 2048, 4096, 6144} {
		flags.ResetGame()
		flags.SetGame(flags.GameFlag(gf))
		for id := uint32(1); id <= 136; id++ {
			o.reset()
			reset()
			clear(p.SpellLvl[:])
			p.SpellLvl[id] = 1
			ret := bookAwardCall("nox_xxx_spellGrantToPlayer_4FB550", uint32(uintptr(unsafe.Pointer(u))), id, 0, 0, 0)
			want, level, bits := uint32(1), uint32(2), uint32(1)<<(id&31)
			if gf&4096 != 0 && single[id] {
				want = 0
				level = 1
				bits = 0
			}
			if ret != want || p.SpellLvl[id] != level || records() != [3]uint32{bits, 0, 0x2468ace0 ^ bits} {
				t.Fatal("quest single-level", gf, id, ret, want, p.SpellLvl[id], level, records())
			}
			message := ""
			if want == 0 {
				message = "Maximum level"
			}
			bookAwardCheckDirect(t, o, message, false)
			rows = append(rows, map[string]any{"flags": gf, "id": id, "return": ret, "level": level, "records": records(), "queue": o.state(), "direct": visibilityEffectsPackets(o.s)})
		}
	}
	spellbookCapture(t, "book-awards-quest-single-level", rows, "b263803b962140d736d45963ca32bdd3ad92aa492a9d11c15e0758c8dc0f9d36")
}
func TestBookAwardsSpellNotificationGates(t *testing.T) {
	o := newReliableReportsOwner(t)
	t.Cleanup(flags.PortTestGameFlags(0))
	bookAwardStrings(t, o.s)
	configure, restore := o.s.PortTestAISpellDefs()
	t.Cleanup(restore)
	t.Cleanup(o.s.PortTestBookSpellOwner())
	u := &o.units[0]
	p := u.UpdateDataPlayer().Player
	ids, reset, _, free := legacy.PortTestPlayerFileRecords()
	t.Cleanup(free)
	p.Prot4636 = ids[0]
	p.Field4792 = 1
	t.Cleanup(o.s.PortTestCombatAudioReset)
	var rows []map[string]any
	for _, gf := range []uint32{0, 2048, 4096, 6144} {
		flags.ResetGame()
		flags.SetGame(flags.GameFlag(gf))
		for _, id := range []uint32{1, 34} {
			for _, sf := range []uint32{0x100, 0x1000, 0x4000, 0x10000, 0x15000} {
				configure([]server.PortTestSpellClassDef{{Index: id, Flags: sf, Valid: true}})
				o.reset()
				reset()
				o.s.PortTestCombatAudioReset()
				clear(p.SpellLvl[:])
				ret := bookAwardCall("nox_xxx_spellGrantToPlayer_4FB550", uint32(uintptr(unsafe.Pointer(u))), id, 1, 0, 0)
				state := o.state()
				want := 4
				if gf&2048 != 0 && (id == 34 || sf&0x15000 != 0) {
					want = 1
				}
				if ret != 1 || len(state.Nodes) != want || len(o.s.PortTestCombatAudioSnapshot()) != 1 {
					t.Fatal("spell notify gate", gf, id, sf, ret, len(state.Nodes), want)
				}
				rows = append(rows, map[string]any{"flags": gf, "id": id, "spellFlags": sf, "return": ret, "queue": state, "direct": visibilityEffectsPackets(o.s)})
			}
		}
	}
	spellbookCapture(t, "book-awards-spell-notification-gates", rows, "e7a22b1cdc1b1b3277bcd211e402b1ee2aaef6aaa607c97e6cee62232df16dca")
}
