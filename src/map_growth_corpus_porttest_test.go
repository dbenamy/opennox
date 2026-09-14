//go:build porttest

package opennox

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
)

// Locked from corrected C, repeated in default, server and highres builds.
var growthCExpected = map[string]string{
	"branch-choices":    "d6ed89ed995f785dcb6a49c478fa3bb0f230255d61566220f0b4e1da69c1a7fe",
	"door-translations": "725d8586be855808d39c36b8e25cee6207b05094c4496e327269e4206d4f52d7",
	"doors":             "663f55684abee8b66634fc91fad8eea5831fb7700658d396596a2e90faa70425",
	"frontiers":         "d32a70cd608e37efe9281e04d35e48facd061cf6ec78a3280d5eeba2516c1cb4",
	"hall-expansion":    "28988fcc83acfe40787de912b9568d4c3466a20bfd5b8ad81ca00d2ceb94805f",
	"initial-layouts":   "7b4934c4dc0ecb073e19a9fcbe115e9b9f523b1d21aaedc909150f07160ad6b4",
	"obstructions":      "aff01ff1f5016f212af4df91d5ce5ca78a8527622d64a2e02b15d1ca0b52153d",
	"recursion-guards":  "d3597f23d62b2fc2774e53986edf3ee55eb027be0da62631da07a98185d312b3",
	"room-expansion":    "d102d5a027e605001f44c896b7c9121735b40e3cbcfaa55363e818ba62c3b1b4",
}

func growthCapture(t *testing.T, label string, cases []legacy.PortTestPaintSpec) []legacy.PortTestPaintResult {
	t.Helper()
	out := growthRun(cases)
	for i, r := range out {
		if !r.Intact || !r.ControlOK {
			t.Fatalf("%s case %d guard/control state", label, i)
		}
	}
	data, err := json.Marshal(out)
	if err != nil {
		t.Fatal(err)
	}
	if prefix := os.Getenv("OPENNOX_MAP_GROWTH_CAPTURE"); prefix != "" {
		if err = os.WriteFile(prefix+"-"+label+".json", data, 0600); err != nil {
			t.Fatal(err)
		}
	}
	hash := fmt.Sprintf("%x", sha256.Sum256(data))
	if want, ok := growthCExpected[label]; ok {
		if hash != want {
			t.Fatalf("%s differs from C: got %s want %s", label, hash, want)
		}
	} else {
		t.Fatalf("missing C baseline for %s", label)
	}
	t.Logf("%s: %d cases %s", label, len(cases), hash)
	return out
}

func TestMapGrowthBranchChoices(t *testing.T) {
	var cases []legacy.PortTestPaintSpec
	for kind := uint32(0); kind <= 6; kind++ {
		for _, dims := range [][2]uint32{{0, 0}, {3, 7}, {7, 3}, {5, 5}, {0xffffffff, 7}, {7, 0xffffffff}} {
			for _, branch := range []int32{-1, 0, 1, 50, 99, 100, 101} {
				for _, room := range []int32{0, 50, 100} {
					for _, seed := range []int{0, 1, 17, 999} {
						s := growthBase()
						s.Seed = seed
						s.Globals["growthInitGrid"] = roomValue(0)
						s.Records[0].Words[24] = uint32(branch)
						s.Records[0].Words[28] = uint32(room)
						s.Records[1].Words[0] = kind
						s.Records[1].Words[12] = dims[0]
						s.Records[1].Words[16] = dims[1]
						s.Actions = []legacy.PortTestPaintAction{paintAction(6, roomArg(2))}
						cases = append(cases, s)
					}
				}
			}
		}
	}
	growthCapture(t, "branch-choices", cases)
}

func TestMapGrowthInitialLayouts(t *testing.T) {
	var cases []legacy.PortTestPaintSpec
	for _, mean := range []uint32{0, 1, 5, 9, 15} {
		for _, spread := range []uint32{0, 1, 5} {
			for seed := 0; seed < 8; seed++ {
				s := growthBase()
				s.Seed = seed
				s.Records[0].Words[32] = mean
				s.Records[0].Words[36] = spread
				s.Actions = []legacy.PortTestPaintAction{paintAction(1, roomArg(1))}
				cases = append(cases, s)
			}
		}
	}
	for _, radius := range []uint32{0, 16, 32} {
		for seed := 0; seed < 16; seed++ {
			s := growthBase()
			s.Seed = seed
			s.Records[0].Words[0] = 1
			s.Records[0].Words[68] = radius
			s.Actions = []legacy.PortTestPaintAction{paintAction(1, roomArg(1))}
			cases = append(cases, s)
		}
	}
	growthCapture(t, "initial-layouts", cases)
}

func TestMapGrowthRecursionGuards(t *testing.T) {
	var cases []legacy.PortTestPaintSpec
	for kind := uint32(1); kind <= 5; kind++ {
		for _, depth := range []int32{0, 1, 2, 9, 100} {
			for _, limit := range []int32{0, 1, depth, depth + 1} {
				s := growthBase()
				s.Globals["growthInitGrid"] = roomValue(0)
				s.Records[0].Words[72] = uint32(limit)
				s.Records[1].Words[0] = kind
				s.Actions = []legacy.PortTestPaintAction{paintAction(3, roomArg(2), roomValue(depth), roomValue(0), roomValue(0), roomValue(0))}
				cases = append(cases, s)
			}
		}
	}
	out := growthCapture(t, "recursion-guards", cases)
	for i, r := range out {
		if r.Steps[0].Return != 0 {
			t.Fatalf("guard case %d did not stop", i)
		}
	}
}

func TestMapGrowthRoomExpansion(t *testing.T) {
	var cases []legacy.PortTestPaintSpec
	for _, depth := range []uint32{1, 2, 3, 4} {
		for _, limit := range []float32{0, 100, 500, 1000} {
			for seed := 0; seed < 8; seed++ {
				s := growthBase()
				s.Seed = seed
				s.Records[0].Words[64] = math.Float32bits(limit)
				s.Records[0].Words[72] = depth
				s.Globals["roomGlobal4"] = roomArg(2)
				roomSetGeometry(&s.Records[1], 1, 7, 7, 0, 0)
				s.Actions = []legacy.PortTestPaintAction{paintAction(3, roomArg(2), roomValue(0), roomValue(0), roomValue(0), roomArg(2))}
				cases = append(cases, s)
			}
		}
	}
	growthCapture(t, "room-expansion", cases)
}

func TestMapGrowthHallExpansion(t *testing.T) {
	var cases []legacy.PortTestPaintSpec
	for kind := uint32(2); kind <= 5; kind++ {
		for _, branch := range []uint32{0, 50, 100} {
			for _, room := range []uint32{0, 50, 100} {
				for seed := 0; seed < 4; seed++ {
					s := growthBase()
					s.Seed = seed
					s.Records[0].Words[64] = math.Float32bits(1000)
					s.Records[0].Words[72] = 3
					s.Records[0].Words[24] = branch
					s.Records[0].Words[28] = room
					s.Globals["roomGlobal4"] = roomArg(2)
					w, h := int32(3), int32(7)
					if kind >= 4 {
						w, h = h, w
					}
					roomSetGeometry(&s.Records[1], int32(kind), w, h, 0, 0)
					s.Actions = []legacy.PortTestPaintAction{paintAction(5, roomArg(2), roomValue(1), roomValue(0), roomValue(3), roomArg(2))}
					cases = append(cases, s)
				}
			}
		}
	}
	growthCapture(t, "hall-expansion", cases)
}
