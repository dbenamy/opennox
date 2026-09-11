//go:build porttest

package legacy

/*
#include <string.h>
#include "GAME4_1.h"
extern void* nox_alloc_tradeSession_2386492;
extern void* nox_alloc_tradeItems_2386496;
extern uint32_t dword_5d4594_2386500;
extern uint32_t dword_5d4594_1565512;
static void* portTestShopAdd(int session, void* obj) {
	float bits;
	memcpy(&bits, &obj, sizeof(bits));
	return nox_xxx_addItemToShopSession_50EE00(session, bits);
}
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
	u     *server.Object
	init  unsafe.Pointer
	alive bool
}
type portTestShopPools struct {
	proxy              *portTestRoamOwnerServer
	restore            func()
	sessions           []unsafe.Pointer
	owned              []*portTestShopOwned
	items              []*portTestShopOwned
	ids                map[uint32]uint32
	steps              []PortTestShopStep
	initialAlive       int
	players            func() [][]uint32
	freePlayers        func()
	prepareProtection  func(uint32)
	snapshotProtection func() ([]uint32, bool)
}

func portTestShopPoolsEnvironment(proxy *portTestRoamOwnerServer) *portTestShopPools {
	p := &portTestShopPools{proxy: proxy}
	prepareProtection, snapshotProtection, freeProtection := portTestPenaltyProtectionEnvironment()
	p.prepareProtection, p.snapshotProtection = prepareProtection, snapshotProtection
	oldSessions, oldItems, oldHead := C.nox_alloc_tradeSession_2386492, C.nox_alloc_tradeItems_2386496, C.dword_5d4594_2386500
	table := unsafe.Slice(memmap.PtrUint32(0x5D4594, 2386364), 32)
	oldTable := append([]uint32(nil), table...)
	C.nox_alloc_tradeSession_2386492, C.nox_alloc_tradeItems_2386496, C.dword_5d4594_2386500 = nil, nil, 0
	if C.nox_xxx_registerShopClasses_50E2A0() == 0 {
		panic("shop fixture pools")
	}
	p.restore = func() {
		p.cleanup()
		freeProtection()
		C.nox_xxx_deleteShopInventories_50E300()
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
	p.owned = append(p.owned, o)
	p.identify(u.CObj(), id)
	if u.InitData != nil {
		p.identify(u.InitData, id+10000)
	}
	return o
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
	C.sub_50E360()
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
	}
	if p.proxy.core.Objs.Alive != p.initialAlive {
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
	p.players, p.freePlayers = portTestSpawnPlayers(p.proxy, []PortTestSpawnPlayer{{Flags: 4}, {Flags: 4}})
	for i := 0; i < 2; i++ {
		pl := p.proxy.life.players[i].UpdateDataPlayer().Player
		pl.GoldVal, pl.ProtPlayerGold = s.spec.Gold[i], 0
	}
	p.prepareProtection(s.spec.Gold[0])
	if s.spec.ProtectedGold {
		p.proxy.life.players[0].UpdateDataPlayer().Player.ProtPlayerGold = 0x40000001
	}
	for i, spec := range s.spec.Items {
		// Actual allocator objects let destruction execute the retained object
		// free path. Price/charge/modifier arithmetic has guarded query fixtures.
		u := p.proxy.core.NewObjectByTypeInd(23)
		u.TypeInd, u.ObjClass, u.ObjSubClass, u.ObjFlags = spec.Type, object.Class(spec.Class), object.SubClass(spec.Subclass), object.Flags(spec.Flags)
		u.Worth, u.NetCode = spec.Worth, uint32(70000+i)
		init, _ := alloc.Make([]byte{}, 20)
		u.InitData = unsafe.Pointer(&init[0])
		p.items = append(p.items, p.own(u, uint32(70000+i)))
	}
	for _, a := range s.spec.Sequence {
		var q unsafe.Pointer
		if a.Op != PortTestShopCreate && a.Op != PortTestShopReset && a.Op != PortTestShopPlayerCleanup {
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
			q = unsafe.Pointer(C.nox_xxx_createShopStruct_50E870())
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
			C.sub_510000(C.int(uintptr(q)))
			p.sessions[a.Session] = nil
		case PortTestShopReset:
			p.reset()
		case PortTestShopAdd:
			rv = uint32(uintptr(C.portTestShopAdd(C.int(uintptr(q)), p.items[a.Item].u.CObj())))
		case PortTestShopFind:
			w := shopTestWords(q, 16)
			rv = uint32(uintptr(unsafe.Pointer(C.sub_50FFE0((*C.uint32_t)(shopTestPointer(w[5])), C.int(a.Value)))))
		case PortTestShopFreeList:
			w := shopTestWords(q, 16)
			rv = uint32(C.sub_50F6B0(C.int(w[5])))
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
			rv = uint32(C.sub_50FD20((*C.uint32_t)(q), C.int(words[2+a.Side])))
		case PortTestShopBalance:
			rv = uint32(uintptr(unsafe.Pointer(C.sub_50FB90((*C.uint32_t)(q)))))
		case PortTestShopPacket:
			player := C.int(words[2+a.Side])
			item := C.int(uintptr(s.item().CObj()))
			switch a.Item {
			case 0:
				rv = uint32(C.sub_50E820(player, item))
			case 1:
				rv = uint32(C.sub_50F2B0(player, (*C.uint32_t)(s.ptr(9))))
			case 2:
				rv = uint32(C.sub_50F450(player))
			case 3:
				rv = uint32(C.nox_xxx_sendEndTradeMsg_50F560(player))
			case 4:
				rv = uint32(C.sub_50F6E0(player))
			case 5:
				rv = uint32(C.sub_50F720(player, (*C.uint32_t)(q)))
			case 6:
				rv = uint32(C.nox_xxx_tradeP2PUpdStuff_50FA00(player, (*C.uint32_t)(q)))
			case 7:
				rv = uint32(C.sub_50FF90(player, C.int(uintptr(q)), item))
			default:
				panic("shop fixture packet")
			}
		case PortTestShopWithdraw:
			rv = uint32(C.nox_xxx_tradeP2PAddOfferMB_50FE20(C.int(uintptr(q)), C.int(a.Value)))
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
			panic("shop fixture sequence operation")
		}
		p.steps = append(p.steps, p.snapshot(rv))
	}
}

func (p *portTestShopPools) snapshot(rv uint32) PortTestShopStep {
	r := PortTestShopStep{Alive: p.proxy.core.Objs.Alive - p.initialAlive}
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
		}
	}
	if len(p.owned) > 64 {
		b, _ := json.Marshal(r.Objects)
		h := sha256.Sum256(b)
		r.ObjectDigest, r.Objects = hex.EncodeToString(h[:]), nil
	}
	if p.players != nil {
		r.Players = p.players()
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
		r.Packets = append(r.Packets, packet)
		q = binary.LittleEndian.Uint32(b[408:])
	}
	r.Return, r.Head = p.normalize(rv), p.normalize(uint32(C.dword_5d4594_2386500))
	for i := range r.Cached {
		r.Cached[i] = p.normalize(*memmap.PtrUint32(0x5D4594, 2386364+uintptr(4*i)))
	}
	return r
}
