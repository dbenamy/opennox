//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"os"
	"path/filepath"
	"testing"
	"unsafe"
)

// Real shipped resources and production widgets/callbacks, with authored shop
// reports and inventory state. This is distinct from a running game scene.
func TestClientShopUIAssetWindows(t *testing.T) {
	path := os.Getenv("OPENNOX_SHOP_UI_ASSETS")
	if path == "" {
		t.Skip("set OPENNOX_SHOP_UI_ASSETS to the original Nox data directory")
	}
	o := newShopUIOwner(t)
	old := legacy.Nox_new_window_from_file
	t.Cleanup(func() { legacy.Nox_new_window_from_file = old })
	legacy.Nox_new_window_from_file = func(name string, fn gui.WindowFunc) *gui.Window {
		if name != "Shop.wnd" && name != "MultMove.wnd" {
			return old(name, fn)
		}
		f, err := os.Open(filepath.Join(path, "window", name))
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		return newWindowFromReader(o.c.GUI, f, fn)
	}
	o.reset(t)
	o.construct(t)
	o.constructShop(t)
	o.initAmount(t)
	var rows []shopUIResult
	rows = append(rows, o.shopCapture(t, len(rows), 37, txptr(unsafe.Pointer(alloc.InternCString16("Merchant"))), txptr(unsafe.Pointer(alloc.InternCString("FixtureGreeting"))), uintptr(o.c.Things.TypeByID("Shopkeeper").Index())))
	if o.shopWindow().Flags.IsHidden() {
		t.Fatal("shipped shop did not open")
	}
	ids, free := alloc.New([4]byte{})
	defer free()
	*ids = [4]byte{255, 255, 255, 255}
	typ := uintptr(o.c.Things.TypeByID("RedApple").Index())
	for _, code := range []uintptr{100, 101} {
		rows = append(rows, o.shopCapture(t, len(rows), 16, typ, code, 7, 100, txptr(unsafe.Pointer(ids))))
	}
	rows = append(rows, o.shopCapture(t, len(rows), 22, 2))
	*o.windowWords["dword_5d4594_1062552"] = 100
	rows = append(rows, o.shopCapture(t, len(rows), 8))
	p, freePoint := alloc.New([2]int32{})
	defer freePoint()
	pos := o.shopWindow().ChildByID(3806).GlobalPos()
	*p = [2]int32{int32(pos.X + 25), int32(pos.Y + 25)}
	rows = append(rows, o.shopCapture(t, len(rows), 32, txptr(unsafe.Pointer(p))))
	w := legacy.Get_nox_gui_itemAmount_dialog_1319228()
	if *o.windowWords["dword_5d4594_1319268"] != 1 || *o.windowWords["dword_5d4594_1319248"] != 2 {
		t.Fatal("shipped buy quantity did not open")
	}
	legacy.PortTestTradeUI(5, txptr(w.C()), 16391, txptr(w.ChildByID(3602).C()), 0)
	if tradeQuantityText(w, 3601) != "2" {
		t.Fatal("shipped quantity up button")
	}
	legacy.PortTestTradeUI(4, txptr(w.C()))
	rows = append(rows, o.shopCapture(t, len(rows), 0))
	legacy.PortTestTradeUI(5, txptr(w.C()), 16391, txptr(w.ChildByID(3604).C()), 0)
	rows = append(rows, o.shopCapture(t, len(rows), 0))
	o.stack(t, 0, 0, 2, "RedApple", 200)
	for _, op := range []int{40, 41} {
		rows = append(rows, o.shopCapture(t, len(rows), op, 200, 7))
		if *o.windowWords["dword_5d4594_1319268"] != 1 {
			t.Fatalf("shipped quantity op%d did not open", op)
		}
		legacy.PortTestTradeUI(4, txptr(w.C()))
		rows = append(rows, o.shopCapture(t, len(rows), 0))
		legacy.PortTestTradeUI(5, txptr(w.C()), 16391, txptr(w.ChildByID(3604).C()), 0)
		rows = append(rows, o.shopCapture(t, len(rows), 0))
	}
	rows = append(rows, o.shopCapture(t, len(rows), 19, 100))
	rows = append(rows, o.shopCapture(t, len(rows), 15))
	shopUICapture(t, "asset-windows", rows, "67eb669972110888bf4157d370444250a43d0864fe2af97b6c4d66bc4309185a")
}
