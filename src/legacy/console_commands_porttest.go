//go:build porttest

package legacy

/*
#include "defs.h"
#include "client__system__parsecmd.h"
extern unsigned int nox_client_consoleIsServer_823684;
extern uint32_t dword_5d4594_825736;
extern nox_playerInfo* nox_console_playerWhoSent_823692;
*/
import "C"
import "unsafe"

func PortTestConsoleContext() func() {
	head := C.dword_5d4594_825736
	C.dword_5d4594_825736 = 0
	mode, sender := C.nox_client_consoleIsServer_823684, C.nox_console_playerWhoSent_823692
	C.nox_console_playerWhoSent_823692 = nil
	return func() {
		C.dword_5d4594_825736 = head
		C.nox_client_consoleIsServer_823684, C.nox_console_playerWhoSent_823692 = mode, sender
	}
}
func PortTestConsoleServer() bool { return C.nox_client_consoleIsServer_823684 != 0 }

func PortTestConsoleRemote(player unsafe.Pointer, action int, text *uint16) int {
	return int(C.nox_xxx_serverHandleClientConsole_443E90((*C.nox_playerInfo)(player), C.char(action), (*C.wchar2_t)(unsafe.Pointer(text))))
}
func PortTestConsoleSender() unsafe.Pointer {
	return unsafe.Pointer(C.nox_console_playerWhoSent_823692)
}
