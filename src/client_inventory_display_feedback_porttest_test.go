//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestClientInventoryDisplaySecondaryFeedback(t *testing.T) {
	o := newInventoryDisplayOwner(t)
	var rows []inventoryDisplayResult
	for _, index := range []int{0, 19, 20, 21, 42, 63, 83} {
		for _, previous := range []bool{false, true} {
			for _, pending := range []uint32{0, 321, 999} {
				for _, code := range []uint32{0, 321, 322, 999} {
					for _, status := range []uint32{0, 1, 2, 255} {
						o.reset(t)
						cells := legacy.PortTestInventoryCells()
						dr := o.item(t, "Bow", 123)
						cells[index].Drawable, cells[index].Count = dr, 2
						cells[index].Codes[0], cells[index].Codes[1] = 321, 322
						old := (index + 1) % 84
						if previous {
							p := o.item(t, "Bow", 500)
							cells[old].Drawable, cells[old].Count, cells[old].Codes[0], cells[old].Alternate = p, 1, 501, 1
							*o.displayWords["dword_5d4594_1062480"] = o.cell(old)
						}
						*o.displayWords["dword_5d4594_1062484"] = pending
						r := o.call(t, len(rows), 15, uintptr(code), uintptr(status), 0)
						found := code == 321 || code == 322
						if found {
							if r.Return != 1 || cells[index].Alternate != 1 || *o.displayWords["dword_5d4594_1062480"] != o.cell(index) || cells[old].Alternate != 0 {
								t.Fatal("successful feedback must replace alternate selection")
							}
							if len(o.console) != 0 || *o.displayWords["dword_5d4594_1062484"] != pending {
								t.Fatal("successful feedback unexpectedly consumed failure state")
							}
						} else if r.Return != 0 {
							t.Fatal("missing feedback return")
						}
						if !found && status == 1 && (len(o.console) != 1 || o.console[0] != "6:System: Cannot use secondary weapon") {
							t.Fatalf("missing secondary rejection message: %v", o.console)
						}
						rows = append(rows, r)
					}
				}
			}
		}
	}
	inventoryDisplayCapture(t, "feedback", rows, "8752e095de6066d9c89d4241db1c6b4c7a703a5043b06a5c60714a22a804581e")
}
