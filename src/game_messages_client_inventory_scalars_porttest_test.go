//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/opennox/v1/legacy"
)

func TestGameMessageClientInventoryScalars(t *testing.T) {
	o := newUIInventoryOwner(t)
	o.reset(t)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	state := o.regions[2]
	table := o.regions[3]
	for i := range state {
		state[i] = byte(i*7 + 11)
	}
	for i := range table {
		table[i] = byte(i*13 + 3)
	}
	initialState, initialTable := bytes.Clone(state), bytes.Clone(table)
	*o.meters.NamedWord("nox_player_netCode_85319C") = 7
	type row struct {
		On, Kind, Code, Index, Return int
		Value, Gold, Armor, Modifier  uint32
	}
	var rows []row
	values := []uint32{0, 1, 0x80000000, 0x007fffff, 0x00800000, 0x3f800000, 0x7f7fffff, 0x7f800000, 0xff800000, 0x7fc12345, 0xffffffff}
	for on := 0; on < 2; on++ {
		for _, kind := range []int{73, 74, 104} {
			codes, indices := []uint16{0}, []byte{0}
			if kind == 104 {
				codes = []uint16{7, 8, 0x8007, 0xffff}
				indices = []byte{0, 1, 127, 128, 255}
			}
			for _, code := range codes {
				for _, index := range indices {
					for _, value := range values {
						binary.LittleEndian.PutUint32(connected, uint32(on))
						copy(state, initialState)
						copy(table, initialTable)
						*o.words[0] = 0x13572468
						wantState, wantTable := bytes.Clone(state), bytes.Clone(table)
						wantGold := uint32(0x13572468)
						data := []byte{byte(kind)}
						if kind == 104 {
							data = binary.LittleEndian.AppendUint16(data, code)
						}
						data = binary.LittleEndian.AppendUint32(data, value)
						if kind == 104 {
							data = append(data, index)
						}
						before := bytes.Clone(data)
						if on != 0 {
							switch kind {
							case 73:
								binary.LittleEndian.PutUint32(wantState[12:], value)
							case 74:
								wantGold = value
							case 104:
								if code&0x7fff == 7 {
									binary.LittleEndian.PutUint32(wantTable[int(index)*4:], value)
								}
							}
						}
						n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(kind), data)
						if n != len(data) || !bytes.Equal(data, before) || !bytes.Equal(state, wantState) || !bytes.Equal(table, wantTable) || *o.words[0] != wantGold {
							t.Fatalf("inventory scalar on%d kind%d code%x index%d value%x ret%d", on, kind, code, index, value, n)
						}
						rows = append(rows, row{on, kind, int(code), int(index), n, value, *o.words[0], binary.LittleEndian.Uint32(state[12:]), binary.LittleEndian.Uint32(table[int(index)*4:])})
					}
				}
			}
		}
	}
	gameMessageCapture(t, "game-client-inventory-scalars", rows)
}
