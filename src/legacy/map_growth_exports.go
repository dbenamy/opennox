package legacy

/*
#include <stdint.h>
extern uint32_t dword_5d4594_1550912;
*/
import "C"
import "unsafe"

func mapGrowthFrontier() *uint32 { return (*uint32)(unsafe.Pointer(&C.dword_5d4594_1550912)) }

//export nox_xxx_mapgen_Doors_4D4790
func nox_xxx_mapgen_Doors_4D4790() *C.float { return (*C.float)(unsafe.Pointer(mapGrowthDoors())) }

//export nox_xxx_mapGenMkSmallRoom_4D4F40
func nox_xxx_mapGenMkSmallRoom_4D4F40(cfg *C.uint32_t) *C.float {
	return (*C.float)(unsafe.Pointer(mapGrowthInitial(mapRoomRaw(unsafe.Pointer(cfg)))))
}

//export sub_4D52F0
func sub_4D52F0() { mapGrowthFrontiers() }
