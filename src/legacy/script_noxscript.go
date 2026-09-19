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

//export nox_script_activatorResolveObjs_51B0C0
func nox_script_activatorResolveObjs_51B0C0() {
	GetServer().NoxScriptC().ActResolveObjs()
}

//export nox_xxx_netGetUnitCodeServ_578AC0
func nox_xxx_netGetUnitCodeServ_578AC0(cobj *nox_object_t) C.uint {
	return C.uint(GetServer().S().GetUnitNetCode(asObjectS(cobj)))
}

//export nox_setImaginaryCaster
func nox_setImaginaryCaster() int { return Nox_setImaginaryCaster() }

//export nox_script_readWriteZzz_541670
func nox_script_readWriteZzz_541670(cpath, cpath2, cdst *C.char) int {
	return Nox_script_readWriteZzz_541670((*byte)(unsafe.Pointer(cpath)), (*byte)(unsafe.Pointer(cpath2)), (*byte)(unsafe.Pointer(cdst)))
}

//export nox_xxx_scriptCallByEventBlock_502490
func nox_xxx_scriptCallByEventBlock_502490(a1 unsafe.Pointer, a2, a3 unsafe.Pointer, eventCode int32) unsafe.Pointer {
	return GetServer().NoxScriptC().ScriptCallback((*server.ScriptCallback)(a1), AsObjectP(a2), AsObjectP(a3), server.ScriptEventType(eventCode))
}

//export nox_script_callByIndex_507310
func nox_script_callByIndex_507310(index int, a2 unsafe.Pointer, a3 unsafe.Pointer) {
	if err := GetServer().S().NoxScriptVM.CallByIndex(index, AsObjectP(a2), AsObjectP(a3)); err != nil {
		scriptLog.Println(err)
	}
}

//export nox_script_objCallbackName_508CB0
func nox_script_objCallbackName_508CB0(obj *nox_object_t, event int) *C.char {
	s, ok := GetServer().S().NoxScriptVM.Nox_script_objCallbackName_508CB0(asObjectS(obj), event)
	if !ok {
		return nil
	}
	return internCStr(s)
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
