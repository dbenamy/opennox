//go:build porttest

package legacy

/*
#include "defs.h"
#include "GAME2_3.h"
*/
import "C"
import (
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"unsafe"
)

func PortTestServerBrowserCollectionOwner() func() {
	head := (*legacyListNode)(unsafe.Pointer(browserListHead()))
	saved := *head
	oldSnapshot := browserSelectedSnapshot
	browserSelectedSnapshot = nil
	init := memmap.PtrUint32(0x5D4594, 1305808)
	oldInit := *init
	listClear(head)
	*init = 1
	return func() {
		sub_49FFA0(0)
		if browserSelectedSnapshot != nil {
			alloc.FreePtr(browserSelectedSnapshot)
		}
		browserSelectedSnapshot = oldSnapshot
		*head = saved
		*init = oldInit
	}
}
func PortTestServerBrowserCollectionClear() { sub_49FFA0(0) }
func PortTestServerBrowserCollectionAdd(p unsafe.Pointer) int {
	return int(nox_wol_servers_addResult_4A0030((*C.nox_gui_server_ent_t)(p)))
}
func PortTestServerBrowserCollectionSnapshot() []unsafe.Pointer {
	var out []unsafe.Pointer
	head := (*legacyListNode)(unsafe.Pointer(browserListHead()))
	for n := listNext(head); n != nil; n = listNext(n) {
		out = append(out, unsafe.Pointer(n))
	}
	return out
}
func PortTestServerBrowserCollectionAt(index int32) unsafe.Pointer {
	return unsafe.Pointer(sub_4A04C0(C.int(index)))
}
func PortTestServerBrowserCollectionID(id int32) unsafe.Pointer {
	return unsafe.Pointer(sub_4A0490(C.int(id)))
}
func PortTestServerBrowserCollectionMissing(addr string, port uint16) int {
	p := CString(addr)
	defer StrFree(p)
	return int(sub_4A0410(p, C.short(port)))
}

// Native re-sort retains the existing records and owns its selection snapshot.
func PortTestServerBrowserCollectionResort() func() {
	sub_4A0390()
	return func() {}
}
