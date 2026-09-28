//go:build porttest

package legacy

import (
	"unsafe"

	"github.com/opennox/opennox/v1/server"
)

var portTestWorldDamageKey byte
var portTestWorldDamageWords [6]uint32

func init() {
	server.PortTestRegisterDamageCallback(unsafe.Pointer(&portTestWorldDamageKey), func(target, owner, source *server.Object, amount, kind int32) int32 {
		portTestWorldDamageWords[0]++
		portTestWorldDamageWords[1] = uint32(uintptr(unsafe.Pointer(target)))
		portTestWorldDamageWords[2] = uint32(uintptr(unsafe.Pointer(owner)))
		portTestWorldDamageWords[3] = uint32(uintptr(unsafe.Pointer(source)))
		portTestWorldDamageWords[4], portTestWorldDamageWords[5] = uint32(amount), uint32(kind)
		return 1
	})
}

// Observe damage arguments independently of the damage implementation.
func PortTestWorldDamageObserver() (unsafe.Pointer, *[6]uint32, func()) {
	old := portTestWorldDamageWords
	portTestWorldDamageWords = [6]uint32{}
	return unsafe.Pointer(&portTestWorldDamageKey), &portTestWorldDamageWords, func() { portTestWorldDamageWords = old }
}
