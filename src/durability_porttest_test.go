//go:build porttest

package opennox

import (
	"math"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
)

const (
	durabilityHalf    = uint64(0x3fe0000000000000) // .5
	durabilityQuarter = uint64(0x3fd0000000000000) // .25
)

// defaultDurability is deliberately integer-only: the shipped threshold values
// are exact binary rationals, so this does not mirror C floating comparisons.
func defaultDurability(cur, max uint16) int {
	if max == 0 {
		return 4
	}
	if cur == max {
		return 0
	}
	if uint32(cur)*2 >= uint32(max) {
		return 1
	}
	if uint32(cur)*4 < uint32(max) {
		return 3
	}
	return 2
}

// changedDurability models comparison order only for deliberately altered
// globals. NaNs, infinities, and signed zero therefore exercise C's ordered
// floating comparisons without using an original-C result as an oracle.
func changedDurability(cur, max uint16, halfBits, quarterBits uint64) int {
	if max == 0 {
		return 4
	}
	if cur == max {
		return 0
	}
	if float64(cur) >= float64(max)*math.Float64frombits(halfBits) {
		return 1
	}
	if float64(cur) < float64(max)*math.Float64frombits(quarterBits) {
		return 3
	}
	return 2
}

func durabilityCases() []legacy.PortTestDurabilityCase {
	// Every uint16 maximum gets values surrounding both rational thresholds,
	// equality, and the uint16 endpoints. Duplicates intentionally cost no
	// semantic ambiguity and keep construction simple.
	out := make([]legacy.PortTestDurabilityCase, 0, 12*65536)
	for m := uint32(0); m <= math.MaxUint16; m++ {
		q, h := m/4, m/2
		vals := []uint32{0, 1, math.MaxUint16, m - 1, m, m + 1, q - 1, q, q + 1, h - 1, h, h + 1}
		for _, v := range vals {
			if v <= math.MaxUint16 {
				out = append(out, legacy.PortTestDurabilityCase{Current: uint16(v), Maximum: uint16(m)})
			}
		}
	}
	return out
}

func checkDurability(t *testing.T, name string, cases []legacy.PortTestDurabilityCase, half, quarter uint64, want func(uint16, uint16) int) {
	t.Helper()
	got := legacy.PortTestDurability(cases, half, quarter)
	if got.AfterCallHalf != half || got.AfterCallQuarter != quarter {
		t.Fatalf("%s: C mutated thresholds: after call %016x/%016x want %016x/%016x", name, got.AfterCallHalf, got.AfterCallQuarter, half, quarter)
	}
	if got.AfterRestoreHalf != got.BeforeHalf || got.AfterRestoreQuarter != got.BeforeQuarter {
		t.Fatalf("%s: thresholds not restored: before %016x/%016x after %016x/%016x", name, got.BeforeHalf, got.BeforeQuarter, got.AfterRestoreHalf, got.AfterRestoreQuarter)
	}
	if len(got.Results) != len(cases) {
		t.Fatalf("%s: result length got %d want %d", name, len(got.Results), len(cases))
	}
	for i, c := range cases {
		if w := want(c.Current, c.Maximum); got.Results[i] != w {
			t.Fatalf("%s case %d current=%d max=%d half=%016x quarter=%016x: got %d want %d", name, i, c.Current, c.Maximum, half, quarter, got.Results[i], w)
		}
	}
}

func TestDurabilityABI(t *testing.T) {
	cases := durabilityCases()
	checkDurability(t, "shipped", cases, durabilityHalf, durabilityQuarter, defaultDurability)

	// Include reversed bands, signed zeros, both infinity signs, quiet and
	// signaling NaNs, and ordinary non-rational values. Equality/max==0 must
	// continue to short-circuit before either comparison under all of them.
	changed := [][2]uint64{
		{0x3fd0000000000000, 0x3fe0000000000000}, // reversed .25/.5
		{0, 0}, {0x8000000000000000, 0x8000000000000000},
		{0x8000000000000000, 0},
		{0x7ff0000000000000, durabilityQuarter}, {0xfff0000000000000, durabilityQuarter},
		{0x7ff8000000001234, durabilityQuarter}, {0x7ff0000000001234, durabilityQuarter},
		{durabilityHalf, 0x7ff8000000005678}, {durabilityHalf, 0x7ff0000000005678},
		{0x3fd5555555555555, 0x3fc5555555555555},
		{durabilityHalf, 0x7ff0000000000000}, {durabilityHalf, 0xfff0000000000000},
		{0x7fefffffffffffff, 1}, {1, 0x7fefffffffffffff},
		{0xbfe0000000000000, 0xbfd0000000000000},
	}
	for _, thresholds := range changed {
		half, quarter := thresholds[0], thresholds[1]
		checkDurability(t, "changed", cases, half, quarter, func(cur, max uint16) int {
			return changedDurability(cur, max, half, quarter)
		})
	}
}
