//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"testing"
	"unsafe"

	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func TestGameMessageClientRayEffects(t *testing.T) {
	c, pix, env := newEffectsFullOwner(t)
	words, restore := legacy.PortTestDrawableEffectGlobals()
	t.Cleanup(restore)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	data, free := alloc.Make([]byte{}, 9)
	defer free()
	endpoints, freeEndpoints := alloc.Make([]int32{}, 4)
	defer freeEndpoints()
	type state struct {
		Calls     []effectsSpawnCall
		Drawables [][]uint32
		Cache     []uint32
		Globals   []uint32
		RNG       [2]int
	}
	type row struct {
		On, Kind, Fail int
		Seed           uint32
		Coords         [4]uint16
		State          state
	}
	var rows []row
	for _, kind := range []byte{125, 140, 141, 142, 143, 144, 145} {
		for on := 0; on < 2; on++ {
			for _, fail := range []int{0, 1, 2} {
				for _, seed := range []uint32{0, 1, 1023} {
					for _, xy := range [][4]uint16{{0, 0, 0, 0}, {48, 70, 65, 91}, {65535, 65534, 65534, 65535}} {
						var expected state
						for phase := 0; phase < 2; phase++ {
							c.resetCase(env, pix, seed, 100)
							c.FailEvery = fail
							binary.LittleEndian.PutUint32(connected, uint32(on))
							blue := c.Things.IndByID("BlueSpark")
							if blue == 0 {
								t.Fatal("ray spark fixture")
							}
							*words[2] = uint32(blue)
							data[0] = kind
							for i, v := range xy {
								binary.LittleEndian.PutUint16(data[1+2*i:], v)
								endpoints[i] = int32(v)
							}
							input := bytes.Clone(data)
							if phase == 0 {
								if on != 0 {
									legacy.PortTestClientEffects(8, nil, nil, [8]int32{}, unsafe.Pointer(&data[0]))
									if kind == 140 || kind == 142 {
										legacy.PortTestClientEffects(5, nil, nil, [8]int32{int32(blue)}, unsafe.Pointer(&endpoints[0]))
									}
									if kind == 125 {
										legacy.PortTestClientEffects(3, nil, nil, [8]int32{int32(xy[2]), int32(xy[3]), 10, int32(blue)}, nil)
									}
								}
							} else {
								if n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(kind), data); n != 9 {
									t.Fatal("ray effect length", n)
								}
							}
							if !bytes.Equal(data, input) {
								t.Fatal("ray input mutation")
							}
							cache := append([]uint32(nil), unsafe.Slice(memmap.PtrUint32(0x5D4594, 1303540), 203)...)
							for i, v := range cache[:96] {
								if v != 0 {
									r := c.refs[(*client.Drawable)(unsafe.Pointer(uintptr(v)))]
									if r == 0 {
										t.Fatal("unknown ray cache owner")
									}
									cache[i] = r
								}
							}
							got := state{append([]effectsSpawnCall(nil), c.Calls...), c.snapshotDrawables(t), cache, env.Snapshot(), [2]int{c.srv.Rand.Logic.Index(), c.srv.Rand.Other.Index()}}
							if phase == 0 {
								expected = got
								continue
							}
							if !reflect.DeepEqual(expected, got) {
								t.Fatalf("ray dispatch kind%d on%d fail%d seed%d xy%v", kind, on, fail, seed, xy)
							}
							if on == 0 && len(c.Calls) != 0 {
								t.Fatal("disconnected ray spawned")
							}
							rows = append(rows, row{on, int(kind), fail, seed, xy, got})
						}
					}
				}
			}
		}
	}
	interactionCapture(t, "game-progress-ray-effects", rows)
}
