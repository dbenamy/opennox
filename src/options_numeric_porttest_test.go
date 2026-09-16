//go:build porttest

package opennox

import (
	"fmt"
	"math"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
)

func TestOptionsNumericSliders(t *testing.T) {
	type record struct {
		Menu                     bool
		Event, ID, Value, Return int
		Sensitivity, Gamma       uint32
		Dirty                    bool
	}
	var records []record
	values := []int{-2147483648, -10000, -1000, -101, -1, 101, 1000, 10000, 2147483647}
	for v := 0; v <= 100; v++ {
		values = append(values, v)
	}
	for _, menu := range []bool{false, true} {
		t.Run(fmt.Sprint(menu), func(t *testing.T) {
			o := newOptionsOwner(t)
			for _, event := range []int{16393, 16396} {
				for _, id := range []int{316, 318, 999} {
					for _, v := range values {
						o.c.SetSensitivity(1)
						setGamma(1)
						configDirty = false
						ret := legacy.PortTestOptionsEvent(menu, o.root, event, o.controls[id], v)
						sens, gamma := o.c.GetSensitivity(), getGamma()
						if ret != 0 {
							t.Fatalf("unexpected return menu=%v event=%d id=%d value=%d: %d", menu, event, id, v, ret)
						}
						if event != 16393 || id != 318 {
							if sens != 1 {
								t.Fatalf("unrelated event changed sensitivity")
							}
						}
						if event != 16393 || id != 316 {
							if gamma != 1 || configDirty {
								t.Fatalf("unrelated event changed gamma/config")
							}
						}
						if event == 16393 && id == 316 {
							want := float32(0.1) + float32(3)*(float32(v)/100)
							if want < gammaMin {
								want = gammaMin
							}
							if want > gammaMax {
								want = gammaMax
							}
							if gamma != want || !configDirty {
								t.Fatalf("gamma %d got=%v want=%v dirty=%v", v, gamma, want, configDirty)
							}
						}
						if event == 16393 && id == 318 && v >= 0 && v <= 100 {
							want := math.Pow(10, float64(v)/50-1)
							if math.Abs(float64(sens)-want) > want*0.000001 {
								t.Fatalf("sensitivity curve %d got=%g want=%g", v, sens, want)
							}
							if v == 50 && sens != 1 {
								t.Fatal("midpoint must be exactly one")
							}
						}
						records = append(records, record{menu, event, id, v, ret, math.Float32bits(sens), math.Float32bits(gamma), configDirty})
					}
				}
			}
		})
	}
	spellbookCapture(t, "options-numeric", records, "608a984491b1c8f501b42d24379ac70fb29a124e7557d42a9b609af8f4064aa5")
}
