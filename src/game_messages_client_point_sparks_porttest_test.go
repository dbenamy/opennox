//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"image"
	"reflect"
	"testing"

	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
)

// Compare message dispatch with the already-qualified particle helper, then
// independently check the message's type, signed coordinates and particle count.
func TestGameMessageClientPointSparks(t *testing.T) {
	c, pix, env := newEffectsFullOwner(t, "BlueSpark", "YellowSpark", "CyanSpark", "VioletSpark")
	words, restore := legacy.PortTestDrawableEffectGlobals()
	t.Cleanup(restore)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	type state struct {
		Calls     []effectsSpawnCall
		Drawables [][]uint32
		RNG       [2]int
	}
	type row struct {
		On, Kind, Fail int
		Seed           uint32
		Pos            image.Point
		State          state
	}
	var rows []row
	for which, name := range []string{"BlueSpark", "YellowSpark", "CyanSpark", "VioletSpark"} {
		id := c.Things.IndByID(name)
		if id == 0 {
			t.Fatal("spark fixture type")
		}
		count, speed, ttl := 25, 500, 25
		if which == 0 {
			count, speed, ttl = 50, 1000, 30
		}
		for on := 0; on < 2; on++ {
			for _, fail := range []int{0, 1, 2} {
				for _, seed := range []uint32{0, 1, 1023} {
					for _, pos := range []image.Point{image.Pt(-32768, -1), image.Pt(0, 0), image.Pt(100, 200), image.Pt(32767, -32768)} {
						var expected state
						for phase := 0; phase < 2; phase++ {
							c.resetCase(env, pix, seed, 100)
							c.FailEvery = fail
							binary.LittleEndian.PutUint32(connected, uint32(on))
							*words[2] = uint32(c.Things.IndByID("BlueSpark"))
							*words[3] = uint32(c.Things.IndByID("VioletSpark"))
							*memmap.PtrUint32(0x5D4594, 1200780) = uint32(c.Things.IndByID("YellowSpark"))
							*memmap.PtrUint32(0x5D4594, 1200784) = uint32(c.Things.IndByID("CyanSpark"))
							// Historical mapped addresses do not own the two named C globals.
							*words[15+(1200776-1200772)/4] = 4
							*words[15+(1200796-1200772)/4] = 4
							if phase == 0 {
								if on != 0 {
									legacy.PortTestClientEffects(2, nil, nil, [8]int32{int32(id), int32(count), int32(speed), int32(ttl), int32(pos.X), int32(pos.Y)}, nil)
								}
							} else {
								data := []byte{byte(129 + which), byte(pos.X), byte(pos.X >> 8), byte(pos.Y), byte(pos.Y >> 8)}
								input := bytes.Clone(data)
								n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(129+which), data)
								if n != 5 || !bytes.Equal(data, input) {
									t.Fatal("point sparks length/input")
								}
							}
							got := state{append([]effectsSpawnCall(nil), c.Calls...), c.snapshotDrawables(t), [2]int{c.srv.Rand.Logic.Index(), c.srv.Rand.Other.Index()}}
							if phase == 0 {
								expected = got
								continue
							}
							if !reflect.DeepEqual(expected, got) {
								t.Fatalf("point sparks helper contract kind%d on%d fail%d seed%d pos%v", 129+which, on, fail, seed, pos)
							}
							want := 0
							if on != 0 {
								want = count
							}
							if len(c.Calls) != want {
								t.Fatal("point sparks count")
							}
							live := 0
							for i, call := range c.Calls {
								if call.Type != id || call.Position != pos {
									t.Fatal("point sparks type/coordinates")
								}
								ok := fail == 0 || (i+1)%fail != 0
								if (call.Ref != 0) != ok {
									t.Fatal("point sparks allocation schedule")
								}
								if ok {
									live++
								}
							}
							if c.Objs.Count != live {
								t.Fatal("point sparks live count")
							}
							rows = append(rows, row{on, 129 + which, fail, seed, pos, got})
						}
					}
				}
			}
		}
	}
	interactionCapture(t, "game-progress-point-sparks", rows)
}
