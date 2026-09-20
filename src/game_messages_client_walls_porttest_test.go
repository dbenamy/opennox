//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"image"
	"reflect"
	"sort"
	"testing"
	"unsafe"

	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/libs/wall"
	"github.com/opennox/opennox/v1/legacy"
)

func TestGameMessageClientMagicWalls(t *testing.T) {
	c, _, _ := newEffectsFullOwner(t)
	t.Cleanup(c.srv.PortTestMinimapWalls())
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	type row struct {
		On, Kind, Existing, Flags, Value, Return int
		Pos                                      image.Point
		Walls                                    [][16]byte
	}
	var rows []row
	sorted := func(v [][16]byte) { sort.Slice(v, func(i, j int) bool { return bytes.Compare(v[i][:], v[j][:]) < 0 }) }
	for on := 0; on < 2; on++ {
		for _, kind := range []int{61, 62} {
			for exists := 0; exists < 2; exists++ {
				for _, flags := range []wall.Flags{0, wall.FlagDoor, wall.FlagBroken} {
					for _, pos := range []image.Point{image.Pt(0, 0), image.Pt(0, 1), image.Pt(254, 254), image.Pt(255, 255), image.Pt(255, 254)} {
						for _, value := range []byte{0, 1, 127, 255} {
							c.srv.Walls.Reset()
							binary.LittleEndian.PutUint32(connected, uint32(on))
							var initial [16]byte
							if exists != 0 {
								w := c.srv.Walls.CreateAtGrid(pos)
								if w == nil {
									t.Fatal("wall allocation")
								}
								w.Dir0, w.Tile1, w.Field2, w.Field3 = 13, 17, 19, 23
								w.Flags4 = flags
								w.Health7 = 29
								w.Field8 = 0x1234
								w.Field10 = 0xabcd
								w.Field12 = 0x87654321
								copy(initial[:], unsafe.Slice((*byte)(w.C()), 16))
							}
							var want [][16]byte
							if exists != 0 {
								want = append(want, initial)
							}
							eligible := exists != 0 && flags == 0 && (pos.X+pos.Y)&1 == 0
							data := []byte{byte(kind), value, value ^ 0xa5, value ^ 0x5a, byte(pos.X), byte(pos.Y)}
							if kind == 62 {
								data = []byte{62, byte(pos.X), byte(pos.Y)}
							}
							before := bytes.Clone(data)
							if on != 0 {
								if kind == 61 {
									var changed [16]byte
									if eligible {
										changed = initial
										want = nil
									} else {
										changed[5], changed[6] = byte(pos.X), byte(pos.Y)
									}
									changed[0], changed[1], changed[2] = value^0xa5, value, value^0x5a
									want = append(want, changed)
								} else if eligible {
									want = nil
								}
							}
							n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(kind), data)
							var got [][16]byte
							for w := c.srv.Walls.IndexByY(pos.Y); w != nil; w = w.NextByY24 {
								if len(got) > 2 {
									t.Fatal("wall row cycle")
								}
								if w.GridPos() == pos {
									var b [16]byte
									copy(b[:], unsafe.Slice((*byte)(w.C()), 16))
									got = append(got, b)
								}
							}
							sorted(got)
							sorted(want)
							if n != len(data) || !bytes.Equal(data, before) || !reflect.DeepEqual(got, want) {
								t.Fatalf("wall on%d kind%d exists%d flags%x pos%v value%d return%d got%x want%x", on, kind, exists, flags, pos, value, n, got, want)
							}
							rows = append(rows, row{on, kind, exists, int(flags), int(value), n, pos, got})
							c.srv.Walls.Reset()
						}
					}
				}
			}
		}
	}
	interactionCapture(t, "game-client-magic-walls", rows)
}
