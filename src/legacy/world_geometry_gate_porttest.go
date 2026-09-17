//go:build porttest

package legacy

/*
#include <stdint.h>
char nox_xxx_updateDoor_53AC50(int a1);
static void* geometryDoorUpdate(void){return (void*)nox_xxx_updateDoor_53AC50;}
*/
import "C"
import "unsafe"

func PortTestGeometryDoorUpdate() unsafe.Pointer { return C.geometryDoorUpdate() }
func PortTestGeometryQueues() (reset func(), snapshot func(map[unsafe.Pointer]uint32) [3]uint32, restore func()) {
	old := [3]uint32{collisionAngleHead, collisionActiveHead, collisionActiveTail}
	reset = func() { collisionAngleHead = 0; collisionActiveHead = 0; collisionActiveTail = 0 }
	reset()
	snapshot = func(ids map[unsafe.Pointer]uint32) (out [3]uint32) {
		for i, v := range [3]uint32{collisionAngleHead, collisionActiveHead, collisionActiveTail} {
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
		collisionAngleHead = old[0]
		collisionActiveHead = old[1]
		collisionActiveTail = old[2]
	}
	return
}
