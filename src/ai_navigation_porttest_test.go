//go:build porttest

package opennox

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/opennox/opennox/v1/legacy"
	"math"
	"testing"
)

func navigationCorpus() []legacy.PortTestRoamSpec {
	bits := math.Float32bits
	var out []legacy.PortTestRoamSpec
	for n := 0; n < 2048; n++ {
		for op := 0; op < 11; op++ {
			o := &legacy.PortTestRoamOwnerSpec{Frame: []uint32{0, 1, 10, 11, 16, 123, 0xfffffff0, 0xffffffff}[n%8], FPS: []uint32{0, 1, 30, 60}[n%4], Aggression: []uint32{bits(0), bits(.5), bits(1)}[n%3], Status: []uint32{0, 0x20, 0x4000, 0x8000, 0xc000, 0x10000, 0x14000}[n%7], Enemy: n%3 == 0, X: bits(100), Y: bits(100)}
			g := &legacy.PortTestNavigationSpec{Op: op, Speed: []uint32{0, bits(.01) - 1, bits(.01), bits(1), bits(3)}[n%5], Multiplier: bits(1.5), TX: bits(103), TY: bits(104), Follow: bits(float32(n%4) * 10), Resume: []uint32{0, bits(.5) - 1, bits(.5), bits(.5) + 1, bits(1), 0x7fc12345}[n%6], FleeRange: bits(float32(n%4) * 10), Cur: uint16(n % 3 * 50), Max: []uint16{0, 100, 65535}[n%3], Previous: []uint32{1, 3, 56}[n%3], PathCount: uint32(n % 3), PathIndex: uint32(n % 3), PathStatus: uint32(n % 3), PathFrame: o.Frame - []uint32{0, 10, 11, 15, 16, 60, 61, 0xffffffff}[n%8], RetryFrame: o.Frame - uint32(n%200), OneShot: uint32(n % 3), Generator: n % 4, NoOwner: n%4 == 0, Food: n%2 == 0, TargetArg: n%2 == 0}
			if n%4 == 0 {
				g.GameFlags = 0x1000
			}
			if n%2 == 0 {
				o.Buffs = 1 << 29
			}
			// Gate deeper spell-policy calls, which remain outside this batch.
			if o.Status&0x20 != 0 {
				o.Buffs |= 1 << 29
			}
			if op == 2 {
				g.TX = bits(float32(100 + n%120))
				g.TY = bits(float32(100 + n%17))
				o.Buffs |= []uint32{0, 1 << 3, 1 << 5, 1 << 28}[n%4]
			}
			out = append(out, legacy.PortTestRoamSpec{Op: 8, Seed: n + 1, Stack: int8(n % 24), Owner: o, Navigation: g})
		}
	}
	// Lifecycle masks are independent of spatial state.
	for _, op := range []int{3, 4, 6} {
		for phase := 1; phase <= 3; phase++ {
			for mask := uint32(0); mask < 16; mask++ {
				out = append(out, legacy.PortTestRoamSpec{Op: 8, Seed: 1, Stack: 1, Owner: &legacy.PortTestRoamOwnerSpec{FPS: 30, Status: mask << 13, X: bits(100), Y: bits(100)}, Navigation: &legacy.PortTestNavigationSpec{Op: op, Phase: phase}})
			}
		}
	}
	// Numeric dodge corpus: threshold neighbors, cancellation, signed zero and nonfinite values.
	vals := []uint32{0, 0x80000000, bits(.01) - 1, bits(.01), bits(1), bits(7.9999) - 1, bits(7.9999), bits(7.9999) + 1, bits(8), bits(100), 0x7f800000, 0xff800000, 0x7fc12345, 0x3f812345, 0x49800001}
	for i, x := range vals {
		for j, y := range vals {
			for k, m := range vals {
				out = append(out, legacy.PortTestRoamSpec{Op: 8, Seed: 1, Stack: 1, Owner: &legacy.PortTestRoamOwnerSpec{FPS: 30, X: 0, Y: 0}, Navigation: &legacy.PortTestNavigationSpec{Op: 2, Speed: []uint32{bits(1), bits(.01), 0x7fc12345}[k%3], Multiplier: m, TX: x, TY: y, Cur: uint16(i), Max: uint16(j)}})
			}
		}
	}
	// Cross timer boundaries independently from path status/count and stack capacity.
	for _, op := range []int{0, 1, 3, 4, 6, 10} {
		for _, frame := range []uint32{0, 123, 0xfffffff0} {
			for _, fps := range []uint32{0, 1, 30} {
				for _, age := range []uint32{0, 10, 11, 15, 16, 60, 61, 149, 150, 0xffffffff} {
					for mode := 0; mode < 4; mode++ {
						for _, stack := range []int8{1, 20, 22, 23} {
							o := &legacy.PortTestRoamOwnerSpec{Frame: frame, FPS: fps, X: bits(100), Y: bits(100), PathMode: byte(mode), Buffs: 1 << 29}
							g := &legacy.PortTestNavigationSpec{Op: op, Speed: bits(1), Multiplier: bits(1.5), TX: bits(103), TY: bits(104), PathStatus: uint32(mode % 3), PathFrame: frame - age, RetryFrame: frame - age, Generator: mode, Previous: 3, Cur: 1, Max: 100, Resume: bits(.5), Food: mode%2 == 0}
							if op == 0 || op == 1 || op == 4 {
								if mode%2 == 0 {
									g.TX = bits(200)
								}
								g.Follow = bits(10)
							}
							out = append(out, legacy.PortTestRoamSpec{Op: 8, Seed: 7, Stack: stack, Owner: o, Navigation: g})
						}
					}
				}
			}
		}
	}
	for n := 0; n < 4096; n++ {
		x := bits(float32(n%301) / 7)
		y := bits(float32(n%131) / 13)
		sp := legacy.PortTestRoamSpec{Op: 8, Seed: 1, Stack: 1, Owner: &legacy.PortTestRoamOwnerSpec{FPS: 30}, Navigation: &legacy.PortTestNavigationSpec{Op: 2, TX: x, TY: y, Speed: bits(float32(n%117+1) / 17), Multiplier: bits(float32(n%37+1) / 19)}}
		if n == 0 {
			sp.Navigation.TX = 0x40ffff2e
			sp.Navigation.TY = 0
		}
		if n == 1 {
			sp.Navigation.TX = 0x3dcccccd
			sp.Navigation.TY = 0x4101999a
			sp.Navigation.Speed = 0x3dcccccd
			sp.Navigation.Multiplier = 0x3dcccccd
		}
		if n == 2 {
			sp.Navigation.Op = 0
			sp.Navigation.TX = 0x41d66cea
			sp.Navigation.TY = 0
			sp.Navigation.Follow = 0x410ef347
			sp.Navigation.Previous = 3
			sp.Owner.Status = 0x4000
			sp.Owner.PathMode = 2
		}
		out = append(out, sp)
	}
	return out
}

func TestAINavigationBaseline(t *testing.T) {
	specs := navigationCorpus()
	got := legacy.PortTestRoam(specs)
	for i, r := range got {
		if !r.Intact {
			t.Fatalf("case %d op=%d guard/read-only mutation", i, specs[i].Navigation.Op)
		}
		sp := specs[i]
		g := sp.Navigation
		if g.Op == 10 {
			selected := uint32(0)
			for j := 0; j+1 < len(r.Trace); j++ {
				if r.Trace[j] == 2489452 {
					selected = r.Trace[j+1]
				}
			}
			want := uint32(0)
			if g.Food {
				want = 100
			}
			if selected != want {
				t.Fatalf("case %d food=%d want=%d", i, selected, want)
			}
		}
		ratio := 1.0
		if g.Max != 0 {
			ratio = float64(g.Cur) / float64(g.Max)
		}
		resume := ratio >= float64(math.Float32frombits(g.Resume))
		cast := sp.Owner.Status&0x20 != 0 && sp.Owner.Buffs&(1<<29) != 0
		want := -1
		switch g.Op {
		case 7:
			want = 0
			if resume {
				want = 1
			}
		case 8:
			want = 0
			if !resume || cast {
				want = 1
			}
		case 9:
			want = 0
			if cast {
				want = 1
			}
		}
		if want >= 0 && r.Return != want {
			t.Fatalf("case %d policy=%d want=%d", i, r.Return, want)
		}
	}
	data, _ := json.Marshal(got)
	hash := fmt.Sprintf("%x", sha256.Sum256(data))
	t.Logf("cases=%d complete-state-sha256=%s", len(got), hash)
	const baseline = "7996e036ac64ebdf8c2786f48052d3f3dadfd0f5fab36c2e6113907a17134f3d"
	if baseline != "" && hash != baseline {
		t.Fatal("original C state differs", hash)
	}
}
