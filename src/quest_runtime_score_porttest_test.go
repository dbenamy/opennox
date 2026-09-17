//go:build porttest

package opennox

import (
	"fmt"
	"math"
	"math/big"
	"testing"
)

func TestQuestRuntimeScorePrimitive(t *testing.T) {
	_ = newQuestRuntimeOwner(t)
	type row struct {
		Input  [4]uint32
		Result uint64
	}
	var rows []row
	defer func() {
		spellbookCapture(t, "quest-runtime-score-primitive", rows, "922f6e231d6bcf5740ad68ac57851ad9c78ccae86267e7b7701e4637e0490e26")
	}()
	values := []uint32{0, 1, 10, 65535, 1000000, 0xffffffff}
	for _, a := range values {
		for _, b := range values {
			for _, c := range values {
				for _, stage := range []uint32{0, 1, 2, 3, 4, 10, 1000} {
					input := [4]uint32{a, b, c, stage}
					t.Run(fmt.Sprint(input), func(t *testing.T) {
						got := questRuntimeCall("sub_4D66E0", nil, input[:]...)
						if a == 0 && b == 0 && c == 0 && got != 0 {
							t.Fatal("zero contributions produced score", got)
						}
						if stage == 1 {
							// Exact rational weighting, independent of C/libm operation order; stage1
							// makes the difficulty multiplier one for every finite exponent.
							numerator := new(big.Int).SetUint64(uint64(a)*100 + uint64(b)*350 + uint64(c))
							exact := new(big.Rat).SetFrac(numerator, big.NewInt(10))
							f, _ := new(big.Float).SetPrec(128).SetRat(exact).Float32()
							if float64(f) < math.MaxInt32 {
								want := uint64(f)
								if got != want {
									t.Fatalf("stage1 score%d want%d", got, want)
								}
							}
						}
						rows = append(rows, row{input, got})
					})
				}
			}
		}
	}
}
func TestQuestRuntimeDirtySlotMask(t *testing.T) {
	o := newQuestRuntimeOwner(t)
	type row struct {
		Mask, Index, After uint32
		Return             uint64
	}
	var rows []row
	defer func() {
		spellbookCapture(t, "quest-runtime-slot-mask", rows, "24f2119712ce7774495c509600feb27966cc31150787c755ceec418b88e0ebe5")
	}()
	for _, mask := range []uint32{0, 1, 0x80000000, 0x55555555, 0xffffffff} {
		for index := uint32(0); index < 256; index++ {
			t.Run(fmt.Sprintf("mask%x/index%d", mask, index), func(t *testing.T) {
				*o.quest["1556300"] = mask
				got := questRuntimeCall("sub_4D79A0", nil, index)
				want := ^(uint32(1) << (index & 31))
				if got != uint64(want) || *o.quest["1556300"] != mask&want {
					t.Fatalf("mask=%x return%x want%x", *o.quest["1556300"], got, want)
				}
				rows = append(rows, row{mask, index, *o.quest["1556300"], got})
			})
		}
	}
}
