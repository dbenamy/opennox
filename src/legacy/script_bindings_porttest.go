//go:build porttest

package legacy

/*
#include "defs.h"
extern unsigned int dword_5d4594_1599628;
int sub_512E80(wchar2_t*);
*/
import "C"
import "unsafe"

// Exercise the production callback-transfer adapter with fixture-owned records.
func PortTestScriptBindingCallback(record, name unsafe.Pointer) int {
	return objectXferScript(record, name)
}

func PortTestScriptBindingRegistryOwner() (*uint32, func()) {
	p := (*uint32)(unsafe.Pointer(&C.dword_5d4594_1599628))
	old := *p
	return p, func() { *p = old }
}

func PortTestScriptBindingIntern(text unsafe.Pointer) int {
	return int(C.sub_512E80((*C.wchar2_t)(text)))
}
