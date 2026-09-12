//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/opennox/v1/legacy"
	"math"
	"testing"
)

func roomWords(t *testing.T, s legacy.PortTestMapRoomStep, ptr uint32) []uint32 {
	t.Helper()
	for _, r := range s.Regions {
		if r.Alive && ptr >= r.ID && ptr-r.ID < uint32(len(r.Words)*4) {
			return r.Words[(ptr-r.ID)/4:]
		}
	}
	t.Fatalf("missing live captured pointer %d", ptr)
	return nil
}
func roomSlot(t *testing.T, s legacy.PortTestMapRoomStep, slot int) []uint32 {
	t.Helper()
	return roomWords(t, s, s.Slots[slot])
}
func TestMapRoomsRounding(t *testing.T) {
	values := []float32{0, float32(math.Copysign(0, -1)), 1, -1, 32.526913, -32.526913, 65.053826, -65.053826, 1e6, -1e6, 1e11, -1e11}
	for _, v := range []float32{16.263456, -16.263456, 48.79037, -48.79037} {
		values = append(values, math.Nextafter32(v, float32(math.Inf(-1))), v, math.Nextafter32(v, float32(math.Inf(1))))
	}
	var cases []legacy.PortTestMapRoomSpec
	for _, x := range values {
		for _, y := range values {
			s := roomSmoke(0)
			s.NoGrid = true
			s.Records[3].Words[0] = math.Float32bits(x)
			s.Records[3].Words[4] = math.Float32bits(y)
			cases = append(cases, s)
		}
	}
	r := mapRoomsHash(t, "map-rooms-rounding", cases)
	for i, pair := range [][2]int32{{0, 0}, {1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
		s := roomSmoke(0)
		s.NoGrid = true
		s.Records[3].Words[0] = math.Float32bits(float32(float64(pair[0]) * 32.526913))
		s.Records[3].Words[4] = math.Float32bits(float32(float64(pair[1]) * 32.526913))
		v := mapRoomsHash(t, fmt.Sprintf("map-rooms-round-contract-%d", i), []legacy.PortTestMapRoomSpec{s})[0].Steps[0]
		w := roomSlot(t, v, 8)
		if int32(w[0]) != pair[0] || int32(w[1]) != pair[1] {
			t.Fatalf("cell coordinate %v => %v", pair, w[:2])
		}
	}
	if len(r) != len(values)*len(values) {
		t.Fatal("rounding corpus length")
	}
}
func TestMapRoomsGrid(t *testing.T) {
	var cases []legacy.PortTestMapRoomSpec
	var coords [][3]int32
	for _, radius := range []int32{0, 1, 4} {
		for x := -radius - 1; x <= radius+1; x++ {
			for y := -radius - 1; y <= radius+1; y++ {
				s := roomSmoke(1)
				s.Radius = int(radius)
				s.Records[3].Words[0] = uint32(x)
				s.Records[3].Words[4] = uint32(y)
				cases = append(cases, s)
				coords = append(coords, [3]int32{radius, x, y})
			}
		}
	}
	r := mapRoomsHash(t, "map-rooms-grid", cases)
	for i, v := range r {
		c := coords[i]
		got := v.Steps[0].Return[0]
		outside := c[1] < -c[0] || c[1] > c[0] || c[2] < -c[0] || c[2] > c[0]
		if outside {
			if got != 0 {
				t.Fatalf("outside grid %v => %d", c, got)
			}
		} else {
			cell := roomWords(t, v.Steps[0], got)
			if int32(cell[1]) != c[1] || int32(cell[2]) != c[2] || cell[4] != 0 {
				t.Fatalf("grid cell %v => %v", c, cell[:5])
			}
		}
	}
	cases = nil
	for _, x := range []int32{-5, -1, 0, 1, 4, 5} {
		for _, w := range []int32{0, 1, 2, 5} {
			s := roomBase()
			roomSetGeometry(&s.Records[1], 1, w, 2, float32(float64(x)*32.526913), 0)
			s.Actions = []legacy.PortTestMapRoomAction{roomAction(4, roomArg(2)), roomAction(6, roomArg(2)), roomAction(5, roomArg(2)), roomAction(6, roomArg(2))}
			cases = append(cases, s)
		}
	}
	r = mapRoomsHash(t, "map-rooms-occupancy-sequence", cases)
	for i, v := range r {
		if v.Steps[3].Return[0] != 0 {
			t.Fatalf("case %d occupancy remains after clear", i)
		}
	}
	positive := r[2*4+2]
	if positive.Steps[1].Return[0] != positive.Steps[1].Slots[2] {
		t.Fatal("registered room not reported by overlap")
	}
}
func TestMapRoomsConnections(t *testing.T) {
	var cases []legacy.PortTestMapRoomSpec
	for dir := int32(0); dir < 4; dir++ {
		for kind := int32(1); kind <= 5; kind++ {
			s := roomBase()
			s.Records[2].Words[0] = uint32(kind)
			for i := 0; i < 10; i++ {
				s.Actions = append(s.Actions, roomAction(20, roomArg(2), roomArg(3), roomValue(dir)))
			}
			s.Actions = append(s.Actions, roomAction(19, roomArg(2), roomValue(dir)))
			cases = append(cases, s)
		}
	}
	r := mapRoomsHash(t, "map-rooms-connections", cases)
	for i, v := range r {
		for j := 0; j < 10; j++ {
			want := uint32(1)
			if j >= 8 {
				want = 0
			}
			if v.Steps[j].Return[0] != want {
				t.Fatalf("case %d link %d => %d", i, j, v.Steps[j].Return[0])
			}
		}
		dir := i / 5
		w := roomSlot(t, v.Steps[9], 2)
		if byte(w[54]>>(8*dir)) != 8 {
			t.Fatal("connection capacity changed")
		}
	}
	cases = nil
	for dir := int32(0); dir < 4; dir++ {
		s := roomBase()
		s.Actions = []legacy.PortTestMapRoomAction{roomAction(25, roomArg(2), roomArg(3), roomValue(dir))}
		cases = append(cases, s)
	}
	r = mapRoomsHash(t, "map-rooms-reciprocal", cases)
	for i, v := range r {
		s := v.Steps[0]
		w := roomSlot(t, s, 2)
		other := roomSlot(t, s, 3)
		if w[22+8*i] != s.Slots[3] {
			t.Fatal("forward link")
		}
		if other[22+8*(i^1)] != s.Slots[2] {
			t.Fatal("connection is not in the opposite direction")
		}
		found := false
		for _, x := range other[22:54] {
			if x == s.Slots[2] {
				found = true
			}
		}
		if !found {
			t.Fatal("reciprocal link")
		}
	}
}
func TestMapRoomsAllocation(t *testing.T) {
	var cases []legacy.PortTestMapRoomSpec
	for kind := int32(1); kind <= 5; kind++ {
		for _, w := range []int32{-1, 0, 1, 2, 10, 100} {
			for _, h := range []int32{-1, 0, 1, 5} {
				s := roomBase()
				s.NoGrid = true
				a := roomAction(21, roomValue(w), roomValue(h))
				if kind != 1 {
					a = roomAction(47, roomValue(kind), roomValue(w), roomValue(h))
				}
				a.Assign = 10
				s.Actions = []legacy.PortTestMapRoomAction{a, roomAction(23, roomArg(10))}
				cases = append(cases, s)
			}
		}
	}
	r := mapRoomsHash(t, "map-rooms-allocation", cases)
	for i, v := range r {
		s := v.Steps[0]
		w := roomSlot(t, s, 10)
		if len(w) != 94 || w[0] != uint32(i/24+1) {
			t.Fatal("room allocation shape/type")
		}
		for off, x := range w {
			if off != 0 && off != 3 && off != 4 && off != 7 && off != 8 && x != 0 {
				t.Fatalf("case %d uninitialized room word %d", i, off)
			}
		}
		for _, region := range v.Steps[1].Regions {
			if region.ID == s.Slots[10] && region.Alive {
				t.Fatal("freed room still live")
			}
		}
	}
}
func TestMapRoomsLists(t *testing.T) {
	var cases []legacy.PortTestMapRoomSpec
	for _, remove := range []int{2, 3} {
		s := roomBase()
		s.Actions = []legacy.PortTestMapRoomAction{roomAction(13, roomArg(2)), roomAction(13, roomArg(3)), roomAction(14, roomArg(remove)), roomAction(11)}
		cases = append(cases, s)
	}
	r := mapRoomsHash(t, "map-rooms-lists", cases)
	for i, v := range r {
		s := v.Steps[1]
		head := roomSlot(t, s, 3)
		old := roomSlot(t, s, 2)
		if head[14] != s.Slots[2] || head[15] != 0 || old[15] != s.Slots[3] || s.Globals[4] != s.Slots[3] {
			t.Fatal("list insertion links")
		}
		last := v.Steps[3]
		keep := 3
		if i == 1 {
			keep = 2
		}
		if last.Return[0] != last.Slots[keep] || roomSlot(t, last, keep)[15] != 0 {
			t.Fatal("list unlink head/back-link")
		}
	}
	s := roomBase()
	s.Actions = []legacy.PortTestMapRoomAction{roomAction(13, roomArg(2)), roomAction(13, roomArg(3)), roomAction(24), roomAction(11)}
	v := mapRoomsHash(t, "map-rooms-free-list", []legacy.PortTestMapRoomSpec{s})[0]
	if v.Steps[3].Return[0] != 0 || v.Steps[2].Globals[4] != 0 {
		t.Fatal("free list retains head")
	}
}
func TestMapRoomsGeometry(t *testing.T) {
	for _, op := range []int{15, 16, 17, 18, 26, 27, 28, 29, 30, 33, 34, 35, 36, 37, 38, 43, 44, 45, 46} {
		var cases []legacy.PortTestMapRoomSpec
		for kind := int32(0); kind <= 6; kind++ {
			for _, x := range []float32{-65.053826, -0.01, 0, 0.01, 32.526913, 130.10765} {
				s := roomSmoke(op)
				roomSetGeometry(&s.Records[1], kind, 4, 4, x, x)
				s.Records[3].Words[0] = math.Float32bits(x)
				s.Records[3].Words[4] = math.Float32bits(x)
				if op == 34 || op == 36 {
					s.Records[1].Refs[368] = roomArg(9)
				}
				if op == 43 {
					s.Records[2].Words[0] = 1
				}
				cases = append(cases, s)
			}
		}
		mapRoomsHash(t, fmt.Sprintf("map-rooms-geometry-%02d", op), cases)
	}
}
func TestMapRoomsExclusions(t *testing.T) {
	var cases []legacy.PortTestMapRoomSpec
	for _, op := range []int{33, 34, 36} {
		for _, x := range []float32{-0.5001, -0.5, -0.4999, 0, 0.5, 1, 1.5, 2, 2.5, 130.10765} {
			s := roomSmoke(op)
			s.Records = append(s.Records, roomRecord(28))
			s.Records[9].Words = map[int]uint32{4: math.Float32bits(1), 8: math.Float32bits(1), 12: math.Float32bits(2), 16: math.Float32bits(2)}
			s.Records[1].Refs[368] = roomArg(10)
			s.Records[8].Words[4] = math.Float32bits(x)
			s.Records[8].Words[8] = math.Float32bits(x)
			s.Records[8].Words[12] = math.Float32bits(x + 1)
			s.Records[8].Words[16] = math.Float32bits(x + 1)
			s.Records[3].Words[0] = math.Float32bits(x)
			s.Records[3].Words[4] = math.Float32bits(x)
			cases = append(cases, s)
		}
	}
	mapRoomsHash(t, "map-rooms-exclusion-boundaries", cases)
	s := roomBase()
	a := roomAction(31, roomArg(2), roomArg(4), roomFloat(2), roomFloat(3))
	a.Assign = 10
	b := a
	b.Assign = 11
	mark := roomAction(-1)
	mark.Writes = []legacy.PortTestMapRoomWrite{{Slot: 10, Words: map[int]uint32{0: 1}}}
	s.Actions = []legacy.PortTestMapRoomAction{a, b, mark, roomAction(32, roomArg(2)), roomAction(36, roomArg(2), roomArg(4)), roomAction(23, roomArg(2))}
	v := mapRoomsHash(t, "map-rooms-exclusion-sequence", []legacy.PortTestMapRoomSpec{s})[0]
	last := v.Steps[3]
	if roomSlot(t, last, 2)[92] != last.Slots[11] || roomSlot(t, last, 11)[6] != 0 {
		t.Fatal("selective exclusion removal links")
	}
	if v.Steps[4].Return[0] != 1 {
		t.Fatal("remaining exclusion does not contain origin")
	}
}
func TestMapRoomsDecoration(t *testing.T) {
	var cases []legacy.PortTestMapRoomSpec
	for flags := uint32(0); flags < 64; flags++ {
		for _, roomFlag := range []uint32{0, 1, 2, 4, 8, 16, 32} {
			s := roomSmoke(51)
			s.Records[4].Words[64] = flags
			s.Records[1].Words[364] = roomFlag
			cases = append(cases, s)
		}
	}
	mapRoomsHash(t, "map-rooms-decoration-flags", cases)
	cases = nil
	for _, op := range []int{49, 50, 52, 53, 54} {
		for _, limit := range []uint32{0, 1, 2, 255} {
			for mode := 0; mode < 5; mode++ {
				s := roomSmoke(op)
				s.Records[4].Words[64] = limit << 8
				switch mode {
				case 1:
					s.Records[4].Words[76] = 99
				case 2:
					s.Records[4].Words[80] = 1
				case 3:
					s.Records[4].Words[64] |= 1 << 24
				case 4:
					s.Records[4].Words[64] |= 1 << 16
				}
				if (op == 49 || op == 50) && mode == 4 { // Zero total weight is the terminating unavailable-decoration branch.
					if op == 49 {
						s.Records[0].Words[192] = 0
					} else {
						s.Records[6].Words[8] = 0
					}
				}
				cases = append(cases, s)
			}
		}
	}
	mapRoomsHash(t, "map-rooms-decoration-boundaries", cases)
	s := roomSmoke(54)
	s.Records[0].Refs[120] = roomArg(6)
	s.Records[4].Words[64] = 1<<24 | 2<<8
	s.Records[5].Words[64] = 1<<24 | 2<<8
	s.Actions = []legacy.PortTestMapRoomAction{roomAction(13, roomArg(2)), roomAction(13, roomArg(3)), roomAction(54, roomArg(1))}
	v := mapRoomsHash(t, "map-rooms-required-decorations", []legacy.PortTestMapRoomSpec{s})[0].Steps[2]
	if v.Return[0] != 1 || roomSlot(t, v, 2)[93] != v.Slots[5] || roomSlot(t, v, 3)[93] != v.Slots[6] {
		t.Fatal("required decoration assignment")
	}
}
func TestMapRoomsRandom(t *testing.T) {
	var cases []legacy.PortTestMapRoomSpec
	for _, seed := range []uint32{0, 1, 2, 17, 12345, 0xffffffff} {
		for _, span := range []int32{0, 1, 2, 10, 100} {
			s := roomBase()
			s.NoGrid = true
			s.Seed = seed
			for i := 0; i < 8; i++ {
				s.Actions = append(s.Actions, roomAction(57, roomValue(-span), roomValue(span)), roomAction(58, roomValue(30), roomValue(span)), roomAction(59, roomFloat(-float32(span)), roomFloat(float32(span))))
			}
			cases = append(cases, s)
		}
	}
	r := mapRoomsHash(t, "map-rooms-random", cases)
	for i, v := range r {
		span := []int32{0, 1, 2, 10, 100}[i%5]
		for j, s := range v.Steps {
			switch j % 3 {
			case 0:
				x := int32(s.Return[0])
				if x < -span || x > span {
					t.Fatal("uniform random bounds")
				}
			case 1:
				x := int32(s.Return[0])
				if x < 30-span || x > 30+span {
					t.Fatal("centered random bounds")
				}
			case 2:
				x := math.Float64frombits(uint64(s.Return[0]) | uint64(s.Return[1])<<32)
				if x < float64(-span) || x > float64(span) {
					t.Fatal("floating random bounds")
				}
			}
		}
	}
}
func TestMapRoomsPointCache(t *testing.T) {
	var cases []legacy.PortTestMapRoomSpec
	for _, delta := range []float32{0, 0.09999, 0.1, 0.10001, 1} {
		s := roomBase()
		s.Actions = []legacy.PortTestMapRoomAction{roomAction(37, roomArg(2), roomArg(4)), roomAction(37, roomArg(2), roomArg(4))}
		s.Actions[1].Writes = []legacy.PortTestMapRoomWrite{{Slot: 4, Words: map[int]uint32{0: math.Float32bits(delta), 4: math.Float32bits(delta)}}}
		cases = append(cases, s)
	}
	r := mapRoomsHash(t, "map-rooms-point-cache-epsilon", cases)
	for i, v := range r {
		count := roomSlot(t, v.Steps[1], 2)[88] & 255
		want := uint32(2)
		if i < 2 {
			want = 1
		}
		if count != want {
			t.Fatalf("point cache case %d count %d want %d", i, count, want)
		}
	}
	s := roomBase()
	for i := 0; i < 18; i++ {
		a := roomAction(37, roomArg(2), roomArg(4))
		a.Writes = []legacy.PortTestMapRoomWrite{{Slot: 4, Words: map[int]uint32{0: math.Float32bits(float32(i)), 4: math.Float32bits(float32(i))}}}
		s.Actions = append(s.Actions, a)
	}
	v := mapRoomsHash(t, "map-rooms-point-cache-capacity", []legacy.PortTestMapRoomSpec{s})[0]
	if roomSlot(t, v.Steps[17], 2)[88]&255 != 16 {
		t.Fatal("point cache exceeds capacity")
	}
}
func TestMapRoomsRandomVariety(t *testing.T) {
	s := roomBase()
	s.NoGrid = true
	s.Seed = 42
	for i := 0; i < 128; i++ {
		s.Actions = append(s.Actions, roomAction(57, roomValue(0), roomValue(3)))
	}
	r := mapRoomsHash(t, "map-rooms-random-variety", []legacy.PortTestMapRoomSpec{s})[0]
	seen := map[uint32]bool{}
	for _, step := range r.Steps {
		seen[step.Return[0]] = true
	}
	if len(seen) != 4 {
		t.Fatalf("seeded random missed bins: %v", seen)
	}
}
func TestMapRoomsCorridorTrimming(t *testing.T) {
	var cases []legacy.PortTestMapRoomSpec
	for kind := int32(2); kind <= 5; kind++ {
		for _, width := range []int32{1, 2, 5, 8} {
			for _, x := range []float32{-65.053826, 0, 32.526913, 130.10765} {
				s := roomSmoke(43)
				roomSetGeometry(&s.Records[1], kind, width, width, x, x)
				roomSetGeometry(&s.Records[2], 1, 4, 4, 0, 0)
				cases = append(cases, s)
			}
		}
	}
	mapRoomsHash(t, "map-rooms-corridor-trimming", cases)
}
func TestMapRoomsGeometryLimits(t *testing.T) {
	for _, op := range []int{15, 16, 19, 26, 33, 38, 39, 40, 41, 42, 43, 44, 45, 46, 52, 55} {
		var cases []legacy.PortTestMapRoomSpec
		for kind := int32(1); kind <= 5; kind++ {
			for _, n := range []int32{-1, 0, 1, 2, 8} {
				s := roomSmoke(op)
				roomSetGeometry(&s.Records[1], kind, n, n, 0, 0)
				if op == 39 || op == 41 || op == 42 {
					s.Actions[0].Args[0] = roomValue(kind)
				}
				if op == 40 {
					s.Actions[0].Args[0] = roomValue((kind - 1) % 4)
				}
				if op == 55 {
					s.Actions[0].Args[1] = roomFloat(float32(n) / 10)
				}
				cases = append(cases, s)
			}
		}
		mapRoomsHash(t, fmt.Sprintf("map-rooms-limits-%02d", op), cases)
	}
}

func TestMapRoomsConstructorBoundaries(t *testing.T) {
	var cases []legacy.PortTestMapRoomSpec
	for _, seed := range []uint32{0, 1, 19} {
		for _, mean := range []uint32{1, 4, 5, 8, 20} {
			for _, spread := range []uint32{0, 1, 2} {
				for _, flag := range []uint32{0, 1} {
					s := roomSmoke(22)
					s.NoGrid = true
					s.Seed = seed
					s.Records[0].Words[32] = mean
					s.Records[0].Words[36] = spread
					s.Records[0].Words[56] = flag
					cases = append(cases, s)
				}
			}
		}
	}
	r := mapRoomsHash(t, "map-rooms-prepare-boundaries", cases)
	for _, v := range r {
		w := roomSlot(t, v.Steps[0], 10)
		if w[0] != 1 || int32(w[3]) < 5 || int32(w[4]) < 5 {
			t.Fatal("prepared room minimum dimensions")
		}
	}
	cases = nil
	for kind := int32(2); kind <= 5; kind++ {
		for _, width := range []int32{-1, 0, 1, 8} {
			for _, spread := range []uint32{0, 1, 2} {
				s := roomSmoke(48)
				s.NoGrid = true
				s.Records[0].Words[4] = 5
				s.Records[0].Words[8] = spread
				s.Actions[0].Args[1] = roomValue(kind)
				s.Actions[0].Args[2] = roomValue(width)
				cases = append(cases, s)
			}
		}
	}
	mapRoomsHash(t, "map-rooms-hall-boundaries", cases)
}
func TestMapRoomsResolveOverlap(t *testing.T) {
	var cases []legacy.PortTestMapRoomSpec
	for _, x := range []float32{0, 32.526913, 65.053826, 97.580739} {
		for _, width := range []int32{1, 2, 4} {
			s := roomBase()
			s.Radius = 8
			roomSetGeometry(&s.Records[1], 1, width, 2, x, 32.526913)
			roomSetGeometry(&s.Records[2], 1, 4, 4, 0, 0)
			s.Actions = []legacy.PortTestMapRoomAction{roomAction(4, roomArg(3)), roomAction(8, roomArg(2), roomArg(3)), roomAction(6, roomArg(2))}
			cases = append(cases, s)
		}
	}
	r := mapRoomsHash(t, "map-rooms-resolve-overlap", cases)
	for i, v := range r {
		if v.Steps[1].Return[0] != 1 || v.Steps[2].Return[0] != 0 {
			t.Fatalf("single obstacle case %d not resolved", i)
		}
	}
	s := roomBase()
	s.Radius = 16
	s.Records = append(s.Records, roomRecord(376))
	roomSetGeometry(&s.Records[1], 1, 2, 2, 97.580739, 130.10765)
	roomSetGeometry(&s.Records[2], 1, 4, 10, 0, 0)
	roomSetGeometry(&s.Records[9], 1, 4, 10, 130.10765, 0)
	s.Actions = []legacy.PortTestMapRoomAction{roomAction(4, roomArg(3)), roomAction(4, roomArg(10)), roomAction(8, roomArg(2), roomArg(3)), roomAction(6, roomArg(2))}
	v := mapRoomsHash(t, "map-rooms-resolve-iteration-limit", []legacy.PortTestMapRoomSpec{s})[0]
	if v.Steps[2].Return[0] != 0 || v.Steps[3].Return[0] == 0 {
		t.Fatal("oscillating overlap must hit the bounded iteration limit")
	}
}
