//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestClientTradeUIQuantityDestroyRecreate(t *testing.T) {
	o := newTradeUIOwner(t)
	o.reset(t)
	o.construct(t)
	for cycle := 0; cycle < 100; cycle++ {
		o.initAmount(t)
		legacy.PortTestTradeUI(7, o.quantityArgs(200, 100, 32)...)
		if *o.windowWords["dword_5d4594_1319268"] != 1 {
			t.Fatal("show recreated quantity dialog")
		}
		legacy.PortTestTradeUI(6)
		o.amountReady = false
		// Match the game frame: reclaim deferred GUI destruction before recreation.
		o.c.GUI.FreeDestroyed()
		for _, n := range []string{"nox_gui_itemAmount_item_1319256", "dword_5d4594_1319232", "dword_5d4594_1319236", "dword_5d4594_1319268"} {
			if *o.windowWords[n] != 0 {
				t.Fatalf("destroy cycle%d retains %s", cycle, n)
			}
		}
		if legacy.Get_nox_gui_itemAmount_dialog_1319228() != nil {
			t.Fatal("destroy retained amount window")
		}
		if legacy.PortTestTradeUI(2) != 0 {
			t.Fatal("cancel destroyed amount dialog")
		}
	}
}
