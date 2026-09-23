//go:build porttest

package legacy

/*
#include <stdint.h>
#include "GAME2.h"
#include "GAME2_2.h"
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
	catalog, context, enabled := dword_5d4594_1045420, dword_5d4594_1045428, dword_5d4594_1045432
	clear(rows)
	dword_5d4594_1045420 = 0
	dword_5d4594_1045428 = 0
	dword_5d4594_1045432 = 0
	return rows, (*uint32)(unsafe.Pointer(&dword_5d4594_1045432)), (*unsafe.Pointer)(unsafe.Pointer(&dword_5d4594_1045420)), func() {
		copy(rows, saved)
		dword_5d4594_1045420 = catalog
		dword_5d4594_1045428 = context
		dword_5d4594_1045432 = enabled
	}
}
func PortTestClientAudioSlot(id int32) unsafe.Pointer {
	return audioAssetSlot(id)
}
func PortTestClientAudioDelay(p unsafe.Pointer) int32 { return audioAssetDelay(p) }
func PortTestClientAudioSample(p unsafe.Pointer, key *byte) int32 {
	return audioAssetSample(p, key)
}
func PortTestClientAudioRecord(f *binfile.MemFile, scratch []byte) int {
	return audioAssetRecord(f, scratch)
}
