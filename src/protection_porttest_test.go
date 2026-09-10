//go:build porttest

package opennox

import (
	"bytes"
	"fmt"
	"math/rand"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
)

// Compute each byte lane independently of the production word-reading loop.
func expectedProtectionChecksum(data []byte) uint32 {
	var lanes [4]byte
	for i, b := range data[:len(data)&^3] {
		lanes[i%4] ^= b
	}
	return uint32(lanes[0]) | uint32(lanes[1])<<8 | uint32(lanes[2])<<16 | uint32(lanes[3])<<24
}

func TestProtectionABI(t *testing.T) {
	rng := rand.New(rand.NewSource(0x56fac0))
	check := func(data []byte) {
		t.Helper()
		before := append([]byte(nil), data...)
		want := expectedProtectionChecksum(data)
		if got := protectBytes(data); got != want {
			t.Fatalf("length %d: Go=%08x expected=%08x", len(data), got, want)
		}
		if got := legacy.PortTestProtectionChecksum(data); got != want {
			t.Fatalf("C ABI length=%d: Go=%08x expected=%08x", len(data), got, want)
		}
		if !bytes.Equal(data, before) {
			t.Fatal("input mutated")
		}
	}
	check(nil)
	// Every short length and alignment, not just aligned whole words.
	for n := 0; n <= 1024; n++ {
		for off := 0; off < 8; off++ {
			buf := make([]byte, n+off)
			rng.Read(buf)
			check(buf[off:])
		}
	}
	for i := 0; i < 1000; i++ {
		buf := make([]byte, rng.Intn(65537))
		rng.Read(buf)
		check(buf)
	}
	// Cross the bridge's chunk boundary with each possible trailing length.
	for _, n := range []int{1<<20 - 1, 1 << 20, 1<<20 + 1, 1<<20 + 3, 1<<20 + 4, 2<<20 + 7} {
		buf := make([]byte, n)
		rng.Read(buf)
		check(buf)
	}
	for _, n := range []uint32{0, 1, 3, 4, 1024, 0x7fffffff, 0xffffffff} {
		if got := legacy.PortTestProtectionNull(n); got != 0 {
			t.Fatalf("null with length %d: %08x", n, got)
		}
	}
}

func FuzzProtectionABI(f *testing.F) {
	f.Add([]byte{})
	f.Add([]byte{1, 2, 3, 128, 9, 8, 7})
	f.Fuzz(func(t *testing.T, data []byte) {
		want := expectedProtectionChecksum(data)
		if got := protectBytes(data); got != want {
			t.Fatalf("Go=%08x expected=%08x", got, want)
		}
		if got := legacy.PortTestProtectionChecksum(data); got != want {
			t.Fatalf("ABI Go=%08x expected=%08x", got, want)
		}
	})
}

var protectionBenchmarkResult uint32

func BenchmarkProtectionChecksum(b *testing.B) {
	for _, n := range []int{64, 256, 4096} {
		data := make([]byte, n)
		for i := range data {
			data[i] = byte(i)
		}
		for _, impl := range []struct {
			name string
			fn   func([]byte) uint32
		}{
			{"Go", protectBytes},
			{"CToGo", legacy.PortTestProtectionChecksum},
		} {
			b.Run(fmt.Sprintf("%d/%s", n, impl.name), func(b *testing.B) {
				b.SetBytes(int64(n))
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					protectionBenchmarkResult = impl.fn(data)
				}
			})
		}
	}
}
