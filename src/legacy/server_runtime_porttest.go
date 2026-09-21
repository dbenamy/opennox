//go:build porttest

package legacy

import (
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

func PortTestRuntimeRejectedClear(head unsafe.Pointer) uintptr { return runtimeRejectedClear(head) }
func PortTestRuntimeEquipmentMask(armor bool, name string) uint32 {
	p, free := alloc.CString16(name)
	defer free()
	return runtimeEquipmentMask(armor, p)
}
func PortTestRuntimeEquipmentLabel(armor bool, mask uint32) (string, bool) {
	p := runtimeEquipmentLabel(armor, mask)
	return alloc.GoString16(p), p != nil
}
func PortTestRuntimeEquipmentLoad(armor bool)      { runtimeEquipmentLoad(armor) }
func PortTestRuntimeModifierIcon(mask byte) uint32 { return runtimeModifierIcon(mask) }
func PortTestRuntimeModifierLabel(mask byte) (string, bool) {
	p := runtimeModifierLabel(mask)
	return alloc.GoString16(p), p != nil
}
func PortTestRuntimeArmorConductivity(u *server.Object) float64 { return runtimeArmorConductivity(u) }
