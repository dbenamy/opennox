//go:build porttest

package legacy

/*
#include "defs.h"
#include "client__system__parsecmd.h"
*/
import "C"
import (
	"unsafe"
)

func PortTestConsoleContext() func() {
	head := interactionMessageHead
	interactionMessageHead = 0
	mode, sender := consoleCommandServer, consoleCommandSender
	consoleCommandSender = nil
	return func() {
		interactionMessageHead = head
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
