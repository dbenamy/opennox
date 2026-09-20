//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"testing"

	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/opennox/v1/legacy"
)

func TestGameMessageClientSessionInventory(t *testing.T) {
	o := newInventoryDisplayOwner(t)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	type row struct {
		On, Kind, Code, Status, Index int
		Previous                      bool
		State                         inventoryDisplayResult
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for _, kind := range []int{224, 225} {
			for _, code := range []int{0, 321, 322, 999, 0x8000, 0x8141, 0x8142, 0xffff} {
				for _, status := range []int{0, 1, 2, 127, 128, 255} {
					for _, index := range []int{0, 83} {
						for _, previous := range []bool{false, true} {
							setup := func() {
								o.reset(t)
								binary.LittleEndian.PutUint32(connected, uint32(on))
								cells := legacy.PortTestInventoryCells()
								dr := o.item(t, "Bow", 123)
								cells[index].Drawable, cells[index].Count = dr, 2
								cells[index].Codes[0], cells[index].Codes[1] = 321, 322
								if previous {
									old := (index + 1) % 84
									p := o.item(t, "Bow", 500)
									cells[old].Drawable, cells[old].Count, cells[old].Codes[0], cells[old].Alternate = p, 1, 501, 1
									*o.displayWords["dword_5d4594_1062480"] = o.cell(old)
								}
								*o.displayWords["dword_5d4594_1062484"] = 999
								*o.displayWords["dword_5d4594_1062488"] = 777
							}
							setup()
							if kind == 224 {
								legacy.PortTestInventoryDisplay(15, uintptr(code&0x7fff), uintptr(status), 0)
							} else {
								legacy.PortTestUIInventoryCall(23, uint32(code&0x7fff), 0, 0)
							}
							want := o.displaySnapshot(t, 0, 15, 0)
							setup()
							data := binary.LittleEndian.AppendUint16([]byte{byte(kind)}, uint16(code))
							if kind == 224 {
								data = append(data, byte(status))
							}
							input := bytes.Clone(data)
							n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(kind), data)
							got := o.displaySnapshot(t, 0, 15, 0)
							if n != len(data) || !bytes.Equal(data, input) || !reflect.DeepEqual(got, want) {
								t.Fatalf("inventory on%d kind%d code%d status%d index%d previous%v", on, kind, code, status, index, previous)
							}
							if kind == 225 && *o.displayWords["dword_5d4594_1062488"] != uint32(code&0x7fff) {
								t.Fatal("quiver code mask")
							}
							if kind == 224 && *o.displayWords["dword_5d4594_1062488"] != 777 {
								t.Fatal("secondary message changed quiver")
							}
							rows = append(rows, row{on, kind, code, status, index, previous, got})
						}
					}
				}
			}
		}
	}
	interactionCapture(t, "game-client-session-inventory", rows)
}
