//go:build porttest

package opennox

import (
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func (o *scoreboardOwner) rankColumn(side, column int) *gui.Window {
	return (*gui.Window)(*memmap.PtrPtr(0x5D4594, 1090060+uintptr(column*8+side*4)))
}
func TestScoreboardColumnsAndHeadings(t *testing.T) {
	o := newScoreboardOwner(t)
	var rows []scoreboardResult
	for mode := 1; mode <= 5; mode++ {
		for side := 0; side < 2; side++ {
			for _, padding := range []int{0, 1, 3} {
				o.resetRank(t)
				o.constructRank(t)
				*o.rankWords["dword_5d4594_1090120"] = uint32(mode)
				got := legacy.PortTestScoreboard(4, uintptr(side), uintptr(padding), 0)
				headings := []string{"player", "", "class", "score", "ping"}
				if mode == 1 {
					headings[2], headings[3], headings[4] = "LivesHeading", "HealthHeading", "class"
				}
				if mode == 5 {
					headings[3] = "rank"
				}
				for col, want := range headings {
					d := (*gui.ScrollListBoxData)(o.rankColumn(side, col).WidgetData)
					if int(d.Field_11_1) != padding+1 {
						t.Fatalf("mode%d side%d column%d count%d", mode, side, col, d.Field_11_1)
					}
					items := unsafe.Slice(d.Items, int(d.Count))
					if text := alloc.GoString16S(items[padding].Text[:]); text != want {
						t.Fatalf("mode%d column%d heading%q want%q", mode, col, text, want)
					}
				}
				rows = append(rows, o.rankCapture(t, 4, got))
				legacy.PortTestScoreboard(8, 0, 0, 0)
				for s := 0; s < 2; s++ {
					for col := 0; col < 5; col++ {
						if (*gui.ScrollListBoxData)(o.rankColumn(s, col).WidgetData).Field_11_1 != 0 {
							t.Fatal("clear retained rows")
						}
					}
				}
				rows = append(rows, o.rankCapture(t, 8, 0))
			}
		}
	}
	o.resetRank(t)
	o.constructRank(t)
	col := o.rankColumn(0, 0)
	for _, text := range []string{"", "plain", "é界😀", "100% complete"} {
		ret := legacy.PortTestScoreboard(9, uintptr(col.C()), 7, uintptr(unsafe.Pointer(alloc.InternCString16(text))))
		d := (*gui.ScrollListBoxData)(col.WidgetData)
		items := unsafe.Slice(d.Items, int(d.Count))
		if got := alloc.GoString16S(items[d.Field_11_1-1].Text[:]); got != text {
			t.Fatalf("insert%q got%q", text, got)
		}
		if ret != 1 {
			t.Fatal("insert return", ret)
		}
		rows = append(rows, o.rankCapture(t, 9, ret))
	}
	ret := legacy.PortTestScoreboard(0, uintptr(col.C()), 9, 0)
	d := (*gui.ScrollListBoxData)(col.WidgetData)
	if got := alloc.GoString16S(unsafe.Slice(d.Items, int(d.Count))[d.Field_11_1-1].Text[:]); got != "InternalError" {
		t.Fatal("missing string fallback", got)
	}
	rows = append(rows, o.rankCapture(t, 0, ret))
	scoreboardCapture(t, "columns-headings", rows, "4068fbafe0517aac01761badcddb0690e0bf6fdb81948421118e9b10fd96c66e")
}
