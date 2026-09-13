//go:build porttest

package opennox

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

func populationBase() legacy.PortTestPaintSpec {
	s := legacy.PortTestPaintSpec{Seed: 12345, Globals: map[string]legacy.PortTestMapRoomArg{"gameFlags": roomValue(0x200000)}}
	for _, n := range []int{1120, 376, 376, 32, 256, 160, 152, 192} {
		s.Records = append(s.Records, roomRecord(n))
	}
	roomSetGeometry(&s.Records[1], 1, 4, 4, 0, 0)
	roomSetGeometry(&s.Records[2], 1, 6, 8, 32.526913, 65.053826)
	s.Records[0].Words[64] = math.Float32bits(1024)
	s.Records[0].Words[68] = 8
	return s
}
func populationRun(cases []legacy.PortTestPaintSpec) []legacy.PortTestPaintResult {
	return legacy.PortTestMapPopulation(cases, func(core *server.Server) (legacy.Server, func()) {
		old := noxServer
		wrapped := &Server{Server: core}
		core.ExtServer = unsafe.Pointer(wrapped)
		noxServer = wrapped
		return wrapped, func() { noxServer = old; core.ExtServer = nil }
	})
}
func populationCapture(t *testing.T, label string, cases []legacy.PortTestPaintSpec) []legacy.PortTestPaintResult {
	t.Helper()
	out := populationRun(cases)
	for i, r := range out {
		if !r.Intact || !r.ControlOK {
			t.Fatalf("case %d guard/control state", i)
		}
	}
	data, err := json.Marshal(out)
	if err != nil {
		t.Fatal(err)
	}
	if prefix := os.Getenv("OPENNOX_MAP_POPULATION_CAPTURE"); prefix != "" {
		if err := os.WriteFile(prefix+"-"+label+".json", data, 0600); err != nil {
			t.Fatal(err)
		}
	}
	t.Logf("%s: %d cases %x", label, len(cases), sha256.Sum256(data))
	return out
}
func TestMapPopulationGlobals(t *testing.T) {
	var cases []legacy.PortTestPaintSpec
	for _, v := range []uint32{0, 1, 0xffffffff, 0x80000000, 0x3f800000, 0x7f800000, 0x7fc00000} {
		s := populationBase()
		for _, n := range []string{"dword_5d4594_1550916", "dword_5d4594_2487576", "dword_5d4594_2487580"} {
			s.Globals[n] = roomValue(int32(v))
		}
		s.Actions = []legacy.PortTestPaintAction{paintAction(0), paintAction(19), paintAction(20)}
		cases = append(cases, s)
	}
	out := populationCapture(t, "globals", cases)
	for i, r := range out {
		if r.Steps[0].Return != cases[i].Globals["dword_5d4594_1550916"].Value || r.Steps[2].Return != r.Steps[0].Return {
			t.Fatal("global return")
		}
	}
}
func TestMapPopulationProgress(t *testing.T) {
	var cases []legacy.PortTestPaintSpec
	for _, flags := range []uint32{0, 0x200000} {
		for _, tick := range []uint32{0, 1, 99, 100, 101, 0x7fffffff, 0x80000000, 0xffffffff} {
			for _, last := range []uint32{0, 1, 100, 0xffffffff} {
				s := populationBase()
				s.Globals["gameFlags"] = roomValue(int32(flags))
				s.Globals["ticks"] = roomValue(int32(tick))
				s.Globals["dword_5d4594_2487568"] = roomValue(int32(last))
				s.Globals["progressThreshold"] = roomValue(100)
				s.Globals["blob2487572"] = roomValue(17)
				s.Globals["guiSequence"] = roomValue(17)
				s.Actions = []legacy.PortTestPaintAction{paintAction(18, roomValue(156))}
				cases = append(cases, s)
			}
		}
	}
	out := populationCapture(t, "progress", cases)
	for i, r := range out {
		g := r.Steps[0].Globals
		if cases[i].Globals["gameFlags"].Value != 0 && g["blob2487572"] != 17 {
			t.Fatal("suppressed progress changed sequence")
		}
	}
}
func TestMapPopulationDistances(t *testing.T) {
	var cases []legacy.PortTestPaintSpec
	for kind := uint32(1); kind <= 6; kind++ {
		for _, dist := range []float32{0, 1, -1, 32.526913, 10000} {
			s := populationBase()
			s.Records[2].Words[0] = kind
			s.Records[1].Refs[88] = roomArg(3)
			s.Records[1].Words[216] = 1
			s.Records[2].Refs[120] = roomArg(2)
			s.Records[2].Words[216] = 256
			s.Actions = []legacy.PortTestPaintAction{paintAction(21, roomArg(2), roomValue(0), roomValue(int32(math.Float32bits(dist)))), paintAction(23, roomArg(2)), paintAction(19), paintAction(20)}
			cases = append(cases, s)
		}
	}
	populationCapture(t, "distances", cases)
}
func TestMapPopulationPrefabCoordinates(t *testing.T) {
	var cases []legacy.PortTestPaintSpec
	for count := uint32(1); count <= 7; count++ {
		for seed := 0; seed < 12; seed++ {
			s := populationBase()
			s.Seed = seed
			s.Records[0].Words[84] = count
			s.Actions = []legacy.PortTestPaintAction{paintAction(26, roomArg(1))}
			for i := uint32(0); i < count && i < 5; i++ {
				s.Actions = append(s.Actions, paintAction(28, roomArg(2), roomArg(4)))
			}
			cases = append(cases, s)
		}
	}
	out := populationCapture(t, "prefab-coordinates", cases)
	for i, r := range out {
		want := len(cases[i].Actions) - 1
		if r.Steps[len(r.Steps)-1].Globals["blob2487608"] != uint32(want) {
			t.Fatal("candidate count")
		}
	}
}
func TestMapPopulationMetadataLookup(t *testing.T) {
	var cases []legacy.PortTestPaintSpec
	for _, name := range []string{"one", "two", "three", "ONE", "missing", ""} {
		s := populationBase()
		paintString(&s.Records[4], 0, name)
		for i, v := range []string{"one", "two", "three"} {
			paintString(&s.Records[7], i*64, v)
		}
		s.Globals["dword_5d4594_2487672"] = roomArg(8)
		s.Globals["dword_5d4594_2487676"] = roomValue(3)
		s.Actions = []legacy.PortTestPaintAction{paintAction(33, roomArg(5)), paintAction(36, roomValue(0)), paintAction(36, roomValue(2))}
		cases = append(cases, s)
	}
	out := populationCapture(t, "metadata-lookup", cases)
	for i, r := range out {
		want := uint32(i)
		if i >= 3 {
			want = 0xffffffff
		}
		if r.Steps[0].Return != want {
			t.Fatal(fmt.Sprintf("lookup %d", i))
		}
	}
}
