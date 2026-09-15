//go:build porttest

package opennox

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"math"
	"os"
	"testing"
)

func hallwayConnection(dir, x, y, width int) legacy.PortTestPaintSpec {
	s := populationBase()
	s.Globals["occupancyEnabled"] = roomValue(1)
	s.Globals["roomGlobal4"] = roomArg(2)
	s.Records[0].Words[68] = 32
	s.Records[0].Refs[80] = roomArg(6)
	roomSetGeometry(&s.Records[1], 1, 4, 4, 0, 0)
	roomSetGeometry(&s.Records[2], 1, 4, 4, float32(float64(x)*32.526913), float32(float64(y)*32.526913))
	s.Records[1].Refs[56] = roomArg(3)
	s.Records[2].Refs[60] = roomArg(2)
	p := &s.Records[5]
	p.Words[60] = math.Float32bits(130.10765)
	p.Words[64] = math.Float32bits(130.10765)
	p.Words[76] = 1
	p.Refs[148] = roomArg(2)
	mx, my := 2, 0
	switch dir {
	case 1:
		my = 3
	case 2:
		mx, my = 3, 2
	case 3:
		mx, my = 0, 2
	}
	p.Words[80+16*dir] = uint32(mx)
	p.Words[84+16*dir] = uint32(my)
	p.Words[88+16*dir] = uint32(width)
	p.Words[92+16*dir] = 1
	s.Actions = []legacy.PortTestPaintAction{paintAction(31, roomArg(1))}
	return s
}

func TestMapHallwaysBentConnection(t *testing.T) {
	for dir := 0; dir < 4; dir++ {
		for _, offset := range []int{-6, 6} {
			for _, width := range []int{1, 2} {
				x, y := offset, -12
				switch dir {
				case 1:
					y = 12
				case 2:
					x, y = 12, offset
				case 3:
					x, y = -12, offset
				}
				out := hallwayRun([]legacy.PortTestPaintSpec{hallwayConnection(dir, x, y, width)})
				if len(out) != 1 || !out[0].Intact || !out[0].ControlOK || out[0].Steps[0].Return != 1 {
					t.Fatalf("dir %d offset %d width %d: bent hallway admission", dir, offset, width)
				}
				halls := 0
				for _, r := range out[0].Steps[0].Records {
					if r.Alive && r.Kind == "input" && len(r.Words) == 94 && r.Words[0] >= 2 && r.Words[0] <= 5 {
						halls++
					}
				}
				if halls != 3 {
					t.Fatalf("dir %d offset %d width %d: hallways got %d want 3", dir, offset, width, halls)
				}
			}
		}
	}
}

func hallwayCapture(t *testing.T, label string, cases []legacy.PortTestPaintSpec) []legacy.PortTestPaintResult {
	t.Helper()
	out := hallwayRun(cases)
	for i, r := range out {
		if !r.Intact || !r.ControlOK {
			t.Fatalf("%s case %d guard/control state", label, i)
		}
	}
	data, err := json.Marshal(out)
	if err != nil {
		t.Fatal(err)
	}
	if prefix := os.Getenv("OPENNOX_MAP_HALLWAYS_CAPTURE"); prefix != "" {
		if err := os.WriteFile(prefix+"-"+label+".json", data, 0600); err != nil {
			t.Fatal(err)
		}
	}
	want, ok := hallwayCExpected[label]
	if !ok {
		t.Fatalf("missing locked C capture for %s", label)
	}
	got := fmt.Sprintf("%x", sha256.Sum256(data))
	if got != want {
		// Preserve intermittent failures even when capture was not requested.
		if os.Getenv("OPENNOX_MAP_HALLWAYS_CAPTURE") == "" {
			const dir = "../build/port-failures"
			if err := os.MkdirAll(dir, 0700); err != nil {
				t.Logf("cannot create mismatch artifact directory: %v", err)
			} else if f, err := os.CreateTemp(dir, "map-hallways-"+label+"-*.json"); err != nil {
				t.Logf("cannot create mismatch artifact: %v", err)
			} else {
				_, writeErr := f.Write(data)
				closeErr := f.Close()
				t.Logf("full mismatch capture: %s (write error: %v; close error: %v)", f.Name(), writeErr, closeErr)
			}
		}
		t.Fatalf("%s complete capture differs from C: got %s want %s", label, got, want)
	}
	t.Logf("%s: %d cases %x", label, len(cases), sha256.Sum256(data))
	return out
}

func TestMapHallwaysRoutes(t *testing.T) {
	var cases []legacy.PortTestPaintSpec
	for dir := 0; dir < 4; dir++ {
		for _, gap := range []int{1, 2, 3, 4, 8, 12} {
			for _, width := range []int{1, 2, 3} {
				for _, offset := range []int{-12, -6, -3, -2, -1, 0, 1, 2, 3, 6, 12} {
					for _, dim := range [][2]int{{4, 4}, {3, 7}, {7, 3}} {
						x, y := offset, -dim[1]-gap
						switch dir {
						case 1:
							y = 4 + gap
						case 2:
							x, y = 4+gap, offset
						case 3:
							x, y = -dim[0]-gap, offset
						}
						s := hallwayConnection(dir, x, y, width)
						roomSetGeometry(&s.Records[2], 1, int32(dim[0]), int32(dim[1]), float32(float64(x)*32.526913), float32(float64(y)*32.526913))
						cases = append(cases, s)
					}
				}
			}
		}
	}
	out := hallwayCapture(t, "routes", cases)
	counts := map[int]int{}
	for i, r := range out {
		step := r.Steps[0]
		hallwayTopology(t, i, step)
		halls := 0
		for _, rec := range step.Records {
			if rec.Alive && rec.Kind == "input" && len(rec.Words) == 94 && rec.Words[0] >= 2 && rec.Words[0] <= 5 {
				halls++
			}
		}
		counts[halls]++
		if (step.Return == 0 && halls != 0) || (step.Return != 0 && (halls < 1 || halls > 3)) {
			t.Fatalf("case %d return %d corridor count %d", i, step.Return, halls)
		}
	}
	t.Logf("admitted hallway counts: %v", counts)
	for _, n := range []int{1, 2, 3} {
		if counts[n] == 0 {
			t.Fatalf("missing %d-segment route coverage", n)
		}
	}
}

func TestMapHallwaysObstructions(t *testing.T) {
	var cases []legacy.PortTestPaintSpec
	for dir := 0; dir < 4; dir++ {
		for _, offset := range []int{-6, 0, 6} {
			for _, width := range []int{1, 2} {
				for _, barrier := range []int{0, 1, 2} {
					x, y := offset, -12
					switch dir {
					case 1:
						y = 12
					case 2:
						x, y = 12, offset
					case 3:
						x, y = -12, offset
					}
					s := hallwayConnection(dir, x, y, width)
					r := roomRecord(376)
					bx, by, bw, bh := int32(-24), int32(-4), int32(48), int32(1)
					switch dir {
					case 1:
						by = 7
					case 2:
						bx, by, bw, bh = 7, -24, 1, 48
					case 3:
						bx, by, bw, bh = -4, -24, 1, 48
					}
					if barrier == 0 {
						bx, by = 24, 24
						bw, bh = 1, 1
					}
					if barrier == 1 {
						if dir < 2 {
							bx, bw = 2, 1
						} else {
							by, bh = 2, 1
						}
					}
					roomSetGeometry(&r, 2, bw, bh, float32(float64(bx)*32.526913), float32(float64(by)*32.526913))
					s.Records = append(s.Records, r)
					s.Records[2].Refs[56] = roomArg(9)
					s.Records[8].Refs[60] = roomArg(3)
					cases = append(cases, s)
				}
			}
		}
	}
	out := hallwayCapture(t, "obstructions", cases)
	rejected := 0
	for i, r := range out {
		step := r.Steps[0]
		hallwayTopology(t, i, step)
		if i%3 == 2 && step.Return != 0 {
			t.Fatalf("case %d crossed complete barrier", i)
		}
		if i%3 == 0 && step.Return != 1 {
			t.Fatalf("case %d unobstructed route rejected", i)
		}
		if step.Return == 0 {
			rejected++
			for _, rec := range step.Records {
				if rec.Alive && rec.Kind == "input" && len(rec.Words) == 94 && rec.Words[0] >= 2 && rec.Words[0] <= 5 && rec.ID != step.Slots[9] {
					t.Fatalf("case %d left tentative corridor alive", i)
				}
			}
		}
	}
	if rejected == 0 {
		t.Fatal("missing rejection coverage")
	}
}

// Every admitted edge must name a live room with a reciprocal opposite edge.
// Successful routes must reach the intended candidate from the replacement prefab.
func hallwayTopology(t *testing.T, index int, step legacy.PortTestPaintStep) {
	t.Helper()
	rows := map[uint32][]uint32{}
	for _, rec := range step.Records {
		if rec.Alive {
			rows[rec.ID] = rec.Words
		}
	}
	for id, row := range rows {
		if len(row) != 94 {
			continue
		}
		for dir := 0; dir < 4; dir++ {
			n := int(row[54] >> uint(8*dir) & 255)
			if n > 8 {
				t.Fatalf("case %d room %x direction %d neighbor count %d", index, id, dir, n)
			}
			for j := 0; j < n; j++ {
				other := row[22+8*dir+j]
				target := rows[other]
				if len(target) != 94 {
					t.Fatalf("case %d dangling room edge", index)
				}
				opposite := dir ^ 1
				back := false
				for k := 0; k < int(target[54]>>uint(8*opposite)&255); k++ {
					if target[22+8*opposite+k] == id {
						back = true
					}
				}
				if !back {
					t.Fatalf("case %d asymmetric room edge", index)
				}
			}
		}
	}
	if step.Return == 0 {
		return
	}
	start := rows[step.Slots[6]][37]
	wanted := step.Slots[3]
	seen := map[uint32]bool{}
	todo := []uint32{start}
	for len(todo) > 0 {
		id := todo[len(todo)-1]
		todo = todo[:len(todo)-1]
		if seen[id] {
			continue
		}
		seen[id] = true
		row := rows[id]
		for dir := 0; dir < 4; dir++ {
			for j := 0; j < int(row[54]>>uint(8*dir)&255); j++ {
				todo = append(todo, row[22+8*dir+j])
			}
		}
	}
	if !seen[wanted] {
		t.Fatalf("case %d admitted route does not reach candidate", index)
	}
}

func TestMapHallwaysCandidateFallback(t *testing.T) {
	var cases []legacy.PortTestPaintSpec
	for dir := 0; dir < 4; dir++ {
		for _, offset := range []int{-6, 6} {
			for _, width := range []int{2, 3} {
				for _, present := range []bool{false, true} {
					for _, reverse := range []bool{false, true} {
						x, y := offset, -12
						nx, ny := -2, -2
						switch dir {
						case 1:
							y = 12
							ny = 6
						case 2:
							x, y = 12, offset
							nx, ny = 6, -2
						case 3:
							x, y = -12, offset
							nx, ny = -2, -2
						}
						s := hallwayConnection(dir, x, y, width)
						next := roomRecord(376)
						roomSetGeometry(&next, 1, 4, 4, float32(float64(x)*32.526913), float32(float64(y)*32.526913))
						next.Refs = map[int]legacy.PortTestMapRoomArg{}
						roomSetGeometry(&s.Records[2], 1, 1, 1, float32(float64(nx)*32.526913), float32(float64(ny)*32.526913))
						s.Records = append(s.Records, next)
						if present {
							s.Records[2].Refs[56] = roomArg(9)
							s.Records[8].Refs[60] = roomArg(3)
							if reverse {
								s.Records[1].Refs[56] = roomArg(9)
								s.Records[8].Refs[60] = roomArg(2)
								s.Records[8].Refs[56] = roomArg(3)
								s.Records[2].Refs[60] = roomArg(9)
								delete(s.Records[2].Refs, 56)
							}
						}
						cases = append(cases, s)
					}
				}
			}
		}
	}
	out := hallwayCapture(t, "fallback", cases)
	for i, r := range out {
		step := r.Steps[0]
		want := uint32(i / 2 % 2)
		if step.Return != want {
			t.Fatalf("case %d candidate fallback got %d want %d", i, step.Return, want)
		}
		for _, rec := range step.Records {
			if rec.ID == step.Slots[3] && rec.Words[54] != 0 {
				t.Fatalf("case %d connected undersized first candidate", i)
			}
			if want == 1 && rec.ID == step.Slots[9] && rec.Words[54] == 0 {
				t.Fatalf("case %d missing fallback candidate connection", i)
			}
		}
	}
}

// Initialize the real startup direction table for hallway topology contracts.
// Restore it so older complete-capture fixtures keep their original inputs.
func hallwayRun(cases []legacy.PortTestPaintSpec) []legacy.PortTestPaintResult {
	var previous [4]uint32
	for i, v := range [4]uint32{1, 0, 3, 2} {
		p := memmap.PtrUint32(0x587000, 254952+uintptr(i)*4)
		previous[i] = *p
		*p = v
	}
	defer func() {
		for i, v := range previous {
			*memmap.PtrUint32(0x587000, 254952+uintptr(i)*4) = v
		}
	}()
	return populationRun(cases)
}
