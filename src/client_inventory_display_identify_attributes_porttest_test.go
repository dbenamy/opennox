//go:build porttest

package opennox

import (
	"image"
	"strings"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/client/gui"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func TestClientInventoryDisplayIdentifyAttributes(t *testing.T) {
	o := newInventoryDisplayOwner(t)
	callbacks := legacy.PortTestInventoryDisplayModifierFunctions()
	pos, free := alloc.Make([]int32{}, 2)
	defer free()
	var rows []inventoryDisplayResult
	for _, kind := range []int{0, 1, 2, 3} {
		for _, mask := range []int{0, 1, 2, 3, 4, 8, 15} {
			for _, detailed := range []bool{false, true} {
				o.reset(t)
				o.inventoryRects()
				_, icon, list := o.identifyWindows(t)
				o.c.Mouse = image.Pt(15, 22)
				if detailed {
					noxflags.SetGame(2048)
				}
				dr := o.item(t, "Bow", 123)
				*(*uint16)(unsafe.Add(dr.C(), 292)), *(*uint16)(unsafe.Add(dr.C(), 294)) = 33, 100
				*(*byte)(unsafe.Add(dr.C(), 298)) = 9
				o.weapon.DamageMin72 = 7
				o.weapon.DamageCoeffOrArmor64 = 0.125
				o.weapon.ReqStrength60 = 10
				o.armor.TypeInd = dr.TypeIDVal
				o.armor.DamageCoeffOrArmor64 = 0.12345
				switch kind {
				case 0:
					dr.ObjSubClass = 0
				case 1:
					dr.ObjClass = 0x2000000
					dr.ObjSubClass = 0
				case 2:
					dr.ObjClass = 0x2000000
					dr.ObjSubClass = 2
				case 3:
					dr.ObjSubClass = 2
				}
				typ := o.c.Things.TypeByID("Bow")
				typ.PrettyName = alloc.InternCString16("Practice bow")
				typ.PrettyImage = uint32(uintptr(o.images[9].C()))
				typ.Desc = alloc.InternCString16([]string{"A weathered bow.", "A spare training bow."}[mask%2])
				for i, m := range o.mods {
					*(*unsafe.Pointer)(unsafe.Add(m.C(), 8)) = unsafe.Pointer(alloc.InternCString16("Fine"))
					*(*unsafe.Pointer)(unsafe.Add(m.C(), 12)) = unsafe.Pointer(alloc.InternCString16("of Storms"))
					*(*unsafe.Pointer)(unsafe.Add(m.C(), 16)) = nil
					if i >= 2 {
						*(*unsafe.Pointer)(unsafe.Add(m.C(), 16)) = unsafe.Pointer(alloc.InternCString16([]string{"Burning attribute", "Sparking attribute"}[i-2]))
					}
					if mask&(1<<i) != 0 {
						*txword(dr, 432+uintptr(i*4)) = uint32(uintptr(m.C()))
					}
				}
				o.mods[0].Attack40 = server.ModifierEffFnc{Fnc: callbacks[4], Valf: 1.25}
				o.mods[0].Defend76 = server.ModifierEffFnc{Fnc: callbacks[3], Valf: 1.5}
				o.mods[1].Defend76 = server.ModifierEffFnc{Fnc: callbacks[2], Valf: 0.5}
				o.mods[2].AttackPreHit52 = server.ModifierEffFnc{Fnc: callbacks[1], Valf: 2.5}
				o.mods[3].AttackPreHit52 = server.ModifierEffFnc{Fnc: callbacks[0], Valf: 1.25}
				*o.displayWords["dword_5d4594_1063116"] = uint32(uintptr(dr.C()))
				r := o.call(t, len(rows), 3, txptr(unsafe.Pointer(&pos[0])), 0, 0)
				d := (*gui.ScrollListBoxData)(list.WidgetData)
				for _, item := range unsafe.Slice(d.Items, int(d.Count)) {
					if item.Text[0] != 0 {
						r.WidgetText = append(r.WidgetText, alloc.GoString16(&item.Text[0]))
					}
				}
				joined := strings.Join(r.WidgetText, "\n")
				for bit, text := range map[int]string{4: "Burning attribute", 8: "Sparking attribute"} {
					if strings.Contains(joined, text) != (mask&bit != 0) {
						t.Fatalf("kind%d mask%x missing or extra attribute %q: %v", kind, mask, text, r.WidgetText)
					}
				}
				if !strings.Contains(joined, "Weight 9") || !strings.Contains(joined, []string{"A weathered bow.", "A spare training bow."}[mask%2]) {
					t.Fatalf("missing item weight/description: %v", r.WidgetText)
				}
				if uint32(uintptr(icon.DrawData().BgImageHnd)) != typ.PrettyImage {
					t.Fatal("identify image")
				}
				rows = append(rows, r)
				r = o.call(t, len(rows), 13, txptr(icon.C()), txptr(icon.DrawData().C()), 0)
				if r.Return != 1 {
					t.Fatal("identify image draw return")
				}
				rows = append(rows, r)
			}
		}
	}
	inventoryDisplayCapture(t, "identify-attributes", rows, "28f961342a02c55af14747d007408e4446a67cc642d91ac31ed3ace7d0157386")
}
