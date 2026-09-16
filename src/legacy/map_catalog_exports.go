package legacy

/*
#include "defs.h"
*/
import "C"
import "unsafe"

//export nox_common_maplist_first_4D09B0
func nox_common_maplist_first_4D09B0() *C.nox_map_list_item {
	return (*C.nox_map_list_item)(unsafe.Pointer(mapCatalogFirst()))
}

//export nox_common_maplist_next_4D09C0
func nox_common_maplist_next_4D09C0(p *C.nox_map_list_item) *C.nox_map_list_item {
	return (*C.nox_map_list_item)(unsafe.Pointer(mapCatalogNext((*Nox_map_list_item)(unsafe.Pointer(p)))))
}

//export sub_4D0D70
func sub_4D0D70() C.int { return C.int(mapCycleEnabled()) }

//export sub_4D0D90
func sub_4D0D90(v C.int) C.int { return C.int(mapCycleSetEnabled(uint32(v))) }

//export nox_xxx_getQuestMapFile_4D0F60
func nox_xxx_getQuestMapFile_4D0F60() *C.char { return (*C.char)(unsafe.Pointer(mapQuestChoose())) }
