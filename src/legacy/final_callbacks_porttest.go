//go:build porttest

package legacy

/*
extern int nox_porttest_final_ptr(void*);
extern int nox_porttest_final_ptr2(void*, void*);
extern int sub_479D00(void);
extern void nox_xxx_updateFlameCleanse_53D510(int);
*/
import "C"

import (
	"unsafe"

	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/legacy/common/ccall"
)

var portTestFinalPtr func(unsafe.Pointer) int32
var portTestFinalPtr2 func(unsafe.Pointer, unsafe.Pointer) int32

//export nox_porttest_final_ptr
func nox_porttest_final_ptr(p unsafe.Pointer) C.int { return C.int(portTestFinalPtr(p)) }

//export nox_porttest_final_ptr2
func nox_porttest_final_ptr2(a, b unsafe.Pointer) C.int { return C.int(portTestFinalPtr2(a, b)) }

// These callbacks observe the foreign boundary; they contain no engine algorithms.
func PortTestObserveFinalPtr(fn func(unsafe.Pointer) int32) (unsafe.Pointer, func()) {
	old := portTestFinalPtr
	portTestFinalPtr = fn
	return C.nox_porttest_final_ptr, func() { portTestFinalPtr = old }
}

func PortTestObserveFinalPtr2(fn func(unsafe.Pointer, unsafe.Pointer) int32) (unsafe.Pointer, func()) {
	old := portTestFinalPtr2
	portTestFinalPtr2 = fn
	return C.nox_porttest_final_ptr2, func() { portTestFinalPtr2 = old }
}

func PortTestFinalMenuCallback() gui.WindowFunc             { return gui.WrapFuncC(Get_sub_4A18E0()) }
func PortTestFinalConversationTooltip() unsafe.Pointer      { return C.sub_479D00 }
func PortTestFinalFlameCallback() unsafe.Pointer            { return C.nox_xxx_updateFlameCleanse_53D510 }
func PortTestFinalPlayerSectionCall(key unsafe.Pointer) int { return ccall.CallIntPtr(key, nil) }
