//go:build porttest

package opennox

import (
	"encoding/binary"
	"fmt"
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"reflect"
	"testing"
	"unsafe"
)

func TestServerPanelsClassSelection(t *testing.T) {
	type row struct {
		Initial   byte
		Index     int
		After     byte
		Selection []int32
	}
	var rows []row
	for _, extra := range []byte{0, 0x18, 0x60, 0xf8} {
		for classes := 0; classes < 8; classes++ {
			t.Run(fmt.Sprintf("%x-%d", extra, classes), func(t *testing.T) {
				o := newServerOptionsOwner(t)
				o.installSubpanels(t, true)
				state := memmap.PtrUint8(0x5D4594, 371616)
				*state = extra | byte(classes)
				raw := legacy.PortTestServerPanelsConstruct("access", o.options, unsafe.Pointer(&o.settings[0]))
				w := (*gui.Window)(unsafe.Pointer(uintptr(uint32(raw))))
				if w == nil {
					t.Fatal("class root")
				}
				list := w.ChildByID(10123)
				d := (*gui.ScrollListBoxData)(list.WidgetData)
				if d.Field_4 == 0 {
					t.Fatal("class list must support multiple selections")
				}
				selection := unsafe.Slice((*int32)(unsafe.Pointer(uintptr(d.Field_12))), int(d.Count)+1)
				var actual, want []int32
				for _, v := range selection {
					if v < 0 {
						break
					}
					actual = append(actual, v)
				}
				for i := 0; i < 3; i++ {
					if classes&(1<<uint(i)) != 0 {
						want = append(want, int32(i))
					}
				}
				if !reflect.DeepEqual(actual, want) {
					t.Fatalf("class constructor %v want %v", actual, want)
				}
				for _, index := range []int{-1, 0, 1, 2, 3, 7, 8, 31, 32, 255} {
					before := extra | byte(classes)
					*state = before
					result := gui.EventRespInt(w.Func94(gui.AsWindowEvent(16400, uintptr(list.C()), uintptr(uint32(index)))))
					expected := before ^ byte(uint32(1)<<uint(uint32(index)&31))
					if result != 0 || *state != expected {
						t.Fatal("class event", index, *state, expected)
					}
					rows = append(rows, row{before, index, *state, actual})
				}
			})
		}
	}
	spellbookCapture(t, "server-panels-class-selection", rows, "cd6a96bd91465b68d0198f69721530af394a8e86006f82a46d5d625ac9c78350")
}

func TestServerPanelsAdmissionFlags(t *testing.T) {
	type row struct {
		ID            uint
		Checked       bool
		Text          string
		Before, After []byte
		Enabled       bool
	}
	var rows []row
	o := newServerOptionsOwner(t)
	o.installSubpanels(t, true)
	raw := legacy.PortTestServerPanelsConstruct("access", o.options, unsafe.Pointer(&o.settings[0]))
	w := (*gui.Window)(unsafe.Pointer(uintptr(uint32(raw))))
	if w == nil {
		t.Fatal("admission flags root")
	}
	state := unsafe.Slice(memmap.PtrUint8(0x5D4594, 371616), 10)
	seed := []byte{0xa5, 0xa5, 0xa5, 0xa5, 0xa5, 0xa5, 0xa5, 0xa5, 0xa5, 0xa5}
	for _, id := range []uint{10102, 10103, 10124, 10125, 10127, 10129, 10131} {
		for _, checked := range []bool{false, true} {
			for _, tc := range []struct {
				text  string
				value int
			}{{"", 0}, {"0", 0}, {"5", 5}, {"14", 14}, {"-1", -1}, {"65536", 65536}} {
				copy(state, seed)
				child := w.ChildByID(id)
				child.DrawData().Field0 &^= 4
				if checked {
					child.DrawData().Field0 |= 4
				}
				edit := w.ChildByID(10104)
				switch id {
				case 10125:
					edit = w.ChildByID(10126)
				case 10127:
					edit = w.ChildByID(10128)
				case 10129:
					edit = w.ChildByID(10130)
				case 10131:
					edit = w.ChildByID(10132)
				}
				serverPanelsSetText(edit, tc.text)
				want := append([]byte(nil), seed...)
				expectedEnabled := edit.Flags.IsEnabled()
				switch id {
				case 10102:
					want[0] ^= 0x10
				case 10103:
					want[0] ^= 0x20
					expectedEnabled = !checked
				case 10124:
					want[2] ^= 0x80
				case 10125:
					expectedEnabled = !checked
					if checked {
						want[1] |= 15
					} else if tc.text != "" {
						want[1] = 0xa0 | byte(tc.value)
					}
				case 10127:
					expectedEnabled = !checked
					if checked {
						want[1] |= 0xf0
					} else if tc.text != "" {
						want[1] = 5 | byte(tc.value)
					} // legacy checkbox path does not shift the edit value.
				case 10129, 10131:
					expectedEnabled = !checked
					off := 5
					if id == 10131 {
						off = 7
					}
					if checked {
						binary.LittleEndian.PutUint16(want[off:], 0xffff)
					} else if tc.text != "" {
						binary.LittleEndian.PutUint16(want[off:], uint16(tc.value))
					}
				}
				if gui.EventRespInt(w.Func94(gui.AsWindowEvent(16391, uintptr(child.C()), 0))) != 0 {
					t.Fatal("admission toggle result")
				}
				if !reflect.DeepEqual(state, want) || edit.Flags.IsEnabled() != expectedEnabled {
					t.Fatalf("admission flag %d/%t/%q state %x want %x enabled %t want %t", id, checked, tc.text, state, want, edit.Flags.IsEnabled(), expectedEnabled)
				}
				rows = append(rows, row{id, checked, tc.text, append([]byte(nil), seed...), append([]byte(nil), state...), edit.Flags.IsEnabled()})
			}
		}
	}
	spellbookCapture(t, "server-panels-admission-flags", rows, "ccfca1fc04e7e7528fec282d0d161181022c24de9b7a0bff3b5a4be02e88ffca")
}
