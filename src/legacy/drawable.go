package legacy

/*
#include "defs.h"
extern int nox_drawable_count;
static void go_nox_drawable_call_sprite_func(void(* fnc)(nox_drawable*, int), nox_drawable* dr, int arg) {
	fnc(dr, arg);
}
*/
import "C"
import (
	"image"
	"unsafe"

	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/legacy/common/ccall"
)

func asDrawable(p *nox_drawable) *client.Drawable {
	return (*client.Drawable)(unsafe.Pointer(p))
}

func AsDrawableP(p unsafe.Pointer) *client.Drawable {
	return (*client.Drawable)(p)
}

type nox_drawable = C.nox_drawable

func nox_xxx_spriteLoadAdd_45A360_drawable(id, x, y int) *nox_drawable {
	return (*nox_drawable)(GetClient().Nox_xxx_spriteLoadAdd_45A360_drawable(id, image.Pt(x, y)).C())
}

func nox_xxx_sprite_45A110_drawable(dr *nox_drawable) {
	GetClient().Cli().Objs.List34Add(asDrawable(dr))
}

func nox_xxx_netSpriteByCodeStatic_45A720(id int) *nox_drawable {
	return (*nox_drawable)(GetClient().Cli().Objs.ByNetCodeStatic(id).C())
}

func nox_xxx_netSpriteByCodeDynamic_45A6F0(id int) *nox_drawable {
	return (*nox_drawable)(GetClient().Cli().Objs.ByNetCodeDynamic(id).C())
}

func nox_xxx_cliRemoveHealthbar_459E30(dr *nox_drawable, v uint8) {
	GetClient().Cli().Objs.RemoveHealthBar(asDrawable(dr), v)
}

func sub_45A670(a1 uint32) {
	GetClient().Sub_45A670(a1)
}

func nox_xxx_spriteTransparentDecay_49B950(dr *nox_drawable, a2 int) {
	GetClient().Cli().Objs.TransparentDecay(asDrawable(dr), a2)
}

func sub_459DD0(dr *nox_drawable, a2 uint8) {
	GetClient().Cli().Objs.MinimapAdd(asDrawable(dr), a2)
}

func nox_xxx_spriteToList_49BC80_drawable(dr *nox_drawable) {
	GetClient().Cli().Objs.List5Add(asDrawable(dr))
}

func nox_xxx_spriteCreate_48E970(typeID int, code uint16, x, y int) *nox_drawable {
	return (*nox_drawable)(GetClient().Nox_xxx_spriteCreate_48E970(typeID, code, x, y).C())
}

func CallDrawFunc(s *client.Drawable, vp *noxrender.Viewport) int {
	return ccall.CallIntPtr2(s.DrawFuncPtr, vp.C(), s.C())
}

func Nox_xxx_spriteGetMB_476F80() *client.Drawable {
	return asDrawable((*nox_drawable)(interactionHoverDrawable))
}

func Nox_xxx_clientGetSpriteAtCursor_476F90() *client.Drawable {
	return asDrawable((*nox_drawable)(interactionUsableCursorDrawable))
}
func Get_dword_5d4594_1096640() *client.Drawable {
	return AsDrawableP(interactionHoverDrawable)
}
func Set_dword_5d4594_1096640(dr *client.Drawable) {
	interactionHoverDrawable = dr.C()
}
func Get_nox_client_spriteUnderCursorXxx_1096644() *client.Drawable {
	return AsDrawableP(interactionUsableCursorDrawable)
}
func Set_nox_client_spriteUnderCursorXxx_1096644(dr *client.Drawable) {
	interactionUsableCursorDrawable = dr.C()
}
func Sub_495B50(fx *client.DrawableFX) {
	combatFXDetach((*combatFX)(fx.C()))
}
func Sub_4523D0(p unsafe.Pointer) {
	audioEventDelete((*audioEvent)(p))
}
func Sub_495FC0(p *client.DrawableFX, dr *client.Drawable) {
	combatFXAttach((*combatFX)(p.C()), dr)
}
func Sub_49C520(dr *client.Drawable) int {
	return bool2int(presentationRayContains(dr))
}
func Sub_45A9B0(a1, a2 *client.Drawable) {
	clientDrawableLoopAudio(a1, a2)
}
func Nox_xxx_unitSpriteCheckAlly_4951F0(id int) bool {
	return combatAllyLookup(uint32(id)) != nil
}
func Nox_xxx_draw_44C650_free_kind(ptr unsafe.Pointer, kind int) {
	spriteDataFreeKind(ptr, kind)
}
