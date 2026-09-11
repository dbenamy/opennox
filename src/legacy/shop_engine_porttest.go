//go:build porttest

package legacy

/*
#include <string.h>
#include "GAME4_1.h"
#include "server__system__trade.h"
int sub_50E7A0(uint32_t* a1, int a2);
uint32_t* nox_xxx_createPlayerShopSession_50E8F0(int a1, int a2);
int sub_50F0F0(int a1, int a2);
int sub_50F1A0(int a1, int a2);
int nox_xxx_servSendShopItems_50F280(int a1, int a2);
uint32_t* nox_xxx_tradeSetPlayer_50F370(uint32_t* a1, int a2);
int sub_50FAE0(int a1, int a2, int a3, int a4, int a5);
int sub_50FD60(uint32_t* a1, int a2);
void sub_510320(int a1, int a2);
int sub_5104F0(int a1, short a2);
int sub_510540(int a1);
int sub_5105D0(int a1);
extern uint32_t dword_5d4594_2488728;
static uint32_t portTestTradePickupTrace[2049];
static uint32_t* portTestTradePickupData(void) {return portTestTradePickupTrace;}
static int portTestTradePickup(int unit, int item, int a3, int a4) {
 unsigned int i=1+4*portTestTradePickupTrace[0]++;
 if(i+3<2049) {portTestTradePickupTrace[i]=unit;portTestTradePickupTrace[i+1]=item;
 portTestTradePickupTrace[i+2]=a3;portTestTradePickupTrace[i+3]=a4;}
 return 1;
}
static void* portTestTradePickupPtr(void) {return portTestTradePickup;}
extern uint32_t dword_5d4594_2386548;
extern uint32_t dword_5d4594_2386552;
extern uint32_t dword_5d4594_2386560;
static int portTestTradeOffer(int session, int unit, void* item) {
 float bits;
 memcpy(&bits, &item, 4);
 return nox_xxx_tradeP2PAddOffer2_50F820_trade(session, unit, bits);
}
*/
import "C"

import (
	"encoding/binary"
	"unsafe"

	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

type PortTestShopEngineSpec struct {
	Cache      [12]uint32
	Names      [2]string
	VendorName string
	ExtraLives uint32
	NoSellType uint32
}

const (
	PortTestTradeRemove = 100 + iota
	PortTestTradeCreatePlayer
	PortTestTradeStart
	PortTestTradeIntro
	PortTestTradePeerIntro
	PortTestTradeSendStock
	PortTestTradeSetPlayer
	PortTestTradeOfferPacket
	PortTestTradeOfferAllowed
	PortTestTradeStockDecrement
	PortTestTradeShortfall
	PortTestTradeStockRemovable
	PortTestTradeIsGem
	PortTestTradeAddOffer
	PortTestTradeBuy
	PortTestTradeBuyMany
	PortTestTradeSellQuote
	PortTestTradeSell
)

func portTestTradeCache() [12]uint32 {
	var out [12]uint32
	for i := range out {
		out[i] = *memmap.PtrUint32(0x5D4594, 2386516+uintptr(4*i))
	}
	out[8], out[9], out[11] = uint32(C.dword_5d4594_2386548), uint32(C.dword_5d4594_2386552), uint32(C.dword_5d4594_2386560)
	return out
}
func portTestTradeSetCache(in [12]uint32) {
	for i, v := range in {
		*memmap.PtrUint32(0x5D4594, 2386516+uintptr(4*i)) = v
	}
	C.dword_5d4594_2386548, C.dword_5d4594_2386552, C.dword_5d4594_2386560 = C.uint32_t(in[8]), C.uint32_t(in[9]), C.uint32_t(in[11])
}
func (p *portTestShopPools) enginePrepare() func() {
	sp := p.proxy.callbacks.shop.spec.Engine
	if sp == nil {
		return func() {}
	}
	restoreStrings := p.proxy.core.PortTestTradeStrings()
	oldPlayerState := Nox_xxx_playerSetState_4FA020
	p.engineStateRequests = nil
	// Player animation is a retained dependency. Capture its requested state;
	// the real freeze/status services still execute, without a full game loop.
	Nox_xxx_playerSetState_4FA020 = func(u *server.Object, state server.PlayerState) bool {
		p.engineStateRequests = append(p.engineStateRequests, p.normalize(uint32(uintptr(u.CObj()))), uint32(state))
		return true
	}
	freeze := memmap.PtrUint32(0x5D4594, 1567712)
	oldFreeze := *freeze
	*freeze = 0
	old := portTestTradeCache()
	oldDropFlag := C.dword_5d4594_2488728
	drop := unsafe.Slice(memmap.PtrUint32(0x587000, 279432), 6)
	oldDrop := append([]uint32(nil), drop...)
	clear(drop)
	if sp.NoSellType != 0 {
		drop[0] = 1
		drop[1] = sp.NoSellType
	}
	C.dword_5d4594_2488728 = 1
	clear(unsafe.Slice((*uint32)(unsafe.Pointer(C.portTestTradePickupData())), 2049))
	portTestTradeSetCache(sp.Cache)
	for i, u := range p.proxy.life.players[:2] {
		alloc.StrCopy16(u.UpdateDataPlayer().Player.NameFinal[:], sp.Names[i])
		u.UpdateDataPlayer().Field80 = sp.ExtraLives
	}
	alloc.StrCopy16(p.proxy.life.players[2].UpdateDataPlayer().Player.NameFinal[:], "Third")
	s := p.proxy.callbacks.shop
	name := unsafe.Pointer(alloc.InternCString("PortTradeVendor"))
	*(*unsafe.Pointer)(s.npc().CObj()) = name
	p.proxy.life.ids[uint32(uintptr(name))] = 50300
	scratch := unsafe.Slice((*byte)(memmap.PtrOff(0x5D4594, 1563460)), 512)
	oldScratch := append([]byte(nil), scratch...)
	alloc.StrCopy(unsafe.Slice((*byte)(unsafe.Add(s.npc().InitData, 1684)), 32), sp.VendorName)
	return func() {
		restoreStrings()
		Nox_xxx_playerSetState_4FA020 = oldPlayerState
		*freeze = oldFreeze
		copy(scratch, oldScratch)
		portTestTradeSetCache(old)
		copy(drop, oldDrop)
		C.dword_5d4594_2488728 = oldDropFlag
	}
}
func (p *portTestShopPools) engineAdopt(q unsafe.Pointer) uint32 {
	known := false
	for _, old := range p.sessions {
		if old == q && q != nil {
			known = true
			break
		}
	}
	p.sessions = append(p.sessions, q)
	if q == nil {
		return 0
	}
	if !known {
		id := uint32(60000 + len(p.sessions) - 1)
		p.identify(q, id)
		w := shopTestWords(q, 16)
		p.own((*server.Object)(shopTestPointer(w[12])), id+1000)
		p.own((*server.Object)(shopTestPointer(w[13])), id+2000)
	}
	return uint32(uintptr(q))
}
func (p *portTestShopPools) engineAction(a PortTestShopAction, q unsafe.Pointer) uint32 {
	var w []uint32
	if q != nil {
		w = shopTestWords(q, 16)
	}
	var u *server.Object
	if q != nil {
		u = (*server.Object)(shopTestPointer(w[2+a.Side]))
	}
	var item *server.Object
	if a.Item >= 0 && a.Item < len(p.items) {
		item = p.items[a.Item].u
	}
	session := C.int(uintptr(q))
	unit := C.int(0)
	if u != nil {
		unit = C.int(uintptr(u.CObj()))
	}
	switch a.Op {
	case PortTestTradeRemove:
		return uint32(C.sub_50E7A0((*C.uint32_t)(q), C.int(uintptr(item.CObj()))))
	case PortTestTradeCreatePlayer:
		return p.engineAdopt(unsafe.Pointer(C.nox_xxx_createPlayerShopSession_50E8F0(C.int(uintptr(p.proxy.life.players[0].CObj())), C.int(uintptr(p.proxy.callbacks.shop.npc().CObj())))))
	case PortTestTradeStart:
		left, right := &p.proxy.life.players[0], p.proxy.callbacks.shop.npc()
		if p.proxy.callbacks.shop.spec.Peer {
			right = &p.proxy.life.players[1]
		}
		if a.Value == 1 {
			right = &p.proxy.life.players[2]
		}
		if a.Value == 2 {
			left = &p.proxy.life.players[2]
		}
		if a.Side != 0 {
			left, right = right, left
		}
		return p.engineAdopt(unsafe.Pointer(C.nox_xxx_servShopStart_50EF10_trade(C.int(uintptr(left.CObj())), C.int(uintptr(right.CObj())))))
	case PortTestTradeIntro:
		return uint32(C.sub_50F0F0(unit, session))
	case PortTestTradePeerIntro:
		return uint32(C.sub_50F1A0(unit, session))
	case PortTestTradeSendStock:
		return uint32(C.nox_xxx_servSendShopItems_50F280(unit, session))
	case PortTestTradeSetPlayer:
		return uint32(uintptr(unsafe.Pointer(C.nox_xxx_tradeSetPlayer_50F370((*C.uint32_t)(q), unit))))
	case PortTestTradeOfferPacket:
		return uint32(C.sub_50FAE0(unit, C.int(w[2+a.Value%2]), session, C.int(uintptr(item.CObj())), C.int(item.Worth)))
	case PortTestTradeOfferAllowed:
		return uint32(C.sub_50FD60((*C.uint32_t)(shopTestPointer(w[8+a.Side])), C.int(uintptr(item.CObj()))))
	case PortTestTradeStockDecrement:
		C.sub_510320(C.int(uintptr(item.CObj())), session)
	case PortTestTradeShortfall:
		return uint32(C.sub_5104F0(unit, C.short(a.Value)))
	case PortTestTradeStockRemovable:
		return uint32(C.sub_510540(C.int(uintptr(item.CObj()))))
	case PortTestTradeIsGem:
		return uint32(C.sub_5105D0(C.int(uintptr(item.CObj()))))
	case PortTestTradeAddOffer:
		return uint32(C.portTestTradeOffer(session, unit, item.CObj()))
	case PortTestTradeBuy:
		C.sub_5100C0_trade(unit, (*C.uint32_t)(q), C.int(a.Value))
	case PortTestTradeBuyMany:
		return uint32(uintptr(unsafe.Pointer(C.sub_510640_trade(unit, session, C.int(item.TypeInd), (*C.float)(shopTestPointer(a.Value))))))
	case PortTestTradeSellQuote:
		return uint32(uintptr(unsafe.Pointer(C.sub_5109C0_trade((*C.int)(u.CObj()), session, (*C.uint32_t)(shopTestPointer(a.Value))))))
	case PortTestTradeSell:
		return uint32(uintptr(unsafe.Pointer(C.sub_510BE0_trade((*C.int)(u.CObj()), session, (*C.uint32_t)(shopTestPointer(a.Value))))))
	default:
		panic("trade fixture operation")
	}
	return 0
}

// The original name packets transmit unused stack tails. Exclude only bytes
// after the terminating name, leaving all defined payload bytes and metadata.
func portTestTradePacketDefined(b []byte) {
	if len(b) < 2 {
		return
	}
	switch binary.LittleEndian.Uint16(b) {
	case 3529:
		if len(b) != 86 {
			panic("shop intro packet size")
		}
		for i := 54; i < len(b); i++ {
			if b[i] == 0 {
				clear(b[i+1:])
				break
			}
		}
	case 3273:
		if len(b) != 52 {
			panic("peer intro packet size")
		}
		for i := 2; i+1 < len(b); i += 2 {
			if binary.LittleEndian.Uint16(b[i:]) == 0 {
				clear(b[i+2:])
				break
			}
		}
	}
}

func (p *portTestShopPools) engineDiscover() {
	if p.proxy.callbacks.shop.spec.Engine == nil {
		return
	}
	add := func(u *server.Object) {
		for _, o := range p.owned {
			if o.alive && o.u == u {
				return
			}
		}
		p.own(u, uint32(76000+len(p.owned)))
	}
	for i := 0; i < 2; i++ {
		for u := p.proxy.life.players[i].InvFirstItem; u != nil; u = u.InvNextItem {
			add(u)
		}
	}
	for _, q := range p.sessions {
		if q == nil {
			continue
		}
		w := shopTestWords(q, 16)
		for n := w[5]; n != 0; {
			v := shopTestWords(shopTestPointer(n), 4)
			add((*server.Object)(shopTestPointer(v[0])))
			n = v[2]
		}
	}
}
func (p *portTestShopPools) engineItem(u *server.Object, sp PortTestShopItem) {
	if p.proxy.callbacks.shop.spec.Engine == nil {
		return
	}
	if sp.Pickup {
		u.Pickup.Ptr = C.portTestTradePickupPtr()
		p.identify(u.Pickup.Ptr, 53000)
	}
	for i, enabled := range sp.Mods {
		if enabled {
			*(*unsafe.Pointer)(unsafe.Add(u.InitData, 4*i)) = unsafe.Add(p.proxy.callbacks.shop.ptr(8), i*int(unsafe.Sizeof(server.ModifierEff{})))
		}
	}
}
func (p *portTestShopPools) enginePickupTrace() []uint32 {
	raw := unsafe.Slice((*uint32)(unsafe.Pointer(C.portTestTradePickupData())), 2049)
	if raw[0] > 512 {
		panic("trade pickup trace capacity")
	}
	trace := append([]uint32(nil), raw[:1+4*raw[0]]...)
	for i := range trace {
		trace[i] = p.normalize(trace[i])
	}
	return trace
}
