//go:build porttest

package legacy

/*
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"

	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

// This bridge starts on the original profile-specific libc path. Conversion
// changes only the three calls below to the centralized legacy helpers.
func portAllocationCalloc(n, size uintptr) unsafe.Pointer {
	return C.calloc(C.size_t(n), C.size_t(size))
}
func portAllocationRealloc(p unsafe.Pointer, size uintptr) unsafe.Pointer {
	return C.realloc(p, C.size_t(size))
}
func portAllocationFree(p unsafe.Pointer) { C.free(p) }

type PortAllocationStep struct {
	Size        int
	Prefix      []byte
	Live        bool
	CountDelta  int
	OldReleased bool
}

type PortAllocationCapture struct {
	Initial           []byte
	InitialLive       bool
	InitialCountDelta int
	Steps             []PortAllocationStep
	FinalLive         bool
	FinalCountDelta   int
}

// PortTestLegacyAllocation checks bounded successful allocations and reallocs.
// It captures preserved bytes only; newly grown bytes have unspecified values.
// Addresses and whether realloc happened in place are deliberately not captured.
func PortTestLegacyAllocation(count, size int, resizes []int) (out PortAllocationCapture) {
	before := alloc.PortTestAllocationCount()
	p := portAllocationCalloc(uintptr(count), uintptr(size))
	if p == nil {
		panic("bounded calloc failed")
	}
	defer func() {
		if p != nil {
			portAllocationFree(p)
		}
	}()
	n := count * size
	out.Initial = append([]byte(nil), unsafe.Slice((*byte)(p), n)...)
	out.InitialLive = alloc.PortTestAllocationLive(p)
	out.InitialCountDelta = alloc.PortTestAllocationCount() - before
	fill := func(p unsafe.Pointer, n, stage int) {
		buf := unsafe.Slice((*byte)(p), n)
		for i := range buf {
			buf[i] = byte(i*37 + stage*53 + 19)
		}
	}
	fill(p, n, 0)
	for stage, next := range resizes {
		old := p
		q := portAllocationRealloc(p, uintptr(next))
		if q == nil {
			panic("bounded realloc failed")
		}
		p = q
		prefix := n
		if next < prefix {
			prefix = next
		}
		out.Steps = append(out.Steps, PortAllocationStep{
			Size:        next,
			Prefix:      append([]byte(nil), unsafe.Slice((*byte)(p), prefix)...),
			Live:        alloc.PortTestAllocationLive(p),
			CountDelta:  alloc.PortTestAllocationCount() - before,
			OldReleased: old == p || !alloc.PortTestAllocationLive(old),
		})
		n = next
		fill(p, n, stage+1)
	}
	last := p
	portAllocationFree(p)
	p = nil
	out.FinalLive = alloc.PortTestAllocationLive(last)
	out.FinalCountDelta = alloc.PortTestAllocationCount() - before
	return out
}

// PortTestLegacyCStringOwnership exercises the real CString/StrFree pair,
// including cgo's special malloc helper and the safe profile's macro route.
func PortTestLegacyCStringOwnership(s string) (out PortAllocationCapture) {
	before := alloc.PortTestAllocationCount()
	p := CString(s)
	ptr := unsafe.Pointer(p)
	out.Initial = append([]byte(nil), unsafe.Slice((*byte)(ptr), len(s)+1)...)
	out.InitialLive = alloc.PortTestAllocationLive(ptr)
	out.InitialCountDelta = alloc.PortTestAllocationCount() - before
	StrFree(p)
	out.FinalLive = alloc.PortTestAllocationLive(ptr)
	out.FinalCountDelta = alloc.PortTestAllocationCount() - before
	return out
}
