package legacy

/*
#include <stdint.h>
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"

	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/internal/protection"
)

func deleteProtectionRecord(id uint32) bool {
	head := protectionHead()
	r := protection.Find(head, uint32(dword_5d4594_2516348), id)
	if r == nil {
		return false
	}
	tail := *(**protection.Record)(unsafe.Pointer(&dword_5d4594_2516352))
	head, tail = protection.Unlink(head, tail, r)
	dword_5d4594_2516344 = C.uint32_t(uintptr(unsafe.Pointer(head)))
	dword_5d4594_2516352 = C.uint32_t(uintptr(unsafe.Pointer(tail)))
	dword_5d4594_2516328 ^= C.uint32_t(r.ID ^ r.Value)
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
	dword_5d4594_2516328 = 0
	*memmap.PtrUint16(0x587000, 311204) = 0
	dword_5d4594_2516348 = 0
	dword_5d4594_2516352 = 0
	dword_5d4594_2516344 = 0
}
