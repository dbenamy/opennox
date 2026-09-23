package legacy

/*
#include <stdint.h>
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"

	"github.com/opennox/opennox/v1/internal/protection"
)

func createProtectionRecord(id, value uint32) C.int {
	r := (*protection.Record)(C.calloc(1, C.size_t(unsafe.Sizeof(protection.Record{}))))
	if !protection.Initialize(r, id, value, uint32(dword_5d4594_2516348), (*uint32)(unsafe.Pointer(&dword_5d4594_2516328))) {
		return 0
	}
	return insertProtectionRecord(r)
}
