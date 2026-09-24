package legacy

/*
#include "defs.h"
#include "client__draw__doordraw.h"
typedef const char* objectDrawCString;
*/
import "C"
import (
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/client/noxrender"
	"unsafe"
)

//export nox_thing_door_draw
func nox_thing_door_draw(a1 *C.uint32_t, a2 *C.nox_drawable) C.int {
	vp := (*noxrender.Viewport)(unsafe.Pointer(a1))
	dr := (*client.Drawable)(unsafe.Pointer(a2))
	return C.int(objectDoorDraw(vp, dr))
}

//export nox_thing_arrow_draw
func nox_thing_arrow_draw(a1 *C.int, a2 *C.nox_drawable) C.int {
	vp := (*noxrender.Viewport)(unsafe.Pointer(a1))
	dr := (*client.Drawable)(unsafe.Pointer(a2))
	return C.int(objectArrowDraw(vp, dr, false))
}

//export nox_thing_weak_arrow_draw
func nox_thing_weak_arrow_draw(a1 *C.int, a2 *C.nox_drawable) C.int {
	vp := (*noxrender.Viewport)(unsafe.Pointer(a1))
	dr := (*client.Drawable)(unsafe.Pointer(a2))
	return C.int(objectArrowDraw(vp, dr, true))
}

//export nox_thing_arrow_tail_link_draw
func nox_thing_arrow_tail_link_draw(a1 *C.uint32_t, a2 *C.nox_drawable) C.int {
	vp := (*noxrender.Viewport)(unsafe.Pointer(a1))
	dr := (*client.Drawable)(unsafe.Pointer(a2))
	return C.int(objectArrowTailDraw(vp, dr, false))
}

//export nox_thing_weak_arrow_tail_link_draw
func nox_thing_weak_arrow_tail_link_draw(a1 *C.uint32_t, a2 *C.nox_drawable) C.int {
	vp := (*noxrender.Viewport)(unsafe.Pointer(a1))
	dr := (*client.Drawable)(unsafe.Pointer(a2))
	return C.int(objectArrowTailDraw(vp, dr, true))
}

//export nox_thing_glyph_draw
func nox_thing_glyph_draw(a1 *C.int, a2 *C.nox_drawable) C.int {
	vp := (*noxrender.Viewport)(unsafe.Pointer(a1))
	dr := (*client.Drawable)(unsafe.Pointer(a2))
	return C.int(objectGlyphDraw(vp, dr))
}

//export nox_thing_summon_effect_draw
func nox_thing_summon_effect_draw(a1 *C.int, a2 *C.nox_drawable) C.int {
	vp := (*noxrender.Viewport)(unsafe.Pointer(a1))
	dr := (*client.Drawable)(unsafe.Pointer(a2))
	return C.int(objectSummonDraw(vp, dr))
}

//export nox_thing_weapon_draw
func nox_thing_weapon_draw(a1 *C.int, a2 *C.nox_drawable) C.int {
	vp := (*noxrender.Viewport)(unsafe.Pointer(a1))
	dr := (*client.Drawable)(unsafe.Pointer(a2))
	return C.int(objectEquipmentDraw(vp, dr, false, false))
}

//export nox_thing_weapon_animate_draw
func nox_thing_weapon_animate_draw(a1 *C.int, a2 *C.nox_drawable) C.int {
	vp := (*noxrender.Viewport)(unsafe.Pointer(a1))
	dr := (*client.Drawable)(unsafe.Pointer(a2))
	return C.int(objectEquipmentDraw(vp, dr, false, true))
}

//export nox_thing_armor_draw
func nox_thing_armor_draw(a1 *C.int, a2 *C.nox_drawable) C.int {
	vp := (*noxrender.Viewport)(unsafe.Pointer(a1))
	dr := (*client.Drawable)(unsafe.Pointer(a2))
	return C.int(objectEquipmentDraw(vp, dr, true, false))
}

//export nox_thing_armor_animate_draw
func nox_thing_armor_animate_draw(a1 *C.int, a2 *C.nox_drawable) C.int {
	vp := (*noxrender.Viewport)(unsafe.Pointer(a1))
	dr := (*client.Drawable)(unsafe.Pointer(a2))
	return C.int(objectEquipmentDraw(vp, dr, true, true))
}

//export nox_thing_spherical_shield_draw
func nox_thing_spherical_shield_draw(a1 *C.int, a2 *C.nox_drawable) C.int {
	vp := (*noxrender.Viewport)(unsafe.Pointer(a1))
	dr := (*client.Drawable)(unsafe.Pointer(a2))
	return C.int(objectShieldDraw(vp, dr))
}

//export nox_thing_monster_gen_draw
func nox_thing_monster_gen_draw(a1 *C.int, a2 *C.nox_drawable) C.int {
	vp := (*noxrender.Viewport)(unsafe.Pointer(a1))
	dr := (*client.Drawable)(unsafe.Pointer(a2))
	return C.int(objectGeneratorDraw(vp, dr))
}

//export nox_thing_pressure_plate_draw
func nox_thing_pressure_plate_draw(a1 *C.uint32_t, a2 *C.nox_drawable) C.int {
	vp := (*noxrender.Viewport)(unsafe.Pointer(a1))
	dr := (*client.Drawable)(unsafe.Pointer(a2))
	return C.int(objectPressureDraw(vp, dr))
}

//export nox_thing_trigger_draw
func nox_thing_trigger_draw(a1 *C.uint32_t, a2 *C.nox_drawable) C.int {
	vp := (*noxrender.Viewport)(unsafe.Pointer(a1))
	dr := (*client.Drawable)(unsafe.Pointer(a2))
	return C.int(objectTriggerDraw(vp, dr))
}

//export nox_thing_black_powder_draw
func nox_thing_black_powder_draw(a1 *C.uint32_t, a2 *C.nox_drawable) C.int {
	vp := (*noxrender.Viewport)(unsafe.Pointer(a1))
	dr := (*client.Drawable)(unsafe.Pointer(a2))
	return C.int(objectPowderDraw(vp, dr))
}

//export nox_thing_base_draw
func nox_thing_base_draw(a1 *C.int, a2 *C.nox_drawable) C.int {
	vp := (*noxrender.Viewport)(unsafe.Pointer(a1))
	dr := (*client.Drawable)(unsafe.Pointer(a2))
	objectEquipmentDraw(vp, dr, false, false)
	return 1
}

//export nox_thing_flag_draw
func nox_thing_flag_draw(a1 *C.int, a2 *C.nox_drawable) C.int {
	vp := (*noxrender.Viewport)(unsafe.Pointer(a1))
	dr := (*client.Drawable)(unsafe.Pointer(a2))
	return C.int(objectFlagDraw(vp, dr))
}

//export nox_things_door_draw_parse
func nox_things_door_draw_parse(obj *C.nox_thing, f *C.nox_memfile, attr *C.char) C.bool {
	o := (*client.ObjectType)(unsafe.Pointer(obj))
	o.DrawFunc = C.nox_thing_door_draw
	o.DrawData = spriteStaticRandomData(asMemfileP(unsafe.Pointer(f)), unsafe.Slice((*byte)(unsafe.Pointer(attr)), 256))
	return C.bool(o.DrawData != nil)
}

//export sub_4B9470
func sub_4B9470(arg *C.objectDrawCString) C.int { return C.int(objectTeamColor(unsafe.Pointer(arg))) }

//export sub_4B9650
func sub_4B9650(typ C.int) *C.uint32_t {
	return (*C.uint32_t)(unsafe.Pointer(uintptr(objectBaseMaterials(int(typ)))))
}

func sub_4BC720(p C.int) C.int {
	return C.int(uintptr(objectGeneratorCountdown((*client.Drawable)(unsafe.Pointer(uintptr(uint32(p)))))))
}

//export nox_xxx_updDrawMonsterGen_4BC920
func nox_xxx_updDrawMonsterGen_4BC920() C.int { return 1 }
