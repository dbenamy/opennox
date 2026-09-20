//go:build porttest

package legacy

/*
#include <stdlib.h>
#include "GAME1_1.h"
*/
import "C"
import (
	"bytes"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"unsafe"
)

type portTestListNode struct {
	Next, Prev   unsafe.Pointer
	Key, Payload uint32
}
type PortTestListNodeState struct {
	Next, Prev   int
	Key, Payload uint32
}
type PortTestLists struct {
	base  unsafe.Pointer
	nodes []portTestListNode
}

func PortTestListsOpen(n int) *PortTestLists {
	p := C.calloc(C.size_t(n), 16)
	if p == nil {
		panic("list fixture allocation")
	}
	f := &PortTestLists{base: p, nodes: unsafe.Slice((*portTestListNode)(p), n)}
	for i := range f.nodes {
		f.nodes[i].Payload = uint32(i)*0x12345 + 7
	}
	return f
}
func (f *PortTestLists) Close() { C.free(f.base) }
func (f *PortTestLists) ptr(i int) unsafe.Pointer {
	if i < 0 {
		return nil
	}
	return unsafe.Pointer(&f.nodes[i])
}
func (f *PortTestLists) id(p unsafe.Pointer) int {
	if p == nil {
		return -1
	}
	d := uintptr(p) - uintptr(f.base)
	if d%16 != 0 || d/16 >= uintptr(len(f.nodes)) {
		panic("list pointer outside fixture")
	}
	return int(d / 16)
}
func (f *PortTestLists) Init(i int, head bool, key uint32) {
	if head {
		listClear((*legacyListNode)(f.ptr(i)))
	} else {
		if unsafe.Pointer(listInit((*legacyListNode)(f.ptr(i)))) != f.ptr(i) {
			panic("node init return identity")
		}
		f.nodes[i].Key = key
	}
}
func (f *PortTestLists) State(heads ...int) []PortTestListNodeState {
	out := make([]PortTestListNodeState, len(f.nodes))
	for i, n := range f.nodes {
		key := n.Key
		for _, h := range heads {
			if i == h {
				if key != uint32(uintptr(f.ptr(h))) {
					panic("list sentinel changed")
				}
				key = 0
			}
		}
		out[i] = PortTestListNodeState{f.id(n.Next), f.id(n.Prev), key, n.Payload}
	}
	return out
}
func (f *PortTestLists) Op(op string, a, b int) int {
	p, q := f.ptr(a), f.ptr(b)
	switch op {
	case "append":
		C.nox_common_list_append_4258E0((*C.nox_list_item_t)(p), (*C.nox_list_item_t)(q))
		return 0
	case "prepend":
		return f.id(unsafe.Pointer(C.sub_425900((*C.uint)(p), (*C.uint)(q))))
	case "ascending":
		return listAscending((*legacyListNode)(p), (*legacyListNode)(q))
	case "descending":
		C.sub_4257F0((*C.int)(p), (*C.uint)(q))
		return 0
	case "remove":
		C.nox_common_list_remove_425920(p)
		return 0
	case "first":
		return f.id(unsafe.Pointer(C.nox_common_list_getFirstSafe_425890((*C.nox_list_item_t)(p))))
	case "next":
		return f.id(unsafe.Pointer(C.nox_common_list_getNextSafe_4258A0((*C.nox_list_item_t)(p))))
	case "raw-next":
		return f.id(unsafe.Pointer(C.nox_common_list_getNext_425940((*C.nox_list_item_t)(p))))
	case "prev":
		return f.id(unsafe.Pointer(uintptr(C.sub_425960(C.int(uintptr(p))))))
	}
	panic(op)
}
func (f *PortTestLists) At(head, index int) int {
	return f.id(unsafe.Pointer(C.sub_4258C0((**C.uint)(f.ptr(head)), C.int(index))))
}

type PortTestGroupState struct {
	ID      uint32
	Name    []uint16
	Team    uint32
	Members []int32
}
type PortTestGroups struct{ saved []byte }

func PortTestGroupsOpen() *PortTestGroups {
	p := unsafe.Slice(memmap.PtrUint8(0x5D4594, 599460), 16)
	f := &PortTestGroups{saved: bytes.Clone(p)}
	clear(p)
	playerGroupsInit()
	return f
}
func (f *PortTestGroups) Init() { playerGroupsInit() }
func (f *PortTestGroups) Free() { playerGroupsFree() }
func (f *PortTestGroups) Close() {
	f.Free()
	copy(unsafe.Slice(memmap.PtrUint8(0x5D4594, 599460), 16), f.saved)
}
func (f *PortTestGroups) Find(id uint32) unsafe.Pointer {
	return unsafe.Pointer(C.sub_425A70(C.int(id)))
}
func (f *PortTestGroups) Exists(id uint32) bool { return playerGroupFind(id) != nil }
func (f *PortTestGroups) Add(id uint32, name []uint16) bool {
	if len(name) > 9 {
		panic("group fixture name too long")
	}
	p, free := alloc.Make([]uint16{}, len(name)+1)
	defer free()
	copy(p, name)
	return playerGroupAdd(id, &p[0]) != nil
}
func (f *PortTestGroups) Member(id uint32, index int32) {
	p := f.Find(id)
	if p == nil {
		panic("missing fixture group")
	}
	playerGroupAddMember((*playerGroup)(p), index)
}
func (f *PortTestGroups) Remove(id uint32, index int32) bool {
	p := f.Find(id)
	if p == nil {
		panic("missing fixture group")
	}
	// The C result can point into a freed group; compare address bits only.
	result := C.sub_425B60(p, C.int(index))
	return uintptr(unsafe.Pointer(result)) == uintptr(p)+40
}
func (f *PortTestGroups) State() (out []PortTestGroupState) {
	head := memmap.PtrOff(0x5D4594, 599460)
	prev := head
	for p := unsafe.Pointer(C.sub_425A50()); p != nil; p = unsafe.Pointer(C.sub_425A60((*C.int)(p))) {
		words := unsafe.Slice((*uint32)(p), 13)
		if words[1] != uint32(uintptr(prev)) {
			panic("group backward link")
		}
		prev = p
		v := PortTestGroupState{ID: words[8], Name: append([]uint16(nil), unsafe.Slice((*uint16)(unsafe.Add(p, 12)), 10)...), Team: words[9]}
		childHead := unsafe.Add(p, 40)
		childPrev := childHead
		for q := unsafe.Pointer(listNext(&(*playerGroup)(p).members)); q != nil; q = unsafe.Pointer(listNext((*legacyListNode)(q))) {
			w := unsafe.Slice((*uint32)(q), 4)
			if w[1] != uint32(uintptr(childPrev)) {
				panic("member backward link")
			}
			childPrev = q
			v.Members = append(v.Members, int32(w[3]))
		}
		if words[11] != uint32(uintptr(childPrev)) {
			panic("member tail link")
		}
		out = append(out, v)
	}
	if *(*uint32)(unsafe.Add(head, 4)) != uint32(uintptr(prev)) {
		panic("group tail link")
	}
	return out
}

func (f *PortTestGroups) Initialized() uint32 { return memmap.Uint32(0x5D4594, 599472) }

func (f *PortTestLists) Pointer(i int) unsafe.Pointer { return f.ptr(i) }
func (f *PortTestGroups) AddLongName(id uint32, name []uint16) {
	p, free := alloc.Make([]uint16{}, len(name)+1)
	defer free()
	copy(p, name)
	if playerGroupAdd(id, &p[0]) == nil {
		panic("duplicate native group fixture")
	}
}
