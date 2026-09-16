//go:build porttest

package legacy

/*
#include "defs.h"
#include "GAME2.h"
#include "GAME2_1.h"
#include "client__gui__guispell.h"
*/
import "C"
import "unsafe"

// PortTestQuickbarInvoke calls the production quickbar owner without duplicating
// any of its algorithms. Pointer arguments use the target's 32-bit ABI.
func PortTestQuickbarInvoke(op string, a [7]uint32) uint32 {
	switch op {
	case "nox_xxx_guiSpellTargetClickSet_45D9D0":
		return uint32(C.nox_xxx_guiSpellTargetClickSet_45D9D0(C.int(a[0])))
	case "nox_xxx_guiSpell_45DA10":
		return uint32(C.nox_xxx_guiSpell_45DA10(C.int(a[0])))
	case "nox_client_invokeSpellSlot_45DA50":
		C.nox_client_invokeSpellSlot_45DA50(C.int(a[0]))
		return 0
	case "nox_xxx_clientStoreLastButton_45DAD0":
		C.nox_xxx_clientStoreLastButton_45DAD0(C.int(a[0]))
		return 0
	case "nox_xxx_clientSendAbil_45DAF0":
		return uint32(C.nox_xxx_clientSendAbil_45DAF0(C.int(a[0])))
	case "nox_xxx_clientSendSpell_45DB20":
		return uint32(C.nox_xxx_clientSendSpell_45DB20((*C.char)(unsafe.Pointer(uintptr(a[0]))), C.int(a[1]), C.char(a[2])))
	case "sub_45DB90":
		return uint32(C.sub_45DB90())
	case "nox_xxx_guiSpellTargetClickCheckSend_45DBB0":
		C.nox_xxx_guiSpellTargetClickCheckSend_45DBB0()
		return 0
	case "nox_xxx_book_45DBE0":
		return uint32(uintptr(unsafe.Pointer(C.nox_xxx_book_45DBE0(unsafe.Pointer(uintptr(a[0])), C.int(a[1]), C.int(a[2])))))
	case "nox_xxx_spellKeyPackSetSpell_45DC40":
		C.nox_xxx_spellKeyPackSetSpell_45DC40(C.int(a[0]), C.int(a[1]), C.int(a[2]))
		return 0
	case "nox_xxx_bookSpellDrop_45DCA0":
		return uint32(C.nox_xxx_bookSpellDrop_45DCA0(C.int(a[0]), C.char(a[1]), C.int(a[2]), C.int(a[3]), (*C.int)(unsafe.Pointer(uintptr(a[4])))))
	case "nox_xxx_updateSpellIcons_45DDF0":
		return uint32(C.nox_xxx_updateSpellIcons_45DDF0(C.int(a[0])))
	case "nox_xxx_updateSpellIconDir_45DE10":
		return uint32(C.nox_xxx_updateSpellIconDir_45DE10(C.int(a[0]), C.int(a[1])))
	case "nox_xxx_spellBoxPointToWnd_45DE60":
		return uint32(C.nox_xxx_spellBoxPointToWnd_45DE60(C.int(a[0]), C.int(a[1]), C.int(a[2])))
	case "nox_xxx_guiSpellSetCursor_45DF60":
		return uint32(C.nox_xxx_guiSpellSetCursor_45DF60(C.int(a[0]), C.char(a[1])))
	case "sub_45DFC0":
		return uint32(C.sub_45DFC0(C.int(a[0])))
	case "nox_xxx_clientUpdateButtonRow_45E110":
		return uint32(C.nox_xxx_clientUpdateButtonRow_45E110(C.int(a[0])))
	case "nox_xxx_buttonsGetSelectedRow_45E180":
		return uint32(C.nox_xxx_buttonsGetSelectedRow_45E180())
	case "nox_xxx_quickbarButtonBookDraw_45EF30":
		return uint32(C.nox_xxx_quickbarButtonBookDraw_45EF30())
	case "sub_45EF40":
		return uint32(C.sub_45EF40())
	case "nox_xxx_quickBarWnd_45EF50":
		return uint32(C.nox_xxx_quickBarWnd_45EF50(C.int(a[0]), C.int(a[1]), C.uint(a[2])))
	case "nox_xxx_clientSwapQuickbarKeys_45F300":
		C.nox_xxx_clientSwapQuickbarKeys_45F300(C.int(a[0]), C.int(a[1]), C.int(a[2]))
		return 0
	case "sub_45F390":
		C.sub_45F390(C.int(a[0]), C.int(a[1]), C.int(a[2]), C.int(a[3]))
		return 0
	case "nox_xxx_quickbarButtonBookWnd_45F450":
		return uint32(C.nox_xxx_quickbarButtonBookWnd_45F450(C.int(a[0]), C.uint(a[1])))
	case "sub_45F500":
		return uint32(C.sub_45F500(C.int(a[0]), C.int(a[1])))
	case "sub_45F520":
		return uint32(C.sub_45F520(C.int(a[0]), C.uint(a[1])))
	case "sub_45F5D0":
		return uint32(C.sub_45F5D0((*C.uint32_t)(unsafe.Pointer(uintptr(a[0])))))
	case "nox_xxx_quickbarTrapUpDownProc_45F630":
		return uint32(C.nox_xxx_quickbarTrapUpDownProc_45F630(C.int(a[0]), C.uint(a[1])))
	case "nox_xxx_quickbarTrapUpDownDraw_45F6F0":
		return uint32(C.nox_xxx_quickbarTrapUpDownDraw_45F6F0((*C.uint32_t)(unsafe.Pointer(uintptr(a[0])))))
	case "nox_xxx_quickbarTrapButtonProc_45F7A0":
		return uint32(C.nox_xxx_quickbarTrapButtonProc_45F7A0((*C.uint32_t)(unsafe.Pointer(uintptr(a[0]))), C.uint(a[1]), C.uint(a[2])))
	case "sub_45F890":
		return uint32(C.sub_45F890())
	case "nox_xxx_quickbar_45F8D0":
		return uint32(C.nox_xxx_quickbar_45F8D0(C.int(a[0]), C.int(a[1]), C.int(a[2]), C.int(a[3])))
	case "sub_45F8F0":
		return uint32(C.sub_45F8F0())
	case "sub_45F900":
		C.sub_45F900()
		return 0
	case "nox_xxx_quickbarTrapProc_45FB90":
		return uint32(C.nox_xxx_quickbarTrapProc_45FB90(C.int(a[0]), C.uint(a[1]), C.int(a[2]), C.int(a[3])))
	case "nox_xxx_quickbarDrawFn_460000":
		return uint32(C.nox_xxx_quickbarDrawFn_460000())
	case "nox_xxx_quickBarInitWindow_4601F0":
		return uint32(C.nox_xxx_quickBarInitWindow_4601F0(C.int(a[0]), C.int(a[1]), C.int(a[2]), C.int(a[3]), C.int(a[4]), C.int(a[5]), C.int(a[6])))
	case "sub_4602F0":
		return uint32(uintptr(unsafe.Pointer(C.sub_4602F0())))
	case "sub_460380":
		return uint32(uintptr(unsafe.Pointer(C.sub_460380())))
	case "nox_client_trapSetNext_4603A0":
		C.nox_client_trapSetNext_4603A0()
		return 0
	case "nox_client_trapSetPrev_4603F0":
		C.nox_client_trapSetPrev_4603F0()
		return 0
	case "nox_client_trapSetSelect_4604B0":
		return uint32(C.nox_client_trapSetSelect_4604B0(C.int(a[0])))
	case "sub_4604E0":
		return uint32(C.sub_4604E0())
	case "nox_client_spellSetNext_4604F0":
		return uint32(C.nox_client_spellSetNext_4604F0())
	case "nox_client_spellSetPrev_460540":
		return uint32(C.nox_client_spellSetPrev_460540())
	case "nox_client_spellSetSelect_460590":
		C.nox_client_spellSetSelect_460590()
		return 0
	case "sub_460630":
		C.sub_460630()
		return 0
	case "nox_xxx_guiSpell_460650":
		return uint32(C.nox_xxx_guiSpell_460650())
	case "sub_460660":
		return uint32(C.sub_460660())
	case "nox_xxx_quickBarClose_4606B0":
		return uint32(C.nox_xxx_quickBarClose_4606B0())
	case "sub_460920":
		C.sub_460920()
		return 0
	case "sub_460940":
		return uint32(C.sub_460940(unsafe.Pointer(uintptr(a[0]))))
	case "sub_460A10":
		return uint32(C.sub_460A10(C.int(a[0]), C.int(a[1]), C.int(a[2]), C.char(a[3])))
	case "sub_460B90":
		return uint32(C.sub_460B90(C.int(a[0])))
	case "sub_460D40":
		return uint32(C.sub_460D40())
	case "sub_460D50":
		return uint32(C.sub_460D50())
	case "nox_xxx_cliPrepareGameplay1_460E60":
		return uint32(C.nox_xxx_cliPrepareGameplay1_460E60())
	case "sub_460EA0":
		return uint32(C.sub_460EA0(C.int(a[0])))
	case "sub_460EB0":
		C.sub_460EB0(C.int(a[0]), C.char(a[1]))
		return 0
	case "sub_461010":
		C.sub_461010()
		return 0
	case "sub_461060":
		C.sub_461060()
		return 0
	case "sub_461090":
		return uint32(uintptr(unsafe.Pointer(C.sub_461090(C.int(a[0]), C.int(a[1])))))
	case "sub_4610D0":
		return uint32(uintptr(unsafe.Pointer(C.sub_4610D0(C.uchar(a[0])))))
	case "sub_461120":
		return uint32(uintptr(unsafe.Pointer(C.sub_461120(C.int(a[0]), C.int(a[1])))))
	case "sub_461160":
		return uint32(C.sub_461160(C.int(a[0])))
	case "sub_4611A0":
		return uint32(C.sub_4611A0())
	case "sub_4611B0":
		return uint32(C.sub_4611B0())
	case "nox_xxx_netAbilityRewardCli_4611E0":
		C.nox_xxx_netAbilityRewardCli_4611E0(C.int(a[0]), C.int(a[1]), (*C.char)(unsafe.Pointer(uintptr(a[2]))))
		return 0
	case "nox_xxx_buttonFindFirstEmptySlot_461250":
		return uint32(C.nox_xxx_buttonFindFirstEmptySlot_461250())
	case "sub_4612A0":
		return uint32(C.sub_4612A0())
	case "nox_xxx_buttonHaveSpellInBarMB_4612D0":
		return uint32(C.nox_xxx_buttonHaveSpellInBarMB_4612D0(C.int(a[0])))
	case "nox_xxx_buttonSetImgMB_461320":
		C.nox_xxx_buttonSetImgMB_461320(C.int(a[0]), (*C.uint32_t)(unsafe.Pointer(uintptr(a[1]))))
		return 0
	case "sub_461360":
		return uint32(C.sub_461360(C.int(a[0])))
	case "sub_461400":
		return uint32(C.sub_461400())
	case "sub_461440":
		return uint32(C.sub_461440(C.int(a[0])))
	case "sub_461450":
		return uint32(C.sub_461450())
	case "nox_xxx_spellPutInBox_45DEB0":
		return uint32(C.nox_xxx_spellPutInBox_45DEB0((*C.int)(unsafe.Pointer(uintptr(a[0]))), C.int(a[1]), C.int(a[2]), C.int(a[3])))
	case "nox_client_buildTrap_45E040":
		C.nox_client_buildTrap_45E040()
		return 0
	case "nox_xxx_quickBarCreate_45E190":
		return uint32(C.nox_xxx_quickBarCreate_45E190())
	case "nox_xxx_quickbarButtonBook_45F3F0":
		return uint32(C.nox_xxx_quickbarButtonBook_45F3F0())
	case "sub_45F480":
		return uint32(C.sub_45F480(C.int(a[0])))
	case "sub_45F9B0":
		return uint32(C.sub_45F9B0((*C.uint32_t)(unsafe.Pointer(uintptr(a[0])))))
	case "nox_xxx_quickbarDraw_45FAC0":
		return uint32(C.nox_xxx_quickbarDraw_45FAC0((*C.uint32_t)(unsafe.Pointer(uintptr(a[0])))))
	case "nox_xxx_quickBarDrawFn_45FBD0":
		return uint32(C.nox_xxx_quickBarDrawFn_45FBD0(C.int(a[0])))
	case "nox_xxx_quickBarWarriorDraw_45FDE0":
		return uint32(C.nox_xxx_quickBarWarriorDraw_45FDE0(C.int(a[0])))
	case "sub_460070":
		return uint32(C.sub_460070())
	case "nox_xxx_quickbarAddTrap_460EC0":
		return uint32(uintptr(unsafe.Pointer(C.nox_xxx_quickbarAddTrap_460EC0(C.int(a[0])))))
	}
	panic("unknown quickbar operation: " + op)
}

func PortTestQuickbarCallbacks() map[string]unsafe.Pointer {
	return map[string]unsafe.Pointer{
		"nox_xxx_guiSpellTargetClickSet_45D9D0":       unsafe.Pointer(C.nox_xxx_guiSpellTargetClickSet_45D9D0),
		"nox_xxx_guiSpell_45DA10":                     unsafe.Pointer(C.nox_xxx_guiSpell_45DA10),
		"nox_client_invokeSpellSlot_45DA50":           unsafe.Pointer(C.nox_client_invokeSpellSlot_45DA50),
		"nox_xxx_clientStoreLastButton_45DAD0":        unsafe.Pointer(C.nox_xxx_clientStoreLastButton_45DAD0),
		"nox_xxx_clientSendAbil_45DAF0":               unsafe.Pointer(C.nox_xxx_clientSendAbil_45DAF0),
		"nox_xxx_clientSendSpell_45DB20":              unsafe.Pointer(C.nox_xxx_clientSendSpell_45DB20),
		"sub_45DB90":                                  unsafe.Pointer(C.sub_45DB90),
		"nox_xxx_guiSpellTargetClickCheckSend_45DBB0": unsafe.Pointer(C.nox_xxx_guiSpellTargetClickCheckSend_45DBB0),
		"nox_xxx_book_45DBE0":                         unsafe.Pointer(C.nox_xxx_book_45DBE0),
		"nox_xxx_spellKeyPackSetSpell_45DC40":         unsafe.Pointer(C.nox_xxx_spellKeyPackSetSpell_45DC40),
		"nox_xxx_bookSpellDrop_45DCA0":                unsafe.Pointer(C.nox_xxx_bookSpellDrop_45DCA0),
		"nox_xxx_updateSpellIcons_45DDF0":             unsafe.Pointer(C.nox_xxx_updateSpellIcons_45DDF0),
		"nox_xxx_updateSpellIconDir_45DE10":           unsafe.Pointer(C.nox_xxx_updateSpellIconDir_45DE10),
		"nox_xxx_spellBoxPointToWnd_45DE60":           unsafe.Pointer(C.nox_xxx_spellBoxPointToWnd_45DE60),
		"nox_xxx_guiSpellSetCursor_45DF60":            unsafe.Pointer(C.nox_xxx_guiSpellSetCursor_45DF60),
		"sub_45DFC0":                                  unsafe.Pointer(C.sub_45DFC0),
		"nox_xxx_clientUpdateButtonRow_45E110":        unsafe.Pointer(C.nox_xxx_clientUpdateButtonRow_45E110),
		"nox_xxx_buttonsGetSelectedRow_45E180":        unsafe.Pointer(C.nox_xxx_buttonsGetSelectedRow_45E180),
		"nox_xxx_quickbarButtonBookDraw_45EF30":       unsafe.Pointer(C.nox_xxx_quickbarButtonBookDraw_45EF30),
		"sub_45EF40":                                  unsafe.Pointer(C.sub_45EF40),
		"nox_xxx_quickBarWnd_45EF50":                  unsafe.Pointer(C.nox_xxx_quickBarWnd_45EF50),
		"nox_xxx_clientSwapQuickbarKeys_45F300":       unsafe.Pointer(C.nox_xxx_clientSwapQuickbarKeys_45F300),
		"sub_45F390":                                  unsafe.Pointer(C.sub_45F390),
		"nox_xxx_quickbarButtonBookWnd_45F450":        unsafe.Pointer(C.nox_xxx_quickbarButtonBookWnd_45F450),
		"sub_45F500":                                  unsafe.Pointer(C.sub_45F500),
		"sub_45F520":                                  unsafe.Pointer(C.sub_45F520),
		"sub_45F5D0":                                  unsafe.Pointer(C.sub_45F5D0),
		"nox_xxx_quickbarTrapUpDownProc_45F630":       unsafe.Pointer(C.nox_xxx_quickbarTrapUpDownProc_45F630),
		"nox_xxx_quickbarTrapUpDownDraw_45F6F0":       unsafe.Pointer(C.nox_xxx_quickbarTrapUpDownDraw_45F6F0),
		"nox_xxx_quickbarTrapButtonProc_45F7A0":       unsafe.Pointer(C.nox_xxx_quickbarTrapButtonProc_45F7A0),
		"sub_45F890":                                  unsafe.Pointer(C.sub_45F890),
		"nox_xxx_quickbar_45F8D0":                     unsafe.Pointer(C.nox_xxx_quickbar_45F8D0),
		"sub_45F8F0":                                  unsafe.Pointer(C.sub_45F8F0),
		"sub_45F900":                                  unsafe.Pointer(C.sub_45F900),
		"nox_xxx_quickbarTrapProc_45FB90":             unsafe.Pointer(C.nox_xxx_quickbarTrapProc_45FB90),
		"nox_xxx_quickbarDrawFn_460000":               unsafe.Pointer(C.nox_xxx_quickbarDrawFn_460000),
		"nox_xxx_quickBarInitWindow_4601F0":           unsafe.Pointer(C.nox_xxx_quickBarInitWindow_4601F0),
		"sub_4602F0":                                  unsafe.Pointer(C.sub_4602F0),
		"sub_460380":                                  unsafe.Pointer(C.sub_460380),
		"nox_client_trapSetNext_4603A0":               unsafe.Pointer(C.nox_client_trapSetNext_4603A0),
		"nox_client_trapSetPrev_4603F0":               unsafe.Pointer(C.nox_client_trapSetPrev_4603F0),
		"nox_client_trapSetSelect_4604B0":             unsafe.Pointer(C.nox_client_trapSetSelect_4604B0),
		"sub_4604E0":                                  unsafe.Pointer(C.sub_4604E0),
		"nox_client_spellSetNext_4604F0":              unsafe.Pointer(C.nox_client_spellSetNext_4604F0),
		"nox_client_spellSetPrev_460540":              unsafe.Pointer(C.nox_client_spellSetPrev_460540),
		"nox_client_spellSetSelect_460590":            unsafe.Pointer(C.nox_client_spellSetSelect_460590),
		"sub_460630":                                  unsafe.Pointer(C.sub_460630),
		"nox_xxx_guiSpell_460650":                     unsafe.Pointer(C.nox_xxx_guiSpell_460650),
		"sub_460660":                                  unsafe.Pointer(C.sub_460660),
		"nox_xxx_quickBarClose_4606B0":                unsafe.Pointer(C.nox_xxx_quickBarClose_4606B0),
		"sub_460920":                                  unsafe.Pointer(C.sub_460920),
		"sub_460940":                                  unsafe.Pointer(C.sub_460940),
		"sub_460A10":                                  unsafe.Pointer(C.sub_460A10),
		"sub_460B90":                                  unsafe.Pointer(C.sub_460B90),
		"sub_460D40":                                  unsafe.Pointer(C.sub_460D40),
		"sub_460D50":                                  unsafe.Pointer(C.sub_460D50),
		"nox_xxx_cliPrepareGameplay1_460E60":          unsafe.Pointer(C.nox_xxx_cliPrepareGameplay1_460E60),
		"sub_460EA0":                                  unsafe.Pointer(C.sub_460EA0),
		"sub_460EB0":                                  unsafe.Pointer(C.sub_460EB0),
		"sub_461010":                                  unsafe.Pointer(C.sub_461010),
		"sub_461060":                                  unsafe.Pointer(C.sub_461060),
		"sub_461090":                                  unsafe.Pointer(C.sub_461090),
		"sub_4610D0":                                  unsafe.Pointer(C.sub_4610D0),
		"sub_461120":                                  unsafe.Pointer(C.sub_461120),
		"sub_461160":                                  unsafe.Pointer(C.sub_461160),
		"sub_4611A0":                                  unsafe.Pointer(C.sub_4611A0),
		"sub_4611B0":                                  unsafe.Pointer(C.sub_4611B0),
		"nox_xxx_netAbilityRewardCli_4611E0":          unsafe.Pointer(C.nox_xxx_netAbilityRewardCli_4611E0),
		"nox_xxx_buttonFindFirstEmptySlot_461250":     unsafe.Pointer(C.nox_xxx_buttonFindFirstEmptySlot_461250),
		"sub_4612A0":                                  unsafe.Pointer(C.sub_4612A0),
		"nox_xxx_buttonHaveSpellInBarMB_4612D0":       unsafe.Pointer(C.nox_xxx_buttonHaveSpellInBarMB_4612D0),
		"nox_xxx_buttonSetImgMB_461320":               unsafe.Pointer(C.nox_xxx_buttonSetImgMB_461320),
		"sub_461360":                                  unsafe.Pointer(C.sub_461360),
		"sub_461400":                                  unsafe.Pointer(C.sub_461400),
		"sub_461440":                                  unsafe.Pointer(C.sub_461440),
		"sub_461450":                                  unsafe.Pointer(C.sub_461450),
		"nox_xxx_spellPutInBox_45DEB0":                unsafe.Pointer(C.nox_xxx_spellPutInBox_45DEB0),
		"nox_client_buildTrap_45E040":                 unsafe.Pointer(C.nox_client_buildTrap_45E040),
		"nox_xxx_quickBarCreate_45E190":               unsafe.Pointer(C.nox_xxx_quickBarCreate_45E190),
		"nox_xxx_quickbarButtonBook_45F3F0":           unsafe.Pointer(C.nox_xxx_quickbarButtonBook_45F3F0),
		"sub_45F480":                                  unsafe.Pointer(C.sub_45F480),
		"sub_45F9B0":                                  unsafe.Pointer(C.sub_45F9B0),
		"nox_xxx_quickbarDraw_45FAC0":                 unsafe.Pointer(C.nox_xxx_quickbarDraw_45FAC0),
		"nox_xxx_quickBarDrawFn_45FBD0":               unsafe.Pointer(C.nox_xxx_quickBarDrawFn_45FBD0),
		"nox_xxx_quickBarWarriorDraw_45FDE0":          unsafe.Pointer(C.nox_xxx_quickBarWarriorDraw_45FDE0),
		"sub_460070":                                  unsafe.Pointer(C.sub_460070),
		"nox_xxx_quickbarAddTrap_460EC0":              unsafe.Pointer(C.nox_xxx_quickbarAddTrap_460EC0),
	}
}
