package legacy

/*
#include "GAME2_1.h"
#include "client__system__parsecmd.h"
*/
import "C"
import "github.com/opennox/opennox/v1/server"

func Nox_xxx_serverHandleClientConsole_443E90(a1 *server.Player, a2 byte, cmd string) {
	consoleCommandRemote(a1, a2, cmd)
}
func Nox_xxx_cmdSayDo_46A4B0(text string, a2 int) {
	cstr, free := CWString(text)
	defer free()
	nox_xxx_cmdSayDo_46A4B0(cstr, C.int(a2))
}
