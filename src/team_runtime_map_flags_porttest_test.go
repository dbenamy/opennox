//go:build porttest

package opennox

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/opennox/libs/object"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func newTeamRuntimeFlagOwner(t *testing.T) (*matchRosterOwner, []*server.Object) {
	o := newMatchRosterOwner(t)
	t.Cleanup(o.s.PortTestAttackTypes(0, nil, "TeamFixtureFlag"))
	colors := unsafe.Slice(memmap.PtrUint32(0x587000, 205224), 8)
	oldColors := append([]uint32(nil), colors...)
	clear(colors)
	t.Cleanup(func() { copy(colors, oldColors) })
	oldList := o.s.Objs.First()
	t.Cleanup(func() { o.s.Objs.SetObjects(oldList) })
	var flags []*server.Object
	for i, name := range []string{"RedFlag", "BlueFlag", "GoldFlag"} {
		text, free := alloc.CString(name)
		t.Cleanup(free)
		colors[2*i], colors[2*i+1] = uint32(uintptr(unsafe.Pointer(text))), uint32(i+1)
		u := o.s.NewObjectByTypeID("TeamFixtureFlag")
		if u == nil {
			t.Fatal("flag factory")
		}
		t.Cleanup(func() { o.s.Objs.FreeObject(u) })
		u.ObjClass = object.Class(0x10000000)
		u.TeamVal.ID = server.TeamID(i + 1)
		header := o.record(t, 4)
		*(*unsafe.Pointer)(header) = unsafe.Pointer(text)
		*(*unsafe.Pointer)(unsafe.Add(u.InitData, 4)) = header
		flags = append(flags, u)
	}
	return o, flags
}

func TestTeamRuntimeMapFlags(t *testing.T) {
	o, flags := newTeamRuntimeFlagOwner(t)
	type row struct {
		Count, TeamCount, Found int
		Colors                  []byte
		FlagRefs                []int
		On                      []bool
		Queue                   legacy.PortTestReliableReportState
	}
	var rows []row
	for count := 0; count <= 3; count++ {
		for teamCount := 0; teamCount <= 3; teamCount++ {
			t.Run(fmt.Sprintf("flags%d/teams%d", count, teamCount), func(t *testing.T) {
				defer noxflags.PortTestGameFlags(0x8001)()
				o.s.Teams.Reset()
				o.s.Teams.ActiveCnt = 0
				for i := 0; i < teamCount; i++ {
					o.s.Teams.Create(server.TeamID(i + 1))
				}
				for i, u := range flags {
					*(*unsafe.Pointer)(unsafe.Add(u.CObj(), 444)) = nil
					if i+1 < count {
						*(*unsafe.Pointer)(unsafe.Add(u.CObj(), 444)) = flags[i+1].CObj()
					}
				}
				o.s.Objs.SetObjects(nil)
				if count > 0 {
					o.s.Objs.SetObjects(flags[0])
				}
				*o.roster["team-cap"] = 99
				found := legacy.PortTestTeamRuntimeMap("scan", 0)
				wantFound := 0
				if count > 0 {
					wantFound = 1
				}
				if found != wantFound || *o.roster["team-cap"] != uint32(count) {
					t.Fatal("map flag count")
				}
				if legacy.PortTestTeamRuntimeMap("assign", 0) != 0 {
					t.Fatal("map assignment return")
				}
				r := row{Count: count, TeamCount: teamCount, Found: found}
				for i := 0; i < teamCount; i++ {
					tm := o.s.Teams.ByID(server.TeamID(i + 1))
					ref := -1
					if i < count {
						if tm.Field_72 != flags[i].CObj() || tm.ColorInd != server.TeamColor(i+1) {
							t.Fatal("flag association or nonzero color")
						}
						ref = i
					} else if tm.Field_72 != nil {
						t.Fatal("team without map flag")
					}
					r.Colors = append(r.Colors, byte(tm.ColorInd))
					r.FlagRefs = append(r.FlagRefs, ref)
				}
				for _, on := range []int{1, 0, 1} {
					for _, u := range flags {
						u.ObjFlags = object.FlagActive
						if on != 0 {
							u.ObjFlags &^= object.FlagEnabled
						} else {
							u.ObjFlags |= object.FlagEnabled
						}
					}
					o.reset()
					if legacy.PortTestTeamRuntimeMap("toggle", on) != 0 {
						t.Fatal("toggle return")
					}
					for i, u := range flags {
						want := on != 0
						if i >= count || i >= teamCount {
							want = !want
						}
						got := u.ObjFlags.Has(object.FlagEnabled)
						if got != want {
							t.Fatal("team flag enabled state")
						}
						r.On = append(r.On, got)
					}
				}
				r.Queue = o.state()
				rows = append(rows, r)
			})
		}
	}
	spellbookCapture(t, "team-runtime-map-flags", rows, "bbbd493f3d389dc65bde4180f98dcb2d9fe4191476d23f63a3f9e12d0971d0a5")
}
