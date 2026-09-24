//go:build porttest

package legacy

import (
	"unsafe"

	"github.com/opennox/opennox/v1/common/memmap"
)

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

// PortTestDurability invokes the Go owner after replacing its two
// threshold doubles by raw bits. It observes that the owner only reads them, then
// restores both globals even on a test panic.
func PortTestDurability(cases []PortTestDurabilityCase, halfBits, quarterBits uint64) (snap PortTestDurabilitySnapshot) {
	if len(cases) == 0 {
		return snap
	}
	half := (*uint64)(unsafe.Pointer(&qword_581450_9544))
	quarter := memmap.PtrUint64(0x581450, 9608)
	oldHalf, oldQuarter := *half, *quarter
	snap.BeforeHalf, snap.BeforeQuarter = oldHalf, oldQuarter
	*half, *quarter = halfBits, quarterBits
	defer func() {
		*half, *quarter = oldHalf, oldQuarter
		snap.AfterRestoreHalf, snap.AfterRestoreQuarter = *half, *quarter
	}()

	cur := make([]uint16, len(cases))
	max := make([]uint16, len(cases))
	out := make([]int, len(cases))
	for i, c := range cases {
		cur[i], max[i] = c.Current, c.Maximum
	}
	for i := range out {
		out[i] = durabilityBand(cur[i], max[i])
	}
	snap.AfterCallHalf, snap.AfterCallQuarter = *half, *quarter
	snap.Results = make([]int, len(out))
	for i, v := range out {
		snap.Results[i] = int(v)
	}
	return snap
}
