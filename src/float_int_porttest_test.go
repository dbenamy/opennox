//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"math"
	"testing"
)

// Decode IEEE754 integer magnitude with integer arithmetic, including the
// x87 invalid-conversion result to be checked against the original C baseline.
func floatIntExpected(bits uint32) int32 {
	exp := int((bits>>23)&255) - 127
	if exp < 0 {
		return 0
	}
	if exp >= 31 {
		return -2147483648
	}
	mantissa := uint32(0x800000) | (bits & 0x7fffff)
	var magnitude uint32
	if exp >= 23 {
		magnitude = mantissa << uint(exp-23)
	} else {
		magnitude = mantissa >> uint(23-exp)
	}
	if bits>>31 != 0 {
		return -int32(magnitude)
	}
	return int32(magnitude)
}
func checkFloatInt(t *testing.T, bits []uint32, control int) {
	t.Helper()
	got := legacy.PortTestFloatInt(bits, control)
	for i, b := range bits {
		full, abs := floatIntExpected(b), floatIntExpected(b&0x7fffffff)
		want := [3]int32{full, int32(int16(full)), int32(int16(abs))}
		if got[i].Values != want || !got[i].GuardsOK || !got[i].ControlOK {
			t.Fatalf("bits=%08x control=%x got=%+v want=%v", b, control, got[i], want)
		}
	}
}
func TestFloatIntBaseline(t *testing.T) {
	var bits []uint32
	// Every exponent/sign, with mantissas at representative rounding boundaries.
	for sign := uint32(0); sign < 2; sign++ {
		for exp := uint32(0); exp < 256; exp++ {
			for _, mantissa := range []uint32{0, 1, 2, 0x3ffffe, 0x3fffff, 0x400000, 0x400001, 0x7ffffe, 0x7fffff} {
				bits = append(bits, sign<<31|exp<<23|mantissa)
			}
		}
	}
	// Every signed16 integer boundary and its adjacent representable floats.
	for n := -65536; n <= 65536; n++ {
		f := float32(n)
		bits = append(bits, math.Float32bits(f), math.Float32bits(math.Nextafter32(f, float32(math.Inf(-1)))), math.Float32bits(math.Nextafter32(f, float32(math.Inf(1)))))
	}
	// Deterministic raw patterns include arbitrary NaNs/subnormals/out-of-range.
	state := uint32(0x6379a52b)
	for i := 0; i < 1000000; i++ {
		state ^= state << 13
		state ^= state >> 17
		state ^= state << 5
		bits = append(bits, state)
	}
	checkFloatInt(t, bits, -1)
	t.Logf("%d float inputs, %d ABI conversions", len(bits), 3*len(bits))
}
func TestFloatIntControlWord(t *testing.T) {
	values := []uint32{0, 0x80000000, 1, 0x80000001, 0x3fc00000, 0xbfc00000, 0x46ffffff, 0x47000000, 0x47000001, 0xc7000001, 0x4effffff, 0x4f000000, 0xcf000000, 0xcf000001, 0x7f800000, 0xff800000, 0x7f800001, 0xffc12345, 0x7fffffff}
	for _, pc := range []int{0, 2, 3} {
		for rc := 0; rc < 4; rc++ {
			checkFloatInt(t, values, pc<<8|rc<<10)
		}
	}
}

func BenchmarkFloatIntCABI(b *testing.B) {
	got := legacy.PortTestFloatIntBenchmark(b.N)
	n := uint64(b.N)
	tail := n % 1024
	want := uint32((n/1024)*523776 + tail*(tail-1)/2)
	if got != want {
		b.Fatalf("checksum got=%x want=%x", got, want)
	}
}
