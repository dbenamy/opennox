//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"testing"
	"unsafe"

	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/opennox/v1/legacy"
)

func TestGameMessageClientItemUpdates(t *testing.T) {
	o := newUIInventoryOwner(t)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	type row struct {
		On, Kind, Mode, Code, Value, Equipped, Return int
		Items                                         [][4]uint16
		Meters                                        [2][2]uint32
		Hidden                                        [2]bool
		Lookup                                        bool
		StackIndex                                    uint32
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for _, kind := range []int{68, 100} {
			for mode := 0; mode < 4; mode++ {
				for _, code := range []uint16{0, 0x1000, 0x1001, 0x9000, 0xffff} {
					for _, value := range []uint16{0, 1, 255, 32767, 32768, 65535} {
						for _, equipped := range []uint32{0, 1, 2} {
							o.reset(t)
							binary.LittleEndian.PutUint32(connected, uint32(on))
							if mode == 0 || mode == 3 {
								o.stack(83, 2)
							}
							if mode == 1 || mode == 3 {
								binary.LittleEndian.PutUint32(o.regions[1], uint32(uintptr(o.items[1].C())))
								o.items[1].NetCode32 = 0x1000
							}
							binary.LittleEndian.PutUint32(o.grid[83*148+132:], equipped)
							for _, i := range []int{5, 6} {
								o.meters.Records[i].Window.SetHidden(true)
							}
							wantMeters := append([]legacy.PortTestMeterRecord(nil), o.meters.Records...)
							var raw, want [][]byte
							for _, dr := range o.items {
								b := unsafe.Slice((*byte)(dr.C()), int(unsafe.Sizeof(*dr)))
								raw = append(raw, b)
								want = append(want, bytes.Clone(b))
							}
							grid := bytes.Clone(o.grid)
							data := []byte{byte(kind), byte(code), byte(code >> 8)}
							current, maximum := value, ^value
							if kind == 68 {
								data = binary.LittleEndian.AppendUint16(data, current)
								data = binary.LittleEndian.AppendUint16(data, maximum)
							} else {
								current, maximum = uint16(byte(current)), uint16(byte(maximum))
								data = append(data, byte(current), byte(maximum))
							}
							before := bytes.Clone(data)
							id := code & 0x7fff
							gridMatch := (mode == 0 || mode == 3) && (id == 0x1000 || id == 0x1001)
							dragMatch := (mode == 1 || mode == 3) && id == 0x1000
							target := -1
							if on != 0 {
								if gridMatch {
									target = 0
								} else if kind == 68 && dragMatch {
									target = 1
								}
							}
							show := on != 0 && kind == 100 && gridMatch && equipped == 1
							if target >= 0 {
								off := 292
								if kind == 100 {
									off = 448
								}
								binary.LittleEndian.PutUint16(want[target][off:], current)
								binary.LittleEndian.PutUint16(want[target][off+2:], maximum)
							}
							if show {
								for _, i := range []int{5, 6} {
									wantMeters[i].Current = uint32(current)
									wantMeters[i].Maximum = uint32(maximum)
								}
							}
							n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(kind), data)
							if n != len(data) || !bytes.Equal(data, before) || !bytes.Equal(o.grid, grid) || !reflect.DeepEqual(o.meters.Records, wantMeters) {
								t.Fatalf("item update on%d kind%d mode%d code%x value%d equipped%d ret%d", on, kind, mode, code, value, equipped, n)
							}
							r := row{On: on, Kind: kind, Mode: mode, Code: int(code), Value: int(value), Equipped: int(equipped), Return: n}
							for i, b := range raw {
								if !bytes.Equal(b, want[i]) {
									t.Fatalf("item %d mutation on%d kind%d mode%d code%x value%d", i, on, kind, mode, code, value)
								}
								r.Items = append(r.Items, [4]uint16{binary.LittleEndian.Uint16(b[292:]), binary.LittleEndian.Uint16(b[294:]), binary.LittleEndian.Uint16(b[448:]), binary.LittleEndian.Uint16(b[450:])})
							}
							for j, i := range []int{5, 6} {
								m := o.meters.Records[i]
								r.Meters[j] = [2]uint32{m.Current, m.Maximum}
								r.Hidden[j] = m.Window.GetFlags().Has(0x10)
								if r.Hidden[j] == show {
									t.Fatal("charge meter visibility")
								}
							}
							lookup, stackIndex := binary.LittleEndian.Uint32(o.regions[0]), binary.LittleEndian.Uint32(o.regions[0][4:])
							wantLookup, wantIndex := uint32(0), uint32(0)
							if on != 0 && gridMatch {
								wantLookup = o.cell(83)
								if id == 0x1001 {
									wantIndex = 1
								}
							}
							if lookup != wantLookup || stackIndex != wantIndex {
								t.Fatal("inventory lookup cell/index")
							}
							r.Lookup = lookup != 0
							r.StackIndex = stackIndex
							rows = append(rows, r)
						}
					}
				}
			}
		}
	}
	interactionCapture(t, "game-client-item-updates", rows)
}
