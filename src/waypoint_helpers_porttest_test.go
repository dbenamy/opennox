//go:build porttest

package opennox

import (
	"testing"

	"github.com/opennox/opennox/v1/legacy"
)

func TestWaypointMasksABI(t *testing.T) {
	flags := []uint32{0, 1, 2, 0x01000000, 0x01000001, 0xffffffff, 0xfffffffe}
	var specs []legacy.PortTestWaypointMask
	for f := 0; f < 256; f++ {
		for m := 0; m < 256; m++ {
			for _, flag := range flags {
				specs = append(specs, legacy.PortTestWaypointMask{Flags: flag, Flags2: byte(f), Mask: byte(m)})
			}
		}
	}
	for m := 0; m < 256; m++ {
		specs = append(specs, legacy.PortTestWaypointMask{Nil: true, Flags: 0xffffffff, Flags2: 255, Mask: byte(m)})
	}
	got := legacy.PortTestWaypointMasks(specs)
	if len(got) != len(specs) {
		t.Fatal("missing mask snapshots")
	}
	for i, s := range specs {
		mask := 0
		for bit := 0; bit < 8; bit++ {
			if (int(s.Flags2)/(1<<bit))%2 == 1 && (int(s.Mask)/(1<<bit))%2 == 1 {
				mask = 1
			}
		}
		enabled := 0
		if !s.Nil && s.Flags%2 == 1 {
			enabled = mask
		}
		if s.Nil {
			mask = -1
		}
		if got[i].Mask != mask || got[i].EnabledMask != enabled || !got[i].Unchanged {
			t.Fatalf("case=%d spec=%+v got=%+v want mask=%d enabled=%d", i, s, got[i], mask, enabled)
		}
	}
}

func TestWaypointLinksABI(t *testing.T) {
	var specs []legacy.PortTestWaypointLink
	for source := -1; source < 3; source++ {
		for next := -1; next < 3; next++ {
			specs = append(specs, legacy.PortTestWaypointLink{Source: source, Next: next})
		}
	}
	got := legacy.PortTestWaypointLinks(specs)
	if len(got) != len(specs) {
		t.Fatal("missing link snapshots")
	}
	for i, s := range specs {
		want := s.Next
		if s.Source < 0 {
			want = -1
		}
		if got[i].First != want || got[i].Second != want || !got[i].Unchanged {
			t.Fatalf("case=%+v got=%+v want=%d", s, got[i], want)
		}
	}
}

func TestWaypointAllocationABI(t *testing.T) {
	got := legacy.PortTestWaypointAllocations(256)
	if len(got) != 256 {
		t.Fatal("missing allocation snapshots")
	}
	for i, raw := range got {
		if len(raw) != 516 {
			t.Fatal("allocation failed or size changed", i, len(raw))
		}
		for offset, value := range raw {
			want := byte(0)
			if offset == 483 {
				want = 1
			}
			if value != want {
				t.Fatalf("allocation=%d offset=%d got=%x want=%x", i, offset, value, want)
			}
		}
	}
}
