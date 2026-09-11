//go:build porttest

package legacy

/*
#include "GAME4_1.h"
*/
import "C"

import (
	"bytes"
	"unsafe"

	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

type PortTestWaypointAppendSpec struct {
	// Mode 0 calls sub_51D300 with Kind. Mode 1 calls sub_51D2C0 and obtains
	// its signed char argument from BlobKind at 0x973F18+35972.
	Mode                  byte
	Count, Kind, BlobKind byte
	ExistingSlot          byte // 0..31, or any larger value for no custom slot
	ExistingTarget        int8 // -1 nil, 0 source, 1 target, 2 alternate
	ExistingKind          byte
	Target                int8 // -1 nil, 0 source, 1 target, 2 alternate
}

type PortTestWaypointAppendResult struct {
	Return                   int
	PointsCnt                byte
	Targets                  [32]int8 // normalized pointer identities
	Kinds                    [32]byte
	SourceOnlyExpectedWrites bool
	OutsideSourceUnchanged   bool // target, alternate, and both guards
	GuardWordsUnchanged      bool
	BlobAfterCall            byte
}

type PortTestWaypointAppendSnapshot struct {
	Results                      []PortTestWaypointAppendResult
	BlobBefore, BlobAfterRestore byte
}

func portTestWaypointAppendPtr(data []server.Waypoint, source *server.Waypoint, ind int8) *server.Waypoint {
	switch ind {
	case -1:
		return nil
	case 0:
		return source
	case 1:
		return &data[2]
	case 2:
		return &data[3]
	default:
		panic("invalid waypoint identity")
	}
}

func portTestWaypointAppendIndex(data []server.Waypoint, source, p *server.Waypoint) int8 {
	if p == nil {
		return -1
	}
	if p == source {
		return 0
	}
	if p == &data[2] {
		return 1
	}
	if p == &data[3] {
		return 2
	}
	return -2
}

// PortTestWaypointAppend invokes both original-C routes against C-owned
// waypoints. Each source starts with poisoned padding and real pointer values
// in every slot, so accidental reads/writes outside the documented fields are
// observable without placing Go pointers in C-visible storage.
func PortTestWaypointAppend(specs []PortTestWaypointAppendSpec) (snap PortTestWaypointAppendSnapshot) {
	data, free := alloc.Make([]server.Waypoint{}, 5) // guards 0/4, source 1, targets 2/3
	defer free()
	raw := unsafe.Slice((*byte)(unsafe.Pointer(&data[0])), len(data)*int(unsafe.Sizeof(server.Waypoint{})))
	wpSize := int(unsafe.Sizeof(server.Waypoint{}))
	srcOff := wpSize
	mode := memmap.PtrUint8(0x973F18, 35972)
	oldMode := *mode
	snap.BlobBefore = oldMode
	defer func() {
		*mode = oldMode
		snap.BlobAfterRestore = *mode
	}()

	before := make([]byte, len(raw))
	snap.Results = make([]PortTestWaypointAppendResult, 0, len(specs))
	for _, spec := range specs {
		for i := range raw {
			raw[i] = 0xa5
		}
		// Use different guard patterns to catch a mistaken range boundary.
		for i := range raw[:wpSize] {
			raw[i] = 0x3c
		}
		for i := range raw[4*wpSize:] {
			raw[4*wpSize+i] = 0xc3
		}
		source := &data[1]
		// Equal indices make pointer identity, rather than a convenient index
		// comparison, necessary for duplicate detection.
		source.Index, data[2].Index, data[3].Index = 42, 42, 42
		source.PointsCnt = spec.Count
		for i := range source.Points {
			source.Points[i].Waypoint = &data[3]
			source.Points[i].Ind = byte(0x40 + i)
		}
		if spec.ExistingSlot < byte(len(source.Points)) {
			source.Points[spec.ExistingSlot].Waypoint = portTestWaypointAppendPtr(data, source, spec.ExistingTarget)
			source.Points[spec.ExistingSlot].Ind = spec.ExistingKind
		}
		copy(before, raw)

		*mode = spec.BlobKind
		target := portTestWaypointAppendPtr(data, source, spec.Target)
		var ret C.int
		if spec.Mode == 0 {
			ret = C.sub_51D300(C.int(uintptr(unsafe.Pointer(source))), C.int(uintptr(unsafe.Pointer(target))), C.char(int8(spec.Kind)))
		} else if spec.Mode == 1 {
			ret = C.sub_51D2C0(C.int(uintptr(unsafe.Pointer(source))), C.int(uintptr(unsafe.Pointer(target))))
		} else {
			panic("invalid append mode")
		}

		res := PortTestWaypointAppendResult{Return: int(ret), PointsCnt: source.PointsCnt, BlobAfterCall: *mode}
		for i, p := range source.Points {
			res.Targets[i] = portTestWaypointAppendIndex(data, source, p.Waypoint)
			res.Kinds[i] = p.Ind
		}
		res.OutsideSourceUnchanged = bytes.Equal(before[:srcOff], raw[:srcOff]) && bytes.Equal(before[srcOff+wpSize:], raw[srcOff+wpSize:])
		res.GuardWordsUnchanged = bytes.Equal(before[:wpSize], raw[:wpSize]) && bytes.Equal(before[4*wpSize:], raw[4*wpSize:])
		res.SourceOnlyExpectedWrites = true
		allowed := func(off int) bool {
			if res.Return == 0 {
				return false
			}
			count := int(spec.Count)
			return off >= 92+8*count && off < 96+8*count || off == 96+8*count || off == 476
		}
		for i, v := range raw[srcOff : srcOff+wpSize] {
			if v != before[srcOff+i] && !allowed(i) {
				res.SourceOnlyExpectedWrites = false
				break
			}
		}
		snap.Results = append(snap.Results, res)
	}
	return snap
}
