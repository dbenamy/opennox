//go:build porttest

package opennox

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/opennox/opennox/v1/internal/binfile"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc/handles"
)

// Corrected C step captures; outer start/retry/save behavior has separate contracts.
var orchestrationCExpected = map[string]string{
	"step-boundaries": "635b72c0db8747b78fa33790d65a35c0c45e5886ce3f3ca7632a29fbbfd5de32",
	"normal-maps":     "c6fc4af3bd0c2935bf7c5c3607302ca0e066950237e2d455f5170a42967ef6ab",
	"ring-maps":       "d8f1e06ff469e2c6b6dce2f9583c375a0c62f936831875547d06c3803acfb93d",
	"theme-failures":  "e0bdd5828f0ebcc7a584699b059b6826c8b19281160eac5351eb9600c4da89c6",
}

func orchestrationTheme(t *testing.T, name, body string) {
	t.Helper()
	f, err := binfile.BinfileOpen("mapgen/"+name+".thm", binfile.WriteOnly)
	if err != nil {
		t.Fatal(err)
	}
	if err = f.SetKey(1); err != nil {
		_ = f.Close()
		t.Fatal(err)
	}
	if _, err = f.Write([]byte(body)); err != nil {
		_ = f.Close()
		t.Fatal(err)
	}
	if err = f.Close(); err != nil {
		t.Fatal(err)
	}
}
func orchestrationCase(name string) legacy.PortTestPaintSpec {
	s := growthBase()
	s.Globals["growthInitGrid"] = roomValue(0)
	paintString(&s.Records[1], 0, name)
	s.Actions = []legacy.PortTestPaintAction{paintAction(0, roomArg(2)), paintAction(1), paintAction(2)}
	return s
}
func orchestrationCapture(t *testing.T, label string, cases []legacy.PortTestPaintSpec) []legacy.PortTestPaintResult {
	t.Helper()
	out := orchestrationRun(cases)
	for i, r := range out {
		if !r.Intact || !r.ControlOK {
			t.Fatalf("%s case %d guards/control", label, i)
		}
	}
	data, err := json.Marshal(out)
	if err != nil {
		t.Fatal(err)
	}
	if prefix := os.Getenv("OPENNOX_MAP_ORCHESTRATION_CAPTURE"); prefix != "" {
		if err = os.WriteFile(prefix+"-"+label+".json", data, 0600); err != nil {
			t.Fatal(err)
		}
	}
	hash := fmt.Sprintf("%x", sha256.Sum256(data))
	if want := orchestrationCExpected[label]; want != hash {
		t.Fatalf("%s hash %s want %s", label, hash, want)
	}
	t.Logf("%s: %d cases %s", label, len(cases), hash)
	return out
}
func TestMapOrchestrationGeneratedMaps(t *testing.T) {
	handles.Init()
	defer handles.Release()
	t.Chdir(t.TempDir())
	if err := os.Mkdir("mapgen", 0700); err != nil {
		t.Fatal(err)
	}
	decor := "DECOR ROOM base WALL_FLOOR PaintWall PaintTile END DECOR HALL hall WALL_FLOOR PaintWall PaintTile END "
	var cases []legacy.PortTestPaintSpec
	for _, size := range []int{200, 400} {
		for _, room := range []int{5, 7} {
			for _, limit := range []int{2, 3} {
				for _, branch := range []int{0, 100} {
					for seed := 0; seed < 8; seed++ {
						name := fmt.Sprintf("normal-%03d", len(cases))
						body := fmt.Sprintf("ALGORITHM_DATA mapSize %d midRoomSize %d roomVariance 0 recursionLimit %d hallBranchRate %d seed %d END ", size, room, limit, branch, seed) + decor
						orchestrationTheme(t, name, body)
						cases = append(cases, orchestrationCase(name))
					}
				}
			}
		}
	}
	out := orchestrationCapture(t, "normal-maps", cases)
	for i, r := range out {
		if r.Steps[2].Return != 1 {
			t.Fatalf("normal map %d return %d", i, r.Steps[2].Return)
		}
	}
	cases = nil
	for _, limit := range []int{1, 2} {
		for seed := 0; seed < 8; seed++ {
			name := fmt.Sprintf("ring-%03d", len(cases))
			body := fmt.Sprintf("ALGORITHM_DATA skeleton HALL_RING mapSize 1050 recursionLimit %d seed %d END ", limit, seed) + decor
			orchestrationTheme(t, name, body)
			cases = append(cases, orchestrationCase(name))
		}
	}
	out = orchestrationCapture(t, "ring-maps", cases)
	for i, r := range out {
		if r.Steps[2].Return != 1 {
			t.Fatalf("ring map %d return %d", i, r.Steps[2].Return)
		}
	}
}
func TestMapOrchestrationThemeFailures(t *testing.T) {
	handles.Init()
	defer handles.Release()
	t.Chdir(t.TempDir())
	if err := os.Mkdir("mapgen", 0700); err != nil {
		t.Fatal(err)
	}
	bodies := []string{"", "unknown ", "DECOR ROOM base END ", "DECOR HALL hall END ", "DECOR ROOM base END unknown ", "ALGORITHM_DATA mapSize 200 END DECOR ROOM base WALL_FLOOR PaintWall "}
	cases := []legacy.PortTestPaintSpec{orchestrationCase("missing")}
	for i, body := range bodies {
		name := fmt.Sprintf("invalid-%02d", i)
		orchestrationTheme(t, name, body)
		cases = append(cases, orchestrationCase(name))
	}
	out := orchestrationCapture(t, "theme-failures", cases)
	for i, r := range out {
		if r.Steps[2].Return != 0 {
			t.Fatalf("invalid theme %d return %d", i, r.Steps[2].Return)
		}
	}
}
