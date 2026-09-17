//go:build porttest

package legacy

/*
#include "GAME2.h"
#include "GAME1.h"
#include "client__gui__servopts__guiserv.h"
extern int nox_server_gameSettingsUpdated;
*/
import "C"

import (
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/common/memmap"
	"unsafe"
)

// Own the private native table/cache; all selection operations call production Go.
func PortTestServerOptionsModes() func() {
	cache := serverOptionsModesLoaded
	table := append([]serverOptionsMode(nil), serverOptionsModes...)
	serverOptionsModesLoaded = false
	return func() { copy(serverOptionsModes, table); serverOptionsModesLoaded = cache }
}
func PortTestServerOptionsModeName(mode uint16) string  { return serverOptionsModeName(mode) }
func PortTestServerOptionsModeFromName(name string) int { return serverOptionsModeFromName(name) }
func PortTestServerOptionsModeLoaded() bool             { return serverOptionsModesLoaded }

func PortTestServerOptionsWords() (map[string]*uint32, func()) {
	words := map[string]*uint32{
		"panel-1523024": &serverConfigRuleWords[0],
		"panel-1523028": &serverConfigRuleWords[1],
		"panel-1523032": &serverConfigRuleWords[2],
		"panel-1523036": &serverConfigRuleWords[3],
		"panel-1523040": &serverConfigRuleWords[4],
		"panel-1523044": &serverConfigRuleWords[5],
		"panel-1523048": &serverConfigRuleWords[6],

		"panel-1045480": serverPanelsWord(1045480),
		"panel-1045484": serverPanelsWord(1045484),
		"panel-1045508": serverPanelsWord(1045508),
		"panel-1045516": serverPanelsWord(1045516),
		"panel-1045520": serverPanelsWord(1045520),
		"panel-1045528": serverPanelsWord(1045528),
		"panel-1045532": serverPanelsWord(1045532),
		"panel-1045536": serverPanelsWord(1045536),
		"panel-1045540": serverPanelsWord(1045540),
		"panel-1045544": serverPanelsWord(1045544),
		"panel-1045548": serverPanelsWord(1045548),
		"panel-1045552": serverPanelsWord(1045552),
		"panel-1045556": serverPanelsWord(1045556),
		"panel-1045576": serverPanelsWord(1045576),
		"panel-1045580": serverPanelsWord(1045580),
		"panel-1045584": serverPanelsWord(1045584),
		"panel-1045588": serverPanelsWord(1045588),
		"panel-1045596": serverPanelsWord(1045596),
		"panel-1309812": serverPanelsWord(1309812),
		"panel-1316704": serverPanelsWord(1316704),
		"panel-1316708": serverPanelsWord(1316708),
		"panel-1316712": serverPanelsWord(1316712),
		"panel-1316972": serverPanelsWord(1316972),

		"settings-updated":      (*uint32)(unsafe.Pointer(&C.nox_server_gameSettingsUpdated)),
		"settings-record-dirty": (*uint32)(unsafe.Pointer(&serverConfigRecordDirty)),
		"root":                  serverOptionsWord(1046492),
		"maps":                  serverOptionsWord(1046496),
		"map-controls":          serverOptionsWord(1046500),
		"limits-panel":          serverOptionsWord(1046504),
		"teams-panel":           serverOptionsWord(1046508),
		"name":                  serverOptionsWord(1046512),
		"score":                 serverOptionsWord(1046516),
		"time":                  serverOptionsWord(1046520),
		"main-panel":            serverOptionsWord(1046524),
		"access":                serverOptionsWord(1046528),
		"players":               serverOptionsWord(1046532),
		"advanced":              serverOptionsWord(1046536),
		"advanced-open":         serverOptionsWord(1046540),
		"tabs2":                 serverOptionsWord(1046356),
		"tabs3":                 serverOptionsWord(1046360),
		"first-open":            &serverOptionsFirstOpen,
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
		return serverOptionsSetup(serverOptionsRecord(ptr))
	case "refresh":
		return serverOptionsRefresh()
	case "apply":
		return int(serverOptionsApply())
	case "tab":
		return serverOptionsTab(value)
	case "tab-order":
		return serverOptionsTabOrder(value)
	case "reset-map":
		return serverOptionsResetMap()
	case "try-close":
		return serverOptionsTryClose()
	case "tooltip-assign":
		return serverOptionsTooltip(false, *(*byte)(ptr))
	case "tooltip-damage":
		return serverOptionsTooltip(true, *(*byte)(ptr))
	case "populate":
		serverOptionsPopulate()
	case "measure":
		serverOptionsMeasure()
	case "selected-mode":
		return serverOptionsSelectedMode()
	case "dirty-set":
		return serverOptionsDirty(value)
	case "dirty-get":
		return int(*serverOptionsWord(1046544))
	case "visible":
		return serverOptionsVisible(value)
	case "open":
		return bool2int(serverOptionsRoot != 0)
	case "limits":
		return serverOptionsLimits(serverOptionsRecord(ptr))
	case "quest":
		return serverOptionsQuest(value)
	case "name":
		return serverOptionsName(text)
	case "settings-read":
		return int(serverOptionsRead(serverOptionsRecord(ptr)))
	case "settings-labels":
		return serverOptionsSettingsLabels(serverOptionsRecord(ptr))
	case "labels":
		return serverOptionsLabels(serverOptionsRecord(ptr))
	case "team-count":
		return serverOptionsTeamCount()
	case "construct":
		return serverOptionsConstruct()
	case "close":
		return int(serverOptionsClose(value))
	default:
		panic(op)
	}
	return 0
}
func PortTestServerOptionsEvent(root, child *gui.Window, code, value int) int {
	return serverOptionsEvent(root, code, uintptr(child.C()), value)
}
func PortTestServerOptionsEdit(root *gui.Window, code, id int) int {
	return serverOptionsEvent(root, 16387, uintptr(code), id)
}
func PortTestServerOptionsMaps(mode int, current string, update bool) {
	serverOptionsMapList(mode, current, update)
}
func PortTestServerOptionsDraw(w *gui.Window) int { return serverOptionsDraw(w, w.DrawData()) }
func PortTestServerOptionsKey(w *gui.Window, event, key, state int) int {
	return serverOptionsKey(w, event, key, state)
}

func PortTestServerOptionsRuleState() func() {
	old := ruleLoaderContext
	oldOnline := Get_dword_5d4594_2650652()
	Set_dword_5d4594_2650652(0)
	restoreTable := portTestRuleTable()
	return func() { restoreTable(); ruleLoaderContext = old; Set_dword_5d4594_2650652(oldOnline) }
}
