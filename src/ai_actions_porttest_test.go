//go:build porttest

package opennox

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"testing"

	"github.com/opennox/libs/prand"
	"github.com/opennox/opennox/v1/legacy"
)

func aiSpec(action int) legacy.PortTestAIActionSpec {
	return legacy.PortTestAIActionSpec{Action: action, Seed: 12345, X: math.Float32bits(300), Y: math.Float32bits(400), TX: math.Float32bits(301), TY: math.Float32bits(401), Speed: math.Float32bits(2.3), Run: math.Float32bits(1.7), Stack: 1}
}

func aiActionCorpus() map[string][]legacy.PortTestAIActionSpec {
	groups := make(map[string][]legacy.PortTestAIActionSpec)
	for d := 0; d < 256; d++ {
		for target := 0; target < 256; target++ {
			sp := aiSpec(2)
			sp.Direction = uint16(d)
			sp.Desired = uint16((d*257 + target*33) & 65535)
			sp.TX = uint32(target)
			groups["angles"] = append(groups["angles"], sp)
		}
	}
	for d := 0; d < 65536; d++ {
		sp := aiSpec(3)
		sp.Direction = uint16(d ^ 0xaaaa)
		sp.Desired = uint16(d ^ 0x5555)
		sp.TX = uint32(d) | 0xabcd0000
		sp.Stack = int8(d % 24)
		groups["set-angle"] = append(groups["set-angle"], sp)
		sp = aiSpec(4)
		sp.Direction = uint16(d)
		sp.Seed = d + 1
		sp.Status = uint32(d&1) * 0x4000
		if d&2 != 0 {
			sp.Flags = 0x400
		}
		sp.X = math.Float32bits(float32(80 + d%5800))
		sp.Y = math.Float32bits(float32(80 + (d*19)%5800))
		groups["walk"] = append(groups["walk"], sp)
	}
	points := [][2]float32{{0, 0}, {1, 0}, {0, 1}, {-1, 0}, {0, -1}, {0.001, -0.001}, {5000, -5000}, {1e20, 1e-20}, {float32(math.Inf(1)), 0}, {math.Float32frombits(0x7fc12345), 1}}
	for a := 0; a < 2; a++ {
		for d := 0; d < 256; d++ {
			for i, p := range points {
				sp := aiSpec(a)
				sp.Direction = uint16(d)
				sp.Desired = uint16(d * 257)
				sp.TX = math.Float32bits(300 + p[0])
				sp.TY = math.Float32bits(400 + p[1])
				sp.Stack = int8(i % 3)
				sp.NilTarget = i%4 == 0
				groups["points"] = append(groups["points"], sp)
			}
		}
	}
	// Seed every confusion branch: draw gates, ranged/melee eligibility, running,
	// normal/empty/full stacks. Audio/combat callees remain their actual C code.
	for seed := 0; seed < 1000; seed++ {
		for variant := 0; variant < 8; variant++ {
			sp := aiSpec(5)
			sp.Seed = seed
			sp.Direction = uint16(seed % 256)
			sp.Flags = 0x400
			if variant&1 != 0 {
				sp.Melee = math.Float32bits(12)
			}
			if variant&2 != 0 {
				sp.Missile = 'x'
			}
			if variant&4 != 0 {
				sp.Status = 0x20
			}
			sp.Stack = []int8{0, 1, 23}[seed%3]
			groups["confusion"] = append(groups["confusion"], sp)
		}
	}
	for a := 0; a < 6; a++ {
		for mode := 1; mode < 4; mode++ {
			for _, stack := range []int8{0, 1, 23} {
				sp := aiSpec(a)
				sp.Mode = mode
				sp.Stack = stack
				groups["lifecycle"] = append(groups["lifecycle"], sp)
			}
		}
	}
	// Demonstrably detects spilling component*30 before adding position.
	sp := aiSpec(4)
	sp.Direction = uint16(-prand.New(sp.Seed).IntClamp(-20, 20))
	sp.X = 0x430a3931
	sp.Y = math.Float32bits(192)
	groups["precision"] = []legacy.PortTestAIActionSpec{sp}
	return groups
}

func TestAIActions(t *testing.T) {
	for name, specs := range aiActionCorpus() {
		t.Run(name, func(t *testing.T) {
			original, restored := legacy.PortTestAIActions(specs, false)
			got, restored2 := legacy.PortTestAIActions(specs, true)
			if !restored || !restored2 || len(got) != len(specs) {
				t.Fatal("fixture restoration/count")
			}
			for i, r := range got {
				if !r.GuardsOK || !r.ReadOnlyOK {
					t.Fatalf("case %d: guard/read-only state changed", i)
				}
				if !reflect.DeepEqual(r, original[i]) {
					t.Fatalf("case %d input=%+v\noriginal=%+v\nregistered=%+v", i, specs[i], original[i], r)
				}
				if name == "set-angle" && (r.Direction != uint16(specs[i].TX&255) || r.Desired != r.Direction) {
					t.Fatalf("set-angle %d: %+v", i, r)
				}
				if name == "walk" {
					draw := prand.New(specs[i].Seed).IntClamp(-20, 20)
					base := uint16(uint8(int(int16(specs[i].Direction)) + draw))
					if r.Direction != base && r.Direction != uint16(uint8(base+64)) {
						t.Fatalf("walk wrap %d", i)
					}
				}
			}
			if name == "precision" && got[0].Direction != 0 {
				t.Fatal("complete-expression rounding must reject probe before water lookup")
			}
			data, _ := json.Marshal(original)
			t.Logf("cases=%d original-state-sha256=%x", len(specs), sha256.Sum256(data))
			if hash := fmt.Sprintf("%x", sha256.Sum256(data)); hash != aiActionBaseline[name] {
				t.Fatalf("original state hash changed: %s", hash)
			}
		})
	}
}

var aiActionBaseline = map[string]string{
	"angles":    "81e19f73333afe70253185d4d356645c4fa6258a395f8693b4ce31a7cbebaace",
	"confusion": "23fd980c3261cd74934e59903679a048cc049494ece9784537f7573604d8896e",
	"lifecycle": "002b25ff52fe361e53b643d2c0def86f517a16614d46812020c800970e95de16",
	"points":    "02ecadfb9febf47d66aebf605df6b073de0139aa4644c7caf8d6658112c139d4",
	"precision": "c381b3ed1cac4cd325a267e4a38bd7f6fab087950b915d6885cc1eb4c6062c8c",
	"set-angle": "5115c1c9828b3b065d17e4da1596c2a67702aa3a657a63c8fa4c887a8244ef59",
	"walk":      "bd84ef08990300ca0a6f48f162e5dec7f27987c21b42b24ab03ce223a509244d",
}
