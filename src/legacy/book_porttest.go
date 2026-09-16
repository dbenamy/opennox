//go:build porttest

package legacy

/*
#include "defs.h"
#include "GAME2.h"
#include "client__gui__guibook.h"
extern nox_window* nox_win_unk1;
extern uint32_t nox_xxx_aNox_cfg_0_587000_132132;
extern float2 obj_5d4594_1046620;
extern uint32_t nox_xxx_aNox_cfg_0_587000_132136;
extern unsigned int nox_player_netCode_85319C;
extern int nox_win_width;
extern int nox_win_height;
extern uint32_t dword_5d4594_1046636;
extern uint32_t dword_5d4594_1046640;
extern uint32_t dword_5d4594_1046648;
extern uint32_t dword_5d4594_1046652;
extern uint32_t dword_5d4594_1046656;
extern uint32_t dword_5d4594_1046852;
extern uint32_t dword_5d4594_1046864;
extern uint32_t dword_5d4594_1046868;
extern uint32_t dword_5d4594_1046872;
extern uint32_t dword_5d4594_1046924;
extern uint32_t dword_5d4594_1046928;
extern uint32_t dword_5d4594_1046932;
extern uint32_t dword_5d4594_1046936;
extern uint32_t dword_5d4594_1046944;
extern uint32_t dword_5d4594_1046948;
extern uint32_t dword_5d4594_1046952;
extern uint32_t dword_5d4594_1046956;
extern uint32_t dword_5d4594_1047512;
extern uint32_t dword_5d4594_1047516;
extern uint32_t dword_5d4594_1047520;
extern uint32_t dword_5d4594_1047524;
extern uint32_t dword_5d4594_1047528;
extern uint32_t dword_5d4594_1047532;
extern uint32_t dword_5d4594_1047536;
extern uint32_t dword_5d4594_1047540;
extern uint32_t dword_8531A0_2576;
int nox_xxx_guiSpellSortFn_45ABC0(const void* a1, const void* a2);
int nox_xxx_bookSetColor_45AC40();
void nox_client_toggleSpellbook_45AC70();
int nox_xxx_bookHideMB_45ACA0(int a1);
int nox_xxx_guiSpellSortList_45ADF0(int a1);
void nox_xxx_book_45B010(int a1);
int nox_xxx_bookWndProc_45B070(int a1, int a2);
int nox_xxx_bookClickSpell_45B1F0();
int nox_xxx_bookClickCreature_45B200();
int nox_xxx_book_45B210(int a1, int a2);
int nox_xxx_bookChildWndProcMB_45B360(uint32_t* a1, unsigned int a2);
int nox_xxx_bookListWndProc_45B5F0(int a1, unsigned int a2, unsigned int a3);
int nox_xxx_bookMoveToPage_45B930(int a1);
int nox_xxx_book_45BD30(int a1, int a2);
int nox_xxx_bookInit_45B9D0();
int nox_xxx_bookDrawIconFn_45CB30(uint32_t* a1);
int nox_xxx_bookWndFn_45CC10(uint32_t* a1, int a2, unsigned int a3);
int sub_45CFC0();
void nox_xxx_netSpellRewardCli_45CFE0(int a1, int a2, int a3, int a4);
void nox_xxx_netGuideRewardCli_45D140(int a1, int a2);
int* nox_xxx_bookSetForward_45D200(int* a1, int a2, int2* a3);
void nox_xxx_abilityReward_45D290(int a1, char* a2, int a3);
int sub_45D320(int a1);
int sub_45D400(int a1);
char* nox_xxx_clientQuestDisableAbility_45D4A0(int a1);
int sub_45D500(int a1);
uint32_t* sub_45D550(uint32_t* a1);
void nox_xxx_bookFillAll_45D570(int a1, int a2);
int sub_45D7D0(int* a1, int* a2);
void sub_45D810();
int sub_45D9B0();
void nox_xxx_bookShowMB_45AD70(int a1);
int nox_xxx_bookDrawList_45BD40(int a1);
int nox_xxx_book_45CF00(uint32_t* a1);
int nox_xxx_bookDrawFn_45C7D0(uint32_t* a1);
void sub_45D870();
*/
import "C"
import "unsafe"

func PortTestBookInvoke(op string, a [4]uint32) uint32 {
	switch op {
	case "nox_xxx_guiSpellSortFn_45ABC0":
		return uint32(C.nox_xxx_guiSpellSortFn_45ABC0(unsafe.Pointer(uintptr(a[0])), unsafe.Pointer(uintptr(a[1]))))
	case "nox_xxx_bookSetColor_45AC40":
		return uint32(C.nox_xxx_bookSetColor_45AC40())
	case "nox_client_toggleSpellbook_45AC70":
		C.nox_client_toggleSpellbook_45AC70()
		return 0
	case "nox_xxx_bookHideMB_45ACA0":
		return uint32(C.nox_xxx_bookHideMB_45ACA0(C.int(a[0])))
	case "nox_xxx_guiSpellSortList_45ADF0":
		return uint32(C.nox_xxx_guiSpellSortList_45ADF0(C.int(a[0])))
	case "nox_xxx_book_45B010":
		C.nox_xxx_book_45B010(C.int(a[0]))
		return 0
	case "nox_xxx_bookWndProc_45B070":
		return uint32(C.nox_xxx_bookWndProc_45B070(C.int(a[0]), C.int(a[1])))
	case "nox_xxx_bookClickSpell_45B1F0":
		return uint32(C.nox_xxx_bookClickSpell_45B1F0())
	case "nox_xxx_bookClickCreature_45B200":
		return uint32(C.nox_xxx_bookClickCreature_45B200())
	case "nox_xxx_book_45B210":
		return uint32(C.nox_xxx_book_45B210(C.int(a[0]), C.int(a[1])))
	case "nox_xxx_bookChildWndProcMB_45B360":
		return uint32(C.nox_xxx_bookChildWndProcMB_45B360((*C.uint)(unsafe.Pointer(uintptr(a[0]))), C.uint(a[1])))
	case "nox_xxx_bookListWndProc_45B5F0":
		return uint32(C.nox_xxx_bookListWndProc_45B5F0(C.int(a[0]), C.uint(a[1]), C.uint(a[2])))
	case "nox_xxx_bookMoveToPage_45B930":
		return uint32(C.nox_xxx_bookMoveToPage_45B930(C.int(a[0])))
	case "nox_xxx_book_45BD30":
		return uint32(C.nox_xxx_book_45BD30(C.int(a[0]), C.int(a[1])))
	case "nox_xxx_bookInit_45B9D0":
		return uint32(C.nox_xxx_bookInit_45B9D0())
	case "nox_xxx_bookDrawIconFn_45CB30":
		return uint32(C.nox_xxx_bookDrawIconFn_45CB30((*C.uint)(unsafe.Pointer(uintptr(a[0])))))
	case "nox_xxx_bookWndFn_45CC10":
		return uint32(C.nox_xxx_bookWndFn_45CC10((*C.uint)(unsafe.Pointer(uintptr(a[0]))), C.int(a[1]), C.uint(a[2])))
	case "sub_45CFC0":
		return uint32(C.sub_45CFC0())
	case "nox_xxx_netSpellRewardCli_45CFE0":
		C.nox_xxx_netSpellRewardCli_45CFE0(C.int(a[0]), C.int(a[1]), C.int(a[2]), C.int(a[3]))
		return 0
	case "nox_xxx_netGuideRewardCli_45D140":
		C.nox_xxx_netGuideRewardCli_45D140(C.int(a[0]), C.int(a[1]))
		return 0
	case "nox_xxx_bookSetForward_45D200":
		return uint32(uintptr(unsafe.Pointer(C.nox_xxx_bookSetForward_45D200((*C.int)(unsafe.Pointer(uintptr(a[0]))), C.int(a[1]), (*C.int2)(unsafe.Pointer(uintptr(a[2])))))))
	case "nox_xxx_abilityReward_45D290":
		C.nox_xxx_abilityReward_45D290(C.int(a[0]), (*C.char)(unsafe.Pointer(uintptr(a[1]))), C.int(a[2]))
		return 0
	case "sub_45D320":
		return uint32(C.sub_45D320(C.int(a[0])))
	case "sub_45D400":
		return uint32(C.sub_45D400(C.int(a[0])))
	case "nox_xxx_clientQuestDisableAbility_45D4A0":
		return uint32(uintptr(unsafe.Pointer(C.nox_xxx_clientQuestDisableAbility_45D4A0(C.int(a[0])))))
	case "sub_45D500":
		return uint32(C.sub_45D500(C.int(a[0])))
	case "sub_45D550":
		return uint32(uintptr(unsafe.Pointer(C.sub_45D550((*C.uint)(unsafe.Pointer(uintptr(a[0])))))))
	case "nox_xxx_bookFillAll_45D570":
		C.nox_xxx_bookFillAll_45D570(C.int(a[0]), C.int(a[1]))
		return 0
	case "sub_45D7D0":
		return uint32(C.sub_45D7D0((*C.int)(unsafe.Pointer(uintptr(a[0]))), (*C.int)(unsafe.Pointer(uintptr(a[1])))))
	case "sub_45D810":
		C.sub_45D810()
		return 0
	case "sub_45D9B0":
		return uint32(C.sub_45D9B0())
	case "nox_xxx_bookShowMB_45AD70":
		C.nox_xxx_bookShowMB_45AD70(C.int(a[0]))
		return 0
	case "nox_xxx_bookDrawList_45BD40":
		return uint32(C.nox_xxx_bookDrawList_45BD40(C.int(a[0])))
	case "nox_xxx_book_45CF00":
		return uint32(C.nox_xxx_book_45CF00((*C.uint)(unsafe.Pointer(uintptr(a[0])))))
	case "nox_xxx_bookDrawFn_45C7D0":
		return uint32(C.nox_xxx_bookDrawFn_45C7D0((*C.uint)(unsafe.Pointer(uintptr(a[0])))))
	case "sub_45D870":
		C.sub_45D870()
		return 0
	}
	panic(op)
}
func PortTestBookCallbacks() map[string]unsafe.Pointer {
	return map[string]unsafe.Pointer{
		"nox_xxx_guiSpellSortFn_45ABC0":            unsafe.Pointer(C.nox_xxx_guiSpellSortFn_45ABC0),
		"nox_xxx_bookSetColor_45AC40":              unsafe.Pointer(C.nox_xxx_bookSetColor_45AC40),
		"nox_client_toggleSpellbook_45AC70":        unsafe.Pointer(C.nox_client_toggleSpellbook_45AC70),
		"nox_xxx_bookHideMB_45ACA0":                unsafe.Pointer(C.nox_xxx_bookHideMB_45ACA0),
		"nox_xxx_guiSpellSortList_45ADF0":          unsafe.Pointer(C.nox_xxx_guiSpellSortList_45ADF0),
		"nox_xxx_book_45B010":                      unsafe.Pointer(C.nox_xxx_book_45B010),
		"nox_xxx_bookWndProc_45B070":               unsafe.Pointer(C.nox_xxx_bookWndProc_45B070),
		"nox_xxx_bookClickSpell_45B1F0":            unsafe.Pointer(C.nox_xxx_bookClickSpell_45B1F0),
		"nox_xxx_bookClickCreature_45B200":         unsafe.Pointer(C.nox_xxx_bookClickCreature_45B200),
		"nox_xxx_book_45B210":                      unsafe.Pointer(C.nox_xxx_book_45B210),
		"nox_xxx_bookChildWndProcMB_45B360":        unsafe.Pointer(C.nox_xxx_bookChildWndProcMB_45B360),
		"nox_xxx_bookListWndProc_45B5F0":           unsafe.Pointer(C.nox_xxx_bookListWndProc_45B5F0),
		"nox_xxx_bookMoveToPage_45B930":            unsafe.Pointer(C.nox_xxx_bookMoveToPage_45B930),
		"nox_xxx_book_45BD30":                      unsafe.Pointer(C.nox_xxx_book_45BD30),
		"nox_xxx_bookInit_45B9D0":                  unsafe.Pointer(C.nox_xxx_bookInit_45B9D0),
		"nox_xxx_bookDrawIconFn_45CB30":            unsafe.Pointer(C.nox_xxx_bookDrawIconFn_45CB30),
		"nox_xxx_bookWndFn_45CC10":                 unsafe.Pointer(C.nox_xxx_bookWndFn_45CC10),
		"sub_45CFC0":                               unsafe.Pointer(C.sub_45CFC0),
		"nox_xxx_netSpellRewardCli_45CFE0":         unsafe.Pointer(C.nox_xxx_netSpellRewardCli_45CFE0),
		"nox_xxx_netGuideRewardCli_45D140":         unsafe.Pointer(C.nox_xxx_netGuideRewardCli_45D140),
		"nox_xxx_bookSetForward_45D200":            unsafe.Pointer(C.nox_xxx_bookSetForward_45D200),
		"nox_xxx_abilityReward_45D290":             unsafe.Pointer(C.nox_xxx_abilityReward_45D290),
		"sub_45D320":                               unsafe.Pointer(C.sub_45D320),
		"sub_45D400":                               unsafe.Pointer(C.sub_45D400),
		"nox_xxx_clientQuestDisableAbility_45D4A0": unsafe.Pointer(C.nox_xxx_clientQuestDisableAbility_45D4A0),
		"sub_45D500":                               unsafe.Pointer(C.sub_45D500),
		"sub_45D550":                               unsafe.Pointer(C.sub_45D550),
		"nox_xxx_bookFillAll_45D570":               unsafe.Pointer(C.nox_xxx_bookFillAll_45D570),
		"sub_45D7D0":                               unsafe.Pointer(C.sub_45D7D0),
		"sub_45D810":                               unsafe.Pointer(C.sub_45D810),
		"sub_45D9B0":                               unsafe.Pointer(C.sub_45D9B0),
		"nox_xxx_bookShowMB_45AD70":                unsafe.Pointer(C.nox_xxx_bookShowMB_45AD70),
		"nox_xxx_bookDrawList_45BD40":              unsafe.Pointer(C.nox_xxx_bookDrawList_45BD40),
		"nox_xxx_book_45CF00":                      unsafe.Pointer(C.nox_xxx_book_45CF00),
		"nox_xxx_bookDrawFn_45C7D0":                unsafe.Pointer(C.nox_xxx_bookDrawFn_45C7D0),
		"sub_45D870":                               unsafe.Pointer(C.sub_45D870),
	}
}
func PortTestBookWords() (map[string]*uint32, func()) {
	words := map[string]*uint32{
		"nox_win_unk1":                     (*uint32)(unsafe.Pointer(&C.nox_win_unk1)),
		"nox_xxx_aNox_cfg_0_587000_132136": (*uint32)(unsafe.Pointer(&C.nox_xxx_aNox_cfg_0_587000_132136)),
		"nox_player_netCode_85319C":        (*uint32)(unsafe.Pointer(&C.nox_player_netCode_85319C)),
		"nox_win_width":                    (*uint32)(unsafe.Pointer(&C.nox_win_width)),
		"nox_win_height":                   (*uint32)(unsafe.Pointer(&C.nox_win_height)),
		"nox_xxx_aNox_cfg_0_587000_132132": (*uint32)(unsafe.Pointer(&C.nox_xxx_aNox_cfg_0_587000_132132)),
		"dword_5d4594_1046636":             (*uint32)(unsafe.Pointer(&C.dword_5d4594_1046636)),
		"dword_5d4594_1046640":             (*uint32)(unsafe.Pointer(&C.dword_5d4594_1046640)),
		"dword_5d4594_1046648":             (*uint32)(unsafe.Pointer(&C.dword_5d4594_1046648)),
		"dword_5d4594_1046652":             (*uint32)(unsafe.Pointer(&C.dword_5d4594_1046652)),
		"dword_5d4594_1046656":             (*uint32)(unsafe.Pointer(&C.dword_5d4594_1046656)),
		"dword_5d4594_1046852":             (*uint32)(unsafe.Pointer(&C.dword_5d4594_1046852)),
		"dword_5d4594_1046864":             (*uint32)(unsafe.Pointer(&C.dword_5d4594_1046864)),
		"dword_5d4594_1046868":             (*uint32)(unsafe.Pointer(&C.dword_5d4594_1046868)),
		"dword_5d4594_1046872":             (*uint32)(unsafe.Pointer(&C.dword_5d4594_1046872)),
		"dword_5d4594_1046924":             (*uint32)(unsafe.Pointer(&C.dword_5d4594_1046924)),
		"dword_5d4594_1046928":             (*uint32)(unsafe.Pointer(&C.dword_5d4594_1046928)),
		"dword_5d4594_1046932":             (*uint32)(unsafe.Pointer(&C.dword_5d4594_1046932)),
		"dword_5d4594_1046936":             (*uint32)(unsafe.Pointer(&C.dword_5d4594_1046936)),
		"dword_5d4594_1046944":             (*uint32)(unsafe.Pointer(&C.dword_5d4594_1046944)),
		"dword_5d4594_1046948":             (*uint32)(unsafe.Pointer(&C.dword_5d4594_1046948)),
		"dword_5d4594_1046952":             (*uint32)(unsafe.Pointer(&C.dword_5d4594_1046952)),
		"dword_5d4594_1046956":             (*uint32)(unsafe.Pointer(&C.dword_5d4594_1046956)),
		"dword_5d4594_1047512":             (*uint32)(unsafe.Pointer(&C.dword_5d4594_1047512)),
		"dword_5d4594_1047516":             (*uint32)(unsafe.Pointer(&C.dword_5d4594_1047516)),
		"dword_5d4594_1047520":             (*uint32)(unsafe.Pointer(&C.dword_5d4594_1047520)),
		"dword_5d4594_1047524":             (*uint32)(unsafe.Pointer(&C.dword_5d4594_1047524)),
		"dword_5d4594_1047528":             (*uint32)(unsafe.Pointer(&C.dword_5d4594_1047528)),
		"dword_5d4594_1047532":             (*uint32)(unsafe.Pointer(&C.dword_5d4594_1047532)),
		"dword_5d4594_1047536":             (*uint32)(unsafe.Pointer(&C.dword_5d4594_1047536)),
		"dword_5d4594_1047540":             (*uint32)(unsafe.Pointer(&C.dword_5d4594_1047540)),
		"dword_8531A0_2576":                (*uint32)(unsafe.Pointer(&C.dword_8531A0_2576)),
	}
	saved := make(map[string]uint32, len(words))
	for name, p := range words {
		saved[name] = *p
	}
	return words, func() {
		for name, p := range words {
			*p = saved[name]
		}
	}
}
func PortTestBookVector() *[2]float32 { return (*[2]float32)(unsafe.Pointer(&C.obj_5d4594_1046620)) }
