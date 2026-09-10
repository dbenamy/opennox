//go:build porttest

package legacy

/*
#include "GAME5_2.h"
*/
import "C"

import (
	"bytes"
	"unsafe"

	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

type PortTestClientCode struct {
	Nil         bool
	Code, Class uint32
}
type PortTestClientCodeResult struct {
	Code      uint32
	Unchanged bool
}

func PortTestClientCodes(specs []PortTestClientCode) []PortTestClientCodeResult {
	data, free := alloc.Make([]client.Drawable{}, 1)
	defer free()
	dr := &data[0]
	out := make([]PortTestClientCodeResult, 0, len(specs))
	for _, spec := range specs {
		dr.NetCode32 = spec.Code
		dr.ObjClass = object.Class(spec.Class)
		raw := unsafe.Slice((*byte)(unsafe.Pointer(dr)), int(unsafe.Sizeof(*dr)))
		before := append([]byte(nil), raw...)
		var arg C.int
		if !spec.Nil {
			arg = C.int(uintptr(unsafe.Pointer(dr)))
		}
		code := uint32(C.nox_xxx_netGetUnitCodeCli_578B00(arg))
		out = append(out, PortTestClientCodeResult{Code: code, Unchanged: bytes.Equal(raw, before)})
	}
	return out
}

func PortTestNetworkBits(value uint32) (clear, flag uint32) {
	return uint32(C.nox_xxx_netClearHighBit_578B30(C.short(value))), uint32(C.nox_xxx_netTestHighBit_578B70(C.uint(value)))
}
