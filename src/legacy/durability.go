package legacy

/*
#include <stdint.h>
#include "GAME5_2.h"
extern uint64_t qword_581450_9544;
*/
import "C"

import (
	"math"

	"github.com/opennox/opennox/v1/common/memmap"
)

func durabilityBand(current, maximum uint16) int {
	if maximum == 0 {
		return 4
	}
	if current == maximum {
		return 0
	}
	// The thresholds remain shared with C. Preserve their ordered comparisons,
	// including equality at half/quarter and nonfinite threshold values.
	if float64(current) >= float64(maximum)*math.Float64frombits(uint64(C.qword_581450_9544)) {
		return 1
	}
	if float64(current) < float64(maximum)*memmap.Float64(0x581450, 9608) {
		return 3
	}
	return 2
}

//export sub_57B190
func sub_57B190(current, maximum C.ushort) C.int {
	return C.int(durabilityBand(uint16(current), uint16(maximum)))
}
