//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/object"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"testing"
)

func TestTeamUIMapWrappers(t *testing.T) {
	type record struct {
		Kind            string
		Count, Return   int
		Loaded, Visible bool
		Scanned         uint32
	}
	var records []record
	for _, kind := range []string{"ctf", "ball"} {
		for _, count := range []int{0, 1, 2, 8, 16} {
			t.Run(fmt.Sprintf("%s/%d", kind, count), func(t *testing.T) {
				o := newTeamUIOwner(t)
				defer noxflags.PortTestGameFlags(0x8000)()
				words, restore := legacy.PortTestMatchRosterGlobals()
				t.Cleanup(restore)
				s := o.c.srv.Server
				if !s.Objs.Init(32) {
					t.Fatal("object pool")
				}
				t.Cleanup(s.Objs.FreeObjects)
				t.Cleanup(s.PortTestAttackTypes(0, nil, "TeamUIFlag", "GameBallStart"))
				old := s.Objs.First()
				t.Cleanup(func() { s.Objs.SetObjects(old) })
				var flags []*server.Object
				for i := 0; i < count; i++ {
					u := s.NewObjectByTypeID("TeamUIFlag")
					if u == nil {
						t.Fatal("flag allocation")
					}
					t.Cleanup(func() { s.Objs.FreeObject(u) })
					u.ObjClass = object.Class(0x10000000)
					flags = append(flags, u)
				}
				s.Objs.SetObjects(nil)
				for i := len(flags) - 1; i >= 0; i-- {
					flags[i].ObjNext = s.Objs.First()
					s.Objs.SetObjects(flags[i])
				}
				got := teamUICall("map-"+kind, 0, 0)
				want := 0
				if count > 0 && kind == "ctf" {
					want = 1
				}
				// No ball start is on the map: the HUD opens, but ball reset returns zero.
				if got != want {
					t.Errorf("wrapper return %d want %d", got, want)
				}
				w := o.window(kind)
				loaded := w != nil
				visible := loaded && !w.Flags.IsHidden()
				if loaded != (count > 0) || visible != (count > 0) || *words["team-cap"] != uint32(count) {
					t.Errorf("loaded=%v visible=%v scanned=%d", loaded, visible, *words["team-cap"])
				}
				records = append(records, record{kind, count, got, loaded, visible, *words["team-cap"]})
			})
		}
	}
	spellbookCapture(t, "team-ui-map-wrappers", records, "9f1eccaf481dc600ed937a07e5a40ebeda5d9ad626300ad74c1284bc761db76f")
}
