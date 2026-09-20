//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"testing"

	"github.com/opennox/libs/strman"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
)

func TestGameMessageNoticeAwards(t *testing.T) {
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
		for kind, notice := range []byte{2, 0} {
			clear(region)
			*index = initial
			printer.lines = nil
			want := [3]string{}
			for _, id := range []uint32{0, 5, 2, 1, 3, 4, 2} {
				previous := *index
				data := []byte{169, notice, 0, 0, 0, 0}
				binary.LittleEndian.PutUint32(data[2:], id)
				before := bytes.Clone(data)
				if n := legacy.PortTestGameNotice(data); n != 6 || !bytes.Equal(data, before) {
					t.Fatal("award notice length/input", n)
				}
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
	gameMessageCapture(t, "game-notice-awards", rows)
}
