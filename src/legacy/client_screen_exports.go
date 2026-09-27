package legacy

/*
#include "defs.h"
*/
import "C"
import (
	"github.com/opennox/opennox/v1/client/noxrender"
	"unsafe"
)

//export nox_client_screenParticleDraw_489700
func nox_client_screenParticleDraw_489700(vp unsafe.Pointer, p *C.nox_screenParticle) C.int {
	return C.int(screenParticleDraw((*noxrender.Viewport)(vp), (*Nox_screenParticle)(unsafe.Pointer(p))))
}
