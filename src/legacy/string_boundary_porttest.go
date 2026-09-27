//go:build porttest

package legacy

/*
#include <stdlib.h>
*/
import "C"
import "unsafe"

// CString uses this same cgo malloc path, including the safe-profile macro.
func portTestStringMalloc(n uintptr) unsafe.Pointer { return C.malloc(C.size_t(n)) }
