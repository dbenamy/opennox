//go:build porttest

package legacy

/*
#include "GAME1.h"
#include "GAME3_2.h"
*/
import "C"
import "unsafe"
import "github.com/opennox/opennox/v1/client/gui"

// Invoke the original configuration operations; tests own mapped fields, players,
// report state and packet queues. Do not duplicate any production algorithm here.
func PortTestServerConfigScalar(op string, a, b int) uint64 {
	switch op {
	case "flags-set":
		return uint64(C.sub_409E40(C.int(a)))
	case "flags-get":
		return uint64(C.nox_xxx_getServerSubFlags_409E60())
	case "flags-add":
		return uint64(C.sub_409E70(C.int(a)))
	case "flags-remove":
		return uint64(C.sub_409EC0(C.int(a)))
	case "flags-toggle":
		return uint64(C.sub_409EF0(C.int(a)))
	case "flags-query":
		return uint64(C.sub_409F40(C.int(a)))
	case "limit-set":
		return uint64(C.nox_xxx_servSetPlrLimit_409F80(C.int(a)))
	case "limit-get":
		return uint64(C.nox_xxx_servGetPlrLimit_409FA0())
	case "score":
		return uint64(C.nox_xxx_servGamedataGet_40A020(C.short(a)))
	case "minutes":
		return uint64(C.sub_40A180(C.short(a)))
	case "timer-set":
		return uint64(C.sub_40A1F0(C.int(a)))
	case "timer-get":
		return uint64(C.sub_40A220())
	case "timer-left":
		return uint64(C.sub_40A230())
	case "timer-init":
		return uint64(C.sub_40A250())
	case "timer-reset":
		return uint64(C.sub_40A310(C.int(a)))
	case "3512-set":
		C.nox_xxx_set3512_40A340(C.int(a))
		return 0
	case "3512-get":
		return uint64(C.nox_xxx_get3512_40A350())
	case "mode-store":
		return uint64(C.sub_40A3C0(C.uint(a)))
	case "respawn-set":
		return uint64(C.nox_xxx_ruleSetNoRespawn_40A5E0(C.int(a)))
	case "respawn-get":
		return uint64(C.nox_server_doPlayersAutoRespawn_40A5F0())
	case "updated-set":
		C.nox_server_gameSettingsUpdated_40A670()
		return 0
	case "updated-get":
		return uint64(C.nox_server_gameDoSwitchMap_40A680())
	case "updated-clear":
		C.nox_server_gameUnsetMapLoad_40A690()
		return 0
	case "rate-dirty-set":
		return uint64(C.sub_40A6A0(C.int(a)))
	case "rate-dirty-get":
		return uint64(C.sub_40A6B0())
	case "rate-get":
		return uint64(C.nox_xxx_rateGet_40A6C0())
	case "rate-set":
		return uint64(C.nox_xxx_rateUpdate_40A6D0(C.int(a)))
	case "connection-rate":
		return uint64(C.sub_40A710(C.int(a)))
	case "special-mode":
		return uint64(C.sub_40A740())
	case "refresh":
		return uint64(C.sub_4161E0())
	case "slot-index":
		return uint64(C.sub_416580())
	case "slot-copy":
		return uint64(C.sub_4165F0(C.int(a), C.int(b)))
	case "record-state":
		return uint64(C.sub_416650())
	case "acquired-get":
		return uint64(C.sub_4169C0())
	case "acquired-set":
		return uint64(C.nox_xxx_cliSetSettingsAcquired_4169D0(C.int(a)))
	default:
		panic(op)
	}
}
func PortTestServerConfigPointer(op string, a int, p unsafe.Pointer) unsafe.Pointer {
	switch op {
	case "name-set":
		return unsafe.Pointer(C.nox_xxx_gameSetServername_40A440((*C.char)(p)))
	case "name-get":
		return unsafe.Pointer(C.nox_xxx_serverOptionsGetServername_40A4C0())
	case "password-set":
		return unsafe.Pointer(C.nox_xxx_sysopSetPass_40A610((*C.wchar2_t)(p)))
	case "password-get":
		return unsafe.Pointer(C.nox_xxx_sysopGetPass_40A630())
	case "slot":
		return unsafe.Pointer(C.nox_xxx_cliGamedataGet_416590(C.int(a)))
	case "slot-current":
		return unsafe.Pointer(C.sub_4165B0())
	case "slot-select":
		return unsafe.Pointer(C.sub_4165D0(C.int(a)))
	case "admission":
		return unsafe.Pointer(C.sub_416630())
	case "settings":
		return C.sub_416640()
	case "open-admission":
		return unsafe.Pointer(C.sub_4169F0())
	default:
		panic(op)
	}
}

func PortTestServerConfigRuleOpen(parent *gui.Window, settings unsafe.Pointer) int {
	return int(C.sub_4CEBA0(C.int(uintptr(parent.C())), (*C.char)(settings)))
}
