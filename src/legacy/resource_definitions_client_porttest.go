//go:build porttest

package legacy

/*
#include "defs.h"
#include "memfile.h"
extern uint64_t qword_581450_9544;
extern uint64_t qword_581450_9552;
*/
import "C"
import (
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/internal/binfile"
	"unsafe"
)

func PortTestResourceLightConstants() (*uint64, *uint64) {
	return (*uint64)(unsafe.Pointer(&C.qword_581450_9544)), (*uint64)(unsafe.Pointer(&C.qword_581450_9552))
}

func PortTestResourceClientParser(kind string, typ *client.ObjectType, f *binfile.MemFile, input unsafe.Pointer) bool {
	return resourceClientParser(kind, typ, f, (*byte)(input))
}
