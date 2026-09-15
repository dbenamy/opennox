//go:build porttest

package legacy

/*
#include "GAME2_2.h"
extern uint32_t dword_5d4594_1193348;
extern uint32_t dword_5d4594_1193352;
*/
import "C"
import "github.com/opennox/opennox/v1/common/memmap"

type PortTestEntryEnvironment struct {
	context, active uint32
	blink           byte
}

func PortTestNewEntryEnvironment() *PortTestEntryEnvironment {
	e := &PortTestEntryEnvironment{uint32(C.dword_5d4594_1193348), uint32(C.dword_5d4594_1193352), *memmap.PtrUint8(0x5D4594, 1193344)}
	C.dword_5d4594_1193348 = 0
	C.dword_5d4594_1193352 = 0
	*memmap.PtrUint8(0x5D4594, 1193344) = 0
	return e
}
func (e *PortTestEntryEnvironment) Context(on bool) {
	C.sub_488BA0()
	if on {
		C.sub_488B60()
	}
}
func (e *PortTestEntryEnvironment) ClearActive() { C.dword_5d4594_1193352 = 0 }
func (e *PortTestEntryEnvironment) State() [3]uint32 {
	context := uint32(0)
	if C.dword_5d4594_1193348 != 0 {
		context = 1
	}
	return [3]uint32{context, uint32(C.dword_5d4594_1193352), uint32(*memmap.PtrUint8(0x5D4594, 1193344))}
}
func (e *PortTestEntryEnvironment) Blink(v byte) { *memmap.PtrUint8(0x5D4594, 1193344) = v }
func (e *PortTestEntryEnvironment) Restore() {
	C.sub_488BA0()
	C.dword_5d4594_1193348 = C.uint32_t(e.context)
	C.dword_5d4594_1193352 = C.uint32_t(e.active)
	*memmap.PtrUint8(0x5D4594, 1193344) = e.blink
}
