//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"testing"
	"unsafe"
)

func TestServerPanelsDrawing(t *testing.T) {
	type row struct {
		Kind     string
		Mode     int
		Position image.Point
		Pixels   string
	}
	var rows []row
	for _, kind := range []string{"weapon", "spell", "general"} {
		for mode := 0; mode < 3; mode++ {
			for _, pos := range []image.Point{image.Pt(5, 7), image.Pt(-3, -2), image.Pt(17, 13)} {
				t.Run(fmt.Sprintf("%s-%d-%v", kind, mode, pos), func(t *testing.T) {
					o := newServerOptionsOwner(t)
					o.installSubpanels(t, true)
					_, restore := legacy.PortTestServerPanelsObjectWords()
					t.Cleanup(restore)
					counts := unsafe.Slice(memmap.PtrUint32(0x5D4594, 1045472), 2)
					old := append([]uint32(nil), counts...)
					t.Cleanup(func() { copy(counts, old) })
					raw := legacy.PortTestServerPanelsConstruct(kind, o.options, unsafe.Pointer(&o.settings[0]))
					if raw == 0 {
						t.Fatal("draw constructor")
					}
					w := (*gui.Window)(unsafe.Pointer(uintptr(uint32(raw))))
					o.options.SetPos(image.Pt(3, 2))
					w.SetPos(pos)
					w.SizeVal = image.Pt(21, 17)
					w.DrawData().BgColorVal = 0x12345678
					if mode == 1 {
						w.DrawData().BgColorVal = 0x80000000
					}
					if mode == 2 {
						w.Flags |= gui.StatusImage
						w.DrawData().SetBackgroundImage(o.images[3])
					}
					for i := range o.pix.Pix {
						o.pix.Pix[i] = 0xffff
					}
					before := effectsPixelHash(o.pix)
					w.Draw()
					after := effectsPixelHash(o.pix)
					wantChange := kind != "general" && mode != 1
					if (after != before) != wantChange {
						t.Fatal("draw path did not match background/image/no-op contract")
					}
					rows = append(rows, row{kind, mode, pos, after})
				})
			}
		}
	}
	spellbookCapture(t, "server-panels-drawing", rows, "da781e6b25048a32c1febfeaec15a8ea7e6d3e3016ad3a3c947987546a2a38b7")
}
