//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestClientShopUIModes(t *testing.T) {
	o := newShopUIOwner(t)
	var rows []shopUIResult
	for _, active := range []uint32{0, 1, 2} {
		for _, old := range []uint32{0, 1, 2, 3, 4, 0xffffffff} {
			for _, next := range []uintptr{0, 1, 2, 3, 4, 0xffffffff} {
				o.reset(t)
				o.construct(t)
				o.constructShop(t)
				*o.windowWords["dword_5d4594_1098624"], *o.windowWords["dword_5d4594_1098628"] = active, old
				if got := legacy.PortTestShopUI(0); got != active {
					t.Fatal("active getter")
				}
				if got := legacy.PortTestShopUI(21); got != old {
					t.Fatal("mode getter")
				}
				if got := legacy.PortTestShopUI(29); (got != 0) != (old == 2) {
					t.Fatal("buy-mode getter")
				}
				rows = append(rows, o.shopCapture(t, len(rows), 22, next))
				if got := legacy.PortTestShopUI(21); got != uint32(next) {
					t.Fatalf("mode transition%d -> %d returned%d", old, next, got)
				}
			}
		}
	}
	shopUICapture(t, "modes", rows, "2b088fe335d717b7e31b737ed30fb76b15d2f3332b92e2ab8afa825503bf12e7")
}
func TestClientShopUICancelRequest(t *testing.T) {
	o := newShopUIOwner(t)
	var rows []shopUIResult
	for _, active := range []uint32{0, 1, 2} {
		for _, mode := range []uint32{0, 1, 2, 3, 4} {
			o.reset(t)
			o.construct(t)
			o.constructShop(t)
			*o.windowWords["dword_5d4594_1098624"], *o.windowWords["dword_5d4594_1098628"] = active, mode
			r := o.shopCapture(t, len(rows), 1)
			rows = append(rows, r)
			m := shopTradeMessages(r)
			if active == 0 {
				if r.Return != 0 || len(m) != 0 {
					t.Fatal("inactive cancel requested close")
				}
			} else if r.Return != 1 || len(m) != 1 || len(m[0]) != 2 || m[0][1] != 18 {
				t.Fatalf("cancel request %v return%d", m, r.Return)
			}
		}
	}
	shopUICapture(t, "cancel-request", rows, "75314d8d499d2219e4c6be20717266c70c6f289c384addf9581f510af2b0f32c")
}
func TestClientShopUIScroll(t *testing.T) {
	o := newShopUIOwner(t)
	var rows []shopUIResult
	for _, button := range []uint{3808, 3809} {
		for _, scroll := range []uint32{0, 1, 49, 50, 125, 299, 300} {
			o.reset(t)
			o.construct(t)
			o.constructShop(t)
			*o.windowWords["dword_5d4594_1107036"] = scroll
			w := o.shopWindow()
			r := o.shopCapture(t, len(rows), 5, txptr(w.C()), 16391, txptr(w.ChildByID(button).C()), 0)
			rows = append(rows, r)
			want := int(scroll)
			if button == 3808 {
				want -= 50
				if want < 0 {
					want = 0
				}
				want = want / 50 * 50
			} else {
				want += 50
				if want > 300 {
					want = 300
				} else {
					want = want / 50 * 50
				}
			}
			if got := *o.windowWords["dword_5d4594_1107036"]; got != uint32(want) {
				t.Fatalf("button%d scroll%d -> %d want%d", button, scroll, got, want)
			}
		}
	}
	for _, value := range []uintptr{0, 1, 125, 299, 300, 301, 0xffffffff} {
		o.reset(t)
		o.construct(t)
		o.constructShop(t)
		w := o.shopWindow()
		rows = append(rows, o.shopCapture(t, len(rows), 5, txptr(w.C()), 16393, 0, value))
		if got := *o.windowWords["dword_5d4594_1107036"]; got != 300-uint32(value) {
			t.Fatalf("slider value%d scroll%d", value, got)
		}
	}
	shopUICapture(t, "scroll", rows, "095b221e2b2afb94e80d169e3733c94899148647bd63175d12d1edd8b7dc4bfe")
}
func TestClientShopUIPanelModes(t *testing.T) {
	o := newShopUIOwner(t)
	var rows []shopUIResult
	for _, active := range []uint32{0, 1, 2} {
		for _, mode := range []uint32{1, 2, 3, 4} {
			for _, id := range []uint{3801, 3802, 3803, 3804, 3810} {
				o.reset(t)
				o.construct(t)
				o.constructShop(t)
				*o.windowWords["dword_5d4594_1098624"], *o.windowWords["dword_5d4594_1098628"] = active, mode
				w := o.shopWindow()
				r := o.shopCapture(t, len(rows), 5, txptr(w.C()), 16391, txptr(w.ChildByID(id).C()), 0)
				rows = append(rows, r)
				want := mode
				if active != 0 {
					switch id {
					case 3802:
						want = 3
					case 3803:
						want = 4
					case 3804:
						want = 2
					}
				}
				if got := *o.windowWords["dword_5d4594_1098628"]; got != want {
					t.Fatalf("button%d mode%d active%d -> %d want%d", id, mode, active, got, want)
				}
			}
		}
	}
	shopUICapture(t, "panel-modes", rows, "0195f932ae8177a9e2efd6039f1ea9298e67ce6259c82ddcefebb4a30e14e681")
}
func TestClientShopUIMouseEvents(t *testing.T) {
	o := newShopUIOwner(t)
	var rows []shopUIResult
	for _, mode := range []uint32{0, 1, 2, 3, 4} {
		for event := 0; event <= 20; event++ {
			for _, inside := range []bool{false, true} {
				o.reset(t)
				o.construct(t)
				o.constructShop(t)
				*o.windowWords["dword_5d4594_1098628"] = mode
				p := o.shopWindow().ChildByID(3806).GlobalPos()
				p.X += 25
				p.Y += 25
				if !inside {
					p.X = 0
					p.Y = 0
				}
				rows = append(rows, o.shopCapture(t, len(rows), 6, txptr(o.shopWindow().C()), uintptr(event), inventoryWindowPoint(p.X, p.Y)))
			}
		}
	}
	shopUICapture(t, "mouse-events", rows, "23e23b26a7deaaa8c33059888e39acd46196a5b17881e5e912af950b47db665b")
}
func TestClientShopUIMissingResource(t *testing.T) {
	o := newShopUIOwner(t)
	o.reset(t)
	o.construct(t)
	o.missingShop = true
	r := o.shopCapture(t, 0, 4)
	if r.Return != 0 || o.shopWindow() != nil {
		t.Fatal("missing shop resource succeeded")
	}
	shopUICapture(t, "missing-resource", []shopUIResult{r}, "077a0e9ea4be91f2722c5f27503563fd934caa47ce8c49a881a7a3404c82b993")
}
