//go:build porttest && !server

package opennox

import (
	"testing"

	"github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
)

func TestFloorEligibility(t *testing.T) {
	coords := []int32{-2147483648, -100, 0, 68, 69, 70, 80, 81, 82, 100, 128, 192, 5784, 5785, 5786, 5830, 5831, 5832, 16777215, 16777216, 2147483647}
	tiles := []int32{-2147483648, -2, -1, 0, 6, 254, 255, 256, 65535, 2147483647}
	fs := []uint32{0, 1, uint32(noxflags.EngineNoFloorRendering), ^uint32(noxflags.EngineNoFloorRendering), ^uint32(0)}
	var inputs []legacy.PortTestFloorInput
	for _, x := range coords {
		for _, y := range coords {
			for _, tile := range tiles {
				for _, f := range fs {
					inputs = append(inputs, legacy.PortTestFloorInput{X: x, Y: y, Tile: tile, Flags: f})
				}
			}
		}
	}
	// Each unrelated flag alone must leave eligibility unchanged.
	for bit := uint(0); bit < 32; bit++ {
		inputs = append(inputs, legacy.PortTestFloorInput{X: 128, Y: 192, Tile: 6, Flags: 1 << bit})
	}
	// Disabled rendering must short-circuit before either viewport or grid access.
	inputs = append(inputs, legacy.PortTestFloorInput{Flags: uint32(noxflags.EngineNoFloorRendering), NilViewport: true})
	got := legacy.PortTestFloorEligibility(inputs)
	if !got.Unchanged || !got.Restored || len(got.Values) != len(inputs) {
		t.Fatalf("fixture state: unchanged=%t restored=%t results=%d/%d", got.Unchanged, got.Restored, len(got.Values), len(inputs))
	}
	for n, in := range inputs {
		i, j := gridLookupIndex(float32(in.X)), gridLookupIndex(float32(in.Y))
		want := 0
		if in.Flags&uint32(noxflags.EngineNoFloorRendering) == 0 && i > 1 && i < 127 && j > 1 && j < 127 && in.Tile != -1 && in.Tile != 255 {
			want = 1
		}
		if got.Values[n] != want {
			t.Fatalf("case %d: %+v got=%d want=%d", n, in, got.Values[n], want)
		}
	}
	t.Logf("floor eligibility cases: %d", len(inputs))
}
