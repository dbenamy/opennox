//go:build porttest

package opennox

import (
	"math"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
)

func TestMapPopulationOrdering(t *testing.T) {
	var cases []legacy.PortTestPaintSpec
	for count := 1; count <= 8; count++ {
		for seed := 0; seed < 32; seed++ {
			for _, spawn := range []bool{false, true} {
				s := populationBase()
				s.Seed = seed
				s.Records[1].Refs[372] = roomArg(3)
				s.Records[2].Refs[92] = roomArg(4)
				s.Records[3].Words[0] = uint32(count)
				s.Records[3].Words[4] = uint32(count)
				s.Records[3].Words[12] = uint32(count)
				s.Records[3].Refs[8] = roomArg(9)
				for i := 0; i < count; i++ {
					s.Records = append(s.Records, roomRecord(100))
					r := &s.Records[8+i]
					r.Words[0] = 0
					r.Words[68] = 0
					r.Words[72] = 0
					if spawn {
						r.Words[68] = 1
						r.Words[72] = 1
					}
					paintString(r, 4, "PaintObject")
					if i+1 < count {
						r.Refs[88] = roomArg(10 + i)
					}
				}
				s.Actions = []legacy.PortTestPaintAction{paintAction(8, roomArg(1), roomArg(2))}
				cases = append(cases, s)
			}
		}
	}
	out := populationCapture(t, "population-ordering", cases)
	for i, r := range out {
		step := r.Steps[0]
		count := int(cases[i].Records[3].Words[12])
		rows := map[uint32][]uint32{}
		for _, rec := range step.Records {
			rows[rec.ID] = rec.Words
		}
		head := uint32(0)
		for slot := 9; slot < 9+count; slot++ {
			ptr := step.Slots[slot]
			if rows[ptr][24] == 0 {
				if head != 0 {
					t.Fatal("multiple ordering heads")
				}
				head = ptr
			}
		}
		seen := map[uint32]bool{}
		prev := uint32(0)
		for p := head; p != 0; p = rows[p][23] {
			if seen[p] || len(rows[p]) != 25 || rows[p][24] != prev {
				t.Fatalf("case %d corrupt ordering", i)
			}
			seen[p] = true
			prev = p
		}
		if len(seen) != count {
			t.Fatalf("case %d ordering lost nodes", i)
		}
		objects := 0
		for _, rec := range step.Records {
			if rec.Kind == "object" {
				objects++
			}
		}
		want := 0
		if cases[i].Records[8].Words[68] != 0 {
			want = count
		}
		if objects != want {
			t.Fatalf("case %d population %d want %d", i, objects, want)
		}
	}
}
func TestMapPopulationRoomObjects(t *testing.T) {
	var cases []legacy.PortTestPaintSpec
	for seed := 0; seed < 32; seed++ {
		for _, tag := range []uint32{0, 3, 4, 5} {
			s := populationBase()
			s.Seed = seed
			s.Records[0].Words[1096] = 1
			s.Records[0].Words[548] = 27
			s.Records[5].Words[0] = tag
			s.Records[5].Words[68] = 2
			s.Records[5].Words[72] = 2
			name := "PaintObject"
			if tag == 5 {
				name = "*"
			}
			paintString(&s.Records[5], 4, name)
			s.Actions = []legacy.PortTestPaintAction{paintAction(10, roomArg(1), roomArg(2), roomArg(6))}
			cases = append(cases, s)
		}
	}
	out := populationCapture(t, "room-objects", cases)
	for i, r := range out {
		count := 0
		for _, rec := range r.Steps[0].Records {
			if rec.Kind == "object" {
				count++
			}
		}
		want := 1
		if cases[i].Records[5].Words[0] == 0 {
			want = 2
		}
		if count != want {
			t.Fatalf("case %d placed objects %d want %d", i, count, want)
		}
	}
}
func TestMapPopulationWaypointConnections(t *testing.T) {
	var cases []legacy.PortTestPaintSpec
	for count := 0; count <= 8; count++ {
		s := populationBase()
		s.Globals["roomGlobal4"] = roomArg(2)
		s.Records[1].Words[352] = uint32(count)
		for i := 0; i < count; i++ {
			x, y := float32(i*32), float32(i*16)
			s.Records[1].Words[224+8*i] = math.Float32bits(x)
			s.Records[1].Words[228+8*i] = math.Float32bits(y)
			s.Actions = append(s.Actions, paintAction(15, legacy.PortTestMapRoomArg{Slot: 2, Offset: 224 + 8*i}))
		}
		s.Actions = append(s.Actions, paintAction(17, roomArg(1)))
		cases = append(cases, s)
	}
	out := populationCapture(t, "waypoint-connections", cases)
	for i, r := range out {
		step := r.Steps[len(r.Steps)-1]
		count := 0
		for _, rec := range step.Records {
			if rec.Kind == "waypoint" {
				count++
				if byte(rec.Words[119]) != byte(max(i-1, 0)) {
					t.Fatalf("case %d waypoint connections", i)
				}
			}
		}
		if count != i {
			t.Fatalf("case %d waypoint count %d", i, count)
		}
	}
}
