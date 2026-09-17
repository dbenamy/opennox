//go:build porttest

package opennox

import (
	"fmt"
	"math"
	"testing"

	"github.com/opennox/libs/types"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

func TestTeamRuntimeNearestFlag(t *testing.T) {
	o := newMatchRosterOwner(t)
	defer noxflags.PortTestGameFlags(1)()
	oldHost := legacy.ClientPlayerNetCode()
	legacy.ClientSetPlayerNetCode(0)
	t.Cleanup(func() { legacy.ClientSetPlayerNetCode(oldHost) })
	var flags [2]*server.Object
	for i := range flags {
		u := o.s.NewObjectByTypeID("PortCreatureMonster")
		if u == nil {
			t.Fatal("flag position owner")
		}
		t.Cleanup(func() { o.s.Objs.FreeObject(u) })
		u.PosVec = types.Ptf(float32(10*i), 0)
		flags[i] = u
	}
	type row struct {
		Name  string
		IDs   []byte
		Queue legacy.PortTestReliableReportState
	}
	var rows []row
	xs := []float32{-100, 0, math.Nextafter32(5, 0), 5, math.Nextafter32(5, 10), 10, 100000, float32(math.Inf(1)), float32(math.NaN())}
	for present := 0; present < 4; present++ {
		for existing := 0; existing < 3; existing++ {
			for _, x := range xs {
				name := fmt.Sprintf("flags%d/existing%d/x%x", present, existing, math.Float32bits(x))
				t.Run(name, func(t *testing.T) {
					o.s.Teams.Reset()
					o.s.Teams.ActiveCnt = 0
					for i := range flags {
						tm := o.s.Teams.Create(server.TeamID(i + 1))
						if present&(1<<i) != 0 {
							tm.Field_72 = flags[i].CObj()
						}
					}
					for i := range o.units {
						u := &o.units[i]
						u.TeamVal = server.ObjectTeam{}
						u.UpdateDataPlayer().Player.NetCodeVal = u.NetCode
						u.UpdateDataPlayer().Player.Field3680 = 0
						u.PosVec = types.Ptf(x, 0)
						if existing != 0 {
							legacy.Nox_xxx_createAtImpl_4191D0(server.TeamID(existing), u.TeamPtr(), 0, int(u.NetCode), 0)
						}
					}
					o.reset()
					if got := legacy.PortTestTeamRuntimeMap("nearest", 0); got != 0 {
						t.Fatal("nearest assignment return")
					}
					want := existing
					// For this geometry the bisector is x=5. Strict comparison gives
					// the first team the exact tie; the initial distance cutoff is 1e9.
					if !math.IsNaN(float64(x)) && !math.IsInf(float64(x), 0) && x != 100000 {
						switch present {
						case 1:
							want = 1
						case 2:
							want = 2
						case 3:
							want = 1
							if x > 5 {
								want = 2
							}
						}
					}
					r := row{Name: name, Queue: o.state()}
					for i := range o.units {
						id := byte(o.units[i].TeamVal.ID)
						if int(id) != want {
							t.Fatalf("nearest team %d want %d", id, want)
						}
						r.IDs = append(r.IDs, id)
					}
					rows = append(rows, r)
				})
			}
		}
	}
	spellbookCapture(t, "team-runtime-nearest", rows, "b5b1a98d8f8ac14f9ca551c93dd1a0d851878645dc9461874078f7a95c4bf2f3")
}
