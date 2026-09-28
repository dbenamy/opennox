//go:build porttest && linux

package alloc

import (
	"sync"
	"syscall"
	"unsafe"
)

// PortTestAllocationObserver observes the pinned fixture thread's engine calls.
// Callbacks run before tracker registration and after tracker removal, matching
// the original libc observer boundary. Realloc and raw malloc are not observed.
type PortTestAllocationObserver struct {
	Context      any
	BeforeCalloc func(num, size uintptr) bool
	Allocated    func(ptr unsafe.Pointer, size uintptr)
	Released     func(ptr unsafe.Pointer)
}

var portTestAllocationObservers sync.Map

func PortTestCurrentAllocationObserver() *PortTestAllocationObserver {
	value, ok := portTestAllocationObservers.Load(syscall.Gettid())
	if !ok {
		return nil
	}
	return value.(*PortTestAllocationObserver)
}

// PortTestObserveAllocations requires a pinned OS thread. Restore on that same
// thread before unpinning; independent threads may each own one observer.
func PortTestObserveAllocations(observer *PortTestAllocationObserver) func() {
	if observer == nil {
		panic("nil allocation observer")
	}
	thread := syscall.Gettid()
	if _, loaded := portTestAllocationObservers.LoadOrStore(thread, observer); loaded {
		panic("allocation observer already active on thread")
	}
	var once sync.Once
	return func() {
		if syscall.Gettid() != thread {
			panic("allocation observer restored on another thread")
		}
		once.Do(func() {
			if !portTestAllocationObservers.CompareAndDelete(thread, observer) {
				panic("allocation observer ownership changed")
			}
		})
	}
}

func RawCalloc(num, size uintptr) unsafe.Pointer {
	observer := PortTestCurrentAllocationObserver()
	if observer != nil && observer.BeforeCalloc != nil && !observer.BeforeCalloc(num, size) {
		return nil
	}
	ptr := rawCalloc(num, size)
	if ptr != nil && observer != nil && observer.Allocated != nil {
		observer.Allocated(ptr, num*size)
	}
	return ptr
}

func RawFree(ptr unsafe.Pointer) {
	if ptr != nil {
		if observer := PortTestCurrentAllocationObserver(); observer != nil && observer.Released != nil {
			observer.Released(ptr)
		}
	}
	rawFree(ptr)
}
