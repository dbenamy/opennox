//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestCharacterCreationConstructionNotifications(t *testing.T) {
	// The resource parser sends numeric IDs before characterShowClass has stored
	// its root/host state. These are not addresses and must not be dereferenced.
	for _, id := range []uint{0, 1, 601, 602, 603, 605, 610, 0xffffffff} {
		if got := legacy.PortTestCharacterClassNewChild(id); got != 0 {
			t.Fatalf("child %d: got %d, want ignored event", id, got)
		}
	}
}
