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
)

type PortTestAliasCall struct {
	Kind              string
	Key1, Key2, Frame uint32
	Slot              int
	Wrapper           bool
}

type PortTestAliasResult struct {
	Return                int
	PointerSame, GuardsOK bool
	Data                  []byte
}

func PortTestAliasTable(initial []byte, calls []PortTestAliasCall) []PortTestAliasResult {
	const size = 255 * 8
	raw, free := alloc.Make([]byte{}, size+32)
	defer free()
	for i := range raw {
		raw[i] = 0xa7
	}
	data := raw[16 : 16+size]
	copy(data, initial)
	guard := append([]byte(nil), raw[:16]...)
	p := unsafe.Pointer(&data[0])
	out := make([]PortTestAliasResult, 0, len(calls))
	for _, c := range calls {
		r := PortTestAliasResult{PointerSame: true}
		switch c.Kind {
		case "select":
			r.Return = int(C.nox_xxx_cliGenerateAlias_57B9A0(C.int(uintptr(p)), C.int(c.Key1), C.int(c.Key2), C.uint(c.Frame)))
		case "write":
			ptr := unsafe.Add(p, c.Slot*8)
			ret := C.sub_57BA10(C.int(uintptr(ptr)), C.short(c.Key1), C.short(c.Key2), C.int(c.Frame))
			r.PointerSame = uint32(ret) == uint32(uintptr(ptr))
		case "reset":
			if c.Wrapper {
				Sub_57B920(p)
			} else {
				r.Return = int(C.sub_57B920(p))
			}
		default:
			panic("unknown alias call")
		}
		r.Data = append([]byte(nil), data...)
		r.GuardsOK = bytes.Equal(raw[:16], guard) && bytes.Equal(raw[16+size:], guard)
		out = append(out, r)
	}
	return out
}
