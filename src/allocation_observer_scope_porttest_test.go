//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestAllocationObserverThreadIsolation(t *testing.T) {
	for i := 0; i < 3; i++ {
		got := legacy.PortTestAllocationObserverIsolation()
		want := [8]int{1, 0, 1, 1, 1, 1, 1, 0}
		if got != want {
			t.Fatalf("foreign allocation/count, own failure/success/count/free/valid, tracker delta: got %v want %v", got, want)
		}
	}
}
