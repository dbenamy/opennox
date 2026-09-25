//go:build porttest

package legacy

import (
	"unsafe"

	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/ccall"
	"github.com/opennox/opennox/v1/server"
)

// PortTestInstallCallbackTables exposes the existing isolated table fixture so
// callers can inspect a selected identity without initializing the game blob.
func PortTestInstallCallbackTables() func() {
	return portTestCallbackTablesEnvironment()
}

// PortTestCallbackTableFunctions returns the callback words from the three
// installed parser tables in table order: strike, die, then dead.
func PortTestCallbackTableFunctions() []unsafe.Pointer {
	var out []unsafe.Pointer
	for _, base := range []uintptr{287096, 287280, 287192} {
		for i := uintptr(0); *memmap.PtrPtr(0x587000, base+8*i) != nil; i++ {
			out = append(out, *memmap.PtrPtr(0x587000, base+8*i+4))
		}
	}
	return out
}

// PortTestMonsterCallbackCallResult preserves the callback's original int ABI.
func PortTestMonsterCallbackCallResult(key unsafe.Pointer, u *server.Object) int32 {
	return int32(ccall.CallIntPtr(key, u.CObj()))
}

// PortTestMonsterCallbackCallDiscard exercises the lifecycle void-call route.
func PortTestMonsterCallbackCallDiscard(key unsafe.Pointer, u *server.Object) {
	ccall.CallVoidPtr(key, u.CObj())
}
