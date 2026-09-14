//go:build porttest

package opennox

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc/handles"
	"github.com/opennox/opennox/v1/server"
)

func themeBase() legacy.PortTestPaintSpec {
	s := legacy.PortTestPaintSpec{Seed: 12345, Globals: map[string]legacy.PortTestMapRoomArg{"themeClock": roomValue(1700000000)}}
	for _, size := range []int{1120, 256, 256, 156, 224, 32, 160, 256} {
		s.Records = append(s.Records, roomRecord(size))
	}
	paintString(&s.Records[7], 0, "theme-fixture.dat")
	s.Globals["themeInputPath"] = roomArg(8)
	return s
}
func themeRun(cases []legacy.PortTestPaintSpec) []legacy.PortTestPaintResult {
	return legacy.PortTestMapTheme(cases, func(core *server.Server) (legacy.Server, func()) {
		old := noxServer
		wrapped := &Server{Server: core}
		core.ExtServer = unsafe.Pointer(wrapped)
		noxServer = wrapped
		return wrapped, func() { noxServer = old; core.ExtServer = nil }
	})
}
func TestMapThemeTokensProbe(t *testing.T) {
	handles.Init()
	defer handles.Release()
	t.Chdir(t.TempDir())
	if err := os.WriteFile("theme-fixture.dat", []byte(" // comment\n first\tsecond \n"), 0600); err != nil {
		t.Fatal(err)
	}
	s := themeBase()
	s.Actions = []legacy.PortTestPaintAction{paintAction(1, roomValue(legacy.PortTestThemeFile), roomValue(legacy.PortTestThemeToken)), paintAction(1, roomValue(legacy.PortTestThemeFile), roomValue(legacy.PortTestThemeToken)), paintAction(1, roomValue(legacy.PortTestThemeFile), roomValue(legacy.PortTestThemeToken))}
	out := themeRun([]legacy.PortTestPaintSpec{s})
	if !out[0].Intact || !out[0].ControlOK {
		t.Fatal("theme fixture state")
	}
	for i, want := range []uint32{1, 1, 0} {
		if out[0].Steps[i].Return != want {
			t.Fatalf("token %d result", i)
		}
	}
}
func TestMapThemeAlgorithmProbe(t *testing.T) {
	handles.Init()
	defer handles.Release()
	t.Chdir(t.TempDir())
	if err := os.WriteFile("theme-fixture.dat", []byte("midHallLength 12 END "), 0600); err != nil {
		t.Fatal(err)
	}
	s := themeBase()
	s.Actions = []legacy.PortTestPaintAction{paintAction(9, roomArg(1), roomValue(legacy.PortTestThemeFile))}
	out := themeRun([]legacy.PortTestPaintSpec{s})
	step := out[0].Steps[0]
	if !out[0].Intact || !out[0].ControlOK || step.Return != 1 {
		t.Fatal("algorithm parsing failed")
	}
	for _, r := range step.Records {
		if r.ID == step.Slots[1] && r.Words[1] != 12 {
			t.Fatalf("hall length got %d want 12", r.Words[1])
		}
	}
}

func TestMapThemeModifierCountersProbe(t *testing.T) {
	handles.Init()
	defer handles.Release()
	t.Chdir(t.TempDir())
	if err := os.WriteFile("theme-fixture.dat", []byte("QUALITY good END MATERIAL iron END PRIMARY_ENCHANTMENT fire END SECONDARY_ENCHANTMENT light END END "), 0600); err != nil {
		t.Fatal(err)
	}
	s := themeBase()
	s.Actions = []legacy.PortTestPaintAction{paintAction(13, roomArg(4), roomValue(legacy.PortTestThemeFile))}
	out := themeRun([]legacy.PortTestPaintSpec{s})
	step := out[0].Steps[0]
	if !out[0].Intact || !out[0].ControlOK || step.Return != 1 {
		t.Fatal("modifier parsing failed")
	}
	for _, r := range step.Records {
		if r.ID == step.Slots[4] {
			for i := 0; i < 4; i++ {
				if r.Words[34+i] != 1 {
					t.Fatalf("modifier slot %d count got %d want 1", i, r.Words[34+i])
				}
			}
		}
	}
}

func TestMapThemePrerequisiteBoundaries(t *testing.T) {
	handles.Init()
	defer handles.Release()
	t.Chdir(t.TempDir())
	for _, n := range []int{1, 2, 3, 10, 59} {
		input := "midHallLength " + strings.Repeat("0", n-1) + "7 END "
		if err := os.WriteFile("theme-fixture.dat", []byte(input), 0600); err != nil {
			t.Fatal(err)
		}
		s := themeBase()
		s.Actions = []legacy.PortTestPaintAction{paintAction(9, roomArg(1), roomValue(legacy.PortTestThemeFile))}
		out := themeRun([]legacy.PortTestPaintSpec{s})
		step := out[0].Steps[0]
		if !out[0].Intact || !out[0].ControlOK || step.Return != 1 {
			t.Fatalf("value length %d parsing", n)
		}
		for _, r := range step.Records {
			if r.ID == step.Slots[1] && r.Words[1] != 7 {
				t.Fatalf("value length %d result", n)
			}
		}
	}
	names := []string{"QUALITY", "MATERIAL", "PRIMARY_ENCHANTMENT", "SECONDARY_ENCHANTMENT"}
	for mask := 0; mask < 16; mask++ {
		for _, count := range []int{1, 3} {
			for _, reverse := range []bool{false, true} {
				var input strings.Builder
				for i := 0; i < 4; i++ {
					slot := i
					if reverse {
						slot = 3 - i
					}
					if mask&(1<<slot) == 0 {
						continue
					}
					input.WriteString(names[slot] + " ")
					for j := 0; j < count; j++ {
						fmt.Fprintf(&input, "slot%d_item%d ", slot, j)
					}
					input.WriteString("END ")
				}
				input.WriteString("END ")
				if err := os.WriteFile("theme-fixture.dat", []byte(input.String()), 0600); err != nil {
					t.Fatal(err)
				}
				s := themeBase()
				s.Actions = []legacy.PortTestPaintAction{paintAction(13, roomArg(4), roomValue(legacy.PortTestThemeFile))}
				out := themeRun([]legacy.PortTestPaintSpec{s})
				step := out[0].Steps[0]
				if !out[0].Intact || !out[0].ControlOK || step.Return != 1 {
					t.Fatalf("mask %d count %d modifier parsing", mask, count)
				}
				rows := map[uint32][]uint32{}
				for _, r := range step.Records {
					if r.Alive {
						rows[r.ID] = r.Words
					}
				}
				row := rows[step.Slots[4]]
				for slot := 0; slot < 4; slot++ {
					want := uint32(0)
					if mask&(1<<slot) != 0 {
						want = uint32(count)
					}
					if row[34+slot] != want {
						t.Fatalf("mask %d slot %d count", mask, slot)
					}
					if want == 0 {
						if row[30+slot] != 0 {
							t.Fatal("empty modifier slot allocated")
						}
						continue
					}
					words := rows[row[30+slot]]
					if len(words) != 15*count {
						t.Fatal("modifier array extent")
					}
					for j := 0; j < count; j++ {
						want := fmt.Sprintf("slot%d_item%d", slot, j)
						for k := 0; k < len(want)+1; k++ {
							got := byte(words[j*15+k/4] >> uint(8*(k%4)))
							expected := byte(0)
							if k < len(want) {
								expected = want[k]
							}
							if got != expected {
								t.Fatalf("mask %d slot %d item %d text", mask, slot, j)
							}
						}
					}
				}
			}
		}
	}
}
