//go:build porttest

package legacy

import (
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/client/noxrender"
	"unsafe"
)

var portTestAdapterKeys [2]byte
var portTestAdapterWords [4]uint32
var portTestAdapterResult int32

func init() {
	for i := range portTestAdapterKeys {
		client.RegisterDrawableDrawCallbackGo(unsafe.Pointer(&portTestAdapterKeys[i]), func(vp *noxrender.Viewport, dr *client.Drawable) int32 {
			portTestAdapterWords[0] = uint32(uintptr(unsafe.Pointer(vp)))
			portTestAdapterWords[1] = uint32(uintptr(unsafe.Pointer(dr)))
			portTestAdapterWords[2] = uint32(i + 1)
			portTestAdapterWords[3]++
			return portTestAdapterResult
		})
	}
}
func PortTestAdapterDrawCallback(which int) unsafe.Pointer {
	if int32(which) == 1 {
		return unsafe.Pointer(&portTestAdapterKeys[0])
	}
	return unsafe.Pointer(&portTestAdapterKeys[1])
}
func PortTestAdapterDrawReset(result int) {
	portTestAdapterWords = [4]uint32{}
	portTestAdapterResult = int32(result)
}
func PortTestAdapterDrawValue(field int) uintptr {
	switch int32(field) {
	case 0, 1, 2, 3:
		return uintptr(portTestAdapterWords[int32(field)])
	}
	return 0
}
