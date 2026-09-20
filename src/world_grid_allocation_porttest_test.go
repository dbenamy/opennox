//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestWorldGridAllocationOwnership(t *testing.T) {
	var rows []legacy.PortTestWorldGridAllocation
	for _, failAt := range []int{0, 1, 2, 3, 65, 129, 130} {
		r := legacy.PortTestWorldGridAllocate(failAt)
		wantRows := 128
		wantRet := 1
		if failAt > 0 && failAt <= 129 {
			wantRet = 0
			wantRows = failAt - 2
			if wantRows < 0 {
				wantRows = 0
			}
		}
		if r.Return != wantRet || r.Rows != wantRows || r.RowsFreed != wantRows || !r.Zero || !r.ValidFrees || r.FinalRemaining != 0 {
			t.Fatalf("grid allocation %+v", r)
		}
		if failAt == 1 {
			if len(r.Sizes) != 0 || r.Remaining != 0 {
				t.Fatalf("outer failure %+v", r)
			}
		} else {
			if len(r.Sizes) != wantRows+1 || r.Sizes[0] != 128*4 || !r.OuterRetained || r.Remaining != 1 {
				t.Fatalf("outer ownership %+v", r)
			}
			for _, sz := range r.Sizes[1:] {
				if sz != 128*44 {
					t.Fatalf("row allocation size%d", sz)
				}
			}
		}
		rows = append(rows, r)
	}
	drawableStateCapture(t, "world-grid-allocation", rows, "1ee00c2f2c174e339681f954e1db52615f75ff9408df434e2787c6e11c952dd8")
}
