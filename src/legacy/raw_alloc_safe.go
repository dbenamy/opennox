//go:build safe

package legacy

import (
	"unsafe"

	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func legacyCalloc(num, size uintptr) unsafe.Pointer {
	ptr, _ := alloc.Calloc(int(num), size)
	return ptr
}
func legacyRealloc(ptr unsafe.Pointer, size uintptr) unsafe.Pointer { return alloc.Realloc(ptr, size) }
func legacyFree(ptr unsafe.Pointer)                                 { alloc.FreePtr(ptr) }
