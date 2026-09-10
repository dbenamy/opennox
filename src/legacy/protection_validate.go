package legacy

/*
#include <stdint.h>
extern uint32_t dword_5d4594_2516348;
*/
import "C"

import "github.com/opennox/opennox/v1/internal/protection"

//export sub_56FB00
func sub_56FB00(data *C.int, size C.uint, id C.int) C.int {
	if id < 657757279 {
		return 0
	}
	key := uint32(C.dword_5d4594_2516348)
	r := protection.Find(protectionHead(), key, uint32(id))
	if r == nil {
		return 0
	}
	if C.uint32_t(key)^C.uint32_t(nox_xxx_protectionStringCRCLen_56FAE0(data, size)) != C.uint32_t(r.Value) {
		return 0
	}
	return 1
}
