//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"testing"
	"unsafe"
)

func TestClientShopUIDrawModes(t *testing.T) {
	o := newShopUIOwner(t)
	var rows []shopUIResult
	p, free := alloc.New([2]int32{})
	defer free()
	for _, mode := range []uint32{0, 1, 2, 3, 4, 5} {
		for _, offset := range [][2]int32{{0, 0}, {17, 23}} {
			o.reset(t)
			o.construct(t)
			o.constructShop(t)
			legacy.PortTestShopUI(37, 0, txptr(unsafe.Pointer(alloc.InternCString("FixtureGreeting"))), 0)
			*o.windowWords["dword_5d4594_1098628"] = mode
			rows = append(rows, o.shopCapture(t, len(rows), 8))
			*p = offset
			for _, op := range []int{9, 35, 36} {
				rows = append(rows, o.shopCapture(t, len(rows), op, txptr(unsafe.Pointer(p))))
			}
		}
	}
	shopUICapture(t, "draw-modes", rows, "0278ece430375cb4cda9120053ad5d10ae84dfa7326ef7af4e5f276e5a664b97")
}
func TestClientShopUIGridDrawHover(t *testing.T) {
	o := newShopUIOwner(t)
	var rows []shopUIResult
	for _, scroll := range []uint32{0, 1, 49, 50, 125, 299, 300, 450, 499} {
		for _, gold := range []uint32{0, 31, 0xffffffff} {
			o.reset(t)
			o.construct(t)
			o.constructShop(t)
			*o.windowWords["dword_5d4594_1098628"] = 2
			*o.windowWords["dword_5d4594_1107036"], *o.windowWords["dword_5d4594_1062552"] = scroll, gold
			cells := legacy.PortTestShopUICells()
			for i := range cells {
				if i%7 == 0 {
					continue
				}
				c := &cells[i]
				c.Drawable = o.item(t, "RedApple", uint32(100+i))
				c.Count = uint32(1 + i%32)
				c.Codes[0] = uint32(900 + i)
				c.Value = []uint32{0, 31, 32, 0xffffffff}[i%4]
			}
			rows = append(rows, o.shopCapture(t, len(rows), 10))
			w := o.shopWindow().ChildByID(3806)
			pos, size := w.GlobalPos(), w.Size()
			for _, off := range [][2]int{{0, 0}, {25, 25}, {49, 49}, {50, 50}, {299, 199}, {size.X, size.Y}} {
				x, y := pos.X+off[0], pos.Y+off[1]
				rows = append(rows, o.shopCapture(t, len(rows), 11, txptr(w.C()), 0, inventoryWindowPoint(x, y)))
				col, row := off[0]/50, (off[1]+int(scroll))/50
				if col > 5 {
					col = 5
				}
				if row > 9 {
					row = 9
				}
				c := &cells[col*10+row]
				if c.Count != 0 && c.Drawable.NetCode32 != c.Codes[0] {
					t.Fatal("shop hover did not select first code")
				}
			}
		}
	}
	shopUICapture(t, "grid-draw-hover", rows, "44463ee20cabfb2527c7ce103ad7bec568936bd5f057ccfe4a3b8b516f2b290c")
}
