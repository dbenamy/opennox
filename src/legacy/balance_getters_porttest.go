//go:build porttest

// Calls the existing production C
// exports and their existing same-package Go implementations; no C algorithm
// or duplicate lookup implementation is introduced.
package legacy

/*
#include "GAME1_1.h"
*/
import "C"

type PortTestBalanceGetterPair struct {
	ScalarC, ScalarGo float64
	IndexC, IndexGo   float64
}

func PortTestBalanceGetters(key string, index int) PortTestBalanceGetterPair {
	// internCStr matches existing production callers and keeps this borrowed
	// C string alive through both immediate calls.
	k := internCStr(key)
	ci := C.int(index)
	return PortTestBalanceGetterPair{
		ScalarC:  float64(C.nox_xxx_gamedataGetFloat_419D40(k)),
		ScalarGo: float64(nox_xxx_gamedataGetFloat_419D40(k)),
		IndexC:   float64(C.nox_xxx_gamedataGetFloatTable_419D70(k, ci)),
		IndexGo:  float64(nox_xxx_gamedataGetFloatTable_419D70(k, int(int32(index)))),
	}
}
