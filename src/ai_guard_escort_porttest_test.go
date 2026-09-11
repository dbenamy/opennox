//go:build porttest

package opennox

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math"
	"testing"

	"github.com/opennox/libs/prand"
	"github.com/opennox/opennox/v1/legacy"
)

func guardEscortCorpus() []legacy.PortTestRoamSpec {
	var out []legacy.PortTestRoamSpec
	aggs := []uint32{0, math.Float32bits(.079999998) - 1, math.Float32bits(.079999998), math.Float32bits(.079999998) + 1, math.Float32bits(.33000001) - 1, math.Float32bits(.33000001), math.Float32bits(.33000001) + 1, math.Float32bits(.66000003) - 1, math.Float32bits(.66000003), math.Float32bits(.66000003) + 1, math.Float32bits(1)}
	for n := 0; n < 2048; n++ {
		for op := 0; op < 8; op++ {
			sp := legacy.PortTestRoamSpec{Op: 7, Seed: n + 1, Stack: int8(n % 24)}
			o := &legacy.PortTestRoamOwnerSpec{Frame: uint32(n % 128), FPS: 30, Aggression: aggs[n%len(aggs)], Status: []uint32{0, 0x200, 0x4000, 0x8000, 0xc000}[n%5], Enemy: n%4 == 0, X: math.Float32bits(float32(100 + n%20)), Y: math.Float32bits(100)}
			g := &legacy.PortTestGuardEscortSpec{Op: op, Kind: uint16(n%3 + 1), Direction: uint16(n % 256), Desired: uint16((n * 17) % 256), NetCode: uint32(n % 32), Sight: math.Float32bits(float32(n%5) * 100), Follow: math.Float32bits(float32(n%7) * 30), TX: math.Float32bits(float32(100 + n%300)), TY: math.Float32bits(float32(100+n%4) * 1), Heard: uint32(n % 2), HeardFrame: uint32(n%128) - uint32(n%100), Precheck: n%3 != 0, Seen: n%7 == 0, SeenPlayer: n%2 == 0, Previous: n%2 != 0, Unresolved: n%4 == 1, NoOwner: n%4 == 1, Name: []string{"**OWNER**", "missing", "**PLAYER**", ""}[n%4]}
			g.UseTerrain = n%2 == 0
			g.Water = n%3 == 0
			if n%13 == 0 {
				o.X = math.Float32bits(-1)
			} // real ray rejection at map boundary
			sp.Owner = o
			sp.GuardEscort = g
			out = append(out, sp)
		}
	}

	// Resolver branches are checked separately across sparse active-player slots.
	for seed := 1; seed <= 256; seed++ {
		for mode := 0; mode < 9; mode++ {
			sp := legacy.PortTestRoamSpec{Op: 7, Seed: seed, Stack: 1}
			sp.Owner = &legacy.PortTestRoamOwnerSpec{Frame: 123, FPS: 30, X: math.Float32bits(100), Y: math.Float32bits(100)}
			g := &legacy.PortTestGuardEscortSpec{Op: 6, Kind: 1, Players: mode % 4}
			switch mode {
			case 0:
				g.Name = "**OWNER**"
			case 1:
				g.Name = "**OWNER**"
				g.NoOwner = true
			case 2, 3, 4:
				g.Name = "**PLAYER**"
				g.Players = mode - 2
			case 5:
				g.Name = "escort-target"
				g.ScriptName = "zone:escort-target"
			case 6:
				g.Name = "zone:escort-target"
				g.ScriptName = g.Name
				g.Pending = true
			case 7:
				g.Name = "ESCORT-TARGET"
				g.ScriptName = "zone:escort-target"
			case 8:
				g.Name = "**PLAYER**"
				g.Players = 3
			}
			sp.GuardEscort = g
			out = append(out, sp)
		}
	}

	for _, frame := range []uint32{0, 1, 15, 16, 0x80000000, 0xfffffff0, 0xffffffff} {
		for _, fps := range []uint32{0, 1, 30, 60} {
			for _, age := range []uint32{0, 1, 89, 90, 180, 0xffffffff} {
				for _, stack := range []int8{0, 19, 20, 21, 22, 23} {
					for _, op := range []int{0, 1, 4, 5} {
						sp := legacy.PortTestRoamSpec{Op: 7, Seed: 7, Stack: stack}
						sp.Owner = &legacy.PortTestRoamOwnerSpec{Frame: frame, FPS: fps, X: math.Float32bits(100), Y: math.Float32bits(100)}
						sp.GuardEscort = &legacy.PortTestGuardEscortSpec{Op: op, Kind: 1, Heard: 1, HeardFrame: frame - age, Precheck: true, NetCode: 1, TX: math.Float32bits(100), TY: math.Float32bits(100)}
						out = append(out, sp)
					}
				}
			}
		}
	}
	for mode := 0; mode < 3; mode++ {
		sp := legacy.PortTestRoamSpec{Op: 7, Seed: 1, Stack: 1}
		o := &legacy.PortTestRoamOwnerSpec{Frame: 0, FPS: 30, X: math.Float32bits(100), Y: math.Float32bits(100)}
		g := &legacy.PortTestGuardEscortSpec{Kind: 1, TX: math.Float32bits(130), TY: math.Float32bits(100.001)}
		switch mode {
		case 0:
			o.X = math.Float32bits(108)
			o.Y = math.Float32bits(100.001)
		case 1:
			g.Op = 1
		case 2:
			o.Frame = 1
			o.Aggression = math.Float32bits(.5)
			o.Enemy = true
			g.Sight = math.Float32bits(30)
			g.TX = math.Float32bits(130) - 1
			g.TY = math.Float32bits(100.03)
		}
		sp.Owner = o
		sp.GuardEscort = g
		out = append(out, sp)
	}
	return out
}
func TestAIGuardEscortBaseline(t *testing.T) {
	specs := guardEscortCorpus()
	got := legacy.PortTestRoam(specs)
	for i, r := range got {
		if !r.Intact {
			t.Fatalf("case %d guard/read-only mutation", i)
		}
		sp := specs[i]
		if sp.GuardEscort.Op == 7 && r.Return != int((sp.Owner.Status>>9)&1) {
			t.Fatalf("case %d look-at-damager must report injury even when stack is full", i)
		}
		g := sp.GuardEscort
		if g.Op == 6 {
			rng := prand.New(sp.Seed)
			want := 0
			switch g.Name {
			case "**OWNER**":
				if !g.NoOwner {
					want = 100
				}
			case "**PLAYER**":
				ind := rng.IntClamp(0, g.Players-1)
				if g.Players != 0 {
					want = 200 + ind
				}
			default:
				if g.ScriptName != "" && (g.Name == g.ScriptName || g.Name == "escort-target") {
					want = 100
				}
			}
			if r.Return != want || r.Logic != rng.Index() {
				t.Fatalf("resolver case %d got=%+v want target=%d RNG=%d", i, r, want, rng.Index())
			}
		}
		if g.Op == 4 || g.Op == 5 {
			fresh := g.Heard != 0 && sp.Owner.Frame-g.HeardFrame < 3*sp.Owner.FPS
			if r.Return != int(boolToGuardInt(fresh)) {
				t.Fatalf("sound freshness case %d", i)
			}
		}
	}
	for i, want := range []int8{3, 3, 5} {
		if got[len(got)-3+i].Stack != want {
			t.Fatalf("precision case %d stack=%d want=%d", i, got[len(got)-3+i].Stack, want)
		}
	}
	data, _ := json.Marshal(got)
	hash := fmt.Sprintf("%x", sha256.Sum256(data))
	t.Logf("cases=%d complete-state-sha256=%s", len(got), hash)
	const baseline = "9133ad56a2008e306b5a1733670e0ff0eddfc7ae1d3e30e48ac48c3b884417a9"
	if baseline != "" && hash != baseline {
		t.Fatal("original C state differs", hash)
	}
}

func boolToGuardInt(v bool) int {
	if v {
		return 1
	}
	return 0
}

func TestAIGuardEscortRepeated(t *testing.T) {
	for _, op := range []int{0, 1} {
		sp := legacy.PortTestRoamSpec{Op: 7, Seed: 7, Stack: 1}
		sp.Owner = &legacy.PortTestRoamOwnerSpec{Repeat: 200000, Frame: 1, FPS: 30, X: math.Float32bits(100), Y: math.Float32bits(100)}
		sp.GuardEscort = &legacy.PortTestGuardEscortSpec{Op: op, Kind: 1, TX: math.Float32bits(100), TY: math.Float32bits(100)}
		got := legacy.PortTestRoam([]legacy.PortTestRoamSpec{sp})
		if !got[0].Intact {
			t.Fatal("guard/read-only state")
		}
		data, _ := json.Marshal(got)
		hash := fmt.Sprintf("%x", sha256.Sum256(data))
		t.Logf("op=%d updates=%d ns/update=%.1f complete-state-sha256=%s", op, sp.Owner.Repeat, float64(got[0].Nanos)/float64(sp.Owner.Repeat), hash)
		if hash != "d87dffe6830205f9ddab9e9d7f0c7fc2739f48c27731e88dc16c20b09a26d176" {
			t.Fatal("repeated original-C state differs", hash)
		}
	}
}
