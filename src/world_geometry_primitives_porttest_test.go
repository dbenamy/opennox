//go:build porttest

package opennox

import (
	"encoding/binary"
	"github.com/opennox/opennox/v1/common/memmap/nox/blobdata"
	"github.com/opennox/opennox/v1/legacy"
	"math"
	"testing"
)

func worldGeometryTables(t *testing.T) map[string]*uint32 {
	t.Helper()
	for off, data := range blobdata.PortTestWorldGeometryTables() {
		copy(serverConfigOwnBytes(t, 0x587000, off, len(data)), data)
	}
	words, restore := legacy.PortTestWorldGeometryGlobals()
	t.Cleanup(restore)
	*words["directionThreshold"] = 6
	*words["objectForce"] = math.Float32bits(40)
	*words["wallForce"] = math.Float32bits(100)
	*words["gameBall"] = 0
	return words
}
func worldGeometryWords(values ...float32) []uint32 {
	out := make([]uint32, len(values))
	for i, v := range values {
		out[i] = math.Float32bits(v)
	}
	return out
}
func worldGeometrySpec(op string) legacy.PortTestWorldGeometrySpec {
	words := make([]uint32, 24)
	for i := range words {
		words[i] = 0x3f000000 + uint32(i)*0x10000
	}
	return legacy.PortTestWorldGeometrySpec{Op: op, Words: words, Offsets: [4]int{0, 4, 8, 16}}
}
func worldGeometryCall(t *testing.T, s legacy.PortTestWorldGeometrySpec) legacy.PortTestWorldGeometryResult {
	t.Helper()
	r := legacy.PortTestWorldGeometry(s)
	if !r.GuardsOK {
		t.Fatal("geometry buffer guards", s.Op)
	}
	return r
}

func TestWorldGeometryIntegerSegments(t *testing.T) {
	for _, tc := range []struct {
		a, b [4]int32
		want int32
	}{
		{[4]int32{0, 0, 10, 10}, [4]int32{0, 10, 10, 0}, 1},
		{[4]int32{0, 0, 10, 0}, [4]int32{0, 1, 10, 1}, 0},
		{[4]int32{0, 0, 10, 10}, [4]int32{10, 10, 20, 20}, 1},
		{[4]int32{0, 0, 0, 0}, [4]int32{1, 1, 2, 2}, 0},
	} {
		s := worldGeometrySpec("segments")
		for i := 0; i < 4; i++ {
			s.Words[i] = uint32(tc.a[i])
			s.Words[4+i] = uint32(tc.b[i])
		}
		if got := worldGeometryCall(t, s).Return; got != tc.want {
			t.Fatal(tc, got)
		}
	}
	type row struct {
		A, B   [4]int32
		Result legacy.PortTestWorldGeometryResult
	}
	var rows []row
	segments := [][4]int32{{0, 0, 0, 0}, {0, 0, 10, 10}, {10, 10, 0, 0}, {0, 10, 10, 0}, {0, 0, 10, 0}, {5, -10, 5, 10}, {-10, -10, 10, 10}, {1, 1, 1, 10}, {-2147483648, 0, 2147483647, 0}, {-65536, 65536, 65536, -65536}, {5888, 0, 0, 5888}}
	for _, a := range segments {
		for _, b := range segments {
			s := worldGeometrySpec("segments")
			for i := 0; i < 4; i++ {
				s.Words[i] = uint32(a[i])
				s.Words[4+i] = uint32(b[i])
			}
			r := worldGeometryCall(t, s)
			for i, v := range s.Words {
				if r.Words[i] != v {
					t.Fatal("segment input mutation")
				}
			}
			rows = append(rows, row{a, b, r})
		}
	}
	spellbookCapture(t, "world-geometry-integer-segments", rows, "cedc1970f7b2c32c16c28d19d79ce48d7cabf2922fd4ab36b243703f3cde9820")
}
func TestWorldGeometryRectangles(t *testing.T) {
	type row struct {
		Op     string
		Point  [2]uint32
		Bounds [4]uint32
		Result legacy.PortTestWorldGeometryResult
	}
	var rows []row
	for _, op := range []string{"rect-int", "rect-float", "rect-float-alt"} {
		for _, bounds := range [][4]int32{{0, 0, 10, 10}, {10, 10, 0, 0}, {-10, -20, 10, 20}, {5, 5, 5, 5}, {-2147483648, -2147483648, 2147483647, 2147483647}} {
			for _, x := range []int32{-21, -10, -1, 0, 5, 10, 11, 21} {
				for _, y := range []int32{-21, -10, 0, 5, 10, 21} {
					s := worldGeometrySpec(op)
					s.Words[0] = uint32(x)
					s.Words[1] = uint32(y)
					for i, v := range bounds {
						s.Words[4+i] = uint32(v)
					}
					if op != "rect-int" {
						s.Words[0] = math.Float32bits(float32(x))
						s.Words[1] = math.Float32bits(float32(y))
						for i, v := range bounds {
							s.Words[4+i] = math.Float32bits(float32(v))
						}
					}
					r := worldGeometryCall(t, s)
					want := x >= bounds[0] && x <= bounds[2] && y >= bounds[1] && y <= bounds[3]
					if (r.Return != 0) != want {
						t.Fatal("rectangle inclusion", op, x, y, bounds, r.Return)
					}
					rows = append(rows, row{op, [2]uint32{s.Words[0], s.Words[1]}, [4]uint32{s.Words[4], s.Words[5], s.Words[6], s.Words[7]}, r})
				}
			}
		}
		if op != "rect-int" {
			for _, bits := range []uint32{0x80000000, 0x00000001, 0x7f800000, 0xff800000, 0x7fc00001} {
				s := worldGeometrySpec(op)
				s.Words[0] = bits
				s.Words[1] = math.Float32bits(5)
				copy(s.Words[4:], worldGeometryWords(0, 0, 10, 10))
				r := worldGeometryCall(t, s)
				rows = append(rows, row{op, [2]uint32{bits, s.Words[1]}, [4]uint32{s.Words[4], s.Words[5], s.Words[6], s.Words[7]}, r})
			}
		}
	}
	spellbookCapture(t, "world-geometry-rectangles", rows, "c5fb9f17a5234f40eb521e1d5f8bfcf3e717e069e427c45f40b185c47c389a02")
}
func TestWorldGeometryDirectionsAndVectors(t *testing.T) {
	words := worldGeometryTables(t)
	type row struct {
		Op        string
		Threshold uint32
		Input     []uint32
		Index     int32
		Result    legacy.PortTestWorldGeometryResult
	}
	var rows []row
	for _, threshold := range []uint32{0, 6, 20} {
		*words["directionThreshold"] = threshold
		for index := int32(0); index < 256; index++ {
			for _, op := range []string{"indexed-direction", "direction4-index", "direction4-angle"} {
				s := worldGeometrySpec(op)
				s.Ints[0] = index
				r := worldGeometryCall(t, s)
				rows = append(rows, row{op, threshold, nil, index, r})
			}
		}
	}
	for x := int32(-1); x <= 1; x++ {
		for y := int32(-1); y <= 1; y++ {
			s := worldGeometrySpec("direction-angle")
			s.Words[0] = uint32(x)
			s.Words[1] = uint32(y)
			r := worldGeometryCall(t, s)
			rows = append(rows, row{s.Op, 0, []uint32{uint32(x), uint32(y)}, 0, r})
		}
	}
	for _, tc := range []struct {
		x, y  float32
		angle int32
	}{{1, 0, 0}, {0, 1, 64}, {-1, 0, 128}, {0, -1, 192}} {
		s := worldGeometrySpec("vector-angle")
		copy(s.Words, worldGeometryWords(tc.x, tc.y))
		if got := worldGeometryCall(t, s).Return; got != tc.angle {
			t.Fatal("cardinal angle", tc, got)
		}
	}
	for _, op := range []string{"vector-angle", "normalize"} {
		for _, x := range []float32{-100000, -1, -0.00001, 0, 0.00001, 1, 3, 100000} {
			for _, y := range []float32{-100000, -4, -1, 0, 1, 4, 100000} {
				s := worldGeometrySpec(op)
				copy(s.Words, worldGeometryWords(x, y))
				r := worldGeometryCall(t, s)
				if op == "normalize" && x == 3 && y == 4 {
					if math.Abs(float64(math.Float32frombits(r.Words[0]))-0.6) > 1e-6 || math.Abs(float64(math.Float32frombits(r.Words[1]))-0.8) > 1e-6 {
						t.Fatal("normalized 3-4-5 vector")
					}
				}
				rows = append(rows, row{op, 0, s.Words[:2], 0, r})
			}
		}
	}
	spellbookCapture(t, "world-geometry-directions-vectors", rows, "30839bd9cb91b4c78f50f745a64b2114e2123d1d889e7c2c9454c1d390baedfe")
}
func TestWorldGeometryShapeAndMapCoordinates(t *testing.T) {
	worldGeometryTables(t)
	type row struct {
		Op     string
		Input  []uint32
		Result legacy.PortTestWorldGeometryResult
	}
	var rows []row
	for _, w := range []float32{-20, 0, 1, 10, 100, 100000} {
		for _, h := range []float32{-10, 0, 1, 30, 100000} {
			s := worldGeometrySpec("box-calc")
			s.Words[0] = 3
			s.Words[3] = math.Float32bits(w)
			s.Words[4] = math.Float32bits(h)
			r := worldGeometryCall(t, s)
			if r.Words[3] != s.Words[3] || r.Words[4] != s.Words[4] {
				t.Fatal("box dimensions changed")
			}
			rows = append(rows, row{s.Op, append([]uint32(nil), s.Words[:13]...), r})
		}
	}
	for _, x := range []float32{-1, 0, 80.49999, 80.5, 80.50001, 100, 5853.4995, 5853.5, 5853.5005, 6000} {
		for _, y := range []float32{0, 80.5, 100, 5853.5, 6000} {
			for _, alias := range []bool{false, true} {
				s := worldGeometrySpec("map-coordinates")
				copy(s.Words, worldGeometryWords(x, y))
				if alias {
					s.Offsets[1] = 0
				}
				r := worldGeometryCall(t, s)
				if r.Return != 1 {
					t.Fatal("coordinate transform return")
				}
				if !alias && x == 0 && y == 0 {
					if math.Float32frombits(r.Words[0]) != 82.5 || math.Float32frombits(r.Words[1]) != 81.5 || r.Words[4] != 0 {
						t.Fatal("coordinate clamp and diagonal", r)
					}
				}
				rows = append(rows, row{s.Op, s.Words[:2], r})
			}
		}
	}
	for _, which := range []int{0, 1} {
		s := worldGeometrySpec("map-coordinates")
		s.Offsets[which] = -1
		if worldGeometryCall(t, s).Return != 0 {
			t.Fatal("nil coordinate transform")
		}
	}
	// Keep the exact shipped direction-table bytes in the baseline audit.
	table := blobdata.PortTestWorldGeometryTables()[230056]
	if len(table) != 36 || binary.LittleEndian.Uint32(table) == 0 && binary.LittleEndian.Uint32(table[32:]) == 0 {
		t.Fatal("uninitialized direction table")
	}
	spellbookCapture(t, "world-geometry-shape-map", rows, "743c3165e1f53ae7ed0eb1ef57dbb9c540676e87e5a5a4efea536fd4ad13259f")
}
