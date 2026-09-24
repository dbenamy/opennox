//go:build porttest

package legacy

/*
#include <stdint.h>
#include <stdlib.h>
*/
import "C"
import (
	"math"
	"unsafe"

	"github.com/opennox/opennox/v1/common/memmap"
)

type PortTestCreateResult struct {
	Result         int
	ID, Value, Sum uint32
	Count          uint16
	LinksValid     bool
}

// PortTestCreate starts an empty manager, so insertion takes no random draw.
func PortTestCreate(id, bits, key, sum uint32, mode int) PortTestCreateResult {
	head, tail, oldKey, oldSum := dword_5d4594_2516344, dword_5d4594_2516352, dword_5d4594_2516348, dword_5d4594_2516328
	count := memmap.PtrUint16(0x587000, 311204)
	oldCount := *count
	defer func() {
		C.free(unsafe.Pointer(uintptr(dword_5d4594_2516344)))
		dword_5d4594_2516344, dword_5d4594_2516352, dword_5d4594_2516348, dword_5d4594_2516328 = head, tail, oldKey, oldSum
		*count = oldCount
	}()
	dword_5d4594_2516344, dword_5d4594_2516352 = 0, 0
	dword_5d4594_2516348, dword_5d4594_2516328 = uint32(key), uint32(sum)
	*count = 0
	var result int
	switch mode {
	case 0:
		result = Nox_xxx_protectionCreateStructForInt_56F280(int(id), int(bits))
	case 1:
		result = Nox_xxx_protectionCreateStructForFloat_56F480(int(id), math.Float32frombits(bits))
	default:
		panic(mode)
	}
	out := PortTestCreateResult{Result: result, Count: *count, Sum: uint32(dword_5d4594_2516328)}
	if p := dword_5d4594_2516344; p != 0 {
		r := (*[4]uint32)(unsafe.Pointer(uintptr(p)))
		out.ID, out.Value = r[0], r[1]
		out.LinksValid = p == dword_5d4594_2516352 && r[2] == 0 && r[3] == 0 && dword_5d4594_2516348 == uint32(key)
	}
	return out
}
