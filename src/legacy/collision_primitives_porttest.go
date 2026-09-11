//go:build porttest

package legacy

/*
#include <stdint.h>
#include "GAME5_2.h"
*/
import "C"

import (
	"unsafe"

	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

const portTestCollisionGuardWords = 4

// PortTestReflectionSpec addresses both two-float arguments by word offset in
// Words. Offsets may be equal or partially overlap, so the fixture exposes the
// exact C load/store behavior without converting raw float bits through Go.
type PortTestReflectionSpec struct {
	Words                        []uint32
	NormalOffset, VelocityOffset int
}

type PortTestReflectionSnapshot struct {
	ReturnRaw         uint32
	ReturnPointerSame bool
	Words             []uint32
	GuardsOK          bool
}

func portTestCollisionWords(words []uint32) (all, data, guards []uint32, free func()) {
	all, free = alloc.Make([]uint32{}, len(words)+2*portTestCollisionGuardWords)
	for i := range all[:portTestCollisionGuardWords] {
		all[i] = 0x91c3d57b
	}
	for i := range all[len(all)-portTestCollisionGuardWords:] {
		all[len(all)-portTestCollisionGuardWords+i] = 0x6ea4b809
	}
	data = all[portTestCollisionGuardWords : portTestCollisionGuardWords+len(words)]
	copy(data, words)
	guards = append([]uint32(nil), all[:portTestCollisionGuardWords]...)
	guards = append(guards, all[len(all)-portTestCollisionGuardWords:]...)
	return all, data, guards, free
}

// PortTestCollisionReflect invokes the original C helper on C-owned words.
func PortTestCollisionReflect(specs []PortTestReflectionSpec) []PortTestReflectionSnapshot {
	out := make([]PortTestReflectionSnapshot, 0, len(specs))
	for _, spec := range specs {
		if spec.NormalOffset < 0 || spec.NormalOffset+2 > len(spec.Words) || spec.VelocityOffset < 0 || spec.VelocityOffset+2 > len(spec.Words) {
			panic("invalid reflection word offset")
		}
		all, data, guard, free := portTestCollisionWords(spec.Words)
		ret := C.nox_xxx_collideReflect_57B810(
			(*C.float)(unsafe.Pointer(&data[spec.NormalOffset])),
			C.int(uintptr(unsafe.Pointer(&data[spec.VelocityOffset]))),
		)
		pre := all[:portTestCollisionGuardWords]
		post := all[len(all)-portTestCollisionGuardWords:]
		guardsOK := true
		for i := range pre {
			guardsOK = guardsOK && pre[i] == guard[i] && post[i] == guard[portTestCollisionGuardWords+i]
		}
		ptr := uint32(uintptr(unsafe.Pointer(&data[spec.VelocityOffset])))
		out = append(out, PortTestReflectionSnapshot{ReturnRaw: uint32(ret), ReturnPointerSame: uint32(ret) == ptr, Words: append([]uint32(nil), data...), GuardsOK: guardsOK})
		free()
	}
	return out
}

// PortTestContainmentSpec addresses the two-word position, eleven-word shape,
// and two-word point in one raw C-owned buffer. The helper has no output
// pointer; Words and Unchanged make unintended writes observable.
type PortTestContainmentSpec struct {
	Words                                    []uint32
	PositionOffset, ShapeOffset, PointOffset int
}

type PortTestContainmentSnapshot struct {
	Return    int
	Words     []uint32
	Unchanged bool
	GuardsOK  bool
}

// PortTestCollisionContainment invokes original C 57B850 without a Go float
// conversion or an inferred arithmetic oracle.
func PortTestCollisionContainment(specs []PortTestContainmentSpec) []PortTestContainmentSnapshot {
	out := make([]PortTestContainmentSnapshot, 0, len(specs))
	for _, spec := range specs {
		if spec.PositionOffset < 0 || spec.PositionOffset+2 > len(spec.Words) || spec.ShapeOffset < 0 || spec.ShapeOffset+11 > len(spec.Words) || spec.PointOffset < 0 || spec.PointOffset+2 > len(spec.Words) {
			panic("invalid containment word offset")
		}
		all, data, guard, free := portTestCollisionWords(spec.Words)
		before := append([]uint32(nil), data...)
		ret := C.nox_xxx_map_57B850(
			(*C.float2)(unsafe.Pointer(&data[spec.PositionOffset])),
			(*C.float)(unsafe.Pointer(&data[spec.ShapeOffset])),
			(*C.float2)(unsafe.Pointer(&data[spec.PointOffset])),
		)
		pre := all[:portTestCollisionGuardWords]
		post := all[len(all)-portTestCollisionGuardWords:]
		guardsOK := true
		for i := range pre {
			guardsOK = guardsOK && pre[i] == guard[i] && post[i] == guard[portTestCollisionGuardWords+i]
		}
		out = append(out, PortTestContainmentSnapshot{Return: int(ret), Words: append([]uint32(nil), data...), Unchanged: equalUint32(data, before), GuardsOK: guardsOK})
		free()
	}
	return out
}

func equalUint32(a, b []uint32) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
