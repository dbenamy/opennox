//go:build porttest && linux

package alloc

/*
#include <stdlib.h>
*/
import "C"
import "unsafe"

// RawAlignedAlloc supplies untracked aligned fixture storage. Release with RawFree.
func RawAlignedAlloc(alignment, size uintptr) unsafe.Pointer {
	return C.aligned_alloc(C.size_t(alignment), C.size_t(size))
}
