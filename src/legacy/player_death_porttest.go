//go:build porttest

package legacy

/*
#include "GAME4_3.h"
#include "GAME5.h"
*/
import "C"

import (
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

// The actual callers ignore these historical return values. Capture owned state,
// messages, created objects and RNG consumption instead of pointer-shaped results.
func PortTestPlayerDeath(op string, victim, killer, assist *server.Object, info *server.Player) {
	raw := func(p unsafe.Pointer) C.int { return C.int(uintptr(p)) }
	var pl unsafe.Pointer
	if info != nil {
		pl = info.C()
	}
	switch op {
	case "death":
		C.nox_xxx_diePlayer_54D2B0(raw(unsafe.Pointer(victim)))
	case "arena":
		C.nox_xxx_playerUpdateScore_54D980(raw(unsafe.Pointer(victim)), raw(unsafe.Pointer(killer)), raw(unsafe.Pointer(assist)), raw(pl))
	case "elimination":
		C.nox_xxx_playerHandleElimDeath_54D7A0(raw(unsafe.Pointer(victim)), raw(unsafe.Pointer(killer)))
	case "kotr":
		C.nox_xxx_playerHandleKotrDeath_54DC40(raw(unsafe.Pointer(victim)), raw(unsafe.Pointer(killer)))
	case "notify":
		C.nox_xxx_netNotifyPlayerDied_54DF00(raw(unsafe.Pointer(victim)))
	default:
		panic(op)
	}
}
func PortTestPlayerCorpseCache() { C.nox_xxx_createCorpse_53FCA0() }
func PortTestPlayerCorpseCreate(pos types.Pointf, angle int32) {
	C.nox_xxx_respawnPlayerImpl_53FBC0((*C.float)(unsafe.Pointer(&pos)), C.int(angle))
}

// Own the real lookup cache so assist resolution cannot retain fixture objects.
func PortTestPlayerDeathLookupOwner() func() {
	oldState, oldInit := netCodeCacheState, netCodeCacheNeedInit
	netCodeCacheState, netCodeCacheNeedInit = netCodeCacheStorage{}, 1
	return func() { netCodeCacheState, netCodeCacheNeedInit = oldState, oldInit }
}
