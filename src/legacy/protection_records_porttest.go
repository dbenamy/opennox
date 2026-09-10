//go:build porttest

package legacy

/*
#include "GAME5_2.h"
extern uint32_t dword_5d4594_2516344;
extern uint32_t dword_5d4594_2516348;
*/
import "C"
import (
	"unsafe"

	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

type PortTestRecordResult struct {
	Lookup, At     int
	Values         [][2]uint32
	Counter        uint32
	LinksUnchanged bool
}

// PortTestRecords owns a temporary C-heap list and restores all touched globals.
// Call serially: the production protection manager is global state.
func PortTestRecords(values [][2]uint32, key uint32, id, index, a, b int, counter uint32) PortTestRecordResult {
	oldHead, oldKey := C.dword_5d4594_2516344, C.dword_5d4594_2516348
	count := memmap.PtrUint32(0x5D4594, 2516360)
	oldCount := *count
	defer func() {
		C.dword_5d4594_2516344, C.dword_5d4594_2516348 = oldHead, oldKey
		*count = oldCount
	}()
	var records [][4]uint32
	if len(values) != 0 {
		var free func()
		records, free = alloc.Make([][4]uint32{}, len(values))
		defer free()
	}
	ptr := func(i int) *C.int {
		if i < 0 || i >= len(records) {
			return nil
		}
		return (*C.int)(unsafe.Pointer(&records[i]))
	}
	for i, v := range values {
		records[i] = [4]uint32{v[0], v[1], uint32(uintptr(unsafe.Pointer(ptr(i + 1)))), uint32(uintptr(unsafe.Pointer(ptr(i - 1))))}
	}
	before := append([][4]uint32(nil), records...)
	C.dword_5d4594_2516344 = C.uint(uintptr(unsafe.Pointer(ptr(0))))
	C.dword_5d4594_2516348 = C.uint(key)
	*count = counter
	findIndex := func(p unsafe.Pointer) int {
		if p == nil {
			return -1
		}
		for i := range records {
			if p == unsafe.Pointer(&records[i]) {
				return i
			}
		}
		return -2 // invalid non-null result must not be confused with a miss
	}
	out := PortTestRecordResult{
		Lookup: findIndex(unsafe.Pointer(C.sub_56F590(C.int(id)))),
		At:     findIndex(unsafe.Pointer(C.sub_56F6F0(C.int(index)))),
	}
	C.sub_56F720(ptr(a), ptr(b))
	out.Counter = *count
	out.LinksUnchanged = C.dword_5d4594_2516344 == C.uint(uintptr(unsafe.Pointer(ptr(0)))) && C.dword_5d4594_2516348 == C.uint(key)
	for i, r := range records {
		out.Values = append(out.Values, [2]uint32{r[0], r[1]})
		out.LinksUnchanged = out.LinksUnchanged && r[2] == before[i][2] && r[3] == before[i][3]
	}
	return out
}
