//go:build porttest

package legacy

/*
#include <stdint.h>
static uint32_t worldDamageWords[6];
static int worldDamageObserve(void* target, void* owner, void* source, int amount, int kind) {
 worldDamageWords[0]++;
 worldDamageWords[1]=(uintptr_t)target; worldDamageWords[2]=(uintptr_t)owner;
 worldDamageWords[3]=(uintptr_t)source; worldDamageWords[4]=amount; worldDamageWords[5]=kind;
 return 1;
}
static void* worldDamageCallback(void) { return worldDamageObserve; }
static uint32_t* worldDamageSnapshot(void) { return worldDamageWords; }
*/
import "C"
import "unsafe"

// Observe the selected collision's raw damage dispatch; no damage algorithm.
func PortTestWorldDamageObserver() (unsafe.Pointer, *[6]uint32, func()) {
	p := (*[6]uint32)(unsafe.Pointer(C.worldDamageSnapshot()))
	old := *p
	*p = [6]uint32{}
	return C.worldDamageCallback(), p, func() { *p = old }
}
