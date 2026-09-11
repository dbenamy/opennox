//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"math"
	"os"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
)

func projectionSpec(kind string, line [4]float32, point [2]float32, length float32) legacy.PortTestLineProjectionSpec {
	w := make([]uint32, 12)
	for i := range w {
		w[i] = 0x91c3d57b
	}
	for i, v := range line {
		w[1+i] = math.Float32bits(v)
	}
	for i, v := range point {
		w[6+i] = math.Float32bits(v)
	}
	return legacy.PortTestLineProjectionSpec{Kind: kind, Words: w, LineOffset: 1, PointOffset: 6, OutputOffset: 9, LengthBits: math.Float32bits(length)}
}

func TestLineProjection(t *testing.T) {
	var specs []legacy.PortTestLineProjectionSpec
	var analytic [][3]uint32
	for _, line := range [][4]float32{{0, 2, 4, 2}, {4, 2, 0, 2}, {2, 0, 2, 4}, {2, 4, 2, 0}} {
		for _, px := range []float32{-2, 0, 1, 4, 6} {
			for _, py := range []float32{-2, 0, 1, 4, 6} {
				for _, kind := range []string{"clamp", "line"} {
					specs = append(specs, projectionSpec(kind, line, [2]float32{px, py}, 4))
					x, y := px, float32(2)
					v := px
					if line[0] == line[2] {
						x, y = 2, py
						v = py
					}
					ret := uint32(0)
					if kind == "line" {
						if v >= 0 && v <= 4 {
							ret = 1
						}
					} else {
						if v < 0 {
							v = 0
						}
						if v > 4 {
							v = 4
						}
						if line[0] == line[2] {
							y = v
						} else {
							x = v
						}
					}
					analytic = append(analytic, [3]uint32{ret, math.Float32bits(x), math.Float32bits(y)})
				}
			}
		}
	}
	edge := []uint32{0, 0x80000000, 1, 0x80000001, 0x00800000, 0x80800000, 0x3f800000, 0xbf800000, 0x7f7fffff, 0xff7fffff, 0x7f800000, 0xff800000, 0x7fc01234, 0xffc01234, 0x7f801234, 0xff801234}
	for _, kind := range []string{"clamp", "line"} {
		fields := []int{1, 2, 3, 4, 6, 7}
		if kind == "clamp" {
			fields = append(fields, -1)
		}
		for _, field := range fields {
			for _, v := range edge {
				for _, out := range []int{1, 2, 3, 4, 6, 7, 9} {
					s := projectionSpec(kind, [4]float32{0, 0, 4, 4}, [2]float32{2, 1}, 32)
					if field < 0 {
						s.LengthBits = v
					} else {
						s.Words[field] = v
					}
					s.OutputOffset = out
					specs = append(specs, s)
				}
			}
		}
	}
	seed := uint32(0x1729ca53)
	next := func() uint32 { seed = seed*1664525 + 1013904223; return seed }
	offsets := []int{1, 2, 3, 4, 6, 7, 9}
	for n := 0; n < 10000; n++ {
		words := make([]uint32, 12)
		for i := range words {
			words[i] = next()
		}
		length := next()
		for _, kind := range []string{"clamp", "line"} {
			s := legacy.PortTestLineProjectionSpec{Kind: kind, Words: words, LineOffset: 1, PointOffset: 6, OutputOffset: offsets[n%len(offsets)], LengthBits: length}
			if n%3 == 0 {
				s.PointOffset = 2
			}
			specs = append(specs, s)
		}
	}
	got := legacy.PortTestLineProjection(specs)
	actual := make([]byte, len(got)*12)
	for i, g := range got {
		if !g.GuardsOK || !g.OutsideOutputUnchanged || g.Return < 0 || g.Return > 1 {
			t.Fatalf("case%d memory/return failed: %+v", i, g)
		}
		values := [3]uint32{uint32(g.Return), g.Words[specs[i].OutputOffset], g.Words[specs[i].OutputOffset+1]}
		if i < len(analytic) && values != analytic[i] {
			t.Fatalf("analytic case%d got%x want%x spec=%+v", i, values, analytic[i], specs[i])
		}
		for j, v := range values {
			binary.LittleEndian.PutUint32(actual[i*12+j*4:], v)
		}
	}
	want, err := os.ReadFile("testdata/porting/line_projection.bin")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(actual, want) {
		for i := range got {
			if (i+1)*12 > len(want) || !bytes.Equal(actual[i*12:(i+1)*12], want[i*12:(i+1)*12]) {
				t.Fatalf("C baseline mismatch case%d actual=%x expected=%x spec=%+v", i, actual[i*12:(i+1)*12], want[min(i*12, len(want)):min((i+1)*12, len(want))], specs[i])
			}
		}
		t.Fatal("C baseline length mismatch")
	}
}
