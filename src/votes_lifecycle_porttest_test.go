//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestPortVotesGUILifecycle(t *testing.T) {
	o, words, w := newVoteGUIOwner(t)
	type row struct {
		Topic, Count, Previous, Choice, PreviousChoice uint32
		Window, Players, Topics, Captured, Stacked     bool
		Return                                         int
	}
	var rows []row
	snapshot := func(rv int) {
		rows = append(rows, row{*words["topic"], *words["count"], *words["previousCount"], *words["choice"], *words["previousChoice"], *words["window"] != 0, *words["players"] != 0, *words["topics"] != 0, o.c.GUI.Captured() != nil, o.c.GUI.StackHead() != nil, rv})
	}
	for pass := 0; pass < 2; pass++ {
		if pass != 0 {
			if legacy.PortTestVoteGUI("init", 0) != 1 {
				t.Fatal("reinitialization")
			}
			w = legacy.PortTestVoteGUIWindow()
		}
		*words["topic"] = 3
		*words["count"] = 7
		*words["previousCount"] = 5
		*words["choice"] = 1
		*words["previousChoice"] = 1
		w.ShowModal()
		w.StackPush()
		w.Capture(true)
		rv := legacy.PortTestVoteGUI("reset-choice", 0)
		snapshot(rv)
		if rv != 0 || *words["choice"] != 0 || *words["previousChoice"] != 0 || *words["count"] != 7 || *words["previousCount"] != 5 || o.c.GUI.Captured() != w || o.c.GUI.StackHead() != w {
			t.Fatalf("reset-choice changed unrelated state: %+v", rows[len(rows)-1])
		}
		rv = legacy.PortTestVoteGUI("close", 0)
		snapshot(rv)
		if rv != 0 || *words["topic"] != 3 || *words["window"] != 0 || *words["players"] != 0 || *words["topics"] != 0 || *words["count"] != 0 || *words["previousCount"] != 0 || o.c.GUI.Captured() != nil || o.c.GUI.StackHead() != nil {
			t.Fatal("window disposal")
		}
		o.c.GUI.FreeDestroyed()
		rv = legacy.PortTestVoteGUI("close", 0)
		snapshot(rv)
		if rv != 0 {
			t.Fatal("repeat disposal")
		}
	}
	spellbookCapture(t, "votes-gui-lifecycle", rows, "25fb1756584bd1ecf2a6ca296409873270af73387827cd64009bfb2963c94441")
}
