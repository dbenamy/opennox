//go:build porttest

package opennox

import (
	"bytes"
	"math/rand"
	"testing"

	"github.com/opennox/opennox/v1/internal/protectionref"
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
	for _, n := range []uint32{0, 1, 3, 4, 1024, 0x7fffffff, 0xffffffff} {
		if got := protectionref.NullWithLength(n); got != 0 {
			t.Fatalf("null with length %d: %08x", n, got)
		}
	}
}
