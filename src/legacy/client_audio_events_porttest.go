//go:build porttest

package legacy

/*
#include <stdint.h>
#include "client/audio/ail/compat_mss.h"
#include "GAME1_3.h"
#include "GAME2.h"
#include "client__audio__audevent.h"
int sub_43DA80();
void sub_43DAD0();
int sub_43DB20();
int sub_43DB30(int a1);
char* sub_43DB40(int a1);
int sub_43DB60();
void sub_43DBA0();
void sub_43DC00();
int sub_43DC10();
int sub_43DC30();
void sub_43EDB0(HSAMPLE a1);
int sub_43EE00(void* a1p);
int sub_43F0E0(uint32_t* a1);
void sub_44D960();
int sub_44D970();
int sub_44D990();
unsigned char sub_450750();
char sub_450760(char a1);
int sub_451850(int a2, void* a3p);
int sub_451920(uint32_t* a2);
void sub_451970();
void sub_4519C0();
int sub_451BE0(int a1);
int sub_451CA0(uint32_t* a1);
int sub_451F30(int a1, int a2);
int sub_451F90(int a1);
int sub_451FE0(int a1);
int sub_452010();
void sub_452050(uint32_t* a1);
int* sub_452120(int a1);
void sub_452190(int a1);
int* sub_4521A0(int a1);
int sub_4521F0();
int***** sub_452230();
uint32_t* nox_xxx_draw_452300(uint32_t* a1);
int sub_4523D0(void* a1p);
int sub_452410(int a1);
int sub_452490(uint32_t* a1);
void sub_452510(int a3);
long long sub_452690(int a3, long long a4, int a5);
int sub_4526D0(int a1);
int sub_4526F0(int a1);
int* sub_452810(int a1, char a2);
void nox_xxx_clientPlaySoundSpecial_452D80(int a1, int a2);
void sub_452DC0(int a1, int a2, int a3);
void sub_452E10(int a1, int a2, int a3);
int sub_452E90(uint32_t* a1, int a2);
int sub_452EB0(int* a1);
int sub_452EE0(int a1, int a2);
unsigned int sub_452F10(int a1, int a2);
int sub_452F50(int a1, int a2);
uint32_t* sub_452F80(int a1, int a2);
int sub_452FA0(int a1);
int sub_452FE0(int a1, int a2);
void sub_453050();
void nox_xxx____setargv_9_453060();
int sub_453070();
int sub_451CF0(uint32_t* a1);
int sub_451DC0(int a1);
int sub_451E80(int a1);
int sub_452580(uint32_t* a1);
int sub_452770(uint32_t* a1);
extern uint32_t dword_587000_122848;
extern uint32_t dword_587000_126996;
extern void* dword_587000_127004;
extern uint32_t dword_587000_93156;
extern uint32_t dword_5d4594_1045420;
extern uint32_t dword_5d4594_1045424;
extern uint32_t dword_5d4594_1045428;
extern uint32_t dword_5d4594_1045432;
extern uint32_t dword_5d4594_1045436;
extern uint32_t dword_5d4594_816368;
extern uint32_t dword_5d4594_816372;
extern uint32_t dword_5d4594_816376;
extern uint32_t dword_5d4594_831092;
static uint64_t nox_porttest_audio_event_call(int op, uint64_t a0, uint64_t a1, uint64_t a2, uint64_t a3) {
switch(op) {
case 0: return (uint64_t)(int64_t)sub_43DA80();
case 1: sub_43DAD0(); return 0;
case 2: return (uint64_t)(int64_t)sub_43DB20();
case 3: return (uint64_t)(int64_t)sub_43DB30((int)a0);
case 4: return (uint64_t)(uintptr_t)sub_43DB40((int)a0);
case 5: return (uint64_t)(int64_t)sub_43DB60();
case 6: sub_43DBA0(); return 0;
case 7: sub_43DC00(); return 0;
case 8: return (uint64_t)(int64_t)sub_43DC10();
case 9: return (uint64_t)(int64_t)sub_43DC30();
case 10: sub_43EDB0((HSAMPLE)a0); return 0;
case 11: return (uint64_t)(int64_t)sub_43EE00((void*)a0);
case 12: return (uint64_t)(int64_t)sub_43F0E0((uint32_t*)a0);
case 13: sub_44D960(); return 0;
case 14: return (uint64_t)(int64_t)sub_44D970();
case 15: return (uint64_t)(int64_t)sub_44D990();
case 16: return (uint64_t)(int64_t)sub_450750();
case 17: return (uint64_t)(int64_t)sub_450760((char)a0);
case 18: return (uint64_t)(int64_t)sub_451850((int)a0,(void*)a1);
case 19: return (uint64_t)(int64_t)sub_451920((uint32_t*)a0);
case 20: sub_451970(); return 0;
case 21: sub_4519C0(); return 0;
case 22: return (uint64_t)(int64_t)sub_451BE0((int)a0);
case 23: return (uint64_t)(int64_t)sub_451CA0((uint32_t*)a0);
case 24: return (uint64_t)(int64_t)sub_451F30((int)a0,(int)a1);
case 25: return (uint64_t)(int64_t)sub_451F90((int)a0);
case 26: return (uint64_t)(int64_t)sub_451FE0((int)a0);
case 27: return (uint64_t)(int64_t)sub_452010();
case 28: sub_452050((uint32_t*)a0); return 0;
case 29: return (uint64_t)(uintptr_t)sub_452120((int)a0);
case 30: sub_452190((int)a0); return 0;
case 31: return (uint64_t)(uintptr_t)sub_4521A0((int)a0);
case 32: return (uint64_t)(int64_t)sub_4521F0();
case 33: return (uint64_t)(uintptr_t)sub_452230();
case 34: return (uint64_t)(uintptr_t)nox_xxx_draw_452300((uint32_t*)a0);
case 35: return (uint64_t)(int64_t)sub_4523D0((void*)a0);
case 36: return (uint64_t)(int64_t)sub_452410((int)a0);
case 37: return (uint64_t)(int64_t)sub_452490((uint32_t*)a0);
case 38: sub_452510((int)a0); return 0;
case 39: return (uint64_t)(int64_t)sub_452690((int)a0,(long long)a1,(int)a2);
case 40: return (uint64_t)(int64_t)sub_4526D0((int)a0);
case 41: return (uint64_t)(int64_t)sub_4526F0((int)a0);
case 42: return (uint64_t)(uintptr_t)sub_452810((int)a0,(char)a1);
case 43: nox_xxx_clientPlaySoundSpecial_452D80((int)a0,(int)a1); return 0;
case 44: sub_452DC0((int)a0,(int)a1,(int)a2); return 0;
case 45: sub_452E10((int)a0,(int)a1,(int)a2); return 0;
case 46: return (uint64_t)(int64_t)sub_452E90((uint32_t*)a0,(int)a1);
case 47: return (uint64_t)(int64_t)sub_452EB0((int*)a0);
case 48: return (uint64_t)(int64_t)sub_452EE0((int)a0,(int)a1);
case 49: return (uint64_t)(int64_t)sub_452F10((int)a0,(int)a1);
case 50: return (uint64_t)(int64_t)sub_452F50((int)a0,(int)a1);
case 51: return (uint64_t)(uintptr_t)sub_452F80((int)a0,(int)a1);
case 52: return (uint64_t)(int64_t)sub_452FA0((int)a0);
case 53: return (uint64_t)(int64_t)sub_452FE0((int)a0,(int)a1);
case 54: sub_453050(); return 0;
case 55: nox_xxx____setargv_9_453060(); return 0;
case 56: return (uint64_t)(int64_t)sub_453070();
case 57: return (uint64_t)(int64_t)sub_451CF0((uint32_t*)a0);
case 58: return (uint64_t)(int64_t)sub_451DC0((int)a0);
case 59: return (uint64_t)(int64_t)sub_451E80((int)a0);
case 60: return (uint64_t)(int64_t)sub_452580((uint32_t*)a0);
case 61: return (uint64_t)(int64_t)sub_452770((uint32_t*)a0);
} return 0;
}
static void* nox_porttest_audio_event_global(int op) {
switch(op) {
case 0: return &dword_587000_122848;
case 1: return &dword_587000_126996;
case 2: return &dword_587000_127004;
case 3: return &dword_587000_93156;
case 4: return &dword_5d4594_1045420;
case 5: return &dword_5d4594_1045424;
case 6: return &dword_5d4594_1045428;
case 7: return &dword_5d4594_1045432;
case 8: return &dword_5d4594_1045436;
case 9: return &dword_5d4594_816368;
case 10: return &dword_5d4594_816372;
case 11: return &dword_5d4594_816376;
case 12: return &dword_5d4594_831092;
} return 0;
}
*/
import "C"

var portTestAudioEventNames = []string{
	"sub_43DA80",
	"sub_43DAD0",
	"sub_43DB20",
	"sub_43DB30",
	"sub_43DB40",
	"sub_43DB60",
	"sub_43DBA0",
	"sub_43DC00",
	"sub_43DC10",
	"sub_43DC30",
	"sub_43EDB0",
	"sub_43EE00",
	"sub_43F0E0",
	"sub_44D960",
	"sub_44D970",
	"sub_44D990",
	"sub_450750",
	"sub_450760",
	"sub_451850",
	"sub_451920",
	"sub_451970",
	"sub_4519C0",
	"sub_451BE0",
	"sub_451CA0",
	"sub_451F30",
	"sub_451F90",
	"sub_451FE0",
	"sub_452010",
	"sub_452050",
	"sub_452120",
	"sub_452190",
	"sub_4521A0",
	"sub_4521F0",
	"sub_452230",
	"nox_xxx_draw_452300",
	"sub_4523D0",
	"sub_452410",
	"sub_452490",
	"sub_452510",
	"sub_452690",
	"sub_4526D0",
	"sub_4526F0",
	"sub_452810",
	"nox_xxx_clientPlaySoundSpecial_452D80",
	"sub_452DC0",
	"sub_452E10",
	"sub_452E90",
	"sub_452EB0",
	"sub_452EE0",
	"sub_452F10",
	"sub_452F50",
	"sub_452F80",
	"sub_452FA0",
	"sub_452FE0",
	"sub_453050",
	"nox_xxx____setargv_9_453060",
	"sub_453070",
	"sub_451CF0",
	"sub_451DC0",
	"sub_451E80",
	"sub_452580",
	"sub_452770",
}
var portTestAudioEventGlobals = []string{
	"dword_587000_122848",
	"dword_587000_126996",
	"dword_587000_127004",
	"dword_587000_93156",
	"dword_5d4594_1045420",
	"dword_5d4594_1045424",
	"dword_5d4594_1045428",
	"dword_5d4594_1045432",
	"dword_5d4594_1045436",
	"dword_5d4594_816368",
	"dword_5d4594_816372",
	"dword_5d4594_816376",
	"dword_5d4594_831092",
}

func PortTestAudioEventCall(name string, args ...uint64) uint64 {
	if len(args) > 4 {
		panic("too many audio event arguments")
	}
	var a [4]uint64
	copy(a[:], args)
	for i, n := range portTestAudioEventNames {
		if name == n {
			return uint64(C.nox_porttest_audio_event_call(C.int(i), C.uint64_t(a[0]), C.uint64_t(a[1]), C.uint64_t(a[2]), C.uint64_t(a[3])))
		}
	}
	panic("unknown audio event operation: " + name)
}
func PortTestAudioEventGlobals() (map[string]*uint32, func()) {
	words := make(map[string]*uint32)
	saved := make(map[string]uint32)
	for i, n := range portTestAudioEventGlobals {
		p := (*uint32)(C.nox_porttest_audio_event_global(C.int(i)))
		words[n] = p
		saved[n] = *p
		*p = 0
	}
	return words, func() {
		for n, p := range words {
			*p = saved[n]
		}
	}
}
