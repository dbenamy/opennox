//go:build porttest

package legacy

import (
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"unsafe"
)

// PortTestGameNotice invokes the production notice dispatcher with fixture-owned data.
func PortTestGameNotice(data []byte) int {
	return clientGameNotice(data)
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
