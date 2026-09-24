package alloc

/*
#include <stdlib.h>
*/
import "C"
import "unsafe"

// RawCalloc allocates zeroed libc memory without registering it in allocs.
func RawCalloc(num, size uintptr) unsafe.Pointer {
	return C.calloc(C.size_t(num), C.size_t(size))
}

// RawRealloc reallocates libc memory without changing allocs.
func RawRealloc(ptr unsafe.Pointer, size uintptr) unsafe.Pointer {
	return C.realloc(ptr, C.size_t(size))
}

// RawFree releases libc memory without consulting or changing allocs.
func RawFree(ptr unsafe.Pointer) { C.free(ptr) }
