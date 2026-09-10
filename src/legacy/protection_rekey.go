package legacy

/*
#include <stdint.h>
extern uint32_t dword_5d4594_2516348;
extern uint32_t dword_5d4594_2516328;
int nox_xxx_protect_56F240(void);
*/
import "C"
import (
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/internal/protection"
)

//export nox_xxx_protectData_56F5C0
func nox_xxx_protectData_56F5C0() C.int {
	frame := GetServer().S().Frame()
	oldKey := uint32(C.dword_5d4594_2516348)
	newKey := uint32(C.nox_xxx_protect_56F240()) ^ frame
	C.dword_5d4594_2516328 = C.uint32_t(^newKey)
	count := int(*memmap.PtrUint16(0x587000, 311204))
	head := protectionHead()
	for i := 0; i < count/4; i++ {
		a := GetServer().S().Rand.Logic.IntClamp(0, count/2)
		b := GetServer().S().Rand.Logic.IntClamp(count/2+1, count-1)
		if a != b {
			swapProtectionRecords(protection.At(head, int32(a)), protection.At(head, int32(b)))
		}
	}
	C.dword_5d4594_2516348 = 0
	C.dword_5d4594_2516328 = C.uint32_t(protection.Rekey(head, oldKey, newKey))
	*memmap.PtrUint32(0x5D4594, 2516364)++
	C.dword_5d4594_2516348 = C.uint32_t(newKey)
	return C.int(newKey)
}
