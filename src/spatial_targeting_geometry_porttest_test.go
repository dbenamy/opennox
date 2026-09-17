//go:build porttest

package opennox

import (
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/legacy"
	"math"
	"testing"
)

func TestSpatialTargetingLocalCoordinates(t *testing.T) {
	type row struct {
		Grid     [2]int32
		Point    [2]uint32
		Quadrant int8
		Triangle int32
	}
	var rows []row
	values := []float32{-513, -257, -256, -255, -46, -23.001, -23, -22.999, -12, -1, -.001, 0, .001, 1, 10.999, 11, 11.499, 11.5, 11.999, 12, 21.999, 22, 22.999, 23, 23.001, 46, 255, 256, 257, 513}
	for _, grid := range [][2]int32{{0, 0}, {6, 4}, {-6, -4}, {1000000, -1000000}, {2147483647, -2147483648}} {
		for _, x := range values {
			for _, y := range values {
				p := types.Pointf{float32(int32(uint32(grid[0])*23)) + x, float32(int32(uint32(grid[1])*23)) + y}
				before := p
				q := legacy.PortTestSpatialQuadrant(&p)
				tri := legacy.PortTestSpatialTriangle(&p, grid)
				if p != before {
					t.Fatal("coordinate helper changed its input")
				}
				if q != 1 && q != 2 && q != 4 && q != 8 {
					t.Fatal("quadrant must select exactly one region", q)
				}
				rows = append(rows, row{grid, motionBits(p), q, tri})
			}
		}
	}
	// Independently classify every integer point in one tile by its two diagonals.
	for x := 0; x < 23; x++ {
		for y := 0; y < 23; y++ {
			p := types.Pointf{float32(x), float32(y)}
			want := int32(0)
			if x > y {
				want = 1
			}
			if x+y >= 22 {
				if x > y {
					want = 2
				} else {
					want = 3
				}
			}
			if got := legacy.PortTestSpatialTriangle(&p, [2]int32{}); got != want {
				t.Fatalf("tile triangle (%d,%d): %d != %d", x, y, got, want)
			}
			quadrants := [2][2]int8{{2, 8}, {1, 4}}
			ix, iy := 0, 0
			if x >= 12 {
				ix = 1
			}
			if y >= 12 {
				iy = 1
			}
			if got := legacy.PortTestSpatialQuadrant(&p); got != quadrants[ix][iy] {
				t.Fatal("tile quadrant", x, y, got)
			}
		}
	}
	spellbookCapture(t, "spatial-targeting-local-coordinates", rows, "2a2645a9588c9050c68ebfef18bbdcda96e4740e2249f20ab9298985bffa813a")
}

func TestSpatialTargetingWallNormals(t *testing.T) {
	o := newCollisionCoreOwner(t)
	o.s.PortTestAIEmptyMap()
	t.Cleanup(o.s.Map.Free)
	configure, guard, free := o.s.PortTestGeometryWalls()
	t.Cleanup(free)
	type row struct {
		Direction          int
		Neighbours, Window bool
		Start, End         [2]float32
		Return             int32
		Normal             [2]uint32
	}
	var rows []row
	starts := [][2]float32{{2, 11}, {11, 2}, {21, 11}, {11, 21}, {0, 0}, {22, 0}, {0, 22}, {22, 22}, {11, 11}, {11.5, 11.5}, {-1, -257}, {257, 513}}
	ends := [][2]float32{{2, 2}, {21, 2}, {21, 21}, {2, 21}, {11, 11}, {11.5, 11.5}, {12, 12}, {-1, -1}, {23, 23}}
	contacts := 0
	for dir := -1; dir < 11; dir++ {
		for _, neighbours := range []bool{false, true} {
			for _, window := range []bool{false, true} {
				configure(dir, neighbours, window)
				for _, start := range starts {
					for _, end := range ends {
						grid := [2]int32{6, 4}
						ray := [4]float32{138 + start[0], 92 + start[1], 138 + end[0], 92 + end[1]}
						before := ray
						normal := types.Pointf{17, -19}
						rv := legacy.PortTestSpatialNormal(&grid, &ray, &normal)
						if grid != ([2]int32{6, 4}) || ray != before || !guard() {
							t.Fatal("normal helper changed input or wall guards")
						}
						if rv == 0 && normal != (types.Pointf{17, -19}) {
							t.Fatal("no contact changed normal")
						}
						if dir == -1 && rv != 0 {
							t.Fatal("empty grid reported contact")
						}
						if rv != 0 && normal != (types.Pointf{17, -19}) {
							contacts++
							if math.Abs(float64(normal.X)) != float64(float32(.70709997)) || math.Abs(float64(normal.Y)) != float64(float32(.70709997)) {
								t.Fatal("contact normal must lie on tile diagonal", normal)
							}
							if dir == 0 && !neighbours && !window && start == ([2]float32{2, 11}) && (normal.X >= 0 || normal.Y >= 0) {
								t.Fatal("northwest side should face northwest", normal)
							}
						}
						rows = append(rows, row{dir, neighbours, window, start, end, rv, motionBits(normal)})
					}
				}
			}
		}
	}
	if contacts == 0 {
		t.Fatal("no real wall contact exercised")
	}
	spellbookCapture(t, "spatial-targeting-wall-normals", rows, "d2857f52cbd68d9c2863672b2013cdf846b70c49c80a8ec9302cc719a065394d")
}
