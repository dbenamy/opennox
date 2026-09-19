//go:build porttest

package legacy

/*
#include "GAME2_3.h"
char* sub_494FF0(void);
int sub_495430(void);
void sub_495B50(void* p);
void sub_495FC0(void* p, nox_drawable* dr);
extern void* nox_alloc_friendList_1203860;
extern void* nox_alloc_healthChange_1301772;
extern uint32_t dword_5d4594_1203864;
extern uint32_t dword_5d4594_1301776;
extern uint32_t dword_5d4594_1301780;
extern uint32_t dword_5d4594_1203832;
extern uint32_t dword_5d4594_1203836;
extern uint32_t dword_5d4594_1203840;
extern uint32_t nox_color_yellow_2589772;
extern unsigned int nox_player_netCode_85319C;
*/
import "C"
import (
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/client/noxrender"
	"unsafe"
)

func PortTestCombatAllyClear() { C.sub_4951C0() }
func PortTestCombatAllyLookup(code uint32) unsafe.Pointer {
	return unsafe.Pointer(C.sub_495020(C.int(code)))
}
func PortTestCombatAllyFree() unsafe.Pointer { return unsafe.Pointer(C.sub_494FF0()) }
func PortTestCombatAllyAdd(code uint32, x, y uint16) bool {
	return C.sub_495060(C.int(code), C.short(x), C.short(y)) != 0
}
func PortTestCombatAllyRemove(code uint32) bool { return C.sub_4950C0(C.int(code)) != 0 }
func PortTestCombatAllyFlag(code uint32, flag byte) bool {
	return C.sub_4950F0(C.int(code), C.char(flag)) != 0
}
func PortTestCombatAllyPair(code uint32, x, y uint16) bool {
	return C.sub_495120(C.int(code), C.short(x), C.short(y)) != 0
}
func PortTestCombatAllyFirst(code uint32, x uint16) bool {
	return C.sub_495150(C.int(code), C.short(x)) != 0
}
func PortTestCombatAllyRead(code uint32, values *[2]uint16, flag *byte) bool {
	return C.sub_495180(C.int(code), (*C.uint16_t)(&values[0]), (*C.uint16_t)(&values[1]), (*C.uint8_t)(flag)) != 0
}
func PortTestCombatIsAlly(code uint32) bool {
	return C.nox_xxx_unitSpriteCheckAlly_4951F0(C.int(code)) != 0
}

func PortTestCombatOverlayGlobals() (map[string]*uint32, map[string]*unsafe.Pointer, func()) {
	words := map[string]*uint32{
		"yellow":    (*uint32)(unsafe.Pointer(&C.nox_color_yellow_2589772)),
		"localCode": (*uint32)(unsafe.Pointer(&C.nox_player_netCode_85319C)),
		"friends":   (*uint32)(unsafe.Pointer(&C.dword_5d4594_1203864)),
		"health":    (*uint32)(unsafe.Pointer(&C.dword_5d4594_1301776)),
		"font":      (*uint32)(unsafe.Pointer(&C.dword_5d4594_1301780)),
		"feedRows":  (*uint32)(unsafe.Pointer(&C.dword_5d4594_1203832)),
		"feedWrite": (*uint32)(unsafe.Pointer(&C.dword_5d4594_1203836)),
		"feedRead":  (*uint32)(unsafe.Pointer(&C.dword_5d4594_1203840)),
	}
	pools := map[string]*unsafe.Pointer{"friend": (*unsafe.Pointer)(unsafe.Pointer(&C.nox_alloc_friendList_1203860)), "health": (*unsafe.Pointer)(unsafe.Pointer(&C.nox_alloc_healthChange_1301772))}
	oldWords := make(map[string]uint32)
	oldPools := make(map[string]unsafe.Pointer)
	for k, p := range words {
		oldWords[k] = *p
		*p = 0
	}
	for k, p := range pools {
		oldPools[k] = *p
		*p = nil
	}
	return words, pools, func() {
		for k, p := range words {
			*p = oldWords[k]
		}
		for k, p := range pools {
			*p = oldPools[k]
		}
	}
}
func PortTestCombatFriendInit() bool { return C.nox_xxx_allocClassListFriends_495980() != 0 }
func PortTestCombatFriendClear()     { C.sub_4959B0() }
func PortTestCombatFriendDestroy()   { C.sub_4959D0() }
func PortTestCombatFriendAdd(code uint32) unsafe.Pointer {
	return unsafe.Pointer(C.nox_xxx_cliAddObjFriend_4959F0(C.int(code)))
}
func PortTestCombatFriendRemove(code uint32)   { C.sub_495A20(C.int(code)) }
func PortTestCombatFriendHas(code uint32) bool { return C.sub_495A80(C.int(code)) != 0 }
func PortTestCombatHealthInit() bool           { return C.nox_xxx_allocArrayHealthChanges_49A5F0() != 0 }
func PortTestCombatHealthClear()               { C.sub_49A630() }
func PortTestCombatHealthDestroy()             { C.sub_49A8C0() }
func PortTestCombatHealthAdd(code uint32, amount int16) {
	C.nox_xxx_cliAddHealthChange_49A650(C.int(code), C.short(amount))
}
func PortTestCombatHealthRemove(p unsafe.Pointer) { C.sub_49A880(C.int(uintptr(p))) }
func PortTestCombatHealthDraw(v *noxrender.Viewport, dr *client.Drawable) {
	C.sub_49A6A0((*C.nox_draw_viewport_t)(v.C()), (*C.nox_drawable)(dr.C()))
}
func PortTestCombatFeedInit()                { C.sub_4958F0() }
func PortTestCombatFeedAdd(p unsafe.Pointer) { C.sub_495210(C.int(uintptr(p))) }
func PortTestCombatFeedDraw()                { C.sub_495430() }
func PortTestCombatFeedRow(p *[6]uint32)     { C.sub_495500((*C.int)(unsafe.Pointer(p))) }
func PortTestCombatEffectAttach(p unsafe.Pointer, dr *client.Drawable) {
	C.sub_495FC0(p, (*C.nox_drawable)(dr.C()))
}
func PortTestCombatEffectDetach(p unsafe.Pointer) { C.sub_495B50(p) }
func PortTestCombatEffectDraw(v *noxrender.Viewport, dr *client.Drawable) {
	C.sub_495BB0((*C.nox_drawable)(dr.C()), (*C.nox_draw_viewport_t)(v.C()))
}
func PortTestCombatEffectKind(kind int, v *noxrender.Viewport, dr *client.Drawable, p unsafe.Pointer) {
	if kind == 1 {
		C.sub_495BF0(C.int(uintptr(dr.C())), C.int(uintptr(p)), C.int(uintptr(v.C())))
	} else {
		C.sub_495D00((*C.uint32_t)(dr.C()), C.int(uintptr(p)), (*C.uint32_t)(v.C()))
	}
}
