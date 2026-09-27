//go:build porttest

package legacy

import (
	"unsafe"

	"github.com/opennox/opennox/v1/server"
)

var portTestModifierKeys [4]byte
var portTestModifierArgs [6]uintptr
var portTestModifierResult int32
var portTestModifierCalls, portTestModifierKind int

func portTestModifierRecord3(kind int, m *server.ModifierEff, a, b *server.Object) {
	portTestModifierArgs[0] = uintptr(unsafe.Pointer(m))
	portTestModifierArgs[1] = uintptr(unsafe.Pointer(a))
	portTestModifierArgs[2] = uintptr(unsafe.Pointer(b))
	portTestModifierCalls++
	portTestModifierKind = kind
}
func init() {
	server.RegisterModifierEffect3(unsafe.Pointer(&portTestModifierKeys[0]), func(m *server.ModifierEff, a, b *server.Object) int32 {
		portTestModifierRecord3(1, m, a, b)
		return portTestModifierResult
	})
	server.RegisterModifierEffect3(unsafe.Pointer(&portTestModifierKeys[1]), func(m *server.ModifierEff, a, b *server.Object) int32 {
		portTestModifierRecord3(2, m, a, b)
		return 0
	})
	server.RegisterModifierEffect5(unsafe.Pointer(&portTestModifierKeys[2]), func(m *server.ModifierEff, a, b, c *server.Object, data unsafe.Pointer) {
		portTestModifierRecord3(3, m, a, b)
		portTestModifierArgs[3], portTestModifierArgs[4] = uintptr(unsafe.Pointer(c)), uintptr(data)
		if data != nil {
			*(*uint32)(data) = ^uint32(portTestModifierResult)
		}
	})
	server.RegisterModifierEffect6(unsafe.Pointer(&portTestModifierKeys[3]), func(m *server.ModifierEff, a, b, c, d *server.Object, data unsafe.Pointer) {
		portTestModifierRecord3(4, m, a, b)
		portTestModifierArgs[3], portTestModifierArgs[4], portTestModifierArgs[5] = uintptr(unsafe.Pointer(c)), uintptr(unsafe.Pointer(d)), uintptr(data)
		if data != nil {
			*(*uint32)(data) = ^uint32(portTestModifierResult)
		}
	})
}
func PortTestModifierObserverReset(value int32) [4]unsafe.Pointer {
	portTestModifierArgs = [6]uintptr{}
	portTestModifierResult, portTestModifierCalls, portTestModifierKind = value, 0, 0
	return [4]unsafe.Pointer{unsafe.Pointer(&portTestModifierKeys[0]), unsafe.Pointer(&portTestModifierKeys[1]), unsafe.Pointer(&portTestModifierKeys[2]), unsafe.Pointer(&portTestModifierKeys[3])}
}
func PortTestModifierObserverSnapshot() ([6]uintptr, int, int) {
	return portTestModifierArgs, portTestModifierCalls, portTestModifierKind
}
func PortTestModifierCall3Result(key unsafe.Pointer, m *server.ModifierEff, a, b *server.Object) int32 {
	return server.CallModifierEffect3Result(key, m, a, b)
}
func PortTestModifierCall3Discard(key unsafe.Pointer, m *server.ModifierEff, a, b *server.Object) {
	server.CallModifierEffect3Discard(key, m, a, b)
}
func PortTestModifierCall5(key unsafe.Pointer, m *server.ModifierEff, a, b, c *server.Object, data unsafe.Pointer) {
	server.CallModifierEffect5(key, m, a, b, c, data)
}
func PortTestModifierCall6(key unsafe.Pointer, m *server.ModifierEff, a, b, c, d *server.Object, data unsafe.Pointer) {
	server.CallModifierEffect6(key, m, a, b, c, d, data)
}
