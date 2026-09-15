//go:build porttest

package legacy

/*
#include "client__draw__animdraw.h"
#include "client__draw__canidraw.h"
#include "client__draw__staticdraw.h"
#include "client__draw__slavedraw.h"
void nox_xxx_draw_44C650_free_kind(void*, int);
void* nox_xxx_draw_44C780(int);
*/
import "C"

import (
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/internal/binfile"
	"unsafe"
)

func PortTestSpriteParse(op int, obj *client.ObjectType, mf *binfile.MemFile, attr unsafe.Pointer, vector unsafe.Pointer) int {
	o := (*C.nox_thing)(obj.C())
	f := (*C.nox_memfile)(mf.C())
	a := (*C.char)(attr)
	switch op {
	case 0:
		if C.nox_things_animate_draw_parse(o, f, a) {
			return 1
		}
		return 0
	case 1:
		if C.nox_things_cond_animate_draw_parse(o, f, a) {
			return 1
		}
		return 0
	case 2:
		if C.nox_things_static_draw_parse(o, f, a) {
			return 1
		}
		return 0
	case 3:
		if C.nox_things_static_random_draw_parse(o, f, a) {
			return 1
		}
		return 0
	case 4:
		return spriteVectorHeader((*client.AnimationVector)(vector), mf)
	case 6:
		return spriteVectorFrames((*client.AnimationVector)(vector), mf)
	case 7:
		if C.nox_things_animate_state_draw_parse(o, f, a) {
			return 1
		}
		return 0
	case 9:
		if C.nox_things_slave_draw_parse(o, f, a) {
			return 1
		}
		return 0
	case 8:
		return int(client.ParseAnimKind(GoString(a)))
	}
	panic("unknown sprite parser")
}
func PortTestSpriteFree(data unsafe.Pointer, kind int) {
	C.nox_xxx_draw_44C650_free_kind(data, C.int(kind))
}

func PortTestSpriteFreeVectorFrames(data unsafe.Pointer) {
	C.nox_xxx_draw_44C780(C.int(uintptr(data) + 4))
}
