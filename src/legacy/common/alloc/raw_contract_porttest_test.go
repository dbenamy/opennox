//go:build porttest && linux && 386

package alloc

import (
	"bytes"
	"runtime"
	"sync"
	"testing"
	"unsafe"
)

func TestRawAllocationBoundaries(t *testing.T) {
	for _, size := range []uintptr{0, 1, 7, 8, 15, 16, 17, 255, 256, 257, 4095, 4096, 4097, 65535, 65536, 65537} {
		p := RawCalloc(1, size)
		if p == nil {
			t.Fatalf("calloc(%d) failed", size)
		}
		if uintptr(p)%16 != 0 {
			t.Fatalf("calloc(%d) alignment", size)
		}
		data := unsafe.Slice((*byte)(p), int(size))
		if !bytes.Equal(data, make([]byte, int(size))) {
			t.Fatalf("calloc(%d) not zero", size)
		}
		for i := range data {
			data[i] = byte(i*31 + 7)
		}
		if q := RawRealloc(p, ^uintptr(0)); q != nil {
			RawFree(q)
			t.Fatal("oversized realloc succeeded")
		}
		for i, v := range data {
			if v != byte(i*31+7) {
				t.Fatalf("failed realloc changed byte %d", i)
			}
		}
		q := RawRealloc(p, size+127)
		if q == nil {
			RawFree(p)
			t.Fatal("bounded realloc failed")
		}
		for i, v := range unsafe.Slice((*byte)(q), int(size)) {
			if v != byte(i*31+7) {
				t.Fatalf("successful realloc changed byte %d", i)
			}
		}
		if got := RawRealloc(q, 0); got != nil {
			RawFree(got)
			t.Fatal("realloc(non-nil,0) should release and return nil")
		}
	}
	for _, pair := range [][2]uintptr{{0, 0}, {0, ^uintptr(0)}, {^uintptr(0), 0}} {
		p := RawCalloc(pair[0], pair[1])
		if p == nil {
			t.Fatalf("zero-product calloc %v returned nil", pair)
		}
		RawFree(p)
	}
	for _, pair := range [][2]uintptr{{^uintptr(0), 1}, {^uintptr(0), 2}, {^uintptr(0)/2 + 1, 2}} {
		if p := RawCalloc(pair[0], pair[1]); p != nil {
			RawFree(p)
			t.Fatalf("oversized calloc %v succeeded", pair)
		}
	}
	p := RawRealloc(nil, 0)
	if p == nil {
		t.Fatal("realloc(nil,0) returned nil")
	}
	RawFree(p)
	p = RawRealloc(nil, 257)
	if p == nil {
		t.Fatal("realloc(nil,257) failed")
	}
	RawFree(p)
	RawFree(nil)
}

func TestRawAllocationThreadHandoff(t *testing.T) {
	type item struct {
		ptr  unsafe.Pointer
		size int
		seed byte
	}
	const workers = 8
	work := make(chan item, 16)
	var producers, consumers sync.WaitGroup
	for n := 0; n < workers; n++ {
		consumers.Add(1)
		go func() {
			defer consumers.Done()
			for x := range work {
				data := unsafe.Slice((*byte)(x.ptr), x.size)
				for i, v := range data {
					if v != byte(i*37)+x.seed {
						t.Errorf("handoff content changed at %d", i)
						break
					}
				}
				RawFree(x.ptr)
			}
		}()
		producers.Add(1)
		go func(seed byte) {
			defer producers.Done()
			for n := 0; n < 128; n++ {
				size := 1 + (n*127+int(seed)*31)%10007
				p := RawCalloc(1, uintptr(size))
				if p == nil {
					t.Error("handoff allocation failed")
					return
				}
				data := unsafe.Slice((*byte)(p), size)
				for i := range data {
					data[i] = byte(i*37) + seed
				}
				work <- item{p, size, seed}
			}
		}(byte(n))
	}
	producers.Wait()
	close(work)
	consumers.Wait()
}

// Deliberately return only the address: engine records retain 32-bit raw words.
//
//go:noinline
func rawAddressContractAllocate() uintptr {
	p := RawCalloc(1, 257)
	if p == nil {
		panic("bounded raw allocation failed")
	}
	data := unsafe.Slice((*byte)(p), 257)
	for i := range data {
		data[i] = byte(i*41 + 11)
	}
	return uintptr(p)
}

func TestRawAllocationAddressLifetime(t *testing.T) {
	address := rawAddressContractAllocate()
	for i := 0; i < 3; i++ {
		runtime.GC()
	}
	p := unsafe.Pointer(address)
	defer RawFree(p)
	for i, v := range unsafe.Slice((*byte)(p), 257) {
		if v != byte(i*41+11) {
			t.Fatalf("raw-address lifetime changed byte %d", i)
		}
	}
}

func TestRawAllocationMixedLifetimes(t *testing.T) {
	type block struct {
		p    unsafe.Pointer
		size int
		seed byte
	}
	var slots [32]block
	sizes := [...]int{0, 1, 17, 48, 256, 1000, 4097, 16385, 65535, 131073, 1048577}
	check := func(x block, n int) {
		t.Helper()
		for i, v := range unsafe.Slice((*byte)(x.p), n) {
			if v != byte(i*19)+x.seed {
				t.Fatalf("live block content changed at %d/%d", i, n)
			}
		}
	}
	defer func() {
		for _, x := range slots {
			RawFree(x.p)
		}
	}()
	random := uint32(12345)
	for step := 0; step < 512; step++ {
		random = random*1664525 + 1013904223
		index := int(random>>16) % len(slots)
		x := slots[index]
		size := sizes[int(random>>8)%len(sizes)]
		if x.p != nil {
			check(x, x.size)
		}
		if x.p != nil && step%3 == 0 {
			if size == 0 {
				size = 1
			}
			p := RawRealloc(x.p, uintptr(size))
			if p == nil {
				t.Fatal("bounded mixed realloc failed")
			}
			oldSize := x.size
			x.p = p
			slots[index] = x
			prefix := oldSize
			if size < prefix {
				prefix = size
			}
			check(x, prefix)
		} else {
			RawFree(x.p)
			slots[index] = block{}
			x.p = RawCalloc(1, uintptr(size))
			if x.p == nil {
				t.Fatal("bounded mixed calloc failed")
			}
			slots[index] = x
			for i, v := range unsafe.Slice((*byte)(x.p), size) {
				if v != 0 {
					t.Fatalf("reused calloc not zero at %d", i)
				}
			}
		}
		x.size = size
		x.seed = byte(step*7 + 3)
		for i := range unsafe.Slice((*byte)(x.p), size) {
			*(*byte)(unsafe.Add(x.p, i)) = byte(i*19) + x.seed
		}
		slots[index] = x
	}
	for _, x := range slots {
		if x.p != nil {
			check(x, x.size)
		}
	}
}

func TestRawAlignedAllocation(t *testing.T) {
	for _, alignment := range []uintptr{16, 256, 4096, 65536} {
		for _, count := range []uintptr{0, 1, 2, 7, 17} {
			size := alignment * count
			p := RawAlignedAlloc(alignment, size)
			if p == nil {
				t.Fatalf("aligned allocation failed: %d/%d", alignment, size)
			}
			if uintptr(p)%alignment != 0 {
				RawFree(p)
				t.Fatalf("unaligned result: %d/%d", alignment, size)
			}
			data := unsafe.Slice((*byte)(p), size)
			for i := range data {
				data[i] = byte(i*13 + 9)
			}
			runtime.GC()
			for i, v := range data {
				if v != byte(i*13+9) {
					RawFree(p)
					t.Fatalf("aligned contents changed: %d", i)
				}
			}
			RawFree(p)
		}
	}
}
