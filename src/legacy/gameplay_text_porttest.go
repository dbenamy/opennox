//go:build porttest

package legacy

import (
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

func gameplayTextInvoke(op int, args [5]uint32) uint32 {
	ptr := func(i int) unsafe.Pointer { return unsafe.Pointer(uintptr(args[i])) }
	obj := func(i int) *server.Object { return (*server.Object)(ptr(i)) }
	switch op {
	case 0, 1:
		if op == 0 && (obj(0) == nil || uint32(obj(0).ObjClass)&4 == 0) {
			return uint32(textFormatLine(obj(0), nil))
		}
		a := gameplayTextFixtureArguments((*uint16)(ptr(1)), args[2:4])
		if op == 0 {
			return uint32(textFormatLine(obj(0), (*uint16)(ptr(1)), a...))
		}
		return uint32(textFormatAll(byte(args[0]), (*uint16)(ptr(1)), a...))
	case 2:
		return uint32(gameplayTextInformation(int(args[0]), int(args[1]), ptr(2)))
	case 3:
		return uint32(gameplayTextInformationAll(int(args[0]), ptr(1)))
	case 4:
		gameplayTextPrivate(obj(0), (*byte)(ptr(1)), byte(args[2]))
		return 0
	case 5:
		return uint32(gameplayTextPrivateAll((*byte)(ptr(0))))
	case 6:
		return uint32(uintptr(GetServer().S().Players.FirstUnit().CObj()))
	case 7:
		return uint32(uintptr(GetServer().S().Players.NextUnit(obj(0)).CObj()))
	case 8:
		return uint32(bool2int(gameplayTextByteEncoding(gameplayTextUnits((*uint16)(ptr(0))))))
	case 9:
		return uint32(gameplayTextChat(obj(0), gameplayTextUnits((*uint16)(ptr(1))), uint16(args[2])))
	default:
		panic("text operation")
	}
}

// Decode the fixture's two raw ABI words before creating Go pointer values.
// Integer patterns must never be stored in pointer slots for GC to scan.
func gameplayTextFixtureArguments(format *uint16, words []uint32) []textFormatArgument {
	out := make([]textFormatArgument, len(words))
	for i, w := range words {
		out[i] = textFormatWord(w)
	}
	units := gameplayTextUnits(format)
	arg := 0
	for i := 0; i < len(units); i++ {
		if units[i] != '%' {
			continue
		}
		i++
		for i < len(units) && (units[i] == '+' || units[i] == '-' || units[i] == '.' || units[i] >= '0' && units[i] <= '9') {
			i++
		}
		if i == len(units) {
			panic("incomplete fixture format")
		}
		if units[i] == '%' || units[i] == '!' {
			continue
		}
		if arg >= len(words) {
			panic("fixture argument limit")
		}
		switch units[i] {
		case 's', 'S':
			out[arg] = textFormatPointer(unsafe.Pointer(uintptr(words[arg])))
		case 'c', 'd', 'i', 'u', 'x', 'X', 'o':
		default:
			panic("unsupported fixture argument type")
		}
		arg++
	}
	return out
}
