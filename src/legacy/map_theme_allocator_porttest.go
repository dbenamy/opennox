//go:build porttest

package legacy

import (
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"runtime"
)

var themeObserverPinned bool
var themeObserverOldClock func() uint32

func themeObserve(active bool, epoch uint32) {
	// Allocation observation belongs to the pinned fixture thread; other
	// threads must not inherit callbacks or failure injection.
	if active && !themeObserverPinned {
		runtime.LockOSThread()
		themeObserverPinned = true
		themeObserverOldClock = mapThemeClock
	}
	if active {
		mapThemeClock = func() uint32 { return epoch }
	}
	allocationTestThemeObserve(active)
	// Fixtures reset observation more than once during teardown. Pin only the
	// active interval, preserving any outer LockOSThread held by the fixture.
	if !active && themeObserverPinned {
		mapThemeClock = themeObserverOldClock
		themeObserverOldClock = nil
		themeObserverPinned = false
		runtime.UnlockOSThread()
	}
}

// PortTestThemeObserverThreadScope exercises activation, repeated activation and
// idempotent cleanup around an actual foreign thread.
func PortTestThemeObserverThreadScope() [4]int {
	before := int(allocationTestThemeState())
	themeObserve(true, 12345)
	defer themeObserve(false, 0)
	themeObserve(true, 12345)
	own := int(allocationTestThemeState())
	other := int(allocationTestOtherThreadThemeState())
	themeObserve(false, 0)
	themeObserve(false, 0)
	return [4]int{before, own, other, int(allocationTestThemeState())}
}

// Exercise the same clock source used by the theme loader.
func portTestThemeClock() uint32 { return mapThemeClock() }

// PortTestAllocationObserverIsolation observes actual allocations on two threads.
// It preserves safe calloc's nil marker and removes only its own marker afterward.
func PortTestAllocationObserverIsolation() [8]int {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if alloc.PortTestAllocationLive(nil) {
		panic("preexisting nil allocation marker")
	}
	before := alloc.PortTestAllocationCount()
	allocationTestGridStart(1)
	defer allocationTestGridStop()
	other := int(allocationTestOtherThreadAllocation())
	otherCount := int(allocationTestGridStat(-1))
	first := legacyCalloc(128, 4)
	firstNil := first == nil
	if first != nil {
		legacyFree(first)
	}
	if alloc.PortTestAllocationLive(nil) {
		alloc.FreePtr(nil)
	}
	second := legacyCalloc(128, 4)
	secondOK := second != nil
	if second != nil {
		legacyFree(second)
	}
	if alloc.PortTestAllocationLive(nil) {
		alloc.FreePtr(nil)
	}
	return [8]int{other, otherCount, bool2int(firstNil), bool2int(secondOK), int(allocationTestGridStat(-1)), int(allocationTestGridStat(-2)), int(allocationTestGridStat(-3)), alloc.PortTestAllocationCount() - before}
}
