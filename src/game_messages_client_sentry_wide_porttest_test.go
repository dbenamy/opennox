//go:build porttest

package opennox

import (
	"encoding/binary"
	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/libs/prand"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"image"
	"slices"
	"testing"
	"unsafe"
)

// The two legacy distance expressions have different signedness. Large wire
// coordinates distinguish a wrapped signed spark distance from unsigned audio distance.
func TestGameMessageClientSentryWideCoordinates(t *testing.T) {
	c, pix, env := newEffectsFullOwner(t, "VioletSpark")
	render := legacy.PortTestNewObjectRenderEnvironment(c.r.Data())
	t.Cleanup(render.Restore)
	words, restore := legacy.PortTestDrawableEffectGlobals()
	t.Cleanup(restore)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	serverConfigOwnBytes(t, 0x852978, 8, 4)
	local, free := alloc.Make([]uint32{}, 5)
	defer free()
	*memmap.PtrPtr(0x852978, 8) = unsafe.Pointer(&local[0])
	var sounds [][2]int
	t.Cleanup(legacy.PortTestClientSoundObserver(func(id, v int) { sounds = append(sounds, [2]int{id, v}) }))
	type row struct {
		On, Pause, Fail, Local int
		Seed                   uint32
		Coords                 [4]uint16
		Calls                  []effectsSpawnCall
		Drawables              [][]uint32
		Sounds                 [][2]int
		RNG                    [2]int
	}
	var rows []row
	vectors := [][4]uint16{{0, 0, 65535, 0}, {0, 0, 65535, 65535}, {0, 0, 40000, 40000}, {0, 0, 46340, 0}, {0, 0, 46341, 0}, {65535, 65535, 0, 0}}
	for on := 0; on < 2; on++ {
		for pause := 0; pause < 2; pause++ {
			for _, fail := range []int{0, 1} {
				for mode := 0; mode < 3; mode++ {
					for _, seed := range []uint32{1, 31, 1023} {
						for index, xy := range vectors {
							c.resetCase(env, pix, seed, 100)
							render.Reset()
							c.FailEvery = fail
							sounds = nil
							*memmap.PtrPtr(0x852978, 8) = unsafe.Pointer(&local[0])
							flags := noxflags.GameFlag(0)
							if pause != 0 {
								flags = noxflags.GamePause
							}
							reset := noxflags.PortTestGameFlags(flags)
							binary.LittleEndian.PutUint32(connected, uint32(on))
							*words[3] = uint32(c.Things.IndByID("VioletSpark"))
							// Mode1 is outside signed-squared range but outside sound range too; mode2
							// squares to zero in the original 32-bit arithmetic.
							local[3] = uint32(xy[2]) - []uint32{0, 46341, 0x80000000}[mode]
							local[4] = uint32(xy[3])
							data := make([]byte, 9)
							data[0] = 149
							for i, v := range xy {
								binary.LittleEndian.PutUint16(data[1+2*i:], v)
							}
							n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(149), data)
							reset()
							if n != 9 {
								t.Fatal("wide sentry length")
							}
							rng := prand.New(int(seed + 1))
							var wantSounds [][2]int
							if on != 0 && rng.Int(0, 100) < 25 && mode != 1 {
								wantSounds = append(wantSounds, [2]int{297, 100})
							}
							if !slices.Equal(sounds, wantSounds) {
								t.Fatalf("wide sentry unsigned sound mode%d sounds%v want%v", mode, sounds, wantSounds)
							}
							count := 0
							if on != 0 && pause == 0 {
								count = 1
							}
							if len(c.Calls) != count {
								t.Fatal("wide sentry pause/connection")
							}
							if count != 0 {
								want := image.Pt(int(xy[2]), int(xy[3]))
								if index == 3 {
									want.X -= 4
								}
								call := c.Calls[0]
								if call.Position != want || (call.Ref != 0) != (fail == 0) {
									t.Fatal("wide sentry signed-distance position")
								}
							}
							rows = append(rows, row{on, pause, fail, mode, seed, xy, slices.Clone(c.Calls), c.snapshotDrawables(t), slices.Clone(sounds), [2]int{c.srv.Rand.Logic.Index(), c.srv.Rand.Other.Index()}})
						}
					}
				}
			}
		}
	}
	interactionCapture(t, "game-progress-sentry-wide", rows)
}
