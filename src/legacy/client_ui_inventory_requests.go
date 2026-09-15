package legacy

/*
#include "defs.h"
#include "GAME2_1.h"
int nox_xxx_netClientSend2_4E53C0(int,const void*,int,int,int);
*/
import "C"

import (
	"encoding/binary"
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/internal/netlist"
	"unsafe"
)

func uiInventoryItemRequest(op byte, dr *client.Drawable) int {
	var msg [3]byte
	msg[0] = op
	binary.LittleEndian.PutUint16(msg[1:], uint16(nox_xxx_netGetUnitCodeCli_578B00(C.int(uiInventoryPointer(unsafe.Pointer(dr))))))
	return bool2int(GetServer().S().NetList.AddToMsgListCli(31, netlist.Kind0, msg[:]))
}
func uiInventoryEquipRequest(dr *client.Drawable) int  { return uiInventoryItemRequest(117, dr) }
func uiInventoryDequipRequest(dr *client.Drawable) int { return uiInventoryItemRequest(118, dr) }
func uiInventoryUse(dr *client.Drawable) {
	if dr != nil {
		uiInventoryItemRequest(116, dr)
	}
}

//export nox_xxx_clientEquip_4623B0
func nox_xxx_clientEquip_4623B0(v C.int) C.int {
	return C.int(uiInventoryEquipRequest(uiInventoryDrawable(uint32(v))))
}

//export nox_xxx_clientDequip_464B70
func nox_xxx_clientDequip_464B70(v C.int) C.int {
	return C.int(uiInventoryDequipRequest(uiInventoryDrawable(uint32(v))))
}

//export nox_xxx_clientUse_465C70
func nox_xxx_clientUse_465C70(v C.int) { uiInventoryUse(uiInventoryDrawable(uint32(v))) }
func uiInventoryTrade(op byte, code uint16) int {
	msg := [4]byte{201, op, byte(code), byte(code >> 8)}
	return bool2int(GetServer().S().NetList.AddToMsgListCli(31, netlist.Kind0, msg[:]))
}

//export nox_xxx_trade_4657B0
func nox_xxx_trade_4657B0(v C.short) C.int { return C.int(uiInventoryTrade(30, uint16(v))) }

//export nox_xxx_clientTrade_465870
func nox_xxx_clientTrade_465870(v C.short) C.int { return C.int(uiInventoryTrade(28, uint16(v))) }

//export nox_xxx_send2ServInvenFail_461630
func nox_xxx_send2ServInvenFail_461630(v C.short) C.int {
	msg := [3]byte{241, byte(v), byte(uint16(v) >> 8)}
	return C.nox_xxx_netClientSend2_4E53C0(31, unsafe.Pointer(&msg[0]), 3, 0, 0)
}

//export nox_xxx_clientDrop_465BE0
func nox_xxx_clientDrop_465BE0(pos *C.int2) C.int {
	dr := uiInventoryDragged()
	if dr == nil {
		return 0
	}
	msg := [7]byte{114}
	binary.LittleEndian.PutUint16(msg[1:], uint16(nox_xxx_netGetUnitCodeCli_578B00(C.int(uiInventoryPointer(dr.C())))))
	binary.LittleEndian.PutUint16(msg[3:], uint16(pos.field_0))
	binary.LittleEndian.PutUint16(msg[5:], uint16(pos.field_4))
	return C.int(bool2int(GetServer().S().NetList.AddToMsgListCli(31, netlist.Kind0, msg[:])))
}

//export nox_xxx_clientKeyEquip_465C30
func nox_xxx_clientKeyEquip_465C30(col, row C.int) C.int {
	uiInventorySetClick(int(col), int(row))
	uiInventoryDragCopy()
	uiInventoryEquipRequest(uiInventoryDragged())
	return C.int(uiInventoryPlace(uiInventoryDragged(), int(col), int(row)))
}
func uiInventorySetDragged(dr *client.Drawable) {
	*memmap.PtrUint32(0x5D4594, 1049848) = uiInventoryPointer(unsafe.Pointer(dr))
}
