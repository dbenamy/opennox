//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"math"
	"reflect"
	"testing"
)

func TestWorldGridTrigTables(t *testing.T) {
	r := legacy.PortTestWorldTables()
	if !r.GuardsOK || !r.WarmUnchanged || !reflect.DeepEqual(r.WarmReturns, []int8{1, 2, 127, -128, -1}) {
		t.Fatalf("table initialization guards/flag: guards%v unchanged%v returns%v", r.GuardsOK, r.WarmUnchanged, r.WarmReturns)
	}
	sineScale, acosScale := math.Float64frombits(r.Constants[1]), math.Float64frombits(r.Constants[0])
	if sineScale <= 0 || acosScale <= 0 || math.IsInf(sineScale, 0) || math.IsInf(acosScale, 0) {
		t.Fatalf("shipped constants %v %v", sineScale, acosScale)
	}
	if r.Sine[0] != 0 || int32(r.Sine[1024]) <= 0 || int32(r.Sine[3072]) >= 0 || r.Acos[0] == 0 || int8(r.Acos[len(r.Acos)-1]) != r.ColdReturn {
		t.Fatal("nontrivial table/return contract")
	}
	for i, v := range r.Sine {
		// A one-unit envelope allows libm endpoint rounding differences; exact C
		// integer words are also captured for the eventual native comparison.
		exact := math.Sin(float64(i)*0.0015339808) * sineScale
		if math.Abs(float64(int32(v))-exact) > 1.00001 {
			t.Fatalf("sine %d got%d expected%v", i, int32(v), exact)
		}
	}
	for i, v := range r.Acos {
		exact := math.Acos(float64(i)*0.00024414062-1) * acosScale
		if math.Abs(float64(int32(v))-exact) > 1.00001 {
			t.Fatalf("acos %d got%d expected%v", i, int32(v), exact)
		}
		if i > 0 && v > r.Acos[i-1] {
			t.Fatalf("acos not monotonic at%d", i)
		}
	}
	drawableStateCapture(t, "world-tables", r, "5fa9a38a75f46830c5615e3ccfbf30db32e451d5e536fa5789f560bf1dbb50c6")
}
