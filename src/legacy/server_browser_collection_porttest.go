//go:build porttest

package legacy

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
		browserListClear(false)
		if browserSelectedSnapshot != nil {
			alloc.FreePtr(browserSelectedSnapshot)
		}
		browserSelectedSnapshot = oldSnapshot
		*head = saved
		*init = oldInit
	}
}
func PortTestServerBrowserCollectionClear() { browserListClear(false) }
func PortTestServerBrowserCollectionAdd(p unsafe.Pointer) int {
	return browserListAdd(p)
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
	return browserListAt(index)
}
func PortTestServerBrowserCollectionID(id int32) unsafe.Pointer {
	return browserListID(id)
}
func PortTestServerBrowserCollectionMissing(addr string, port uint16) int {
	p := CString(addr)
	defer StrFree(p)
	return browserListMissing(GoString(p), int16(port))
}

// Native re-sort retains the existing records and owns its selection snapshot.
func PortTestServerBrowserCollectionResort() func() {
	browserListResort()
	return func() {}
}
