//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/things"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"image"
	"testing"
	"unsafe"
)

func TestClientPresentationBookReward(t *testing.T) {
	o, bar, configure := newSpellbookAdditionOwner(t)
	oldAbilities := o.c.srv.abilities.defs
	t.Cleanup(func() { o.c.srv.abilities.defs = oldAbilities })
	o.c.srv.abilities.defs[1] = AbilityDef{name: "Berserk", field24: 1}
	guides := serverConfigOwnBytes(t, 0x5D4594, 740076, 41*7*4)
	clear(guides)
	*memmap.PtrUint32(0x5D4594, 740076+8*4) = 1
	clear(serverConfigOwnBytes(t, 0x5D4594, 1217504, 4))
	timestamp := memmap.PtrUint32(0x5D4594, 1217504)
	type record struct {
		Kind, Auto, Height int
		Coop               bool
		Delta, Timestamp   uint32
		Book               spellbookResult
		Particles          [][]uint32
		RNG                [2]int
		Slot               uint32
	}
	var rows []record
	for _, kind := range []int{2, 3, 4} {
		for _, auto := range []int{0, 1, 2} {
			for _, coop := range []bool{false, true} {
				for _, delta := range []uint32{0, 1, 2, 0xffffffff} {
					for _, height := range []int{480, 601} {
						func() {
							prepareSpellbookAddition(t, o, bar, 640)
							configure([]server.PortTestSpellClassDef{{Index: 1, Flags: uint32(things.SpellClassAny), Valid: true}})
							class := byte(1)
							if kind == 3 {
								class = 0
							}
							if kind == 4 {
								class = 2
							}
							*(*byte)(unsafe.Add(unsafe.Pointer(&o.players[0]), 2251)) = class
							*(*uint32)(unsafe.Add(unsafe.Pointer(&o.players[0]), 3700)) = 1
							*(*uint32)(unsafe.Add(unsafe.Pointer(&o.players[0]), 4248)) = 1
							if coop {
								noxflags.SetGame(noxflags.GameModeCoop)
							}
							*o.words["nox_win_height"] = uint32(height)
							seq := uint32(o.c.GetInputSeq())
							*timestamp = seq - delta
							particles, free := legacy.PortTestEffectsScreenParticles(768)
							defer free()
							legacy.PortTestPresentationBookReward(kind, 1, auto)
							run := !coop || delta >= 2
							wantCount := 0
							if run {
								wantCount = 600
							}
							// In solo play autofill finishes immediately, producing its extra 50 particles.
							if run && kind != 4 && auto == 1 && !coop {
								wantCount += 50
							}
							if len(particles()) != wantCount {
								t.Fatalf("book particles kind%d auto%d coop%v delta%d: %d want%d", kind, auto, coop, delta, len(particles()), wantCount)
							}
							wantTime := seq - delta
							if run {
								wantTime = seq
							}
							if *timestamp != wantTime {
								t.Fatal("book throttle timestamp")
							}
							if run && o.bookWindow().Off != image.Pt(5, height/3) {
								t.Fatal("book reward placement")
							}
							rows = append(rows, record{kind, auto, height, coop, delta, *timestamp, o.bookSnapshot(fmt.Sprintf("kind%d-auto%d-coop%v-delta%d-height%d", kind, auto, coop, delta, height), 0), particles(), [2]int{o.c.srv.Rand.Logic.Index(), o.c.srv.Rand.Other.Index()}, bar[0]})
						}()
					}
				}
			}
		}
	}
	spellbookCapture(t, "client-presentation-book-reward", rows, "ca478b6c0369733893d64dd6a1eff38171702cfe904cc65f88d7e6f465189fea")
}
