package legacy

/*
#include <stdint.h>
extern uint32_t dword_5d4594_2516344;
extern uint32_t dword_5d4594_2516348;
*/
import "C"
import (
	"unsafe"

	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/internal/protection"
)

// The C manager stores payloads at words 0/1 and links at words 2/3.
var (
	_ [16 - unsafe.Sizeof(protection.Record{})]byte
	_ [unsafe.Sizeof(protection.Record{}) - 16]byte
	_ [4 - unsafe.Alignof(protection.Record{})]byte
	_ [unsafe.Alignof(protection.Record{}) - 4]byte
	_ [0 - unsafe.Offsetof(protection.Record{}.ID)]byte
	_ [4 - unsafe.Offsetof(protection.Record{}.Value)]byte
	_ [unsafe.Offsetof(protection.Record{}.Value) - 4]byte
	_ [8 - unsafe.Offsetof(protection.Record{}.Next)]byte
	_ [unsafe.Offsetof(protection.Record{}.Next) - 8]byte
	_ [12 - unsafe.Offsetof(protection.Record{}.Prev)]byte
	_ [unsafe.Offsetof(protection.Record{}.Prev) - 12]byte
)

func protectionHead() *protection.Record {
	return *(**protection.Record)(unsafe.Pointer(&C.dword_5d4594_2516344))
}

//export sub_56F590
func sub_56F590(id C.int) *C.uint32_t {
	p := protection.Find(protectionHead(), uint32(C.dword_5d4594_2516348), uint32(id))
	return (*C.uint32_t)(unsafe.Pointer(p))
}

//export sub_56F6F0
func sub_56F6F0(index C.int) *C.uint32_t {
	return (*C.uint32_t)(unsafe.Pointer(protection.At(protectionHead(), int32(index))))
}

//export sub_56F720
func sub_56F720(a, b *C.int) {
	if protection.Swap((*protection.Record)(unsafe.Pointer(a)), (*protection.Record)(unsafe.Pointer(b))) {
		*memmap.PtrUint32(0x5D4594, 2516360)++
	}
}
