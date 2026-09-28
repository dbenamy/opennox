//go:build porttest

package legacy

import (
	"github.com/opennox/opennox/v1/internal/cryptfile"
	"unsafe"

	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

// Exercise the Go adapter and its separate recognized/error results.
func PortTestPrefabReadSection(bounds unsafe.Pointer, name string, initial uint32) (int, uint32) {
	s, free := alloc.CString(name)
	defer free()
	err := uint32(initial)
	ok := portTestInvoke_nox_xxx_mapReadSection_426EA0(bounds, (*int8)(unsafe.Pointer(s)), &err)
	return int(ok), uint32(err)
}

func PortTestPrefabObjectNode(obj unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(nox_xxx_unitAddToList_5048A0(int32(uintptr(obj))))
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
func PortTestPrefabReleasePayload(p unsafe.Pointer) { legacyFree(p) }

func PortTestPrefabRawAllocation(size int) unsafe.Pointer { return mapRoomCalloc(1, uintptr(size)) }
func PortTestPrefabClearSecrets()                         { worldSecretClear() }
func PortTestPrefabSecretList() (*uint32, func()) {
	p := paintGlobals()["secretWalls"]
	old := *p
	*p = 0
	return p, func() { *p = old }
}

// Fixture-native copies preserve the original wrapper ABI conversions.
func portTestInvoke_nox_xxx_mapReadSection_426EA0(a1 unsafe.Pointer, cname *int8, cerr *uint32) int {
	ok, err := Nox_xxx_mapReadSection(cryptfile.Global(), a1, GoStringP(unsafe.Pointer(cname)))
	*cerr = uint32(bool2int(err != nil))
	if err != nil {
		mapLog.Println(err)
	}
	return bool2int(ok)
}
