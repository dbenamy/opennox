//go:build porttest

package legacy

/*
#include <stdlib.h>
*/
import "C"
import "unsafe"

// The C object reader allocates names with libc calloc, outside alloc's tracker.
func PortTestObjectXferFreeName(p unsafe.Pointer) { C.free(p) }
