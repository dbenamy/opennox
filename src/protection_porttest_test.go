//go:build porttest

package opennox

import (
	"bytes"
	"fmt"
	"math/rand"
	"testing"

	"github.com/opennox/opennox/v1/internal/protectionref"
	"github.com/opennox/opennox/v1/legacy"
)

func TestProtectionCReference(t *testing.T) {
	rng := rand.New(rand.NewSource(0x56fac0))
	check := func(data []byte) {
		t.Helper()
		before := append([]byte(nil), data...)
		want := protectionref.Checksum(data)
		if got := protectBytes(data); got != want {
			t.Fatalf("length %d: Go=%08x C=%08x", len(data), got, want)
		}
		if got := protectionref.Nullable(data); got != want {
			t.Fatalf("nullable reference=%08x, direct=%08x", got, want)
		}
		for _, nullable := range []bool{false, true} {
			if got := legacy.PortTestProtectionChecksum(data, nullable); got != want {
				t.Fatalf("C ABI nullable=%v length=%d: Go=%08x C=%08x", nullable, len(data), got, want)
			}
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
		if got := protectionref.NullWithLength(n) | legacy.PortTestProtectionNull(n); got != 0 {
			t.Fatalf("null with length %d: %08x", n, got)
		}
	}
}

func FuzzProtectionCReference(f *testing.F) {
	f.Add([]byte{})
	f.Add([]byte{1, 2, 3, 128, 9, 8, 7})
	f.Fuzz(func(t *testing.T, data []byte) {
		want := protectionref.Checksum(data)
		if got := protectBytes(data); got != want {
			t.Fatalf("Go=%08x C=%08x", got, want)
		}
		for _, nullable := range []bool{false, true} {
			if got := legacy.PortTestProtectionChecksum(data, nullable); got != want {
				t.Fatalf("ABI nullable=%v Go=%08x C=%08x", nullable, got, want)
			}
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
			{"CReference", protectionref.Checksum},
			{"CToGo", func(p []byte) uint32 { return legacy.PortTestProtectionChecksum(p, false) }},
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
