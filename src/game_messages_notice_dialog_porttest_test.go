//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"slices"
	"testing"

	"github.com/opennox/libs/spell"
	"github.com/opennox/libs/strman"
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

func TestGameMessageNoticeSpellAndDialog(t *testing.T) {
	o := newInventoryTransactionOwner(t)
	o.reset(t)
	printer, index, region := teamRuntimeJoinTextOwner(t, o.c.srv.Server)
	set, restore := o.c.srv.PortTestMeterStrings(
		strman.Entry{ID: "guimsg.c:systemmsg", Vals: []strman.Variant{{Str: "System: %s"}}},
		strman.Entry{ID: "plyrspel.c:SpellCastSuccess", Vals: []strman.Variant{{Str: "Cast %s"}}},
		strman.Entry{ID: "Noxworld.c:ErrChangedClass", Vals: []strman.Variant{{Str: "Changed Ω"}}},
		strman.Entry{ID: "guiserv.c:Notice", Vals: []strman.Variant{{Str: "Notice é"}}},
	)
	defer restore()
	set(0)
	defs, restoreDefs := o.c.srv.PortTestAISpellDefs()
	defer restoreDefs()
	defs([]server.PortTestSpellClassDef{{Index: 1, Valid: true}, {Index: 2, Valid: true}})
	titles := []string{"", "Flame Ω", "Frost é"}
	for i := 1; i < 3; i++ {
		o.c.srv.Spells.DefByInd(spell.ID(i)).Title = titles[i]
	}
	var rows []map[string]any
	for _, id := range []uint32{1, 2, 1} {
		data := []byte{169, 1, 0, 0, 0, 0}
		binary.LittleEndian.PutUint32(data[2:], id)
		before := bytes.Clone(data)
		clear(region)
		*index = 2
		printer.lines = nil
		o.sounds = nil
		n := legacy.PortTestGameNotice(data)
		if n != 6 || !bytes.Equal(data, before) || !slices.Equal(printer.lines, []string{"System: Cast " + titles[id]}) || *index != 0 || len(o.sounds) != 0 {
			t.Fatal("spell notice", id, n, printer.lines)
		}
		rows = append(rows, map[string]any{"id": id, "return": n, "text": slices.Clone(printer.lines), "ring": bytes.Clone(region)})
	}
	oldDialog := legacy.Nox_xxx_dialogMsgBoxCreate_449A10
	defer func() { legacy.Nox_xxx_dialogMsgBoxCreate_449A10 = oldDialog }()
	calls := 0
	legacy.Nox_xxx_dialogMsgBoxCreate_449A10 = func(w *gui.Window, title, text string, flags gui.DialogFlags, a, b func()) {
		calls++
		if w != nil || title != "Notice é" || text != "Changed Ω" || flags != 33 || a != nil || b != nil {
			t.Fatal("class-change dialog contract", title, text, flags)
		}
	}
	for i := 0; i < 3; i++ {
		data := []byte{169, 17}
		before := bytes.Clone(data)
		prior := bytes.Clone(region)
		printer.lines = nil
		o.sounds = nil
		n := legacy.PortTestGameNotice(data)
		if n != 2 || calls != i+1 || !bytes.Equal(data, before) || !bytes.Equal(region, prior) || len(printer.lines) != 0 || len(o.sounds) != 0 {
			t.Fatal("class-change notice return/effects", n, calls)
		}
		rows = append(rows, map[string]any{"dialog": i, "return": n, "calls": calls})
	}
	gameMessageCapture(t, "game-notice-spell-dialog", rows)
}
