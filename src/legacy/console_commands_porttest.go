//go:build porttest

package legacy

/*
#include "defs.h"
#include "client__system__parsecmd.h"
extern uint32_t dword_5d4594_825736;
*/
import "C"
import (
	"unsafe"
)

func PortTestConsoleContext() func() {
	head := C.dword_5d4594_825736
	C.dword_5d4594_825736 = 0
	mode, sender := consoleCommandServer, consoleCommandSender
	consoleCommandSender = nil
	return func() {
		C.dword_5d4594_825736 = head
		consoleCommandServer, consoleCommandSender = mode, sender
	}
}
func PortTestConsoleServer() bool { return consoleCommandServer }

func PortTestConsoleRemote(player unsafe.Pointer, action int, text *uint16) int {
	return int(C.nox_xxx_serverHandleClientConsole_443E90((*C.nox_playerInfo)(player), C.char(action), (*C.wchar2_t)(unsafe.Pointer(text))))
}
func PortTestConsoleSender() unsafe.Pointer {
	return unsafe.Pointer(consoleCommandSender)
}
