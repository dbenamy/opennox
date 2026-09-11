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
	"unsafe"

	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/internal/protection"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func portTestResourceProtection(p *portTestShopPools, sp *PortTestResourceSpec) (func() []uint32, func()) {
	type guarded struct {
		Left   [2]uint32
		Record protection.Record
		Right  [2]uint32
	}
	records, free := alloc.Make([]guarded{}, 7)
	head, tail, key, sum, seq := C.dword_5d4594_2516344, C.dword_5d4594_2516352, C.dword_5d4594_2516348, C.dword_5d4594_2516328, C.dword_5d4594_2516356
	count := memmap.PtrUint16(0x587000, 311204)
	oldCount := *count
	swaps, rekeys := memmap.PtrUint32(0x5D4594, 2516360), memmap.PtrUint32(0x5D4594, 2516364)
	oldSwaps, oldRekeys := *swaps, *rekeys
	oldRNG := protectionRandom
	values := []uint32{uint32(sp.HP), p.proxy.callbacks.shop.spec.Gold[0], uint32(sp.MaxHP), uint32(sp.Mana), uint32(sp.MaxMana), 0x1234, 0x5678}
	offsets := []uintptr{4584, 4588, 4592, 4596, 4600, 4632, 4632}
	const encode = uint32(0x12345678)
	var total uint32
	for i := range records {
		b := &records[i]
		b.Left = [2]uint32{0x12345678, 0xa5a5a5a5}
		b.Right = [2]uint32{0x5a5a5a5a, 0xabcdef01}
		id := uint32(0x40000010 + i)
		b.Record = protection.Record{ID: id ^ encode, Value: values[i] ^ encode}
		if i > 0 {
			b.Record.Prev = &records[i-1].Record
		}
		if i+1 < len(records) {
			b.Record.Next = &records[i+1].Record
		}
		p.identify(unsafe.Pointer(&b.Record), uint32(54100+i))
		total ^= b.Record.Value
		player := 0
		if i == 6 {
			player = 1
		}
		dst := (*uint32)(unsafe.Add(unsafe.Pointer(p.proxy.life.players[player].UpdateDataPlayer().Player), offsets[i]))
		*dst = 0
		if sp.Protected {
			*dst = id
		}
	}
	C.dword_5d4594_2516344 = C.uint32_t(uintptr(unsafe.Pointer(&records[0].Record)))
	C.dword_5d4594_2516352 = C.uint32_t(uintptr(unsafe.Pointer(&records[len(records)-1].Record)))
	C.dword_5d4594_2516348, C.dword_5d4594_2516328, C.dword_5d4594_2516356 = C.uint32_t(encode), C.uint32_t(total), 0
	*count, *swaps, *rekeys = uint16(len(records)), 0, 0
	protectionRandom.Seed(1)
	snapshot := func() []uint32 {
		key := uint32(C.dword_5d4594_2516348)
		out := []uint32{key, uint32(C.dword_5d4594_2516328), uint32(C.dword_5d4594_2516356), uint32(*count), *swaps, *rekeys, p.normalize(uint32(C.dword_5d4594_2516344)), p.normalize(uint32(C.dword_5d4594_2516352))}
		for i := range records {
			b := &records[i]
			if b.Left != [2]uint32{0x12345678, 0xa5a5a5a5} || b.Right != [2]uint32{0x5a5a5a5a, 0xabcdef01} {
				panic("resource protection guard")
			}
			out = append(out, b.Record.ID^key, b.Record.Value^key, p.normalize(uint32(uintptr(unsafe.Pointer(b.Record.Next)))), p.normalize(uint32(uintptr(unsafe.Pointer(b.Record.Prev)))))
		}
		state, interval := portTestRandomState(protectionRandom)
		for i := 0; i < len(state); i += 4 {
			out = append(out, binary.LittleEndian.Uint32(state[i:]))
		}
		return append(out, interval[:]...)
	}
	restore := func() {
		C.dword_5d4594_2516344, C.dword_5d4594_2516352 = head, tail
		C.dword_5d4594_2516348, C.dword_5d4594_2516328, C.dword_5d4594_2516356 = key, sum, seq
		*count, *swaps, *rekeys = oldCount, oldSwaps, oldRekeys
		protectionRandom = oldRNG
		free()
	}
	return snapshot, restore
}
