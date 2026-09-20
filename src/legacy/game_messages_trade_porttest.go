//go:build porttest

package legacy

import (
	"bytes"
	"encoding/binary"
)

type PortTestServerTradeSpec struct {
	Dispatch bool
	Subtype  byte
	Code     uint16
	Count    byte
}

func (p *portTestShopPools) portTestServerTrade(a PortTestShopAction, session *shopSession) uint32 {
	sp := a.ServerMessage
	u := session.Units[a.Side]
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
	} else {
		switch sp.Subtype {
		case 14:
			shopCancelTrade(session)
		case 16:
			shopWithdraw(session, uint32(sp.Code))
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
