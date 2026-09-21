package legacy

/*
#include <stdlib.h>
*/
import "C"
import "unsafe"

func runtimeRejectedClear(head unsafe.Pointer) uintptr {
	first := listNext((*legacyListNode)(head))
	result := uintptr(unsafe.Pointer(first))
	for p := first; p != nil; {
		next := listNext(p)
		listRemove(p)
		C.free(unsafe.Pointer(p))
		p = next
	}
	return result
}
