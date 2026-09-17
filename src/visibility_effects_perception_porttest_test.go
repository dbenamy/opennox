//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/object"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/memmap/nox/blobdata"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

func TestVisibilityEffectsTargetSelection(t *testing.T) {
	_, u, others, calls, setup := visibilitySeenOwner(t)
	u.TeamPtr().ID = 1
	var rows []visibilitySeenRow
	defer func() {
		spellbookCapture(t, "visibility-effects-target-selection", rows, "ccedbdfaee008b670bbc391d926a11b1e1ca9084e1e139aca494efc0cc784190")
	}()
	for n := 0; n <= 16; n++ {
		for _, preferred := range []int{-1, 0, 7, 15, 17} {
			for _, mode := range []int{0, 1, 2, 3} {
				name := fmt.Sprintf("n%d/preferred%d/mode%d", n, preferred, mode)
				t.Run(name, func(t *testing.T) {
					setup(n)
					for i, v := range others {
						v.TeamPtr().ID = 2
						v.HealthData.Cur = 80
						v.HealthData.Max = 100
						v.PosVec = types.Pointf{float32(18 - i), 0}
						if mode == 1 && i%2 == 0 {
							v.TeamPtr().ID = 1
						}
						if mode == 2 && i%2 == 0 {
							v.HealthData.Cur = 0
						}
						if mode == 3 {
							v.TeamPtr().ID = 1
						}
					}
					if preferred >= 0 {
						*(**server.Object)(unsafe.Add(u.UpdateData, 1216)) = others[preferred]
					}
					*(**server.Object)(unsafe.Add(u.UpdateData, 1196)) = others[17]
					rv := legacy.PortTestVisibilityEffects(19, u, nil, nil, nil, [5]int32{}, nil, "")
					r := visibilitySeenCapture(name, u, rv, *calls)
					want := uint32(0)
					for i := 0; i < n; i++ {
						if i == preferred {
							want = others[i].NetCode
							break
						}
						eligible := mode != 3 && !(mode == 1 && i%2 == 0) && !(mode == 2 && i%2 == 0)
						if eligible {
							want = others[i].NetCode
						}
					}
					if r.Target != want || r.Count != byte(n) || len(r.Calls) != 0 {
						t.Fatalf("selection %+v want%d", r, want)
					}
					rows = append(rows, r)
				})
			}
		}
	}
}
func TestVisibilityEffectsPerceptionCandidate(t *testing.T) {
	s, u, others, calls, setup := visibilitySeenOwner(t)
	s.PortTestAIEmptyMap()
	t.Cleanup(s.Map.Free)
	for off, data := range blobdata.PortTestCombatTables() {
		dst := unsafe.Slice((*byte)(memmap.PtrOff(0x587000, off)), len(data))
		old := append([]byte(nil), dst...)
		copy(dst, data)
		t.Cleanup(func() { copy(dst, old) })
	}
	u.PosVec = types.Pointf{100, 100}
	u.Direction1 = 0
	u.TeamPtr().ID = 1
	v := others[0]
	v.TeamPtr().ID = 2
	var rows []visibilitySeenRow
	defer func() {
		spellbookCapture(t, "visibility-effects-candidates", rows, "f89f2ce15f5fd5167381de92d7935136e0c6d13a722a84125796b9e63d6cb756")
	}()
	for _, status := range []uint32{0, 0x100, 0x400, 0x500} {
		for _, kind := range []string{"front", "behind", "side", "self", "simple", "destroyed", "dead", "friendly", "known", "invisible"} {
			name := fmt.Sprintf("status%x/%s", status, kind)
			t.Run(name, func(t *testing.T) {
				setup(0)
				objectXferSetWord(u.UpdateData, 1440, status)
				v.ObjClass = object.ClassMonster
				v.ObjFlags = 0
				v.Buffs = 0
				v.TeamPtr().ID = 2
				v.PosVec = types.Pointf{150, 100}
				candidate := v
				switch kind {
				case "behind":
					v.PosVec = types.Pointf{50, 100}
				case "side":
					v.PosVec = types.Pointf{100, 150}
				case "self":
					candidate = u
				case "simple":
					v.ObjClass = object.ClassSimple
				case "destroyed":
					v.ObjFlags = object.Flags(0x20)
				case "dead":
					v.ObjFlags = object.Flags(0x8000)
				case "friendly":
					v.TeamPtr().ID = 1
				case "known":
					setup(1)
					objectXferSetWord(u.UpdateData, 1440, status)
				case "invisible":
					v.Buffs = 1 << server.ENCHANT_INVISIBLE
				}
				rv := legacy.PortTestVisibilityEffects(20, candidate, u, nil, nil, [5]int32{}, nil, "")
				r := visibilitySeenCapture(name, u, rv, *calls)
				added := kind != "self" && kind != "simple" && kind != "destroyed" && kind != "dead" && kind != "known" && kind != "invisible"
				if kind == "friendly" && status&0x400 == 0 {
					added = false
				}
				if (kind == "behind" || kind == "side") && status&0x100 == 0 {
					added = false
				}
				count := byte(0)
				if added || kind == "known" {
					count = 1
				}
				if r.Count != count {
					t.Fatalf("count%d want%d", r.Count, count)
				}
				callCount := 0
				if added {
					callCount = 2
				}
				if len(r.Calls) != callCount {
					t.Fatalf("calls%v", r.Calls)
				}
				rows = append(rows, r)
			})
		}
	}
}
func TestVisibilityEffectsSightMaintenance(t *testing.T) {
	s, u, others, calls, setup := visibilitySeenOwner(t)
	s.PortTestAIEmptyMap()
	t.Cleanup(s.Map.Free)
	s.SetTickRate(30)
	s.SetFrame(100)
	frameCopy := memmap.PtrUint32(0x5D4594, 2487684)
	old := *frameCopy
	*frameCopy = 0
	t.Cleanup(func() { *frameCopy = old })
	u.TeamPtr().ID = 1
	u.PosVec = types.Pointf{100, 100}
	u.PrevPos = u.PosVec
	v := others[0]
	v.TeamPtr().ID = 2
	var rows []visibilitySeenRow
	defer func() {
		spellbookCapture(t, "visibility-effects-sight-maintenance", rows, "8435200c73359f3374ac76ee7be8c068104fb1b6a7c304e61e3b7f6487a4bf79")
	}()
	for _, kind := range []string{"visible", "destroyed", "dead", "invisible", "range-edge", "range-outside", "moved", "actor-dead"} {
		t.Run(kind, func(t *testing.T) {
			setup(1)
			u.ObjFlags = 0
			u.PrevPos = u.PosVec
			v.ObjFlags = 0
			v.Buffs = 0
			v.PosVec = types.Pointf{150, 100}
			objectXferSetWord(u.UpdateData, 1204, 99)
			objectXferSetWord(u.UpdateData, 1208, 1000)
			objectXferSetWord(u.UpdateData, 1212, 99)
			*(**server.Object)(unsafe.Add(u.UpdateData, 1196)) = v
			switch kind {
			case "destroyed":
				v.ObjFlags = object.Flags(0x20)
			case "dead":
				v.ObjFlags = object.Flags(0x8000)
			case "invisible":
				v.Buffs = 1 << server.ENCHANT_INVISIBLE
			case "range-edge":
				v.PosVec = types.Pointf{380, 100}
			case "range-outside":
				v.PosVec = types.Pointf{381, 100}
			case "moved":
				u.PrevPos = types.Pointf{68, 100}
			case "actor-dead":
				u.ObjFlags = object.Flags(0x8000)
			}
			rv := legacy.PortTestVisibilityEffects(17, u, nil, nil, nil, [5]int32{}, nil, "")
			r := visibilitySeenCapture(kind, u, rv, *calls)
			removed := kind != "visible" && kind != "range-edge" && kind != "actor-dead"
			want := byte(1)
			if removed {
				want = 0
			}
			if r.Count != want {
				t.Fatalf("count%d want%d", r.Count, want)
			}
			if removed {
				if len(r.Calls) != 2 || r.Target != 0 || r.Previous != v.NetCode {
					t.Fatalf("lost state%+v", r)
				}
			} else if len(r.Calls) != 0 || r.Target != v.NetCode {
				t.Fatalf("kept state%+v", r)
			}
			rows = append(rows, r)
		})
	}
}
