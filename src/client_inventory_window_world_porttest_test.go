//go:build porttest && !server

package opennox

import (
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/common/memmap"
	"image"
	"testing"
)

// The production world finder is client-only. Use its actual spatial index and
// visibility owner rather than substituting a finder in server builds.
func TestClientInventoryWindowWorldSelection(t *testing.T) {
	o := newInventoryWindowOwner(t)
	var rows []inventoryWindowResult
	for _, distance := range []int{-1, 0, 74, 75, 76, 100} {
		o.reset(t)
		o.construct(t)
		o.openInput()
		*o.windowWords["dword_5d4594_1049864"] = 4
		pos := image.Pt(320, 240)
		var indexed *client.Drawable
		if distance >= 0 {
			pos.X += distance
			dr := o.item(t, "RedApple", 100)
			dr.SetPos(pos)
			*txword(dr, 276) = 0
			o.c.Objs.AddIndex2D(dr)
			indexed = dr
			if got := o.c.Nox_drawable_find(pos, 20); got != dr {
				t.Fatal("actual world finder did not select fixture drawable")
			}
		} else if o.c.Nox_drawable_find(pos, 20) != nil {
			t.Fatal("empty world unexpectedly contains selectable object")
		}
		ret := o.invoke(1, txptr(o.mainWindow().C()), 6, inventoryWindowPoint(pos.X, pos.Y), 0)
		// The shared drawable snapshot expects unlinked objects. Remove only the
		// fixture's index membership after the real selection, before recording.
		if indexed != nil {
			o.c.Objs.Index2DRemove(indexed, indexed.GetExt())
		}
		if ret != 1 {
			t.Fatalf("selection return%d", ret)
		}
		rows = append(rows, o.capture(t, distance, 0, 7, 0, 0, 0, 0))
		if *o.windowWords["dword_5d4594_1049864"] != 0 || *memmap.PtrUint32(0x5D4594, 1049848) != 0 {
			t.Fatal("world selection retained selection mode or began drag")
		}
	}
	inventoryWindowCapture(t, "world-selection", rows, "3eb2bf33a5bb83a7c9f5b024df2e59345f415934c4e0d8f3fdd73c327e01ba7b")
}
