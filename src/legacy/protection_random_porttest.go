//go:build porttest

package legacy

/*
#include <stdint.h>
#include "GAME5_2.h"
extern uint32_t dword_5d4594_2516372;
extern uint32_t dword_5d4594_2516380;
*/
import "C"

import (
	"math"
	"unsafe"

	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/internal/protection"
)

type PortTestRandomOp struct {
	Mode int
	A, B uint32
}
type PortTestRandomSnapshot struct {
	Random protection.Random
	Result uint64
}

func PortTestProtectionRandom(seed uint32, ops []PortTestRandomOp) []PortTestRandomSnapshot {
	state := (*[5]float64)(memmap.PtrOff(0x5D4594, 2516388))
	max := memmap.PtrUint32(0x5D4594, 2516376)
	oldState := *state
	oldRange := [3]uint32{uint32(C.dword_5d4594_2516372), uint32(C.dword_5d4594_2516380), *max}
	constants := unsafe.Slice((*float64)(memmap.PtrOff(0x581450, 11344)), 5)
	oldConstants := append([]float64(nil), constants...)
	defer func() {
		*state = oldState
		C.dword_5d4594_2516372, C.dword_5d4594_2516380, *max = C.uint32_t(oldRange[0]), C.uint32_t(oldRange[1]), oldRange[2]
		copy(constants, oldConstants)
	}()
	copy(constants, []float64{0x1p-32, 2111111111, 1492, 1776, 5115})
	snapshot := func(result uint64) PortTestRandomSnapshot {
		return PortTestRandomSnapshot{Random: protection.Random{State: *state, Span: uint32(C.dword_5d4594_2516372), Min: uint32(C.dword_5d4594_2516380), Max: *max}, Result: result}
	}
	C.sub_56FF00(C.int(seed))
	out := []PortTestRandomSnapshot{snapshot(0)}
	for _, op := range ops {
		var result uint64
		switch op.Mode {
		case 0:
			result = math.Float64bits(float64(C.nox_xxx_unkDoubleSmth_56FE30()))
		case 1:
			result = uint64(uint32(C.sub_56FF80(C.int(op.A), C.int(op.B))))
		case 2:
			result = uint64(uint32(C.nox_xxx_protect_56F240()))
		case 3:
			C.sub_56FF00(C.int(op.A))
		default:
			panic("invalid random fixture mode")
		}
		out = append(out, snapshot(result))
	}
	return out
}
