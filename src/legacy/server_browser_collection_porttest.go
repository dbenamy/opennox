//go:build porttest

package legacy

/*
#include "defs.h"
#include "GAME2_3.h"
extern nox_list_item_t nox_gui_wol_servers_list;
*/
import "C"
import (
	"github.com/opennox/opennox/v1/common/memmap"
	"unsafe"
)

func PortTestServerBrowserCollectionOwner() func() {
	head := (*legacyListNode)(unsafe.Pointer(&C.nox_gui_wol_servers_list))
	saved := *head
	init := memmap.PtrUint32(0x5D4594, 1305808)
	oldInit := *init
	listClear(head)
	*init = 1
	return func() { C.sub_49FFA0(0); *head = saved; *init = oldInit }
}
func PortTestServerBrowserCollectionClear() { C.sub_49FFA0(0) }
func PortTestServerBrowserCollectionAdd(p unsafe.Pointer) int {
	return int(C.nox_wol_servers_addResult_4A0030((*C.nox_gui_server_ent_t)(p)))
}
func PortTestServerBrowserCollectionSnapshot() []unsafe.Pointer {
	var out []unsafe.Pointer
	head := (*legacyListNode)(unsafe.Pointer(&C.nox_gui_wol_servers_list))
	for n := listNext(head); n != nil; n = listNext(n) {
		out = append(out, unsafe.Pointer(n))
	}
	return out
}
func PortTestServerBrowserCollectionAt(index int32) unsafe.Pointer {
	return unsafe.Pointer(C.sub_4A04C0(C.int(index)))
}
func PortTestServerBrowserCollectionID(id int32) unsafe.Pointer {
	return unsafe.Pointer(C.sub_4A0490(C.int(id)))
}
func PortTestServerBrowserCollectionMissing(addr string, port uint16) int {
	p := CString(addr)
	defer StrFree(p)
	return int(C.sub_4A0410(p, C.short(port)))
}

// Original C detaches and leaks the old list. The fixture owns those allocations
// and releases them only after checking the selected record's lifetime.
func PortTestServerBrowserCollectionResort() func() {
	old := PortTestServerBrowserCollectionSnapshot()
	C.sub_4A0390()
	return func() {
		for _, p := range old {
			StrFree((*byte)(p))
		}
	}
}
