//go:build porttest && linux

package alloc

import "unsafe"

// RawAlignedAlloc supplies untracked aligned fixture storage. Release with RawFree.
func RawAlignedAlloc(alignment, size uintptr) unsafe.Pointer {
	return rawAllocate(size, alignment)
}
