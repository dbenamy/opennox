package legacy

/*
#include "client__gui__guiquit.h"
#include "GAME1_1.h"
#include "GAME1_3.h"
#include "GAME3_2.h"
#include "GAME4_1.h"
extern unsigned int dword_5d4594_2650652;
extern uint32_t dword_5d4594_830272;
int* nox_xxx_guiServerOptionsHide_4597E0(int a1);
*/
import "C"
import (
	"unsafe"

	"github.com/opennox/opennox/v1/client/gui"
)

var (
	Sub_446380 func()
	Sub_445B40 func() int
)

//export nox_xxx____setargv_4_44B000
func nox_xxx____setargv_4_44B000() {
	C.dword_5d4594_830272 = 1
}

//export sub_446380
func sub_446380() { Sub_446380() }

//export sub_445B40
func sub_445B40() int { return Sub_445B40() }

func Sub_4D70B0() {
	questRuntimeSettings()
}

func Sub_509CB0() {
	matchRosterForget()
}

func Sub_41F4B0() {
	onlineSessionListCleanup()
}

func Sub_41EC30() {
	onlineSessionListCleanup()
}

func Sub_446490(v int) {
	sessionMOTDFree(v)
}

func Nox_xxx_guiServerOptionsHide_4597E0(v int) {
	serverOptionsClose(v)
}

func Sub_445C40() {
	sessionQuitToggle()
}

func Set_nox_wnd_quitMenu_825760(win *gui.Window) {
	sessionQuitRoot = win
}

func Nox_xxx____setargv_4_44B000() {
	nox_xxx____setargv_4_44B000()
}

func Get_dword_5d4594_2650652() int {
	return int(C.dword_5d4594_2650652)
}

func Set_dword_5d4594_2650652(v int) {
	C.dword_5d4594_2650652 = C.uint(v)
}

func Sub_41CEE0(p unsafe.Pointer, a2 int) {
	playerFileClientWrite(p, a2)
}
