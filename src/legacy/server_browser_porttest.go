//go:build porttest

package legacy

/*
#include "GAME3.h"
#include "client__shell__noxworld.h"
*/
import "C"
import (
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"unsafe"
)

func PortTestServerBrowserMode(mode uint16) string {
	return alloc.GoString16((*uint16)(unsafe.Pointer(nox_gui_wol_gameModeString_43BCB0(C.short(mode)))))
}
func PortTestServerBrowserHit(point, record unsafe.Pointer) int {
	return int(sub_4A2560((*C.uint32_t)(point), C.int(uintptr(record))))
}
func PortTestServerBrowserClamp(x, y int32, out unsafe.Pointer) uintptr {
	return uintptr(unsafe.Pointer(sub_4A2830(C.int(x), C.int(y), (*C.uint32_t)(out))))
}
func PortTestServerBrowserCount(point, head unsafe.Pointer) int {
	return int(sub_4A25C0((*C.uint32_t)(point), (*C.int)(head)))
}
func PortTestServerBrowserList(items []unsafe.Pointer) (unsafe.Pointer, func()) {
	head, free := alloc.New(legacyListNode{})
	listClear(head)
	var old []legacyListNode
	for _, p := range items {
		n := (*legacyListNode)(p)
		old = append(old, *n)
		listInit(n)
		listAppend(head, n)
	}
	return unsafe.Pointer(head), func() {
		for i, p := range items {
			*(*legacyListNode)(p) = old[i]
		}
		free()
	}
}

func PortTestServerBrowserPopup(parent, point, head unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(uintptr(uint32(sub_4A2610(int32(uintptr(parent)), (*uint32)(point), (*int32)(head)))))
}
func PortTestServerBrowserPopupShown() bool { return sub_4A28B0() != 0 }
func PortTestServerBrowserPopupAt(i int) unsafe.Pointer {
	return unsafe.Pointer(uintptr(uint32(sub_4A28C0(int32(i)))))
}
func PortTestServerBrowserPopupClose() { sub_4A2890() }
