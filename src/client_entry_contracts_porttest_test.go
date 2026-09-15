//go:build porttest

package opennox

import (
	"github.com/opennox/libs/client/seat"
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestClientEntryFocusedDestroyContract(t *testing.T) {
	o := newEntryOwner(t)
	o.create(t, 0, 8, nil)
	o.c.GUI.Focus(o.win)
	if !o.seat.enabled || o.env.State()[1] == 0 {
		t.Fatal("entry focus did not enable text input")
	}
	o.win.Destroy()
	o.c.GUI.FreeDestroyed()
	o.win = nil
	if o.seat.enabled || o.env.State()[1] != 0 {
		t.Fatal("destroyed entry retained active text-input ownership")
	}
}
func TestClientEntryInputOwnerContract(t *testing.T) {
	o := newEntryOwner(t)
	d := o.create(t, 0, 8, nil)
	o.c.GUI.Focus(o.win)
	if o.key(30, 2) != 1 || d.Text[0] != 'a' || d.Field_1052 != 1 {
		t.Fatal("real keyboard mapping did not append a")
	}
	o.key(14, 2)
	if d.Text[0] != 0 || d.Field_1052 != 0 {
		t.Fatal("backspace did not remove one code unit")
	}
	o.c.GUI.Focus(nil)
	if o.seat.enabled || o.env.State()[1] != 0 {
		t.Fatal("focus loss did not release text input")
	}
	d = o.create(t, 6, 8, nil)
	o.c.GUI.Focus(o.win)
	o.seat.send(&seat.TextEditEvent{Text: "compose"})
	o.key(30, 2)
	if d.Text[256] != 'c' || d.Field_1052>>16 != 7 {
		t.Fatal("real input composition was not copied")
	}
	o.seat.send(&seat.TextInputEvent{Text: "水"})
	if d.Text[0] != 0x6c34 || d.Field_1052&0xffff != 1 {
		t.Fatal("real committed text did not reach entry")
	}
	if d.Field_1048 == 0 {
		t.Fatal("composition listbox missing")
	}
}
func TestClientEntryConstructorContract(t *testing.T) {
	o := newEntryOwner(t)
	d := o.create(t, 0, 8, func(_ *gui.WindowData, d *gui.EntryFieldData) {
		for i := range d.Text {
			d.Text[i] = 0x1234
		}
		d.Field_1040 = 300
		d.Field_1044 = 0xffffffff
		d.Field_1052 = 0xffffffff
	})
	if d.Field_1040 != 256 || d.Field_1044 != 0 || d.Field_1052 != 0 || d.Field_1048 != 0 {
		t.Fatal("constructor scalar initialization")
	}
	for i, v := range d.Text {
		want := uint16(0x1234)
		if i%256 < 128 {
			want = 0
		}
		if v != want {
			t.Fatalf("constructor code unit %d: %x != %x", i, v, want)
		}
	}
}

func TestClientEntryLongCompositionContract(t *testing.T) {
	o := newEntryOwner(t)
	d := o.create(t, 6, 8, nil)
	o.c.GUI.Focus(o.win)
	text := ""
	for i := 0; i < 400; i++ {
		text += "水"
	}
	o.seat.send(&seat.TextEditEvent{Text: text})
	o.key(30, 2)
	if d.Field_1052>>16 != 255 || d.Text[511] != 0 || d.Field_1040 != 256 || d.Field_1048 == 0 {
		t.Fatal("long composition exceeded its 256-unit buffer")
	}
	for _, v := range d.Text[256:511] {
		if v != 0x6c34 {
			t.Fatal("long composition prefix changed")
		}
	}
	legacy.NoxInputOnChar('a')
	if d.Field_1052>>16 != 255 || d.Text[511] != 0 || d.Text[0] != 'a' || d.Field_1052&0xffff != 1 {
		t.Fatal("character path did not bound composition")
	}
}
func TestClientEntryDestroyOtherContract(t *testing.T) {
	o := newEntryOwner(t)
	o.create(t, 0, 8, nil)
	old := o.win
	o.win = nil
	o.create(t, 0, 8, nil)
	o.c.GUI.Focus(o.win)
	active := o.env.State()[1]
	old.Destroy()
	o.c.GUI.FreeDestroyed()
	if !o.seat.enabled || o.env.State()[1] != active {
		t.Fatal("destroying inactive entry disturbed active owner")
	}
}

func TestClientEntryCompositionTailContract(t *testing.T) {
	o := newEntryOwner(t)
	d := o.create(t, 6, 8, nil)
	o.c.GUI.Focus(o.win)
	for path := 0; path < 2; path++ {
		for i := 256; i < 512; i++ {
			d.Text[i] = 0xabcd
		}
		o.seat.send(&seat.TextEditEvent{Text: "ab"})
		if path == 0 {
			o.key(30, 2)
		} else {
			legacy.NoxInputOnChar('x')
		}
		if d.Text[256] != 'a' || d.Text[257] != 'b' || d.Text[258] != 0 || d.Field_1052>>16 != 2 {
			t.Fatal("composition copy prefix/count")
		}
		for i := 259; i < 512; i++ {
			if d.Text[i] != 0xabcd {
				t.Fatalf("composition copy changed unused unit %d", i)
			}
		}
	}
}
