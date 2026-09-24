//go:build porttest

package legacy

/*
#include "defs.h"
#include "GAME1_2.h"
#include "GAME1_3.h"
#include "GAME2_3.h"
#include "client__shell__noxworld.h"
*/
import "C"
import (
	"github.com/opennox/opennox/v1/client/gui"
	"unsafe"
)

func PortTestServerBrowserWords() (map[string]*uint32, func()) {
	words := map[string]*uint32{
		"dword_587000_87408":               (*uint32)(unsafe.Pointer(&browserUI.listMode)),
		"dword_587000_87412":               (*uint32)(unsafe.Pointer(&browserUI.region)),
		"dword_5d4594_1305788":             (*uint32)(unsafe.Pointer(&browserUI.polygonsReady)),
		"dword_5d4594_1307716":             (*uint32)(unsafe.Pointer(&browserUI.popup)),
		"dword_5d4594_1307720":             (*uint32)(unsafe.Pointer(&browserUI.popupCount)),
		"dword_5d4594_528252":              (*uint32)(unsafe.Pointer(&onlineRetryPending)),
		"dword_5d4594_528256":              (*uint32)(unsafe.Pointer(&onlineRetryActive)),
		"dword_5d4594_814548":              (*uint32)(unsafe.Pointer(&browserUI.connectionState)),
		"dword_5d4594_814624":              (*uint32)(unsafe.Pointer(&browserUI.selected)),
		"dword_5d4594_814984":              (*uint32)(unsafe.Pointer(&browserUI.mapWindow)),
		"dword_5d4594_814988":              (*uint32)(unsafe.Pointer(&browserUI.overview)),
		"dword_5d4594_814992":              (*uint32)(unsafe.Pointer(&browserUI.filter)),
		"dword_5d4594_814996":              (*uint32)(unsafe.Pointer(&browserUI.label)),
		"dword_5d4594_815000":              (*uint32)(unsafe.Pointer(&browserUI.detailPanel)),
		"dword_5d4594_815004":              (*uint32)(unsafe.Pointer(&browserUI.detailList)),
		"dword_5d4594_815016":              (*uint32)(unsafe.Pointer(&browserUI.playersColumn)),
		"dword_5d4594_815020":              (*uint32)(unsafe.Pointer(&browserUI.modeColumn)),
		"dword_5d4594_815024":              (*uint32)(unsafe.Pointer(&browserUI.mapColumn)),
		"dword_5d4594_815028":              (*uint32)(unsafe.Pointer(&browserUI.pingColumn)),
		"dword_5d4594_815032":              (*uint32)(unsafe.Pointer(&browserUI.statusColumn)),
		"dword_5d4594_815044":              (*uint32)(unsafe.Pointer(&browserUI.transition)),
		"dword_5d4594_815052":              (*uint32)(unsafe.Pointer(&browserUI.hosting)),
		"dword_5d4594_815056":              (*uint32)(unsafe.Pointer(&browserUI.hasSelection)),
		"dword_5d4594_815096":              (*uint32)(unsafe.Pointer(&browserUI.pendingKicked)),
		"dword_5d4594_815100":              (*uint32)(unsafe.Pointer(&browserUI.pendingTimeout)),
		"dword_5d4594_815104":              (*uint32)(unsafe.Pointer(&browserUI.retry)),
		"nox_client_connError_814552":      (*uint32)(unsafe.Pointer(&browserUI.connectionError)),
		"nox_game_createOrJoin_815048":     (*uint32)(unsafe.Pointer(&browserUI.creating)),
		"nox_wnd_xxx_815040":               (*uint32)(unsafe.Pointer(&browserUI.animation)),
		"nox_wol_server_result_cnt_815088": (*uint32)(unsafe.Pointer(&browserUI.resultCount)),
		"nox_wol_servers_sorting_166704":   (*uint32)(unsafe.Pointer(&browserUI.sort)),
		"nox_wol_wnd_gameList_815012":      (*uint32)(unsafe.Pointer(&browserUI.gameList)),
		"nox_wol_wnd_world_814980":         (*uint32)(unsafe.Pointer(&browserUI.world)),
	}
	old := make(map[string]uint32)
	for name, p := range words {
		old[name] = *p
		*p = 0
	}
	return words, func() {
		for name, p := range words {
			*p = old[name]
		}
	}
}
func PortTestServerBrowserClocks() (map[string]*uint64, func()) {
	words := map[string]*uint64{
		"qword_5d4594_814956": (*uint64)(unsafe.Pointer(&browserUI.connectionDeadline)),
		"qword_5d4594_815068": (*uint64)(unsafe.Pointer(&browserUI.refreshDeadline)),
	}
	old := make(map[string]uint64)
	for name, p := range words {
		old[name] = *p
		*p = 0
	}
	return words, func() {
		for name, p := range words {
			*p = old[name]
		}
	}
}
func PortTestServerBrowserStateSet(v int32) int32 { return int32(sub_43AF90(C.int(v))) }
func PortTestServerBrowserErrorSet(v int32)       { nox_client_setConnError_43AFA0(C.int(v)) }
func PortTestServerBrowserStateGet(which int) uint32 {
	switch which {
	case 0:
		return uint32(C.sub_43AF30())
	case 1:
		return uint32(sub_43AF40())
	case 2:
		return uint32(sub_43AF80())
	case 3:
		return uint32(sub_43B6D0())
	case 4:
		return uint32(nox_client_getServerAddr_43B300())
	case 5:
		return uint32(nox_client_getServerPort_43B320())
	case 6:
		return uint32(sub_43B340())
	}
	panic(which)
}
func PortTestServerBrowserSortClick(id int) { nox_wol_servers_sortBtnHandler_4A0290(int32(id)) }
func PortTestServerBrowserFormat(addr string, port uint16, out unsafe.Pointer) int {
	p := CString(addr)
	defer StrFree(p)
	return int(nox_sprintAddrPort_43BC80(p, C.ushort(port), (*C.char)(out)))
}

func PortTestServerBrowserRow(record unsafe.Pointer) {
	nox_gui_wol_newServerLine_43B7C0((*C.nox_gui_server_ent_t)(record))
}
func PortTestServerBrowserTrimName(text unsafe.Pointer, width byte) uintptr {
	return uintptr(unsafe.Pointer(sub_43BC10((*C.wchar2_t)(text), C.uchar(width))))
}

func PortTestServerBrowserDetails(record unsafe.Pointer) {
	nox_client_gui_serverInfoBlock_4394D0(C.int(uintptr(record)))
}

func PortTestServerBrowserTick() int { return int(sub_438770()) }
func PortTestServerBrowserNotification(which int, window unsafe.Pointer, id uint32) int {
	if which == 0 {
		return int(nox_xxx_windowMultiplayerSub_439E70(int32(uintptr(window)), uint32(22), (*int32)(unsafe.Pointer(uintptr(id))), int32(0)))
	}
	return int(sub_439D00((*int32)(window), int32(22), uint32(id), int32(0)))
}

func PortTestServerBrowserHostDescription() uintptr { return uintptr(unsafe.Pointer(sub_43AA70())) }

func PortTestServerBrowserColumnEvent(w unsafe.Pointer, code int, a, b uint32, input bool) int {
	if input {
		return int(sub_438EF0((*gui.Window)(w), int32(code), uint32(a), int32(b)))
	}
	return int(sub_439050(uint32(uintptr(w)), uint32(code), unsafe.Pointer(uintptr(a)), uint32(b)))
}
func PortTestServerBrowserInfoPosition(x, y int32, out *[2]uint32) {
	browserInfoPosition(x, y, out)
}

func PortTestServerBrowserEvent(w unsafe.Pointer, code int, a, b uint32) int {
	return int(nox_xxx_windowMultiplayerSub_439E70(int32(uintptr(w)), uint32(code), (*int32)(unsafe.Pointer(uintptr(a))), int32(b)))
}
