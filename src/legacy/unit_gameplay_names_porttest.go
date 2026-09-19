//go:build porttest

package legacy

import (
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

func PortTestUnitNPCName(u *server.Object) string { return alloc.GoString16(unitNPCName(u)) }
func PortTestUnitItemName(u *server.Object) (string, unsafe.Pointer) {
	p := unitItemName(u)
	return alloc.GoString16(p), unsafe.Pointer(p)
}
