package legacy

/*
#include <stdint.h>
extern uint32_t dword_5d4594_2516356;
*/
import "C"

//export sub_56F250
func sub_56F250() C.int {
	var result C.int
	for i := 0; i < 7; i++ {
		result = createProtectionRecord(uint32(C.dword_5d4594_2516356), 0)
		// Reserved slots consume an ID even if record allocation fails.
		C.dword_5d4594_2516356++
	}
	return result
}

//export nox_xxx_protectionCreateInt_56F400
func nox_xxx_protectionCreateInt_56F400(value C.int) C.int {
	if createProtectionRecord(uint32(C.dword_5d4594_2516356), uint32(value)) == 0 {
		return 0
	}
	id := C.dword_5d4594_2516356
	C.dword_5d4594_2516356++
	return C.int(id)
}
