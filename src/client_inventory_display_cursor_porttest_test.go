//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"image"
	"testing"
	"unsafe"
)

func TestClientInventoryDisplayIdentifyCursor(t *testing.T) {
	o := newInventoryDisplayOwner(t)
	pos, free := alloc.Make([]int32{}, 2)
	defer free()
	var rows []inventoryDisplayResult
	for _, v := range []struct {
		pos    image.Point
		cursor uint32
	}{{image.Pt(315, 14), 6}, {image.Pt(254, 14), 7}, {image.Pt(254, 64), 6}, {image.Pt(254, 114), 7}, {image.Pt(12, 16), 0}, {image.Pt(-1, -1), 7}} {
		o.reset(t)
		o.inventoryRects()
		o.identifyWindows(t)
		o.c.Mouse = v.pos.Add(image.Pt(4, 6))
		r := o.call(t, len(rows), 3, txptr(unsafe.Pointer(&pos[0])), 0, 0)
		if uint32(o.c.Cursor) != v.cursor {
			t.Fatalf("identify mouse%v cursor%d want%d", v.pos, o.c.Cursor, v.cursor)
		}
		rows = append(rows, r)
	}
	inventoryDisplayCapture(t, "identify-cursor", rows, "3a86c59824cfae227add1f69081a847b4b27a0ad464b48841e1ae52de66ca80f")
}
