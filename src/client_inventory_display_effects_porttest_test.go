//go:build porttest

package opennox

import (
	"fmt"
	"testing"
	"unsafe"

	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func TestClientInventoryDisplayStatEffects(t *testing.T) {
	o := newInventoryDisplayOwner(t)
	pos, free := alloc.Make([]int32{}, 2)
	defer free()
	var rows []inventoryDisplayResult
	for _, language := range []int{0, 6, 8} {
		o.language(language)
		strMan = o.c.srv.Strings()
		for _, buff := range []uint32{0, 0x10, 0x200, 0x210} {
			for _, speed := range []float32{-0.002, -0.00005, 0, 0.00005, 0.002} {
				for _, weight := range []byte{0, 1, 255} {
					o.reset(t)
					noxflags.SetGame(2048)
					*memmap.PtrUint32(0x5D4594, 1062540) = buff
					*memmap.PtrFloat32(0x5D4594, 1063100) = speed
					*memmap.PtrFloat32(0x5D4594, 1062548) = 0.4567
					*memmap.PtrUint32(0x5D4594, 1062544) = 4321
					*(*uint16)(unsafe.Add(unsafe.Pointer(&o.players[0]), 3652)) = 50
					dr := o.item(t, "RedApple", 123)
					*(*byte)(unsafe.Add(dr.C(), 298)) = weight
					cells := legacy.PortTestInventoryCells()
					cells[0].Drawable, cells[0].Count = dr, 2
					cells[83].Drawable, cells[83].Count = dr, 3
					r := o.call(t, len(rows), 7, txptr(unsafe.Pointer(&pos[0])), 0, 0)
					expected := map[string]bool{"Experience 4321 / 400": false, fmt.Sprintf("%d / 50", int(weight)*5): false, "457 / 1000": false}
					for _, text := range r.Text {
						if _, ok := expected[text.Text]; ok {
							expected[text.Text] = true
						}
					}
					for text, seen := range expected {
						if !seen {
							t.Fatalf("missing stat label %q", text)
						}
					}
					rows = append(rows, r)
				}
			}
		}
	}
	inventoryDisplayCapture(t, "stat-effects", rows, "4b5ded6c8abb77d36a6f8b9ce9bcf2a902abcb24388ef5a1709cd19f32d21eb9")
}
