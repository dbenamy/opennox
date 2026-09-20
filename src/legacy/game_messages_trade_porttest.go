//go:build porttest

package legacy

import (
	"bytes"
	"encoding/binary"
	"unsafe"
)

type PortTestServerTradeSpec struct {
	Dispatch                         bool
	NoSession, Found, Saving, Vendor bool
	Flags                            uint32
	Subtype                          byte
	Code                             uint16
	Count                            byte
}

func (p *portTestShopPools) portTestServerTrade(a PortTestShopAction, session *shopSession) uint32 {
	sp := a.ServerMessage
	u := session.Units[a.Side]
	if sp.NoSession {
		old := u.UpdateDataPlayer().Trade70
		u.UpdateDataPlayer().Trade70 = nil
		defer func() { u.UpdateDataPlayer().Trade70 = old }()
	}
	if sp.Subtype == 15 || sp.Subtype == 16 {
		target := p.items[0].u
		oldList, oldPending := p.proxy.core.Objs.List, p.proxy.core.Objs.Pending
		oldNext, oldExtent := target.ObjNext, target.Extent
		p.proxy.core.Objs.List, p.proxy.core.Objs.Pending = target, nil
		target.ObjNext, target.Extent = nil, 23
		restore := PortTestCreatureXferLookupOwner()
		defer func() {
			restore()
			target.ObjNext, target.Extent = oldNext, oldExtent
			p.proxy.core.Objs.List, p.proxy.core.Objs.Pending = oldList, oldPending
		}()
	}
	gold := session.Gold
	var stock []uint32
	for n := session.Stock; n != nil; n = n.Next {
		stock = append(stock, uint32(uintptr(n.Object.CObj())))
	}
	if sp.Dispatch {
		length := 4
		switch sp.Subtype {
		case 14, 17, 18:
			length = 2
		case 23, 25:
			length = 5
		}
		data := make([]byte, length)
		data[0], data[1] = 201, sp.Subtype
		if length >= 4 {
			binary.LittleEndian.PutUint16(data[2:], sp.Code)
		}
		if length == 5 {
			data[4] = sp.Count
		}
		before := bytes.Clone(data)
		pl := u.UpdateDataPlayer().Player
		if n := Nox_xxx_netOnPacketRecvServ_51BAD0_net_sdecode_switch(pl.PlayerIndex(), data, pl, u, u.UpdateData); n != length || !bytes.Equal(data, before) {
			panic("trade message length/input")
		}
	} else if !sp.NoSession {
		switch sp.Subtype {
		case 14:
			shopCancelTrade(session)
		case 15:
			if sp.Found && tradeAddOffer(session, u, p.items[a.Item].u) == 1 {
				inventoryRemove(u, p.items[a.Item].u)
			}
		case 16:
			code := uint32(sp.Code)
			if sp.Code == 0x8017 {
				code = p.items[0].u.NetCode
			}
			shopWithdraw(session, code)
		case 17:
			shopAccept(session, u)
		case 18:
			shopExit(session)
		case 22:
			tradeBuy(u, session, uint32(sp.Code))
		case 23:
			tradeBuyMany(u, session, int32(sp.Code), uint32(sp.Count))
		case 24:
			tradeSell(u, session, uint32(sp.Code))
		case 25:
			shopSell(u, session, int32(sp.Code), uint32(sp.Count))
		case 26:
			shopRepair(u, session, uint32(sp.Code))
		case 28:
			tradeSaleQuote(u, session, uint32(sp.Code))
		case 30:
			shopRepairQuote(u, session, uint32(sp.Code))
		default:
			panic("unconfigured direct trade contract")
		}
	}
	// Match the existing fixture's ownership accounting after real destruction.
	live := false
	for q := shopHead(); q != nil; q = q.Next {
		if q == session {
			live = true
		}
	}
	if !live {
		for _, g := range gold {
			p.markFreed(uint32(uintptr(g.CObj())))
		}
		for _, v := range stock {
			p.markFreed(v)
		}
		p.sessions[a.Session] = nil
	}
	return 0
}

func (p *portTestShopPools) portTestServerTradeOpening(a PortTestShopAction) uint32 {
	sp := a.ServerMessage
	u, target := &p.proxy.life.players[0], p.proxy.callbacks.shop.npc()
	pl := u.UpdateDataPlayer().Player
	oldFlags, oldSaving := pl.Field3680, Nox_xxx_gameGet_4DB1B0
	pl.Field3680 = sp.Flags
	Nox_xxx_gameGet_4DB1B0 = func() bool { return sp.Saving }
	oldList, oldPending := p.proxy.core.Objs.List, p.proxy.core.Objs.Pending
	oldNext, oldCode, oldExtent, oldSubclass := target.ObjNext, target.NetCode, target.Extent, target.ObjSubClass
	p.proxy.core.Objs.List, p.proxy.core.Objs.Pending = target, nil
	target.ObjNext, target.NetCode, target.Extent = nil, 7, 23
	target.ObjSubClass &^= 8
	if sp.Vendor {
		target.ObjSubClass |= 8
	}
	restore := PortTestCreatureXferLookupOwner()
	defer func() {
		restore()
		if pl.Field3680 != sp.Flags {
			panic("trade opening changed admission flags")
		}
		pl.Field3680, Nox_xxx_gameGet_4DB1B0 = oldFlags, oldSaving
		p.proxy.core.Objs.List, p.proxy.core.Objs.Pending = oldList, oldPending
		target.ObjNext, target.NetCode, target.Extent, target.ObjSubClass = oldNext, oldCode, oldExtent, oldSubclass
	}()
	if sp.Dispatch {
		data := []byte{201, 21, 0, 0}
		binary.LittleEndian.PutUint16(data[2:], sp.Code)
		before := bytes.Clone(data)
		if n := Nox_xxx_netOnPacketRecvServ_51BAD0_net_sdecode_switch(pl.PlayerIndex(), data, pl, u, u.UpdateData); n != 4 || !bytes.Equal(data, before) {
			panic("trade opening length/input")
		}
	} else if sp.Found && sp.Vendor && !sp.Saving && sp.Flags&3 == 0 {
		tradeStart(u, target)
	}
	q := u.UpdateDataPlayer().Trade70
	want := sp.Found && sp.Vendor && !sp.Saving && sp.Flags&3 == 0
	if (q != nil) != want {
		panic("trade opening admission differs from contract")
	}
	if q != nil {
		p.engineAdopt(unsafe.Pointer(q))
	}
	return 0
}
