//go:build porttest

package legacy

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
	return int(nox_xxx_serverHandleClientConsole_443E90((*nox_playerInfo)(player), int8(action), (*uint16)(unsafe.Pointer(text))))
}
func PortTestConsoleSender() unsafe.Pointer {
	return unsafe.Pointer(consoleCommandSender)
}
