//go:build porttest

package opennox

import (
	"testing"
	"unsafe"

	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func newObjectXferOwner(t *testing.T) *server.Server {
	t.Helper()
	t.Cleanup(noxflags.PortTestGameFlags(0x200000))
	core := new(server.Server)
	owners := core.PortTestMapPaintingOwners([2]unsafe.Pointer{})
	owners.ObjectXferTypes()
	old, oldGet := noxServer, legacy.GetServer
	wrapped := &Server{Server: core}
	core.ExtServer = unsafe.Pointer(wrapped)
	noxServer = wrapped
	legacy.GetServer = func() legacy.Server { return wrapped }
	oldTOC := objectTypeCode16ByInd
	objectTypeCode16ByInd = make([]uint16, 64)
	for i := 1; i < core.Types.Count(); i++ {
		objectTypeCode16ByInd[i] = uint16(1000 + i)
	}
	objectTypeCode16ByInd[63] = 900 // Deliberately stale type entry for factory rejection.
	original := cryptfile.Global()
	cryptfile.SetGlobal(nil)
	t.Cleanup(func() {
		cryptfile.Close()
		cryptfile.SetGlobal(original)
		objectTypeCode16ByInd = oldTOC
		legacy.GetServer = oldGet
		noxServer = old
		core.ExtServer = nil
		owners.Close()
	})
	return core
}

// Common records use a simple real factory object with no UseData allocation.
// FreeObject only clears IDPtr and Field189; release these owned buffers first.
func newObjectXferSimple(t *testing.T, s *server.Server) *server.Object {
	t.Helper()
	u := s.NewObjectByTypeInd(1)
	if u == nil {
		t.Fatal("object allocation")
	}
	t.Cleanup(func() {
		if u.IDPtr != nil {
			legacy.PortTestObjectXferFreeName(u.IDPtr)
			u.IDPtr = nil
		}
		if u.Field189 != nil {
			alloc.FreePtr(u.Field189)
			u.Field189 = nil
		}
		s.Objs.FreeObject(u)
	})
	return u
}

// The typed fixture does not hand these buffers to the painting owner's tracker:
// production FreeObject owns UseData; this fixture owns its other allocations.
func newObjectXferTyped(t *testing.T, s *server.Server, name string) *server.Object {
	t.Helper()
	u := s.NewObjectByTypeID("Port" + name)
	if u == nil {
		t.Fatal("typed allocation: " + name)
	}
	trackObjectXferTyped(t, s, u)
	return u
}

func trackObjectXferTyped(t *testing.T, s *server.Server, u *server.Object) {
	t.Cleanup(func() {
		if u.IDPtr != nil {
			legacy.PortTestObjectXferFreeName(u.IDPtr)
			u.IDPtr = nil
		}
		for _, pp := range []*unsafe.Pointer{&u.InitData, &u.CollideData, &u.UpdateData, &u.Field189} {
			if *pp != nil {
				alloc.FreePtr(*pp)
				*pp = nil
			}
		}
		s.Objs.FreeObject(u)
	})
}
