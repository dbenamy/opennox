//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/legacy"
	"math"
	"testing"
)

func TestWorldCollisionsMassExchange(t *testing.T) {
	o := newReliableReportsOwner(t)
	a, b := &o.units[0], &o.units[1]
	type row struct {
		Name          string
		Before, After [4]uint32
	}
	var rows []row
	defer func() {
		spellbookCapture(t, "world-collisions-mass", rows, "cffe0c709a796471df9eda3fa6700d90909f0d61f813a6cadf7f08367641e523")
	}()
	for _, m1 := range []float32{0, 1.0 / 256, 1, 2, 1000, 1e20} {
		for _, m2 := range []float32{0, 1.0 / 256, 1, 2, 1000, 1e20} {
			if m1+m2 == 0 {
				continue
			}
			for _, v := range [][4]float32{{0, 0, 0, 0}, {1, 2, 3, 4}, {-3.5, 8, 10, -6.25}, {1000001, -0.00001, -999999, 0.00003}, {1e20, -1e20, -1e20, 1e20}} {
				name := fmt.Sprintf("mass%g/%g/velocity%v", m1, m2, v)
				t.Run(name, func(t *testing.T) {
					a.Mass = m1
					b.Mass = m2
					a.VelVec = types.Pointf{v[0], v[1]}
					b.VelVec = types.Pointf{v[2], v[3]}
					before := [4]uint32{}
					for i, x := range v {
						before[i] = math.Float32bits(x)
					}
					legacy.PortTestWorldCollision(0, a, b, nil)
					after := [4]float32{a.VelVec.X, a.VelVec.Y, b.VelVec.X, b.VelVec.Y}
					bits := [4]uint32{}
					for i, x := range after {
						bits[i] = math.Float32bits(x)
					}
					if m1 == m2 {
						if after != [4]float32{v[2], v[3], v[0], v[1]} {
							t.Fatalf("equal masses did not exchange velocity: %v", after)
						}
					}
					// Independently check momentum, scaling tolerance by each contribution to
					// allow cancellation and the production single-precision stores.
					for axis := 0; axis < 2; axis++ {
						p1 := float64(m1) * float64(v[axis])
						p2 := float64(m2) * float64(v[axis+2])
						got := float64(m1)*float64(after[axis]) + float64(m2)*float64(after[axis+2])
						scale := math.Abs(p1) + math.Abs(p2) + 1
						if math.Abs(got-p1-p2) > scale*1e-6 {
							t.Fatalf("momentum axis%d=%g want%g", axis, got, p1+p2)
						}
					}
					rows = append(rows, row{name, before, bits})
				})
			}
		}
	}
	// The helper explicitly accepts either missing collision object.
	a.VelVec = types.Pointf{3, 4}
	legacy.PortTestWorldCollision(0, a, nil, nil)
	legacy.PortTestWorldCollision(0, nil, a, nil)
	if a.VelVec != (types.Pointf{3, 4}) {
		t.Fatal("nil collision changed velocity")
	}
}
