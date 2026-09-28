//go:build porttest && linux

package legacy

import (
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"runtime"
	"unsafe"
)

type allocationTestState struct {
	restore              func()
	theme                bool
	grid                 bool
	fail, count, freed   int
	valid                bool
	pointers             [130]unsafe.Pointer
	sizes                [130]uintptr
	live                 [130]bool
	events               unsafe.Pointer
	capacity, eventCount int
}

func allocationTestCurrent(create bool) *allocationTestState {
	if observer := alloc.PortTestCurrentAllocationObserver(); observer != nil {
		state, ok := observer.Context.(*allocationTestState)
		if !ok {
			panic("unexpected fixture allocation observer")
		}
		return state
	}
	if !create {
		return nil
	}
	state := new(allocationTestState)
	observer := &alloc.PortTestAllocationObserver{Context: state, BeforeCalloc: state.beforeCalloc, Allocated: state.allocated, Released: state.released}
	state.restore = alloc.PortTestObserveAllocations(observer)
	return state
}

func (s *allocationTestState) stopIfIdle() {
	if !s.theme && !s.grid && s.events == nil {
		s.restore()
		s.restore = nil
	}
}

func (s *allocationTestState) beforeCalloc(num, size uintptr) bool {
	if s.grid && s.fail > 0 && num == 128 && (size == 4 || size == 44) {
		s.fail--
		if s.fail == 0 {
			return false
		}
	}
	return true
}

func (s *allocationTestState) allocated(ptr unsafe.Pointer, size uintptr) {
	if s.grid {
		if s.count >= len(s.pointers) {
			s.valid = false
		} else {
			s.pointers[s.count], s.sizes[s.count], s.live[s.count] = ptr, size, true
			s.count++
		}
	}
	if s.theme {
		themeTestAllocated(ptr, size)
	}
}

func (s *allocationTestState) released(ptr unsafe.Pointer) {
	if s.events != nil {
		if s.eventCount < s.capacity {
			unsafe.Slice((*uintptr)(s.events), s.capacity)[s.eventCount] = uintptr(ptr)
		}
		s.eventCount++
	}
	if s.grid {
		found := false
		for i := 0; i < s.count; i++ {
			if s.pointers[i] == ptr && s.live[i] {
				s.live[i] = false
				s.freed++
				found = true
				break
			}
		}
		if !found {
			s.valid = false
		}
	}
	if s.theme {
		themeTestReleased(ptr)
	}
}

func allocationTestGridStart(fail int) {
	s := allocationTestCurrent(true)
	s.grid, s.fail, s.count, s.freed, s.valid = true, fail, 0, 0, true
}
func allocationTestGridStop() {
	if s := allocationTestCurrent(false); s != nil {
		s.grid, s.fail = false, 0
		s.stopIfIdle()
	}
}
func allocationTestGridStat(index int) int {
	s := allocationTestCurrent(false)
	if s == nil {
		return 0
	}
	switch index {
	case -1:
		return s.count
	case -2:
		return s.freed
	case -3:
		return bool2int(s.valid)
	default:
		return int(s.sizes[index])
	}
}
func allocationTestGridContains(ptr unsafe.Pointer) int {
	s := allocationTestCurrent(false)
	if s == nil {
		return 0
	}
	for i := 0; i < s.count; i++ {
		if s.pointers[i] == ptr && s.live[i] {
			return 1
		}
	}
	return 0
}

func allocationTestResourceStart(events unsafe.Pointer, capacity int) {
	s := allocationTestCurrent(true)
	s.events, s.capacity, s.eventCount = events, capacity, 0
}
func allocationTestResourceStop() int {
	s := allocationTestCurrent(false)
	if s == nil {
		return 0
	}
	count := s.eventCount
	s.events = nil
	s.stopIfIdle()
	return count
}
func allocationTestThemeObserve(active bool) {
	s := allocationTestCurrent(active)
	if s == nil {
		return
	}
	s.theme = active
	if !active {
		s.stopIfIdle()
	}
}
func allocationTestThemeState() int {
	s := allocationTestCurrent(false)
	return bool2int(s != nil && s.theme)
}

// The caller is pinned, so this pinned goroutine necessarily uses another thread.
func allocationTestOtherThread(fn func() int) int {
	done := make(chan int, 1)
	go func() {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
		done <- fn()
	}()
	return <-done
}
func allocationTestOtherThreadThemeState() int {
	return allocationTestOtherThread(allocationTestThemeState)
}
func allocationTestOtherThreadAllocation() int {
	return allocationTestOtherThread(func() int {
		ptr := alloc.RawCalloc(128, 4)
		if ptr != nil {
			*(*byte)(ptr) = 0x5a
		}
		alloc.RawFree(ptr)
		return bool2int(ptr != nil)
	})
}
