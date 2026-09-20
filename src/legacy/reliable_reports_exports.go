package legacy

/*
#include "defs.h"
*/
import "C"
import (
	"github.com/opennox/opennox/v1/common/ntype"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

//export sub_4E4F30
func sub_4E4F30(to C.int) C.int { return C.int(ReliableResetSequence(ntype.PlayerInd(to))) }

//export nox_xxx_playerResetImportantCtr_4E4F40
func nox_xxx_playerResetImportantCtr_4E4F40(to C.int) C.int { return C.int(reliableResetRate(int(to))) }

//export nox_xxx_netSendPacket1_4E5390
func nox_xxx_netSendPacket1_4E5390(to, data, size, related, priority C.int) C.int {
	return C.int(reliableEnqueue(int(to), unsafe.Slice((*byte)(unsafe.Pointer(uintptr(uint32(data)))), int(size)), (*server.Object)(unsafe.Pointer(uintptr(uint32(related)))), int(priority), 1))
}

//export nox_xxx_netClientSend2_4E53C0
func nox_xxx_netClientSend2_4E53C0(to C.int, data unsafe.Pointer, size, related, priority C.int) C.int {
	return C.int(reliableClientSend(int(to), unsafe.Slice((*byte)(data), int(size)), (*server.Object)(unsafe.Pointer(uintptr(uint32(related)))), int(priority)))
}

func nox_xxx_netSendPacket0_4E5420(to C.int, data unsafe.Pointer, size, related, priority C.int) C.int {
	return C.int(reliableEnqueue(int(to), unsafe.Slice((*byte)(data), int(size)), (*server.Object)(unsafe.Pointer(uintptr(uint32(related)))), int(priority), 0))
}

//export nox_net_importantACK_4E55A0
func nox_net_importantACK_4E55A0(to, frame C.int) C.int {
	return C.int(reliableACK(int(to), uint32(frame)))
}

//export sub_4E55F0
func sub_4E55F0(to C.uchar) C.int { return C.int(reliableRemoveRecipient(byte(to))) }
