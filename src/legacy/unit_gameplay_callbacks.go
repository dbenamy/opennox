package legacy

/*
#include "GAME4_3.h"
*/
import "C"

import (
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

func unitUndeadUpdate(u *server.Object) {
	target := *(**server.Object)(u.CollideData)
	if target != nil && *(*byte)(unsafe.Add(target.CObj(), 88))&1 != 0 || GetServer().S().Frame()-u.Field34 > 70 {
		GetServer().DelayedDelete(u)
	}
}
func unitRead(u, it *server.Object, warp bool) bool {
	if u.ObjClass&4 == 0 {
		return true
	}
	stamp := equipmentWord(it.UseData.Ptr, 256)
	core := GetServer().S()
	if *stamp != 0 && core.Frame()-*stamp <= 3*uint32(core.TickRate()) {
		return true
	}
	if !core.MapTraceVision(u, it) {
		return true
	}
	if !warp {
		gameplayTextPrivate(u, (*byte)(it.UseData.Ptr), 1)
	} else if questRuntimeWord(1556120) != 0 {
		threshold := uint32(Nox_server_questNextStageThreshold_4D74F0(int(int32(questRuntimeStage()))))
		gameplayTextInformation(int(u.UpdateDataPlayer().Player.PlayerInd), 21, unsafe.Pointer(&threshold))
	} else {
		gameplayTextPrivate(u, alloc.InternCString("GeneralPrint:WarpClosed"), 1)
	}
	*stamp = core.Frame()
	return true
}

//export nox_xxx_updateUndeadKiller_53E190
func nox_xxx_updateUndeadKiller_53E190(a C.int) {
	unitUndeadUpdate((*server.Object)(unsafe.Pointer(uintptr(uint32(a)))))
}

//export nox_xxx_useRead_53F7C0
func nox_xxx_useRead_53F7C0(a, b C.int) C.int {
	return C.int(bool2int(unitRead((*server.Object)(unsafe.Pointer(uintptr(uint32(a)))), (*server.Object)(unsafe.Pointer(uintptr(uint32(b)))), false)))
}

//export sub_53F830
func sub_53F830(a, b C.int) C.int {
	return C.int(bool2int(unitRead((*server.Object)(unsafe.Pointer(uintptr(uint32(a)))), (*server.Object)(unsafe.Pointer(uintptr(uint32(b)))), true)))
}
