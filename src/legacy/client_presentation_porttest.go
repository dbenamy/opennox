//go:build porttest

package legacy

/*
#include "GAME2_2.h"
#include "GAME2_3.h"
void* sub_49BB80(char);
int sub_49C520(nox_drawable*);
extern uint32_t dword_5d4594_1303508;
*/
import "C"
import (
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/client/noxrender"
	"unsafe"
)

func PortTestPresentationWallY(p *[8]byte) int32 {
	return int32(C.sub_476080((*C.uchar)(unsafe.Pointer(p))))
}
func PortTestPresentationDrawableY(dr *client.Drawable) int32 {
	return int32(C.sub_4761B0((*C.nox_drawable)(dr.C())))
}
func PortTestPresentationBake(v *noxrender.Viewport, dr *client.Drawable) {
	C.sub_476AE0((*C.nox_draw_viewport_t)(v.C()), (*C.nox_drawable)(dr.C()))
}
func PortTestPresentationCopy(dst, src unsafe.Pointer, n uint32) {
	C.sub_476D70((*C.uint32_t)(dst), (*C.int)(src), C.uint(n))
}
func PortTestPresentationPhonemeMark(index int)       { C.nox_client_setPhonemeFrame_476E00(C.int(index)) }
func PortTestPresentationPhonemeInit() unsafe.Pointer { return unsafe.Pointer(C.sub_476E20()) }
func PortTestPresentationPhonemeDraw()                { C.sub_476E90() }
func PortTestPresentationShieldInit() bool            { return C.nox_xxx_loadReflSheild_499360() != 0 }
func PortTestPresentationShieldDestroy()              { C.sub_499450() }
func PortTestPresentationShieldDraw(v *noxrender.Viewport, dr *client.Drawable) {
	C.nox_xxx_drawShield_499810((*C.nox_draw_viewport_t)(v.C()), (*C.nox_drawable)(dr.C()))
}
func PortTestPresentationTurnUndead(pos *[2]int16) {
	C.nox_xxx_fxDrawTurnUndead_499880((*C.short)(unsafe.Pointer(pos)))
}
func PortTestPresentationChantStart(id byte) { C.sub_49BB80(C.char(id)) }
func PortTestPresentationChantClear()        { C.sub_49BBB0() }
func PortTestPresentationChantTick()         { C.sub_49BBC0() }
func PortTestPresentationChantOwner() (*uint32, func()) {
	p := (*uint32)(unsafe.Pointer(&C.dword_5d4594_1303508))
	old := *p
	*p = 0
	return p, func() { *p = old }
}
func PortTestPresentationRayAdd(p *[7]byte) {
	C.nox_xxx_clientAddRayEffect_49C160(C.int(uintptr(unsafe.Pointer(p))))
}
func PortTestPresentationRayRemove(p *[7]byte) {
	C.nox_xxx_clientRemoveRayEffect_49C450(C.int(uintptr(unsafe.Pointer(p))))
}
func PortTestPresentationSpareClear() { C.nox_xxx_spriteDeleteSomeList_49C4B0() }
func PortTestPresentationRayClear()   { C.nox_xxx_sprite_49C4F0() }
func PortTestPresentationRayContains(dr *client.Drawable) bool {
	return C.sub_49C520((*C.nox_drawable)(dr.C())) != 0
}
func PortTestPresentationBookReward(kind, id, auto int) {
	C.nox_xxx_bookRewardCli_499CF0((*C.int)(unsafe.Pointer(uintptr(kind))), C.int(id), C.int(auto))
}
func PortTestPresentationBubble(kind, x, y int, z int16, a, b, c, d, e byte, lifetime int) {
	C.sub_499F60(C.int(kind), C.int(x), C.int(y), C.short(z), C.char(a), C.char(b), C.char(c), C.char(d), C.char(e), C.int(lifetime))
}
func PortTestPresentationEquip(op byte, id, mask uint32, mods *[4]byte) {
	C.nox_xxx_clientEquip_49A3D0(C.char(op), C.int(id), C.int(mask), C.int(uintptr(unsafe.Pointer(mods))))
}
