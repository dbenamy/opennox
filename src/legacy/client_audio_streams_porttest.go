//go:build porttest

package legacy

/*
#include "GAME2_2.h"
#include "GAME3_1.h"
extern void* dword_587000_155144;
extern int nox_porttest_audio_stream_callback(int op, void* obj);
static int nox_porttest_audio_cb0(void* p) {return nox_porttest_audio_stream_callback(0,p);}
static int nox_porttest_audio_cb1(void* p) {return nox_porttest_audio_stream_callback(1,p);}
static int nox_porttest_audio_cb2(void* p) {return nox_porttest_audio_stream_callback(2,p);}
static int nox_porttest_audio_cb3(void* p) {return nox_porttest_audio_stream_callback(3,p);}
static int nox_porttest_audio_cb4(void* p) {return nox_porttest_audio_stream_callback(4,p);}
static int nox_porttest_audio_cb5(void* p) {return nox_porttest_audio_stream_callback(5,p);}
static int nox_porttest_audio_cb6(void* p) {return nox_porttest_audio_stream_callback(6,p);}
static int nox_porttest_audio_cb7(void* p) {return nox_porttest_audio_stream_callback(7,p);}
static int nox_porttest_audio_cb8(void* p) {return nox_porttest_audio_stream_callback(8,p);}
static int nox_porttest_audio_cb9(void* p) {return nox_porttest_audio_stream_callback(9,p);}
static int nox_porttest_audio_cb10(void* p) {return nox_porttest_audio_stream_callback(10,p);}
static int nox_porttest_audio_cb11(void* p) {return nox_porttest_audio_stream_callback(11,p);}
static int nox_porttest_audio_cb12(void* p) {return nox_porttest_audio_stream_callback(12,p);}
static int nox_porttest_audio_cb13(void* p) {return nox_porttest_audio_stream_callback(13,p);}

static void* nox_porttest_audio_callback_addr(int op) {
 switch(op) {
 case 0: return nox_porttest_audio_cb0;
 case 1: return nox_porttest_audio_cb1;
 case 2: return nox_porttest_audio_cb2;
 case 3: return nox_porttest_audio_cb3;
 case 4: return nox_porttest_audio_cb4;
 case 5: return nox_porttest_audio_cb5;
 case 6: return nox_porttest_audio_cb6;
 case 7: return nox_porttest_audio_cb7;
 case 8: return nox_porttest_audio_cb8;
 case 9: return nox_porttest_audio_cb9;
 case 10: return nox_porttest_audio_cb10;
 case 11: return nox_porttest_audio_cb11;
 case 12: return nox_porttest_audio_cb12;
 case 13: return nox_porttest_audio_cb13;
 default: return 0;
 }
}
static uint32_t nox_porttest_audio_stream_call(int op, uintptr_t a0, uintptr_t a1, uintptr_t a2, uintptr_t a3) {
 switch (op) {
case 0: return (uint32_t)(uintptr_t)sub_486640((void*)a0, (int)a1);
case 1: return (uint32_t)(uintptr_t)sub_4866D0((uint32_t*)a0, (int)a1);
case 2: return (uint32_t)(uintptr_t)sub_486AA0((uint32_t*)a0, (int)a1, (uint32_t*)a2);
case 3: return (uint32_t)(uintptr_t)sub_486B60((int)a0, (int)a1);
case 4: return (uint32_t)(uintptr_t)sub_486DB0((int)a0, (char*)a1, (signed int)a2);
case 5: return (uint32_t)(uintptr_t)sub_486E00((int)a0);
case 6: return (uint32_t)(uintptr_t)sub_486E30((int)a0, (uint32_t*)a1);
case 7: return (uint32_t)(uintptr_t)sub_486E90((int)a0);
case 8: return (uint32_t)(uintptr_t)sub_486FA0((int)a0);
case 9: return (uint32_t)(uintptr_t)sub_486FE0((int)a0);
case 10: sub_487030((void*)a0); return 0;
case 11: sub_487050((uint32_t*)a0); return 0;
case 12: sub_487070((void*)a0); return 0;
case 13: sub_487090((uint32_t**)a0); return 0;
case 14: sub_4870A0(); return 0;
case 15: return (uint32_t)(uintptr_t)sub_4870E0((int*)a0);
case 16: return (uint32_t)(uintptr_t)sub_487100((int**)a0);
case 17: return (uint32_t)(uintptr_t)sub_487150((int)a0, (const void*)a1);
case 18: return (uint32_t)(uintptr_t)sub_4871C0((int)a0, (int)a1, (const void*)a2);
case 19: sub_4872C0((void*)a0); return 0;
case 20: return (uint32_t)(uintptr_t)sub_487310((uint32_t*)a0);
case 21: return (uint32_t)(uintptr_t)sub_487360((int)a0, (int**)a1, (int*)a2);
case 22: return (uint32_t)(uintptr_t)sub_4873C0((int)a0);
case 23: return (uint32_t)(uintptr_t)sub_487590((int)a0, (const void*)a1);
case 24: return (uint32_t)(uintptr_t)sub_4875B0((int*)a0);
case 25: return (uint32_t)(uintptr_t)sub_4875D0((int**)a0);
case 26: return (uint32_t)(uintptr_t)sub_4875F0();
case 27: sub_487680((void*)a0); return 0;
case 28: return (uint32_t)(uintptr_t)sub_4876A0((uint32_t**)a0);
case 29: return (uint32_t)(uintptr_t)sub_487750((int)a0);
case 30: return (uint32_t)(uintptr_t)sub_487790((int)a0, (int)a1);
case 31: return (uint32_t)(uintptr_t)sub_4877D0((int)a0, (int*)a1);
case 32: return (uint32_t)(uintptr_t)sub_4877F0((int**)a0);
case 33: return (uint32_t)(uintptr_t)sub_487810((int)a0, (int)a1);
case 34: return (uint32_t)(uintptr_t)sub_487910((int)a0, (int)a1);
case 35: return (uint32_t)(uintptr_t)sub_487970((int)a0, (int)a1);
case 36: sub_487C30((uint32_t*)a0); return 0;
case 37: return (uint32_t)(uintptr_t)sub_487C50((int)a0, (uint32_t*)a1);
case 38: return (uint32_t)(uintptr_t)sub_487C80((int)a0);
case 39: return (uint32_t)(uintptr_t)sub_487D00((uint32_t*)a0);
case 40: return (uint32_t)(uintptr_t)sub_487D30((uint32_t*)a0, (int)a1, (int)a2);
case 41: return (uint32_t)(uintptr_t)sub_487D60((int)a0);
case 42: return (uint32_t)(uintptr_t)sub_4BD280((int)a0, (int)a1);
case 43: sub_4BD2D0((void*)a0); return 0;
case 44: return (uint32_t)(uintptr_t)sub_4BD2E0((uint32_t**)a0);
case 45: return (uint32_t)(uintptr_t)sub_4BD300((uint32_t*)a0, (int)a1);
case 46: return (uint32_t)(uintptr_t)sub_4BD340((int)a0, (int)a1, (int)a2, (int)a3);
case 47: sub_4BD3C0((void*)a0); return 0;
case 48: return (uint32_t)(uintptr_t)sub_4BD420((int)a0, (int)a1);
case 49: return (uint32_t)(uintptr_t)sub_4BD470((uint32_t**)a0, (int)a1);
case 50: return (uint32_t)(uintptr_t)sub_4BD600((int)a0);
case 51: return (uint32_t)(uintptr_t)sub_4BD650((int)a0);
case 52: return (uint32_t)(uintptr_t)sub_4BD660((int)a0);
case 53: return (uint32_t)(uintptr_t)sub_4BD680((int)a0);
case 54: return (uint32_t)(uintptr_t)sub_4BD690((int)a0);
case 55: return (uint32_t)(uintptr_t)sub_4BD710((int)a0);
case 56: return (uint32_t)(uintptr_t)sub_4BD720((int)a0);
case 57: sub_4BD7A0((void*)a0); return 0;
case 58: return (uint32_t)(uintptr_t)sub_4BD7C0((uint32_t*)a0);
case 59: sub_4BD840((int)a0); return 0;
case 60: return (uint32_t)(uintptr_t)sub_4BD8C0((int)a0);
case 61: return (uint32_t)(uintptr_t)sub_4BD940((int)a0);
case 62: return (uint32_t)(uintptr_t)sub_4BD9B0((uint32_t*)a0);
case 63: sub_4BDA60((void*)a0); return 0;
case 64: return (uint32_t)(uintptr_t)sub_4BDA80((int)a0);
case 65: return (uint32_t)(uintptr_t)sub_4BDB20((int)a0);
case 66: return (uint32_t)(uintptr_t)sub_4BDB30((int)a0);
case 67: return (uint32_t)(uintptr_t)sub_4BDB40((int)a0);
case 68: sub_4BDB90((uint32_t*)a0, (uint32_t*)a1); return 0;
case 69: return (uint32_t)(uintptr_t)sub_4BDC00((int)a0);
 }
 return 0;
}
*/
import "C"
import "unsafe"

var portTestAudioStreamNames = []string{
	"sub_486640",
	"sub_4866D0",
	"sub_486AA0",
	"sub_486B60",
	"sub_486DB0",
	"sub_486E00",
	"sub_486E30",
	"sub_486E90",
	"sub_486FA0",
	"sub_486FE0",
	"sub_487030",
	"sub_487050",
	"sub_487070",
	"sub_487090",
	"sub_4870A0",
	"sub_4870E0",
	"sub_487100",
	"sub_487150",
	"sub_4871C0",
	"sub_4872C0",
	"sub_487310",
	"sub_487360",
	"sub_4873C0",
	"sub_487590",
	"sub_4875B0",
	"sub_4875D0",
	"sub_4875F0",
	"sub_487680",
	"sub_4876A0",
	"sub_487750",
	"sub_487790",
	"sub_4877D0",
	"sub_4877F0",
	"sub_487810",
	"sub_487910",
	"sub_487970",
	"sub_487C30",
	"sub_487C50",
	"sub_487C80",
	"sub_487D00",
	"sub_487D30",
	"sub_487D60",
	"sub_4BD280",
	"sub_4BD2D0",
	"sub_4BD2E0",
	"sub_4BD300",
	"sub_4BD340",
	"sub_4BD3C0",
	"sub_4BD420",
	"sub_4BD470",
	"sub_4BD600",
	"sub_4BD650",
	"sub_4BD660",
	"sub_4BD680",
	"sub_4BD690",
	"sub_4BD710",
	"sub_4BD720",
	"sub_4BD7A0",
	"sub_4BD7C0",
	"sub_4BD840",
	"sub_4BD8C0",
	"sub_4BD940",
	"sub_4BD9B0",
	"sub_4BDA60",
	"sub_4BDA80",
	"sub_4BDB20",
	"sub_4BDB30",
	"sub_4BDB40",
	"sub_4BDB90",
	"sub_4BDC00",
}

func PortTestAudioStreamCall(name string, args ...uint32) uint32 {
	if len(args) > 4 {
		panic("too many audio stream arguments")
	}
	var a [4]uint32
	copy(a[:], args)
	for op, n := range portTestAudioStreamNames {
		if n == name {
			return uint32(C.nox_porttest_audio_stream_call(C.int(op), C.uintptr_t(a[0]), C.uintptr_t(a[1]), C.uintptr_t(a[2]), C.uintptr_t(a[3])))
		}
	}
	panic("unknown audio stream operation: " + name)
}

var portTestAudioStreamCallback func(int, unsafe.Pointer) int

//export nox_porttest_audio_stream_callback
func nox_porttest_audio_stream_callback(op C.int, obj unsafe.Pointer) C.int {
	if portTestAudioStreamCallback == nil {
		panic("audio callback without an owner")
	}
	return C.int(portTestAudioStreamCallback(int(op), obj))
}

func PortTestAudioStreamGlobalOwner(p unsafe.Pointer, fn func(int, unsafe.Pointer) int) ([]unsafe.Pointer, func()) {
	old, callback := C.dword_587000_155144, portTestAudioStreamCallback
	C.dword_587000_155144, portTestAudioStreamCallback = p, fn
	return []unsafe.Pointer{
		C.nox_porttest_audio_callback_addr(0),
		C.nox_porttest_audio_callback_addr(1),
		C.nox_porttest_audio_callback_addr(2),
		C.nox_porttest_audio_callback_addr(3),
		C.nox_porttest_audio_callback_addr(4),
		C.nox_porttest_audio_callback_addr(5),
		C.nox_porttest_audio_callback_addr(6),
		C.nox_porttest_audio_callback_addr(7),
		C.nox_porttest_audio_callback_addr(8),
		C.nox_porttest_audio_callback_addr(9),
		C.nox_porttest_audio_callback_addr(10),
		C.nox_porttest_audio_callback_addr(11),
		C.nox_porttest_audio_callback_addr(12),
		C.nox_porttest_audio_callback_addr(13),
	}, func() { C.dword_587000_155144, portTestAudioStreamCallback = old, callback }
}
