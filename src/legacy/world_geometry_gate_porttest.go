//go:build porttest

package legacy

/*
#include <stdint.h>
extern uint32_t dword_5d4594_2491552;
extern uint32_t dword_5d4594_2488604, dword_5d4594_2488608;
char nox_xxx_updateDoor_53AC50(int a1);
static void* geometryDoorUpdate(void){return (void*)nox_xxx_updateDoor_53AC50;}
*/
import "C"
import "unsafe"

func PortTestGeometryDoorUpdate() unsafe.Pointer { return C.geometryDoorUpdate() }
func PortTestGeometryQueues() (reset func(), snapshot func(map[unsafe.Pointer]uint32) [3]uint32, restore func()) {
	old := [3]C.uint32_t{C.dword_5d4594_2491552, C.dword_5d4594_2488604, C.dword_5d4594_2488608}
	reset = func() { C.dword_5d4594_2491552 = 0; C.dword_5d4594_2488604 = 0; C.dword_5d4594_2488608 = 0 }
	reset()
	snapshot = func(ids map[unsafe.Pointer]uint32) (out [3]uint32) {
		for i, v := range [3]C.uint32_t{C.dword_5d4594_2491552, C.dword_5d4594_2488604, C.dword_5d4594_2488608} {
			if v != 0 {
				var ok bool
				out[i], ok = ids[unsafe.Pointer(uintptr(v))]
				if !ok {
					panic("geometry queue outside fixture")
				}
			}
		}
		return
	}
	restore = func() {
		C.dword_5d4594_2491552 = old[0]
		C.dword_5d4594_2488604 = old[1]
		C.dword_5d4594_2488608 = old[2]
	}
	return
}
