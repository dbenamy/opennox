//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
	"unsafe"
)

func TestClientScreenEffectsParticleLifecycle(t *testing.T) {
	o := newObjectDrawingOwner(t)
	o.c.GUI = gui.New(o.c.Render())
	c := o.c
	env := legacy.PortTestNewScreenEnvironment()
	defer env.Restore()
	type result struct {
		Case, Step, Capacity int
		Input                [10]int32
		Created              int
		Pixels               string
		Particles            [][]uint32
		Render, Globals      []uint32
		Other                int
	}
	var out []result
	variants := [][9]int32{
		{48, 48, 0, 0, 0, 1, 0, 0, 0}, // default decay without gravity/timer
		{48, 48, 0, 0, 0, 4, 0, 2, 2},
		{48, 48, 0, 0, 1, 1, 1, 1, 1}, // grow and spawn children
		{48, 48, 0, 0, -1, 4, 1, 1, 1},
		{48, 48, 1, -1, 0, 1, 1, 2, 2}, // shrink to deletion
		{48, 48, -2, 2, 0, 4, 2, 2, 1},
		{48, 48, 0, 0, 0, 2, 255, 1, 1},
		{0, 48, 0, 0, 0, 2, 2, 2, 1},
		{1, 1, -1, 0, 0, 2, 2, 1, 1},
		{95, 95, 1, 1, 0, 3, 2, 2, 2},
		{96, 48, 0, 0, 0, 3, 0, 0, 1},
		{-1, 48, 0, 0, 0, 3, 0, 0, 2},
	}
	id := 0
	for _, capacity := range []int{0, 1, 2, 7} {
		for _, kind := range []int32{-1, 0, 1, 2, 3, 4, 5} {
			for _, v := range variants {
				for _, seed := range []uint32{1, 17, 31} {
					id++
					o.reset(seed, 120)
					env.Reset()
					c.GUI.ValXXX, c.GUI.ValYYY = 0, 0
					snapshot, free := legacy.PortTestEffectsScreenParticles(capacity)
					a := [10]int32{kind}
					copy(a[1:], v[:])
					created := 0
					for i := 0; i < capacity+2; i++ {
						args := a
						args[1] += int32(i)
						if legacy.PortTestScreenParticleCreate(args) != nil {
							created++
						}
					}
					for step := -1; step < 5; step++ {
						if step >= 0 {
							clear(o.pix.Pix)
							legacy.PortTestScreenParticlesDraw(c.Viewport())
						}
						globals := append(env.State(), uint32(c.GUI.ValXXX), uint32(c.GUI.ValYYY))
						out = append(out, result{id, step, capacity, a, created, effectsPixelHash(o.pix), snapshot(), append([]uint32(nil), unsafe.Slice((*uint32)(unsafe.Pointer(c.r.Data())), 264)...), globals, c.srv.Rand.Other.Index()})
					}
					free()
				}
			}
		}
	}
	effectsCapture(t, "screen-effects-particle-lifecycle", out, len(out), "d21c6e420b68c92a842a92948104828a13d69eb8e612c4c06039c105f7c83123")
}
func TestClientScreenEffectsParticleContracts(t *testing.T) {
	o := newObjectDrawingOwner(t)
	o.c.GUI = gui.New(o.c.Render())
	env := legacy.PortTestNewScreenEnvironment()
	defer env.Restore()
	legacy.PortTestScreenParticlesDraw(nil)
	if o.c.GUI.ValYYY != 0 {
		t.Fatal("nil viewport changed render flag")
	}
	snapshot, free := legacy.PortTestEffectsScreenParticles(1)
	defer free()
	a := [10]int32{0, 48, 48, 0, 0, 0, 2, 0, 0, 0}
	p := legacy.PortTestScreenParticleCreate(a)
	if p == nil || p.Field_40 != [4]byte{2, 3, 2, 3} {
		t.Fatal("stationary particle default decay")
	}
	p.Field_32 |= 0xa1b2c300
	a[0], a[5], a[7], a[9] = 4, 1, 1, 2
	q := legacy.PortTestScreenParticleCreate(a)
	if q != p || q.Field_32 != 0xa1b2c302 || len(snapshot()) != 1 {
		t.Fatal("full pool did not reuse tail with byte-width state")
	}
	// Actual final shrink tick deletes and returns before position integration.
	q.Field_40 = [4]byte{1, 1, 2, 1}
	if legacy.PortTestScreenParticleDraw(q, o.c.Viewport()) != 0 || len(snapshot()) != 0 {
		t.Fatal("final particle shrink did not release node")
	}
	a[1] = 0
	p = legacy.PortTestScreenParticleCreate(a)
	if legacy.PortTestScreenParticleDraw(p, o.c.Viewport()) != 0 || len(snapshot()) != 0 {
		t.Fatal("strict viewport edge did not release node")
	}
}
