//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/opennox/v1/client/gui"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestTeamUIPlayerListMissingResource(t *testing.T) {
	o := newTeamUIOwner(t)
	o.missing = true
	if got := teamUICall("players-construct", 0, 0); got != 0 || o.window("players") != nil {
		t.Fatalf("missing resource accepted: %d", got)
	}
	o.missing = false
	o.openPlayers(t)
}
func TestTeamUIPlayerListLifecycle(t *testing.T) {
	for _, destroy := range []int{0, 1} {
		t.Run(fmt.Sprint(destroy), func(t *testing.T) {
			o := newTeamUIOwner(t)
			w := o.openPlayers(t)
			legacy.PortTestTeamUI("player-add", nil, nil, 101, 0, "Player")
			teamUICall("players-destroy", 0, destroy)
			if o.window("players") != nil || len(legacy.PortTestTeamUIRows(false)) != 0 || len(legacy.PortTestTeamUIRows(true)) != 0 {
				t.Fatal("list owner survives close")
			}
			if w.Flags.Has(gui.StatusDestroyed) != (destroy != 0) {
				t.Error("window destruction flag")
			}
			for _, op := range []string{"players-refresh-if-open", "team-clear", "player-remove"} {
				if got := teamUICall(op, 101, 0); got != 0 {
					t.Errorf("closed %s returned %d", op, got)
				}
			}
			teamUICall("players-destroy", 0, destroy)
			o.openPlayers(t)
		})
	}
}
func TestTeamUIHeadlessRoster(t *testing.T) {
	type record struct {
		Mode     int
		Headless bool
		Names    []string
		Enabled  bool
		Join     bool
	}
	var records []record
	for _, mode := range []int{0, 1, 128, 129, 4096, 4097} {
		for _, headless := range []bool{false, true} {
			t.Run(fmt.Sprintf("%d/%v", mode, headless), func(t *testing.T) {
				o := newTeamUIOwner(t)
				defer noxflags.PortTestGameFlags(noxflags.GameFlag(mode))()
				oldEngine := noxflags.GetEngine()
				t.Cleanup(func() { noxflags.ResetEngine(); noxflags.SetEngine(oldEngine) })
				noxflags.UnsetEngine(noxflags.EngineNoRendering)
				if headless {
					noxflags.SetEngine(noxflags.EngineNoRendering)
				}
				w := o.openPlayers(t)
				for _, i := range []int{0, 31} {
					p := &o.players[i]
					p.Active = 1
					p.PlayerInd = byte(i)
					p.NetCodeVal = uint32(100 + i)
					p.SetName(fmt.Sprintf("Player %d", i))
				}
				teamUICall("players-refresh-if-open", 0, 0)
				names := teamUIRowNames(w.ChildByID(10501))
				want := 2
				if headless {
					want = 1
				}
				if len(names) != want || names[0] != "Player 0" {
					t.Fatalf("headless roster %v", names)
				}
				enabled := w.ChildByID(10501).Flags.IsEnabled()
				if enabled != (mode&4096 == 0) || w.ChildByID(10502).Flags.IsEnabled() != enabled {
					t.Error("quest list enable state")
				}
				w.ChildByID(10502).Func94(&gui.RawEvent{Event: 16403, Arg1: 0})
				legacy.PortTestTeamUIEvent(w, w.ChildByID(10502), 16400, 0)
				join := w.ChildByID(10507).Flags.IsEnabled()
				if join != (!headless || mode&1 == 0) {
					t.Error("headless host join control")
				}
				records = append(records, record{mode, headless, names, enabled, join})
			})
		}
	}
	spellbookCapture(t, "team-ui-headless-roster", records, "caa2e97fc696fc8fd3008724dfeea433a753d16339e6cb44d33e78bbcb11fd63")
}
