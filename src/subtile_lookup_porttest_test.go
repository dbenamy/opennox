//go:build porttest

package opennox

import (
	"reflect"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
)

func subtileI32(v uint32) int32   { return int32(v) }
func subtileAdd(a, b int32) int32 { return int32(uint32(a) + uint32(b)) }
func subtileNorm(w, h byte, edge int32) int32 {
	wi, hi := int32(w), int32(h)
	if wi == 3 && hi == 3 {
		return edge
	}
	if edge == 0 {
		return 0
	}
	if edge <= wi-2 {
		return 1
	}
	if edge == wi-1 {
		return 2
	}
	b := wi + 2*hi - 4
	if edge < b {
		return 3 + int32((uint32(wi)^uint32(edge))&1)
	}
	if edge == b {
		return 5
	}
	c := 2*(hi+wi) - 6
	if edge > c {
		if edge == c+1 {
			return 7
		}
		return subtileAdd(edge, 2*(6-hi-wi))
	}
	return 6
}
func subtilePredicate(x, y, cat int32) int32 {
	if cat < 0 || cat >= 12 {
		return 0
	}
	diagonal := int32(uint32(46) - uint32(y))
	a, b, c, d := x <= y, x >= y, x >= diagonal, x <= diagonal
	regions := [12]bool{a && c, a, a && d, c, d, b && c, b, b && d, a || c, a || d, b || d, b || c}
	if regions[cat] {
		return 1
	}
	return 0
}
func checkSubtile(t *testing.T, specs []legacy.PortTestSubtileLookupSpec) []legacy.PortTestSubtileLookupResult {
	t.Helper()
	got := legacy.PortTestSubtileLookup(specs)
	if !got.GuardsRestored || len(got.Results) != len(specs) || !reflect.DeepEqual(got.TableBefore, got.TableAfterRestore) {
		t.Fatal("fixture restore/results")
	}
	for i, r := range got.Results {
		if !r.TableUnchanged || !r.PointUnchanged || !r.NodesUnchanged || !r.GuardsUnchanged {
			t.Fatalf("case %d mutated state: %+v", i, r)
		}
	}
	return got.Results
}

func TestSubtilePredicateCABI(t *testing.T) {
	var specs []legacy.PortTestSubtileLookupSpec
	for x := int32(-1); x <= 47; x++ {
		for y := int32(-1); y <= 47; y++ {
			for cat := int32(-2); cat <= 13; cat++ {
				specs = append(specs, legacy.PortTestSubtileLookupSpec{Mode: 0, X: x, Y: y, Category: cat})
			}
		}
	}
	ext := []int32{subtileI32(0x80000000), -1, 0, 1, 45, 46, subtileI32(0x7fffffff)}
	for _, x := range ext {
		for _, y := range ext {
			for _, cat := range []int32{-1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, subtileI32(0x80000000), subtileI32(0x7fffffff)} {
				specs = append(specs, legacy.PortTestSubtileLookupSpec{Mode: 0, X: x, Y: y, Category: cat})
			}
		}
	}
	t.Logf("%d point predicates", len(specs))
	got := checkSubtile(t, specs)
	for i, s := range specs {
		if got[i].Return != subtilePredicate(s.X, s.Y, s.Category) {
			t.Fatalf("predicate %d %+v got=%d", i, s, got[i].Return)
		}
	}
}

func TestSubtileLookupCABI(t *testing.T) {
	// One normal row produces all twelve normalized categories independently;
	// values repeat so the last matching node rule is observable.
	rows := []legacy.PortTestSubtileRow{{Index: 0, Width: 5, Height: 6}}
	var nodes []legacy.PortTestSubtileNode
	for edge := int32(0); edge <= 21; edge++ {
		c := subtileNorm(5, 6, edge)
		if c >= 0 && c < 12 {
			nodes = append(nodes, legacy.PortTestSubtileNode{Value: 100 + edge, Row: 0, Edge: edge})
		}
	}
	// Same normalized category as the preceding default node; if it matches,
	// the C traversal must choose this final node rather than the first match.
	nodes = append(nodes, legacy.PortTestSubtileNode{Value: 999, Row: 0, Edge: 21})
	specs := []legacy.PortTestSubtileLookupSpec{{Mode: 1, X: 12, Y: 12, Fallback: -77, Rows: rows, Nodes: nodes}, {Mode: 1, X: 12, Y: 12, Fallback: 0x12345678, NilList: true, NilPoint: true}, {Mode: 1, X: 12, Y: 12, Fallback: -9, Rows: rows, Nodes: nil}}
	// All 64 physical rows: 3x3 retains raw categories and verifies C list link
	// traversal does not substitute an index or active-count lookup.
	for row := 0; row < 64; row++ {
		specs = append(specs, legacy.PortTestSubtileLookupSpec{Mode: 1, X: int32(row % 46), Y: int32((row * 3) % 46), Fallback: -1, Rows: []legacy.PortTestSubtileRow{{Index: row, Width: 3, Height: 3}}, Nodes: []legacy.PortTestSubtileNode{{Value: int32(0x24680000 + 17*row), Row: int32(row), Edge: int32(row % 12)}}})
	}
	// Every possible last matching position, including no match and a
	// nonmatching tail. Values are raw words distinct from node/row indices.
	for length := 0; length <= 31; length++ {
		for last := -1; last < length; last++ {
			for _, fallback := range []int32{0, -1, -2147483648, 2147483647} {
				var ordered []legacy.PortTestSubtileNode
				var defs []legacy.PortTestSubtileRow
				for i := 0; i < length; i++ {
					row := int32((i*7 + length) % 64)
					category := int32(6) // at(0,46), x>=y is false
					if i == last || (i < last && i%2 == 0) {
						category = 1
					}
					ordered = append(ordered, legacy.PortTestSubtileNode{Value: int32(uint32(0xa5000000) + uint32(i)*173 + uint32(length)), Row: row, Edge: category})
					defs = append(defs, legacy.PortTestSubtileRow{Index: int(row), Width: 3, Height: 3})
				}
				specs = append(specs, legacy.PortTestSubtileLookupSpec{Mode: 1, X: 0, Y: 46, Fallback: fallback, Rows: defs, Nodes: ordered})
			}
		}
	}
	specs = append(specs, legacy.PortTestSubtileLookupSpec{Mode: 1, NilList: true, Fallback: -2147483648})
	t.Logf("%d list lookups", len(specs))
	got := checkSubtile(t, specs)
	model := func(s legacy.PortTestSubtileLookupSpec) int32 {
		if s.NilList || len(s.Nodes) == 0 {
			return s.Fallback
		}
		out := s.Fallback
		dims := map[int32][2]byte{}
		for _, r := range s.Rows {
			dims[int32(r.Index)] = [2]byte{r.Width, r.Height}
		}
		for _, n := range s.Nodes {
			d := dims[n.Row]
			if subtilePredicate(s.X, s.Y, subtileNorm(d[0], d[1], n.Edge)) != 0 {
				out = n.Value
			}
		}
		return out
	}
	for i, s := range specs {
		if got[i].Return != model(s) {
			t.Fatalf("lookup %d got=%d want=%d spec=%+v", i, got[i].Return, model(s), s)
		}
	}
}
