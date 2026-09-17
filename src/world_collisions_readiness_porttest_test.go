//go:build porttest

package opennox

import (
	"fmt"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestWorldCollisionsQuestReadiness(t *testing.T) {
	o := newWorldCollisionOwner(t)
	oldAllow, oldInfinite, oldInc := questAllowDefault, questLevelWarpInfinite, questLevelWarpInc
	t.Cleanup(func() { questAllowDefault, questLevelWarpInfinite, questLevelWarpInc = oldAllow, oldInfinite, oldInc })
	oldEngine := noxflags.GetEngine()
	t.Cleanup(func() { noxflags.ResetEngine(); noxflags.SetEngine(oldEngine) })
	questLevelWarpInfinite = false
	questLevelWarpInc = 5
	var rows []struct {
		Name   string
		Return uint32
	}
	defer func() {
		spellbookCapture(t, "world-collisions-quest-readiness", rows, "685d3180c749b4387013f194d5c3841008d2650471e0c8e7f74980743861f8aa")
	}()
	for _, op := range []int{5, 6} {
		for _, active := range []int{0, 1, 3, 4, 7} {
			for _, ready := range []int{0, 1, 3, 4, 7} {
				for _, participation := range []uint32{0, 1, 2} {
					for _, hosting := range []bool{false, true} {
						for _, headless := range []bool{false, true} {
							for _, allow := range []bool{false, true} {
								for _, stage := range []uint32{0, 4, 5, 19, 20} {
									name := fmt.Sprintf("op%d/active%d/ready%d/participation%d/hosting%t/headless%t/allow%t/stage%d", op, active, ready, participation, hosting, headless, allow, stage)
									t.Run(name, func(t *testing.T) {
										flags := noxflags.GameFlag(0)
										if hosting {
											flags = 1
										}
										restore := noxflags.PortTestGameFlags(flags)
										defer restore()
										noxflags.UnsetEngine(noxflags.EngineNoRendering)
										if headless {
											noxflags.SetEngine(noxflags.EngineNoRendering)
										}
										questAllowDefault = allow
										*memmap.PtrUint32(0x587000, 202028) = stage
										for i := range o.units {
											u := &o.units[i]
											pl := u.UpdateDataPlayer().Player
											pl.Active = 0
											pl.PlayerUnit = nil
											if active&(1<<uint(i)) != 0 {
												pl.Active = 1
												pl.PlayerUnit = u
											}
											objectXferSetWord(pl.C(), 4792, participation)
											objectXferSetWord(pl.C(), 4696, uint32(i*10))
											value := uint32(0)
											if ready&(1<<uint(i)) != 0 {
												value = uint32(uintptr(u.CObj()))
											}
											objectXferSetWord(u.UpdateData, 312, value)
											objectXferSetWord(u.UpdateData, 316, value)
										}
										considered := active
										if hosting && headless {
											considered &^= 4
										}
										if participation == 0 {
											considered = 0
										}
										want := uint32(0)
										allReady := considered != 0 && considered&ready == considered
										threshold := (stage/5 + 1) * 5
										if stage >= 20 {
											threshold = stage
										}
										eligible := allow
										for i := 0; i < 3; i++ {
											if considered&(1<<uint(i)) != 0 && uint32(i*10) >= threshold {
												eligible = true
											}
										}
										if allReady && (op == 6 || eligible) {
											want = 1
										}
										rv := legacy.PortTestWorldCollision(op, nil, nil, nil)
										if rv != want {
											t.Fatalf("return%d want%d", rv, want)
										}
										rows = append(rows, struct {
											Name   string
											Return uint32
										}{name, rv})
									})
								}
							}
						}
					}
				}
			}
		}
	}
}
