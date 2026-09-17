//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"image"
	"testing"
	"unsafe"
)

func TestServerOptionsDrawing(t *testing.T) {
	type row struct {
		Hidden, Inside        bool
		Mouse                 image.Point
		AfterHidden, Captured bool
		Pixels                string
	}
	var rows []row
	for _, hidden := range []bool{false, true} {
		for _, tc := range []struct {
			p      image.Point
			inside bool
		}{{image.Pt(10, 10), true}, {image.Pt(70, 70), false}, {image.Pt(4, 10), false}, {image.Pt(10, 4), false}} {
			t.Run(fmt.Sprintf("%t-%v", hidden, tc.p), func(t *testing.T) {
				o := newServerOptionsOwner(t)
				o.c.Mouse = tc.p
				o.options.SetPos(image.Pt(5, 5))
				o.options.SizeVal = image.Pt(60, 70)
				list := o.options.ChildByID(10120)
				list.SetPos(image.Pt(0, 0))
				list.SizeVal = image.Pt(40, 40)
				list.SetHidden(hidden)
				list.Capture(true)
				for i := range o.pix.Pix {
					o.pix.Pix[i] = 0xffff
				}
				before := effectsPixelHash(o.pix)
				if legacy.PortTestServerOptionsDraw(o.options) != 1 {
					t.Fatal("draw return")
				}
				after := effectsPixelHash(o.pix)
				if after == before {
					t.Fatal("background did not render")
				}
				wantHidden := hidden || !tc.inside
				if list.GetFlags().IsHidden() != wantHidden {
					t.Fatal("dropdown dismissal")
				}
				captured := o.c.GUI.Captured() == list
				if captured != (hidden || tc.inside) {
					t.Fatal("dropdown capture")
				}
				rows = append(rows, row{hidden, tc.inside, tc.p, list.GetFlags().IsHidden(), captured, after})
			})
		}
	}
	spellbookCapture(t, "server-options-drawing", rows, "0dd395d62e6ea918949ca90927e6b1b3f88cc1d85c9b3c5c3b081e9c3c4b7df8")
}

func TestServerOptionsTooltips(t *testing.T) {
	newServerOptionsOwner(t)
	cursor := unsafe.Slice(memmap.PtrUint16(0x5D4594, 1096676), 256)
	old := append([]uint16(nil), cursor...)
	t.Cleanup(func() { copy(cursor, old) })
	type row struct {
		Kind  string
		Flags byte
		Text  string
	}
	var rows []row
	for _, kind := range []string{"assign", "damage"} {
		for flags := 0; flags < 256; flags++ {
			value := byte(flags)
			if got := legacy.PortTestServerOptions("tooltip-"+kind, unsafe.Pointer(&value), 0, ""); got != 1 {
				t.Fatal("tooltip return")
			}
			want := kind + " off"
			if flags&4 != 0 {
				want = kind + " on"
			}
			text := alloc.GoString16(&cursor[0])
			if text != want {
				t.Fatalf("%s flags %x: %q", kind, flags, text)
			}
			rows = append(rows, row{kind, value, text})
		}
	}
	spellbookCapture(t, "server-options-tooltips", rows, "3c9414088fd9de2adc5977b173adab4928c58b0a8643a1d2af15b9156a1d739a")
}
