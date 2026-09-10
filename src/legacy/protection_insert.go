package legacy

/*
#include <stdint.h>
extern uint32_t dword_5d4594_2516344;
extern uint32_t dword_5d4594_2516352;
*/
import "C"
import (
	"unsafe"

	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/internal/protection"
)

func insertProtectionRecord(r *protection.Record) C.int {
	count := memmap.PtrUint16(0x587000, 311204)
	if *count == 0 {
		(*count)++
		C.dword_5d4594_2516352 = C.uint32_t(uintptr(unsafe.Pointer(r)))
		C.dword_5d4594_2516344 = C.uint32_t(uintptr(unsafe.Pointer(r)))
		return 1
	}

	target := GetServer().S().Rand.Logic.IntClamp(0, int(*count)-1)
	before := protectionHead()
	for i := 0; i < target; i++ {
		before = before.Next
	}
	head := protection.InsertBefore(protectionHead(), before, r)
	C.dword_5d4594_2516344 = C.uint32_t(uintptr(unsafe.Pointer(head)))
	(*count)++
	return 1
}
