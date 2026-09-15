//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"os"
	"path/filepath"
	"testing"
	"unsafe"
)

func (o *shopUIOwner) installMessageDialog(t *testing.T) {
	o.messageOwner = true
	t.Helper()
	old, oldParent, oldFocus, oldCapture := nox_gui_curDialog_830224, dword_5d4594_830228, dword_5d4594_830232, dword_5d4594_830236
	oldAccept, oldCancel := func_5D4594_830220, func_5d4594_830216
	flags := *memmap.PtrUint32(0x5D4594, 830240)
	nox_gui_curDialog_830224, dword_5d4594_830228, dword_5d4594_830232, dword_5d4594_830236 = nil, nil, nil, nil
	func_5D4594_830220, func_5d4594_830216 = nil, nil
	t.Cleanup(func() {
		sub_44A400()
		nox_gui_curDialog_830224, dword_5d4594_830228, dword_5d4594_830232, dword_5d4594_830236 = old, oldParent, oldFocus, oldCapture
		func_5D4594_830220, func_5d4594_830216 = oldAccept, oldCancel
		*memmap.PtrUint32(0x5D4594, 830240) = flags
	})
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, "window"), 0700); err != nil {
		t.Fatal(err)
	}
	s := "FONT = small; WINDOW 4000 0 0 200 160 USER; STATUS = ENABLED+ABOVE; CHILD "
	for id := 4001; id <= 4010; id++ {
		typ, data := "PUSHBUTTON", "STYLE = MOUSETRACK;"
		if id == 4003 || id == 4008 {
			typ, data = "ENTRYFIELD", "DATA = 10 100 0 2;"
		}
		if id == 4004 || id == 4005 {
			typ, data = "STATICTEXT", "DATA = 1 0 WindowDir:Blank;"
		}
		s += fmt.Sprintf("WINDOW %d 5 %d 180 20 %s; STATUS = ENABLED+HIDDEN; %s END ", id, (id-4001)*10, typ, data)
	}
	s += "END END"
	if err := os.WriteFile(filepath.Join(dir, "window", "dlg.wnd"), []byte(s), 0600); err != nil {
		t.Fatal(err)
	}
	t.Chdir(dir)
}

func TestClientShopUINotEnoughGold(t *testing.T) {
	o := newShopUIOwner(t)
	o.installMessageDialog(t)
	var rows []shopUIResult
	for _, amount := range []uintptr{0, 1, 7, 32768, 0x7fffffff, 0x80000000, 0xffffffff} {
		sub_44A400()
		o.reset(t)
		o.construct(t)
		o.constructShop(t)
		rows = append(rows, o.shopCapture(t, len(rows), 38, amount))
		if nox_gui_curDialog_830224 == nil {
			t.Fatal("actual shop information dialog did not open")
		}
		want := fmt.Sprintf("Need %d more gold.", int32(amount))
		if got := tradeQuantityText(nox_gui_curDialog_830224, 4004); got != want {
			t.Fatalf("information text %q want%q", got, want)
		}
		if got := tradeQuantityText(nox_gui_curDialog_830224, 4005); got != "Shop information" {
			t.Fatalf("information title %q", got)
		}
		nox_gui_curDialog_830224.Func94(&WindowEvent0x4007{Win: nox_gui_curDialog_830224.ChildByID(4001)})
		if nox_gui_curDialog_830224 != nil {
			t.Fatal("OK did not close message")
		}
	}
	shopUICapture(t, "gold-dialog", rows, "cfd53b45d46df280376abed4a2e144db637a044e1f215ce20b13b416517dddb1")
}

func TestClientShopUIQuantityAllocationFailure(t *testing.T) {
	o := newShopUIOwner(t)
	for _, op := range []int{40, 41} {
		o.reset(t)
		o.construct(t)
		o.constructShop(t)
		o.initAmount(t)
		c := &legacy.PortTestInventoryCells()[0]
		c.Drawable = o.item(t, "RedApple", 123)
		c.Count = 2
		c.Codes[0], c.Codes[1] = 123, 124
		// The quantity owner preserves an existing dialog on failed allocation;
		// opening a new one must not leave shop requests permanently pending.
		o.exhausted(t, func() { legacy.PortTestShopUI(op, 123, 0) })
		name := "dword_5d4594_1098616"
		if op == 41 {
			name = "dword_5d4594_1098620"
		}
		if *o.shopWords[name] != 0 {
			t.Errorf("operation%d retained pending flag after quantity allocation failure", op)
		}
		if *o.windowWords["dword_5d4594_1319268"] != 0 {
			t.Error("failed quantity allocation opened dialog")
		}
		legacy.PortTestShopUI(op, 123, 0)
		if *o.shopWords[name] != 1 || *o.windowWords["dword_5d4594_1319268"] != 1 {
			t.Errorf("operation%d could not retry after allocation recovered", op)
		}
		legacy.PortTestTradeUI(2)
		if *o.shopWords[name] != 0 {
			t.Errorf("operation%d cancel retained pending flag", op)
		}

	}
}

func TestClientShopUIBuyPriceOverflow(t *testing.T) {
	o := newShopUIOwner(t)
	o.installMessageDialog(t)
	p, free := alloc.New([2]int32{})
	defer free()
	for _, price := range []uint32{0, 1, 0x08000000, 0x80000000, 0xffffffff} {
		sub_44A400()
		o.reset(t)
		o.construct(t)
		o.constructShop(t)
		o.initAmount(t)
		c := &legacy.PortTestShopUICells()[0]
		c.Drawable = o.item(t, "RedApple", 123)
		c.Count = 32
		c.Value = price
		for j := range c.Codes {
			c.Codes[j] = uint32(123 + j)
		}
		// Fixture player gold is explicitly zero: free goods allow the full stack;
		// every positive unit price is unaffordable regardless of product overflow.
		*o.windowWords["dword_5d4594_1062552"] = 0
		pos := o.shopWindow().ChildByID(3806).GlobalPos()
		*p = [2]int32{int32(pos.X + 25), int32(pos.Y + 25)}
		legacy.PortTestShopUI(32, txptr(unsafe.Pointer(p)))
		if price == 0 {
			if *o.windowWords["dword_5d4594_1319248"] != 32 || *o.windowWords["dword_5d4594_1319268"] != 1 {
				t.Error("free stock did not offer full stack")
			}
		} else {
			if *o.windowWords["dword_5d4594_1319268"] != 0 {
				t.Errorf("price%#x opened quantity despite zero gold", price)
			}
			if nox_gui_curDialog_830224 == nil {
				t.Errorf("price%#x did not show insufficient gold", price)
			}
		}
	}
}
