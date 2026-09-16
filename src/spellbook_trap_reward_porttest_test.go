//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/things"
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

func TestSpellbookTrapReward(t *testing.T) {
	o := newSpellbookOwner(t)
	configure, restore := o.c.srv.Server.PortTestAISpellDefs()
	t.Cleanup(restore)
	configure([]server.PortTestSpellClassDef{{Index: 34, Flags: uint32(things.SpellClassAny), Valid: true}})
	words, restore := legacy.PortTestBookTrapWords()
	t.Cleanup(restore)
	oldLoad := legacy.Nox_xxx_gLoadImg
	t.Cleanup(func() { legacy.Nox_xxx_gLoadImg = oldLoad })
	legacy.Nox_xxx_gLoadImg = func(name string) *noxrender.Image {
		for i, n := range []string{"QuickBarTrap", "QuickBarTrapHit", "QuickBarBomber", "QuickBarBomberHit"} {
			if name == n {
				o.loads = append(o.loads, name)
				return o.images[i]
			}
		}
		return oldLoad(name)
	}
	var rows []struct {
		Book      spellbookResult
		Height    uint32
		Particles [][]uint32
		RNG       [2]int
	}
	for class := 0; class < 3; class++ {
		for _, rank := range []uint32{0, 1, 2} {
			for _, notify := range []uint32{0, 1} {
				for _, enabled := range []bool{false, true} {
					func() {
						o.resetBook(t)
						o.bookCall("nox_xxx_bookInit_45B9D0")
						*(*byte)(unsafe.Add(unsafe.Pointer(&o.players[0]), 2251)) = byte(class)
						for i := 0; i < 3; i++ {
							status := gui.StatusFlags(0)
							if enabled {
								status = 8
							}
							w := o.c.GUI.NewWindowRaw(nil, status, 200+i*36, 430, 30, 30, nil)
							*words[i] = uint32(uintptr(w.C()))
						}
						*words[3] = 0xffffffff
						particles, free := legacy.PortTestEffectsScreenParticles(128)
						defer free()
						ret := o.bookCall("nox_xxx_netSpellRewardCli_45CFE0", 34, rank, notify, 1)
						burst := class != 0 && rank == 1 && notify != 0 && !enabled
						want := 0
						if burst {
							want = 50
						}
						if len(particles()) != want {
							t.Fatal("trap reward particle condition")
						}
						if burst && *words[3] != 481 || !burst && *words[3] != 0xffffffff {
							t.Fatal("trap animation starting height")
						}
						value := *(*uint32)(unsafe.Add(unsafe.Pointer(&o.players[0]), 3832))
						if class != 0 && value != rank || class == 0 && value != 0 {
							t.Fatal("trap knowledge class check")
						}
						rows = append(rows, struct {
							Book      spellbookResult
							Height    uint32
							Particles [][]uint32
							RNG       [2]int
						}{o.bookSnapshot(fmt.Sprintf("class%d-rank%d-notify%d-enabled%v", class, rank, notify, enabled), ret), *words[3], particles(), [2]int{o.c.srv.Rand.Logic.Index(), o.c.srv.Rand.Other.Index()}})
					}()
				}
			}
		}
	}
	spellbookCapture(t, "trap-reward", rows, "29207f4112ef3aec75695a466de5ece07a6e4a8e7464552d041c9f34a191f839")
}
