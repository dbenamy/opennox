//go:build porttest

package legacy

/*
#include <stddef.h>
#include <stdint.h>
#include "GAME5_2.h"

extern uint64_t qword_581450_9544;

static double* portTestDurabilityQuarter(void) {
	return getMemDoublePtr(0x581450, 9608);
}

// Keeping the loop in C makes the exhaustive boundary fixture one cgo call.
static void portTestDurabilityBatch(const uint16_t* current, const uint16_t* maximum, int* out, size_t n) {
	for (size_t i = 0; i < n; i++) {
		out[i] = sub_57B190(current[i], maximum[i]);
	}
}
*/
import "C"

import "unsafe"

type PortTestDurabilityCase struct {
	Current uint16
	Maximum uint16
}

type PortTestDurabilitySnapshot struct {
	Results                               []int
	BeforeHalf, BeforeQuarter             uint64
	AfterCallHalf, AfterCallQuarter       uint64
	AfterRestoreHalf, AfterRestoreQuarter uint64
}

// PortTestDurability invokes the original C classifier after replacing its two
// threshold doubles by raw bits. It observes that C only reads them, then
// restores both globals even on a test panic.
func PortTestDurability(cases []PortTestDurabilityCase, halfBits, quarterBits uint64) (snap PortTestDurabilitySnapshot) {
	if len(cases) == 0 {
		return snap
	}
	half := (*uint64)(unsafe.Pointer(&C.qword_581450_9544))
	quarter := (*uint64)(unsafe.Pointer(C.portTestDurabilityQuarter()))
	oldHalf, oldQuarter := *half, *quarter
	snap.BeforeHalf, snap.BeforeQuarter = oldHalf, oldQuarter
	*half, *quarter = halfBits, quarterBits
	defer func() {
		*half, *quarter = oldHalf, oldQuarter
		snap.AfterRestoreHalf, snap.AfterRestoreQuarter = *half, *quarter
	}()

	cur := make([]uint16, len(cases))
	max := make([]uint16, len(cases))
	out := make([]C.int, len(cases))
	for i, c := range cases {
		cur[i], max[i] = c.Current, c.Maximum
	}
	C.portTestDurabilityBatch(
		(*C.uint16_t)(unsafe.Pointer(unsafe.SliceData(cur))),
		(*C.uint16_t)(unsafe.Pointer(unsafe.SliceData(max))),
		(*C.int)(unsafe.Pointer(unsafe.SliceData(out))), C.size_t(len(out)),
	)
	snap.AfterCallHalf, snap.AfterCallQuarter = *half, *quarter
	snap.Results = make([]int, len(out))
	for i, v := range out {
		snap.Results[i] = int(v)
	}
	return snap
}
