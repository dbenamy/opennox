//go:build porttest

package opennox

import (
	"errors"
	"os"
	"reflect"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc/handles"
	"github.com/opennox/opennox/v1/server"
)

// These offsets are frozen from the original 386 C layout probe, not from
// the replacement aliases. Object/drawable append Go handles beyond the C
// prefix. Player keeps C field offsets/size but its existing Go alignment is 4
// (the old C declaration has alignment 1). The probe record lives in
// docs/porting/native-layout-types-original-layout.json.
func TestNativeLayoutOwners(t *testing.T) {
	cases := []struct {
		name   string
		typ    reflect.Type
		size   uintptr
		align  int
		fields map[string]uintptr
	}{
		{"Object", reflect.TypeOf(server.Object{}), 780, 4, map[string]uintptr{"TypeInd": 4, "NetCode": 36, "PosVec": 56, "Collide": 696, "Xfer": 704, "Damage": 716, "Update": 744, "UpdateData": 748, "ScriptPickup": 764, "objectHandle": 776, "serverHandle": 772}},
		{"Drawable", reflect.TypeOf(client.Drawable{}), 516, 4, map[string]uintptr{"PosVec": 12, "Shape": 44, "TypeIDVal": 108, "Buffs": 124, "Field_127": 508, "clientHandle": 512}},
		{"Player", reflect.TypeOf(server.Player{}), 4828, 4, map[string]uintptr{"PlayerUnit": 2056, "NetCodeVal": 2060, "PlayerInd": 2064, "Active": 2092, "Field2096Buf": 2096}},
		{"Window", reflect.TypeOf(gui.Window{}), 404, 4, map[string]uintptr{"WidgetData": 32, "drawData": 36, "drawFunc": 380, "parent": 396, "Field100Ptr": 400}},
		{"Team", reflect.TypeOf(server.Team{}), 80, 4, map[string]uintptr{"Lessons": 52, "Field_72": 72}},
		{"Waypoint", reflect.TypeOf(server.Waypoint{}), 516, 4, map[string]uintptr{"PosVec": 8, "Points": 92, "PointsCnt": 476, "WpNext": 484, "Field16": 512}},
		{"ScreenParticle", reflect.TypeOf(legacy.Nox_screenParticle{}), 52, 4, map[string]uintptr{"Draw_fnc": 0, "Field_44": 44, "Field_48": 48}},
		{"Thing", reflect.TypeOf(client.ObjectType{}), 128, 4, map[string]uintptr{"Name": 0, "PrettyName": 4, "Health": 124}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.typ.Size() != c.size || c.typ.Align() != c.align {
				t.Fatalf("size/alignment = %d/%d, want %d/%d", c.typ.Size(), c.typ.Align(), c.size, c.align)
			}
			for name, want := range c.fields {
				f, ok := c.typ.FieldByName(name)
				if !ok || f.Offset != want {
					t.Errorf("field %s offset = %d (present %v), want %d", name, f.Offset, ok, want)
				}
			}
		})
	}
}

func TestNativeFileHandleIdentity(t *testing.T) {
	handles.Init()
	t.Cleanup(handles.Release)
	raw, files := prefabScriptsFiles(t, prefabScriptsWords(0x12345678, 0x89abcdef), prefabScriptsWords(0xfedcba98, 0x76543210))
	first, second := legacy.NewFileHandle(files[0]), legacy.NewFileHandle(files[1])
	if first == nil || second == nil || first == second || unsafe.Pointer(first) != raw[0] || unsafe.Pointer(second) != raw[1] {
		t.Fatal("file identities were not retained or distinct")
	}
	if legacy.NewFileHandle(files[0]) != first || files[0].Handle != raw[0] || !handles.IsValid(uintptr(raw[0])) || !handles.IsValid(uintptr(raw[1])) {
		t.Fatal("registered file identity changed")
	}
	read := func(p unsafe.Pointer) uint32 { return uint32(legacy.PortTestPrefabScriptsCall(0, p, nil, nil, 0)) }
	if read(raw[0]) != 0x12345678 || read(unsafe.Pointer(legacy.NewFileHandle(files[0]))) != 0x89abcdef || read(raw[1]) != 0xfedcba98 {
		t.Fatal("handle lookup changed file identity or cursor ownership")
	}
	legacy.Nox_fs_close(nil)
	legacy.Nox_fs_close(first)
	legacy.Nox_fs_close(first)
	if _, err := files[0].File.Stat(); !errors.Is(err, os.ErrClosed) {
		t.Fatalf("closing handle did not close its file: %v", err)
	}
	if legacy.NewFileHandle(files[0]) != first {
		t.Fatal("closed file lost its retained handle identity")
	}
	if read(raw[1]) != 0x76543210 {
		t.Fatal("closing one handle affected the other file")
	}
}
