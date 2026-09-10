//go:build porttest

package legacy

/*
#include "GAME5_2.h"
*/
import "C"

import (
	"github.com/opennox/opennox/v1/server"
)

type PortTestSpellClassCall struct {
	PlayerClass uint32
	Spell       uint32
}

type PortTestSpellClassResult struct {
	PlayerClass uint32
	Spell       uint32
	Result      int
}

// PortTestSpellClass invokes the original C ABI against a server containing
// only the requested spell definitions. Missing IDs are intentionally absent.
func PortTestSpellClass(defs []server.PortTestSpellClassDef, calls []PortTestSpellClassCall, allowAll bool) []PortTestSpellClassResult {
	core := server.PortTestSpellClassServer(defs)
	core.Spells.AllowAll = allowAll
	oldGet := GetServer
	GetServer = func() Server { return &portTestRandomServer{core: core} }
	defer func() { GetServer = oldGet }()
	out := make([]PortTestSpellClassResult, 0, len(calls))
	for _, call := range calls {
		out = append(out, PortTestSpellClassResult{
			PlayerClass: call.PlayerClass,
			Spell:       call.Spell,
			Result:      int(C.nox_xxx_playerCheckSpellClass_57AEA0(C.int(int32(call.PlayerClass)), C.int(int32(call.Spell)))),
		})
	}
	return out
}
