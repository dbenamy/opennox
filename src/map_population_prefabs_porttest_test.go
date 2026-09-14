//go:build porttest

package opennox

import (
	"fmt"
	"math"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
)

func populationCacheBase() legacy.PortTestPaintSpec {
	s := populationBase()
	s.Globals["decodedCache"] = roomValue(1)
	s.Globals["dword_5d4594_1599576"] = roomArg(7)
	s.Globals["dword_5d4594_1599596"] = roomValue(1)
	paintString(&s.Records[6], 0, "fixture")
	s.Records[6].Words[64] = math.Float32bits(130.10765)
	s.Records[6].Words[68] = math.Float32bits(260.2153)
	paintString(&s.Records[5], 0, "fixture")
	s.Records[0].Words[12] = 4
	s.Records[0].Words[16] = 4
	return s
}
func TestMapPopulationPrefabMetadata(t *testing.T) {
	var cases []legacy.PortTestPaintSpec
	for _, name := range []string{"fixture", "FIXTURE", "missing"} {
		for mask := 0; mask < 16; mask++ {
			s := populationCacheBase()
			paintString(&s.Records[5], 0, name)
			s.Globals["dword_5d4594_2487656"] = roomValue(6)
			s.Globals["blob2487660"] = roomValue(7)
			s.Globals["blob2487664"] = roomValue(8)
			s.Globals["blob2487668"] = roomValue(9)
			s.Globals["dword_5d4594_1599540"] = roomArg(4)
			previous := 0
			for dir := 0; dir < 4; dir++ {
				if mask&(1<<dir) == 0 {
					continue
				}
				for j := 0; j < 2; j++ {
					slot := 100 + len(s.Objects)
					obj := legacy.PortTestPaintObject{Slot: slot, Type: 6 + dir, Words: map[int]uint32{56: math.Float32bits(float32(46 * (dir + j + 1))), 60: math.Float32bits(float32(46 * (dir + 2)))}, Refs: map[int]legacy.PortTestMapRoomArg{}}
					if previous == 0 {
						s.Records[3].Refs[0] = roomArg(slot)
					} else {
						s.Objects[len(s.Objects)-1].Refs[444] = roomArg(slot)
					}
					previous = slot
					s.Objects = append(s.Objects, obj)
				}
			}
			s.Actions = []legacy.PortTestPaintAction{paintAction(29, roomArg(1), roomArg(6))}
			cases = append(cases, s)
		}
	}
	out := populationCapture(t, "prefab-metadata", cases)
	for i, r := range out {
		step := r.Steps[0]
		if i >= 32 {
			if step.Return != 0 || step.Globals["prefabLoadCount"] != 0 {
				t.Fatal("missing prefab load")
			}
			continue
		}
		if step.Globals["prefabLoadCount"] != 1 {
			t.Fatal("loader request count")
		}
		for _, rec := range step.Records {
			if rec.ID != step.Slots[6] {
				continue
			}
			for dir := 0; dir < 4; dir++ {
				want := uint32(0)
				if (i%16)&(1<<dir) != 0 {
					want = 2
				}
				if rec.Words[(88+16*dir)/4] != want {
					t.Fatalf("case %d exit %d count", i, dir)
				}
			}
		}
	}
}
func TestMapPopulationPrefabRooms(t *testing.T) {
	var cases []legacy.PortTestPaintSpec
	for seed := 0; seed < 32; seed++ {
		s := populationCacheBase()
		s.Seed = seed
		s.Globals["occupancyEnabled"] = roomValue(1)
		s.Records[0].Words[68] = 32
		s.Records[0].Words[84] = 1
		s.Actions = []legacy.PortTestPaintAction{paintAction(26, roomArg(1)), paintAction(27, roomArg(1), roomArg(6))}
		cases = append(cases, s)
	}
	out := populationCapture(t, "prefab-rooms", cases)
	for i, r := range out {
		step := r.Steps[1]
		if step.Return != 1 || step.Globals["roomGlobal4"] == 0 {
			t.Fatalf("case %d no prefab room", i)
		}
		occupied := 0
		for _, rec := range step.Records {
			if rec.Kind == "occupancy" {
				for j := 4; j < len(rec.Words); j += 5 {
					if rec.Words[j] != 0 {
						occupied++
					}
				}
			}
		}
		if occupied == 0 {
			t.Fatalf("case %d empty occupancy", i)
		}
	}
}
func TestMapPopulationMetadataAllocation(t *testing.T) {
	var cases []legacy.PortTestPaintSpec
	for _, name := range []string{"first-second-n", "first-second-s", "first-second-ew", "first-second-c", "first-second-nsewc", "plain", "first-second-N", "first-second-S", "first-second-EW", "first-second-C", "first-second-NSEWC"} {
		s := populationCacheBase()
		paintString(&s.Records[6], 0, name)
		s.Actions = []legacy.PortTestPaintAction{paintAction(34), paintAction(35)}
		cases = append(cases, s)
	}
	out := populationCapture(t, "metadata-allocation", cases)
	for i, r := range out {
		found := false
		for _, rec := range r.Steps[0].Records {
			if rec.Kind == "prefabMetadata" {
				found = true
				if !rec.Alive {
					t.Fatal("metadata prematurely freed")
				}
				want := []uint32{0, 0, 0, 0, 0, 0, 1, 2, 12, 16, 31}[i]
				if rec.Words[15] != want {
					t.Fatalf("case %d prefab flags %d want %d", i, rec.Words[15], want)
				}
			}
		}
		if !found {
			t.Fatal("metadata not allocated")
		}
		for _, rec := range r.Steps[1].Records {
			if rec.Kind == "prefabMetadata" && rec.Alive {
				t.Fatalf("case %d metadata not freed", i)
			}
		}
	}
}

func TestMapPopulationPrefabSelection(t *testing.T) {
	var cases []legacy.PortTestPaintSpec
	for n := 0; n <= 7; n++ {
		for required := 0; required < 3; required++ {
			for seed := 0; seed < 8; seed++ {
				s := populationCacheBase()
				s.Seed = seed
				s.Globals["occupancyEnabled"] = roomValue(1)
				s.Records[0].Words[68] = 32
				s.Records[0].Words[84] = uint32(n)
				s.Records[6] = roomRecord(76 * max(n, 1))
				s.Globals["dword_5d4594_1599596"] = roomValue(int32(n))
				for i := 0; i < n; i++ {
					slot := 9 + i
					s.Records = append(s.Records, roomRecord(160))
					r := &s.Records[slot-1]
					name := fmt.Sprintf("prefab%d", i)
					paintString(r, 0, name)
					paintString(&s.Records[6], i*76, name)
					s.Records[6].Words[i*76+64] = math.Float32bits(65.053826)
					s.Records[6].Words[i*76+68] = math.Float32bits(65.053826)
					if required == 1 || (required == 2 && i%2 == 0) {
						r.Words[72] = 1
					}
					if i+1 < n {
						r.Refs[156] = roomArg(slot + 1)
					}
				}
				if n > 0 {
					s.Records[0].Refs[80] = roomArg(9)
				}
				s.Actions = []legacy.PortTestPaintAction{paintAction(25, roomArg(1))}
				cases = append(cases, s)
			}
		}
	}
	out := populationCapture(t, "prefab-selection", cases)
	for i, r := range out {
		n := int(cases[i].Records[0].Words[84])
		placed := 0
		for _, rec := range r.Steps[0].Records {
			if rec.Kind == "input" && len(rec.Words) == 40 && rec.Words[19] != 0 {
				placed++
			}
		}
		if placed > 5 || placed > n {
			t.Fatalf("case %d prefab limit", i)
		}
		if n == 0 && r.Steps[0].Return != 1 {
			t.Fatal("empty prefab list")
		}
		if n == 1 && (placed != 1 || r.Steps[0].Return != 1) {
			t.Fatalf("case %d single prefab not placed", i)
		}
	}
}
func TestMapPopulationPrefabReplacement(t *testing.T) {
	var cases []legacy.PortTestPaintSpec
	for seed := 0; seed < 16; seed++ {
		s := populationCacheBase()
		s.Seed = seed
		s.Globals["occupancyEnabled"] = roomValue(1)
		s.Records[0].Words[68] = 32
		s.Records[0].Words[84] = 1
		s.Records[0].Refs[80] = roomArg(6)
		s.Actions = []legacy.PortTestPaintAction{paintAction(25, roomArg(1)), paintAction(31, roomArg(1))}
		cases = append(cases, s)
	}
	out := populationCapture(t, "prefab-replacement", cases)
	for i, r := range out {
		if r.Steps[0].Return != 1 || r.Steps[1].Return != 1 {
			t.Fatalf("case %d replacement failed", i)
		}
		freed := 0
		for _, rec := range r.Steps[1].Records {
			if rec.Kind == "input" && !rec.Alive {
				freed++
			}
		}
		if freed != 1 {
			t.Fatalf("case %d freed rooms %d", i, freed)
		}
	}
}
