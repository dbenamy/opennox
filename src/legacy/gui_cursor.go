package legacy

import "C"
import "github.com/opennox/opennox/v1/client/gui"

//export nox_client_setCursorType_477610
func nox_client_setCursorType_477610(v int) int {
	GetClient().Nox_client_setCursorType(gui.Cursor(v))
	return v
}
