//go:build porttest

package legacy

/*
#include <math.h>
*/
import "C"

// These observers call the same original libc entrypoints as the placement owners.
func portTestPlacementSin(x float64) float64 { return float64(C.sin(C.double(x))) }
func portTestPlacementCos(x float64) float64 { return float64(C.cos(C.double(x))) }
