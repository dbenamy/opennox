//go:build porttest

package legacy

/*
#include <stdlib.h>
#include "GAME5_2.h"
extern uint32_t dword_5d4594_2516344;
extern uint32_t dword_5d4594_2516352;
extern uint32_t dword_5d4594_2516348;
extern uint32_t dword_5d4594_2516328;
extern uint32_t dword_5d4594_2516356;
*/
import "C"
import (
	"unsafe"

	"github.com/opennox/opennox/v1/common/memmap"
)

type PortTestRemovalSnapshot struct {
	Result             int
	Handle             uint32
	Indices            []int
	Values             [][2]uint32
	Sum, Key, Sequence uint32
	Count              uint16
	LinksValid         bool
}

func PortTestRemove(values [][2]uint32, key, sum uint32, count uint16, ids []uint32) []PortTestRemovalSnapshot {
	oldHead, oldTail, oldKey, oldSum, oldSeq := C.dword_5d4594_2516344, C.dword_5d4594_2516352, C.dword_5d4594_2516348, C.dword_5d4594_2516328, C.dword_5d4594_2516356
	counter := memmap.PtrUint16(0x587000, 311204)
	oldCount := *counter
	defer func() {
		Sub_56F3B0()
		C.dword_5d4594_2516344, C.dword_5d4594_2516352, C.dword_5d4594_2516348, C.dword_5d4594_2516328, C.dword_5d4594_2516356 = oldHead, oldTail, oldKey, oldSum, oldSeq
		*counter = oldCount
	}()
	C.dword_5d4594_2516344, C.dword_5d4594_2516352 = 0, 0
	C.dword_5d4594_2516348, C.dword_5d4594_2516328, C.dword_5d4594_2516356 = C.uint(key), C.uint(sum), 0xdeadbeef
	*counter = count
	records := make([]*[4]uint32, len(values))
	for i, v := range values {
		r := (*[4]uint32)(C.calloc(1, 16))
		if r == nil {
			panic("fixture allocation failed")
		}
		records[i] = r
		*r = [4]uint32{v[0], v[1], 0, 0}
	}
	ptr := func(i int) uint32 {
		if i < 0 || i >= len(records) {
			return 0
		}
		return uint32(uintptr(unsafe.Pointer(records[i])))
	}
	for i, r := range records {
		r[2], r[3] = ptr(i+1), ptr(i-1)
	}
	C.dword_5d4594_2516344, C.dword_5d4594_2516352 = C.uint(ptr(0)), C.uint(ptr(len(records)-1))
	snapshot := func() PortTestRemovalSnapshot {
		out := PortTestRemovalSnapshot{Sum: uint32(C.dword_5d4594_2516328), Key: uint32(C.dword_5d4594_2516348), Sequence: uint32(C.dword_5d4594_2516356), Count: *counter, LinksValid: true}
		p := uint32(C.dword_5d4594_2516344)
		prev := uint32(0)
		for p != 0 {
			index := -1
			for i := range records {
				if ptr(i) == p {
					index = i
					break
				}
			}
			if index < 0 || len(out.Indices) >= len(records) {
				out.LinksValid = false
				break
			}
			r := records[index]
			out.Indices = append(out.Indices, index)
			out.Values = append(out.Values, [2]uint32{r[0], r[1]})
			out.LinksValid = out.LinksValid && r[3] == prev
			prev, p = p, r[2]
		}
		out.LinksValid = out.LinksValid && prev == uint32(C.dword_5d4594_2516352)
		return out
	}
	var out []PortTestRemovalSnapshot
	for _, id := range ids {
		handle := C.int(id)
		result := int(C.sub_56F4F0(&handle))
		s := snapshot()
		s.Result = result
		s.Handle = uint32(handle)
		out = append(out, s)
	}
	Sub_56F3B0()
	out = append(out, snapshot())
	return out
}
