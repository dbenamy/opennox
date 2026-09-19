//go:build porttest

package blobdata

import (
	"encoding/binary"
	"math"
)

// PortTestUnitActionNames decodes the shipped pointers before runtime relocation.
func PortTestUnitActionNames() []string {
	out := make([]string, 72)
	for i := range out {
		ptr := binary.LittleEndian.Uint32(data587000[261768+4*i:])
		if i == 39 {
			if ptr != 0x5D4594+2488508 {
				panic("unexpected empty action pointer")
			}
			continue // This entry points to zero-initialized writable storage.
		}
		off := int(ptr - 0x587000)
		end := off
		for data587000[end] != 0 {
			end++
		}
		out[i] = string(data587000[off:end])
	}
	return out
}

func PortTestUnitExperienceCoefficient() float32 {
	return math.Float32frombits(binary.LittleEndian.Uint32(data587000[206148:]))
}
