//go:build porttest

package opennox

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math/rand"
	"testing"

	"github.com/opennox/libs/prand"
	"github.com/opennox/opennox/v1/legacy"
)

func roamCorpus() []legacy.PortTestRoamSpec {
	rng := rand.New(rand.NewSource(545790))
	var out []legacy.PortTestRoamSpec
	for n := 0; n < 16384; n++ {
		sp := legacy.PortTestRoamSpec{Seed: n + 1, Index: byte(n % 16), Insert: byte(n % 34), Count: byte(n % 33), Mask: byte(n % 256), Stack: int8(n % 24)}
		for i := range sp.History {
			sp.History[i] = byte(rng.Intn(34))
		}
		for i := range sp.Neighbors {
			sp.Neighbors[i] = byte(rng.Intn(34))
		}
		for i := 1; i < 34; i++ {
			sp.Flags[i] = byte(rng.Intn(256))
			sp.Enabled[i] = rng.Intn(4) != 0
		}
		// Deliberately exercise the history-only fallback and empty history paths.
		if n%4 == 0 {
			for i := range sp.Neighbors {
				sp.Neighbors[i] = sp.History[i%16]
			}
		}
		if n%4 == 1 {
			clear(sp.History[:])
		}
		for op := 0; op < 6; op++ {
			sp.Op = op
			out = append(out, sp)
		}
	}
	return out
}

func roamExpected(sp legacy.PortTestRoamSpec) (history [16]byte, index uint32, arg uint32, field2 uint32, ret int, stack int8, logic int, changed bool) {
	history, index, arg, field2, stack = sp.History, uint32(sp.Index), 33, 99, sp.Stack
	rng := prand.New(sp.Seed)
	eligible := func(id byte) bool { return id != 0 && sp.Enabled[id] && sp.Flags[id]&sp.Mask != 0 }
	insert := func(id byte) {
		index = (index + 1) % 16
		history[index] = id
		for i, v := range history {
			if uint32(i) != index && v == id {
				history[i] = 0
			}
		}
	}
	choose := func() byte {
		var candidates []byte
		for _, id := range sp.Neighbors[:sp.Count] {
			seen := false
			for _, old := range history {
				if old == id {
					seen = true
					break
				}
			}
			if eligible(id) && !seen {
				candidates = append(candidates, id)
			}
		}
		if len(candidates) > 0 {
			return candidates[rng.IntClamp(0, len(candidates)-1)]
		}
		for k := 1; k <= 16; k++ {
			id := history[(int(index)+k)%16]
			if !eligible(id) {
				continue
			}
			for _, neighbor := range sp.Neighbors[:sp.Count] {
				if neighbor == id {
					return id
				}
			}
		}
		return 0
	}
	switch sp.Op {
	case 0:
		history[0] = 0
		index = 0
	case 1:
		arg = 0
	case 2:
		insert(sp.Insert)
	case 3:
		for k := 1; k < 16; k++ {
			id := history[(int(index)-k+16)%16]
			if eligible(id) {
				ret = int(id)
				break
			}
		}
	case 4:
		ret = int(choose())
	case 5:
		field2 = 0
		if sp.Count != 0 {
			arg = uint32(choose())
		}
		if sp.Count != 0 && arg != 0 {
			insert(byte(arg))
			ret = 1
		} else {
			changed = true
			index = 0
			if stack > 0 {
				stack--
			}
		}
	}
	logic = rng.Index()
	return
}

func TestAIRoamBaseline(t *testing.T) {
	specs := roamCorpus()
	got := legacy.PortTestRoam(specs)
	for i, r := range got {
		h, index, arg, f, ret, stack, logic, changed := roamExpected(specs[i])
		if !r.Intact || r.History != h || r.Index != index || r.Arg != arg || r.Field2 != f || r.Return != ret || r.Stack != stack || r.Logic != logic || r.Other != prand.New(specs[i].Seed+1).Index() || r.Changed != changed {
			t.Fatalf("case %d spec=%+v\ngot=%+v\nwant history=%v index=%d arg=%d field2=%d ret=%d stack=%d logic=%d changed=%t", i, specs[i], r, h, index, arg, f, ret, stack, logic, changed)
		}
	}
	data, _ := json.Marshal(got)
	hash := fmt.Sprintf("%x", sha256.Sum256(data))
	t.Logf("cases=%d complete-state-sha256=%s", len(got), hash)
	const baseline = "277a2ccf25e87d42404880e8fabc5276d10fec2c85fcdf4f472028c7c61c4628"
	if baseline != "" && hash != baseline {
		t.Fatalf("state differs from original C: %s", hash)
	}
}
