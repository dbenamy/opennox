//go:build porttest

package opennox

import (
	"os"
	"testing"

	"github.com/opennox/opennox/v1/internal/binfile"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc/handles"
)

func TestMapThemeFullFileProbe(t *testing.T) {
	handles.Init()
	defer handles.Release()
	t.Chdir(t.TempDir())
	if err := os.Mkdir("mapgen", 0700); err != nil {
		t.Fatal(err)
	}
	f, err := binfile.BinfileOpen("mapgen/theme-fixture.thm", binfile.WriteOnly)
	if err != nil {
		t.Fatal(err)
	}
	if err := f.SetKey(1); err != nil {
		_ = f.Close()
		t.Fatal(err)
	}
	if _, err := f.Write([]byte("DECOR ROOM base END DECOR HALL hall END ")); err != nil {
		_ = f.Close()
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
	s := themeBase()
	s.Globals["themeInputPath"] = roomValue(0)
	paintString(&s.Records[1], 0, "theme-fixture")
	s.Actions = []legacy.PortTestPaintAction{paintAction(0, roomArg(1), roomArg(2)), paintAction(32, roomArg(1))}
	out := themeRun([]legacy.PortTestPaintSpec{s})
	step := out[0].Steps[0]
	if !out[0].Intact || !out[0].ControlOK || step.Return != 1 {
		t.Fatal("full theme parsing")
	}
	for _, r := range step.Records {
		if r.ID == step.Slots[1] {
			if r.Words[1] != 5 || r.Words[19] != 1700000000 || r.Words[23] != 1 || r.Words[31] != 1 {
				t.Fatal("theme defaults or decoration counts")
			}
			for _, start := range []int{24, 32} {
				for i := 0; i < 6; i++ {
					if r.Words[start+i] != 1000 {
						t.Fatal("theme frequency totals")
					}
				}
			}
		}
	}
	for _, r := range out[0].Steps[1].Records {
		if r.Kind == "themeAllocation" && r.Alive {
			t.Fatal("theme cleanup left a decoration alive")
		}
	}
}
