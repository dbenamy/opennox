//go:build porttest

package blobdata

// PortTestClientEffectsCurveTable returns the production Hermite coefficients.
// The headless fixture initializes this bounded table instead of overwriting
// every live blob through the full game startup sequence.
func PortTestClientEffectsCurveTable() []byte {
	return append([]byte(nil), data581450[9872:9936]...)
}

type PortTestClientEffectsTable struct {
	Base, Offset uintptr
	Data         []byte
}

func PortTestClientEffectsTables() []PortTestClientEffectsTable {
	var out []PortTestClientEffectsTable
	for _, r := range []struct {
		base, offset uintptr
		size         int
	}{
		{0x581450, 9872, 64},
		{0x587000, 178204, 16},
		{0x587000, 180460, 64},
		{0x587000, 192088, 2048},
		{0x587000, 194136, 2048},
		{0x587000, 155956, 256},
	} {
		data := data587000
		if r.base == 0x581450 {
			data = data581450
		}
		out = append(out, PortTestClientEffectsTable{r.base, r.offset, append([]byte(nil), data[int(r.offset):int(r.offset)+r.size]...)})
	}
	return out
}
