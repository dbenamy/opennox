//go:build porttest

package legacy

/*
#include "GAME3_2.h"
#include "GAME3_3.h"
extern nox_server_netCodeCacheStruct nox_server_netCodeCache;
extern uint32_t nox_server_needInitNetCodeCache;
static uint32_t objectLookupInvoke(int op,uint32_t* a) {
 switch(op) {
 case 0:return (uint32_t)nox_xxx_getObjectByScrName_4DA4F0((char*)a[0]);
 case 1:return sub_4DA5C0(a[0],(const char*)a[1]);
 case 2:return sub_4DA660(a[0],(const char*)a[1]);
 case 3:return (uint32_t)nox_server_getObjectFromNetCode_4ECCB0(a[0]);
 case 4:return nox_server_netCodeCache_lookupObj_4ECD90(a[0]);
 case 5:sub_4ECDE0((uint32_t*)a[0],a[1]);return 0;
 case 6:sub_4ECE10((uint32_t*)a[0],a[1]);return 0;
 case 7:nox_server_netCodeCache_initArray_4ECE50();return 0;
 case 8:nox_server_netCodeCache_addObj_4ECEA0(a[0]);return 0;
 case 9:return nox_server_netCodeCache_nextUnused_4ECEF0();
 case 10:return sub_4ECF10(a[0]);
 case 11:sub_4ECFA0((nox_object_t*)a[0]);return 0;
 case 12:sub_4ECFE0();return 0;
 default:return 0xDEADBEEF;
 }
}
*/
import "C"

import (
	"unsafe"

	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/server"
)

type PortTestLookupObject struct {
	Name                  *string
	Code, ScriptID, Flags uint32
	Inventory             []int
}
type PortTestObjectLookupSpec struct {
	Objects                 []PortTestLookupObject
	Names                   []string
	Main, Pending, Missiles []int
	CacheValues             []int
	Uninitialized           bool
	InitFlag                uint32
}
type PortTestObjectLookupResult struct {
	NeedInit uint32
	Heads    [4]uint32
	Nodes    [16][3]uint32
}
type portTestObjectLookup struct {
	objects []*server.Object
	names   []unsafe.Pointer
	restore func()
}

func (p *portTestShopPools) objectLookupEnsure() {
	if p.reportLookup != nil || p.reports == nil || p.reports.Lookup == nil {
		return
	}
	sp := p.reports.Lookup
	st := &portTestObjectLookup{}
	p.reportLookup = st
	core := p.proxy.core
	oldMain, oldPending, oldMissiles := core.Objs.List, core.Objs.Pending, core.Objs.MissileList
	st.restore = func() { core.Objs.List, core.Objs.Pending, core.Objs.MissileList = oldMain, oldPending, oldMissiles }
	for i, o := range sp.Objects {
		ptr := p.objectiveRegion(772)
		clear(unsafe.Slice((*byte)(ptr), 772))
		u := (*server.Object)(ptr)
		u.NetCode = o.Code
		u.ScriptIDVal = int(o.ScriptID)
		u.ObjFlags = object.Flags(o.Flags)
		if o.Name != nil {
			u.IDPtr = p.objectiveString(*o.Name)
		}
		p.identify(ptr, 940001+uint32(i))
		st.objects = append(st.objects, u)
	}
	ref := func(id int) *server.Object {
		if id == 0 {
			return nil
		}
		if id < 1 || id > len(st.objects) {
			panic("lookup object reference")
		}
		return st.objects[id-1]
	}
	for i, o := range sp.Objects {
		var first *server.Object
		for j := len(o.Inventory) - 1; j >= 0; j-- {
			u := ref(o.Inventory[j])
			u.InvNextItem = first
			first = u
		}
		st.objects[i].InvFirstItem = first
	}
	list := func(ids []int) *server.Object {
		var head *server.Object
		for i := len(ids) - 1; i >= 0; i-- {
			u := ref(ids[i])
			u.ObjNext = head
			head = u
		}
		return head
	}
	core.Objs.List, core.Objs.Pending, core.Objs.MissileList = list(sp.Main), list(sp.Pending), list(sp.Missiles)
	for _, name := range sp.Names {
		st.names = append(st.names, p.objectiveString(name))
	}
	C.nox_server_netCodeCache = C.nox_server_netCodeCacheStruct{}
	C.nox_server_needInitNetCodeCache = 1
	if !sp.Uninitialized {
		C.nox_server_netCodeCache_initArray_4ECE50()
	}
	if sp.InitFlag != 0 {
		C.nox_server_needInitNetCodeCache = C.uint32_t(sp.InitFlag)
	}
	for i := 0; i < 16; i++ {
		p.identify(unsafe.Pointer(&C.nox_server_netCodeCache.objArray[i]), 920000+uint32(i))
	}
	p.identify(unsafe.Pointer(&C.nox_server_netCodeCache.firstUsedObject), 930000)
	p.identify(unsafe.Pointer(&C.nox_server_netCodeCache.firstFreeObject), 930001)
	if len(sp.CacheValues) > 16 {
		panic("lookup cache seed bounds")
	}
	for i, id := range sp.CacheValues {
		C.nox_server_netCodeCache.objArray[i].value = ref(id).CObj()
	}
}
func (p *portTestShopPools) objectLookupArg(a PortTestGameplayReportArg) uint32 {
	st := p.reportLookup
	if st == nil {
		panic("lookup fixture not prepared")
	}
	var ptr unsafe.Pointer
	var size int
	switch a.Kind {
	case "lookup-object":
		if a.Ref == 0 {
			if a.Offset != 0 {
				panic("lookup null offset")
			}
			return 0
		}
		if a.Ref < 1 || a.Ref > len(st.objects) {
			panic("lookup object argument")
		}
		ptr = st.objects[a.Ref-1].CObj()
		size = 772
	case "lookup-name":
		if a.Ref < 0 || a.Ref >= len(st.names) {
			panic("lookup name argument")
		}
		ptr = st.names[a.Ref]
		size = len(p.reports.Lookup.Names[a.Ref]) + 1
	case "lookup-node":
		if a.Ref < 0 || a.Ref >= 16 {
			panic("lookup node argument")
		}
		ptr = unsafe.Pointer(&C.nox_server_netCodeCache.objArray[a.Ref])
		size = 12
	case "lookup-head":
		if a.Ref == 0 {
			ptr = unsafe.Pointer(&C.nox_server_netCodeCache.firstUsedObject)
		} else if a.Ref == 1 {
			ptr = unsafe.Pointer(&C.nox_server_netCodeCache.firstFreeObject)
		} else {
			panic("lookup head argument")
		}
		size = 8
	default:
		panic("lookup argument kind")
	}
	if a.Offset < 0 || a.Offset >= size {
		panic("lookup argument offset")
	}
	return uint32(uintptr(unsafe.Add(ptr, a.Offset)))
}
func objectLookupInvoke(op int, args [5]uint32) uint32 {
	if op == 13 {
		// Fixture-only object mutations between real cache operations.
		off := int(args[1])
		if off != 16 && off != 36 && off != 44 {
			panic("lookup mutation field")
		}
		*(*uint32)(unsafe.Add(unsafe.Pointer(uintptr(args[0])), off)) = args[2]
		return 0
	}
	if op < 0 || op > 12 {
		panic("lookup operation")
	}
	return uint32(C.objectLookupInvoke(C.int(op), (*C.uint32_t)(unsafe.Pointer(&args[0]))))
}
func (p *portTestShopPools) objectLookupSnapshot() *PortTestObjectLookupResult {
	if p.reportLookup == nil {
		return nil
	}
	out := &PortTestObjectLookupResult{NeedInit: uint32(C.nox_server_needInitNetCodeCache)}
	norm := func(v unsafe.Pointer) uint32 { return p.normalize(uint32(uintptr(v))) }
	out.Heads = [4]uint32{norm(unsafe.Pointer(C.nox_server_netCodeCache.firstFreeObject)), norm(unsafe.Pointer(C.nox_server_netCodeCache.lastFreeObject)), norm(unsafe.Pointer(C.nox_server_netCodeCache.firstUsedObject)), norm(unsafe.Pointer(C.nox_server_netCodeCache.lastUsedObject))}
	for i := range out.Nodes {
		n := C.nox_server_netCodeCache.objArray[i]
		out.Nodes[i] = [3]uint32{norm(n.value), norm(n.next), norm(n.prev)}
	}
	return out
}
