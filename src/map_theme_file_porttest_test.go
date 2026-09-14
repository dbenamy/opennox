//go:build porttest

package opennox

import (
	"fmt"
	"os"
	"strings"
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

// Exercise the outer parser through real keyed theme files, including missing
// files, partial sections and settings validation after otherwise valid parsing.
func TestMapThemeFullFiles(t *testing.T) {
	handles.Init()
	defer handles.Release()
	t.Chdir(t.TempDir())
	if err := os.Mkdir("mapgen", 0700); err != nil {
		t.Fatal(err)
	}
	valid := "DECOR ROOM base END DECOR HALL hall END "
	texts := []string{"", "unknown ", "DECOR ROOM base END ", "DECOR HALL hall END ", valid, valid + "DECOR BACKDROP back END ", "// comment\n" + strings.ToLower(valid)}
	sections := []string{
		"ALGORITHM_DATA midHallLength 12 END ",
		"AMBIENT_LIGHT 1 2 255 ",
		"EXIT OBJECT exit NORTH LINKDATA link END ",
		"SPELL_SET FIREBALL MAGIC_MISSILE END ",
		"WEAPON_SET WEAPON TEMPLATE QUALITY good END END WEAPON tier sword QUALITY good better END END END ",
		"ARMOR_SET ARMOR tier armor MATERIAL iron END END END ",
		"PREFAB MUST_OCCUR AREAMAP area FOREACH PaintObject CONTAINS * OBJECT item END END ",
		"DECOR TEMPLATE template WALL_FLOOR wall floor EDGING edge floor DECOR_SET 0 1 OBJECT item DENSITY 0.1 * * CONTAINS * SPELL FIREBALL END END END DECOR ROOM copy COPY template END ",
	}
	for _, section := range sections {
		texts = append(texts, section+valid, valid+section)
		words := strings.Fields(section)
		for _, n := range []int{1, 2, len(words) - 1} {
			texts = append(texts, valid+strings.Join(words[:n], " ")+" ")
		}
	}
	texts = append(texts, strings.Join(sections, "")+valid, "IF 100% "+valid+" ELSE unknown ENDIF ", "IF 0% unknown ELSE "+valid+" ENDIF ", valid+"DECOR ROOM forbidden OCCUR_CONSTRAINT START FREQUENCY HARDLY_EVER END ")
	var specs []legacy.PortTestPaintSpec
	for i, text := range texts {
		name := fmt.Sprintf("full-%04d", i)
		f, err := binfile.BinfileOpen("mapgen/"+name+".thm", binfile.WriteOnly)
		if err != nil {
			t.Fatal(err)
		}
		if err = f.SetKey(1); err != nil {
			_ = f.Close()
			t.Fatal(err)
		}
		if _, err = f.Write([]byte(text)); err != nil {
			_ = f.Close()
			t.Fatal(err)
		}
		if err = f.Close(); err != nil {
			t.Fatal(err)
		}
		s := themeBase()
		s.Globals["themeInputPath"] = roomValue(0)
		paintString(&s.Records[1], 0, name)
		s.Records[0].Words[1116] = 0x12345678
		s.Actions = []legacy.PortTestPaintAction{paintAction(0, roomArg(1), roomArg(2)), paintAction(32, roomArg(1))}
		specs = append(specs, s)
	}
	s := themeBase()
	s.Globals["themeInputPath"] = roomValue(0)
	paintString(&s.Records[1], 0, "missing")
	s.Actions = []legacy.PortTestPaintAction{paintAction(0, roomArg(1), roomArg(2))}
	specs = append(specs, s)
	out := themeCheckCapture(t, "full-files", themeRun(specs))
	for i, r := range out[:len(texts)] {
		for _, row := range r.Steps[0].Records {
			if row.ID == r.Steps[0].Slots[1] && row.Words[279] != 0x12345678 {
				t.Fatalf("case %d changed bytes outside config", i)
			}
		}
	}
}
