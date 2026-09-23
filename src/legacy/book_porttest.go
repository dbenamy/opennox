//go:build porttest

package legacy

/*
#include "defs.h"
#include "GAME2.h"
#include "client__gui__guibook.h"

void nox_client_toggleSpellbook_45AC70();
int nox_xxx_bookHideMB_45ACA0(int a1);
int nox_xxx_bookClickSpell_45B1F0();
int nox_xxx_bookClickCreature_45B200();
int sub_45CFC0();
int* nox_xxx_bookSetForward_45D200(int* a1, int a2, int2* a3);
void nox_xxx_abilityReward_45D290(int a1, char* a2, int a3);
int sub_45D500(int a1);
void nox_xxx_bookFillAll_45D570(int a1, int a2);
int sub_45D9B0();
int nox_xxx_book_45CF00(uint32_t* a1);
void sub_45D870();
*/
import "C"
import "unsafe"
import "image"

func PortTestBookInvoke(op string, a [4]uint32) uint32 {
	switch op {
	case "nox_xxx_guiSpellSortFn_45ABC0":
		return uint32(bookCompare(*(*uint32)(unsafe.Pointer(uintptr(a[0]))), *(*uint32)(unsafe.Pointer(uintptr(a[1])))))
	case "nox_xxx_bookSetColor_45AC40":
		return uint32(bookSetColor())
	case "nox_client_toggleSpellbook_45AC70":
		C.nox_client_toggleSpellbook_45AC70()
		return 0
	case "nox_xxx_bookHideMB_45ACA0":
		return uint32(C.nox_xxx_bookHideMB_45ACA0(C.int(a[0])))
	case "nox_xxx_guiSpellSortList_45ADF0":
		return uint32(bookSort(int(a[0])))
	case "nox_xxx_book_45B010":
		bookOpen(int(a[0]))
		return 0
	case "nox_xxx_bookWndProc_45B070":
		return uint32(bookForward(int(a[1])))
	case "nox_xxx_bookClickSpell_45B1F0":
		return uint32(C.nox_xxx_bookClickSpell_45B1F0())
	case "nox_xxx_bookClickCreature_45B200":
		return uint32(C.nox_xxx_bookClickCreature_45B200())
	case "nox_xxx_book_45B210":
		return uint32(bookBackward(int(a[1])))
	case "nox_xxx_bookChildWndProcMB_45B360":
		return uint32(bookTab(bookWindow(a[0]), a[1]))
	case "nox_xxx_bookListWndProc_45B5F0":
		return uint32(bookListEvents(bookWindow(a[0]), a[1], bookPoint(a[2])))
	case "nox_xxx_bookMoveToPage_45B930":
		return uint32(bookMoveToPage(int(a[0])))
	case "nox_xxx_book_45BD30":
		return 1
	case "nox_xxx_bookInit_45B9D0":
		return uint32(bookInit())
	case "nox_xxx_bookDrawIconFn_45CB30":
		return uint32(bookDrawIcon(bookWindow(a[0])))
	case "nox_xxx_bookWndFn_45CC10":
		return uint32(bookIconEvents(bookWindow(a[0]), int(a[1]), bookPoint(a[2])))
	case "sub_45CFC0":
		return uint32(C.sub_45CFC0())
	case "nox_xxx_netSpellRewardCli_45CFE0":
		nox_xxx_netSpellRewardCli_45CFE0(C.int(a[0]), C.int(a[1]), C.int(a[2]), C.int(a[3]))
		return 0
	case "nox_xxx_netGuideRewardCli_45D140":
		nox_xxx_netGuideRewardCli_45D140(C.int(a[0]), C.int(a[1]))
		return 0
	case "nox_xxx_bookSetForward_45D200":
		return uint32(uintptr(unsafe.Pointer(C.nox_xxx_bookSetForward_45D200((*C.int)(unsafe.Pointer(uintptr(a[0]))), C.int(a[1]), (*C.int2)(unsafe.Pointer(uintptr(a[2])))))))
	case "nox_xxx_abilityReward_45D290":
		C.nox_xxx_abilityReward_45D290(C.int(a[0]), (*C.char)(unsafe.Pointer(uintptr(a[1]))), C.int(a[2]))
		return 0
	case "sub_45D320":
		return uint32(sub_45D320(C.int(a[0])))
	case "sub_45D400":
		return uint32(sub_45D400(C.int(a[0])))
	case "nox_xxx_clientQuestDisableAbility_45D4A0":
		return uint32(uintptr(unsafe.Pointer(nox_xxx_clientQuestDisableAbility_45D4A0(C.int(a[0])))))
	case "sub_45D500":
		return uint32(C.sub_45D500(C.int(a[0])))
	case "sub_45D550":
		return uint32(bookIconPosition((*image.Point)(unsafe.Pointer(uintptr(a[0])))))
	case "nox_xxx_bookFillAll_45D570":
		C.nox_xxx_bookFillAll_45D570(C.int(a[0]), C.int(a[1]))
		return 0
	case "sub_45D7D0":
		return uint32(bookPath(AsPoint(unsafe.Pointer(uintptr(a[0]))), AsPoint(unsafe.Pointer(uintptr(a[1])))))
	case "sub_45D810":
		bookStopAddition()
		return 0
	case "sub_45D9B0":
		return uint32(C.sub_45D9B0())
	case "nox_xxx_bookShowMB_45AD70":
		bookShow(int(a[0]))
		return 0
	case "nox_xxx_bookDrawList_45BD40":
		return uint32(bookDrawList(bookWindow(a[0])))
	case "nox_xxx_book_45CF00":
		return uint32(C.nox_xxx_book_45CF00((*C.uint)(unsafe.Pointer(uintptr(a[0])))))
	case "nox_xxx_bookDrawFn_45C7D0":
		return uint32(bookDrawAddition(bookWindow(a[0])))
	case "sub_45D870":
		C.sub_45D870()
		return 0
	}
	panic(op)
}
func PortTestBookCallbacks() map[string]unsafe.Pointer {
	return map[string]unsafe.Pointer{
		"nox_xxx_guiSpellSortFn_45ABC0":            nil,
		"nox_xxx_bookSetColor_45AC40":              nil,
		"nox_client_toggleSpellbook_45AC70":        unsafe.Pointer(C.nox_client_toggleSpellbook_45AC70),
		"nox_xxx_bookHideMB_45ACA0":                unsafe.Pointer(C.nox_xxx_bookHideMB_45ACA0),
		"nox_xxx_guiSpellSortList_45ADF0":          nil,
		"nox_xxx_book_45B010":                      nil,
		"nox_xxx_bookWndProc_45B070":               nil,
		"nox_xxx_bookClickSpell_45B1F0":            unsafe.Pointer(C.nox_xxx_bookClickSpell_45B1F0),
		"nox_xxx_bookClickCreature_45B200":         unsafe.Pointer(C.nox_xxx_bookClickCreature_45B200),
		"nox_xxx_book_45B210":                      nil,
		"nox_xxx_bookChildWndProcMB_45B360":        nil,
		"nox_xxx_bookListWndProc_45B5F0":           nil,
		"nox_xxx_bookMoveToPage_45B930":            nil,
		"nox_xxx_book_45BD30":                      nil,
		"nox_xxx_bookInit_45B9D0":                  nil,
		"nox_xxx_bookDrawIconFn_45CB30":            nil,
		"nox_xxx_bookWndFn_45CC10":                 nil,
		"sub_45CFC0":                               unsafe.Pointer(C.sub_45CFC0),
		"nox_xxx_netSpellRewardCli_45CFE0":         nil, // Preserve stable capture IDs after retiring the C export.
		"nox_xxx_netGuideRewardCli_45D140":         nil,
		"nox_xxx_bookSetForward_45D200":            unsafe.Pointer(C.nox_xxx_bookSetForward_45D200),
		"nox_xxx_abilityReward_45D290":             unsafe.Pointer(C.nox_xxx_abilityReward_45D290),
		"sub_45D320":                               nil,
		"sub_45D400":                               nil,
		"nox_xxx_clientQuestDisableAbility_45D4A0": nil,
		"sub_45D500":                               unsafe.Pointer(C.sub_45D500),
		"sub_45D550":                               nil,
		"nox_xxx_bookFillAll_45D570":               unsafe.Pointer(C.nox_xxx_bookFillAll_45D570),
		"sub_45D7D0":                               nil,
		"sub_45D810":                               nil,
		"sub_45D9B0":                               unsafe.Pointer(C.sub_45D9B0),
		"nox_xxx_bookShowMB_45AD70":                nil,
		"nox_xxx_bookDrawList_45BD40":              nil,
		"nox_xxx_book_45CF00":                      unsafe.Pointer(C.nox_xxx_book_45CF00),
		"nox_xxx_bookDrawFn_45C7D0":                nil,
		"sub_45D870":                               unsafe.Pointer(C.sub_45D870),
	}
}
func PortTestBookWords() (map[string]*uint32, func()) {
	words := map[string]*uint32{
		"nox_win_unk1":                     (*uint32)(unsafe.Pointer(&legacyGlobals.nox_win_unk1)),
		"nox_xxx_aNox_cfg_0_587000_132136": (*uint32)(unsafe.Pointer(&nox_xxx_aNox_cfg_0_587000_132136)),
		"nox_player_netCode_85319C":        (*uint32)(unsafe.Pointer(&nox_player_netCode_85319C)),
		"nox_win_width":                    (*uint32)(unsafe.Pointer(&nox_win_width)),
		"nox_win_height":                   (*uint32)(unsafe.Pointer(&nox_win_height)),
		"nox_xxx_aNox_cfg_0_587000_132132": (*uint32)(unsafe.Pointer(&nox_xxx_aNox_cfg_0_587000_132132)),
		"dword_5d4594_1046636":             (*uint32)(unsafe.Pointer(&dword_5d4594_1046636)),
		"dword_5d4594_1046640":             (*uint32)(unsafe.Pointer(&dword_5d4594_1046640)),
		"dword_5d4594_1046648":             (*uint32)(unsafe.Pointer(&dword_5d4594_1046648)),
		"dword_5d4594_1046652":             (*uint32)(unsafe.Pointer(&dword_5d4594_1046652)),
		"dword_5d4594_1046656":             (*uint32)(unsafe.Pointer(&dword_5d4594_1046656)),
		"dword_5d4594_1046852":             (*uint32)(unsafe.Pointer(&dword_5d4594_1046852)),
		"dword_5d4594_1046864":             (*uint32)(unsafe.Pointer(&dword_5d4594_1046864)),
		"dword_5d4594_1046868":             (*uint32)(unsafe.Pointer(&dword_5d4594_1046868)),
		"dword_5d4594_1046872":             (*uint32)(unsafe.Pointer(&dword_5d4594_1046872)),
		"dword_5d4594_1046924":             (*uint32)(unsafe.Pointer(&dword_5d4594_1046924)),
		"dword_5d4594_1046928":             (*uint32)(unsafe.Pointer(&dword_5d4594_1046928)),
		"dword_5d4594_1046932":             (*uint32)(unsafe.Pointer(&dword_5d4594_1046932)),
		"dword_5d4594_1046936":             (*uint32)(unsafe.Pointer(&dword_5d4594_1046936)),
		"dword_5d4594_1046944":             (*uint32)(unsafe.Pointer(&dword_5d4594_1046944)),
		"dword_5d4594_1046948":             (*uint32)(unsafe.Pointer(&dword_5d4594_1046948)),
		"dword_5d4594_1046952":             (*uint32)(unsafe.Pointer(&dword_5d4594_1046952)),
		"dword_5d4594_1046956":             (*uint32)(unsafe.Pointer(&dword_5d4594_1046956)),
		"dword_5d4594_1047512":             (*uint32)(unsafe.Pointer(&dword_5d4594_1047512)),
		"dword_5d4594_1047516":             (*uint32)(unsafe.Pointer(&dword_5d4594_1047516)),
		"dword_5d4594_1047520":             (*uint32)(unsafe.Pointer(&dword_5d4594_1047520)),
		"dword_5d4594_1047524":             (*uint32)(unsafe.Pointer(&dword_5d4594_1047524)),
		"dword_5d4594_1047528":             (*uint32)(unsafe.Pointer(&dword_5d4594_1047528)),
		"dword_5d4594_1047532":             (*uint32)(unsafe.Pointer(&dword_5d4594_1047532)),
		"dword_5d4594_1047536":             (*uint32)(unsafe.Pointer(&dword_5d4594_1047536)),
		"dword_5d4594_1047540":             (*uint32)(unsafe.Pointer(&dword_5d4594_1047540)),
		"dword_8531A0_2576":                (*uint32)(unsafe.Pointer(&dword_8531A0_2576)),
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
func PortTestBookVector() *[2]float32 { return &bookVector }
