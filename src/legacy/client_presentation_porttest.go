//go:build porttest

package legacy

/*
#include "GAME2_3.h"
*/
import "C"
import (
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/client/noxrender"
	"image"
	"unsafe"
)

func PortTestPresentationWallY(p *[8]byte) int32                          { return presentationWallY(p) }
func PortTestPresentationDrawableY(dr *client.Drawable) int32             { return presentationDrawableY(dr) }
func PortTestPresentationBake(v *noxrender.Viewport, dr *client.Drawable) { presentationBake(v, dr) }
func PortTestPresentationCopy(dst, src unsafe.Pointer, n uint32)          { presentationCopy(dst, src, n) }
func PortTestPresentationPhonemeMark(index int)                           { presentationPhonemeMark(index) }
func PortTestPresentationPhonemeInit() unsafe.Pointer                     { return presentationPhonemeInit().C() }
func PortTestPresentationPhonemeDraw()                                    { presentationPhonemeDraw() }
func PortTestPresentationShieldInit() bool                                { return presentationShieldInit() }
func PortTestPresentationShieldDestroy()                                  { presentationShieldDestroy() }
func PortTestPresentationShieldDraw(v *noxrender.Viewport, dr *client.Drawable) {
	presentationShieldDraw(v, dr)
}
func PortTestPresentationTurnUndead(pos *[2]int16) {
	C.nox_xxx_fxDrawTurnUndead_499880((*C.short)(unsafe.Pointer(pos)))
}
func PortTestPresentationChantStart(id byte) { presentationChantStart(id) }
func PortTestPresentationChantClear()        { presentationChantClear() }
func PortTestPresentationChantTick()         { presentationChantTick() }
func PortTestPresentationChantOwner() (*uint32, func()) {
	old := presentationChantTree
	presentationChantTree = nil
	return (*uint32)(unsafe.Pointer(&presentationChantTree)), func() { presentationChantTree = old }
}
func PortTestPresentationRayAdd(p *[7]byte) {
	C.nox_xxx_clientAddRayEffect_49C160(C.int(uintptr(unsafe.Pointer(p))))
}
func PortTestPresentationRayRemove(p *[7]byte) {
	C.nox_xxx_clientRemoveRayEffect_49C450(C.int(uintptr(unsafe.Pointer(p))))
}
func PortTestPresentationSpareClear()                          { presentationTransientClear() }
func PortTestPresentationRayClear()                            { presentationRayClear() }
func PortTestPresentationRayContains(dr *client.Drawable) bool { return presentationRayContains(dr) }
func PortTestPresentationBookReward(kind, id, auto int)        { presentationBookReward(kind, id, auto) }
func PortTestPresentationBubble(kind, x, y int, z int16, a, b, c, d, e byte, lifetime int) {
	presentationBubble(kind, image.Pt(x, y), z, a, b, c, d, e, lifetime)
}
func PortTestPresentationEquip(op byte, id, mask uint32, mods *[4]byte) {
	C.nox_xxx_clientEquip_49A3D0(C.char(op), C.int(id), C.int(mask), C.int(uintptr(unsafe.Pointer(mods))))
}
