//go:build porttest

package legacy

import (
	"math"
	"unsafe"
)

func PortTestTextFormat(kind int, dst *uint16, count uint32, f *uint16, a, b uint32, bits uint64, w1, w2 *uint16, narrow *byte) int {
	var args []textFormatArgument
	switch kind {
	case 1:
		args = []textFormatArgument{textFormatWord(a), textFormatWord(b)}
	case 2:
		args = []textFormatArgument{textFormatReal(math.Float64frombits(bits))}
	case 3:
		args = []textFormatArgument{textFormatPointer(unsafe.Pointer(w1)), textFormatPointer(unsafe.Pointer(w2))}
	case 4:
		args = []textFormatArgument{textFormatPointer(unsafe.Pointer(narrow))}
	case 5:
		args = []textFormatArgument{textFormatPointer(unsafe.Pointer(w1)), textFormatPointer(unsafe.Pointer(narrow)), textFormatWord(a), textFormatReal(math.Float64frombits(bits))}
	}
	return textFormatInto(dst, int(count), f, args...)
}
