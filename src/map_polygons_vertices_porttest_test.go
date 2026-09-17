//go:build porttest

package opennox

import (
	"encoding/binary"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"math"
	"reflect"
	"testing"
	"unsafe"
)

func TestMapPolygonsNearestVertices(t *testing.T) {
	o := newMapPolygonsOwner(t)
	for _, v := range []struct {
		index int
		x, y  float32
	}{{1, 0, 0}, {3, 10, 0}, {5, 0, 10}, {7, -10, 0}, {9, 10, 0}, {11, 16777216, 16777216}} {
		legacy.PortTestMapPolygonVertex(v.x, v.y, v.index, 0)
	}
	// Inactive vertices must not participate even when closer than every live one.
	legacy.PortTestMapPolygonVertex(1, 1, 2, 0)
	binary.LittleEndian.PutUint32(o.vertex(2)[12:], 0)
	for _, tc := range []struct {
		x, y, limit float32
		want        int
	}{{0, 0, 0, 0}, {0, 0, 1, 1}, {1, 0, 1, 0}, {1, 0, 1.0001, 1}, {5, 0, 26, 1}, {10, 0, 1, 3}, {-9, 0, 2, 7}, {1, 1, 3, 1}, {500, 500, 10, 0}} {
		if got := o.vertexID(legacy.PortTestMapPolygonVertexNearest(tc.x, tc.y, tc.limit)); got != tc.want {
			t.Fatal("nearest contract", tc, got)
		}
	}
	var rows []struct {
		X, Y, Limit uint32
		Index       int
	}
	inputs := []float32{-16777216, -10, -0.5, 0, 0.5, 5, 9.999999, 10, 10.000001, 16777216, math.Float32frombits(0x7f800000), math.Float32frombits(0x7fc00001)}
	for _, x := range inputs {
		for _, y := range []float32{-10, 0, 0.5, 10, 16777216} {
			for _, limit := range []float32{-1, 0, 0.25, 1, 25, 25.000002, 100, 900, math.Float32frombits(0x7f800000)} {
				p := legacy.PortTestMapPolygonVertexNearest(x, y, limit)
				rows = append(rows, struct {
					X, Y, Limit uint32
					Index       int
				}{math.Float32bits(x), math.Float32bits(y), math.Float32bits(limit), o.vertexID(p)})
			}
		}
	}
	spellbookCapture(t, "map-polygons-nearest-vertices", rows, "878f3cd54d120a0ae23ff3531f64ed4fdb122873a46401ed5df1acbabb9c6eaf")
}
func TestMapPolygonsVertexSelection(t *testing.T) {
	var rows []any
	for _, editor := range []bool{false, true} {
		t.Run(map[bool]string{false: "runtime", true: "editor"}[editor], func(t *testing.T) {
			o := newMapPolygonsOwner(t)
			if editor {
				noxflags.PortTestGameFlags(0x200000)
			}
			for i, xy := range [][2]float32{{10.5, 20.5}, {10.5, 20.5}, {11, 21}, {100, 120}, {10.5, 20.5}} {
				p := legacy.PortTestMapPolygonVertexSelect(xy[0], xy[1])
				count := binary.LittleEndian.Uint16(o.control)
				if i == 0 && (o.vertexID(p) != 1 || count != 1) {
					t.Fatal("initial selection", o.vertexID(p), count)
				}
				if i == 1 && editor && (p != nil || count != 1) {
					t.Fatal("duplicate editor selection", count)
				}
				if i == 1 && !editor && (o.vertexID(p) != 2 || count != 2) {
					t.Fatal("runtime insertion", o.vertexID(p), count)
				}
				ids := make([]uint32, int(count))
				for j := range ids {
					ids[j] = binary.LittleEndian.Uint32(o.selection[4*j:])
				}
				rows = append(rows, []any{editor, i, o.vertexID(p), ids, *o.words["vertices"]})
			}
			first := legacy.PortTestMapPolygonPointer("vertex-get", nil, 1)
			if legacy.PortTestMapPolygonScalar("selected", first) != 1 {
				t.Fatal("selection membership")
			}
			binary.LittleEndian.PutUint16(o.control, 0)
			if legacy.PortTestMapPolygonScalar("selected", first) != 0 {
				t.Fatal("empty selection")
			}
		})
	}
	spellbookCapture(t, "map-polygons-vertex-selection", rows, "b4ff24bd1afdbbef8a5adc33b279fe0459df7b87ac133e0761ebd654f758659b")
}
func TestMapPolygonsConstructionAndRemap(t *testing.T) {
	var rows []any
	for _, editor := range []bool{false, true} {
		for _, count := range []int{0, 1, 2, 3, 4, 7} {
			t.Run(map[bool]string{false: "runtime", true: "editor"}[editor]+string(rune('0'+count)), func(t *testing.T) {
				o := newMapPolygonsOwner(t)
				if editor {
					noxflags.PortTestGameFlags(0x200000)
				}
				points := [][2]float32{{10, 20}, {-3, 40}, {25, -7}, {30, 50}, {11, 22}, {12, 23}, {13, 24}}
				for i := 0; i < count; i++ {
					legacy.PortTestMapPolygonVertex(points[i][0], points[i][1], i+1, 0)
					binary.LittleEndian.PutUint32(o.selection[4*i:], uint32(i+1))
				}
				binary.LittleEndian.PutUint16(o.control, uint16(count))
				legacy.PortTestMapPolygonScalar("construct", nil)
				if count < 3 {
					if *o.words["polygons"] != 1 || int(binary.LittleEndian.Uint16(o.control)) != count {
						t.Fatal("incomplete selection changed")
					}
					rows = append(rows, []any{editor, count, "incomplete"})
					return
				}
				p := o.polygon(1)
				if binary.LittleEndian.Uint16(o.control) != 0 || binary.LittleEndian.Uint16(p[128:]) != uint16(count) || *o.words["polygons"] != 2 {
					t.Fatal("constructed record")
				}
				bounds := [4]int32{int32(binary.LittleEndian.Uint32(p[88:])), int32(binary.LittleEndian.Uint32(p[92:])), int32(binary.LittleEndian.Uint32(p[96:])), int32(binary.LittleEndian.Uint32(p[100:]))}
				want := [4]int32{-3, -7, 25, 40}
				if count >= 4 {
					want = [4]int32{-3, -7, 30, 50}
				}
				if bounds != want {
					t.Fatal("polygon bounds", bounds, want)
				}
				meta := *(*unsafe.Pointer)(unsafe.Pointer(&p[0]))
				if (meta != nil) != editor {
					t.Fatal("metadata allocation", editor)
				}
				if meta != nil {
					for _, b := range unsafe.Slice((*byte)(meta), 256) {
						if b != 0 {
							t.Fatal("metadata not initialized")
						}
					}
				}
				vertices := unsafe.Slice((*uint32)(*(*unsafe.Pointer)(unsafe.Pointer(&p[108]))), count)
				for i, v := range vertices {
					if v != uint32(i+1) {
						t.Fatal("vertex copy", i, v)
					}
				}
				vertices[0] = 50
				vertices[1] = 60
				vertices[2] = 99
				legacy.PortTestMapPolygonRemapAdd(1, 50)
				legacy.PortTestMapPolygonRemapAdd(2, 60)
				legacy.PortTestMapPolygonRemapAdd(3, 50)
				legacy.PortTestMapPolygonScalar("remap", o.pointer(1))
				if !reflect.DeepEqual(vertices[:3], []uint32{3, 2, 99}) {
					t.Fatal("newest remap wins; unknown stays", vertices)
				}
				record := append([]byte(nil), p...)
				clear(record[:4])
				clear(record[108:112])
				rows = append(rows, []any{editor, count, record, append([]uint32(nil), vertices...)})
			})
		}
	}
	spellbookCapture(t, "map-polygons-construction-remap", rows, "1a9735f5b52561323d965f919bb3cbd9d29166e0be657757b5bca5d96e31b49c")
}
