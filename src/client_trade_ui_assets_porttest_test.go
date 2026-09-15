//go:build porttest

package opennox

import (
	"os"
	"path/filepath"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/legacy"
)

// This uses the shipped window descriptions with the real GUI and renderer,
// but authored trade reports; it is not a two-client multiplayer scenario.
func TestClientTradeUIAssetWindows(t *testing.T) {
	path := os.Getenv("OPENNOX_TRADE_UI_ASSETS")
	if path == "" {
		t.Skip("set OPENNOX_TRADE_UI_ASSETS to the original Nox data directory")
	}
	o := newTradeUIOwner(t)
	old := legacy.Nox_new_window_from_file
	t.Cleanup(func() { legacy.Nox_new_window_from_file = old })
	legacy.Nox_new_window_from_file = func(name string, fn gui.WindowFunc) *gui.Window {
		file := ""
		switch name {
		case "Trade.wnd":
			file = "trade.wnd"
		case "MultMove.wnd":
			file = name
		default:
			return old(name, fn)
		}
		f, err := os.Open(filepath.Join(path, "window", file))
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		return newWindowFromReader(o.c.GUI, f, fn)
	}
	o.reset(t)
	o.construct(t)
	o.constructTrade(t)
	o.initAmount(t)
	var rows []tradeUIResult
	// Start, add both offers, draw and accept through the actual widget callback.
	data := o.record(t, 52)
	data[0], data[1] = 201, 12
	data[2] = 'P'
	rows = append(rows, o.tradeCapture(t, len(rows), 22, txptr(unsafe.Pointer(&data[0]))))
	if o.tradeWindow().Flags.IsHidden() {
		t.Fatal("shipped trade window did not open")
	}
	add := o.record(t, 15)
	for _, side := range []byte{1, 0} {
		tradeUIAddRecord(add, side, uint16(o.c.Things.TypeByID("RedApple").Index()), uint16(100+side), 7)
		rows = append(rows, o.tradeCapture(t, len(rows), 27, txptr(unsafe.Pointer(&add[0]))))
	}
	rows = append(rows, o.tradeCapture(t, len(rows), 16))
	w := o.tradeWindow()
	rows = append(rows, o.tradeCapture(t, len(rows), 14, txptr(w.C()), 16391, txptr(w.ChildByID(3708).C()), 0))
	rows = append(rows, o.tradeCapture(t, len(rows), 24))
	// Both price layouts and button callbacks on the shipped quantity window.
	for _, priced := range []uintptr{0, 1} {
		legacy.PortTestTradeUI(9, priced, 7)
		rows = append(rows, o.tradeCapture(t, len(rows), 7, o.quantityArgs(200, 100, 32)...))
		amount := legacy.Get_nox_gui_itemAmount_dialog_1319228()
		rows = append(rows, o.tradeCapture(t, len(rows), 5, txptr(amount.C()), 16391, txptr(amount.ChildByID(3602).C()), 0))
		if tradeQuantityText(amount, 3601) != "2" {
			t.Fatal("shipped quantity up button")
		}
		rows = append(rows, o.tradeCapture(t, len(rows), 4, txptr(amount.C())))
		rows = append(rows, o.tradeCapture(t, len(rows), 5, txptr(amount.C()), 16391, txptr(amount.ChildByID(3604).C()), 0))
		if len(o.callbacks) != int(priced)+1 || o.callbacks[len(o.callbacks)-1][5] != 2 {
			t.Fatal("shipped quantity accept callback")
		}
	}
	tradeUICapture(t, "asset-windows", rows, "e47882abd1467ccd8fc71735234118d0b3b6f2bfd5dfa68ef1961bf892611a94")
}
