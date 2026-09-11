//go:build porttest

package opennox

import (
	"github.com/opennox/libs/prand"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestEdgeNormalizationThresholds(t *testing.T) {
	checks := 0
	for w := 0; w < 256; w++ {
		var specs []legacy.PortTestEdgeDirectSpec
		for h := 0; h < 256; h++ {
			values := []int32{-2147483648, -1, 0, 1, 2147483647}
			for _, bounds := range [][2]int{{w - 3, w}, {w + 2*h - 6, w + 2*h - 2}, {2*(w+h) - 8, 2*(w+h) - 3}} {
				for v := bounds[0]; v <= bounds[1]; v++ {
					values = append(values, int32(v))
				}
			}
			seen := make(map[int32]bool)
			for _, value := range values {
				if seen[value] {
					continue
				}
				seen[value] = true
				specs = append(specs, legacy.PortTestEdgeDirectSpec{Width: byte(w), Height: byte(h), Index: int32((w + h) % 64), Edge: value, Normalize: true})
			}
		}
		seed := 137 + w
		got := legacy.PortTestEdgeMapping(seed, specs, nil)
		if !got.Restored || !got.TableUnchanged || !got.MappingUnchanged || !got.GuardsUnchanged || len(got.Direct) != len(specs) {
			t.Fatal("fixture state")
		}
		logic, other := prand.New(seed).Index(), prand.New(seed+1).Index()
		for i, s := range specs {
			want := edgeNormalize(s.Width, s.Height, s.Edge)
			r := got.Direct[i]
			if r.Result != want || r.LogicIndex != logic || r.OtherIndex != other {
				t.Fatalf("w/h/edge/index=%d/%d/%d/%d got=%+v expected=%d", s.Width, s.Height, s.Edge, s.Index, r, want)
			}
			checks++
		}
	}
	t.Logf("%d normalization threshold calls", checks)
}
