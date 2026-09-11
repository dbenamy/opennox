//go:build porttest

package legacy

/*
#include "GAME5.h"
*/
import "C"
import (
	"bytes"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"strings"
	"unsafe"
)

type PortTestCallbackLoaderResult struct {
	Kind, Name           string
	Empty                bool
	Success, WantSuccess bool
	Slot, WantSlot       int
	Intact, FieldCorrect bool
}

func PortTestCallbackLoaders() []PortTestCallbackLoaderResult {
	restore := portTestCallbackTablesEnvironment()
	defer restore()
	var out []PortTestCallbackLoaderResult
	for kind, base := range []uintptr{287096, 287280, 287192} {
		field := []uintptr{236, 228, 232}[kind]
		var names []string
		var funcs []uint32
		for i := uintptr(0); *memmap.PtrPtr(0x587000, base+8*i) != nil; i++ {
			names = append(names, alloc.GoString((*byte)(*memmap.PtrPtr(0x587000, base+8*i))))
			funcs = append(funcs, *memmap.PtrUint32(0x587000, base+8*i+4))
		}
		cases := append(append([]string(nil), names...), "NULL", "null", "NuLl", "unknown", "", strings.ToLower(names[0]))
		for _, empty := range []bool{false, true} {
			for ci, name := range cases {
				r := PortTestCallbackLoaderResult{Kind: []string{"strike", "die", "dead"}[kind], Name: name, Empty: empty, Slot: -1, WantSlot: -1, Intact: true}
				b, free := alloc.Make([]byte{}, int(unsafe.Sizeof(server.MonsterDef{}))+16)
				for i := range b {
					b[i] = 0xa5
				}
				def := unsafe.Pointer(&b[8])
				before := bytes.Clone(b)
				old := *memmap.PtrUint32(0x587000, base)
				if empty {
					*memmap.PtrUint32(0x587000, base) = 0
				}
				table := unsafe.Slice((*byte)(memmap.PtrOff(0x587000, 287096)), 640)
				beforeTable := bytes.Clone(table)
				cs, freeCS := alloc.CString(name)
				switch kind {
				case 0:
					r.Success = C.nox_xxx_monsterLoadStrikeFn_549040(C.int(uintptr(def)), (*C.char)(unsafe.Pointer(cs))) != 0
				case 1:
					r.Success = C.nox_xxx_monsterLoadDieFn_5490E0(C.int(uintptr(def)), (*C.char)(unsafe.Pointer(cs))) != 0
				case 2:
					r.Success = C.nox_xxx_monsterLoadDeadFn_549180(C.int(uintptr(def)), (*C.char)(unsafe.Pointer(cs))) != 0
				}
				freeCS()
				r.Intact = bytes.Equal(table, beforeTable)
				*memmap.PtrUint32(0x587000, base) = old
				want := uint32(0xa5a5a5a5)
				if strings.EqualFold(name, "NULL") {
					r.WantSuccess = true
					want = 0
				} else if !empty && ci < len(names) {
					r.WantSuccess = true
					r.WantSlot = ci
					want = funcs[ci]
				}
				actual := *(*uint32)(unsafe.Add(def, field))
				r.FieldCorrect = actual == want
				for i, p := range funcs {
					if p == actual {
						r.Slot = i
						break
					}
				}
				// Only the selected callback word is writable; all other bytes and guards match.
				copy(b[8+int(field):12+int(field)], before[8+int(field):12+int(field)])
				r.Intact = r.Intact && bytes.Equal(b, before)
				out = append(out, r)
				free()
			}
		}
	}
	return out
}
