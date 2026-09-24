package legacy

/*
#include <stdint.h>
typedef uint16_t wchar2_t;
*/
import "C"
import "unsafe"

func sub_409E40(a1 C.int) C.int { return C.int(serverConfigFlagsSet(int32(a1))) }

//export sub_409F40
func sub_409F40(a1 C.int) C.int { return C.int(serverConfigFlagsQuery(int32(a1))) }

func nox_xxx_servSetPlrLimit_409F80(a1 C.int) C.int { return C.int(serverConfigLimitSet(int32(a1))) }

func sub_40A1F0(a1 C.int) C.int { return C.int(serverConfigTimerSet(int32(a1))) }

func sub_40A310(a1 C.int) C.longlong { return C.longlong(serverConfigTimerReset(int32(a1))) }

func nox_xxx_set3512_40A340(a1 C.int) { serverConfig3512Set(int32(a1)) }

func nox_xxx_get3512_40A350() C.int { return C.int(serverConfig3512Get()) }

//export nox_xxx_cliGamedataGet_416590
func nox_xxx_cliGamedataGet_416590(a1 C.int) *C.char {
	return (*C.char)(unsafe.Pointer(serverConfigSlot(int32(a1))))
}

//export sub_4165B0
func sub_4165B0() *C.char { return (*C.char)(unsafe.Pointer(serverConfigSlotCurrent())) }

//export sub_4165D0
func sub_4165D0(a1 C.int) *C.char {
	return (*C.char)(unsafe.Pointer(serverConfigSlotSelect(int32(a1))))
}

func sub_4165F0(a1 C.int, a2 C.int) C.int { return C.int(serverConfigSlotCopy(int32(a1), int32(a2))) }

//export sub_416640
func sub_416640() unsafe.Pointer { return (unsafe.Pointer)(unsafe.Pointer(serverConfigSettings())) }

func sub_4169C0() C.int { return C.int(serverConfigAcquiredGet()) }

func nox_xxx_cliSetSettingsAcquired_4169D0(a1 C.int) C.int {
	return C.int(serverConfigAcquiredSet(int32(a1)))
}
