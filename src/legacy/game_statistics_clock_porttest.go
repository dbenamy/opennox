//go:build porttest

package legacy

// PortTestStatisticsClock reuses the existing map-theme time/allocation observer.
// Disable allocation callbacks while the statistics fixture owns the observer;
// actual allocations and frees still execute through the original allocator.
func PortTestStatisticsClock() (set func(uint32), restore func()) {
	oldAlloc, oldFree := themeTestOnAllocate, themeTestOnRelease
	themeTestOnAllocate, themeTestOnRelease = nil, nil
	return func(value uint32) { themeObserve(true, value) }, func() {
		themeObserve(false, 0)
		themeTestOnAllocate, themeTestOnRelease = oldAlloc, oldFree
	}
}
