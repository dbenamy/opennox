//go:build porttest

package legacy

/*
#include "GAME3.h"
extern uint32_t dword_587000_171388;
*/
import "C"

func PortTestCharacterDefaultOwner(value uint32) func() {
	old := C.dword_587000_171388
	C.dword_587000_171388 = C.uint32_t(value)
	return func() { C.dword_587000_171388 = old }
}
func PortTestCharacterSetup() uint32 { return uint32(C.sub_4A5E90()) }
