//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"reflect"
	"testing"
	"unsafe"
)

func TestGameMessageClientSessionShop(t *testing.T) {
	o := newShopUIOwner(t)
	o.installMessageDialog(t)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	data, free := alloc.Make([]byte{}, 86)
	t.Cleanup(free)
	type row struct {
		On, Kind, Active int
		Value            uint32
		Return           int
		State            shopUIResult
	}
	var rows []row
	lengths := map[int]int{2: 2, 8: 18, 9: 4, 13: 86, 27: 4, 29: 8, 31: 8}
	for on := 0; on < 2; on++ {
		for active := 0; active < 2; active++ {
			for _, kind := range []int{2, 8, 9, 13, 27, 29, 31} {
				for _, value := range []uint32{0, 1, 0x80000000, 0xffffffff} {
					clear(data)
					data[0], data[1] = 201, byte(kind)
					binary.LittleEndian.PutUint16(data[2:], 123)
					binary.LittleEndian.PutUint32(data[4:], value)
					switch kind {
					case 8:
						binary.LittleEndian.PutUint16(data[2:], uint16(o.c.Things.TypeByID("RedApple").Index()))
						binary.LittleEndian.PutUint16(data[4:], 123)
						binary.LittleEndian.PutUint32(data[6:], value)
						binary.LittleEndian.PutUint32(data[10:], ^value)
						copy(data[14:], []byte{255, 255, 255, 255})
					case 13:
						binary.LittleEndian.PutUint16(data[2:], uint16(o.c.Things.TypeByID("Shopkeeper").Index()))
						clear(data[4:])
						alloc.StrCopyZero16(unsafe.Slice((*uint16)(unsafe.Pointer(&data[4])), 25), "Merchant Ω")
						copy(data[54:], "FixtureGreeting")
					case 27:
						binary.LittleEndian.PutUint16(data[2:], uint16(value))
					}
					setup := func() {
						sub_44A400()
						o.reset(t)
						o.construct(t)
						o.constructShop(t)
						o.initAmount(t)
						binary.LittleEndian.PutUint32(connected, uint32(on))
						*o.windowWords["dword_5d4594_1098624"] = uint32(active)
						cell := &legacy.PortTestInventoryCells()[0]
						cell.Drawable = o.item(t, "RedApple", 123)
						cell.Count = 2
						cell.Codes[0], cell.Codes[1] = 123, 124
						legacy.PortTestShopUI(16, uintptr(o.c.Things.TypeByID("RedApple").Index()), 123, 7, 31, 0)
					}
					direct := func() {
						switch kind {
						case 2:
							legacy.PortTestShopUI(15)
						case 8:
							legacy.PortTestShopUI(16, uintptr(binary.LittleEndian.Uint16(data[2:])), uintptr(binary.LittleEndian.Uint16(data[4:])), uintptr(value), uintptr(uint16(^value)), txptr(unsafe.Pointer(&data[14])))
						case 9:
							legacy.PortTestShopUI(19, 123)
						case 13:
							legacy.PortTestShopUI(37, txptr(unsafe.Pointer(&data[4])), txptr(unsafe.Pointer(&data[54])), uintptr(binary.LittleEndian.Uint16(data[2:])))
						case 27:
							legacy.PortTestShopUI(38, uintptr(uint16(value)))
						case 29:
							legacy.PortTestShopUI(40, 123, uintptr(value))
						case 31:
							legacy.PortTestShopUI(41, 123, uintptr(value))
						}
					}
					setup()
					if on != 0 {
						direct()
					}
					want := o.shopCapture(t, 0, 0)
					setup()
					input := bytes.Clone(data)
					n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(201), data[:lengths[kind]])
					got := o.shopCapture(t, 0, 0)
					if n != lengths[kind] || !bytes.Equal(input, data) || !reflect.DeepEqual(got, want) {
						t.Fatal("shop dispatch", on, active, kind, value, n)
					}
					rows = append(rows, row{on, kind, active, value, n, got})
				}
			}
		}
	}
	interactionCapture(t, "game-client-session-shop", rows)
}
