//go:build porttest

package legacy

/*
#include <stdint.h>
static uintptr_t nox_porttest_tooltip_words[4];
static void nox_porttest_tooltip_probe(uintptr_t window, uintptr_t draw, uintptr_t argument) {
    nox_porttest_tooltip_words[0]++;
    nox_porttest_tooltip_words[1] = window;
    nox_porttest_tooltip_words[2] = draw;
    nox_porttest_tooltip_words[3] = argument;
}
static void* nox_porttest_tooltip_address(void) { return (void*)nox_porttest_tooltip_probe; }
static uintptr_t* nox_porttest_tooltip_state(void) { return nox_porttest_tooltip_words; }
*/
import "C"
import "unsafe"

// PortTestGUITooltipProbe observes the foreign callback ABI independently of
// the native tooltip implementation. The returned cleanup restores its state.
func PortTestGUITooltipProbe() (unsafe.Pointer, *[4]uint32, func()) {
	words := (*[4]uint32)(unsafe.Pointer(C.nox_porttest_tooltip_state()))
	saved := *words
	*words = [4]uint32{}
	return C.nox_porttest_tooltip_address(), words, func() { *words = saved }
}
