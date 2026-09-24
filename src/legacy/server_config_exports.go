package legacy

/*
#include <stdint.h>
typedef uint16_t wchar2_t;
*/
import "C"
import "unsafe"

//export nox_xxx_cliGamedataGet_416590
func nox_xxx_cliGamedataGet_416590(a1 C.int) *C.char {
	return (*C.char)(unsafe.Pointer(serverConfigSlot(int32(a1))))
}

//export sub_4165B0
func sub_4165B0() *C.char { return (*C.char)(unsafe.Pointer(serverConfigSlotCurrent())) }

//export sub_416640
func sub_416640() unsafe.Pointer { return (unsafe.Pointer)(unsafe.Pointer(serverConfigSettings())) }
