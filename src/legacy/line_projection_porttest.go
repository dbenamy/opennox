//go:build porttest

package legacy

/*
#include <stdint.h>
#include "GAME5_2.h"
void sub_57C790(float4* a1, float2* a2, float2* a3, float a4);
*/
import "C"

import (
	"math"
	"unsafe"

	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

const portTestProjectionGuardWords = 4

// PortTestLineProjectionSpec uses word offsets into one C-owned buffer. Line,
// point, and output may overlap in any valid arrangement. LengthBits applies
// only to clamp and is passed through Float32frombits without arithmetic.
type PortTestLineProjectionSpec struct {
	Kind                                  string // clamp (57C790) or line (57C8A0)
	Words                                 []uint32
	LineOffset, PointOffset, OutputOffset int
	LengthBits                            uint32
}

type PortTestLineProjectionSnapshot struct {
	Return                 int
	Words                  []uint32
	OutsideOutputUnchanged bool
	GuardsOK               bool
}

func portTestProjectionWords(words []uint32) (all, data, guards []uint32, free func()) {
	all, free = alloc.Make([]uint32{}, len(words)+2*portTestProjectionGuardWords)
	for i := range all[:portTestProjectionGuardWords] {
		all[i] = 0x2db6a41f
	}
	for i := range all[len(all)-portTestProjectionGuardWords:] {
		all[len(all)-portTestProjectionGuardWords+i] = 0xc91857e3
	}
	data = all[portTestProjectionGuardWords : portTestProjectionGuardWords+len(words)]
	copy(data, words)
	guards = append([]uint32(nil), all[:portTestProjectionGuardWords]...)
	guards = append(guards, all[len(all)-portTestProjectionGuardWords:]...)
	return all, data, guards, free
}

// PortTestLineProjection calls the original C helpers and returns raw storage
// after the call. It deliberately supplies no numeric oracle.
func PortTestLineProjection(specs []PortTestLineProjectionSpec) []PortTestLineProjectionSnapshot {
	out := make([]PortTestLineProjectionSnapshot, 0, len(specs))
	for _, spec := range specs {
		if spec.LineOffset < 0 || spec.LineOffset+4 > len(spec.Words) || spec.PointOffset < 0 || spec.PointOffset+2 > len(spec.Words) || spec.OutputOffset < 0 || spec.OutputOffset+2 > len(spec.Words) {
			panic("invalid projection word offset")
		}
		all, data, guards, free := portTestProjectionWords(spec.Words)
		before := append([]uint32(nil), data...)
		line := (*C.float4)(unsafe.Pointer(&data[spec.LineOffset]))
		point := (*C.float2)(unsafe.Pointer(&data[spec.PointOffset]))
		result := 0
		switch spec.Kind {
		case "clamp":
			C.sub_57C790(line, point, (*C.float2)(unsafe.Pointer(&data[spec.OutputOffset])), C.float(math.Float32frombits(spec.LengthBits)))
		case "line":
			result = int(C.nox_xxx_mathPointOnTheLine_57C8A0(line, point, (*C.float2)(unsafe.Pointer(&data[spec.OutputOffset]))))
		default:
			panic("unknown projection kind")
		}
		outsideUnchanged := true
		for i := range data {
			if i != spec.OutputOffset && i != spec.OutputOffset+1 && data[i] != before[i] {
				outsideUnchanged = false
			}
		}
		pre := all[:portTestProjectionGuardWords]
		post := all[len(all)-portTestProjectionGuardWords:]
		guardsOK := true
		for i := range pre {
			guardsOK = guardsOK && pre[i] == guards[i] && post[i] == guards[portTestProjectionGuardWords+i]
		}
		out = append(out, PortTestLineProjectionSnapshot{Return: result, Words: append([]uint32(nil), data...), OutsideOutputUnchanged: outsideUnchanged, GuardsOK: guardsOK})
		free()
	}
	return out
}
