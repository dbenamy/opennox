//go:build porttest

package legacy

import (
	"testing"
	"unsafe"
)

// Force address reuse independently of libc or native allocator size classes.
// A saved slot identifies its bound record, including an interior alias, even
// when a later record occupies the same address.
func TestPaintSlotIdentityOnAddressReuse(t *testing.T) {
	var storage [64]byte
	p := unsafe.Pointer(&storage[0])
	f := &paintTestFixture{mapRoomTestFixture: &mapRoomTestFixture{slots: map[int]*mapRoomTestRegion{}}}
	old := f.register(p, 32, "input", false)
	alias := &mapRoomTestRegion{ptr: unsafe.Add(p, 4), size: 28, id: old.id + 4, kind: "alias", alive: true}
	f.slots[1], f.slots[2], f.slots[3] = old, alias, nil
	old.alive = false
	replacement := f.register(p, 32, "replacement", false)
	f.slots[4] = replacement
	if got := f.norm(mapRoomRaw(p)); got != replacement.id {
		t.Fatalf("forced reuse setup: address identifies %d, want %d", got, replacement.id)
	}
	got := f.snapshotSlots()
	for slot, want := range map[int]uint32{1: old.id, 2: alias.id, 3: 0, 4: replacement.id} {
		if got[slot] != want {
			t.Errorf("slot %d: identity %d, want bound identity %d", slot, got[slot], want)
		}
	}
}
