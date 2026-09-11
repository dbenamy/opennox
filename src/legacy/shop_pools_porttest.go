//go:build porttest

package legacy

/*
#include <string.h>
#include "GAME4_1.h"
extern void* nox_alloc_tradeSession_2386492;
extern void* nox_alloc_tradeItems_2386496;
extern uint32_t dword_5d4594_2386500;
extern uint32_t dword_5d4594_1565512;

*/
import "C"

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"unsafe"

	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/ntype"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

const (
	PortTestShopCreate = iota
	PortTestShopDestroy
	PortTestShopReset
	PortTestShopAdd
	PortTestShopFind
	PortTestShopFreeList
	PortTestShopPacket
	PortTestShopSet
	PortTestShopGold
	PortTestShopOffer
	PortTestShopTotal
	PortTestShopBalance
	PortTestShopAccept
	PortTestShopWithdraw
	PortTestShopCancel
	PortTestShopPlayerCleanup
	PortTestShopInventory
	PortTestShopRepairQuote
	PortTestShopRepair
	PortTestShopSell
	PortTestShopLoad
	PortTestShopLookup
	PortTestShopDetach
)

type PortTestShopAction struct {
	Op, Session, Item int
	Side              int
	Value             uint32
}
type PortTestShopPacketResult struct {
	Recipient, Ordered byte
	A4, A5             uint32
	Sequence           [32]uint16
	Data               []byte
}
type PortTestShopStep struct {
	EffectsUseData           []uint32   `json:",omitempty"`
	TemporaryUpdatesData     []uint32   `json:",omitempty"`
	EquipmentData            []uint32   `json:",omitempty"`
	InventoryData            []uint32   `json:",omitempty"`
	ResourceData             [][]uint32 `json:",omitempty"`
	ResourceMessages         [][]byte   `json:",omitempty"`
	EngineStateRequests      []uint32   `json:",omitempty"`
	EngineMessages           [][]byte   `json:",omitempty"`
	EngineState              []uint32   `json:",omitempty"`
	ObjectData               [][]uint32 `json:",omitempty"`
	Cached                   [32]uint32
	Return, Head             uint32
	Alive                    int
	Sessions, Nodes, Objects [][]uint32
	ObjectDigest             string                     `json:",omitempty"`
	Players                  [][]uint32                 `json:",omitempty"`
	Packets                  []PortTestShopPacketResult `json:",omitempty"`
	Protection               []uint32                   `json:",omitempty"`
}
type portTestShopOwned struct {
	health            unsafe.Pointer
	initSize, useSize int
	undefinedTail     bool
	u                 *server.Object
	init              unsafe.Pointer
	alive             bool
}
type portTestShopPools struct {
	equipment           *portTestEquipment
	effectsUse          *portTestEffectsUse
	temporary           *portTestTemporaryUpdates
	inventory           *portTestInventory
	resources           *portTestResources
	engineStateRequests []uint32
	proxy               *portTestRoamOwnerServer
	restore             func()
	sessions            []unsafe.Pointer
	owned               []*portTestShopOwned
	items               []*portTestShopOwned
	ids                 map[uint32]uint32
	steps               []PortTestShopStep
	initialAlive        int
	players             func() [][]uint32
	freePlayers         func()
	prepareProtection   func(uint32)
	snapshotProtection  func() ([]uint32, bool)
}

func portTestShopPoolsEnvironment(proxy *portTestRoamOwnerServer) *portTestShopPools {
	p := &portTestShopPools{proxy: proxy}
	prepareProtection, snapshotProtection, freeProtection := portTestPenaltyProtectionEnvironment()
	p.prepareProtection, p.snapshotProtection = prepareProtection, snapshotProtection
	oldSessions, oldItems, oldHead := C.nox_alloc_tradeSession_2386492, C.nox_alloc_tradeItems_2386496, C.dword_5d4594_2386500
	table := unsafe.Slice(memmap.PtrUint32(0x5D4594, 2386364), 32)
	oldTable := append([]uint32(nil), table...)
	C.nox_alloc_tradeSession_2386492, C.nox_alloc_tradeItems_2386496, C.dword_5d4594_2386500 = nil, nil, 0
	if shopInit() == 0 {
		panic("shop fixture pools")
	}
	p.restore = func() {
		p.cleanup()
		freeProtection()
		shopFree()
		C.nox_alloc_tradeSession_2386492, C.nox_alloc_tradeItems_2386496, C.dword_5d4594_2386500 = oldSessions, oldItems, oldHead
		copy(table, oldTable)
	}
	return p
}

func shopTestWords(ptr unsafe.Pointer, n int) []uint32 { return unsafe.Slice((*uint32)(ptr), n) }
func shopTestPointer(v uint32) unsafe.Pointer          { return unsafe.Pointer(uintptr(v)) }

func (p *portTestShopPools) normalize(v uint32) uint32 {
	if id, ok := p.ids[v]; ok {
		return id
	}
	if id, ok := p.proxy.life.ids[v]; ok {
		return id
	}
	return v
}
func (p *portTestShopPools) identify(ptr unsafe.Pointer, id uint32) {
	p.ids[uint32(uintptr(ptr))] = id
	p.proxy.life.ids[uint32(uintptr(ptr))] = id
}
func (p *portTestShopPools) own(u *server.Object, id uint32) *portTestShopOwned {
	o := &portTestShopOwned{u: u, init: u.InitData, alive: true}
	if typ := p.proxy.core.Types.ByInd(int(u.TypeInd)); typ != nil {
		o.initSize, o.useSize = int(typ.InitDataSize), int(typ.UseDataSize)
	}
	o.health = unsafe.Pointer(u.HealthData)
	p.owned = append(p.owned, o)
	p.identify(u.CObj(), id)
	if u.InitData != nil {
		p.identify(u.InitData, id+10000)
	}
	if u.HealthData != nil {
		p.identify(unsafe.Pointer(u.HealthData), id+20000)
	}
	if u.UseData.Ptr != nil {
		p.identify(u.UseData.Ptr, id+30000)
	}
	if u.Xfer != nil {
		p.proxy.life.ids[uint32(uintptr(u.Xfer))] = 52000 + uint32(u.TypeInd)
	}
	return o
}
func (p *portTestShopPools) observeDelete(u *server.Object) {
	// Creations remain owned by the enclosing lifecycle fixture after deletion.
	if p.proxy.callbacks.shop.spec.TemporaryUpdates != nil {
		for _, created := range p.proxy.life.created {
			if created == u {
				return
			}
		}
	}
	if _, ok := p.ids[uint32(uintptr(u.CObj()))]; !ok {
		p.own(u, uint32(75000+len(p.owned)))
	}
}
func (p *portTestShopPools) markFreed(v uint32) {
	for _, o := range p.owned {
		if o.alive && uint32(uintptr(o.u.CObj())) == v {
			o.alive = false
			return
		}
	}
	panic("shop fixture free of unowned object")
}
func (p *portTestShopPools) markStockFreed(w []uint32) {
	for n := w[5]; n != 0; {
		item := shopTestWords(shopTestPointer(n), 4)
		p.markFreed(item[0])
		n = item[2]
	}
}
func (p *portTestShopPools) reset() {
	for q := uint32(C.dword_5d4594_2386500); q != 0; {
		w := shopTestWords(shopTestPointer(q), 16)
		if w[4] != 0 {
			p.markStockFreed(w)
		}
		q = w[14]
	}
	shopReset()
	for i := range p.sessions {
		p.sessions[i] = nil
	}
}
func (p *portTestShopPools) cleanup() {
	if p.ids == nil {
		return
	}
	p.reset()
	if p.freePlayers != nil {
		p.freePlayers()
		p.freePlayers, p.players = nil, nil
	}
	for _, o := range p.owned {
		if o.alive {
			o.u.InvFirstItem, o.u.InvNextItem, o.u.Field125, o.u.InvHolder = nil, nil, nil, nil
			o.u.InitData = nil
			p.proxy.core.Objs.FreeObject(o.u)
		}
		if o.init != nil {
			alloc.FreePtr(o.init)
		}
		if o.health != nil {
			alloc.FreePtr(o.health)
		}
	}
	expected := p.initialAlive
	if p.proxy.callbacks.shop.spec.EffectsUse != nil {
		// Projectiles are captured and freed later by the enclosing lifecycle trace.
		seen := make(map[*server.Object]bool)
		for _, u := range p.proxy.life.created {
			if seen[u] {
				panic("effects fixture duplicate created object")
			}
			seen[u] = true
			for _, o := range p.owned {
				if o.u == u {
					panic("effects fixture duplicate ownership")
				}
			}
			expected++
		}
	}
	if p.proxy.core.Objs.Alive != expected {
		panic("shop fixture object leak")
	}
	p.owned, p.items, p.sessions, p.ids = nil, nil, nil, nil
}
func (p *portTestShopPools) prepare() {
	p.cleanup()
	p.initialAlive = p.proxy.core.Objs.Alive
	p.ids = make(map[uint32]uint32)
	p.steps = nil
}
func (p *portTestShopPools) run() {
	s := p.proxy.callbacks.shop
	players := []PortTestSpawnPlayer{{Flags: 4}, {Flags: 4}}
	if s.spec.Engine != nil {
		players = append(players, PortTestSpawnPlayer{Flags: 4})
	}
	p.players, p.freePlayers = portTestSpawnPlayers(p.proxy, players)
	for i := 0; i < 2; i++ {
		pl := p.proxy.life.players[i].UpdateDataPlayer().Player
		pl.GoldVal, pl.ProtPlayerGold = s.spec.Gold[i], 0
	}
	p.prepareProtection(s.spec.Gold[0])
	if s.spec.ProtectedGold {
		p.proxy.life.players[0].UpdateDataPlayer().Player.ProtPlayerGold = 0x40000001
	}
	defer p.enginePrepare()()
	defer p.resourcePrepare()()
	defer p.inventoryPrepare()()
	defer p.equipmentPrepare()()
	defer p.effectsUsePrepare()()
	defer p.temporaryPrepare()()
	for i, spec := range s.spec.Items {
		// Actual allocator objects let destruction execute the retained object
		// free path. Price/charge/modifier arithmetic has guarded query fixtures.
		u := p.proxy.core.NewObjectByTypeInd(23)
		u.TypeInd, u.ObjClass, u.ObjSubClass, u.ObjFlags = spec.Type, object.Class(spec.Class), object.SubClass(spec.Subclass), object.Flags(spec.Flags)
		u.Worth, u.NetCode = spec.Worth, uint32(70000+i)
		init, _ := alloc.Make([]byte{}, 20)
		u.InitData = unsafe.Pointer(&init[0])
		if spec.Health {
			hp, _ := alloc.New(server.HealthData{})
			hp.Cur, hp.Max = spec.HP, spec.MaxHP
			u.HealthData = hp
		}
		if spec.Use != [128]byte{} {
			data, _ := alloc.Make([]byte{}, 128)
			copy(data, spec.Use[:])
			u.UseData.Ptr = unsafe.Pointer(&data[0])
		}
		p.engineItem(u, spec)
		o := p.own(u, uint32(70000+i))
		o.initSize = 20
		if u.UseData.Ptr != nil {
			o.useSize = 128
		}
		p.items = append(p.items, o)
	}
	p.inventoryItems()
	p.equipmentItems()
	p.effectsUseItems()
	p.temporaryItems()
	for _, a := range s.spec.Sequence {
		var q unsafe.Pointer
		if a.Op != PortTestShopCreate && a.Op != PortTestShopReset && a.Op != PortTestShopPlayerCleanup && a.Op != PortTestTradeCreatePlayer && a.Op != PortTestTradeStart && a.Op < 200 {
			q = p.sessions[a.Session]
			if q == nil {
				panic("shop fixture stale session")
			}
		}
		var rv uint32
		var words []uint32
		if q != nil {
			words = shopTestWords(q, 16)
		}
		switch a.Op {
		case PortTestShopCreate:
			q = unsafe.Pointer(shopCreate())
			p.sessions = append(p.sessions, q)
			if q != nil {
				id := uint32(60000 + len(p.sessions) - 1)
				p.identify(q, id)
				w := shopTestWords(q, 16)
				// Capture the constructor before installing participants/mode.
				p.own((*server.Object)(shopTestPointer(w[12])), id+1000)
				p.own((*server.Object)(shopTestPointer(w[13])), id+2000)
				rv = uint32(uintptr(q))
				p.steps = append(p.steps, p.snapshot(rv))
				w[2], w[3] = uint32(uintptr(s.npc().CObj())), uint32(uintptr(p.proxy.life.players[0].CObj()))
				if s.spec.Session == 3 {
					w[2], w[3] = w[3], w[2]
				}
				if s.spec.Peer {
					w[2], w[3] = uint32(uintptr(p.proxy.life.players[0].CObj())), uint32(uintptr(p.proxy.life.players[1].CObj()))
				}
				w[4] = a.Value
				for _, v := range w[2:4] {
					u := (*server.Object)(shopTestPointer(v))
					if u.Class().Has(object.ClassPlayer) {
						*(*unsafe.Pointer)(unsafe.Add(u.UpdateData, 280)) = q
					}
				}
				continue
			}
		case PortTestShopDestroy:
			w := shopTestWords(q, 16)
			p.markStockFreed(w)
			p.markFreed(w[12])
			p.markFreed(w[13])
			shopDestroy((*shopSession)(q))
			p.sessions[a.Session] = nil
		case PortTestShopReset:
			p.reset()
		case PortTestShopAdd:
			rv = uint32(uintptr(unsafe.Pointer(shopAdd((*shopSession)(q), p.items[a.Item].u))))
		case PortTestShopFind:
			w := shopTestWords(q, 16)
			rv = uint32(uintptr(unsafe.Pointer(shopFind((*shopItem)(shopTestPointer(w[5])), a.Value))))
		case PortTestShopFreeList:
			w := shopTestWords(q, 16)
			rv = uint32(shopFreeList((*shopItem)(shopTestPointer(w[5]))))
			w[5] = 0
		case PortTestShopSet:
			if a.Item < 6 || a.Item > 11 || a.Item == 8 || a.Item == 9 {
				panic("shop fixture set scalar")
			}
			words[a.Item] = a.Value
		case PortTestShopGold:
			gold := (*server.Object)(shopTestPointer(words[12+a.Side]))
			*(*uint32)(gold.InitData) = a.Value
		case PortTestShopOffer:
			n := alloc.AsClass(C.nox_alloc_tradeItems_2386496).NewObject()
			if n == nil {
				panic("shop fixture offer allocation")
			}
			w := shopTestWords(n, 4)
			w[0], w[1], w[2] = uint32(uintptr(p.items[a.Item].u.CObj())), a.Value, words[8+a.Side]
			if w[2] != 0 {
				shopTestWords(shopTestPointer(w[2]), 4)[3] = uint32(uintptr(n))
			}
			words[8+a.Side] = uint32(uintptr(n))
		case PortTestShopTotal:
			rv = uint32(shopTotal((*shopSession)(q), (*server.Object)(shopTestPointer(words[2+a.Side]))))
		case PortTestShopBalance:
			rv = shopBalance((*shopSession)(q))
		case PortTestShopPacket:
			player := C.int(words[2+a.Side])
			item := C.int(uintptr(s.item().CObj()))
			switch a.Item {
			case 0:
				rv = shopSendCode(objectFromInt(player), objectFromInt(item), 2505)
			case 1:
				rv = shopSendItem(objectFromInt(player), (*shopItem)(s.ptr(9)))
			case 2:
				rv = uint32(shopSendShort(objectFromInt(player), 457, 0))
			case 3:
				rv = uint32(shopSendShort(objectFromInt(player), 713, 1))
			case 4:
				rv = uint32(shopSendShort(objectFromInt(player), 1993, 1))
			case 5:
				rv = shopSendAcceptance(objectFromInt(player), (*shopSession)(q))
			case 6:
				rv = shopSendGold(objectFromInt(player), (*shopSession)(q))
			case 7:
				rv = uint32(shopSendCode(objectFromInt(player), objectFromInt(item), 1481))
			default:
				panic("shop fixture packet")
			}
		case PortTestShopWithdraw:
			rv = uint32(C.nox_xxx_tradeP2PAddOfferMB_50FE20(C.int(uintptr(q)), C.int(a.Value)))
		case PortTestShopInventory:
			owner := (*server.Object)(shopTestPointer(words[2+a.Side]))
			it := p.items[a.Item].u
			it.InvHolder, it.InvNextItem, it.Field125 = owner, owner.InvFirstItem, nil
			if it.InvNextItem != nil {
				it.InvNextItem.Field125 = it
			}
			owner.InvFirstItem = it
		case PortTestShopRepairQuote:
			rv = uint32(uintptr(unsafe.Pointer(C.sub_5108D0(C.int(words[2+a.Side]), C.int(uintptr(q)), C.int(a.Value)))))
		case PortTestShopRepair:
			rv = uint32(uintptr(unsafe.Pointer(C.sub_510AE0((*C.int)(shopTestPointer(words[2+a.Side])), C.int(uintptr(q)), (*C.uint32_t)(shopTestPointer(a.Value))))))
		case PortTestShopSell:
			C.sub_510D10((*C.int)(shopTestPointer(words[2+a.Side])), C.int(uintptr(q)), C.int(a.Item), C.uint(a.Value))
		case PortTestShopLookup:
			rv = uint32(C.sub_510DE0(C.int(words[2+a.Side]), C.int(a.Value)))
		case PortTestShopDetach:
			rv = uint32(shopDetach((*shopSession)(q), (*server.Object)(shopTestPointer(words[2+a.Side]))))
		case PortTestShopLoad:
			shopLoad((*shopSession)(q))
			for n := words[5]; n != 0; {
				w := shopTestWords(shopTestPointer(n), 4)
				if _, ok := p.ids[w[0]]; !ok {
					u := (*server.Object)(shopTestPointer(w[0]))
					o := p.own(u, uint32(75000+len(p.owned)))
					// The original loader fills four modifier words in a five-word
					// local. Only its uninitialized fifth word is outside the oracle.
					if u.Class()&0x13001000 != 0 && u.InitData != nil {
						for _, v := range shopTestWords(u.InitData, 4) {
							if v != 0 {
								o.undefinedTail = true
							}
						}
						if u.Class()&0x1000 != 0 && u.SubClass()&0x47f0000 != 0 {
							o.undefinedTail = true
						}
					}
				}
				n = w[2]
			}
		case PortTestShopAccept, PortTestShopCancel:
			gold := [2]uint32{words[12], words[13]}
			// Remember stock ownership before a destructor can invalidate nodes.
			var stock []uint32
			for n := words[5]; n != 0; {
				w := shopTestWords(shopTestPointer(n), 4)
				stock = append(stock, w[0])
				n = w[2]
			}
			if a.Op == PortTestShopAccept {
				C.nox_xxx_tradeAccept_50F5A0(C.int(uintptr(q)), C.int(words[2+a.Side]))
			} else {
				C.nox_xxx_shopCancelSession_510DC0(q)
			}
			live := false
			for n := uint32(C.dword_5d4594_2386500); n != 0; n = shopTestWords(shopTestPointer(n), 16)[14] {
				if n == uint32(uintptr(q)) {
					live = true
				}
			}
			if !live {
				p.markFreed(gold[0])
				p.markFreed(gold[1])
				for _, obj := range stock {
					p.markFreed(obj)
				}
				p.sessions[a.Session] = nil
			}
		case PortTestShopPlayerCleanup:
			q = *memmap.PtrPtr(0x5D4594, 2386364+uintptr(4*a.Item))
			if q != nil {
				w := shopTestWords(q, 16)
				p.markStockFreed(w)
				p.markFreed(w[12])
				p.markFreed(w[13])
				for i, old := range p.sessions {
					if old == q {
						p.sessions[i] = nil
					}
				}
			}
			C.sub_510E20(C.int(a.Item))
		default:
			if a.Op >= 600 {
				rv = p.temporaryAction(a)
			} else if a.Op >= 500 {
				rv = p.effectsUseAction(a)
			} else if a.Op >= 400 {
				rv = p.equipmentAction(a)
			} else if a.Op >= 300 {
				rv = p.inventoryAction(a)
			} else if a.Op >= 200 {
				rv = p.resourceAction(a)
			} else {
				rv = p.engineAction(a, q)
			}
		}
		p.engineDiscover()
		p.steps = append(p.steps, p.snapshot(rv))
	}
}

func (p *portTestShopPools) snapshot(rv uint32) PortTestShopStep {
	r := PortTestShopStep{Alive: p.proxy.core.Objs.Alive - p.initialAlive}
	r.ResourceData, r.ResourceMessages = p.resourceSnapshot()
	r.InventoryData = p.inventorySnapshot()
	r.EquipmentData = p.equipmentSnapshot()
	r.EffectsUseData = p.effectsUseSnapshot()
	r.TemporaryUpdatesData = p.temporarySnapshot()
	var sessions, nodes []unsafe.Pointer
	seen := make(map[unsafe.Pointer]bool)
	for q := shopTestPointer(uint32(C.dword_5d4594_2386500)); q != nil; {
		if len(sessions) == 64 || seen[q] {
			panic("shop fixture session cycle")
		}
		seen[q] = true
		sessions = append(sessions, q)
		w := shopTestWords(q, 16)
		for _, first := range []uint32{w[5], w[8], w[9]} {
			for n := shopTestPointer(first); n != nil; {
				if len(nodes) == 500 || seen[n] {
					panic("shop fixture item cycle")
				}
				seen[n] = true
				if _, ok := p.ids[uint32(uintptr(n))]; !ok {
					p.identify(n, uint32(90000+len(p.ids)))
				}
				nodes = append(nodes, n)
				n = shopTestPointer(shopTestWords(n, 4)[2])
			}
		}
		q = shopTestPointer(w[14])
	}
	normalized := func(q unsafe.Pointer, n int) []uint32 {
		w := append([]uint32(nil), shopTestWords(q, n)...)
		for i := range w {
			w[i] = p.normalize(w[i])
		}
		return w
	}
	for _, q := range sessions {
		r.Sessions = append(r.Sessions, normalized(q, 16))
	}
	for _, n := range nodes {
		r.Nodes = append(r.Nodes, normalized(n, 4))
	}
	for _, o := range p.owned {
		if o.alive {
			w := append([]uint32{p.normalize(uint32(uintptr(o.u.CObj())))}, normalized(o.u.CObj(), 193)...)
			if o.init != nil {
				w = append(w, normalized(o.init, 1)...)
			}
			r.Objects = append(r.Objects, w)
			if p.proxy.callbacks.shop.spec.CaptureData {
				d := []uint32{p.normalize(uint32(uintptr(o.u.CObj()))), uint32(o.initSize), uint32(o.useSize)}
				if o.init != nil {
					init := normalized(o.init, o.initSize/4)
					if o.undefinedTail && len(init) > 4 {
						init[4] = 0
					}
					d = append(d, init...)
				}
				if o.u.UseData.Ptr != nil {
					d = append(d, normalized(o.u.UseData.Ptr, o.useSize/4)...)
				}
				if o.health != nil {
					d = append(d, normalized(o.health, 2)...)
				}
				r.ObjectData = append(r.ObjectData, d)
			}
		}
	}
	if len(p.owned) > 64 {
		b, _ := json.Marshal(r.Objects)
		h := sha256.Sum256(b)
		r.ObjectDigest, r.Objects = hex.EncodeToString(h[:]), nil
	}
	if p.players != nil {
		r.Players = p.players()
		if p.inventory != nil {
			for i := range p.proxy.life.players {
				if *(*uint32)(unsafe.Add(p.proxy.life.players[i].CObj(), 772)) != p.inventory.playerHandle {
					panic("inventory player server association changed")
				}
				// Runtime server registration is an identity, not gameplay state.
				r.Players[3*i][193] = 55004
			}
		}
	}
	if p.proxy.callbacks.shop.spec.ProtectedGold {
		var intact bool
		r.Protection, intact = p.snapshotProtection()
		if !intact {
			panic("shop protection fixture guards")
		}
	}
	for q := uint32(C.dword_5d4594_1565512); q != 0; {
		b := unsafe.Slice((*byte)(shopTestPointer(q)), 416)
		packet := PortTestShopPacketResult{Recipient: b[250], Ordered: b[184], A4: binary.LittleEndian.Uint32(b[404:]), A5: binary.LittleEndian.Uint32(b[180:]), Data: bytes.Clone(b[251 : 251+int(b[401])])}
		for i := range packet.Sequence {
			packet.Sequence[i] = binary.LittleEndian.Uint16(b[186+2*i:])
		}
		if p.proxy.callbacks.shop.spec.Engine != nil {
			portTestTradePacketDefined(packet.Data)
		}
		r.Packets = append(r.Packets, packet)
		q = binary.LittleEndian.Uint32(b[408:])
	}
	if p.proxy.callbacks.shop.spec.Engine != nil {
		r.EngineStateRequests = append([]uint32(nil), p.engineStateRequests...)
		cache := portTestTradeCache()
		r.EngineState = append(r.EngineState, cache[:]...)
		r.EngineState = append(r.EngineState, p.enginePickupTrace()...)
		for _, ind := range []ntype.PlayerInd{1, 7, 31} {
			r.EngineMessages = append(r.EngineMessages, p.proxy.core.NetList.CopyPacketsA(ind, 1))
		}
		r.EngineState = append(r.EngineState, normalized(p.proxy.callbacks.shop.npc().InitData, 431)...)
	}
	r.Return, r.Head = p.normalize(rv), p.normalize(uint32(C.dword_5d4594_2386500))
	for i := range r.Cached {
		r.Cached[i] = p.normalize(*memmap.PtrUint32(0x5D4594, 2386364+uintptr(4*i)))
	}
	return r
}
