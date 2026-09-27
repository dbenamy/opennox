//go:build porttest

package legacy

import "unsafe"

var porttestImageEnd func(unsafe.Pointer)
var porttestImageEndKey byte

func init() {
	imageAnimationEndCallbacks[unsafe.Pointer(&porttestImageEndKey)] = func(ref *ImageRef) { porttestImageEnd(ref.C()) }
}

// PortTestObserveImageEnd supplies an observer, not an animation implementation.
func PortTestObserveImageEnd(fn func(unsafe.Pointer)) (unsafe.Pointer, func()) {
	old := porttestImageEnd
	porttestImageEnd = fn
	return unsafe.Pointer(&porttestImageEndKey), func() { porttestImageEnd = old }
}
