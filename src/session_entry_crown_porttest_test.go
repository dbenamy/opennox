//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/object"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

func TestSessionEntryCrownCleanup(t *testing.T) {
	o := newMatchRosterOwner(t)
	t.Cleanup(o.s.PortTestRewardTypes([]string{"Crown", "OtherItem"}, nil, true, 0, 0))
	serverConfigOwnBytes(t, 0x5D4594, 1523076, 4)
	cache := memmap.PtrUint32(0x5D4594, 1523076)
	oldList, oldDeleted := o.s.Objs.First(), o.s.Objs.DeletedList
	t.Cleanup(func() { o.s.Objs.SetObjects(oldList); o.s.Objs.DeletedList = oldDeleted })
	type row struct {
		Mask     int
		Warm     bool
		Deleted  [3]bool
		Frames   [3]uint32
		Cache    uint32
		Respawns []int
	}
	var rows []row
	for mask := 0; mask < 8; mask++ {
		for _, warm := range []bool{false, true} {
			t.Run(fmt.Sprintf("mask%d/warm%v", mask, warm), func(t *testing.T) {
				defer noxflags.PortTestGameFlags(1)()
				var units, listed []*server.Object
				for i, name := range []string{"Crown", "OtherItem", "Crown"} {
					u := o.s.NewObjectByTypeID(name)
					if u == nil {
						t.Fatal("item allocation")
					}
					units = append(units, u)
					t.Cleanup(func() { u.DeletedNext = nil; o.s.Objs.FreeObject(u) })
					u.ObjFlags = object.FlagActive
					u.ObjClass = 0
					u.TeamVal = server.ObjectTeam{}
					if mask&(1<<i) != 0 {
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
				t.Cleanup(func() { o.s.Objs.SetObjects(nil); o.s.Objs.DeletedList = nil })
				restore := legacy.PortTestTeamRuntimeRespawns(listed)
				defer restore()
				*cache = 0
				if warm {
					*cache = uint32(o.s.Types.ByID("Crown").Ind())
				}
				o.reset()
				legacy.PortTestSessionEntryScalar("crown-clear", 0)
				r := row{Mask: mask, Warm: warm, Cache: *cache}
				if *cache == 0 || *cache != uint32(o.s.Types.ByID("Crown").Ind()) {
					t.Fatal("crown cache")
				}
				for i, u := range units {
					r.Deleted[i] = u.Flags().Has(object.FlagDestroyed)
					r.Frames[i] = u.DeletedAt
					want := i != 1 && mask&(1<<i) != 0
					if r.Deleted[i] != want || want && u.DeletedAt != o.s.Frame() {
						t.Fatalf("item%d deletion mismatch", i)
					}
				}
				for _, u := range legacy.PortTestSessionEntryRespawns() {
					id := -1
					for i, p := range units {
						if p == u {
							id = i
						}
					}
					if id != 1 {
						t.Fatalf("unexpected respawn item%d", id)
					}
					r.Respawns = append(r.Respawns, id)
				}
				wantRespawns := 0
				if mask&2 != 0 {
					wantRespawns = 1
				}
				if len(r.Respawns) != wantRespawns {
					t.Fatal("respawn count")
				}
				rows = append(rows, r)
			})
		}
	}
	spellbookCapture(t, "session-entry-crown-clear", rows, "50efcabe4578c8031ca445e427e5c3c6804123cc6b634160b2aabf1a5c14036a")
}
