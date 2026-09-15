//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestClientWindowAliasedPositionContract(t *testing.T) {
	o := newWindowHelperOwner(t)
	o.create()
	// The C ABI permits the same storage for both output pointers. Its writes
	// occur in order: local Y=12, then parent's X=3 and Y=4, yielding19.
	ret, xy := legacy.PortTestWindowHelper(1, o.win, 0, 0, nil, 1)
	if ret != 0 || xy != [2]uint32{19, 0x89abcdef} {
		t.Fatalf("aliased position: ret%d xy%x", ret, xy)
	}
}
