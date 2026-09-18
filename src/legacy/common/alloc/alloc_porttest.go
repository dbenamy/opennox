//go:build porttest

package alloc

import "unsafe"

// PortTestAllocationLive observes ownership without reading released storage.
func PortTestAllocationLive(ptr unsafe.Pointer) bool {
	allocMu.Lock()
	defer allocMu.Unlock()
	_, ok := allocs[ptr]
	return ok
}

// PortTestAllocationCount lets a sequential owner check complete cleanup of
// records which were rejected before they could enter its public list.
func PortTestAllocationCount() int {
	allocMu.Lock()
	defer allocMu.Unlock()
	return len(allocs)
}
