//go:build porttest

package legacy

/*
#include <stdint.h>
#include "GAME5_2.h"
int nox_xxx_netGetUnitByExtent_4ED020(int a1);
*/
import "C"

import (
	"bytes"
	"unsafe"

	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

type PortTestExtentSpec struct {
	Extent, NetCode, Flags uint32
}

type PortTestExtentSnapshot struct {
	DynamicResult    uint32
	FoundIndex       int
	ObjectsUnchanged bool
}

// PortTestNetworkExtent calls the C dynamic-code and extent-lookup routines
// against an isolated C-backed server list. With noServer, only the dynamic
// no-high-bit path may be used; Lookup is deliberately skipped.
func PortTestNetworkExtent(specs []PortTestExtentSpec, codes []uint32, noServer bool) []PortTestExtentSnapshot {
	var objects []server.Object
	if len(specs) != 0 {
		var free func()
		objects, free = alloc.Make([]server.Object{}, len(specs))
		defer free()
	}
	for i, spec := range specs {
		objects[i].Extent = spec.Extent
		objects[i].NetCode = spec.NetCode
		objects[i].ObjFlags = object.Flags(spec.Flags)
		if i+1 < len(objects) {
			objects[i].ObjNext = &objects[i+1]
		}
	}

	oldGet := GetServer
	var core *server.Server
	if !noServer {
		core = new(server.Server)
		if len(objects) != 0 {
			core.Objs.List = &objects[0]
		}
		GetServer = func() Server { return &portTestRandomServer{core: core} }
	} else {
		GetServer = func() Server { panic("unexpected server access for an unmarked code") }
	}
	defer func() { GetServer = oldGet }()

	findIndex := func(raw uint32) int {
		if raw == 0 {
			return -1
		}
		for i := range objects {
			if uint32(uintptr(unsafe.Pointer(&objects[i]))) == raw {
				return i
			}
		}
		return -2
	}
	var raw []byte
	if len(objects) > 0 {
		raw = unsafe.Slice((*byte)(unsafe.Pointer(&objects[0])), len(objects)*int(unsafe.Sizeof(server.Object{})))
	}
	before := append([]byte(nil), raw...)
	unchanged := func() bool { return bytes.Equal(raw, before) }

	out := make([]PortTestExtentSnapshot, 0, len(codes))
	for _, code := range codes {
		s := PortTestExtentSnapshot{
			DynamicResult:    uint32(C.nox_xxx_packetDynamicUnitCode_578B40(C.int(code))),
			FoundIndex:       -3,
			ObjectsUnchanged: unchanged(),
		}
		if !noServer {
			s.FoundIndex = findIndex(uint32(C.nox_xxx_netGetUnitByExtent_4ED020(C.int(code))))
			s.ObjectsUnchanged = s.ObjectsUnchanged && unchanged()
		}
		out = append(out, s)
	}
	return out
}
