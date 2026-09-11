//go:build porttest

package opennox

import (
	"bytes"
	"math"
	"os"
	"reflect"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
)

func TestCollisionReflection(t *testing.T) {
	edge := []uint32{0, 0x80000000, 1, 0x80000001, 0x00800000, 0x80800000, 0x3f800000, 0xbf800000, 0x7f7fffff, 0xff7fffff, 0x7f800000, 0xff800000, 0x7fc01234, 0xffc01234, 0x7f801234, 0xff801234}
	var specs []legacy.PortTestReflectionSpec
	for _, x := range edge {
		for _, y := range edge {
			for _, vx := range edge {
				for _, vy := range edge {
					specs = append(specs, legacy.PortTestReflectionSpec{Words: []uint32{0x13579bdf, x, y, 0x2468ace0, vx, vy, 0xfedcba98}, NormalOffset: 1, VelocityOffset: 4})
				}
			}
		}
	}
	for _, x := range edge {
		for _, y := range edge {
			specs = append(specs, legacy.PortTestReflectionSpec{Words: []uint32{0x13579bdf, x, y, 0xfedcba98}, NormalOffset: 1, VelocityOffset: 1})
			for _, z := range edge {
				for _, off := range [][2]int{{1, 2}, {2, 1}} {
					specs = append(specs, legacy.PortTestReflectionSpec{Words: []uint32{0x13579bdf, x, y, z, 0xfedcba98}, NormalOffset: off[0], VelocityOffset: off[1]})
				}
			}
		}
	}
	nan := func(v uint32) bool { return v&0x7fffffff > 0x7f800000 }
	quiet := func(v uint32) uint32 {
		if nan(v) {
			return v | 0x00400000
		}
		return v
	}
	got := legacy.PortTestCollisionReflect(specs)
	for i, s := range specs {
		x, y := s.Words[s.NormalOffset], s.Words[s.NormalOffset+1]
		ax, ay := x&0x7fffffff, y&0x7fffffff
		invalid := nan(x) || nan(y) || (ax == 0 && ay == 0x7f800000) || (ay == 0 && ax == 0x7f800000)
		swap := !invalid && (ax == 0 || ay == 0 || (x^y)&0x80000000 != 0)
		vx, vy := s.Words[s.VelocityOffset], s.Words[s.VelocityOffset+1]
		want := append([]uint32(nil), s.Words...)
		if swap {
			want[s.VelocityOffset], want[s.VelocityOffset+1] = vy, quiet(vx)
		} else {
			want[s.VelocityOffset], want[s.VelocityOffset+1] = quiet(vy)^0x80000000, quiet(vx)^0x80000000
		}
		if !got[i].GuardsOK || !got[i].ReturnPointerSame || !reflect.DeepEqual(got[i].Words, want) {
			t.Fatalf("case%d normal=%x/%x velocity=%x/%x offsets=%d/%d got=%x want=%x guards=%v ptr=%v", i, x, y, vx, vy, s.NormalOffset, s.VelocityOffset, got[i].Words, want, got[i].GuardsOK, got[i].ReturnPointerSame)
		}
	}
}

func collisionDiamond(x, y, r, px, py float32) legacy.PortTestContainmentSpec {
	w := make([]uint32, 17)
	w[0], w[1] = math.Float32bits(x), math.Float32bits(y)
	w[7], w[8], w[9], w[10], w[11], w[12] = 0, math.Float32bits(-r), math.Float32bits(-r), 0, 0, math.Float32bits(r)
	w[13], w[14] = math.Float32bits(px), math.Float32bits(py)
	// The first five shape words and trailing words are deliberately non-finite;
	// none belongs to the six box coordinates read by this helper.
	for _, i := range []int{2, 3, 4, 5, 6, 15, 16} {
		w[i] = 0x7f801234
	}
	return legacy.PortTestContainmentSpec{Words: w, PositionOffset: 0, ShapeOffset: 2, PointOffset: 13}
}

func TestCollisionContainment(t *testing.T) {
	var specs []legacy.PortTestContainmentSpec
	var analytic []int
	for _, pos := range [][2]float32{{0, 0}, {100, -50}, {1 << 24, -(1 << 24)}} {
		for dx := -17; dx <= 17; dx++ {
			for dy := -17; dy <= 17; dy++ {
				px, py := pos[0]+float32(dx), pos[1]+float32(dy)
				specs = append(specs, collisionDiamond(pos[0], pos[1], 16, px, py))
				want := 0
				if math.Abs(float64(px)-float64(pos[0]))+math.Abs(float64(py)-float64(pos[1])) < 16 {
					want = 1
				}
				analytic = append(analytic, want)
			}
		}
	}
	for _, r := range []float32{math.SmallestNonzeroFloat32, 1, 16, 1 << 24, math.MaxFloat32} {
		for _, p := range []float32{0, r, math.Nextafter32(r, 0), math.Nextafter32(r, float32(math.Inf(1))), -r} {
			specs = append(specs, collisionDiamond(0, 0, r, p, 0))
		}
	}
	seed := uint32(0xac7e4921)
	next := func() uint32 { seed = seed*1664525 + 1013904223; return seed }
	for n := 0; n < 50000; n++ {
		w := make([]uint32, 17)
		for i := range w {
			w[i] = next()
		}
		s := legacy.PortTestContainmentSpec{Words: w, PositionOffset: 0, ShapeOffset: 2, PointOffset: 13}
		if n%4 == 0 {
			s.PositionOffset = 3
			s.PointOffset = 9
		} // Input aliasing is read-only.
		specs = append(specs, s)
	}
	got := legacy.PortTestCollisionContainment(specs)
	actual := make([]byte, (len(got)+7)/8)
	for i, g := range got {
		if !g.GuardsOK || !g.Unchanged || !reflect.DeepEqual(g.Words, specs[i].Words) || g.Return < 0 || g.Return > 1 {
			t.Fatalf("case%d changed storage or invalid result: %+v", i, g)
		}
		if i < len(analytic) && g.Return != analytic[i] {
			t.Fatalf("analytic case%d got%d want%d", i, g.Return, analytic[i])
		}
		actual[i/8] |= byte(g.Return) << uint(i%8)
	}
	want, err := os.ReadFile("testdata/porting/collision_containment.bin")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(actual, want) {
		for i, g := range got {
			if i/8 >= len(want) || g.Return != int(want[i/8]>>uint(i%8)&1) {
				t.Fatalf("C baseline mismatch case%d result%d spec=%+v", i, g.Return, specs[i])
			}
		}
		t.Fatal("C baseline size or trailing-bit mismatch")
	}
}
