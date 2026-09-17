package legacy

/*
#include <stdint.h>
typedef uint16_t wchar2_t;
*/
import "C"
import "unsafe"

//export sub_409E40
func sub_409E40(a1 C.int) C.int { return C.int(serverConfigFlagsSet(int32(a1))) }

//export sub_409E70
func sub_409E70(a1 C.int) C.int { return C.int(serverConfigFlagsAdd(int32(a1))) }

//export sub_409EC0
func sub_409EC0(a1 C.int) C.int { return C.int(serverConfigFlagsRemove(int32(a1))) }

//export sub_409F40
func sub_409F40(a1 C.int) C.int { return C.int(serverConfigFlagsQuery(int32(a1))) }

//export nox_xxx_servSetPlrLimit_409F80
func nox_xxx_servSetPlrLimit_409F80(a1 C.int) C.int { return C.int(serverConfigLimitSet(int32(a1))) }

//export nox_xxx_servGetPlrLimit_409FA0
func nox_xxx_servGetPlrLimit_409FA0() C.int { return C.int(serverConfigLimitGet()) }

//export nox_xxx_servGamedataGet_40A020
func nox_xxx_servGamedataGet_40A020(a1 C.short) C.short { return C.short(serverConfigScore(int16(a1))) }

//export sub_40A180
func sub_40A180(a1 C.short) C.uchar { return C.uchar(serverConfigMinutes(int16(a1))) }

//export sub_40A1F0
func sub_40A1F0(a1 C.int) C.int { return C.int(serverConfigTimerSet(int32(a1))) }

//export sub_40A220
func sub_40A220() C.int { return C.int(serverConfigTimerGet()) }

//export sub_40A250
func sub_40A250() C.longlong { return C.longlong(serverConfigTimerInit()) }

//export sub_40A310
func sub_40A310(a1 C.int) C.longlong { return C.longlong(serverConfigTimerReset(int32(a1))) }

//export nox_xxx_set3512_40A340
func nox_xxx_set3512_40A340(a1 C.int) { serverConfig3512Set(int32(a1)) }

//export nox_xxx_get3512_40A350
func nox_xxx_get3512_40A350() C.int { return C.int(serverConfig3512Get()) }

//export sub_40A3C0
func sub_40A3C0(a1 C.uint) C.uint { return C.uint(serverConfigModeStore(uint32(a1))) }

//export nox_xxx_gameSetServername_40A440
func nox_xxx_gameSetServername_40A440(a1 *C.char) *C.char {
	return (*C.char)(unsafe.Pointer(serverConfigNameSet((*byte)(unsafe.Pointer(a1)))))
}

//export nox_xxx_serverOptionsGetServername_40A4C0
func nox_xxx_serverOptionsGetServername_40A4C0() *C.char {
	return (*C.char)(unsafe.Pointer(serverConfigNameGet()))
}

//export nox_xxx_sysopSetPass_40A610
func nox_xxx_sysopSetPass_40A610(a1 *C.wchar2_t) *C.wchar2_t {
	return (*C.wchar2_t)(unsafe.Pointer(serverConfigPasswordSet((*uint16)(unsafe.Pointer(a1)))))
}

//export nox_server_gameSettingsUpdated_40A670
func nox_server_gameSettingsUpdated_40A670() { serverConfigUpdatedSet() }

//export nox_server_gameUnsetMapLoad_40A690
func nox_server_gameUnsetMapLoad_40A690() { serverConfigUpdatedClear() }

//export nox_xxx_rateUpdate_40A6D0
func nox_xxx_rateUpdate_40A6D0(a1 C.int) C.int { return C.int(serverConfigRateSet(int32(a1))) }

//export sub_40A710
func sub_40A710(a1 C.int) C.int { return C.int(serverConfigConnectionRate(int32(a1))) }

//export sub_40A740
func sub_40A740() C.int { return C.int(serverConfigSpecialMode()) }

//export sub_4161E0
func sub_4161E0() C.int { return C.int(serverConfigRefresh()) }

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

//export sub_4165F0
func sub_4165F0(a1 C.int, a2 C.int) C.int { return C.int(serverConfigSlotCopy(int32(a1), int32(a2))) }

//export sub_416640
func sub_416640() unsafe.Pointer { return (unsafe.Pointer)(unsafe.Pointer(serverConfigSettings())) }

//export sub_416770
func sub_416770(a1 C.int, a2 *C.wchar2_t, a3 *C.char) *C.int {
	return (*C.int)(unsafe.Pointer(serverConfigBlockedAdd(int32(a1), (*uint16)(unsafe.Pointer(a2)), (*byte)(unsafe.Pointer(a3)))))
}

//export sub_4169C0
func sub_4169C0() C.int { return C.int(serverConfigAcquiredGet()) }

//export nox_xxx_cliSetSettingsAcquired_4169D0
func nox_xxx_cliSetSettingsAcquired_4169D0(a1 C.int) C.int {
	return C.int(serverConfigAcquiredSet(int32(a1)))
}
