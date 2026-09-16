//go:build porttest

package opennox

import (
	"os"
	"path/filepath"
	"testing"

	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/legacy"
)

func TestObjectXferFailedCallbackOwnership(t *testing.T) {
	core := newObjectXferOwner(t)
	typ := core.Types.ByID("PortInvisibleLight")
	// Rejection happens before payload access. No side-buffer allocation is needed
	// for this type, so the contract observes ownership of the object-pool slot.
	typ.Weight = 255
	parent := newObjectXferSimple(t, core)
	before := core.Objs.Alive
	var stream mapDrawableStream
	stream.u16(uint16(1000 + typ.Ind()))
	stream.u32(2)
	stream.u16(61)
	path := filepath.Join(t.TempDir(), "rejected.bin")
	if err := os.WriteFile(path, stream.Bytes(), 0600); err != nil {
		t.Fatal(err)
	}
	if err := cryptfile.OpenGlobal(path, cryptfile.ReadOnly, -1); err != nil {
		t.Fatal(err)
	}
	defer cryptfile.Close()
	if got := legacy.Nox_xxx_xfer_4F3E30(60, parent, 1); got != 0 {
		t.Fatal("bad child accepted")
	}
	if parent.InvFirstItem != nil || core.Objs.Alive != before {
		t.Fatalf("failed child retained: inventory=%p live=%d want=%d", parent.InvFirstItem, core.Objs.Alive, before)
	}
}

func TestObjectXferRejectedInventoryOwnership(t *testing.T) {
	core := newObjectXferOwner(t)
	t.Cleanup(noxflags.PortTestGameFlags(0))
	typ := core.Types.ByInd(1)
	typ.Weight = 255
	parent := core.NewObjectByTypeInd(1)
	child := core.NewObjectByTypeInd(1)
	if parent == nil || child == nil {
		t.Fatal("allocation")
	}
	parent.InvFirstItem = child
	child.InvHolder = parent
	typ.SetAllowed(false)
	// Placement owns and disposes both objects on admission rejection. No test
	// destructor frees them a second time; the real pool must return to zero live.
	if got := legacy.Nox_xxx_servMapLoadPlaceObj_4F3F50(parent, 0, nil); got != 0 {
		t.Fatal("disallowed placement accepted")
	}
	if core.Objs.Alive != 0 {
		t.Fatalf("rejected inventory live=%d want=0", core.Objs.Alive)
	}
}
