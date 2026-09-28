//go:build porttest

package legacy

import (
	"unsafe"

	"github.com/opennox/opennox/v1/server"
)

type PortTestDamageForwardState struct {
	Count    int
	Callback int
	Words    [5]uint32
}

var portTestDamageForwardKeys [2]byte
var portTestDamageForwardState PortTestDamageForwardState
var portTestDamageForwardResult int32

func init() {
	for index := range portTestDamageForwardKeys {
		server.PortTestRegisterDamageCallback(unsafe.Pointer(&portTestDamageForwardKeys[index]), func(a, b, c *server.Object, d, e int32) int32 {
			portTestDamageForwardState.Count++
			portTestDamageForwardState.Callback = index + 1
			portTestDamageForwardState.Words = [5]uint32{uint32(uintptr(unsafe.Pointer(a))), uint32(uintptr(unsafe.Pointer(b))), uint32(uintptr(unsafe.Pointer(c))), uint32(d), uint32(e)}
			return portTestDamageForwardResult
		})
	}
}
func PortTestDamageForwardCallback(which int) unsafe.Pointer {
	switch which {
	case 0, 1:
		return unsafe.Pointer(&portTestDamageForwardKeys[which])
	default:
		panic("unknown damage forwarding callback")
	}
}
func PortTestDamageForwardReset(result int32) {
	portTestDamageForwardState = PortTestDamageForwardState{}
	portTestDamageForwardResult = result
}
func PortTestDamageForwardSnapshot() PortTestDamageForwardState { return portTestDamageForwardState }
