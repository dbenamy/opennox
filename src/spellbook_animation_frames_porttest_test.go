//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/things"
	"github.com/opennox/opennox/v1/client/gui"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

func TestSpellbookAnimationFrames(t *testing.T) {
	o, bar, configure := newSpellbookAdditionOwner(t)
	pause, restore := legacy.PortTestBookPauseOwner()
	t.Cleanup(restore)
	oldTicks, oldFrame, oldLoading := nox_gameTicks_371764, nox_gameFrame_371772, dword_5d4594_1563080
	t.Cleanup(func() {
		nox_gameTicks_371764, nox_gameFrame_371772, dword_5d4594_1563080 = oldTicks, oldFrame, oldLoading
	})
	dword_5d4594_1563080 = false
	var rows []struct {
		Book      spellbookResult
		Particles [][]uint32
		RNG       [2]int
		Slot      uint32
		Flags     uint32
	}
	for _, width := range []uint32{640, 750, 1000} {
		for _, kind := range []uint32{2, 3, 4} {
			for _, finishEarly := range []bool{false, true} {
				func() {
					prepareSpellbookAddition(t, o, bar, width)
					configure([]server.PortTestSpellClassDef{{Index: 1, Flags: uint32(things.SpellClassAny), Valid: true}})
					pause()
					noxflags.SetGame(noxflags.GameModeCoop | noxflags.GamePause)
					o.c.Inp.Tick()
					particles, freeParticles := legacy.PortTestEffectsScreenParticles(512)
					defer freeParticles()
					record := func(step string, ret uint32) {
						rows = append(rows, struct {
							Book      spellbookResult
							Particles [][]uint32
							RNG       [2]int
							Slot      uint32
							Flags     uint32
						}{o.bookSnapshot(fmt.Sprintf("width%d-kind%d-early%v-%s", width, kind, finishEarly, step), ret), particles(), [2]int{o.c.srv.Rand.Logic.Index(), o.c.srv.Rand.Other.Index()}, bar[0], uint32(noxflags.GetGame())})
					}
					o.bookCall("nox_xxx_bookFillAll_45D570", kind, 1)
					if *o.words["dword_5d4594_1047520"] != 1 || bar[0] != 0 || len(particles()) != 0 {
						t.Fatal("cooperative addition must remain pending")
					}
					record("start", 0)
					wait := 2 * int(o.c.srv.TickRate())
					if wait <= 0 {
						t.Fatal("positive animation wait")
					}
					for tick := 0; tick < wait; tick++ {
						before := *legacy.PortTestBookVector()
						ret := o.bookCall("nox_xxx_bookDrawFn_45C7D0", *o.words["dword_5d4594_1046956"])
						if ret != 1 || len(particles()) != 0 || *legacy.PortTestBookVector() != before {
							t.Fatal("initial delay must not advance animation")
						}
						if tick == 0 || tick == wait-1 {
							record(fmt.Sprintf("wait%d", tick), ret)
						}
						o.c.Inp.Tick()
					}
					for frame := 0; frame < 300; frame++ {
						ret := o.bookCall("nox_xxx_bookDrawFn_45C7D0", *o.words["dword_5d4594_1046956"])
						record(fmt.Sprintf("frame%d", frame), ret)
						if frame == 0 && len(particles()) != 52 {
							t.Fatal("first frame particle burst")
						}
						if finishEarly && frame == 3 {
							o.bookCall("sub_45D870")
							record("early-finish", 0)
						}
						if *o.words["dword_5d4594_1047520"] == 0 {
							break
						}
						if frame == 299 {
							t.Fatal("animation failed to finish")
						}
						o.c.Inp.Tick()
					}
					if bar[0] != 1 || noxflags.HasGame(noxflags.GamePause) {
						t.Fatal("completion must add slot and release pause")
					}
					win := (*gui.Window)(unsafe.Pointer(uintptr(*o.words["dword_5d4594_1046956"])))
					if !win.GetFlags().IsHidden() {
						t.Fatal("completed animation window remains visible")
					}
				}()
			}
		}
	}
	spellbookCapture(t, "animation-frames", rows, "d134eee7d7194447696a8b210bf73897152457466e25fb75933b8b345671b01d")
}
