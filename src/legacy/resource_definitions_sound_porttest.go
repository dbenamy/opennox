//go:build porttest

package legacy

/*
#include <stdlib.h>
#include "GAME1.h"
#include "GAME1_1.h"
extern uint32_t dword_5d4594_588120;
*/
import "C"
import (
	"bytes"
	"github.com/opennox/opennox/v1/internal/binfile"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"unsafe"
)

func PortTestResourceSoundOwner() func() {
	old := C.dword_5d4594_588120
	C.dword_5d4594_588120 = 0
	files.Lock()
	oldFiles := files.byHandle
	files.byHandle = make(map[unsafe.Pointer]*binfile.File)
	files.Unlock()
	return func() {
		PortTestResourceSoundClear()
		C.dword_5d4594_588120 = old
		files.Lock()
		defer files.Unlock()
		for _, f := range files.byHandle {
			_ = f.Close()
		}
		files.byHandle = oldFiles
	}
}
func PortTestResourceSoundClear() {
	for p := unsafe.Pointer(uintptr(C.dword_5d4594_588120)); p != nil; {
		next := *(*unsafe.Pointer)(unsafe.Add(p, 76))
		C.free(*(*unsafe.Pointer)(p))
		C.free(p)
		p = next
	}
	C.dword_5d4594_588120 = 0
}

type PortTestResourceSoundRow struct {
	Name  string
	Words [20]uint32
}

func PortTestResourceSounds() []PortTestResourceSoundRow {
	var out []PortTestResourceSoundRow
	for p := unsafe.Pointer(uintptr(C.dword_5d4594_588120)); p != nil; p = *(*unsafe.Pointer)(unsafe.Add(p, 76)) {
		if len(out) > 100 {
			panic("sound list cycle")
		}
		row := PortTestResourceSoundRow{Name: alloc.GoString(*(**byte)(p))}
		copy(row.Words[:], unsafe.Slice((*uint32)(unsafe.Add(p, 4)), 20))
		if row.Words[18] != 0 {
			row.Words[18] = uint32(len(out) + 2)
		}
		out = append(out, row)
	}
	return out
}
func PortTestResourceSoundLookup(name string) int {
	p := C.nox_xxx_getDefaultSoundSet_424350(internCStr(name))
	if p == nil {
		return 0
	}
	i := 1
	for it := unsafe.Pointer(uintptr(C.dword_5d4594_588120)); it != nil; it = *(*unsafe.Pointer)(unsafe.Add(it, 76)) {
		if it == unsafe.Pointer(p) {
			return i
		}
		i++
	}
	panic("sound lookup outside owner")
}
func PortTestResourceTokens(f *binfile.Binfile) []PortTestMonsterToken {
	h := NewFileHandle(f.File)
	defer nox_fs_close(h)
	b, free := alloc.Make([]byte{}, 272)
	defer free()
	var out []PortTestMonsterToken
	for i := 0; i < 1024; i++ {
		for j := range b {
			b[j] = 0xa5
		}
		rv := int(C.nox_xxx_parseString_409470(h, (*C.uint8_t)(unsafe.Pointer(&b[8]))))
		r := PortTestMonsterToken{Return: rv, Buffer: bytes.Clone(b[8:264]), Intact: true}
		for _, v := range append(bytes.Clone(b[:8]), b[264:]...) {
			r.Intact = r.Intact && v == 0xa5
		}
		out = append(out, r)
		if rv == 0 {
			return out
		}
	}
	panic("token loop failed to end")
}
