//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestClientScreenEffectsDistancePrimitives(t *testing.T) {
	env := legacy.PortTestNewEffectsEnvironment()
	defer env.Restore()
	type result struct {
		Op     int
		Input  [4]int32
		Return uint32
	}
	var out []result
	capture := func(op int, a [4]int32) { out = append(out, result{op, a, legacy.PortTestScreenPrimitive(op, a)}) }
	for v := uint32(0); v < 65536; v++ {
		capture(0, [4]int32{int32(v)})
	}
	for shift := uint(16); shift < 32; shift++ {
		for d := int32(-3); d <= 3; d++ {
			capture(0, [4]int32{int32(uint32(1<<shift) + uint32(d))})
		}
	}
	// Cover both sides of every bucket at all larger table scales.
	for shift := uint(10); shift <= 24; shift += 2 {
		for bucket := uint32(64); bucket < 256; bucket++ {
			for _, offset := range []uint32{0xffffffff, 0, (1 << shift) - 1, 1 << shift} {
				capture(0, [4]int32{int32((bucket << shift) + offset)})
			}
		}
	}
	var random uint32 = 0x81234abc
	for i := 0; i < 65536; i++ {
		random = random*1664525 + 1013904223
		capture(0, [4]int32{int32(random)})
	}
	for _, v := range []uint32{0xfffffff0, 0xfffffffe, 0xffffffff} {
		capture(0, [4]int32{int32(v)})
	}
	for _, x := range []int32{-2147483648, -65536, -5888, -32768, -257, -256, -255, -17, -16, -15, -1, 0, 1, 15, 16, 17, 255, 256, 257, 5888, 32767, 65535, 2147483647} {
		for _, y := range []int32{-65536, -5888, -257, -1, 0, 1, 257, 5888, 65535} {
			capture(1, [4]int32{x, y})
			for _, origin := range []int32{0, 37, -32768, 2147483647} {
				capture(2, [4]int32{origin, origin, origin + x, origin + y})
			}
		}
	}
	// The table's exact power-of-four anchors define scale across every branch.
	if legacy.PortTestScreenPrimitive(0, [4]int32{}) != 0 {
		t.Fatal("zero distance")
	}
	for shift := uint(0); shift < 32; shift += 2 {
		if legacy.PortTestScreenPrimitive(0, [4]int32{int32(uint32(1) << shift)}) != uint32(1)<<(shift/2) {
			t.Fatalf("sqrt scale anchor %d", shift)
		}
	}
	effectsCapture(t, "screen-effects-distance", out, len(out), "1e2c756f945ea9f98383bf55d6b20e6e4c9350847694e786445d6ac4d116fe83")
}
