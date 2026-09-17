//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/prand"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"testing"
	"unsafe"
)

func TestVisibilityEffectsSeenCooldown(t *testing.T) {
	s, u, others, calls, setup := visibilitySeenOwner(t)
	u.TeamPtr().ID = 1
	v := others[0]
	sounds, free := alloc.Make([]uint32{}, 19)
	sounds[17] = 317
	t.Cleanup(free)
	type row struct {
		State        visibilitySeenRow
		Cooldown     uint32
		Logic, Other int
		Sounds       []uint32
	}
	var rows []row
	defer func() {
		spellbookCapture(t, "visibility-effects-seen-cooldown", rows, "d5dc781cd3d9b23eb7b2c00c86875baed9ce4e294b1c15d9d6d956ee547ba4e7")
	}()
	for _, frame := range []uint32{0, 1, 99, 100, 101, 0xfffffffe, 0xffffffff} {
		for _, last := range []uint32{0, 100, 0xffffffff} {
			for _, enemy := range []bool{false, true} {
				for _, sound := range []bool{false, true} {
					name := fmt.Sprintf("frame%x/last%x/enemy%t/sound%t", frame, last, enemy, sound)
					t.Run(name, func(t *testing.T) {
						setup(0)
						s.SetFrame(frame)
						s.SetTickRate(30)
						s.Rand.Logic = prand.New(12345)
						s.Rand.Other = prand.New(54321)
						s.PortTestCombatAudioReset()
						v.TeamPtr().ID = 1
						if enemy {
							v.TeamPtr().ID = 2
						}
						objectXferSetWord(u.UpdateData, 536, last)
						if sound {
							*(*unsafe.Pointer)(unsafe.Add(u.UpdateData, 488)) = unsafe.Pointer(&sounds[0])
						}
						rv := legacy.PortTestVisibilityEffects(21, u, v, nil, nil, [5]int32{}, nil, "")
						r := row{State: visibilitySeenCapture(name, u, rv, *calls), Cooldown: objectXferGetWord(u.UpdateData, 536), Logic: s.Rand.Logic.Index(), Other: s.Rand.Other.Index()}
						expectedRNG := prand.New(12345)
						want := last
						triggered := frame > last && enemy
						if triggered {
							want = frame + uint32(expectedRNG.IntClamp(60, 120))
						}
						if r.Cooldown != want || r.Logic != expectedRNG.Index() || r.Other != prand.New(54321).Index() {
							t.Fatalf("cooldown/random%+v want%x", r, want)
						}
						events := s.PortTestCombatAudioSnapshot()
						n := 0
						if triggered && sound {
							n = 1
						}
						if len(events) != n {
							t.Fatalf("audio%v", events)
						}
						for _, ev := range events {
							if ev.ID != 317 || ev.Obj != u || ev.Kind != 0 || ev.Code != 0 || ev.ByPos {
								t.Fatalf("event%+v", ev)
							}
							r.Sounds = append(r.Sounds, uint32(ev.ID))
						}
						if r.State.Count != 1 || len(r.State.Calls) != 2 {
							t.Fatalf("seen%+v", r.State)
						}
						rows = append(rows, r)
					})
				}
			}
		}
	}
}
