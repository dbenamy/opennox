//go:build porttest

package legacy

import (
	"github.com/opennox/opennox/v1/client/gui"
	"image"
	"unsafe"
)

func nox_client_invAlterWeapon_4672C0() { uiInventoryAlterWeapon() }

func nox_client_makePlayerStatsDlg_463880(pos unsafe.Pointer) {
	v := unsafe.Slice((*int32)(pos), 2)
	uiInventoryStats(image.Pt(int(v[0]), int(v[1])))
}

func nox_client_toggleSpellbook_45AC70() { bookToggle() }

func nox_xxx_abilityReward_45D290(id int32, notify *int8, autoAdd int32) {
	bookAbilityReward(int(id), uintptr(unsafe.Pointer(notify)), int(autoAdd))
}

func nox_xxx_bookFillAll_45D570(kind, id int32) { bookAdd(int(kind), int(id)) }

func nox_xxx_bookHideMB_45ACA0(reset int32) int32 { return int32(bookHide(int(reset))) }

func nox_xxx_bookSetForward_45D200(kind unsafe.Pointer, id int32, pos unsafe.Pointer) unsafe.Pointer {
	return unsafe.Pointer(bookSetForward(uintptr(kind), int(id), AsPoint(pos)))
}

func nox_xxx_cliInventorySpriteUpd_465A30() { uiInventoryDragCopy() }

func nox_xxx_clientDequip_464B70(v int32) int32 {
	return int32(uiInventoryDequipRequest(uiInventoryDrawable(uint32(v))))
}

func nox_xxx_clientEquip_4623B0(v int32) int32 {
	return int32(uiInventoryEquipRequest(uiInventoryDrawable(uint32(v))))
}

func nox_xxx_clientSetAltWeapon_461550(v int32) int32 {
	return int32(uiInventorySetAlternate(uiInventoryCellRef(uint32(v))))
}

func nox_xxx_clientTrade_465870(v int16) int32 { return int32(uiInventoryTrade(28, uint16(v))) }

func nox_xxx_clientUse_465C70(v int32) { uiInventoryUse(uiInventoryDrawable(uint32(v))) }

func nox_xxx_guiDrawInventoryTray_4643B0(ax, ay int32) int32 {
	return int32(uiInventoryDrawTray(int32(ax), int32(ay)))
}

func nox_xxx_trade_4657B0(v int16) int32 { return int32(uiInventoryTrade(30, uint16(v))) }

func nox_xxx_wndGetHandle_4676A0() *gui.Window {
	return (*gui.Window)(unsafe.Pointer(legacyGlobals.dword_5d4594_1062452))
}

func sub_45CFC0() int32 { return int32(bookShown()) }

func sub_45D500(show int32) int32 { return int32(bookTemporaryShow(int(show))) }

func sub_45D870() { bookFinishAddition() }

func sub_45D9B0() int32 { return int32(*bookWord(1047520)) }

func sub_4615C0() int { return int(uintptr(unsafe.Pointer(uiInventoryCurrentWeapon()))) }

func sub_461B50() *uint8 { return (*uint8)(unsafe.Pointer(uiInventoryCompact())) }

func sub_461EF0(code int) *int8 {
	return (*int8)(unsafe.Pointer(uiInventoryFindCode(uint32(code))))
}

func sub_462740() int32 { return int32(uiInventoryCloseIdentify()) }

func sub_4627F0(p unsafe.Pointer) int32 {
	v := unsafe.Slice((*uint32)(p), 2)
	return int32(uiInventoryIdentify(image.Pt(int(v[0]), int(v[1]))))
}

func sub_463370(w, pos, out unsafe.Pointer) int32 {
	p := (*[2]int32)(pos)
	pt := image.Pt(int(p[0]), int(p[1])).Sub(uiWindowPosition((*gui.Window)(w)))
	v := unsafe.Slice((*int32)(out), 2)
	v[0], v[1] = int32(pt.X), int32(pt.Y)
	return int32(pt.Y)
}

func sub_4649B0(v, col, row int32) int32 {
	return int32(uiInventoryPlace(uiInventoryDrawable(uint32(v)), int(col), int(row)))
}

func sub_464B40(col, row int32) int32 {
	return int32(bool2int(uiInventoryValidCell(int(col), int(row))))
}

func sub_467650() int32 { return int32(uiInventoryRepairMode()) }

func sub_467810(col, row int) int {
	if col < 0 || col >= 4 || row < 0 || row >= 20 {
		return 0
	}
	return int(uiInventoryGrid()[row+21*col].Count)
}

func sub_467870(col, row int) *int8 {
	if col < 0 || col >= 4 || row < 0 || row >= 20 {
		return nil
	}
	return (*int8)(unsafe.Pointer(&uiInventoryGrid()[row+21*col].Codes[0]))
}

func sub_467BB0() int32 { return int32(uiInventoryOpenWindow()) }

func sub_467C10() int32 { return int32(uiInventoryCloseWindow()) }

func sub_467C80() int32 { return int32(bool2int(uiInventoryWindowOpenState())) }

func sub_472310() *uint8 { return (*uint8)(unsafe.Pointer(uiMeterRefreshPotions())) }
