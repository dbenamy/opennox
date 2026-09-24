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
*/
import "C"
import (
	"github.com/opennox/opennox/v1/internal/cryptfile"
	"unsafe"

	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

// Exercise the Go adapter and its separate recognized/error results.
func PortTestPrefabReadSection(bounds unsafe.Pointer, name string, initial uint32) (int, uint32) {
	s, free := alloc.CString(name)
	defer free()
	err := C.uint(initial)
	ok := portTestInvoke_nox_xxx_mapReadSection_426EA0(bounds, (*C.char)(unsafe.Pointer(s)), &err)
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
func PortTestPrefabClearSecrets()                         { worldSecretClear() }
func PortTestPrefabSecretList() (*uint32, func()) {
	p := paintGlobals()["secretWalls"]
	old := *p
	*p = 0
	return p, func() { *p = old }
}

// Fixture-native copies preserve the original wrapper ABI conversions.
func portTestInvoke_nox_xxx_mapReadSection_426EA0(a1 unsafe.Pointer, cname *C.char, cerr *C.uint) int {
	ok, err := Nox_xxx_mapReadSection(cryptfile.Global(), a1, GoString(cname))
	*cerr = C.uint(bool2int(err != nil))
	if err != nil {
		mapLog.Println(err)
	}
	return bool2int(ok)
}
