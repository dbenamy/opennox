//go:build porttest

package opennox

import (
	"crypto/sha256"
	"testing"
	"unsafe"
)

// The original C bridge checks capacity and passes &buf[0], erasing a nonempty
// slice's length. Match full-buffer behavior for every accepted shorter view.
func TestFloorAssetsScratchCapacity(t *testing.T) {
	o := newFloorAssetsOwner(t)
	storage := o.scratch
	type result struct {
		Return, Consumed  int
		Counts            [2]uint32
		Metadata, Scratch [32]byte
		Allocated         bool
	}
	run := func(op, length, capacity int) (r result, panicked bool) {
		o.reset()
		name := "BufferContract"
		payload := floorAssetsFloor(name, [3]byte{17, 129, 254}, [3]byte{1, 1, 1}, 1)
		if op == 1 || op == 3 || op == 5 {
			payload = floorAssetsEdge(name, [3]byte{1, 1, 1}, 1, 0, 0x454e4420)
		}
		if op == 4 {
			*o.count = 1
			copy(o.defs[0].NameBuf[:], name)
		}
		if op == 5 {
			*o.edgeCount = 1
			copy(o.edges[:32], name)
		}
		beforeDefs, beforeEdges := o.metadata()
		beforeHash := sha256.Sum256(append(beforeDefs, beforeEdges...))
		beforeScratch := sha256.Sum256(storage)
		beforeCounts := [2]uint32{*o.count, *o.edgeCount}
		func() {
			o.scratch = storage[:length:capacity]
			defer func() {
				o.scratch = storage
				if recover() != nil {
					panicked = true
				}
			}()
			r.Return, r.Consumed = o.invoke(t, op, payload)
		}()
		defs, edges := o.metadata()
		r.Counts = [2]uint32{*o.count, *o.edgeCount}
		r.Metadata = sha256.Sum256(append(defs, edges...))
		r.Scratch = sha256.Sum256(storage)
		r.Allocated = o.defs[0].Data32 != nil || *(*unsafe.Pointer)(unsafe.Pointer(&o.edges[32])) != nil
		if panicked && (r.Metadata != beforeHash || r.Scratch != beforeScratch || r.Counts != beforeCounts || r.Allocated) {
			t.Fatal("rejected scratch view changed owner state")
		}
		return
	}
	for _, op := range []int{0, 1, 3, 4, 5} {
		want, panicked := run(op, len(storage), cap(storage))
		if panicked || want.Return != 1 {
			t.Fatal("full scratch baseline failed")
		}
		for _, capacity := range []int{len(storage) - 1, len(storage)} {
			for _, length := range []int{0, 1, 31, 64, len(storage)} {
				if length > capacity {
					continue
				}
				got, panicked := run(op, length, capacity)
				wantPanic := capacity < len(storage) || length == 0
				if panicked != wantPanic || (!wantPanic && got != want) {
					t.Fatalf("op%d length%d capacity%d panic%v/%v result%+v want%+v", op, length, capacity, panicked, wantPanic, got, want)
				}
			}
		}
	}
}
