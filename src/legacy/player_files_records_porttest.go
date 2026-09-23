//go:build porttest

package legacy

/*
#include <stdint.h>
*/
import "C"
import (
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/internal/protection"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"unsafe"
)

// PortTestPlayerFileRecords owns the actual spell/guide bookkeeping list, so
// section restoration exercises the real award functions and checksum updates.
func PortTestPlayerFileRecords() (ids [2]uint32, reset func(), snapshot func() [3]uint32, release func()) {
	oldHead, oldKey, oldSum := dword_5d4594_2516344, dword_5d4594_2516348, dword_5d4594_2516328
	records, free := alloc.Make([]protection.Record{}, 2)
	ids = [2]uint32{0x30000001, 0x30000002}
	const key = uint32(0x13579abc)
	reset = func() {
		records[0] = protection.Record{ID: ids[0] ^ key, Value: key, Next: &records[1]}
		records[1] = protection.Record{ID: ids[1] ^ key, Value: key, Prev: &records[0]}
		dword_5d4594_2516344 = C.uint32_t(uintptr(unsafe.Pointer(&records[0])))
		dword_5d4594_2516348 = C.uint32_t(key)
		dword_5d4594_2516328 = 0x2468ace0
	}
	snapshot = func() [3]uint32 {
		return [3]uint32{records[0].Value ^ key, records[1].Value ^ key, uint32(dword_5d4594_2516328)}
	}
	release = func() {
		dword_5d4594_2516344, dword_5d4594_2516348, dword_5d4594_2516328 = oldHead, oldKey, oldSum
		free()
	}
	reset()
	return
}

// PortTestPlayerFileEnchants owns the actual iteration table and live count.
func PortTestPlayerFileEnchants(ids []uint32) func() {
	if len(ids) > 29 {
		panic("enchant table capacity")
	}
	table := unsafe.Slice(memmap.PtrUint32(0x587000, 66000), 29)
	old := append([]uint32(nil), table...)
	count := bookEnchantN
	clear(table)
	copy(table, ids)
	bookEnchantN = int32(len(ids))
	return func() { copy(table, old); bookEnchantN = count }
}
