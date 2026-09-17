//go:build porttest

package opennox

import (
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/memmap/nox/blobdata"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

type collisionCoreOwner struct {
	absScratch *uint32
	*worldCollisionOwner
	words                              map[string]*uint32
	callback                           unsafe.Pointer
	resetHits, resetCalls, resetQueues func()
	hits                               func(map[unsafe.Pointer]uint32) [][5]uint32
	calls                              func(map[unsafe.Pointer]uint32) [][4]uint32
	queues                             func(map[unsafe.Pointer]uint32) [3]uint32
}

func newCollisionCoreOwner(t *testing.T) *collisionCoreOwner {
	t.Helper()
	o := &collisionCoreOwner{worldCollisionOwner: newWorldCollisionOwner(t)}
	t.Cleanup(noxflags.PortTestGameFlags(0))
	worldGeometryTables(t)
	// The retained absolute-value helper uses an originally relocated scratch pointer.
	o.absScratch = memmap.PtrUint32(0x5D4594, 527672)
	absSlot := memmap.PtrPtr(0x587000, 55744)
	oldAbs, oldSlot := *o.absScratch, *absSlot
	*o.absScratch, *absSlot = 0, unsafe.Pointer(o.absScratch)
	t.Cleanup(func() { *o.absScratch, *absSlot = oldAbs, oldSlot })
	normal := unsafe.Slice((*byte)(memmap.PtrOff(0x587000, 289928)), 32)
	old := append([]byte(nil), normal...)
	copy(normal, blobdata.PortTestCollisionCoreNormals())
	t.Cleanup(func() { copy(normal, old) })
	nonzero := false
	for _, v := range normal {
		nonzero = nonzero || v != 0
	}
	if !nonzero {
		t.Fatal("missing circle/box normals")
	}
	t.Cleanup(o.s.PortTestRewardTypes([]string{"Trigger", "BlackPowder", "TelekinesisHand", "SmallFist", "MediumFist", "LargeFist", "Meteor", "Spike", "PeriodicSpike"}, nil, true, 0, 0))
	var restore func()
	o.words, restore = legacy.PortTestCollisionCoreGlobals()
	t.Cleanup(restore)
	o.resetHits, o.hits, restore = legacy.PortTestWorldGeometryHitOwner()
	t.Cleanup(restore)
	o.callback, o.resetCalls, o.calls, restore = legacy.PortTestCollisionCoreObserver()
	t.Cleanup(restore)
	o.resetQueues, o.queues, restore = legacy.PortTestGeometryQueues()
	t.Cleanup(restore)
	return o
}
func collisionCoreID(t *testing.T, ids map[unsafe.Pointer]uint32, raw uint32) uint32 {
	t.Helper()
	if raw == 0 || raw == 0xffffffff {
		return raw
	}
	id, ok := ids[unsafe.Pointer(uintptr(raw))]
	if !ok {
		t.Fatalf("pointer outside collision owner: %x", raw)
	}
	return id
}
func collisionCoreGuarded(t *testing.T, o *collisionCoreOwner, n int) unsafe.Pointer {
	t.Helper()
	raw := unsafe.Slice((*byte)(o.record(t, n+16)), n+16)
	for i := 0; i < 8; i++ {
		raw[i] = 0xa5
		raw[n+8+i] = 0x5a
	}
	t.Cleanup(func() {
		for i := 0; i < 8; i++ {
			if raw[i] != 0xa5 || raw[n+8+i] != 0x5a {
				t.Error("collision fixture guard changed")
			}
		}
	})
	return unsafe.Pointer(&raw[8])
}
func collisionCoreIDs(objects ...*server.Object) map[unsafe.Pointer]uint32 {
	ids := map[unsafe.Pointer]uint32{}
	for i, u := range objects {
		ids[u.CObj()] = uint32(1001 + i)
	}
	return ids
}
