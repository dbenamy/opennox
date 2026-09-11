//go:build porttest

package opennox

import (
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

func waypointAppendKind(spec legacy.PortTestWaypointAppendSpec) byte {
	if spec.Mode == 1 {
		return spec.BlobKind
	}
	return spec.Kind
}

// initialWaypointEdge mirrors fixture initialization. In particular, every
// unspecified slot points to alternate target 2, so target-2 tests prove that
// the oracle scans all eligible slots instead of just ExistingSlot.
func initialWaypointEdge(spec legacy.PortTestWaypointAppendSpec, slot int) (int8, byte) {
	if slot == int(spec.ExistingSlot) && slot < 32 {
		return spec.ExistingTarget, spec.ExistingKind
	}
	return 2, byte(0x40 + slot)
}

// appendModel deliberately represents the x86 comparison, not Go byte
// equality: stored Ind is movzbl, while incoming C char is movsbl.
func appendModel(spec legacy.PortTestWaypointAppendSpec) (ret int, count byte, target int8, kind byte) {
	count = spec.Count
	kind = waypointAppendKind(spec)
	target = spec.Target
	if count >= 31 || target == 0 {
		return 0, count, target, kind
	}
	incoming := int(int8(kind))
	for slot := 0; slot < int(count); slot++ {
		entryTarget, entryKind := initialWaypointEdge(spec, slot)
		if entryTarget == target && int(entryKind) == incoming {
			return 0, count, target, kind
		}
	}
	return 1, count + 1, target, kind
}

func TestWaypointAppendABI(t *testing.T) {
	if unsafe.Sizeof(server.Waypoint{}) != 516 || unsafe.Offsetof(server.Waypoint{}.Points) != 92 || unsafe.Offsetof(server.Waypoint{}.PointsCnt) != 476 || unsafe.Sizeof(server.WaypointSub{}) != 8 || unsafe.Offsetof(server.WaypointSub{}.Waypoint) != 0 || unsafe.Offsetof(server.WaypointSub{}.Ind) != 4 {
		t.Fatalf("unexpected waypoint layout: wp size=%d points=%d count=%d sub size=%d ptr=%d ind=%d", unsafe.Sizeof(server.Waypoint{}), unsafe.Offsetof(server.Waypoint{}.Points), unsafe.Offsetof(server.Waypoint{}.PointsCnt), unsafe.Sizeof(server.WaypointSub{}), unsafe.Offsetof(server.WaypointSub{}.Waypoint), unsafe.Offsetof(server.WaypointSub{}.Ind))
	}
	var specs []legacy.PortTestWaypointAppendSpec
	addBoth := func(base legacy.PortTestWaypointAppendSpec) {
		k := base.Kind
		base.Mode, base.BlobKind = 0, k+1 // direct must ignore distinct blob state
		specs = append(specs, base)
		base.Mode, base.Kind, base.BlobKind = 1, k+1, k // wrapper must ignore distinct explicit kind
		specs = append(specs, base)
	}
	// Exact apparent duplicates cover every incoming byte and every count.
	// ExistingSlot cycles through every physical slot, including 29 (last C
	// scan slot), 30 (last insertion slot), and 31/outside-count slots.
	for kind := 0; kind < 256; kind++ {
		for count := 0; count < 256; count++ {
			k := byte(kind)
			addBoth(legacy.PortTestWaypointAppendSpec{Count: byte(count), Kind: k, ExistingSlot: byte(kind+count) & 31, ExistingTarget: 1, ExistingKind: k, Target: 1})
		}
	}
	for kind := 0; kind < 256; kind++ {
		k := byte(kind)
		for _, count := range []byte{0, 1, 2, 15, 30, 31, 32, 255} {
			for slot := byte(0); slot < 32; slot++ {
				// Test an exact duplicate at every position for each kind.
				addBoth(legacy.PortTestWaypointAppendSpec{Count: count, Kind: k, ExistingSlot: slot, ExistingTarget: 1, ExistingKind: k, Target: 1})
				// Same kind but alternate target proves duplicate detection is
				// pointer-based even though all waypoint Index values are equal.
				addBoth(legacy.PortTestWaypointAppendSpec{Count: count, Kind: k, ExistingSlot: slot, ExistingTarget: 2, ExistingKind: k, Target: 1})
				// Target 2 has every unspecified slot initialized to it, making
				// the independent oracle account for every scanned slot.
				addBoth(legacy.PortTestWaypointAppendSpec{Count: count, Kind: k, ExistingSlot: slot, ExistingTarget: 1, ExistingKind: k ^ 1, Target: 2})
			}
			// Nil and self cover exact high-byte behavior and the self guard.
			addBoth(legacy.PortTestWaypointAppendSpec{Count: count, Kind: k, ExistingSlot: count & 31, ExistingTarget: -1, ExistingKind: k, Target: -1})
			addBoth(legacy.PortTestWaypointAppendSpec{Count: count, Kind: k, ExistingSlot: count & 31, ExistingTarget: 1, ExistingKind: k, Target: 0})
		}
	}
	snap := legacy.PortTestWaypointAppend(specs)
	if snap.BlobAfterRestore != snap.BlobBefore {
		t.Fatalf("blob kind not restored: before=%x after=%x", snap.BlobBefore, snap.BlobAfterRestore)
	}
	got := snap.Results
	if len(got) != len(specs) {
		t.Fatalf("result count got %d want %d", len(got), len(specs))
	}
	for i, spec := range specs {
		wantRet, wantCount, target, kind := appendModel(spec)
		res := got[i]
		if res.Return != wantRet || res.PointsCnt != wantCount || res.BlobAfterCall != spec.BlobKind || !res.SourceOnlyExpectedWrites || !res.OutsideSourceUnchanged || !res.GuardWordsUnchanged {
			t.Fatalf("case %d spec=%+v got=%+v want ret/count=%d/%d", i, spec, res, wantRet, wantCount)
		}
		for slot := 0; slot < 32; slot++ {
			if wantRet != 0 && slot == int(spec.Count) {
				if res.Targets[slot] != target || res.Kinds[slot] != kind {
					t.Fatalf("case %d append slot=%d got target/kind=%d/%d want %d/%d", i, slot, res.Targets[slot], res.Kinds[slot], target, kind)
				}
				continue
			}
			wantTarget, wantKind := initialWaypointEdge(spec, slot)
			if res.Targets[slot] != wantTarget || res.Kinds[slot] != wantKind {
				t.Fatalf("case %d changed slot=%d got target/kind=%d/%d want %d/%d", i, slot, res.Targets[slot], res.Kinds[slot], wantTarget, wantKind)
			}
		}
	}
}
