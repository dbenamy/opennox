//go:build porttest

package opennox

import (
	"math"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
)

func TestServerRuntimeDoubleConversion(t *testing.T) {
	var inputs []uint64
	for sign := uint64(0); sign < 2; sign++ {
		for exponent := uint64(0); exponent < 2048; exponent++ {
			for _, mantissa := range []uint64{0, 1, 2, 0x7ffffffffffff, 0x8000000000000, 0xffffffffffffe, 0xfffffffffffff} {
				inputs = append(inputs, sign<<63|exponent<<52|mantissa)
			}
		}
	}
	for _, v := range []float64{-2147483649, -2147483648, -2147483647, -65536, -1, -.5, 0, .5, 1, 65536, 2147483647, 2147483648, 2147483649} {
		inputs = append(inputs, math.Float64bits(v), math.Float64bits(math.Nextafter(v, math.Inf(-1))), math.Float64bits(math.Nextafter(v, math.Inf(1))))
	}
	state := uint64(0x731f6258ba946d0e)
	for i := 0; i < 8192; i++ {
		state ^= state << 13
		state ^= state >> 7
		state ^= state << 17
		inputs = append(inputs, state)
	}
	type row struct {
		Bits    uint64
		Control int
		Value   int32
	}
	var rows []row
	check := func(bits []uint64, control int) {
		got, unchanged := legacy.PortTestDoubleInt(bits, control)
		if !unchanged || len(got) != len(bits) {
			t.Fatal("double conversion changed the control word")
		}
		for i, b := range bits {
			value := math.Float64frombits(b)
			want := int32(-2147483648)
			if !math.IsNaN(value) && value >= -2147483648 && value < 2147483648 {
				want = int32(math.Trunc(value))
			}
			if got[i] != want {
				t.Fatalf("double bits=%016x control=%x: got %d want %d", b, control, got[i], want)
			}
			rows = append(rows, row{b, control, got[i]})
		}
	}
	check(inputs, -1)
	controls := []uint64{0, 1, 0x8000000000000000, 0x8000000000000001, 0x3ff8000000000000, 0xbff8000000000000, 0x41dfffffffffffff, 0x41e0000000000000, 0xc1e0000000000000, 0x7ff0000000000000, 0xfff0000000000000, 0x7ff0000000000001, 0x7ff8123456789012}
	for _, pc := range []int{0, 2, 3} {
		for rc := 0; rc < 4; rc++ {
			check(controls, pc<<8|rc<<10)
		}
	}
	interactionCapture(t, "server-runtime-double-conversion", rows)
}
