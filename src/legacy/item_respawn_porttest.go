//go:build porttest

package legacy

/*
#include "GAME3_3.h"
extern uint32_t dword_5d4594_1568024;
extern uint32_t dword_5d4594_1568028;
extern uint32_t nox_xxx_respawnAllow_587000_205200;
void sub_4EC720();
*/
import "C"
import (
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

func PortTestItemRespawnGlobals() (*uint32, func()) {
	oldCache := C.dword_5d4594_1568028
	return (*uint32)(unsafe.Pointer(&C.nox_xxx_respawnAllow_587000_205200)), func() { C.dword_5d4594_1568028 = oldCache }
}
func PortTestItemRespawn(op string, u *server.Object) uintptr {
	switch op {
	case "add":
		return uintptr(unsafe.Pointer(C.nox_xxx_respawnAdd_4EC5E0(asObjectC(u))))
	case "remove":
		C.sub_4EC6A0(C.int(uintptr(u.CObj())))
	case "reset":
		C.sub_4EC5B0()
	case "tick":
		C.sub_4EC720()
	}
	return 0
}
func PortTestItemRespawnRecords() []*[15]uint32 {
	var out []*[15]uint32
	for p := uint32(C.dword_5d4594_1568024); p != 0; {
		if len(out) > 1024 {
			panic("respawn list cycle")
		}
		r := (*[15]uint32)(unsafe.Pointer(uintptr(p)))
		out = append(out, r)
		p = r[13]
	}
	return out
}
func PortTestItemRespawnSameTeam(a, b *server.Object) bool {
	return C.nox_xxx_unitsHaveSameTeam_4EC520(asObjectC(a), asObjectC(b)) != 0
}
func PortTestItemRespawnDropCrown(u *server.Object, stamp uint32) {
	C.sub_4ED050(C.int(uintptr(u.CObj())), C.int(stamp))
}
