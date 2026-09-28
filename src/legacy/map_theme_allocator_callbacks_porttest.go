//go:build porttest

package legacy

import "unsafe"

var themeTestOnAllocate func(unsafe.Pointer, int)
var themeTestOnRelease func(unsafe.Pointer)

func themeTestAllocated(ptr unsafe.Pointer, size uintptr) {
	if themeTestOnAllocate != nil {
		themeTestOnAllocate(ptr, int(size))
	}
}

func themeTestReleased(ptr unsafe.Pointer) {
	if themeTestOnRelease != nil {
		themeTestOnRelease(ptr)
	}
}
