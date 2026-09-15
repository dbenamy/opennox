//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/memmap/nox/blobdata"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"testing"
	"unsafe"
)

func (o *inventoryDisplayOwner) inventoryRects() {
	copy(unsafe.Slice((*byte)(memmap.PtrOff(0x587000, 136192)), 192), blobdata.PortTestInventoryDisplayGeometry())
}
func TestClientInventoryDisplayItemTooltips(t *testing.T) {
	o := newInventoryDisplayOwner(t)
	pos, free := alloc.Make([]int32{}, 2)
	defer free()
	var rows []inventoryDisplayResult
	for _, scroll := range []uint32{0, 49, 50, 800, 850} {
		for _, mode := range []byte{0, 1} {
			for _, point := range [][2]int32{{-1, -1}, {211, 216}, {263, 13}, {264, 13}, {265, 14}, {313, 62}, {314, 13}, {315, 14}, {363, 62}, {364, 63}, {513, 212}, {514, 213}, {515, 214}, {20, 20}} {
				o.reset(t)
				o.inventoryRects()
				*memmap.PtrUint8(0x5D4594, 1049869) = mode
				*o.displayWords["dword_5d4594_1062512"] = scroll
				dr := o.item(t, "RedApple", 999)
				for i := range legacy.PortTestInventoryCells() {
					c := &legacy.PortTestInventoryCells()[i]
					c.Drawable, c.Count, c.Codes[0] = dr, 1, uint32(100+i)
				}
				pos[0], pos[1] = point[0], point[1]
				r := o.call(t, len(rows), 11, txptr(o.parent.C()), txptr(unsafe.Pointer(&pos[0])), 0)
				if (point == [2]int32{-1, -1} || point == [2]int32{515, 214}) && r.Return != 0 {
					t.Fatal("outside item tooltip")
				}
				rows = append(rows, r)
			}
		}
	}
	inventoryDisplayCapture(t, "item-tooltips", rows, "5aa276565fae6c6b45bf3165198bf78f317a2f730a9e6b026dcfb5f24d3bc0fb")
}
