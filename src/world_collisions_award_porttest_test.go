//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/object"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestWorldCollisionsSpellAward(t *testing.T) {
	o := newWorldCollisionOwner(t)
	configure, restore := o.s.PortTestMeterStrings()
	t.Cleanup(restore)
	configure(0)
	a := newObjectXferSimple(t, o.s)
	b := &o.units[1]
	other := newObjectXferSimple(t, o.s)
	data := o.record(t, 4)
	a.CollideData = data
	t.Cleanup(func() { a.CollideData = nil })
	var rows []struct {
		Name          string
		Return, Level uint32
		Sounds        []int
		Queue         int
	}
	defer func() {
		spellbookCapture(t, "world-collisions-spell-award", rows, "27cbf6527811a747cb52e1e4efd9fb60efc21929d0d3f3a889ab9fd66dd6abe5")
	}()
	for _, id := range []int{1, 19, 34, 50, 136} {
		for _, level := range []uint32{0, 1, 2, 3, 4, 5} {
			for _, quest := range []bool{false, true} {
				for _, player := range []bool{false, true} {
					name := fmt.Sprintf("id%d/level%d/quest%t/player%t", id, level, quest, player)
					t.Run(name, func(t *testing.T) {
						o.reset()
						o.s.PortTestCombatAudioReset()
						flags := noxflags.GameFlag(0)
						if quest {
							flags = noxflags.GameModeQuest
						}
						defer noxflags.PortTestGameFlags(flags)()
						objectXferSetWord(data, 0, uint32(id))
						p := b.UpdateDataPlayer().Player.C()
						objectXferSetWord(p, 3696+4*id, level)
						target := b
						if !player {
							target = other
							target.ObjClass = object.ClassSimple
						}
						rv := legacy.PortTestWorldCollision(16, a, target, nil)
						awarded := player && level != 5 && (!quest || level != 3 && (id != 19 && id != 34 || level == 0))
						want := level
						if awarded {
							want++
							if want > 5 {
								want = 5
							}
							if quest && want > 3 {
								want = 3
							}
						}
						if (rv == 1) != awarded || rv > 1 || objectXferGetWord(p, 3696+4*id) != want {
							t.Fatalf("award return%d level%d want%d success%t", rv, objectXferGetWord(p, 3696+4*id), want, awarded)
						}
						var sounds []int
						for _, e := range o.s.PortTestCombatAudioSnapshot() {
							sounds = append(sounds, int(e.ID))
						}
						if (len(sounds) == 1) != awarded || len(sounds) > 1 {
							t.Fatal("award sound count")
						}
						if awarded && sounds[0] != 226 {
							t.Fatal("award sound")
						}
						rows = append(rows, struct {
							Name          string
							Return, Level uint32
							Sounds        []int
							Queue         int
						}{name, rv, objectXferGetWord(p, 3696+4*id), sounds, len(o.state().Nodes)})
					})
				}
			}
		}
	}
	if legacy.PortTestWorldCollision(16, a, nil, nil) != 0 {
		t.Fatal("nil award")
	}
}
