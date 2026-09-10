//go:build porttest

package legacy

/*
#include "GAME5_2.h"
*/
import "C"
import (
	"bytes"
	"unsafe"

	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/legacy/common/alloc/memguard"
	"github.com/opennox/opennox/v1/server"
)

type PortTestProtectionObject struct {
	Nil, GuardObject, WithHealth                     bool
	NetCode                                          uint32
	TypeInd, HP                                      uint16
	TypePresent, InitPresent, NamePresent, GuardInit bool
	InitSize                                         uint32
	Init, Name                                       []byte
}

type PortTestProtectionObjectResult struct {
	State           PortTestRekeySnapshot
	Digest          uint32
	Results         []uint32
	ObjectUnchanged bool
}

func PortTestObjectProtection(initial [][2]uint32, key, sum, id uint32, spec PortTestProtectionObject, mode, repeat int) PortTestProtectionObjectResult {
	var obj *server.Object
	var watched [][]byte
	if spec.GuardObject {
		buf, free := memguard.New(int(unsafe.Sizeof(server.Object{})))
		defer free()
		obj = (*server.Object)(unsafe.Pointer(&buf[0]))
	} else if !spec.Nil {
		buf, free := alloc.Make([]server.Object{}, 1)
		defer free()
		obj = &buf[0]
		obj.NetCode, obj.TypeInd = spec.NetCode, spec.TypeInd
		watched = append(watched, unsafe.Slice((*byte)(unsafe.Pointer(obj)), int(unsafe.Sizeof(*obj))))
		if spec.WithHealth {
			h, free := alloc.Make([]server.HealthData{}, 1)
			defer free()
			h[0].Cur = spec.HP
			obj.HealthData = &h[0]
			watched = append(watched, unsafe.Slice((*byte)(unsafe.Pointer(obj.HealthData)), int(unsafe.Sizeof(h[0]))))
		}
		if spec.InitPresent {
			if spec.GuardInit {
				b, free := memguard.New(4)
				defer free()
				obj.InitData = unsafe.Pointer(&b[0])
			} else {
				b, free := alloc.Make(spec.Init, max(1, len(spec.Init)))
				defer free()
				obj.InitData = unsafe.Pointer(&b[0])
				watched = append(watched, b)
			}
		}
		if spec.NamePresent {
			b, free := alloc.Make(spec.Name, len(spec.Name)+1)
			defer free()
			obj.IDPtr = unsafe.Pointer(&b[0])
			watched = append(watched, b)
		}
	}
	before := make([][]byte, len(watched))
	for i, b := range watched {
		before[i] = bytes.Clone(b)
	}
	out := PortTestProtectionObjectResult{ObjectUnchanged: true}
	out.State = portTestRekeyOperation(initial, key, sum, 0x87654321, 0xffffffff, 0xffffffff, 1, 0x12345678, 123, func() uint32 {
		if spec.TypePresent {
			restore := GetServer().S().PortTestObjectInitSize(spec.TypeInd, spec.InitSize)
			defer restore()
		}
		if !spec.GuardObject {
			out.Digest = objectProtectionChecksum(obj)
		}
		var result uint32
		for i := 0; i < repeat; i++ {
			if mode == 0 {
				result = uint32(C.nox_xxx_protect_56FBF0(C.int(id), asObjectC(obj)))
			} else {
				result = uint32(C.nox_xxx_protect_56FC50(C.int(id), asObjectC(obj)))
			}
			out.Results = append(out.Results, result)
		}
		return result
	})
	for i, b := range watched {
		out.ObjectUnchanged = out.ObjectUnchanged && bytes.Equal(b, before[i])
	}
	return out
}
