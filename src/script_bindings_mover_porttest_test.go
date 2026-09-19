//go:build porttest

package opennox

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/opennox/libs/object"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func TestScriptBindingsMover(t *testing.T) {
	o := newWorldCollisionOwner(t)
	saved := append([]server.Object(nil), o.units...)
	data, free := alloc.Make([]server.MoverUpdateData{}, 3)
	t.Cleanup(free)
	wp, free := alloc.New(server.Waypoint{})
	t.Cleanup(free)
	oldList, oldUpdates := o.s.Objs.List, o.s.Objs.UpdatableList
	oldMover := legacy.Get_dword_5d4594_2386836()
	legacy.Set_dword_5d4594_2386836(1)
	t.Cleanup(func() {
		copy(o.units, saved)
		o.s.Objs.List = oldList
		o.s.Objs.UpdatableList = oldUpdates
		legacy.Set_dword_5d4594_2386836(oldMover)
	})
	type state struct {
		Flags     object.Flags
		Velocity  types.Pointf
		Data      server.MoverUpdateData
		Updatable uint32
	}
	type row struct {
		Direct, Disabled bool
		Mask             int
		Index            uint32
		Repeat           int
		States           []state
		Order            []int
	}
	var rows []row
	for _, direct := range []bool{false, true} {
		for _, disabled := range []bool{false, true} {
			for mask := 0; mask < 4; mask++ {
				for _, index := range []uint32{0, 1, 0x80000000, 0xffffffff} {
					copy(o.units, saved)
					o.s.Objs.UpdatableList = nil
					for i := range o.units {
						u := &o.units[i]
						u.ObjClass = object.ClassSimple
						u.ObjFlags = 0
						u.TypeInd = 1
						u.Extent = 1234
						u.NetCode = 9876
						u.VelVec = types.Ptf(12.5, -6.25)
						u.IsUpdatable = 0
						u.UpdatablePrev = nil
						u.UpdatableNext = nil
						u.UpdateData = unsafe.Pointer(&data[i])
						u.ObjNext = nil
						if i < 2 {
							u.ObjNext = &o.units[i+1]
						}
						data[i] = server.MoverUpdateData{Field_0: 0x7f, Field_1: 1.5, Field_2: -7, Field_3: 33, Field_4: 44, Field_5: 55, Field_6: 66, Field_8: 9999}
						if i > 0 && mask&(1<<(i-1)) != 0 {
							data[i].Field_8 = 1234
						}
					}
					if !direct {
						o.units[0].TypeInd = 2
					}
					if disabled {
						o.units[0].ObjFlags = 0x8000
					}
					o.s.Objs.List = &o.units[1]
					*wp = server.Waypoint{Index: index, PosVec: types.Ptf(77, 88)}
					for repeat := 0; repeat < 2; repeat++ {
						legacy.Nox_server_scriptMoveTo_5123C0(&o.units[0], wp)
						var states []state
						var wantOrder []int
						for i := range o.units {
							u := &o.units[i]
							selected := !disabled && (direct && i == 0 || !direct && i > 0 && mask&(1<<(i-1)) != 0)
							if selected {
								if u.ObjFlags&0x1000000 == 0 || u.VelVec != (types.Pointf{}) || data[i].Field_0 != 0 || uint32(data[i].Field_2) != index || u.IsUpdatable != 1 {
									t.Fatalf("selected mover %d", i)
								}
								wantOrder = append([]int{i}, wantOrder...)
							} else if u.VelVec != types.Ptf(12.5, -6.25) || data[i].Field_0 != 0x7f || data[i].Field_2 != -7 || u.IsUpdatable != 0 {
								t.Fatal("unselected mover changed")
							}
							if data[i].Field_1 != 1.5 || data[i].Field_3 != 33 || data[i].Field_4 != 44 || data[i].Field_5 != 55 || data[i].Field_6 != 66 || data[i].Field_7 != nil {
								t.Fatal("unrelated mover data changed")
							}
							states = append(states, state{u.ObjFlags, u.VelVec, data[i], uint32(u.IsUpdatable)})
						}
						var order []int
						var prev *server.Object
						for u := o.s.Objs.UpdatableList; u != nil; u = u.UpdatableNext {
							if len(order) >= 3 || u.UpdatablePrev != prev {
								t.Fatal("updatable list cycle/backlink")
							}
							found := -1
							for i := range o.units {
								if u == &o.units[i] {
									found = i
								}
							}
							if found < 0 {
								t.Fatal("unowned updatable")
							}
							order = append(order, found)
							prev = u
						}
						if !reflect.DeepEqual(order, wantOrder) {
							t.Fatal("updatable order", order, wantOrder)
						}
						rows = append(rows, row{direct, disabled, mask, index, repeat, states, order})
					}
				}
			}
		}
	}
	spellbookCapture(t, "script-bindings-mover", rows, "841935d04d97a13aba2bab78937b01d53b2f6df6e3477e573214485d06301d3c")
}
