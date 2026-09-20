//go:build porttest

package legacy

/*
#include "client__network__inform.h"
*/
import "C"
import "unsafe"

// PortTestGameNotice invokes the original notice dispatcher with fixture-owned data.
func PortTestGameNotice(data []byte) int {
	return int(C.nox_client_handlePacketInform_4C9BF0(C.int(uintptr(unsafe.Pointer(&data[0])))))
}
