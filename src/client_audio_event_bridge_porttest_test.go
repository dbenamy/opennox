//go:build porttest

package opennox

import (
	"bytes"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/legacy"
)

// Check the original C-facing lookup's address result without reading through
// out-of-slot pointers. High argument words and multiplication overflow are part
// of the existing 386 interface, not bounds checks added by this conversion.
func TestClientAudioEventMusicSlotAddress(t *testing.T) {
	words, restore := legacy.PortTestAudioEventGlobals()
	defer restore()
	memory := serverConfigOwnBytes(t, 0x5D4594, 815516, 656)
	for i := range memory {
		memory[i] = byte(i*37 + 11)
	}
	before := bytes.Clone(memory)
	count, level := words["dword_5d4594_816368"], words["dword_5d4594_816372"]
	cases := []struct {
		arg   uint64
		delta int
	}{
		{0, 0}, {1, 16}, {5, 80}, {6, 96},
		{0xffffffff, -16}, {0x7fffffff, -16}, {0x80000000, 0},
		{0x10000003, 48}, {0xfffffff0, -256},
		{0x100000001, 16}, {0xffffffff00000005, 80},
	}
	for l := uint32(0); l < 4; l++ {
		*count, *level = 0x98765432, l
		for _, tc := range cases {
			got := legacy.PortTestAudioEventCall("sub_43DB40", tc.arg)
			want := uint64(uintptr(unsafe.Pointer(&memory[256+96*int(l)+tc.delta])))
			if got != want {
				t.Fatalf("level %d, argument %#x: address %#x, want %#x", l, tc.arg, got, want)
			}
			if *count != 0x98765432 || *level != l || !bytes.Equal(memory, before) {
				t.Fatalf("lookup mutated memory or counters: level %d, argument %#x", l, tc.arg)
			}
		}
	}
}
