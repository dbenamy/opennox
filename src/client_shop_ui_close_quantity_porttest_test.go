//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"testing"
	"unsafe"
)

func TestClientShopUICloseDuringQuantity(t *testing.T) {
	o := newShopUIOwner(t)
	var rows []shopUIResult
	p, free := alloc.New([2]int32{})
	defer free()
	for _, closeOp := range []int{15, 13} {
		for _, kind := range []int{32, 40, 41, 0} {
			o.reset(t)
			o.construct(t)
			o.constructShop(t)
			o.initAmount(t)
			*o.windowWords["dword_5d4594_1098624"] = 1
			typ := uintptr(o.c.Things.TypeByID("RedApple").Index())
			switch kind {
			case 32:
				c := &legacy.PortTestShopUICells()[0]
				c.Drawable = o.item(t, "RedApple", 123)
				c.Count = 2
				c.Codes[0], c.Codes[1] = 123, 124
				pos := o.shopWindow().ChildByID(3806).GlobalPos()
				*p = [2]int32{int32(pos.X + 25), int32(pos.Y + 25)}
				legacy.PortTestShopUI(32, txptr(unsafe.Pointer(p)))
			case 40, 41:
				o.stack(t, 0, 0, 2, "RedApple", 123)
				legacy.PortTestShopUI(kind, 123, 0)
			case 0:
				legacy.PortTestTradeUI(9, 0, 0)
				legacy.PortTestTradeUI(7, txptr(unsafe.Pointer(alloc.InternCString16("Unrelated inventory quantity"))), 200, 100, 123, typ, 0, 1, 0, 0, 0)
			}
			before := *o.windowWords["nox_gui_itemAmount_item_1319256"]
			if before == 0 || *o.windowWords["dword_5d4594_1319268"] != 1 {
				t.Fatal("quantity setup failed")
			}
			dr := (*client.Drawable)(unsafe.Pointer(uintptr(before)))
			rows = append(rows, o.shopCapture(t, len(rows), closeOp))
			if closeOp == 13 {
				o.shopReady = false
			}
			if kind == 0 {
				if *o.windowWords["nox_gui_itemAmount_item_1319256"] != before || *o.windowWords["dword_5d4594_1319268"] != 1 {
					t.Errorf("close%d altered another owner's quantity dialog", closeOp)
				}
			} else {
				if *o.windowWords["nox_gui_itemAmount_item_1319256"] != 0 || *o.windowWords["dword_5d4594_1319268"] != 0 || o.objects[o.identities[dr]].Live {
					t.Errorf("close%d kind%d retained shop quantity drawable/dialog", closeOp, kind)
				}
				if *o.shopWords["dword_5d4594_1098616"] != 0 || *o.shopWords["dword_5d4594_1098620"] != 0 {
					t.Errorf("close%d kind%d retained pending shop request", closeOp, kind)
				}
			}
		}
	}
	shopUICapture(t, "close-quantity", rows, "7be4a66be4354d0127737e77b8661254dea295c17f4ff14d969dd56c019bbecd")
}
