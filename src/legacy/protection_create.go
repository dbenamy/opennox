package legacy

/*
#include <stdint.h>
#include <stdlib.h>
extern uint32_t dword_5d4594_2516348;
extern uint32_t dword_5d4594_2516328;
int sub_56F2F0(uint32_t* record);
*/
import "C"
import (
	"unsafe"

	"github.com/opennox/opennox/v1/internal/protection"
)

func createProtectionRecord(id, value uint32) C.int {
	r := (*protection.Record)(C.calloc(1, C.size_t(unsafe.Sizeof(protection.Record{}))))
	if !protection.Initialize(r, id, value, uint32(C.dword_5d4594_2516348), (*uint32)(unsafe.Pointer(&C.dword_5d4594_2516328))) {
		return 0
	}
	// The existing C insertion routine takes ownership and preserves RNG use.
	return C.sub_56F2F0((*C.uint32_t)(unsafe.Pointer(r)))
}

//export nox_xxx_protectionCreateStructForInt_56F280
func nox_xxx_protectionCreateStructForInt_56F280(id, value C.int) C.int {
	return createProtectionRecord(uint32(id), uint32(value))
}
