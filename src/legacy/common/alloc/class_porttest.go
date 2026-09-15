//go:build porttest

package alloc

// PortTestCounts observes the actual pool lists without changing ownership.
func (al *Class) PortTestCounts() [3]int {
	return [3]int{al.active.Count(), al.static.Count(), al.dynamic.Count()}
}
