//go:build porttest

package legacy

/*
#include "defs.h"
#include "GAME3.h"
extern nox_gui_animation* nox_wnd_xxx_1309740;
extern uint32_t nox_xxx_normalWndBits_587000_172880;
void sub_4AA650();
int nox_porttest_options_done();
extern void* dword_5d4594_1309720;
extern uint32_t dword_5d4594_1309728;
extern uint32_t dword_5d4594_1309732;
extern uint32_t dword_5d4594_1309736;
extern uint32_t dword_5d4594_1309820;
extern uint32_t dword_5d4594_1309824;
extern uint32_t dword_5d4594_1309828;
extern uint32_t dword_5d4594_1309832;
extern uint32_t dword_5d4594_1309836;
extern uint32_t dword_587000_126996;
extern uint32_t dword_587000_122848;
extern uint32_t dword_587000_93156;
extern uint32_t dword_5d4594_831092;
extern uint32_t dword_5d4594_816376;
extern void* dword_587000_127004;
extern void* dword_587000_122852;
extern void* dword_587000_93164;

*/
import "C"
import (
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/legacy/timer"
	"unsafe"
)

// PortTestOptionsEvent invokes the installed production owner, initially C.
func PortTestOptionsEvent(menu bool, root *gui.Window, event int, child *gui.Window, value int) int {
	if menu {
		return int(C.sub_4AABE0(C.int(uintptr(root.C())), C.int(event), (*C.int)(child.C()), C.int(value)))
	}
	return int(C.nox_xxx_windowOptionsProc_4ADF30(C.int(uintptr(root.C())), C.int(event), (*C.int)(child.C()), C.int(value)))
}

// PortTestOptionsWords isolates live UI/audio globals from their backing blobs.
func PortTestOptionsWords() (map[int]*uint32, func()) {
	words := map[int]*uint32{
		172880:  (*uint32)(unsafe.Pointer(&C.nox_xxx_normalWndBits_587000_172880)),
		1309720: (*uint32)(unsafe.Pointer(&C.dword_5d4594_1309720)),
		1309728: (*uint32)(unsafe.Pointer(&C.dword_5d4594_1309728)),
		1309732: (*uint32)(unsafe.Pointer(&C.dword_5d4594_1309732)),
		1309736: (*uint32)(unsafe.Pointer(&C.dword_5d4594_1309736)),
		1309820: (*uint32)(unsafe.Pointer(&C.dword_5d4594_1309820)),
		1309824: (*uint32)(unsafe.Pointer(&C.dword_5d4594_1309824)),
		1309828: (*uint32)(unsafe.Pointer(&C.dword_5d4594_1309828)),
		1309832: (*uint32)(unsafe.Pointer(&C.dword_5d4594_1309832)),
		1309836: (*uint32)(unsafe.Pointer(&C.dword_5d4594_1309836)),
		126996:  (*uint32)(unsafe.Pointer(&C.dword_587000_126996)),
		122848:  (*uint32)(unsafe.Pointer(&C.dword_587000_122848)),
		93156:   (*uint32)(unsafe.Pointer(&C.dword_587000_93156)),
		831092:  (*uint32)(unsafe.Pointer(&C.dword_5d4594_831092)),
		816376:  (*uint32)(unsafe.Pointer(&C.dword_5d4594_816376)),
	}
	old := make(map[int]uint32)
	for n, p := range words {
		old[n] = *p
		*p = 0
	}
	return words, func() {
		for n, p := range words {
			*p = old[n]
		}
	}
}
func PortTestOptionsTimers() ([3]*timer.Timer, func()) {
	slots := [3]*unsafe.Pointer{&C.dword_587000_127004, &C.dword_587000_122852, &C.dword_587000_93164}
	var old [3]unsafe.Pointer
	var out [3]*timer.Timer
	var frees [3]func()
	for i, p := range slots {
		old[i] = *p
		out[i], frees[i] = alloc.New(timer.Timer{})
		*p = unsafe.Pointer(out[i])
	}
	return out, func() {
		for i, p := range slots {
			*p = old[i]
			frees[i]()
		}
	}
}

func PortTestOptionsConstruct(menu bool) int {
	if menu {
		return int(C.nox_game_showOptions_4AA6B0())
	}
	return int(C.nox_game_initOptionsInGame_4ADAD0())
}
func PortTestOptionsAction(op, arg int) int {
	switch op {
	case 0:
		return int(uintptr(unsafe.Pointer(C.sub_4AAA70())))
	case 1:
		return int(C.sub_4ADA40())
	case 2:
		return int(C.sub_4AD9B0(C.int(arg)))
	case 3:
		return int(C.sub_4AE3D0())
	case 4:
		return int(C.sub_4AE3B0())
	case 5:
		return int(C.sub_4AB0C0())
	case 6:
		C.sub_4AA650()
		return 0
	}
	panic(op)
}
func PortTestOptionsDraw(w *gui.Window) int { return int(C.sub_4ADEF0((*C.uint32_t)(w.C()), 0)) }
func PortTestOptionsAnimWord() *uint32      { return (*uint32)(unsafe.Pointer(&C.nox_wnd_xxx_1309740)) }

var portTestOptionsDoneCount int

//export nox_porttest_options_done
func nox_porttest_options_done() C.int { portTestOptionsDoneCount++; return 1 }
func PortTestOptionsDone() (unsafe.Pointer, func() int, func()) {
	old := portTestOptionsDoneCount
	portTestOptionsDoneCount = 0
	return C.nox_porttest_options_done, func() int { return portTestOptionsDoneCount }, func() { portTestOptionsDoneCount = old }
}
