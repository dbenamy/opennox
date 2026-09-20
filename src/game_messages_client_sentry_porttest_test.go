//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/libs/prand"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"image"
	"math"
	"slices"
	"testing"
	"unsafe"
)

func TestGameMessageClientSentry(t *testing.T) {
	c, pix, env := newEffectsFullOwner(t, "VioletSpark")
	render := legacy.PortTestNewObjectRenderEnvironment(c.r.Data())
	t.Cleanup(render.Restore)
	words, restore := legacy.PortTestDrawableEffectGlobals()
	t.Cleanup(restore)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	serverConfigOwnBytes(t, 0x852978, 8, 4)
	local, free := alloc.Make([]uint32{}, 5)
	defer free()
	var sounds [][2]int
	t.Cleanup(legacy.PortTestClientSoundObserver(func(id, v int) { sounds = append(sounds, [2]int{id, v}) }))
	type row struct {
		On, Pause, Local, Fail, Queued int
		Seed                           uint32
		Coords                         [4]uint16
		Calls                          []effectsSpawnCall
		Drawables                      [][]uint32
		Queue                          []uint32
		Sounds                         [][2]int
		RNG                            [2]int
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for pause := 0; pause < 2; pause++ {
			for loc := 0; loc < 4; loc++ {
				for _, fail := range []int{0, 1, 2} {
					for _, queued := range []int{0, 31, 32} {
						for _, seed := range []uint32{1, 31, 1023} {
							for _, xy := range [][4]uint16{{0, 0, 0, 0}, {100, 200, 300, 400}, {65535, 65534, 65534, 65535}} {
								c.resetCase(env, pix, seed, 100)
								render.Reset()
								c.FailEvery = fail
								sounds = nil
								flags := noxflags.GameFlag(0)
								if pause != 0 {
									flags = noxflags.GamePause
								}
								reset := noxflags.PortTestGameFlags(flags)
								binary.LittleEndian.PutUint32(connected, uint32(on))
								*words[3] = uint32(c.Things.IndByID("VioletSpark"))
								*memmap.PtrPtr(0x852978, 8) = nil
								if loc != 0 {
									local[3] = uint32(xy[2]) + uint32([]int{0, 0, 599, 600}[loc])
									local[4] = uint32(xy[3])
									*memmap.PtrPtr(0x852978, 8) = unsafe.Pointer(&local[0])
								}
								for i := 0; i < queued; i++ {
									legacy.PortTestObjectRenderBeam(1, c.Viewport(), [4]int32{1, 2, 3, 4})
								}
								before := slices.Clone(render.State())
								data := make([]byte, 9)
								data[0] = 149
								for i, v := range xy {
									binary.LittleEndian.PutUint16(data[1+2*i:], v)
								}
								input := bytes.Clone(data)
								n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(149), data)
								reset()
								if n != 9 || !bytes.Equal(data, input) {
									t.Fatal("sentry input/length")
								}
								wantQueue := slices.Clone(before)
								if on != 0 && queued < 32 {
									wantQueue[1]++
									wantQueue[18+queued*2] = binary.LittleEndian.Uint32(data[1:])
									wantQueue[19+queued*2] = binary.LittleEndian.Uint32(data[5:])
								}
								if !slices.Equal(render.State(), wantQueue) {
									t.Fatal("sentry beam queue/count")
								}
								rng := prand.New(int(seed + 1))
								var wantSounds [][2]int
								if on != 0 {
									if rng.Int(0, 100) < 25 && loc != 0 && loc != 3 {
										distance := []int{0, 0, 599, 600}[loc]
										wantSounds = append(wantSounds, [2]int{297, 100 * (600 - distance) / 600})
									}
								}
								if !slices.Equal(sounds, wantSounds) {
									t.Fatal("sentry proximity sound")
								}
								count := 0
								if on != 0 && pause == 0 {
									count = 1
								}
								if len(c.Calls) != count {
									t.Fatal("sentry pause gate")
								}
								for _, call := range c.Calls {
									dx, dy := int(xy[2])-int(xy[0]), int(xy[3])-int(xy[1])
									length := int(math.Sqrt(float64(dx*dx + dy*dy)))
									if length == 0 {
										length = 1
									}
									pos := image.Pt(int(xy[2])-4*dx/length, int(xy[3])-4*dy/length)
									if call.Type != int(*words[3]) || call.Position != pos || (call.Ref != 0) != (fail != 1) {
										t.Fatal("sentry spark type/position/allocation")
									}
									if call.Ref != 0 {
										angle, speed, ttl, vel := rng.Int(0, 255), rng.Int(1, 1500), rng.Int(5, 20), rng.Int(-4, 4)
										for dr, ref := range c.refs {
											if ref != call.Ref {
												continue
											}
											raw := unsafe.Slice((*byte)(dr.C()), 512)
											if *txword(dr, 432) != uint32(pos.X)<<12 || *txword(dr, 436) != uint32(pos.Y)<<12 || *txword(dr, 440) != uint32(speed) || *txword(dr, 444) != 100 || *txword(dr, 448) != 100+uint32(ttl) || raw[299] != byte(angle) || int(int8(raw[296])) != vel || binary.LittleEndian.Uint16(raw[104:]) != 22 || dr.ObjFlags&0x400000 == 0 {
												t.Fatal("sentry spark fields")
											}
										}
									}
								}
								if c.srv.Rand.Other.Index() != rng.Index() {
									t.Fatal("sentry random consumption")
								}
								rows = append(rows, row{on, pause, loc, fail, queued, seed, xy, append([]effectsSpawnCall(nil), c.Calls...), c.snapshotDrawables(t), render.State(), slices.Clone(sounds), [2]int{c.srv.Rand.Logic.Index(), c.srv.Rand.Other.Index()}})
							}
						}
					}
				}
			}
		}
	}
	interactionCapture(t, "game-progress-sentry", rows)
}
