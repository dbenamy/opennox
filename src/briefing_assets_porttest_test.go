//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"golang.org/x/image/font/basicfont"
	"os"
	"path/filepath"
	"testing"
	"unsafe"
)

// Shipped window geometry and real presentation owners with authored content.
// This is an integration fixture, separate from a fresh running game scene.
func TestBriefingShippedWindow(t *testing.T) {
	path := os.Getenv("OPENNOX_BRIEFING_ASSETS")
	if path == "" {
		t.Skip("set OPENNOX_BRIEFING_ASSETS to the original Nox data directory")
	}
	o := newBriefingOwner(t)
	font, free := o.c.Render().GetFonts().PortTestWindowFont(basicfont.Face7x13, "large")
	t.Cleanup(free)
	o.c.dataRefs[uint32(txptr(font))] = 0xf0200001
	old := legacy.Nox_new_window_from_file
	t.Cleanup(func() { legacy.Nox_new_window_from_file = old })
	legacy.Nox_new_window_from_file = func(name string, fn gui.WindowFunc) *gui.Window {
		if name != "Briefing.wnd" {
			return old(name, fn)
		}
		f, err := os.Open(filepath.Join(path, "window", name))
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		return newWindowFromReader(o.c.GUI, f, fn)
	}
	var rows []briefingResult
	for _, mode := range []int{1, 2, 4} {
		o.resetBriefing(t)
		o.constructBriefing(t)
		*memmap.PtrUint32(0x5D4594, 832468) = 17
		*memmap.PtrPtr(0x5D4594, 832464) = unsafe.Pointer(alloc.InternCString16("A custom briefing."))
		legacy.PortTestBriefing(9, 0, 0, 0)
		legacy.PortTestBriefing(15, 254, 1, uintptr(mode))
		if o.c.GUI.Captured() == nil || o.briefWindow().Flags.IsHidden() || memmap.Uint32(0x5D4594, 832472) != uint32(mode) {
			t.Fatal("shipped briefing mode did not open")
		}
		op := 1
		if mode == 2 {
			op = 2
		} else if mode == 4 {
			op = 3
		}
		for _, frame := range []uint32{30, 31} {
			rows = append(rows, o.drawBriefing(t, op, frame, 0))
		}
	}
	briefingCapture(t, "shipped-window", rows, "ce442082b9c3b2c3d8ea55ec3ce3104f46d125b6373dbacb730d248419edcfa8")
}
