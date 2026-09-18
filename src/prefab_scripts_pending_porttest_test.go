//go:build porttest

package opennox

import (
	"fmt"
	"math"
	"testing"
	"unsafe"

	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func TestPrefabScriptsPendingReferences(t *testing.T) {
	grid := (*[2]uint32)(memmap.PtrOff(0x5D4594, 739980))
	oldGrid := *grid
	defer func() { *grid = oldGrid }()
	*grid = [2]uint32{17, 23}
	type rec struct {
		Name   string
		Extent uint32
		Old    int
		Data   []uint32
	}
	type row struct {
		Name    string
		Return  uint64
		Records []rec
	}
	var rows []row
	for _, base := range []uint32{0, 1, 0x7ffffffe, 0xfffffffe} {
		for _, found := range []bool{false, true} {
			for _, wide := range []bool{false, true} {
				t.Run(fmt.Sprintf("base%x/found%t/wide%t", base, found, wide), func(t *testing.T) {
					s := newObjectXferOwner(t)
					names := []string{"Elevator", "ElevatorShaft", "Transporter", "Hole", "Exit", "Mover", "Glyph"}
					var objs []*server.Object
					for _, name := range names {
						objs = append(objs, newObjectXferTyped(t, s, name))
					}
					objs = append(objs, newObjectXferSimple(t, s), newObjectXferSimple(t, s))
					names = append(names, "Target", "DuplicateTarget")
					for i, u := range objs {
						u.Extent = uint32((i + 1) * 100)
						if i >= 7 {
							u.Extent = 900
						}
						if i+1 < len(objs) {
							u.ObjNext = objs[i+1]
						}
						if i > 0 {
							u.ObjPrev = objs[i-1]
						}
					}
					s.Objs.Pending = objs[0]
					defer func() { s.Objs.Pending = nil }()
					wp := s.WPs.Nox_xxx_waypointNewNotMap_579970(1234, types.Ptf(46, 92))
					wp.Field1 = 17
					defer func() { s.WPs.List = s.WPs.Pending; s.WPs.Pending = nil; s.Nox_xxx_waypointDeleteAll_579DD0() }()
					word := func(p unsafe.Pointer, off int) *uint32 { return (*uint32)(unsafe.Add(p, off)) }
					target := uint32(999)
					if found {
						target = 900
					}
					*word(objs[0].UpdateData, 8) = target
					*word(objs[1].UpdateData, 8) = target
					*word(objs[2].UpdateData, 16) = target
					*word(objs[5].UpdateData, 8) = 999
					if found {
						*word(objs[5].UpdateData, 8) = 17
					}
					*word(objs[5].UpdateData, 32) = target
					*word(objs[3].CollideData, 8) = 7
					*word(objs[3].CollideData, 12) = 0xfffffff9
					for _, p := range []unsafe.Pointer{unsafe.Add(objs[4].CollideData, 80), unsafe.Add(objs[6].InitData, 28)} {
						*word(p, 0) = math.Float32bits(16777216)
						*word(p, 4) = math.Float32bits(-0.25)
					}
					bounds, free := alloc.New([4]uint32{})
					defer free()
					bounds[0], bounds[1] = 400, 500
					if wide {
						bounds[0], bounds[1] = 0x7fffffff, 0x80000000
					}
					dx, dy := bounds[0]-23*grid[0], bounds[1]-23*grid[1]
					ret := legacy.PortTestPrefabScriptsCall(15, unsafe.Pointer(bounds), nil, nil, base)
					if ret != uint64(base+uint32(len(objs))) {
						t.Fatal("next extent", ret)
					}
					for i, u := range objs {
						old := int((i + 1) * 100)
						if i >= 7 {
							old = 900
						}
						if u.Extent != base+uint32(i) || u.ScriptIDVal != old {
							t.Fatal("pending indices", i, u.Extent, u.ScriptIDVal)
						}
					}
					targetID := uint32(0)
					var targetPtr unsafe.Pointer
					if found {
						targetID = base + 7
						targetPtr = objs[7].CObj()
					}
					for i, off := range []int{4, 4, 12} {
						p := objs[i].UpdateData
						if *(*unsafe.Pointer)(unsafe.Add(p, off)) != targetPtr || *word(p, off+4) != targetID {
							t.Fatal("linked target/first duplicate", i)
						}
					}
					waypointID := uint32(0)
					if found {
						waypointID = 1234
					}
					if *word(objs[5].UpdateData, 8) != waypointID || *word(objs[5].UpdateData, 32) != targetID {
						t.Fatal("mover references")
					}
					if *word(objs[3].CollideData, 8) != 7+dx || *word(objs[3].CollideData, 12) != uint32(0xfffffff9)+dy {
						t.Fatal("hole coordinate adjustment")
					}
					wantX := math.Float32bits(float32(float64(16777216) + float64(int32(dx))))
					wantY := math.Float32bits(float32(-0.25 + float64(int32(dy))))
					for _, p := range []unsafe.Pointer{unsafe.Add(objs[4].CollideData, 80), unsafe.Add(objs[6].InitData, 28)} {
						if *word(p, 0) != wantX || *word(p, 4) != wantY {
							t.Fatal("floating coordinate adjustment")
						}
					}
					r := row{Name: t.Name(), Return: ret}
					for i, u := range objs {
						var data []uint32
						switch i {
						case 0, 1, 2:
							p := u.UpdateData
							off := 4
							if i == 2 {
								off = 12
							}
							normalized := uint32(0)
							if *(*unsafe.Pointer)(unsafe.Add(p, off)) != nil {
								normalized = 8
							}
							data = []uint32{normalized, *word(p, off+4)}
						case 3:
							data = []uint32{*word(u.CollideData, 8), *word(u.CollideData, 12)}
						case 4:
							data = []uint32{*word(u.CollideData, 80), *word(u.CollideData, 84)}
						case 5:
							data = []uint32{*word(u.UpdateData, 8), *word(u.UpdateData, 32)}
						case 6:
							data = []uint32{*word(u.InitData, 28), *word(u.InitData, 32)}
						}
						r.Records = append(r.Records, rec{names[i], u.Extent, u.ScriptIDVal, data})
					}
					rows = append(rows, r)
				})
			}
		}
	}
	spellbookCapture(t, "prefab-scripts-pending", rows, "9e5398ec1b43f8f876605a9d257ba3885e3c684cd242ea348aea0520bd1543e8")
}
