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
	"github.com/opennox/opennox/v1/common/ntype"
	"github.com/opennox/opennox/v1/legacy/common/ccall"
)

func asDrawable(p *nox_drawable) *client.Drawable {
	return (*client.Drawable)(unsafe.Pointer(p))
}

func AsDrawableP(p unsafe.Pointer) *client.Drawable {
	return (*client.Drawable)(p)
}

type nox_drawable = C.nox_drawable

//export nox_xxx_sprite_49AA00_drawable
func nox_xxx_sprite_49AA00_drawable(d *nox_drawable) {
	GetClient().Cli().Objs.AddIndex2D(asDrawable(d))
}

//export nox_xxx_forEachSprite_49AB00
func nox_xxx_forEachSprite_49AB00(a1 *C.int4, cfnc unsafe.Pointer, data int) {
	if cfnc == nil {
		return
	}
	rect := image.Rect(int(a1.field_0), int(a1.field_4), int(a1.field_8), int(a1.field_C))
	GetClient().Cli().Objs.EachInRect(rect, func(dr *client.Drawable) {
		C.go_nox_drawable_call_sprite_func((*[0]byte)(cfnc), (*nox_drawable)(dr.C()), C.int(data))
	})
}

//export nox_drawable_find_49ABF0
func nox_drawable_find_49ABF0(pt *C.nox_point, r int) *nox_drawable {
	return (*nox_drawable)(GetClient().Nox_drawable_find(image.Point{X: int(pt.x), Y: int(pt.y)}, r).C())
}

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

//export nox_xxx_spriteDelete_45A4B0
func nox_xxx_spriteDelete_45A4B0(dr *nox_drawable) int {
	return GetClient().Nox_xxx_spriteDelete_45A4B0(asDrawable(dr))
}

//export nox_new_drawable_for_thing
func nox_new_drawable_for_thing(i int) *nox_drawable {
	return (*nox_drawable)(GetClient().Nox_new_drawable_for_thing(i).C())
}

//export nox_xxx_cliDestroyObj_45A9A0
func nox_xxx_cliDestroyObj_45A9A0(dr *nox_drawable) {
	GetClient().Cli().Nox_xxx_cliDestroyObj_45A9A0(asDrawable(dr))
}

//export nox_xxx_spriteToSightDestroyList_49BAB0_drawable
func nox_xxx_spriteToSightDestroyList_49BAB0_drawable(dr *nox_drawable) {
	GetClient().Cli().Objs.List6Add(asDrawable(dr))
}

//export sub_45A060
func sub_45A060() *nox_drawable {
	return (*nox_drawable)(GetClient().Cli().Objs.FirstList1().C())
}

//export nox_xxx_cliFirstMinimapObj_459EB0
func nox_xxx_cliFirstMinimapObj_459EB0() *nox_drawable {
	return (*nox_drawable)(GetClient().Cli().Objs.FirstMinimapList().C())
}

//export nox_xxx_cliGetSpritePlayer_45A000
func nox_xxx_cliGetSpritePlayer_45A000() *nox_drawable {
	return (*nox_drawable)(GetClient().Cli().Objs.FirstPlayerList().C())
}

//export sub_49BCD0
func sub_49BCD0(dr *nox_drawable) {
	GetClient().Cli().Objs.List5Delete(asDrawable(dr))
}

//export nox_xxx_client_4984B0_drawable
func nox_xxx_client_4984B0_drawable(dr *nox_drawable) int {
	return bool2int(GetClient().Nox_xxx_client_4984B0_drawable(asDrawable(dr)))
}

//export sub_4992B0
func sub_4992B0(a, b int) int {
	return GetClient().Cli().Sight.Sub_4992B0(a, b)
}

//export sub_498C20
func sub_498C20(a, b *C.nox_point, a3 int32) int32 {
	p1, p2 := (*ntype.Point32)(unsafe.Pointer(a)), (*ntype.Point32)(unsafe.Pointer(b))
	var v1, v2 *image.Point
	if p1 != nil {
		v := p1.Point()
		v1 = &v
	}
	if p2 != nil {
		v := p2.Point()
		v2 = &v
	}
	return GetClient().Cli().Sight.Sub_498C20(v1, v2, a3)
}

//export sub_499290
func sub_499290(a1 int) C.nox_point {
	p := GetClient().Cli().Sight.Sub_499290(a1)
	return C.nox_point{x: C.int(p.X), y: C.int(p.Y)}
}

//export nox_xxx_updateSpritePosition_49AA90
func nox_xxx_updateSpritePosition_49AA90(dr *nox_drawable, x, y int) {
	GetClient().Cli().Nox_xxx_updateSpritePosition_49AA90(asDrawable(dr), x, y)
}

func nox_xxx_spriteCreate_48E970(typeID int, code uint16, x, y int) *nox_drawable {
	return (*nox_drawable)(GetClient().Nox_xxx_spriteCreate_48E970(typeID, code, x, y).C())
}

//export nox_xxx_spriteLoadError_4356E0
func nox_xxx_spriteLoadError_4356E0() {
	GetClient().Cli().Nox_xxx_spriteLoadError_4356E0()
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
