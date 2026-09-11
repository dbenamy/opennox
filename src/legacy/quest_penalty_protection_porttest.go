//go:build porttest

package legacy

/*
#include <stdint.h>
extern uint32_t dword_5d4594_2516344, dword_5d4594_2516352;
extern uint32_t dword_5d4594_2516348, dword_5d4594_2516328, dword_5d4594_2516356;
*/
import "C"
import (
	"encoding/binary"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/internal/protection"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"unsafe"
)

func portTestPenaltyProtectionEnvironment() (prepare func(uint32), snapshot func() ([]uint32, bool), restore func()) {
	type guarded struct {
		Left   [2]uint32
		Record protection.Record
		Right  [2]uint32
	}
	b, free := alloc.New(guarded{})
	b.Left = [2]uint32{0x12345678, 0xa5a5a5a5}
	b.Right = [2]uint32{0x5a5a5a5a, 0xabcdef01}
	head, tail, key, sum, seq := C.dword_5d4594_2516344, C.dword_5d4594_2516352, C.dword_5d4594_2516348, C.dword_5d4594_2516328, C.dword_5d4594_2516356
	count := memmap.PtrUint16(0x587000, 311204)
	oldCount := *count
	swaps, rekeys := memmap.PtrUint32(0x5D4594, 2516360), memmap.PtrUint32(0x5D4594, 2516364)
	oldSwaps, oldRekeys := *swaps, *rekeys
	rng := protectionRandom
	prepare = func(gold uint32) {
		const key = uint32(0x12345678)
		b.Record = protection.Record{ID: 0x40000001 ^ key, Value: gold ^ key}
		C.dword_5d4594_2516344 = C.uint32_t(uintptr(unsafe.Pointer(&b.Record)))
		C.dword_5d4594_2516352 = C.dword_5d4594_2516344
		C.dword_5d4594_2516348 = C.uint32_t(key)
		C.dword_5d4594_2516328 = C.uint32_t(gold ^ key)
		C.dword_5d4594_2516356 = 0
		*count, *swaps, *rekeys = 1, 0, 0
		protectionRandom.Seed(1)
	}
	snapshot = func() ([]uint32, bool) {
		key := uint32(C.dword_5d4594_2516348)
		out := []uint32{b.Record.ID ^ key, b.Record.Value ^ key, key, uint32(C.dword_5d4594_2516328), uint32(C.dword_5d4594_2516356), uint32(*count), *swaps, *rekeys}
		state, interval := portTestRandomState(protectionRandom)
		for i := 0; i < len(state); i += 4 {
			out = append(out, binary.LittleEndian.Uint32(state[i:]))
		}
		out = append(out, interval[:]...)
		intact := b.Left == [2]uint32{0x12345678, 0xa5a5a5a5} && b.Right == [2]uint32{0x5a5a5a5a, 0xabcdef01} && b.Record.Next == nil && b.Record.Prev == nil
		return out, intact
	}
	restore = func() {
		C.dword_5d4594_2516344, C.dword_5d4594_2516352 = head, tail
		C.dword_5d4594_2516348, C.dword_5d4594_2516328, C.dword_5d4594_2516356 = key, sum, seq
		*count, *swaps, *rekeys = oldCount, oldSwaps, oldRekeys
		protectionRandom = rng
		free()
	}
	return
}
