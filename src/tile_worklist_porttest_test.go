//go:build porttest

package opennox

import (
	"reflect"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
)

const worklistWords = 1500
const worklistNil = ^uint16(0)

type worklistState struct {
	count, overflow uint32
	queue           []uint32
}

func worklistI32(bits uint32) int32 { return int32(bits) }

func worklistRefGet(ref uint16, sep *[3]uint32, st *worklistState) (uint32, bool) {
	switch {
	case ref < 3:
		return sep[ref], true
	case ref == 3:
		return st.count, true
	case ref == 4:
		return st.overflow, true
	case ref == worklistNil:
		return 0, false
	default:
		return st.queue[int(ref)-5], true
	}
}
func worklistRefSet(ref uint16, v uint32, sep *[3]uint32, st *worklistState) {
	switch {
	case ref < 3:
		sep[ref] = v
	case ref == 3:
		st.count = v
	case ref == 4:
		st.overflow = v
	case ref != worklistNil:
		st.queue[int(ref)-5] = v
	}
}

// C scans with signed count but caps with unsigned count. Keep those distinct.
func modelEnqueue(st *worklistState, s legacy.PortTestTileWorklistSpec) {
	if s.X <= 0 || s.X >= 127 || s.Y <= 0 || s.Y >= 127 || s.Flags&3 == 0 {
		return
	}
	match := s.Flags&2 != 0 && s.Key == s.Field2 || s.Flags&1 != 0 && s.Key == s.Field1
	if !match || s.Flags&1 != 0 && s.Y == 1 || s.Flags&2 != 0 && s.X == 1 {
		return
	}
	if int32(st.count) > 0 {
		for i := uint32(0); i < st.count; i++ {
			j := 3 * i
			if st.queue[j] == uint32(s.X) && st.queue[j+1] == uint32(s.Y) && st.queue[j+2] == uint32(s.Flags) {
				return
			}
		}
	}
	if st.count >= 500 {
		st.overflow = 1
		return
	}
	j := 3 * st.count
	st.queue[j], st.queue[j+1], st.queue[j+2] = uint32(s.X), uint32(s.Y), uint32(s.Flags)
	st.count++
}

// This models the C stores/reloads in order, including output aliases.
func modelPop(st *worklistState, s legacy.PortTestTileWorklistSpec) (int, [3]uint32, [3]bool) {
	sep := [3]uint32{0x01020304, 0x11223344, 0x55667788}
	read := func() ([3]uint32, [3]bool) {
		a, ao := worklistRefGet(s.PopX, &sep, st)
		b, bo := worklistRefGet(s.PopY, &sep, st)
		c, co := worklistRefGet(s.PopZ, &sep, st)
		return [3]uint32{a, b, c}, [3]bool{ao, bo, co}
	}
	if int32(st.count) <= 0 {
		o, ok := read()
		return 0, o, ok
	}
	st.count--
	worklistRefSet(s.PopX, st.queue[3*st.count], &sep, st)
	worklistRefSet(s.PopY, st.queue[3*st.count+1], &sep, st)
	worklistRefSet(s.PopZ, st.queue[3*st.count+2], &sep, st)
	o, ok := read()
	return 1, o, ok
}

func checkWorklistRun(t *testing.T, initialCount, initialOverflow uint32, queue []uint32, specs []legacy.PortTestTileWorklistSpec) legacy.PortTestTileWorklistSnapshot {
	t.Helper()
	st := worklistState{initialCount, initialOverflow, append([]uint32(nil), queue...)}
	got := legacy.PortTestTileWorklist(initialCount, initialOverflow, queue, specs)
	if len(got.Results) != len(specs) || !got.GridPointerRestored || !reflect.DeepEqual(got.Before, got.AfterRestore) {
		t.Fatalf("fixture did not restore original globals: before=%+v after=%+v ptr=%v", got.Before, got.AfterRestore, got.GridPointerRestored)
	}
	for i, s := range specs {
		ret := 0
		var out [3]uint32
		var readable [3]bool
		if s.Op == 0 {
			modelEnqueue(&st, s)
		} else {
			ret, out, readable = modelPop(&st, s)
		}
		r := got.Results[i]
		if r.Return != ret || r.Count != st.count || r.Overflow != st.overflow || !reflect.DeepEqual(r.Queue, st.queue) || !r.GridUnchanged || !r.GridPointerUnchanged || !r.QueueGuardsUnchanged || !r.OutputGuardsOK || (s.Op == 1 && (r.Outputs != out || r.OutputReadable != readable)) {
			t.Fatalf("case %d spec=%+v got ret/count/overflow/out/ok=%d/%x/%x/%x/%v want=%d/%x/%x/%x/%v", i, s, r.Return, r.Count, r.Overflow, r.Outputs, r.OutputReadable, ret, st.count, st.overflow, out, readable)
		}
	}
	return got
}

func TestTileWorklistCABI(t *testing.T) {
	q := make([]uint32, worklistWords)
	for i := range q {
		q[i] = uint32(i%499 + 1)
	}
	// Full coordinate Cartesian boundary grid, all low-bit flag combinations
	// plus raw signed extrema, and all four key-match conditions.
	coords := []int32{-1, 0, 1, 2, 126, 127, 128, worklistI32(0x80000000), worklistI32(0x7fffffff)}
	flags := []int32{0, 1, 2, 3, 4, 5, 6, 7, worklistI32(0x80000000), worklistI32(0x7fffffff), worklistI32(0xffffffff), worklistI32(0x80000001)}
	var specs []legacy.PortTestTileWorklistSpec
	for _, x := range coords {
		for _, y := range coords {
			for _, f := range flags {
				// neither, field1-only, field2-only, both. For zero low bits these
				// deliberately exercise the nil-grid no-dereference path as well.
				inputs := []legacy.PortTestTileWorklistSpec{
					{Op: 0, X: x, Y: y, Flags: f, Key: 0x33, Field1: 0x11, Field2: 0x22, NilGrid: f&3 == 0},
					{Op: 0, X: x, Y: y, Flags: f, Key: 0x11, Field1: 0x11, Field2: 0x22, NilGrid: f&3 == 0},
					{Op: 0, X: x, Y: y, Flags: f, Key: 0x22, Field1: 0x11, Field2: 0x22, NilGrid: f&3 == 0},
					{Op: 0, X: x, Y: y, Flags: f, Key: 0x55, Field1: 0x55, Field2: 0x55, NilGrid: f&3 == 0},
				}
				for _, in := range inputs {
					specs = append(specs, in, legacy.PortTestTileWorklistSpec{Op: 1, PopX: 0, PopY: 1, PopZ: 2})
				}
			}
		}
	}
	// Explicit safe out-of-range nil-grid calls with flag bits set prove bounds
	// happen before the grid load.
	specs = append(specs,
		legacy.PortTestTileWorklistSpec{Op: 0, X: 0, Y: 2, Flags: 3, Key: 1, NilGrid: true}, legacy.PortTestTileWorklistSpec{Op: 1, PopX: 0, PopY: 1, PopZ: 2},
		legacy.PortTestTileWorklistSpec{Op: 0, X: 127, Y: 2, Flags: 3, Key: 1, NilGrid: true}, legacy.PortTestTileWorklistSpec{Op: 1, PopX: 0, PopY: 1, PopZ: 2})
	for _, key := range []uint32{0, 0x7fffffff, 0x80000000, 0xffffffff} {
		for _, f := range []int32{1, 2, 3, worklistI32(0x80000001), worklistI32(0x7fffffff), -1} {
			for match := 0; match < 4; match++ {
				a, b := key^1, key^2
				if match&1 != 0 {
					a = key
				}
				if match&2 != 0 {
					b = key
				}
				specs = append(specs, legacy.PortTestTileWorklistSpec{X: 4, Y: 5, Flags: f, Key: key, Field1: a, Field2: b}, legacy.PortTestTileWorklistSpec{Op: 1, PopX: 0, PopY: 1, PopZ: 2})
			}
		}
	}
	if len(specs) != 7972 {
		t.Fatalf("matrix size %d", len(specs))
	}
	checkWorklistRun(t, 0, 0xa5a5a5a5, q, specs)
}

func TestTileWorklistCapacityDuplicatesAndLIFO(t *testing.T) {
	q := make([]uint32, worklistWords)
	for i := 0; i < 500; i++ {
		q[3*i], q[3*i+1], q[3*i+2] = uint32(i%124+2), uint32(i/124+2), uint32(i+9)
	}
	// Transition to full then test duplicates at a middle and final active slot,
	// followed by a unique qualifying request and a pop that preserves overflow.
	q[3*200], q[3*200+1], q[3*200+2] = 40, 41, 5
	checkWorklistRun(t, 499, 0x91a2b3c4, q, []legacy.PortTestTileWorklistSpec{
		{Op: 0, X: 44, Y: 45, Flags: 1, Key: 8, Field1: 9}, // gate miss preserves poison
		{Op: 0, X: 44, Y: 45, Flags: 1, Key: 9, Field1: 9}, // 499 -> 500
		{Op: 0, X: 44, Y: 45, Flags: 1, Key: 8, Field1: 9}, // full, gate miss
		{Op: 0, X: 127, Y: 45, Flags: 1, NilGrid: true},    // full, bounds miss
		{Op: 0, X: 40, Y: 41, Flags: 5, Key: 9, Field1: 9}, // middle duplicate
		{Op: 0, X: 44, Y: 45, Flags: 1, Key: 9, Field1: 9}, // final (newly appended) duplicate
		{Op: 0, X: 46, Y: 47, Flags: 1, Key: 9, Field1: 9}, // unique overflow
		{Op: 1, PopX: 0, PopY: 1, PopZ: 2},
		{X: 46, Y: 47, Flags: 1, Key: 9, Field1: 9}, // append after overflow remains sticky
	})
	// Key is excluded from duplicate identity; all flag bits remain significant.
	checkWorklistRun(t, 0, 0x77, q, []legacy.PortTestTileWorklistSpec{
		{X: 9, Y: 10, Flags: 1, Key: 0x80000000, Field1: 0x80000000},
		{X: 9, Y: 10, Flags: 1, Key: 0xffffffff, Field1: 0xffffffff},
		{X: 9, Y: 10, Flags: 5, Key: 0xffffffff, Field1: 0xffffffff},
		{Op: 1, PopX: 0, PopY: 1, PopZ: 2}, {Op: 1, PopX: 0, PopY: 1, PopZ: 2},
		{Op: 1, PopX: worklistNil, PopY: worklistNil, PopZ: worklistNil},
	})
	// A high-bit count never scans or pops, but unsigned capacity is already
	// full. Nil outputs are valid only on that no-write pop path.
	for _, n := range []uint32{0x80000000, 0xffffffff} {
		checkWorklistRun(t, n, 0x91a2b3c4, q, []legacy.PortTestTileWorklistSpec{
			{Op: 0, X: 4, Y: 4, Flags: 1, Key: 2, Field1: 1},
			{Op: 0, X: 4, Y: 4, Flags: 0, Key: 1, NilGrid: true},
			{Op: 0, X: 4, Y: 4, Flags: 1, Key: 1, Field1: 1},
			{Op: 1, PopX: worklistNil, PopY: worklistNil, PopZ: worklistNil},
		})
	}
	checkWorklistRun(t, 0, 0x77, q, []legacy.PortTestTileWorklistSpec{{Op: 1, PopX: worklistNil, PopY: worklistNil, PopZ: worklistNil}})
}

func TestTileWorklistPopAliases(t *testing.T) {
	q := make([]uint32, worklistWords)
	for i := range q {
		q[i] = uint32(i%499 + 1)
	}
	copy(q[:9], []uint32{2, 3, 4, 1, 5, 6, 1, 7, 8}) // safe values for count reloads
	// Separate destinations, every pair sharing a destination, count, overflow,
	// and current record queue words all appear. Runs are independent so aliases
	// that alter count cannot make a later case unsafe.
	refs := []uint16{0, 1, 2, 3, 4, 5, 8, 11, 12, 13}
	for _, x := range refs {
		for _, y := range refs {
			for _, z := range refs {
				checkWorklistRun(t, 3, 0x77, q, []legacy.PortTestTileWorklistSpec{{Op: 1, PopX: x, PopY: y, PopZ: z}})
			}
		}
	}
	// Plain LIFO has an explicit oracle in addition to the generic state model.
	lifo := []legacy.PortTestTileWorklistSpec{{Op: 1, PopX: 0, PopY: 1, PopZ: 2}, {Op: 1, PopX: 0, PopY: 1, PopZ: 2}, {Op: 1, PopX: 0, PopY: 1, PopZ: 2}, {Op: 1, PopX: 0, PopY: 1, PopZ: 2}}
	got := checkWorklistRun(t, 3, 0xa5a5a5a5, q, lifo)
	want := [][3]uint32{{1, 7, 8}, {1, 5, 6}, {2, 3, 4}, {0x01020304, 0x11223344, 0x55667788}}
	for i, r := range got.Results {
		wantRet := 0
		if i < 3 {
			wantRet = 1
		}
		if r.Return != wantRet || r.Outputs != want[i] || r.OutputReadable != [3]bool{true, true, true} {
			t.Fatalf("plain LIFO %d: %+v", i, r)
		}
	}
}
