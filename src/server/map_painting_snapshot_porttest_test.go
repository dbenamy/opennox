//go:build porttest

package server

import (
	"testing"
	"unsafe"
)

func TestMapPaintingWallDefinitionSnapshotInvalidation(t *testing.T) {
	s := new(Server)
	owner := s.PortTestMapPaintingOwners([2]unsafe.Pointer{})
	defer owner.Close()
	snapshot := func() [32]byte { return owner.WallSnapshot(func(v uint32) uint32 { return v }).Definitions }
	initial := snapshot()
	if snapshot() != initial {
		t.Fatal("unchanged wall definitions changed digest")
	}
	// Include inactive storage, not only the three fixture definitions.
	s.Walls.defs[79].Flags32 ^= 0x80
	changed := snapshot()
	if changed == initial {
		t.Fatal("inactive definition mutation was not captured")
	}
	if snapshot() != changed {
		t.Fatal("unchanged mutated definitions changed digest")
	}
	s.Walls.defs[79].Flags32 ^= 0x80
	if snapshot() != initial {
		t.Fatal("restored wall definitions changed digest")
	}
	owner.Reset(7, 2, 0)
	if snapshot() == initial {
		t.Fatal("reset variation limits were not captured")
	}
	owner.Reset(7, 3, 4)
	if snapshot() != initial {
		t.Fatal("restored variation limits changed digest")
	}
}
