//go:build porttest

package opennox

import (
	"math"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
)

func TestMapGrowthFillProbe(t *testing.T) {
	s := growthBase()
	s.Records[0].Words[64] = math.Float32bits(1000)
	s.Records[0].Words[72] = 2
	s.Globals["roomGlobal4"] = roomArg(2)
	roomSetGeometry(&s.Records[1], 1, 5, 5, 0, 0)
	s.Actions = []legacy.PortTestPaintAction{paintAction(4, roomArg(2), roomValue(1), roomValue(0), roomValue(0), roomArg(2))}
	out := growthRun([]legacy.PortTestPaintSpec{s})
	r := out[0]
	if !r.Intact || !r.ControlOK || r.Steps[0].Return != 1 {
		t.Fatal("single-level room fill failed")
	}
	var mask uint32
	attempts := 0
	for _, row := range r.Steps[0].Records {
		if row.Kind == "growthReleased" && len(row.Words) == 95 {
			kind := row.Words[1]
			if kind < 2 || kind > 5 {
				t.Fatalf("room fill created invalid hallway kind %d", kind)
			}
			mask |= 1 << (kind - 2)
			attempts++
		}
	}
	if mask != 15 || attempts != 8 {
		t.Fatalf("room fill attempted directions mask %x, attempts %d; want f and 8", mask, attempts)
	}
}

func TestMapGrowthFillDirections(t *testing.T) {
	var cases []legacy.PortTestPaintSpec
	for blocked := 0; blocked < 16; blocked++ {
		for seed := 0; seed < 8; seed++ {
			s := growthBase()
			s.Seed = seed
			s.Records[0].Words[64] = math.Float32bits(1000)
			s.Records[0].Words[72] = 2
			s.Globals["roomGlobal4"] = roomArg(2)
			roomSetGeometry(&s.Records[1], 1, 5, 5, 0, 0)
			roomSetGeometry(&s.Records[2], 1, 5, 5, 500, 500)
			for dir := 0; dir < 4; dir++ {
				if blocked&(1<<dir) != 0 {
					s.Records[1].Words[216] |= 1 << uint(8*dir)
					s.Records[1].Refs[88+32*dir] = roomArg(3)
				}
			}
			s.Actions = []legacy.PortTestPaintAction{paintAction(4, roomArg(2), roomValue(1), roomValue(0), roomValue(0), roomArg(2))}
			cases = append(cases, s)
		}
	}
	out := growthRun(cases)
	for i, r := range out {
		if !r.Intact || !r.ControlOK || r.Steps[0].Return != 1 {
			t.Fatalf("room fill case %d", i)
		}
		mask, attempts := uint32(0), 0
		for _, row := range r.Steps[0].Records {
			if row.Kind == "growthReleased" && len(row.Words) == 95 {
				kind := row.Words[1]
				if kind < 2 || kind > 5 {
					t.Fatalf("case %d invalid hallway kind %d", i, kind)
				}
				mask |= 1 << (kind - 2)
				attempts++
			}
		}
		want := uint32(15 ^ (i / 8))
		wantAttempts := 8
		if want == 0 {
			wantAttempts = 0
		}
		if mask != want || attempts != wantAttempts {
			t.Fatalf("case %d got directions %x / %d attempts, want %x / %d", i, mask, attempts, want, wantAttempts)
		}
	}
}

func TestMapGrowthMergeProbe(t *testing.T) {
	var cases []legacy.PortTestPaintSpec
	for _, named := range []uint32{0, 100} {
		s := growthBase()
		s.Seed = 0
		s.Records[0].Words[64] = math.Float32bits(1000)
		s.Records[0].Words[72] = 2
		s.Records[0].Words[48] = 100
		s.Globals["roomGlobal4"] = roomArg(2)
		s.Globals["growthMergeRate"] = roomValue(int32(named))
		roomSetGeometry(&s.Records[1], 1, 5, 5, 0, 0)
		roomSetGeometry(&s.Records[2], 1, 5, 5, 0, float32(-8*32.526913))
		s.Records[1].Refs[56] = roomArg(3)
		s.Records[2].Refs[60] = roomArg(2)
		s.Actions = []legacy.PortTestPaintAction{paintAction(4, roomArg(2), roomValue(1), roomValue(0), roomValue(0), roomArg(2))}
		cases = append(cases, s)
	}
	out := growthRun(cases)
	counts := make([]uint32, len(out))
	for i, r := range out {
		if !r.Intact || !r.ControlOK || r.Steps[0].Return != 1 {
			t.Fatalf("merge fixture %d", i)
		}
		for _, row := range r.Steps[0].Records {
			if row.ID == r.Steps[0].Slots[2] {
				counts[i] = row.Words[54] & 255
			}
		}
	}
	t.Logf("configured merge 100, named values 0/100: north connections %v", counts)
	if counts[0] == 0 || counts[0] != counts[1] {
		t.Fatal("configured merge rate does not control room merging")
	}
}

func TestMapGrowthMergeSettings(t *testing.T) {
	var cases []legacy.PortTestPaintSpec
	for dir := 0; dir < 4; dir++ {
		for _, rate := range []uint32{0, 50, 100} {
			for seed := 0; seed < 8; seed++ {
				for _, named := range []uint32{0, 100} {
					s := growthBase()
					s.Seed = seed
					s.Records[0].Words[64] = math.Float32bits(1000)
					s.Records[0].Words[72] = 2
					s.Records[0].Words[48] = rate
					s.Globals["roomGlobal4"] = roomArg(2)
					s.Globals["growthMergeRate"] = roomValue(int32(named))
					x, y := float32(0), float32(-8*32.526913)
					switch dir {
					case 1:
						y = float32(8 * 32.526913)
					case 2:
						x, y = float32(8*32.526913), 0
					case 3:
						x, y = float32(-8*32.526913), 0
					}
					roomSetGeometry(&s.Records[1], 1, 5, 5, 0, 0)
					roomSetGeometry(&s.Records[2], 1, 5, 5, x, y)
					s.Records[1].Refs[56] = roomArg(3)
					s.Records[2].Refs[60] = roomArg(2)
					s.Actions = []legacy.PortTestPaintAction{paintAction(4, roomArg(2), roomValue(1), roomValue(0), roomValue(0), roomArg(2))}
					cases = append(cases, s)
				}
			}
		}
	}
	out := growthRun(cases)
	counts := make([]uint32, len(out))
	for i, r := range out {
		if !r.Intact || !r.ControlOK || r.Steps[0].Return != 1 {
			t.Fatalf("merge case %d", i)
		}
		dir := i / 48
		rate := (i % 48) / 16
		for _, row := range r.Steps[0].Records {
			if row.ID == r.Steps[0].Slots[2] {
				counts[i] = (row.Words[54] >> uint(8*dir)) & 255
			}
		}
		if rate == 0 && counts[i] != 0 {
			t.Fatalf("case %d merged with zero rate", i)
		}
		if rate == 2 && counts[i] == 0 {
			t.Fatalf("case %d failed to merge with full rate", i)
		}
		if i%2 == 1 && (counts[i] != counts[i-1] || r.RandomTail != out[i-1].RandomTail) {
			t.Fatalf("case %d separate named word changed configured merge behavior", i)
		}
	}
}
