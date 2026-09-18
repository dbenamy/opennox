//go:build porttest

package legacy

/*
#include "defs.h"
#include "GAME4_2.h"
extern uint32_t dword_5d4594_588120;
*/
import "C"
import (
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

func PortTestCreatureXferAdjust(value, delta uint32) (uint32, uint32) {
	ret := creatureXferAdjust(&value, delta)
	return value, ret
}
func PortTestCreatureXferHelper(op int, u *server.Object, entry unsafe.Pointer, delta uint32) uint32 {
	switch op {
	case 0:
		return uint32(creatureXferAction(u))
	case 1:
		return creatureXferArgument(entry, delta)
	case 2:
		return uint32(creatureXferBuffs(u))
	case 3:
		return creatureXferVoice(u)
	case 4:
		creatureXferDefaults(u)
		return 0
	case 5:
		return uint32(creatureXferEquipment(u))
	case 6:
		return creatureXferPostload(u)
	default:
		panic("unknown creature xfer helper")
	}
}

// Supply actual linked definitions for the production definition lookup, restoring the
// previous owner rather than replacing lookup behavior.
func PortTestCreatureXferDefinitions(head unsafe.Pointer) func() {
	old := monsterDefinitions
	monsterDefinitions = (*server.MonsterDef)(head)
	return func() { monsterDefinitions = old }
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
