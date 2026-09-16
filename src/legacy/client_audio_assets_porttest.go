//go:build porttest

package legacy

/*
#include <stdint.h>
#include "GAME2.h"
#include "GAME2_2.h"
extern uint32_t dword_5d4594_1045420;
extern uint32_t dword_5d4594_1045428;
extern uint32_t dword_5d4594_1045432;
*/
import "C"
import (
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/internal/binfile"
	"unsafe"
)

func PortTestClientAudioAssetsOwner() ([]byte, *uint32, *unsafe.Pointer, func()) {
	rows := unsafe.Slice(memmap.PtrUint8(0x5D4594, 840628), 200*1023)
	saved := append([]byte(nil), rows...)
	catalog, context, enabled := C.dword_5d4594_1045420, C.dword_5d4594_1045428, C.dword_5d4594_1045432
	clear(rows)
	C.dword_5d4594_1045420 = 0
	C.dword_5d4594_1045428 = 0
	C.dword_5d4594_1045432 = 0
	return rows, (*uint32)(unsafe.Pointer(&C.dword_5d4594_1045432)), (*unsafe.Pointer)(unsafe.Pointer(&C.dword_5d4594_1045420)), func() {
		copy(rows, saved)
		C.dword_5d4594_1045420 = catalog
		C.dword_5d4594_1045428 = context
		C.dword_5d4594_1045432 = enabled
	}
}
func PortTestClientAudioSlot(id int32) unsafe.Pointer {
	return unsafe.Pointer(C.nox_xxx_draw_452270(C.int(id)))
}
func PortTestClientAudioDelay(p unsafe.Pointer) int32 { return int32(C.sub_4522A0(C.int(uintptr(p)))) }
func PortTestClientAudioSample(p unsafe.Pointer, key *byte) int32 {
	return int32(C.sub_486A10(C.int(uintptr(p)), unsafe.Pointer(key)))
}
func PortTestClientAudioRecord(f *binfile.MemFile, scratch []byte) int {
	return int(C.sub_452BD0(C.int(uintptr(f.C())), (*C.char)(unsafe.Pointer(&scratch[0]))))
}
