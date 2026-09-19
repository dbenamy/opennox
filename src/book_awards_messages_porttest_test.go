//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
	"testing"
	"unsafe"

	"github.com/opennox/libs/strman"
	flags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

func TestBookAwardsClientMessages(t *testing.T) {
	o := newReliableReportsOwner(t)
	printer, index, region := teamRuntimeJoinTextOwner(t, o.s)
	success := bookAwardWords(t, 0x587000, 217092, 6)
	errors := bookAwardWords(t, 0x587000, 216380, 6)
	var entries []strman.Entry
	entries = append(entries, strman.Entry{ID: "guimsg.c:systemmsg", Vals: []strman.Variant{{Str: "System: %s"}}})
	texts := []string{"", "Hello", "Café Ω", "Ability learned", "Please retry", "The end"}
	for i, s := range texts {
		key := fmt.Sprintf("Result%d", i)
		success[i] = bookAwardString(t, key)
		errors[i] = bookAwardString(t, key)
		for _, file := range []string{"Ability.c:", "plyrspel.c:"} {
			entries = append(entries, strman.Entry{ID: strman.ID(file + key), Vals: []strman.Variant{{Str: s}}})
		}
	}
	configure, restore := o.s.PortTestMeterStrings(entries...)
	t.Cleanup(restore)
	configure(0)
	var rows []map[string]any
	for _, initial := range []uint32{0, 1, 2} {
		for kind, op := range []string{"nox_xxx_abilGetSuccess_4FB960_ability", "nox_xxx_abilGetError_4FB0B0_magic_plyrspel"} {
			clear(region)
			*index = initial
			printer.lines = nil
			want := [3]string{}
			for _, id := range []uint32{0, 5, 2, 1, 3, 4, 2} {
				previous := *index
				bookAwardCall(op, id)
				at := (previous + 1) % 3
				want[at] = texts[id]
				if *index != at {
					t.Fatal("message ring", *index, at)
				}
				var got [3]string
				for i := range got {
					got[i] = legacy.GoWStringP(memmap.PtrOff(0x5D4594, 823804+644*uintptr(i)))
				}
				if got != want || binary.LittleEndian.Uint32(region[644*at+636:]) != 273 || region[644*at+640] != 0 {
					t.Fatal("centered message", got, want)
				}
				if last := printer.lines[len(printer.lines)-1]; last != "System: "+texts[id] {
					t.Fatal("console", last, texts[id])
				}
				rows = append(rows, map[string]any{"initial": initial, "kind": kind, "id": id, "head": *index, "text": got, "region": bytes.Clone(region), "console": append([]string(nil), printer.lines...)})
			}
		}
	}
	spellbookCapture(t, "book-awards-client-messages", rows, "7887f0dc5c68830813a51897ff64bfae0078add1736f404c6def4efe4e49bb10")
}
func TestBookAwardsShopClosure(t *testing.T) {
	o := newReliableReportsOwner(t)
	t.Cleanup(flags.PortTestGameFlags(6144))
	bookAwardStrings(t, o.s)
	configure, restore := o.s.PortTestAISpellDefs()
	t.Cleanup(restore)
	t.Cleanup(o.s.PortTestBookSpellOwner())
	configure([]server.PortTestSpellClassDef{{Index: 1, Flags: 0x100, Valid: true}})
	u := &o.units[0]
	p := u.UpdateDataPlayer().Player
	ids, resetRecords, _, free := legacy.PortTestPlayerFileRecords()
	t.Cleanup(free)
	p.Prot4636 = ids[0]
	p.Field4792 = 0
	t.Cleanup(o.s.PortTestCombatAudioReset)
	var rows []map[string]any
	for _, notify := range []uint32{0, 1, 2} {
		for _, auto := range []uint32{0, 1, 2} {
			o.reset()
			resetRecords()
			o.s.PortTestCombatAudioReset()
			clear(p.SpellLvl[:])
			snapshot, release := legacy.PortTestBookAwardShop(u, &o.units[1])
			t.Cleanup(release)
			got := bookAwardCall("nox_xxx_spellGrantToPlayer_4FB550", uint32(uintptr(unsafe.Pointer(u))), 1, notify, auto, 0)
			closed := notify == 1 && auto == 1
			want := [4]uint32{1, 1, 1, 0}
			if closed {
				want = [4]uint32{0, 0, 0, 1}
			}
			state := snapshot()
			queue := o.state()
			if got != 1 || state != want {
				t.Fatal("shop closure", notify, auto, got, state, want)
			}
			closePackets := 0
			for _, n := range queue.Nodes {
				if reflect.DeepEqual(n.Data, []byte{201, 2}) {
					closePackets++
				}
			}
			expected := 0
			if closed {
				expected = 1
			}
			if closePackets != expected {
				t.Fatal("shop close packet", notify, auto, closePackets, expected, queue.Nodes)
			}
			rows = append(rows, map[string]any{"notify": notify, "auto": auto, "return": got, "state": state, "queue": queue, "direct": visibilityEffectsPackets(o.s)})
		}
	}
	spellbookCapture(t, "book-awards-shop-closure", rows, "8b1a99658f3b7d10a968b7f1a208298c1b8a6d02c5f2049ea59e40dafd84c2d2")
}
