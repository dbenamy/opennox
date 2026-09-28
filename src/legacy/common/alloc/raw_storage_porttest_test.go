//go:build porttest && linux && 386

package alloc

import (
	"runtime"
	"testing"
	"unsafe"
)

func rawStorageSnapshot(t *testing.T) (live, idle int) {
	t.Helper()
	for i := range rawHeap.classes {
		rawHeap.classes[i].Lock()
	}
	defer func() {
		for i := len(rawHeap.classes) - 1; i >= 0; i-- {
			rawHeap.classes[i].Unlock()
		}
	}()
	spans := make(map[*rawSpan]bool)
	visit := func(_, value any) bool {
		s := value.(*rawSpan)
		if spans[s] {
			return true
		}
		spans[s] = true
		if !s.mapped {
			t.Fatal("unmapped span still registered")
		}
		live += s.live
		if s.live == 0 {
			idle += len(s.data)
		}
		return true
	}
	rawHeap.pages.Range(visit)
	rawHeap.exact.Range(visit)
	linked := make(map[*rawSpan]bool)
	for i := range rawHeap.classes {
		for s := rawHeap.classes[i].partial; s != nil; s = s.next {
			if linked[s] || !spans[s] {
				t.Fatal("duplicate/cyclic or unregistered partial span")
			}
			linked[s] = true
		}
	}
	return live, idle
}

// Bounded idle spans reuse realistic working sets without repeated mapping calls.
// All live ownership must retire.
func TestRawStorageReleaseBudget(t *testing.T) {
	before, _ := rawStorageSnapshot(t)
	var pointers []unsafe.Pointer
	defer func() {
		for _, p := range pointers {
			RawFree(p)
		}
	}()
	for i := 0; i < 65537; i++ {
		p := RawCalloc(1, 16)
		if p == nil {
			t.Fatal("small allocation failed")
		}
		pointers = append(pointers, p)
	}
	for _, size := range []uintptr{65536, 262144, 1048576, 1048593} {
		for i := 0; i < 17; i++ {
			p := RawCalloc(1, size)
			if p == nil {
				t.Fatal("large allocation failed")
			}
			pointers = append(pointers, p)
		}
	}
	// Dropping reuse hints during GC must not affect live unmanaged payloads.
	for i, p := range pointers {
		*(*byte)(p) = byte(i + 1)
	}
	runtime.GC()
	runtime.GC()
	extra := RawCalloc(1, 65536)
	if extra == nil {
		t.Fatal("allocation after GC failed")
	}
	RawFree(extra)
	for i, p := range pointers {
		if *(*byte)(p) != byte(i+1) {
			t.Fatalf("GC/reuse changed live block %d", i)
		}
	}
	for parity := 0; parity < 2; parity++ {
		for i := parity; i < len(pointers); i += 2 {
			RawFree(pointers[i])
			pointers[i] = nil
		}
	}
	live, idle := rawStorageSnapshot(t)
	if live != before {
		t.Fatalf("live backend ownership %d, want %d", live, before)
	}
	if idle > 18*1024*1024 {
		t.Fatalf("idle storage retained %d bytes, budget 18 MiB", idle)
	}
}

func TestRawStorageSubstantialShrink(t *testing.T) {
	p := RawCalloc(1, 2*1024*1024+17)
	if p == nil {
		t.Fatal("large allocation failed")
	}
	defer func() { RawFree(p) }()
	*(*byte)(p) = 0x5a
	next := RawRealloc(p, 17)
	if next == nil {
		t.Fatal("bounded shrink failed")
	}
	p = next
	if *(*byte)(p) != 0x5a {
		t.Fatal("shrink changed prefix")
	}
	s := rawFind(p)
	c := &rawHeap.classes[s.class]
	c.Lock()
	capacity := s.stride
	if capacity == 0 {
		capacity = uintptr(len(s.data))
	}
	c.Unlock()
	if capacity > 4096 {
		t.Fatalf("17-byte allocation still owns %d-byte block after shrink", capacity)
	}
}

// Independent guard-byte checks for the exact-byte clearing primitive, including
// unaligned starts, zero length, and boundaries around common vector/page sizes.
func TestRawClearBoundaries(t *testing.T) {
	rawClear(nil, 0)
	for _, n := range []int{0, 1, 2, 3, 4, 7, 8, 15, 16, 17, 31, 32, 33, 63, 64, 65, 127, 128, 129, 255, 256, 257, 4095, 4096, 4097, 65535, 65536, 65537} {
		for _, offset := range []int{0, 1, 3, 7, 15, 16, 31} {
			data := make([]byte, n+96)
			for i := range data {
				data[i] = 0xa5
			}
			start := 32 + offset
			rawClear(unsafe.Pointer(&data[start]), uintptr(n))
			for i, v := range data {
				want := byte(0xa5)
				if i >= start && i < start+n {
					want = 0
				}
				if v != want {
					t.Fatalf("size=%d offset=%d byte=%d: got %02x want %02x", n, offset, i, v, want)
				}
			}
		}
	}
}
