//go:build porttest

package legacy

import (
	"runtime"
	"unsafe"
)

// Sprite data uses raw libc allocations rather than the tracked alloc package.
func PortTestResourceAllocate(size uintptr) unsafe.Pointer { return legacyCalloc(1, uintptr(size)) }
func PortTestResourceRelease(ptr unsafe.Pointer)           { legacyFree(ptr) }

// Observe engine frees on the current thread, recording raw addresses before release.
// Normalize recorded addresses after stopping observation; the addresses are never
// dereferenced after release. Existing theme/grid observers keep their own state.
func PortTestObserveResourceFrees(pointers []unsafe.Pointer, free func()) []uint32 {
	ids := make(map[uintptr]uint32, len(pointers))
	for i, p := range pointers {
		if p == nil || ids[uintptr(p)] != 0 {
			panic("invalid resource allocation registry")
		}
		ids[uintptr(p)] = uint32(i + 1)
	}
	capacity := len(pointers) + 16
	events := legacyCalloc(uintptr(capacity), uintptr(unsafe.Sizeof(uintptr(0))))
	if events == nil {
		panic("resource observer allocation")
	}
	defer legacyFree(events)
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	allocationTestResourceStart(events, capacity)
	stopped := false
	defer func() {
		if !stopped {
			allocationTestResourceStop()
		}
	}()
	free()
	count := int(allocationTestResourceStop())
	stopped = true
	if count > capacity {
		panic("resource free observer overflow")
	}
	order := make([]uint32, 0, len(pointers))
	for _, p := range unsafe.Slice((*uintptr)(events), count) {
		if id := ids[p]; id != 0 {
			order = append(order, id)
		}
	}
	return order
}
