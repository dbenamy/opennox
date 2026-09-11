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

func roamOwnerCorpus() []legacy.PortTestRoamSpec {
	var out []legacy.PortTestRoamSpec
	aggs := []uint32{0, math.Float32bits(.33000001) - 1, math.Float32bits(.33000001), math.Float32bits(.33000001) + 1, math.Float32bits(.66000003) - 1, math.Float32bits(.66000003), math.Float32bits(.66000003) + 1, math.Float32bits(1), 0x7fc12345}
	positions := []uint32{math.Float32bits(92) - 1, math.Float32bits(92), math.Float32bits(92) + 1, math.Float32bits(100), math.Float32bits(108) - 1, math.Float32bits(108), math.Float32bits(108) + 1, math.Float32bits(140)}
	for n := 0; n < 1024; n++ {
		for mode := 0; mode < 12; mode++ {
			sp := legacy.PortTestRoamSpec{Op: 6, Seed: n + 1, Index: byte(n % 16), Count: byte(n % 5), Mask: 0x80, Stack: int8([]int{0, 1, 21, 22, 23}[n%5])}
			for i := 1; i < 34; i++ {
				sp.Enabled[i] = true
				sp.Flags[i] = 255
			}
			for i := range sp.History {
				if (n+i)%3 != 0 {
					sp.History[i] = byte((n+i)%4 + 1)
				}
			}
			sp.Neighbors = [32]byte{2, 3, 2, 4}
			o := &legacy.PortTestRoamOwnerSpec{Frame: uint32(32 + n%33), FPS: 30, Aggression: aggs[n%len(aggs)], Current: 1, X: math.Float32bits(100), Y: math.Float32bits(100), WX: positions[n%len(positions)], WY: math.Float32bits(100)}
			switch mode {
			case 1:
				o.PathMode = 1
			case 2:
				o.ExistingPath = true
			case 3:
				o.PathMode = 3
			case 4:
				o.Current = 0
			case 5:
				o.Current = 0
				o.Register = true
				sp.Count = 2
			case 6:
				o.Current = 0
				o.Fallback = 1
			case 7:
				o.Enemy = true
			case 8:
				o.Frame = 32
				o.Status = 0x20000
				o.Buffs = 1
			case 9:
				sp.Enabled[1] = false
				o.Register = true
			case 10:
				o.PathMode = 1
				clear(sp.History[:])
				sp.Mask = byte(n)
			case 11:
				sp.Count = 0
				o.WX = math.Float32bits(100)
			}
			sp.Owner = o
			out = append(out, sp)
		}
	}

	// The double sum is greater than 64; rounding the sum to float32 loses dy².
	sp := legacy.PortTestRoamSpec{Op: 6, Seed: 1, Mask: 128, Stack: 1}
	sp.Enabled[1] = true
	sp.Flags[1] = 255
	sp.Owner = &legacy.PortTestRoamOwnerSpec{Frame: 123, FPS: 30, Current: 1, X: math.Float32bits(100), Y: math.Float32bits(100), WX: math.Float32bits(108), WY: math.Float32bits(100.001)}
	out = append(out, sp)
	return out
}
func TestAIRoamOwnerBaseline(t *testing.T) {
	specs := roamOwnerCorpus()
	got := legacy.PortTestRoam(specs)
	for i, r := range got {
		if !r.Intact {
			t.Fatalf("case %d guard/read-only mutation: %+v", i, specs[i])
		}
		sp := specs[i]
		o := sp.Owner
		aggression := float64(math.Float32frombits(o.Aggression))
		reacts := aggression > .66000003 || (aggression < .66000003 && aggression > .33000001)
		marked := false
		for j := 0; j < len(r.Changes); j += 2 {
			if r.Changes[j] == 4096+516 && r.Changes[j+1] == 1 {
				marked = true
			}
		}
		if marked != reacts {
			t.Fatalf("case %d retaliation threshold: bits=%x marked=%t want=%t", i, o.Aggression, marked, reacts)
		}
		if o.Enemy && aggression > .66000003 {
			if len(r.Trace) != 0 || r.Stack != min(sp.Stack+1, 23) || r.Changed != (sp.Stack < 23) {
				t.Fatalf("case %d fight interruption: %+v", i, r)
			}
		}
		if o.Status&0x20000 != 0 && o.Buffs&1 != 0 && o.Frame&31 == 0 && !o.Enemy {
			rng := prand.New(sp.Seed)
			if rng.IntClamp(0, 100) < 10 {
				if sp.Stack <= 20 {
					rng.IntClamp(3, 10)
				}
				if len(r.Trace) != 0 || r.Stack != min(sp.Stack+3, 23) || r.Logic != rng.Index() {
					t.Fatalf("case %d idle push/capacity/RNG: %+v", i, r)
				}
			}
		}
	}
	last := got[len(got)-1]
	if last.Stack != 1 || len(last.Trace) != 3 {
		t.Fatal("arrival must preserve double sum above 64", last)
	}
	data, _ := json.Marshal(got)
	hash := fmt.Sprintf("%x", sha256.Sum256(data))
	t.Logf("cases=%d complete-state-sha256=%s", len(got), hash)
	const baseline = "3adafc8c4c9214e93e850c3eb3df43727ffae81be105f9496e0406dc31ee2856"
	if baseline != "" && hash != baseline {
		t.Fatalf("original-C state differs: %s", hash)
	}
}

func TestAIRoamOwnerRepeated(t *testing.T) {
	for _, mode := range []byte{0, 3} {
		sp := legacy.PortTestRoamSpec{Op: 6, Seed: 7, Mask: 128, Stack: 1}
		sp.Enabled[1] = true
		sp.Flags[1] = 255
		sp.Owner = &legacy.PortTestRoamOwnerSpec{Repeat: 200000, Frame: 123, FPS: 30, Current: 1, PathMode: mode, X: math.Float32bits(100), Y: math.Float32bits(100), WX: math.Float32bits(140), WY: math.Float32bits(100)}
		got := legacy.PortTestRoam([]legacy.PortTestRoamSpec{sp})
		if !got[0].Intact {
			t.Fatal("read-only/guard state")
		}
		data, _ := json.Marshal(got)
		hash := fmt.Sprintf("%x", sha256.Sum256(data))
		t.Logf("mode=%d updates=%d ns/update=%.1f complete-state-sha256=%s", mode, sp.Owner.Repeat, float64(got[0].Nanos)/float64(sp.Owner.Repeat), hash)
		want := map[byte]string{0: "94e8490623daefaa7c2d503e2cf0c443286e88a65ccabd465c6f520a3eec5725", 3: "758c1268461630573729d1d03fbd000992a2be97def593ad34dad6ffa062052f"}
		if hash != want[mode] {
			t.Fatal("repeated state differs from original C", hash)
		}
	}
}
