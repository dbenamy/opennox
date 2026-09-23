//go:build porttest

package legacy

import (
	"bytes"
	"encoding/binary"
	"unsafe"

	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

// PortTestJournalQueue owns the existing important-message allocator and queue.
// It only observes queued game reports; no transport or replacement serializer.
func PortTestJournalQueue() (reset, restore func()) {
	words := []*uint32{(*uint32)(&dword_5d4594_1565512), (*uint32)(&dword_5d4594_1565516), (*uint32)(&dword_5d4594_1565520), (*uint32)(&dword_5d4594_2649712), memmap.PtrUint32(0x5D4594, 1565508), (*uint32)(&dword_8531A0_2572)}
	old := make([]uint32, len(words))
	for i, p := range words {
		old[i] = *p
		*p = 0
	}
	seq := unsafe.Slice((*byte)(memmap.PtrOff(0x5D4594, 1565524)), 64)
	before := bytes.Clone(seq)
	reset = func() {
		if p := *memmap.PtrPtr(0x5D4594, 1565508); p != nil {
			alloc.AsClass(p).FreeAllObjects()
		}
		*words[0] = 0
		*words[1] = 0
		*words[3] = 0xffffffff
		*words[5] = 0x3def3def
		clear(seq)
	}
	restore = func() {
		if p := *memmap.PtrPtr(0x5D4594, 1565508); p != nil {
			alloc.AsClass(p).Free()
		}
		for i, p := range words {
			*p = old[i]
		}
		copy(seq, before)
	}
	reset()
	return
}
func PortTestJournalPackets() []PortTestShopPacketResult {
	var out []PortTestShopPacketResult
	seen := map[uint32]bool{}
	for q := uint32(dword_5d4594_1565512); q != 0; {
		if seen[q] {
			panic("journal fixture queue cycle")
		}
		seen[q] = true
		b := unsafe.Slice((*byte)(unsafe.Pointer(uintptr(q))), 416)
		v := PortTestShopPacketResult{Recipient: b[250], Ordered: b[184], A4: binary.LittleEndian.Uint32(b[404:]), A5: binary.LittleEndian.Uint32(b[180:]), Data: bytes.Clone(b[251 : 251+int(b[401])])}
		for i := range v.Sequence {
			v.Sequence[i] = binary.LittleEndian.Uint16(b[186+2*i:])
		}
		out = append(out, v)
		q = binary.LittleEndian.Uint32(b[408:])
	}
	return out
}
