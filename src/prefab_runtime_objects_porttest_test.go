//go:build porttest

package opennox

import (
	"fmt"
	"math"
	"reflect"
	"testing"
	"unsafe"

	"github.com/opennox/libs/object"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func prefabRuntimePaths(t *testing.T, words map[string]*uint32) {
	t.Helper()
	for _, key := range []string{"path", "alternate"} {
		p, free := alloc.Make([]byte{}, 2048)
		t.Cleanup(free)
		*words[key] = uint32(uintptr(unsafe.Pointer(&p[0])))
	}
}
func TestPrefabRuntimeObjectList(t *testing.T) {
	s := newObjectXferOwner(t)
	words, restore := legacy.PortTestPrefabRuntimeGlobals()
	defer restore()
	type row struct {
		Name   string
		Return uint64
		Nodes  [][3]uint32
		Links  [][3]uint32
	}
	var rows []row
	for _, order := range [][3]int{{0, 1, 2}, {0, 2, 1}, {1, 0, 2}, {1, 2, 0}, {2, 0, 1}, {2, 1, 0}} {
		objects := make([]*server.Object, 4)
		nodes := make([]unsafe.Pointer, 3)
		ids := map[uint32]uint32{0: 0}
		alive := [3]bool{true, true, true}
		for i := range objects {
			objects[i] = s.NewObjectByTypeInd(1)
			if objects[i] == nil {
				t.Fatal("object allocation")
			}
			ids[uint32(uintptr(objects[i].CObj()))] = uint32(100 + i)
		}
		norm := func(v uint32) uint32 {
			n, ok := ids[v]
			if !ok {
				t.Fatalf("unknown list pointer %x", v)
			}
			return n
		}
		snapshot := func(name string, ret uint64) {
			r := row{Name: name, Return: ret}
			for p := *words["objects"]; p != 0; p = *(*uint32)(unsafe.Pointer(uintptr(p + 4))) {
				a := *(*[3]uint32)(unsafe.Pointer(uintptr(p)))
				for i := range a {
					a[i] = norm(a[i])
				}
				r.Nodes = append(r.Nodes, a)
			}
			for i, u := range objects[:3] {
				if alive[i] {
					a := *(*[2]uint32)(unsafe.Add(u.CObj(), 444))
					r.Links = append(r.Links, [3]uint32{uint32(100 + i), norm(a[0]), norm(a[1])})
				}
			}
			rows = append(rows, r)
		}
		for i, u := range objects[:3] {
			nodes[i] = legacy.PortTestPrefabObjectNode(u.CObj())
			if nodes[i] == nil {
				t.Fatal("cache allocation")
			}
			ids[uint32(uintptr(nodes[i]))] = uint32(200 + i)
		}
		snapshot(fmt.Sprint(order, "/created"), 0)
		if got := legacy.PortTestPrefabCall(28, [6]uint32{}); norm(uint32(got)) != 202 {
			t.Fatal("cache head")
		}
		for i, u := range objects[:3] {
			wantObj, wantNode := uint32(0), uint32(0)
			if i > 0 {
				wantObj = uint32(100 + i - 1)
				wantNode = uint32(200 + i - 1)
			}
			if norm(uint32(legacy.PortTestPrefabCall(27, [6]uint32{uint32(uintptr(u.CObj()))}))) != wantObj || norm(uint32(legacy.PortTestPrefabCall(29, [6]uint32{uint32(uintptr(nodes[i]))}))) != wantNode {
				t.Fatal("cache iterators")
			}
		}
		if legacy.PortTestPrefabCall(27, [6]uint32{}) != 0 || legacy.PortTestPrefabCall(29, [6]uint32{}) != 0 {
			t.Fatal("nil iterators")
		}
		for _, state := range [][3]uint32{{0, 0, 0}, {0, 0, 1}, {0, 0, 2}, {1, 0, 0}, {math.MaxUint32, math.MaxUint32, 0}} {
			*words["loaded"], *words["selected"], *words["placed"] = state[0], state[1], state[2]
			got := legacy.PortTestPrefabCall(26, [6]uint32{})
			want := uint32(0)
			if state[0] == state[1] && state[0] != math.MaxUint32 && state[2] != 1 {
				want = 102
			}
			if norm(uint32(got)) != want {
				t.Fatal("cache admission")
			}
			snapshot(fmt.Sprint(order, "/admission", state), uint64(norm(uint32(got))))
		}
		if legacy.PortTestPrefabCall(30, [6]uint32{}) != 0 || legacy.PortTestPrefabCall(30, [6]uint32{uint32(uintptr(objects[3].CObj()))}) != 0 {
			t.Fatal("absent removal")
		}
		for _, i := range order {
			ret := legacy.PortTestPrefabCall(30, [6]uint32{uint32(uintptr(objects[i].CObj()))})
			if ret != 1 {
				t.Fatal("existing removal")
			}
			alive[i] = false
			snapshot(fmt.Sprint(order, "/remove", i), ret)
		}
		if *words["objects"] != 0 {
			t.Fatal("list not empty")
		}
		s.Objs.FreeObject(objects[3])
	}
	spellbookCapture(t, "prefab-runtime-object-list", rows, "513f6d2b61450d490e6bb1516d3811cddd9fd791966aa5c8db1e6f7a32f12b32")
}

func TestPrefabRuntimeObjectPlacement(t *testing.T) {
	s := newObjectXferOwner(t)
	words, restore := legacy.PortTestPrefabRuntimeGlobals()
	defer restore()
	prefabRuntimePaths(t, words)
	type row struct {
		Name      string
		Positions [][6]uint32
		Pending   []uint32
		Flags     []uint32
	}
	var rows []row
	for _, offset := range [][2]int32{{0, 0}, {-100, 99}, {1, -1}, {2147483647, -2147483648}} {
		objects := make([]*server.Object, 3)
		ids := map[*server.Object]uint32{}
		positions := []types.Pointf{{0, 0}, {1.5, -2.25}, {2000000000, -2000000000}}
		for i := range objects {
			u := s.NewObjectByTypeInd(1)
			if u == nil {
				t.Fatal("object allocation")
			}
			objects[i] = u
			ids[u] = uint32(i + 1)
			u.PosVec = positions[i]
			u.VelVec = types.Pointf{9, 10}
			u.ForceVec = types.Pointf{11, 12}
			legacy.PortTestPrefabObjectNode(u.CObj())
		}
		if legacy.PortTestPrefabCall(25, [6]uint32{uint32(offset[0]), uint32(offset[1])}) != 1 {
			t.Fatal("placement return")
		}
		r := row{Name: fmt.Sprint(offset)}
		for i, u := range objects {
			// Source offsets are converted to float before adding to each stored position.
			want := types.Pointf{float32(float64(float32(offset[0])) + float64(positions[i].X)), float32(float64(float32(offset[1])) + float64(positions[i].Y))}
			if u.PosVec != want || u.PrevPos != want || u.NewPos != want || u.VelVec != (types.Pointf{}) || u.ForceVec != (types.Pointf{}) {
				t.Fatalf("offset %v object%d placement got %v want %v", offset, i, u.PosVec, want)
			}
			if !u.ObjFlags.Has(object.FlagActive|object.FlagPending) || uint32(u.ObjFlags)&0x80000000 == 0 {
				t.Fatal("placement flags")
			}
			r.Positions = append(r.Positions, [6]uint32{math.Float32bits(u.PosVec.X), math.Float32bits(u.PosVec.Y), math.Float32bits(u.PrevPos.X), math.Float32bits(u.PrevPos.Y), math.Float32bits(u.NewPos.X), math.Float32bits(u.NewPos.Y)})
			r.Flags = append(r.Flags, uint32(u.ObjFlags))
		}
		for u := s.Objs.Pending; u != nil; u = u.ObjNext {
			r.Pending = append(r.Pending, ids[u])
		}
		if !reflect.DeepEqual(r.Pending, []uint32{1, 2, 3}) {
			t.Fatal("actual pending order")
		}
		rows = append(rows, r)
		*words["placed"] = 1
		noxServer.Nox_xxx_free503F40()
		s.Objs.Pending = nil
		for _, u := range objects {
			u.ObjNext = nil
			u.ObjPrev = nil
			s.Objs.FreeObject(u)
		}
	}
	spellbookCapture(t, "prefab-runtime-object-placement", rows, "926acfa9d64ae80722835b7e0577e65d9bfc77d1856ef291ab5416fdc4b91c96")
}
