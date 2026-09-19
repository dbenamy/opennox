//go:build porttest

package legacy

import (
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"unsafe"
)

type PortTestMOTDLineResult struct {
	Text           string
	Next           int
	Guard          bool
	InputUnchanged bool
}

func PortTestMOTDLine(input string) PortTestMOTDLineResult {
	src, freeSrc := alloc.CString(input)
	defer freeSrc()
	base, freeOut := alloc.Malloc(uintptr(len(input) + 65))
	defer freeOut()
	raw := unsafe.Slice((*byte)(base), len(input)+65)
	out := unsafe.Add(base, 32)
	for i := range raw {
		raw[i] = 0xa5
	}
	next := sessionMOTDLine(src, (*byte)(out))
	r := PortTestMOTDLineResult{Text: alloc.GoString((*byte)(out)), Next: -1, Guard: true, InputUnchanged: alloc.GoString(src) == input}
	if next != nil {
		r.Next = int(uintptr(unsafe.Pointer(next)) - uintptr(unsafe.Pointer(src)))
	}
	for _, b := range raw[:32] {
		if b != 0xa5 {
			r.Guard = false
		}
	}
	for _, b := range raw[32+len(r.Text)+1:] {
		if b != 0xa5 {
			r.Guard = false
		}
	}
	return r
}
