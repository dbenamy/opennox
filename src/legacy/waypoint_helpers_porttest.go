//go:build porttest

package legacy

/*
#include <stdlib.h>
#include "GAME5_2.h"
*/
import "C"

import (
	"bytes"
	"unsafe"

	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

type PortTestWaypointMask struct {
	Nil          bool
	Flags        uint32
	Flags2, Mask byte
}
type PortTestWaypointMaskResult struct {
	Mask, EnabledMask int
	Unchanged         bool
}

func PortTestWaypointMasks(specs []PortTestWaypointMask) []PortTestWaypointMaskResult {
	wp, free := alloc.New(server.Waypoint{})
	defer free()
	raw := unsafe.Slice((*byte)(unsafe.Pointer(wp)), int(unsafe.Sizeof(*wp)))
	out := make([]PortTestWaypointMaskResult, 0, len(specs))
	for _, spec := range specs {
		wp.Flags, wp.Flags2 = spec.Flags, spec.Flags2
		before := append([]byte(nil), raw...)
		arg := C.int(uintptr(unsafe.Pointer(wp)))
		if spec.Nil {
			arg = 0
		}
		mask := -1
		if !spec.Nil {
			mask = 0
			if wp.HasFlag2Mask(spec.Mask) {
				mask = 1
			}
		}
		enabled := bool2int(waypointEnabledMask(waypointFromRaw(arg), spec.Mask))
		out = append(out, PortTestWaypointMaskResult{Mask: mask, EnabledMask: enabled, Unchanged: bytes.Equal(raw, before)})
	}
	return out
}

type PortTestWaypointLink struct{ Source, Next int }
type PortTestWaypointLinkResult struct {
	First, Second int
	Unchanged     bool
}

func PortTestWaypointLinks(specs []PortTestWaypointLink) []PortTestWaypointLinkResult {
	data, free := alloc.Make([]server.Waypoint{}, 3)
	defer free()
	raw := unsafe.Slice((*byte)(unsafe.Pointer(&data[0])), len(data)*int(unsafe.Sizeof(server.Waypoint{})))
	index := func(v uint32) int {
		if v == 0 {
			return -1
		}
		for i := range data {
			if uint32(uintptr(unsafe.Pointer(&data[i]))) == v {
				return i
			}
		}
		return -2
	}
	out := make([]PortTestWaypointLinkResult, 0, len(specs))
	for _, spec := range specs {
		var next *server.Waypoint
		if spec.Next >= 0 {
			next = &data[spec.Next]
		}
		var arg C.int
		if spec.Source >= 0 {
			wp := &data[spec.Source]
			wp.WpNext = next
			wp.WpPrev = &data[(spec.Source+1)%len(data)]
			arg = C.int(uintptr(unsafe.Pointer(wp)))
		}
		before := append([]byte(nil), raw...)
		a := index(uint32(C.nox_xxx_waypointNext_579870(arg)))
		b := index(uint32(C.sub_5798A0(arg)))
		out = append(out, PortTestWaypointLinkResult{First: a, Second: b, Unchanged: bytes.Equal(raw, before)})
	}
	return out
}

func PortTestWaypointAllocations(count int) [][]byte {
	out := make([][]byte, 0, count)
	for i := 0; i < count; i++ {
		p := C.sub_579E70()
		if p == nil {
			out = append(out, nil)
			continue
		}
		raw := unsafe.Slice((*byte)(unsafe.Pointer(p)), int(unsafe.Sizeof(server.Waypoint{})))
		out = append(out, append([]byte(nil), raw...))
		// Poison before freeing, to exercise zeroing when the allocator reuses storage.
		for j := range raw {
			raw[j] = 0xa5
		}
		C.free(unsafe.Pointer(p))
	}
	return out
}
