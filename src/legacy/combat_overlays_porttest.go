//go:build porttest

package legacy

/*
#include <stdint.h>
extern uint32_t nox_color_yellow_2589772;
extern unsigned int nox_player_netCode_85319C;
*/
import "C"
import (
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/client/noxrender"
	"unsafe"
)

func PortTestCombatOverlayGlobals() (map[string]*uint32, map[string]*unsafe.Pointer, func()) {
	words := map[string]*uint32{
		"yellow":    (*uint32)(unsafe.Pointer(&C.nox_color_yellow_2589772)),
		"localCode": (*uint32)(unsafe.Pointer(&C.nox_player_netCode_85319C)),
		"friends":   (*uint32)(unsafe.Pointer(&combatFriendHead)),
		"health":    (*uint32)(unsafe.Pointer(&combatHealthHead)),
		"font":      (*uint32)(unsafe.Pointer(&combatHealthFont)),
		"feedRows":  (*uint32)(unsafe.Pointer(&combatFeedRows)),
		"feedWrite": (*uint32)(unsafe.Pointer(&combatFeedWrite)),
		"feedRead":  (*uint32)(unsafe.Pointer(&combatFeedRead)),
	}
	pools := map[string]*unsafe.Pointer{"friend": (*unsafe.Pointer)(unsafe.Pointer(&combatFriendPool)), "health": (*unsafe.Pointer)(unsafe.Pointer(&combatHealthPool))}
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
func PortTestCombatAllyClear() { combatAllyClear() }
func PortTestCombatAllyLookup(code uint32) unsafe.Pointer {
	return unsafe.Pointer(combatAllyLookup(code))
}
func PortTestCombatAllyFree() unsafe.Pointer               { return unsafe.Pointer(combatAllyFree()) }
func PortTestCombatAllyAdd(code uint32, x, y uint16) bool  { return combatAllyAdd(code, x, y) }
func PortTestCombatAllyRemove(code uint32) bool            { return combatAllyRemove(code) }
func PortTestCombatAllyFlag(code uint32, flag byte) bool   { return combatAllyFlag(code, flag) }
func PortTestCombatAllyPair(code uint32, x, y uint16) bool { return combatAllyPair(code, x, y) }
func PortTestCombatAllyFirst(code uint32, x uint16) bool   { return combatAllyFirst(code, x) }
func PortTestCombatAllyRead(code uint32, pair *[2]uint16, flag *byte) bool {
	return combatAllyRead(code, &pair[0], &pair[1], flag)
}
func PortTestCombatIsAlly(code uint32) bool { return combatAllyLookup(code) != nil }
func PortTestCombatFriendInit() bool        { return combatFriendInit() }
func PortTestCombatFriendClear()            { combatFriendClear() }
func PortTestCombatFriendDestroy()          { combatFriendDestroy() }
func PortTestCombatFriendAdd(code uint32) unsafe.Pointer {
	return unsafe.Pointer(combatFriendAdd(code))
}
func PortTestCombatFriendRemove(code uint32)                              { combatFriendRemove(code) }
func PortTestCombatFriendHas(code uint32) bool                            { return combatFriendHas(code) }
func PortTestCombatHealthInit() bool                                      { return combatHealthInit() }
func PortTestCombatHealthClear()                                          { combatHealthClear() }
func PortTestCombatHealthDestroy()                                        { combatHealthDestroy() }
func PortTestCombatHealthAdd(code uint32, amount int16)                   { combatHealthAdd(code, amount) }
func PortTestCombatHealthRemove(p unsafe.Pointer)                         { combatHealthRemove((*combatHealth)(p)) }
func PortTestCombatHealthDraw(v *noxrender.Viewport, dr *client.Drawable) { combatHealthDraw(v, dr) }
func PortTestCombatFeedInit()                                             { combatFeedInit() }
func PortTestCombatFeedAdd(p unsafe.Pointer)                              { combatFeedAdd(p) }
func PortTestCombatFeedDraw()                                             { combatFeedDraw() }
func PortTestCombatFeedRow(p *[6]uint32)                                  { combatFeedRow(p) }
func PortTestCombatEffectAttach(p unsafe.Pointer, dr *client.Drawable) {
	combatFXAttach((*combatFX)(p), dr)
}
func PortTestCombatEffectDetach(p unsafe.Pointer)                         { combatFXDetach((*combatFX)(p)) }
func PortTestCombatEffectDraw(v *noxrender.Viewport, dr *client.Drawable) { combatFXDraw(v, dr) }
func PortTestCombatEffectKind(kind int, v *noxrender.Viewport, dr *client.Drawable, p unsafe.Pointer) {
	combatFXTrail(v, dr, (*combatFX)(p), kind != 1)
}
