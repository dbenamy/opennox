//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/script"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc/handles"
	"github.com/opennox/opennox/v1/server"
	"reflect"
	"testing"
	"unsafe"
)

type visibilitySeenRow struct {
	Name             string
	Return           uint32
	Count            byte
	Slots            []uint32
	Target, Previous uint32
	Calls            [][4]uint32
}

func visibilitySeenOwner(t *testing.T) (s *server.Server, u *server.Object, others []*server.Object, calls *[][4]uint32, setup func(int)) {
	t.Cleanup(handles.PortTestInit())
	s = newCreatureXferOwner(t)
	// These legacy type caches belong to the current server's type table.
	for _, off := range []uintptr{2488532, 2488536} {
		p := memmap.PtrUint32(0x5D4594, off)
		old := *p
		*p = 0
		t.Cleanup(func() { *p = old })
	}
	s.NoxScriptVM.Init(s)
	noxServer.noxScript.Init(noxServer)
	s.SetFrame(0)
	out := memmap.PtrUint32(0x5D4594, 1599076)
	old := *out
	t.Cleanup(func() { *out = old })
	u = newCreatureXferObject(t, s, "Monster")
	u.NetCode = 100
	for i := 0; i < 18; i++ {
		v := newCreatureXferObject(t, s, "Monster")
		v.NetCode = uint32(i + 1)
		v.PosVec = types.Pointf{float32(i + 1), 0}
		others = append(others, v)
	}
	calls = new([][4]uint32)
	count := (*byte)(unsafe.Add(u.UpdateData, 1129))
	u.OnUnitSeeEnemy(func(script.Unit) { *calls = append(*calls, [4]uint32{1, 14, uint32(*count), 0}) })
	u.OnUnitLostEnemy(func(script.Unit) { *calls = append(*calls, [4]uint32{1, 15, uint32(*count), 0}) })
	indices := [2]int{}
	for i, event := range []uint32{14, 15} {
		indices[i] = s.NoxScriptVM.AsFuncIndex(fmt.Sprintf("visibility%d", event), func() {
			*calls = append(*calls, [4]uint32{2, event, s.NoxScriptVM.Caller().NetCode, s.NoxScriptVM.Trigger().NetCode})
		})
	}
	setup = func(n int) {
		clear(unsafe.Slice((*byte)(u.UpdateData), int(unsafe.Sizeof(server.MonsterUpdateData{}))))
		for i, off := range []uintptr{1232, 1296} {
			cb := (*server.ScriptCallback)(unsafe.Add(u.UpdateData, off))
			cb.Func = int32(indices[i])
		}
		*count = byte(n)
		slots := unsafe.Slice((**server.Object)(unsafe.Add(u.UpdateData, 1132)), 16)
		for i := 0; i < n; i++ {
			slots[i] = others[i]
		}
		*calls = nil
	}
	return
}
func visibilitySeenCapture(name string, u *server.Object, rv uint32, calls [][4]uint32) visibilitySeenRow {
	id := func(v *server.Object) uint32 {
		if v == nil {
			return 0
		}
		return v.NetCode
	}
	r := visibilitySeenRow{Name: name, Return: rv, Count: *(*byte)(unsafe.Add(u.UpdateData, 1129)), Target: id(*(**server.Object)(unsafe.Add(u.UpdateData, 1196))), Previous: objectXferGetWord(u.UpdateData, 1200), Calls: append([][4]uint32(nil), calls...)}
	for _, v := range unsafe.Slice((**server.Object)(unsafe.Add(u.UpdateData, 1132)), 16) {
		r.Slots = append(r.Slots, id(v))
	}
	return r
}
func TestVisibilityEffectsSeenRemoval(t *testing.T) {
	_, u, others, calls, setup := visibilitySeenOwner(t)
	var rows []visibilitySeenRow
	defer func() {
		spellbookCapture(t, "visibility-effects-seen-removal", rows, "e593003d2d14121191be94d98b5138771c80a9967ae0507a155e083feb02de6d")
	}()
	for n := 0; n <= 16; n++ {
		for target := 0; target <= n; target++ {
			for _, selected := range []bool{false, true} {
				for _, op := range []int{18, 22, 23} {
					if op == 18 && target == n {
						continue
					}
					name := fmt.Sprintf("n%d/target%d/selected%t/op%d", n, target, selected, op)
					t.Run(name, func(t *testing.T) {
						setup(n)
						v := others[target]
						if selected {
							*(**server.Object)(unsafe.Add(u.UpdateData, 1196)) = v
						}
						rv := legacy.PortTestVisibilityEffects(op, u, v, nil, nil, [5]int32{int32(target)}, nil, "")
						present := target < n
						wantCount := n
						wantReturn := uint32(0)
						removed := present && op != 23
						if op == 23 {
							if present {
								wantReturn = 1
							}
						} else if !present {
							wantReturn = uint32(n)
						} else {
							wantCount--
							wantReturn = uint32(wantCount)
							if target < wantCount {
								wantReturn = uint32(uintptr(unsafe.Add(u.UpdateData, 1132+4*wantCount)))
							}
						}
						if rv != wantReturn {
							t.Fatalf("raw return %x want %x", rv, wantReturn)
						}
						// Pointer returns are asserted against the owned buffer before normalization.
						norm := rv
						if removed && target < wantCount {
							norm = 0xF0000000 | uint32(wantCount)
						}
						r := visibilitySeenCapture(name, u, norm, *calls)
						if int(r.Count) != wantCount {
							t.Fatalf("count%d want%d", r.Count, wantCount)
						}
						wantSlots := make([]uint32, 16)
						for i := 0; i < n; i++ {
							wantSlots[i] = uint32(i + 1)
						}
						if removed {
							copy(wantSlots[target:n-1], wantSlots[target+1:n])
						}
						if !reflect.DeepEqual(r.Slots, wantSlots) {
							t.Fatalf("slots%v want%v", r.Slots, wantSlots)
						}
						var wantCalls [][4]uint32
						if removed {
							wantCalls = [][4]uint32{{1, 15, uint32(n), 0}, {2, 15, v.NetCode, u.NetCode}}
						}
						if !reflect.DeepEqual(r.Calls, wantCalls) {
							t.Fatalf("calls%v want%v", r.Calls, wantCalls)
						}
						wantTarget, wantPrevious := uint32(0), uint32(0)
						if selected {
							wantTarget = v.NetCode
							if removed {
								wantTarget = 0
								wantPrevious = v.NetCode
							}
						}
						if r.Target != wantTarget || r.Previous != wantPrevious {
							t.Fatalf("selection %d/%d want%d/%d", r.Target, r.Previous, wantTarget, wantPrevious)
						}
						rows = append(rows, r)
					})
				}
			}
		}
	}
}
func TestVisibilityEffectsSeenInsertion(t *testing.T) {
	_, u, others, calls, setup := visibilitySeenOwner(t)
	var rows []visibilitySeenRow
	defer func() {
		spellbookCapture(t, "visibility-effects-seen-insertion", rows, "44ecd46e0b25eb4252ad53528425223ba67843270e885bd797f43b50ecb539a8")
	}()
	for n := 0; n <= 16; n++ {
		for _, distance := range []float32{0, 1, 15, 16, 17, 100} {
			name := fmt.Sprintf("n%d/distance%g", n, distance)
			t.Run(name, func(t *testing.T) {
				setup(n)
				v := others[17]
				v.PosVec = types.Pointf{distance, 0}
				rv := legacy.PortTestVisibilityEffects(21, u, v, nil, nil, [5]int32{}, nil, "")
				r := visibilitySeenCapture(name, u, rv, *calls)
				added := n < 16 || distance < 16
				wantCount := n
				if added && n < 16 {
					wantCount++
				}
				if int(r.Count) != wantCount {
					t.Fatalf("count%d want%d", r.Count, wantCount)
				}
				wantSlots := make([]uint32, 16)
				for i := 0; i < n; i++ {
					wantSlots[i] = uint32(i + 1)
				}
				if added {
					wantSlots[wantCount-1] = 18
				}
				if !reflect.DeepEqual(r.Slots, wantSlots) {
					t.Fatalf("slots%v want%v", r.Slots, wantSlots)
				}
				var wantCalls [][4]uint32
				if added {
					if n == 16 {
						wantCalls = append(wantCalls, [4]uint32{1, 15, 16, 0}, [4]uint32{2, 15, 16, 100})
					}
					wantCalls = append(wantCalls, [4]uint32{1, 14, uint32(wantCount), 0}, [4]uint32{2, 14, 18, 100})
				}
				if !reflect.DeepEqual(r.Calls, wantCalls) {
					t.Fatalf("calls%v want%v", r.Calls, wantCalls)
				}
				rows = append(rows, r)
			})
		}
	}
}
