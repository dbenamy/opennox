//go:build porttest

package opennox

import (
	"image"
	"testing"

	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
)

func TestJournalMeasuredHeight(t *testing.T) {
	o := newJournalOwner(t)
	var rows []journalResult
	for _, flags := range []uint16{0, 1, 2, 3, 4, 8, 15, 0x8000, 0xffff} {
		for _, names := range [][]string{{}, {"Short"}, {"Wrapped"}, {"Multiline"}, {"Empty"}, {"Unicode"}, {"Short", "Wrapped", "Multiline", "Empty", "Unicode"}} {
			o.resetJournal(t)
			for _, name := range names {
				o.call(t, 0, 31, name, flags)
			}
			font := o.c.r.GetFonts().AsFont(nil)
			height := -o.c.r.FontHeight(nil)
			prefix := ""
			switch flags {
			case 2:
				prefix = "Quest:"
			case 4:
				prefix = "Complete:"
			case 8:
				prefix = "Hint:"
			}
			// Independent expected layout uses the fixture's authored strings and public renderer.
			texts := map[string]string{"Short": "A short entry.", "Wrapped": "This journal entry has enough words to wrap over several lines in the real renderer.", "Multiline": "First line\nSecond line\nThird line", "Empty": "", "Unicode": "A tale of café and Ω."}
			for _, name := range names {
				height += o.c.r.FontHeight(nil) + o.c.r.GetStringSizeWrapped(font, prefix+" "+texts[name], 240).Y
			}
			if height < 0 {
				height = 0
			}
			legacy.PortTestJournal(9, 0, 0, 0)
			r := o.snapshot(t, 9, 0, true)
			rows = append(rows, r)
			if r.Height != uint32(height) {
				t.Fatalf("journal height flags%d names%v got%d want%d", flags, names, r.Height, height)
			}
		}
	}
	journalCapture(t, "measured-height", rows, "68432c8d7189658bf1c6e1426f29d2dfff49c30092c8f357275d8d1bce4d1b64")
}
func TestJournalDrawOrderAndScroll(t *testing.T) {
	o := newJournalOwner(t)
	var rows []journalResult
	for _, flags := range []uint16{0, 1, 2, 3, 4, 8, 15, 0xffff} {
		for _, scroll := range []int{0, 1, 12, 13, 14, 40, 80, 149, 150, 151, 400} {
			for _, pos := range []image.Point{{0, 0}, {17, 23}} {
				o.resetJournal(t)
				for _, name := range []string{"Short", "Wrapped", "Multiline", "Unicode"} {
					o.call(t, 0, 31, name, flags)
				}
				legacy.PortTestJournal(9, 0, 0, 0)
				before := memmap.Uint32(0x5D4594, 1064848)
				legacy.PortTestJournal(10, uintptr(pos.X), uintptr(pos.Y), uintptr(scroll))
				r := o.snapshot(t, 10, 0, true)
				rows = append(rows, r)
				if r.Height != before {
					t.Fatal("drawing mutated cached height")
				}
			}
		}
	}
	journalCapture(t, "draw-scroll", rows, "50af65c7e3ab957fce8d6302258f3d79fbe875c122b06ed5a0005400e1fb9300")
}
