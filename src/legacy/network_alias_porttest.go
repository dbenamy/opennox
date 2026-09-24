//go:build porttest

package legacy

import (
	"bytes"
	"github.com/opennox/opennox/v1/server"
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
			r.Return = int(portTestInvoke_nox_xxx_cliGenerateAlias_57B9A0(int32(uintptr(p)), int32(c.Key1), int32(c.Key2), uint32(c.Frame)))
		case "write":
			ptr := unsafe.Add(p, c.Slot*8)
			ret := portTestInvoke_sub_57BA10(int32(uintptr(ptr)), int16(c.Key1), int16(c.Key2), int32(c.Frame))
			r.PointerSame = uint32(ret) == uint32(uintptr(ptr))
		case "reset":
			if c.Wrapper {
				Sub_57B920(p)
			} else {
				r.Return = int(portTestInvoke_sub_57B920(p))
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

// Fixture-native copies preserve the original wrapper ABI conversions.
func portTestInvoke_nox_xxx_cliGenerateAlias_57B9A0(ptr, key1, key2 int32, frame uint32) int8 {
	table := (*[255]server.PlayerNetData)(unsafe.Pointer(uintptr(uint32(ptr))))
	return int8(selectNetworkAlias(table, int32(key1), int32(key2), uint32(frame)))
}

func portTestInvoke_sub_57B920(ptr unsafe.Pointer) int32 {
	resetNetworkAliases((*[255]server.PlayerNetData)(ptr))
	return 0
}

func portTestInvoke_sub_57BA10(ptr int32, key1, key2 int16, frame int32) int32 {
	rec := (*server.PlayerNetData)(unsafe.Pointer(uintptr(uint32(ptr))))
	*rec = server.PlayerNetData{Field0: uint16(key1), Field2: uint16(key2), Frame4: uint32(frame)}
	return ptr
}
