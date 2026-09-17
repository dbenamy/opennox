//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/object"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/ntype"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"slices"
	"testing"
	"unsafe"
)

func TestTeamRuntimeCrownSetup(t *testing.T) {
	o := newMatchRosterOwner(t)
	t.Cleanup(o.s.PortTestRewardTypes([]string{"Crown", "TeamFixtureFlag"}, nil, true, 0, 0))
	cache := memmap.PtrUint32(0x5D4594, 527652)
	oldCache := *cache
	t.Cleanup(func() { *cache = oldCache })
	oldList, oldDeleted := o.s.Objs.First(), o.s.Objs.DeletedList
	t.Cleanup(func() { o.s.Objs.SetObjects(oldList); o.s.Objs.DeletedList = oldDeleted })
	type row struct {
		Name       string
		Result     int
		Deleted    []bool
		OwnerWords []uint32
		TeamCrown  []uint32
		Minimap    [][]uint32
		Queue      legacy.PortTestReliableReportState
	}
	var rows []row
	for mask := 0; mask < 8; mask++ {
		for _, teamMode := range []bool{false, true} {
			for _, withFlag := range []bool{false, true} {
				name := fmt.Sprintf("crowns%d/teams%t/flag%t", mask, teamMode, withFlag)
				t.Run(name, func(t *testing.T) {
					defer noxflags.PortTestGameFlags(1)()
					noxflags.UnsetGamePlay(^noxflags.GameplayFlag(0))
					if teamMode {
						noxflags.SetGamePlay(4)
					}
					o.s.Teams.Reset()
					o.s.Teams.ActiveCnt = 0
					o.s.Teams.Create(1)
					o.s.Teams.Create(2)
					for i := range o.units {
						o.units[i].TeamVal = server.ObjectTeam{}
					}
					var objects, listed []*server.Object
					for i := 0; i < 4; i++ {
						typ := "Crown"
						if i == 3 {
							typ = "TeamFixtureFlag"
						}
						u := o.s.NewObjectByTypeID(typ)
						if u == nil {
							t.Fatal("crown factory")
						}
						objects = append(objects, u)
						t.Cleanup(func() { u.TeamVal = server.ObjectTeam{}; u.DeletedNext = nil; o.s.Objs.FreeObject(u) })
						u.NetCode = uint32(500 + i)
						u.ObjFlags = object.FlagActive
						if i == 3 {
							u.ObjClass = object.Class(0x10000000)
						}
						objectXferSetWord(u.UpdateData, 4, 99)
						if i == 1 {
							legacy.Nox_xxx_createAtImpl_4191D0(1, u.TeamPtr(), 0, int(u.NetCode), 0)
						}
						if i == 2 {
							u.TeamVal.ID = 3
						}
						if i < 3 && mask&(1<<i) != 0 || i == 3 && withFlag {
							listed = append(listed, u)
						}
					}
					for i, u := range listed {
						*(*unsafe.Pointer)(unsafe.Add(u.CObj(), 444)) = nil
						if i+1 < len(listed) {
							*(*unsafe.Pointer)(unsafe.Add(u.CObj(), 444)) = listed[i+1].CObj()
						}
					}
					o.s.Objs.SetObjects(nil)
					if len(listed) > 0 {
						o.s.Objs.SetObjects(listed[0])
					}
					o.s.Objs.DeletedList = nil
					t.Cleanup(func() {
						o.s.Objs.SetObjects(nil)
						o.s.Objs.DeletedList = nil
						for _, slot := range []ntype.PlayerInd{1, 3, 7, 31} {
							for _, u := range objects {
								o.s.Players.Nox_xxx_netUnmarkMinimapObj_417300(slot, u, 0xffffffff)
							}
						}
					})
					restoreRespawns := legacy.PortTestTeamRuntimeRespawns(listed)
					defer restoreRespawns()
					for tm := o.s.Teams.First(); tm != nil; tm = o.s.Teams.Next(tm) {
						objectXferSetWord(tm.C(), 76, 0xa5a5a5a5)
					}
					*cache = 0
					o.reset()
					result := legacy.PortTestTeamRuntimeMap("crown", 0)
					wantResult := 0
					if mask != 0 {
						wantResult = 1
					}
					if result != wantResult {
						t.Fatal("crown setup result")
					}
					if *cache != uint32(o.s.Types.ByID("Crown").Ind()) || *cache == 0 {
						t.Fatal("crown type cache")
					}
					r := row{Name: name, Result: result, Queue: o.state()}
					for i, u := range objects {
						present := i < 3 && mask&(1<<i) != 0 || i == 3 && withFlag
						deleted := present && (i == 3 || !teamMode && i != 0 || teamMode && i == 0)
						if u.ObjFlags.Has(object.FlagDestroyed) != deleted {
							t.Fatal("crown/flag removal", i)
						}
						r.Deleted = append(r.Deleted, deleted)
						r.OwnerWords = append(r.OwnerWords, objectXferGetWord(u.UpdateData, 4))
						wantOwner := uint32(99)
						if present && i < 3 {
							wantOwner = 0
						}
						if r.OwnerWords[i] != wantOwner {
							t.Fatal("crown owner reset")
						}
					}
					for i := 1; i <= 2; i++ {
						p := objectXferGetWord(o.s.Teams.ByID(server.TeamID(i)).C(), 76)
						ref := uint32(0)
						if p != 0 {
							if p != uint32(uintptr(objects[1].CObj())) {
								t.Fatal("team crown pointer")
							}
							ref = 501
						}
						want := uint32(0)
						if i == 1 && teamMode && mask&2 != 0 {
							want = 501
						}
						if ref != want {
							t.Fatal("team crown association")
						}
						r.TeamCrown = append(r.TeamCrown, ref)
					}
					var wantMinimap []uint32
					if !teamMode && mask&1 != 0 {
						wantMinimap = []uint32{500}
					}
					if teamMode && mask&2 != 0 {
						wantMinimap = []uint32{501}
					}
					for _, slot := range []ntype.PlayerInd{1, 3, 7, 31} {
						var entries []uint32
						pl := o.s.Players.ByInd(slot)
						if head := pl.Field4580; head != nil {
							for p := head; ; p = p.Field8 {
								entries = append(entries, p.Field4.NetCode)
								if p.Field8 == head {
									break
								}
								if len(entries) > 4 {
									t.Fatal("minimap cycle")
								}
							}
						}
						if !slices.Equal(entries, wantMinimap) {
							t.Fatal("crown minimap", slot, entries, wantMinimap)
						}
						r.Minimap = append(r.Minimap, entries)
					}
					rows = append(rows, r)
				})
			}
		}
	}
	spellbookCapture(t, "team-runtime-crown-setup", rows, "733a50b842de312861f5662dab6630b9f6df7cc638f1e6b2cc1fb38c7147816d")
}
