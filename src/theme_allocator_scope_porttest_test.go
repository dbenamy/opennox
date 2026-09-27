//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestThemeAllocatorObserverThreadScope(t *testing.T) {
	for i := 0; i < 3; i++ {
		if got := legacy.PortTestThemeObserverThreadScope(); got != [4]int{0, 1, 0, 0} {
			t.Fatalf("observer state before/fixture/foreign/after: %v", got)
		}
	}
}
