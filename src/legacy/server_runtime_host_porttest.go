//go:build porttest

package legacy

import (
	"fmt"
	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

type PortTestRuntimeHostSpec struct {
	Host, Unit, Marker, Missing, Cached, Reverse bool
	Count, MonsterMask, MonitorMask              int
}

func (p *portTestShopPools) runtimeHostContract() []uint32 {
	sp := p.proxy.callbacks.shop.spec.TemporaryUpdates.World.Objectives.Attack.Controls.RuntimeHost
	core := p.proxy.core
	host := core.Players.ByIndRaw(31)
	savedHost := *host
	host.Active = byte(bool2int(sp.Host))
	host.PlayerInd = 31
	host.PlayerUnit = nil
	if sp.Unit {
		host.PlayerUnit = p.resources.unit
	}
	defer func() { head := host.Field4580; *host = savedHost; host.Field4580 = head }()
	marker := p.items[0].u
	saved := []server.Object{*marker, *p.items[1].u, *p.items[2].u, *p.resources.unit}
	defer func() {
		*marker = saved[0]
		*p.items[1].u = saved[1]
		*p.items[2].u = saved[2]
		*p.resources.unit = saved[3]
	}()
	oldPlayer := p.resources.unit.UpdateDataPlayer().Player
	p.resources.unit.UpdateDataPlayer().Player = host
	defer func() { p.resources.unit.UpdateDataPlayer().Player = oldPlayer }()
	var missing []string
	if sp.Missing {
		missing = []string{"SaveGameLocation"}
	}
	restore := core.PortTestRewardTypes([]string{"SaveGameLocation"}, missing, true, 0, 0)
	defer restore()
	typ := core.Types.IndByID("SaveGameLocation")
	cache := memmap.PtrUint32(0x5D4594, 1563124)
	oldCache := *cache
	defer func() { *cache = oldCache }()
	*cache = 0
	if sp.Cached {
		*cache = 77
	}
	marker.TypeInd = uint16(typ)
	if sp.Cached {
		marker.TypeInd = 77
	}
	// A missing definition yields index zero, which cannot match these live objects.
	if marker.TypeInd == 0 {
		marker.TypeInd = 78
	}
	marker.ObjClass = object.ClassSimple
	marker.Field129 = nil
	*(*uint32)(unsafe.Add(marker.CObj(), 44)) = 0xaabbccdd
	p.resources.unit.Field129 = nil
	order := []*server.Object{p.items[1].u, p.items[2].u}
	if sp.Reverse {
		order[0], order[1] = order[1], order[0]
	}
	for i, child := range order {
		child.ObjClass = object.ClassSimple
		if sp.MonsterMask&(1<<i) != 0 {
			child.ObjClass = object.ClassMonster
		}
		child.ObjSubClass = 0x10203004
		*(*byte)(unsafe.Add(child.UpdateData, 1440)) = 0
		if sp.MonitorMask&(1<<i) != 0 {
			*(*byte)(unsafe.Add(child.UpdateData, 1440)) = 0x80
		}
		child.ObjOwner = nil
		child.Field128 = nil
	}
	for i := sp.Count - 1; i >= 0; i-- {
		order[i].ObjOwner = marker
		order[i].Field128 = marker.Field129
		marker.Field129 = order[i]
	}
	oldList := core.Objs.List
	defer func() { core.Objs.List = oldList }()
	core.Objs.List = p.resources.unit
	p.resources.unit.ObjNext = nil
	if sp.Marker {
		p.resources.unit.ObjNext = marker
		marker.ObjNext = nil
	}
	Nox_xxx_monstersAllBelongToHost_4DB6A0()
	moved := sp.Host && sp.Unit && sp.Marker && (sp.Cached || !sp.Missing)
	out := []uint32{*cache, p.normalize(uint32(uintptr(marker.Field129.CObj()))), p.normalize(uint32(uintptr(p.resources.unit.Field129.CObj()))), *(*uint32)(unsafe.Add(marker.CObj(), 44))}
	expectedCache := uint32(0)
	if sp.Cached {
		expectedCache = 77
	} else if sp.Host && sp.Unit {
		expectedCache = uint32(typ)
	}
	if *cache != expectedCache {
		panic("host ownership type cache")
	}
	for i, child := range order {
		owner := (*server.Object)(nil)
		if i < sp.Count {
			owner = marker
			if moved {
				owner = p.resources.unit
			}
		}
		sub := uint32(0x10203004)
		if moved && i < sp.Count && sp.MonsterMask&(1<<i) != 0 && sp.MonitorMask&(1<<i) != 0 {
			sub |= 0x80
		}
		if child.ObjOwner != owner || uint32(child.ObjSubClass) != sub {
			panic(fmt.Sprintf("host child %d: owner=%p/%p subclass=%x/%x", i, child.ObjOwner, owner, child.ObjSubClass, sub))
		}
		out = append(out, p.normalize(uint32(uintptr(child.ObjOwner.CObj()))), p.normalize(uint32(uintptr(child.Field128.CObj()))), uint32(child.ObjSubClass))
	}
	if moved {
		if marker.Field129 != nil || out[3] != 0 {
			panic("host marker not emptied")
		}
		next := p.resources.unit.Field129
		for i := sp.Count - 1; i >= 0; i-- {
			if next != order[i] {
				panic("host ownership linked order")
			}
			next = next.Field128
		}
		if next != nil {
			panic("host ownership list tail")
		}
	} else if out[3] != 0xaabbccdd {
		panic("inactive host transfer changed marker")
	}
	return out
}
