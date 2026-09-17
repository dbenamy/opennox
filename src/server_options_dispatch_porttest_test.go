//go:build porttest

package opennox

import (
	"encoding/binary"
	"fmt"
	"github.com/opennox/libs/client/keybind"
	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/opennox/v1/client/gui"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
	"unsafe"
)

func TestServerOptionsIgnoredEvents(t *testing.T) {
	o := newServerOptionsOwner(t)
	for _, event := range []int{-1, 0, 21, 23, 16384, 16390, 16392, 16399, 16401} {
		want := 1
		if event == 23 {
			want = 0
		}
		if got := legacy.PortTestServerOptionsEvent(o.options, nil, event, 0); got != want {
			t.Fatalf("event %d returned %d", event, got)
		}
	}
	if got := legacy.PortTestServerOptionsEdit(o.options, 0, 9999); got != 0 {
		t.Fatal("unknown edit control")
	}
}

func TestServerOptionsCycleControl(t *testing.T) {
	type row struct {
		Flags                     uint32
		Checked, Initial, Enabled bool
		Dirty                     uint32
	}
	var rows []row
	for _, flags := range []uint32{0, 1, 128, 129, 16384, 32768, 49153} {
		for _, checked := range []bool{false, true} {
			for _, initial := range []bool{false, true} {
				t.Run(fmt.Sprintf("%x-%t-%t", flags, checked, initial), func(t *testing.T) {
					o := newServerOptionsOwner(t)
					defer noxflags.PortTestGameFlags(noxflags.GameFlag(flags))()
					child := o.options.ChildByID(10122)
					if checked {
						child.DrawData().Field0 |= 4
					} else {
						child.DrawData().Field0 &^= 4
					}
					for _, id := range []uint{10183, 10197} {
						w := o.options.ChildByID(id)
						if initial {
							w.Flags |= gui.StatusEnabled
						} else {
							w.Flags &^= gui.StatusEnabled
						}
					}
					if got := legacy.PortTestServerOptionsEvent(o.options, child, 16391, 0); got != 1 {
						t.Fatal("checkbox dispatch return")
					}
					want := checked
					if flags&49152 != 0 {
						want = initial
					}
					for _, id := range []uint{10183, 10197} {
						if o.options.ChildByID(id).Flags.IsEnabled() != want {
							t.Fatal("cycle panel enablement")
						}
					}
					if *o.optionWords["dirty"] != 1 {
						t.Fatal("cycle edit not marked dirty")
					}
					rows = append(rows, row{flags, checked, initial, want, *o.optionWords["dirty"]})
				})
			}
		}
	}
	spellbookCapture(t, "server-options-cycle-control", rows, "6c181982a123d50d20a96007d54fca66cd22e6b583818e34395f59ca095808e6")
}

func TestServerOptionsCancel(t *testing.T) {
	type row struct {
		Flags       uint32
		Mode, After uint16
		ID          uint
		Open        bool
	}
	var rows []row
	for _, flags := range []uint32{0, 1, 128, 129} {
		for _, mode := range []uint16{0x20, 0x8000, 0xc120, 0x100} {
			for _, id := range []uint{10146, 10149} {
				t.Run(fmt.Sprintf("%x-%x-%d", flags, mode, id), func(t *testing.T) {
					o := newServerOptionsOwner(t)
					defer noxflags.PortTestGameFlags(noxflags.GameFlag(flags))()
					binary.LittleEndian.PutUint16(o.settings[52:], mode)
					if got := legacy.PortTestServerOptionsEvent(o.options, o.options.ChildByID(id), 16391, 0); got != 1 {
						t.Fatal("cancel return")
					}
					expectedPacket := id == 10146 && flags&128 != 0 && mode&0x8000 != 0
					packets := o.packets()
					if expectedPacket {
						if len(packets) != 1 || len(packets[0].Data) != 2 || packets[0].Data[0] != byte(netmsg.MSG_TEAM_MSG) || packets[0].Data[1] != 7 {
							t.Fatalf("cancel team message: %+v", packets)
						}
					} else if len(packets) != 0 {
						t.Fatalf("unexpected cancel messages: %+v", packets)
					}
					want := mode
					if id == 10146 && flags&128 != 0 {
						want &= 0x3fff
					}
					after := binary.LittleEndian.Uint16(o.settings[52:])
					if after != want || *o.optionWords["root"] != 0 {
						t.Fatalf("cancel mode %x want %x root %x", after, want, *o.optionWords["root"])
					}
					rows = append(rows, row{flags, mode, after, id, false})
				})
			}
		}
	}
	spellbookCapture(t, "server-options-cancel", rows, "1854ecbf95825829afdfef1e13869bbb75c92148a7093c45b79c300d591135aa")
}

func TestServerOptionsDropdownOpen(t *testing.T) {
	o := newServerOptionsOwner(t)
	legacy.PortTestServerOptionsModeName(0)
	list := o.options.ChildByID(10120)
	list.SetHidden(true)
	for i := 0; i < 2; i++ {
		if got := legacy.PortTestServerOptionsEvent(o.options, o.options.ChildByID(10119), 16391, 0); got != 1 {
			t.Fatal("dropdown button return")
		}
		if list.Flags.IsHidden() || o.c.GUI.Captured() != list {
			t.Fatal("dropdown visibility/capture")
		}
		if n := (*gui.ScrollListBoxData)(list.WidgetData).Field_11_0; n != 6 {
			t.Fatalf("dropdown count %d", n)
		}
	}
}

func TestServerOptionsCycleCheckboxClick(t *testing.T) {
	o := newServerOptionsOwner(t)
	o.options.SetFunc94(func(w *gui.Window, e gui.WindowEvent) gui.WindowEventResp {
		a, b := e.EventArgsC()
		return gui.RawEventResp(legacy.PortTestServerOptionsEvent(w, (*gui.Window)(unsafe.Pointer(a)), e.EventCode(), int(b)))
	})
	child := o.options.ChildByID(10122)
	if child.DrawData().Window != o.options {
		t.Fatal("checkbox event owner")
	}
	for i := 0; i < 4; i++ {
		child.Func93(gui.WindowKeyPress{Key: keybind.KeySpace, Pressed: true})
		checked := child.DrawData().Field0&4 != 0
		if checked != (i%2 == 0) {
			t.Fatal("checkbox did not toggle")
		}
		for _, id := range []uint{10183, 10197} {
			if o.options.ChildByID(id).Flags.IsEnabled() == checked {
				t.Fatal("new-game controls disagree with checkbox after click")
			}
		}
	}
}
