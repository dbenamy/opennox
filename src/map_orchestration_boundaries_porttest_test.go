//go:build porttest

package opennox

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc/handles"
)

func TestMapOrchestrationNameBoundaries(t *testing.T) {
	handles.Init()
	defer handles.Release()
	for _, name := range []string{"", "a", strings.Repeat("x", 63)} {
		s := orchestrationCase(name)
		s.Actions = s.Actions[:2]
		r := orchestrationRun([]legacy.PortTestPaintSpec{s})[0]
		if !r.Intact || !r.ControlOK || r.Steps[0].Return != uint32(len(name)+1) {
			t.Fatalf("name length %d", len(name))
		}
		var pointer uint32
		for _, rec := range r.Steps[1].Records {
			if rec.Kind == "mapName" {
				pointer = rec.ID
				for i := 0; i <= len(name); i++ {
					want := byte(0)
					if i < len(name) {
						want = name[i]
					}
					if byte(rec.Words[i/4]>>uint(8*(i%4))) != want {
						t.Fatalf("name content %d/%d", len(name), i)
					}
				}
			}
		}
		if pointer == 0 || r.Steps[1].Return != pointer {
			t.Fatal("name getter identity")
		}
	}
}

func TestMapOrchestrationStepBoundaries(t *testing.T) {
	handles.Init()
	defer handles.Release()
	t.Chdir(t.TempDir())
	if err := os.Mkdir("mapgen", 0700); err != nil {
		t.Fatal(err)
	}
	var cases []legacy.PortTestPaintSpec
	for kind := 0; kind < 3; kind++ {
		for seed := 0; seed < 8; seed++ {
			name := fmt.Sprintf("boundary-%d-%d", kind, seed)
			body := fmt.Sprintf("ALGORITHM_DATA mapSize 200 midRoomSize 5 roomVariance 0 recursionLimit 2 seed %d END ", seed) + "DECOR ROOM base WALL_FLOOR PaintWall PaintTile END DECOR HALL hall WALL_FLOOR PaintWall PaintTile END "
			if kind == 1 {
				body += "DECOR BACKDROP background WALL_FLOOR PaintWall PaintTile END "
			}
			if kind == 2 {
				body += "DECOR ROOM impossible MUST_OCCUR ROOM_SIZE_CONSTRAINT 99999 99999 END "
			}
			orchestrationTheme(t, name, body)
			cases = append(cases, orchestrationCase(name))
		}
	}
	out := orchestrationCapture(t, "step-boundaries", cases)
	for i, r := range out {
		want := uint32(1)
		if i >= 16 {
			want = 2
		}
		if r.Steps[2].Return != want {
			t.Fatalf("boundary %d return %d want %d", i, r.Steps[2].Return, want)
		}
	}
	for seed := 0; seed < 8; seed++ {
		if len(out[8+seed].Steps[2].Cells) <= len(out[seed].Steps[2].Cells) {
			t.Fatalf("background did not add floor cells, seed %d", seed)
		}
	}
}
