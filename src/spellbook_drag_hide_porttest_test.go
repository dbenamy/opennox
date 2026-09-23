//go:build porttest

package opennox

import "testing"

// TestSpellbookDragHideContract exercises the book's real hide/export path and
// observes the actual Client drag state rather than replacing callback funcs.
func TestSpellbookDragHideContract(t *testing.T) {
	cases := []struct {
		name string
		typ  int
	}{
		{name: "none", typ: 0},
		{name: "book", typ: 1},
		{name: "quickbar", typ: 2},
		{name: "negative", typ: -1},
		{name: "min-int32", typ: int(-2147483648)},
		{name: "max-int32", typ: int(2147483647)},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			o := newSpellbookOwner(t)
			o.resetBook(t)
			if got := o.bookCall("nox_xxx_bookInit_45B9D0"); got != 1 {
				t.Fatalf("init returned %d", got)
			}
			o.bookCall("nox_xxx_book_45B010", 1)
			if o.bookCall("sub_45CFC0") != 1 || o.bookWindow().GetFlags().IsHidden() {
				t.Fatal("fixture did not show book")
			}
			if !o.bookWindow().Capture(true) {
				t.Fatal("fixture could not capture book")
			}
			const spellID = uint32(0x12345678)
			o.c.dragndropSpellSet(spellID, tc.typ)

			if got := o.bookCall("nox_xxx_bookHideMB_45ACA0", 0); got != 1 {
				t.Fatalf("hide returned %d", got)
			}
			if o.bookCall("sub_45CFC0") != 0 || o.c.GUI.Captured() != nil {
				t.Fatal("hide did not hide the book and release its capture")
			}
			wantID, wantType := spellID, tc.typ
			if tc.typ == 1 {
				wantID, wantType = 0, 0
			}
			if o.c.dragndrapSpell != wantID || o.c.dragndropSpellType != wantType {
				t.Fatalf("drag after hide = (%#x,%d), want (%#x,%d)", o.c.dragndrapSpell, o.c.dragndropSpellType, wantID, wantType)
			}

			// A repeated hide exits at the already-hidden guard before reading or
			// clearing the drag state.
			o.c.dragndropSpellSet(spellID, 1)
			if got := o.bookCall("nox_xxx_bookHideMB_45ACA0", 0); got != 0 {
				t.Fatalf("already-hidden hide returned %d", got)
			}
			if o.c.dragndrapSpell != spellID || o.c.dragndropSpellType != 1 {
				t.Fatal("already-hidden guard changed drag state")
			}
		})
	}

	// The active-addition guard also exits before the type getter/conditional
	// clear, while retaining the book capture.
	o := newSpellbookOwner(t)
	o.resetBook(t)
	if got := o.bookCall("nox_xxx_bookInit_45B9D0"); got != 1 {
		t.Fatalf("blocked case init returned %d", got)
	}
	o.bookCall("nox_xxx_book_45B010", 1)
	if !o.bookWindow().Capture(true) {
		t.Fatal("blocked case could not capture book")
	}
	const spellID = uint32(0x87654321)
	o.c.dragndropSpellSet(spellID, 1)
	*o.words["dword_5d4594_1047520"] = 1
	if got := o.bookCall("nox_xxx_bookHideMB_45ACA0", 0); got != 0 {
		t.Fatalf("blocked hide returned %d", got)
	}
	if o.c.dragndrapSpell != spellID || o.c.dragndropSpellType != 1 || o.c.GUI.Captured() != o.bookWindow() {
		t.Fatal("blocked hide changed drag state or capture")
	}
}
