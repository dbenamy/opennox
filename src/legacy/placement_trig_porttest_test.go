//go:build porttest

package legacy

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strconv"
	"testing"

	"github.com/opennox/libs/prand"
)

func TestPlacementTrigRNGDomain(t *testing.T) {
	// Independently captured twice from the original 386 libc. The inventory owner
	// tries at most 64 angles; the generator's 32 angles are a prefix of this domain.
	const want = "0e09c10d0957707d8160e4ab4c5fbec4f008f77c428a66bba4a762418def66e9"
	h := sha256.New()
	var buf [16]byte
	for seed := 0; seed < 4096; seed++ {
		rng := prand.New(seed)
		angle := float32(rng.FloatClamp(float64(float32(-3.1415927)), float64(float32(3.1415927))))
		if rng.Index() != (seed+1)%4096 {
			t.Fatalf("seed %d consumed unexpected RNG index %d", seed, rng.Index())
		}
		for attempt := 0; attempt < 64; attempt++ {
			next := float64(angle) + 1.8849558
			angle = float32(next)
			if math.Abs(next) > 128 || math.Abs(float64(angle)) > 128 {
				t.Fatalf("placement angle escaped qualified range: seed=%d attempt=%d", seed, attempt)
			}
			binary.LittleEndian.PutUint64(buf[:8], math.Float64bits(portTestPlacementCos(next)))
			binary.LittleEndian.PutUint64(buf[8:], math.Float64bits(portTestPlacementSin(float64(angle))))
			h.Write(buf[:])
		}
	}
	if got := fmt.Sprintf("%x", h.Sum(nil)); got != want {
		t.Fatalf("all placement trig values = %s, want original %s", got, want)
	}
}

func TestPlacementTrigBoundaries(t *testing.T) {
	var rows []struct{ Input, Sin, Cos string }
	data, err := os.ReadFile("testdata/placement_trig_boundaries.json")
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(data, &rows); err != nil {
		t.Fatal(err)
	}
	if len(rows) != 70 {
		t.Fatalf("boundary capture cases=%d, want70", len(rows))
	}
	bits := func(s string) uint64 {
		v, err := strconv.ParseUint(s, 16, 64)
		if err != nil {
			t.Fatal(err)
		}
		return v
	}
	for _, r := range rows {
		x := math.Float64frombits(bits(r.Input))
		if got := math.Float64bits(portTestPlacementSin(x)); got != bits(r.Sin) {
			t.Errorf("sin(%s)=%016x want%s", r.Input, got, r.Sin)
		}
		if got := math.Float64bits(portTestPlacementCos(x)); got != bits(r.Cos) {
			t.Errorf("cos(%s)=%016x want%s", r.Input, got, r.Cos)
		}
	}
}
