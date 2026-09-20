//go:build porttest

package legacy

/*
#include "defs.h"
#include "GAME2.h"
#include "GAME2_1.h"
#include "client__gui__guispell.h"
*/
import "C"
import (
	"github.com/opennox/opennox/v1/client/gui"
	"image"
	"unsafe"
)

// PortTestQuickbarInvoke calls the production quickbar owner without duplicating
// any of its algorithms. Pointer arguments use the target's 32-bit ABI.
func PortTestQuickbarInvoke(op string, a [7]uint32) uint32 {
	bar := func(p uint32) *quickbarRecord { return (*quickbarRecord)(unsafe.Pointer(uintptr(p))) }
	switch op {
	case "nox_xxx_guiSpellTargetClickSet_45D9D0":
		return uint32(quickbarQueueTarget(a[0]))
	case "nox_xxx_guiSpell_45DA10":
		return uint32(quickbarQueueInstant(a[0]))
	case "nox_client_invokeSpellSlot_45DA50":
		quickbarInvokeSlot(int(a[0]))
		return 0
	case "nox_xxx_clientStoreLastButton_45DAD0":
		quickbarLastButton(int(a[0]))
		return 0
	case "nox_xxx_clientSendAbil_45DAF0":
		return uint32(quickbarSendAbility(a[0]))
	case "nox_xxx_clientSendSpell_45DB20":
		return uint32(quickbarSendSpells((*uint32)(unsafe.Pointer(uintptr(a[0]))), int(a[1]), byte(a[2])))
	case "sub_45DB90":
		return uint32(quickbarResetFlash())
	case "nox_xxx_guiSpellTargetClickCheckSend_45DBB0":
		quickbarSendPending()
		return 0
	case "nox_xxx_book_45DBE0":
		return uint32(quickbarBookSlot(uintptr(a[0]), a[1], int(a[2])))
	case "nox_xxx_spellKeyPackSetSpell_45DC40":
		quickbarSetSlot(bar(a[0]), a[1], int(a[2]))
		return 0
	case "nox_xxx_bookSpellDrop_45DCA0":
		return uint32(quickbarDrop(a[0], byte(a[1]), image.Pt(int(a[2]), int(a[3])), bar(a[4])))
	case "nox_xxx_updateSpellIcons_45DDF0":
		return uint32(quickbarDirections(bar(a[0])))
	case "nox_xxx_updateSpellIconDir_45DE10":
		return uint32(quickbarDirection(bar(a[1]), int(a[0])))
	case "nox_xxx_spellBoxPointToWnd_45DE60":
		return uint32(quickbarPoint(bar(a[0]), image.Pt(int(a[1]), int(a[2]))))
	case "nox_xxx_guiSpellSetCursor_45DF60":
		return uint32(quickbarSpellCursor(a[0], byte(a[1])))
	case "sub_45DFC0":
		return uint32(quickbarAbilityCursor(a[0]))
	case "nox_xxx_clientUpdateButtonRow_45E110":
		return uint32(quickbarSelectRow(int(a[0])))
	case "nox_xxx_buttonsGetSelectedRow_45E180":
		return uint32(quickbarMain().Selected)
	case "nox_xxx_quickbarButtonBookDraw_45EF30":
		return uint32(1)
	case "sub_45EF40":
		return uint32(0)
	case "nox_xxx_quickBarWnd_45EF50":
		return uint32(quickbarSlotEvent(bookWindow(a[0]), a[1], a[2]))
	case "nox_xxx_clientSwapQuickbarKeys_45F300":
		quickbarSwap(bar(a[0]), int(a[1]), int(a[2]))
		return 0
	case "sub_45F390":
		quickbarSwapRows((*[5]quickbarSlot)(unsafe.Pointer(uintptr(a[0]))), (*[5]quickbarSlot)(unsafe.Pointer(uintptr(a[1]))), int(a[2]), int(a[3]))
		return 0
	case "nox_xxx_quickbarButtonBookWnd_45F450":
		return uint32(quickbarBookEvent(bookWindow(a[0]), a[1], 0))
	case "sub_45F500":
		return uint32(quickbarDirectionLit(bar(a[1]), int(a[0])))
	case "sub_45F520":
		return uint32(quickbarDirectionEvent(bookWindow(a[0]), a[1], 0))
	case "sub_45F5D0":
		return uint32(quickbarKeyLabel(bookWindow(a[0]), nil))
	case "nox_xxx_quickbarTrapUpDownProc_45F630":
		return uint32(quickbarRowEvent(bookWindow(a[0]), a[1], 0))
	case "nox_xxx_quickbarTrapUpDownDraw_45F6F0":
		return uint32(quickbarDrawControl(bookWindow(a[0]), bookWindow(a[0]).DrawData()))
	case "nox_xxx_quickbarTrapButtonProc_45F7A0":
		return uint32(quickbarTrapButtonEvent(bookWindow(a[0]), a[1], a[2]))
	case "sub_45F890":
		return uint32(quickbarTrapTargets())
	case "nox_xxx_quickbar_45F8D0":
		return uint32(quickbarGenericEvent(bookWindow(a[0]), a[1], a[2]))
	case "sub_45F8F0":
		quickbarModifier()
		return 1
	case "sub_45F900":
		quickbarModifier()
		return 0
	case "nox_xxx_quickbarTrapProc_45FB90":
		return uint32(quickbarTrapEvent(bookWindow(a[0]), a[1], a[2]))
	case "nox_xxx_quickbarDrawFn_460000":
		return uint32(quickbarDrawSlide(nil, nil))
	case "nox_xxx_quickBarInitWindow_4601F0":
		return uint32(quickbarInitWindow(bar(a[0]), int(a[1]), int(a[2]), int(a[3]), int(a[4]), gui.WrapFuncC(unsafe.Pointer(uintptr(a[5]))), gui.WrapDrawFuncC(unsafe.Pointer(uintptr(a[6])))))
	case "sub_4602F0":
		return uint32(quickbarClearSlots())
	case "sub_460380":
		return uint32(quickbarClearAbilities())
	case "nox_client_trapSetNext_4603A0":
		quickbarTrapNext()
		return 0
	case "nox_client_trapSetPrev_4603F0":
		quickbarTrapPrevious()
		return 0
	case "nox_client_trapSetSelect_4604B0":
		return uint32(quickbarTrapSelect(int(a[0])))
	case "sub_4604E0":
		return uint32(*quickbarByte(1048140))
	case "nox_client_spellSetNext_4604F0":
		return uint32(quickbarMoveRow(1))
	case "nox_client_spellSetPrev_460540":
		return uint32(quickbarMoveRow(-1))
	case "nox_client_spellSetSelect_460590":
		quickbarTimedRow()
		return 0
	case "sub_460630":
		quickbarRememberMouseSequence()
		return 0
	case "nox_xxx_guiSpell_460650":
		return uint32(*quickbarWord(1047928))
	case "sub_460660":
		return uint32(quickbarCancelCapture())
	case "nox_xxx_quickBarClose_4606B0":
		return uint32(quickbarCloseExpanded())
	case "sub_460920":
		quickbarExpand()
		return 0
	case "sub_460940":
		return uint32(quickbarSave())
	case "sub_460A10":
		return uint32(quickbarSaveRows(bar(a[0]), int(a[1]), int(a[2]), byte(a[3])))
	case "sub_460B90":
		return uint32(quickbarVisible(a[0] != 0))
	case "sub_460D40":
		return uint32(bool2int(*quickbarWord(1049508) != 0))
	case "sub_460D50":
		return uint32(quickbarDestroy())
	case "nox_xxx_cliPrepareGameplay1_460E60":
		return uint32(quickbarPrepare())
	case "sub_460EA0":
		return uint32(quickbarVisible(a[0] != 0))
	case "sub_460EB0":
		quickbarSetFlash(int(a[0]), byte(a[1]))
		return 0
	case "sub_461010":
		quickbarCloseTrap()
		return 0
	case "sub_461060":
		quickbarToggleTrap()
		return 0
	case "sub_461090":
		return uint32(quickbarAbilityState(a[0], a[1]))
	case "sub_4610D0":
		return uint32(quickbarResetAbility(byte(a[0])))
	case "sub_461120":
		return uint32(quickbarAbilityFlags(a[0], a[1]))
	case "sub_461160":
		return uint32(quickbarAbilityAvailable(a[0]))
	case "sub_4611A0":
		return uint32(*quickbarWord(1047932))
	case "sub_4611B0":
		return uint32(quickbarSendPendingAbility())
	case "nox_xxx_netAbilityRewardCli_4611E0":
		quickbarAbilityReward(int(a[0]), int(a[1]), uintptr(a[2]))
		return 0
	case "nox_xxx_buttonFindFirstEmptySlot_461250":
		return uint32(quickbarFirstEmpty(false))
	case "sub_4612A0":
		return uint32(quickbarFirstEmpty(true))
	case "nox_xxx_buttonHaveSpellInBarMB_4612D0":
		return uint32(quickbarContains(a[0]))
	case "nox_xxx_buttonSetImgMB_461320":
		quickbarSlotPosition(int(a[0]), (*[2]int32)(unsafe.Pointer(uintptr(a[1]))))
		return 0
	case "sub_461360":
		return uint32(quickbarRemove(a[0]))
	case "sub_461400":
		return uint32(quickbarRestoreSlots())
	case "sub_461440":
		*quickbarWord(1049688) = a[0]
		return a[0]
	case "sub_461450":
		return uint32(*quickbarWord(1049688))
	case "nox_xxx_spellPutInBox_45DEB0":
		return uint32(quickbarPut(bar(a[0]), a[1], image.Pt(int(a[2]), int(a[3]))))
	case "nox_client_buildTrap_45E040":
		quickbarBuildTrap()
		return 0
	case "nox_xxx_quickBarCreate_45E190":
		return uint32(quickbarCreate())
	case "nox_xxx_quickbarButtonBook_45F3F0":
		return uint32(quickbarBookTooltip())
	case "sub_45F480":
		return uint32(quickbarDirectionTooltip(bookWindow(a[0])))
	case "sub_45F9B0":
		return uint32(quickbarRowLabel(bookWindow(a[0]), false))
	case "nox_xxx_quickbarDraw_45FAC0":
		return uint32(quickbarRowLabel(bookWindow(a[0]), true))
	case "nox_xxx_quickBarDrawFn_45FBD0":
		return uint32(quickbarDrawAbility(bookWindow(a[0]), nil))
	case "nox_xxx_quickBarWarriorDraw_45FDE0":
		return uint32(quickbarDrawSpell(bookWindow(a[0]), nil))
	case "sub_460070":
		return uint32(quickbarRevealTrap())
	case "nox_xxx_quickbarAddTrap_460EC0":
		return uint32(quickbarAddTrap(int(a[0])))
	}
	panic("unknown quickbar operation: " + op)
}

func PortTestQuickbarCallbacks() map[string]unsafe.Pointer {
	return map[string]unsafe.Pointer{
		"nox_xxx_guiSpellTargetClickSet_45D9D0":       nil,
		"nox_xxx_guiSpell_45DA10":                     nil,
		"nox_client_invokeSpellSlot_45DA50":           nil,
		"nox_xxx_clientStoreLastButton_45DAD0":        nil,
		"nox_xxx_clientSendAbil_45DAF0":               nil,
		"nox_xxx_clientSendSpell_45DB20":              nil,
		"sub_45DB90":                                  nil,
		"nox_xxx_guiSpellTargetClickCheckSend_45DBB0": nil,
		"nox_xxx_book_45DBE0":                         unsafe.Pointer(C.nox_xxx_book_45DBE0),
		"nox_xxx_spellKeyPackSetSpell_45DC40":         nil,
		"nox_xxx_bookSpellDrop_45DCA0":                nil,
		"nox_xxx_updateSpellIcons_45DDF0":             nil,
		"nox_xxx_updateSpellIconDir_45DE10":           nil,
		"nox_xxx_spellBoxPointToWnd_45DE60":           nil,
		"nox_xxx_guiSpellSetCursor_45DF60":            nil,
		"sub_45DFC0":                                  nil,
		"nox_xxx_clientUpdateButtonRow_45E110":        unsafe.Pointer(C.nox_xxx_clientUpdateButtonRow_45E110),
		"nox_xxx_buttonsGetSelectedRow_45E180":        unsafe.Pointer(C.nox_xxx_buttonsGetSelectedRow_45E180),
		"nox_xxx_quickbarButtonBookDraw_45EF30":       nil,
		"sub_45EF40":                                  nil,
		"nox_xxx_quickBarWnd_45EF50":                  nil,
		"nox_xxx_clientSwapQuickbarKeys_45F300":       nil,
		"sub_45F390":                                  nil,
		"nox_xxx_quickbarButtonBookWnd_45F450":        nil,
		"sub_45F500":                                  nil,
		"sub_45F520":                                  nil,
		"sub_45F5D0":                                  nil,
		"nox_xxx_quickbarTrapUpDownProc_45F630":       nil,
		"nox_xxx_quickbarTrapUpDownDraw_45F6F0":       nil,
		"nox_xxx_quickbarTrapButtonProc_45F7A0":       nil,
		"sub_45F890":                                  nil,
		"nox_xxx_quickbar_45F8D0":                     nil,
		"sub_45F8F0":                                  nil,
		"sub_45F900":                                  nil,
		"nox_xxx_quickbarTrapProc_45FB90":             nil,
		"nox_xxx_quickbarDrawFn_460000":               nil,
		"nox_xxx_quickBarInitWindow_4601F0":           nil,
		"sub_4602F0":                                  unsafe.Pointer(C.sub_4602F0),
		"sub_460380":                                  nil,
		"nox_client_trapSetNext_4603A0":               nil,
		"nox_client_trapSetPrev_4603F0":               nil,
		"nox_client_trapSetSelect_4604B0":             unsafe.Pointer(C.nox_client_trapSetSelect_4604B0),
		"sub_4604E0":                                  unsafe.Pointer(C.sub_4604E0),
		"nox_client_spellSetNext_4604F0":              nil,
		"nox_client_spellSetPrev_460540":              nil,
		"nox_client_spellSetSelect_460590":            nil,
		"sub_460630":                                  nil,
		"nox_xxx_guiSpell_460650":                     nil,
		"sub_460660":                                  unsafe.Pointer(C.sub_460660),
		"nox_xxx_quickBarClose_4606B0":                unsafe.Pointer(C.nox_xxx_quickBarClose_4606B0),
		"sub_460920":                                  nil,
		"sub_460940":                                  unsafe.Pointer(C.sub_460940),
		"sub_460A10":                                  nil,
		"sub_460B90":                                  nil,
		"sub_460D40":                                  nil,
		"sub_460D50":                                  nil,
		"nox_xxx_cliPrepareGameplay1_460E60":          nil,
		"sub_460EA0":                                  unsafe.Pointer(C.sub_460EA0),
		"sub_460EB0":                                  nil,
		"sub_461010":                                  nil,
		"sub_461060":                                  nil,
		"sub_461090":                                  nil,
		"sub_4610D0":                                  nil,
		"sub_461120":                                  nil,
		"sub_461160":                                  nil,
		"sub_4611A0":                                  nil,
		"sub_4611B0":                                  nil,
		"nox_xxx_netAbilityRewardCli_4611E0":          nil,
		"nox_xxx_buttonFindFirstEmptySlot_461250":     nil,
		"sub_4612A0":                                  nil,
		"nox_xxx_buttonHaveSpellInBarMB_4612D0":       nil,
		"nox_xxx_buttonSetImgMB_461320":               nil,
		"sub_461360":                                  nil,
		"sub_461400":                                  nil,
		"sub_461440":                                  nil,
		"sub_461450":                                  nil,
		"nox_xxx_spellPutInBox_45DEB0":                nil,
		"nox_client_buildTrap_45E040":                 nil,
		"nox_xxx_quickBarCreate_45E190":               nil,
		"nox_xxx_quickbarButtonBook_45F3F0":           unsafe.Pointer(C.nox_xxx_quickbarButtonBook_45F3F0),
		"sub_45F480":                                  unsafe.Pointer(C.sub_45F480),
		"sub_45F9B0":                                  nil,
		"nox_xxx_quickbarDraw_45FAC0":                 nil,
		"nox_xxx_quickBarDrawFn_45FBD0":               nil,
		"nox_xxx_quickBarWarriorDraw_45FDE0":          nil,
		"sub_460070":                                  nil,
		"nox_xxx_quickbarAddTrap_460EC0":              nil,
	}
}
