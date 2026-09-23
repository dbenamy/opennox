//go:build porttest

package legacy

/*
#include <stdint.h>

static uint32_t ptDamageForwardWords[5];
static int32_t ptDamageForwardResult;
static int ptDamageForwardCount;
static int ptDamageForwardCallback;

static int32_t ptDamageForwardRecord(int callback, uintptr_t a, uintptr_t b, uintptr_t c, uintptr_t d, uintptr_t e) {
	ptDamageForwardCount++;
	ptDamageForwardCallback = callback;
	ptDamageForwardWords[0] = (uint32_t)a;
	ptDamageForwardWords[1] = (uint32_t)b;
	ptDamageForwardWords[2] = (uint32_t)c;
	ptDamageForwardWords[3] = (uint32_t)d;
	ptDamageForwardWords[4] = (uint32_t)e;
	return ptDamageForwardResult;
}

static int32_t ptDamageForwardA(uintptr_t a, uintptr_t b, uintptr_t c, uintptr_t d, uintptr_t e) {
	return ptDamageForwardRecord(1, a, b, c, d, e);
}
static int32_t ptDamageForwardB(uintptr_t a, uintptr_t b, uintptr_t c, uintptr_t d, uintptr_t e) {
	return ptDamageForwardRecord(2, a, b, c, d, e);
}
static void* ptDamageForwardAPtr(void) { return (void*)ptDamageForwardA; }
static void* ptDamageForwardBPtr(void) { return (void*)ptDamageForwardB; }
static void ptDamageForwardReset(int32_t result) {
	ptDamageForwardCount = 0;
	ptDamageForwardCallback = 0;
	for (int i = 0; i < 5; i++) ptDamageForwardWords[i] = 0;
	ptDamageForwardResult = result;
}
static int ptDamageForwardGetCount(void) { return ptDamageForwardCount; }
static int ptDamageForwardGetCallback(void) { return ptDamageForwardCallback; }
static uint32_t ptDamageForwardGetWord(int i) { return ptDamageForwardWords[i]; }
*/
import "C"
import "unsafe"

type PortTestDamageForwardState struct {
	Count    int
	Callback int
	Words    [5]uint32
}

func PortTestDamageForwardCallback(which int) unsafe.Pointer {
	switch which {
	case 0:
		return C.ptDamageForwardAPtr()
	case 1:
		return C.ptDamageForwardBPtr()
	default:
		panic("unknown damage forwarding callback")
	}
}

func PortTestDamageForwardReset(result int32) {
	C.ptDamageForwardReset(C.int32_t(result))
}

func PortTestDamageForwardSnapshot() (out PortTestDamageForwardState) {
	out.Count = int(C.ptDamageForwardGetCount())
	out.Callback = int(C.ptDamageForwardGetCallback())
	for i := range out.Words {
		out.Words[i] = uint32(C.ptDamageForwardGetWord(C.int(i)))
	}
	return out
}
