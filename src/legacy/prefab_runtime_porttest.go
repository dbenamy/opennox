//go:build porttest

package legacy

/*
#include <stdint.h>
#include <stdlib.h>
#include "GAME1.h"
extern uint32_t* nox_xxx_tileAllocTileInCoordList_5040A0(int, int, float);
extern uint32_t* sub_504290(char, char);
extern uint32_t* sub_5044B0(int, float, float);
extern uint32_t* nox_xxx_unitAddToList_5048A0(int);
extern int nox_xxx_mapReadSection_426EA0(void*, char*, unsigned int*);
*/
import "C"
import (
	"unsafe"

	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

// Exercise the real C entrypoint and its separate recognized/error results.
func PortTestPrefabReadSection(bounds unsafe.Pointer, name string, initial uint32) (int, uint32) {
	s, free := alloc.CString(name)
	defer free()
	err := C.uint(initial)
	ok := C.nox_xxx_mapReadSection_426EA0(bounds, (*C.char)(unsafe.Pointer(s)), &err)
	return int(ok), uint32(err)
}

func PortTestPrefabObjectNode(obj unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(nox_xxx_unitAddToList_5048A0(C.int(uintptr(obj))))
}

func PortTestPrefabLoadedWord() *uint32 { return prefabGlobal(prefabLoaded) }

func PortTestPrefabCacheNode(kind int) unsafe.Pointer {
	switch kind {
	case 0:
		return unsafe.Pointer(nox_xxx_tileAllocTileInCoordList_5040A0(2, 3, 0))
	case 1:
		return unsafe.Pointer(sub_504290(2, 3))
	case 2:
		return unsafe.Pointer(sub_5044B0(17, 2, 3))
	}
	panic("prefab cache kind")
}

// These payloads still have C allocation ownership during the C baseline.
func PortTestPrefabReleasePayload(p unsafe.Pointer) { C.free(p) }

func PortTestPrefabRawAllocation(size int) unsafe.Pointer { return mapRoomCalloc(1, uintptr(size)) }
func PortTestPrefabClearSecrets()                         { C.sub_410730() }
func PortTestPrefabSecretList() (*uint32, func()) {
	p := paintGlobals()["secretWalls"]
	old := *p
	*p = 0
	return p, func() { *p = old }
}
