//go:build porttest

package legacy

/*
#include "server__dbase__objdb.h"
#include "server__object__objutil.h"
*/
import "C"

import (
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

func PortTestUnitNPCName(u *server.Object) string {
	return GoWString((*wchar2_t)(unsafe.Pointer(C.sub_4E39F0_obj_db(asObjectC(u)))))
}
func PortTestUnitItemName(u *server.Object) (string, unsafe.Pointer) {
	p := C.nox_xxx_itemGetName_4E77E0_obj_util(C.int(uintptr(u.CObj())))
	return GoWString((*wchar2_t)(unsafe.Pointer(p))), unsafe.Pointer(p)
}
