package legacy

/*
#include <stdint.h>
*/
import "C"

import (
	"time"

	"github.com/opennox/opennox/v1/common/memmap"
)

func Sub_56F1C0() { initializeProtection() }

func initializeProtection() uint32 {
	protectionRandom.Seed(uint32(time.Now().Unix()))
	dword_5d4594_2516352 = 0
	dword_5d4594_2516348 = C.uint32_t(GetServer().S().Frame())
	dword_5d4594_2516344 = 0
	*memmap.PtrUint16(0x587000, 311204) = 0
	dword_5d4594_2516356 = 657757279
	dword_5d4594_2516348 ^= C.uint32_t(protectionRandom.Draw())
	dword_5d4594_2516328 = ^dword_5d4594_2516348
	*memmap.PtrUint32(0x5D4594, 2516340) = uint32(nox_xxx_protectionCreateInt_56F400(0))
	sub_56F250()
	result := uint32(nox_xxx_protectionCreateInt_56F400(1))
	*memmap.PtrUint32(0x5D4594, 2516332) = result
	return result
}
