//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/object"
	"github.com/opennox/libs/prand"
	"github.com/opennox/opennox/v1/internal/netlist"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestObjectReportsMonsterAnimation(t *testing.T) {
	s, _, _ := objectReportsPlayers(t)
	t.Cleanup(s.PortTestObjectReportAnimations(5, 2))
	u := newCreatureXferObject(t, s, "Monster")
	u.HealthData = nil
	u.Direction1 = 0
	ud := u.UpdateDataMonster()
	var rows []struct {
		Name   string
		Return int
		Packet []byte
		Random int
	}
	defer func() {
		spellbookCapture(t, "object-reports-monster-animation", rows, "d306a19e408363b8a29b4f74492a65fe926fa217d9a9582d8c8a143d43d9b4a0")
	}()
	animations := map[int]byte{9: 12, 16: 1, 17: 3, 18: 7, 19: 7, 20: 7, 21: 5, 23: 5, 22: 6, 24: 13, 30: 9, 31: 10, 33: 14, 35: 14, 34: 15}
	walking := map[int]bool{7: true, 8: true, 10: true, 13: true, 29: true, 36: true, 37: true}
	for act := -1; act < 38; act++ {
		for _, moving := range []bool{false, true} {
			for _, npc := range []bool{false, true} {
				for _, confused := range []bool{false, true} {
					for _, seed := range []int{1, 7, 12345} {
						name := fmt.Sprintf("act%d/moving%t/npc%t/confused%t/seed%d", act, moving, npc, confused, seed)
						t.Run(name, func(t *testing.T) {
							s.NetList.ResetAll()
							s.Rand.Logic = prand.New(seed)
							expected := prand.New(seed)
							ud.AIStackInd = -1
							if act >= 0 {
								ud.AIStackInd = 0
								ud.AIStack[0].Action = uint32(act)
							}
							ud.StatusFlags = 0
							if moving {
								ud.StatusFlags = 0x4000
							}
							u.ObjSubClass = 0
							if npc {
								u.ObjSubClass = object.SubClass(0x10)
							}
							ud.Field523_2 = 0
							if confused {
								ud.Field523_2 = 1
							}
							ud.Field120_1 = 91
							rv := legacy.PortTestObjectReports(3, nil, u, 1, 0, 0, nil)
							got := s.NetList.CopyPacketsA(1, netlist.Kind2)
							anim := byte(8)
							if v, ok := animations[act]; ok {
								anim = v
							}
							if walking[act] {
								anim = 12
								if moving {
									anim = 13
								}
							}
							frame := byte(91)
							if npc && confused && expected.IntClamp(0, 10) >= 8 {
								anim = 14
								frame = byte(expected.IntClamp(0, 5))
							}
							if rv != 1 || len(got) != 11 {
								t.Fatalf("return%d packet%x", rv, got)
							}
							if got[9] != (64|anim) || got[10] != frame || s.Rand.Logic.Index() != expected.Index() {
								t.Fatalf("packet%x anim%d frame%d rng%d want%d", got, anim, frame, s.Rand.Logic.Index(), expected.Index())
							}
							rows = append(rows, struct {
								Name   string
								Return int
								Packet []byte
								Random int
							}{name, rv, got, s.Rand.Logic.Index()})
						})
					}
				}
			}
		}
	}
}
