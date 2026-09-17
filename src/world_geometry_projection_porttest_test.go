//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"math"
	"testing"
)

func TestWorldGeometryEdgeProjection(t *testing.T) {
	for _, tc := range []struct {
		point    [2]int32
		distance float32
		want     int32
		after    [2]int32
	}{
		{[2]int32{5, 1}, 1, 1, [2]int32{5, 0}},
		{[2]int32{5, 2}, 1, 0, [2]int32{5, 2}},
		{[2]int32{-1, -1}, 2, 0, [2]int32{-1, -1}},
	} {
		s := worldGeometrySpec("edge-project")
		s.Words[0] = uint32(tc.point[0])
		s.Words[1] = uint32(tc.point[1])
		copy(s.Words[4:], []uint32{0, 0, 10, 0})
		s.Floats[0] = math.Float32bits(tc.distance)
		r := worldGeometryCall(t, s)
		if r.Return != tc.want || int32(r.Words[0]) != tc.after[0] || int32(r.Words[1]) != tc.after[1] {
			t.Fatal("edge projection", tc, r)
		}
	}
	type row struct {
		Segment  [4]int32
		Point    [2]int32
		Distance uint32
		Alias    bool
		Result   legacy.PortTestWorldGeometryResult
	}
	var rows []row
	for _, seg := range [][4]int32{{0, 0, 10, 0}, {0, 0, 0, 10}, {0, 0, 10, 10}, {10, 10, 0, 0}, {3, 3, 3, 3}, {-10, -10, 10, 10}} {
		for _, pt := range [][2]int32{{0, 0}, {5, 1}, {5, 5}, {10, 10}, {11, 11}, {-1, -1}, {-2, 5}, {20, 5}} {
			for _, distance := range []float32{-1, 0, 0.5, 1, 2, 10, float32(math.Inf(1))} {
				for _, alias := range []bool{false, true} {
					s := worldGeometrySpec("edge-project")
					s.Words[0] = uint32(pt[0])
					s.Words[1] = uint32(pt[1])
					for i, v := range seg {
						s.Words[4+i] = uint32(v)
					}
					s.Floats[0] = math.Float32bits(distance)
					if alias {
						s.Offsets[0] = 4
					}
					rows = append(rows, row{seg, pt, s.Floats[0], alias, worldGeometryCall(t, s)})
				}
			}
		}
	}
	spellbookCapture(t, "world-geometry-edge-projection", rows, "915f386757318804927e4050b8edf96b39a7f71aa8b0bd1b5aac41c9e04e07ff")
}
func TestWorldGeometryWallBounds(t *testing.T) {
	diamond := [8]int32{230, 115, 115, 230, 345, 230, 230, 345}
	for _, tc := range []struct {
		x, y int32
		want int32
	}{{230, 230, 1}, {0, 0, 0}, {500, 500, 0}, {230, 115, 1}, {230, 345, 1}} {
		s := worldGeometrySpec("wall-point")
		s.Words[0] = uint32(tc.x)
		s.Words[1] = uint32(tc.y)
		for i, v := range diamond {
			s.Words[4+i] = uint32(v)
		}
		if got := worldGeometryCall(t, s).Return; got != tc.want {
			t.Fatal("wall diamond", tc, got)
		}
	}
	type row struct {
		Op     string
		Input  []uint32
		Alias  int
		Result legacy.PortTestWorldGeometryResult
	}
	var rows []row
	for _, shape := range [][8]int32{diamond, {0, -1, -2, 0, 6000, 0, 0, 6001}, {0, 30, 40, 0, 10, 0, 0, 20}, {0, -2147483648, -2147483648, 0, 2147483647, 0, 0, 2147483647}} {
		for _, alias := range []int{0, 2, 4, 12} {
			s := worldGeometrySpec("wall-bounds")
			for i, v := range shape {
				s.Words[i] = uint32(v)
			}
			s.Offsets[1] = alias
			rows = append(rows, row{s.Op, append([]uint32(nil), s.Words[:8]...), alias, worldGeometryCall(t, s)})
		}
	}
	for _, x := range []int32{0, 57, 58, 114, 115, 138, 229, 230, 345, 346, 5888, 5889} {
		for _, y := range []int32{0, 57, 58, 114, 115, 230, 345, 346, 5888, 5889} {
			s := worldGeometrySpec("wall-point")
			s.Words[0] = uint32(x)
			s.Words[1] = uint32(y)
			for i, v := range diamond {
				s.Words[4+i] = uint32(v)
			}
			rows = append(rows, row{s.Op, append([]uint32(nil), s.Words[:12]...), -1, worldGeometryCall(t, s)})
		}
	}
	spellbookCapture(t, "world-geometry-wall-bounds", rows, "560caea44ae7150061bde76fcc7eff18ee56e9cdbda29e0ab991ae0512046bca")
}
func TestWorldGeometryDiagonalProjection(t *testing.T) {
	type row struct {
		Op     string
		Input  []uint32
		Flags  [2]int32
		Range  [2]uint32
		Result legacy.PortTestWorldGeometryResult
	}
	var rows []row
	for _, op := range []string{"project-positive", "project-negative"} {
		s := worldGeometrySpec(op)
		copy(s.Words, worldGeometryWords(0, 0))
		copy(s.Words[4:], worldGeometryWords(-100, -100))
		if op == "project-negative" {
			s.Words[5] = math.Float32bits(100)
		}
		s.Floats[0] = math.Float32bits(2)
		s.Floats[1] = math.Float32bits(10)
		if r := worldGeometryCall(t, s); r.Return != 0 || r.Words[8] != s.Words[8] || r.Words[9] != s.Words[9] {
			t.Fatal("unclamped projection")
		}
		s.Ints[0] = 1
		r := worldGeometryCall(t, s)
		wantY := float32(2)
		if op == "project-negative" {
			wantY = 21
		}
		if r.Return != 1 || math.Float32frombits(r.Words[8]) != 2 || math.Float32frombits(r.Words[9]) != wantY {
			t.Fatal("clamped projection", op, r)
		}
		for _, x := range []float32{-30, -1, 0, 1, 5, 10, 23, 30} {
			for _, y := range []float32{-30, 0, 5, 23, 30} {
				for _, low := range []int32{0, 1} {
					for _, high := range []int32{0, 1} {
						for _, limits := range [][2]float32{{0, 23}, {2, 10}, {10, 2}, {5, 5}} {
							s := worldGeometrySpec(op)
							copy(s.Words, worldGeometryWords(0, 0))
							copy(s.Words[4:], worldGeometryWords(x, y))
							s.Ints[0] = low
							s.Ints[1] = high
							s.Floats[0] = math.Float32bits(limits[0])
							s.Floats[1] = math.Float32bits(limits[1])
							r := worldGeometryCall(t, s)
							rows = append(rows, row{op, append([]uint32(nil), s.Words[:6]...), [2]int32{low, high}, [2]uint32{s.Floats[0], s.Floats[1]}, r})
						}
					}
				}
			}
		}
	}
	spellbookCapture(t, "world-geometry-diagonal-projection", rows, "7e8c343b28a1832b20104c765b8853e3bb2dacf08f7b8b71bbf995039428dceb")
}
func TestWorldGeometryQuadrantContract(t *testing.T) {
	for _, tc := range []struct {
		x, y float32
		want int32
	}{{0, 20, 8}, {0, 0, 1}, {20, 20, 4}, {20, 0, 2}} {
		s := worldGeometrySpec("quadrant")
		copy(s.Words, worldGeometryWords(tc.x, tc.y))
		copy(s.Words[4:], worldGeometryWords(0, 0))
		if got := worldGeometryCall(t, s).Return; got != tc.want {
			t.Errorf("quadrant (%g,%g): %d want %d", tc.x, tc.y, got, tc.want)
		}
	}
	type row struct {
		X, Y   uint32
		Result int32
	}
	var rows []row
	edge := float32(16.263456)
	for _, x := range []float32{-100, 0, math.Nextafter32(edge, float32(math.Inf(-1))), edge, math.Nextafter32(edge, float32(math.Inf(1))), 100} {
		for _, y := range []float32{-1, float32(math.Copysign(0, -1)), 0, 1} {
			s := worldGeometrySpec("quadrant")
			copy(s.Words, worldGeometryWords(x, y))
			copy(s.Words[4:], worldGeometryWords(0, 0))
			r := worldGeometryCall(t, s)
			var want int32
			if float64(x) <= 16.263456 {
				want = 8
				if y <= 0 {
					want = 1
				}
			} else {
				want = 4
				if y <= 0 {
					want = 2
				}
			}
			if r.Return != want {
				t.Fatal("quadrant threshold", x, y, r.Return, want)
			}
			rows = append(rows, row{math.Float32bits(x), math.Float32bits(y), r.Return})
		}
	}
	spellbookCapture(t, "world-geometry-quadrant", rows, "b59b359161c396159327ba59b9bb2505c831f32a93b2609693084e1831534b91")

}
func TestWorldGeometryRectangleCrossings(t *testing.T) {
	s := worldGeometrySpec("rectangle-crossings")
	copy(s.Words, worldGeometryWords(-5, 5, 15, 5))
	copy(s.Words[4:], worldGeometryWords(0, 0, 10, 10))
	s.Ints[0] = 2
	r := worldGeometryCall(t, s)
	if r.Return != 2 || math.Float32frombits(r.Words[8]) != 10 || math.Float32frombits(r.Words[9]) != 5 || math.Float32frombits(r.Words[10]) != 0 || math.Float32frombits(r.Words[11]) != 5 {
		t.Fatal("crossing count/order", r)
	}
	s = worldGeometrySpec("clipped-center")
	copy(s.Words, worldGeometryWords(-5, 5, 15, 5))
	copy(s.Words[4:], worldGeometryWords(-5, 5, 15, 5))
	copy(s.Words[8:], worldGeometryWords(0, 0, 10, 10))
	r = worldGeometryCall(t, s)
	if r.Return != 1 || math.Float32frombits(r.Words[16]) != 5 || math.Float32frombits(r.Words[17]) != 5 {
		t.Fatal("clipped midpoint", r)
	}
	type row struct {
		Op            string
		Line          [4]float32
		Count, Extend int32
		Result        legacy.PortTestWorldGeometryResult
	}
	var rows []row
	for _, line := range [][4]float32{{-5, 5, 15, 5}, {15, 5, -5, 5}, {5, -5, 5, 15}, {2, 2, 4, 4}, {0, 0, 10, 10}, {-1, -1, 0, 0}, {5, 5, 5, 5}, {20, 20, 30, 30}, {-5, 0, 15, 0}} {
		for _, op := range []string{"rectangle-crossings", "horizontal-crossing", "vertical-crossing", "clipped-center"} {
			for _, count := range []int32{0, 1, 2, 4} {
				for _, extend := range []int32{0, 1} {
					s := worldGeometrySpec(op)
					copy(s.Words, worldGeometryWords(line[:]...))
					copy(s.Words[4:], worldGeometryWords(0, 0, 10, 10))
					s.Ints[0] = count
					s.Ints[1] = extend
					if op == "horizontal-crossing" || op == "vertical-crossing" {
						s.Ints[0] = extend
						s.Floats = [3]uint32{math.Float32bits(0), math.Float32bits(10), math.Float32bits(5)}
					}
					if op == "clipped-center" {
						copy(s.Words[4:], worldGeometryWords(min(line[0], line[2]), min(line[1], line[3]), max(line[0], line[2]), max(line[1], line[3])))
						copy(s.Words[8:], worldGeometryWords(0, 0, 10, 10))
					}
					rows = append(rows, row{op, line, count, extend, worldGeometryCall(t, s)})
				}
			}
		}
	}
	spellbookCapture(t, "world-geometry-rectangle-crossings", rows, "f3d15dbfd32afcce662b0c10b36952bf324e98d0834aec5dbfc6088ebdd60f43")
}
