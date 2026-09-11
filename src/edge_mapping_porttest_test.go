//go:build porttest

package opennox

import (
	"testing"

	"github.com/opennox/libs/prand"
	"github.com/opennox/opennox/v1/legacy"
)

func edgeI32(bits uint32) int32 { return int32(bits) }
func edgeAdd(a, b int32) int32  { return int32(uint32(a) + uint32(b)) }
func edgeTransform(r *prand.Rand, w, h byte, edge int32) int32 {
	wi, hi := int(w), int(h)
	if wi == 3 && hi == 3 {
		return edge
	}
	switch edge {
	case 0:
		return 0
	case 1:
		return int32(r.IntClamp(1, wi-2))
	case 2:
		return int32(wi - 1)
	case 3:
		return int32(wi + 2*r.IntClamp(0, hi-3))
	case 4:
		return int32(wi + 2*r.IntClamp(0, hi-3) + 1)
	case 5:
		return int32(wi + 2*hi - 4)
	case 6:
		return int32(r.IntClamp(1, wi-2) + wi + 2*hi - 4)
	case 7:
		return int32(2*(wi+hi) - 5)
	default:
		return edgeAdd(edge, int32(2*(wi+hi)-12))
	}
}
func edgeNormalize(w, h byte, v int32) int32 {
	wi, hi := int32(w), int32(h)
	if wi == 3 && hi == 3 {
		return v
	}
	if v == 0 {
		return 0
	}
	if v <= wi-2 {
		return 1
	}
	if v == wi-1 {
		return 2
	}
	m := wi + 2*hi - 4
	if v < m {
		if (uint32(wi)^uint32(v))&1 != 0 {
			return 4
		}
		return 3
	}
	if v == m {
		return 5
	}
	if v > 2*(hi+wi)-6 {
		if v == 2*(hi+wi)-5 {
			return 7
		}
		return edgeAdd(v, 2*(6-hi-wi))
	}
	return 6
}

func TestEdgeMappingDirectCABI(t *testing.T) {
	edges := []int32{0, 1, 2, 3, 4, 5, 6, 7, -1, 8, edgeI32(0x80000000), edgeI32(0x7fffffff)}
	specs := make([]legacy.PortTestEdgeDirectSpec, 0, 256*256*len(edges))
	for w := 0; w < 256; w++ {
		for h := 0; h < 256; h++ {
			for _, e := range edges {
				specs = append(specs, legacy.PortTestEdgeDirectSpec{Width: byte(w), Height: byte(h), Edge: e})
			}
		}
	}
	t.Logf("%d direct operations", len(specs))
	seed := 137
	got := legacy.PortTestEdgeMapping(seed, specs, nil)
	if len(got.Direct) != len(specs) || !got.Restored || !got.TableUnchanged || !got.MappingUnchanged || !got.GuardsUnchanged || string(got.TableBefore) != string(got.TableAfterRestore) || len(got.MappingBefore) != len(got.MappingAfterRestore) || got.LogicBefore != prand.New(seed).Index() || got.OtherBefore != prand.New(seed+1).Index() {
		t.Fatalf("fixture state: table=%v mapping=%v guards=%v restored=%v", got.TableUnchanged, got.MappingUnchanged, got.GuardsUnchanged, got.Restored)
	}
	want := prand.New(seed)
	for i, s := range specs {
		r := got.Direct[i]
		exp := edgeTransform(want, s.Width, s.Height, s.Edge)
		if r.Result != exp || r.LogicIndex != want.Index() || r.OtherIndex != prand.New(seed+1).Index() {
			t.Fatalf("direct %d w/h/e=%d/%d/%d got=%+v want=%d/%d", i, s.Width, s.Height, s.Edge, r, exp, want.Index())
		}
	}
	if got.LogicAfter != want.Index() || got.OtherAfter != prand.New(seed+1).Index() {
		t.Fatalf("final RNG %d/%d want %d/%d", got.LogicAfter, got.OtherAfter, want.Index(), prand.New(seed+1).Index())
	}
}

func TestEdgeMappingDependentCABI(t *testing.T) {
	var specs []legacy.PortTestEdgeMapSpec
	// Every physical row, normalized class, mapping column and transform branch.
	for index := int32(0); index < 64; index++ {
		for _, dims := range [][2]byte{{3, 3}, {5, 6}, {255, 255}} {
			for class := int32(0); class < 12; class++ {
				current := class
				if dims != [2]byte{3, 3} {
					w, h := int32(dims[0]), int32(dims[1])
					switch class {
					case 0:
						current = 0
					case 1:
						current = 1
					case 2:
						current = w - 1
					case 3:
						current = w
					case 4:
						current = w + 1
					case 5:
						current = w + 2*h - 4
					case 6:
						current = w + 2*h - 3
					case 7:
						current = 2*(w+h) - 5
					default:
						current = class + 2*(w+h) - 12
					}
				}
				if n := edgeNormalize(dims[0], dims[1], current); n != class {
					t.Fatalf("test inverse %v %d -> %d", dims, class, n)
				}
				for column := int32(0); column < 12; column++ {
					for _, mapped := range []uint32{0, 1, 2, 3, 4, 5, 6, 7, 8, 255, 0x80000000, 0x7fffffff, 0xffffffff} {
						specs = append(specs, legacy.PortTestEdgeMapSpec{Width: dims[0], Height: dims[1], Index: index, Current: current, Category: column, Normalized: class, Mapping: mapped})
					}
				}
			}
		}
	}
	t.Logf("%d dependent operations", len(specs))
	seed := 91
	got := legacy.PortTestEdgeMapping(seed, nil, specs)
	if len(got.Mapped) != len(specs) || !got.Restored || !got.TableUnchanged || !got.MappingUnchanged || !got.GuardsUnchanged || string(got.TableBefore) != string(got.TableAfterRestore) || len(got.MappingBefore) != len(got.MappingAfterRestore) {
		t.Fatalf("fixture state: table=%v mapping=%v guards=%v restored=%v", got.TableUnchanged, got.MappingUnchanged, got.GuardsUnchanged, got.Restored)
	}
	r := prand.New(seed)
	for i, s := range specs {
		class := edgeNormalize(s.Width, s.Height, s.Current)
		if class < 0 || class >= 12 {
			t.Fatalf("bad test class %d for %+v", class, s)
		}
		exp := s.Current
		ret := 1
		if s.Mapping == 255 {
			ret = 0
		} else {
			exp = edgeTransform(r, s.Width, s.Height, int32(s.Mapping))
		}
		g := got.Mapped[i]
		if g.Return != ret || g.Current != uint32(exp) || g.LogicIndex != r.Index() || g.OtherIndex != prand.New(seed+1).Index() || !g.RecordGuardsOK {
			t.Fatalf("mapped %d class=%d spec=%+v got=%+v want ret/current/index=%d/%x/%d", i, class, s, g, ret, uint32(exp), r.Index())
		}
	}
}
