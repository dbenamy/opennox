package legacy

/*
#include "GAME1.h"
#include "GAME1_2.h"
#include "GAME1_1.h"
#include "GAME1_3.h"
#include "GAME2.h"
#include "GAME2_1.h"
#include "GAME2_2.h"
#include "GAME3.h"
#include "GAME3_1.h"
#include "GAME3_2.h"
#include "GAME3_3.h"
#include "GAME5.h"
#include "GAME5_2.h"
#include "GAME2_3.h"
#include "common__system__team.h"
#include "common__net_list.h"
#include "client__gui__guimsg.h"
#include "client__drawable__drawable.h"

*/
import "C"
import (
	"github.com/opennox/opennox/v1/common/memmap"
	"unsafe"

	"github.com/opennox/opennox/v1/legacy/common/ccall"
	"github.com/opennox/opennox/v1/server"
)

var (
	Nox_client_getIntroScreenDuration_44E3B0 func() int
	Nox_client_getBriefDuration              func() int
	Nox_game_exit_xxx2                       func()
	Sub_470510                               func()
	Sub_4703F0                               func()
	Nox_xxx_cliDrawConnectedLoop_43B360      func() int
	Nox_client_guiXxxDestroy_4A24A0          func() int
	Nox_client_quit_4460C0                   func()
)

//export nox_client_getIntroScreenDuration_44E3B0
func nox_client_getIntroScreenDuration_44E3B0() int {
	return Nox_client_getIntroScreenDuration_44E3B0()
}

//export nox_client_getBriefDuration
func nox_client_getBriefDuration() int {
	return Nox_client_getBriefDuration()
}

//export nox_game_SetCliDrawFunc
func nox_game_SetCliDrawFunc(fnc unsafe.Pointer) {
	if fnc == nil {
		GetClient().SetDrawFunc(nil)
	} else {
		GetClient().SetDrawFunc(func() bool {
			return ccall.CallIntVoid(fnc) != 0
		})
	}
}

//export sub_43DE40
func sub_43DE40(fnc unsafe.Pointer) int {
	if fnc == nil {
		GetServer().SetUpdateFunc2(nil)
	} else {
		GetServer().SetUpdateFunc2(func() bool {
			return ccall.CallIntVoid(fnc) != 0
		})
	}
	return 1
}

func nox_game_exit_xxx2() {
	Nox_game_exit_xxx2()
}

func sub_470510() {
	Sub_470510()
}

//export sub_4703F0
func sub_4703F0() {
	Sub_4703F0()
}

//export nox_xxx_cliDrawConnectedLoop_43B360
func nox_xxx_cliDrawConnectedLoop_43B360() int {
	return Nox_xxx_cliDrawConnectedLoop_43B360()
}

//export nox_client_guiXxxDestroy_4A24A0
func nox_client_guiXxxDestroy_4A24A0() int {
	return Nox_client_guiXxxDestroy_4A24A0()
}

//export nox_client_quit_4460C0
func nox_client_quit_4460C0() {
	Nox_client_quit_4460C0()
}
func Sub_43DB60() {
	audioEventMusicEnter()
}
func Nox_xxx_mapGenStart_4D4320() int {
	return int(mapOrchestrationStart())
}
func Nox_xxx_servResetPlayers_4D23C0() {
	sessionResetPlayers()
}
func Nox_xxx_gameLoopMemDump_413E30() {

}

func Nox_xxx_getRandomName_4358A0() string {
	return clientRandomName()
}
func Sub_4E4EF0() {
	reliableResetRates()
}
func Sub_48D740() {
	clientSequenceInit()
}
func Sub_473930() {
	clientAnimationCachesLoad()
}
func Sub_43DBA0() {
	audioEventMusicLeave()
}
func Nox_xxx_getHostInfoPtr_431770() *server.PlayerInfo {
	return (*server.PlayerInfo)(memmap.PtrOff(0x5D4594, 807172))
}
func Sub_41FA40() string {
	return "" // The former account selector remains unset.
}
func Sub_43AF40() int {
	return int(sub_43AF40())
}
func Sub_43AA70() {
	sub_43AA70()
}
func Nox_server_gameDoSwitchMap_40A680() int {
	return int(C.int(serverConfigUpdatedGet()))
}
func Sub_459D60() int {
	return int(*serverOptionsWord(1046544))
}
func Sub_459DA0() int {
	return bool2int(serverOptionsRoot != 0)
}
func Sub_4DF020() {
	sessionBroadcastSettings()
}
func Sub_4161E0() {
	serverConfigRefresh()
}
func Sub_473960() {
	clientAnimationCachesClear()
}
func Sub_48D800() {
	chatBubbleDestroy()
}
func Sub_49A8C0() {
	combatHealthDestroy()
}
func Sub_4E4DE0() {
	reliableInit()
}
func Sub_48D760() {
	clientSequenceFree()
}
func Sub_417CF0() {
	GetServer().TeamsRemoveActive(false)
}
func Sub_499450() {
	presentationShieldDestroy()
}
func Sub_4959D0() {
	combatFriendDestroy()
}

func Nox_xxx_netSavePlayer_41CE00() {
	playerFileSaveRequest()
}
func Sub_4D39F0(a1 string) {
	prefabScriptGeneration(a1)
}
func Sub_48D4B0(a1 int) {
	sub_48D4B0(C.int(a1))
}
func Nox_xxx_set3512_40A340(a1 int) {
	serverConfig3512Set(int32(a1))
}
func Sub_459D50(a1 int) {
	serverOptionsDirty(a1)
}
func Nox_xxx_gameSetMapPath_409D70(a1 string) {
	sessionSetMapPath((*byte)(unsafe.Pointer(internCStr(a1))))
}
func Nox_xxx_gui_43E1A0(a1 int) {
	clientModalWindow(a1)
}
func Nox_xxx_printCentered_445490(str string) {
	wstr, free := CWString(str)
	defer free()
	nox_xxx_printCentered_445490(wstr)
}
func Nox_xxx_mapValidateMB_4CF470(a1 string, a2 uint32) int {
	return int(sessionMapValidate((*byte)(unsafe.Pointer(internCStr(a1))), a2))
}
func Nox_xxx_copyServerIPAndPort_431790(a1 string) {
	clientServerAddressCopy(a1)
}
