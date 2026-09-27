//go:build porttest

package legacy

/*
extern void nox_porttest_image_end(void*);
*/
import "C"
import "unsafe"

var porttestImageEnd func(unsafe.Pointer)

//export nox_porttest_image_end
func nox_porttest_image_end(ref unsafe.Pointer) { porttestImageEnd(ref) }

// PortTestObserveImageEnd supplies a foreign ABI observer, not an animation implementation.
func PortTestObserveImageEnd(fn func(unsafe.Pointer)) (unsafe.Pointer, func()) {
	old := porttestImageEnd
	porttestImageEnd = fn
	return C.nox_porttest_image_end, func() { porttestImageEnd = old }
}
