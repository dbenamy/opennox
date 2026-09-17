//go:build porttest

package legacy

/*
#include "defs.h"
#include "GAME4_2.h"
extern void* nox_monsterBin_head_2386924;
extern uint32_t dword_5d4594_588120;
*/
import "C"
import (
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

func PortTestCreatureXferAdjust(value, delta uint32) (uint32, uint32) {
	v := C.int(value)
	ret := C.nox_xxx_AssignIfGreater_52A420(&v, C.int(delta))
	return uint32(v), uint32(ret)
}
func PortTestCreatureXferHelper(op int, u *server.Object, entry unsafe.Pointer, delta uint32) uint32 {
	addr := C.int(uintptr(u.CObj()))
	switch op {
	case 0:
		return uint32(C.nox_xxx_XFer_ActionData_529CE0(addr))
	case 1:
		return uint32(C.sub_52A440(addr, C.int(uintptr(entry)), C.int(delta)))
	case 2:
		return uint32(C.nox_xxx_XFer_ReadMonsterBuffs_52AAB0((*C.uint32_t)(u.CObj())))
	case 3:
		return uint32(C.nox_xxx_readNPCVoiceSet_52AD10(addr))
	case 4:
		C.nox_xxx_monsterOnSpawnSpellcaster_529BC0(addr)
		return 0
	case 5:
		return uint32(C.sub_52BA70(addr))
	case 6:
		return uint32(C.sub_52BAF0(addr))
	default:
		panic("unknown creature xfer helper")
	}
}

// Supply actual linked definitions for the production C lookup, restoring the
// previous owner rather than replacing lookup behavior.
func PortTestCreatureXferDefinitions(head unsafe.Pointer) func() {
	old := C.nox_monsterBin_head_2386924
	C.nox_monsterBin_head_2386924 = head
	return func() { C.nox_monsterBin_head_2386924 = old }
}

func PortTestCreatureXferVoiceSets(head unsafe.Pointer) func() {
	old := C.dword_5d4594_588120
	C.dword_5d4594_588120 = C.uint32_t(uintptr(head))
	return func() { C.dword_5d4594_588120 = old }
}
func PortTestCreatureXferLookupOwner() func() {
	old, init := netCodeCacheState, netCodeCacheNeedInit
	netCodeCacheState = netCodeCacheStorage{}
	netCodeCacheNeedInit = 1
	return func() { netCodeCacheState, netCodeCacheNeedInit = old, init }
}
