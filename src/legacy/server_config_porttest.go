//go:build porttest

package legacy

import "unsafe"
import "github.com/opennox/opennox/v1/client/gui"

// Invoke the native configuration operations; tests own mapped fields, players,
// report state and packet queues. Do not duplicate any production algorithm here.
func PortTestServerConfigScalar(op string, a, b int) uint64 {
	switch op {
	case "flags-set":
		return uint64(int32(serverConfigFlagsSet(int32(int32(a)))))
	case "flags-get":
		return uint64(int32(serverConfigFlagsGet()))
	case "flags-add":
		return uint64(int32(serverConfigFlagsAdd(int32(int32(a)))))
	case "flags-remove":
		return uint64(int32(serverConfigFlagsRemove(int32(int32(a)))))
	case "flags-toggle":
		return uint64(int32(serverConfigFlagsToggle(int32(int32(a)))))
	case "flags-query":
		return uint64(int32(serverConfigFlagsQuery(int32(int32(a)))))
	case "limit-set":
		return uint64(int32(serverConfigLimitSet(int32(int32(a)))))
	case "limit-get":
		return uint64(int32(serverConfigLimitGet()))
	case "score":
		return uint64(int16(serverConfigScore(int16(int16(a)))))
	case "minutes":
		return uint64(byte(serverConfigMinutes(int16(int16(a)))))
	case "timer-set":
		return uint64(int32(serverConfigTimerSet(int32(int32(a)))))
	case "timer-get":
		return uint64(int32(serverConfigTimerGet()))
	case "timer-left":
		return uint64(int32(serverConfigTimerLeft()))
	case "timer-init":
		return uint64(int64(serverConfigTimerInit()))
	case "timer-reset":
		return uint64(int64(serverConfigTimerReset(int32(int32(a)))))
	case "3512-set":
		serverConfig3512Set(int32(int32(a)))
		return 0
	case "3512-get":
		return uint64(int32(serverConfig3512Get()))
	case "mode-store":
		return uint64(uint32(serverConfigModeStore(uint32(uint32(a)))))
	case "respawn-set":
		return uint64(int32(serverConfigRespawnSet(int32(int32(a)))))
	case "respawn-get":
		return uint64(int32(serverConfigRespawnGet()))
	case "updated-set":
		serverConfigUpdatedSet()
		return 0
	case "updated-get":
		return uint64(int32(serverConfigUpdatedGet()))
	case "updated-clear":
		serverConfigUpdatedClear()
		return 0
	case "rate-dirty-set":
		return uint64(int32(serverConfigRateDirtySet(int32(int32(a)))))
	case "rate-dirty-get":
		return uint64(int32(serverConfigRateDirtyGet()))
	case "rate-get":
		return uint64(int32(serverConfigRateGet()))
	case "rate-set":
		return uint64(int32(serverConfigRateSet(int32(int32(a)))))
	case "connection-rate":
		return uint64(int32(serverConfigConnectionRate(int32(int32(a)))))
	case "special-mode":
		return uint64(int32(serverConfigSpecialMode()))
	case "refresh":
		return uint64(int32(serverConfigRefresh()))
	case "slot-index":
		return uint64(int32(serverConfigSlotIndex()))
	case "slot-copy":
		return uint64(int32(serverConfigSlotCopy(int32(int32(a)), int32(int32(b)))))
	case "record-state":
		return uint64(int32(serverConfigRecordState()))
	case "acquired-get":
		return uint64(int32(serverConfigAcquiredGet()))
	case "acquired-set":
		return uint64(int32(serverConfigAcquiredSet(int32(int32(a)))))
	default:
		panic(op)
	}
}
func PortTestServerConfigPointer(op string, a int, p unsafe.Pointer) unsafe.Pointer {
	switch op {
	case "name-set":
		return unsafe.Pointer((*byte)(unsafe.Pointer(serverConfigNameSet((*byte)(unsafe.Pointer((*byte)(p)))))))
	case "name-get":
		return unsafe.Pointer((*byte)(unsafe.Pointer(serverConfigNameGet())))
	case "password-set":
		return unsafe.Pointer((*uint16)(unsafe.Pointer(serverConfigPasswordSet((*uint16)(unsafe.Pointer((*uint16)(p)))))))
	case "password-get":
		return unsafe.Pointer((*uint16)(unsafe.Pointer(serverConfigPasswordGet())))
	case "slot":
		return unsafe.Pointer((*byte)(unsafe.Pointer(serverConfigSlot(int32(int32(a))))))
	case "slot-current":
		return unsafe.Pointer((*byte)(unsafe.Pointer(serverConfigSlotCurrent())))
	case "slot-select":
		return unsafe.Pointer((*byte)(unsafe.Pointer(serverConfigSlotSelect(int32(int32(a))))))
	case "admission":
		return unsafe.Pointer((*byte)(unsafe.Pointer(serverConfigAdmissionData())))
	case "settings":
		return (unsafe.Pointer)(unsafe.Pointer(serverConfigSettings()))
	case "open-admission":
		return unsafe.Pointer((*byte)(unsafe.Pointer(serverConfigOpenAdmission())))
	default:
		panic(op)
	}
}

func PortTestServerConfigRuleOpen(parent *gui.Window, settings unsafe.Pointer) int {
	return int(int32(serverConfigRuleOpen((*gui.Window)(unsafe.Pointer(uintptr(uint32(int32(uintptr(parent.C())))))), (*byte)(unsafe.Pointer((*byte)(settings))))))
}
