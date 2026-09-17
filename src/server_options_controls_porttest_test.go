//go:build porttest

package opennox

import (
	"encoding/binary"
	"fmt"
	"image"
	"strings"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func TestServerOptionsDropdown(t *testing.T) {
	o := newServerOptionsOwner(t)
	legacy.PortTestServerOptionsModeName(0)
	want := []string{"Capture flag", "Arena battle", "Last survivor", "King of the realm", "Flag ball", "Chat lobby"}
	for repeat := 0; repeat < 3; repeat++ {
		o.call("populate", 0, "")
		d := (*gui.ScrollListBoxData)(o.options.ChildByID(10120).WidgetData)
		if int(d.Field_11_0) != len(want) {
			t.Fatalf("mode rows %d", d.Field_11_0)
		}
		for i, name := range want {
			p := o.event(10120, 16406, uintptr(i), 0)
			if got := alloc.GoString16((*uint16)(unsafe.Pointer(uintptr(p)))); got != name {
				t.Fatalf("row %d %q", i, got)
			}
		}
	}
	w := o.options.ChildByID(10120)
	w.Off = image.Pt(31, 17)
	w.EndPos = image.Pt(290, 199)
	o.call("measure", 0, "")
	maxWidth := 0
	for _, s := range want {
		if width := o.c.Render().GetStringSizeWrapped(w.DrawData().Font(), s, 0).X; width > maxWidth {
			maxWidth = width
		}
	}
	height := len(want)*(o.c.Render().FontHeight(w.DrawData().Font())+1) + 2
	if w.SizeVal != image.Pt(maxWidth+7, height) || w.Off.X != 290-maxWidth-7 || w.EndPos != image.Pt(290, 17+height) {
		t.Fatalf("dropdown geometry %+v %+v %+v", w.SizeVal, w.Off, w.EndPos)
	}
}

func TestServerOptionsLimitsAndNames(t *testing.T) {
	type row struct {
		Score, Time         int
		ScoreText, TimeText string
	}
	var rows []row
	for _, score := range []int{0, 1, 255, 256, 32767, 32768, 65535} {
		for _, minutes := range []int{0, 1, 127, 128, 254, 255} {
			t.Run(fmt.Sprintf("%d-%d", score, minutes), func(t *testing.T) {
				o := newServerOptionsOwner(t)
				binary.LittleEndian.PutUint16(o.settings[54:], uint16(score))
				o.settings[56] = byte(minutes)
				o.call("limits", 0, "")
				got := row{score, minutes, o.entry(10134), o.entry(10135)}
				if got.ScoreText != fmt.Sprint(score) || got.TimeText != fmt.Sprint(minutes) {
					t.Fatal(got)
				}
				rows = append(rows, got)
			})
		}
	}
	spellbookCapture(t, "server-options-limits", rows, "7dabecf8dcca1aa651ad2845cbc975e67ef67144acd938436bb8c310c549e712")
	t.Run("names", func(t *testing.T) {
		o := newServerOptionsOwner(t)
		for _, name := range []string{"", "Nox", "123456789012345", "1234567890123456", strings.Repeat("x", 90)} {
			o.call("name", 0, name)
			want := name
			if len(want) > 15 {
				want = want[:15]
			}
			if got := o.entry(10101); got != want {
				t.Fatalf("name %q -> %q", name, got)
			}
		}
	})
}

func TestServerOptionsDirtyVisibility(t *testing.T) {
	o := newServerOptionsOwner(t)
	for _, value := range []int{0, 1, -1, 2, 2147483647, -2147483648} {
		if o.call("dirty-set", value, "") != value || o.call("dirty-get", 0, "") != value {
			t.Fatal("dirty value narrowed")
		}
	}
	if o.call("open", 0, "") != 1 {
		t.Fatal("owned root not open")
	}
	for _, value := range []int{1, 0, -1, 0} {
		o.call("visible", value, "")
		if o.options.GetFlags().IsHidden() != (value != 0) {
			t.Fatal("visibility")
		}
	}
	*o.optionWords["root"] = 0
	if o.call("open", 0, "") != 0 || o.call("visible", 1, "") != 0 {
		t.Fatal("closed visibility")
	}
}

func TestServerOptionsReadSettings(t *testing.T) {
	type row struct {
		Name, Score, Time string
		Checked           bool
		Record            []byte
	}
	var rows []row
	for _, score := range []string{"", "0", "1", "-1", "65535", "65536", " 12", "17tail"} {
		for _, minutes := range []string{"", "0", "255", "256", "-1"} {
			for _, checked := range []bool{false, true} {
				t.Run(fmt.Sprintf("score%q-time%q-checked%t", score, minutes, checked), func(t *testing.T) {
					o := newServerOptionsOwner(t)
					legacy.PortTestServerOptionsModeName(0)
					for i := range o.settings {
						o.settings[i] = byte(0x80 + i)
					}
					o.text(10101, "Server name")
					o.text(10134, score)
					o.text(10135, minutes)
					o.event(10119, 16385, uintptr(unsafe.Pointer(alloc.InternCString16("Capture flag"))), 0)
					d := o.options.ChildByID(10122).DrawData()
					if checked {
						d.Field0 |= 4
					} else {
						d.Field0 &^= 4
					}
					for i := 0; i < 20; i++ {
						*memmap.PtrUint8(0x5D4594, 1045488+uintptr(i)) = byte(3*i + 7)
					}
					*memmap.PtrUint32(0x5D4594, 1045452) = 0x12345678
					*memmap.PtrUint32(0x5D4594, 1045456) = 0x89abcdef
					o.call("settings-read", 0, "")
					if o.settings[0] != 0 || binary.LittleEndian.Uint16(o.settings[52:]) != 0x20 || binary.LittleEndian.Uint32(o.settings[44:]) != 0x12345678 || binary.LittleEndian.Uint32(o.settings[48:]) != 0x89abcdef {
						t.Fatalf("settings fields %x", o.settings)
					}
					wantScore := map[string]uint16{"": 0xb7b6, "0": 0, "1": 1, "-1": 65535, "65535": 65535, "65536": 0, " 12": 12, "17tail": 17}[score]
					wantTime := map[string]byte{"": 0xb8, "0": 0, "255": 255, "256": 0, "-1": 255}[minutes]
					if binary.LittleEndian.Uint16(o.settings[54:]) != wantScore || o.settings[56] != wantTime {
						t.Fatalf("numeric fields: %x", o.settings[54:57])
					}
					if o.settings[57] != byte(map[bool]int{false: 0, true: 1}[checked]) {
						t.Fatal("checkbox")
					}
					for i := 0; i < 20; i++ {
						if o.settings[24+i] != byte(3*i+7) {
							t.Fatal("spell mask")
						}
					}
					if string(o.settings[9:20]) != "Server name" || o.settings[20] != 0 {
						t.Fatal("server name")
					}
					rows = append(rows, row{"Server name", score, minutes, checked, append([]byte(nil), o.settings...)})
				})
			}
		}
	}
	spellbookCapture(t, "server-options-settings-read", rows, "c53b5f690aaa4fd4c8515fb5f3cdbdf7febf76e8e9f5a53f28dd55afa4deab25")
}
