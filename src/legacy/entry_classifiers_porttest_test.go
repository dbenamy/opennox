//go:build porttest

package legacy

import (
	"crypto/sha256"
	"fmt"
	"os"
	"testing"
)

// The capture consists of digit then alphanumeric bitsets, least-significant
// bit first, over all 65,536 possible legacy text units. Only truth is observable.
func TestEntryClassifiers(t *testing.T) {
	const want = "f5c39db5e866885891bbb1fe8d0b2244d9fbcea7d232cdd4b7e6f90dd03ff6e0"
	bits := make([]byte, 2*8192)
	for i := 0; i < 65536; i++ {
		v := uint16(i)
		if uiEntryDigit(v) {
			bits[i/8] |= 1 << uint(i%8)
		}
		if uiEntryAlnum(v) {
			bits[8192+i/8] |= 1 << uint(i%8)
		}
	}
	for _, c := range []uint16{'0', '5', '9'} {
		if !uiEntryDigit(c) || !uiEntryAlnum(c) {
			t.Fatalf("ASCII digit %q rejected", c)
		}
	}
	for _, c := range []uint16{'a', 'Z'} {
		if uiEntryDigit(c) || !uiEntryAlnum(c) {
			t.Fatalf("ASCII letter %q classification", c)
		}
	}
	for _, c := range []uint16{0, ' ', '!', '_'} {
		if uiEntryDigit(c) || uiEntryAlnum(c) {
			t.Fatalf("ASCII non-alphanumeric %q accepted", c)
		}
	}
	got := fmt.Sprintf("%x", sha256.Sum256(bits))
	if path := os.Getenv("OPENNOX_ENTRY_CLASSIFIERS_CAPTURE"); path != "" {
		if err := os.WriteFile(path, bits, 0600); err != nil {
			t.Fatal(err)
		}
	}
	t.Logf("131072 boolean classifications: %s", got)
	if got != want {
		t.Fatalf("classifier capture %s want %s", got, want)
	}
}
