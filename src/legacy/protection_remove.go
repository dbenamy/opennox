package legacy

/*
#include <stdint.h>
#include <stdlib.h>
extern uint32_t dword_5d4594_2516344;
extern uint32_t dword_5d4594_2516352;
extern uint32_t dword_5d4594_2516348;
extern uint32_t dword_5d4594_2516328;
*/
import "C"
import (
	"unsafe"

	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/internal/protection"
)

func deleteProtectionRecord(id uint32) bool {
	head := protectionHead()
	r := protection.Find(head, uint32(C.dword_5d4594_2516348), id)
	if r == nil {
		return false
	}
	tail := *(**protection.Record)(unsafe.Pointer(&C.dword_5d4594_2516352))
	head, tail = protection.Unlink(head, tail, r)
	C.dword_5d4594_2516344 = C.uint32_t(uintptr(unsafe.Pointer(head)))
	C.dword_5d4594_2516352 = C.uint32_t(uintptr(unsafe.Pointer(tail)))
	C.dword_5d4594_2516328 ^= C.uint32_t(r.ID ^ r.Value)
	*memmap.PtrUint16(0x587000, 311204)--
	C.free(unsafe.Pointer(r))
	return true
}

//export sub_56F4F0
func sub_56F4F0(id *C.int) C.int {
	if deleteProtectionRecord(uint32(*id)) {
		*id = 0
		return 1
	}
	return 0
}

func freeProtectionRecords() {
	for p := protectionHead(); p != nil; {
		next := p.Next
		C.free(unsafe.Pointer(p))
		p = next
	}
	C.dword_5d4594_2516328 = 0
	*memmap.PtrUint16(0x587000, 311204) = 0
	C.dword_5d4594_2516348 = 0
	C.dword_5d4594_2516352 = 0
	C.dword_5d4594_2516344 = 0
}
