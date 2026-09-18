//go:build porttest

package opennox

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/opennox/libs/object"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func TestPrefabScriptsObjectNames(t *testing.T) {
	nameBuf := unsafe.Slice((*byte)(memmap.PtrOff(0x5D4594, 2489164)), 256)
	saved := append([]byte(nil), nameBuf...)
	defer copy(nameBuf, saved)
	type row struct {
		Name, ID, Waypoint   string
		Flags, WaypointFlags uint32
		Callbacks            []string
		Repeated             bool
	}
	var rows []row
	specs := []struct {
		name, xfer string
		class      object.Class
		events     []int
	}{
		{"simple", "", object.ClassSimple, []int{14}},
		{"trigger", "TriggerXfer", object.ClassTrigger, []int{14, 1, 2, 0}},
		{"monster", "MonsterXfer", object.ClassMonster, []int{14, 3, 5, 4, 6, 7, 8, 9, 10, 11}},
		{"hole", "HoleXfer", object.ClassHole, []int{14, 12}},
		{"generator", "MonsterGeneratorXfer", object.ClassMonsterGenerator, []int{14, 15, 16, 18, 17}},
	}
	for _, sp := range specs {
		for _, marked := range []bool{false, true} {
			for nameMode := 0; nameMode < 3; nameMode++ {
				t.Run(fmt.Sprintf("%s/marked%t/name%d", sp.name, marked, nameMode), func(t *testing.T) {
					s := newObjectXferOwner(t)
					u := newObjectXferSimple(t, s)
					oldClass := u.ObjClass
					u.ObjClass = sp.class
					update, fu := alloc.Make([]byte{}, 4096)
					defer fu()
					collide, fc := alloc.Make([]byte{}, 88)
					defer fc()
					u.UpdateData = unsafe.Pointer(&update[0])
					u.CollideData = unsafe.Pointer(&collide[0])
					defer func() { u.ObjClass = oldClass; u.UpdateData = nil; u.CollideData = nil }()
					if sp.xfer != "" {
						s.Types.ByInd(1).Xfer = server.PortTestPrefabScriptXfer(sp.xfer)
					}
					// Handler identity is read from the real type, not this instance's Xfer field.
					u.Xfer = nil
					id := ""
					if nameMode == 2 {
						id = "Door"
					}
					if nameMode != 0 {
						u.IDPtr = legacy.PortTestPrefabRawAllocation(len(id) + 1)
						copy(unsafe.Slice((*byte)(u.IDPtr), len(id)+1), id)
					}
					for _, ev := range sp.events {
						name := fmt.Sprintf("Callback%d", ev)
						if nameMode == 0 && ev%2 == 0 {
							name = ""
						}
						p, free := alloc.CString(name)
						legacy.PortTestPrefabScriptsCall(14, u.CObj(), unsafe.Pointer(p), nil, uint32(ev))
						free()
					}
					wp := s.NewWaypoint(types.Ptf(46, 92))
					defer s.Nox_xxx_waypointDeleteAll_579DD0()
					wpName := "Route"
					if nameMode == 0 {
						wpName = ""
					}
					copy(wp.NameBuf[:], wpName)
					if marked {
						u.ObjFlags |= object.Flags(0x80000000)
						wp.Flags |= 0x80000000
					}
					beforeFlags := uint32(u.ObjFlags)
					_, restore := legacy.PortTestObjectXferEditorList()
					defer restore()
					legacy.PortTestPrefabObjectNode(u.CObj())
					xy := [2]int32{int32(nameMode-1) * 46, int32(nameMode) * 92}
					if ret := legacy.PortTestPrefabScriptsCall(11, unsafe.Pointer(&xy), nil, nil, 7); ret != 0 {
						t.Fatal("name adjustment return", ret)
					}
					wantID := id
					if marked && nameMode != 0 {
						wantID += "%7"
					}
					wantWP := wpName
					if marked && wpName != "" {
						wantWP += "%7"
					}
					if u.ID() != wantID || wp.ID() != wantWP || uint32(u.ObjFlags) != beforeFlags&0x7fffffff || wp.Flags&0x80000000 != 0 {
						t.Fatalf("identity/flags ID%q want%q wp%q want%q", u.ID(), wantID, wp.ID(), wantWP)
					}
					var callbacks []string
					for _, ev := range sp.events {
						want := fmt.Sprintf("Callback%d", ev)
						if nameMode == 0 && ev%2 == 0 {
							want = ""
						} else if marked {
							want += fmt.Sprintf("%%7%%%d%%%d", xy[0], xy[1])
						}
						got, ok := s.NoxScriptVM.Nox_script_objCallbackName_508CB0(u, ev)
						if !ok || got != want {
							t.Fatalf("event%d got%q/%t want%q", ev, got, ok, want)
						}
						callbacks = append(callbacks, got)
					}
					legacy.PortTestPrefabScriptsCall(11, nil, nil, nil, 9)
					repeated := u.ID() == wantID && wp.ID() == wantWP
					for i, ev := range sp.events {
						got, ok := s.NoxScriptVM.Nox_script_objCallbackName_508CB0(u, ev)
						repeated = repeated && ok && got == callbacks[i]
					}
					if !repeated {
						t.Fatal("second pass renamed unmarked records")
					}
					rows = append(rows, row{t.Name(), u.ID(), wp.ID(), uint32(u.ObjFlags), wp.Flags, callbacks, repeated})
				})
			}
		}
	}
	spellbookCapture(t, "prefab-scripts-object-names", rows, "1039c556f4a46d91f3d2be32adf3b571b904e2a8676e35b83d87ffbeed7117ad")
}
