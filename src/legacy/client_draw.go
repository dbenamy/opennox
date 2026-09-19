//go:build !server

package legacy

/*
#include "defs.h"
#include "GAME1.h"
#include "GAME1_1.h"
#include "GAME1_2.h"
#include "GAME1_3.h"
#include "GAME2.h"
#include "GAME2_1.h"
#include "GAME2_2.h"
#include "GAME2_3.h"
#include "GAME3_1.h"
#include "client__draw__glowdraw.h"
#include "client__gui__guiggovr.h"
void  nox_xxx_cliLight16_469140(nox_drawable* dr, nox_draw_viewport_t* vp);
void nox_xxx_clientDrawAll_436100_draw_A();
void nox_xxx_clientDrawAll_436100_draw_B();
void nox_xxx_drawAllMB_475810_draw_A(nox_draw_viewport_t* vp);
void nox_xxx_drawAllMB_475810_draw_C(nox_draw_viewport_t* vp, int v36, int v7);
int sub_436F50();
*/
import "C"
import (
	"unsafe"

	"github.com/opennox/libs/types"

	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/server"
)

func Nox_xxx_clientDrawAll_436100_draw_A() {
	C.nox_xxx_clientDrawAll_436100_draw_A()
}

func Nox_xxx_drawMinimapAndLines_4738E0() {
	minimapDrawAndMessages()
}

func Nox_xxx_clientDrawAll_436100_draw_B() {
	C.nox_xxx_clientDrawAll_436100_draw_B()
}

func Sub_436F50() {
	C.sub_436F50()
}

func Sub_437100() {
	C.sub_437100()
}

func Sub_470DE0() {
	uiMeterHeartbeat()
}

func Nox_xxx_clientEnumHover_476FA0() {
	C.nox_xxx_clientEnumHover_476FA0()
}

func Nox_xxx_spriteDeleteSomeList_49C4B0() {
	presentationTransientClear()
}

func Sub_49BBC0() {
	presentationChantTick()
}

func Nox_xxx_polygonDrawColor_421B80() {
	mapPolygonColor()
}

func Nox_xxx_cliToggleObsWindow_4357A0() {
	C.nox_xxx_cliToggleObsWindow_4357A0()
}

func Nox_xxx_motd_4467F0() {
	C.nox_xxx_motd_4467F0()
}

func Sub_42EBA0() int {
	return int(C.sub_42EBA0())
}

func Sub_49B6E0() {
	C.sub_49B6E0()
}

func Get_nox_thing_glow_orb_draw() unsafe.Pointer {
	return C.nox_thing_glow_orb_draw
}

func Nox_xxx_drawAllMB_475810_draw_B(vp *noxrender.Viewport) int {
	if noxflags.HasEngine(noxflags.EngineNoFloorRendering) {
		return 0
	}
	tile := tileAtPoint(types.Pointf{X: float32(vp.World.Max.X), Y: float32(vp.World.Max.Y)})
	if tile == -1 || tile == 255 {
		return 0
	}
	return 1
}
func Sub_4C5060(vp *noxrender.Viewport) {
	objectRenderBeamDraw(vp)
}
func Nox_xxx_drawWalls_473C10(vp *noxrender.Viewport, a2 *server.Wall) {
	worldWallDraw(vp, a2)
}
func Sub_4761B0(dr *client.Drawable) int {
	return int(presentationDrawableY(dr))
}
func Sub_476080(a1 unsafe.Pointer) int {
	return int(presentationWallY((*[8]byte)(a1)))
}
func Sub_459DB0(dr *client.Drawable) int {
	return int(C.sub_459DB0((*nox_drawable)(dr.C())))
}
func Sub_49A6A0(vp *noxrender.Viewport, dr *client.Drawable) {
	combatHealthDraw(vp, dr)
}
func Nox_xxx_sprite_4756E0_drawable(dr *client.Drawable) int {
	return bool2int(worldWallStaticPass(dr))
}
func Nox_xxx_sprite_475740_drawable(dr *client.Drawable) int {
	return bool2int(worldWallDynamicPass(dr))
}
func Nox_xxx_sprite_4757A0_drawable(dr *client.Drawable) int {
	return bool2int(worldWallSpecialPass(dr))
}
func Sub_4757D0_drawable(dr *client.Drawable) int {
	return bool2int(worldWallInactivePass(dr))
}
func Nox_xxx_tileDrawImpl_4826A0(vp *noxrender.Viewport) {
	tileCompositionFull(vp)
}
func Nox_xxx_tileDrawMB_481C20_A(vp *noxrender.Viewport, a2 int) {
	tileCompositionHorizontal(vp, a2)
}
func Nox_xxx_tileDrawMB_481C20_B(vp *noxrender.Viewport, a2 int) {
	tileCompositionVertical(vp, a2)
}
func Nox_xxx_tileCheckRedrawMB_482570(vp *noxrender.Viewport) int {
	return tileCompositionRedraw(vp)
}
