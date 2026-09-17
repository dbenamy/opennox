//go:build porttest

package legacy

/*
#include "GAME2.h"
#include "GAME1.h"
#include "client__gui__servopts__guiserv.h"
extern unsigned int dword_5d4594_1046492;
extern uint32_t dword_5d4594_1046496;
extern uint32_t dword_5d4594_1046500;
extern uint32_t dword_5d4594_1046504;
extern uint32_t dword_5d4594_1046508;
extern nox_window* dword_5d4594_1046512;
extern uint32_t dword_5d4594_1046516;
extern uint32_t dword_5d4594_1046520;
extern uint32_t dword_5d4594_1046524;
extern uint32_t dword_5d4594_1046528;
extern uint32_t dword_5d4594_1046532;
extern uint32_t dword_5d4594_1046536;
extern uint32_t dword_5d4594_1046540;
extern uint32_t dword_5d4594_1046356;
extern uint32_t dword_5d4594_1046360;
extern uint32_t dword_587000_129656;
extern uint32_t dword_5d4594_1045480;
extern uint32_t dword_5d4594_1045484;
extern uint32_t dword_5d4594_1045508;
extern uint32_t dword_5d4594_1045516;
extern uint32_t dword_5d4594_1045520;
extern uint32_t dword_5d4594_1045528;
extern uint32_t dword_5d4594_1045532;
extern uint32_t dword_5d4594_1045536;
extern uint32_t dword_5d4594_1045540;
extern uint32_t dword_5d4594_1045544;
extern uint32_t dword_5d4594_1045548;
extern uint32_t dword_5d4594_1045552;
extern uint32_t dword_5d4594_1045556;
extern uint32_t dword_5d4594_1045576;
extern uint32_t dword_5d4594_1045580;
extern uint32_t dword_5d4594_1045584;
extern uint32_t dword_5d4594_1045588;
extern uint32_t dword_5d4594_1045596;
extern uint32_t dword_5d4594_1309812;
extern uint32_t dword_5d4594_1316704;
extern uint32_t dword_5d4594_1316708;
extern uint32_t dword_5d4594_1316712;
extern uint32_t dword_5d4594_1316972;
extern uint32_t dword_5d4594_1523024;
extern uint32_t dword_5d4594_1523028;
extern uint32_t dword_5d4594_1523032;
extern uint32_t dword_5d4594_1523036;
extern uint32_t dword_5d4594_1523040;
extern uint32_t dword_5d4594_1523044;
extern uint32_t dword_5d4594_1523048;
extern int nox_server_gameSettingsUpdated;
extern uint32_t dword_5d4594_371692;
extern int nox_gui_gamemode_loaded_1046548;
extern unsigned char nox_gui_gamemodes[];
void nox_gui_gamemode_load_457410(void);
void nox_xxx_windowServerOptionsFillGametypeList_4596A0(void);
int sub_459650(wchar2_t* title);
void sub_457A10(void);
*/
import "C"

import (
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"unsafe"
)

// Access only the original table/cache owner. No selection algorithm is copied.
func PortTestServerOptionsModes() func() {
	cache := C.nox_gui_gamemode_loaded_1046548
	table := unsafe.Slice((*byte)(unsafe.Pointer(&C.nox_gui_gamemodes)), 8*16)
	saved := append([]byte(nil), table...)
	C.nox_gui_gamemode_loaded_1046548 = 0
	return func() { copy(table, saved); C.nox_gui_gamemode_loaded_1046548 = cache }
}
func PortTestServerOptionsModeName(mode uint16) string {
	return alloc.GoString16((*uint16)(unsafe.Pointer(C.nox_xxx_guiServerOptionsGetGametypeName_4573C0(C.short(mode)))))
}
func PortTestServerOptionsModeFromName(name string) int {
	return int(C.sub_459650((*C.wchar2_t)(unsafe.Pointer(alloc.InternCString16(name)))))
}
func PortTestServerOptionsModeLoaded() bool { return C.nox_gui_gamemode_loaded_1046548 != 0 }

func PortTestServerOptionsWords() (map[string]*uint32, func()) {
	words := map[string]*uint32{
		"panel-1523024": (*uint32)(unsafe.Pointer(&C.dword_5d4594_1523024)),
		"panel-1523028": (*uint32)(unsafe.Pointer(&C.dword_5d4594_1523028)),
		"panel-1523032": (*uint32)(unsafe.Pointer(&C.dword_5d4594_1523032)),
		"panel-1523036": (*uint32)(unsafe.Pointer(&C.dword_5d4594_1523036)),
		"panel-1523040": (*uint32)(unsafe.Pointer(&C.dword_5d4594_1523040)),
		"panel-1523044": (*uint32)(unsafe.Pointer(&C.dword_5d4594_1523044)),
		"panel-1523048": (*uint32)(unsafe.Pointer(&C.dword_5d4594_1523048)),

		"panel-1045480": (*uint32)(unsafe.Pointer(&C.dword_5d4594_1045480)),
		"panel-1045484": (*uint32)(unsafe.Pointer(&C.dword_5d4594_1045484)),
		"panel-1045508": (*uint32)(unsafe.Pointer(&C.dword_5d4594_1045508)),
		"panel-1045516": (*uint32)(unsafe.Pointer(&C.dword_5d4594_1045516)),
		"panel-1045520": (*uint32)(unsafe.Pointer(&C.dword_5d4594_1045520)),
		"panel-1045528": (*uint32)(unsafe.Pointer(&C.dword_5d4594_1045528)),
		"panel-1045532": (*uint32)(unsafe.Pointer(&C.dword_5d4594_1045532)),
		"panel-1045536": (*uint32)(unsafe.Pointer(&C.dword_5d4594_1045536)),
		"panel-1045540": (*uint32)(unsafe.Pointer(&C.dword_5d4594_1045540)),
		"panel-1045544": (*uint32)(unsafe.Pointer(&C.dword_5d4594_1045544)),
		"panel-1045548": (*uint32)(unsafe.Pointer(&C.dword_5d4594_1045548)),
		"panel-1045552": (*uint32)(unsafe.Pointer(&C.dword_5d4594_1045552)),
		"panel-1045556": (*uint32)(unsafe.Pointer(&C.dword_5d4594_1045556)),
		"panel-1045576": (*uint32)(unsafe.Pointer(&C.dword_5d4594_1045576)),
		"panel-1045580": (*uint32)(unsafe.Pointer(&C.dword_5d4594_1045580)),
		"panel-1045584": (*uint32)(unsafe.Pointer(&C.dword_5d4594_1045584)),
		"panel-1045588": (*uint32)(unsafe.Pointer(&C.dword_5d4594_1045588)),
		"panel-1045596": (*uint32)(unsafe.Pointer(&C.dword_5d4594_1045596)),
		"panel-1309812": (*uint32)(unsafe.Pointer(&C.dword_5d4594_1309812)),
		"panel-1316704": (*uint32)(unsafe.Pointer(&C.dword_5d4594_1316704)),
		"panel-1316708": (*uint32)(unsafe.Pointer(&C.dword_5d4594_1316708)),
		"panel-1316712": (*uint32)(unsafe.Pointer(&C.dword_5d4594_1316712)),
		"panel-1316972": (*uint32)(unsafe.Pointer(&C.dword_5d4594_1316972)),

		"settings-updated":      (*uint32)(unsafe.Pointer(&C.nox_server_gameSettingsUpdated)),
		"settings-record-dirty": (*uint32)(unsafe.Pointer(&C.dword_5d4594_371692)),
		"root":                  (*uint32)(unsafe.Pointer(&C.dword_5d4594_1046492)),
		"maps":                  (*uint32)(unsafe.Pointer(&C.dword_5d4594_1046496)),
		"map-controls":          (*uint32)(unsafe.Pointer(&C.dword_5d4594_1046500)),
		"limits-panel":          (*uint32)(unsafe.Pointer(&C.dword_5d4594_1046504)),
		"teams-panel":           (*uint32)(unsafe.Pointer(&C.dword_5d4594_1046508)),
		"name":                  (*uint32)(unsafe.Pointer(&C.dword_5d4594_1046512)),
		"score":                 (*uint32)(unsafe.Pointer(&C.dword_5d4594_1046516)),
		"time":                  (*uint32)(unsafe.Pointer(&C.dword_5d4594_1046520)),
		"main-panel":            (*uint32)(unsafe.Pointer(&C.dword_5d4594_1046524)),
		"access":                (*uint32)(unsafe.Pointer(&C.dword_5d4594_1046528)),
		"players":               (*uint32)(unsafe.Pointer(&C.dword_5d4594_1046532)),
		"advanced":              (*uint32)(unsafe.Pointer(&C.dword_5d4594_1046536)),
		"advanced-open":         (*uint32)(unsafe.Pointer(&C.dword_5d4594_1046540)),
		"tabs2":                 (*uint32)(unsafe.Pointer(&C.dword_5d4594_1046356)),
		"tabs3":                 (*uint32)(unsafe.Pointer(&C.dword_5d4594_1046360)),
		"first-open":            (*uint32)(unsafe.Pointer(&C.dword_587000_129656)),
		"dirty":                 memmap.PtrUint32(0x5D4594, 1046544),
	}
	saved := make(map[string]uint32)
	for name, p := range words {
		saved[name] = *p
		*p = 0
	}
	return words, func() {
		for name, p := range words {
			*p = saved[name]
		}
	}
}
func PortTestServerOptions(op string, ptr unsafe.Pointer, value int, text string) int {
	switch op {
	case "setup":
		return int(C.sub_457B60(C.int(uintptr(ptr))))
	case "refresh":
		return int(C.sub_459C30())
	case "apply":
		return int(C.sub_459150())
	case "tab":
		return int(C.sub_4593B0(C.int(value)))
	case "tab-order":
		return int(C.sub_459560(C.int(value)))
	case "reset-map":
		return int(C.sub_459700())
	case "try-close":
		return int(C.nox_xxx_guiServerOptionsTryHide_4574D0())
	case "tooltip-assign":
		return int(C.nox_xxx_options_457AA0(0, (*C.uint8_t)(ptr)))
	case "tooltip-damage":
		return int(C.nox_xxx_options_457B00(0, (*C.uint8_t)(ptr)))
	case "populate":
		C.nox_xxx_windowServerOptionsFillGametypeList_4596A0()
	case "measure":
		C.sub_457A10()
	case "selected-mode":
		return int(C.sub_459C10())
	case "dirty-set":
		return int(C.sub_459D50(C.int(value)))
	case "dirty-get":
		return int(C.sub_459D60())
	case "visible":
		return int(C.sub_459D80(C.int(value)))
	case "open":
		return int(C.sub_459DA0())
	case "limits":
		return int(C.sub_457460(C.int(uintptr(ptr))))
	case "quest":
		return int(C.sub_457F30(C.int(value)))
	case "name":
		return int(C.sub_459A40((*C.char)(unsafe.Pointer(alloc.InternCString(text)))))
	case "settings-read":
		return int(uintptr(unsafe.Pointer(C.sub_459AA0(ptr))))
	case "settings-labels":
		return int(C.sub_459880(C.int(uintptr(ptr))))
	case "labels":
		return int(C.sub_4580E0(C.int(uintptr(ptr))))
	case "team-count":
		return int(C.sub_459CD0())
	case "construct":
		return int(C.nox_xxx_guiServerOptsLoad_457500())
	case "close":
		return int(uintptr(unsafe.Pointer(C.nox_xxx_guiServerOptionsHide_4597E0(C.int(value)))))
	default:
		panic(op)
	}
	return 0
}
func PortTestServerOptionsEvent(root, child *gui.Window, code, value int) int {
	return int(C.nox_xxx_guiServerOptionsProcPre_4585D0(C.int(uintptr(root.C())), C.uint(code), C.int(uintptr(child.C())), C.int(value)))
}
func PortTestServerOptionsEdit(root *gui.Window, code, id int) int {
	return int(C.nox_xxx_guiServerOptionsProcPre_4585D0(C.int(uintptr(root.C())), 16387, C.int(code), C.int(id)))
}

func PortTestServerOptionsMaps(mode int, current string, update bool) {
	C.nox_client_guiserv_updateMapList_458230(C.int(mode), (*C.char)(unsafe.Pointer(alloc.InternCString(current))), C.bool(update))
}
func PortTestServerOptionsDraw(w *gui.Window) int {
	return int(C.nox_xxx_windowServerOptionsDrawProc_458500((*C.uint32_t)(w.C()), C.int(uintptr(unsafe.Pointer(w.DrawData())))))
}
func PortTestServerOptionsKey(w *gui.Window, event, key, state int) int {
	return int(C.nox_xxx_guiServerOptionsProc_458590(C.int(uintptr(w.C())), C.int(event), C.int(key), C.int(state)))
}

func PortTestServerOptionsRuleState() func() {
	old := ruleLoaderContext
	oldOnline := Get_dword_5d4594_2650652()
	Set_dword_5d4594_2650652(0)
	restoreTable := portTestRuleTable()
	return func() { restoreTable(); ruleLoaderContext = old; Set_dword_5d4594_2650652(oldOnline) }
}
