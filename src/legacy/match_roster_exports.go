package legacy

/*
#include "defs.h"
typedef const void match_roster_const_data;
*/
import "C"
import "unsafe"
import "github.com/opennox/opennox/v1/server"

//export nox_xxx_netGuiGameSettings_4DD9B0
func nox_xxx_netGuiGameSettings_4DD9B0(mode C.char, data *C.match_roster_const_data, to C.int) C.int {
	return C.int(matchRosterGUISettings(byte(mode), unsafe.Pointer(data), int(to)))
}

//export nox_xxx_netGameSettings_4DEF00
func nox_xxx_netGameSettings_4DEF00() C.int { return C.int(matchRosterSettings()) }

//export nox_xxx_newPlayerSendAllPlayers_4DE300
func nox_xxx_newPlayerSendAllPlayers_4DE300(to C.int) *C.char {
	matchRosterSendPlayers(int(to))
	return nil
}

//export nox_xxx_sendAllPlayerIDs_4DE270
func nox_xxx_sendAllPlayerIDs_4DE270(pl C.int) C.int {
	return C.int(matchRosterPlayerIDs((*server.Player)(unsafe.Pointer(uintptr(uint32(pl))))))
}

//export nox_xxx_servMinimapRevealFlag_4DE380
func nox_xxx_servMinimapRevealFlag_4DE380(to C.int) C.int { return C.int(matchRosterMinimap(int(to))) }

//export nox_xxx_netSendSimpleObject2_4DF360
func nox_xxx_netSendSimpleObject2_4DF360(to C.int, u *nox_object_t) C.int {
	return C.int(matchRosterSimpleObject(int(to), asObjectS(u)))
}

//export sub_4DF2E0
func sub_4DF2E0(to C.int) { matchRosterTeamRoster(int(to)) }

//export sub_4E8310
func sub_4E8310() *C.char { return (*C.char)(matchRosterFlagBase()) }

//export sub_4E8320
func sub_4E8320(index C.uchar) *C.uchar { return (*C.uchar)(matchRosterFlagRecord(byte(index))) }

//export sub_5095E0
func sub_5095E0() C.int { return C.int(matchRosterWinner(true)) }
