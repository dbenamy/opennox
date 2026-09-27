package gui

import (
	"unsafe"

	"github.com/opennox/opennox/v1/legacy/common/ccall"
)

var animationCallbacksGo = make(map[unsafe.Pointer]func() int)

// RegisterAnimationCallbackGo binds a stable identity during initialization.
func RegisterAnimationCallbackGo(key unsafe.Pointer, fn func() int) {
	if key == nil || fn == nil {
		panic("invalid animation callback")
	}
	if _, ok := animationCallbacksGo[key]; ok {
		panic("animation callback already registered")
	}
	animationCallbacksGo[key] = fn
}

// CallAnimationCallback preserves the original raw-C fallback, including nil.
func CallAnimationCallback(key unsafe.Pointer) int {
	if fn := animationCallbacksGo[key]; fn != nil {
		return fn()
	}
	return ccall.CallIntVoid(key)
}
