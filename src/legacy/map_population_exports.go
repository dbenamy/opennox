package legacy

/*
#include <stdint.h>
#include "defs.h"
*/
import "C"
import "unsafe"

//export nox_xxx_mapGenSpellIdByName_51E1D0
func nox_xxx_mapGenSpellIdByName_51E1D0(a0 *C.char) C.int {
	return C.int(mapPopulationSpellID(mapRoomRaw(unsafe.Pointer(a0))))
}

//export nox_xxx_mapgen_522340
func nox_xxx_mapgen_522340(a0 C.int, a1 C.int) { mapPopulationRoom(uint32(a0), uint32(a1)) }

//export nox_xxx_mapGenFinishPopulate_5228B0_mapgen_populate
func nox_xxx_mapGenFinishPopulate_5228B0_mapgen_populate(a0 C.int) { mapPopulationFinish(uint32(a0)) }

//export sub_522C80
func sub_522C80(a0 *C.float) *C.float {
	return (*C.float)(mapRoomPointer(mapPopulationWaypoint(mapRoomRaw(unsafe.Pointer(a0)))))
}

//export sub_522D30
func sub_522D30(a0 C.int) *C.float {
	return (*C.float)(mapRoomPointer(mapPopulationHallwayWaypoints(uint32(a0))))
}

//export nox_xxx_mapGenTryNextRoom_522F40
func nox_xxx_mapGenTryNextRoom_522F40(a0 *C.uint32_t) *C.uchar {
	return (*C.uchar)(mapRoomPointer(mapPopulationRoomWaypoints(mapRoomRaw(unsafe.Pointer(a0)))))
}

//export nox_xxx_mapGenSetFlags_5235F0
func nox_xxx_mapGenSetFlags_5235F0(a0 C.char) { mapPopulationProgress(byte(a0)) }

//export sub_5259F0
func sub_5259F0(a0 C.int, a1 C.int, a2 C.float) {
	mapPopulationDistance(uint32(a0), uint32(a1), float32(a2))
}

//export sub_525AF0
func sub_525AF0(a0 C.int) *C.float {
	return (*C.float)(mapRoomPointer(mapPopulationThemes(uint32(a0))))
}

//export nox_xxx_mapGen_InPrefab1_525D20
func nox_xxx_mapGen_InPrefab1_525D20(a0 C.int) C.int {
	return C.int(mapPopulationSelectPrefabs(uint32(a0)))
}

//export nox_xxx_mapGen_InPrefab2_5266F0
func nox_xxx_mapGen_InPrefab2_5266F0(a0 C.int) C.int {
	return C.int(mapPopulationConnectPrefabs(uint32(a0)))
}

//export nox_xxx_mapGenPlacePrefabs_526830
func nox_xxx_mapGenPlacePrefabs_526830(a0 C.int) C.int {
	return C.int(mapPopulationApplyPrefabs(uint32(a0)))
}

//export sub_526950
func sub_526950() *C.char { return (*C.char)(mapRoomPointer(mapPopulationMetadataInit())) }
