//go:build porttest

package blobdata

// PortTestClientEffectsCurveTable returns the production Hermite coefficients.
// The headless fixture initializes this bounded table instead of overwriting
// every live blob through the full game startup sequence.
func PortTestClientEffectsCurveTable() []byte {
	return append([]byte(nil), data581450[9872:9936]...)
}
