//go:build porttest

package legacy

/*
#include "client__network__inform.h"
*/
import "C"
import (
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"unsafe"
)

// PortTestGameNotice invokes the original notice dispatcher with fixture-owned data.
func PortTestGameNotice(data []byte) int {
	return int(C.nox_client_handlePacketInform_4C9BF0(C.int(uintptr(unsafe.Pointer(&data[0])))))
}

// PortTestGameSecretWall owns the actual secret-list head for a borrowed wall.
// The dispatcher still resolves the wall through the production lookup.
func PortTestGameSecretWall(w unsafe.Pointer) func() {
	old := worldSecretHead
	node, free := alloc.Make([]uint32{}, 8)
	node[3] = uint32(uintptr(w))
	worldSecretHead = nil
	worldSecretInsert(unsafe.Pointer(&node[0]))
	return func() { worldSecretHead = old; free() }
}
