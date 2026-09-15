//go:build porttest

package opennox

import (
	"math"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func TestClientInventoryDisplayElementModifiers(t *testing.T) {
	o := newInventoryDisplayOwner(t)
	callbacks := legacy.PortTestInventoryDisplayModifierFunctions()
	var rows []inventoryDisplayResult
	for _, class := range []uint32{0, 1, 0x1000, 0x1000000, 0x2000000, 0x10000000, 0x13001000} {
		for _, pair := range [][2]int{{-1, -1}, {0, -1}, {-1, 0}, {1, -1}, {-1, 1}, {0, 1}, {1, 0}, {0, 0}, {1, 1}, {2, 0}, {2, 1}} {
			for _, value := range []float32{-2.5, 0, 0.125, 1, 123.75} {
				o.reset(t)
				dr := o.item(t, "Bow", 123)
				*txword(dr, 112) = class
				for slot, index := range pair {
					if index < 0 {
						continue
					}
					m := o.mods[slot+2]
					m.AttackPreHit52 = server.ModifierEffFnc{Fnc: callbacks[index], Valf: value + float32(slot)*0.25}
					*txword(dr, 440+uintptr(slot*4)) = uint32(uintptr(m.C()))
				}
				for op := 1; op <= 2; op++ {
					want := float64(0)
					if class&0x13001000 != 0 {
						for slot, index := range pair {
							if index == op-1 {
								want = float64(value + float32(slot)*0.25)
								break
							}
						}
					}
					r := o.call(t, len(rows), op, txptr(dr.C()), 0, 0)
					if r.Return != math.Float64bits(want) {
						t.Fatalf("class%x pair%v op%d got%g want%g", class, pair, op, math.Float64frombits(r.Return), want)
					}
					rows = append(rows, r)
				}
			}
		}
	}
	for op := 1; op <= 2; op++ {
		o.reset(t)
		r := o.call(t, len(rows), op, 0, 0, 0)
		if r.Return != 0 {
			t.Fatal("nil modifier result")
		}
		rows = append(rows, r)
	}
	inventoryDisplayCapture(t, "elements", rows, "a9b3b69ce9ccf826df3884a19a4ffd2c13ed0cc8604ed456315144c59b39232c")
}

func TestClientInventoryDisplayScaledDurability(t *testing.T) {
	o := newInventoryDisplayOwner(t)
	callbacks := legacy.PortTestInventoryDisplayModifierFunctions()
	values, free := alloc.Make([]float32{}, 2)
	defer free()
	for _, class := range []uint32{0, 0x1000, 0x1000000, 0x2000000, 0x10000000} {
		for _, callback := range []int{-1, 0, 2, 3} {
			for _, health := range [][2]uint16{{0, 0}, {0, 100}, {1, 3}, {123, 65535}, {65535, 65535}} {
				for _, scale := range []float32{-1, 0, 0.25, 0.5, 1, 1.125, 2} {
					o.reset(t)
					dr := o.item(t, "Bow", 123)
					*txword(dr, 112) = class
					*(*uint16)(unsafe.Add(dr.C(), 292)), *(*uint16)(unsafe.Add(dr.C(), 294)) = health[0], health[1]
					if callback >= 0 {
						m := o.mods[1]
						m.Defend76 = server.ModifierEffFnc{Fnc: callbacks[callback], Valf: scale}
						*txword(dr, 436) = uint32(uintptr(m.C()))
					}
					want := [2]float32{float32(health[0]), float32(health[1])}
					ret := o.norm(uint32(uintptr(dr.C())))
					if class != 0 && callback == 2 {
						for i := range want {
							want[i] = float32(int32(want[i] * scale))
						}
						ret = uint32(int32(want[1]))
					}
					r := o.call(t, 0, 5, txptr(dr.C()), txptr(unsafe.Pointer(&values[0])), txptr(unsafe.Pointer(&values[1])))
					if values[0] != want[0] || values[1] != want[1] || uint32(r.Return) != ret {
						t.Fatalf("class%x callback%d health%v scale%g got%v/%x want%v/%x", class, callback, health, scale, values, r.Return, want, ret)
					}
				}
			}
		}
	}
}

func TestClientInventoryDisplayScalarState(t *testing.T) {
	o := newInventoryDisplayOwner(t)
	for _, v := range []uint32{0, 1, 0x7fffffff, 0x80000000, 0xffffffff} {
		o.reset(t)
		r := o.call(t, 0, 6, uintptr(v), 0, 0)
		if uint32(r.Return) != v || memmap.Uint32(0x5D4594, 1050012) != v {
			t.Fatalf("scalar%x return%x state%x", v, r.Return, memmap.Uint32(0x5D4594, 1050012))
		}
	}
}
