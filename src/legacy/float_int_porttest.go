//go:build porttest

package legacy

import (
	"github.com/opennox/opennox/v1/legacy/common/fpenv"
	"math"
	"runtime"
)

type PortTestFloatIntResult struct {
	Native    [3]int32
	ControlOK bool
}

// PortTestFloatInt checks native conversion and narrowing with hostile x87 controls.
// cwMask selects PC/RC bits only; -1 preserves the original control word.
func PortTestFloatInt(bits []uint32, cwMask int) (out []PortTestFloatIntResult) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	before := fpenv.Control()
	defer fpenv.SetControl(before)
	control := before
	if cwMask >= 0 {
		control = (before &^ 0x0f00) | uint16(cwMask&0x0f00)
	}
	fpenv.SetControl(control)
	out = make([]PortTestFloatIntResult, len(bits))
	for i, b := range bits {
		out[i] = PortTestFloatIntResult{Native: [3]int32{floatToInt32(math.Float32frombits(b)), int32(int16(floatToInt32(math.Float32frombits(b)))), int32(int16(floatToInt32(math.Float32frombits(b & 0x7fffffff))))}, ControlOK: fpenv.Control() == control}
	}
	return out
}

func PortTestFloatIntBenchmark(n int) uint32 {
	var sum uint32
	for i := 0; i < n; i++ {
		sum += uint32(floatToInt32(float32(i&1023) + 0.75))
	}
	return sum
}

// PortTestDoubleInt exercises native conversion under the frozen precision/rounding controls.
func PortTestDoubleInt(bits []uint64, cwMask int) ([]int32, bool) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	before := fpenv.Control()
	defer fpenv.SetControl(before)
	control := before
	if cwMask >= 0 {
		control = (before &^ 0x0f00) | uint16(cwMask&0x0f00)
	}
	fpenv.SetControl(control)
	out := make([]int32, len(bits))
	for i, b := range bits {
		out[i] = doubleToInt32(math.Float64frombits(b))
	}
	return out, fpenv.Control() == control
}
