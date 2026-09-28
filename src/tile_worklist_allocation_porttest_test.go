//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestTileWorklistAllocationCleanup(t *testing.T) {
	tracked := legacy.PortTestAllocationUsesTracker()
	// Fail the outer table and each of its 128 rows; 130 is beyond the last call.
	for failAt := 0; failAt <= 130; failAt++ {
		r := legacy.PortTestWorklistAllocation(failAt)
		failure := failAt >= 1 && failAt <= 129
		count := 129
		if failure {
			count = failAt - 1
		}
		if r.Success == failure || !r.Zero || !r.ValidFrees || r.Remaining != 0 || r.Freed != count || len(r.Sizes) != count {
			t.Fatalf("allocation %d: %+v", failAt, r)
		}
		for i, size := range r.Sizes {
			want := 128 * 44
			if i == 0 {
				want = 128 * 4
			}
			if size != want {
				t.Fatalf("allocation %d block %d: size %d want %d", failAt, i, size, want)
			}
		}
		delta := 0
		if tracked {
			delta = 129
			if failure {
				delta = 1
			}
		}
		if r.NilTracked != (tracked && failure) || r.TrackedDuring-r.TrackedBefore != delta || r.TrackedAfter != r.TrackedBefore {
			t.Fatalf("allocation %d tracking: %+v", failAt, r)
		}
	}
}
