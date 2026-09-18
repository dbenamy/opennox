//go:build porttest

package legacy

import (
	"sort"
	"unsafe"
)

// Observe real allocations/releases, without replacing either operation.
func PortTestPrefabAllocationBalance(run func()) []int {
	oldAlloc, oldFree := themeTestOnAllocate, themeTestOnRelease
	live := map[unsafe.Pointer]int{}
	themeTestOnAllocate = func(p unsafe.Pointer, n int) { live[p] = n }
	themeTestOnRelease = func(p unsafe.Pointer) { delete(live, p) }
	themeObserve(true, 0)
	defer func() { themeObserve(false, 0); themeTestOnAllocate = oldAlloc; themeTestOnRelease = oldFree }()
	run()
	var out []int
	for _, n := range live {
		out = append(out, n)
	}
	sort.Ints(out)
	return out
}
