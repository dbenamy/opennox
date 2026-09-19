//go:build porttest

package legacy

/*
#include "defs.h"
#include "GAME2_2.h"
*/
import "C"

import (
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

func PortTestSessionFilterVersion() uint32 { return uint32(C.NOX_CLIENT_VERS_CODE) }
func PortTestSessionFilter(mode uint32, filters [11]uint32, record [169]byte) (int, bool) {
	modePtr := memmap.PtrUint32(0x5D4594, 1193372)
	savedMode := *modePtr
	data := unsafe.Slice(memmap.PtrUint32(0x5D4594, 1193388), 11)
	var saved [11]uint32
	copy(saved[:], data)
	defer func() { *modePtr = savedMode; copy(data, saved[:]) }()
	*modePtr = mode
	copy(data, filters[:])
	p, free := alloc.Malloc(169)
	defer free()
	copy(unsafe.Slice((*byte)(p), 169), record[:])
	v := int(C.nox_xxx_checkSomeFlagsOnJoin_4899C0((*C.nox_gui_server_ent_t)(p)))
	unchanged := true
	for i, b := range unsafe.Slice((*byte)(p), 169) {
		if b != record[i] {
			unchanged = false
		}
	}
	return v, unchanged
}

// PortTestSessionFilterRules supplies the real rule tables/server lookup owner.
func PortTestSessionFilterRules() func() {
	core, _ := server.PortTestRuleServerSetup()
	oldGet := GetServer
	GetServer = func() Server { return &portTestRandomServer{core: core} }
	restoreTable := portTestRuleTable()
	oldOnline, oldContext := Get_dword_5d4594_2650652(), ruleLoaderContext
	Set_dword_5d4594_2650652(0)
	restoreFlags := noxflags.PortTestGameFlags(0)
	return func() {
		restoreFlags()
		Set_dword_5d4594_2650652(oldOnline)
		ruleLoaderContext = oldContext
		restoreTable()
		GetServer = oldGet
	}
}
