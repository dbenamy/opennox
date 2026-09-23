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

// Primitive entry dispatch only; the production C callbacks perform the work.
// Object-valued returns are represented by fixture identities 1 and 2.
func PortTestWorldCollision(op int, a, b *server.Object, normal *types.Pointf) uint32 {
	ai, bi := C.int(uintptr(unsafe.Pointer(a))), C.int(uintptr(unsafe.Pointer(b)))
	p := (*C.float)(unsafe.Pointer(normal))
	var rv uint32
	switch op {
	case 0:
		worldCollideMass(a, b)
	case 1:
		C.nox_xxx_collideDoor_4E8AC0(ai, bi)
	case 2:
		rv = uint32(C.nox_xxx_collidePickup_4E8DF0(ai, bi))
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
		C.nox_xxx_collideExit_4E9090(ai, bi, C.int(uintptr(unsafe.Pointer(normal))))
	case 8:
		C.nox_xxx_spellFlyCollide_4E9500(ai, bi, p)
	case 9:
		C.nox_xxx_collideChest_4E9C40((*C.uint32_t)(unsafe.Pointer(a)), bi)
	case 10:
		C.sub_4EAAA0(ai)
	case 11:
		C.sub_4EAAD0(ai, bi)
	case 12:
		rv = uint32(C.nox_xxx_collidePentagram_4EAB20(ai))
	case 13:
		C.nox_xxx_collideSign_4EAB40(ai, bi)
	case 14:
		C.nox_xxx_collideTrapDoor_4EAB60(ai, bi)
	case 15:
		C.sub_4EACA0(ai, bi)
	case 16:
		rv = uint32(C.nox_xxx_collideSpellPedestal_4EAD20(ai, bi))
	case 17:
		C.nox_xxx_collideUndeadKiller_4EBD40(ai, bi, C.int(uintptr(unsafe.Pointer(normal))))
	case 18:
		C.nox_xxx_collideMonsterGen_4EBE10(ai, bi)
	case 19:
		C.sub_4EBE40(ai, bi)
	case 20:
		C.nox_xxx_collideAnkhQuest_4EBF40(ai, bi)
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

func PortTestWorldInversionCallback() unsafe.Pointer { return C.nox_xxx_inversionEffect_4E03D0 }
