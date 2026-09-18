//go:build porttest

package opennox

import (
	"reflect"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

func TestPrefabScriptsPendingGroupReferences(t *testing.T) {
	s := newObjectXferOwner(t)
	s.MapGroups.Init()
	defer func() { prefabGroupReset(s); s.MapGroups.Free() }()
	words, restore := legacy.PortTestPrefabRuntimeGlobals()
	defer restore()
	*words["instance"] = 7
	var objs []*server.Object
	for i, id := range []int{100, 900, 900} {
		u := newObjectXferSimple(t, s)
		u.ScriptIDVal = id
		u.Extent = uint32(41 + i)
		if i > 0 {
			objs[i-1].ObjNext = u
			u.ObjPrev = objs[i-1]
		}
		objs = append(objs, u)
	}
	s.Objs.Pending = objs[0]
	defer func() { s.Objs.Pending = nil }()
	s.MapGroups.Sub504600("Doors", 11, uint8(server.MapGroupObjects))
	for _, id := range []uint32{100, 999, 900} {
		if s.MapGroups.Sub5046A0([]uint32{id}, 11) != 1 {
			t.Fatal("group item allocation")
		}
	}
	if ret := noxServer.Sub504720(0, 0); ret != 1 {
		t.Fatal("group mapping return", ret)
	}
	got := prefabGroupSnapshot(s.MapGroups.GetFirstMapGroup())
	want := prefabGroupRecord{ID: 11, Kind: byte(server.MapGroupObjects), Name: "Doors%7", Items: [][2]uint32{{42, 0}, {999, 0}, {41, 0}}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("group mapping got%+v want%+v", got, want)
	}
	spellbookCapture(t, "prefab-scripts-pending-groups", []prefabGroupRecord{got}, "2ccbfeeed37252eff482d6cb951a0645bbd7a861f001ccb46431d9a637481aa5")
}
