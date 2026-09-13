//go:build porttest

package opennox

import (
	"fmt"
	"math"
	"testing"

	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/legacy"
)

func paintRecord(t *testing.T, step legacy.PortTestPaintStep, slot int) []uint32 {
	t.Helper()
	p := step.Slots[slot]
	for _, r := range step.Records {
		if r.Alive && p >= r.ID && p-r.ID < uint32(len(r.Words)*4) {
			return r.Words[(p-r.ID)/4:]
		}
	}
	t.Fatalf("missing record for slot %d identity %d", slot, p)
	return nil
}
func paintPos(s *legacy.PortTestPaintSpec, slot int, x, y float32) {
	s.Records[slot-1].Words[0] = math.Float32bits(x)
	s.Records[slot-1].Words[4] = math.Float32bits(y)
}
func paintGlobal(name string, value int32) map[string]legacy.PortTestMapRoomArg {
	return map[string]legacy.PortTestMapRoomArg{name: roomValue(value)}
}
func TestMapPaintingCoordinates(t *testing.T) {
	values := []float32{-1e8, -8192, -4096, -65.053826, -0.01, 0, 0.01, 32.526913, 130.10765, 4096, 8192, 1e8}
	var cases []legacy.PortTestPaintSpec
	for _, x := range values {
		for _, y := range values {
			for _, alias := range []bool{false, true} {
				s := paintSmoke(4)
				paintPos(&s, 4, x, y)
				if alias {
					s.Actions[0].Args[1] = roomArg(4)
				}
				cases = append(cases, s)
			}
		}
	}
	for _, which := range []int{0, 1, 2} {
		s := paintSmoke(4)
		if which != 1 {
			s.Actions[0].Args[0] = roomValue(0)
		}
		if which != 0 {
			s.Actions[0].Args[1] = roomValue(0)
		}
		cases = append(cases, s)
	}
	got := paintingHash(t, "map-painting-coordinates", cases)
	for i, r := range got {
		if i >= len(got)-3 {
			if r.Steps[0].Return != 0 {
				t.Fatal("nil coordinate admission")
			}
			continue
		}
		slot := 8
		if i%2 == 1 {
			slot = 4
		}
		w := paintRecord(t, r.Steps[0], slot)
		x, y := math.Float32frombits(w[0]), math.Float32frombits(w[1])
		if x <= 80.5 || x >= 5853.5 || y <= 80.5 || y >= 5853.5 {
			t.Fatalf("unclamped coordinate %v %v", x, y)
		}
	}
	s := paintSmoke(4)
	r := paintingHash(t, "map-painting-coordinate-origin", []legacy.PortTestPaintSpec{s})[0]
	w := paintRecord(t, r.Steps[0], 8)
	if w[0] != math.Float32bits(2957) || w[1] != math.Float32bits(2956) {
		t.Fatal("origin transform")
	}
}
func TestMapPaintingDirectionTables(t *testing.T) {
	var cases []legacy.PortTestPaintSpec
	for a := 0; a < 13; a++ {
		for b := 0; b < 13; b++ {
			s := paintSmoke(3)
			s.Actions[0] = paintAction(3, roomValue(int32(a)), roomValue(int32(b)))
			cases = append(cases, s)
		}
	}
	r := paintingHash(t, "map-painting-direction-union", cases)
	for i := 0; i < 13; i++ {
		if r[i*13+i].Steps[0].Return != uint32(i) {
			t.Fatalf("direction %d self composition", i)
		}
	}
	if r[1].Steps[0].Return != 2 || r[13].Steps[0].Return != 2 {
		t.Fatal("crossing direction composition")
	}
	cases = nil
	for _, op := range []int{5, 30, 31, 32} {
		for v := -2; v <= 17; v++ {
			s := paintSmoke(op)
			s.Actions[0].Args[0] = roomValue(int32(v))
			cases = append(cases, s)
		}
	}
	paintingHash(t, "map-painting-direction-setters", cases)
}
func TestMapPaintingSubtilePool(t *testing.T) {
	var cases []legacy.PortTestPaintSpec
	for _, n := range []int{1, 9, 10, 11, 20, 31} {
		s := paintBase()
		for i := 0; i < n; i++ {
			a := paintAction(0, roomValue(int32(i+1)), roomValue(int32(i*3)), roomValue(int32(i%2)), roomValue(int32(i%8)))
			a.Assign = 50 + i
			s.Actions = append(s.Actions, a)
		}
		for i := 0; i < n; i += 2 {
			s.Actions = append(s.Actions, paintAction(1, roomArg(50+i)))
		}
		for i := 0; i < (n+1)/2; i++ {
			a := paintAction(0, roomValue(99), roomValue(int32(i)), roomValue(1), roomValue(2))
			a.Assign = 100 + i
			s.Actions = append(s.Actions, a)
		}
		cases = append(cases, s)
	}
	out := paintingHash(t, "map-painting-subtile-pool", cases)
	for ci, n := range []int{1, 9, 10, 11, 20, 31} {
		r := out[ci]
		for i := 0; i < n; i++ {
			w := paintRecord(t, r.Steps[i], 50+i)
			if w[0] != uint32(i+1) || w[1] != uint32(i*3) || w[4] != 0 {
				t.Fatal("subtile allocation initialization")
			}
		}
		firstReuse := n + (n+1)/2
		want := 50 + (n-1)/2*2
		if r.Steps[firstReuse].Return != r.Steps[n-1].Slots[want] {
			t.Fatalf("pool %d is not LIFO", n)
		}
	}
	s := paintBase()
	for i := 0; i < 3; i++ {
		a := paintAction(0, roomValue(int32(i+1)), roomValue(0), roomValue(0), roomValue(0))
		a.Assign = 50 + i
		s.Actions = append(s.Actions, a)
	}
	a := paintAction(2, roomArg(9))
	a.Writes = []legacy.PortTestMapRoomWrite{{Slot: 9, Refs: map[int]legacy.PortTestMapRoomArg{16: roomArg(50)}}, {Slot: 50, Refs: map[int]legacy.PortTestMapRoomArg{16: roomArg(51)}}, {Slot: 51, Refs: map[int]legacy.PortTestMapRoomArg{16: roomArg(52)}}}
	s.Actions = append(s.Actions, a)
	r := paintingHash(t, "map-painting-subtile-chain-free", []legacy.PortTestPaintSpec{s})[0]
	last := r.Steps[3]
	if paintRecord(t, last, 9)[4] != 0 || last.Globals["dword_5d4594_588084"] != last.Slots[52] {
		t.Fatal("chain recycling")
	}
}
func TestMapPaintingTileRouting(t *testing.T) {
	var cases []legacy.PortTestPaintSpec
	coords := []int{-1, 0, 1, 2, 64, 126, 127}
	sub := [][2]int{{0, 0}, {0, 45}, {45, 0}, {45, 45}, {23, 23}, {22, 23}, {23, 22}, {22, 24}, {24, 22}}
	for _, x := range coords {
		for _, y := range coords {
			for _, p := range sub {
				s := paintSmoke(8)
				s.Actions[0] = paintAction(8, roomValue(int32(x)), roomValue(int32(y)), roomValue(int32(p[0])), roomValue(int32(p[1])), roomArg(13))
				cases = append(cases, s)
			}
		}
	}
	out := paintingHash(t, "map-painting-tile-routing", cases)
	for i, s := range cases {
		x, y := int32(s.Actions[0].Args[0].Value), int32(s.Actions[0].Args[1].Value)
		if x <= 1 || y <= 1 || x >= 127 || y >= 127 {
			if out[i].Steps[0].Return != 0 || len(out[i].Steps[0].Cells) != 0 {
				t.Fatal("tile bounds changed grid")
			}
		} else if len(out[i].Steps[0].Cells) != 1 {
			t.Fatal("routing must touch one cell")
		}
	}
	cases = nil
	for _, side := range []int{0, 1, 2, 3} {
		for _, anchor := range []int{0, 1} {
			for _, stored := range []int{0, 1, 2} {
				for _, fixed := range []int{0, 1} {
					for _, x := range []int{2, 64, 126} {
						s := paintSmoke(9)
						s.Actions[0].Args[0] = roomValue(int32(x))
						s.Actions[0].Args[3] = roomValue(int32(side))
						s.Records[12].Words[8] = uint32(anchor)
						s.Records[12].Words[12] = 66
						s.Records[12].Words[16] = 66
						s.Records[12].Words[20] = uint32(stored)
						s.Globals["tileFlag"] = roomValue(int32(fixed))
						s.Globals["dword_5d4594_3835348"] = roomValue(7)
						cases = append(cases, s)
					}
				}
			}
		}
	}
	paintingHash(t, "map-painting-tile-anchors", cases)
	cases = nil
	values := []float32{-46, -1, 0, 46, 80.5, 92, 2932.5, 2956, 5853.5, 6000}
	for _, x := range values {
		for _, y := range values {
			s := paintSmoke(7)
			paintPos(&s, 11, x, y)
			cases = append(cases, s)
		}
	}
	paintingHash(t, "map-painting-world-tile-coordinates", cases)
}
func TestMapPaintingFloorRectangles(t *testing.T) {
	var cases []legacy.PortTestPaintSpec
	for _, x := range []float32{-65.053826, 0, 32.526913} {
		for _, w := range []int32{-1, 0, 1, 2, 4, 8} {
			for _, h := range []int32{-1, 0, 1, 2, 4, 8} {
				s := paintSmoke(13)
				paintPos(&s, 4, x, x)
				s.Actions[0].Args[2] = roomValue(w)
				s.Actions[0].Args[3] = roomValue(h)
				cases = append(cases, s)
			}
		}
	}
	out := paintingHash(t, "map-painting-floor-rectangles", cases)
	for i, s := range cases {
		w, h := int32(s.Actions[0].Args[2].Value), int32(s.Actions[0].Args[3].Value)
		if w <= 0 || h <= 0 {
			if len(out[i].Steps[0].Cells) != 0 {
				t.Fatal("empty rectangle paints")
			}
		} else if len(out[i].Steps[0].Cells) == 0 {
			t.Fatal("positive rectangle does not paint")
		}
	}
	cases = nil
	for _, op := range []int{6, 14, 21, 22} {
		for _, n := range []int32{-1, 0, 1, 2, 4, 8} {
			for _, x := range []float32{-1e6, -32.526913, 0, 0.01, 65.053826, 1e6} {
				s := paintSmoke(op)
				paintPos(&s, 4, x, x)
				if op == 14 {
					s.Actions[0].Args[2] = roomValue(n)
				} else if op == 21 || op == 22 {
					s.Actions[0].Args[1] = roomValue(n)
				}
				cases = append(cases, s)
			}
		}
	}
	paintingHash(t, "map-painting-floor-lines", cases)
}
func TestMapPaintingWalls(t *testing.T) {
	var cases []legacy.PortTestPaintSpec
	for _, existing := range []bool{false, true} {
		for _, merge := range []int{0, 1} {
			for _, cycle := range []int{0, 1} {
				for _, dir := range []int{0, 1, 2, 3, 7, 8, 12} {
					for _, variation := range []int{0, 3, 4, 255} {
						s := paintSmoke(34)
						s.Globals["wallDir"] = roomValue(int32(dir))
						s.Globals["wallVariation"] = roomValue(int32(variation))
						s.Globals["dword_5d4594_3835368"] = roomValue(int32(merge))
						s.Globals["dword_5d4594_3835372"] = roomValue(int32(cycle))
						if existing {
							s.Walls = []legacy.PortTestPaintWall{{X: 128, Y: 128, Words: map[int]uint32{0: 1}}}
						}
						cases = append(cases, s)
					}
				}
			}
		}
	}
	out := paintingHash(t, "map-painting-wall-updates", cases)
	for _, r := range out {
		if r.Steps[0].Return != 1 || len(r.Steps[0].Walls.ByPos) != 1 {
			t.Fatal("wall creation/update")
		}
	}
	cases = nil
	for _, op := range []int{11, 12, 33, 34, 35} {
		for _, x := range []float32{-1e6, -65.053826, -0.01, 0, 0.01, 65.053826, 1e6} {
			for _, n := range []int32{-1, 0, 1, 4} {
				s := paintSmoke(op)
				paintPos(&s, 4, x, x)
				if op == 11 || op == 12 {
					s.Actions[0].Args[1] = roomValue(n)
				}
				cases = append(cases, s)
			}
		}
	}
	paintingHash(t, "map-painting-wall-coordinates", cases)
	s := paintSmoke(34)
	s.Actions = append(s.Actions, paintAction(33, roomArg(4), roomArg(8)), paintAction(35, roomArg(4)), paintAction(33, roomArg(4), roomArg(8)))
	r := paintingHash(t, "map-painting-wall-lifecycle", []legacy.PortTestPaintSpec{s})[0]
	if r.Steps[1].Return != 1 || r.Steps[2].Return != 1 || r.Steps[3].Return != 0 || len(r.Steps[2].Walls.ByPos) != 0 {
		t.Fatal("wall lookup/delete lifecycle")
	}
	s = paintSmoke(35)
	s.Walls[0].Data = roomArg(14)
	s.Globals["secretWalls"] = roomArg(14)
	r = paintingHash(t, "map-painting-secret-wall-removal", []legacy.PortTestPaintSpec{s})[0]
	if r.Steps[0].Globals["secretWalls"] != 0 {
		t.Fatal("secret wall head retained")
	}
	for _, record := range r.Steps[0].Records {
		if record.ID == r.Steps[0].Slots[14] && record.Alive {
			t.Fatal("secret record still alive")
		}
	}
}
func TestMapPaintingRooms(t *testing.T) {
	for _, op := range []int{15, 19, 20, 24, 25} {
		var cases []legacy.PortTestPaintSpec
		for kind := int32(0); kind <= 6; kind++ {
			for _, size := range [][2]int32{{1, 1}, {2, 3}, {4, 4}, {8, 6}} {
				for dir := int32(0); dir < 4; dir++ {
					s := paintSmoke(op)
					roomSetGeometry(&s.Records[1], kind, size[0], size[1], 0, 0)
					if op == 20 {
						s.Actions[0].Args[2] = roomValue(dir)
					}
					if op == 25 {
						s.Actions[0].Args[3] = roomValue(dir)
					}
					if op == 19 || op == 24 {
						s.Records[1].Words[216] = 1 << uint(dir*8)
						s.Records[1].Refs[88+int(dir)*32] = roomArg(3)
					}
					cases = append(cases, s)
				}
			}
		}
		paintingHash(t, fmt.Sprintf("map-painting-rooms-%02d", op), cases)
	}
	for _, op := range []int{16, 17, 18} {
		var cases []legacy.PortTestPaintSpec
		for w := int32(1); w <= 8; w++ {
			for h := int32(1); h <= 8; h++ {
				s := paintSmoke(op)
				s.Records[7].Words[0] = uint32(w)
				s.Records[7].Words[4] = uint32(h)
				cases = append(cases, s)
			}
		}
		paintingHash(t, fmt.Sprintf("map-painting-patterns-%02d", op), cases)
	}
	s := paintBase()
	roomSetGeometry(&s.Records[2], 1, 4, 4, 130.10765, 0)
	s.Records[1].Words[216] = 1 << 16
	s.Records[1].Refs[152] = roomArg(3)
	s.Records[2].Words[216] = 1 << 24
	s.Records[2].Refs[184] = roomArg(2)
	s.Actions = []legacy.PortTestPaintAction{paintAction(19, roomArg(1), roomArg(2)), paintAction(19, roomArg(1), roomArg(3)), paintAction(24, roomArg(1), roomArg(2)), paintAction(24, roomArg(1), roomArg(3))}
	paintingHash(t, "map-painting-connected-rooms", []legacy.PortTestPaintSpec{s})
}
func TestMapPaintingObjectPlacement(t *testing.T) {
	var cases []legacy.PortTestPaintSpec
	for _, op := range []int{26, 27, 28, 29, 39, 40} {
		for _, x := range []float32{-4096, -65.053826, -0.01, 0, 0.01, 65.053826, 4096} {
			for _, n := range []int32{2, 3, 4, 7} {
				s := paintSmoke(op)
				paintPos(&s, 4, x, x)
				if op >= 26 && op <= 29 {
					s.Actions[0].Args[2] = roomValue(n)
				}
				cases = append(cases, s)
			}
		}
	}
	out := paintingHash(t, "map-painting-object-placement", cases)
	for i, r := range out {
		step := r.Steps[0]
		want := 1
		if op := cases[i].Actions[0].Op; op == 27 || op == 29 {
			want = 2
		}
		count := 0
		for _, rec := range step.Records {
			if rec.Kind != "object" {
				continue
			}
			count++
			w := rec.Words
			if w[4]&uint32(object.FlagActive) == 0 || w[4]&uint32(object.FlagPending) != 0 {
				t.Fatal("object was not promoted from pending")
			}
			for _, off := range []int{14, 15} {
				if math.Mod(float64(math.Float32frombits(w[off])), 23) != 0 {
					t.Fatal("door position is not snapped")
				}
			}
		}
		if count != want || step.Objects[6] != 0 || step.Objects[5] == 0 {
			t.Fatalf("case %d: objects %d want %d", i, count, want)
		}
	}
	cases = nil
	for _, typ := range []int{1, 2, 3, 4} {
		for dir := int32(-1); dir <= 9; dir++ {
			s := paintSmoke(41)
			s.Objects[0].Type = typ
			s.Actions[0].Args[1] = roomValue(dir)
			cases = append(cases, s)
		}
	}
	oriented := paintingHash(t, "map-painting-object-orientation", cases)
	angles := map[int]uint32{0: 128, 1: 160, 2: 192, 3: 96, 5: 224, 6: 64, 7: 32, 8: 0}
	doors := map[int]uint32{1: 0, 3: 24, 5: 8, 7: 16}
	for i, capture := range oriented {
		step := capture.Steps[0]
		typ, dir := i/11+1, i%11-1
		wantReturn := uint32(0)
		if typ == 4 {
			wantReturn = 1
			if paintRecord(t, step, 100)[31]&0xffff != angles[dir] || paintRecord(t, step, 1100)[94] != angles[dir] {
				t.Fatalf("monster direction %d", dir)
			}
		} else if typ == 2 {
			if angle, ok := doors[dir]; ok {
				wantReturn = 1
				data := paintRecord(t, step, 1100)
				if data[1] != angle || data[2] != angle || data[3] != angle {
					t.Fatalf("door direction %d", dir)
				}
			}
		}
		if step.Return != wantReturn {
			t.Fatalf("orientation type %d direction %d return %d", typ, dir, step.Return)
		}
	}
	cases = nil
	for _, typ := range []int{1, 2, 3} {
		for _, v := range []int32{0, 1, 127, 128, 255} {
			s := paintSmoke(42)
			s.Objects[0].Type = typ
			s.Actions[0].Args[1] = roomValue(v)
			cases = append(cases, s)
		}
	}
	paintingHash(t, "map-painting-spellbook-type-check", cases)
	cases = nil
	for _, name := range []string{"NONE", "none", "PaintDoor", "paintobject", "PaintBook", "missing", ""} {
		s := paintSmoke(38)
		paintString(&s.Records[9], 0, name)
		s.Actions = append(s.Actions, paintAction(39, roomArg(4)))
		cases = append(cases, s)
	}
	paintingHash(t, "map-painting-object-name-gates", cases)
}

func TestMapPaintingSubtileSequences(t *testing.T) {
	var cases []legacy.PortTestPaintSpec
	for _, merge := range []int32{0, 1} {
		for edge := int32(0); edge < 8; edge++ {
			for next := int32(0); next < 8; next++ {
				s := paintBase()
				add := func(tile, variation, border, e int32) legacy.PortTestPaintAction {
					return paintAction(46, roomArg(9), roomValue(tile), roomValue(variation), roomValue(border), roomValue(e), roomValue(merge))
				}
				s.Actions = []legacy.PortTestPaintAction{add(1, 0, 0, edge), add(1, 0, 0, edge), add(1, 0, 0, next), add(2, 1, 1, edge), add(255, 0, 0, 0), add(1, 0, 255, 0), add(1, 0, 0, next)}
				cases = append(cases, s)
			}
		}
	}
	out := paintingHash(t, "map-painting-subtile-sequences", cases)
	for i, c := range out {
		if c.Steps[0].Return != 1 || paintRecord(t, c.Steps[0], 9)[4] == 0 {
			t.Fatal("subtile not added")
		}
		if i < 64 && c.Steps[1].Return != 0 {
			t.Fatal("duplicate subtile admitted")
		}
		if c.Steps[4].Return != 1 || c.Steps[5].Return != 1 {
			t.Fatal("subtile removal failed")
		}
	}
	cases = nil
	for _, base := range []uint32{0, 1, 255} {
		for _, tile := range []int32{0, 1, 255} {
			for _, border := range []int32{0, 1, 255} {
				s := paintSmoke(46)
				s.Records[8].Words[0] = base
				s.Actions[0].Args[1], s.Actions[0].Args[3] = roomValue(tile), roomValue(border)
				cases = append(cases, s)
			}
		}
	}
	paintingHash(t, "map-painting-subtile-none", cases)
}

func TestMapPaintingBorderPropagation(t *testing.T) {
	for _, op := range []int{43, 44, 47} {
		var cases []legacy.PortTestPaintSpec
		for _, n := range []int32{1, 2, 4, 8} {
			for _, border := range []int32{0, 1, 255} {
				for _, mode := range []int32{0, 1} {
					s := paintSmoke(op)
					s.Records[13].Words[0], s.Records[13].Words[4] = 2957, 2956
					s.Records[12].Words[24] = uint32(border)
					a := s.Actions[0]
					a.Globals = map[string]legacy.PortTestMapRoomArg{"tile": roomValue(1), "dword_5d4594_3835356": roomValue(border), "dword_5d4594_3835352": roomValue(mode)}
					s.Actions = []legacy.PortTestPaintAction{paintAction(13, roomArg(1), roomArg(4), roomValue(n), roomValue(n)), a}
					cases = append(cases, s)
				}
			}
		}
		paintingHash(t, fmt.Sprintf("map-painting-border-propagation-%02d", op), cases)
	}
	var cases []legacy.PortTestPaintSpec
	for _, x := range []int32{0, 1, 2, 64, 126, 127} {
		for _, side := range []int32{1, 2, 3} {
			for edge := int32(0); edge < 8; edge++ {
				s := paintSmoke(45)
				s.Actions[0] = paintAction(45, roomValue(x), roomValue(x), roomValue(side), roomValue(0), roomArg(13), roomValue(edge))
				s.Globals["dword_5d4594_3835352"] = roomValue(1)
				cases = append(cases, s)
			}
		}
	}
	paintingHash(t, "map-painting-border-neighbors", cases)
}

func TestMapPaintingWallCorners(t *testing.T) {
	var cases []legacy.PortTestPaintSpec
	for mask := 0; mask < 16; mask++ {
		s := paintSmoke(23)
		for i, xy := range [][2]int{{127, 127}, {129, 129}, {127, 129}, {129, 127}} {
			if mask&(1<<i) != 0 {
				s.Walls = append(s.Walls, legacy.PortTestPaintWall{X: xy[0], Y: xy[1]})
			}
		}
		cases = append(cases, s)
	}
	out := paintingHash(t, "map-painting-wall-neighbor-masks", cases)
	directions := map[int]uint32{3: 1, 5: 10, 6: 9, 7: 6, 9: 7, 10: 8, 11: 4, 12: 0, 13: 3, 14: 5, 15: 2}
	for mask, c := range out {
		if dir, ok := directions[mask]; ok && c.Steps[0].Globals["wallDir"] != dir {
			t.Fatalf("neighbor mask %d direction %d", mask, c.Steps[0].Globals["wallDir"])
		}
	}
	cases = nil
	for mask := 0; mask < 16; mask++ {
		for _, flags := range []uint32{0, 4, 8, 128, 140} {
			for _, op := range []int{36, 37} {
				s := paintSmoke(op)
				s.Cells = []legacy.PortTestPaintCell{
					{X: 64, Y: 64, Words: map[int]uint32{4: 255, 24: 255}},
					{X: 64, Y: 63, Words: map[int]uint32{24: 255}},
					{X: 63, Y: 64, Words: map[int]uint32{4: 255}},
				}
				if mask&1 != 0 {
					s.Cells[2].Words[4] = 0
				}
				if mask&2 != 0 {
					s.Cells[0].Words[4] = 0
				}
				if mask&4 != 0 {
					s.Cells[1].Words[24] = 0
				}
				if mask&8 != 0 {
					s.Cells[0].Words[24] = 0
				}
				s.Walls = []legacy.PortTestPaintWall{{X: 128, Y: 128, Words: map[int]uint32{0: 0x03010c, 4: flags | 128<<8 | 128<<16 | 80<<24}}}
				cases = append(cases, s)
			}
		}
	}
	paintingHash(t, "map-painting-wall-floor-masks", cases)
}

func TestMapPaintingLayoutSelection(t *testing.T) {
	var cases []legacy.PortTestPaintSpec
	for _, count := range []uint32{1, 2, 3} {
		for seed := 0; seed < 32; seed++ {
			s := paintSmoke(10)
			s.Seed = seed
			s.Records[4].Words[88] = count
			s.Records[5].Refs[124] = roomArg(7)
			s.Records[6].Refs[124] = roomArg(1)
			cases = append(cases, s)
		}
	}
	out := paintingHash(t, "map-painting-layout-selection", cases)
	for i, c := range out {
		step := c.Steps[0]
		allowed := []int{6, 7, 1}[:i/32+1]
		found := false
		for _, slot := range allowed {
			found = found || step.Return == step.Slots[slot]
		}
		if !found {
			t.Fatal("selected outside layout list")
		}
	}
}

func TestMapPaintingConnectedPatterns(t *testing.T) {
	var cases []legacy.PortTestPaintSpec
	for kind := int32(1); kind <= 5; kind++ {
		for pattern := uint32(1); pattern <= 3; pattern++ {
			for _, size := range []int32{3, 6, 10} {
				s := paintSmoke(15)
				roomSetGeometry(&s.Records[1], kind, size, size, 0, 0)
				s.Records[5].Refs[120] = roomArg(7)
				s.Records[6].Words[0] = pattern
				cases = append(cases, s)
			}
		}
	}
	paintingHash(t, "map-painting-connected-patterns", cases)
}

func TestMapPaintingObjectGuardsAndRepeat(t *testing.T) {
	var cases []legacy.PortTestPaintSpec
	for _, op := range []int{40, 41, 42} {
		s := paintSmoke(op)
		s.Actions[0].Args[0] = roomValue(0)
		cases = append(cases, s)
	}
	out := paintingHash(t, "map-painting-null-objects", cases)
	for _, c := range out {
		if c.Steps[0].Return != 0 {
			t.Fatal("null object accepted")
		}
	}
	cases = nil
	for _, typ := range []int{1, 2, 3} {
		s := paintSmoke(40)
		s.Objects[0].Type = typ
		a := s.Actions[0]
		a.Writes = []legacy.PortTestMapRoomWrite{{Slot: 4, Words: map[int]uint32{0: math.Float32bits(65.053826), 4: math.Float32bits(32.526913)}}}
		s.Actions = append(s.Actions, a)
		cases = append(cases, s)
	}
	out = paintingHash(t, "map-painting-repeat-object-move", cases)
	for _, c := range out {
		a, b := paintRecord(t, c.Steps[0], 100), paintRecord(t, c.Steps[1], 100)
		if b[10] != a[10]+1 || b[14] == a[14] || b[16] != a[16] || b[17] != a[17] {
			t.Fatal("active-object move sequencing")
		}
	}
}
