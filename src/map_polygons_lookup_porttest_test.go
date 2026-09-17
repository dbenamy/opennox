//go:build porttest

package opennox

import (
	"encoding/binary"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

// The legacy finite ray counts both edges at a shared vertex. Keep that
// limitation explicit, separate from ordinary interior/exterior contracts.
func TestMapPolygonsContainmentContracts(t *testing.T) {
	o := newMapPolygonsOwner(t)
	p := o.construct(t, [][2]float32{{100, 200}, {110, 200}, {110, 210}, {100, 210}})
	o.construct(t, [][2]float32{{105, 200}, {110, 205}, {105, 210}, {100, 205}})
	o.construct(t, [][2]float32{{130, 200}, {140, 200}, {140, 210}, {135, 205}, {130, 210}})
	binary.LittleEndian.PutUint32(o.polygon(2)[116:], 42)
	binary.LittleEndian.PutUint32(o.polygon(3)[124:], 17)
	for _, tc := range []struct{ x, y, contains, first, script int }{{2, 2, 1, 1, 0}, {5, 5, 1, 1, 2}, {35, 2, 0, 3, 3}, {-5, 5, 0, 0, 0}, {50, 50, 0, 0, 0}} {
		point := [2]int32{int32(tc.x + 100), int32(tc.y + 200)}
		if got := legacy.PortTestMapPolygonContains(p, &point); got != tc.contains {
			t.Fatal("containment", tc, got)
		}
		if got := o.polygonID(legacy.PortTestMapPolygonFind("point", &point, 0, 0)); got != tc.first {
			t.Fatal("first containing polygon", tc, got)
		}
		if got := o.polygonID(legacy.PortTestMapPolygonFind("script", &point, 0, 0)); got != tc.script {
			t.Fatal("script-filtered polygon", tc, got)
		}
	}
	point := [2]int32{105, 205}
	if got := o.polygonID(legacy.PortTestMapPolygonFind("point", &point, 2, 0)); got != 2 {
		t.Fatal("cached overlap takes precedence", got)
	}
}

func TestMapPolygonsRayVertexLimitation(t *testing.T) {
	o := newMapPolygonsOwner(t)
	p := o.construct(t, [][2]float32{{0, 0}, {10, 0}, {10, 10}, {0, 10}})
	for phase, want := range []int{0, 1, 0, 1} {
		binary.LittleEndian.PutUint32(o.control[8:], uint32(phase))
		point := [2]int32{2, 2}
		if got := legacy.PortTestMapPolygonContains(p, &point); got != want {
			t.Fatalf("phase %d: got %d, want %d", phase, got, want)
		}
	}
}

func TestMapPolygonsContainmentAndCache(t *testing.T) {
	o := newMapPolygonsOwner(t)
	o.construct(t, [][2]float32{{0, 0}, {10, 0}, {10, 10}, {0, 10}})
	o.construct(t, [][2]float32{{5, 0}, {10, 5}, {5, 10}, {0, 5}})
	o.construct(t, [][2]float32{{30, 0}, {40, 0}, {40, 10}, {35, 5}, {30, 10}})
	binary.LittleEndian.PutUint32(o.polygon(2)[116:], 42)
	binary.LittleEndian.PutUint32(o.polygon(3)[124:], 17)
	point := [2]int32{5, 5}
	before := binary.LittleEndian.Uint32(o.control[8:])
	if legacy.PortTestMapPolygonContains(nil, &point) != 0 || binary.LittleEndian.Uint32(o.control[8:]) != before {
		t.Fatal("nil polygon")
	}
	type row struct {
		Op     string
		Point  [2]int32
		Cache  int
		Phase  uint32
		Result int
		After  uint32
	}
	var rows []row
	for _, x := range []int32{-1, 0, 1, 2, 5, 9, 10, 11, 29, 30, 34, 35, 36, 39, 40, 41} {
		for _, y := range []int32{-1, 0, 1, 4, 5, 6, 9, 10, 11} {
			for _, cache := range []int{0, 1, 2, 3, -559023410} {
				for _, phase := range []uint32{0, 1, 2, 3, 0xfffffffe} {
					for _, op := range []string{"point", "script", "index"} {
						binary.LittleEndian.PutUint32(o.control[8:], phase)
						point := [2]int32{x, y}
						result := 0
						if op == "index" {
							result = legacy.PortTestMapPolygonIndex(&point, cache)
						} else {
							result = o.polygonID(legacy.PortTestMapPolygonFind(op, &point, cache, 0))
						}
						if point != [2]int32{x, y} {
							t.Fatal("containment mutated input", point)
						}
						rows = append(rows, row{op, point, cache, phase, result, binary.LittleEndian.Uint32(o.control[8:])})
					}
				}
			}
		}
	}
	spellbookCapture(t, "map-polygons-containment-cache", rows, "acdec3c27122a1de1c3d634b33241f6a49a2d13459d7f9b3a523995a3f76bd18")
}
func TestMapPolygonsEdgeProjection(t *testing.T) {
	o := newMapPolygonsOwner(t)
	p := o.construct(t, [][2]float32{{0, 0}, {10, 0}, {10, 10}, {0, 10}})
	o.construct(t, [][2]float32{{30, 0}, {40, 0}, {40, 10}, {30, 10}})
	for _, tc := range []struct {
		input, output [2]int32
		distance      float32
		result        int
	}{{[2]int32{5, 1}, [2]int32{5, 0}, 1, 1}, {[2]int32{5, 5}, [2]int32{5, 5}, 1, 0}, {[2]int32{-2, 5}, [2]int32{0, 5}, 2, 1}, {[2]int32{-1, -1}, [2]int32{-1, -1}, 2, 0}} {
		point := tc.input
		got := legacy.PortTestMapPolygonEdge(p, &point, tc.distance)
		if got != tc.result || point != tc.output {
			t.Fatal("edge projection", tc, got, point)
		}
	}
	type row struct {
		Input, After         [2]int32
		Distance             float32
		Cache, Direct, Found int
	}
	var rows []row
	for _, x := range []int32{-3, -1, 0, 1, 5, 9, 10, 11, 13, 29, 30, 35, 40, 41} {
		for _, y := range []int32{-1, 0, 1, 5, 9, 10, 11} {
			for _, distance := range []float32{-1, 0, 0.5, 1, 2, 5} {
				for _, cache := range []int{0, 1, 2, -559023410} {
					input := [2]int32{x, y}
					point := input
					direct := legacy.PortTestMapPolygonEdge(p, &point, distance)
					rows = append(rows, row{input, point, distance, cache, direct, -1})
					point = input
					found := o.polygonID(legacy.PortTestMapPolygonFind("edge", &point, cache, distance))
					rows = append(rows, row{input, point, distance, cache, -1, found})
				}
			}
		}
	}
	spellbookCapture(t, "map-polygons-edge-projection", rows, "4c9a362464b7954459d8bf61b8cf8841a3bf8d085b451f4018f5246a8268306f")
}
