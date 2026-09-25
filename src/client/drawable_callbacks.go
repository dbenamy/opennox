package client

import (
	"runtime"
	"unsafe"

	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/legacy/common/ccall"
)

var drawableDrawCallbacksGo = make(map[unsafe.Pointer]func(*noxrender.Viewport, *Drawable) int32)

// RegisterDrawableDrawCallbackGo binds a stable identity during initialization.
// Drawable and object-type records retain their original pointer-sized fields.
func RegisterDrawableDrawCallbackGo(key unsafe.Pointer, fn func(*noxrender.Viewport, *Drawable) int32) {
	if key == nil || fn == nil {
		panic("invalid drawable callback")
	}
	if _, ok := drawableDrawCallbacksGo[key]; ok {
		panic("drawable callback already registered")
	}
	drawableDrawCallbacksGo[key] = fn
}

// CallDrawableDrawResult preserves the complete signed C-int result. Callers
// require a configured callback; foreign addresses keep the original C ABI.
func CallDrawableDrawResult(key unsafe.Pointer, vp *noxrender.Viewport, dr *Drawable) int32 {
	var result int32
	if fn := drawableDrawCallbacksGo[key]; fn != nil {
		result = fn(vp, dr)
	} else {
		result = int32(ccall.CallIntPtr2(key, vp.C(), dr.C()))
	}
	runtime.KeepAlive(vp)
	runtime.KeepAlive(dr)
	return result
}

// CallDrawableDrawDiscard preserves the void convention used for shield draws.
func CallDrawableDrawDiscard(key unsafe.Pointer, vp *noxrender.Viewport, dr *Drawable) {
	if fn := drawableDrawCallbacksGo[key]; fn != nil {
		fn(vp, dr)
	} else {
		ccall.CallVoidPtr2(key, vp.C(), dr.C())
	}
	runtime.KeepAlive(vp)
	runtime.KeepAlive(dr)
}
