//go:build porttest

package opennox

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
)

func TestClientEffectsCurveSegments(t *testing.T) {
	type result struct {
		Spec     legacy.PortTestEffectsCurveSpec
		Segments []legacy.PortTestEffectsCurveSegment
	}
	var out []result
	points := [][4][2]int32{
		{{0, 0}, {64, 32}, {0, 0}, {0, 0}},
		{{-64, 32}, {64, -32}, {16, -8}, {-16, 8}},
		{{7, -11}, {7, -11}, {0, 0}, {0, 0}},
		{{0, 0}, {1, -1}, {-1, 1}, {1, -1}},
		{{32767, -32768}, {-32768, 32767}, {65535, -65536}, {-65536, 65535}},
		{{16777215, -16777217}, {16777217, -16777215}, {1, -1}, {-1, 1}},
		{{2147483647, -2147483648}, {-2147483648, 2147483647}, {2147483647, 2147483647}, {-2147483648, -2147483648}},
	}
	for pi, p := range points {
		for _, steps := range []int{-1, 0, 1, 2, 3, 4, 7, 8, 16, 31, 32} {
			for _, shift := range []float32{-0.25, 0, 0.125, 0.5, 1, 2} {
				sp := legacy.PortTestEffectsCurveSpec{Points: p, Steps: steps, Shift: shift, Token: int32(0x13570000 + len(out))}
				got := legacy.PortTestEffectsCurve(sp)
				wantCount := steps
				if wantCount < 0 {
					wantCount = 0
				}
				if len(got) != wantCount {
					t.Fatalf("curve callback count %d want%d", len(got), wantCount)
				}
				for i, g := range got {
					if g.Token != sp.Token {
						t.Fatal("curve userdata changed")
					}
					wantFrom := p[0]
					if i > 0 {
						wantFrom = got[i-1].To
					}
					if g.From != wantFrom {
						t.Fatal("curve segments disconnected")
					}
					if pi == 2 && g.To != p[0] {
						t.Fatal("constant Hermite curve changed position")
					}
				}
				if pi < 4 && shift == 0 && steps > 0 && steps&(steps-1) == 0 && got[len(got)-1].To != p[1] {
					t.Fatalf("Hermite curve pi%d steps%d shift%v endpoint %v want%v segments%v", pi, steps, shift, got[len(got)-1].To, p[1], got)
				}
				if pi == 0 && shift == 0 && steps == 2 && got[0].To != [2]int32{32, 16} {
					t.Fatal("symmetric Hermite midpoint")
				}
				out = append(out, result{sp, got})
			}
		}
	}
	data, err := json.Marshal(out)
	if err != nil {
		t.Fatal(err)
	}
	if prefix := os.Getenv("OPENNOX_CLIENT_EFFECTS_CAPTURE"); prefix != "" {
		if err := os.WriteFile(prefix+"-curve-segments.json", data, 0600); err != nil {
			t.Fatal(err)
		}
	}
	hash := fmt.Sprintf("%x", sha256.Sum256(data))
	t.Logf("curve-segments: %d cases %s", len(out), hash)
	const want = "1284f083a8ee4329cc30e12eff4a3e56f949c56ce2c90867f26babf763fc07e7"
	if hash != want {
		t.Fatalf("curve hash %s want repeated original-C capture %s", hash, want)
	}
}
