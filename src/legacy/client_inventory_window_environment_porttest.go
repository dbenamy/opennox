//go:build porttest

package legacy

/*
#include "defs.h"
#include "GAME2_1.h"
#include "client__gui__guiinv.h"




*/
import "C"

import "unsafe"

func PortTestInventoryWindowWords() (map[string]*uint32, func()) {
	words := map[string]*uint32{
		"nox_gui_itemAmount_item_1319256":                   (*uint32)(unsafe.Pointer(&legacyGlobals.nox_gui_itemAmount_item_1319256)),
		"nox_gui_itemAmount_dialog_1319228":                 (*uint32)(unsafe.Pointer(&legacyGlobals.nox_gui_itemAmount_dialog_1319228)),
		"dword_587000_183460":                               (*uint32)(unsafe.Pointer(&dword_587000_183460)),
		"dword_587000_183456":                               (*uint32)(unsafe.Pointer(&dword_587000_183456)),
		"dword_5d4594_1320964":                              (*uint32)(unsafe.Pointer(&dword_5d4594_1320964)),
		"dword_5d4594_1319268":                              (*uint32)(unsafe.Pointer(&dword_5d4594_1319268)),
		"dword_5d4594_1319264":                              (*uint32)(unsafe.Pointer(&dword_5d4594_1319264)),
		"dword_5d4594_1319260":                              (*uint32)(unsafe.Pointer(&dword_5d4594_1319260)),
		"dword_5d4594_1319248":                              (*uint32)(unsafe.Pointer(&dword_5d4594_1319248)),
		"dword_5d4594_1319236":                              (*uint32)(unsafe.Pointer(&dword_5d4594_1319236)),
		"dword_5d4594_1319232":                              (*uint32)(unsafe.Pointer(&dword_5d4594_1319232)),
		"dword_5d4594_1319056":                              (*uint32)(unsafe.Pointer(&interactionKeyState)),
		"nox_wnd_quitMenu_825760":                           (*uint32)(unsafe.Pointer(&sessionQuitRoot)),
		"dword_5d4594_1321228":                              (*uint32)(unsafe.Pointer(&dword_5d4594_1321228)),
		"dword_5d4594_1309820":                              (*uint32)(unsafe.Pointer(&dword_5d4594_1309820)),
		"dword_5d4594_2523804":                              (*uint32)(unsafe.Pointer(&dword_5d4594_2523804)),
		"dword_5d4594_1047520":                              (*uint32)(unsafe.Pointer(&dword_5d4594_1047520)),
		"dword_587000_136184":                               (*uint32)(unsafe.Pointer(&dword_587000_136184)),
		"dword_5d4594_1049796_inventory_click_column_index": (*uint32)(unsafe.Pointer(&dword_5d4594_1049796_inventory_click_column_index)),
		"dword_5d4594_1049800_inventory_click_row_index":    (*uint32)(unsafe.Pointer(&dword_5d4594_1049800_inventory_click_row_index)),
		"dword_5d4594_1049804":                              (*uint32)(unsafe.Pointer(&dword_5d4594_1049804)),
		"dword_5d4594_1049808":                              (*uint32)(unsafe.Pointer(&dword_5d4594_1049808)),
		"dword_5d4594_1049844":                              (*uint32)(unsafe.Pointer(&dword_5d4594_1049844)),
		"dword_5d4594_1049856":                              (*uint32)(unsafe.Pointer(&dword_5d4594_1049856)),
		"dword_5d4594_1049864":                              (*uint32)(unsafe.Pointer(&dword_5d4594_1049864)),
		"dword_5d4594_1049976":                              (*uint32)(unsafe.Pointer(&dword_5d4594_1049976)),
		"dword_5d4594_1049992":                              (*uint32)(unsafe.Pointer(&dword_5d4594_1049992)),
		"dword_5d4594_1049996":                              (*uint32)(unsafe.Pointer(&dword_5d4594_1049996)),
		"dword_5d4594_1050008":                              (*uint32)(unsafe.Pointer(&dword_5d4594_1050008)),
		"dword_5d4594_1062452":                              (*uint32)(unsafe.Pointer(&legacyGlobals.dword_5d4594_1062452)),
		"dword_5d4594_1062456":                              (*uint32)(unsafe.Pointer(&dword_5d4594_1062456)),
		"dword_5d4594_1062468":                              (*uint32)(unsafe.Pointer(&dword_5d4594_1062468)),
		"dword_5d4594_1062476":                              (*uint32)(unsafe.Pointer(&dword_5d4594_1062476)),
		"dword_5d4594_1062480":                              (*uint32)(unsafe.Pointer(&dword_5d4594_1062480)),
		"dword_5d4594_1062488":                              (*uint32)(unsafe.Pointer(&dword_5d4594_1062488)),
		"dword_5d4594_1062492":                              (*uint32)(unsafe.Pointer(&dword_5d4594_1062492)),
		"dword_5d4594_1062496":                              (*uint32)(unsafe.Pointer(&dword_5d4594_1062496)),
		"dword_5d4594_1062508":                              (*uint32)(unsafe.Pointer(&dword_5d4594_1062508)),
		"dword_5d4594_1062512":                              (*uint32)(unsafe.Pointer(&dword_5d4594_1062512)),
		"dword_5d4594_1062516":                              (*uint32)(unsafe.Pointer(&dword_5d4594_1062516)),
		"dword_5d4594_1062520":                              (*uint32)(unsafe.Pointer(&dword_5d4594_1062520)),
		"dword_5d4594_1062524":                              (*uint32)(unsafe.Pointer(&dword_5d4594_1062524)),
		"dword_5d4594_1062528":                              (*uint32)(unsafe.Pointer(&dword_5d4594_1062528)),
		"dword_5d4594_1062552":                              (*uint32)(unsafe.Pointer(&dword_5d4594_1062552)),
		"dword_5d4594_1062556":                              (*uint32)(unsafe.Pointer(&dword_5d4594_1062556)),
		"dword_5d4594_1062560":                              (*uint32)(unsafe.Pointer(&dword_5d4594_1062560)),
		"dword_5d4594_1062564":                              (*uint32)(unsafe.Pointer(&dword_5d4594_1062564)),
		"dword_5d4594_1063116":                              (*uint32)(unsafe.Pointer(&dword_5d4594_1063116)),
		"dword_5d4594_1063120":                              (*uint32)(unsafe.Pointer(&dword_5d4594_1063120)),
		"dword_5d4594_1063636":                              (*uint32)(unsafe.Pointer(&dword_5d4594_1063636)),
		"dword_5d4594_1098624":                              (*uint32)(unsafe.Pointer(&dword_5d4594_1098624)),
		"dword_5d4594_1098628":                              (*uint32)(unsafe.Pointer(&dword_5d4594_1098628)),
		"dword_5d4594_1107036":                              (*uint32)(unsafe.Pointer(&dword_5d4594_1107036)),
		"nox_win_unk5":                                      (*uint32)(unsafe.Pointer(&legacyGlobals.nox_win_unk5)),
	}
	old := make(map[string]uint32, len(words))
	for n, p := range words {
		old[n] = *p
		*p = 0
	}
	return words, func() {
		for n, p := range words {
			*p = old[n]
		}
	}
}

func PortTestInventoryWindowCallbacks() []unsafe.Pointer {
	return []unsafe.Pointer{
		clientUICallbackKey(clientUICallbackID_sub_462740),
		nil,
		inventoryCallbackKey(inventoryCallbackButton),
		inventoryCallbackKey(inventoryCallbackAlt),
		inventoryCallbackKey(inventoryCallbackStatus),
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		uiAmountNativeKey(uiAmountDrop),
		nil,
		nil,
		nil,
		nil,
		nil,
		inventoryCallbackKey(inventoryCallbackHover),
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		clientUICallbackKey(clientUICallbackID_sub_467650),
		nil,
		clientUICallbackKey(clientUICallbackID_sub_467BB0),
		clientUICallbackKey(clientUICallbackID_sub_467C10),
		nil,
		clientUICallbackKey(clientUICallbackID_sub_467C80),
		nil,
		nil,
		inventoryCallbackKey(inventoryCallbackDrawAlt), PortTestMeterInventoryWeaponDrawCallback(), inventoryCallbackKey(inventoryCallbackDrawCur), inventoryCallbackKey(inventoryCallbackMode), inventoryCallbackKey(inventoryCallbackIdentify),
	}
}
