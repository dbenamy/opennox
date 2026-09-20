//go:build porttest

package opennox

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"reflect"
	"testing"
	"unsafe"

	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/libs/prand"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/memmap/nox/blobdata"
	"github.com/opennox/opennox/v1/legacy"
)

func TestGameMessageClientSessionRewards(t *testing.T) {
	q := newQuickbarOwner(t)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	oldAbilities := q.c.srv.abilities.defs
	t.Cleanup(func() { q.c.srv.abilities.defs = oldAbilities })
	q.c.srv.abilities.defs[1] = AbilityDef{name: "Berserk", field24: 1}
	guides := unsafe.Slice(memmap.PtrUint32(0x5D4594, 740076), 41*7)
	oldGuides := append([]uint32(nil), guides...)
	t.Cleanup(func() { copy(guides, oldGuides) })
	clear(guides)
	for i := 1; i <= 40; i++ {
		guides[7*i+1] = 1
	}
	family := serverConfigOwnBytes(t, 0x587000, 132100, 32)
	copy(family, blobdata.PortTestBookGuideFamily())
	*memmap.PtrPtr(0x587000, 132124) = memmap.PtrOff(0x587000, 132100)
	timestamp := serverConfigOwnBytes(t, 0x5D4594, 1217504, 4)
	type state struct {
		Quickbar  quickbarResult
		Particles [32]byte
		RNG       [2]int
		Timestamp uint32
	}
	type row struct {
		On, Kind, ID, Value, Known, Return int
		State                              state
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for _, kind := range []int{205, 209} {
			ids := []int{0, 1, 5, 6, 255}
			if kind == 209 {
				ids = []int{1, 24, 40}
			}
			for _, id := range ids {
				for _, known := range []int{0, 1, 7} {
					values := []int{0, 1, 2, 7, 127, 128, 129, 255}
					if kind == 205 && id == 1 {
						values = make([]int, 256)
						for i := range values {
							values[i] = i
						}
					}
					for _, value := range values {
						run := func(dispatch bool) (state, int) {
							q.reset(t)
							q.c.srv.Rand.Logic, q.c.srv.Rand.Other = prand.New(17), prand.New(31)
							noxflags.UnsetGame(noxflags.GetGame())
							noxflags.SetGame(noxflags.GameFlag(2))
							binary.LittleEndian.PutUint32(connected, uint32(on))
							clear(timestamp)
							class := byte(0)
							if kind == 209 {
								class = 2
							}
							*(*byte)(unsafe.Add(unsafe.Pointer(&q.players[0]), 2251)) = class
							levels := unsafe.Slice((*uint32)(unsafe.Add(unsafe.Pointer(&q.players[0]), 4244)), 41)
							for i := range levels {
								levels[i] = uint32(known)
							}
							for i := 1; i <= 5; i++ {
								*memmap.PtrUint32(0x5D4594, 1047764+uintptr(24*i)+16) = uint32(known)
							}
							particles, free := legacy.PortTestEffectsScreenParticles(768)
							defer free()
							n := 3
							if dispatch {
								data := []byte{byte(kind), byte(id), byte(value)}
								input := bytes.Clone(data)
								n = legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(kind), data)
								if !bytes.Equal(data, input) {
									t.Fatal("reward changed input")
								}
							} else if on != 0 {
								if kind == 205 {
									q.call("nox_xxx_netAbilityRewardCli_4611E0", uint32(id), uint32(value&127), uint32(value>>7))
								} else {
									q.bookCall("nox_xxx_netGuideRewardCli_45D140", uint32(id), uint32(value&127))
								}
							}
							if dispatch && kind == 205 {
								for slot := 1; slot <= 5; slot++ {
									want := uint32(known)
									if on != 0 && slot == id {
										want = uint32(value & 127)
									}
									if memmap.Uint32(0x5D4594, 1047764+uintptr(24*slot)+16) != want {
										t.Fatal("ability rank mask/isolation", on, id, value, slot)
									}
								}
								// The signed char local tests the incoming high bit.
								if value < 128 && len(particles()) != 0 {
									t.Fatal("ability message unexpectedly notified", on, id, value, known)
								}
							}
							if dispatch && kind == 209 {
								for slot := range levels {
									want := uint32(known)
									related := id == 24 && (slot == 7 || slot == 8 || slot == 25 || slot == 26)
									if on != 0 && (slot == id || related) {
										want = 1
									}
									if levels[slot] != want {
										t.Fatal("guide knowledge/family", on, id, value, slot)
									}
								}
							}
							particleData, err := json.Marshal(particles())
							if err != nil {
								t.Fatal(err)
							}
							return state{q.snapshot("reward", 0), sha256.Sum256(particleData), [2]int{q.c.srv.Rand.Logic.Index(), q.c.srv.Rand.Other.Index()}, binary.LittleEndian.Uint32(timestamp)}, n
						}
						want, _ := run(false)
						got, n := run(true)
						if n != 3 || !reflect.DeepEqual(got, want) {
							t.Fatalf("reward on%d kind%d id%d value%d known%d return%d quickbar%v particles%v rng%v/%v timestamp%d/%d", on, kind, id, value, known, n, reflect.DeepEqual(got.Quickbar, want.Quickbar), got.Particles == want.Particles, got.RNG, want.RNG, got.Timestamp, want.Timestamp)
						}
						rows = append(rows, row{on, kind, id, value, known, n, got})
					}
				}
			}
		}
	}
	interactionCapture(t, "game-client-session-rewards", rows)
}
