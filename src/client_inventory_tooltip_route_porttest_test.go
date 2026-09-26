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

// This follows the callbacks installed by the real inventory constructor rather
// than calling the tooltip C adapters directly. The owner still checks each
// stored callback against the fixture's established normalization ID.
func TestClientInventoryTooltipWindowRoute(t *testing.T) {
	o := newInventoryWindowOwner(t)
	o.reset(t)
	o.construct(t)

	callbacks := legacy.PortTestInventoryWindowCallbacks()
	check := func(name string, w *gui.Window, callbackIndex int) {
		t.Helper()
		if w == nil || w.TooltipFuncPtr == nil {
			t.Fatalf("%s tooltip callback is not installed", name)
		}
		if w.TooltipFuncPtr != callbacks[callbackIndex] {
			t.Fatalf("%s callback identity %p want table[%d]=%p", name, w.TooltipFuncPtr, callbackIndex, callbacks[callbackIndex])
		}
		wantID := uint32(0xed200000) + uint32(callbackIndex)
		if gotID := o.c.callbackRefs[w.TooltipFuncPtr]; gotID != wantID {
			t.Fatalf("%s normalized callback ID %#x want %#x", name, gotID, wantID)
		}
	}
	text := func() string { return alloc.GoString16(&o.cursorText[0]) }
	clearText := func() { clear(o.cursorText) }

	// The current-weapon button is stored as nox_win_unk5 with its button as
	// Field100; state 2 selects the close-inventory tooltip.
	current := (*gui.Window)(unsafe.Pointer(uintptr(*o.windowWords["nox_win_unk5"])))
	button := current.Field100()
	check("button", button, 2)
	*memmap.PtrUint8(0x5D4594, 1049868) = 2
	button.TooltipFunc(0)
	if got := text(); got != "Close" {
		t.Fatalf("button tooltip %q", got)
	}

	// The alternate cell is empty after reset, so it uses the secondary-slot
	// fallback tooltip.
	alternate := (*gui.Window)(unsafe.Pointer(uintptr(*o.windowWords["dword_5d4594_1062468"])))
	check("alternate", alternate, 3)
	clearText()
	alternate.TooltipFunc(0)
	if got := text(); got != "Secondary weapon" {
		t.Fatalf("alternate tooltip %q", got)
	}

	// The actual main hover callback consumes one packed uint16,uint16 point.
	// The adjacent special tray includes x313; the main grid starts at x314.
	// Check both independently named items and the outer top/right boundaries.
	main := o.mainWindow()
	check("main hover", main, 19)
	o.stack(t, 0, 0, 1, "RedApple", 100)
	itemText := alloc.GoString16(o.c.Things.TypeByID("RedApple").PrettyName)
	specialText := alloc.GoString16(o.c.Things.TypeByID("Gold").PrettyName)
	if itemText == "" || specialText == "" {
		t.Fatal("missing fixture item name")
	}
	for _, tc := range []struct {
		x, y int
		want string
	}{{313, 13, specialText}, {314, 12, ""}, {314, 13, itemText}, {514, 13, ""}} {
		clearText()
		main.TooltipFunc(inventoryWindowPoint(tc.x, tc.y))
		got := text()
		if got != tc.want {
			t.Fatalf("main hover at (%d,%d) text %q, want %q", tc.x, tc.y, got, tc.want)
		}
	}

	// Status is the main window's previous sibling. The x coordinate is packed
	// through the same callback word; extra bit 0 selects the first effect label.
	status := main.Prev()
	check("status", status, 4)
	*memmap.PtrUint32(0x5D4594, 1062540) = 0
	*memmap.PtrUint8(0x5D4594, 1062536) = 1
	clearText()
	status.TooltipFunc(inventoryWindowPoint(39, 20))
	if got := text(); got != "Status effect 0" {
		t.Fatalf("status tooltip at packed x=39: %q", got)
	}
	clearText()
	status.TooltipFunc(inventoryWindowPoint(40, 20))
	if got := text(); got != "" {
		t.Fatalf("status tooltip at packed x=40: %q", got)
	}
	clearText()
	status.TooltipFunc(inventoryWindowPoint(0xffff, 0xffff))
	if got := text(); got != "" {
		t.Fatalf("out-of-range packed status point produced %q", got)
	}

	// Mode buttons share sub_466E20; the constructor's Stats button keeps the
	// fifth callback's independent normalization identity.
	mode := main.ChildByID(9107)
	check("mode", mode, 37)
	clearText()
	mode.TooltipFunc(0)
	if got := text(); got != "Statistics" {
		t.Fatalf("mode tooltip %q", got)
	}
}
