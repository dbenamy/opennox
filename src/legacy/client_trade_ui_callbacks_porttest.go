//go:build porttest

package legacy

import "unsafe"

var portTestTradeObserve func([7]uint32)
var portTestTradeCallbackKeys [2]byte

func init() {
	for i := range portTestTradeCallbackKeys {
		uiAmountCallbacks[unsafe.Pointer(&portTestTradeCallbackKeys[i])] = func(a1, a2, a3, a4, a5 uintptr) {
			pos := (*[2]int32)(unsafe.Pointer(a1))
			portTestTradeObserve([7]uint32{uint32(i), uint32(pos[0]), uint32(pos[1]), uint32(a2), uint32(a3), uint32(a4), uint32(a5)})
		}
	}
}

func PortTestTradeUIObserve(fn func([7]uint32)) (unsafe.Pointer, unsafe.Pointer, func()) {
	old := portTestTradeObserve
	portTestTradeObserve = fn
	return unsafe.Pointer(&portTestTradeCallbackKeys[0]), unsafe.Pointer(&portTestTradeCallbackKeys[1]), func() { portTestTradeObserve = old }
}
