//go:build porttest

package opennox

import (
	"fmt"
	"image"
	"reflect"
	"testing"
	"unsafe"

	"github.com/opennox/libs/object"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

func TestPrefabRuntimeGroupTraversal(t *testing.T) {
	s := newObjectXferOwner(t)
	s.MapGroups.Init()
	defer s.MapGroups.Free()
	ids := map[uint32]uint32{}
	for i := 0; i < 3; i++ {
		u := s.NewObjectByTypeInd(1)
		if u == nil {
			t.Fatal("object allocation")
		}
		u.Extent = uint32(101 + i)
		u.ObjNext = s.Objs.List
		s.Objs.List = u
		if i == 2 {
			u.ObjFlags |= object.FlagDestroyed
		}
		ids[uint32(uintptr(u.CObj()))] = u.Extent
	}
	defer func() {
		for s.Objs.List != nil {
			u := s.Objs.List
			s.Objs.List = u.ObjNext
			s.Objs.FreeObject(u)
		}
	}()
	for i := 0; i < 2; i++ {
		w := s.NewWaypoint(types.Ptf(float32(i)*100, 0))
		w.Index = uint32(201 + i)
		ids[uint32(uintptr(unsafe.Pointer(w)))] = w.Index
	}
	defer s.Nox_xxx_waypointDeleteAll_579DD0()
	for i, p := range []image.Point{image.Pt(50, 50), image.Pt(51, 51), image.Pt(60, 60)} {
		w := s.Walls.CreateAtGrid(p)
		if w == nil {
			t.Fatal("wall allocation")
		}
		ids[uint32(uintptr(unsafe.Pointer(w)))] = uint32(301 + i)
	}
	group := func(id uint32, kind byte, items [][2]uint32) *server.MapGroup {
		if s.MapGroups.MapLoadAddGroup57C0C0(fmt.Sprint(id), id, kind) != 1 {
			t.Fatal("group allocation")
		}
		for i := len(items) - 1; i >= 0; i-- {
			if s.MapGroups.Sub57C130(items[i][:], id) != 1 {
				t.Fatal("group item")
			}
		}
		return s.MapGroups.GroupByInd(int(id))
	}
	groups := map[uint32]*server.MapGroup{0: nil}
	groups[10] = group(10, 0, [][2]uint32{{101, 0}, {999, 0}, {102, 0}, {101, 0}, {103, 0}})
	groups[20] = group(20, 1, [][2]uint32{{201, 0}, {999, 0}, {202, 0}, {201, 0}})
	groups[50] = group(50, 2, [][2]uint32{{60, 60}})
	groups[30] = group(30, 2, [][2]uint32{{50, 50}, {51, 51}, {99, 99}})
	groups[40] = group(40, 3, [][2]uint32{{10, 0}, {20, 0}, {30, 0}, {999, 0}, {10, 0}})
	groups[41] = group(41, 3, [][2]uint32{{40, 0}, {20, 0}})
	groups[42] = group(42, 4, nil)
	groups[43] = group(43, 0, nil)
	type row struct {
		Name  string
		Calls [][2]uint32
	}
	var rows []row
	for _, id := range []uint32{0, 10, 20, 30, 40, 41, 42, 43, 50} {
		for _, expected := range []int32{-1, 0, 1, 2, 3, 255} {
			for _, data := range []uint32{0, 1, 0x80000000, 0xffffffff} {
				callback := legacy.PortTestPrefabObserver()
				var raw uint32
				if groups[id] != nil {
					raw = uint32(uintptr(groups[id].C()))
				}
				legacy.PortTestPrefabCall(0, [6]uint32{raw, uint32(expected), uint32(uintptr(callback)), data})
				calls := legacy.PortTestPrefabCalls()
				for i := range calls {
					v, ok := ids[calls[i][0]]
					if !ok {
						t.Fatal("callback on unknown owner")
					}
					calls[i][0] = v
				}
				var wantIDs []uint32
				switch id {
				case 10:
					if expected == 0 {
						wantIDs = []uint32{101, 102, 101}
					}
				case 20:
					if expected == 1 {
						wantIDs = []uint32{201, 202, 201}
					}
				case 30:
					if expected == 2 {
						wantIDs = []uint32{301, 302, 303}
					}
				case 50:
					if expected == 2 {
						wantIDs = []uint32{303}
					}
				case 40, 41:
					switch expected {
					case 0:
						wantIDs = []uint32{101, 102, 101, 101, 102, 101}
					case 1:
						wantIDs = []uint32{201, 202, 201}
					case 2:
						wantIDs = []uint32{301, 302, 303}
					}
					if id == 41 && expected == 1 {
						wantIDs = append(wantIDs, 201, 202, 201)
					}
				}
				var want [][2]uint32
				for _, v := range wantIDs {
					want = append(want, [2]uint32{v, data})
				}
				if !reflect.DeepEqual(calls, want) {
					t.Fatalf("group%d/type%d calls%v want%v", id, expected, calls, want)
				}
				rows = append(rows, row{fmt.Sprintf("group%d/type%d/data%x", id, expected, data), calls})
			}
		}
	}
	spellbookCapture(t, "prefab-runtime-group-traversal", rows, "db758da6df8102b5bc297afdf42452240aeaf9ce8236d875aaa0385dc3955006")
}
