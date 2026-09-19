//go:build porttest

package legacy

/*
#include "GAME4_3.h"
*/
import "C"

import (
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"unsafe"
)

func PortTestUnitActionName(id int32, restricted bool) (string, bool) {
	var p *C.char
	if restricted {
		p = C.sub_5345B0(C.int(id))
	} else {
		p = C.sub_534650(C.int(id))
	}
	return C.GoString(p), p != nil
}
func PortTestUnitActionIndex(name string, restricted bool) int32 {
	p, free := alloc.CString(name)
	defer free()
	if restricted {
		return int32(C.nox_xxx_actionNByNameMB_5345F0((*C.char)(unsafe.Pointer(p))))
	}
	return int32(C.nox_xxx_actionByName_534670((*C.char)(unsafe.Pointer(p))))
}
