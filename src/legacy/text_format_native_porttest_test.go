//go:build porttest

package legacy

import (
	"slices"
	"testing"
	"unsafe"
)

func TestTextFormatNativeConsumers(t *testing.T) {
	format := []uint16{'%', 's', 0}
	text := []uint16{'A', 0xd800, 'B', 0}
	arg := textFormatPointer(unsafe.Pointer(&text[0]))
	for _, capacity := range []int{0, 1, 2, 3, 4, 5} {
		dst := make([]uint16, capacity+2)
		for i := range dst {
			dst[i] = 0xa55a
		}
		n := textFormatBuffer(dst[1:1+capacity], &format[0], arg)
		if n != 3 || dst[0] != 0xa55a || dst[len(dst)-1] != 0xa55a {
			t.Fatalf("capacity%d: %x/%d", capacity, dst, n)
		}
		if capacity > 0 {
			end := min(3, capacity-1)
			if !slices.Equal(dst[1:1+end], text[:end]) || dst[1+end] != 0 {
				t.Fatalf("termination capacity%d: %x", capacity, dst)
			}
		}
	}
	// Temporary text no longer depends on the old 256-word C stack buffer.
	long := make([]uint16, 1025)
	for i := 0; i < 1024; i++ {
		long[i] = 0xd800
	}
	got := textFormatUnits(&format[0], textFormatPointer(unsafe.Pointer(&long[0])))
	if !slices.Equal(got, long[:1024]) {
		t.Fatal("long/raw UTF16 temporary")
	}
	nulFormat := []uint16{'A', '%', 'c', 'B', 0}
	if got := textFormatUnits(&nulFormat[0], textFormatWord(0)); !slices.Equal(got, []uint16{'A'}) {
		t.Fatalf("embedded NUL: %x", got)
	}
	// Original line delivery guards the recipient before touching a format.
	if textFormatLine(nil, nil) != 0 {
		t.Fatal("nil recipient")
	}
}

func TestTextFormatFixtureWordTypes(t *testing.T) {
	f := []uint16{'%', 'u', ' ', '%', 'd', 0}
	for _, word := range []uint32{0, 1, 255, 0x7fffffff, 0x80000000, 0xffffffff} {
		args := gameplayTextFixtureArguments(&f[0], []uint32{word, word})
		for _, a := range args {
			if a.ptr != nil || a.word != word {
				t.Fatalf("integer converted to pointer: %x", word)
			}
		}
	}
	f = []uint16{'%', 's', ' ', '%', 'u', 0}
	text := []uint16{'x', 0}
	args := gameplayTextFixtureArguments(&f[0], []uint32{uint32(uintptr(unsafe.Pointer(&text[0]))), 1})
	if args[0].ptr != unsafe.Pointer(&text[0]) || args[1].ptr != nil {
		t.Fatal("typed fixture arguments")
	}
}
