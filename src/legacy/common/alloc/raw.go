package alloc

/*
#include <stdlib.h>
*/
import "C"
import "unsafe"

// rawCalloc allocates zeroed libc memory without registering it in allocs.
func rawCalloc(num, size uintptr) unsafe.Pointer {
	return C.calloc(C.size_t(num), C.size_t(size))
}

// RawRealloc reallocates libc memory without changing allocs.
func RawRealloc(ptr unsafe.Pointer, size uintptr) unsafe.Pointer {
	return C.realloc(ptr, C.size_t(size))
}

// rawFree releases libc memory without consulting or changing allocs.
func rawFree(ptr unsafe.Pointer) { C.free(ptr) }

// RawMalloc allocates untracked libc memory. Like cgo CString's allocator,
// cgo's special malloc wrapper terminates the process on allocation failure.
func RawMalloc(size uintptr) unsafe.Pointer { return C.malloc(C.size_t(size)) }
