//go:build porttest

package opennox

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"unsafe"

	"github.com/opennox/libs/datapath"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc/handles"
	"github.com/opennox/opennox/v1/server"
)

func TestMapOrchestrationOuterContracts(t *testing.T) {
	for _, tc := range []struct {
		name     string
		alt      bool
		files    string
		save     int
		theme    string
		want     uint32
		attempts int
	}{
		{"main-empty", false, "", 1, "success", 1, 1},
		{"main-both", false, "both", 1, "success", 1, 1},
		{"main-current", false, "current", 1, "success", 1, 1},
		{"main-blend", false, "blend", 1, "success", 1, 1},
		{"main-remove-failure", false, "directory", 1, "success", 0, 1},
		{"main-save-failure", false, "both", 0, "success", 0, 1},
		{"alt-both", true, "both", 1, "success", 1, 1},
		{"alt-current", true, "current", 1, "success", 0, 1},
		{"alt-empty", true, "", 1, "success", 0, 1},
		{"alt-blend", true, "blend", 1, "success", 0, 1},
		{"alt-remove-failure", true, "directory", 1, "success", 0, 1},
		{"main-abort", false, "", 1, "missing", 0, 1},
		{"alt-abort", true, "", 1, "missing", 0, 1},
		{"main-retries", false, "", 1, "retry", 0, 100},
		{"alt-retries", true, "", 1, "retry", 0, 100},
	} {
		t.Run(tc.name, func(t *testing.T) {
			handles.Init()
			defer handles.Release()
			dir := t.TempDir()
			t.Chdir(dir)
			oldData := datapath.Data()
			datapath.SetData(dir)
			defer datapath.SetData(oldData)
			must := func(err error) {
				t.Helper()
				if err != nil {
					t.Fatal(err)
				}
			}
			must(os.Mkdir("mapgen", 0700))
			must(os.Mkdir("Maps", 0700))
			if tc.files == "both" || tc.files == "current" {
				must(os.WriteFile("nc.obj", []byte("old"), 0600))
			}
			if tc.files == "both" || tc.files == "blend" {
				must(os.WriteFile("blend.obj", []byte("new"), 0600))
			}
			if tc.files == "directory" {
				must(os.Mkdir("nc.obj", 0700))
				must(os.WriteFile("nc.obj/child", []byte("keep"), 0600))
			}
			decor := "DECOR ROOM base WALL_FLOOR PaintWall PaintTile END DECOR HALL hall WALL_FLOOR PaintWall PaintTile END "
			body := "ALGORITHM_DATA mapSize 200 midRoomSize 5 roomVariance 0 recursionLimit 2 seed 17 END " + decor
			if tc.theme == "retry" {
				body += "DECOR ROOM impossible MUST_OCCUR ROOM_SIZE_CONSTRAINT 99999 99999 END "
			}
			if tc.theme != "missing" {
				orchestrationTheme(t, "outer", body)
			}
			oldSwitch, oldInfo := legacy.Nox_xxx_mapSwitchLevel_4D12E0, legacy.Nox_xxx_mapGenMakeInfo_4D5DB0
			defer func() {
				legacy.Nox_xxx_mapSwitchLevel_4D12E0 = oldSwitch
				legacy.Nox_xxx_mapGenMakeInfo_4D5DB0 = oldInfo
			}()
			var events []string
			attempts := 0
			var core *server.Server
			var original *server.Object
			legacy.Nox_xxx_mapSwitchLevel_4D12E0 = func(v bool) {
				if !v {
					t.Error("map switch argument")
				}
				events = append(events, "switch")
			}
			legacy.Nox_xxx_mapGenMakeInfo_4D5DB0 = func(p unsafe.Pointer) {
				events = append(events, "metadata")
				if p != memmap.PtrOff(0x973F18, 2408) {
					t.Error("metadata target")
				}
				*(*uint32)(p) = 0x1234
			}
			op := 3
			if tc.alt {
				op = 4
			}
			spec := orchestrationCase("outer")
			spec.Actions = []legacy.PortTestPaintAction{paintAction(0, roomArg(2)), paintAction(op)}
			var returned uint32
			var restored, saved, flag bool
			svc := &legacy.PortTestMapOrchestrationServices{
				Progress: func(s string) {
					if s == "theme" {
						attempts++
						for _, v := range unsafe.Slice(memmap.PtrUint8(0x973F18, 2408), 1464) {
							if v != 0 {
								t.Error("metadata not cleared before generation")
								break
							}
						}
						if noxflags.GetGame()&0x400000 == 0 {
							t.Error("generation flag missing during attempt")
						}
						if !tc.alt && memmap.Uint32(0x5D4594, 1550924) != uint32(uintptr(unsafe.Pointer(original))) {
							t.Error("original object list not saved before attempt")
						}
					}
				},
				Before: func(s *server.Server) {
					core = s
					original = s.NewObjectByTypeInd(1)
					if original == nil {
						t.Fatal("original object")
					}
					original.ScriptIDVal = 999
					s.Objs.SetObjects(original)
					for i := range unsafe.Slice(memmap.PtrUint8(0x973F18, 2408), 1464) {
						*memmap.PtrUint8(0x973F18, 2408+uintptr(i)) = 0xa5
					}
				},
				Save: func(path string, flags int) int {
					events = append(events, "save")
					if filepath.Clean(strings.ReplaceAll(path, "\\", "/")) != filepath.Join(dir, "Maps", "$outer", "$outer.map") || flags != 1 {
						t.Errorf("map save request %q %d", path, flags)
					}
					if memmap.Uint32(0x973F18, 2408) != 0x1234 {
						t.Error("save before metadata")
					}
					for o := core.Objs.First(); o != nil; o = o.ObjNext {
						if o.ScriptIDVal != 0 {
							t.Error("save object script ID")
						}
					}
					return tc.save
				},
				After: func(s *server.Server, action int, ret uint32) {
					if action != op {
						return
					}
					returned = ret
					restored = s.Objs.First() == original
					saved = memmap.Uint32(0x5D4594, 1550924) == uint32(uintptr(unsafe.Pointer(original)))
					flag = noxflags.GetGame()&0x400000 != 0
				},
			}
			out := legacy.PortTestMapOrchestration([]legacy.PortTestPaintSpec{spec}, func(s *server.Server) (legacy.Server, func()) {
				old := noxServer
				wrapped := &Server{Server: s}
				s.ExtServer = unsafe.Pointer(wrapped)
				noxServer = wrapped
				return wrapped, func() { noxServer = old; s.ExtServer = nil }
			}, svc)
			if !out[0].Intact || !out[0].ControlOK {
				t.Fatal("outer guards/control")
			}
			if returned != tc.want || attempts != tc.attempts {
				t.Fatalf("return/attempts %d/%d want %d/%d", returned, attempts, tc.want, tc.attempts)
			}
			var wantEvents []string
			if !tc.alt {
				wantEvents = append(wantEvents, "switch")
			}
			reachesInfo := tc.theme == "success" && tc.files != "directory" && (!tc.alt || tc.files == "both")
			if reachesInfo {
				wantEvents = append(wantEvents, "metadata")
				if !tc.alt {
					wantEvents = append(wantEvents, "save")
					if tc.save != 0 {
						wantEvents = append(wantEvents, "switch")
					}
				}
			}
			if !reflect.DeepEqual(events, wantEvents) {
				t.Errorf("events %v want %v", events, wantEvents)
			}
			if flag != (tc.want == 0) {
				t.Errorf("generation flag %v", flag)
			}
			if tc.theme == "success" && tc.files != "directory" {
				moved := tc.files == "both" || (!tc.alt && tc.files == "blend")
				data, err := os.ReadFile("nc.obj")
				if moved {
					if err != nil || string(data) != "new" {
						t.Errorf("replacement file %q %v", data, err)
					}
				} else if !os.IsNotExist(err) {
					t.Errorf("current file remains: %q %v", data, err)
				}
				blend, err := os.ReadFile("blend.obj")
				if tc.alt && tc.files == "blend" {
					if err != nil || string(blend) != "new" {
						t.Error("failed remove changed blend file")
					}
				} else if !os.IsNotExist(err) {
					t.Error("blend file remains after move")
				}
			}
			if !tc.alt {
				if restored != (tc.want == 1) || saved != (tc.want == 0) {
					t.Errorf("object restoration %v saved %v", restored, saved)
				}
			} else if saved {
				t.Error("alternate start detached original objects")
			}
		})
	}
}
