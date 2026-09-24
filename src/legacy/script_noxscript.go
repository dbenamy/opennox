package legacy

/*
#include "server__script__script.h"
#include "server__script__internal.h"
int nox_xxx_gameIsSwitchToSolo_4DB240();
*/
import "C"
import (
	"unsafe"

	"github.com/opennox/opennox/v1/internal/binfile"
	"github.com/opennox/opennox/v1/server"
	"github.com/opennox/opennox/v1/server/noxscript"
)

var (
	Nox_setImaginaryCaster         func() int
	Nox_script_readWriteZzz_541670 func(cpath, cpath2, cdst *byte) int
)

type NoxScript interface {
	noxscript.VM
	ActResolveObjs()
	ScriptToObject(h int) *server.Object
	ScriptCallback(b *server.ScriptCallback, caller, trigger *server.Object, eventCode server.ScriptEventType) unsafe.Pointer
}

func nox_xxx_netGetUnitCodeServ_578AC0(cobj *nox_object_t) C.uint {
	return C.uint(GetServer().S().GetUnitNetCode(asObjectS(cobj)))
}

func nox_xxx_scriptCallByEventBlock_502490(a1 unsafe.Pointer, a2, a3 unsafe.Pointer, eventCode int32) unsafe.Pointer {
	return GetServer().NoxScriptC().ScriptCallback((*server.ScriptCallback)(a1), AsObjectP(a2), AsObjectP(a3), server.ScriptEventType(eventCode))
}

func Sub_516570() {
	monsterControlChapter()
}
func Nox_xxx_script_511C50(a1 int) *server.Object {
	return monsterCacheFind(int32(a1))
}
func Nox_xxx_scriptPrepareFoundUnit_511D70(a1 *server.Object) {
	monsterCachePrepare(a1)
}
func Nox_script_readWriteWww_5417C0(a1, a2, a3 *binfile.File) {
	prefabScriptMerge(a1, a2, a3)
}
