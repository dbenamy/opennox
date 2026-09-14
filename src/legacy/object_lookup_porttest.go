//go:build porttest

package legacy

import (
	"unsafe"

	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
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
	netCodeCacheState = netCodeCacheStorage{}
	netCodeCacheNeedInit = 1
	if !sp.Uninitialized {
		netCodeCacheInit()
	}
	if sp.InitFlag != 0 {
		netCodeCacheNeedInit = sp.InitFlag
	}
	for i := 0; i < 16; i++ {
		p.identify(unsafe.Pointer(&netCodeCacheState.nodes[i]), 920000+uint32(i))
	}
	p.identify(unsafe.Pointer(&netCodeCacheState.used.first), 930000)
	p.identify(unsafe.Pointer(&netCodeCacheState.free.first), 930001)
	if len(sp.CacheValues) > 16 {
		panic("lookup cache seed bounds")
	}
	for i, id := range sp.CacheValues {
		netCodeCacheState.nodes[i].value = ref(id)
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
		ptr = unsafe.Pointer(&netCodeCacheState.nodes[a.Ref])
		size = 12
	case "lookup-head":
		if a.Ref == 0 {
			ptr = unsafe.Pointer(&netCodeCacheState.used.first)
		} else if a.Ref == 1 {
			ptr = unsafe.Pointer(&netCodeCacheState.free.first)
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
	ptr := func(v uint32) unsafe.Pointer { return unsafe.Pointer(uintptr(v)) }
	obj := func(v uint32) *server.Object { return (*server.Object)(ptr(v)) }
	result := func(u *server.Object) uint32 { return uint32(uintptr(u.CObj())) }
	name := func(v uint32) string { return alloc.GoString((*byte)(ptr(v))) }
	switch op {
	case 0:
		return result(objectLookupByName(name(args[0])))
	case 1:
		return result(objectLookupNameAt(obj(args[0]), name(args[1]), true))
	case 2:
		return result(objectLookupNameAt(obj(args[0]), name(args[1]), false))
	case 3:
		return result(objectLookupByNetCode(args[0]))
	case 4:
		return result(netCodeCacheLookup(args[0]))
	case 5:
		(*netCodeCacheList)(ptr(args[0])).prepend((*netCodeCacheNode)(ptr(args[1])))
	case 6:
		(*netCodeCacheList)(ptr(args[0])).remove((*netCodeCacheNode)(ptr(args[1])))
	case 7:
		netCodeCacheInit()
	case 8:
		netCodeCacheAdd(obj(args[0]))
	case 9:
		return uint32(uintptr(unsafe.Pointer(netCodeCacheNextUnused())))
	case 10:
		return result(objectLookupByScriptID(args[0]))
	case 11:
		netCodeCacheInvalidate(obj(args[0]))
	case 12:
		netCodeCacheFlush()
	}
	return 0
}
func (p *portTestShopPools) objectLookupSnapshot() *PortTestObjectLookupResult {
	if p.reportLookup == nil {
		return nil
	}
	out := &PortTestObjectLookupResult{NeedInit: uint32(netCodeCacheNeedInit)}
	norm := func(v unsafe.Pointer) uint32 { return p.normalize(uint32(uintptr(v))) }
	out.Heads = [4]uint32{norm(unsafe.Pointer(netCodeCacheState.free.first)), norm(unsafe.Pointer(netCodeCacheState.free.last)), norm(unsafe.Pointer(netCodeCacheState.used.first)), norm(unsafe.Pointer(netCodeCacheState.used.last))}
	for i := range out.Nodes {
		n := netCodeCacheState.nodes[i]
		out.Nodes[i] = [3]uint32{norm(unsafe.Pointer(n.value)), norm(unsafe.Pointer(n.towardHead)), norm(unsafe.Pointer(n.towardTail))}
	}
	return out
}
