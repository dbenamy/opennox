//go:build porttest

package legacy

/*
#include <stdlib.h>
#include "GAME1.h"
extern void* dword_5d4594_251560;
*/
import "C"

import "unsafe"

type PortTestSecretOperation struct{ Op, ID int }
type PortTestSecretResult struct {
	Return    int
	List      []int
	PayloadOK bool
}

// Nodes and wall records have C allocation ownership, as in map loading.
// Removed pointer values are compared as identities, never dereferenced.
func PortTestSecretWalls(ids []uint16, ops []PortTestSecretOperation) (out []PortTestSecretResult) {
	old := C.dword_5d4594_251560
	C.dword_5d4594_251560 = nil
	nodes := make([]unsafe.Pointer, len(ids))
	walls := make([]unsafe.Pointer, len(ids))
	live := make([]bool, len(ids))
	linked := make([]bool, len(ids))
	for i, id := range ids {
		nodes[i] = C.calloc(1, 32)
		walls[i] = C.calloc(1, 64)
		live[i] = true
		if nodes[i] == nil || walls[i] == nil {
			panic("secret fixture allocation")
		}
		w := unsafe.Slice((*uint32)(nodes[i]), 8)
		w[1] = uint32(i + 1)
		w[2] = 0x12345678
		w[3] = uint32(uintptr(walls[i]))
		w[7] = 0x87654321
		*(*uint16)(unsafe.Add(walls[i], 10)) = id
	}
	defer func() {
		for i, p := range nodes {
			if live[i] {
				C.free(p)
			}
			C.free(walls[i])
		}
		C.dword_5d4594_251560 = old
	}()
	identify := func(p unsafe.Pointer, list []unsafe.Pointer) int {
		if p == nil {
			return -1
		}
		for i, q := range list {
			if p == q {
				return i
			}
		}
		return -2
	}
	for _, op := range ops {
		r := PortTestSecretResult{Return: -1, PayloadOK: true}
		switch op.Op {
		case 0:
			if !live[op.ID] || linked[op.ID] {
				panic("invalid secret fixture insertion")
			}
			r.Return = identify(unsafe.Pointer(C.nox_xxx_wallSecretBlock_410760((*C.uint32_t)(nodes[op.ID]))), nodes)
			linked[op.ID] = true
		case 1:
			r.Return = identify(unsafe.Pointer(C.sub_4107A0(nodes[op.ID])), nodes)
			if r.Return >= 0 {
				live[r.Return] = false
				linked[r.Return] = false
			}
		case 2:
			r.Return = identify(unsafe.Pointer(uintptr(C.sub_410550(C.short(op.ID)))), walls)
		case 3:
			r.Return = identify(unsafe.Pointer(C.sub_410730()), nodes)
			for i := range linked {
				if linked[i] {
					live[i] = false
					linked[i] = false
				}
			}
		case 4:
			r.Return = identify(unsafe.Pointer(uintptr(C.nox_xxx_wallSecretNext_410790(nil))), nodes)
		}
		p := C.nox_xxx_wallSecretGetFirstWall_410780()
		for p != nil {
			i := identify(p, nodes)
			if i < 0 || !live[i] || len(r.List) >= len(nodes) {
				panic("invalid secret list")
			}
			r.List = append(r.List, i)
			w := unsafe.Slice((*uint32)(p), 8)
			r.PayloadOK = r.PayloadOK && w[1] == uint32(i+1) && w[2] == 0x12345678 && w[3] == uint32(uintptr(walls[i])) && w[7] == 0x87654321
			p = unsafe.Pointer(uintptr(C.nox_xxx_wallSecretNext_410790((*C.int)(p))))
		}
		out = append(out, r)
	}
	return out
}
