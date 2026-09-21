//go:build porttest

package legacy

/*
#include "GAME1.h"
#include "GAME5_2.h"
#include "common__object__modifier.h"
#include "common__object__armrlook.h"
#include "common__object__weaplook.h"
*/
import "C"

import (
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

func PortTestRuntimeRejectedClear(head unsafe.Pointer) uintptr {
	return uintptr(unsafe.Pointer(C.sub_57ADF0((*C.int)(head))))
}
func PortTestRuntimeEquipmentMask(armor bool, name string) uint32 {
	p, free := alloc.CString16(name)
	defer free()
	if armor {
		return uint32(C.sub_415DA0((*C.wchar2_t)(p)))
	}
	return uint32(C.sub_415960((*C.wchar2_t)(p)))
}
func PortTestRuntimeEquipmentLabel(armor bool, mask uint32) (string, bool) {
	var value C.int
	if armor {
		value = C.sub_415E80(C.int(mask))
	} else {
		value = C.sub_4159F0(C.int(mask))
	}
	p := (*uint16)(unsafe.Pointer(uintptr(uint32(value))))
	return alloc.GoString16(p), p != nil
}
func PortTestRuntimeEquipmentLoad(armor bool) {
	if armor {
		C.nox_xxx_loadLook_415D50()
	} else {
		C.nox_xxx_loadModifyers_4158C0()
	}
}
func PortTestRuntimeModifierIcon(mask byte) uint32 { return uint32(C.sub_413420(C.char(mask))) }
func PortTestRuntimeModifierLabel(mask byte) (string, bool) {
	p := (*uint16)(unsafe.Pointer(C.sub_413480(C.char(mask))))
	return alloc.GoString16(p), p != nil
}
func PortTestRuntimeArmorConductivity(u *server.Object) float64 {
	return float64(C.sub_415BD0(C.int(uintptr(u.CObj()))))
}
