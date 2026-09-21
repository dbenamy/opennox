//go:build porttest

package legacy

/*
#include <stdint.h>
#include "GAME1_3.h"
#include "GAME2_2.h"
extern uint32_t dword_5d4594_815748;
extern uint32_t dword_5d4594_816412;
*/
import "C"
import "unsafe"

func PortTestClientResourceWords() (map[string]*uint32, func()) {
	words := map[string]*uint32{
		"frameCount":  (*uint32)(unsafe.Pointer(&C.dword_5d4594_815748)),
		"modalWindow": (*uint32)(unsafe.Pointer(&C.dword_5d4594_816412)),
	}
	old := make(map[string]uint32)
	for k, p := range words {
		old[k] = *p
	}
	return words, func() {
		for k, p := range words {
			*p = old[k]
		}
	}
}
func PortTestClientFrameAverage()      { C.sub_43CEB0() }
func PortTestClientShellState() int32  { return int32(C.sub_43BDB0()) }
func PortTestClientBindingCount() byte { return byte(C.sub_47DBC0()) }
