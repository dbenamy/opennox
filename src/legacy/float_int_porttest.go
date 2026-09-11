//go:build porttest

package legacy

/*
#include <stdint.h>
#include "GAME1_1.h"
static unsigned short porttest_float_cw(void) {
 unsigned short cw; __asm__ volatile("fnstcw %0" : "=m"(cw)); return cw;
}
static void porttest_float_set_cw(unsigned short cw) {
 __asm__ volatile("fldcw %0" : : "m"(cw));
}
static uint32_t porttest_float_int_bench(unsigned n) {
 uint32_t sum = 0;
 for (unsigned i = 0; i < n; i++) sum += nox_float2int((float)(i & 1023) + 0.75f);
 return sum;
}
static void porttest_float_int_calls(uint32_t bits, int *out) {
 union { uint32_t bits; float value; } input;
 input.bits = bits;
 out[0] = nox_float2int(input.value);
 out[1] = nox_float2int16(input.value);
 out[2] = nox_float2int16_abs(input.value);
}
*/
import "C"

import (
	"math"
	"runtime"
	"unsafe"
)

type PortTestFloatIntResult struct {
	Native              [3]int32
	Values              [3]int32
	GuardsOK, ControlOK bool
}

// PortTestFloatInt batches live ABI calls with normal C short-to-int promotion.
// cwMask selects PC/RC bits only; -1 preserves the original control word.
func PortTestFloatInt(bits []uint32, cwMask int) (out []PortTestFloatIntResult) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	before := C.porttest_float_cw()
	defer C.porttest_float_set_cw(before)
	control := before
	if cwMask >= 0 {
		control = (before &^ 0x0f00) | C.ushort(cwMask&0x0f00)
	}
	C.porttest_float_set_cw(control)
	out = make([]PortTestFloatIntResult, len(bits))
	for i, b := range bits {
		words := [5]C.int{0x12345678, 0x34567812, 0x45678123, 0x56781234, 0x76543210}
		C.porttest_float_int_calls(C.uint32_t(b), (*C.int)(unsafe.Pointer(&words[1])))
		out[i] = PortTestFloatIntResult{Values: [3]int32{int32(words[1]), int32(words[2]), int32(words[3])}, Native: [3]int32{floatToInt32(math.Float32frombits(b)), int32(int16(floatToInt32(math.Float32frombits(b)))), int32(int16(floatToInt32(math.Float32frombits(b & 0x7fffffff))))}, GuardsOK: words[0] == 0x12345678 && words[4] == 0x76543210, ControlOK: C.porttest_float_cw() == control}
	}
	return out
}

func PortTestFloatIntBenchmark(n int) uint32 { return uint32(C.porttest_float_int_bench(C.uint(n))) }
