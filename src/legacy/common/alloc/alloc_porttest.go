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
