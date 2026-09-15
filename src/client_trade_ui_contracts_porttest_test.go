//go:build porttest

package opennox

import (
	"image"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func TestClientTradeUIConstruction(t *testing.T) {
	o := newTradeUIOwner(t)
	o.reset(t)
	o.construct(t)
	o.constructTrade(t)
	for _, id := range []uint{3702, 3703, 3704, 3705, 3708, 3709, 3710, 3711, 3712, 3713} {
		if o.tradeWindow().ChildByID(id) == nil {
			t.Fatalf("missing trade child %d", id)
		}
	}
	if got := o.tradeWindow().DrawData().Tooltip(); got != "Trade" {
		t.Fatalf("trade tooltip %q", got)
	}
	if !o.tradeWindow().Flags.IsHidden() {
		t.Fatal("new trade window is visible")
	}
	for side, cells := range legacy.PortTestTradeUICells() {
		for i, cell := range cells {
			if cell.Drawable != nil || cell.Count != 0 {
				t.Fatalf("uncleared trade cell %d/%d", side, i)
			}
		}
	}
	r := o.tradeCapture(t, 0, 16)
	tradeUICapture(t, "construction", []tradeUIResult{r}, "5dfb48677d9837083b527e5e0f575229b7df02ae641a310bb4670de5446839b4")
}

func TestClientTradeUIGridRegions(t *testing.T) {
	o := newTradeUIOwner(t)
	o.reset(t)
	o.construct(t)
	o.constructTrade(t)
	point, free := alloc.Make([]int32{}, 2)
	defer free()
	grids := legacy.PortTestTradeUICells()
	checks := 0
	for _, offset := range []image.Point{{0, 0}, {-20, 30}} {
		o.tradeWindow().SetPos(offset)
		for side, op := range []int{13, 18} {
			id := uint(3704 + side)
			origin := o.tradeWindow().ChildByID(id).GlobalPos()
			for y := -1; y <= 101; y++ {
				for x := -1; x <= 101; x++ {
					point[0], point[1] = int32(origin.X+x), int32(origin.Y+y)
					got := legacy.PortTestTradeUI(op, txptr(unsafe.Pointer(&point[0])))
					want := uint32(0)
					if x >= 0 && x <= 100 && y >= 0 && y <= 100 {
						col, row := 0, 0
						if x > 50 {
							col = 1
						}
						if y > 50 {
							row = 1
						}
						want = uint32(txptr(unsafe.Pointer(&grids[side][row+2*col])))
					}
					if got != want {
						t.Fatalf("side%d offset%v point%d,%d got%#x want%#x", side, offset, x, y, got, want)
					}
					checks++
				}
			}
		}
	}
	t.Logf("%d exact grid pointer/boundary contracts", checks)
}

func TestClientTradeUIEmptyAndOutsideHover(t *testing.T) {
	o := newTradeUIOwner(t)
	o.reset(t)
	o.construct(t)
	o.constructTrade(t)
	for _, pos := range []image.Point{{0, 0}, {150, 60}, {500, 400}, {320, 220}, {32767, 32767}, {10, 40}, {60, 90}, {180, 40}, {280, 140}} {
		if got := legacy.PortTestTradeUI(17, txptr(o.tradeWindow().C()), 0, inventoryWindowPoint(pos.X, pos.Y)); got != 1 {
			t.Fatalf("hover %v return%d", pos, got)
		}
		if len(o.events) != 3 {
			t.Fatalf("empty/outside hover changed drawable ownership: %v", o.events)
		}
	}
}
