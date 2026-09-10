//go:build porttest

package legacy

/*
#include <stdint.h>
#include <stdlib.h>
extern uint32_t dword_5d4594_2516344;
extern uint32_t dword_5d4594_2516352;
extern uint32_t dword_5d4594_2516348;
extern uint32_t dword_5d4594_2516328;
*/
import "C"
import (
	"math"
	"unsafe"

	"github.com/opennox/libs/prand"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/internal/protection"
	"github.com/opennox/opennox/v1/server"
)

type portTestRandomServer struct {
	Server
	core *server.Server
}

func (s *portTestRandomServer) S() *server.Server { return s.core }

type PortTestInsertSnapshot struct {
	Result                 int
	Values                 [][2]uint32
	Sum, Key               uint32
	Count                  uint16
	LogicIndex, OtherIndex int
	LinksValid             bool
}

func PortTestInsert(initial, values [][2]uint32, key, sum uint32, seed int) []PortTestInsertSnapshot {
	oldGet := GetServer
	core := new(server.Server)
	core.Rand.Logic, core.Rand.Other = prand.New(seed), prand.New(seed+1)
	GetServer = func() Server { return &portTestRandomServer{core: core} }
	head, tail, oldKey, oldSum := C.dword_5d4594_2516344, C.dword_5d4594_2516352, C.dword_5d4594_2516348, C.dword_5d4594_2516328
	count := memmap.PtrUint16(0x587000, 311204)
	oldCount := *count
	defer func() {
		Sub_56F3B0()
		C.dword_5d4594_2516344, C.dword_5d4594_2516352, C.dword_5d4594_2516348, C.dword_5d4594_2516328 = head, tail, oldKey, oldSum
		*count = oldCount
		GetServer = oldGet
	}()
	C.dword_5d4594_2516344, C.dword_5d4594_2516352 = 0, 0
	C.dword_5d4594_2516348, C.dword_5d4594_2516328 = C.uint(key), C.uint(sum)
	*count = uint16(len(initial))
	var first, last *protection.Record
	for _, v := range initial {
		r := (*protection.Record)(C.calloc(1, C.size_t(unsafe.Sizeof(protection.Record{}))))
		if r == nil {
			panic("fixture allocation failed")
		}
		*r = protection.Record{ID: v[0] ^ key, Value: v[1] ^ key, Prev: last}
		if last == nil {
			first = r
		} else {
			last.Next = r
		}
		last = r
	}
	C.dword_5d4594_2516344 = C.uint(uintptr(unsafe.Pointer(first)))
	C.dword_5d4594_2516352 = C.uint(uintptr(unsafe.Pointer(last)))
	var out []PortTestInsertSnapshot
	for i, v := range values {
		var result int
		if i%2 == 0 {
			result = Nox_xxx_protectionCreateStructForInt_56F280(int(v[0]), int(v[1]))
		} else {
			result = Nox_xxx_protectionCreateStructForFloat_56F480(int(v[0]), math.Float32frombits(v[1]))
		}
		s := PortTestInsertSnapshot{Result: result, Sum: uint32(C.dword_5d4594_2516328), Key: uint32(C.dword_5d4594_2516348), Count: *count, LogicIndex: core.Rand.Logic.Index(), OtherIndex: core.Rand.Other.Index(), LinksValid: true}
		var prev uint32
		for p := protectionHead(); p != nil; p = p.Next {
			if len(s.Values) > len(initial)+i {
				s.LinksValid = false
				break
			}
			s.LinksValid = s.LinksValid && uint32(uintptr(unsafe.Pointer(p.Prev))) == prev
			s.Values = append(s.Values, [2]uint32{p.ID ^ key, p.Value ^ key})
			prev = uint32(uintptr(unsafe.Pointer(p)))
		}
		s.LinksValid = s.LinksValid && prev == uint32(C.dword_5d4594_2516352)
		out = append(out, s)
	}
	return out
}
