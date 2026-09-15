package legacy

/*
#include "defs.h"
#include "memfile.h"
*/
import "C"
import (
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/client/noxrender"
	"unsafe"
)

//export nox_thing_animate_draw
func nox_thing_animate_draw(vp *C.uint32_t, dr *C.nox_drawable) C.int {
	return C.int(spriteAnimateDraw((*noxrender.Viewport)(unsafe.Pointer(vp)), (*client.Drawable)(unsafe.Pointer(dr))))
}

//export nox_thing_cond_animate_draw
func nox_thing_cond_animate_draw(vp *C.uint32_t, dr *C.nox_drawable) C.int {
	return C.int(spriteConditionalDraw((*noxrender.Viewport)(unsafe.Pointer(vp)), (*client.Drawable)(unsafe.Pointer(dr))))
}

//export nox_thing_static_draw
func nox_thing_static_draw(vp *C.uint32_t, dr *C.nox_drawable) C.int {
	return C.int(spriteStaticDraw((*noxrender.Viewport)(unsafe.Pointer(vp)), (*client.Drawable)(unsafe.Pointer(dr))))
}

//export nox_thing_static_random_draw
func nox_thing_static_random_draw(vp *C.uint32_t, dr *C.nox_drawable) C.int {
	return C.int(spriteSlaveDraw((*noxrender.Viewport)(unsafe.Pointer(vp)), (*client.Drawable)(unsafe.Pointer(dr))))
}

//export nox_thing_slave_draw
func nox_thing_slave_draw(vp *C.int, dr *C.nox_drawable) C.int {
	return C.int(spriteSlaveDraw((*noxrender.Viewport)(unsafe.Pointer(vp)), (*client.Drawable)(unsafe.Pointer(dr))))
}

//export nox_thing_boulder_draw
func nox_thing_boulder_draw(vp *C.int, dr *C.nox_drawable) C.int {
	return C.int(spriteBoulderDraw((*noxrender.Viewport)(unsafe.Pointer(vp)), (*client.Drawable)(unsafe.Pointer(dr))))
}

//export nox_things_animate_draw_parse
func nox_things_animate_draw_parse(obj *C.nox_thing, f *C.nox_memfile, attr *C.char) C.bool {
	return C.bool(spriteParseAnimate((*client.ObjectType)(unsafe.Pointer(obj)), asMemfile(f), unsafe.Slice((*byte)(unsafe.Pointer(attr)), 256)))
}

//export nox_things_cond_animate_draw_parse
func nox_things_cond_animate_draw_parse(obj *C.nox_thing, f *C.nox_memfile, attr *C.char) C.bool {
	return C.bool(spriteParseConditional((*client.ObjectType)(unsafe.Pointer(obj)), asMemfile(f), unsafe.Slice((*byte)(unsafe.Pointer(attr)), 256)))
}

//export nox_things_static_draw_parse
func nox_things_static_draw_parse(obj *C.nox_thing, f *C.nox_memfile, attr *C.char) C.bool {
	return C.bool(spriteParseStatic((*client.ObjectType)(unsafe.Pointer(obj)), asMemfile(f), unsafe.Slice((*byte)(unsafe.Pointer(attr)), 256)))
}

//export nox_things_static_random_draw_parse
func nox_things_static_random_draw_parse(obj *C.nox_thing, f *C.nox_memfile, attr *C.char) C.bool {
	return C.bool(spriteParseRandom((*client.ObjectType)(unsafe.Pointer(obj)), asMemfile(f), unsafe.Slice((*byte)(unsafe.Pointer(attr)), 256), false))
}

//export nox_things_slave_draw_parse
func nox_things_slave_draw_parse(obj *C.nox_thing, f *C.nox_memfile, attr *C.char) C.bool {
	return C.bool(spriteParseRandom((*client.ObjectType)(unsafe.Pointer(obj)), asMemfile(f), unsafe.Slice((*byte)(unsafe.Pointer(attr)), 256), true))
}

//export nox_things_animate_state_draw_parse
func nox_things_animate_state_draw_parse(obj *C.nox_thing, f *C.nox_memfile, attr *C.char) C.bool {
	return C.bool(spriteParseState((*client.ObjectType)(unsafe.Pointer(obj)), asMemfile(f)))
}
