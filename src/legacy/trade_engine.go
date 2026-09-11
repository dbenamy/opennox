package legacy

/*
#include "GAME1.h"
#include "GAME3_2.h"
#include "GAME3_3.h"
#include "GAME4.h"
#include "GAME4_3.h"
#include "server__dbase__objdb.h"
#include "common__strman.h"
extern uint32_t dword_5d4594_2386548;
extern uint32_t dword_5d4594_2386552;
extern uint32_t dword_5d4594_2386560;
extern void* nox_alloc_tradeItems_2386496;
static int tradeSendLine(int unit, wchar2_t* text) {
 return nox_xxx_netSendLineMessage_4D9EB0(unit, text);
}
static void tradeFormatName(wchar2_t* dst, wchar2_t* format, wchar2_t* name) {
 nox_swprintf(dst, format, name);
}
*/
import "C"

import (
	"encoding/binary"
	"unsafe"

	"github.com/opennox/libs/strman"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

const tradeSource = `C:\NoxPost\src\Server\System\Trade.c`

func tradeString(key string) *C.wchar2_t {
	return (*C.wchar2_t)(unsafe.Pointer(internWStr(GetServer().S().Strings().GetStringInFile(strman.ID(key), tradeSource))))
}
func tradeLine(u *server.Object, key string) uint32 {
	return uint32(C.tradeSendLine(C.int(uintptr(u.CObj())), tradeString(key)))
}
func tradeNamedLine(u *server.Object, key string, name *uint16) {
	var text [128]uint16
	C.tradeFormatName((*C.wchar2_t)(unsafe.Pointer(&text[0])), tradeString(key), (*C.wchar2_t)(unsafe.Pointer(name)))
	C.tradeSendLine(C.int(uintptr(u.CObj())), (*C.wchar2_t)(unsafe.Pointer(&text[0])))
}
func tradeRemoveStock(s *shopSession, u *server.Object) uint32 {
	for n := s.Stock; n != nil; n = n.Next {
		if n.Object != u {
			continue
		}
		if n.Next != nil {
			n.Next.Prev = n.Prev
		}
		if n.Prev != nil {
			n.Prev.Next = n.Next
		}
		if s.Stock == n {
			s.Stock = n.Next
		}
		player := s.Units[0]
		if player.ObjClass&4 == 0 {
			player = s.Units[1]
		}
		shopSendCode(player, u, 2505)
		shopFreeItem(n)
		return 1
	}
	return 0
}
func tradeCreatePlayer(player, vendor *server.Object) *shopSession {
	var s *shopSession
	if noxflags.HasGame(noxflags.GameModeQuest) {
		ind := uint8(player.UpdateDataPlayer().Player.PlayerInd)
		s = shopCached()[ind]
		shopCached()[ind] = nil
	}
	cached := s != nil
	if !cached {
		s = shopCreate()
		if s == nil {
			return nil
		}
	}
	s.Units = [2]*server.Object{player, vendor}
	s.Kind = 1
	if !cached {
		shopLoad(s)
	}
	return s
}
func tradeCopyWide(b []byte, src *uint16, limit int) {
	for i := 0; i < limit; i++ {
		v := *(*uint16)(unsafe.Add(unsafe.Pointer(src), 2*i))
		if v == 0 {
			break
		}
		binary.LittleEndian.PutUint16(b[2*i:], v)
	}
}
func tradeIntro(s *shopSession) uint32 {
	vendor, player := s.Units[0], s.Units[1]
	if vendor.ObjClass&4 != 0 {
		vendor, player = player, vendor
	}
	var b [86]byte
	b[0], b[1] = 0xc9, 13
	binary.LittleEndian.PutUint16(b[2:], vendor.TypeInd)
	tradeCopyWide(b[4:54], (*uint16)(unsafe.Pointer(C.sub_4E39F0_obj_db(asObjectC(vendor)))), 24)
	copy(b[54:], alloc.GoString((*byte)(unsafe.Add(vendor.InitData, 1684))))
	return shopSend(player, b[:], 1)
}
func tradePeerIntro(u *server.Object, s *shopSession) uint32 {
	var name *uint16
	other := s.Units[0]
	if u == other {
		other = s.Units[1]
	}
	if other.ObjClass&4 != 0 {
		name = &other.UpdateDataPlayer().Player.NameFinal[0]
	} else {
		// The original alternate-side NPC path also names Units[1].
		name = (*uint16)(unsafe.Pointer(C.sub_4E39F0_obj_db(asObjectC(s.Units[1]))))
	}
	var b [52]byte
	b[0], b[1] = 0xc9, 12
	tradeCopyWide(b[2:], name, 24)
	return shopSend(u, b[:], 1)
}
func tradeSendStock(u *server.Object, s *shopSession) uint32 {
	result := uint32(uintptr(unsafe.Pointer(s)))
	for n := s.Stock; n != nil; n = n.Next {
		result = shopSendItem(u, n)
	}
	return result
}
func tradeSetPlayer(s *shopSession, u *server.Object) *shopSession {
	s.Active = 1
	if u.ObjClass&4 != 0 {
		u.UpdateDataPlayer().Trade70 = s
	}
	return s
}
func tradeOfferPacket(u, owner, item *server.Object, cost uint32) uint32 {
	var b [15]byte
	b[0], b[1] = 0xc9, 4
	if u == owner {
		b[2] = 1
	}
	binary.LittleEndian.PutUint16(b[3:], item.TypeInd)
	binary.LittleEndian.PutUint16(b[5:], uint16(item.NetCode))
	binary.LittleEndian.PutUint32(b[7:], cost)
	for i := 11; i < 15; i++ {
		b[i] = 0xff
	}
	if uint32(item.ObjClass)&0x13001000 != 0 {
		for i, p := range unsafe.Slice((*unsafe.Pointer)(item.InitData), 4) {
			if p != nil {
				b[11+i] = *(*byte)(unsafe.Add(p, 4))
			}
		}
	}
	return shopSend(u, b[:], 1)
}
func tradeOfferAllowed(first *shopItem, item *server.Object) bool {
	var kinds [4]uint32
	count := 0
	for n := first; n != nil; n = n.Next {
		typ := uint32(n.Object.TypeInd)
		found := false
		for _, v := range kinds {
			if v == typ {
				found = true
				break
			}
		}
		if !found {
			kinds[count] = typ
			count++
		}
	}
	if count == 0 {
		return true
	}
	for _, v := range kinds[:count] {
		if v == uint32(item.TypeInd) && uint32(item.ObjClass)&0x13001000 == 0 {
			return true
		}
	}
	return count < 4
}
func tradeStockDecrement(item *server.Object, s *shopSession) {
	if item == nil || s == nil {
		return
	}
	data := shopVendor(s).InitData
	index := shopStockIndex(s, item)
	if index < 0 {
		return
	}
	count := int(*(*byte)(data))
	entries := unsafe.Slice((*shopStockEntry)(unsafe.Add(data, 4)), count)
	entries[index].Count--
	if entries[index].Count == 0 {
		copy(entries[index:], entries[index+1:])
		*(*byte)(data)--
	}
}
func tradeShortfall(u *server.Object, amount uint32) uint32 {
	var b [4]byte
	b[0], b[1] = 0xc9, 27
	binary.LittleEndian.PutUint16(b[2:], uint16(amount))
	return shopSend(u, b[:], 0)
}
func tradeStockRemovable(item *server.Object) bool {
	if !noxflags.HasGame(noxflags.GameModeQuest) {
		return true
	}
	cache := unsafe.Slice(memmap.PtrUint32(0x5D4594, 2386520), 4)
	if cache[0] == 0 {
		for i, name := range []string{"Diamond", "Ruby", "Emerald", "AnkhTradable"} {
			cache[i] = uint32(GetServer().S().Types.IndByID(name))
		}
	}
	for _, typ := range cache {
		if uint32(item.TypeInd) == typ {
			return false
		}
	}
	return true
}
func tradeIsGem(item *server.Object) bool {
	cache := unsafe.Slice(memmap.PtrUint32(0x5D4594, 2386536), 3)
	if cache[0] == 0 {
		for i, name := range []string{"Diamond", "Ruby", "Emerald"} {
			cache[i] = uint32(GetServer().S().Types.IndByID(name))
		}
	}
	for _, typ := range cache {
		if uint32(item.TypeInd) == typ {
			return true
		}
	}
	return false
}
func tradeStart(left, right *server.Object) *shopSession {
	lp, rp := left.ObjClass&4 != 0, right.ObjClass&4 != 0
	if lp {
		if old := left.UpdateDataPlayer().Trade70; old != nil {
			if old.Units[1] != right {
				tradeLine(left, "StarterAlreadyTrading")
			}
			return nil
		}
	}
	if rp && right.UpdateDataPlayer().Trade70 != nil {
		if lp {
			tradeNamedLine(left, "OtherAlreadyTrading", &right.UpdateDataPlayer().Player.NameFinal[0])
		}
		return nil
	}
	var s *shopSession
	if lp && !rp {
		s = tradeCreatePlayer(left, right)
	} else if !lp && rp {
		s = tradeCreatePlayer(right, left)
	} else {
		s = shopCreate()
		if s != nil {
			s.Units = [2]*server.Object{left, right}
		}
	}
	if s == nil {
		return nil
	}
	s.Field4 = GetServer().S().Frame()
	tradeSetPlayer(s, left)
	tradeSetPlayer(s, right)
	if s.Kind != 0 {
		if lp {
			tradeIntro(s)
		}
		if rp {
			tradeIntro(s)
		}
	} else {
		if lp {
			tradePeerIntro(left, s)
		}
		if rp {
			tradePeerIntro(right, s)
		}
	}
	if s.Kind != 0 {
		player := s.Units[0]
		if player.ObjClass&4 == 0 {
			player = s.Units[1]
		}
		tradeSendStock(player, s)
		if noxflags.HasGame(2048) {
			C.nox_xxx_unitFreeze_4E79C0(asObjectC(player), 0)
		}
	}
	return s
}
func tradeAddOffer(s *shopSession, u, item *server.Object) uint32 {
	gold := memmap.PtrUint32(0x5D4594, 2386516)
	if *gold == 0 {
		*gold = uint32(GetServer().S().Types.IndByID("Gold"))
	}
	side := 0
	if s.Units[0] != u {
		if s.Units[1] != u {
			return 0
		}
		side = 1
	}
	if !tradeOfferAllowed(s.Offers[side], item) {
		return 0
	}
	n := (*shopItem)(alloc.AsClass(C.nox_alloc_tradeItems_2386496).NewObject())
	if n == nil {
		for _, u := range s.Units {
			if u.ObjClass&4 != 0 {
				tradeLine(u, "TradeMaxObjectsReached")
			}
		}
		return 0
	}
	n.Object, n.Value = item, uint32(shopPrice(1, s, item))
	n.Next = s.Offers[side]
	if n.Next != nil {
		n.Next.Prev = n
	}
	s.Offers[side] = n
	shopBalance(s)
	s.Accepted = [2]uint32{}
	if s.Units[0].ObjClass&4 != 0 {
		if s.Units[1].ObjClass&4 == 0 && s.Totals[1] <= s.Totals[0] {
			s.Accepted[1] = 1
		}
	} else if s.Totals[0] <= s.Totals[1] {
		s.Accepted[0] = 1
	}
	for _, recipient := range s.Units {
		if recipient.ObjClass&4 != 0 {
			shopSendGold(recipient, s)
			tradeOfferPacket(recipient, u, item, n.Value)
			shopSendAcceptance(recipient, s)
		}
	}
	return 1
}
func tradeInventoryCount(u, item *server.Object) int32 {
	return int32(C.nox_xxx_inventoryCountObjects_4E7D30(C.int(uintptr(u.CObj())), C.int(item.TypeInd)))
}
func tradeLimit(key string) uint32 {
	return uint32(floatToInt32(float32(GetServer().S().Balance.Float(key))))
}
func tradePriority(u *server.Object, key string) {
	C.nox_xxx_netPriMsgToPlayer_4DA2C0(asObjectC(u), internCStr(key), 0)
}
func tradeReportGold(u *server.Object) {
	C.sub_4D8870(C.int(uint8(u.UpdateDataPlayer().Player.PlayerInd)), C.int(uintptr(u.CObj())))
}
func tradeBuy(u *server.Object, s *shopSession, code uint32) {
	available := shopGetGold(u)
	if C.dword_5d4594_2386548 == 0 {
		C.dword_5d4594_2386548 = C.uint32_t(GetServer().S().Types.IndByID("AnkhTradable"))
	}
	var item *server.Object
	for n := s.Stock; n != nil; n = n.Next {
		if n.Object != nil && n.Object.NetCode == code {
			item = n.Object
			break
		}
	}
	if item == nil {
		return
	}
	cost := uint32(shopPrice(1, s, item))
	if cost > available {
		tradeShortfall(u, cost-available)
		return
	}
	if item.ObjClass&0x10 != 0 {
		limit := int32(3)
		if noxflags.HasGame(6144) {
			limit = 9
		}
		if tradeInventoryCount(u, item) >= limit {
			tradeLine(u, "pickup.c:MaxSameItem")
			return
		}
	}
	if uint32(item.TypeInd) == uint32(C.dword_5d4594_2386548) && u.UpdateDataPlayer().Field80 >= tradeLimit("MaxExtraLives") {
		tradePriority(u, "pickup.c:MaxTradableAnkhsReached")
		GetServer().S().Audio.EventObj(925, u, 0, 0)
		return
	}
	quest := noxflags.HasGame(noxflags.GameModeQuest)
	if quest && item.ObjClass&0x1000 != 0 && item.ObjSubClass&0x200000 != 0 && tradeInventoryCount(u, item) >= int32(tradeLimit("ForceOfNatureStaffLimit")) {
		tradePriority(u, "pickup.c:MaxSameItem")
		GetServer().S().Audio.EventObj(925, u, 0, 0)
		return
	}
	delivered := item
	if quest && (tradeIsGem(item) || uint32(item.TypeInd) == uint32(C.dword_5d4594_2386548)) {
		delivered = GetServer().S().NewObjectByTypeInd(int(item.TypeInd))
	}
	if delivered.ObjClass&0x110 != 0 || delivered.Pickup.Ptr == nil {
		shopPut(u, delivered)
		GetServer().S().Audio.EventObj(307, u, 2, u.NetCode)
	} else {
		delivered.CallPickup(u, 1, 1)
	}
	tradeStockDecrement(item, s)
	if tradeStockRemovable(item) {
		tradeRemoveStock(s, item)
	}
	shopSubGold(u, cost)
	tradeReportGold(u)
}
func tradeBuyMany(u *server.Object, s *shopSession, typ int32, count uint32) uint32 {
	available := shopGetGold(u)
	if C.dword_5d4594_2386552 == 0 {
		C.dword_5d4594_2386552 = C.uint32_t(GetServer().S().Types.IndByID("AnkhTradable"))
	}
	for bought := uint32(0); bought < count; {
		n := s.Stock
		for n != nil && int32(n.Object.TypeInd) != typ {
			n = n.Next
		}
		if n == nil {
			return 0
		}
		item := n.Object
		if item == nil {
			return uint32(uintptr(unsafe.Pointer(n)))
		}
		cost := uint32(shopPrice(1, s, item))
		if cost > available {
			return tradeShortfall(u, cost-available)
		}
		if item.ObjClass&0x10 != 0 {
			limit := int32(3)
			if noxflags.HasGame(6144) {
				limit = 9
			}
			if tradeInventoryCount(u, item) >= limit {
				return tradeLine(u, "pickup.c:MaxSameItem")
			}
		}
		if uint32(item.TypeInd) == uint32(C.dword_5d4594_2386552) && u.UpdateDataPlayer().Field80 >= tradeLimit("MaxExtraLives") {
			tradePriority(u, "pickup.c:MaxTradableAnkhsReached")
			GetServer().S().Audio.EventObj(925, u, 0, 0)
			return uint32(uintptr(unsafe.Pointer(n)))
		}
		quest := noxflags.HasGame(noxflags.GameModeQuest)
		if quest && item.ObjClass&0x1000 != 0 && item.ObjSubClass&0x200000 != 0 && tradeInventoryCount(u, item) >= int32(tradeLimit("ForceOfNatureStaffLimit")) {
			tradePriority(u, "pickup.c:MaxSameItem")
			GetServer().S().Audio.EventObj(925, u, 0, 0)
			return uint32(uintptr(unsafe.Pointer(n)))
		}
		delivered := item
		if quest && (tradeIsGem(item) || uint32(item.TypeInd) == uint32(C.dword_5d4594_2386552)) {
			delivered = GetServer().S().NewObjectByTypeInd(int(item.TypeInd))
		}
		if delivered.Pickup.Ptr != nil {
			delivered.CallPickup(u, 1, 1)
		} else {
			shopPut(u, delivered)
		}
		tradeStockDecrement(item, s)
		if tradeStockRemovable(item) {
			tradeRemoveStock(s, item)
		}
		shopSubGold(u, cost)
		tradeReportGold(u)
		bought++
		if bought >= count {
			return bought
		}
	}
	return 0
}
func tradeSaleBlocked(u, item *server.Object, glyph uint32) bool {
	var key string
	if C.nox_xxx_ItemIsDroppable_53EBF0(C.int(uintptr(item.CObj()))) == 1 {
		key = "CantSellQuestItem"
	} else if uint32(item.TypeInd) == glyph {
		key = "CantSellItem"
	} else {
		return false
	}
	tradeLine(u, key)
	GetServer().S().Audio.EventObj(925, u, 2, u.NetCode)
	return true
}
func tradeSaleQuote(u *server.Object, s *shopSession, code uint32) uint32 {
	glyph := memmap.PtrUint32(0x5D4594, 2386556)
	if *glyph == 0 {
		*glyph = uint32(GetServer().S().Types.IndByID("Glyph"))
	}
	for item := u.InvFirstItem; item != nil; item = item.InvNextItem {
		if item.NetCode != code {
			continue
		}
		if tradeSaleBlocked(u, item, *glyph) {
			return code
		}
		var b [8]byte
		b[0], b[1] = 0xc9, 29
		binary.LittleEndian.PutUint16(b[2:], uint16(code))
		binary.LittleEndian.PutUint32(b[4:], uint32(shopPrice(0, s, item)))
		return shopSend(u, b[:], 0)
	}
	return code
}
func tradeSell(u *server.Object, s *shopSession, code uint32) uint32 {
	shopGetGold(u)
	if C.dword_5d4594_2386560 == 0 {
		C.dword_5d4594_2386560 = C.uint32_t(GetServer().S().Types.IndByID("Glyph"))
	}
	if u.InvFirstItem == nil {
		return uint32(C.dword_5d4594_2386560)
	}
	for item := u.InvFirstItem; item != nil; item = item.InvNextItem {
		if item.NetCode != code {
			continue
		}
		if tradeSaleBlocked(u, item, uint32(C.dword_5d4594_2386560)) {
			return code
		}
		C.sub_4ED0C0(asObjectC(u), asObjectC(item))
		GetServer().DelayedDelete(item)
		shopAddGold(u, uint32(shopPrice(0, s, item)))
		tradeReportGold(u)
		GetServer().S().Audio.EventObj(307, u, 2, u.NetCode)
		break
	}
	return code
}
