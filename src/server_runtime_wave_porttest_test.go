//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"math"
	"testing"

	"github.com/opennox/opennox/v1/common/memmap/nox/blobdata"
	"github.com/opennox/opennox/v1/legacy"
)

func TestServerRuntimeMeterWave(t *testing.T) {
	coefficients := serverConfigOwnBytes(t, 0x581450, 9760, 16)
	found := false
	for _, r := range blobdata.PortTestMeterTables() {
		if r.Base == 0x581450 && r.Offset == 9760 {
			copy(coefficients, r.Data)
			found = true
		}
	}
	if !found {
		t.Fatal("missing shipped wave coefficients")
	}
	shipped := bytes.Clone(coefficients)
	table := serverConfigOwnBytes(t, 0x5D4594, 1309840, 1280+16)
	_, scale := legacy.PortTestResourceLightConstants()
	oldScale := *scale
	t.Cleanup(func() { *scale = oldScale })
	type row struct {
		Scale        uint64
		ZeroAngle    bool
		Coefficients [2]uint64
		Values       []uint32
	}
	var rows []row
	for _, zeroAngle := range []bool{false, true} {
		copy(coefficients, shipped)
		if zeroAngle {
			binary.LittleEndian.PutUint64(coefficients, 0)
		}
		for _, value := range []float64{0, 1, -1, 32768, 65536, -65536, 1000000, 2147483648, 4294967296} {
			for i := range table {
				table[i] = 0xa5
			}
			*scale = math.Float64bits(value)
			legacy.Sub_4AEE30()
			r := row{Scale: *scale, ZeroAngle: zeroAngle, Coefficients: [2]uint64{binary.LittleEndian.Uint64(coefficients), binary.LittleEndian.Uint64(coefficients[8:])}}
			nonzero := false
			for i := 0; i < 320; i++ {
				v := binary.LittleEndian.Uint32(table[4*i:])
				r.Values = append(r.Values, v)
				nonzero = nonzero || v != 0
				if (zeroAngle || value == 0) && v != 0 {
					t.Fatal("zero wave produced a value", i, v)
				}
			}
			if !zeroAngle && math.Abs(value) >= 32768 && !nonzero {
				t.Fatal("nontrivial wave is empty")
			}
			if *scale != math.Float64bits(value) {
				t.Fatal("wave modified its scale")
			}
			for _, v := range table[1280:] {
				if v != 0xa5 {
					t.Fatal("wave exceeded its 320 entries")
				}
			}
			rows = append(rows, r)
		}
	}
	interactionCapture(t, "server-runtime-meter-wave", rows)
}
