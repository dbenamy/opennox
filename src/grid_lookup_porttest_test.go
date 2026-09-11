//go:build porttest

package opennox

import (
	"math"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
)

// floatIntExpected is supplied by the converter ABI baseline tests. It accepts a
// float32 raw word and returns the observed x87 FISTP int32 word.
//
// Grid arithmetic deliberately does not call C helpers: C stores each expression
// below to float32 before passing it to the converter.
func gridLookupIndex(v float32) int32 {
	scaled := float32((float64(v) + 11.5) * 0.021739131)
	return floatIntExpected(math.Float32bits(scaled))
}

func gridLookupLocal(v float32) int32 {
	whole := float32(float64(v) + 11.5)
	return floatIntExpected(math.Float32bits(whole)) % 46
}

func gridFallback1(i, j int32) int32 {
	return int32(uint32(0x1e000000) | (uint32(i) << 8) | uint32(j))
}

func gridFallback6(i, j int32) int32 {
	return int32(uint32(0xbe000000) | (uint32(i) << 8) | uint32(j))
}

func gridListValue(base uint32, n int32) int32 { return int32(base + uint32(n)) }

// gridLookupWant independently models 411160's float32 spill points, four-way
// tile selection, and the fixture's category nodes. It never invokes a C helper.
func gridLookupWant(xb, yb uint32, listCount int) int32 {
	x, y := math.Float32frombits(xb), math.Float32frombits(yb)
	i, j := gridLookupIndex(x), gridLookupIndex(y)
	u, v := gridLookupLocal(x), gridLookupLocal(y)
	if i <= 1 || i >= 127 || j <= 1 || j >= 127 {
		return -1
	}

	var fallback int32
	var base uint32
	var px, py int32
	north := u+v >= 46
	if u <= v {
		if north {
			fallback, base, px, py = gridFallback6(i, j), 0xca000000, u, v-23
		} else {
			fallback, base, px, py = gridFallback1(i-1, j), 0x2a000000, u+23, v
		}
	} else if north {
		fallback, base, px, py = gridFallback1(i, j), 0x2a000000, u-23, v
	} else {
		fallback, base, px, py = gridFallback6(i, j-1), 0xca000000, u, v+23
	}
	out := fallback
	for cat := int32(0); cat < int32(listCount); cat++ {
		if subtilePredicate(px, py, cat) != 0 {
			out = gridListValue(base, cat) // fixture nodes are ordered 0..listCount-1: last match wins.
		}
	}
	return out
}

func checkGridLookup(t *testing.T, label string, inputs [][2]uint32, listCount int, wrapper bool) {
	t.Helper()
	got := legacy.PortTestGridLookup(inputs, listCount, wrapper)
	if len(got.Results) != len(inputs) || !got.GridUnchanged || !got.TableUnchanged || !got.NodesUnchanged || !got.InputGuardsOK || !got.PointerUnchanged || !got.Restored {
		t.Fatalf("%s fixture state: results=%d/%d grid=%t table=%t nodes=%t guards=%t pointer=%t restored=%t", label, len(got.Results), len(inputs), got.GridUnchanged, got.TableUnchanged, got.NodesUnchanged, got.InputGuardsOK, got.PointerUnchanged, got.Restored)
	}
	for n, in := range inputs {
		want := gridLookupWant(in[0], in[1], listCount)
		if got.Results[n] != want {
			t.Fatalf("%s case=%d x=%08x y=%08x list=%d wrapper=%t got=%08x want=%08x", label, n, in[0], in[1], listCount, wrapper, uint32(got.Results[n]), uint32(want))
		}
	}
}

func gridBits(v float32) uint32 { return math.Float32bits(v) }

func gridCoordinate(cell, local int32) uint32 {
	return gridBits(float32(cell*46+local) - 11.5)
}

func gridThresholdBits(cell int32) []uint32 {
	// Around x where the scaled expression would equal cell. Feed adjacent raw
	// float32 coordinates, not an already rounded expected index.
	center := float32(float64(cell)/0.021739131 - 11.5)
	out := make([]uint32, 0, 17)
	for v := math.Nextafter32(math.Nextafter32(math.Nextafter32(math.Nextafter32(center, float32(math.Inf(-1))), float32(math.Inf(-1))), float32(math.Inf(-1))), float32(math.Inf(-1))); len(out) < 17; v = math.Nextafter32(v, float32(math.Inf(1))) {
		out = append(out, gridBits(v))
	}
	return out
}

func TestGridLookupCABI(t *testing.T) {
	// All 128 physical grid rows/columns, using local points covering each
	// diagonal branch and equality boundary. The oracle derives the actual i/j;
	// it does not assume decimal 1/46 produces the nominal cell.
	localBoundaries := [][2]int32{{0, 0}, {0, 45}, {45, 0}, {45, 45}, {22, 22}, {22, 23}, {23, 22}, {23, 23}, {1, 44}, {44, 1}, {30, 20}, {20, 30}}
	physical := make([][2]uint32, 0, 128*128*len(localBoundaries))
	for i := int32(0); i < 128; i++ {
		for j := int32(0); j < 128; j++ {
			for _, p := range localBoundaries {
				physical = append(physical, [2]uint32{gridCoordinate(i, p[0]), gridCoordinate(j, p[1])})
			}
		}
	}
	checkGridLookup(t, "physical", physical, 0, false)

	// Exhaust the 46x46 local point space for every prefix of the ordered list;
	// this catches both category matching and last-matching traversal.
	local := make([][2]uint32, 0, 46*46)
	for u := int32(0); u < 46; u++ {
		for v := int32(0); v < 46; v++ {
			local = append(local, [2]uint32{gridCoordinate(64, u), gridCoordinate(64, v)})
		}
	}
	for count := 0; count <= 12; count++ {
		checkGridLookup(t, "locals", local, count, false)
	}

	// Float32 ULPs at the index thresholds exercise the C double-PC53 expression
	// followed by the observed float32 spill. Both dimensions vary independently.
	thresholds := make([]uint32, 0, 128)
	for _, cell := range []int32{-1, 0, 1, 2, 3, 63, 64, 125, 126, 127, 128} {
		thresholds = append(thresholds, gridThresholdBits(cell)...)
	}
	ulp := make([][2]uint32, 0, len(thresholds)*len(thresholds))
	for _, x := range thresholds {
		for _, y := range thresholds {
			ulp = append(ulp, [2]uint32{x, y})
		}
	}
	checkGridLookup(t, "ulp", ulp, 0, false)

	// Raw IEEE cases include zeros, subnormals, extremes, infinities, and NaNs.
	// floatIntExpected supplies the independently baseline-pinned FISTP result.
	raw := []uint32{0, 0x80000000, 1, 0x80000001, 0x007fffff, 0x807fffff, 0x00800000, 0x80800000, 0x3f800000, 0xbf800000, 0x43000000, 0xc47a0000, 0x4f000000, 0xcf000000, 0x7f7fffff, 0xff7fffff, 0x7f800000, 0xff800000, 0x7fc12345, 0xffc54321, 0x7f800001, 0xff800001}
	rawPairs := make([][2]uint32, 0, len(raw)*len(raw))
	for _, x := range raw {
		for _, y := range raw {
			rawPairs = append(rawPairs, [2]uint32{x, y})
		}
	}
	checkGridLookup(t, "raw", rawPairs, 12, false)

	// Target ordinary map coordinates densely, then vary complete raw words.
	var state uint32 = 0x6d2b79f5
	random := make([][2]uint32, 60_000)
	for n := range random {
		state = state*1664525 + 1013904223
		x := state
		state = state*1664525 + 1013904223
		y := state
		if n < 30_000 {
			random[n] = [2]uint32{gridBits(float32(x%6_400_001)/1024 - 100), gridBits(float32(y%6_400_001)/1024 - 100)}
		} else {
			random[n] = [2]uint32{x, y}
		}
	}
	checkGridLookup(t, "random", random, 0, false)
	// The public Go wrapper allocates/copies its input. A representative subset
	// exercises each branch, boundary, and invalid result through that route.
	wrapper := append([][2]uint32(nil), local...)
	wrapper = append(wrapper, ulp[:min(len(ulp), 512)]...)
	wrapper = append(wrapper, rawPairs...)
	checkGridLookup(t, "wrapper-fallback", wrapper, 0, true)
	checkGridLookup(t, "wrapper-list", wrapper, 12, true)

	t.Logf("grid cases: physical=%d locals=%d ulp=%d raw=%d random=%d wrapper=%d total=%d", len(physical), len(local)*13, len(ulp), len(rawPairs), len(random), len(wrapper), len(physical)+len(local)*13+len(ulp)+len(rawPairs)+len(random)+2*len(wrapper))
}

func BenchmarkGridLookup(b *testing.B) {
	for _, wrapper := range []bool{false, true} {
		name := "C-caller"
		if wrapper {
			name = "Go-wrapper"
		}
		b.Run(name, func(b *testing.B) {
			got := legacy.PortTestGridLookupBenchmark(b.N, wrapper)
			want := uint32(0x1e000204) * uint32(b.N)
			if got.Checksum != want || !got.GridUnchanged || !got.TableUnchanged || !got.NodesUnchanged || !got.InputGuardsOK || !got.PointerUnchanged || !got.Restored {
				b.Fatalf("checksum=%x want=%x or fixture changed", got.Checksum, want)
			}
		})
	}
}
