//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/client/gui"
	"testing"
)

func TestClientEntryFilterPrecedence(t *testing.T) {
	o := newEntryOwner(t)
	d := o.create(t, 0, 8, func(_ *gui.WindowData, d *gui.EntryFieldData) {
		d.Field_1028 = 1
		d.Field_1032 = 1
	})
	o.c.GUI.Focus(o.win)
	// With both filters active, a letter must be rejected while a digit is kept.
	if o.key(30, 2) != 1 || d.Text[0] != 0 || d.Field_1052 != 0 {
		t.Fatal("digit filter did not take precedence")
	}
	if o.key(2, 1) != 1 || d.Text[0] != 0 || d.Field_1052 != 0 {
		t.Fatal("non-press state appended text")
	}
	if o.key(2, 2) != 1 || d.Text[0] != '1' || d.Text[1] != 0 || d.Field_1052 != 1 {
		t.Fatal("both filters rejected digit")
	}
	// Removing only the digit filter restores the alphanumeric behavior.
	d.Field_1028 = 0
	if o.key(30, 2) != 1 || d.Text[1] != 'a' || d.Text[2] != 0 || d.Field_1052 != 2 {
		t.Fatal("alphanumeric filter rejected letter")
	}
	if o.key(57, 2) != 1 || d.Text[2] != 0 || d.Field_1052 != 2 {
		t.Fatal("alphanumeric filter accepted space")
	}
}
