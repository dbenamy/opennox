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

func pathExecutionCorpus() []legacy.PortTestRoamSpec {
	bits := math.Float32bits
	point := func(x, y float32) [2]uint32 { return [2]uint32{bits(x), bits(y)} }
	var out []legacy.PortTestRoamSpec
	for n := 0; n < 4096; n++ {
		for op := 0; op < 4; op++ {
			sp := legacy.PortTestRoamSpec{Op: 9, Seed: n + 1, Stack: 1}
			o := &legacy.PortTestRoamOwnerSpec{Frame: []uint32{0, 10, 11, 123, 0xfffffff0, 0xffffffff}[n%6], FPS: 30, Status: []uint32{0, 0x4000}[n%2], X: bits(100), Y: bits(100), WX: bits(120), WY: bits(100), PathMode: []byte{1, 4, 5, 6}[n%4]}
			p := &legacy.PortTestPathSpec{Op: op, Wall: n % 3, Count: uint32(n % 5), Index: uint32(n % 5), Status: uint32(n % 3), Frame: o.Frame - []uint32{0, 10, 11, 60, 0xffffffff}[n%5], Speed: bits(float32(n%11+1) / 7), Multiplier: bits(float32(n%7+1) / 3), Target: point(float32(105+n%150), 100), Cache: point(float32(100+n%90), 100), WaypointCache: point(float32(n%151), 100), ObjFlags: []uint32{0, 0x4000}[n%2]}
			for i := range p.Path {
				p.Path[i] = point(float32(100+(n+i*23)%180), float32(100+(n+i*17)%21-10))
			}
			// Mix near points, visible far points, and points beyond the real wall.
			if n%4 == 0 {
				p.Path[0] = point(103, 104)
			}
			if n%7 == 0 && p.Count != 0 {
				p.Path[p.Count-1] = point(100, 100)
			}
			if n%3 == 0 {
				p.Index = 0
			}
			if n%4 == 0 {
				p.WaypointCount = 2
				p.WaypointIndex = uint32(n % 2)
				p.Waypoints[0], p.Waypoints[1] = 1, 2
			}
			if op == 3 {
				p.Start, p.End = byte(n%3), byte(n%2)
			} // endpoint failures/same node; graph success has its own fixture.
			sp.Owner, sp.Path = o, p
			out = append(out, sp)
		}
	}
	// Cross exact distance thresholds, timer wrap and ray/definition flags.
	for _, delta := range []uint32{bits(8) - 1, bits(8), bits(8) + 1, bits(50) - 1, bits(50), bits(50) + 1, bits(100) - 1, bits(100), bits(100) + 1} {
		for wall := 0; wall < 3; wall++ {
			for _, age := range []uint32{0, 10, 11, 0xffffffff} {
				sp := legacy.PortTestRoamSpec{Op: 9, Seed: 1, Stack: 1, Owner: &legacy.PortTestRoamOwnerSpec{Frame: 0, FPS: 30, X: 0, Y: bits(100), WX: bits(120), WY: bits(100), PathMode: 4}}
				p := &legacy.PortTestPathSpec{Op: 1, Wall: wall, Count: 2, Status: 0, Frame: 0 - age, Speed: bits(1), Multiplier: bits(1.3), Target: [2]uint32{delta, bits(100)}, Cache: point(0, 100), WaypointCache: point(0, 100)}
				p.Path[0], p.Path[1] = point(20, 100), point(120, 100)
				sp.Path = p
				out = append(out, sp)
			}
		}
	}
	for end := byte(1); end <= 33; end++ {
		for _, op := range []int{1, 3} {
			sp := legacy.PortTestRoamSpec{Op: 9, Seed: 1, Stack: 1, Owner: &legacy.PortTestRoamOwnerSpec{Frame: 123, FPS: 30, X: bits(100), Y: bits(100), WX: bits(120), WY: bits(100), PathMode: 4}, Path: &legacy.PortTestPathSpec{Op: op, Graph: true, Wall: 1, Start: 1, End: end, Speed: bits(1), Multiplier: bits(1.3), Target: point(200, 100)}}
			out = append(out, sp)
		}
	}
	for n := 0; n < 1024; n++ {
		sp := legacy.PortTestRoamSpec{Op: 9, Seed: 1, Stack: 1, Owner: &legacy.PortTestRoamOwnerSpec{Frame: 123, FPS: 30, Status: uint32(n%2) * 0x4000, X: []uint32{bits(.1), bits(100), bits(100.1)}[n%3], Y: bits(.03)}, Path: &legacy.PortTestPathSpec{Op: 0, Count: 2, Speed: []uint32{bits(1.1), bits(.1), 0, 0x80000000, 0x7f800000, 0x7fc12345}[n%6], Multiplier: bits(float32(n%37+1) / 19)}}
		sp.Path.Path[0], sp.Path.Path[1] = point(200+float32(n%51)/7, 110+float32(n%41)/11), point(500, 500)
		out = append(out, sp)
	}
	for mode := 0; mode < 3; mode++ {
		sp := legacy.PortTestRoamSpec{Op: 9, Seed: 1, Stack: 1, Owner: &legacy.PortTestRoamOwnerSpec{FPS: 30, Status: 0x4000}, Path: &legacy.PortTestPathSpec{Op: 0, Count: 2, Speed: bits(.1), Multiplier: bits(.1)}}
		switch mode {
		case 0:
			sp.Path.Path[0] = [2]uint32{0xc0bd2f55, 0xc0ac7858}
			sp.Path.Path[1] = [2]uint32{0x40076d3a, 0xc0f6e255}
		case 1:
			sp.Path.Path[0] = [2]uint32{0x3dcccccd, 0x4101999a}
			sp.Path.Path[1] = point(100, 100)
		case 2:
			sp.Path.Op = 1
			sp.Path.Count = 0
			sp.Path.Target = [2]uint32{0x40ffae15, 0}
			sp.Owner.PathMode = 1
		}
		out = append(out, sp)
	}
	return out
}
func TestAIPathExecutionBaseline(t *testing.T) {
	specs := pathExecutionCorpus()
	got := legacy.PortTestRoam(specs)
	for i, r := range got {
		want := uint32(1)
		if specs[i].Path.Wall == 1 {
			want = 0
		}
		if len(r.Trace) < 2 || r.Trace[0] != 5 || r.Trace[1] != want {
			t.Fatalf("wall fixture %d trace=%v", i, r.Trace)
		}
		if specs[i].Path.Op == 3 && specs[i].Path.Graph {
			want := min(int(specs[i].Path.End), 16)
			if specs[i].Path.End == 1 {
				want = 0
			}
			if r.Return != want {
				t.Fatalf("owner graph %d count=%d want=%d", i, r.Return, want)
			}
		}
		if !r.Intact {
			t.Fatalf("case %d op=%d guards/read-only mutation", i, specs[i].Path.Op)
		}
	}
	last := len(got) - 3
	chosen := uint32(0)
	for j := 0; j < len(got[last].Changes); j += 2 {
		if got[last].Changes[j] == 4096+268 {
			chosen = got[last].Changes[j+1]
		}
	}
	if chosen != 1 {
		t.Fatal("float32 nearest threshold must choose the second candidate")
	}
	if got[last+2].Return != 0 {
		t.Fatal("unrounded sqrt+epsilon above eight must continue path processing")
	}
	b, _ := json.Marshal(got)
	hash := fmt.Sprintf("%x", sha256.Sum256(b))
	t.Logf("cases=%d complete-state-sha256=%s", len(got), hash)
	const baseline = "f2ae3862c6ea077add1164716de01de0ec26b2e5ccf0757517d3e3daf15c8195"
	if baseline != "" && hash != baseline {
		t.Fatal("original C state differs", hash)
	}
}

func TestAIPathRepeated(t *testing.T) {
	bits := math.Float32bits
	for _, op := range []int{0, 1} {
		sp := legacy.PortTestRoamSpec{Op: 9, Seed: 1, Stack: 1, Owner: &legacy.PortTestRoamOwnerSpec{Repeat: 200000, Frame: 123, FPS: 30, X: bits(100), Y: bits(100)}, Path: &legacy.PortTestPathSpec{Op: op, Count: 2, Frame: 123, Speed: bits(1), Multiplier: bits(1.3), Target: [2]uint32{bits(200), bits(100)}, Cache: [2]uint32{bits(200), bits(100)}}}
		sp.Path.Path[0], sp.Path.Path[1] = [2]uint32{bits(200), bits(100)}, [2]uint32{bits(500), bits(500)}
		got := legacy.PortTestRoam([]legacy.PortTestRoamSpec{sp})
		if !got[0].Intact {
			t.Fatal("repeated path guard/read-only state")
		}
		b, _ := json.Marshal(got)
		hash := fmt.Sprintf("%x", sha256.Sum256(b))
		t.Logf("op=%d updates=%d ns/update=%.1f complete-state-sha256=%s", op, sp.Owner.Repeat, float64(got[0].Nanos)/float64(sp.Owner.Repeat), hash)
		const baseline = "64e97f956b16ea81f659b3bf00ec23fdc11038c26102a63341bfbbefe8f4af4d"
		if baseline != "" && hash != baseline {
			t.Fatal("original C repeated path state differs", hash)
		}
	}
}
