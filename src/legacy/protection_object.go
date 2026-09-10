package legacy

/*
#include "defs.h"
#include <stdint.h>
extern uint32_t dword_5d4594_2516348;
extern uint32_t dword_5d4594_2516328;
*/
import "C"
import (
	"unsafe"

	"github.com/opennox/opennox/v1/internal/protection"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func objectProtectionChecksum(obj *server.Object) uint32 {
	if obj == nil {
		return 0
	}
	var sum uint32
	sum ^= obj.NetCode
	sum ^= uint32(obj.TypeInd)
	if obj.HealthData != nil {
		sum ^= uint32(obj.HealthData.Cur)
	}
	if obj.InitData != nil {
		typ := GetServer().S().Types.ByInd(int(obj.TypeInd))
		if typ != nil {
			if size := int32(typ.InitDataSize); size > 0 {
				sum ^= protection.Checksum(unsafe.Slice((*byte)(obj.InitData), int(size)))
			}
		}
	}
	if obj.IDPtr != nil {
		n := alloc.StrLen((*byte)(obj.IDPtr))
		if n != 0 {
			sum ^= protection.Checksum(unsafe.Slice((*byte)(obj.IDPtr), n))
		}
	}
	return sum
}

func toggleProtectionObject(id int32, obj *server.Object) uint32 {
	if id < 657757279 {
		return uint32(id)
	}
	key := uint32(C.dword_5d4594_2516348)
	r := protection.Find(protectionHead(), key, uint32(id))
	if r == nil {
		return uint32(id)
	}
	digest := objectProtectionChecksum(obj)
	r.Value ^= digest
	C.dword_5d4594_2516328 ^= C.uint32_t(digest)
	return uint32(C.dword_5d4594_2516328)
}

//export nox_xxx_protect_56FBF0
func nox_xxx_protect_56FBF0(id C.int, obj *nox_object_t) C.int {
	return C.int(toggleProtectionObject(int32(id), asObjectS(obj)))
}

//export nox_xxx_protect_56FC50
func nox_xxx_protect_56FC50(id C.int, obj *nox_object_t) C.int {
	return C.int(toggleProtectionObject(int32(id), asObjectS(obj)))
}
