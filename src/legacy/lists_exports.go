package legacy

/*
#include "defs.h"
*/
import "C"
import "unsafe"

//export nox_common_list_clear_425760
func nox_common_list_clear_425760(p *C.nox_list_item_t) {
	listClear((*legacyListNode)(unsafe.Pointer(p)))
}

//export sub_425770
func sub_425770(p unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(listInit((*legacyListNode)(p)))
}

//export sub_425790
func sub_425790(p *C.int, q *C.uint) C.int {
	return C.int(listAscending((*legacyListNode)(unsafe.Pointer(p)), (*legacyListNode)(unsafe.Pointer(q))))
}

//export sub_4257F0
func sub_4257F0(p *C.int, q *C.uint) {
	listDescending((*legacyListNode)(unsafe.Pointer(p)), (*legacyListNode)(unsafe.Pointer(q)))
}

//export nox_common_list_getFirstSafe_425890
func nox_common_list_getFirstSafe_425890(p *C.nox_list_item_t) *C.nox_list_item_t {
	return (*C.nox_list_item_t)(unsafe.Pointer(listNext((*legacyListNode)(unsafe.Pointer(p)))))
}

//export nox_common_list_getNextSafe_4258A0
func nox_common_list_getNextSafe_4258A0(p *C.nox_list_item_t) *C.nox_list_item_t {
	return (*C.nox_list_item_t)(unsafe.Pointer(listNext((*legacyListNode)(unsafe.Pointer(p)))))
}

//export sub_4258C0
func sub_4258C0(p **C.uint, index C.int) *C.uint {
	return (*C.uint)(unsafe.Pointer(listAt((*legacyListNode)(unsafe.Pointer(p)), int(index))))
}

//export nox_common_list_append_4258E0
func nox_common_list_append_4258E0(p, q *C.nox_list_item_t) {
	listAppend((*legacyListNode)(unsafe.Pointer(p)), (*legacyListNode)(unsafe.Pointer(q)))
}

//export sub_425900
func sub_425900(p, q *C.uint) *C.uint {
	return (*C.uint)(unsafe.Pointer(listPrepend((*legacyListNode)(unsafe.Pointer(p)), (*legacyListNode)(unsafe.Pointer(q)))))
}

//export nox_common_list_remove_425920
func nox_common_list_remove_425920(p unsafe.Pointer) { listRemove((*legacyListNode)(p)) }

//export nox_common_list_getNext_425940
func nox_common_list_getNext_425940(p *C.nox_list_item_t) *C.nox_list_item_t {
	return (*C.nox_list_item_t)(unsafe.Pointer(listNext((*legacyListNode)(unsafe.Pointer(p)))))
}

//export sub_425960
func sub_425960(p C.int) C.int {
	return C.int(uintptr(unsafe.Pointer(listPrev((*legacyListNode)(unsafe.Pointer(uintptr(p)))))))
}

//export sub_425A50
func sub_425A50() *C.int { return (*C.int)(unsafe.Pointer(playerGroupsFirst())) }

//export sub_425A60
func sub_425A60(p *C.int) *C.int {
	return (*C.int)(unsafe.Pointer(playerGroupsNext((*playerGroup)(unsafe.Pointer(p)))))
}

//export sub_425A70
func sub_425A70(id C.int) *C.int { return (*C.int)(unsafe.Pointer(playerGroupFind(uint32(id)))) }

//export sub_425B60
func sub_425B60(p unsafe.Pointer, index C.int) *C.char {
	return (*C.char)(playerGroupRemoveMember((*playerGroup)(p), int32(index)))
}
