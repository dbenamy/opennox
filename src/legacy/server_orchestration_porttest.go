//go:build porttest

package legacy

/*
#include "defs.h"
void sub_4D2160();
void sub_4D22B0();
void sub_4DBA30(int);
int nox_xxx_mapLoadRequired_4DCC80();
void sub_4E4170();
int sub_4EDD70();
int sub_4EF660(nox_object_t*);
void sub_4F1F20();
void nox_xxx_updateUnits_51B100_D();
bool sub_57B140();
extern uint32_t dword_5d4594_1563096;
extern uint32_t dword_5d4594_1568300;
extern uint32_t dword_5d4594_2488728;
extern uint32_t dword_5d4594_1568280, dword_5d4594_1568288;
int nox_xxx_initChest_4F0400(int);
static void* orchestrationChestInit(void) { return nox_xxx_initChest_4F0400; }
*/
import "C"
import "github.com/opennox/opennox/v1/server"
import "unsafe"

func PortTestServerOrchestration(op string, u *server.Object, arg int32) uint32 {
	switch op {
	case "flags":
		C.sub_4D2160()
	case "players":
		C.sub_4D22B0()
	case "restore":
		C.sub_4DBA30(C.int(arg))
	case "load-state":
		return uint32(C.nox_xxx_mapLoadRequired_4DCC80())
	case "difficulty":
		C.sub_4E4170()
	case "drop-flags":
		return uint32(C.sub_4EDD70())
	case "reset-player":
		return uint32(C.sub_4EF660(asObjectC(u)))
	case "rewards":
		C.sub_4F1F20()
	case "walls":
		C.nox_xxx_updateUnits_51B100_D()
	case "timeout":
		return uint32(bool2int(bool(C.sub_57B140())))
	}
	return 0
}

func PortTestServerOrchestrationGlobals() (map[string]*uint32, func()) {
	words := map[string]*uint32{"ankh-marker": (*uint32)(unsafe.Pointer(&C.dword_5d4594_1568280)), "selected-marker": (*uint32)(unsafe.Pointer(&C.dword_5d4594_1568288)), "drop-table": (*uint32)(unsafe.Pointer(&C.dword_5d4594_2488728)), "restore-cleanup": (*uint32)(unsafe.Pointer(&C.dword_5d4594_1563096)), "reward-marker": (*uint32)(unsafe.Pointer(&C.dword_5d4594_1568300))}
	saved := map[string]uint32{}
	for k, p := range words {
		saved[k] = *p
		*p = 0
	}
	return words, func() {
		for k, p := range words {
			*p = saved[k]
		}
	}
}

func PortTestServerOrchestrationChestInit() unsafe.Pointer { return C.orchestrationChestInit() }
