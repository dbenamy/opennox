//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
	"unsafe"
)

func TestClientListboxMiddleInsertionContract(t *testing.T) {
	o := newListboxOwner(t)
	// Sixteen real rows keep the original mis-scaled source/destination addresses
	// inside this allocation: shifting indices1/2 incorrectly touches rows4/5/8/9.
	// This reproduces the order error without an out-of-range memory access.
	o.create(t, 8, 70, 65, func(_ *gui.WindowData, d *gui.ScrollListBoxData) { d.Count = 16 })
	for _, s := range []string{"one", "two", "three"} {
		v := []uint16{}
		for _, c := range s {
			v = append(v, uint16(c))
		}
		if o.textEvent(16397, v, -1) != 1 {
			t.Fatal("append failed")
		}
	}
	o.event(16402, 1, 0)
	if o.textEvent(16397, []uint16{'X'}, -1) != 1 {
		t.Fatal("middle insertion failed")
	}
	d := o.data()
	if d.Field_11_0 != 4 {
		t.Fatal("insert row count")
	}
	rows := unsafe.Slice(d.Items, int(d.Count))
	for i, want := range []string{"one", "X", "two", "three"} {
		got := ""
		for _, v := range rows[i].Text {
			if v == 0 {
				break
			}
			got += string(rune(v))
		}
		if got != want {
			t.Fatalf("row %d after middle insertion: %q, want %q", i, got, want)
		}
	}
}

func TestClientListboxTextBoundsContract(t *testing.T) {
	o := newListboxOwner(t)
	o.create(t, 8, 70, 65, nil)
	o.textEvent(16397, []uint16{'A'}, -1)
	o.textEvent(16397, []uint16{'B'}, -1)
	rows := unsafe.Slice(o.data().Items, int(o.data().Count))
	for i := range rows[1].Text {
		rows[1].Text[i] = 0x7777
	}
	rows[1].Text[255] = 0
	before := rows[1]
	text := make([]uint16, 400)
	for i := range text {
		text[i] = 0xd800 + uint16(i%32)
	}
	o.textEvent(16407, text, 0)
	if rows[0].Text[255] != 0 || rows[1] != before {
		t.Fatal("row replacement exceeded text storage")
	}
	for i := 0; i < 255; i++ {
		if rows[0].Text[i] != text[i] {
			t.Fatal("row replacement lost raw code units")
		}
	}
	words := unsafe.Slice((*uint32)(o.win.C()), 101)
	prior := append([]uint32(nil), words...)
	o.textEvent(16385, text, 0)
	label := unsafe.Slice((*uint16)(unsafe.Add(o.win.DrawData().C(), 72)), 64)
	if label[63] != 0 {
		t.Fatal("listbox label terminator")
	}
	for i := 0; i < 63; i++ {
		if label[i] != text[i] {
			t.Fatal("label lost raw code units")
		}
	}
	for i, v := range words {
		if (i < 27 || i >= 59) && v != prior[i] {
			t.Fatalf("label setter changed unrelated window word %d", i)
		}
	}
}
func TestClientListboxSelectionCapacityContract(t *testing.T) {
	o := newListboxOwner(t)
	for _, capacity := range []uint16{1, 2, 8} {
		o.create(t, 8, 70, 65, func(_ *gui.WindowData, d *gui.ScrollListBoxData) { d.Count = capacity; d.Field_4 = 1 })
		for i := 0; i < int(capacity); i++ {
			o.textEvent(16397, []uint16{uint16('A' + i)}, -1)
			o.event(16405, uint32(i), 0)
		}
		got := o.selection()
		if len(got) != int(capacity)+1 || got[capacity] != -1 {
			t.Fatal("full selection lacks sentinel")
		}
		for i := 0; i < int(capacity); i++ {
			if got[i] != int32(i) {
				t.Fatal("full selection order")
			}
		}
		o.event(16405, 0, 0)
		got = o.selection()
		for i := 0; i < int(capacity)-1; i++ {
			if got[i] != int32(i+1) {
				t.Fatal("selection removal shifted incorrectly")
			}
		}
		if got[capacity-1] != -1 || got[capacity] != -1 {
			t.Fatal("selection removal lost sentinel")
		}
	}
}
func TestClientListboxScrollBoundsContract(t *testing.T) {
	o := newListboxOwner(t)
	o.create(t, 8, 70, 65, func(_ *gui.WindowData, d *gui.ScrollListBoxData) { d.Count = 2 })
	for _, v := range []uint16{'A', 'B'} {
		o.textEvent(16397, []uint16{v}, -1)
	}
	d := o.data()
	d.Field_13_1 = 1000
	if legacy.PortTestListboxScrollIndex(d) != 0 {
		t.Fatal("past-end full list should use original fallback index")
	}
	rows := append([]gui.ScrollListBoxItem(nil), unsafe.Slice(d.Items, int(d.Count))...)
	count := d.Field_11_0
	if o.event(16411, 0xffffffff, 0) != 0 || d.Field_11_0 != count {
		t.Fatal("negative head removal was not rejected")
	}
	for i, v := range unsafe.Slice(d.Items, int(d.Count)) {
		if v != rows[i] {
			t.Fatal("negative head removal changed rows")
		}
	}
	// One wrapped row taller than the viewport cannot be made to fit by scrolling.
	// The callback must return once scrolling stops making progress.
	o.create(t, 8, 30, 20, func(_ *gui.WindowData, d *gui.ScrollListBoxData) { d.Field_1 = 1 })
	text := []uint16{}
	for i := 0; i < 50; i++ {
		text = append(text, 'A')
	}
	if o.textEvent(16397, text, -1) != 1 || o.data().Field_11_0 != 1 {
		t.Fatal("oversized auto-scroll row was not added")
	}
}
func TestClientListboxLayoutContract(t *testing.T) {
	var d gui.ScrollListBoxData
	var row gui.ScrollListBoxItem
	if unsafe.Sizeof(d) != 56 || unsafe.Offsetof(d.Field_12) != 48 || unsafe.Sizeof(row) != 524 || unsafe.Offsetof(row.Text) != 4 {
		t.Fatal("listbox C layout changed")
	}
}

func TestClientListboxNarrowDrawContract(t *testing.T) {
	o := newListboxOwner(t)
	for mode := 0; mode < 2; mode++ {
		o.resetRender(uint32(mode+1), 120)
		o.create(t, gui.StatusFlags(8|0x4000|mode*128), 10, 65, func(_ *gui.WindowData, d *gui.ScrollListBoxData) { d.Field_3 = 1 })
		o.textEvent(16397, []uint16{'a', 'b', 'c'}, -1)
		o.win.Draw()
		if o.data().Items.Text[0] != 'a' || o.data().Items.Text[1] != 'b' || o.data().Items.Text[2] != 'c' {
			t.Fatal("display clipping changed stored text")
		}
	}
}

func TestClientListboxFullSelectionEmptyClickContract(t *testing.T) {
	o := newListboxOwner(t)
	for _, capacity := range []uint16{1, 2} {
		o.create(t, 8, 70, 65, func(_ *gui.WindowData, d *gui.ScrollListBoxData) { d.Count = capacity; d.Field_4 = 1 })
		for i := 0; i < int(capacity); i++ {
			o.textEvent(16397, []uint16{uint16('A' + i)}, -1)
			o.event(16405, uint32(i), 0)
		}
		before := o.selection()
		o.notices = nil
		for _, code := range []int{6, 7} {
			if gui.EventRespInt(o.win.Func93(&gui.RawEvent{Event: code, Arg1: 20 | 70<<16})) != 1 {
				t.Fatal("empty-area click result")
			}
			after := o.selection()
			for i, v := range before {
				if after[i] != v {
					t.Fatal("empty-area click changed full selection")
				}
			}
		}
		if len(o.notices) != 2 || o.notices[0].Code != 16400 || o.notices[0].B != 0xffffffff {
			t.Fatal("empty-area click notification")
		}
	}
}
