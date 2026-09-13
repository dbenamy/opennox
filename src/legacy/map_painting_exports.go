package legacy

/*
#include <stdint.h>
#include "defs.h"
*/
import "C"
import (
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

//export nox_xxx_tileListAddNewSubtile_422160
func nox_xxx_tileListAddNewSubtile_422160(a1 C.int, a2 C.int, a3 C.int, a4 C.int) *C.int {
	return (*C.int)(unsafe.Pointer(uintptr(mapRoomRaw(unsafe.Pointer(mapPaintSubtileNew(int32(uint32(a1)), int32(uint32(a2)), int32(uint32(a3)), int32(uint32(a4))))))))
}

//export nox_xxx_tileFreeTile_422200
func nox_xxx_tileFreeTile_422200(a1 C.int) C.int {
	return C.int(mapPaintSubtileClear(mapPaintNode(uint32(a1))))
}

//export nox_xxx_wall_42A6C0
func nox_xxx_wall_42A6C0(a1 C.uchar, a2 C.uchar) C.uchar {
	return C.uchar(mapPaintWallCompose(byte(uint32(a1)), byte(uint32(a2))))
}

//export nox_xxx_mapGenFixCoords_4D3D90
func nox_xxx_mapGenFixCoords_4D3D90(a1 *C.float2, a2 *C.float2) C.int {
	return C.int(mapPaintTransform((*types.Pointf)(mapRoomPointer(uint32(uintptr(unsafe.Pointer(a1))))), (*types.Pointf)(mapRoomPointer(uint32(uintptr(unsafe.Pointer(a2)))))))
}

//export sub_51D8F0
func sub_51D8F0(a1 *C.float2) C.int {
	return C.int(mapPaintWorldFloor((*types.Pointf)(mapRoomPointer(uint32(uintptr(unsafe.Pointer(a1)))))))
}

//export sub_5245A0
func sub_5245A0(a1 C.int, a2 *C.float, a3 C.int, a4 C.int) *C.float {
	return (*C.float)(unsafe.Pointer(uintptr(mapPaintRect(mapRoomPointer(uint32(a1)), (*types.Pointf)(mapRoomPointer(uint32(uintptr(unsafe.Pointer(a2))))), int32(uint32(a3)), int32(uint32(a4))))))
}

//export nox_xxx_gen_524E00
func nox_xxx_gen_524E00(a1 C.int, a2 C.int) {
	mapPaintRoomWalls(mapRoomPointer(uint32(a1)), (*mapRoom)(mapRoomPointer(uint32(a2))))
}

//export sub_526C40
func sub_526C40(a1 C.int) C.int { return C.int(mapPaintAfterWalls(int32(uint32(a1)))) }

//export sub_527030
func sub_527030(a1 *C.float2) C.int {
	return C.int(mapPaintEraseWall((*types.Pointf)(mapRoomPointer(uint32(uintptr(unsafe.Pointer(a1)))))))
}

//export nox_xxx_mapGenGetObjID_527940
func nox_xxx_mapGenGetObjID_527940(a1 *C.char) C.int {
	return C.int(mapPaintSelectObject((*C.char)(mapRoomPointer(uint32(uintptr(unsafe.Pointer(a1)))))))
}

//export nox_xxx_mapGenPlaceObj_5279B0
func nox_xxx_mapGenPlaceObj_5279B0(a1 *C.float2) *C.float {
	return (*C.float)(unsafe.Pointer(uintptr(mapRoomRaw(unsafe.Pointer(mapPaintPlaceObject((*types.Pointf)(mapRoomPointer(uint32(uintptr(unsafe.Pointer(a1)))))))))))
}

//export nox_xxx_mapGenMoveObject_527A10
func nox_xxx_mapGenMoveObject_527A10(a1 *C.float, a2 *C.float2) *C.float {
	return (*C.float)(unsafe.Pointer(uintptr(mapRoomRaw(unsafe.Pointer(mapPaintMoveObject((*server.Object)(mapRoomPointer(uint32(uintptr(unsafe.Pointer(a1))))), (*types.Pointf)(mapRoomPointer(uint32(uintptr(unsafe.Pointer(a2)))))))))))
}

//export nox_xxx_mapGenOrientObj_527C60
func nox_xxx_mapGenOrientObj_527C60(a1 C.int, a2 C.int) C.int {
	return C.int(mapPaintOrientObject((*server.Object)(mapRoomPointer(uint32(a1))), int32(uint32(a2))))
}

//export nox_xxx_mapGenFinishSpellbook_527DB0
func nox_xxx_mapGenFinishSpellbook_527DB0(a1 C.int, a2 C.char) C.int {
	return C.int(mapPaintFinishBook((*server.Object)(mapRoomPointer(uint32(a1))), byte(uint32(a2))))
}

//export nox_xxx_tileSubtile_544310
func nox_xxx_tileSubtile_544310(a1 *C.float2) C.int {
	return C.int(mapPaintWorldBorder((*types.Pointf)(mapRoomPointer(uint32(uintptr(unsafe.Pointer(a1)))))))
}
