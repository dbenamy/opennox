package legacy

/*
#include <stdint.h>
void sub_56FF00(int seed);
int nox_xxx_protect_56F240(void);
extern uint32_t dword_5d4594_2516344;
extern uint32_t dword_5d4594_2516352;
extern uint32_t dword_5d4594_2516348;
extern uint32_t dword_5d4594_2516328;
extern uint32_t dword_5d4594_2516356;
*/
import "C"

import (
	"time"

	"github.com/opennox/opennox/v1/common/memmap"
)

func Sub_56F1C0() { initializeProtection() }

func initializeProtection() uint32 {
	C.sub_56FF00(C.int(time.Now().Unix()))
	C.dword_5d4594_2516352 = 0
	C.dword_5d4594_2516348 = C.uint32_t(GetServer().S().Frame())
	C.dword_5d4594_2516344 = 0
	*memmap.PtrUint16(0x587000, 311204) = 0
	C.dword_5d4594_2516356 = 657757279
	C.dword_5d4594_2516348 ^= C.uint32_t(C.nox_xxx_protect_56F240())
	C.dword_5d4594_2516328 = ^C.dword_5d4594_2516348
	*memmap.PtrUint32(0x5D4594, 2516340) = uint32(nox_xxx_protectionCreateInt_56F400(0))
	sub_56F250()
	result := uint32(nox_xxx_protectionCreateInt_56F400(1))
	*memmap.PtrUint32(0x5D4594, 2516332) = result
	return result
}
