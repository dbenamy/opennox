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

func TestGameMessageClientSessionTrade(t *testing.T) {
	o := newTradeUIOwner(t)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	data := o.record(t, 52)
	type row struct {
		On, Kind, Active int
		Value            uint32
		Return           int
		State            tradeUIResult
	}
	var rows []row
	lengths := map[int]int{1: 2, 3: 3, 4: 15, 5: 4, 6: 14, 7: 2, 12: 52}
	operations := map[int]int{1: 24, 3: 32, 4: 27, 5: 35, 6: 31, 7: 33, 12: 22}
	for on := 0; on < 2; on++ {
		for _, active := range []uint32{0, 1, 2} {
			for _, kind := range []int{1, 3, 4, 5, 6, 7, 12} {
				for _, value := range []uint32{0, 1, 3, 0x80000000, 0xffffffff} {
					clear(data)
					data[0], data[1] = 201, byte(kind)
					switch kind {
					case 3:
						data[2] = byte(value)
					case 4:
						tradeUIAddRecord(data, byte(value&1), uint16(o.c.Things.TypeByID("RedApple").Index()), 999, value)
					case 5:
						binary.LittleEndian.PutUint16(data[2:], uint16(value))
					case 6:
						for i, v := range []uint32{value, ^value, value >> 1} {
							binary.LittleEndian.PutUint32(data[2+4*i:], v)
						}
					case 12:
						alloc.StrCopyZero16(unsafe.Slice((*uint16)(unsafe.Pointer(&data[2])), 25), "Partner Ω")
					}
					setup := func() {
						o.reset(t)
						o.construct(t)
						o.constructTrade(t)
						binary.LittleEndian.PutUint32(connected, uint32(on))
						*o.windowWords["dword_5d4594_1320964"] = active
						*o.tradeWords["dword_5d4594_1320944"], *o.tradeWords["dword_5d4594_1320948"] = 9, 10
						o.tradeCell(t, 0, 0, 2, "RedApple", 1)
						o.tradeCell(t, 1, 0, 2, "RedApple", 101)
					}
					setup()
					if on != 0 {
						legacy.PortTestTradeUI(operations[kind], txptr(unsafe.Pointer(&data[0])))
					}
					want := o.tradeCapture(t, 0, 20)
					setup()
					input := bytes.Clone(data)
					n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(201), data[:lengths[kind]])
					got := o.tradeCapture(t, 0, 20)
					if n != lengths[kind] || !bytes.Equal(input, data) || !reflect.DeepEqual(got, want) {
						t.Fatal("trade dispatch", on, active, kind, value, n)
					}
					rows = append(rows, row{on, kind, int(active), value, n, got})
				}
			}
		}
	}
	interactionCapture(t, "game-client-session-trade", rows)
}
