//go:build porttest

package opennox

import (
	"math"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
)

func TestWorldGridNumericScratch(t *testing.T) {
	inputs := []uint32{0, 0x80000000, 1, 0x80000001, 0x007fffff, 0x00800000, 0x7f7fffff, 0xff7fffff, 0x7f800000, 0xff800000, 0x7fc12345, 0xffc12345, 0x7f812345, 0xff812345}
	for _, f := range []float32{0.5, 1, 1.5, 2.5, 3.5, 127.5, 65535.5, 4194304, 8388607, 8388608, 16777216} {
		b := math.Float32bits(f)
		for _, u := range []uint32{b - 1, b, b + 1} {
			inputs = append(inputs, u, u|0x80000000)
		}
	}
	seed := uint32(0x41ac7209)
	for i := 0; i < 4096; i++ {
		seed = seed*1664525 + 1013904223
		inputs = append(inputs, seed)
	}
	rows := legacy.PortTestWorldNumeric(inputs)
	for _, r := range rows {
		f := math.Float32frombits(r.Input)
		nan := math.IsNaN(float64(f))
		if !r.PointerOK {
			t.Fatalf("scratch pointer input %08x", r.Input)
		}
		if f < 0 {
			if r.Rounded != 0 || r.RoundWords != [3]uint32{0x13572468, 0x87654321, 0xfedcba98} {
				t.Fatalf("negative round changed scratch: %+v", r)
			}
		} else if !nan {
			want := math.Float32bits(float32(float64(f) + 8388608))
			if r.Rounded != want&0x7fffff || r.RoundWords != [3]uint32{1, want, want & 0x7fffff} {
				t.Fatalf("round contract: %+v want %08x", r, want)
			}
		}
		if !nan {
			want := r.Input & 0x7fffffff
			if r.Redirect {
				want = r.Input
			}
			if r.Scratch != want || r.Abs != math.Float64bits(float64(math.Float32frombits(want))) {
				t.Fatalf("absolute scratch contract: %+v", r)
			}
		}
		if r.Redirect && r.Pointed != 0x7f812345 {
			t.Fatalf("redirected mask: %+v", r)
		}
		if !r.Redirect && r.Pointed != r.Scratch {
			t.Fatalf("scratch alias: %+v", r)
		}
	}
	// Frozen after matching original-C captures in separate processes.
	drawableStateCapture(t, "world-numeric", rows, "3dca5eb8c03908873ee1400086c07f389559e3add4b6cb962b9742a7cd47b8ac")
}
