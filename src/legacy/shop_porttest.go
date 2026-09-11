//go:build porttest

package legacy

/*
#include <string.h>
#include "GAME4_1.h"
static int portTestShopPrice(int mode, int session, void* obj) {
	float bits;
	memcpy(&bits, &obj, sizeof(bits));
	return nox_xxx_shopGetItemCost_50E3D0(mode, session, bits);
}
*/
import "C"

import (
	"bytes"
	"encoding/binary"
	"unsafe"

	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

type PortTestShopItem struct {
	Pickup                        bool
	Type                          uint16
	Class, Subclass, Flags, Worth uint32
	Mods                          [4]bool
	ModPrice                      [4]int32
	Use                           [128]byte
	Health                        bool
	HP, MaxHP                     uint16
	QuestFlag                     byte
}
type PortTestShopStock struct {
	Type, Reward uint32
	Count        byte
	Mods         [4]bool
}
type PortTestShopSpec struct {
	Engine            *PortTestShopEngineSpec
	CaptureData       bool
	Load              *PortTestShopLoadSpec
	Peer              bool
	Gold              [2]uint32
	ProtectedGold     bool
	Sequence          []PortTestShopAction
	Items             []PortTestShopItem
	Op, Mode, Session int
	Item              PortTestShopItem
	Buy, Sell         uint32
	Balance           map[string]float64
	GuideWorth        int
	SpellPrices       map[int]int
	Cache             [3]uint32
	Stock             []PortTestShopStock
	Entry             int
	NilItem, NilEntry bool
	Priority          uint32
}
type PortTestShopResult struct {
	Sequence []PortTestShopStep `json:",omitempty"`
	Memory   [][]uint32
	Cache    [3]uint32
	Intact   bool
}
type portTestShopState struct {
	pools          *portTestShopPools
	load           *portTestShopLoad
	configure      func(map[string]float64, int, map[int]int)
	blocks, before [][]byte
	spec           *PortTestShopSpec
}

// Each slab has eight guard bytes on both ends. Objects reserve the Go EXT
// tail too; only the 772-byte C ABI prefix is observed as object contents.
func (s *portTestShopState) ptr(i int) unsafe.Pointer { return unsafe.Pointer(&s.blocks[i][8]) }
func (s *portTestShopState) data(i int) []byte        { b := s.blocks[i]; return b[8 : len(b)-8] }
func (s *portTestShopState) item() *server.Object     { return (*server.Object)(s.ptr(0)) }
func (s *portTestShopState) npc() *server.Object      { return (*server.Object)(s.ptr(1)) }

func portTestShopEnvironment(proxy *portTestRoamOwnerServer) func() {
	configure, freeServer := proxy.core.PortTestShopEnvironment()
	s := &portTestShopState{configure: configure}
	proxy.callbacks.shop = s
	s.pools = portTestShopPoolsEnvironment(proxy)
	s.load = portTestShopLoadEnvironment(proxy)
	sizes := []int{int(unsafe.Sizeof(server.Object{})), int(unsafe.Sizeof(server.Object{})), 64, 1724, 20, 128, 2200, 8, 4 * int(unsafe.Sizeof(server.ModifierEff{})), 16}
	var frees []func()
	for _, n := range sizes {
		b, free := alloc.Make([]byte{}, n+16)
		s.blocks = append(s.blocks, b)
		frees = append(frees, free)
	}
	var old [3]uint32
	for i := range old {
		old[i] = *memmap.PtrUint32(0x5D4594, 2386504+uintptr(4*i))
	}
	return func() {
		s.pools.restore()
		s.load.restore()
		for i, v := range old {
			*memmap.PtrUint32(0x5D4594, 2386504+uintptr(4*i)) = v
		}
		for _, free := range frees {
			free()
		}
		freeServer()
	}
}

func portTestShopPrepare(proxy *portTestRoamOwnerServer, sp *PortTestShopSpec) {
	s := proxy.callbacks.shop
	s.spec = sp
	s.pools.prepare()
	s.configure(sp.Balance, sp.GuideWorth, sp.SpellPrices)
	s.load.prepare(sp.Load)
	for i, b := range s.blocks {
		clear(b)
		for j := 0; j < 8; j++ {
			b[j], b[len(b)-8+j] = 0xa5, 0x5a
		}
		proxy.life.ids[uint32(uintptr(s.ptr(i)))] = uint32(50000 + i)
	}
	u, npc := s.item(), s.npc()
	*u, *npc = *proxy.combat.actor, *proxy.combat.actor
	u.InvFirstItem, u.InvNextItem, u.Field125, u.InvHolder = nil, nil, nil, nil
	u.TypeInd, u.ObjClass, u.ObjSubClass, u.ObjFlags = sp.Item.Type, object.Class(sp.Item.Class), object.SubClass(sp.Item.Subclass), object.Flags(sp.Item.Flags)
	u.Worth, u.NetCode = sp.Item.Worth, 0x12345678
	u.InitData, u.UseData.Ptr, u.UpdateData = s.ptr(4), s.ptr(5), s.ptr(6)
	u.HealthData = nil
	if sp.Item.Health {
		u.HealthData = (*server.HealthData)(s.ptr(7))
		u.HealthData.Cur, u.HealthData.Max = sp.Item.HP, sp.Item.MaxHP
	}
	copy(s.data(5), sp.Item.Use[:])
	s.data(6)[4] = sp.Item.QuestFlag
	for i, enabled := range sp.Item.Mods {
		p := unsafe.Add(s.ptr(8), i*int(unsafe.Sizeof(server.ModifierEff{})))
		*(*uint32)(unsafe.Add(p, 4)) = uint32(40 + i)
		(*server.ModifierEff)(p).Price20 = sp.Item.ModPrice[i]
		proxy.life.ids[uint32(uintptr(p))] = uint32(50100 + i)
		if enabled {
			*(*unsafe.Pointer)(unsafe.Add(u.InitData, 4*i)) = p
		}
	}
	npc.ObjClass, npc.ObjSubClass, npc.InitData = object.ClassMonster, 8, s.ptr(3)
	// The vendor queries consume only init data. Do not inherit unused health
	// and monster-update pointers from the surrounding actor fixture.
	npc.HealthData, npc.UpdateData = nil, nil
	*(*uint32)(unsafe.Add(npc.InitData, 1716)), *(*uint32)(unsafe.Add(npc.InitData, 1720)) = sp.Buy, sp.Sell
	if len(sp.Stock) > 60 {
		panic("shop stock fixture capacity")
	}
	s.data(3)[0] = byte(len(sp.Stock))
	for i, entry := range sp.Stock {
		b := s.data(3)[4+i*28 : 4+(i+1)*28]
		binary.LittleEndian.PutUint32(b, entry.Type)
		b[4] = entry.Count
		binary.LittleEndian.PutUint32(b[8:], entry.Reward)
		for j, enabled := range entry.Mods {
			if enabled {
				*(*unsafe.Pointer)(unsafe.Pointer(&b[12+4*j])) = unsafe.Add(s.ptr(8), j*int(unsafe.Sizeof(server.ModifierEff{})))
			}
		}
	}
	session := unsafe.Slice((*uint32)(s.ptr(2)), 16)
	session[2], session[3] = uint32(uintptr(npc.CObj())), uint32(uintptr(proxy.life.players[0].CObj()))
	if sp.Session == 3 {
		session[2], session[3] = session[3], session[2]
	}
	if sp.Session >= 2 {
		session[4] = 1
	}
	*(*unsafe.Pointer)(s.ptr(9)) = u.CObj()
	*(*uint32)(unsafe.Add(s.ptr(9), 4)) = sp.Priority
	for i, v := range sp.Cache {
		*memmap.PtrUint32(0x5D4594, 2386504+uintptr(4*i)) = v
	}
	s.before = nil
	for _, b := range s.blocks {
		s.before = append(s.before, bytes.Clone(b))
	}
}

func portTestShopCall(proxy *portTestRoamOwnerServer) uint32 {
	s := proxy.callbacks.shop
	sp := s.spec
	u := s.item().CObj()
	session := C.int(uintptr(s.ptr(2)))
	switch sp.Op {
	case 0:
		if sp.Session == 0 {
			session = 0
		}
		return uint32(C.portTestShopPrice(C.int(sp.Mode), session, u))
	case 1:
		entry := unsafe.Add(s.ptr(3), 4+28*sp.Entry)
		if sp.NilItem {
			u = nil
		}
		if sp.NilEntry {
			entry = nil
		}
		return uint32(shopStockMatches((*server.Object)(u), (*shopStockEntry)(entry)))
	case 2:
		return uint32(C.nox_xxx_getSomeShopData_5103A0(session, C.int(uintptr(u))))
	case 3:
		return uint32(shopStockKey((*shopItem)(s.ptr(9))))
	case 4:
		s.pools.run()
		return 0
	default:
		panic("shop fixture operation")
	}
}

func portTestShopTrace(proxy *portTestRoamOwnerServer, normalize func(uint32) uint32) *PortTestShopResult {
	s := proxy.callbacks.shop
	r := &PortTestShopResult{Intact: true, Sequence: s.pools.steps}
	defer s.pools.cleanup()
	for i, b := range s.blocks {
		if s.spec.Engine == nil {
			r.Intact = r.Intact && bytes.Equal(b, s.before[i])
		} else {
			r.Intact = r.Intact && bytes.Equal(b[:8], s.before[i][:8]) && bytes.Equal(b[len(b)-8:], s.before[i][len(b)-8:])
		}
		data := s.data(i)
		if i < 2 {
			data = data[:772]
		}
		var words []uint32
		for off := 0; off+4 <= len(data); off += 4 {
			words = append(words, normalize(binary.LittleEndian.Uint32(data[off:])))
		}
		r.Memory = append(r.Memory, words)
	}
	for i := range r.Cache {
		r.Cache[i] = *memmap.PtrUint32(0x5D4594, 2386504+uintptr(4*i))
	}
	return r
}
