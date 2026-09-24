//go:build porttest

package legacy

/*
#include "defs.h"
#include "GAME3_3.h"
#include "GAME3_2.h"

*/
import "C"

import (
	"bytes"
	"fmt"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

// Primitive entry dispatch through native collision owners and retained C helpers.
// Object-valued returns are represented by fixture identities 1 and 2.
func PortTestWorldCollision(op int, a, b *server.Object, normal *types.Pointf) uint32 {
	ai, bi := C.int(uintptr(unsafe.Pointer(a))), C.int(uintptr(unsafe.Pointer(b)))
	var rv uint32
	switch op {
	case 0:
		worldCollideMass(a, b)
	case 1:
		rv = server.PortTestCollisionResult(collisionKey(collisionIdentityDoor), a, b, normal)
	case 2:
		rv = server.PortTestCollisionResult(collisionKey(collisionIdentityPickup), a, b, normal)
	case 3:
		if unsafe.Pointer(C.sub_4E8E50()) != memmap.PtrOff(0x5D4594, 1567844) {
			panic("pending map buffer")
		}
	case 4:
		rv = uint32(C.sub_4E8E60())
	case 5:
		if bool(C.nox_server_questMaybeWarp_4E8F60()) {
			rv = 1
		}
	case 6:
		rv = uint32(C.sub_4E9010())
	case 7:
		rv = server.PortTestCollisionResult(collisionKey(collisionIdentityExit), a, b, normal)
	case 8:
		rv = server.PortTestCollisionResult(collisionKey(collisionIdentitySpellProjectile), a, b, normal)
	case 9:
		rv = server.PortTestCollisionResult(collisionKey(collisionIdentityChest), a, b, normal)
	case 10:
		rv = server.PortTestCollisionResult(collisionKey(collisionIdentityBarrel), a, b, normal)
	case 11:
		rv = server.PortTestCollisionResult(collisionKey(collisionIdentityAudioEvent), a, b, normal)
	case 12:
		rv = server.PortTestCollisionResult(collisionKey(collisionIdentityPentagram), a, b, normal)
	case 13:
		rv = server.PortTestCollisionResult(collisionKey(collisionIdentitySign), a, b, normal)
	case 14:
		rv = server.PortTestCollisionResult(collisionKey(collisionIdentityTrapDoor), a, b, normal)
	case 15:
		rv = server.PortTestCollisionResult(collisionKey(collisionIdentityTeleport), a, b, normal)
	case 16:
		rv = server.PortTestCollisionResult(collisionKey(collisionIdentityAwardSpell), a, b, normal)
	case 17:
		rv = server.PortTestCollisionResult(collisionKey(collisionIdentityUndeadKiller), a, b, normal)
	case 18:
		rv = server.PortTestCollisionResult(collisionKey(collisionIdentityMonsterGenerator), a, b, normal)
	case 19:
		rv = server.PortTestCollisionResult(collisionKey(collisionIdentitySoulGate), a, b, normal)
	case 20:
		rv = server.PortTestCollisionResult(collisionKey(collisionIdentityAnkh), a, b, normal)
	default:
		panic("world collision operation")
	}
	if a != nil && rv == uint32(ai) {
		return 1
	}
	if b != nil && rv == uint32(bi) {
		return 2
	}
	return rv
}
func PortTestWorldCollisionGlobals() (map[string]*uint32, *uint64, func()) {
	words := map[string]*uint32{"glyph": (*uint32)(unsafe.Pointer(&dword_5d4594_1567960)), "extensions": (*uint32)(unsafe.Pointer(&gameex_flags))}
	words["soulFrame"] = (*uint32)(unsafe.Pointer(&dword_5d4594_1556136))
	words["warpOpen"] = memmap.PtrUint32(0x5D4594, 1556120)
	words["settingsUpdated"] = (*uint32)(unsafe.Pointer(&legacyGlobals.nox_server_gameSettingsUpdated))
	words["savePortal"] = &orchestrationRestoreCleanup
	words["directionX"] = (*uint32)(unsafe.Pointer(&dword_5d4594_1565628))
	words["directionY"] = (*uint32)(unsafe.Pointer(&dword_5d4594_1565632))
	for _, off := range []uintptr{1565652, 1565656, 1565636, 1567708, 1565640} {
		words[fmt.Sprint(off)] = memmap.PtrUint32(0x5D4594, off)
	}
	old := map[string]uint32{}
	for k, p := range words {
		old[k] = *p
		*p = 0
	}
	ticks := (*uint64)(unsafe.Pointer(&qword_5d4594_1567940))
	oldTicks := *ticks
	*ticks = 0
	buf := unsafe.Slice((*byte)(memmap.PtrOff(0x5D4594, 1567844)), 96)
	saved := bytes.Clone(buf)
	clear(buf)
	return words, ticks, func() {
		for k, p := range words {
			*p = old[k]
		}
		*ticks = oldTicks
		copy(buf, saved)
	}
}

func PortTestWorldInversionCallback() unsafe.Pointer { return modifierKey(modifierIDInversionEffect) }
