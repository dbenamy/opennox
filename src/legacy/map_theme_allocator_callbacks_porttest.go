//go:build porttest

package legacy

/*
#include <stddef.h>
*/
import "C"
import "unsafe"

var themeTestOnAllocate func(unsafe.Pointer, int)
var themeTestOnRelease func(unsafe.Pointer)

//export themeTestAllocated
func themeTestAllocated(ptr unsafe.Pointer, size C.size_t) {
	if themeTestOnAllocate != nil {
		themeTestOnAllocate(ptr, int(size))
	}
}

//export themeTestReleased
func themeTestReleased(ptr unsafe.Pointer) {
	if themeTestOnRelease != nil {
		themeTestOnRelease(ptr)
	}
}
