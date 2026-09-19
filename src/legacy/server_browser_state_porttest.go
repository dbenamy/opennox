//go:build porttest

package legacy

/*
#include "defs.h"
#include "GAME1_2.h"
#include "GAME1_3.h"
#include "GAME2_3.h"
#include "client__shell__noxworld.h"
void nox_client_gui_serverInfoBlock_4394D0(int a1);
extern uint32_t dword_587000_87408;
extern int dword_587000_87412;
extern uint32_t dword_5d4594_1305788;
extern uint32_t dword_5d4594_1307716;
extern uint32_t dword_5d4594_1307720;
extern uint32_t dword_5d4594_528252;
extern uint32_t dword_5d4594_528256;
extern unsigned int dword_5d4594_814548;
extern void* dword_5d4594_814624;
extern void* dword_5d4594_814984;
extern uint32_t dword_5d4594_814988;
extern uint32_t dword_5d4594_814992;
extern void* dword_5d4594_814996;
extern void* dword_5d4594_815000;
extern nox_window* dword_5d4594_815004;
extern uint32_t dword_5d4594_815016;
extern uint32_t dword_5d4594_815020;
extern uint32_t dword_5d4594_815024;
extern uint32_t dword_5d4594_815028;
extern uint32_t dword_5d4594_815032;
extern uint32_t dword_5d4594_815044;
extern uint32_t dword_5d4594_815052;
extern uint32_t dword_5d4594_815056;
extern uint32_t dword_5d4594_815096;
extern uint32_t dword_5d4594_815100;
extern int dword_5d4594_815104;
extern uint32_t nox_client_connError_814552;
extern uint32_t nox_game_createOrJoin_815048;
extern nox_gui_animation* nox_wnd_xxx_815040;
extern uint32_t nox_wol_server_result_cnt_815088;
extern uint32_t nox_wol_servers_sorting_166704;
extern uint32_t nox_wol_wnd_gameList_815012;
extern nox_window* nox_wol_wnd_world_814980;
extern uint64_t qword_5d4594_814956;
extern uint64_t qword_5d4594_815068;
*/
import "C"
import "unsafe"

func PortTestServerBrowserWords() (map[string]*uint32, func()) {
	words := map[string]*uint32{
		"dword_587000_87408":               (*uint32)(unsafe.Pointer(&C.dword_587000_87408)),
		"dword_587000_87412":               (*uint32)(unsafe.Pointer(&C.dword_587000_87412)),
		"dword_5d4594_1305788":             (*uint32)(unsafe.Pointer(&C.dword_5d4594_1305788)),
		"dword_5d4594_1307716":             (*uint32)(unsafe.Pointer(&C.dword_5d4594_1307716)),
		"dword_5d4594_1307720":             (*uint32)(unsafe.Pointer(&C.dword_5d4594_1307720)),
		"dword_5d4594_528252":              (*uint32)(unsafe.Pointer(&C.dword_5d4594_528252)),
		"dword_5d4594_528256":              (*uint32)(unsafe.Pointer(&C.dword_5d4594_528256)),
		"dword_5d4594_814548":              (*uint32)(unsafe.Pointer(&C.dword_5d4594_814548)),
		"dword_5d4594_814624":              (*uint32)(unsafe.Pointer(&C.dword_5d4594_814624)),
		"dword_5d4594_814984":              (*uint32)(unsafe.Pointer(&C.dword_5d4594_814984)),
		"dword_5d4594_814988":              (*uint32)(unsafe.Pointer(&C.dword_5d4594_814988)),
		"dword_5d4594_814992":              (*uint32)(unsafe.Pointer(&C.dword_5d4594_814992)),
		"dword_5d4594_814996":              (*uint32)(unsafe.Pointer(&C.dword_5d4594_814996)),
		"dword_5d4594_815000":              (*uint32)(unsafe.Pointer(&C.dword_5d4594_815000)),
		"dword_5d4594_815004":              (*uint32)(unsafe.Pointer(&C.dword_5d4594_815004)),
		"dword_5d4594_815016":              (*uint32)(unsafe.Pointer(&C.dword_5d4594_815016)),
		"dword_5d4594_815020":              (*uint32)(unsafe.Pointer(&C.dword_5d4594_815020)),
		"dword_5d4594_815024":              (*uint32)(unsafe.Pointer(&C.dword_5d4594_815024)),
		"dword_5d4594_815028":              (*uint32)(unsafe.Pointer(&C.dword_5d4594_815028)),
		"dword_5d4594_815032":              (*uint32)(unsafe.Pointer(&C.dword_5d4594_815032)),
		"dword_5d4594_815044":              (*uint32)(unsafe.Pointer(&C.dword_5d4594_815044)),
		"dword_5d4594_815052":              (*uint32)(unsafe.Pointer(&C.dword_5d4594_815052)),
		"dword_5d4594_815056":              (*uint32)(unsafe.Pointer(&C.dword_5d4594_815056)),
		"dword_5d4594_815096":              (*uint32)(unsafe.Pointer(&C.dword_5d4594_815096)),
		"dword_5d4594_815100":              (*uint32)(unsafe.Pointer(&C.dword_5d4594_815100)),
		"dword_5d4594_815104":              (*uint32)(unsafe.Pointer(&C.dword_5d4594_815104)),
		"nox_client_connError_814552":      (*uint32)(unsafe.Pointer(&C.nox_client_connError_814552)),
		"nox_game_createOrJoin_815048":     (*uint32)(unsafe.Pointer(&C.nox_game_createOrJoin_815048)),
		"nox_wnd_xxx_815040":               (*uint32)(unsafe.Pointer(&C.nox_wnd_xxx_815040)),
		"nox_wol_server_result_cnt_815088": (*uint32)(unsafe.Pointer(&C.nox_wol_server_result_cnt_815088)),
		"nox_wol_servers_sorting_166704":   (*uint32)(unsafe.Pointer(&C.nox_wol_servers_sorting_166704)),
		"nox_wol_wnd_gameList_815012":      (*uint32)(unsafe.Pointer(&C.nox_wol_wnd_gameList_815012)),
		"nox_wol_wnd_world_814980":         (*uint32)(unsafe.Pointer(&C.nox_wol_wnd_world_814980)),
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
		"qword_5d4594_814956": (*uint64)(unsafe.Pointer(&C.qword_5d4594_814956)),
		"qword_5d4594_815068": (*uint64)(unsafe.Pointer(&C.qword_5d4594_815068)),
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
func PortTestServerBrowserStateSet(v int32) int32 { return int32(C.sub_43AF90(C.int(v))) }
func PortTestServerBrowserErrorSet(v int32)       { C.nox_client_setConnError_43AFA0(C.int(v)) }
func PortTestServerBrowserStateGet(which int) uint32 {
	switch which {
	case 0:
		return uint32(C.sub_43AF30())
	case 1:
		return uint32(C.sub_43AF40())
	case 2:
		return uint32(C.sub_43AF80())
	case 3:
		return uint32(C.sub_43B6D0())
	case 4:
		return uint32(C.nox_client_getServerAddr_43B300())
	case 5:
		return uint32(C.nox_client_getServerPort_43B320())
	case 6:
		return uint32(C.sub_43B340())
	}
	panic(which)
}
func PortTestServerBrowserSortClick(id int) { C.nox_wol_servers_sortBtnHandler_4A0290(C.int(id)) }
func PortTestServerBrowserFormat(addr string, port uint16, out unsafe.Pointer) int {
	p := CString(addr)
	defer StrFree(p)
	return int(C.nox_sprintAddrPort_43BC80(p, C.ushort(port), (*C.char)(out)))
}

func PortTestServerBrowserRow(record unsafe.Pointer) {
	C.nox_gui_wol_newServerLine_43B7C0((*C.nox_gui_server_ent_t)(record))
}
func PortTestServerBrowserTrimName(text unsafe.Pointer, width byte) uintptr {
	return uintptr(unsafe.Pointer(C.sub_43BC10((*C.wchar2_t)(text), C.uchar(width))))
}

func PortTestServerBrowserDetails(record unsafe.Pointer) {
	C.nox_client_gui_serverInfoBlock_4394D0(C.int(uintptr(record)))
}

func PortTestServerBrowserTick() int { return int(C.sub_438770()) }
func PortTestServerBrowserNotification(which int, window unsafe.Pointer, id uint32) int {
	if which == 0 {
		return int(C.nox_xxx_windowMultiplayerSub_439E70(C.int(uintptr(window)), 22, (*C.int)(unsafe.Pointer(uintptr(id))), 0))
	}
	return int(C.sub_439D00((*C.int)(window), 22, C.uint(id), 0))
}

func PortTestServerBrowserHostDescription() uintptr { return uintptr(unsafe.Pointer(C.sub_43AA70())) }
