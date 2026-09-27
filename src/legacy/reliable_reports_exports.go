package legacy

/*
#include "defs.h"
*/
import "C"
import (
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

func nox_xxx_netClientSend2_4E53C0(to C.int, data unsafe.Pointer, size, related, priority C.int) C.int {
	return C.int(reliableClientSend(int(to), unsafe.Slice((*byte)(data), int(size)), (*server.Object)(unsafe.Pointer(uintptr(uint32(related)))), int(priority)))
}

func nox_xxx_netSendPacket0_4E5420(to C.int, data unsafe.Pointer, size, related, priority C.int) C.int {
	return C.int(reliableEnqueue(int(to), unsafe.Slice((*byte)(data), int(size)), (*server.Object)(unsafe.Pointer(uintptr(uint32(related)))), int(priority), 0))
}

func nox_net_importantACK_4E55A0(to, frame C.int) C.int {
	return C.int(reliableACK(int(to), uint32(frame)))
}
