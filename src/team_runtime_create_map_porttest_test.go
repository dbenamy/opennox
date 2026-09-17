//go:build porttest

package opennox

import (
	"fmt"
	"testing"
	"unicode/utf16"
	"unsafe"

	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

func TestTeamRuntimeCreateMapTeams(t *testing.T) {
	o, flags := newTeamRuntimeFlagOwner(t)
	defer noxflags.PortTestGameFlags(1)()
	type team struct {
		ID, Color byte
		Name      string
		Flag      int
		Named     uint32
	}
	type row struct {
		Count int
		Teams []team
		Queue legacy.PortTestReliableReportState
	}
	var rows []row
	for count := 0; count <= 3; count++ {
		t.Run(fmt.Sprint(count), func(t *testing.T) {
			o.s.Teams.Reset()
			o.s.Teams.ActiveCnt = 0
			noxflags.UnsetGamePlay(^noxflags.GameplayFlag(0))
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
			o.reset()
			if legacy.PortTestTeamRuntimeMap("create", 0) != 0 || o.s.Teams.Count() != 2 || !noxflags.HasGamePlay(4) {
				t.Fatal("map team creation")
			}
			r := row{Count: count, Queue: o.state()}
			for i := 0; i < 2; i++ {
				tm := o.s.Teams.ByID(server.TeamID(i + 1))
				ref := -1
				if i < count {
					title := utf16.Encode([]rune(o.s.Teams.TeamTitle(server.TeamColor(i + 1))))
					if len(title) > 20 {
						title = title[:20]
					}
					if tm.Field_72 != flags[i].CObj() || tm.ColorInd != server.TeamColor(i+1) || tm.Name() != string(utf16.Decode(title)) || objectXferGetWord(tm.C(), 68) != 1 {
						t.Fatalf("map flag team title=%q color=%d named=%d linked=%t", tm.Name(), tm.ColorInd, objectXferGetWord(tm.C(), 68), tm.Field_72 == flags[i].CObj())
					}
					ref = i
				} else if tm.Field_72 != nil || objectXferGetWord(tm.C(), 68) != 0 {
					t.Fatal("map team without flag")
				}
				r.Teams = append(r.Teams, team{byte(tm.ID()), byte(tm.ColorInd), tm.Name(), ref, objectXferGetWord(tm.C(), 68)})
			}
			want := count
			if want > 2 {
				want = 2
			}
			if len(r.Queue.Nodes) != want {
				t.Fatal("map team announcements")
			}
			rows = append(rows, r)
		})
	}
	spellbookCapture(t, "team-runtime-create-map", rows, "3f3d83fd1694d2cb1393622cd87e2d2c88ea7b71f8aedb97ddffe34a64f31ce0")
}
