//go:build porttest

package legacy

/*
#include "GAME3.h"
#include "client__shell__selcolor.h"
extern uint32_t dword_587000_171388;
*/
import "C"

func PortTestCharacterConstruct(color bool) int {
	if color {
		return int(C.nox_game_showSelColor_4A5D00())
	}
	return int(C.nox_game_showSelClass_4A4840())
}
func PortTestCharacterDefaultSet(value int32) int32   { return int32(C.sub_4A7A60(C.int(value))) }
func PortTestCharacterAdmissionSet(value int32) int32 { return int32(C.sub_4A7A70(C.int(value))) }
func PortTestCharacterDefaultGet() uint32             { return uint32(C.dword_587000_171388) }
