//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/object"
	"github.com/opennox/libs/prand"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestVisibilityEffectsSpatialScan(t *testing.T) {
	s, u, others, calls, setup := visibilitySeenOwner(t)
	s.PortTestAIEmptyMap()
	t.Cleanup(s.Map.Free)
	s.SetTickRate(30)
	copyFrame := memmap.PtrUint32(0x5D4594, 2487684)
	old := *copyFrame
	*copyFrame = 0
	t.Cleanup(func() { *copyFrame = old })
	u.PosVec = types.Pointf{100, 100}
	u.PrevPos = u.PosVec
	u.TeamPtr().ID = 1
	type row struct {
		State        visibilitySeenRow
		Times        [4]uint32
		Logic, Other int
	}
	var rows []row
	defer func() {
		spellbookCapture(t, "visibility-effects-spatial-scan", rows, "194c6546cbb13135f451b5af00ce110ec7292c538e635a3308ea91deff2c5813")
	}()
	for _, n := range []int{0, 1, 3} {
		for _, force := range []bool{false, true} {
			for _, due := range []bool{false, true} {
				name := fmt.Sprintf("n%d/force%t/due%t", n, force, due)
				t.Run(name, func(t *testing.T) {
					s.PortTestAIEmptyMap()
					setup(0)
					s.SetFrame(100)
					s.Rand.Logic = prand.New(12345)
					s.Rand.Other = prand.New(54321)
					objectXferSetWord(u.UpdateData, 1440, 0x500)
					objectXferSetWord(u.UpdateData, 536, 100)
					objectXferSetWord(u.UpdateData, 1208, 101)
					if due {
						objectXferSetWord(u.UpdateData, 1208, 100)
					}
					*copyFrame = 0
					if force {
						*copyFrame = 100
					}
					for i := 0; i < n; i++ {
						v := others[i]
						v.TeamPtr().ID = 2
						v.PosVec = types.Pointf{float32(150 + i*25), 100}
						v.NewPos = v.PosVec
						v.ObjFlags = object.FlagActive
						s.Map.AddObjectToIndex(v)
					}
					rv := legacy.PortTestVisibilityEffects(17, u, nil, nil, nil, [5]int32{}, nil, "")
					r := row{State: visibilitySeenCapture(name, u, rv, *calls), Logic: s.Rand.Logic.Index(), Other: s.Rand.Other.Index()}
					for i, off := range []int{1204, 1208, 1212, 536} {
						r.Times[i] = objectXferGetWord(u.UpdateData, off)
					}
					scanned := due || force
					count := 0
					if scanned {
						count = n
					}
					if r.State.Count != byte(count) || len(r.State.Calls) != count*2 {
						t.Fatalf("scan %+v", r)
					}
					if scanned {
						expectedRNG := prand.New(12345)
						expectedDeadline := uint32(100 + expectedRNG.IntClamp(5, 10))
						if r.Logic != expectedRNG.Index() || r.Times[0] != 100 || r.Times[2] != 100 || r.Times[1] != expectedDeadline {
							t.Fatalf("timestamps%v", r.Times)
						}
						if n > 0 && r.State.Target != 1 {
							t.Fatalf("nearest target%d", r.State.Target)
						}
					} else if r.Times[0] != 0 || r.Times[1] != 101 {
						t.Fatalf("deferred timestamps%v", r.Times)
					}
					if r.Times[3] != 100 || r.Other != prand.New(54321).Index() {
						t.Fatalf("unrelated cooldown/random changed%+v", r)
					}
					rows = append(rows, r)
				})
			}
		}
	}
}
func TestVisibilityEffectsGlobalRemoval(t *testing.T) {
	s, u, others, calls, setup := visibilitySeenOwner(t)
	v := others[0]
	u.TeamPtr().ID = 1
	v.TeamPtr().ID = 2
	s.Objs.AddToUpdatable(u)
	t.Cleanup(func() { s.Objs.RemoveFromUpdatable(u) })
	var rows []visibilitySeenRow
	defer func() {
		spellbookCapture(t, "visibility-effects-global-removal", rows, "e8d0482e1df182a9b56425629ffaef2953e3499e85ff72beebb47f8e6fefb87b")
	}()
	for _, kind := range []string{"monster", "destroyed", "simple", "absent"} {
		t.Run(kind, func(t *testing.T) {
			setup(1)
			u.ObjClass = object.ClassMonster
			u.ObjFlags = 0
			target := v
			switch kind {
			case "destroyed":
				u.ObjFlags = object.FlagDestroyed
			case "simple":
				u.ObjClass = object.ClassSimple
			case "absent":
				target = others[17]
			}
			rv := legacy.PortTestVisibilityEffects(24, target, nil, nil, nil, [5]int32{}, nil, "")
			r := visibilitySeenCapture(kind, u, rv, *calls)
			want := byte(1)
			if kind == "monster" {
				want = 0
			}
			if r.Count != want || rv != 0 {
				t.Fatalf("state%+v", r)
			}
			rows = append(rows, r)
		})
	}
}
