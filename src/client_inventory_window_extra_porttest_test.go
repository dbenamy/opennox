//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"testing"
	"unsafe"
)

func TestClientInventoryWindowSliderRouting(t *testing.T) {
	o := newInventoryWindowOwner(t)
	var rows []inventoryWindowResult
	for _, op := range []int{21, 22} {
		for _, drag := range []bool{false, true} {
			for _, event := range []uintptr{5, 6, 8, 19, 20, 99} {
				o.reset(t)
				o.construct(t)
				o.openInput()
				slider := (*gui.Window)(unsafe.Pointer(uintptr(*o.windowWords["dword_5d4594_1062508"])))
				target := slider
				if op == 21 {
					target = slider.Field100()
				}
				if target == nil {
					t.Fatal("slider has no thumb")
				}
				if drag {
					o.stack(t, 0, 0, 1, "RedApple", 100)
					o.invoke(1, txptr(o.mainWindow().C()), 5, inventoryWindowPoint(339, 38), 0)
				}
				rows = append(rows, o.capture(t, len(rows), 0, op, txptr(target.C()), event, inventoryWindowPoint(339, 38), 0))
			}
		}
	}
	inventoryWindowCapture(t, "slider-routing", rows, "ec5719cb7c23180418e23708951905005bb463f03c1f8cb84135c7c3a9af82ab")
}

func TestClientInventoryWindowMapButton(t *testing.T) {
	o := newInventoryWindowOwner(t)
	var rows []inventoryWindowResult
	point, free := alloc.Make([]int32{}, 2)
	defer free()
	point[0], point[1] = 280, 130
	for _, state := range []byte{0, 1, 2, 3} {
		for _, busy := range []uint32{0, 1, 2} {
			for _, initial := range []uint32{0, 1, 6, 7} {
				o.reset(t)
				o.construct(t)
				*memmap.PtrUint8(0x5D4594, 1049868) = state
				*memmap.PtrUint32(0x5D4594, 1096672) = busy
				*memmap.PtrUint32(0x5D4594, 1096424) = initial
				rows = append(rows, o.capture(t, len(rows), 0, 10, txptr(o.mainWindow().C()), txptr(unsafe.Pointer(&point[0])), 0, 0))
				want := initial
				if state == 2 && busy != 1 {
					want ^= 1
				}
				if got := *memmap.PtrUint32(0x5D4594, 1096424); got != want {
					t.Fatalf("state%d busy%d map%d want%d", state, busy, got, want)
				}
			}
		}
	}
	inventoryWindowCapture(t, "map-button", rows, "55c0eba602930c5e940558e213aebedd0f0dded2abb2661dc611dfa5811d85a9")
}

func TestClientInventoryWindowPaperDollDispatch(t *testing.T) {
	o := newInventoryWindowOwner(t)
	var rows []inventoryWindowResult
	for _, gender := range []byte{0, 1} {
		for _, sameColor := range []bool{false, true} {
			o.reset(t)
			o.construct(t)
			o.openInput()
			p := unsafe.Pointer(&o.players[0])
			*(*uint32)(p) = 0
			*(*uint32)(unsafe.Add(p, 4)) = 0
			*(*byte)(unsafe.Add(p, 2252)) = gender
			for i, off := range []uintptr{2292, 2296, 2300, 2304, 2308, 2312} {
				*(*uint32)(unsafe.Add(p, off)) = uint32(0x112233 + i*0x10101)
			}
			if sameColor {
				*(*uint32)(unsafe.Add(p, 2296)) = *(*uint32)(unsafe.Add(p, 2292))
			}
			for i := 0; i < 4; i++ {
				*memmap.PtrPtr(0x973A20, 16+uintptr(i*4)) = unsafe.Pointer(o.images[i].C())
			}
			rows = append(rows, o.capture(t, len(rows), 0, 5, txptr(o.mainWindow().C()), 0, 0, 0))
		}
	}
	inventoryWindowCapture(t, "paperdoll-dispatch", rows, "bb883452586f528e1f36a570f3b065c57615f099b7071eaaf3f20d310b83fbc9")
}
